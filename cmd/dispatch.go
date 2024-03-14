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
	eventEncoding "github.com/dwethmar/goblin/pkg/es/event/encoding"
	eventkv "github.com/dwethmar/goblin/pkg/es/event/kv"
	"github.com/dwethmar/goblin/pkg/kv/bbolt"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var (
	aggregateID *string
)

// dispatchCmd represents the run command
var dispatchCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Dispatch a command",
	Long:  `A longer description`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
		logger := slog.New(logHandler)

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

		bucket := []byte("events")
		eventStore := eventkv.New(bbolt.New(bucket, db), &eventEncoding.Decoder{}, &eventEncoding.Encoder{})

		g, err := game.New(cmd.Context(), game.Options{
			Logger:     logger,
			EventStore: eventStore,
		})
		if err != nil {
			return err
		}

		if err := g.DispatchStringCommand(cmd.Context(), *aggregateID, args[0], args[1:]...); err != nil {
			return fmt.Errorf("could not dispatch command: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dispatchCmd)
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	aggregateID = dispatchCmd.Flags().StringP("aggregate-id", "a", "", "The aggregate id to use")
}
