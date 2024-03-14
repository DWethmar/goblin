package es

type AggregateStore interface {
	Get(aggregateType, aggregateID string) (Aggregate, error)
	Save(a ...Aggregate) error
}
