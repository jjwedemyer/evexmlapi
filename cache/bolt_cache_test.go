package cache

import (
	"reflect"
	"testing"
)

func TestRead_bolt(t *testing.T) {
	ca := NewBoltCache("", 0600, []byte("eve"), nil)
	setup(ca, 1)
	data, err := ca.Read(key)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(data, value) {
		t.Errorf("wanted(%q)\n got(%q).", value, data)
	}
	ca.clear()
}

func TestRead_bolt_expired(t *testing.T) {
	ca := NewBoltCache("", 0600, []byte("eve"), nil)
	setup(ca, -1)
	data, err := ca.Read(key)
	if err != nil {
		t.Error(err)
	}
	if data != nil {
		t.Error("Cache record did not expire correctly.")
	}
	ca.clear()
}

