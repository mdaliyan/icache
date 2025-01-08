package icache

import (
	"testing"
	"time"
)

type item struct {
	ID string
}

func newItem(id string) item {
	return item{
		ID: id,
	}
}

var (
	item1 = newItem("1")
	item2 = newItem("2")
	item3 = newItem("3")
)

func TestNewCache(t *testing.T) {
	p := NewPot[item]()
	p.Set(item1.ID, item1)

	i1, err := p.Get(item1.ID)
	assertNoError(t, err, "Item1 must exist")
	assertEqual(t, item1, i1)
}

func TestClose(t *testing.T) {
	t.Run("pot without ttl", func(t *testing.T) {
		p := new(pot[int])
		p.init()
		time.Sleep(time.Millisecond * 10)

		assertNotNil(t, p.closed)
		err := p.Close()
		assertEqual(t, err, ErrNotClosable, "close() should throw error")
	})

	t.Run("pot with ttl", func(t *testing.T) {
		p := new(pot[int])
		WithTTL(time.Second)(p)
		p.init()
		time.Sleep(time.Millisecond * 10)

		err := p.Close()
		time.Sleep(time.Millisecond * 10)
		assertIsNil(t, err, "close() should not throw error")
	})
}

func TestDrop(t *testing.T) {

	t.Run("by key", func(t *testing.T) {
		p := NewPot[item]()
		p.Set(item1.ID, item1)
		assertEqual(t, 1, p.Len())
		p.Set(item2.ID, item2)
		p.Set(item3.ID, item3)
		assertEqual(t, 3, p.Len())
		p.Drop(item1.ID)
		assertEqual(t, 2, p.Len())
		p.Purge()
		assertEqual(t, 0, p.Len())
	})

	t.Run("one by tag", func(t *testing.T) {
		p := NewPot[item](WithTTL(time.Minute))

		tags1 := []string{"A", "B", "C", "D", "E"}
		p.Set("1", newItem("1"), tags1...)
		assertIsTrue(t, p.Exists("1"))

		tags2 := []string{"A", "B"}
		p.Set("2", newItem("2"), tags2...)

		p.DropTags("C")
		assertIsFalse(t, p.Exists("1"))
		assertIsTrue(t, p.Exists("2"))
	})

	t.Run("multiple by tag", func(t *testing.T) {
		p := NewPot[item]()

		p.Set("1", newItem("1"), "A")
		p.Set("2", newItem("2"), "A", "B")
		p.Set("3", newItem("3"), "B")
		assertIsTrue(t, p.Exists("1"))
		assertIsTrue(t, p.Exists("2"))

		p.DropTags("A")
		assertIsFalse(t, p.Exists("1"))
		assertIsFalse(t, p.Exists("2"))
		assertIsTrue(t, p.Exists("3"))

		entriesWithTagA, err := p.GetByTag("A")
		assertError(t, err, "tag A should contain no items")
		assertEqual(t, 0, len(entriesWithTagA), "tag A should contain no items")

		entriesWithTagB, err := p.GetByTag("B")
		assertNoError(t, err, "tag B should contain one item")
		assertEqual(t, 1, len(entriesWithTagB), "tag B should contain one item")
	})
}

func TestWithTTLOption(t *testing.T) {

	t.Run("initiation", func(t *testing.T) {
		p := new(pot[int])
		WithTTL(time.Second)(p)
		p.init()
		assertEqual(t, time.Second, p.ttl)
	})

	t.Run("overwriting the entry should change ttl", func(t *testing.T) {
		p := NewPot[item](WithTTL(time.Minute))
		p.Set(item1.ID, item1)
		expiresAt, err := p.ExpireTime(item1.ID)
		assertNoError(t, err)

		time.Sleep(time.Microsecond * 50)

		nilExpiresAt, err := p.ExpireTime(item2.ID)
		assertError(t, err, "non-existing entry should have no expiration time")
		assertNotNil(t, nilExpiresAt, "expiration time of missing entry should be nil")

		p.Set(item1.ID, item1)
		newExpiresAt, err := p.ExpireTime(item1.ID)
		assertNoError(t, err, "entry should have expiration time")
		assertNotEqual(t, expiresAt, newExpiresAt, "expiration should be changed")
	})

	t.Run("entries should expire after the given time", func(t *testing.T) {

		p := NewPot[item](WithTTL(3 * time.Second))

		p.Set(item1.ID, item1)

		time.Sleep(time.Second)
		p.Set(item2.ID, item2)

		assertIsTrue(t, p.Exists(item1.ID), "item1 must exist after 1 second")
		assertIsTrue(t, p.Exists(item2.ID), "item2 must exist after insertion")

		time.Sleep(2500 * time.Millisecond)

		p.Set(item3.ID, item3)
		assertIsFalse(t, p.Exists(item1.ID), "item1 must not exist after 2 seconds")
		assertIsTrue(t, p.Exists(item2.ID), "item2 must exist after 1 second")
		assertIsTrue(t, p.Exists(item3.ID), "item3 must exist after insertion")

		time.Sleep(1500 * time.Millisecond)

		assertIsFalse(t, p.Exists(item1.ID), "item1 must not exist after 3 seconds")
		assertIsFalse(t, p.Exists(item2.ID), "item2 must not exist after 2 seconds")
		assertIsTrue(t, p.Exists(item3.ID), "item3 must exist after 1 second")
	})
}
