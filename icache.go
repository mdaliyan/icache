package icache

import (
	"errors"
	"time"
)

// ErrNotFound means the pot couldn't find the entry with the given key
var ErrNotFound = errors.New("entry not found")

// ErrNotClosable means the pot doesn't have invalidation strategy.
var ErrNotClosable = errors.New("pot cannot be closed")

// Pot holds your cached data
type Pot[T any] interface {

	// Purge invalidates all entries
	Purge()

	// Len returns count of the entries
	Len() (l int)

	// Drop invalidates the entry with the given key
	Drop(key ...string)

	// DropTags invalidates all entries with the given tags
	DropTags(tags ...string)

	// Exists checks if the entry with the given key exists
	Exists(key string) bool

	// Set stores the variable in the given key entry
	// tags can be set or ignored
	Set(k string, v T, tags ...string)

	// Get restores the T previously stored as the given key
	// returns ErrNotFound if the key doesn't exist.
	Get(key string) (v T, err error)

	// ExpireTime returns expire time of the given entry
	ExpireTime(key string) (t *time.Time, err error)

	// Close stops the invalidation functionality goroutines.
	// - After closing, the pot resets and next entries you add won't expire. so don't use the pot after closing it.
	// - Pots without invalidation strategy cannot be closed.
	// - Panics if the pot is already closed.
	Close() error
}

// NewPot creates new cache Pot with the given options.
func NewPot[T any](options ...Option) Pot[T] {
	p := new(pot[T])
	for _, option := range options {
		option(p)
	}
	p.init()
	return p
}
