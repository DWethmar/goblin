package es

// EventStore is the interface that wraps the basic event store methods.
type EventStore interface {
	Add(events []*Event) error
	List(aggregateID string) ([]*Event, error)
}
