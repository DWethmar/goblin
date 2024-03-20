package chunk

import "github.com/dwethmar/goblin/pkg/aggr"

var (
	_ aggr.Command = &CreateCommand{}
)

const CreateCommandType = "chunk.create"

type CreateCommand struct {
	ChunkID string
	X, Y    int
	Width   int
	Height  int
}

func (c *CreateCommand) AggregateID() string   { return c.ChunkID }
func (c *CreateCommand) CommandType() string   { return CreateCommandType }
func (c *CreateCommand) AggregateType() string { return AggregateType }
