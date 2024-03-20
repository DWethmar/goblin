package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/aggr/aggrstore"
	eventEncoding "github.com/dwethmar/goblin/pkg/aggr/event/encoding"
	eventkv "github.com/dwethmar/goblin/pkg/aggr/event/kv"
	"github.com/dwethmar/goblin/pkg/aggr/replay"
	"github.com/dwethmar/goblin/pkg/domain/actor"
	actorMemory "github.com/dwethmar/goblin/pkg/domain/actor/memory"
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

	// Log the database stats every 5 seconds
	bbolt.Stats(ctx, db, time.Second*5, c.Logger)

	bboltDB := bbolt.New([]byte("events"), db)
	eventStore := eventkv.New(bboltDB, &eventEncoding.Decoder{}, &eventEncoding.Encoder{})

	actorRepo := actorMemory.NewRepository()
	// Create the event bus and add event handlers
	eventBus := aggr.NewEventBus()
	eventBus.Subscribe(actor.ActorEventsMatcher, actor.ActorSinkHandler(actorRepo))

	// Create the agregate factory and register agregates
	aggregateFactory := aggrstore.NewFactory(actor.RegisterFactory)

	// Create the command bus
	aggregateStore := aggrstore.NewStore(eventStore, aggregateFactory)
	commandBus := aggr.NewCommandBus(aggregateStore, eventBus)

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

	close := func() error { return db.Close() }

	return g, close, nil
}
