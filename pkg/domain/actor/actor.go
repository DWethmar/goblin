package actor

import (
	"fmt"

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
			},
		}, nil
	}

	return nil, nil
}

func (a *Actor) HandleEvent(event *es.Event) error {
	switch event.Type {
	case CreatedEventType:
		if StateCreated.Is(a.state) {
			return fmt.Errorf("actor already created")
		}

		data, ok := event.Data.(*CreatedEventData)
		if !ok {
			return fmt.Errorf("expected *CreatedEventData, got %T", event.Data)
		}

		a.ID = event.AggregateID
		a.Name = data.Name
		a.state = StateCreated
	}

	return nil
}

func (a *Actor) Deleted() bool { return StateDeleted.Is(a.state) }
