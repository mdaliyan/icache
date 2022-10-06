# icache

![example workflow](https://github.com/mdaliyan/icache/actions/workflows/test.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/mdaliyan/icache/badge.svg?branch=master)](https://coveralls.io/github/mdaliyan/icache?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mdaliyan/icache)](https://goreportcard.com/report/github.com/mdaliyan/icache)
[![godoc](https://godoc.org/github.com/mdaliyan/icache.svg?status.svg)](https://godoc.org/github.com/mdaliyan/icache)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat)](https://raw.githubusercontent.com/labstack/echo/master/LICENSE)

icache is a no-dependency, generic support cache library for Go with high concurrent access performance.
icache _doean't require serialization_, so the benchmark report you see at the bottom of the page
is what you get at the end, there is no need to marshal and unmarshal anything and waste your
valuable resources.

## Features
- **Thread safe**: icache stores the data in multiple shards and uses a separate mutex for each shard.
- **Generics**: Using generics makes it type-safe.
- **Tags**: icache supports tags for entries and entries with the shared tags can be dropped at the same time.
- **Pointer friendly**: icache stores data as is. any object can be stored in icache, even pointers.
- **TTL**: ttl can be set to any thing from 1 seconds to infinity.

## Installation

icache v2 requires go v1.8.0 and above.
```bash
go get github.com/mdaliyan/icache/v2
```

### previous version for Go < 1.18

icache v1 is compatible with Go < `1.18`, and it also doean't require serialization.
it uses reflections to make sure about your data type.

Follow instructions at [version v1.x.x](https://github.com/mdaliyan/icache/tree/v1) for installation and usage. 
V1 will be maintained separately.

## Usage

```go 
// make a new Pot:
// - to store user structs 
// - with expiration time of 1 Hour
var pot = icache.NewPot[user](
          icache.WithTTL(time.Hour),
    )

var User = user{
    ID:   "foo",
    Name: "John Doe",
    Age:  30,
}

// set user to "foo" key
pot.Set("foo", User)

// get the user previously set to "foo" key into user1 
user1, err := pot.Get("foo")

// you also can add tags to your entries
pot.Set("foo", User, "tag1", "tag2")
pot.Set("faa", User, "tag1")
pot.Set("fom", User, "tag3")

// and delete multiple entries at once
pot.DropTags("tag1")
```

## Invalidation Strategies
- By **Key**
- By **Tag**
- By **TTL**: Using the `WithTTL()` option.

## Benchmarks
```bash
goos: darwin
goarch: amd64
pkg: github.com/mdaliyan/icache/v2
cpu: Intel(R) Core(TM) i7-8557U CPU @ 1.70GHz
BenchmarkICache
BenchmarkICache-8             	9570999	    118.1 ns/op	   0 B/op	  0 allocs/op
BenchmarkICacheConcurrent
BenchmarkICacheConcurrent-8   	6117471	    191.4 ns/op	   0 B/op	  0 allocs/op
```

## Plans for later
- Different invalidation Strategies can be added if it's needed. 
- Maybe Remote sync
