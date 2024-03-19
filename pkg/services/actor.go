package services

import (
	"context"

	"github.com/dwethmar/goblin/pkg/domain/actor"
	"github.com/dwethmar/goblin/pkg/es"
)

type Actors struct {
	repo       actor.Repository
	commandBus *es.CommandBus
}

func (a *Actors) Create(_ context.Context, id string, name string, x, y int) error {
	cmd := &actor.CreateCommand{
		ActorID: id,
		Name:    name,
		X:       x,
		Y:       y,
	}

	return a.commandBus.Dispatch(cmd)
}

func (a *Actors) Get(ctx context.Context, id string) (*actor.Actor, error) {
	return a.repo.Get(ctx, id)
}

func (a *Actors) List(ctx context.Context, offset, limit int) ([]*actor.Actor, error) {
	return a.repo.List(ctx, offset, limit)
}

func NewActorService(repo actor.Repository, commandBus *es.CommandBus) *Actors {
	return &Actors{
		repo:       repo,
		commandBus: commandBus,
	}
}
