package actor

import (
	"context"
	"testing"

	"github.com/dwethmar/goblin/pkg/es"
)

func TestActorSinkHandler(t *testing.T) {
	t.Run("ActorSinkHandler", func(t *testing.T) {
		ctx := context.Background()
		actor := &Actor{
			ID: "1",
		}
		repo := &RepositoryMock{
			GetFunc: func(ctx context.Context, id string) (*Actor, error) {
				return actor, nil
			},
			UpdateFunc: func(ctx context.Context, a *Actor) (*Actor, error) {
				return a, nil
			},
		}

		event := &es.Event{
			AggregateID: "1",
			Type:        CreatedEventType,
			Data:        &CreatedEventData{Name: "test"},
		}
		handler := ActorSinkHandler(ctx, repo)
		if err := handler(event); err != nil {
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