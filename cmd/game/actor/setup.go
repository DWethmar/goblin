package actor

import "github.com/dwethmar/tards/pkg/es"

func RegisterFactory(f *es.AggregateFactory) error {
	f.Register(AggregateType, func(aggregateID string) *es.Aggregate {
		actor := &Actor{}
		return &es.Aggregate{
			ID:             aggregateID,
			Type:           AggregateType,
			Model:          actor,
			Events:         []*es.Event{},
			Created:        false,
			CommandHandler: &CommandHandler{Actor: actor},
			EventHandler:   &EventHandler{Actor: actor},
		}
	})
	return nil
}
