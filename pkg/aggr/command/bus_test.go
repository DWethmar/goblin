package command

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/dwethmar/goblin/pkg/aggr"
)

func TestBus_Dispatch(t *testing.T) {
	t.Run("should dispatch command", func(t *testing.T) {
		a := &aggr.Aggregate{
			Model: &aggr.MockModel{
				CommandHandler: aggr.CommandHandlerFunc(func(_ context.Context, c aggr.Command) (*aggr.Event, error) {
					return &aggr.Event{
						AggregateID: c.AggregateID(),
						Type:        "test",
						Data:        "test",
						Version:     1,
						Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					}, nil
				}),
				EventHandler: aggr.EventHandlerFunc(func(_ context.Context, e *aggr.Event) error { return nil }),
			},
		}

		aggregateStore := &aggr.MockAggregateStore{
			GetFunc:  func(ctx context.Context, aggregateType, aggregateID string) (*aggr.Aggregate, error) { return a, nil },
			SaveFunc: func(_ context.Context, _ ...*aggr.Aggregate) error { return nil },
		}

		eventBus := &aggr.EventBus{}
		Bus := NewBus(aggregateStore, eventBus)
		err := Bus.HandleCommand(context.TODO(), &aggr.MockCommand{
			ID: "1",
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("should get aggregate and clear its events", func(t *testing.T) {
		a := &aggr.Aggregate{
			Model: &aggr.MockModel{
				CommandHandler: aggr.CommandHandlerFunc(func(_ context.Context, c aggr.Command) (*aggr.Event, error) {
					return &aggr.Event{
						AggregateID: c.AggregateID(),
						Type:        "test",
						Data:        "test",
						Version:     1,
						Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					}, nil
				}),
				EventHandler: aggr.EventHandlerFunc(func(_ context.Context, e *aggr.Event) error { return nil }),
			},
		}

		aggregateStore := &aggr.MockAggregateStore{
			GetFunc:  func(ctx context.Context, aggregateType, aggregateID string) (*aggr.Aggregate, error) { return a, nil },
			SaveFunc: func(_ context.Context, _ ...*aggr.Aggregate) error { return nil },
		}

		eventBus := &aggr.EventBus{}
		bus := NewBus(aggregateStore, eventBus)
		err := bus.HandleCommand(context.TODO(), &aggr.MockCommand{
			ID: "1",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(a.Model.AggregateEvents()) != 0 {
			t.Errorf("unexpected events: %v", a.Model.AggregateEvents())
		}
	})
}

func TestNewBus(t *testing.T) {
	t.Run("should return new instance of Bus", func(t *testing.T) {
		aggregateStore := &aggr.MockAggregateStore{}
		eventBus := &aggr.EventBus{}
		got := NewBus(aggregateStore, eventBus)
		want := &Bus{
			aggregateStore: aggregateStore,
			eventBus:       eventBus,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NewBus() = %v, want %v", got, want)
		}
	})
}
