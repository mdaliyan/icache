# icache

[![Build Status](https://travis-ci.org/mdaliyan/icache.svg?branch=master)](https://travis-ci.org/mdaliyan/icache)
[![godoc](https://godoc.org/github.com/mdaliyan/icache.svg?status.svg)](https://godoc.org/github.com/mdaliyan/icache)
[![Coverage Status](https://coveralls.io/repos/github/mdaliyan/icache/badge.svg?branch=master)](https://coveralls.io/github/mdaliyan/icache?branch=master)


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
	"github.com/mdaliyan/icache"
	"time"
	"fmt"
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

func main() {

	var err error

	// make na new Pot with the default expiration time of 1 Hour
	// set expiretion time to 0 to disable expiration totally
	var pot = iCache.NewPot(time.Hour)

	var U = User{
		ID:   "foo",
		Name: "John Doe",
		Age:  30,
		Contact: Contact{
			Phone:   "+11111111",
			Address: "localhost",
		},
	}

	// Set the value of the key "foo" to "John Doe" User{}
	pot.Set("foo", U)

	// to get the User{} with "foo" key back you need to pass a pointer
	// to User{}
	var u2 User
	err = pot.Get("foo", &u2)
	fmt.Println(err, u2)

	// if you pass a mismatched type pointer you simply get an error
	var u3 string
	err = pot.Get("foo", &u3)

	fmt.Println(err, u3)

}

```
