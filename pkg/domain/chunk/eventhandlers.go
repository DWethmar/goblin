package chunk

import (
	"errors"

	"github.com/dwethmar/goblin/pkg/aggr"
)

func HandleCreatedEvent(c *Chunk, e *aggr.Event) error {
	d, ok := e.Data.(*CreatedEventData)
	if !ok {
		return errors.New("invalid event data")
	}

	c.X = d.X
	c.Y = d.Y
	c.Version = e.Version
	return nil
}
