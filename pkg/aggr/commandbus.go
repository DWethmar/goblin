package aggr

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrCommandHandlerNotFound = errors.New("command handler not found")
)

type CommandBus struct {
	aggregateStore AggregateStore
	eventBus       *EventBus
}

func (b *CommandBus) Dispatch(ctx context.Context, commands ...Command) error {
	aggregates := make([]*Aggregate, 0, len(commands))
	events := make([]*Event, 0, len(commands))

	for _, command := range commands {
		// cache aggregate to avoid multiple Get calls
		var aggregate *Aggregate
		for _, a := range aggregates {
			if a.AggregateID() == command.AggregateID() {
				aggregate = a
				break
			}
		}

		// get aggregate from store if not cached
		if aggregate == nil {
			var err error
			aggregate, err = b.aggregateStore.Get(ctx, command.AggregateType(), command.AggregateID())
			if err != nil {
				return fmt.Errorf("failed to get aggregate: %w", err)
			}

			if aggregate == nil {
				return fmt.Errorf("aggregate is nil")
			}

			aggregates = append(aggregates, aggregate)
		}

		event, err := aggregate.HandleCommand(command)
		if err != nil {
			return fmt.Errorf("failed to dispatch command: %w", err)
		}

		if err := aggregate.HandleEvent(ctx, event); err != nil {
			return fmt.Errorf("failed to dispatch command: %w", err)
		}

		events = append(events, event)
	}

	if err := b.aggregateStore.Save(ctx, aggregates...); err != nil {
		return fmt.Errorf("failed to save aggregate: %w", err)
	}

	for _, event := range events {
		if err := b.eventBus.Publish(ctx, event); err != nil {
			return fmt.Errorf("failed to publish event: %w", err)
		}
	}

	return nil
}

// NewCommandBus returns a new instance of CommandBus.
func NewCommandBus(aggregateStore AggregateStore, eventBus *EventBus) *CommandBus {
	return &CommandBus{
		aggregateStore: aggregateStore,
		eventBus:       eventBus,
	}
}
