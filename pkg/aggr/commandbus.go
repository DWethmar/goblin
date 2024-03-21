package aggr

import (
	"context"
)

// CommandBus applies commands to aggregates and publishes events.
type CommandBus interface {
	HandleCommand(ctx context.Context, cmd Command) error
}

type MockCommandBus struct {
	HandleCommandFunc func(ctx context.Context, cmd Command) error
}

func (m *MockCommandBus) HandleCommand(ctx context.Context, cmd Command) error {
	return m.HandleCommandFunc(ctx, cmd)
}

var NoopCommandBus = &MockCommandBus{
	HandleCommandFunc: func(ctx context.Context, cmd Command) error {
		return nil
	},
}
