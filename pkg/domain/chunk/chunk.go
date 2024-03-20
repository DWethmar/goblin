package chunk

import (
	"context"
	"errors"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/domain"
	"github.com/dwethmar/goblin/pkg/matrix"
)

var _ aggr.Model = &Chunk{}

const AggregateType = "chunk"

var (
	ErrNilCommand          = errors.New("command is nil")
	ErrUnknownCommandType  = errors.New("unknown command type")
	ErrChunkIsDeleted      = errors.New("chunk is deleted")
	ErrChunkDoesNotExist   = errors.New("chunk does not exist")
	ErrChunkAlreadyCreated = errors.New("chunk already created")
)

type Chunk struct {
	ID      string
	Version int
	X       int
	Y       int
	Tiles   matrix.Matrix

	state  domain.State
	events []*aggr.Event
}

func New(id string, x, y int) *Chunk {
	return &Chunk{
		ID:      id,
		X:       x,
		Y:       y,
		Tiles:   [][]int{},
		Version: 0,
	}
}

func (c *Chunk) AggregateID() string   { return c.ID }
func (c *Chunk) AggregateVersion() int { return c.Version }

func (c *Chunk) HandleCommand(cmd aggr.Command) (*aggr.Event, error) {
	if cmd == nil {
		return nil, ErrNilCommand
	}

	if domain.StateDeleted.Is(c.state) {
		return nil, ErrChunkIsDeleted
	}

	switch v := cmd.(type) {
	case *CreateCommand:
		return CreateChunkCommandHandler(c, v)
	}

	return nil, ErrUnknownCommandType
}

func (c *Chunk) HandleEvent(_ context.Context, event *aggr.Event) error {
	switch event.Type {
	case CreatedEventType:
		return HandleCreatedEvent(c, event)
	}

	return nil
}

func (c *Chunk) AggregateEvents() []*aggr.Event { return c.events }
func (c *Chunk) ClearAggregateEvents()          { c.events = []*aggr.Event{} }
func (c *Chunk) Deleted() bool                  { return domain.StateDeleted.Is(c.state) }
