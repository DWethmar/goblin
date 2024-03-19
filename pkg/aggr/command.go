package aggr

type Command interface {
	// CommandType returns the type of the command.
	CommandType() string
	// AggregateID returns the unique identifier of the aggregate.
	AggregateID() string
	// AggregateType returns the type of the aggregate.
	AggregateType() string
}
