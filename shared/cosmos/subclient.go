package cosmos

import (
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"math/big"

	"github.com/ChainSafe/log15"
)

type SubClient struct {
	log       log15.Logger
	rpcClient *rpc.Client
	subKey    string
}

func NewSubClient(log log15.Logger, rpcClient *rpc.Client, subKey string) *SubClient {
	return &SubClient{log: log, rpcClient: rpcClient, subKey: subKey}
}

func (sc *SubClient) GetRpcClient() *rpc.Client {
	return sc.rpcClient
}

func (sc *SubClient) GetSubkey() string {
	return sc.subKey
}


func (sc *SubClient) bond(val *big.Int) error {
	return nil
}
func (sc *SubClient) unbond(val *big.Int) error {
	return nil
}
