package game

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dwethmar/tards/cmd/game/actor"
	"github.com/dwethmar/tards/pkg/es"
	"github.com/dwethmar/tards/pkg/es/file"
)

type Options struct {
	Logger *slog.Logger
}

func Run(opt Options) error {
	logger := opt.Logger

	// eventRepo := memory.NewEventRepository()
	eventRepo := file.NewEventRepository(
		map[string]func(*file.LogEntry) *es.Event{
			actor.CreatedEventType: func(le *file.LogEntry) *es.Event {
				var data actor.CreatedEventData
				if err := json.Unmarshal(le.Data, &data); err != nil {
					logger.Error("unmarshal", "error", err)
					return nil
				}

				return &es.Event{
					AggregateID: le.AggregateID,
					Type:        le.Type,
					Data:        data,
				}
			},
		},
	)

	aggregateFactory := es.NewAggregateFactory()

	// Register the factories
	actor.RegisterFactory(aggregateFactory)

	aggregateStore := es.NewAggregateStore(eventRepo, aggregateFactory)

	eventBus := es.NewEventBus()

	commandBus := es.NewCommandBus(aggregateStore, eventBus)

	err := commandBus.Dispatch(&actor.CreateCommand{
		ActorID: "99",
		Name:    "test",
	})
	if err != nil {
		return fmt.Errorf("dispatching command: %w", err)
	}

	actor, err := aggregateStore.Get(actor.AggregateType, "99")
	if err != nil {
		return fmt.Errorf("getting actor: %w", err)
	}

	logger.Info("actor", "model", actor.Model, "aggregate", actor)

	return nil
}
