package icache

import "sync"

type shards[T any] [shardsCount]*shard[T]

func (s *shards[T]) Purge() {
	for i := 0; i < shardsCount; i++ {
		s[i] = &shard[T]{
			entries: entries[T]{},
		}
	}
}

func (s shards[T]) EntriesLen() (l int) {
	for _, shard := range s {
		l += shard.Len()
	}
	return
}

type shard[T any] struct {
	entries entries[T]
	rw      sync.RWMutex
}

func (s *shard[T]) Len() (l int) {
	s.rw.RLock()
	l = len(s.entries)
	s.rw.RUnlock()
	return
}

func (s *shard[T]) EntryExists(key uint64) (ok bool) {
	s.rw.RLock()
	_, ok = s.entries[key]
	s.rw.RUnlock()
	return
}

func (s *shard[T]) GetEntry(key uint64) (ent *entry[T], ok bool) {
	s.rw.RLock()
	ent, ok = s.entries[key]
	s.rw.RUnlock()
	return
}

func (s *shard[T]) SetEntry(key uint64, ent *entry[T]) {
	s.rw.Lock()
	e, ok := s.entries[key]
	if ok {
		*e = *ent
	} else {
		s.entries[key] = ent
	}
	s.rw.Unlock()
}

func (s *shard[T]) DropEntry(keys uint64) {
	s.rw.Lock()
	s.entries[keys] = nil
	delete(s.entries, keys)
	s.rw.Unlock()
}
