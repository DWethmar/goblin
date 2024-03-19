package encoding

import (
	"bytes"
	"encoding/gob"

	"github.com/dwethmar/goblin/pkg/aggr"
)

type Decoder struct{}

func (d *Decoder) Decode(input []byte) (*aggr.Event, error) {
	var event aggr.Event
	buf := bytes.NewBuffer(input)
	if err := gob.NewDecoder(buf).Decode(&event); err != nil {
		return nil, err
	}
	return &event, nil
}
