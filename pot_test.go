package icache

import (
	"testing"
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

var U = user{
	ID:   "0",
	Name: "John",
	Age:  30,
	Contact: contact{
		Phone:   "+11111111",
		Address: "localhost",
	},
}

func TestUpdateExpireTime(t *testing.T) {
	p := NewPot[user](time.Minute)
	p.Set("1", U)
	expiresAt, err := p.ExpireTime("1")
	assertNoError(t, err, "entry should have expiration time")
	time.Sleep(time.Microsecond * 50)

	nilExpiresAt, err := p.ExpireTime("2")
	assertError(t, err, "entry should have no expiration time")
	assertNotNil(t, nilExpiresAt, "expiration time of missing entry should be nil")

	p.Set("1", U)
	newExpiresAt, err := p.ExpireTime("1")
	assertNoError(t, err, "entry should have expiration time")
	assertNotEqual(t, expiresAt, newExpiresAt, "expiration should be changed")

}

func TestDrop(t *testing.T) {
	p := NewPot[user](0)
	p.Set("1", U)
	assertEqual(t, 1, p.Len())
	p.Set("2", U)
	p.Set("3", U)
	assertEqual(t, 3, p.Len())
	p.Drop("1")
	assertEqual(t, 2, p.Len())
	p.Purge()
	assertEqual(t, 0, p.Len())
}

func TestNewCache(t *testing.T) {
	p := NewPot[user](0)
	p.Set(U.ID, U)
	u1, err := p.Get(U.ID)
	assertNoError(t, err, "cachedUser1 should be found")
	u1.Name = "Jodie"

	cachedUser2, err := p.Get(U.ID)
	assertNoError(t, err, "cachedUser2 should be found")
	assertEqual(t, "John", cachedUser2.Name)
}

func TestAutoExpired(t *testing.T) {
	p := NewPot[user](time.Second * 2)
	user1 := user{Name: "john", ID: "1"}
	user2 := user{Name: "jack", ID: "2"}
	user3 := user{Name: "jane", ID: "3"}

	p.Set(user1.ID, user1)
	time.Sleep(time.Millisecond * 50)
	p.Set(user1.ID, user1)
	time.Sleep(time.Millisecond * 1500)

	p.Set(user3.ID, user3)

	assertIsTrue(t, p.Exists(user1.ID), "user1 should be found")
	_, err := p.Get(user1.ID)
	assertNoError(t, err, "user1 should be found")
	assertIsFalse(t, p.Exists(user2.ID), "user2 should not be found")
	_, err = p.Get(user2.ID)
	assertEqual(t, err, NotFoundErr, "user2 should not be found")
	p.Set(user2.ID, user2)

	time.Sleep(time.Second * 2)

	assertIsTrue(t, p.Exists(user2.ID), "user2 should be found")
	_, err = p.Get(user2.ID)
	assertNoError(t, err, "user2 should be found")
	_, err = p.Get(user1.ID)
	assertError(t, err, "user1 should be expired nowUint")
	assertIsFalse(t, p.Exists(user1.ID), "user1 should be expired after 2 seconds")

	time.Sleep(time.Second * 2)
}
