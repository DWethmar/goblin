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
	"syscall"

	"github.com/dwethmar/goblin/cmd/game"
	eventEncoding "github.com/dwethmar/goblin/pkg/es/event/encoding"
	eventkv "github.com/dwethmar/goblin/pkg/es/event/kv"
	"github.com/dwethmar/goblin/pkg/kv/bbolt"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var (
	pipeFile *string
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

		g, err := game.New(ctx, game.Options{
			Logger:     logger,
			EventStore: eventStore,
		})
		if err != nil {
			return fmt.Errorf("creating game: %w", err)
		}

		done := make(chan struct{}, 1)
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() { // handle signals and begin shutdown
			<-sigs
			done <- struct{}{}
		}()

		go func() {
			<-done
			cancel()
		}()

		go func() {
			if err := pipeGameCmds(ctx, g); err != nil {
				logger.Error("pipeGameCmds", "err", err)
			}
		}()

		<-done
		return nil
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	pipeFile = execCmd.Flags().StringP("pipe-file", "p", "goblinput", "The pipe file to use")
}
