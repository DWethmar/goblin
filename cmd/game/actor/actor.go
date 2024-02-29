package actor

import (
	"errors"

	"github.com/dwethmar/tards/pkg/es"
)

const AggregateType = "actor"

type Model struct {
	Name string
}

func HandleCommand(cmd es.Command, a *es.Aggregate) (*es.Event, error) {
	_, ok := a.Model.(*Model)
	if !ok {
		return nil, errors.New("invalid model")
	}

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

func HandleEvent(event *es.Event, a *es.Aggregate) error {
	if event == nil {
		return errors.New("event is nil")
	}

	if a == nil {
		return errors.New("aggregate is nil")
	}

	m, ok := a.Model.(*Model)
	if !ok {
		return errors.New("invalid model")
	}

	switch event.Type {
	case CreatedEventType:
		data, ok := event.Data.(CreatedEventData)
		if !ok {
			return errors.New("invalid event data")
		}

		a.Created = true
		m.Name = data.Name
	}

	return nil
}
