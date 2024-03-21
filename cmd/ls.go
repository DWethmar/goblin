/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	header bool // show header
	limit  int  // max number of aggregates to show
	offset int  // offset of aggregates to show
)

// lsActorsCmd represents the ls actors command
var lsActorsCmd = &cobra.Command{
	Use:   "actors",
	Short: "list actors",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
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
		defer close()

		actorService := g.ActorService
		actors, err := actorService.List(ctx, offset, limit)
		if err != nil {
			return fmt.Errorf("listing actors: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		if header {
			fmt.Fprintln(w, "ID\tNAME\tX\tY")
		}
		for _, a := range actors {
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\n", a.ID, a.Name, a.X, a.Y)
		}

		return w.Flush()
	},
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list aggregates",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.PersistentFlags().BoolVarP(&header, "header", "H", false, "show header")
	lsCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 10, "max number of aggregates to show")
	lsCmd.PersistentFlags().IntVarP(&offset, "offset", "o", 0, "offset of aggregates to show")
	lsCmd.AddCommand(lsActorsCmd)
}
