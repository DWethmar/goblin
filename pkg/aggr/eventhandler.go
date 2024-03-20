package aggr

import "context"

type EventHandler interface {
	// HandleEvent handles the event.
	HandleEvent(context.Context, *Event) error
}

type EventHandlerFunc func(context.Context, *Event) error

func (f EventHandlerFunc) HandleEvent(ctx context.Context, e *Event) error {
	return f(ctx, e)
}
