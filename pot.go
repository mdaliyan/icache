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
	defer p.windowRW.Unlock()

	p.window = nil
	p.tags.purge()
	p.shards.Purge()
}

func (p *pot[T]) Len() int {
	return p.shards.EntriesLen()
}

func (p *pot[T]) Exists(key string) (ok bool) {
	k, shard := keyGen(key)
	exists := p.shards[shard].EntryExists(k)

	return exists
}

func (p *pot[T]) ExpireTime(key string) (t *time.Time, err error) {
	e, ok := p.getEntry(key)
	if !ok {
		return nil, ErrNotFound
	}

	ti := time.UnixMilli(e.expiresAt)
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

func (p *pot[T]) GetByTag(tag string) ([]T, error) {
	entries := p.tags.getEntriesWithTags(tagKeyGen(tag)...)
	if len(entries) == 0 {
		return nil, ErrNotFound
	}
	result := make([]T, len(entries))
	for i, e := range entries {
		result[i] = e.data
	}
	return result, nil
}

func (p *pot[T]) Get(key string) (v T, err error) {
	e, ok := p.getEntry(key)
	if !ok {
		return v, ErrNotFound
	}

	if e.deleted {
		p.dropEntry(e)
		return v, ErrNotFound
	}

	return e.data, nil
}

func (p *pot[T]) Set(key string, v T, tags ...string) {
	expireTime := time.Now().Add(p.ttl).UnixNano()
	k, shard := keyGen(key)
	e, found := p.shards[shard].GetEntry(k)
	if found {
		p.dropEntry(e)
	}
	e = &entry[T]{
		key:       k,
		shard:     shard,
		expiresAt: expireTime,
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
}

func (p *pot[T]) DropTags(tags ...string) {
	entriesToDrop := p.tags.getEntriesWithTags(tagKeyGen(tags...)...)
	for _, e := range entriesToDrop {
		p.dropEntry(e)
	}
}

func (p *pot[T]) Drop(keys ...string) {
	for _, key := range keys {
		if e, ok := p.getEntry(key); ok {
			p.dropEntry(e)
		}
	}
}

func (p *pot[T]) dropExpiredEntries(at time.Time) {
	now := at.UnixNano()
	p.windowRW.Lock()
	defer p.windowRW.Unlock()

	var expiredWindows int
	for _, e := range p.window {
		if e == nil {
			expiredWindows++
			continue
		}
		if e.expiresAt >= now { // not expired yet
			break
		}
		e.deleted = true
		p.dropEntry(e)
		expiredWindows++
	}
	if expiredWindows > 0 {
		remaining := len(p.window) - expiredWindows
		newWindow := make(entrySlice[T], remaining)
		copy(newWindow, p.window[expiredWindows:])
		p.window = newWindow
	}
}

func (p *pot[T]) dropEntry(e *entry[T]) {
	e.deleted = true
	for _, tag := range e.tags {
		p.tags.dropTagIfNoOtherEntriesExist(tag)
	}
	p.shards[e.shard].DropEntry(e.key)
}
