package eventkv

import (
	"fmt"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/aggr/event"
	"github.com/dwethmar/goblin/pkg/kv"
)

var _ event.Repository = &EventRepository{}

type EventRepository struct {
	kv           kv.DB
	eventDecoder event.Decoder
	eventEncoder event.Encoder
}

func eventID(aggregateID string, version uint) []byte {
	return []byte(fmt.Sprintf("%s-%d", aggregateID, version))
}

func eventIDPrefix(aggregateID string) []byte {
	return []byte(fmt.Sprintf("%s-", aggregateID))
}

func (r *EventRepository) Add(events []*aggr.Event) error {
	entries := make([]kv.Entry, len(events))
	for i, event := range events {
		b, err := r.eventEncoder.Encode(event)
		if err != nil {
			return fmt.Errorf("encoding event: %w", err)
		}

		entries[i] = kv.Entry{
			Key:   eventID(event.AggregateID, event.Version),
			Value: b,
		}
	}

	if err := r.kv.Put(entries...); err != nil {
		return fmt.Errorf("putting event: %w", err)
	}

	return nil
}

func (r *EventRepository) List(aggregateID string) ([]*aggr.Event, error) {
	var events []*aggr.Event
	if err := r.kv.IterateWithPrefix(eventIDPrefix(aggregateID), func(k, v []byte) error {
		event, err := r.eventDecoder.Decode(v)
		if err != nil {
			return fmt.Errorf("decoding event: %w", err)
		}
		events = append(events, event)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("iterating over events: %w", err)
	}

	return events, nil
}

// All returns all events in the store.
func (r *EventRepository) All(errCh chan<- error) <-chan *aggr.Event {
	outCh := make(chan *aggr.Event)
	go func() {
		defer close(outCh)
		if err := r.kv.Iterate(func(k, v []byte) error {
			event, err := r.eventDecoder.Decode(v)
			if err == nil {
				outCh <- event
			}
			return err
		}); err != nil {
			errCh <- fmt.Errorf("iterating over events: %w", err)
		}
	}()
	return outCh
}

func New(kv kv.DB,
	eventDecoder event.Decoder,
	eventEncoder event.Encoder,
) *EventRepository {
	return &EventRepository{
		kv:           kv,
		eventDecoder: eventDecoder,
		eventEncoder: eventEncoder,
	}
}
