package game

import (
	"context"
	"log/slog"

	"github.com/dwethmar/goblin/pkg/services"
)

type Options struct {
	Logger       *slog.Logger
	ActorService *services.Actors
}

type Game struct {
	logger       *slog.Logger
	actorService *services.Actors
}

// ActorService returns the actor service
func (g *Game) ActorService() *services.Actors { return g.actorService }

func New(ctx context.Context, opt Options) (*Game, error) {
	return &Game{
		logger:       opt.Logger,
		actorService: opt.ActorService,
	}, nil
}
