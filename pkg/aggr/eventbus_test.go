package aggr

import (
	"context"
	"reflect"
	"testing"
)

func TestEventBus_Subscribe(t *testing.T) {
	t.Run("should subscribe to event", func(t *testing.T) {
		bus := NewEventBus()

		matcher := MatchEvents{"event1"}

		bus.Subscribe(matcher, EventHandlerFunc(func(ctx context.Context, event *Event) error {
			return nil
		}))

		if len(bus.handlers) != 1 {
			t.Errorf("Expected 1 handler, got %d", len(bus.handlers))
		}
	})

	t.Run("should unsubscribe from event", func(t *testing.T) {
		bus := NewEventBus()

		matcher := MatchEvents{"event1"}

		unsubscribe := bus.Subscribe(matcher, EventHandlerFunc(func(ctx context.Context, event *Event) error {
			return nil
		}))

		unsubscribe()

		if len(bus.handlers) != 0 {
			t.Errorf("Expected 0 handler, got %d", len(bus.handlers))
		}
	})
}

func TestEventBus_HandleEvent(t *testing.T) {
	type fields struct {
		handlers []*handlerMatcherPair
	}
	type args struct {
		ctx   context.Context
		event *Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := &EventBus{
				handlers: tt.fields.handlers,
			}
			if err := bus.HandleEvent(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("EventBus.HandleEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEventBus(t *testing.T) {
	tests := []struct {
		name string
		want *EventBus
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventBus(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventBus() = %v, want %v", got, tt.want)
			}
		})
	}
}
