package substrate

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

var (
	multiEndError = errors.New("multiEnd")

	EventEraIsOldError                = errors.New("EventEraIsOldError")
	BondStateNotEraUpdatedError       = errors.New("BondStateNotEraUpdatedError")
	BondStateNotBondReportedError     = errors.New("BondStateNotBondReportedError")
	BondStateNotActiveReportedError   = errors.New("BondStateNotActiveReportedError")
	BondStateNotWithdrawReportedError = errors.New("BondStateNotWithdrawReportedError")
	BondStateNotTransferReportedError = errors.New("BondStateNotTransferReportedError")
)

func (l *listener) processLiquidityBondEvent(evt *submodel.ChainEvent) (*submodel.BondFlow, error) {
	data, err := submodel.LiquidityBondEventData(evt)
	if err != nil {
		return nil, err
	}

	symBz, err := types.EncodeToBytes(data.Symbol)
	if err != nil {
		return nil, err
	}

	br := new(submodel.BondRecord)
	exist, err := l.conn.QueryStorage(config.RTokenSeriesModuleId, config.StorageBondRecords, symBz, data.BondId[:], br)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("failed to get bondrecord, symbol: %s, bondId: %s", data.Symbol, data.BondId.Hex())
	}

	bsk := submodel.BondStateKey{BlockHash: br.Blockhash, TxHash: br.Txhash}
	bskBz, err := types.EncodeToBytes(bsk)
	if err != nil {
		return nil, err
	}

	var bs submodel.BondState
	exist, err = l.conn.QueryStorage(config.RTokenSeriesModuleId, config.StorageBondStates, symBz, bskBz, &bs)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("failed to get bondstate, symbol: %s, bondId: %s, BlockHash: %s, TxHash: %s",
			data.Symbol, data.BondId.Hex(), hexutil.Encode(br.Blockhash), hexutil.Encode(br.Txhash))
	}

	l.log.Info("BondRecord", "bonder", hexutil.Encode(br.Bonder[:]), "symbol", br.Symbol,
		"pubkey", hexutil.Encode(br.Pubkey), "pool", hexutil.Encode(br.Pool), "blockHash", hexutil.Encode(br.Blockhash),
		"txHash", hexutil.Encode(br.Txhash), "amount", br.Amount.Int, "BondState", bs)

	if br.Bonder != data.AccountId {
		return nil, fmt.Errorf("bonder not matched: %s, %s", hexutil.Encode(br.Bonder[:]), hexutil.Encode(data.AccountId[:]))
	}

	return &submodel.BondFlow{
		Symbol: data.Symbol,
		BondId: data.BondId,
		Record: br,
		Reason: submodel.BondReasonDefault,
		State:  bs,
	}, nil
}

func (l *listener) processEraPoolUpdatedEvt(evt *submodel.ChainEvent) (*submodel.MultiEventFlow, error) {
	data, err := submodel.EraPoolUpdatedData(evt)
	if err != nil {
		return nil, err
	}

	snap, err := l.snapshot(data.Symbol, data.ShotId)
	if err != nil {
		return nil, err
	}

	curEra, err := l.conn.CurrentChainEra(snap.Symbol)
	if err != nil {
		if err.Error() != fmt.Sprintf("era of rsymbol %s not exist", snap.Symbol) {
			l.log.Error("failed to get CurrentChainEra", "error", err)
			return nil, err
		}
	}
	if curEra != snap.Era {
		return nil, EventEraIsOldError
	}

	if snap.BondState != submodel.EraUpdated {
		l.log.Warn("processEraPoolUpdatedEvt: bondState not EraUpdated",
			"symbol", snap.Symbol, "pool", hexutil.Encode(snap.Pool), "BondState", snap.BondState)
		return nil, BondStateNotEraUpdatedError
	}

	th, sub, err := l.thresholdAndSubAccounts(snap.Symbol, snap.Pool)
	if err != nil {
		return nil, err
	}

	data.LastVoterFlag = l.conn.IsLastVoter(data.LastVoter)
	data.Snap = snap

	return &submodel.MultiEventFlow{
		EventId:     config.EraPoolUpdatedEventId,
		Symbol:      snap.Symbol,
		EventData:   data,
		Threshold:   th,
		SubAccounts: sub,
	}, nil
}

