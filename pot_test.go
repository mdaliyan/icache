package icache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	a := assert.New(t)
	p := NewPot(time.Minute)
	p.Set("1", U)
	expiresAt, err := p.ExpireTime("1")
	a.NoError(err, "entry should have expiration time")
	time.Sleep(time.Microsecond * 50)

	nilExpiresAt, err := p.ExpireTime("2")
	a.Error(err, "entry should have no expiration time")
	a.Nil(nilExpiresAt, "expiration time of missing entry should be nil")

	p.Set("1", U)
	newExpiresAt, err := p.ExpireTime("1")
	a.NoError(err, "entry should have expiration time")
	a.NotEqual(expiresAt, newExpiresAt, "expiration should be changed")

}

func TestGetError(t *testing.T) {
	a := assert.New(t)
	p := NewPot(0)
	p.Set("1", U)
	var cachedUser1 string
	a.Error(p.Get("1", &cachedUser1), "type mismatch error")

	p.Set("2", &U)
	a.Error(p.Get("2", cachedUser1), "needs to pass a pointer")
	a.Error(p.Get("2", nil), "needs to pass a pointer")

}

func TestDrop(t *testing.T) {
	a := assert.New(t)
	p := NewPot(0)
	p.Set("1", U)
	a.Equal(1, p.Len())
	p.Set("2", U)
	p.Set("3", U)
	a.Equal(3, p.Len())
	p.Drop("1")
	a.Equal(2, p.Len())
	p.Purge()
	a.Equal(0, p.Len())
}

func TestNewCache(t *testing.T) {
	a := assert.New(t)
	p := NewPot(0)
	p.Set(U.ID, U)
	var cachedUser1 user
	a.NoError(p.Get(U.ID, &cachedUser1), "cachedUser1 should be found")
	cachedUser1.Name = "Jodie"

	var cachedUser2 user
	a.NoError(p.Get(U.ID, &cachedUser2), "cachedUser2 should be found")
	a.Equal("John", cachedUser2.Name)
}

func TestAutoExpired(t *testing.T) {
	a := assert.New(t)
	p := NewPot(time.Second * 2)
	user1 := user{Name: "john", ID: "1"}
	user2 := user{Name: "jack", ID: "2"}
	user3 := user{Name: "jane", ID: "3"}

	p.Set(user1.ID, user1)
	time.Sleep(time.Millisecond * 50)
	p.Set(user1.ID, user1)
	time.Sleep(time.Millisecond * 1500)

	p.Set(user3.ID, user3)

	var cachedUser user
	a.True(p.Exists(user1.ID), "user1 should be found")
	a.NoError(p.Get(user1.ID, &cachedUser), "user1 should be found")
	a.False(p.Exists(user2.ID), "user2 should not be found")
	a.Error(p.Get(user2.ID, &cachedUser), "user2 should not be found")
	p.Set(user2.ID, user2)

	time.Sleep(time.Second * 2)

	a.True(p.Exists(user2.ID), "user2 should be found")
	a.NoError(p.Get(user2.ID, &cachedUser), "user2 should be found")
	a.Error(p.Get(user1.Name, &user1), "user1 should be expired nowUint")
	a.False(p.Exists(user1.ID), "user1 should be expired after 2 seconds")

	time.Sleep(time.Second * 2)
}
