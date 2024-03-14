package es

type Aggregate interface {
	CommandHandler
	EventHandler
	AggregateID() string
	AggregateEvents() []*Event
	ClearAggregateEvents()
	AggregateVersion() int
}
