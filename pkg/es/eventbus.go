package es

import "fmt"

// handlerMatcherPair struct to associate a matcher with an event handler.
type handlerMatcherPair struct {
	matcher Matcher
	handler EventHandler
}

type EventBus struct {
	handlers []handlerMatcherPair
}

// Subscribe adds a new matcher and handler to the EventBus.
func (bus *EventBus) Subscribe(matcher Matcher, handler EventHandler) {
	bus.handlers = append(bus.handlers, handlerMatcherPair{matcher: matcher, handler: handler})
}

// Publish publishes the event to the EventBus.
func (bus *EventBus) Publish(event *Event) error {
	for _, pair := range bus.handlers {
		if pair.matcher.Match(event) {
			if err := pair.handler.HandleEvent(event); err != nil {
				return fmt.Errorf("handle event: %w", err)
			}
		}
	}

	return nil
}

// NewEventBus returns a new EventBus.
func NewEventBus() *EventBus {
	return &EventBus{}
}
