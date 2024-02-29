package es

type Event struct {
	AggregateID string
	Type        string
	Data        interface{}
}
