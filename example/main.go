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
