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

func (s *shards) GetShard(key uint64) (shard *shard) {
	return s[key]
}

func (s shards) EntriesLen() (l int) {
	for _, shard := range s {
		l += shard.Len()
	}
	return
}

type shard struct {
	entries     entries
	entriesLock sync.RWMutex
}

func (s *shard) Len() (l int) {
	s.entriesLock.Lock()
	l = len(s.entries)
	s.entriesLock.Unlock()
	return
}

func (s *shard) EntryExists(key uint64) (ok bool) {
	s.entriesLock.Lock()
	_, ok = s.entries[key]
	s.entriesLock.Unlock()
	return
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
