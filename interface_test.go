package icache_test

import (
	"reflect"
)

type Data struct {
	ID   string
	Name string
	Age  int
}

var d = Data{
	ID:   "0",
	Name: "John",
	Age:  30,
}

func dataSender(i interface{}) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		panic("wrong")
	}
	v.Addr()
}

//func TestSetDataTroughInterface(t *testing.T) {
//	var i *Data
//	dataSender(i)
//	assert.Equal(t, "John", i.Name)
//}
