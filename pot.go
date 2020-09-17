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
}

func (p *pot) init(TTL time.Duration) {
	p.ttl = TTL
	p.Purge()
	if p.ttl > 1 {
		go func() {
			for {
				p.dropExpiredEntries()
				time.Sleep(time.Second)
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
	p.tags.dropTags(tags...)
	// p.tags.dropTags(TagKeyGen(tags)...)
}

func (p *pot) Drop(keys ...string) {
	for _, key := range keys {
		k, shard := keyGen(key)
		p.shards.GetShard(shard).DropEntries(k)
	}
}

func (p *pot) Exists(key string) (ok bool) {
	k, shard := keyGen(key)
	return p.shards.GetShard(shard).EntryExists(k)
}

func (p *pot) ExpireTime(key string) (t *time.Time, err error) {
	k, shardID := keyGen(key)
	ent, ok := p.shards.GetShard(shardID).GetEntry(k)
	if !ok {
		return nil, NotFoundErr
	}
	ti := time.Unix(ent.expiresAt, 0)
	return &ti, nil
}

func (p *pot) Get(key string, i interface{}) (err error) {
	k, shard := keyGen(key)
	ent, ok := p.shards.GetShard(shard).GetEntry(k)
	if !ok {
		return NotFoundErr
	}

	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return NonPointerErr
	}
	if ent.kind != v.String()[2:] {
		vKind := v.String()[2:]
		return fmt.Errorf(`requested entry type does not match: "%s" != "%s"`, ent.kind[:len(ent.kind)-7], vKind[:len(vKind)-7])
	}
	v.Elem().Set(ent.value)

	return nil
}

func (p *pot) Set(key string, i interface{}, tags ...string) {
	k, shard := keyGen(key)
	var entry = &entry{
		key:       k,
		shard:     shard,
		expiresAt: time.Now().Add(p.ttl).UnixNano(),
		tags:      tags, //TagKeyGen(tags),
	}

	var v reflect.Value
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		v = reflect.ValueOf(i).Elem()
	} else {
		v = reflect.ValueOf(i)
	}
	entry.value = v
	entry.kind = v.String()[1:]

	if p.ttl > 0 {
		p.windowRW.Lock()
		p.window = append(p.window, entry)
		p.windowRW.Unlock()
	}

	p.tags.add(entry)
	p.shards.GetShard(shard).SetEntry(k, entry)

	return
}

func (p *pot) dropExpiredEntries() {
	var expiredEntries entrySlice
	p.windowRW.Lock()
	now := time.Now().UnixNano()
	var expiredWindows int
	for _, entry := range p.window {
		if entry == nil {
			expiredWindows++
			continue
		}
		if now > entry.expiresAt {
			expiredWindows++
			expiredEntries = append(expiredEntries, entry)
		} else {
			break
		}
	}
	p.window = p.window[expiredWindows:]
	p.windowRW.Unlock()

	// fmt.Println(p.entries)
	// fmt.Println("time window:", p.window, "--->", expired)
	p.dropEntries(expiredEntries)
}

func (p *pot) dropEntries(entries entrySlice) {
	for _, entry := range entries {
		p.shards.GetShard(entry.shard).DropEntries(entry.key)
		p.tags.drop(entry)
		entry = nil
	}
}
