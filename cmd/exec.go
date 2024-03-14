/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
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
	pipeFile    *string
)

// execCmd represents the run command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "exec a command to the game",
	RunE: func(cmd *cobra.Command, args []string) error {
		logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		logger := slog.New(logHandler)

		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

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
			return fmt.Errorf("creating game: %w", err)
		}

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt)
		go func() {
			<-signalCh
			cancel()
			logger.Info("Shutting down")
		}()

		go func() {
			if err := pipeGameCmds(ctx, g); err != nil {
				logger.Error("pipeGameCmds", "err", err)
			}
		}()

		<-ctx.Done()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	aggregateID = execCmd.Flags().StringP("aggregate-id", "a", "", "The aggregate id to use")
	pipeFile = execCmd.Flags().StringP("pipe-file", "p", "goblinput", "The pipe file to use")
}
