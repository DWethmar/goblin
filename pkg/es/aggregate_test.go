package es

import "errors"

var _ Aggregate = &MockAggregate{}

type MockAggregate struct {
	ID                 string
	Version            int
	Events             []*Event
	CommandHandlerFunc func(Command) (*Event, error)
	EventHandlerFunc   func(*Event) error
}

func (a *MockAggregate) AggregateID() string {
	return a.ID
}

func (a *MockAggregate) AggregateVersion() int {
	return a.Version
}

func (a *MockAggregate) AggregateEvents() []*Event {
	return a.Events
}

func (a *MockAggregate) ClearAggregateEvents() {
	a.Events = []*Event{}
}

func (a *MockAggregate) HandleCommand(c Command) (*Event, error) {
	if a.CommandHandlerFunc == nil {
		return nil, errors.New("CommandHandlerFunc is not set")
	}
	return a.CommandHandlerFunc(c)
}

func (a *MockAggregate) HandleEvent(e *Event) error {
	if a.EventHandlerFunc == nil {
		return errors.New("EventHandlerFunc is not set")
	}
	return a.EventHandlerFunc(e)
}
