package icache

import (
	"math/rand"
	"os"
	"runtime"
	"testing"
	"time"
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
	icache = NewPot[user](time.Hour)
	for i := 0; i < 10000; i++ {
		id := randomString()
		ids = append(ids, id)
		U.ID = id
		U.Age = rand.Intn(70)
		icache.Set(id, U)
	}
	idsLen = len(ids) - 1
	os.Exit(m.Run())
}

func randomID() string {
	return ids[rand.Intn(idsLen)]
}

var idsLen int
var ids []string
var icache Pot[user]

func BenchmarkICache(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		get()
	}
}

var ut user

func get() {
	ut, _ = icache.Get(randomID())
}

func BenchmarkICacheConcurrent(b *testing.B) {
	b.SetParallelism(100)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			get()
		}
	})
}
