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
	t.Run("invalid version", func(t *testing.T) {
		a := &Aggregate{}
		err := a.HandleEvent(&Event{Version: 2})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("model is nil", func(t *testing.T) {
		a := &Aggregate{}
		err := a.HandleEvent(&Event{Version: 1})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("model is not nil", func(t *testing.T) {
		a := &Aggregate{
			Model: &mockModel{
				handleEvent: func(event *Event) error {
					return nil
				},
			},
		}
		err := a.HandleEvent(&Event{Version: 1})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("version is updated", func(t *testing.T) {
		a := &Aggregate{
			Model: &mockModel{
				handleEvent: func(event *Event) error {
					return nil
				},
			},
			Version: 0,
		}
		a.HandleEvent(&Event{Version: 1})
		if a.Version != 1 {
			t.Errorf("unexpected version: %d", a.Version)
		}
	})

	t.Run("failed to handle event", func(t *testing.T) {
		a := &Aggregate{
			Model: &mockModel{
				handleEvent: func(event *Event) error {
					return errors.New("failed to handle event")
				},
			},
		}
		err := a.HandleEvent(&Event{Version: 1})
		if err == nil {
			t.Error("expected error")
		}
	})
}
