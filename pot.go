package iCache

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type pot struct {
	shards         shards
	timeWindow     []expireTime
	timeWindowLock sync.RWMutex
	ttl            time.Duration
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
	p.timeWindowLock.Lock()
	p.timeWindow = []expireTime{}
	for i := 0; i < 256; i++ {
		p.shards[i] = &shard{
			entries: map[uint64]*entry{},
		}
	}
	p.timeWindowLock.Unlock()
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

func (p *pot) Get(key string, i interface{}) (err error) {
	k, shard := keyGen(key)
	ent, ok := p.shards.GetShard(shard).GetEntry(k)
	if !ok {
		return errors.New("not found")
	}

	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("need to be a pointer")
	}
	if ent.kind != v.String()[2:] {
		return errors.New("requested entry type does not match: \"" + ent.kind + "!=" + v.String()[2:] + "\"")
	}
	v.Elem().Set(ent.value)

	return nil
}

func (p *pot) Set(key string, i interface{}) (err error) {
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

	p.shards.GetShard(shard).SetEntry(k, entry)

	if p.ttl > 0 {
		p.timeWindowLock.Lock()
		p.timeWindow = append(p.timeWindow, expireTime{
			key:       k,
			shard:     shard,
			expiresAt: time.Now().Add(p.ttl).UnixNano(),
		})
		p.timeWindowLock.Unlock()
	}
	return
}

func (p *pot) dropExpiredEntries() {
	var expired []expireTime
	p.timeWindowLock.Lock()
	now := time.Now().UnixNano()
	for _, entry := range p.timeWindow {
		if now > entry.expiresAt {
			expired = append(expired, entry)
		} else {
			break
		}
	}
	p.timeWindow = p.timeWindow[len(expired):]
	p.timeWindowLock.Unlock()

	// fmt.Println(p.entries)
	// fmt.Println("time window:", p.timeWindow, "--->", expired)
	for _, entry := range expired {
		p.shards.GetShard(entry.shard).DropEntries(entry.key)
	}

}
