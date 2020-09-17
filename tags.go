package icache

import (
	`sort`
	`sync`
)

type tags struct {
	pairs map[string]entrySlice
	rw    sync.RWMutex
	pot   *pot
}

func (t *tags) purge(p *pot) {
	t.rw.Lock()
	t.pairs = make(map[string]entrySlice)
	t.pot = p
	t.rw.Unlock()
}

func (t *tags) add(e *entry) {
	if e == nil || e.tags == nil {
		return
	}
	t.rw.Lock()
	for _, tag := range e.tags {
		var infos, ok = t.pairs[tag]
		if !ok {
			t.pairs[tag] = entrySlice{e}
			continue
		}
		idx := sort.Search(len(infos), func(i int) bool {
			return e.key >= infos[i].key
		})
		if infos[idx].key != e.key {
			t.pairs[tag] = append(infos[:idx], append(entrySlice{e},infos[idx:]...)...)
		}
	}
	t.rw.Unlock()
}

func (t *tags) drop(e *entry) {
	if e == nil || e.tags == nil {
		return
	}
	t.rw.Lock()
	for _, tag := range e.tags {
		var entries = t.pairs[tag]
		idx := sort.Search(len(entries), func(i int) bool {
			return e.key >= entries[i].key
		})
		if entries[idx].key == e.key {
			t.pairs[tag] = append(entries[:idx], entries[idx+1:]...)
		}
		if len(t.pairs[tag]) == 0 {
			delete(t.pairs, tag)
		}
	}
	t.rw.Unlock()
}

func (t *tags) getEntries(tag string) (entries entrySlice) {
	t.rw.RLock()
	entries, _ = t.pairs[tag]
	t.rw.RUnlock()
	return
}

func (t *tags) dropTags(tags ...string) {
	for _, tag := range tags {
		t.rw.Lock()
		entries, _ := t.pairs[tag]
		t.rw.Unlock()
		t.pot.dropEntries(entries)
	}
}
