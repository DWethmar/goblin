package memory

import (
	"sync"

	"github.com/dwethmar/goblin/pkg/es"
	"github.com/dwethmar/goblin/pkg/es/event"
)

var _ event.Store = &EventRepository{}

type EventRepository struct {
	eventsMux sync.Mutex
	events    []*es.Event
}

// Create implements es.EventStore.
func (r *EventRepository) Add(events []*es.Event) error {
	r.eventsMux.Lock()
	defer r.eventsMux.Unlock()

	r.events = append(r.events, events...)
	return nil
}

// List implements es.EventStore.
func (r *EventRepository) List(aggregateID string) ([]*es.Event, error) {
	r.eventsMux.Lock()
	defer r.eventsMux.Unlock()

	events := make([]*es.Event, 0)
	for _, event := range r.events {
		if event.AggregateID == aggregateID {
			events = append(events, event)
		}
	}

	return events, nil
}

func (r *EventRepository) All(err chan<- error) <-chan *es.Event {
	outCh := make(chan *es.Event)
	go func() {
		r.eventsMux.Lock()
		defer r.eventsMux.Unlock()
		for _, event := range r.events {
			outCh <- event
		}
		close(outCh)
	}()
	return outCh
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		eventsMux: sync.Mutex{},
		events:    make([]*es.Event, 0),
	}
}
