package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dwethmar/goblin/pkg/domain/actor"
	actorMemory "github.com/dwethmar/goblin/pkg/domain/actor/memory"
	"github.com/dwethmar/goblin/pkg/es"
	"github.com/dwethmar/goblin/pkg/es/aggregate"
	eventEncoding "github.com/dwethmar/goblin/pkg/es/event/encoding"
	eventkv "github.com/dwethmar/goblin/pkg/es/event/kv"
	"github.com/dwethmar/goblin/pkg/es/replay"
	"github.com/dwethmar/goblin/pkg/game"
	"github.com/dwethmar/goblin/pkg/kv/bbolt"
	"github.com/dwethmar/goblin/pkg/services"
)

type Config struct {
	Logger *slog.Logger
	Game   string
}

func SetupGame(ctx context.Context, c Config) (*game.Game, func() error, error) {
	db, err := bbolt.Connect(fmt.Sprintf("./.tmp/%s.db", Game))
	if err != nil {
		return nil, nil, fmt.Errorf("connecting to db: %w", err)
	}

	eventStore := eventkv.New(bbolt.New([]byte("events"), db), &eventEncoding.Decoder{}, &eventEncoding.Encoder{})

	actorRepo := actorMemory.NewRepository()
	// Create the event bus and add event handlers
	eventBus := es.NewEventBus()
	eventBus.Subscribe(actor.ActorEventMatcher, actor.ActorSinkHandler(ctx, actorRepo))

	// Create the agregate factory and register agregates
	aggregateFactory := aggregate.NewFactory()
	actor.RegisterFactory(aggregateFactory)

	// Create the command bus
	aggregateStore := aggregate.NewStore(eventStore, aggregateFactory)
	commandBus := es.NewCommandBus(aggregateStore, eventBus)

	// replay all events and rebuild the state
	replayer := replay.New(c.Logger, eventStore, eventBus)
	if err := replayer.Replay(ctx); err != nil {
		return nil, nil, fmt.Errorf("replaying events: %w", err)
	}

	g, err := game.New(ctx, game.Options{
		Logger:       c.Logger,
		ActorService: services.NewActorService(actorRepo, commandBus),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("creating game: %w", err)
	}

	close := func() error {
		return db.Close()
	}

	return g, close, nil
}
