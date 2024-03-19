package bbolt

import (
	"fmt"
	"os"
	"path"

	bolt "go.etcd.io/bbolt"
)

// Connect to a bbolt database
func Connect(filepath string) (*bolt.DB, error) {
	dirName := path.Dir(filepath)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			return nil, fmt.Errorf("creating dir: %w", merr)
		}
	}

	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	return db, nil
}
