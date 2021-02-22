package utils

import "math/big"

const Base = 10

func FromString(src string) (*big.Int, bool) {
	return big.NewInt(0).SetString(src, Base)
}
