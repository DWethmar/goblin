package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// CmdContext is a context for commands.
type CmdContext struct {
	Logger      *slog.Logger
	AggregateID string
}

var (
	ErrInvalidUseCommand         = errors.New("use command is invalid, expected: use <id>")
	ErrInvalidCreateCommand      = errors.New("create command is invalid, expected: create <type>")
	ErrInvalidCreateActorCommand = errors.New("create actor command is invalid, expected: create actor <name> <x> <y>")
	ErrInvalidMoveCommand        = errors.New("move command is invalid, expected: move <x> <y>")
)

var cmds = map[string]func(ctx context.Context, g *Game, s *CmdContext, args []string) error{
	"use": func(ctx context.Context, g *Game, s *CmdContext, args []string) error {
		if len(args) < 1 {
			return ErrInvalidUseCommand
		}
		s.AggregateID = args[0]
		fmt.Printf("using aggregate: %s\n", s.AggregateID)
		return nil
	},
	"create": func(ctx context.Context, g *Game, s *CmdContext, args []string) error {
		if len(args) < 1 {
			return ErrInvalidCreateCommand
		}

		subject := args[0]
		switch subject {
		case "actor":
			if len(args) < 4 {
				return ErrInvalidCreateActorCommand
			}
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
	"move": func(ctx context.Context, g *Game, s *CmdContext, args []string) error {
		if len(args) < 2 {
			return ErrInvalidMoveCommand
		}

		x, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("move command is invalid, x is not a number: %s", args[1])
		}

		y, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("move command is invalid, y is not a number: %s", args[2])
		}

		return g.actorService.Move(ctx, s.AggregateID, x, y)
	},
}

func (g *Game) StringCommand(ctx context.Context, s *CmdContext, str string) error {
	if str == "" {
		g.logger.Debug("empty command")
		return nil
	}

	args := strings.Split(str, " ")
	if len(args) < 1 {
		return fmt.Errorf("command is invalid")
	}

	cmd := args[0]
	if f, ok := cmds[cmd]; ok {
		return f(ctx, g, s, args[1:])
	}

	return fmt.Errorf("unknown command: %s", cmd)
}
