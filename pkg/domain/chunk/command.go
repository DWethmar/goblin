package chunk

import (
	"time"

	"github.com/dwethmar/goblin/pkg/aggr"
)

var (
	_ aggr.Command = &CreateCommand{}
)

const CreateCommandType = "chunk.create"

type CreateCommand struct {
	ChunkID   string
	X, Y      int
	Width     int
	Height    int
	Tiles     []Tile
	timestamp time.Time
}

func (c *CreateCommand) AggregateID() string         { return c.ChunkID }
func (c *CreateCommand) AggregateType() string       { return AggregateType }
func (c *CreateCommand) CommandTimestamp() time.Time { return c.timestamp }
