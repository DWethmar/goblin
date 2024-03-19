package aggr

type CommandHandler interface {
	// HandleCommand applies a command to the aggregate.
	HandleCommand(cmd Command) (*Event, error)
}

type CommandHandlerFunc func(cmd Command) (*Event, error)

func (f CommandHandlerFunc) HandleCommand(cmd Command) (*Event, error) {
	return f(cmd)
}
