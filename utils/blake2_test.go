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
