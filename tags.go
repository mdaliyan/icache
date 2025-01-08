package icache

type tags[T any] struct {
	pairs map[uint64]entries[T]
}

func (t *tags[T]) purge() {
	t.pairs = make(map[uint64]entries[T])
}

func (t *tags[T]) add(e *entry[T]) {
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

	if len(entries) == 0 {
		delete(t.pairs, tag)
	}
}

func (t *tags[T]) getEntriesWithTags(tags ...uint64) entrySlice[T] {
	var results entrySlice[T]
	for _, tag := range tags {
		entries, ok := t.pairs[tag]
		if !ok {
			continue
		}
		for _, e := range entries {
			if e == nil {
				continue
			}
			e.rw.RLock()
			if !e.deleted {
				results = append(results, e)
			}
			e.rw.RUnlock()
		}
	}
	return results
}
