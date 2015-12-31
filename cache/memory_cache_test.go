package cache

import (
	"reflect"
	"testing"
)

func TestRead_memory(t *testing.T) {
	mc := NewMemoryCache()
	mc.Store(key, value, 1000)
	data, _ := mc.Read(key)
	if !reflect.DeepEqual(data, value) {
		t.Error("MemoryCache key not read correctly.")
	}
}
func TestRead_memory_expired(t *testing.T) {
	mc := NewMemoryCache()
	mc.Store(key, value, -1)
	data, _ := mc.Read(key)
	if data != nil {
		t.Error("MemoryCache record did not expire correctly.")
	}
}
