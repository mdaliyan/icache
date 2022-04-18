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
	p := NewPot[user](time.Minute)
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

func TestDrop(t *testing.T) {
	a := assert.New(t)
	p := NewPot[user](0)
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
	p := NewPot[user](0)
	p.Set(U.ID, U)
	u1, err := p.Get(U.ID)
	a.NoError(err, "cachedUser1 should be found")
	u1.Name = "Jodie"

	cachedUser2, err := p.Get(U.ID)
	a.NoError(err, "cachedUser2 should be found")
	a.Equal("John", cachedUser2.Name)
}

func TestAutoExpired(t *testing.T) {
	a := assert.New(t)
	p := NewPot[user](time.Second * 2)
	user1 := user{Name: "john", ID: "1"}
	user2 := user{Name: "jack", ID: "2"}
	user3 := user{Name: "jane", ID: "3"}

	p.Set(user1.ID, user1)
	time.Sleep(time.Millisecond * 50)
	p.Set(user1.ID, user1)
	time.Sleep(time.Millisecond * 1500)

	p.Set(user3.ID, user3)

	a.True(p.Exists(user1.ID), "user1 should be found")
	_, err := p.Get(user1.ID)
	a.NoError(err, "user1 should be found")
	a.False(p.Exists(user2.ID), "user2 should not be found")
	_, err = p.Get(user2.ID)
	a.Equal(err, NotFoundErr, "user2 should not be found")
	p.Set(user2.ID, user2)

	time.Sleep(time.Second * 2)

	a.True(p.Exists(user2.ID), "user2 should be found")
	_, err = p.Get(user2.ID)
	a.NoError(err, "user2 should be found")
	_, err = p.Get(user1.ID)
	a.Error(err, "user1 should be expired nowUint")
	a.False(p.Exists(user1.ID), "user1 should be expired after 2 seconds")

	time.Sleep(time.Second * 2)
}
