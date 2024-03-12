package eventkv

import (
	"fmt"

	"github.com/dwethmar/goblin/pkg/es"
	"github.com/dwethmar/goblin/pkg/es/event"
	"github.com/dwethmar/goblin/pkg/kv"
)

type Repository struct {
	kv           kv.DB
	eventDecoder event.Decoder
	eventEncoder event.Encoder
}

func eventID(aggregateID string, version int) []byte {
	return []byte(fmt.Sprintf("%s-%d", aggregateID, version))
}

func eventIDPrefix(aggregateID string) []byte {
	return []byte(fmt.Sprintf("%s-", aggregateID))
}

func (r *Repository) Add(events []*es.Event) error {
	for _, event := range events {
		b, err := r.eventEncoder.Encode(event)
		if err != nil {
			return fmt.Errorf("encoding event: %w", err)
		}

		if err := r.kv.Put(eventID(event.AggregateID, event.Version), b); err != nil {
			return fmt.Errorf("putting event: %w", err)
		}
	}

	return nil
}

func (r *Repository) List(aggregateID string) ([]*es.Event, error) {
	var events []*es.Event
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

func New(kv kv.DB,
	eventDecoder event.Decoder,
	eventEncoder event.Encoder,
) *Repository {
	return &Repository{
		kv:           kv,
		eventDecoder: eventDecoder,
		eventEncoder: eventEncoder,
	}
}
