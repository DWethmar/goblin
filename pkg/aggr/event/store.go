package event

import "github.com/dwethmar/goblin/pkg/aggr"

// EventStore is the interface that wraps the basic event store methods.
type Store interface {
	Add(events []*aggr.Event) error
	List(aggregateID string) ([]*aggr.Event, error)
	All(err chan<- error) <-chan *aggr.Event
}
