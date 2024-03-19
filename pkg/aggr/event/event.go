package event

import "github.com/dwethmar/goblin/pkg/aggr"

type Decoder interface {
	Decode([]byte) (*aggr.Event, error)
}

type Encoder interface {
	Encode(*aggr.Event) ([]byte, error)
}
