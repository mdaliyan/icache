package icache

import (
	"sync"
)

type tags struct {
	pairs map[uint64]entries
	rw    sync.RWMutex
	pot   *pot
}

func (t *tags) purge(p *pot) {
	t.rw.Lock()
	t.pairs = make(map[uint64]entries)
	t.pot = p
	t.rw.Unlock()
}

func (t *tags) add(e *entry) {
	if e == nil || e.tags == nil {
		return
	}
	t.rw.Lock()
	for _, tag := range e.tags {
		var _, ok = t.pairs[tag]
		if ok {
			t.pairs[tag][e.key] = e
			continue
		}
		tags := make(entries)
		tags[e.key] = e
		t.pairs[tag] = tags
	}
	t.rw.Unlock()
}

func (t *tags) drop(e *entry) {
	if e == nil || e.tags == nil {
		return
	}
	t.rw.Lock()
	for _, tag := range e.tags {
		var _, ok = t.pairs[tag]
		if ok {
			continue
		}
		delete(t.pairs[tag], e.key)
		if len(t.pairs[tag]) == 0 {
			delete(t.pairs, tag)
		}
	}
	t.rw.Unlock()
}

func (t *tags) dropTags(tags ...uint64) {
	for _, tag := range tags {
		var entries entrySlice
		t.rw.RLock()
		if _, ok := t.pairs[tag]; !ok {
			t.rw.RUnlock()
			continue
		}
		for _, e := range t.pairs[tag] {
			e.deleted = true
			entries = append(entries, e)
		}
		t.rw.RUnlock()

		t.rw.Lock()
		delete(t.pairs, tag)
		t.rw.Unlock()
		t.pot.dropEntries(entries...)
	}
}
