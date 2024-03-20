package cmd

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/dwethmar/goblin/pkg/game"
)

// ExecInput reads commands from the reader and executes the string commands on the game.
func ExecInput(ctx context.Context, r io.Reader, s *game.Session) error {
	reader := bufio.NewReader(r)
	inputChan := make(chan string)
	errChan := make(chan error)

	go func() {
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				// Send error to main goroutine and exit
				errChan <- err
				return
			}
			// Send input to main goroutine
			inputChan <- input
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			if err == io.EOF {
				return nil
			}
			return err
		case input := <-inputChan:
			if err := s.Game.StringCommand(ctx, s, strings.Trim(input, "\n")); err != nil {
				return err
			}
		}
	}
}
