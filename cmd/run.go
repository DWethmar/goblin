/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"

	"github.com/dwethmar/goblin/cmd/game"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a command on the goblin game",
	Long:  `A longer description`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := slog.Default()
		logger.Info("run command", "args", args)

		if err := game.Run(game.Options{
			Logger: logger,
			Path:   ".tmp/goblin.db",
		}); err != nil {
			return fmt.Errorf("running game: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
