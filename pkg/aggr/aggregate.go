package aggr

import (
	"time"
)

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

var _ Model = &MockModel{}

// MockModel is a mock implementation of the Model interface.
type MockModel struct {
	EventHandler
	CommandHandler
	ID        string
	Type      string
	Version   uint
	Events    []*Event
	Timestamp time.Time
}

func (a *MockModel) AggregateID() string         { return a.ID }
func (a *MockModel) AggregateType() string       { return a.Type }
func (a *MockModel) AggregateVersion() uint      { return a.Version }
func (a *MockModel) AggregateEvents() []*Event   { return a.Events }
func (a *MockModel) ClearAggregateEvents()       { a.Events = []*Event{} }
func (a *MockModel) ApplyEvent(e *Event)         { a.Events = append(a.Events, e) }
func (a *MockModel) CommandTimestamp() time.Time { return a.Timestamp }
