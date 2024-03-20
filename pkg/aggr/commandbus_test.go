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
		aggregateStore := &MockAggregateStore{
			GetFunc: func(ctx context.Context, aggregateType, aggregateID string) (*Aggregate, error) {
				if ctx != context.TODO() {
					t.Errorf("unexpected context: %v", ctx)
				}

				return &Aggregate{
					Model: &MockAggregate{
						ID: "1",
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
				}, nil

			},
			SaveFunc: func(_ context.Context, _ ...*Aggregate) error { return nil },
		}

		eventBus := &EventBus{}
		commandBus := NewCommandBus(aggregateStore, eventBus)
		e, err := commandBus.HandleCommand(context.TODO(), &MockCommand{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
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
