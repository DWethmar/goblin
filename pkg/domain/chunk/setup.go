package chunk

import (
	"encoding/gob"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/aggr/aggrstore"
)

func init() {
	gob.Register(&CreatedEventData{})
}

func RegisterFactory(f *aggrstore.Factory) {
	f.Register(AggregateType, func(aggregateID string) *aggr.Aggregate {
		return &aggr.Aggregate{
			Model: New(aggregateID, 0, 0),
		}
	})
}
