package es

import (
	"reflect"
	"testing"
)

func TestCommandBus_Dispatch(t *testing.T) {
	t.Run("should dispatch command", func(t *testing.T) {
		aggregateStore := &MockAggregateStore{
			GetFunc: func(aggregateType, aggregateID string) (*Aggregate, error) {
				return &Aggregate{
					Model: &mockModel{
						handleCommand: func(command Command) (*Event, error) {
							return &Event{}, nil
						},
						handleEvent: func(event *Event) error {
							return nil
						},
					},
				}, nil
			},
			SaveFunc: func(p0 *Aggregate) error {
				return nil
			},
		}

		eventBus := &EventBus{}
		commandBus := NewCommandBus(aggregateStore, eventBus)
		err := commandBus.Dispatch(&MockCommand{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestNewCommandBus(t *testing.T) {
	type args struct {
		aggregateStore AggregateStore
		eventBus       *EventBus
	}
	tests := []struct {
		name string
		args args
		want *CommandBus
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommandBus(tt.args.aggregateStore, tt.args.eventBus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommandBus() = %v, want %v", got, tt.want)
			}
		})
	}
}
