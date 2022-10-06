package icache

import "sync"

type shards [shardsCount]*shard

func (s *shards) Purge() {
	for i := 0; i < shardsCount; i++ {
		s[i] = &shard{
			entries: entries{},
		}
	}
}

func (s shards) EntriesLen() (l int) {
	for _, shard := range s {
		l += shard.Len()
	}
	return
}

type shard struct {
	entries entries
	rw      sync.RWMutex
}

func (s *shard) Len() (l int) {
	s.rw.RLock()
	l = len(s.entries)
	s.rw.RUnlock()
	return
}

func (s *shard) EntryExists(key uint64) (ok bool) {
	s.rw.RLock()
	_, ok = s.entries[key]
	s.rw.RUnlock()
	return
}

func (s *shard) GetEntry(key uint64) (ent *entry, ok bool) {
	s.rw.RLock()
	ent, ok = s.entries[key]
	s.rw.RUnlock()
	return
}

func (s *shard) SetEntry(key uint64, ent *entry) {
	s.rw.Lock()
	e, ok := s.entries[key]
	if ok {
		e.key = ent.key
		e.shard = ent.shard
		e.value = ent.value
		e.expiresAt = ent.expiresAt
		e.kind = ent.kind
		e.tags = ent.tags
		e.deleted = ent.deleted
	} else {
		s.entries[key] = ent
	}
	s.rw.Unlock()
}

func (s *shard) DropEntry(keys uint64) {
	s.rw.Lock()
	s.entries[keys] = nil
	delete(s.entries, keys)
	s.rw.Unlock()
}
