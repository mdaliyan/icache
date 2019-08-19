package iCache

import (
	"reflect"
	"time"
)

type Pot interface {
	Purge()
	Len() (l float64)
	Drop(key ...string)
	Exists(key string) bool
	Get(key string, i interface{}) (err error)
	Set(k string, i interface{}) (err error)
}

type Config struct {
	Type       Type
	MultiShard bool
	TTL        time.Duration
}

type Type int

const (
	Value Type = iota
	Pointer
)

func NewPot(config Config) Pot {
	pot := new(pot)
	pot.init(config)
	return pot
}

type expireTime struct {
	Key       uint64
	ExpiresAt int64
}

type entry struct {
	Value     reflect.Value
	Interface interface{}
	Kind      string
}
