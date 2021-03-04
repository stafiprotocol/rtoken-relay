package utils

import (
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
