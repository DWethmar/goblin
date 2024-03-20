package chunk

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("chunk not found")
)

type Reader interface {
	Get(ctx context.Context, id string) (*Chunk, error)
	List(ctx context.Context, limit, offset int) ([]*Chunk, error)
}

type Writer interface {
	Create(ctx context.Context, c *Chunk) (*Chunk, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, c *Chunk) (*Chunk, error)
}

type Repository interface {
	Reader
	Writer
}
