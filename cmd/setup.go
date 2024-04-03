package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/aggr/command"
	"github.com/dwethmar/goblin/pkg/aggr/es"
	eventencoding "github.com/dwethmar/goblin/pkg/aggr/event/encoding"
	eventkv "github.com/dwethmar/goblin/pkg/aggr/event/kv"
	"github.com/dwethmar/goblin/pkg/domain/actor"
	actormemory "github.com/dwethmar/goblin/pkg/domain/actor/memory"
	"github.com/dwethmar/goblin/pkg/domain/chunk"
	chunkmemory "github.com/dwethmar/goblin/pkg/domain/chunk/memory"
	"github.com/dwethmar/goblin/pkg/domain/replay"
	"github.com/dwethmar/goblin/pkg/game"
	"github.com/dwethmar/goblin/pkg/kv/bbolt"
	"github.com/dwethmar/goblin/pkg/services"
)

type Config struct {
	Logger     *slog.Logger
	Game       string
	LogDBStats bool
}

func SetupGame(ctx context.Context, c Config) (*game.Game, func() error, error) {
	db, err := bbolt.Connect(fmt.Sprintf("./.tmp/%s.db", c.Game))
	if err != nil {
		return nil, nil, fmt.Errorf("connecting to db: %w", err)
	}

	if c.LogDBStats {
		// Log the database stats every 5 seconds
		bbolt.Stats(ctx, db, time.Second*5, c.Logger)
	}

	bboltDB := bbolt.New([]byte("events"), db)
	eventStore := eventkv.New(bboltDB, &eventencoding.Decoder{}, &eventencoding.Encoder{})

	// Create the event bus and add event handlers
	eventBus := aggr.NewEventBus()

	actorRepo := actormemory.NewRepository()
	eventBus.Subscribe(actor.MatchAllEvents, actor.ActorSinkHandler(actorRepo))

	chunkRepo := chunkmemory.NewRepository()
	eventBus.Subscribe(chunk.MatchAllEvents, chunk.ChunkSinkHandler(chunkRepo))

	// Create the agregate factory and register agregates
	aggregateFactory := aggr.NewFactory(actor.RegisterFactory, chunk.RegisterFactory)

	// Create the command bus
	aggregateStore := es.NewAggregateStore(eventStore, aggregateFactory)
	commandBus := command.NewCommandBus(aggregateStore, eventBus)

	// replay all events and rebuild the state
	replayer := replay.New(c.Logger, eventStore, eventBus)
	if err := replayer.Replay(ctx, actorRepo, chunkRepo); err != nil {
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
