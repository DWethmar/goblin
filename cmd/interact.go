/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwethmar/goblin/pkg/game"
	"github.com/spf13/cobra"
)

var (
	aggregateId string
	filePath    string
)

// interactCmd represents the run command
var interactCmd = &cobra.Command{
	Use:   "interact",
	Short: "interact with a game",
	RunE: func(cmd *cobra.Command, args []string) error {
		if Game == "" {
			return fmt.Errorf("game is required")
		}

		logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		logger := slog.New(logHandler)
		logger.Info("exec", "game", Game)

		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		g, close, err := SetupGame(ctx, Config{
			Logger: logger,
			Game:   Game,
		})
		if err != nil {
			return fmt.Errorf("setting up game: %w", err)
		}
		defer close() // Ensure game resources are always cleaned up.

		done := make(chan struct{}, 1)
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigs
			cancel() // Trigger context cancellation.
		}()

		go func() {
			defer func() { done <- struct{}{} }() // Signal completion on exit.

			var r io.Reader
			if filePath != "" {
				f, err := os.Open(filePath)
				if err != nil {
					logger.Error("opening file", "err", err)
					return
				}
				defer f.Close()
				r = f
			} else {
				r = os.Stdin
			}

			s := &game.Session{
				Logger:      logger,
				Game:        g,
				AggregateID: aggregateId,
			}

			if err := ExecInput(ctx, r, s); err != nil {
				if err != context.Canceled {
					logger.Error("executing input", "err", err)
				}
				return
			}
		}()

		<-done     // Wait for processing to complete or be cancelled.
		return nil // At this point, defer functions will run, ensuring cleanup.
	},
	ValidArgs: []string{"aggregate-id"},
}

func init() {
	rootCmd.AddCommand(interactCmd)
	interactCmd.PersistentFlags().StringVarP(&aggregateId, "aggregate", "a", "", "The aggregate id")
	interactCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "The file to read commands from")
}
