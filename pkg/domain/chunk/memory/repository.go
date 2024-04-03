package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/dwethmar/goblin/pkg/domain/chunk"
)

var _ chunk.Repository = &Repository{}

type Repository struct {
	mutex        sync.RWMutex
	chunks       map[string]*chunk.Chunk
	chunksSorted []*chunk.Chunk
}

// Create implements chunk.Repository
func (r *Repository) Create(ctx context.Context, c *chunk.Chunk) (*chunk.Chunk, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.chunks[c.ID]; ok {
		return nil, errors.New("actor already exists")
	}

	r.chunks[c.ID] = c
	r.chunksSorted = append(r.chunksSorted, c)
	return c, nil
}

// Get implements chunk.Repository
func (r *Repository) Get(ctx context.Context, id string) (*chunk.Chunk, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	a, ok := r.chunks[id]
	if !ok {
		return nil, chunk.ErrNotFound
	}

	return a, nil
}

// Delete implements chunk.Repository
func (r *Repository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.chunks, id)
	for i, a := range r.chunksSorted {
		if a.ID == id {
			r.chunksSorted = append(r.chunksSorted[:i], r.chunksSorted[i+1:]...)
			break
		}
	}

	return nil
}

// Update implements chunk.Repository
func (r *Repository) Update(ctx context.Context, c *chunk.Chunk) (*chunk.Chunk, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.chunks[c.ID]; !ok {
		return nil, chunk.ErrNotFound
	}

	r.chunks[c.ID] = c
	return c, nil
}

// List implements chunk.Repository
func (r *Repository) List(ctx context.Context, limit, offset int) ([]*chunk.Chunk, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if offset > len(r.chunksSorted) {
		return []*chunk.Chunk{}, nil // empty list
	}

	if offset+limit > len(r.chunksSorted) {
		return r.chunksSorted[offset:], nil
	}

	return r.chunksSorted[offset : offset+limit], nil
}

func NewRepository() *Repository {
	return &Repository{
		mutex:        sync.RWMutex{},
		chunks:       make(map[string]*chunk.Chunk),
		chunksSorted: []*chunk.Chunk{},
	}
}
