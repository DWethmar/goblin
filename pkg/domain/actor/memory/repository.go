package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/dwethmar/goblin/pkg/domain/actor"
)

var _ actor.Repository = &Repository{}

type Repository struct {
	mutex        sync.RWMutex
	actors       map[string]*actor.Actor
	actorsSorted []*actor.Actor
}

// Create implements actor.Repository.
func (r *Repository) Create(ctx context.Context, a *actor.Actor) (*actor.Actor, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.actors[a.ID]; ok {
		return nil, errors.New("actor already exists")
	}

	r.actors[a.ID] = a
	r.actorsSorted = append(r.actorsSorted, a)
	return a, nil
}

// Get implements actor.Repository.
func (r *Repository) Get(ctx context.Context, id string) (*actor.Actor, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	a, ok := r.actors[id]
	if !ok {
		return nil, actor.ErrNotFound
	}

	return a, nil
}

// Delete implements actor.Repository.
func (r *Repository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.actors, id)
	for i, a := range r.actorsSorted {
		if a.ID == id {
			r.actorsSorted = append(r.actorsSorted[:i], r.actorsSorted[i+1:]...)
			break
		}
	}

	return nil
}

// Update implements actor.Repository.
func (r *Repository) Update(ctx context.Context, a *actor.Actor) (*actor.Actor, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.actors[a.ID]; !ok {
		return nil, actor.ErrNotFound
	}

	r.actors[a.ID] = a
	return a, nil
}

// List implements actor.Repository.
func (r *Repository) List(ctx context.Context, offset, limit int) ([]*actor.Actor, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if offset > len(r.actorsSorted) {
		return []*actor.Actor{}, nil
	}

	if offset+limit > len(r.actorsSorted) {
		return r.actorsSorted[offset:], nil
	}

	return r.actorsSorted[offset : offset+limit], nil
}

func NewRepository() *Repository {
	return &Repository{
		mutex:        sync.RWMutex{},
		actors:       make(map[string]*actor.Actor),
		actorsSorted: make([]*actor.Actor, 0),
	}
}
