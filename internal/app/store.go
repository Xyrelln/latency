package app

import (
	"bytes"
	"fmt"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

const defaultBucket = "latency"

type store struct {
	db *bolt.DB
}

func newStore(path string) (*store, error) {
	dbPath := filepath.Join(path, "latency.db")
	log.Infof("create db on path: %s", dbPath)
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Errorf("db open failed: %v", err)
		return nil, err
	}

	return &store{db}, nil
}

func (s *store) get(key []byte) ([]byte, error) {
	var result []byte

	s.db.View(func(tx *bolt.Tx) error {
		// b := tx.Bucket([]byte(defaultBucket))
		b, err := tx.CreateBucketIfNotExists([]byte(defaultBucket))
		if err != nil {
			return fmt.Errorf("create bucket %s failed, err: %v", defaultBucket, err)
		}
		result = b.Get(key)

		return nil
	})
	return result, nil
}

func (s *store) set(key, val []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(defaultBucket))
		if err != nil {
			return fmt.Errorf("create bucket %s failed, err: %v", defaultBucket, err)
		}
		return b.Put(key, val)
	})
}

func (s *store) del(key []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		// b := tx.Bucket([]byte(defaultBucket))
		b, err := tx.CreateBucketIfNotExists([]byte(defaultBucket))
		if err != nil {
			return fmt.Errorf("create bucket %s failed, err: %v", defaultBucket, err)
		}
		return b.Delete(key)
	})
}

func (s *store) list(prefix []byte) ([][]byte, error) {
	var items [][]byte

	err := s.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(defaultBucket)).Cursor()

		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			items = append(items, v)
		}
		return nil
	})

	return items, err
}

func (s *store) close() {
	if s == nil {
		return
	}
	s.db.Close()
}
