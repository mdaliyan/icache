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
	rw           sync.RWMutex
	Entries      entries
	ExpiredDates expiredDates
	path         string
	flushTime    int64
	flushLen     float64
	inited       bool
}

var now int64

func NewDiskPot(path string) (Cache *Pot, err error) {
	Cache = new(Pot)
	Cache.Entries = entries{}
	Cache.ExpiredDates = expiredDates{}
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
	c.rw.Lock()
	c.Entries = entries{}
	c.ExpiredDates = expiredDates{}
	c.inited = true
	c.rw.Unlock()
}

func (c *Pot) panicIfNotInitialized() {
	if ! c.inited {
		panic("iCache should be initialized before use")
	}
}

func (c *Pot) Len() (l float64) {

	c.rw.Lock()
	l = float64(len(c.Entries))
	c.rw.Unlock()
	return l
}

func (c *Pot) Drop(key string) {
	k := keyGen(key)
	c.rw.Lock()
	c.Entries[k] = nil
	delete(c.Entries, k)
	c.rw.Unlock()
}

func (c *Pot) Exists(key string) bool {
	c.rw.Lock()
	_, ok := c.Entries[keyGen(key)]
	c.rw.Unlock()
	return ok
}

func (c *Pot) Flush() (err error, took time.Duration) {
	if c.path == "" {
		return errors.New("disk cache is not enabled"), 0
	}
	st := time.Now()
	c.rw.Lock()
	c.flushTime = now
	c.flushLen = float64(len(c.Entries))
	temp := *c
	c.rw.Unlock()
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
	c.rw.Lock()
	ent, ok := c.Entries[k]
	c.rw.Unlock()

	if ! ok {
		return errors.New("not found")
	}
	if ent.Ttl > 0 {
		if now > ent.ExpiresAt {
			c.rw.Lock()
			c.Entries[k] = nil
			delete(c.Entries, k)
			c.rw.Unlock()
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
	c.rw.Lock()
	c.Entries[keyGen(k)] = &entry{
		Data:      v,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	}
	if ttl > 0 {
		c.ExpiredDates[keyGen(k)] = ExpiresAt
	}
	c.rw.Unlock()
}

func (c *Pot) deleteAllExpired() {
	for k, expiresAt := range c.ExpiredDates {
		if now > expiresAt {
			c.rw.Lock()
			c.Entries[k] = nil
			delete(c.ExpiredDates, k)
			delete(c.Entries, k)
			c.rw.Unlock()
		}
	}
}

func init() {
	go func() {
		for {
			now = time.Now().Unix()
			time.Sleep(time.Second)
		}
	}()
}
