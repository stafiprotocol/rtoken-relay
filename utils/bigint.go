package utils

import (
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"math/big"
)

const Base = 10

func StringToBigint(src string) (*big.Int, bool) {
	return big.NewInt(0).SetString(src, Base)
}

func AddU128(a, b types.U128) types.U128 {
	c := big.NewInt(0).Add(a.Int, b.Int)
	return types.NewU128(*c)
}
