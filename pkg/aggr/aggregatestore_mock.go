package aggr

import "context"

type MockAggregateStore struct {
	GetFunc  func(ctx context.Context, aggregateType, aggregateID string) (*Aggregate, error)
	SaveFunc func(context.Context, ...*Aggregate) error
}

func (m MockAggregateStore) Get(ctx context.Context, aggregateType, aggregateID string) (*Aggregate, error) {
	return m.GetFunc(ctx, aggregateType, aggregateID)
}

func (m MockAggregateStore) Save(ctx context.Context, p0 ...*Aggregate) error {
	return m.SaveFunc(ctx, p0...)
}
