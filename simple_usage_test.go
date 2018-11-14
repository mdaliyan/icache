package iCache

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"sync"
	"github.com/allegro/bigcache"
	"github.com/coocood/freecache"
	"encoding/json"
	"time"
)

type User struct {
	Name string
	Age  int8
}

var U = User{
	Name: "Ali",
	Age:  8,
}

func TestNewCache(t *testing.T) {
	a := assert.New(t)
	c := NewPot()
	c.Set("user1", U, 0)
	var u User
	c.Get("user1", &u)
	u.Name = "jsddjkd"

	var u2 User
	c.Get("user1", &u2)
	a.Equal("Ali", u2.Name)
}

func Benchmark_GR_InterfaceCache(b *testing.B) {
	c := NewPot()
	c.Set("userID", U, 0)
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
	c := NewPot()
	c.Set("userID", U, 0)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var ut User
		c.Get("userID", &ut)
	}
}

func Benchmark_GR_Freecache(b *testing.B) {
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

func BenchmarkFreecache(b *testing.B) {
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

func Benchmark_GR_Bigcache(b *testing.B) {
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

func BenchmarkBigcache(b *testing.B) {
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
