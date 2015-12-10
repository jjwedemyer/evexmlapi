package cache

type Cache interface {
	Read(key string) []byte
	Write(key string, value []byte, duration int64)	
}

type RecordCache struct {
	value []byte
	expireTime int64
}