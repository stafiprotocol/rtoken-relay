package core

import (
	"testing"

	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

func TestRSymbol_Encode(t *testing.T) {
	a := RDOT
	enc, err := types.EncodeToHexString(a)
	if err != nil {
		panic(err)
	}

	var r RSymbol
	err = types.DecodeFromHexString(enc, &r)
	if err != nil {
		panic(err)
	}

}
