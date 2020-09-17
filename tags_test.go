package icache

import (
	`math/rand`
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	a := assert.New(t)
	p := new(pot)
	p.init(time.Minute)

	tags1 := []string{"A", "B", "C", "D", "E"}
	p.Set("1", newItem("1"), tags1...)
	a.True(p.Exists("1"))

	tags2 := []string{"A", "B"}
	p.Set("2", newItem("2"), tags2...)

	p.DropTags("C")
	a.False(p.Exists("1"))
	a.True(p.Exists("2"))
}

func TestMultiEntryTagDrop(t *testing.T) {
	// var i item
	a := assert.New(t)
	p := new(pot)
	p.init(time.Minute)

	p.Set("1", newItem("1"), "A")
	p.Set("2", newItem("2"), "A", "B")
	a.True(p.Exists("1"))
	a.True(p.Exists("2"))

	p.DropTags("A")
	a.False(p.Exists("1"))
	a.False(p.Exists("2"))
}
