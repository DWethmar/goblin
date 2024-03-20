package chunk

import (
	"context"
	"errors"
	"fmt"

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
	Version uint
	X       int
	Y       int
	Width   int
	Height  int
	Tiles   matrix.Matrix

	state  domain.State
	events []*aggr.Event
}

func New(id string, x, y int) *Chunk {
	return &Chunk{
		ID:      id,
		X:       x,
		Y:       y,
		Tiles:   nil,
		Version: 0,
		state:   domain.StateDraft,
	}
}

func (c *Chunk) AggregateID() string    { return c.ID }
func (c *Chunk) AggregateVersion() uint { return c.Version }

func (c *Chunk) HandleCommand(cmd aggr.Command) (*aggr.Event, error) {
	if cmd == nil {
		return nil, ErrNilCommand
	}

	// if state is draft and command is not create, return error
	if domain.StateDraft.Is(c.state) {
		if _, ok := cmd.(*CreateCommand); !ok {
			return nil, ErrChunkAlreadyCreated
		}
	}

	// if state is deleted, return error
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
	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	switch event.Type {
	case CreatedEventType:
		return HandleCreatedEvent(c, event)
	}

	c.Version = event.Version
	c.events = append(c.events, event)

	return nil
}

func (c *Chunk) AggregateEvents() []*aggr.Event { return c.events }
func (c *Chunk) ClearAggregateEvents()          { c.events = []*aggr.Event{} }
func (c *Chunk) Deleted() bool                  { return domain.StateDeleted.Is(c.state) }
