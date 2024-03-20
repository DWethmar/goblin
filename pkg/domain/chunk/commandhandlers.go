package chunk

import (
	"errors"
	"time"

	"github.com/dwethmar/goblin/pkg/aggr"
)

func CreateChunkCommandHandler(c *Chunk, cmd *CreateCommand) (*aggr.Event, error) {
	if cmd.ChunkID == "" {
		return nil, errors.New("chunk id can't be empty")
	}

	return &aggr.Event{
		AggregateID: cmd.ChunkID,
		Type:        CreatedEventType,
		Data: &CreatedEventData{
			X: cmd.X,
			Y: cmd.Y,
		},
		Version:   c.Version + 1,
		CreatedAt: time.Now(),
	}, nil
}
