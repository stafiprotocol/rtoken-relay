package substrate

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
)

func (l *listener) processLiquidityBondEvent(evt *substrate.ChainEvent) (*core.BondFlow, error) {
	evtData, err := liquidityBondEventData(evt)
	if err != nil {
		return nil, err
	}

	bondKey := &core.BondKey{Rsymbol: evtData.Rsymbol, BondId: evtData.BondId}
	bk, err := types.EncodeToBytes(bondKey)
	if err != nil {
		return nil, err
	}

	br := new(core.BondRecord)
	exist, err := l.conn.QueryStorage(config.LiquidityBondModuleId, config.StorageBondRecords, bk, nil, br)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("unable to get bondrecord by bondkey: %+v", bondKey)
	}

	l.log.Info("BondRecord", "bonder", hexutil.Encode(br.Bonder[:]), "rsymbol", br.Rsymbol,
		"pubkey", hexutil.Encode(br.Pubkey), "pool", hexutil.Encode(br.Pool), "blockhash", hexutil.Encode(br.Blockhash),
		"txhash", hexutil.Encode(br.Txhash), "amount", br.Amount.Int)

	if br.Bonder != evtData.AccountId {
		return nil, fmt.Errorf("bonder not matched: %s, %s", hexutil.Encode(br.Bonder[:]), hexutil.Encode(evtData.AccountId[:]))
	}

	return &core.BondFlow{Key: bondKey, Record: br, Reason: core.BondReasonDefault}, nil
}

func (l *listener) processEraPoolUpdatedEvt(evt *substrate.ChainEvent) (*core.EraPoolUpdatedFlow, error) {
	data, err := eraPoolUpdatedData(evt)
	if err != nil {
		return nil, err
	}

	pk := &core.PoolKey{Rsymbol: data.Rsymbol, Pool: data.Pool}
	pkBz, err := types.EncodeToBytes(pk)
	if err != nil {
		return nil, err
	}

	var threshold uint16
	exist, err := l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageMultiThresholds, pkBz, nil, &threshold)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("unable to get threshold of pool: %s, rsymbol: %s", pk.Rsymbol, hexutil.Encode(pk.Pool))
	}

	subs := make([]types.Bytes, 0)
	exist, err = l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageSubAccounts, pkBz, nil, &subs)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("unable to get subAccounts of pool: %s, rsymbol: %s", pk.Rsymbol, hexutil.Encode(pk.Pool))
	}

	return &core.EraPoolUpdatedFlow{
		Evt:           data,
		LastVoterFlag: l.conn.IsLastVoter(data.LastVoter),
		Threshold:     threshold,
		SubAccounts:   subs,
	}, nil
}
