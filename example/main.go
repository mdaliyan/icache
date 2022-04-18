package main

import (
	"fmt"
	"github.com/mdaliyan/icache/v2"
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

	// make na new Pot with the default expiration time of 1 Hour to store user structs
	// set expiretion time to 0 to disable expiration totally
	var pot = icache.NewPot[user](time.Hour)

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

	// get the user{} with "foo" key
	u2, err := pot.Get("foo")
	fmt.Println(err, u2)

	// you also can add tags to your entries
	pot.Set("foo", U, "tag1", "tag2")
	pot.Set("faa", U, "tag1")
	pot.Set("fom", U, "tag3")

	// and delete multiple entries at once
	pot.DropTags("tag1")
}
