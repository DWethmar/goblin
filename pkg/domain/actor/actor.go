package actor

import (
	"context"
	"fmt"

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

	state  State
	events []*aggr.Event
}

func (a *Actor) AggregateID() string   { return a.ID }
func (a *Actor) AggregateVersion() int { return a.Version }

func (a *Actor) HandleCommand(cmd aggr.Command) (*aggr.Event, error) {
	if StateDeleted.Is(a.state) {
		return nil, fmt.Errorf("actor deleted")
	}

	// if state is draft and command is not create, return error
	if StateDraft.Is(a.state) {
		if _, ok := cmd.(*CreateCommand); !ok {
			return nil, fmt.Errorf("actor does not exist")
		}
	}

	// if state is created and command is create, return error
	if StateCreated.Is(a.state) {
		if _, ok := cmd.(*CreateCommand); ok {
			return nil, fmt.Errorf("actor already created")
		}
	}

	switch c := cmd.(type) {
	case *CreateCommand:
		return CreateCommandHandler(a, c)
	case *DestroyCommand:
		return DestroyCommandHandler(a, c)
	case *MoveCommand:
		return MoveCommandHandler(a, c)
	}

	return nil, nil
}

func (a *Actor) HandleEvent(_ context.Context, event *aggr.Event) error {
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
