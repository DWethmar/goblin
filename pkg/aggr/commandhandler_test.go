package aggr

import (
	"context"
	"testing"
)

func TestCommandHandlerFunc_HandleCommand(t *testing.T) {
	t.Run("should call function", func(t *testing.T) {
		var called bool
		f := CommandHandlerFunc(func(_ context.Context, _ Command) (*Event, error) {
			called = true
			return nil, nil
		})

		_, err := f.HandleCommand(context.Background(), nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !called {
			t.Error("expected function to be called")
		}
	})
}
