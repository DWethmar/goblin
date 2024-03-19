package aggr

import "time"

type Event struct {
	AggregateID string
	Type        string
	Data        interface{}
	Version     int
	CreatedAt   time.Time
}
