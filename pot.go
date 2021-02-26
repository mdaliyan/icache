package icache

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

var NotFoundErr = errors.New("not found")
var NonPointerErr = errors.New("second parameter needs to be a pointer")

type pot struct {
	shards   shards
	window   entrySlice
	windowRW sync.RWMutex
	tags     tags
	ttl      time.Duration
	tick     *time.Ticker
	closed   chan bool
}

func (p *pot) reset() {
	p.ttl = 0
	p.Purge()
	p.closed = make(chan bool)
	if p.tick != nil {
		p.tick.Stop()
		p.tick = nil
	}
}

func (p *pot) init(TTL time.Duration) {
	p.reset()
	p.ttl = TTL
	if p.ttl > 1 {
		p.tick = time.NewTicker(time.Second)
		go func() {
			for {
				select {
				case <-p.closed:
					p.reset()
					p.closed = nil
					return
				case t := <-p.tick.C:
					p.dropExpiredEntries(t)
				}
			}
		}()
	}
}

func (p *pot) Purge() {
	p.windowRW.Lock()
	p.window = nil
	p.tags.purge(p)
	p.shards.Purge()
	p.windowRW.Unlock()
}

func (p *pot) Len() int {
	return p.shards.EntriesLen()
}

func (p *pot) DropTags(tags ...string) {
	p.tags.dropTags(TagKeyGen(tags)...)
}

func (p *pot) Drop(keys ...string) {
	for _, key := range keys {
		e, ok := p.getEntry(key)
		if !ok {
			continue
		}
		e.deleted = true
		p.tags.drop(e)
		p.shards[e.shard].DropEntry(e.key)
		e = nil
	}
}

func (p *pot) Exists(key string) (ok bool) {
	k, shard := keyGen(key)
	return p.shards[shard].EntryExists(k)
}

func (p *pot) ExpireTime(key string) (t *time.Time, err error) {
	k, shardID := keyGen(key)
	ent, ok := p.shards[shardID].GetEntry(k)
	if !ok {
		return nil, NotFoundErr
	}
	ti := time.Unix(ent.expiresAt, 0)
	return &ti, nil
}

func (p *pot) getEntry(key string) (*entry, bool) {
	k, shard := keyGen(key)
	e, ok := p.shards[shard].GetEntry(k)
	if e == nil {
		p.shards[shard].DropEntry(k)
		ok = false
	}
	return e, ok
}

func (p *pot) Get(key string, i interface{}) (err error) {
	e, ok := p.getEntry(key)
	if !ok {
		return NotFoundErr
	}
	e.rw.RLock()
	defer e.rw.RUnlock()
	if e.deleted {
		return NotFoundErr
	}

	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return NonPointerErr
	}
	if e.kind != v.String()[2:] {
		vKind := v.String()[2:]
		return fmt.Errorf(`requested entry type does not match: "%s" != "%s"`, e.kind[:len(e.kind)-7], vKind[:len(vKind)-7])
	}
	v.Elem().Set(e.value)

	return nil
}

func (p *pot) Set(key string, i interface{}, tags ...string) {
	k, shard := keyGen(key)
	e, ok := p.shards[shard].GetEntry(k)
	if !ok {
		e = &entry{}
	}
	e.rw.Lock()
	defer e.rw.Unlock()

	var v reflect.Value
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		v = reflect.ValueOf(i).Elem()
	} else {
		v = reflect.ValueOf(i)
	}
	e.key = k
	e.shard = shard
	e.expiresAt = time.Now().Add(p.ttl).UnixNano()
	e.tags = TagKeyGen(tags)
	e.value = v
	e.deleted = false
	e.kind = v.String()[1:]

	if p.ttl > 0 {
		p.windowRW.Lock()
		p.window = append(p.window, e)
		p.windowRW.Unlock()
	}

	p.tags.add(e)
	p.shards[shard].SetEntry(k, e)

	return
}

func (p *pot) dropExpiredEntries(at time.Time) {
	var expiredEntries entrySlice
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

func (p *pot) dropEntries(entries ...*entry) {
	for _, entry := range entries {
		entry.rw.Lock()
		if !entry.deleted {
			entry.rw.Unlock()
			continue
		}
		p.tags.drop(entry)
		p.shards[entry.shard].DropEntry(entry.key)
		entry.rw.Unlock()
		entry = nil
	}
}
