package chunk

import (
	"errors"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/domain"
	"github.com/dwethmar/goblin/pkg/matrix"
)

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
