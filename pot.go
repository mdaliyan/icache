package iCache

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

func NewPot(ttl time.Duration) (Cache *Pot) {
	Cache = new(Pot)
	Cache.init(ttl)
	return
}

type Pot struct {
	entriesLock      sync.RWMutex
	entries          map[uint64]*entry
	expireTimes      []expireTime
	expiredTimesLock sync.RWMutex
	ttl              time.Duration
}

func (c *Pot) init(ttl time.Duration) {
	c.Purge()
	c.ttl = ttl
	if ttl > 1 {
		go func() {
			for {
				c.dropExpiredEntries()
				time.Sleep(time.Second)
			}
		}()
	}
}

func (c *Pot) Purge() {
	c.entriesLock.Lock()
	c.expiredTimesLock.Lock()

	c.entries = map[uint64]*entry{}
	c.expireTimes = []expireTime{}

	c.expiredTimesLock.Unlock()
	c.entriesLock.Unlock()
}

func (c *Pot) Len() (l float64) {
	c.entriesLock.Lock()
	l = float64(len(c.entries))
	c.entriesLock.Unlock()
	return l
}

func (c *Pot) Drop(key string) {
	c.entriesLock.Lock()
	c.dropByUint64(keyGen(key))
	c.entriesLock.Unlock()
}

func (c *Pot) dropByUint64(k uint64) {
	c.entries[k] = nil
	delete(c.entries, k)
}

func (c *Pot) Exists(key string) bool {
	c.entriesLock.Lock()
	_, ok := c.entries[keyGen(key)]
	c.entriesLock.Unlock()
	return ok
}

func (c *Pot) Get(key string, i interface{}) (err error) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("need to be a pointer")
	}
	k := keyGen(key)
	c.entriesLock.Lock()
	ent, ok := c.entries[k]
	c.entriesLock.Unlock()
	if ! ok {
		return errors.New("not found")
	}
	v.Elem().Set(ent.Data)
	return nil
}

func (c *Pot) Set(k string, i interface{}) {
	var v reflect.Value
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		v = reflect.ValueOf(i).Elem()
	} else {
		v = reflect.ValueOf(i)
	}

	key := keyGen(k)

	c.entriesLock.Lock()
	c.entries[key] = &entry{
		Data: v,
	}
	c.entriesLock.Unlock()

	if c.ttl > 0 {
		c.expiredTimesLock.Lock()
		c.expireTimes = append(c.expireTimes, expireTime{
			Key:       key,
			ExpiresAt: time.Now().Add(c.ttl).Unix(),
		})
		c.expiredTimesLock.Unlock()
	}
}

func (c *Pot) dropExpiredEntries() {
	var expired []uint64
	c.expiredTimesLock.Lock()
	for _, entrie := range c.expireTimes {
		if now() > entrie.ExpiresAt {
			expired = append(expired, entrie.Key)
		} else {
			break
		}
	}
	c.expireTimes = c.expireTimes[len(expired):]
	c.expiredTimesLock.Unlock()

	c.entriesLock.Lock()
	for _, k := range expired {
		c.dropByUint64(k)
	}
	c.entriesLock.Unlock()
}
