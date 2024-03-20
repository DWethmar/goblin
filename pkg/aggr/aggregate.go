package aggr

type Model interface {
	CommandHandler
	EventHandler
	AggregateID() string
	AggregateEvents() []*Event
	ClearAggregateEvents()
	AggregateVersion() uint
}

// Aggregate is the interface that wraps the basic methods for an aggregate.
type Aggregate struct {
	Model
}
