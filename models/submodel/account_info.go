package submodel

import "github.com/stafiprotocol/go-substrate-rpc-client/types"

type AccountInfo struct {
	Nonce     uint32
	Consumers uint32
	Providers uint32
	Data      struct {
		Free       types.U128
		Reserved   types.U128
		MiscFrozen types.U128
		FreeFrozen types.U128
	}
}
