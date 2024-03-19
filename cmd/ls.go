/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/dwethmar/goblin/cmd/game"
	eventEncoding "github.com/dwethmar/goblin/pkg/es/event/encoding"
	eventkv "github.com/dwethmar/goblin/pkg/es/event/kv"
	"github.com/dwethmar/goblin/pkg/kv/bbolt"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list aggregates",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		logger := slog.New(logHandler)
		logger.Info("exec", "game", Game)

		if Game == "" {
			return fmt.Errorf("game is required")
		}

		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		gamePath := fmt.Sprintf("./.tmp/%s.db", Game)

		db, err := bbolt.Connect(gamePath)
		if err != nil {
			return fmt.Errorf("connecting to db: %w", err)
		}
		defer db.Close()

		bucket := []byte("events")
		eventStore := eventkv.New(bbolt.New(bucket, db), &eventEncoding.Decoder{}, &eventEncoding.Encoder{})

		_, err = game.New(ctx, game.Options{
			Logger:     logger,
			EventStore: eventStore,
		})
		if err != nil {
			return fmt.Errorf("creating game: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
