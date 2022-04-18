package icache

import "sync"

type entries[T any] map[uint64]*entry[T]

type entrySlice[T any] []*entry[T]

type entry[T any] struct {
	key       uint64
	shard     uint64
	data      T
	expiresAt int64
	tags      []uint64
	deleted   bool
	rw        sync.RWMutex
}
