package aggr

import "context"

type CommandHandler interface {
	// HandleCommand applies a command to the aggregate.
	HandleCommand(ctx context.Context, cmd Command) (*Event, error)
}

type CommandHandlerFunc func(context.Context, Command) (*Event, error)

func (f CommandHandlerFunc) HandleCommand(ctx context.Context, cmd Command) (*Event, error) {
	return f(ctx, cmd)
}
