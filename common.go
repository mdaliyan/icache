package icache

import (
	"reflect"
	"time"
)

// Pot holds your cached data
type Pot interface {
	Purge()
	Len() (l int)
	Drop(key ...string)
	Exists(key string) bool
	Set(k string, i interface{})
	Get(key string, i interface{}) (err error)
	ExpireTime(key string) (t *time.Time, err error)
}

// NewPot creates new Pot with a given ttl duration
func NewPot(TTL time.Duration) Pot {
	pot := new(pot)
	pot.init(TTL)
	return pot
}

type expireTime struct {
	key       string
	expiresAt int64
}

type entries map[uint64]*entry

type entry struct {
	value     reflect.Value
	expiresAt int64
	kind      string
}
