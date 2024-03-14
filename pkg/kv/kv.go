package kv

type Entry struct {
	Key   []byte
	Value []byte
}

type DB interface {
	Get([]byte) ([]byte, error)
	Put(...Entry) error
	Iterate(func(k, v []byte) error) error
	IterateWithPrefix(prefix []byte, f func(k, v []byte) error) error
	Delete([]byte) error
}
