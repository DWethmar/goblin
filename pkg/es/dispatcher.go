package es

type CommandHandler interface {
	// Dispatch applies a command to the aggregate.
	Dispatch(Command) (*Event, error)
}

type DispatcherFunc func(Command) (*Event, error)

func (f DispatcherFunc) Dispatch(c Command) (*Event, error) {
	return f(c)
}
