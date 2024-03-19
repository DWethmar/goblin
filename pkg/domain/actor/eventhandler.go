package actor

import (
	"context"
	"errors"
	"fmt"

	"github.com/dwethmar/goblin/pkg/aggr"
)

// ActorEventMatcher is a matcher that can be used to match events for the actor.
var ActorEventMatcher = aggr.MatchEvents{
	CreatedEventType,
	DestroyedEventType,
	MovedEventType,
}

// ActorSinkHandler returns a handler that can be used to handle events from
// the event store and update the actor repository.
func ActorSinkHandler(ctx context.Context, repo Repository) aggr.EventHandlerFunc {
	return aggr.EventHandlerFunc(func(event *aggr.Event) error {
		a, err := repo.Get(ctx, event.AggregateID)
		if err != nil {
			if errors.Is(err, ErrNotFound) { // actor not found, create it
				a = &Actor{
					ID: event.AggregateID,
				}

				if _, err = repo.Create(ctx, a); err != nil {
					return fmt.Errorf("create actor: %w", err)
				}
			} else {
				return fmt.Errorf("get actor: %w", err)
			}
		}

		if err := a.HandleEvent(event); err != nil {
			return fmt.Errorf("apply event: %w", err)
		}

		if a.Deleted() {
			if err := repo.Delete(ctx, a.ID); err != nil {
				return fmt.Errorf("delete actor: %w", err)
			}
		} else {
			if _, err := repo.Update(ctx, a); err != nil {
				return fmt.Errorf("update actor: %w", err)
			}
		}

		return nil
	})
}
