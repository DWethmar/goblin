package actor

import (
	"fmt"
	"time"

	"github.com/dwethmar/goblin/pkg/es"
)

const AggregateType = "actor"

type State int

func (s State) Is(v State) bool { return s == v }

const (
	StateDraft State = iota
	StateCreated
	StateDeleted
)

type Actor struct {
	ID    string
	Name  string
	X, Y  int
	state State
}

func (a *Actor) AggregateID() string { return a.ID }

func (a *Actor) HandleCommand(cmd es.Command) (*es.Event, error) {
	if StateDeleted.Is(a.state) {
		return nil, fmt.Errorf("actor deleted")
	}

	if _, ok := cmd.(*CreateCommand); !ok && StateCreated.Is(a.state) {
		return nil, fmt.Errorf("actor is already created")
	}

	switch c := cmd.(type) {
	case *CreateCommand:
		return &es.Event{
			AggregateID: c.ActorID,
			Type:        CreatedEventType,
			Data: &CreatedEventData{
				Name: c.Name,
				X:    c.X,
				Y:    c.Y,
			},
			Created: time.Now(),
		}, nil
	case *MoveCommand:
		return &es.Event{
			AggregateID: c.ActorID,
			Type:        MovedEventType,
			Data: &MovedEventData{
				X: c.X,
				Y: c.Y,
			},
			Created: time.Now(),
		}, nil
	}

	return nil, nil
}

func (a *Actor) HandleEvent(event *es.Event) error {
	switch event.Type {
	case CreatedEventType:
		data, ok := event.Data.(*CreatedEventData)
		if !ok {
			return fmt.Errorf("expected *CreatedEventData, got %T", event.Data)
		}

		a.ID = event.AggregateID
		a.Name = data.Name
		a.state = StateCreated

	case MovedEventType:
		data, ok := event.Data.(*MovedEventData)
		if !ok {
			return fmt.Errorf("expected *MovedEventData, got %T", event.Data)
		}

		a.X = data.X
		a.Y = data.Y
	}

	return nil
}

func (a *Actor) Deleted() bool { return StateDeleted.Is(a.state) }
