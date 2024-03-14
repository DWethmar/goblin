package actor

import "github.com/dwethmar/goblin/pkg/es"

var (
	_ es.Command = &CreateCommand{}
	_ es.Command = &MoveCommand{}
)

const CreateCommandType = "actor.create"

type CreateCommand struct {
	ActorID string
	Name    string
	X, Y    int
}

func (c *CreateCommand) AggregateID() string   { return c.ActorID }
func (c *CreateCommand) CommandType() string   { return CreateCommandType }
func (c *CreateCommand) AggregateType() string { return AggregateType }

const MoveCommandType = "actor.move"

type MoveCommand struct {
	ActorID string
	X, Y    int
}

func (c *MoveCommand) AggregateID() string   { return c.ActorID }
func (c *MoveCommand) CommandType() string   { return MoveCommandType }
func (c *MoveCommand) AggregateType() string { return AggregateType }
