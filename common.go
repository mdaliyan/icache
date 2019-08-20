package iCache

import (
	"reflect"
	"time"
)

type Pot interface {
	Purge()
	Len() (l int)
	Drop(key ...string)
	Exists(key string) bool
	Get(key string, i interface{}) (err error)
	Set(k string, i interface{}) (err error)
}

func NewPot(TTL time.Duration) Pot {
	pot := new(pot)
	pot.init(TTL)
	return pot
}

type expireTime struct {
	Key       uint64
	Shard       uint64
	ExpiresAt int64
}

type entry struct {
	Value     reflect.Value
	Interface interface{}
	Kind      string
}
