package icache

import (
	"reflect"
	`sync`
	"time"
)

// Pot holds your cached data
type Pot interface {
	Purge()
	Len() (l int)
	Drop(key ...string)
	DropTags(tags ...string)
	Exists(key string) bool
	Set(k string, i interface{}, tags ...string)
	Get(key string, i interface{}) (err error)
	ExpireTime(key string) (t *time.Time, err error)
}

// NewPot creates new Pot with a given ttl duration
func NewPot(TTL time.Duration) Pot {
	pot := new(pot)
	pot.init(TTL)
	return pot
}

type entries map[uint64]*entry
type entrySlice []*entry

type entry struct {
	key       uint64
	shard     uint64
	value     reflect.Value
	expiresAt int64
	kind      string
	tags      []uint64
	deleted   bool
	rw sync.RWMutex
}
