package utils

import "math/big"

const Base = 10

func StringToBigint(src string) (*big.Int, bool) {
	return big.NewInt(0).SetString(src, Base)
}
