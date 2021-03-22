package cosmos

import (
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"math/big"

	"github.com/ChainSafe/log15"
)

//one pool address with one poolClient
type PoolClient struct {
	log        log15.Logger
	rpcClient  *rpc.Client
	subKeyName string //subKey is one of pubKeys of multiSig pool address,subKeyName is subKey`s name in keyring
}

func NewSubClient(log log15.Logger, rpcClient *rpc.Client, subKey string) *PoolClient {
	return &PoolClient{log: log, rpcClient: rpcClient, subKeyName: subKey}
}

func (pc *PoolClient) GetRpcClient() *rpc.Client {
	return pc.rpcClient
}

func (pc *PoolClient) GetSubKey() string {
	return pc.subKeyName
}

func (pc *PoolClient) bond(val *big.Int) error {
	return nil
}
func (pc *PoolClient) unbond(val *big.Int) error {
	return nil
}
