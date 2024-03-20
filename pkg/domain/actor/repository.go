package actor

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("actor not found")
)

type Reader interface {
	Get(ctx context.Context, id string) (*Actor, error)
	List(ctx context.Context, limit, offset int) ([]*Actor, error)
}

type Writer interface {
	Create(ctx context.Context, a *Actor) (*Actor, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, a *Actor) (*Actor, error)
}

type Repository interface {
	Reader
	Writer
}
