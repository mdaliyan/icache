package iCache

import (
	"reflect"
	"sync"
	"time"
)

type Pot interface {
	Purge()
	Len() (l float64)
	Drop(key string)
	Exists(key string) bool
	Get(key string, i interface{}) (err error)
	Set(k string, i interface{})
}

type Config struct {
	Type       Type
	MultiShard bool
	TTL        time.Duration
}

type Type int

const (
	Value Type = iota
	Interface
)

func NewPot(config Config) (Cache Pot) {
	pot := new(valuePot)
	pot.init(config.TTL)
	return pot
}

var nowUint int64
var nowLock sync.RWMutex

func now() (n int64) {
	nowLock.Lock()
	n = nowUint
	nowLock.Unlock()
	return
}

func init() {
	go func() {
		for {
			nowLock.Lock()
			nowUint = time.Now().Unix()
			nowLock.Unlock()
			time.Sleep(time.Second)
		}
	}()
}

type expireTime struct {
	Key       uint64
	ExpiresAt int64
}

type entry struct {
	Data reflect.Value
}
