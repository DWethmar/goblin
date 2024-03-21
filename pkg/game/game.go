package game

import (
	"context"
	"log/slog"

	"github.com/dwethmar/goblin/pkg/domain/actor"
)

// Service defines an interface for any service that can process game instructions.
type ActorService interface {
	Create(ctx context.Context, aggregateID, name string, x, y int) error
	Move(ctx context.Context, aggregateID string, x, y int) error
	List(ctx context.Context, offset, limit int) ([]*actor.Actor, error)
}

type Options struct {
	Logger       *slog.Logger
	ActorService ActorService
}

type Game struct {
	*InstructionProcessor
	ActorService ActorService
}

func New(ctx context.Context, opt Options) (*Game, error) {
	return &Game{
		InstructionProcessor: &InstructionProcessor{
			Logger:       opt.Logger,
			ActorService: opt.ActorService,
		},
		ActorService: opt.ActorService,
	}, nil
}
