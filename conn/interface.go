package conn

import "github.com/stafiprotocol/go-substrate-rpc-client/types"

type Chain interface {
	TransferVerify(record *BondRecord) (BondReason, error)
	CurrentEra() (types.U32, error)
	BondWork(ck *ChunkKey) error
}
