package json

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/chaos-io/chaos/config/encoder"
)

type jsonEncoder struct{}

func (j jsonEncoder) Encode(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func (j jsonEncoder) Decode(d []byte, v interface{}) error {
	return jsoniter.Unmarshal(d, v)
}

func (j jsonEncoder) String() string {
	return "json"
}

func NewEncoder() encoder.Encoder {
	return jsonEncoder{}
}
