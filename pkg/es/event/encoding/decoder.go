package encoding

import (
	"bytes"
	"encoding/gob"

	"github.com/dwethmar/goblin/pkg/es"
)

type Decoder struct{}

func (d *Decoder) Decode(input []byte) (*es.Event, error) {
	var event es.Event
	buf := bytes.NewBuffer(input)
	if err := gob.NewDecoder(buf).Decode(&event); err != nil {
		return nil, err
	}
	return &event, nil
}
