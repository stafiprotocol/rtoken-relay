package cosmos

import (
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"math/big"

	"github.com/ChainSafe/log15"
)

type SubClient struct {
	log       log15.Logger
	rpcClient *rpc.Client
	subKeys   []string
}

func NewSubClient(log log15.Logger, rpcClient *rpc.Client, subKeys []string) *SubClient {
	return &SubClient{log: log, rpcClient: rpcClient, subKeys: subKeys}
}

func (sc *SubClient) GetRpcClient() *rpc.Client {
	return sc.rpcClient
}

func (sc *SubClient) bond(val *big.Int) error {
	return nil
}
func (sc *SubClient) unbond(val *big.Int) error {
	return nil
}