func (l *listener) processBondReportedEvt(evt *submodel.ChainEvent) (*submodel.BondReportedFlow, error) {
	flow, err := submodel.EventBondReported(evt)
	if err != nil {
		return nil, err
	}

	snap, err := l.snapshot(flow.Symbol, flow.ShotId)
	if err != nil {
		return nil, err
	}

	curEra, err := l.conn.CurrentChainEra(snap.Symbol)
	if err != nil {
		if err.Error() != fmt.Sprintf("era of rsymbol %s not exist", snap.Symbol) {
			l.log.Error("failed to get CurrentChainEra", "error", err)
			return nil, err
		}
	}
	if curEra != snap.Era {
		return nil, EventEraIsOldError
	}

	if snap.BondState != submodel.BondReported {
		l.log.Warn("processBondReportedEvt: bondState not BondReported",
			"symbol", snap.Symbol, "pool", hexutil.Encode(snap.Pool), "BondState", snap.BondState)
		return nil, BondStateNotBondReportedError
	}

	_, sub, err := l.thresholdAndSubAccounts(snap.Symbol, snap.Pool)
	if err != nil {
		return nil, err
	}

	flow.LastVoterFlag = l.conn.IsLastVoter(flow.LastVoter)
	flow.Snap = snap
	flow.LastEra = snap.Era - 1
	flow.SubAccounts = sub

	return flow, nil
}

func (l *listener) processActiveReportedEvt(evt *submodel.ChainEvent) (*submodel.MultiEventFlow, error) {
	flow, err := submodel.EventActiveReported(evt)
	if err != nil {
		return nil, err
	}

	snap, err := l.snapshot(flow.Symbol, flow.ShotId)
	if err != nil {
		return nil, err
	}

	curEra, err := l.conn.CurrentChainEra(snap.Symbol)
	if err != nil {
		if err.Error() != fmt.Sprintf("era of rsymbol %s not exist", snap.Symbol) {
			l.log.Error("failed to get CurrentChainEra", "error", err)
			return nil, err
		}
	}
	if curEra != snap.Era {
		return nil, EventEraIsOldError
	}

	if snap.BondState != submodel.ActiveReported {
		l.log.Warn("processActiveReportedEvt: bondState not ActiveReported",
			"symbol", snap.Symbol, "pool", hexutil.Encode(snap.Pool), "BondState", snap.BondState)
		return nil, BondStateNotActiveReportedError
	}

	th, sub, err := l.thresholdAndSubAccounts(snap.Symbol, snap.Pool)
	if err != nil {
		return nil, err
	}

	flow.LastVoterFlag = l.conn.IsLastVoter(flow.LastVoter)
	flow.Snap = snap

	return &submodel.MultiEventFlow{
		EventId:     config.ActiveReportedEventId,
		Symbol:      snap.Symbol,
		EventData:   flow,
		Threshold:   th,
		SubAccounts: sub,
	}, nil
}

func (l *listener) processWithdrawReportedEvt(evt *submodel.ChainEvent) (*submodel.MultiEventFlow, error) {
	flow, err := submodel.EventWithdrawReported(evt)
	if err != nil {
		return nil, err
	}

	snap, err := l.snapshot(flow.Symbol, flow.ShotId)
	if err != nil {
		return nil, err
	}

	curEra, err := l.conn.CurrentChainEra(snap.Symbol)
	if err != nil {
		if err.Error() != fmt.Sprintf("era of rsymbol %s not exist", snap.Symbol) {
			l.log.Error("failed to get CurrentChainEra", "error", err)
			return nil, err
		}
	}
	if curEra != snap.Era {
		return nil, EventEraIsOldError
	}

	if snap.BondState != submodel.WithdrawReported {
		l.log.Warn("processWithdrawReportedEvt: bondState not WithdrawReported",
			"symbol", snap.Symbol, "pool", hexutil.Encode(snap.Pool), "BondState", snap.BondState)
		return nil, BondStateNotWithdrawReportedError
	}

	receives, total, err := l.unbondings(snap.Symbol, snap.Pool, snap.Era)
	if err != nil {
		return nil, err
	}

	th, sub, err := l.thresholdAndSubAccounts(snap.Symbol, snap.Pool)
	if err != nil {
		return nil, err
	}

	flow.LastVoterFlag = l.conn.IsLastVoter(flow.LastVoter)
	flow.Snap = snap
	flow.Receives = receives
	flow.TotalAmount = total

	return &submodel.MultiEventFlow{
		EventId:     config.WithdrawReportedEventId,
		Symbol:      snap.Symbol,
		EventData:   flow,
		Threshold:   th,
		SubAccounts: sub,
	}, nil
}

func (l *listener) processTransferReportedEvt(evt *submodel.ChainEvent) (*submodel.TransferReportedFlow, error) {
	flow, err := submodel.EventTransferReported(evt)
	if err != nil {
		return nil, err
	}

	snap, err := l.snapshot(flow.Symbol, flow.ShotId)
	if err != nil {
		return nil, err
	}

	if snap.BondState != submodel.TransferReported {
		l.log.Warn("processTransferReportedEvt: bondState not TransferReported",
			"symbol", snap.Symbol, "pool", hexutil.Encode(snap.Pool), "BondState", snap.BondState)
		return nil, BondStateNotWithdrawReportedError
	}

	receives, total, err := l.unbondings(snap.Symbol, snap.Pool, snap.Era)
	if err != nil {
		return nil, err
	}

	flow.Snap = snap
	flow.Receives = receives
	flow.TotalAmount = total

	return flow, nil
}

