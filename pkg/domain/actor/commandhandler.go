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

	return &aggr.Event{
		AggregateID: cmd.ActorID,
		EventType:   CreatedEventType,
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

	return &aggr.Event{
		AggregateID: cmd.ActorID,
		EventType:   DestroyedEventType,
		Data:        nil,
		Version:     a.Version + 1,
		Timestamp:   cmd.CommandTimestamp(),
	}, nil
}

func MoveCommandHandler(a *Actor, cmd *MoveCommand) (*aggr.Event, error) {
	return &aggr.Event{
		AggregateID: cmd.ActorID,
		EventType:   MovedEventType,
		Data: &MovedEventData{
			X: cmd.X,
			Y: cmd.Y,
		},
		Version:   a.Version + 1,
		Timestamp: cmd.CommandTimestamp(),
	}, nil
}
