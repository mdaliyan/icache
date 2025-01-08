package icache

import "sync"

type shards[T any] [shardsCount]*shard[T]

func (s *shards[T]) Purge() {
	for i := 0; i < shardsCount; i++ {
		s[i] = nil
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

func (s *shard[T]) Len() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.entries)
}

func (s *shard[T]) EntryExists(key uint64) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()
	e, ok := s.entries[key]
	return ok && !e.deleted
}

func (s *shard[T]) GetEntry(key uint64) (ent *entry[T], ok bool) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	ent, ok = s.entries[key]
	return
}

func (s *shard[T]) SetEntry(key uint64, ent *entry[T]) {
	if _, ok := s.GetEntry(key); ok {
		s.DropEntry(key)
	}
	s.rw.Lock()
	defer s.rw.Unlock()

	s.entries[key] = ent
}

func (s *shard[T]) DropEntry(key uint64) {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.entries[key] = nil
	delete(s.entries, key)
}
