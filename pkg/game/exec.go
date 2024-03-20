package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// Session is a game session that holds the game and the current aggregate id.
type Session struct {
	Logger      *slog.Logger
	Game        *Game
	AggregateID string
}

var (
	ErrInvalidUseCommand         = errors.New("use command is invalid, expected: use <id>")
	ErrInvalidCreateCommand      = errors.New("create command is invalid, expected: create <type>")
	ErrInvalidCreateActorCommand = errors.New("create actor command is invalid, expected: create actor <name> <x> <y>")
	ErrInvalidMoveCommand        = errors.New("move command is invalid, expected: move <x> <y>")
)

func useCommand(ctx context.Context, s *Session, args []string) error {
	if len(args) < 1 {
		return ErrInvalidUseCommand
	}
	s.AggregateID = args[0]
	s.Logger.Info("Using aggregate", "aggregate", s.AggregateID) // Fix: Add missing final value argument
	return nil
}

func createCommand(ctx context.Context, s *Session, args []string) error {
	if len(args) < 1 {
		return ErrInvalidCreateCommand
	}

	switch subject := args[0]; subject {
	case "actor":
		return createActor(ctx, s, args)
	default:
		return fmt.Errorf("create command is invalid, unknown type: %s", subject)
	}
}

func createActor(ctx context.Context, s *Session, args []string) error {
	if len(args) < 4 {
		return ErrInvalidCreateActorCommand
	}
	id, name := s.AggregateID, args[1]
	x, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("create actor command is invalid, x is not a number: %s", args[2])
	}

	y, err := strconv.Atoi(args[3])
	if err != nil {
		return fmt.Errorf("create actor command is invalid, y is not a number: %s", args[3])
	}

	return s.Game.actorService.Create(ctx, id, name, x, y)
}

func moveCommand(ctx context.Context, s *Session, args []string) error {
	if len(args) < 2 {
		return ErrInvalidMoveCommand
	}

	x, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("move command is invalid, x is not a number: %s", args[0])
	}

	y, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("move command is invalid, y is not a number: %s", args[1])
	}

	return s.Game.actorService.Move(ctx, s.AggregateID, x, y)
}

var cmds = map[string]func(ctx context.Context, s *Session, args []string) error{
	"use":    useCommand,
	"create": createCommand,
	"move":   moveCommand,
}

func (g *Game) StringCommand(ctx context.Context, s *Session, str string) error {
	if str == "" {
		s.Logger.Debug("Empty command")
		return nil
	}

	args := strings.Fields(str)
	if len(args) < 1 {
		return errors.New("command is invalid")
	}

	cmd, args := strings.ToLower(args[0]), args[1:]
	if f, ok := cmds[cmd]; ok {
		return f(ctx, s, args)
	}

	return fmt.Errorf("unknown command: %s", cmd)
}
