package iCache

import "sync"

type shards [256]*shard

func (c *shards) GetShard(key uint64) (shard *shard) {
	return c[uint8(key)]
}

type shard struct {
	entries     map[uint64]*entry
	entriesLock sync.RWMutex
}

func (s *shard) GetEntry(key uint64) (ent *entry, ok bool) {
	s.entriesLock.Lock()
	ent, ok = s.entries[key]
	s.entriesLock.Unlock()
	return
}

func (s *shard) SetEntry(key uint64, ent *entry) {
	s.entriesLock.Lock()
	s.entries[key] = ent
	s.entriesLock.Unlock()
}

func (s *shard) DropEntries(keys ...uint64) {
	s.entriesLock.Lock()
	for _, k := range keys {
		s.entries[k] = nil
		delete(s.entries, k)
	}
	s.entriesLock.Unlock()
}