func (l *listener) processNominationUpdated(evt *submodel.ChainEvent) (*submodel.MultiEventFlow, error) {
	flow, err := submodel.EventNominationUpdated(evt)
	if err != nil {
		return nil, err
	}

	th, sub, err := l.thresholdAndSubAccounts(flow.Symbol, flow.Pool)
	if err != nil {
		return nil, err
	}

	flow.LastVoterFlag = l.conn.IsLastVoter(flow.LastVoter)

	return &submodel.MultiEventFlow{
		EventId:     config.NominationUpdatedEventId,
		Symbol:      flow.Symbol,
		EventData:   flow,
		Threshold:   th,
		SubAccounts: sub,
	}, nil
}

func (l *listener) processNewMultisigEvt(evt *submodel.ChainEvent) (*submodel.EventNewMultisig, error) {
	data, err := submodel.EventNewMultisigData(evt)
	if err != nil {
		return nil, err
	}

	mul := new(submodel.Multisig)
	exist, err := l.conn.QueryStorage(config.MultisigModuleId, config.StorageMultisigs, data.ID[:], data.CallHash[:], mul)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, multiEndError
	}

	data.TimePoint = submodel.NewOptionTimePoint(mul.When)
	data.Approvals = mul.Approvals
	data.CallHashStr = hexutil.Encode(data.CallHash[:])
	return data, nil
}

func (l *listener) processMultisigExecutedEvt(evt *submodel.ChainEvent) (*submodel.EventMultisigExecuted, error) {
	data, err := submodel.EventMultisigExecutedData(evt)
	if err != nil {
		return nil, err
	}
	data.CallHashStr = hexutil.Encode(data.CallHash[:])
	return data, nil
}

func (l *listener) snapshot(symbol core.RSymbol, shotId types.Hash) (*submodel.PoolSnapshot, error) {
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return nil, err
	}
	bz, err := types.EncodeToBytes(shotId)
	if err != nil {
		return nil, err
	}
	snap := new(submodel.PoolSnapshot)
	exist, err := l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageSnapshots, symBz, bz, snap)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("snap of shotId: %s not exist", hexutil.Encode(shotId[:]))
	}

	return snap, nil
}

func (l *listener) thresholdAndSubAccounts(symbol core.RSymbol, pool []byte) (uint16, []types.Bytes, error) {
	poolBz, err := types.EncodeToBytes(pool)
	if err != nil {
		return 0, nil, err
	}
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return 0, nil, err
	}

	var threshold uint16
	exist, err := l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageMultiThresholds, symBz, poolBz, &threshold)
	if err != nil {
		return 0, nil, err
	}
	if !exist {
		return 0, nil, fmt.Errorf("threshold of pool: %s, symbol: %s not exist", symbol, hexutil.Encode(pool))
	}

	subs := make([]types.Bytes, 0)
	exist, err = l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageSubAccounts, symBz, poolBz, &subs)
	if err != nil {
		return 0, nil, err
	}
	if !exist {
		return 0, nil, fmt.Errorf("subAccounts of pool: %s, symbol: %s not exist", symbol, hexutil.Encode(pool))
	}

	return threshold, subs, nil
}

func (l *listener) unbondings(symbol core.RSymbol, pool []byte, era uint32) ([]*submodel.Receive, types.U128, error) {
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return nil, types.U128{}, err
	}

	puk := &submodel.PoolUnbondKey{Pool: pool, Era: era}
	pkbz, err := types.EncodeToBytes(puk)
	if err != nil {
		return nil, types.U128{}, err
	}

	unbonds := make([]submodel.Unbonding, 0)
	exist, err := l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StoragePoolUnbonds, symBz, pkbz, &unbonds)
	if err != nil {
		return nil, types.U128{}, err
	}
	if !exist {
		return nil, types.U128{}, fmt.Errorf("pool unbonds not exist, symbol: %s, pool: %s, era: %d", symbol, hexutil.Encode(pool), era)
	}

	amounts := make(map[string]types.U128)
	for _, ub := range unbonds {
		rec := hexutil.Encode(ub.Recipient)
		acc, ok := amounts[rec]
		if !ok {
			amounts[rec] = ub.Value
		} else {
			amounts[rec] = utils.AddU128(acc, ub.Value)
		}
	}

	receives := make([]*submodel.Receive, 0)
	total := types.NewU128(*big.NewInt(0))
	for k, v := range amounts {
		r, err := hexutil.Decode(k)
		if err != nil {
			return nil, types.U128{}, fmt.Errorf("hexutil.Decode err %s,k: %v", err, k)
		}
		rec := &submodel.Receive{Recipient: r, Value: types.NewUCompact(v.Int)}
		receives = append(receives, rec)
		total = utils.AddU128(total, v)
	}

	return receives, total, nil
}
