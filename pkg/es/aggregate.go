package es

import (
	"errors"
	"fmt"
)

type Aggregate struct {
	ID             string
	Type           string
	Version        int
	Events         []*Event
	Model          interface{}
	Created        bool
	CommandHandler CommandHandler
	EventHandler   EventHandler
}

func (a *Aggregate) HandleCommand(command Command) (*Event, error) {
	event, err := a.CommandHandler.HandleCommand(command)
	if err != nil {
		return nil, fmt.Errorf("failed to dispatch command: %w", err)
	}

	return event, nil
}

func (a *Aggregate) HandleEvent(event *Event) error {
	if a.EventHandler == nil {
		return errors.New("event handler is nil")
	} else {
		if err := a.EventHandler.HandleEvent(event); err != nil {
			return fmt.Errorf("failed to handle event: %w", err)
		}

		a.Events = append(a.Events, event)
		a.Version++
	}

	return nil
}
