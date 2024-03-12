/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/dwethmar/goblin/cmd/game"
	"github.com/dwethmar/goblin/pkg/domain/actor"
	actorMemory "github.com/dwethmar/goblin/pkg/domain/actor/memory"
	"github.com/dwethmar/goblin/pkg/es"
	"github.com/dwethmar/goblin/pkg/es/aggregate"
	eventEncoding "github.com/dwethmar/goblin/pkg/es/event/gobenc"
	eventkv "github.com/dwethmar/goblin/pkg/es/event/kv"
	"github.com/dwethmar/goblin/pkg/kv/bbolt"
	"github.com/dwethmar/goblin/pkg/services"
	"github.com/spf13/cobra"

	bolt "go.etcd.io/bbolt"
)

// dispatchCmd represents the run command
var dispatchCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Dispatch a command",
	Long:  `A longer description`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := slog.Default()
		logger.Info("run command", "args", args)

		dirName := "./.tmp"
		if _, err := os.Stat(dirName); os.IsNotExist(err) {
			if err := os.Mkdir(dirName, 0700); err != nil {
				return fmt.Errorf("creating dir: %w", err)
			}
		}

		dnName := "goblin.db"
		db, err := bolt.Open(path.Join(dirName, dnName), 0600, nil)
		if err != nil {
			return fmt.Errorf("opening db: %w", err)
		}
		defer db.Close()

		// Create the bucket if it does not exist
		bucket := []byte("events")
		if err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucket)
			return err
		}); err != nil {
			return fmt.Errorf("creating bucket: %w", err)
		}

		actorRepo := actorMemory.NewRepository()
		eventStore := eventkv.New(bbolt.New(bucket, db), &eventEncoding.Decoder{}, &eventEncoding.Encoder{})

		// Create the agregate factory and register agregates
		aggregateFactory := aggregate.NewFactory()
		actor.RegisterFactory(aggregateFactory)

		// Create the event bus and add event handlers
		eventBus := es.NewEventBus()
		eventBus.Subscribe(actor.ActorEventMatcher, actor.ActorSinkHandler(ctx, actorRepo))

		aggregateStore := aggregate.NewStore(eventStore, aggregateFactory)
		commandBus := es.NewCommandBus(aggregateStore, eventBus)

		game.New(game.Options{
			Logger:       logger,
			ActorService: services.NewActorService(actorRepo, commandBus),
		})

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dispatchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
