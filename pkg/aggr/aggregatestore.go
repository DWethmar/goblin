package aggr

import (
	"context"
)

type AggregateStore interface {
	// Get returns the aggregate by its type and ID. An Aggregate is always returned, even if it does not exist yet.
	Get(ctx context.Context, aggregateType, aggregateID string) (*Aggregate, error)
	Save(ctx context.Context, a ...*Aggregate) error
}
