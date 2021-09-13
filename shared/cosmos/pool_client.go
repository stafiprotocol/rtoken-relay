package cosmos

import (
	"errors"
	subTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"sync"

	"github.com/ChainSafe/log15"
)

//one pool address with one poolClient
type PoolClient struct {
	eraBlockNumber   int64
	log              log15.Logger
	rpcClient        *rpc.Client
	subKeyName       string //subKey is one of pubKeys of multiSig pool address,subKeyName is subKey`s name in keyring
	mtx              sync.RWMutex
	cachedUnsignedTx map[string]*WrapUnsignedTx //map[hash(unsignedTx)]unsignedTx
}

type WrapUnsignedTx struct {
	UnsignedTx []byte
	Key        string
	SnapshotId subTypes.Hash
	Era        uint32
	Bond       subTypes.U128
	Unbond     subTypes.U128
	Type       submodel.OriginalTx
}

func NewPoolClient(log log15.Logger, rpcClient *rpc.Client, subKey string, eraBLockNumber int64) *PoolClient {
	return &PoolClient{
		eraBlockNumber:   eraBLockNumber,
		log:              log,
		rpcClient:        rpcClient,
		subKeyName:       subKey,
		cachedUnsignedTx: make(map[string]*WrapUnsignedTx)}
}

func (pc *PoolClient) GetRpcClient() *rpc.Client {
	return pc.rpcClient
}

func (pc *PoolClient) GetSubKeyName() string {
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

func (pc *PoolClient) CachedUnsignedTxNumber() int {
	return len(pc.cachedUnsignedTx)
}

func (pc *PoolClient) GetHeightByEra(era uint32) int64 {
	return int64(era) * pc.eraBlockNumber
}

func (pc *PoolClient) GetCurrentEra() (int64, uint32, error) {
	height, err := pc.GetRpcClient().GetCurrentBLockHeight()
	if err != nil {
		return 0, 0, err
	}
	if pc.eraBlockNumber == 0 {
		panic("eraBlockNumber is zero")
	}
	era := uint32(height / pc.eraBlockNumber)
	return height, era, nil
}
