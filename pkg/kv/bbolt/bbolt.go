package bbolt

import (
	"bytes"
	"fmt"

	"github.com/dwethmar/goblin/pkg/kv"
	bolt "go.etcd.io/bbolt"
)

var _ kv.DB = &DB{}

type DB struct {
	bucket []byte
	db     *bolt.DB
}

func (d *DB) Get(k []byte) ([]byte, error) {
	var v []byte
	return v, d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(d.bucket)
		if b == nil {
			return nil
		}
		v = b.Get(k)
		return nil
	})
}

func (d *DB) Put(i ...kv.Entry) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(d.bucket)
		if err != nil {
			return err
		}
		for _, e := range i {
			if err := b.Put(e.Key, e.Value); err != nil {
				return fmt.Errorf("putting key: %w", err)
			}
		}

		return nil
	})
}

func (d *DB) Delete(k []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(d.bucket)
		if b == nil {
			return nil
		}
		return b.Delete(k)
	})
}

func (d *DB) Iterate(f func(k []byte, v []byte) error) error {
	return d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(d.bucket)
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := f(k, v); err != nil {
				return fmt.Errorf("iterating error: %w", err)
			}
		}

		return nil
	})
}

// Iterate iterates over all the key-value pairs in the database.
func (d *DB) IterateWithPrefix(prefix []byte, f func(k []byte, v []byte) error) error {
	return d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(d.bucket)
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			if err := f(k, v); err != nil {
				return err
			}
		}

		return nil
	})
}

func New(bucket []byte, db *bolt.DB) *DB {
	return &DB{
		bucket: bucket,
		db:     db,
	}
}
