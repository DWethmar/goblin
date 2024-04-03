package chunk

import (
	"context"
	"fmt"
	"testing"
)

func TestMerge(t *testing.T) {
	t.Run("merge repos", func(t *testing.T) {
		targetChunks := []*Chunk{}
		target := &MockRepository{
			CreateFunc: func(ctx context.Context, a *Chunk) (*Chunk, error) {
				targetChunks = append(targetChunks, a)
				return a, nil
			},
		}

		source1Chunks := []*Chunk{}
		for i := 0; i < 150; i++ {
			source1Chunks = append(source1Chunks, &Chunk{ID: fmt.Sprintf("%d", i)})
		}

		source1 := &MockRepository{
			ListFunc: func(ctx context.Context, limit, offset int) ([]*Chunk, error) {
				if offset >= len(source1Chunks) {
					return []*Chunk{}, nil
				}

				if offset+limit > len(source1Chunks) {
					return source1Chunks[offset:], nil
				}

				return source1Chunks[offset : offset+limit], nil
			},
		}

		source2Chunks := []*Chunk{}
		for i := 150; i < 250; i++ {
			source2Chunks = append(source2Chunks, &Chunk{ID: fmt.Sprintf("%d", i)})
		}

		source2 := &MockRepository{
			ListFunc: func(ctx context.Context, limit, offset int) ([]*Chunk, error) {
				if offset >= len(source2Chunks) {
					return []*Chunk{}, nil
				}

				if offset+limit > len(source2Chunks) {
					return source2Chunks[offset:], nil
				}

				return source2Chunks[offset : offset+limit], nil
			},
		}

		if err := Merge(context.Background(), target, source1, source2); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(targetChunks) != 250 {
			t.Fatalf("expected 250 Chunks, got %d", len(targetChunks))
		}
	})
}
