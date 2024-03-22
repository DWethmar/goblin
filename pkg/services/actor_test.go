package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/clock"
	"github.com/dwethmar/goblin/pkg/domain/actor"
)

func TestActors_Create(t *testing.T) {
	type fields struct {
		clock       *clock.Clock
		actorReader actor.Reader
		commandBus  aggr.CommandBus
	}
	type args struct {
		ctx  context.Context
		id   string
		name string
		x    int
		y    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create actor should not return error",
			fields: fields{
				clock: clock.New(),
				actorReader: &actor.MockRepository{
					GetFunc: func(_ context.Context, _ string) (*actor.Actor, error) {
						return nil, errors.New("could nto find actor")
					},
				},
				commandBus: &aggr.MockCommandBus{
					HandleCommandFunc: func(_ context.Context, _ aggr.Command) error {
						return nil
					},
				},
			},
			args: args{
				ctx:  context.Background(),
				id:   "1",
				name: "test",
				x:    0,
				y:    0,
			},
			wantErr: false,
		},
		{
			name: "create actor already exists",
			fields: fields{
				clock: clock.New(),
				actorReader: &actor.MockRepository{
					GetFunc: func(_ context.Context, _ string) (*actor.Actor, error) {
						return &actor.Actor{}, nil
					},
				},
			},
			args: args{
				ctx:  context.Background(),
				id:   "1",
				name: "test",
				x:    0,
				y:    0,
			},
			wantErr: true,
		},
		{
			name: "if the existance of the actor could not be checked, an error should be returned",
			fields: fields{
				clock: clock.New(),
				actorReader: &actor.MockRepository{
					GetFunc: func(_ context.Context, _ string) (*actor.Actor, error) {
						return nil, nil
					},
				},
				commandBus: nil,
			},
			args: args{
				ctx:  context.Background(),
				id:   "1",
				name: "test",
				x:    0,
				y:    0,
			},
			wantErr: true,
		},
		{
			name: "error should be retunrned if actor could not be created",
			fields: fields{
				clock: clock.New(),
				actorReader: &actor.MockRepository{
					GetFunc: func(_ context.Context, _ string) (*actor.Actor, error) {
						return nil, errors.New("could nto find actor")
					},
				},
				commandBus: &aggr.MockCommandBus{
					HandleCommandFunc: func(_ context.Context, _ aggr.Command) error {
						return errors.New("could not create actor")
					},
				},
			},
			args: args{
				ctx:  context.Background(),
				id:   "1",
				name: "test",
				x:    0,
				y:    0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Actors{
				clock:       tt.fields.clock,
				actorReader: tt.fields.actorReader,
				commandBus:  tt.fields.commandBus,
			}
			if err := a.Create(tt.args.ctx, tt.args.id, tt.args.name, tt.args.x, tt.args.y); (err != nil) != tt.wantErr {
				t.Errorf("Actors.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestActors_Move(t *testing.T) {
	type fields struct {
		clock       *clock.Clock
		actorReader actor.Reader
		commandBus  aggr.CommandBus
	}
	type args struct {
		ctx context.Context
		id  string
		x   int
		y   int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Actors{
				clock:       tt.fields.clock,
				actorReader: tt.fields.actorReader,
				commandBus:  tt.fields.commandBus,
			}
			if err := a.Move(tt.args.ctx, tt.args.id, tt.args.x, tt.args.y); (err != nil) != tt.wantErr {
				t.Errorf("Actors.Move() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestActors_Get(t *testing.T) {
	type fields struct {
		clock       *clock.Clock
		actorReader actor.Reader
		commandBus  aggr.CommandBus
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *actor.Actor
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Actors{
				clock:       tt.fields.clock,
				actorReader: tt.fields.actorReader,
				commandBus:  tt.fields.commandBus,
			}
			got, err := a.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Actors.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Actors.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActors_List(t *testing.T) {
	type fields struct {
		clock       *clock.Clock
		actorReader actor.Reader
		commandBus  aggr.CommandBus
	}
	type args struct {
		ctx    context.Context
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*actor.Actor
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Actors{
				clock:       tt.fields.clock,
				actorReader: tt.fields.actorReader,
				commandBus:  tt.fields.commandBus,
			}
			got, err := a.List(tt.args.ctx, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Actors.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Actors.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewActorService(t *testing.T) {
	type args struct {
		repo       actor.Repository
		commandBus aggr.CommandBus
	}
	tests := []struct {
		name string
		args args
		want *Actors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewActorService(tt.args.repo, tt.args.commandBus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewActorService() = %v, want %v", got, tt.want)
			}
		})
	}
}
