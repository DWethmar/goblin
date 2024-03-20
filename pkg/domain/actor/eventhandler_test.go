package actor

import (
	"context"
	"testing"
	"time"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/google/go-cmp/cmp"
)

func TestActorSinkHandler(t *testing.T) {
	t.Run("ActorSinkHandler", func(t *testing.T) {
		ctx := context.TODO()
		actor := &Actor{
			ID: "1",
		}
		repo := &MockRepository{
			GetFunc: func(ctx context.Context, id string) (*Actor, error) {
				if ctx != context.TODO() {
					t.Errorf("ActorSinkHandler() ctx = %v, want %v", ctx, context.TODO())
				}
				if id != "1" {
					t.Errorf("ActorSinkHandler() id = %v, want 1", id)
				}
				return actor, nil
			},
			UpdateFunc: func(ctx context.Context, a *Actor) (*Actor, error) {
				if ctx != context.TODO() {
					t.Errorf("ActorSinkHandler() ctx = %v, want %v", ctx, context.TODO())
				}
				if diff := cmp.Diff(actor, a, cmp.AllowUnexported(Actor{})); diff != "" {
					t.Errorf("ActorSinkHandler() UpdateFunc() mismatch (-want +got):\n%s", diff)
				}
				return a, nil
			},
		}

		event := &aggr.Event{
			AggregateID: "1",
			Type:        CreatedEventType,
			Data:        &CreatedEventData{Name: "test"},
			Timestamp:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		handler := ActorSinkHandler(repo)
		if err := handler(ctx, event); err != nil {
			t.Errorf("ActorSinkHandler() error = %v", err)
			return
		}
		actor, err := repo.Get(ctx, "1")
		if err != nil {
			t.Errorf("ActorSinkHandler() error = %v", err)
			return
		}
		if actor.ID != "1" {
			t.Errorf("ActorSinkHandler() actor.ID = %v, want 1", actor.ID)
		}
	})
}
