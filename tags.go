package icache

import (
	"sync"
)

type tags[T any] struct {
	rw    sync.RWMutex
	pairs map[uint64]entries[T]
}

func (t *tags[T]) purge() {
	t.rw.Lock()
	defer t.rw.Unlock()
	t.pairs = make(map[uint64]entries[T])
}

func (t *tags[T]) add(e *entry[T]) {
	t.rw.Lock()
	defer t.rw.Unlock()
	if e == nil || e.tags == nil {
		return
	}
	for _, tag := range e.tags {
		var _, ok = t.pairs[tag]
		if ok {
			t.pairs[tag][e.key] = e
			continue
		}
		tags := make(entries[T])
		tags[e.key] = e
		t.pairs[tag] = tags
	}
}

func (t *tags[T]) dropTagIfNoOtherEntriesExist(tag uint64) {
	entries := t.getEntriesWithTags(tag)

	t.rw.Lock()
	defer t.rw.Unlock()
	if len(entries) == 0 {
		t.pairs[tag] = nil
		delete(t.pairs, tag)
	}
}

func (t *tags[T]) getEntriesWithTags(tags ...uint64) entrySlice[T] {
	t.rw.RLock()
	defer t.rw.RUnlock()

	var results entrySlice[T]
	for _, tag := range tags {
		entries, ok := t.pairs[tag]
		if !ok {
			continue
		}
		for _, e := range entries {
			e.rw.RLock()
			if !e.deleted {
				results = append(results, e)
			}
			e.rw.RUnlock()
		}
	}
	return results
}
