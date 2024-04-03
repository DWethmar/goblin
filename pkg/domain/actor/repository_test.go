package actor

import (
	"context"
	"fmt"
	"testing"
)

func TestMerge(t *testing.T) {
	t.Run("merge repos", func(t *testing.T) {
		targetActors := []*Actor{}
		target := &MockRepository{
			CreateFunc: func(ctx context.Context, a *Actor) (*Actor, error) {
				targetActors = append(targetActors, a)
				return a, nil
			},
		}

		source1Actors := []*Actor{}
		for i := 0; i < 150; i++ {
			source1Actors = append(source1Actors, &Actor{ID: fmt.Sprintf("%d", i)})
		}

		source1 := &MockRepository{
			ListFunc: func(ctx context.Context, limit, offset int) ([]*Actor, error) {
				if offset >= len(source1Actors) {
					return []*Actor{}, nil
				}

				if offset+limit > len(source1Actors) {
					return source1Actors[offset:], nil
				}

				return source1Actors[offset : offset+limit], nil
			},
		}

		source2Actors := []*Actor{}
		for i := 150; i < 250; i++ {
			source2Actors = append(source2Actors, &Actor{ID: fmt.Sprintf("%d", i)})
		}

		source2 := &MockRepository{
			ListFunc: func(ctx context.Context, limit, offset int) ([]*Actor, error) {
				if offset >= len(source2Actors) {
					return []*Actor{}, nil
				}

				if offset+limit > len(source2Actors) {
					return source2Actors[offset:], nil
				}

				return source2Actors[offset : offset+limit], nil
			},
		}

		if err := Merge(context.Background(), target, source1, source2); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(targetActors) != 250 {
			t.Fatalf("expected 250 actors, got %d", len(targetActors))
		}
	})
}
