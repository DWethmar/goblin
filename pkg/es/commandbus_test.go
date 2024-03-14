package es

import (
	"reflect"
	"testing"
	"time"
)

func TestCommandBus_Dispatch(t *testing.T) {
	t.Run("should dispatch command", func(t *testing.T) {
		aggregateStore := &MockAggregateStore{
			GetFunc: func(aggregateType, aggregateID string) (Aggregate, error) {
				return &MockAggregate{
					ID: "1",
					CommandHandlerFunc: func(c Command) (*Event, error) {
						return &Event{
							AggregateID: "1",
							Type:        "test",
							Data:        "test",
							Version:     1,
							CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						}, nil
					},
					EventHandlerFunc: func(e *Event) error { return nil },
				}, nil
			},
			SaveFunc: func(p0 Aggregate) error { return nil },
		}

		eventBus := &EventBus{}
		commandBus := NewCommandBus(aggregateStore, eventBus)
		err := commandBus.Dispatch(&MockCommand{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestNewCommandBus(t *testing.T) {
	t.Run("should return new instance of CommandBus", func(t *testing.T) {
		aggregateStore := &MockAggregateStore{}
		eventBus := &EventBus{}
		got := NewCommandBus(aggregateStore, eventBus)
		want := &CommandBus{
			aggregateStore: aggregateStore,
			eventBus:       eventBus,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewCommandBus() = %v, want %v", got, want)
		}
	})
}
