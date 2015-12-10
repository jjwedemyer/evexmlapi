package cache

import (
	"testing"
	"reflect"
	"time"
)

var (
	key = "key";
	value = []byte("value");
)

func TestRead(t *testing.T) {
	mc := MemoryCache{}.New()	
	mc.Write(key, value, 1000)
	if !reflect.DeepEqual(mc.Read(key),value)  {
		t.Error("MemoryCache key not read correctly.")
	}
}
func TestRead_expired(t *testing.T) {
	mc := MemoryCache{}.New()	
	mc.Write(key, value, 1)
	duration := time.Second * time.Duration(2)
	time.Sleep(duration)
	if mc.Read(key) != nil  {
		t.Error("MemoryCache record did not expire correctly.")
	}
}