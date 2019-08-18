package iCache

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type valuePot struct {
	entriesLock    sync.RWMutex
	entries        map[uint64]*entry
	timeWindow     []expireTime
	timeWindowLock sync.RWMutex
	ttl            time.Duration
}

func (c *valuePot) init(ttl time.Duration) {
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

func (c *valuePot) Purge() {
	c.entriesLock.Lock()
	c.timeWindowLock.Lock()

	c.entries = map[uint64]*entry{}
	c.timeWindow = []expireTime{}

	c.timeWindowLock.Unlock()
	c.entriesLock.Unlock()
}

func (c *valuePot) Len() (l float64) {
	c.entriesLock.Lock()
	l = float64(len(c.entries))
	c.entriesLock.Unlock()
	return l
}

func (c *valuePot) Drop(key string) {
	c.entriesLock.Lock()
	c.dropByUint64(keyGen(key))
	c.entriesLock.Unlock()
}

func (c *valuePot) dropByUint64(k uint64) {
	c.entries[k] = nil
	delete(c.entries, k)
}

func (c *valuePot) Exists(key string) bool {
	c.entriesLock.Lock()
	_, ok := c.entries[keyGen(key)]
	c.entriesLock.Unlock()
	return ok
}

func (c *valuePot) Get(key string, i interface{}) (err error) {
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

func (c *valuePot) Set(k string, i interface{}) {
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
		c.timeWindowLock.Lock()
		c.timeWindow = append(c.timeWindow, expireTime{
			Key:       key,
			ExpiresAt: time.Now().Add(c.ttl).Unix(),
		})
		c.timeWindowLock.Unlock()
	}
}

func (c *valuePot) dropExpiredEntries() {
	var expired []uint64
	c.timeWindowLock.Lock()
	for _, entrie := range c.timeWindow {
		if now() > entrie.ExpiresAt {
			expired = append(expired, entrie.Key)
		} else {
			break
		}
	}
	c.timeWindow = c.timeWindow[len(expired):]
	c.timeWindowLock.Unlock()

	// fmt.Println("time window:", c.timeWindow, "--->", expired)
	c.entriesLock.Lock()
	for _, k := range expired {
		c.dropByUint64(k)
	}
	c.entriesLock.Unlock()
}
