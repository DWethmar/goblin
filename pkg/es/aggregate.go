package es

type Model interface {
	CommandHandler
	EventHandler
	AggregateID() string
	AggregateEvents() []*Event
	ClearAggregateEvents()
	AggregateVersion() int
}

// Aggregate is the interface that wraps the basic methods for an aggregate.
type Aggregate struct {
	Model
}
