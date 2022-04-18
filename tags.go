package icache

import (
	"sync"
)

type tags[T any] struct {
	pairs map[uint64]entries[T]
	rw    sync.RWMutex
	pot   *pot[T]
}

func (t *tags[T]) purge() {
	t.rw.Lock()
	t.pairs = make(map[uint64]entries[T])
	t.rw.Unlock()
}

func (t *tags[T]) add(e *entry[T]) {
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
		tags := make(entries[T])
		tags[e.key] = e
		t.pairs[tag] = tags
	}
	t.rw.Unlock()
}

func (t *tags[T]) drop(e *entry[T]) {
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

func (t *tags[T]) getEntries(tag uint64) (entries entrySlice[T]) {
	t.rw.RLock()
	if _, ok := t.pairs[tag]; !ok {
		t.rw.RUnlock()
		return
	}
	for _, e := range t.pairs[tag] {
		entries = append(entries, e)
	}
	t.rw.RUnlock()
	return
}

func (t *tags[T]) dropTags(tags ...uint64) {
	for _, tag := range tags {
		var entries entrySlice[T]
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
