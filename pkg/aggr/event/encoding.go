package event

import "github.com/dwethmar/goblin/pkg/aggr"

// Decoder and Encoder are used to encode and decode events.
type Decoder interface {
	Decode([]byte) (*aggr.Event, error)
}

// Encoder is used to encode events.
type Encoder interface {
	Encode(*aggr.Event) ([]byte, error)
}
