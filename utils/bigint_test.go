package utils

import (
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

var oneEth = big.NewInt(1000000000000000000)

func TestFromString(t *testing.T) {
	a := "32000000000000000000"
	b, ok := FromString(a)
	assert.Equal(t, true, ok)

	x := big.NewInt(32)
	x.Mul(x, oneEth)
	assert.Equal(t, 0, b.Cmp(x))
}

func TestBlake2(t *testing.T) {
	a := types.U32(666)
	x, err := Blake2Hash(a)
	assert.NoError(t, err)
	assert.Equal(t, x.Hex(), "0x2fbd238e1c2a4ac8aa0f961daa5772b34ccb202cb9c499ac18d47b1434cb0b99")
}
