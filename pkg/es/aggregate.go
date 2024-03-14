package es

// Aggregate is the interface that wraps the basic methods for an aggregate.
type Aggregate interface {
	CommandHandler
	EventHandler
	AggregateID() string
	AggregateEvents() []*Event
	ClearAggregateEvents()
	AggregateVersion() int
}
