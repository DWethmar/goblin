package memory

import (
	"sync"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/aggr/event"
)

var _ event.Repository = &EventRepository{}

type EventRepository struct {
	eventsMux sync.Mutex
	events    []*aggr.Event
}

// Create implements aggr.EventStore.
func (r *EventRepository) Add(events []*aggr.Event) error {
	r.eventsMux.Lock()
	defer r.eventsMux.Unlock()
	r.events = append(r.events, events...)
	return nil
}

// List implements aggr.EventStore.
func (r *EventRepository) List(aggregateID string) ([]*aggr.Event, error) {
	r.eventsMux.Lock()
	defer r.eventsMux.Unlock()
	events := make([]*aggr.Event, 0)
	for _, event := range r.events {
		if event.AggregateID == aggregateID {
			events = append(events, event)
		}
	}

	return events, nil
}

func (r *EventRepository) All(err chan<- error) <-chan *aggr.Event {
	outCh := make(chan *aggr.Event)
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
		events:    make([]*aggr.Event, 0),
	}
}
