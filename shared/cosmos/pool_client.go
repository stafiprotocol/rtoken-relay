package cosmos

import (
	"errors"
	"fmt"
	"sync"

	subTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
)

//one pool address with one poolClient
type PoolClient struct {
	eraSeconds       int64
	eraFactor        int64
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

func NewPoolClient(rpcClient *rpc.Client, subKey string, eraSeconds, eraFactor int64) *PoolClient {
	return &PoolClient{
		eraSeconds:       eraSeconds,
		eraFactor:        eraFactor,
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

func (pc *PoolClient) GetHeightByEra(era uint32) (int64, error) {

	targetTimestamp := (int64(era) + pc.eraFactor) * pc.eraSeconds
	blockNumber, timestamp, err := pc.GetRpcClient().GetCurrentBLockAndTimestamp()
	if err != nil {
		return 0, err
	}
	seconds := timestamp - targetTimestamp
	if seconds < 0 {
		return 0, fmt.Errorf("timestamp can not less than targetTimestamp")
	}

	tmpTargetBlock := blockNumber - seconds/7

	block, err := pc.GetRpcClient().QueryBlock(tmpTargetBlock)
	if err != nil {
		return 0, err
	}

	findDuTime := block.Block.Header.Time.Unix() - targetTimestamp

	if findDuTime == 0 {
		return block.Block.Height, nil
	}

	if findDuTime > 7 || findDuTime < -7 {
		tmpTargetBlock -= findDuTime / 7

		block, err = pc.GetRpcClient().QueryBlock(tmpTargetBlock)
		if err != nil {
			return 0, err
		}
	}

	var afterBlockNumber int64
	var preBlockNumber int64
	if block.Block.Header.Time.Unix() > targetTimestamp {
		afterBlockNumber = block.Block.Height

		for {
			block, err := pc.GetRpcClient().QueryBlock(afterBlockNumber - 1)
			if err != nil {
				return 0, err
			}
			if block.Block.Time.Unix() > targetTimestamp {
				afterBlockNumber = block.Block.Height
			} else {
				preBlockNumber = block.Block.Height
				break
			}
		}

	} else {
		preBlockNumber = block.Block.Height
		for {
			block, err := pc.GetRpcClient().QueryBlock(preBlockNumber + 1)
			if err != nil {
				return 0, err
			}
			if block.Block.Time.Unix() > targetTimestamp {
				afterBlockNumber = block.Block.Height
				break
			} else {
				preBlockNumber = block.Block.Height
			}
		}
	}

	return afterBlockNumber, nil
}

func (pc *PoolClient) GetCurrentEra() (uint32, error) {
	_, timestamp, err := pc.GetRpcClient().GetCurrentBLockAndTimestamp()
	if err != nil {
		return 0, err
	}

	if pc.eraSeconds <= 0 {
		panic("eraSeconds is zero")
	}
	era := timestamp/pc.eraSeconds - pc.eraFactor
	if era < 0 {
		panic("era can not less than zero")
	}

	return uint32(era), nil
}
