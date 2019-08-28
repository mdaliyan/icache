package icache

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type pot struct {
	shards   shards
	window   []expireTime
	windowRW sync.RWMutex
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
	p.shards.Purge()
	p.windowRW.Unlock()
}

func (p *pot) Len() int {
	return p.shards.EntriesLen()
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
		return nil, errors.New("not found")
	}
	ti := time.Unix(ent.expiresAt, 0)
	return &ti, nil
}

func (p *pot) Get(key string, i interface{}) (err error) {
	k, shard := keyGen(key)
	ent, ok := p.shards.GetShard(shard).GetEntry(k)
	if !ok {
		return errors.New("not found")
	}

	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("second parameter needs to be a pointer")
	}
	if ent.kind != v.String()[2:] {
		vKind := v.String()[2:]
		return fmt.Errorf(`requested entry type does not match: "%s" != "%s"`, ent.kind[:len(ent.kind)-7], vKind[:len(vKind)-7])
	}
	v.Elem().Set(ent.value)

	return nil
}

func (p *pot) Set(key string, i interface{}) {
	var entry = &entry{}
	k, shard := keyGen(key)

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
		entry.expiresAt = time.Now().Add(p.ttl).UnixNano()
		p.window = append(p.window, expireTime{
			key:       k,
			shard:     shard,
			expiresAt: entry.expiresAt,
		})
		p.windowRW.Unlock()
	}

	p.shards.GetShard(shard).SetEntry(k, entry)

	return
}

func (p *pot) dropExpiredEntries() {
	var expired []expireTime
	p.windowRW.Lock()
	now := time.Now().UnixNano()
	var expiredWindows int
	for _, timeWindow := range p.window {
		if now > timeWindow.expiresAt {
			expiredWindows++
			ent, ok := p.shards.GetShard(timeWindow.shard).GetEntry(timeWindow.key)
			if ok && timeWindow.expiresAt == ent.expiresAt {
				expired = append(expired, timeWindow)
			}
		} else {
			break
		}
	}
	p.window = p.window[expiredWindows:]
	p.windowRW.Unlock()

	// fmt.Println(p.entries)
	// fmt.Println("time window:", p.window, "--->", expired)
	for _, entry := range expired {
		p.shards.GetShard(entry.shard).DropEntries(entry.key)
	}

}
