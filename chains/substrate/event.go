package substrate

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
)

var multiEndError = errors.New("multiEnd")

func (l *listener) processLiquidityBondEvent(evt *substrate.ChainEvent) (*core.BondFlow, error) {
	evtData, err := substrate.LiquidityBondEventData(evt)
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

func (l *listener) processEraPoolUpdatedEvt(evt *substrate.ChainEvent) (*core.MultisigFlow, error) {
	data, err := substrate.EraPoolUpdatedData(evt)
	if err != nil {
		return nil, err
	}

	snap, err := l.snapshot(data.ShotId)
	if err != nil {
		return nil, err
	}

	mc, err := l.multiCall(snap.Rsymbol, snap.Pool)
	if err != nil {
		return nil, err
	}

	data.LastVoterFlag = l.conn.IsLastVoter(data.LastVoter)
	data.Snap = snap

	return &core.MultisigFlow{
		HeadFlow: data,
		MulCall:  mc,
	}, nil
}

func (l *listener) processNewMultisigEvt(evt *substrate.ChainEvent) (*core.MultisigFlow, error) {
	data, err := substrate.EventNewMultisig(evt)
	if err != nil {
		return nil, err
	}

	mul := new(core.Multisig)
	exist, err := l.conn.QueryStorage(config.MultisigModuleId, config.StorageMultisigs, data.ID[:], data.CallHash[:], mul)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, multiEndError
	}

	return &core.MultisigFlow{
		MulCall:  &core.MultisigCall{TimePoint: core.NewOptionTimePoint(mul.When)},
		CallHash: hexutil.Encode(data.CallHash[:]),
		NewMul:   data,
		Multisig: mul,
	}, nil
}

func (l *listener) processMultisigExecutedEvt(evt *substrate.ChainEvent) (*core.MultisigFlow, error) {
	data, err := substrate.EventMultisigExecuted(evt)
	if err != nil {
		return nil, err
	}

	return &core.MultisigFlow{
		CallHash:    hexutil.Encode(data.CallHash[:]),
		MulExecuted: data,
	}, nil
}

func (l *listener) processBondReportEvt(evt *substrate.ChainEvent) (*core.BondReportFlow, error) {
	flow, err := substrate.EventBondReport(evt)
	if err != nil {
		return nil, err
	}

	flow.LastVoterFlag = l.conn.IsLastVoter(flow.LastVoter)
	flow.LastEra = flow.Era - 1

	return flow, nil
}

func (l *listener) processWithdrawUnbondEvt(evt *substrate.ChainEvent) (*core.MultisigFlow, error) {
	data, err := substrate.EventWithdrawUnbond(evt)
	if err != nil {
		return nil, err
	}

	snap, err := l.snapshot(data.ShotId)
	if err != nil {
		return nil, err
	}

	mc, err := l.multiCall(snap.Rsymbol, snap.Pool)
	if err != nil {
		return nil, err
	}

	data.LastVoterFlag = l.conn.IsLastVoter(data.LastVoter)

	return &core.MultisigFlow{
		HeadFlow: data,
		MulCall:  mc,
	}, nil
}

func (l *listener) snapshot(shotId types.Hash) (*core.EraPoolSnapshot, error) {
	bz, err := types.EncodeToBytes(shotId)
	snap := new(core.EraPoolSnapshot)
	exist, err := l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageSnapshots, bz, nil, snap)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("unable to get snap of shotId: %s", hexutil.Encode(shotId[:]))
	}

	return snap, nil
}

func (l *listener) multiCall(symbol core.RSymbol, pool []byte) (*core.MultisigCall, error) {
	pk := &core.PoolKey{Rsymbol: symbol, Pool: pool}
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
		return nil, fmt.Errorf("unable to get threshold of pool: %s, rsymbol: %s", symbol, hexutil.Encode(pool))
	}

	subs := make([]types.Bytes, 0)
	exist, err = l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageSubAccounts, pkBz, nil, &subs)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("unable to get subAccounts of pool: %s, rsymbol: %s", symbol, hexutil.Encode(pool))
	}

	return &core.MultisigCall{
		Threshold:   threshold,
		SubAccounts: subs,
	}, nil
}
