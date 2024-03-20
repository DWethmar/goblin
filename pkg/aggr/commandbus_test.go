package aggr

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCommandBus_Dispatch(t *testing.T) {
	t.Run("should dispatch command", func(t *testing.T) {
		a := &Aggregate{
			Model: &MockAggregate{
				CommandHandlerFunc: func(_ context.Context, c Command) (*Event, error) {
					return &Event{
						AggregateID: c.AggregateID(),
						Type:        "test",
						Data:        "test",
						Version:     1,
						Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					}, nil
				},
				EventHandlerFunc: func(_ context.Context, e *Event) error { return nil },
			},
		}

		aggregateStore := &MockAggregateStore{
			GetFunc:  func(ctx context.Context, aggregateType, aggregateID string) (*Aggregate, error) { return a, nil },
			SaveFunc: func(_ context.Context, _ ...*Aggregate) error { return nil },
		}

		eventBus := &EventBus{}
		commandBus := NewCommandBus(aggregateStore, eventBus)
		e, err := commandBus.HandleCommand(context.TODO(), &MockCommand{
			aggregateID: "1",
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if e == nil {
			t.Errorf("expected event, got nil")
		}

		want := &Event{
			AggregateID: "1",
			Type:        "test",
			Data:        "test",
			Version:     1,
			Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		if diff := cmp.Diff(e, want); diff != "" {
			t.Errorf("unexpected event (-got +want):\n%s", diff)
		}
	})

	t.Run("should get aggregate and clear its events", func(t *testing.T) {
		a := &Aggregate{
			Model: &MockAggregate{
				CommandHandlerFunc: func(_ context.Context, c Command) (*Event, error) {
					return &Event{
						AggregateID: c.AggregateID(),
						Type:        "test",
						Data:        "test",
						Version:     1,
						Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					}, nil
				},
				EventHandlerFunc: func(_ context.Context, e *Event) error { return nil },
			},
		}

		aggregateStore := &MockAggregateStore{
			GetFunc:  func(ctx context.Context, aggregateType, aggregateID string) (*Aggregate, error) { return a, nil },
			SaveFunc: func(_ context.Context, _ ...*Aggregate) error { return nil },
		}

		eventBus := &EventBus{}
		commandBus := NewCommandBus(aggregateStore, eventBus)
		_, err := commandBus.HandleCommand(context.TODO(), &MockCommand{
			aggregateID: "1",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(a.Model.AggregateEvents()) != 0 {
			t.Errorf("unexpected events: %v", a.Model.AggregateEvents())
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
