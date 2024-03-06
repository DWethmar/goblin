package event

import "github.com/dwethmar/goblin/pkg/es"

type Decoder interface {
	Decode([]byte) (*es.Event, error)
}

type Encoder interface {
	Encode(*es.Event) ([]byte, error)
}
