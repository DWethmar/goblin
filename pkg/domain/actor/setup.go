package actor

import (
	"encoding/gob"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/aggr/aggrstore"
)

func init() {
	gob.Register(&CreatedEventData{})
	gob.Register(&MovedEventData{})
}

func RegisterFactory(f *aggrstore.Factory) {
	f.Register(AggregateType, func(aggregateID string) *aggr.Aggregate {
		return &aggr.Aggregate{
			Model: &Actor{
				ID: aggregateID,
			},
		}
	})
}
