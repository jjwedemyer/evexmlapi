package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type FileCache struct {
	directory string
	prefix    string
	cache
}

// NewFileCache returns a new FileCache struct
// with default directory and prefix.
func NewFileCache(dir string, prefix string) *FileCache {
	if dir == "" {
		dir = "./cache_files"
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err1 := os.MkdirAll(dir, os.ModePerm)
		if err1 != nil {
			log.Panic(err1)
		}
	}
	return &FileCache{cache: cache{}, directory: dir, prefix: prefix}
}

// Read looks for a valid cache
func (c FileCache) Read(key string) ([]byte, error) {
	file := c.path(key)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, nil
	}
	record, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	val := new(RecordCache)
	err = json.Unmarshal(record, val)
	if err != nil {
		return nil, err
	}
	if val.ExpireTime > c.now() {
		return val.Value, nil
	}
	os.Remove(file)
	return nil, nil
}

// Store caches the value as JSON with the expire time
func (c *FileCache) Store(key string, value []byte, duration int64) error {
	expireT := duration + c.now()
	record := RecordCache{Value: value, ExpireTime: expireT}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.path(key), data, os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		return err
	}
	return nil
}

func (c FileCache) clear() {
	os.RemoveAll(c.directory)	
}

func (c *FileCache) path(urlPath string) string {
	hash := sha256.Sum256([]byte(urlPath))
	hashStr := fmt.Sprintf("%x", hash)
	file := c.prefix + hashStr
	return c.directory + "/" + file
}
