package icache

import (
	"fmt"
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

func assertValueIsTrue(t *testing.T, value bool) {
	assertValueIs(t, true, value)
}
func assertValueIsFalse(t *testing.T, value bool) {
	assertValueIs(t, false, value)
}

func assertValueIs(t *testing.T, expected, actual bool) {
	if expected != actual {
		fmt.Printf("expected: %v\n", expected)
		fmt.Printf("actual: %v\n", actual)
		fmt.Printf("actual: %v\n", actual)
		t.Fail()
	}
}

func TestSingleEntryTagDrop(t *testing.T) {
	// var i item
	p := new(pot)
	p.init(time.Minute)

	tags1 := []string{"A", "B", "C", "D", "E"}
	p.Set("1", newItem("1"), tags1...)
	assertValueIsTrue(t, p.Exists("1"))

	tags2 := []string{"A", "B"}
	p.Set("2", newItem("2"), tags2...)

	p.DropTags("C")
	assertValueIsFalse(t, p.Exists("1"))
	assertValueIsTrue(t, p.Exists("2"))
}

func TestMultiEntryTagDrop(t *testing.T) {
	// var i item
	p := new(pot)
	p.init(time.Minute)

	p.Set("1", newItem("1"), "A")
	p.Set("2", newItem("2"), "A", "B")
	p.Set("e", newItem("6"), "A", "B")
	assertValueIsTrue(t, p.Exists("1"))
	assertValueIsTrue(t, p.Exists("2"))

	p.DropTags("A")
	assertValueIsFalse(t, p.Exists("1"))
	assertValueIsFalse(t, p.Exists("2"))
}
