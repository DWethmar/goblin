package actor

import "github.com/dwethmar/tards/pkg/es"

const CreateCommandType = "create"

var _ es.Command = &CreateCommand{}

type CreateCommand struct {
	ActorID string
	Name    string
}

func (c *CreateCommand) AggregateID() string   { return c.ActorID }
func (c *CreateCommand) CommandType() string   { return CreateCommandType }
func (c *CreateCommand) AggregateType() string { return AggregateType }

const MoveCommandType = "move"

var _ es.Command = &CreateCommand{}

type MoveCommand struct {
	ActorID string
	Name    string
}

func (c *MoveCommand) AggregateID() string   { return c.ActorID }
func (c *MoveCommand) CommandType() string   { return MoveCommandType }
func (c *MoveCommand) AggregateType() string { return AggregateType }
