package aggr

import "time"

type Command interface {
	// AggregateID returns the unique identifier of the aggregate.
	AggregateID() string
	// AggregateType returns the type of the aggregate.
	AggregateType() string
	// CommandTimestamp returns the time when the command was created.
	CommandTimestamp() time.Time
}

type MockCommand struct {
	ID        string
	Type      string
	Timestamp time.Time
}

func (c MockCommand) AggregateID() string         { return c.ID }
func (c MockCommand) AggregateType() string       { return c.Type }
func (c MockCommand) CommandTimestamp() time.Time { return c.Timestamp }
