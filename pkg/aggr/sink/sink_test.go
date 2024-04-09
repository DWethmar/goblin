package sink

import (
	"context"
	"sync"
	"testing"

	"github.com/dwethmar/goblin/pkg/aggr"
)

func TestSink(t *testing.T) {
	t.Run("should distribute events to multiple handlers based on the aggregate ID", func(t *testing.T) {
		var a, b, c int
		// Arrange
		ctx := context.Background()
		in := make(chan *aggr.Event)
		handlers := []aggr.EventHandler{
			aggr.EventHandlerFunc(func(ctx context.Context, event *aggr.Event) error {
				a++
				return nil
			}),
			aggr.EventHandlerFunc(func(ctx context.Context, event *aggr.Event) error {
				b++
				return nil
			}),
			aggr.EventHandlerFunc(func(ctx context.Context, event *aggr.Event) error {
				c++
				return nil
			}),
		}

		errCh := make(chan error, 1)

		// Act
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			eventSink := &EventSink{
				in:       in,
				handlers: handlers,
				errCh:    errCh,
				groupIDFunc: func(aggregateID string) uint {
					if aggregateID == "test1" {
						return 0
					}

					if aggregateID == "test2" {
						return 1
					}

					return 2
				},
			}
			eventSink.Sink(ctx)
		}()

		// three events with the same aggregate ID
		in <- &aggr.Event{
			AggregateID: "test1",
		}

		in <- &aggr.Event{
			AggregateID: "test1",
		}

		in <- &aggr.Event{
			AggregateID: "test1",
		}

		// two events with the same aggregate ID
		in <- &aggr.Event{
			AggregateID: "test2",
		}

		in <- &aggr.Event{
			AggregateID: "test2",
		}

		// one event with a different aggregate ID
		in <- &aggr.Event{
			AggregateID: "test3",
		}

		close(in)
		wg.Wait()

		// Assert
		if a != 3 {
			t.Errorf("expected a to be 3, got %d", a)
		}

		if b != 2 {
			t.Errorf("expected b to be 2, got %d", b)
		}

		if c != 1 {
			t.Errorf("expected c to be 1, got %d", c)
		}

		// Assert
		select {
		case err := <-errCh:
			t.Errorf("expected err to be nil, got %v", err)
		default:
		}
	})
}
