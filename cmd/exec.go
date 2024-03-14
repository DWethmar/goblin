/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"strings"
	"syscall"

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
			return fmt.Errorf("creating game: %w", err)
		}

		r := os.Stdin

		if p := *pipeFile; p != "" {
			os.Remove(p)
			err := syscall.Mkfifo(p, 0666)
			if err != nil {
				log.Fatal("Make named pipe file error:", err)
			}
			file, err := os.OpenFile(p, os.O_CREATE, os.ModeNamedPipe)
			if err != nil {
				log.Fatal("Open named pipe file error:", err)
			}
			r = file
			defer file.Close()
			defer os.Remove(p)
		} else {
			fmt.Println("No pipe file specified, using stdin")
		}

		return ReadEvent(r, func(i string) error {
			args := strings.Split(i, " ")
			if len(args) < 2 {
				logger.Error("invalid command", "command", i)
			}

			if err := g.DispatchStringCommand(cmd.Context(), *aggregateID, args[0], args[1:]...); err != nil {
				logger.Error("dispatching command", "error", err)
			}

			return nil
		})
	},
}

func ReadEvent(r io.Reader, f func(i string) error) error {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				continue
			}

			return fmt.Errorf("reading line: %w", err)
		}

		if err := f(string(line)); err != nil {
			return fmt.Errorf("processing line: %w", err)
		}
	}
}

func init() {
	rootCmd.AddCommand(execCmd)
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	aggregateID = execCmd.Flags().StringP("aggregate-id", "a", "", "The aggregate id to use")
	pipeFile = execCmd.Flags().StringP("pipe-file", "p", "goblinput", "The pipe file to use")
}
