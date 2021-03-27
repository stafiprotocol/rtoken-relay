package cosmos

import (
	"errors"
	subTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"math/big"
	"sync"

	"github.com/ChainSafe/log15"
)

//todo test only will update this on release
const EraBlockNumber = int64(10 * 60 / 6) //6hours 6*60*60/6

//one pool address with one poolClient
type PoolClient struct {
	log              log15.Logger
	rpcClient        *rpc.Client
	subKeyName       string //subKey is one of pubKeys of multiSig pool address,subKeyName is subKey`s name in keyring
	mtx              sync.RWMutex
	cachedUnsignedTx map[string]*WrapUnsignedTx //map[hash(unsignedTx)]unsignedTx
}

type WrapUnsignedTx struct {
	UnsignedTx []byte
	Hash       string
	SnapshotId subTypes.Hash
	Era        uint32
	Type       core.OriginalTx
}

func NewPoolClient(log log15.Logger, rpcClient *rpc.Client, subKey string) *PoolClient {
	return &PoolClient{
		log:              log,
		rpcClient:        rpcClient,
		subKeyName:       subKey,
		cachedUnsignedTx: make(map[string]*WrapUnsignedTx)}
}

func (pc *PoolClient) GetRpcClient() *rpc.Client {
	return pc.rpcClient
}

func (pc *PoolClient) GetSubKey() string {
	return pc.subKeyName
}
func (pc *PoolClient) CacheUnsignedTx(key string, tx *WrapUnsignedTx) {
	pc.mtx.Lock()
	pc.cachedUnsignedTx[key] = tx
	pc.mtx.Unlock()
}
func (pc *PoolClient) GetWrappedUnsignedTx(key string) (*WrapUnsignedTx, error) {
	pc.mtx.RLock()
	defer pc.mtx.RUnlock()
	if tx, exist := pc.cachedUnsignedTx[key]; exist {
		return tx, nil
	}
	return nil, errors.New("unsignedTx of this key not exist")
}

func (pc *PoolClient) RemoveUnsignedTx(key string) {
	pc.mtx.Lock()
	delete(pc.cachedUnsignedTx, key)
	pc.mtx.Unlock()
}

func (pc *PoolClient) GetHeightByEra(era uint32) int64 {
	return int64(era) * EraBlockNumber
}

func (pc *PoolClient) GetCurrentEra() (int64, uint32, error) {
	height, err := pc.GetRpcClient().GetCurrentBLockHeight()
	if err != nil {
		return 0, 0, err
	}
	era := uint32(height / EraBlockNumber)
	return height, era, nil
}
func (pc *PoolClient) bond(val *big.Int) error {
	return nil
}
func (pc *PoolClient) unbond(val *big.Int) error {
	return nil
}
