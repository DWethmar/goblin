package actor

import (
	"time"

	"github.com/dwethmar/goblin/pkg/aggr"
)

var (
	_ aggr.Command = &CreateCommand{}
	_ aggr.Command = &MoveCommand{}
)

const CreateCommandType = "actor.create"

type CreateCommand struct {
	ActorID   string
	Name      string
	X, Y      int
	Timestamp time.Time
}

func (c *CreateCommand) AggregateID() string   { return c.ActorID }
func (c *CreateCommand) CommandType() string   { return CreateCommandType }
func (c *CreateCommand) AggregateType() string { return AggregateType }

const DestroyCommandType = "actor.destroy"

type DestroyCommand struct {
	ActorID   string
	Timestamp time.Time
}

func (c *DestroyCommand) AggregateID() string   { return c.ActorID }
func (c *DestroyCommand) CommandType() string   { return DestroyCommandType }
func (c *DestroyCommand) AggregateType() string { return AggregateType }

const MoveCommandType = "actor.move"

type MoveCommand struct {
	ActorID   string
	X, Y      int
	Timestamp time.Time
}

func (c *MoveCommand) AggregateID() string   { return c.ActorID }
func (c *MoveCommand) CommandType() string   { return MoveCommandType }
func (c *MoveCommand) AggregateType() string { return AggregateType }
