package substrate

import (
	"errors"
)

const (
	ChainTypeStafi    = "stafi"
	ChainTypePolkadot = "polkadot"

	AddressTypeAccountId    = "AccountId"
	AddressTypeMultiAddress = "MultiAddress"
)

var (
	TerminatedError           = errors.New("terminated")
	BondEqualToUnbondError    = errors.New("BondEqualToUnbondError")
	BondSmallerThanLeastError = errors.New("BondSmallerThanLeastError")
)
