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
