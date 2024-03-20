package cmd

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/domain/actor"
	"github.com/dwethmar/goblin/pkg/game"
	"github.com/dwethmar/goblin/pkg/services"
)

func TestExecInput(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		r := io.Reader(strings.NewReader(strings.Join([]string{
			"use 1",
			"create test 1 1",
		}, "\n")))

		aggregateStore := &aggr.MockAggregateStore{
			GetFunc: func(ctx context.Context, _, id string) (*aggr.Aggregate, error) {
				return &aggr.Aggregate{
					Model: &actor.Actor{
						ID: id,
					},
				}, nil
			},
			SaveFunc: func(_ context.Context, _ ...*aggr.Aggregate) error { return nil },
		}

		commandBus := aggr.NewCommandBus(aggregateStore, aggr.NewEventBus())

		g, err := game.New(ctx, game.Options{
			Logger:       slog.Default(),
			ActorService: services.NewActorService(nil, commandBus),
		})
		if err != nil {
			t.Errorf("New() error = %v, want nil", err)
			return
		}

		s := &game.Session{
			Logger:      slog.Default(),
			Game:        g,
			AggregateID: "1",
		}

		if err := ExecInput(ctx, r, s); err != nil {
			t.Errorf("ExecInput() error = %v, want nil", err)
		}
	})

	t.Run("should return error if command fails", func(t *testing.T) {
		ctx := context.Background()
		r := io.Reader(strings.NewReader(strings.Join([]string{
			"move 0 0\n",
		}, "")))

		aggregateStore := &aggr.MockAggregateStore{
			GetFunc: func(ctx context.Context, _, id string) (*aggr.Aggregate, error) {
				return nil, errors.ErrUnsupported
			},
		}

		commandBus := aggr.NewCommandBus(aggregateStore, aggr.NewEventBus())

		g, err := game.New(ctx, game.Options{
			Logger:       slog.Default(),
			ActorService: services.NewActorService(nil, commandBus),
		})
		if err != nil {
			t.Errorf("New() error = %v, want nil", err)
			return
		}

		s := &game.Session{
			Logger:      slog.Default(),
			Game:        g,
			AggregateID: "1",
		}

		if err := ExecInput(ctx, r, s); !errors.Is(err, errors.ErrUnsupported) {
			t.Errorf("ExecInput() error = %v, want %v", err, errors.ErrUnsupported)
		}
	})
}
