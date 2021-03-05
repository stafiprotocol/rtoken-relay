package conn

import "github.com/stafiprotocol/go-substrate-rpc-client/types"

type EvtEraPoolUpdated struct {
	Symbol RSymbol
	NewEra types.U32
	Pool   types.Bytes
	Bond   types.U128
	Unbond types.U128
}
