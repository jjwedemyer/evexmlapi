package cache

import (
	
	"os"
	"reflect"
	"testing"
)

var (
	key   = "key12345"
	value = []byte("value")
)

func TestRead_file(t *testing.T) {
	ca := setup(1)	
	data, err := ca.Read(key)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(data, value) {
		t.Errorf("wanted(%q)\n got(%q).", value, data)
	}
	ca.clear()	
}

func TestRead_file_expired(t *testing.T) {
	ca := setup(-1)	
	data, err := ca.Read(key)
	if err != nil {
		t.Error(err)
	}
	if data != nil {
		t.Error("Cache record did not expire correctly.")
	}
	ca.clear()
}

func setup(duration int64) *FileCache {
	ca := NewFileCache()
	os.MkdirAll(ca.directory, os.ModePerm)
	ch2 := make(chan int)
	go func() {
		ch2 <- 1		
		_ = ca.Store(key, value, duration)
	}()
	<-ch2
	return ca
}

