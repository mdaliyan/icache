package icache

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type item struct {
	ID string
}

func newItem(id string) item {
	return item{
		ID: id,
	}
}

func TestSingleEntryTagDrop(t *testing.T) {
	// var i item
	p := NewPot[item](time.Minute)

	tags1 := []string{"A", "B", "C", "D", "E"}
	p.Set("1", newItem("1"), tags1...)
	assertIsTrue(t, p.Exists("1"))

	tags2 := []string{"A", "B"}
	p.Set("2", newItem("2"), tags2...)

	p.DropTags("C")
	assertIsFalse(t, p.Exists("1"))
	assertIsTrue(t, p.Exists("2"))
}

func TestMultiEntryTagDrop(t *testing.T) {
	// var i item
	p := NewPot[item](time.Minute)

	p.Set("1", newItem("1"), "A")
	p.Set("2", newItem("2"), "A", "B")
	p.Set("e", newItem("6"), "A", "B")
	assertIsTrue(t, p.Exists("1"))
	assertIsTrue(t, p.Exists("2"))

	p.DropTags("A")
	assertIsFalse(t, p.Exists("1"))
	assertIsFalse(t, p.Exists("2"))
}
