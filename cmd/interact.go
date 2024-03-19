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
		defer close()

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
			s := &game.CmdContext{
				Logger:      logger,
				AggregateID: aggregateId,
			}

			var r io.Reader
			if filePath != "" {
				f, err := os.Open(filePath)
				if err != nil {
					logger.Error("opening file", "err", err)
					done <- struct{}{}
					return
				}
				defer f.Close()
				r = f
			} else {
				r = os.Stdin
				fmt.Print("Enter cmd: ")
			}

			if err := ExecLines(ctx, r, g, s); err != nil {
				logger.Error("exec lines", "err", err)
			}

			done <- struct{}{}
		}()

		<-done
		return nil
	},
	ValidArgs: []string{"aggregate-id"},
}

func init() {
	rootCmd.AddCommand(interactCmd)
	interactCmd.PersistentFlags().StringVarP(&aggregateId, "aggregate", "a", "", "The aggregate id")
	interactCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "The file to read commands from")
}
