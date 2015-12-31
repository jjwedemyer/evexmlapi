package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type FileCache struct {
	directory string
	prefix    string
	cache
}

// NewFileCache returns a new FileCache struct
// with default directory and prefix.
func NewFileCache() *FileCache {
	dir := "/tmp/eve-xmlapi-go"
	prefix := ""
	return &FileCache{cache: cache{}, directory: dir, prefix: prefix}
}

// Read looks for a valid cache
func (c FileCache) Read(urlPath string) ([]byte, error) {
	file := c.path(urlPath)
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
func (c *FileCache) Store(urlPath string, value []byte, duration int64) error {
	expireT := duration + c.now()
	record := RecordCache{Value: value, ExpireTime: expireT}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.path(urlPath), data, os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		return err
	}
	return nil
}

func (c FileCache) clear() {
	os.RemoveAll(c.directory)
	os.Mkdir(c.directory, os.ModePerm)
}

func (c *FileCache) path(urlPath string) string {
	hash := sha256.Sum256([]byte(urlPath))
	hashStr := fmt.Sprintf("%x", hash)
	file := c.prefix + hashStr
	return c.directory + "/" + file
}
