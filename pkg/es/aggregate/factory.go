package aggregate

import (
	"errors"
	"fmt"

	"github.com/dwethmar/goblin/pkg/es"
)

var (
	ErrAggregateTypeAlreadyRegistered = errors.New("aggregate type already registered")
)

type FactoryFunc func(aggregateID string) *es.Aggregate

type Factory struct {
	createFuncs map[string]FactoryFunc
}

func (f *Factory) Register(aggregateType string, createFunc FactoryFunc) error {
	if _, ok := f.createFuncs[aggregateType]; ok {
		return ErrAggregateTypeAlreadyRegistered
	}
	f.createFuncs[aggregateType] = createFunc
	return nil
}

func (f *Factory) Create(aggregateType string, aggregateID string) (*es.Aggregate, error) {
	createFunc, ok := f.createFuncs[aggregateType]
	if !ok {
		return nil, fmt.Errorf("aggregate type %q not found", aggregateType)
	}
	return createFunc(aggregateID), nil
}

func NewFactory() *Factory {
	return &Factory{
		createFuncs: make(map[string]FactoryFunc),
	}
}
