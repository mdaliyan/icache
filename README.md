# icache

[![Build Status](https://travis-ci.org/mdaliyan/icache.svg?branch=master)](https://travis-ci.org/mdaliyan/icache)
[![Coverage Status](https://coveralls.io/repos/github/mdaliyan/icache/badge.svg?branch=master)](https://coveralls.io/github/mdaliyan/icache?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mdaliyan/icache?style=flat)](https://goreportcard.com/report/github.com/mdaliyan/icache)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/mdaliyan/icache) 
[![godoc](https://godoc.org/github.com/mdaliyan/icache.svg?status.svg)](https://godoc.org/github.com/mdaliyan/icache)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat)](https://raw.githubusercontent.com/labstack/echo/master/LICENSE)


A cache library for Go with zero GC overhead and high concurrent performance.

go-cache is an in-memory key:value store/cache similar to memcached that
is suitable for applications running on a single process. Its major
advantage is that you don't need to marshal or unmarshal your data or do
type assertion as it doesn't serialize your data and stores values of
your variables so they are thread-safe and you won't face data-race
problem.

Any object can be stored, for a given duration or forever, and the cache
can be safely used by multiple goroutines.

# Installation

```bash
go get github.com/mdaliyan/icache
```

# Usage

```go
package main

import (
	"fmt"
	"github.com/mdaliyan/icache"
	"time"
)

type user struct {
	ID      string
	Name    string
	Age     int
	Contact contact
}

type contact struct {
	Phone   string
	Address string
}

func main() {

	var err error

	// make na new Pot with the default expiration time of 1 Hour
	// set expiretion time to 0 to disable expiration totally
	var pot = icache.NewPot(time.Hour)

	var U = user{
		ID:   "foo",
		Name: "John Doe",
		Age:  30,
		Contact: contact{
			Phone:   "+11111111",
			Address: "localhost",
		},
	}

	// Set the value of the key "foo" to "John Doe" user{}
	pot.Set("foo", U)

	// to get the user{} with "foo" key back you need to pass a pointer
	// to user{}
	var u2 user
	err = pot.Get("foo", &u2)
	fmt.Println(err, u2)

	// if you pass a mismatched type pointer you simply get an error
	var u3 string
	err = pot.Get("foo", &u3)

	fmt.Println(err, u3)

}
```
