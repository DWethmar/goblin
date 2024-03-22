package actor

import (
	"errors"
	"fmt"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/domain"
)

func HandleCreatedEvent(a *Actor, e *aggr.Event) error {
	d, ok := e.Data.(*CreatedEventData)
	if !ok {
		return errors.New("invalid event data")
	}

	a.Name = d.Name
	a.X = d.X
	a.Y = d.Y
	a.Version = e.Version
	a.State = domain.StateCreated
	return nil
}

func HandleDestroyedEvent(a *Actor, _ *aggr.Event) error {
	a.State = domain.StateDeleted
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
