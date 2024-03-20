package actor

import (
	"context"
	"errors"
)

type MockRepository struct {
	GetFunc    func(ctx context.Context, id string) (*Actor, error)
	CreateFunc func(ctx context.Context, a *Actor) (*Actor, error)
	DeleteFunc func(ctx context.Context, id string) error
	UpdateFunc func(ctx context.Context, a *Actor) (*Actor, error)
	ListFunc   func(ctx context.Context, limit, offset int) ([]*Actor, error)
}

func (r *MockRepository) Get(ctx context.Context, id string) (*Actor, error) {
	if r.GetFunc == nil {
		return nil, errors.New("GetFunc not implemented")
	}
	return r.GetFunc(ctx, id)
}

func (r *MockRepository) Create(ctx context.Context, a *Actor) (*Actor, error) {
	if r.CreateFunc == nil {
		return nil, errors.New("CreateFunc not implemented")
	}
	return r.CreateFunc(ctx, a)
}

func (r *MockRepository) Delete(ctx context.Context, id string) error {
	if r.DeleteFunc == nil {
		return errors.New("DeleteFunc not implemented")
	}
	return r.DeleteFunc(ctx, id)
}

func (r *MockRepository) Update(ctx context.Context, a *Actor) (*Actor, error) {
	if r.UpdateFunc == nil {
		return nil, errors.New("UpdateFunc not implemented")
	}
	return r.UpdateFunc(ctx, a)
}

func (r *MockRepository) List(ctx context.Context, limit, offset int) ([]*Actor, error) {
	if r.ListFunc == nil {
		return nil, errors.New("ListFunc not implemented")
	}
	return r.ListFunc(ctx, limit, offset)
}
