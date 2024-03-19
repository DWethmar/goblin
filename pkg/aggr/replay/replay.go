// Package replay provides a way to replay events from the event store to the event bus.
package replay

import (
	"context"
	"log/slog"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/aggr/event"
)

type Replayer struct {
	logger     *slog.Logger
	eventStore event.Store
	eventBus   *aggr.EventBus
}

func (r *Replayer) Replay(ctx context.Context) error {
	errCh := make(chan error)
	eventCh := r.eventStore.All(errCh)
	r.logger.DebugContext(ctx, "replaying events")

	for {
		select {
		case event, ok := <-eventCh:
			if !ok {
				r.logger.DebugContext(ctx, "replay done")
				return nil
			}

			r.logger.DebugContext(ctx, "replaying event", "type", event.Type, "version", event.Version, "aggregate", event.AggregateID)

			if err := r.eventBus.Publish(event); err != nil {
				return err
			}
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func New(logger *slog.Logger, eventStore event.Store, eventBus *aggr.EventBus) *Replayer {
	return &Replayer{
		logger:     logger,
		eventStore: eventStore,
		eventBus:   eventBus,
	}
}
