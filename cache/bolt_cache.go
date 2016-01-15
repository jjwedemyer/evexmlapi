package cache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

type BoltCache struct {
	db             *bolt.DB
	requestsBucket []byte
	path           string
	cache
}

func NewBoltCache(path string, port os.FileMode, bucketName []byte, options *bolt.Options) *BoltCache {
	if path == "" {
		path = "./cache.db"
	}
	if options == nil {
		options = &bolt.Options{Timeout: 1 * time.Second}
	}
	db, err := bolt.Open(path, port, options)
	if err != nil {
		log.Panic(err)
	}	
	return &BoltCache{db: db, requestsBucket: bucketName, path: path}
}

func (c BoltCache) Read(key string) ([]byte, error) {
	var value []byte
	err := c.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(c.requestsBucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", c.requestsBucket)
		}

		var buffer bytes.Buffer
		buffer.Write(bucket.Get([]byte(key)))

		record := buffer.Bytes()
		val := new(RecordCache)
		err := json.Unmarshal(record, val)
		if err != nil {
			return err
		}
		if val.ExpireTime > c.now() {
			value = val.Value
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *BoltCache) Store(key string, value []byte, duration int64) error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(c.requestsBucket)
		if err != nil {
			return err
		}
		expireT := duration + c.now()
		record := RecordCache{Value: value, ExpireTime: expireT}
		data, err := json.Marshal(record)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(key), data)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (c *BoltCache) clear() {
	c.Close()
	os.Remove(c.path)
}

func (c *BoltCache) Close(){
	if c.db != nil {
		c.db.Close()
	}
}