package aggr

import (
	"errors"
	"fmt"
)

var (
	ErrorAggregateTypeNotFound = errors.New("aggregate type not found")
)

type Factory struct {
	createFuncs map[string]func(aggregateID string) *Aggregate
}

func (f *Factory) Register(aggregateType string, ff func(aggregateID string) *Aggregate) {
	f.createFuncs[aggregateType] = ff
}

func (f *Factory) Create(aggregateType string, aggregateID string) (*Aggregate, error) {
	createFunc, ok := f.createFuncs[aggregateType]
	if !ok {
		return nil, fmt.Errorf("aggregate type %q not found: %w", aggregateType, ErrorAggregateTypeNotFound)
	}
	return createFunc(aggregateID), nil
}

func NewFactory(opt ...func(f *Factory)) *Factory {
	f := &Factory{
		createFuncs: make(map[string]func(aggregateID string) *Aggregate),
	}

	for _, o := range opt {
		o(f)
	}

	return f
}
