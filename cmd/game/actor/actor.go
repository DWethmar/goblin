package actor

import (
	"fmt"

	"github.com/dwethmar/tards/pkg/es"
)

const AggregateType = "actor"

type Actor struct {
	Name string
}

type CommandHandler struct {
	Actor *Actor
}

func (h *CommandHandler) HandleCommand(cmd es.Command) (*es.Event, error) {
	switch c := cmd.(type) {
	case *CreateCommand:
		return &es.Event{
			AggregateID: c.ActorID,
			Type:        CreatedEventType,
			Data: CreatedEventData{
				Name: c.Name,
			},
		}, nil
	}

	return nil, nil
}

type EventHandler struct {
	Actor *Actor
}

func (h *EventHandler) HandleEvent(event *es.Event) error {
	switch event.Type {
	case CreatedEventType:
		data, ok := event.Data.(CreatedEventData)
		if !ok {
			return fmt.Errorf("invalid event data type: %T", event.Data)
		}

		h.Actor.Name = data.Name
	}

	return nil
}
