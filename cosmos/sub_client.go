package cosmos

import (
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/cosmos/rpc"
	"math/big"
)

type SubClient struct {
	Log       log15.Logger
	RpcClient rpc.Client
}

func (sc *SubClient) bond(val *big.Int) error {
	return nil
}
func (sc *SubClient) unbond(val *big.Int) error {
	return nil
}
