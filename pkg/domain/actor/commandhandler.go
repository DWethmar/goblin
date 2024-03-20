package actor

import (
	"fmt"
	"time"

	"github.com/dwethmar/goblin/pkg/aggr"
)

func CreateCommandHandler(a *Actor, cmd *CreateCommand) (*aggr.Event, error) {
	if cmd.ActorID == "" {
		return nil, fmt.Errorf("actor id can't be empty")
	}

	if cmd.Name == "" {
		return nil, fmt.Errorf("name can't be empty")
	}

	return &aggr.Event{
		AggregateID: cmd.ActorID,
		Type:        CreatedEventType,
		Data: &CreatedEventData{
			Name: cmd.Name,
			X:    cmd.X,
			Y:    cmd.Y,
		},
		Version:   a.Version + 1,
		CreatedAt: time.Now(),
	}, nil
}

func DestroyCommandHandler(a *Actor, cmd *DestroyCommand) (*aggr.Event, error) {
	return &aggr.Event{
		AggregateID: cmd.ActorID,
		Type:        DestroyedEventType,
		Data:        nil,
		Version:     a.Version + 1,
		CreatedAt:   time.Now(),
	}, nil
}

func MoveCommandHandler(a *Actor, cmd *MoveCommand) (*aggr.Event, error) {
	return &aggr.Event{
		AggregateID: cmd.ActorID,
		Type:        MovedEventType,
		Data: &MovedEventData{
			X: cmd.X,
			Y: cmd.Y,
		},
		Version:   a.Version + 1,
		CreatedAt: time.Now(),
	}, nil
}
