package memory

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/dwethmar/goblin/pkg/domain/actor"
	"github.com/google/go-cmp/cmp"
)

func TestRepository_Create(t *testing.T) {
	type fields struct {
		actors       map[string]*actor.Actor
		actorsSorted []*actor.Actor
	}
	type args struct {
		ctx context.Context
		a   *actor.Actor
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *actor.Actor
		wantErr bool
	}{
		{
			name: "create actor",
			fields: fields{
				actors:       map[string]*actor.Actor{},
				actorsSorted: []*actor.Actor{},
			},
			args: args{
				ctx: context.Background(),
				a: &actor.Actor{
					ID:   "1",
					Name: "test",
				},
			},
			want: &actor.Actor{
				ID:   "1",
				Name: "test",
			},
			wantErr: false,
		},
		{
			name: "create actor already exists",
			fields: fields{
				actors: map[string]*actor.Actor{
					"1": {
						ID:   "1",
						Name: "test",
					},
				},
				actorsSorted: []*actor.Actor{},
			},
			args: args{
				ctx: context.Background(),
				a: &actor.Actor{
					ID:   "1",
					Name: "test",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				mutex:        sync.RWMutex{},
				actors:       tt.fields.actors,
				actorsSorted: tt.fields.actorsSorted,
			}
			got, err := r.Create(tt.args.ctx, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(actor.Actor{})); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestRepository_Get(t *testing.T) {
	type fields struct {
		actors       map[string]*actor.Actor
		actorsSorted []*actor.Actor
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
		{
			name: "get actor",
			fields: fields{
				actors: map[string]*actor.Actor{
					"1": {
						ID:   "1",
						Name: "test",
					},
				},
				actorsSorted: []*actor.Actor{},
			},
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			want: &actor.Actor{
				ID:   "1",
				Name: "test",
			},
		},
		{
			name: "get actor not found",
			fields: fields{
				actors: map[string]*actor.Actor{},
			},
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				mutex:        sync.RWMutex{},
				actors:       tt.fields.actors,
				actorsSorted: tt.fields.actorsSorted,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	t.Run("delete actor", func(t *testing.T) {
		r := &Repository{
			mutex: sync.RWMutex{},
			actors: map[string]*actor.Actor{
				"1": {
					ID:   "1",
					Name: "test",
				},
			},
			actorsSorted: []*actor.Actor{},
		}
		err := r.Delete(context.Background(), "1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if _, ok := r.actors["1"]; ok {
			t.Errorf("actor not deleted")
		}
	})
}

func TestRepository_Update(t *testing.T) {
	type fields struct {
		actors       map[string]*actor.Actor
		actorsSorted []*actor.Actor
	}
	type args struct {
		ctx context.Context
		a   *actor.Actor
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *actor.Actor
		wantErr bool
	}{
		{
			name: "update actor",
			fields: fields{
				actors: map[string]*actor.Actor{
					"1": {
						ID:   "1",
						Name: "test",
					},
				},
				actorsSorted: []*actor.Actor{},
			},
			args: args{
				ctx: context.Background(),
				a: &actor.Actor{
					ID:   "1",
					Name: "test2",
				},
			},
			want: &actor.Actor{
				ID:   "1",
				Name: "test2",
			},
		},
		{
			name: "update actor not found",
			fields: fields{
				actors: map[string]*actor.Actor{},
			},
			args: args{
				ctx: context.Background(),
				a: &actor.Actor{
					ID:   "1",
					Name: "test2",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				mutex:        sync.RWMutex{},
				actors:       tt.fields.actors,
				actorsSorted: tt.fields.actorsSorted,
			}
			got, err := r.Update(tt.args.ctx, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_List(t *testing.T) {
	type fields struct {
		actors       map[string]*actor.Actor
		actorsSorted []*actor.Actor
	}
	type args struct {
		ctx    context.Context
		limit  int
		offset int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*actor.Actor
		wantErr bool
	}{
		{
			name: "list actors",
			fields: fields{
				actors: map[string]*actor.Actor{
					"1": {
						ID:   "1",
						Name: "test",
					},
					"2": {
						ID:   "2",
						Name: "test",
					},
				},
				actorsSorted: []*actor.Actor{
					{
						ID:   "1",
						Name: "test",
					},
					{
						ID:   "2",
						Name: "test",
					},
				},
			},
			args: args{
				ctx:    context.Background(),
				limit:  0,
				offset: 0,
			},
			want: []*actor.Actor{
				{
					ID:   "1",
					Name: "test",
				},
				{
					ID:   "2",
					Name: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				mutex:        sync.RWMutex{},
				actors:       tt.fields.actors,
				actorsSorted: tt.fields.actorsSorted,
			}
			got, err := r.List(tt.args.ctx, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRepository(t *testing.T) {
	tests := []struct {
		name string
		want *Repository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
