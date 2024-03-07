package es

type MockAggregateStore struct {
	GetFunc  func(aggregateType, aggregateID string) (*Aggregate, error)
	SaveFunc func(*Aggregate) error
}

func (m MockAggregateStore) Get(aggregateType, aggregateID string) (*Aggregate, error) {
	return m.GetFunc(aggregateType, aggregateID)
}

func (m MockAggregateStore) Save(p0 *Aggregate) error {
	return m.SaveFunc(p0)
}
