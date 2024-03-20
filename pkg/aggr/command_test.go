package aggr

import "time"

type MockCommand struct {
	aggregateID   string
	aggregateType string
	timestamp     time.Time
}

func (c MockCommand) AggregateID() string         { return c.aggregateID }
func (c MockCommand) AggregateType() string       { return c.aggregateType }
func (c MockCommand) CommandTimestamp() time.Time { return c.timestamp }
