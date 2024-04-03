package actor

import (
	"context"
	"errors"
	"fmt"
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

// Merge copies all actors from the sources to the target repository
func Merge(ctx context.Context, target Repository, sources ...Repository) error {
	for _, source := range sources {
		var offset int
		for {
			actors, err := source.List(ctx, 100, offset)
			if err != nil {
				return fmt.Errorf("list actors: %w", err)
			}

			if len(actors) == 0 {
				break
			}

			for _, a := range actors {
				fmt.Printf("Creating actor %s\n", a.ID)
				if _, err := target.Create(ctx, a); err != nil {
					return fmt.Errorf("create actor: %w", err)
				}
			}

			offset += len(actors)
		}
	}

	return nil
}
