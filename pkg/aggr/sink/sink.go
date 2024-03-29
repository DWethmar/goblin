// Package sink provides a simple event sink that distributes events to multiple handlers based on the aggregate ID.
package sink

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/conc"
)

// aggregateGroupID returns a group ID based on the aggregate ID.
func aggregateGroupID(aggregateID string) int {
	h := fnv.New32a()
	h.Write([]byte(aggregateID))
	return int(h.Sum32())
}

// EventSink distributes events to the same handler based on the aggregate ID.
type EventSink struct {
	in          <-chan *aggr.Event
	handlers    []aggr.EventHandler
	errCh       chan error
	groupIDFunc func(aggregateID string) int
}

func New(in <-chan *aggr.Event, handlers []aggr.EventHandler) *EventSink {
	return &EventSink{
		in:          in,
		handlers:    handlers,
		errCh:       make(chan error, 1),
		groupIDFunc: aggregateGroupID,
	}
}

func (s *EventSink) Sink(ctx context.Context) {
	taskCh := make(chan conc.Task[*aggr.Event], 100)

	go func() {
		defer close(taskCh)
		for event := range s.in {
			taskCh <- conc.Task[*aggr.Event]{
				Group: s.groupIDFunc(event.AggregateID),
				Value: event,
			}
		}
	}()

	groups, err := conc.GroupTasks(taskCh, len(s.handlers))
	if err != nil {
		s.errCh <- fmt.Errorf("failed to group tasks: %w", err)
		return
	}

	wg := sync.WaitGroup{}
	for i, group := range groups {
		wg.Add(1)
		go func(i int, in <-chan *aggr.Event) {
			defer wg.Done()
			for event := range in {
				if err := s.handlers[i].HandleEvent(ctx, event); err != nil {
					s.errCh <- fmt.Errorf("failed to handle event: %w", err)
				}
			}
		}(i, group)
	}

	wg.Wait()
}

func (s *EventSink) Errors() <-chan error {
	return s.errCh
}
