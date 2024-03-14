package aggregate

import (
	"fmt"

	"github.com/dwethmar/goblin/pkg/es"
	"github.com/dwethmar/goblin/pkg/es/event"
)

var _ es.AggregateStore = &Store{}

type Store struct {
	eventStore event.Store
	factory    *Factory
}

func (s *Store) Get(aggregateType, aggregateID string) (es.Aggregate, error) {
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
		if err := aggregate.HandleEvent(event); err != nil {
			return nil, fmt.Errorf("failed to handle event on aggregate: %w", err)
		}
	}

	return aggregate, nil
}

func (s *Store) Save(a ...es.Aggregate) error {
	var events []*es.Event
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
