package utils

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlakeTwo256(t *testing.T) {
	a := "0x736466736466"
	b, err := hexutil.Decode(a)
	assert.NoError(t, err)

	x := BlakeTwo256(b)
	assert.Equal(t, hexutil.Encode(x[:]), "0x66af11336fceaf1aa1b37b9fff097e744e5b10f488d2f60f2608a4c4eb878b9e")

}

func TestBlakeTwo256_1(t *testing.T) {
	a := "0x08010b00a0724e1809"
	b, err := hexutil.Decode(a)
	assert.NoError(t, err)

	x := BlakeTwo256(b)
	assert.Equal(t, hexutil.Encode(x[:]), "0xba6c8ec1798285f8f312523e2353ebe8468fab4b55afe1a788a64a65f8bcc72c")
}
