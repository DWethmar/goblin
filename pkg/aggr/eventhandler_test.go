package aggr

import (
	"context"
	"testing"
)

func TestEventHandlerFunc_HandleEvent(t *testing.T) {
	t.Run("should call function", func(t *testing.T) {
		var called bool
		f := EventHandlerFunc(func(_ context.Context, _ *Event) error {
			called = true
			return nil
		})

		err := f.HandleEvent(context.Background(), &Event{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !called {
			t.Error("expected function to be called")
		}
	})
}
