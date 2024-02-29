package game

import (
	"fmt"
	"log/slog"

	"github.com/dwethmar/tards/cmd/game/actor"
	"github.com/dwethmar/tards/pkg/es"
	"github.com/dwethmar/tards/pkg/es/memory"
)

type Options struct {
	Logger *slog.Logger
}

func Run(opt Options) error {
	logger := opt.Logger

	eventRepo := memory.NewEventRepository()

	aggregateFactory := es.NewAggregateFactory()

	// Register the factories
	actor.RegisterFactory(aggregateFactory)

	aggregateStore := es.NewAggregateStore(eventRepo, aggregateFactory)

	eventBus := es.NewEventBus()

	commandBus := es.NewCommandBus(aggregateStore, eventBus)

	err := commandBus.Dispatch(&actor.CreateCommand{
		ActorID: "1",
		Name:    "test",
	})
	if err != nil {
		return fmt.Errorf("dispatching command: %w", err)
	}

	actor, err := aggregateStore.Get(actor.AggregateType, "1")
	if err != nil {
		return fmt.Errorf("getting actor: %w", err)
	}

	logger.Info("actor", "actor", actor)

	return nil
}
