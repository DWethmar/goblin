package game

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dwethmar/goblin/pkg/domain/actor"
	actorMemory "github.com/dwethmar/goblin/pkg/domain/actor/memory"
	"github.com/dwethmar/goblin/pkg/es"
	"github.com/dwethmar/goblin/pkg/es/aggregate"
	"github.com/dwethmar/goblin/pkg/es/event"
	"github.com/dwethmar/goblin/pkg/services"
)

type Options struct {
	Logger     *slog.Logger
	EventStore event.Store
}

type Game struct {
	logger       *slog.Logger
	actorService *services.Actors
}

func (g *Game) DispatchStringCommand(ctx context.Context, agregateID, cmd string, args ...string) error {
	if agregateID == "" {
		return fmt.Errorf("aggregate id is required")
	}

	switch cmd {
	case "crta":
		return g.CreateActor(agregateID, args[0])
	}

	return fmt.Errorf("unknown command: %s", cmd)
}

func New(ctx context.Context, opt Options) (*Game, error) {
	logger := opt.Logger

	actorRepo := actorMemory.NewRepository()
	// Create the event bus and add event handlers
	eventBus := es.NewEventBus()
	eventBus.Subscribe(actor.ActorEventMatcher, actor.ActorSinkHandler(ctx, actorRepo))

	// Create the agregate factory and register agregates
	aggregateFactory := aggregate.NewFactory()
	actor.RegisterFactory(aggregateFactory)

	aggregateStore := aggregate.NewStore(opt.EventStore, aggregateFactory)
	commandBus := es.NewCommandBus(aggregateStore, eventBus)

	return &Game{
		logger:       logger,
		actorService: services.NewActorService(actorRepo, commandBus),
	}, nil
}
