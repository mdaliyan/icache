package iCache

import (
	"encoding/gob"
	"errors"
	"log"
	"os"
	"reflect"
	"sync"
	"time"
)

type entries map[uint64]*entry
type expiredDates map[uint64]int64

type entry struct {
	Data      reflect.Value
	ExpiresAt int64
	Ttl       int64
}

type Pot struct {
	entriesLock      sync.RWMutex
	Entries          entries
	ExpiredDates     expiredDates
	expiredDatesLock sync.RWMutex
	path             string
	flushTime        int64
	flushLen         float64
	inited           bool
}

var now int64

func NewDiskPot(path string) (Cache *Pot, err error) {
	Cache = new(Pot)
	Cache.Init()
	Cache.path = path
	if file, err := os.Open(path); err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(Cache)
		file.Close()
	}
	go func() {
		sleepTime := time.Second * 20
		var took time.Duration = 0
		var loopErr error
		for {
			time.Sleep(sleepTime + took*100)
			if Cache.Len() > Cache.flushLen*1.05 || now-Cache.flushTime > 300 {
				if loopErr, took = Cache.Flush(); loopErr != nil {
					log.Fatal("error on write to disk at", Cache.path, ":", loopErr.Error())
				}
			}
		}
	}()
	return
}

func NewPot() (Cache *Pot) {
	Cache = new(Pot)
	Cache.Init()
	return
}

func (c *Pot) Init() {
	c.Purge()
	go func() {
		for {
			c.deleteAllExpired()
			time.Sleep(time.Second)
		}
	}()
}

func (c *Pot) Purge() {
	c.entriesLock.Lock()
	c.Entries = entries{}
	c.entriesLock.Unlock()

	c.expiredDatesLock.Lock()
	c.ExpiredDates = expiredDates{}
	c.expiredDatesLock.Unlock()

	c.inited = true
}

func (c *Pot) panicIfNotInitialized() {
	if ! c.inited {
		panic("iCache should be initialized before use")
	}
}

func (c *Pot) Len() (l float64) {
	c.entriesLock.Lock()
	l = float64(len(c.Entries))
	c.entriesLock.Unlock()
	return l
}

func (c *Pot) Drop(key string) {
	k := keyGen(key)
	c.entriesLock.Lock()
	c.Entries[k] = nil
	delete(c.Entries, k)
	c.entriesLock.Unlock()
}

func (c *Pot) dropByUint64(k uint64) {
	c.entriesLock.Lock()
	c.expiredDatesLock.Lock()
	c.Entries[k] = nil
	delete(c.Entries, k)
	delete(c.ExpiredDates, k)
	c.expiredDatesLock.Unlock()
	c.entriesLock.Unlock()
}

func (c *Pot) Exists(key string) bool {
	c.entriesLock.Lock()
	_, ok := c.Entries[keyGen(key)]
	c.entriesLock.Unlock()
	return ok
}

func (c *Pot) Flush() (err error, took time.Duration) {
	if c.path == "" {
		return errors.New("disk cache is not enabled"), 0
	}
	st := time.Now()
	c.entriesLock.Lock()
	c.flushTime = now
	c.flushLen = float64(len(c.Entries))
	temp := *c
	c.entriesLock.Unlock()
	file, err := os.Create(temp.path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(temp)
		file.Close()
	}
	took = time.Since(st)
	return
}

func (c *Pot) Get(key string, i interface{}) (err error) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("need to be a pointer")
	}
	k := keyGen(key)
	c.entriesLock.Lock()
	ent, ok := c.Entries[k]
	c.entriesLock.Unlock()

	if ! ok {
		return errors.New("not found")
	}
	if ent.Ttl > 0 {
		if now > ent.ExpiresAt {
			c.dropByUint64(k)
			return errors.New("expired")
		}
		// ent.ExpiresAt = time.Now().Unix() + ent.Ttl
	}
	v.Elem().Set(ent.Data)

	return nil
}

func (c *Pot) Set(k string, i interface{}, ttl time.Duration) {
	ExpiresAt := time.Now().Add(ttl).Unix()
	var v reflect.Value
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		v = reflect.ValueOf(i).Elem()
	} else {
		v = reflect.ValueOf(i)
	}
	c.entriesLock.Lock()
	c.Entries[keyGen(k)] = &entry{
		Data:      v,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	}
	if ttl > 0 {
		c.expiredDatesLock.Lock()
		c.ExpiredDates[keyGen(k)] = ExpiresAt
		c.expiredDatesLock.Unlock()
	}
	c.entriesLock.Unlock()
}

func (c *Pot) deleteAllExpired() {
	var expired []uint64
	for k, expiresAt := range c.ExpiredDates {
		if now > expiresAt {
			expired = append(expired, k)
		}
	}
	c.entriesLock.Lock()
	for _, k := range expired {
		c.Entries[k] = nil
		delete(c.ExpiredDates, k)
		delete(c.Entries, k)
	}
	c.entriesLock.Unlock()
}

func init() {
	go func() {
		for {
			now = time.Now().Unix()
			time.Sleep(time.Second)
		}
	}()
}
