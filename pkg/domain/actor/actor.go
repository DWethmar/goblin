package actor

import (
	"fmt"
	"time"

	"github.com/dwethmar/goblin/pkg/aggr"
)

var _ aggr.Model = &Actor{}

const AggregateType = "actor"

type State int

func (s State) Is(v State) bool { return s == v }

const (
	StateDraft State = iota
	StateCreated
	StateDeleted
)

type Actor struct {
	ID      string
	Version int
	Name    string
	X, Y    int
	state   State

	events []*aggr.Event
}

func (a *Actor) AggregateID() string   { return a.ID }
func (a *Actor) AggregateVersion() int { return a.Version }

func (a *Actor) HandleCommand(cmd aggr.Command) (*aggr.Event, error) {
	if StateDeleted.Is(a.state) {
		return nil, fmt.Errorf("actor deleted")
	}

	if StateCreated.Is(a.state) {
		switch cmd.(type) {
		case *CreateCommand:
			return nil, fmt.Errorf("actor already created")
		}
	}

	switch c := cmd.(type) {
	case *CreateCommand:
		if c.Name == "" {
			return nil, fmt.Errorf("name can't be empty")
		}

		return &aggr.Event{
			AggregateID: c.ActorID,
			Type:        CreatedEventType,
			Data: &CreatedEventData{
				Name: c.Name,
				X:    c.X,
				Y:    c.Y,
			},
			Version:   a.Version + 1,
			CreatedAt: time.Now(),
		}, nil
	case *MoveCommand:
		return &aggr.Event{
			AggregateID: c.ActorID,
			Type:        MovedEventType,
			Data: &MovedEventData{
				X: c.X,
				Y: c.Y,
			},
			Version:   a.Version + 1,
			CreatedAt: time.Now(),
		}, nil
	}

	return nil, nil
}

func (a *Actor) HandleEvent(event *aggr.Event) error {
	switch event.Type {
	case CreatedEventType:
		data, ok := event.Data.(*CreatedEventData)
		if !ok {
			return fmt.Errorf("expected *CreatedEventData, got %T", event.Data)
		}

		a.ID = event.AggregateID
		a.Name = data.Name
		a.state = StateCreated
	case DestroyedEventType:
		_, ok := event.Data.(*DestroyedEventData)
		if !ok {
			return fmt.Errorf("expected *DestroyedEventData, got %T", event.Data)
		}
		a.state = StateDeleted
	case MovedEventType:
		data, ok := event.Data.(*MovedEventData)
		if !ok {
			return fmt.Errorf("expected *MovedEventData, got %T", event.Data)
		}
		a.X = data.X
		a.Y = data.Y
	}

	a.Version = event.Version
	a.events = append(a.events, event)
	return nil
}

func (a *Actor) AggregateEvents() []*aggr.Event { return a.events }
func (a *Actor) ClearAggregateEvents()          { a.events = []*aggr.Event{} }
func (a *Actor) Deleted() bool                  { return StateDeleted.Is(a.state) }
