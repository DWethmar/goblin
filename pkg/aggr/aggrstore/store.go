package aggrstore

import (
	"context"
	"fmt"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/aggr/event"
)

var _ aggr.AggregateStore = &Store{}

type Store struct {
	eventStore event.Store
	factory    *Factory
}

func (s *Store) Get(ctx context.Context, aggregateType, aggregateID string) (*aggr.Aggregate, error) {
	events, err := s.eventStore.List(aggregateID)
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

func (s *Store) Save(ctx context.Context, a ...*aggr.Aggregate) error {
	var events []*aggr.Event
	for _, aggregate := range a {
		events = append(events, aggregate.AggregateEvents()...)
	}

	if len(events) == 0 {
		return nil
	}

	if err := s.eventStore.Add(events); err != nil {
		return fmt.Errorf("failed to save events: %w", err)
	}

	return nil
}

func NewStore(eventStore event.Store, factory *Factory) *Store {
	return &Store{
		eventStore: eventStore,
		factory:    factory,
	}
}
