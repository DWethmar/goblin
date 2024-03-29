/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dwethmar/goblin/pkg/game"
	"github.com/spf13/cobra"
)

var (
	// Game is the game name
	Game string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goblin",
	Short: "Goblin is a game",
	Long:  `Goblin is a game`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		gameName, _ := cmd.Flags().GetString("game")
		if err := game.ValidateName(gameName); err != nil {
			return fmt.Errorf("validating game name: %w", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&Game, "game", "g", "goblin", "The game name")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
