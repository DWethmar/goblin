package es

type Matcher interface {
	// Match returns true if the event matches the matcher.
	Match(Event) bool
}

type MatcherFunc func(Event) bool

func (f MatcherFunc) Match(event Event) bool {
	return f(event)
}

type MatchEventType string

func (t MatchEventType) Match(event Event) bool {
	return event.Type == string(t)
}
