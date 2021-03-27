package core

import (
	"fmt"
	"github.com/stafiprotocol/go-substrate-rpc-client/scale"
)

type RSymbol string

const (
	RFIS = RSymbol("RFIS")
	RDOT = RSymbol("RDOT")
	RKSM = RSymbol("RKSM")
)

func (r *RSymbol) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		//*r = RFIS
		return fmt.Errorf("RSymbol decode error: %d", b)
	case 1:
		*r = RDOT
	case 2:
		*r = RKSM
	default:
		return fmt.Errorf("RSymbol decode error: %d", b)
	}

	return nil
}

func (r RSymbol) Encode(encoder scale.Encoder) error {
	switch r {
	case RFIS:
		//return encoder.PushByte(0)
		return fmt.Errorf("RFIS not supported")
	case RDOT:
		return encoder.PushByte(1)
	case RKSM:
		return encoder.PushByte(2)
	default:
		return fmt.Errorf("RSymbol %s not supported", r)
	}
}
