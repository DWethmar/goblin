// Package replay provides a way to replay events from the event store to the event bus.
package replay

import (
	"context"
	"log/slog"
	"sync"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/aggr/event"
	"github.com/dwethmar/goblin/pkg/aggr/sink"
)

type Replayer struct {
	logger     *slog.Logger
	eventStore event.Repository
	eventBus   *aggr.EventBus
}

func New(
	logger *slog.Logger,
	eventStore event.Repository,
	eventBus *aggr.EventBus,
) *Replayer {
	return &Replayer{
		logger:     logger,
		eventStore: eventStore,
		eventBus:   eventBus,
	}
}

func (r *Replayer) Replay(ctx context.Context) error {
	errCn := make(chan error)
	eventSink := sink.New(r.eventStore.All(errCn), []aggr.EventHandler{r.eventBus})

	done := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		eventSink.Sink(ctx)
		done <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-done:
			return nil
		case err := <-errCn:
			return err
		case err := <-eventSink.Errors():
			return err
		}
	}
}
