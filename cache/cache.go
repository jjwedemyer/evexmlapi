package cache

import "time"

type Cache interface {
	Read(key string) ([]byte, error)
	Store(key string, value []byte, duration int64) error
}

type RecordCache struct {
	Value      []byte	`json: "value"`
	ExpireTime int64	`json: "expireTime"`
}

type cache struct{}

func (c cache) now() int64 {
	return time.Now().Unix()
}
