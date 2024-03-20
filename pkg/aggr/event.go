package aggr

import (
	"errors"
	"time"
)

var (
	ErrAggregateIDEmpty = errors.New("aggregate id can't be empty")
	ErrTypeEmpty        = errors.New("type can't be empty")
	ErrTimestampZero    = errors.New("timestamp can't be zero")
)

type Event struct {
	AggregateID string
	Type        string
	Data        interface{}
	Version     uint
	Timestamp   time.Time
}

func (e *Event) Validate() error {
	if e.AggregateID == "" {
		return ErrAggregateIDEmpty
	}

	if e.Type == "" {
		return ErrTypeEmpty
	}

	if e.Timestamp.IsZero() {
		return ErrTimestampZero
	}

	return nil
}
