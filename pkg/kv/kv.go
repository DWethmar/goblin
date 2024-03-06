package kv

type DB interface {
	Get([]byte) ([]byte, error)
	Put([]byte, []byte) error
	IterateWithPrefix(prefix []byte, f func(k, v []byte) error) error
	Delete([]byte) error
}
