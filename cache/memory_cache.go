package cache

import "time"

type MemoryCache struct {
	_cache map[string]RecordCache
}

func (mc MemoryCache) New() *MemoryCache {
	return &MemoryCache{_cache: make(map[string]RecordCache)}
}

func (mc MemoryCache) Read(key string) []byte {
	val, ok := mc._cache[key]
	if ok && val.expireTime > time.Now().Unix() {		
		return val.value
	}
	return nil
}

func (mc *MemoryCache) Write(key string, value []byte, duration int64) {
	expireT := duration + time.Now().Unix()
	mc._cache[key] = RecordCache{value: value, expireTime: expireT}
}

