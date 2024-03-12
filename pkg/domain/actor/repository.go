package actor

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("actor not found")
)

type Repository interface {
	Get(ctx context.Context, id string) (*Actor, error)
	Create(ctx context.Context, a *Actor) (*Actor, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, a *Actor) (*Actor, error)
	List(ctx context.Context, limit, offset int) ([]*Actor, error)
}
