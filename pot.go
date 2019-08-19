package iCache

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type pot struct {
	shards         shards
	entriesLock    sync.RWMutex
	entries        map[uint64]*entry
	timeWindow     []expireTime
	timeWindowLock sync.RWMutex
	ttl            time.Duration
	typ            Type
	multiShard     bool
}

func (p *pot) init(config Config) {
	p.ttl = config.TTL
	p.typ = config.Type
	p.multiShard = config.MultiShard
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

	p.entriesLock.Lock()
	p.timeWindowLock.Lock()

	p.timeWindow = []expireTime{}
	if p.multiShard {
		for i := 0; i < 256; i++ {
			p.shards[i] = &shard{
				entries: map[uint64]*entry{},
			}
		}
	} else {
		p.entries = map[uint64]*entry{}
	}

	p.timeWindowLock.Unlock()
	p.entriesLock.Unlock()
}

func (p *pot) getEntry(key uint64) (ent *entry, ok bool) {
	if p.multiShard {
		return p.shards.GetShard(key).GetEntry(key)
	} else {
		p.entriesLock.Lock()
		ent, ok = p.entries[key]
		p.entriesLock.Unlock()
	}
	return
}

func (p *pot) setEntry(key uint64, ent *entry) {
	if p.multiShard {
		p.shards.GetShard(key).SetEntry(key, ent)
	} else {
		p.entriesLock.Lock()
		p.entries[key] = ent
		p.entriesLock.Unlock()
	}
}

func (p *pot) dropEntries(keys ...uint64) {
	if p.multiShard {
		for _, key := range keys {
			p.shards.GetShard(key).DropEntries(key)
		}
	} else {
		p.entriesLock.Lock()
		for _, k := range keys {
			p.entries[k] = nil
			delete(p.entries, k)
		}
		p.entriesLock.Unlock()
	}
}

func (p *pot) Len() (l float64) {
	if p.multiShard {
	} else {
		p.entriesLock.Lock()
		l = float64(len(p.entries))
		p.entriesLock.Unlock()
	}
	return l
}

func (p *pot) Drop(keys ...string) {
	var ks []uint64
	for _, k := range keys {
		ks = append(ks, keyGen(k))
	}
	p.dropEntries(ks...)
}

func (p *pot) Exists(key string) (ok bool) {
	_, ok = p.getEntry(keyGen(key))
	return
}

func (p *pot) Get(key string, i interface{}) (err error) {

	k := keyGen(key)

	ent, ok := p.getEntry(k)
	if ! ok {
		return errors.New("not found")
	}

	switch p.typ {
	case Value:
		v := reflect.ValueOf(i)
		if v.Kind() != reflect.Ptr || v.IsNil() {
			return errors.New("need to be a pointer")
		}
		if ent.Kind != v.String()[2:] {
			return errors.New("requested entry type does not match: \"" + ent.Kind + "!=" + v.String()[2:] + "\"")
		}
		v.Elem().Set(ent.Value)
	}

	return nil
}

func (p *pot) Set(k string, i interface{}) (err error) {
	var entry = &entry{}
	key := keyGen(k)

	typ := reflect.TypeOf(i)

	switch p.typ {
	case Value:
		var v reflect.Value
		if typ.Kind() == reflect.Ptr {
			v = reflect.ValueOf(i).Elem()
		} else {
			v = reflect.ValueOf(i)
		}
		entry.Value = v
		entry.Kind = v.String()[1:]
	case Pointer:
		if typ.Kind() != reflect.Ptr {
			return errors.New("need to be a pointer")
		}
	}

	p.setEntry(key, entry)

	if p.ttl > 0 {
		p.timeWindowLock.Lock()
		p.timeWindow = append(p.timeWindow, expireTime{
			Key:       key,
			ExpiresAt: time.Now().Add(p.ttl).UnixNano(),
		})
		p.timeWindowLock.Unlock()
	}
	return
}

func (p *pot) dropExpiredEntries() {
	var expired []uint64
	p.timeWindowLock.Lock()
	now := time.Now().UnixNano()
	for _, entry := range p.timeWindow {
		if now > entry.ExpiresAt {
			expired = append(expired, entry.Key)
		} else {
			break
		}
	}
	p.timeWindow = p.timeWindow[len(expired):]
	p.timeWindowLock.Unlock()

	// fmt.Println(p.entries)
	// fmt.Println("time window:", p.timeWindow, "--->", expired)
	p.dropEntries(expired...)

}
