package actor

import (
	"context"
	"errors"
	"fmt"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/domain"
)

// ActorEventsMatcher is a matcher that can be used to match all actor events.
var ActorEventsMatcher = aggr.MatchEvents{
	CreatedEventType,
	DestroyedEventType,
	MovedEventType,
}

// ActorSinkHandler returns a handler that can be used to handle events from
// the event store and update the actor repository.
func ActorSinkHandler(repo Repository) aggr.EventHandlerFunc {
	return aggr.EventHandlerFunc(func(ctx context.Context, event *aggr.Event) error {
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

		if err := a.HandleEvent(ctx, event); err != nil {
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

func HandleCreatedEvent(a *Actor, e *aggr.Event) error {
	d, ok := e.Data.(*CreatedEventData)
	if !ok {
		return errors.New("invalid event data")
	}

	a.Name = d.Name
	a.X = d.X
	a.Y = d.Y
	a.Version = e.Version
	a.state = domain.StateCreated
	return nil
}

func HandleDestroyedEvent(a *Actor, _ *aggr.Event) error {
	a.state = domain.StateDeleted
	return nil
}

func HandleMovedEvent(a *Actor, e *aggr.Event) error {
	d, ok := e.Data.(*MovedEventData)
	if !ok {
		return fmt.Errorf("expected *MovedEventData, got %T", e.Data)
	}
	a.X = d.X
	a.Y = d.Y
	return nil
}
