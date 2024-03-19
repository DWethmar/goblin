package encoding

import (
	"testing"

	"github.com/dwethmar/goblin/pkg/aggr"
	"github.com/google/go-cmp/cmp"
)

func TestDecoder_Decode(t *testing.T) {
	t.Run("should decode event", func(t *testing.T) {
		decoder := &Decoder{}
		event := &aggr.Event{
			AggregateID: "aggregateID",
			Type:        "type",
			Data:        []byte("data"),
		}

		encoder := &Encoder{}
		encoded, encodeErr := encoder.Encode(event)
		if encodeErr != nil {
			t.Fatalf("encoding failed: %v", encodeErr)
		}

		decoded, decodeErr := decoder.Decode(encoded)
		if decodeErr != nil {
			t.Fatalf("decoding failed: %v", decodeErr)
		}

		if diff := cmp.Diff(event, decoded); diff != "" {
			t.Errorf("unexpected event (-want +got):\n%s", diff)
		}
	})
}
