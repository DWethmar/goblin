package bbolt

import (
	"bytes"

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
		b, err := tx.CreateBucketIfNotExists(d.bucket)
		if err != nil {
			return err
		}

		v = b.Get(k)
		return nil
	})
}

func (d *DB) Put(k []byte, v []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(d.bucket)
		if err != nil {
			return err
		}

		return b.Put(k, v)
	})
}

func (d *DB) Delete(k []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(d.bucket)
		if err != nil {
			return err
		}

		if b == nil {
			return nil
		}
		return b.Delete(k)
	})
}

// Iterate iterates over all the key-value pairs in the database.
func (d *DB) IterateWithPrefix(prefix []byte, f func(k []byte, v []byte) error) error {
	return d.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(d.bucket).Cursor()
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
