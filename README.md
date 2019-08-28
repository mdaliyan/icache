# icache

[![Build Status](https://travis-ci.org/mdaliyan/icache.svg?branch=master)](https://travis-ci.org/mdaliyan/icache)
[![Coverage Status](https://coveralls.io/repos/github/mdaliyan/icache/badge.svg?branch=master)](https://coveralls.io/github/mdaliyan/icache?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mdaliyan/icache?style=flat)](https://goreportcard.com/report/github.com/mdaliyan/icache)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/mdaliyan/icache) 
[![godoc](https://godoc.org/github.com/mdaliyan/icache.svg?status.svg)](https://godoc.org/github.com/mdaliyan/icache)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat)](https://raw.githubusercontent.com/labstack/echo/master/LICENSE)

icache is a cache library for Go with high concurrent access performance. 
it is an in-memory key:value store/cache similar to memcache that
is suitable for applications running on a single machine. Its major
advantage is that it doesn't serialize your data and only stores values of
your variables so they are thread-safe and you won't face data-race problem.

Any object can be stored, for a given duration or forever, and the cache
can be safely used by multiple goroutines.

# Installation

```bash
go get github.com/mdaliyan/icache
```

# Usage

```go 
// make na new Pot with the default expiration time of 1 Hour
// set ttl to 0 to disable expiration entirely
var pot = icache.NewPot(time.Hour) 


var U = user{
    ID:   "foo",
    Name: "John Doe",
    Age:  30,
}

// set user to "foo" key
pot.Set("foo", U)


// get the user previously set to "foo" key into u2 
var u2 user
err = pot.Get("foo", &u2)


// if you pass a mismatched type you simply get a "type mismatch" error
var u3 string
err = pot.Get("foo", &u3)

```

