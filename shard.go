package icache

import "sync"

type shards struct {
	Map sync.Map
	len uint64
}

func (s *shards) Purge() {
	s.Map.Range(func(key interface{}, _ interface{}) bool {
		s.Map.Delete(key)
		return true
	})
}

func (s *shards) Len() (l int) {
	s.Map.Range(func(_, _ interface{}) bool {
		l++
		return true
	})
	return
}

func (s *shards) EntryExists(key string) (ok bool) {
	_, ok = s.Map.Load(key)
	return
}

func (s *shards) GetEntry(key string) (*entry, bool) {
	e, ok := s.Map.Load(key)
	if ok {
		return e.(*entry), ok
	}
	return nil, ok
}

func (s *shards) SetEntry(key string, ent *entry) {
	s.Map.Store(key, ent)
}

func (s *shards) DropEntries(keys ...string) {
	for _, k := range keys {
		s.Map.Delete(k)
	}
}
