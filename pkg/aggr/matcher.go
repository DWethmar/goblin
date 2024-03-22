package aggr

type Matcher interface {
	// Match returns true if the event matches the matcher.
	Match(*Event) bool
}

type MatcherFunc func(*Event) bool

func (f MatcherFunc) Match(event *Event) bool {
	return f(event)
}

// MatchEvents is a matcher that matches events by types
type MatchEvents []string

// Match returns true if the event type is in the list of events.
func (m MatchEvents) Match(event *Event) bool {
	for _, e := range m {
		if e == event.EventType {
			return true
		}
	}
	return false
}
