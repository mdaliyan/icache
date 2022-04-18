package icache

import (
	"time"
)

// Pot holds your cached data
type Pot[T any] interface {
	Purge()
	Len() (l int)
	Drop(key ...string)
	DropTags(tags ...string)
	Exists(key string) bool
	Set(k string, v T, tags ...string)
	Get(key string) (v T, err error)
	ExpireTime(key string) (t *time.Time, err error)
}

// NewPot creates new Pot with a given ttl duration
func NewPot[T any](TTL time.Duration) Pot[T] {
	pot := new(pot[T])
	pot.init(TTL)
	return pot
}
