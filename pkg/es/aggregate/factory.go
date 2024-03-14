package aggregate

import (
	"errors"
	"fmt"

	"github.com/dwethmar/goblin/pkg/es"
)

var (
	ErrAggregateTypeAlreadyRegistered = errors.New("aggregate type already registered")
)

type Factory struct {
	createFuncs map[string]func(aggregateID string) es.Aggregate
}

func (f *Factory) Register(aggregateType string, ff func(aggregateID string) es.Aggregate) error {
	if _, ok := f.createFuncs[aggregateType]; ok {
		return ErrAggregateTypeAlreadyRegistered
	}
	f.createFuncs[aggregateType] = ff
	return nil
}

func (f *Factory) Create(aggregateType string, aggregateID string) (es.Aggregate, error) {
	createFunc, ok := f.createFuncs[aggregateType]
	if !ok {
		return nil, fmt.Errorf("aggregate type %q not found", aggregateType)
	}
	return createFunc(aggregateID), nil
}

func NewFactory() *Factory {
	return &Factory{
		createFuncs: make(map[string]func(aggregateID string) es.Aggregate),
	}
}
