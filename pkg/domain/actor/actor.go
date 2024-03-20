package actor

import (
	"context"
	"errors"
	"fmt"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/domain"
)

var _ aggr.Model = &Actor{}

const AggregateType = "actor"

var (
	ErrNilCommand          = errors.New("command is nil")
	ErrUnknownCommandType  = errors.New("unknown command type")
	ErrActorIsDeleted      = errors.New("actor is deleted")
	ErrActorDoesNotExist   = errors.New("actor does not exist")
	ErrActorAlreadyCreated = errors.New("actor already created")
)

type Actor struct {
	ID      string
	Version int
	Name    string
	X, Y    int

	state  domain.State
	events []*aggr.Event
}

func New(id, name string, x, y int) *Actor {
	return &Actor{
		ID:     id,
		Name:   name,
		state:  domain.StateDraft,
		X:      x,
		Y:      y,
		events: []*aggr.Event{},
	}
}

func (a *Actor) AggregateID() string   { return a.ID }
func (a *Actor) AggregateVersion() int { return a.Version }

func (a *Actor) HandleCommand(cmd aggr.Command) (*aggr.Event, error) {
	if cmd == nil {
		return nil, ErrNilCommand
	}

	if domain.StateDeleted.Is(a.state) {
		return nil, ErrActorIsDeleted
	}

	// if state is draft and command is not create, return error
	if domain.StateDraft.Is(a.state) {
		if _, ok := cmd.(*CreateCommand); !ok {
			return nil, ErrActorDoesNotExist
		}
	}

	// if state is created and command is create, return error
	if domain.StateCreated.Is(a.state) {
		if _, ok := cmd.(*CreateCommand); ok {
			return nil, ErrActorAlreadyCreated
		}
	}

	switch c := cmd.(type) {
	case *CreateCommand:
		return CreateCommandHandler(a, c)
	case *DestroyCommand:
		return DestroyCommandHandler(a, c)
	case *MoveCommand:
		return MoveCommandHandler(a, c)
	default:
		return nil, fmt.Errorf("unknown command type: %T: %w", cmd, ErrUnknownCommandType)
	}
}

func (a *Actor) HandleEvent(_ context.Context, event *aggr.Event) error {
	var err error
	switch event.Type {
	case CreatedEventType:
		err = HandleCreatedEvent(a, event)
	case DestroyedEventType:
		err = HandleDestroyedEvent(a, event)
	case MovedEventType:
		err = HandleMovedEvent(a, event)
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}

	if err != nil {
		return err
	}

	a.Version = event.Version
	a.events = append(a.events, event)
	return nil
}

func (a *Actor) AggregateEvents() []*aggr.Event { return a.events }
func (a *Actor) ClearAggregateEvents()          { a.events = []*aggr.Event{} }
func (a *Actor) Deleted() bool                  { return domain.StateDeleted.Is(a.state) }
