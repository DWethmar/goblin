package actor

import (
	"fmt"

	"github.com/dwethmar/goblin/pkg/aggr"
)

func CreateCommandHandler(a *Actor, cmd *CreateCommand) (*aggr.Event, error) {
	if cmd.ActorID == "" {
		return nil, fmt.Errorf("actor id can't be empty")
	}

	if cmd.Name == "" {
		return nil, fmt.Errorf("name can't be empty")
	}

	if cmd.CommandTimestamp().IsZero() {
		return nil, fmt.Errorf("destroyed at can't be zero")
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
		Timestamp: cmd.CommandTimestamp(),
	}, nil
}

func DestroyCommandHandler(a *Actor, cmd *DestroyCommand) (*aggr.Event, error) {
	if cmd.ActorID == "" {
		return nil, fmt.Errorf("actor id can't be empty")
	}

	if cmd.CommandTimestamp().IsZero() {
		return nil, fmt.Errorf("destroyed at can't be zero")
	}

	return &aggr.Event{
		AggregateID: cmd.ActorID,
		Type:        DestroyedEventType,
		Data:        nil,
		Version:     a.Version + 1,
		Timestamp:   cmd.CommandTimestamp(),
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
		Timestamp: cmd.CommandTimestamp(),
	}, nil
}
