package icache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestKeyGenerator(t *testing.T) {

	t.Run("KeyGen", func(t *testing.T) {
		hash, shardID := keyGen("icache")
		assertEqual(t, uint64(16773551877005858910), hash)
		assertEqual(t, uint64(94), shardID)
		hash, shardID = keyGen("icache2")
		assertEqual(t, uint64(1192961860816945028), hash)
		assertEqual(t, uint64(132), shardID)
	})

	t.Run("TagKeyGen", func(t *testing.T) {
		hash := tagKeyGen("icache", "icache2")
		assertEqual(t, uint64(16773551877005858910), hash[0])
		assertEqual(t, uint64(1192961860816945028), hash[1])
	})
}

func assertEqual(t *testing.T, expected, actual interface{}, msg ...string) {
	if expected != actual {
		failAssertion(t, false, expected, actual, msg...)
	}
}

func assertNotEqual(t *testing.T, expected, actual interface{}, msg ...string) {
	if expected == actual {
		failAssertion(t, true, expected, actual, msg...)
	}
}

func assertIsTrue(t *testing.T, value bool, msg ...string) {
	assertEqual(t, true, value, msg...)
}

func assertIsFalse(t *testing.T, value bool, msg ...string) {
	assertEqual(t, false, value, msg...)
}

func assertError(t *testing.T, err error, msg ...string) {
	assertNotNil(t, err, msg...)
}

func assertNoError(t *testing.T, err error, msg ...string) {
	assertIsNil(t, err, msg...)
}

func assertIsNil(t *testing.T, i interface{}, msg ...string) {
	if i != nil {
		failAssertion(t, false, nil, i, msg...)
	}
}

func assertNotNil(t *testing.T, i interface{}, msg ...string) {
	if i == nil {
		failAssertion(t, true, nil, i, msg...)
	}
}

func failAssertion(t *testing.T, not bool, expected, actual interface{}, msg ...string) {
	if len(msg) != 0 {
		fmt.Println("assertion: ", msg)
	}
	if not {
		fmt.Print("not expected:")
	} else {
		fmt.Print("expected:")
	}
	fmt.Printf("%v\n", expected)
	fmt.Printf("actual: %v\n", actual)
	t.Fail()
}
