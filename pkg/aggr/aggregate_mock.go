package aggr

import (
	"context"
	"errors"
	"time"
)

var _ Model = &MockAggregate{}

type MockAggregate struct {
	ID                 string
	Type               string
	Version            uint
	Events             []*Event
	Timestamp          time.Time
	CommandHandlerFunc func(context.Context, Command) (*Event, error)
	EventHandlerFunc   func(context.Context, *Event) error
}

func (a *MockAggregate) AggregateID() string         { return a.ID }
func (a *MockAggregate) AggregateType() string       { return a.Type }
func (a *MockAggregate) AggregateVersion() uint      { return a.Version }
func (a *MockAggregate) AggregateEvents() []*Event   { return a.Events }
func (a *MockAggregate) ClearAggregateEvents()       { a.Events = []*Event{} }
func (a *MockAggregate) ApplyEvent(e *Event)         { a.Events = append(a.Events, e) }
func (a *MockAggregate) CommandTimestamp() time.Time { return a.Timestamp }

func (a *MockAggregate) HandleCommand(ctx context.Context, c Command) (*Event, error) {
	if a.CommandHandlerFunc == nil {
		return nil, errors.New("CommandHandlerFunc is not set")
	}
	return a.CommandHandlerFunc(ctx, c)
}

func (a *MockAggregate) HandleEvent(ctx context.Context, e *Event) error {
	if a.EventHandlerFunc == nil {
		return errors.New("EventHandlerFunc is not set")
	}
	return a.EventHandlerFunc(ctx, e)
}
