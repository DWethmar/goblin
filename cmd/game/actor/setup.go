package actor

import (
	"github.com/dwethmar/goblin/pkg/es"
	"github.com/dwethmar/goblin/pkg/es/aggregate"
)

func RegisterFactory(f *aggregate.Factory) error {
	f.Register(AggregateType, func(aggregateID string) *es.Aggregate {
		actor := &Actor{}
		return &es.Aggregate{
			ID:     aggregateID,
			Type:   AggregateType,
			Model:  actor,
			Events: []*es.Event{},
		}
	})
	return nil
}
