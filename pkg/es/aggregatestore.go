package es

import "fmt"

type AggregateStore struct {
	eventRepo EventStore
	factory   *AggregateFactory
}

func (s *AggregateStore) Get(aggregateType, aggregateID string) (*Aggregate, error) {
	events, err := s.eventRepo.List(aggregateID)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("failed to handle event: %w", err)
		}
	}

	return aggregate, nil
}

func (s *AggregateStore) Save(a *Aggregate) error {
	if len(a.Events) == 0 {
		return nil
	}

	if err := s.eventRepo.Add(a.Events); err != nil {
		return fmt.Errorf("failed to save events: %w", err)
	}

	return nil
}

func NewAggregateStore(eventRepo EventStore, factory *AggregateFactory) *AggregateStore {
	return &AggregateStore{
		eventRepo: eventRepo,
		factory:   factory,
	}
}
