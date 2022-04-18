package icache

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, expected, actual interface{}, msg ...string) {
	if expected != actual {
		failAssertion(t, false, expected, actual)
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
	fmt.Println(fmt.Sprintf("%v", expected))
	fmt.Println(fmt.Sprintf("actual: %v", actual))
	t.Fail()
}
