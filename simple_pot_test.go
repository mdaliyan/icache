package iCache_test

import (
	"encoding/json"
	"github.com/allegro/bigcache"
	"github.com/coocood/freecache"
	. "github.com/mdaliyan/icache"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type User struct {
	Name string
	ID   string
}

var U = User{
	ID:   "0",
	Name: "John",
}

func TestNewCache(t *testing.T) {
	a := assert.New(t)
	p := NewPot(Config{})
	p.Set(U.ID, U)
	var cachedUser1 User
	a.NoError(p.Get(U.ID, &cachedUser1), "cachedUser1 should be found")
	cachedUser1.Name = "Jodie"

	var cachedUser2 User
	a.NoError(p.Get(U.ID, &cachedUser2), "cachedUser2 should be found")
	a.Equal("John", cachedUser2.Name)
}

func TestAutoExpired(t *testing.T) {
	a := assert.New(t)
	p := NewPot(Config{TTL:time.Second * 2})
	user1 := User{Name: "john", ID: "1"}
	user2 := User{Name: "jack", ID: "2"}
	user3 := User{Name: "jane", ID: "3"}

	p.Set(user1.ID, user1)
	time.Sleep(time.Millisecond*1500)

	p.Set(user3.ID, user3)

	var cachedUser User
	a.NoError(p.Get(user1.ID, &cachedUser), "first user should be found")
	a.Error(p.Get(user2.ID, &cachedUser), "second user should not be found")
	p.Set(user2.ID, user2)

	time.Sleep(time.Second*2)

	a.NoError(p.Get(user2.ID, &cachedUser), "second user should be found")
	a.Error(p.Get(user1.Name, &user1), "user1 should be expired nowUint")

	time.Sleep(time.Second*2)
}

func Benchmark_GR_InterfaceCache(b *testing.B) {
	c := NewPot(Config{TTL:time.Minute})
	c.Set("userID", U)
	wg := sync.WaitGroup{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			var ut User
			c.Get("userID", &ut)
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkInterfaceCache(b *testing.B) {
	c := NewPot(Config{TTL:time.Minute})
	c.Set("userID", U)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var ut User
		c.Get("userID", &ut)
	}
}

func Benchmark_GR_FreeCache(b *testing.B) {
	c := freecache.NewCache(100 * 100)
	key := []byte("userID")
	by, _ := json.Marshal(U)
	c.Set(key, by, 0)
	wg := sync.WaitGroup{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			var ut User
			byt, _ := c.Get(key)
			json.Unmarshal(byt, &ut)
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkFreeCache(b *testing.B) {
	c := freecache.NewCache(100 * 100)
	key := []byte("userID")
	by, _ := json.Marshal(U)
	c.Set(key, by, 0)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var ut User
		byt, _ := c.Get(key)
		json.Unmarshal(byt, &ut)
	}
}

func Benchmark_GR_BigCache(b *testing.B) {
	c, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	by, _ := json.Marshal(U)
	c.Set("userID", by)
	wg := sync.WaitGroup{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			var ut User
			byt, _ := c.Get("userID")
			json.Unmarshal(byt, &ut)
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkBigCache(b *testing.B) {
	c, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	by, _ := json.Marshal(U)
	c.Set("userID", by)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var ut User
		byt, _ := c.Get("userID")
		json.Unmarshal(byt, &ut)
	}
}
