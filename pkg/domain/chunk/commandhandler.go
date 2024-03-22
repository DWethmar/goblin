package chunk

import (
	"errors"
	"fmt"
	"time"

	"github.com/dwethmar/goblin/pkg/aggr"
)

func CreateChunkCommandHandler(c *Chunk, cmd *CreateCommand) (*aggr.Event, error) {
	if cmd.ChunkID == "" {
		return nil, errors.New("chunk id can't be empty")
	}

	// Check if any tiles are out of bounds
	for _, t := range cmd.Tiles {
		if t.X < 0 || t.X >= cmd.Width || t.Y < 0 || t.Y >= cmd.Height {
			return nil, fmt.Errorf("tile out of bounds: %v", t)
		}
	}

	return &aggr.Event{
		AggregateID: cmd.ChunkID,
		EventType:   CreatedEventType,
		Data: &CreatedEventData{
			X:      cmd.X,
			Y:      cmd.Y,
			Width:  cmd.Width,
			Height: cmd.Height,
			Tiles:  cmd.Tiles,
		},
		Version:   c.Version + 1,
		Timestamp: time.Now(),
	}, nil
}
