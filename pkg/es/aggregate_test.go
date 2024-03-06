package es

import (
	"errors"
	"testing"
)

type mockModel struct {
	handleCommand func(command Command) (*Event, error)
	handleEvent   func(event *Event) error
}

func (m *mockModel) HandleCommand(command Command) (*Event, error) {
	return m.handleCommand(command)
}

func (m *mockModel) HandleEvent(event *Event) error {
	return m.handleEvent(event)
}

func TestAggregate_HandleCommand(t *testing.T) {
	t.Run("model is nil", func(t *testing.T) {
		a := &Aggregate{}
		_, err := a.HandleCommand(&MockCommand{})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("model is not nil", func(t *testing.T) {
		a := &Aggregate{
			Model: &mockModel{
				handleCommand: func(command Command) (*Event, error) {
					return &Event{}, nil
				},
			},
		}
		_, err := a.HandleCommand(&MockCommand{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("failed to handle command", func(t *testing.T) {
		a := &Aggregate{
			Model: &mockModel{
				handleCommand: func(command Command) (*Event, error) {
					return nil, errors.New("failed to handle command")
				},
			},
		}
		_, err := a.HandleCommand(&MockCommand{})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("event is nil", func(t *testing.T) {
		a := &Aggregate{
			Model: &mockModel{
				handleCommand: func(command Command) (*Event, error) {
					return nil, nil
				},
			},
		}

		_, err := a.HandleCommand(&MockCommand{})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("event is versioned", func(t *testing.T) {
		a := &Aggregate{
			Model: &mockModel{
				handleCommand: func(command Command) (*Event, error) {
					return &Event{}, nil
				},
			},
			Version: 1,
		}
		event, _ := a.HandleCommand(&MockCommand{})
		if event.Version != 2 {
			t.Errorf("unexpected version: %d", event.Version)
		}
	})
}

func TestAggregate_HandleEvent(t *testing.T) {
	type fields struct {
		ID      string
		Type    string
		Events  []*Event
		Model   AggregateModel
		Version int
		Created bool
	}
	type args struct {
		event *Event
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
			a := &Aggregate{
				ID:      tt.fields.ID,
				Type:    tt.fields.Type,
				Events:  tt.fields.Events,
				Model:   tt.fields.Model,
				Version: tt.fields.Version,
				Created: tt.fields.Created,
			}
			if err := a.HandleEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Aggregate.HandleEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
