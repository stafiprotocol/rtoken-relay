package cosmos

import (
	"errors"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"math/big"
	"sync"

	"github.com/ChainSafe/log15"
)

//one pool address with one poolClient
type PoolClient struct {
	log              log15.Logger
	rpcClient        *rpc.Client
	subKeyName       string //subKey is one of pubKeys of multiSig pool address,subKeyName is subKey`s name in keyring
	mtx              sync.RWMutex
	cachedUnsignedTx map[string][]byte //map[hash(unsignedTx)]unsignedTx
}

func NewSubClient(log log15.Logger, rpcClient *rpc.Client, subKey string) *PoolClient {
	return &PoolClient{
		log:              log,
		rpcClient:        rpcClient,
		subKeyName:       subKey,
		cachedUnsignedTx: make(map[string][]byte)}
}

func (pc *PoolClient) GetRpcClient() *rpc.Client {
	return pc.rpcClient
}

func (pc *PoolClient) GetSubKey() string {
	return pc.subKeyName
}
func (pc *PoolClient) CacheUnsignedTx(key string, tx []byte) {
	pc.mtx.Lock()
	pc.cachedUnsignedTx[key] = tx
	pc.mtx.Unlock()
}
func (pc *PoolClient) GetUnsignedTx(key string) ([]byte, error) {
	pc.mtx.RLock()
	defer pc.mtx.RUnlock()
	if tx, exist := pc.cachedUnsignedTx[key]; exist {
		return tx, nil
	}
	return nil, errors.New("unsignedTx of this key not exist")
}

func (pc *PoolClient) bond(val *big.Int) error {
	return nil
}
func (pc *PoolClient) unbond(val *big.Int) error {
	return nil
}
