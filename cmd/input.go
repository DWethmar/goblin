package cmd

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/dwethmar/goblin/pkg/game"
)

func ExecInput(ctx context.Context, r io.Reader, g *game.Game, s *game.CmdContext) error {
	reader := bufio.NewReader(r)
	for {
		if ctx.Err() != nil {
			return nil
		}
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
		if err := g.StringCommand(ctx, s, strings.Trim(input, "\n")); err != nil {
			return err
		}
	}
}
