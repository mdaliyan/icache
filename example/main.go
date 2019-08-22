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
