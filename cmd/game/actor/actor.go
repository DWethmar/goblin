package actor

import (
	"fmt"

	"github.com/dwethmar/goblin/pkg/es"
)

const AggregateType = "actor"

type Actor struct {
	Name    string
	created bool
}

func (a *Actor) HandleCommand(cmd es.Command) (*es.Event, error) {
	switch c := cmd.(type) {
	case *CreateCommand:
		if a.created {
			return nil, fmt.Errorf("actor already created")
		}

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

func (a *Actor) HandleEvent(event *es.Event) error {
	switch event.Type {
	case CreatedEventType:
		data, ok := event.Data.(CreatedEventData)
		if !ok {
			return fmt.Errorf("invalid event data type: %T", event.Data)
		}

		a.Name = data.Name
		a.created = true
	}

	return nil
}
