package cmd

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/dwethmar/goblin/pkg/game"
)

func ExecLines(ctx context.Context, r io.Reader, g *game.Game, s *game.State) error {
	reader := bufio.NewReader(r)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
		if err := g.ExecStringCommand(ctx, s, strings.Trim(input, "\n")); err != nil {
			return err
		}
	}
}
