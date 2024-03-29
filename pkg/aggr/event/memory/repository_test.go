package memory

import (
	"reflect"
	"sync"
	"testing"

	"github.com/dwethmar/goblin/pkg/aggr"
)

func TestEventRepository_Add(t *testing.T) {
	type fields struct {
		eventsMux sync.Mutex
		events    []*aggr.Event
	}
	type args struct {
		events []*aggr.Event
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
			r := &EventRepository{
				eventsMux: tt.fields.eventsMux,
				events:    tt.fields.events,
			}
			if err := r.Add(tt.args.events); (err != nil) != tt.wantErr {
				t.Errorf("EventRepository.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventRepository_List(t *testing.T) {
	type fields struct {
		events []*aggr.Event
	}
	type args struct {
		aggregateID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*aggr.Event
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EventRepository{
				eventsMux: sync.Mutex{},
				events:    tt.fields.events,
			}
			got, err := r.List(tt.args.aggregateID)
			if (err != nil) != tt.wantErr {
				t.Errorf("EventRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventRepository_All(t *testing.T) {
	type fields struct {
		events []*aggr.Event
	}
	type args struct {
		err chan<- error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan *aggr.Event
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &EventRepository{
				eventsMux: sync.Mutex{},
				events:    tt.fields.events,
			}
			if got := r.All(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventRepository.All() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEventRepository(t *testing.T) {
	tests := []struct {
		name string
		want *EventRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
