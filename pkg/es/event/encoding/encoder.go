package encoding

import (
	"bytes"
	"encoding/gob"

	"github.com/dwethmar/goblin/pkg/es"
)

type Encoder struct{}

func (e *Encoder) Encode(event *es.Event) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(event); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
