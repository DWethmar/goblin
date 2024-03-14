package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"syscall"

	"github.com/dwethmar/goblin/cmd/game"
)

func readEvent(ctx context.Context, r io.Reader, errCh chan<- error) chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		reader := bufio.NewReader(r)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						continue
					}
					errCh <- fmt.Errorf("reading from input: %w", err)
					return
				}
				ch <- strings.TrimSpace(line)
			}
		}
	}()
	return ch
}

func pipeGameCmds(ctx context.Context, g *game.Game) error {
	r := os.Stdin
	if p := *pipeFile; p != "" {
		os.Remove(p)
		err := syscall.Mkfifo(p, 0666)
		if err != nil {
			return fmt.Errorf("creating named pipe: %w", err)
		}
		file, err := os.OpenFile(p, os.O_CREATE, os.ModeNamedPipe)
		if err != nil {
			return fmt.Errorf("opening named pipe: %w", err)
		}
		r = file
		defer file.Close()
		defer os.Remove(p)
	}

	errCh := make(chan error)
	defer close(errCh)
	eventCh := readEvent(ctx, r, errCh)

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			return err
		case event := <-eventCh:
			if err := g.ExecStringCommand(ctx, event); err != nil {
				slog.ErrorContext(ctx, "executing command", "err", err)
			}
		}
	}
}
