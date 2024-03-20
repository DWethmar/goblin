// Package clock provides a simple clock type that can be used to get the current time.
// Its main purpose is to make it easier to test code that depends on the current time.
package clock

import "time"

// Clock can be used to get the current time.
type Clock struct {
	NowFunc func() time.Time
}

func (c *Clock) Now() time.Time { return c.NowFunc() }

func New() *Clock {
	return &Clock{
		NowFunc: time.Now,
	}
}
