package file

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dwethmar/tards/pkg/es"
)

var _ es.EventStore = &EventRepository{}

type LogEntry struct {
	AggregateID string          `json:"aggregate_id"`
	Type        string          `json:"type"`
	Data        json.RawMessage `json:"data"`
}

type EventRepository struct {
	LogFile         string
	DataUnmarshaler map[string]func(*LogEntry) *es.Event
}

// Create implements es.EventStore.
func (r *EventRepository) Add(events []*es.Event) error {
	f, err := os.OpenFile(r.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	for _, event := range events {
		data, err := json.Marshal(event.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal event data: %w", err)
		}

		entry := LogEntry{
			AggregateID: event.AggregateID,
			Type:        event.Type,
			Data:        data,
		}

		if err := json.NewEncoder(f).Encode(entry); err != nil {
			return fmt.Errorf("failed to write event: %w", err)
		}
	}

	return nil
}

// List implements es.EventStore.
func (r *EventRepository) List(aggregateID string) ([]*es.Event, error) {
	// check if file exists
	if _, err := os.Stat(r.LogFile); os.IsNotExist(err) {
		return nil, nil
	}

	f, err := os.Open(r.LogFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	events := make([]*es.Event, 0)
	dec := json.NewDecoder(f)
	for {
		var entry LogEntry

		if err := dec.Decode(&entry); err != nil {
			break
		}

		if entry.AggregateID != aggregateID {
			continue
		}

		unmarshal, ok := r.DataUnmarshaler[entry.Type]
		if !ok {
			return nil, fmt.Errorf("no unmarshaler for type %q", entry.Type)
		}

		event := unmarshal(&entry)
		events = append(events, event)
	}

	return events, nil
}

func NewEventRepository(
	dataUnmarshaler map[string]func(*LogEntry) *es.Event,
) *EventRepository {
	return &EventRepository{
		LogFile:         "events.log",
		DataUnmarshaler: dataUnmarshaler,
	}
}
