package submodel

import (
	"github.com/stafiprotocol/go-substrate-rpc-client/scale"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

type option struct {
	hasValue bool
}

// IsNone returns true if the value is missing
func (o option) IsNone() bool {
	return !o.hasValue
}

// IsNone returns true if a value is present
func (o option) IsSome() bool {
	return o.hasValue
}

type OptionTimePoint struct {
	option
	value types.TimePoint
}

func NewOptionTimePoint(value types.TimePoint) *OptionTimePoint {
	return &OptionTimePoint{option{true}, value}
}

func NewOptionTimePointEmpty() *OptionTimePoint {
	return &OptionTimePoint{option: option{false}}
}

func (o OptionTimePoint) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

func (o *OptionTimePoint) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

// SetSome sets a value
func (o *OptionTimePoint) SetSome(value types.TimePoint) {
	o.hasValue = true
	o.value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionTimePoint) SetNone() {
	o.hasValue = false
	o.value = types.TimePoint{}
}
