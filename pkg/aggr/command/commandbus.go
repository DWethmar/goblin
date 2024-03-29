package command

import (
	"context"
	"fmt"

	"github.com/dwethmar/goblin/pkg/aggr"
)

var _ aggr.CommandBus = &CommandBus{}

// CommandBus applies commands to aggregates and publishes events.
type CommandBus struct {
	aggregateStore aggr.AggregateStore
	eventBus       aggr.EventHandler
}

func (b *CommandBus) HandleCommand(ctx context.Context, cmd aggr.Command) error {
	a, err := b.aggregateStore.Get(ctx, cmd.AggregateType(), cmd.AggregateID())
	if err != nil {
		return fmt.Errorf("failed to get aggregate from store: %w", err)
	}

	if a == nil {
		return fmt.Errorf("aggregate is nil")
	}

	defer a.ClearAggregateEvents()

	event, err := a.HandleCommand(ctx, cmd)
	if err != nil {
		return fmt.Errorf("aggregate failed to handle command: %w", err)
	}

	if err := a.HandleEvent(ctx, event); err != nil {
		return fmt.Errorf("aggregate failed to handle event: %w", err)
	}

	if err := b.aggregateStore.Save(ctx, a); err != nil {
		return fmt.Errorf("aggregate failed to save: %w", err)
	}

	if err := b.eventBus.HandleEvent(ctx, event); err != nil {
		return fmt.Errorf("event bus failed to handle event: %w", err)
	}

	return nil
}

// NewCommandBus returns a new instance of CommandBus.
func NewCommandBus(aggregateStore aggr.AggregateStore, eventBus aggr.EventHandler) *CommandBus {
	return &CommandBus{
		aggregateStore: aggregateStore,
		eventBus:       eventBus,
	}
}
