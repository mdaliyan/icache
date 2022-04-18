# icache

![example workflow](https://github.com/mdaliyan/icache/actions/workflows/test.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/mdaliyan/icache/badge.svg?branch=master)](https://coveralls.io/github/mdaliyan/icache?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mdaliyan/icache)](https://goreportcard.com/report/github.com/mdaliyan/icache)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/mdaliyan/icache) 
[![godoc](https://godoc.org/github.com/mdaliyan/icache.svg?status.svg)](https://godoc.org/github.com/mdaliyan/icache)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat)](https://raw.githubusercontent.com/labstack/echo/master/LICENSE)

icache is a no dependency generic cache library for Go with high concurrent access performance. 
Its major advantages over the other libraries is that it doesn't serialize the data and only 
holds the values as they are. As the result you won't need to unmarshal anything when you get
the values. This saves you time and resources.

Any object can be stored in icache, even pointers, for a given duration or forever,
and the cache can be safely used by multiple goroutines.

# Installation

```bash
go get github.com/mdaliyan/icache/v2
```

# Usage

```go 
// make a new Pot:
// - to store user structs 
// - with global expiration time of 1 Hour
// * set ttl to 0 to disable expiration entirely
var pot = icache.NewPot[user](time.Hour) 

var U = user{
    ID:   "foo",
    Name: "John Doe",
    Age:  30,
}

// set user to "foo" key
pot.Set("foo", U)

// get the user previously set to "foo" key into u2 
u2, err = pot.Get("foo")

// you also can add tags to your entries
pot.Set("foo", U, "tag1", "tag2")
pot.Set("faa", U, "tag1")
pot.Set("fom", U, "tag3")

// and delete multiple entries at once
pot.DropTags("tag1")
```

I might add MGet method later

# Benchmarking
```bash
goos: darwin
goarch: amd64
pkg: github.com/mdaliyan/icache/v2
cpu: Intel(R) Core(TM) i7-8557U CPU @ 1.70GHz
BenchmarkICache
BenchmarkICache-8             	 9570999	       118.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkICacheConcurrent
BenchmarkICacheConcurrent-8   	 6117471	       191.4 ns/op	       0 B/op	       0 allocs/op
```
