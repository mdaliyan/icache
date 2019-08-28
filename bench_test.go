package icache

import (
	"encoding/json"
	"math/rand"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/allegro/bigcache"
	"github.com/coocood/freecache"
)

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
	runtime.GOMAXPROCS(runtime.NumCPU())
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
	var ut user
	icache.Get(randomID(), &ut)
}

func getFromFreeCache() {
	var ut user
	byt, err := freeCache.Get([]byte(randomID()))
	if err == nil {
		json.Unmarshal(byt, &ut)
	}
}

func getFromBigCache() {
	var ut user
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
