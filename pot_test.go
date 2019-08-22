package iCache_test

import (
	"encoding/json"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/allegro/bigcache"
	"github.com/coocood/freecache"
	. "github.com/mdaliyan/icache"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID      string
	Name    string
	Age     int
	Contact Contact
}

type Contact struct {
	Phone   string
	Address string
}

var U = User{
	ID:   "0",
	Name: "John",
	Age:  30,
	Contact: Contact{
		Phone:   "+11111111",
		Address: "localhost",
	},
}

func TestGetError(t *testing.T) {
	a := assert.New(t)
	p := NewPot(0)
	p.Set("1", U)
	var cachedUser1 string
	a.Error(p.Get("1", &cachedUser1), "type mismatch error")

	p.Set("2", &U)
	a.Error(p.Get("2", cachedUser1), "needs to pass a pointer")
	a.Error(p.Get("2", nil), "needs to pass a pointer")

}

func TestDrop(t *testing.T) {
	a := assert.New(t)
	p := NewPot(0)
	p.Set("1", U)
	a.Equal(1, p.Len())
	p.Set("2", U)
	p.Set("3", U)
	a.Equal(3, p.Len())
	p.Drop("1")
	a.Equal(2, p.Len())
	p.Purge()
	a.Equal(0, p.Len())
}

func TestNewCache(t *testing.T) {
	a := assert.New(t)
	p := NewPot(0)
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
	p := NewPot(time.Second * 2)
	user1 := User{Name: "john", ID: "1"}
	user2 := User{Name: "jack", ID: "2"}
	user3 := User{Name: "jane", ID: "3"}

	p.Set(user1.ID, user1)
	time.Sleep(time.Millisecond * 1500)

	p.Set(user3.ID, user3)

	var cachedUser User
	a.True(p.Exists(user1.ID), "user1 should be found")
	a.NoError(p.Get(user1.ID, &cachedUser), "user1 should be found")
	a.False(p.Exists(user2.ID), "user2 should not be found")
	a.Error(p.Get(user2.ID, &cachedUser), "user2 should not be found")
	p.Set(user2.ID, user2)

	time.Sleep(time.Second * 2)

	a.True(p.Exists(user2.ID), "user2 should be found")
	a.NoError(p.Get(user2.ID, &cachedUser), "user2 should be found")
	a.Error(p.Get(user1.Name, &user1), "user1 should be expired nowUint")
	a.False(p.Exists(user1.ID), "user1 should be expired after 2 seconds")

	time.Sleep(time.Second * 2)
}

// =============================================================================================
//	Benchmarks
// =============================================================================================

func randomString() string {
	n := rand.Intn(10) + 8
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func TestMain(m *testing.M) {
	icache = NewPot(time.Hour)
	bigCache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	freeCache = freecache.NewCache(100 * 100)
	for i := 0; i < 10000; i++ {
		id := randomString()
		ids = append(ids, id)
		U.ID = id
		U.Age = rand.Intn(70)
		Ujson, _ := json.Marshal(U)
		icache.Set(id, U)
		freeCache.Set([]byte(id), Ujson, int(time.Hour.Seconds()))
		bigCache.Set(id, Ujson)
	}
	idsLen = len(ids) - 1
	os.Exit(m.Run())
}

func randomID() string {
	return ids[rand.Intn(idsLen)]
}

var idsLen int
var ids []string
var icache Pot
var freeCache *freecache.Cache
var bigCache *bigcache.BigCache

func Benchmark_iCache(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getFromICache()
	}
}

func Benchmark_FreeCache(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getFromFreeCache()
	}
}

func Benchmark_BigCache(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getFromBigCache()
	}
}

func getFromICache() {
	var ut User
	icache.Get(randomID(), &ut)
}

func getFromFreeCache() {
	var ut User
	byt, err := freeCache.Get([]byte(randomID()))
	if err == nil {
		json.Unmarshal(byt, &ut)
	}
}

func getFromBigCache() {
	var ut User
	byt, err := bigCache.Get(randomID())
	if err == nil {
		json.Unmarshal(byt, &ut)
	}
}

func Benchmark_iCache_Concurrent(b *testing.B) {
	b.SetParallelism(100)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			getFromICache()
		}
	})
}

func Benchmark_FreeCache_Concurrent(b *testing.B) {
	b.SetParallelism(100)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			getFromFreeCache()
		}
	})
}

func Benchmark_BigCache_Concurrent(b *testing.B) {
	b.SetParallelism(100)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			getFromBigCache()
		}
	})
}
