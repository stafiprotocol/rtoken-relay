package conn

import (
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"math/big"
)

type Chain interface {
	TransferVerify(record *BondRecord) (BondReason, error)
	CurrentEra() (types.U32, error)
	BondWork(evtData *EvtEraPoolUpdated) (*big.Int, error)
}
