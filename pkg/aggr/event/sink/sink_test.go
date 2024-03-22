package sink

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/dwethmar/goblin/pkg/aggr"
)

func TestNewWorkerSink(t *testing.T) {
	t.Run("NewWorkerSink", func(t *testing.T) {
		numWorkers := 4
		handler := aggr.EventHandlerFunc(func(_ context.Context, e *aggr.Event) error { return nil })
		ws := NewWorkerSink([]aggr.EventHandler{handler, handler, handler, handler})

		if ws.numWorkers != numWorkers {
			t.Errorf("Expected numWorkers to be %d, got %d", numWorkers, ws.numWorkers)
		}

		if len(ws.workerChannels) != numWorkers {
			t.Errorf("Expected workerChannels to have %d elements, got %d", numWorkers, len(ws.workerChannels))
		}

		if len(ws.handlers) != numWorkers {
			t.Errorf("Expected handlers to have %d elements, got %d", numWorkers, len(ws.handlers))
		}
	})
}

func TestWorkerSink_ProcessEvent(t *testing.T) {
	t.Run("ProcessEvent", func(t *testing.T) {
		handler := aggr.EventHandlerFunc(func(_ context.Context, e *aggr.Event) error { return nil })
		ws := NewWorkerSink([]aggr.EventHandler{handler, handler, handler, handler})

		event := &aggr.Event{
			AggregateID: "1",
		}

		ws.ProcessEvent(context.Background(), event)
	})

	t.Run("check if every event is there", func(t *testing.T) {
		expected := 10
		processedEvents := make(map[string]bool)
		mu := &sync.Mutex{}

		handler := aggr.EventHandlerFunc(func(_ context.Context, e *aggr.Event) error {
			mu.Lock()
			processedEvents[e.AggregateID] = true
			mu.Unlock()
			return nil
		})

		ws := NewWorkerSink([]aggr.EventHandler{handler, handler})

		for i := 0; i < int(expected); i++ {
			ws.ProcessEvent(context.Background(), &aggr.Event{
				AggregateID: fmt.Sprintf("%d", i),
			})
		}

		ws.Close()
		mu.Lock()
		if len(processedEvents) != int(expected) {
			t.Errorf("Expected %d unique events to be processed, got %d", expected, len(processedEvents))
		}
		mu.Unlock()
	})

	t.Run("check if amount of event handler calls is the same", func(t *testing.T) {
		expected := 1000000
		var actual1 int32 = 0
		var actual2 int32 = 0
		var actual3 int32 = 0

		handler1 := aggr.EventHandlerFunc(func(_ context.Context, e *aggr.Event) error {
			atomic.AddInt32(&actual1, 1)
			return nil
		})

		handler2 := aggr.EventHandlerFunc(func(_ context.Context, e *aggr.Event) error {
			atomic.AddInt32(&actual2, 1)
			return nil
		})

		handler3 := aggr.EventHandlerFunc(func(_ context.Context, e *aggr.Event) error {
			atomic.AddInt32(&actual3, 1)
			return nil
		})

		ws := NewWorkerSink([]aggr.EventHandler{
			handler1, handler1, handler1, handler1,
			handler2, handler2, handler2, handler2,
			handler3, handler3, handler3, handler3,
		})

		for i := 0; i < int(expected); i++ {
			ws.ProcessEvent(context.Background(), &aggr.Event{
				AggregateID: fmt.Sprintf("%d", i),
			})
		}

		ws.Close()

		actual := actual1 + actual2 + actual3
		if int(actual) != int(expected) {
			t.Errorf("Expected %d event handler calls, got %d", expected, actual)
		}
	})
}

func TestWorkerSink_getWorkerIndex(t *testing.T) {
	type fields struct {
		workerChannels []chan *aggr.Event
		numWorkers     int
		handlers       []aggr.EventHandler
		wg             sync.WaitGroup
		errCh          chan error
	}
	type args struct {
		aggregateID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := &WorkerSink{
				workerChannels: tt.fields.workerChannels,
				numWorkers:     tt.fields.numWorkers,
				handlers:       tt.fields.handlers,
				wg:             tt.fields.wg,
				errCh:          tt.fields.errCh,
			}
			if got := ws.getWorkerIndex(tt.args.aggregateID); got != tt.want {
				t.Errorf("WorkerSink.getWorkerIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkerSink_Errors(t *testing.T) {
	type fields struct {
		workerChannels []chan *aggr.Event
		numWorkers     int
		handlers       []aggr.EventHandler
		wg             sync.WaitGroup
		errCh          chan error
	}
	tests := []struct {
		name   string
		fields fields
		want   <-chan error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := &WorkerSink{
				workerChannels: tt.fields.workerChannels,
				numWorkers:     tt.fields.numWorkers,
				handlers:       tt.fields.handlers,
				wg:             tt.fields.wg,
				errCh:          tt.fields.errCh,
			}
			if got := ws.Errors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WorkerSink.Errors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkerSink_Close(t *testing.T) {
	type fields struct {
		workerChannels []chan *aggr.Event
		numWorkers     int
		handlers       []aggr.EventHandler
		wg             sync.WaitGroup
		errCh          chan error
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := &WorkerSink{
				workerChannels: tt.fields.workerChannels,
				numWorkers:     tt.fields.numWorkers,
				handlers:       tt.fields.handlers,
				wg:             tt.fields.wg,
				errCh:          tt.fields.errCh,
			}
			ws.Close()
		})
	}
}
