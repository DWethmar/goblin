// Package es provides a way to store and retrieve aggregates
// with eventsourcing.
package es

import (
	"context"
	"fmt"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/aggr/event"
)

var _ aggr.AggregateStore = &AggregateStore{}

type AggregateStore struct {
	eventRepository event.Repository
	factory         *aggr.Factory
}

func (s *AggregateStore) Get(ctx context.Context, aggregateType, aggregateID string) (*aggr.Aggregate, error) {
	events, err := s.eventRepository.List(aggregateID)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	aggregate, err := s.factory.Create(aggregateType, aggregateID)
	if err != nil {
		return nil, fmt.Errorf("failed to create aggregate: %w", err)
	}

	if len(events) == 0 {
		return aggregate, nil
	}

	for _, event := range events {
		if err := aggregate.HandleEvent(ctx, event); err != nil {
			return nil, fmt.Errorf("failed to handle event on aggregate: %w", err)
		}
	}

	return aggregate, nil
}

func (s *AggregateStore) Save(ctx context.Context, a ...*aggr.Aggregate) error {
	var events []*aggr.Event
	for _, aggregate := range a {
		events = append(events, aggregate.AggregateEvents()...)
	}

	if len(events) == 0 {
		return nil
	}

	if err := s.eventRepository.Add(events); err != nil {
		return fmt.Errorf("failed to save events: %w", err)
	}

	return nil
}

func NewAggregateStore(eventRepository event.Repository, factory *aggr.Factory) *AggregateStore {
	return &AggregateStore{
		eventRepository: eventRepository,
		factory:         factory,
	}
}
