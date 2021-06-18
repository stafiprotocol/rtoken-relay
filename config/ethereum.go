package config

import "math/big"

var (
	AmountBase = *big.NewInt(1000000000000000000) // 18 zero

	Call         uint8 = 0
	DelegateCall uint8 = 1
)
