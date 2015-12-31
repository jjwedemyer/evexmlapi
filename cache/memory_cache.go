package cache

type MemoryCache struct {
	_cache map[string]RecordCache
	cache  cache
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{_cache: make(map[string]RecordCache), cache: cache{}}
}

func (c MemoryCache) Read(urlPath string) ([]byte, error) {
	val, ok := c._cache[urlPath]
	if ok && val.ExpireTime > c.cache.now() {
		return val.Value, nil
	}
	return nil, nil
}

func (c *MemoryCache) Store(urlPath string, value []byte, duration int64) error {
	expireT := duration + c.cache.now()
	c._cache[urlPath] = RecordCache{Value: value, ExpireTime: expireT}
	return nil
}
