package iCache

import (
	"reflect"
	"sync"
	"time"
)

var nowUint int64
var nowLock sync.RWMutex

func now() (n int64){
	nowLock.Lock()
	n = nowUint
	nowLock.Unlock()
	return
}

func init() {
	go func() {
		for {
			nowLock.Lock()
			nowUint = time.Now().Unix()
			nowLock.Unlock()
			time.Sleep(time.Second)
		}
	}()
}

type expireTime struct {
	Key       uint64
	ExpiresAt int64
}

type entry struct {
	Data reflect.Value
}