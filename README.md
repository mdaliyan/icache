# icache

[![Build Status](https://travis-ci.org/mdaliyan/icache.svg?branch=master)](https://travis-ci.org/mdaliyan/icache)
[![Coverage Status](https://coveralls.io/repos/github/mdaliyan/icache/badge.svg?branch=master)](https://coveralls.io/github/mdaliyan/icache?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mdaliyan/icache?style=flat)](https://goreportcard.com/report/github.com/mdaliyan/icache)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/mdaliyan/icache) 
[![godoc](https://godoc.org/github.com/mdaliyan/icache.svg?status.svg)](https://godoc.org/github.com/mdaliyan/icache)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat)](https://raw.githubusercontent.com/labstack/echo/master/LICENSE)

icache is a cache library for Go with high concurrent access performance. 
Its major advantages over the other libraries is that it doesn't serialize
your data and only stores values. As the result you won't need to unmarshal
anything when you get the values from it. This saves you time and resources.

Any object can be stored, for a given duration or forever, and the cache
can be safely used by multiple goroutines.

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

// you also can add tags to your entries
pot.Set("foo", U, "tag1", "tag2")
pot.Set("faa", U, "tag1")
pot.Set("fom", U, "tag3")

// and delete multiple entries at once
pot.DropTags("tag1")
```

I might add MGet method later
