package game

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/dwethmar/goblin/pkg/domain/actor"
	actorMemory "github.com/dwethmar/goblin/pkg/domain/actor/memory"
	"github.com/dwethmar/goblin/pkg/es"
	"github.com/dwethmar/goblin/pkg/es/aggregate"
	"github.com/dwethmar/goblin/pkg/es/event"
	"github.com/dwethmar/goblin/pkg/es/replay"
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

func (g *Game) ExecStringCommand(ctx context.Context, cmdStr string) error {
	args := strings.Split(cmdStr, " ")
	fmt.Printf("args: %v\n", args)

	if len(args) < 2 {
		return fmt.Errorf("command is invalid")
	}

	cmd := args[0]

	switch cmd {
	case "create":
		switch args[1] {
		case "actor":
			if len(args) < 4 {
				return fmt.Errorf("command is invalid")
			}
			id := args[2]
			name := args[3]
			return g.actorService.Create(id, name)
		}
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

	// replay all events and rebuild the state
	replayer := replay.New(logger, opt.EventStore, eventBus)
	if err := replayer.Replay(ctx); err != nil {
		return nil, fmt.Errorf("replaying events: %w", err)
	}

	return &Game{
		logger:       logger,
		actorService: services.NewActorService(actorRepo, commandBus),
	}, nil
}
