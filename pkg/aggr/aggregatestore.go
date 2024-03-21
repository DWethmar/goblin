package aggr

import (
	"context"
)

type AggregateStore interface {
	// Get returns the aggregate by its type and ID. An Aggregate is always returned, even if it does not exist yet.
	Get(ctx context.Context, aggregateType, aggregateID string) (*Aggregate, error)
	Save(ctx context.Context, a ...*Aggregate) error
}

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

// NoopRepository is a mock repository that returns an empty aggregate.
var NoopRepository = &MockAggregateStore{
	GetFunc: func(ctx context.Context, aggregateType, id string) (*Aggregate, error) {
		return &Aggregate{
			Model: &MockModel{
				ID:   id,
				Type: aggregateType,
			},
		}, nil
	},
	SaveFunc: func(_ context.Context, _ ...*Aggregate) error { return nil },
}
