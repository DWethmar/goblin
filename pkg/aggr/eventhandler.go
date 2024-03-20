package aggr

import "context"

// EventHandler handles events.
type EventHandler interface {
	HandleEvent(context.Context, *Event) error
}

// EventHandlerFunc is a function that implements the EventHandler interface.
type EventHandlerFunc func(context.Context, *Event) error

func (f EventHandlerFunc) HandleEvent(ctx context.Context, e *Event) error {
	return f(ctx, e)
}
