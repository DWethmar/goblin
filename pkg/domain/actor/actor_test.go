package actor

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/dwethmar/goblin/pkg/domain"
	"github.com/google/go-cmp/cmp"
)

func TestActor_AggregateID(t *testing.T) {
	type fields struct {
		ID string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should return correct ID",
			fields: fields{
				ID: "123",
			},
			want: "123",
		},
		{
			name: "should return correct empty ID",
			fields: fields{
				ID: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Actor{
				ID: tt.fields.ID,
			}
			if got := a.AggregateID(); got != tt.want {
				t.Errorf("Actor.AggregateID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActor_AggregateVersion(t *testing.T) {
	type fields struct {
		Version uint
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		{
			name: "should return correct version",
			fields: fields{
				Version: 123,
			},
			want: 123,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Actor{
				Version: tt.fields.Version,
			}
			if got := a.AggregateVersion(); got != tt.want {
				t.Errorf("Actor.AggregateVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActor_HandleCommand(t *testing.T) {
	type fields struct {
		ID      string
		Version uint
		Name    string
		X       int
		Y       int
		state   domain.State
		events  []*aggr.Event
	}
	type args struct {
		ctx context.Context
		cmd aggr.Command
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *aggr.Event
		wantErr bool
		Err     error
	}{
		{
			name: "should return error if command is nil",
			fields: fields{
				state: domain.StateDraft,
			},
			args: args{
				cmd: nil,
			},
			want:    nil,
			wantErr: true,
			Err:     ErrNilCommand,
		},
		{
			name: "should return error if command is not handled",
			fields: fields{
				state: domain.StateCreated,
			},
			args: args{
				cmd: &aggr.MockModel{
					Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			want:    nil,
			wantErr: true,
			Err:     ErrUnknownCommandType,
		},
		{
			name: "should return error if actor is deleted",
			fields: fields{
				state: domain.StateDeleted,
			},
			args: args{
				cmd: &CreateCommand{},
			},
			want:    nil,
			wantErr: true,
			Err:     ErrActorIsDeleted,
		},
		{
			name: "should return error if actor is draft and command is not create",
			fields: fields{
				state: domain.StateDraft,
			},
			args: args{
				cmd: &MoveCommand{},
			},
			want:    nil,
			wantErr: true,
			Err:     ErrActorDoesNotExist,
		},
		{
			name: "should return error if actor is created and command is create",
			fields: fields{
				state: domain.StateCreated,
			},
			args: args{
				cmd: &CreateCommand{},
			},
			want:    nil,
			wantErr: true,
			Err:     ErrActorAlreadyCreated,
		},
		{
			name: "should return error if deletig a draft actor",
			fields: fields{
				state: domain.StateDeleted,
			},
			args: args{
				cmd: &MoveCommand{},
			},
			want:    nil,
			wantErr: true,
			Err:     ErrActorIsDeleted,
		},
		{
			name: "should return no error when deleting an created actor",
			fields: fields{
				ID:    "123",
				state: domain.StateCreated,
			},
			args: args{
				cmd: &DestroyCommand{
					ActorID:   "123",
					Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			want: &aggr.Event{
				AggregateID: "123",
				EventType:   DestroyedEventType,
				Data:        nil,
				Version:     1,
				Timestamp:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Actor{
				ID:      tt.fields.ID,
				Version: tt.fields.Version,
				Name:    tt.fields.Name,
				X:       tt.fields.X,
				Y:       tt.fields.Y,
				State:   tt.fields.state,
				Events:  tt.fields.events,
			}
			got, err := a.HandleCommand(tt.args.ctx, tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Actor.HandleCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && !errors.Is(err, tt.Err) {
				t.Errorf("Actor.HandleCommand() error = %v, wantErr %v", err, tt.Err)
				return
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Actor.HandleCommand() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestActor_HandleEvent(t *testing.T) {
	type fields struct {
		ID      string
		Version uint
		Name    string
		X       int
		Y       int
		state   domain.State
		events  []*aggr.Event
	}
	type args struct {
		ctx   context.Context
		event *aggr.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test casaggr.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Actor{
				ID:      tt.fields.ID,
				Version: tt.fields.Version,
				Name:    tt.fields.Name,
				X:       tt.fields.X,
				Y:       tt.fields.Y,
				State:   tt.fields.state,
				Events:  tt.fields.events,
			}
			if err := a.HandleEvent(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Actor.HandleEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestActor_AggregateEvents(t *testing.T) {
	type fields struct {
		ID      string
		Version uint
		Name    string
		X       int
		Y       int
		state   domain.State
		events  []*aggr.Event
	}
	tests := []struct {
		name   string
		fields fields
		want   []*aggr.Event
	}{
		// TODO: Add test casaggr.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Actor{
				ID:      tt.fields.ID,
				Version: tt.fields.Version,
				Name:    tt.fields.Name,
				X:       tt.fields.X,
				Y:       tt.fields.Y,
				State:   tt.fields.state,
				Events:  tt.fields.events,
			}
			if got := a.AggregateEvents(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Actor.AggregateEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActor_ClearAggregateEvents(t *testing.T) {
	type fields struct {
		ID      string
		Version uint
		Name    string
		X       int
		Y       int
		state   domain.State
		events  []*aggr.Event
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test casaggr.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Actor{
				ID:      tt.fields.ID,
				Version: tt.fields.Version,
				Name:    tt.fields.Name,
				X:       tt.fields.X,
				Y:       tt.fields.Y,
				State:   tt.fields.state,
				Events:  tt.fields.events,
			}
			a.ClearAggregateEvents()
		})
	}
}

func TestActor_Deleted(t *testing.T) {
	type fields struct {
		ID      string
		Version uint
		Name    string
		X       int
		Y       int
		state   domain.State
		events  []*aggr.Event
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test casaggr.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Actor{
				ID:      tt.fields.ID,
				Version: tt.fields.Version,
				Name:    tt.fields.Name,
				X:       tt.fields.X,
				Y:       tt.fields.Y,
				State:   tt.fields.state,
				Events:  tt.fields.events,
			}
			if got := a.Deleted(); got != tt.want {
				t.Errorf("Actor.Deleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("should return new instance of Actor", func(t *testing.T) {
		got := New("1", "test", 1, 1)
		want := &Actor{
			ID:     "1",
			Name:   "test",
			State:  domain.StateDraft,
			X:      1,
			Y:      1,
			Events: []*aggr.Event{},
		}
		if diff := cmp.Diff(got, want, cmp.AllowUnexported(Actor{})); diff != "" {
			t.Errorf("New() mismatch (-want +got):\n%s", diff)
		}
	})
}
