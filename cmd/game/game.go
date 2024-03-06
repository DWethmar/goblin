package game

import (
	"fmt"
	"log/slog"

	"github.com/dwethmar/goblin/cmd/game/actor"
	"github.com/dwethmar/goblin/pkg/es"
	"github.com/dwethmar/goblin/pkg/es/aggregate"
	eventEncoding "github.com/dwethmar/goblin/pkg/es/event/gobenc"
	eventkv "github.com/dwethmar/goblin/pkg/es/event/kv"
	kvbolt "github.com/dwethmar/goblin/pkg/kv/bbolt"

	bolt "go.etcd.io/bbolt"
)

type Options struct {
	Logger *slog.Logger
	Path   string
}

func Run(opt Options) error {
	logger := opt.Logger
	bucket := []byte("events")
	db, err := bolt.Open(opt.Path, 0600, nil)
	if err != nil {
		return fmt.Errorf("opening db: %w", err)
	}
	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	}); err != nil {
		return fmt.Errorf("creating bucket: %w", err)
	}

	eventStore := eventkv.New(kvbolt.New(bucket, db), &eventEncoding.Decoder{}, &eventEncoding.Encoder{})

	aggregateFactory := aggregate.NewFactory()

	// Register the factories
	actor.RegisterFactory(aggregateFactory)

	aggregateStore := aggregate.NewStore(eventStore, aggregateFactory)

	eventBus := es.NewEventBus()

	commandBus := es.NewCommandBus(aggregateStore, eventBus)

	if err := commandBus.Dispatch(&actor.CreateCommand{
		ActorID: "99",
		Name:    "test",
	}); err != nil {
		return fmt.Errorf("dispatching command: %w", err)
	}

	actor, err := aggregateStore.Get(actor.AggregateType, "99")
	if err != nil {
		return fmt.Errorf("getting actor: %w", err)
	}

	logger.Info("actor", "model", actor.Model, "aggregate", actor)

	return nil
}
