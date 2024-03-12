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
	ActorService *services.Actors
}

func New(opt Options) *Game {
	logger := opt.Logger

	g := &Game{
		logger:       logger,
		ActorService: opt.ActorService,
	}

	if err := g.ActorService.CreateActor("1", "test"); err != nil {
		logger.Error("creating actor", "err", err)
	}

	a, err := g.ActorService.GetActor(context.Background(), "1")
	if err != nil {
		logger.Error("get actor", "err", err)
	}

	logger.Info("actor", "actor", a)

	return g
}
