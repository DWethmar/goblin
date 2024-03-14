package replay

import (
	"context"

	"github.com/dwethmar/goblin/pkg/es"
	"github.com/dwethmar/goblin/pkg/es/event"
)

type Replayer struct {
	eventStore event.Store
	eventBus   *es.EventBus
}

func (r *Replayer) Replay(ctx context.Context) error {
	errCh := make(chan error)
	eventCh := r.eventStore.All(errCh)

	for {
		select {
		case event, ok := <-eventCh:
			if !ok {
				return nil
			}

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

func New(eventStore event.Store, eventBus *es.EventBus) *Replayer {
	return &Replayer{
		eventStore: eventStore,
		eventBus:   eventBus,
	}
}
