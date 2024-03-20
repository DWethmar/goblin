package chunk

import (
	"context"
	"errors"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/domain"
	"github.com/dwethmar/goblin/pkg/matrix"
)

// MatchAllEvents is a matcher that can be used to match all chunk events.
var MatchAllEvents = aggr.MatchEvents{
	CreatedEventType,
}

func ChunkSinkHandler(repo Repository) aggr.EventHandlerFunc {
	return aggr.EventHandlerFunc(func(ctx context.Context, event *aggr.Event) error {
		c, err := repo.Get(ctx, event.AggregateID)
		if err != nil {
			if errors.Is(err, ErrNotFound) { // chunk not found, create it
				c = &Chunk{
					ID: event.AggregateID,
				}

				if c, err = repo.Create(ctx, c); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		if err := c.HandleEvent(ctx, event); err != nil {
			return err
		}

		if c.Deleted() {
			if err := repo.Delete(ctx, c.ID); err != nil {
				return err
			}
		} else {
			if _, err := repo.Update(ctx, c); err != nil {
				return err
			}
		}

		return nil
	})
}

func HandleCreatedEvent(c *Chunk, e *aggr.Event) error {
	d, ok := e.Data.(*CreatedEventData)
	if !ok {
		return errors.New("invalid event data")
	}

	c.X = d.X
	c.Y = d.Y
	c.Width = d.Width
	c.Height = d.Height
	c.Tiles = matrix.New(d.Width, d.Height, 0)
	for _, t := range d.Tiles {
		c.Tiles.Set(t.X, t.Y, t.Value)
	}
	c.Version = e.Version
	c.state = domain.StateCreated
	return nil
}
