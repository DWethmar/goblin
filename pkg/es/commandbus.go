package es

import (
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

func (b *CommandBus) Dispatch(command Command) error {
	aggregate, err := b.aggregateStore.Get(command.AggregateType(), command.AggregateID())
	if err != nil {
		return fmt.Errorf("failed to get aggregate: %w", err)
	}

	event, err := aggregate.HandleCommand(command)
	if err != nil {
		return fmt.Errorf("failed to dispatch command: %w", err)
	}

	if event == nil {
		return nil
	}

	if err := aggregate.HandleEvent(event); err != nil {
		return fmt.Errorf("failed to handle event: %w", err)
	}

	if err := b.aggregateStore.Save(aggregate); err != nil {
		return fmt.Errorf("failed to save aggregate: %w", err)
	}

	for _, event := range aggregate.Events {
		if err := b.eventBus.Publish(event); err != nil {
			return fmt.Errorf("failed to publish event: %w", err)
		}
	}

	clear(aggregate.Events)

	return nil
}

// NewCommandBus returns a new instance of CommandBus.
func NewCommandBus(aggregateStore AggregateStore, eventBus *EventBus) *CommandBus {
	return &CommandBus{
		aggregateStore: aggregateStore,
		eventBus:       eventBus,
	}
}
