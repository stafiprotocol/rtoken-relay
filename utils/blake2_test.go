package utils

import (
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlake2Hash(t *testing.T) {
	a := types.NewU32(100)
	h, err := Blake2Hash(a)
	assert.NoError(t, err)
	assert.Equal(t, h.Hex(), "0x4cd2c1c389ef163004139d22c65f8b340110ae9d21e7a78b04de446bb85958c1")

	a = types.U32(666)
	x, err := Blake2Hash(a)
	assert.NoError(t, err)
	assert.Equal(t, x.Hex(), "0x564d10a173f92fd4dcd3a4767fd44b40bcba4799f66fb32e429cb71568a0a04d")
}
