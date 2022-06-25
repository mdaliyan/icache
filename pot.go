package icache

import (
	"sync"
	"time"
)

type pot[T any] struct {
	shards   shards[T]
	window   entrySlice[T]
	windowRW sync.RWMutex
	tags     tags[T]
	ttl      time.Duration
	tick     *time.Ticker
	closed   chan bool
}

func (p *pot[T]) setTTL(TTL time.Duration) {
	p.ttl = TTL
}

func (p *pot[T]) init() {
	p.Purge()
	p.tags.pot = p
	p.closed = make(chan bool)
	if p.ttl < 1 {
		return
	}
	p.tick = time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case t := <-p.tick.C: // triggered every second
				p.dropExpiredEntries(t)
			case <-p.closed: // triggered when the pot is closed
				return
			}
		}
	}()
}

func (p *pot[T]) Close() error {
	if p.ttl < 1 {
		return ErrNotClosable
	}
	p.tick.Stop()
	close(p.closed)
	p.Purge()
	return nil
}

func (p *pot[T]) Purge() {
	p.windowRW.Lock()
	p.window = nil
	p.tags.purge()
	p.shards.Purge()
	p.windowRW.Unlock()
}

func (p *pot[T]) Len() int {
	return p.shards.EntriesLen()
}

func (p *pot[T]) Exists(key string) (ok bool) {
	k, shard := keyGen(key)
	return p.shards[shard].EntryExists(k)
}

func (p *pot[T]) ExpireTime(key string) (t *time.Time, err error) {
	k, shardID := keyGen(key)
	ent, ok := p.shards[shardID].GetEntry(k)
	if !ok {
		return nil, ErrNotFound
	}
	ti := time.Unix(ent.expiresAt, 0)
	return &ti, nil
}

func (p *pot[T]) getEntry(key string) (*entry[T], bool) {
	k, shard := keyGen(key)
	e, ok := p.shards[shard].GetEntry(k)
	if e == nil {
		p.shards[shard].DropEntry(k)
		ok = false
	}
	return e, ok
}

func (p *pot[T]) Get(key string) (v T, err error) {
	e, ok := p.getEntry(key)
	if !ok {
		return v, ErrNotFound
	}
	e.rw.RLock()
	defer e.rw.RUnlock()
	if e.deleted {
		return v, ErrNotFound
	}

	return e.data, nil
}

func (p *pot[T]) Set(key string, v T, tags ...string) {
	k, shard := keyGen(key)
	e, found := p.shards[shard].GetEntry(k)
	if found {
		p.dropEntry(e)
	}
	e = &entry[T]{
		key:       k,
		shard:     shard,
		expiresAt: time.Now().Add(p.ttl).UnixNano(),
		tags:      tagKeyGen(tags...),
		data:      v,
		deleted:   false,
	}

	if p.ttl > 0 {
		p.windowRW.Lock()
		p.window = append(p.window, e)
		p.windowRW.Unlock()
	}

	p.tags.add(e)
	p.shards[shard].SetEntry(k, e)

	return
}

func (p *pot[T]) DropTags(tags ...string) {
	p.tags.dropTags(tagKeyGen(tags...)...)
}

func (p *pot[T]) Drop(keys ...string) {
	for _, key := range keys {
		if e, ok := p.getEntry(key); ok {
			p.dropEntry(e)
		}
	}
}

func (p *pot[T]) dropExpiredEntries(at time.Time) {
	var expiredEntries entrySlice[T]
	now := at.UnixNano()

	p.windowRW.Lock()
	var expiredWindows int
	for _, entry := range p.window {
		if entry == nil {
			expiredWindows++
			continue
		}
		entry.rw.Lock()
		if now > entry.expiresAt {
			expiredWindows++
			entry.deleted = true
			expiredEntries = append(expiredEntries, entry)
		} else {
			entry.rw.Unlock()
			break
		}
		entry.rw.Unlock()
	}
	p.window = p.window[expiredWindows:]
	p.windowRW.Unlock()

	p.dropEntries(expiredEntries...)
}

func (p *pot[T]) dropEntries(entries ...*entry[T]) {
	for _, e := range entries {
		e.rw.Lock()
		if !e.deleted {
			e.rw.Unlock()
			continue
		}
		p.dropEntry(e)
	}
}

func (p *pot[T]) dropEntry(e *entry[T]) {
	e.deleted = true
	p.tags.drop(e)
	p.shards[e.shard].DropEntry(e.key)
	e = nil
}
