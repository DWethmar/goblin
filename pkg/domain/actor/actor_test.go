package actor

import (
	"reflect"
	"testing"

	"github.com/dwethmar/goblin/pkg/aggr"
)

func TestState_Is(t *testing.T) {
	type args struct {
		v State
	}
	tests := []struct {
		name string
		s    State
		args args
		want bool
	}{
		{"test1", StateDraft, args{StateDraft}, true},
		{"test2", StateDraft, args{StateCreated}, false},
		{"test3", StateDraft, args{StateDeleted}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Is(tt.args.v); got != tt.want {
				t.Errorf("State.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActor_AggregateID(t *testing.T) {
	type fields struct {
		ID      string
		Version int
		Name    string
		X       int
		Y       int
		state   State
		events  []*aggr.Event
	}
	tests := []struct {
		name   string
		fields fields
		want   string
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
				state:   tt.fields.state,
				events:  tt.fields.events,
			}
			if got := a.AggregateID(); got != tt.want {
				t.Errorf("Actor.AggregateID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActor_AggregateVersion(t *testing.T) {
	type fields struct {
		ID      string
		Version int
		Name    string
		X       int
		Y       int
		state   State
		events  []*aggr.Event
	}
	tests := []struct {
		name   string
		fields fields
		want   int
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
				state:   tt.fields.state,
				events:  tt.fields.events,
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
		Version int
		Name    string
		X       int
		Y       int
		state   State
		events  []*aggr.Event
	}
	type args struct {
		cmd aggr.Command
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *aggr.Event
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
				state:   tt.fields.state,
				events:  tt.fields.events,
			}
			got, err := a.HandleCommand(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Actor.HandleCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Actor.HandleCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActor_HandleEvent(t *testing.T) {
	type fields struct {
		ID      string
		Version int
		Name    string
		X       int
		Y       int
		state   State
		events  []*aggr.Event
	}
	type args struct {
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
				state:   tt.fields.state,
				events:  tt.fields.events,
			}
			if err := a.HandleEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Actor.HandleEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestActor_AggregateEvents(t *testing.T) {
	type fields struct {
		ID      string
		Version int
		Name    string
		X       int
		Y       int
		state   State
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
				state:   tt.fields.state,
				events:  tt.fields.events,
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
		Version int
		Name    string
		X       int
		Y       int
		state   State
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
				state:   tt.fields.state,
				events:  tt.fields.events,
			}
			a.ClearAggregateEvents()
		})
	}
}

func TestActor_Deleted(t *testing.T) {
	type fields struct {
		ID      string
		Version int
		Name    string
		X       int
		Y       int
		state   State
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
				state:   tt.fields.state,
				events:  tt.fields.events,
			}
			if got := a.Deleted(); got != tt.want {
				t.Errorf("Actor.Deleted() = %v, want %v", got, tt.want)
			}
		})
	}
}
