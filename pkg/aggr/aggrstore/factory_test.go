package aggrstore

import (
	"errors"
	"testing"

	"github.com/dwethmar/goblin/pkg/aggr"
)

func TestFactory_Register(t *testing.T) {
	t.Run("Register", func(t *testing.T) {
		f := NewFactory()
		f.Register("test", func(aggregateID string) *aggr.Aggregate {
			return nil
		})

		if _, ok := f.createFuncs["test"]; !ok {
			t.Errorf("Register failed")
		}
	})
}

func TestFactory_Create(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		f := NewFactory()
		f.Register("test", func(aggregateID string) *aggr.Aggregate {
			return &aggr.Aggregate{
				Model: &aggr.MockAggregate{
					ID: aggregateID,
				},
			}
		})

		agg, err := f.Create("test", "123")
		if err != nil {
			t.Errorf("Create failed")
		}

		if id := agg.AggregateID(); id != "123" {
			t.Errorf("Create failed, expected 123, got %s", id)
		}
	})

	t.Run("aggregate type not found", func(t *testing.T) {
		f := NewFactory()
		_, err := f.Create("test", "123")
		if !errors.Is(err, ErrorAggregateTypeNotFound) {
			t.Errorf("Create failed, expected: %v, got %v", ErrorAggregateTypeNotFound, err)
		}
	})
}

func TestNewFactory(t *testing.T) {
	t.Run("NewFactory", func(t *testing.T) {
		f := NewFactory()
		if f == nil {
			t.Errorf("NewFactory failed")
		}
	})

	t.Run("NewFactory with options", func(t *testing.T) {
		f := NewFactory(func(f *Factory) {
			f.Register("test", func(aggregateID string) *aggr.Aggregate {
				return nil
			})
		})

		if _, ok := f.createFuncs["test"]; !ok {
			t.Errorf("NewFactory with options failed")
		}
	})
}
