package es

type EventHandler interface {
	// HandleEvent handles the event.
	HandleEvent(*Event) error
}

type EventHandlerFunc func(*Event) error

func (f EventHandlerFunc) HandleEvent(e *Event) error {
	return f(e)
}
