package kv

type MockKV struct {
	GetFn               func([]byte) ([]byte, error)
	PutFn               func(...Entry) error
	IterateFn           func(func(k, v []byte) error) error
	IterateWithPrefixFn func(prefix []byte, f func(k, v []byte) error) error
	DeleteFn            func([]byte) error
}

func (m *MockKV) Get(k []byte) ([]byte, error) {
	return m.GetFn(k)
}

func (m *MockKV) Put(i ...Entry) error {
	return m.PutFn(i...)
}

func (m *MockKV) Iterate(f func(k, v []byte) error) error {
	return m.IterateFn(f)
}

func (m *MockKV) IterateWithPrefix(prefix []byte, f func(k, v []byte) error) error {
	return m.IterateWithPrefixFn(prefix, f)
}

func (m *MockKV) Delete(k []byte) error {
	return m.DeleteFn(k)
}
