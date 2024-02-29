package actor

import "github.com/dwethmar/tards/pkg/es"

func RegisterFactory(f *es.AggregateFactory) error {
	f.Register(AggregateType, func(aggregateID string) *es.Aggregate {
		return &es.Aggregate{}
	})
	return nil
}
