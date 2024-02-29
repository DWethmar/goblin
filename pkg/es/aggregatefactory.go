package es

import (
	"errors"
	"fmt"
)

var (
	ErrAggregateTypeAlreadyRegistered = errors.New("aggregate type already registered")
)

type AggregateFactoryFunc func(aggregateID string) *Aggregate

type AggregateFactory struct {
	createFuncs map[string]AggregateFactoryFunc
}

func (f *AggregateFactory) Register(aggregateType string, createFunc AggregateFactoryFunc) error {
	if _, ok := f.createFuncs[aggregateType]; ok {
		return ErrAggregateTypeAlreadyRegistered
	}
	f.createFuncs[aggregateType] = createFunc
	return nil
}

func (f *AggregateFactory) Create(aggregateType string, aggregateID string) (*Aggregate, error) {
	createFunc, ok := f.createFuncs[aggregateType]
	if !ok {
		return nil, fmt.Errorf("aggregate type %q not found", aggregateType)
	}
	return createFunc(aggregateID), nil
}

func NewAggregateFactory() *AggregateFactory {
	return &AggregateFactory{
		createFuncs: make(map[string]AggregateFactoryFunc),
	}
}
