package es

import "fmt"

type Aggregate struct {
	ID             string
	Type           string
	Version        int
	Events         []*Event
	Model          interface{}
	Created        bool
	commandHandler CommandHandler
	eventHandler   EventHandler
}

func (a *Aggregate) HandleCommand(command Command) (*Event, error) {
	event, err := a.commandHandler.Dispatch(command)
	if err != nil {
		return nil, fmt.Errorf("failed to dispatch command: %w", err)
	}

	return event, nil
}

func (a *Aggregate) HandleEvent(event *Event) error {
	if a.eventHandler != nil {
		if err := a.eventHandler.HandleEvent(event); err != nil {
			return err
		}

		a.Events = append(a.Events, event)
		a.Version++
	}

	return nil
}
