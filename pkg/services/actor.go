package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/clock"
	"github.com/dwethmar/goblin/pkg/domain/actor"
)

type Actors struct {
	clock       *clock.Clock
	actorReader actor.Reader
	commandBus  aggr.CommandBus
}

func (a *Actors) Create(ctx context.Context, id string, name string, x, y int) error {
	//check if actor already exists
	if r, err := a.actorReader.Get(ctx, id); err == nil {
		if r != nil {
			return fmt.Errorf("actor with id %s already exists", id)
		}

		return errors.New("error getting actor")
	}

	if err := a.commandBus.HandleCommand(ctx, &actor.CreateCommand{
		ActorID:   id,
		Name:      name,
		X:         x,
		Y:         y,
		Timestamp: a.clock.Now(),
	}); err != nil {
		return fmt.Errorf("error dispatching create actor command: %w", err)
	}
	return nil
}

func (a *Actors) Move(ctx context.Context, id string, x, y int) error {
	if err := a.commandBus.HandleCommand(ctx, &actor.MoveCommand{
		ActorID:   id,
		X:         x,
		Y:         y,
		Timestamp: a.clock.Now(),
	}); err != nil {
		return fmt.Errorf("error dispatching move actor command: %w", err)
	}
	return nil
}

func (a *Actors) Get(ctx context.Context, id string) (*actor.Actor, error) {
	return a.actorReader.Get(ctx, id)
}

func (a *Actors) List(ctx context.Context, limit, offset int) ([]*actor.Actor, error) {
	return a.actorReader.List(ctx, limit, offset)
}

func NewActorService(repo actor.Repository, commandBus aggr.CommandBus) *Actors {
	return &Actors{
		clock:       clock.New(time.UTC),
		actorReader: repo,
		commandBus:  commandBus,
	}
}
