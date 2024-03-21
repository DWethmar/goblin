package services

import (
	"context"
	"reflect"
	"testing"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/clock"
	"github.com/dwethmar/goblin/pkg/domain/actor"
)

func TestActors_Create(t *testing.T) {
	t.Run("create actor", func(t *testing.T) {

	})
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
