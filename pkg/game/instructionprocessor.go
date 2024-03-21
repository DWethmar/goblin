package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

var (
	ErrInvalidUseInstruction         = errors.New("use instruction is invalid, expected: use <id>")
	ErrInvalidCreateInstruction      = errors.New("create instruction is invalid, expected: create <type>")
	ErrInvalidCreateActorInstruction = errors.New("create actor instruction is invalid, expected: create actor <name> <x> <y>")
	ErrInvalidMoveInstruction        = errors.New("move instruction is invalid, expected: move <x> <y>")
)

// InstructionProcessor is a service that processes game instructions.
type InstructionProcessor struct {
	Logger       *slog.Logger
	AggregateID  string
	ActorService ActorService
}

func (p *InstructionProcessor) use(instruction string) error {
	var id string
	if _, err := fmt.Sscanf(instruction, "use %s", &id); err != nil || id == "" {
		p.Logger.Error("Error scanning use instruction", "error", err)
		return fmt.Errorf("%w: %s", ErrInvalidUseInstruction, instruction)
	}
	p.AggregateID = id
	p.Logger.Info("Using aggregate", "aggregate", p.AggregateID)
	return nil
}

func (p *InstructionProcessor) create(ctx context.Context, in string) error {
	var subject string
	if _, err := fmt.Sscanf(in, "create %s", &subject); err != nil {
		p.Logger.Error("Error scanning create instruction", "error", err)
		return fmt.Errorf("%w: %s", ErrInvalidCreateInstruction, in)
	}

	switch subject {
	case "actor":
		return p.createActor(ctx, in)
	default:
		return fmt.Errorf("create instruction is invalid, unknown type: %s", subject)
	}
}

func (p *InstructionProcessor) createActor(ctx context.Context, in string) error {
	var name string
	var x, y int
	if _, err := fmt.Sscanf(in, "create actor %s %d %d", &name, &x, &y); err != nil {
		p.Logger.Error("Error scanning create actor instruction", "error", err)
		return fmt.Errorf("%w: %s", ErrInvalidCreateActorInstruction, in)
	}
	return p.ActorService.Create(ctx, p.AggregateID, name, x, y)
}

func (p *InstructionProcessor) move(ctx context.Context, in string) error {
	var x, y int
	if _, err := fmt.Sscanf(in, "move %d %d", &x, &y); err != nil {
		p.Logger.Error("Error scanning move instruction", "error", err)
		return fmt.Errorf("%w: %s", ErrInvalidMoveInstruction, in)
	}
	return p.ActorService.Move(ctx, p.AggregateID, x, y)
}

// Process interprets and processes a given instruction for the game.
func (p *InstructionProcessor) Process(ctx context.Context, in string) error {
	if in == "" {
		p.Logger.Debug("Empty instruction")
		return nil
	}

	cmd := strings.Fields(in)[0]
	switch cmd {
	case "use":
		return p.use(in)
	case "create":
		return p.create(ctx, in)
	case "move":
		return p.move(ctx, in)
	default:
		return fmt.Errorf("unknown instruction: %s", cmd)
	}
}
