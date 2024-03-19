/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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

		g, close, err := SetupGame(ctx, Config{
			Logger: logger,
			Game:   Game,
		})
		if err != nil {
			return fmt.Errorf("setting up game: %w", err)
		}
		defer close()

		actorService := g.ActorService()
		actors, err := actorService.List(ctx, 0, 10)
		if err != nil {
			return fmt.Errorf("listing actors: %w", err)
		}

		if len(actors) == 0 {
			fmt.Println("no actors found")
			return nil
		}

		for _, a := range actors {
			fmt.Printf("%s: %s\n", a.ID, a.Name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
