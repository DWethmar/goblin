package game

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type State struct {
	Logger      *slog.Logger
	AggregateID string
}

var cmds = map[string]func(ctx context.Context, g *Game, s *State, args []string) error{
	"use": func(ctx context.Context, g *Game, s *State, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("use command is invalid, expected: use <id>")
		}
		s.AggregateID = args[0]
		fmt.Printf("using aggregate: %s\n", s.AggregateID)
		return nil
	},
	"create": func(ctx context.Context, g *Game, s *State, args []string) error {
		if len(args) < 3 {
			return fmt.Errorf("create command is invalid, expected: create <name> <x> <y>")
		}

		switch args[0] {
		case "actor":
			var (
				id   = s.AggregateID
				name = args[1]
			)
			x, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("create command is invalid, x is not a number: %s", args[2])
			}

			y, err := strconv.Atoi(args[3])
			if err != nil {
				return fmt.Errorf("create command is invalid, y is not a number: %s", args[3])
			}

			return g.actorService.Create(ctx, id, name, x, y)
		}

		return fmt.Errorf("create command is invalid, unknown type: %s", args[0])
	},
}

func (g *Game) ExecStringCommand(ctx context.Context, s *State, cmdStr string) error {
	args := strings.Split(cmdStr, " ")
	if len(args) < 1 {
		return fmt.Errorf("command is invalid")
	}

	cmd := args[0]
	if f, ok := cmds[cmd]; ok {
		return f(ctx, g, s, args[1:])
	}

	return fmt.Errorf("unknown command: %s", cmd)
}
