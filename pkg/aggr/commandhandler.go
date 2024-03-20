package aggr

import "context"

type CommandHandler interface {
	// HandleCommand applies a command to the aggregate.
	HandleCommand(ctx context.Context, cmd Command) (*Event, error)
}

type CommandHandlerFunc func(cmd Command) (*Event, error)

func (f CommandHandlerFunc) HandleCommand(cmd Command) (*Event, error) {
	return f(cmd)
}
