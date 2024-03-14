package es

import (
	"errors"
	"fmt"
)

type AggregateModel interface {
	CommandHandler
	EventHandler
	AggregateID() string
}

type Aggregate struct {
	Type    string
	Events  []*Event
	Model   AggregateModel
	Version int
	Created bool
}

func (a *Aggregate) AggregateID() string { return a.Model.AggregateID() }

// HandleCommand handles a command and returns an event
func (a *Aggregate) HandleCommand(command Command) (*Event, error) {
	if a.Model == nil {
		return nil, errors.New("model is nil")
	}

	event, err := a.Model.HandleCommand(command)
	if err != nil {
		return nil, fmt.Errorf("failed to handle command: %w", err)
	}

	if event == nil {
		return nil, fmt.Errorf("event is nil")
	}

	event.Version = a.Version + 1
	return event, nil
}

// HandleEvent handles an event and updates the aggregate
func (a *Aggregate) HandleEvent(event *Event) error {
	if event.Version != a.Version+1 {
		return fmt.Errorf("invalid version: %d, expected: %d", event.Version, a.Version+1)
	}

	if a.Model == nil {
		return errors.New("model is nil")
	}

	if err := a.Model.HandleEvent(event); err != nil {
		return fmt.Errorf("failed to handle event on model: %w", err)
	}

	a.Version = event.Version
	a.Events = append(a.Events, event)

	return nil
}
