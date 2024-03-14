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

func (a *Actors) Create(id string, name string) error {
	cmd := &actor.CreateCommand{
		ActorID: id,
		Name:    name,
	}

	return a.commandBus.Dispatch(cmd)
}

func (a *Actors) Get(ctx context.Context, id string) (*actor.Actor, error) {
	return a.repo.Get(ctx, id)
}

func NewActorService(repo actor.Repository, commandBus *es.CommandBus) *Actors {
	return &Actors{
		repo:       repo,
		commandBus: commandBus,
	}
}
