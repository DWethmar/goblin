package clock

import (
	"testing"
)

func TestClock_Now(t *testing.T) {
	t.Run("get current time", func(t *testing.T) {
		c := New()
		now := c.Now()
		if now.IsZero() {
			t.Errorf("Clock.Now() = %v, want non zero", now)
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("create new clock", func(t *testing.T) {
		c := New()
		if c == nil {
			t.Errorf("New() = %v, want non-nil", c)
		}

		// Check if NowFunc is set and returns a non-zero time
		if c.NowFunc == nil {
			t.Errorf("New() NowFunc is nil, want non-nil")
		}
	})
}
