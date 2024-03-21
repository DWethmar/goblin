package command

import (
	"context"
	"fmt"

	"github.com/dwethmar/goblin/pkg/aggr"
)

// Bus applies commands to aggregates and publishes events.
type Bus struct {
	aggregateStore aggr.AggregateStore
	eventBus       aggr.EventHandler
}

func (b *Bus) HandleCommand(ctx context.Context, cmd aggr.Command) error {
	a, err := b.aggregateStore.Get(ctx, cmd.AggregateType(), cmd.AggregateID())
	if err != nil {
		return fmt.Errorf("failed to get aggregate from store: %w", err)
	}

	if a == nil {
		return fmt.Errorf("aggregate is nil")
	}

	defer func() {
		a.ClearAggregateEvents()
	}()

	event, err := a.HandleCommand(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to use command on aggregate: %w", err)
	}

	if err := a.HandleEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to apply event on aggregate: %w", err)
	}

	if err := b.aggregateStore.Save(ctx, a); err != nil {
		return fmt.Errorf("failed to save aggregate: %w", err)
	}

	if err := b.eventBus.HandleEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// NewBus returns a new instance of CommandBus.
func NewBus(aggregateStore aggr.AggregateStore, eventBus aggr.EventHandler) *Bus {
	return &Bus{
		aggregateStore: aggregateStore,
		eventBus:       eventBus,
	}
}
