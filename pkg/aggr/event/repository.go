package event

import "github.com/dwethmar/goblin/pkg/aggr"

// Repository is used to store events.
type Repository interface {
	Add(events []*aggr.Event) error
	List(aggregateID string) ([]*aggr.Event, error)
	All(err chan<- error) <-chan *aggr.Event
}
