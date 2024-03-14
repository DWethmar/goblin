package actor

import (
	"encoding/gob"

	"github.com/dwethmar/goblin/pkg/es"
	"github.com/dwethmar/goblin/pkg/es/aggregate"
)

func init() {
	gob.Register(&CreatedEventData{})
}

func RegisterFactory(f *aggregate.Factory) error {
	f.Register(AggregateType, func(aggregateID string) es.Aggregate {
		return &Actor{
			ID: aggregateID,
		}
	})
	return nil
}
