package event

import "github.com/dwethmar/goblin/pkg/es"

// EventStore is the interface that wraps the basic event store methods.
type Store interface {
	Add(events []*es.Event) error
	List(aggregateID string) ([]*es.Event, error)
	All(err chan<- error) <-chan *es.Event
}
