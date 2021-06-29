package config

import "math/big"

var (
	AmountBase = *big.NewInt(1000000000000000000) // 18 zero

)

type CallEnum uint8

var (
	Call         = CallEnum(0)
	DelegateCall = CallEnum(1)
)

type TxHashState uint8

var (
	HashStateUnsubmit = TxHashState(0)
	HashStateFail     = TxHashState(1)
	HashStateSuccess  = TxHashState(2)
)
