package aggr

import (
	"context"
	"fmt"
)

// CommandBus applies commands to aggregates and publishes events.
type CommandBus struct {
	aggregateStore AggregateStore
	eventBus       EventHandler
}

func (b *CommandBus) HandleCommand(ctx context.Context, cmd Command) (*Event, error) {
	a, err := b.aggregateStore.Get(ctx, cmd.AggregateType(), cmd.AggregateID())
	if err != nil {
		return nil, fmt.Errorf("failed to get aggregate from store: %w", err)
	}

	if a == nil {
		return nil, fmt.Errorf("aggregate is nil")
	}

	defer func() {
		a.ClearAggregateEvents()
	}()

	event, err := a.HandleCommand(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to use command on aggregate: %w", err)
	}

	if err := a.HandleEvent(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to apply event on aggregate: %w", err)
	}

	if err := b.aggregateStore.Save(ctx, a); err != nil {
		return nil, fmt.Errorf("failed to save aggregate: %w", err)
	}

	if err := b.eventBus.HandleEvent(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to publish event: %w", err)
	}

	return event, nil
}

// NewCommandBus returns a new instance of CommandBus.
func NewCommandBus(aggregateStore AggregateStore, eventBus EventHandler) *CommandBus {
	return &CommandBus{
		aggregateStore: aggregateStore,
		eventBus:       eventBus,
	}
}
