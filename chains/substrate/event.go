package substrate

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common/hexutil"
	bncCmnTypes "github.com/stafiprotocol/go-sdk/common/types"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

var (
	multiEndError             = errors.New("multiEnd")
	eraSnapShotsNotExistError = errors.New("eraSnapShots not exist")

	EventEraIsOldError                = errors.New("EventEraIsOldError")
	BondStateNotEraUpdatedError       = errors.New("BondStateNotEraUpdatedError")
	BondStateNotBondReportedError     = errors.New("BondStateNotBondReportedError")
	BondStateNotActiveReportedError   = errors.New("BondStateNotActiveReportedError")
	BondStateNotWithdrawReportedError = errors.New("BondStateNotWithdrawReportedError")
	BondStateNotTransferReportedError = errors.New("BondStateNotTransferReportedError")
	ErrNotCared                       = errors.New("not care this symbol")
)

func (l *listener) processLiquidityBondEvent(evt *submodel.ChainEvent) (*submodel.BondFlow, error) {
	data, err := submodel.LiquidityBondEventData(evt)
	if err != nil {
		return nil, err
	}
	if !l.cared(data.Symbol) {
		return nil, ErrNotCared
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

	if l.cared(br.Symbol) {
		l.log.Info("BondRecord", "bonder", hexutil.Encode(br.Bonder[:]), "symbol", br.Symbol,
			"pubkey", hexutil.Encode(br.Pubkey), "pool", hexutil.Encode(br.Pool), "blockHash", hexutil.Encode(br.Blockhash),
			"txHash", hexutil.Encode(br.Txhash), "amount", br.Amount.Int, "BondState", bs)
	}

	if br.Bonder != data.AccountId {
		return nil, fmt.Errorf("bonder not matched: %s, %s", hexutil.Encode(br.Bonder[:]), hexutil.Encode(data.AccountId[:]))
	}

	return &submodel.BondFlow{
		Symbol:      data.Symbol,
		BondId:      data.BondId,
		Record:      br,
		Reason:      submodel.BondReasonDefault,
		State:       bs,
		VerifyTimes: 0,
	}, nil
}

func (l *listener) processEraPoolUpdatedEvt(evt *submodel.ChainEvent) (*submodel.MultiEventFlow, error) {
	data, err := submodel.EraPoolUpdatedData(evt)
	if err != nil {
		return nil, err
	}

	if !l.cared(data.Symbol) {
		return nil, ErrNotCared
	}

	snap, err := l.snapshot(data.Symbol, data.ShotId)
	if err != nil {
		return nil, err
	}

	curEra, err := l.conn.CurrentChainEra(snap.Symbol)
	if err != nil {
		if err.Error() != fmt.Sprintf("era of symbol %s not exist", snap.Symbol) {
			if l.cared(snap.Symbol) {
				l.log.Error("failed to get CurrentChainEra", "error", err)
			}
			return nil, err
		}
	}
	if curEra != snap.Era {
		return nil, EventEraIsOldError
	}

	if snap.BondState != submodel.EraUpdated {
		if l.cared(snap.Symbol) {
			l.log.Warn("processEraPoolUpdatedEvt: bondState not EraUpdated",
				"symbol", snap.Symbol, "pool", hexutil.Encode(snap.Pool), "BondState", snap.BondState)
		}
		return nil, BondStateNotEraUpdatedError
	}

	th, sub, err := l.thresholdAndSubAccounts(snap.Symbol, snap.Pool)
	if err != nil {
		return nil, err
	}

	var validatorId interface{}
	var leastBond *big.Int
	if data.Symbol == core.RMATIC || data.Symbol == core.RBNB {
		validatorId, err = l.validatorId(snap.Symbol, snap.Pool)
		if err != nil {
			return nil, err
		}

		leastBond, err = l.leastBond(snap.Symbol)
		if err != nil {
			return nil, err
		}
	}

	data.LastVoterFlag = l.conn.IsLastVoter(data.LastVoter)
	data.Snap = snap
	data.LeastBond = leastBond

	return &submodel.MultiEventFlow{
		EventId:     config.EraPoolUpdatedEventId,
		Symbol:      snap.Symbol,
		EventData:   data,
		Threshold:   th,
		SubAccounts: sub,
		ValidatorId: validatorId,
	}, nil
}

func (l *listener) processBondReportedEvt(evt *submodel.ChainEvent) (*submodel.BondReportedFlow, error) {
	flow, err := submodel.EventBondReported(evt)
	if err != nil {
		return nil, err
	}

	if !l.cared(flow.Symbol) {
		return nil, ErrNotCared
	}

	snap, err := l.snapshot(flow.Symbol, flow.ShotId)
	if err != nil {
		return nil, err
	}

	curEra, err := l.conn.CurrentChainEra(snap.Symbol)
	if err != nil {
		if err.Error() != fmt.Sprintf("era of symbol %s not exist", snap.Symbol) {
			if l.cared(snap.Symbol) {
				l.log.Error("failed to get CurrentChainEra", "error", err)
			}
			return nil, err
		}
	}
	if curEra != snap.Era {
		return nil, EventEraIsOldError
	}

	if snap.BondState != submodel.BondReported {
		if l.cared(snap.Symbol) {
			l.log.Warn("processBondReportedEvt: bondState not BondReported",
				"symbol", snap.Symbol, "pool", hexutil.Encode(snap.Pool), "BondState", snap.BondState)
		}
		return nil, BondStateNotBondReportedError
	}

	_, sub, err := l.thresholdAndSubAccounts(snap.Symbol, snap.Pool)
	if err != nil {
		return nil, err
	}

	var validatorId interface{}
	var leastBond *big.Int
	if snap.Symbol == core.RMATIC || snap.Symbol == core.RBNB {
		validatorId, err = l.validatorId(snap.Symbol, snap.Pool)
		if err != nil {
			return nil, err
		}

		leastBond, err = l.leastBond(snap.Symbol)
		if err != nil {
			return nil, err
		}
	}

	flow.LastEra = snap.Era - 1
	if snap.Symbol == core.RBNB {
		_, err = l.eraSnapShots(snap.Symbol, flow.LastEra)
		if err != nil {
			if err.Error() != eraSnapShotsNotExistError.Error() {
				return nil, err
			}
			flow.LastEra = 0
		}
	}

	flow.LastVoterFlag = l.conn.IsLastVoter(flow.LastVoter)
	flow.Snap = snap
	flow.SubAccounts = sub
	flow.ValidatorId = validatorId
	flow.LeastBond = leastBond

	return flow, nil
}

func (l *listener) processActiveReportedEvt(evt *submodel.ChainEvent) (*submodel.MultiEventFlow, error) {
	flow, err := submodel.EventActiveReported(evt)
	if err != nil {
		return nil, err
	}
	if !l.cared(flow.Symbol) {
		return nil, ErrNotCared
	}

	//turn to processActiveReportedEvtAsWithdrawReported
	if flow.Symbol == core.RATOM || flow.Symbol == core.RBNB {
		return l.processActiveReportedEvtAsWithdrawReported(evt)
	}

	snap, err := l.snapshot(flow.Symbol, flow.ShotId)
	if err != nil {
		return nil, err
	}

	curEra, err := l.conn.CurrentChainEra(snap.Symbol)
	if err != nil {
		if err.Error() != fmt.Sprintf("era of symbol %s not exist", snap.Symbol) {
			if l.cared(snap.Symbol) {
				l.log.Error("failed to get CurrentChainEra", "error", err)
			}
			return nil, err
		}
	}
	if curEra != snap.Era {
		return nil, EventEraIsOldError
	}

	if snap.BondState != submodel.ActiveReported {
		if l.cared(snap.Symbol) {
			l.log.Warn("processActiveReportedEvt: bondState not ActiveReported",
				"symbol", snap.Symbol, "pool", hexutil.Encode(snap.Pool), "BondState", snap.BondState)
		}
		return nil, BondStateNotActiveReportedError
	}

	th, sub, err := l.thresholdAndSubAccounts(snap.Symbol, snap.Pool)
	if err != nil {
		return nil, err
	}

	var valdatorId interface{}
	if snap.Symbol == core.RMATIC {
		valdatorId, err = l.validatorId(snap.Symbol, snap.Pool)
		if err != nil {
			return nil, err
		}
	}

	flow.LastVoterFlag = l.conn.IsLastVoter(flow.LastVoter)
	flow.Snap = snap

	return &submodel.MultiEventFlow{
		EventId:     config.ActiveReportedEventId,
		Symbol:      snap.Symbol,
		EventData:   flow,
		Threshold:   th,
		SubAccounts: sub,
		ValidatorId: valdatorId,
	}, nil
}

func (l *listener) processActiveReportedEvtAsWithdrawReported(evt *submodel.ChainEvent) (*submodel.MultiEventFlow, error) {
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
		if err.Error() != fmt.Sprintf("era of symbol %s not exist", snap.Symbol) {
			if l.cared(snap.Symbol) {
				l.log.Error("failed to get CurrentChainEra", "error", err)
			}
			return nil, err
		}
	}
	if curEra != snap.Era {
		return nil, EventEraIsOldError
	}

	if snap.BondState != submodel.ActiveReported {
		if l.cared(snap.Symbol) {
			l.log.Warn("processWithdrawReportedEvt: bondState not WithdrawReported",
				"symbol", snap.Symbol, "pool", hexutil.Encode(snap.Pool), "BondState", snap.BondState)
		}
		return nil, BondStateNotActiveReportedError
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
	if !l.cared(flow.Symbol) {
		return nil, ErrNotCared
	}
	snap, err := l.snapshot(flow.Symbol, flow.ShotId)
	if err != nil {
		return nil, err
	}

	curEra, err := l.conn.CurrentChainEra(snap.Symbol)
	if err != nil {
		if err.Error() != fmt.Sprintf("era of symbol %s not exist", snap.Symbol) {
			if l.cared(snap.Symbol) {
				l.log.Error("failed to get CurrentChainEra", "error", err)
			}
			return nil, err
		}
	}
	if curEra != snap.Era {
		return nil, EventEraIsOldError
	}

	if snap.BondState != submodel.WithdrawReported {
		if l.cared(snap.Symbol) {
			l.log.Warn("processWithdrawReportedEvt: bondState not WithdrawReported",
				"symbol", snap.Symbol, "pool", hexutil.Encode(snap.Pool), "BondState", snap.BondState)
		}
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

	var valdatorId interface{}
	if snap.Symbol == core.RMATIC {
		valdatorId, err = l.validatorId(snap.Symbol, snap.Pool)
		if err != nil {
			return nil, err
		}
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
		ValidatorId: valdatorId,
	}, nil
}

func (l *listener) processTransferReportedEvt(evt *submodel.ChainEvent) (*submodel.TransferReportedFlow, error) {
	flow, err := submodel.EventTransferReported(evt)
	if err != nil {
		return nil, err
	}
	if !l.cared(flow.Symbol) {
		return nil, ErrNotCared
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
	if !l.cared(flow.Symbol) {
		return nil, ErrNotCared
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

func (l *listener) processValidatorUpdated(evt *submodel.ChainEvent) (*submodel.MultiEventFlow, error) {
	flow, err := submodel.EventValidatorUpdated(evt)
	if err != nil {
		return nil, err
	}
	if !l.cared(flow.Symbol) {
		return nil, ErrNotCared
	}

	return &submodel.MultiEventFlow{
		EventId:   config.ValidatorUpdatedEventId,
		Symbol:    flow.Symbol,
		EventData: flow,
	}, nil
}

func (l *listener) processSignatureEnoughEvt(evt *submodel.ChainEvent) (*submodel.SubmitSignatures, error) {
	data, err := submodel.SignatureEnoughData(evt)
	if err != nil {
		return nil, err
	}
	if !l.cared(data.RSymbol) {
		return nil, ErrNotCared
	}

	symBz, err := types.EncodeToBytes(data.RSymbol)
	if err != nil {
		return nil, err
	}

	sk := submodel.SignaturesKey{
		Era:        data.Era,
		Pool:       data.Pool,
		TxType:     data.TxType,
		ProposalId: data.ProposalId,
	}

	skBz, err := types.EncodeToBytes(sk)
	if err != nil {
		return nil, err
	}

	var sigs []types.Bytes
	exist, err := l.conn.QueryStorage(config.RTokenSeriesModuleId, config.StorageSignatures, symBz, skBz, &sigs)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("unable to get signatures: signature key %+v ", sk)
	}

	th, err := l.threshold(data.RSymbol, data.Pool)
	if err != nil {
		return nil, fmt.Errorf("unable to get threshold: pool %s ", hex.EncodeToString(data.Pool))
	}

	//check sigs
	if len(sigs) < int(th) {
		return nil, fmt.Errorf("sigs len < threshold,sigs len: %d ,threshold: %d", len(sigs), th)
	}

	return &submodel.SubmitSignatures{
		Symbol:     data.RSymbol,
		Era:        types.NewU32(data.Era),
		Pool:       data.Pool,
		TxType:     data.TxType,
		ProposalId: data.ProposalId,
		Signature:  sigs[:],
		Threshold:  uint32(th),
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
	if symbol == core.RBNB {
		return 0, nil, nil
	}

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

func (l *listener) threshold(symbol core.RSymbol, pool []byte) (uint16, error) {
	return l.conn.threshold(symbol, pool)
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
	recipients := make([]string, 0)
	for _, ub := range unbonds {
		rec := hexutil.Encode(ub.Recipient)
		acc, ok := amounts[rec]
		if !ok {
			amounts[rec] = ub.Value
			recipients = append(recipients, rec)
		} else {
			amounts[rec] = utils.AddU128(acc, ub.Value)
		}
	}

	sort.Strings(recipients)
	receives := make([]*submodel.Receive, 0)
	total := types.NewU128(*big.NewInt(0))
	for _, rec := range recipients {
		v := amounts[rec]
		r, _ := hexutil.Decode(rec)
		rec := &submodel.Receive{Recipient: r, Value: types.NewUCompact(v.Int)}
		receives = append(receives, rec)
		total = utils.AddU128(total, v)
	}

	return receives, total, nil
}

func (l *listener) validatorId(symbol core.RSymbol, pool []byte) (interface{}, error) {
	poolBz, err := types.EncodeToBytes(pool)
	if err != nil {
		return nil, err
	}
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return nil, err
	}

	validatorIds := make([]types.Bytes, 0)
	exist, err := l.conn.QueryStorage(config.RTokenSeriesModuleId, config.StorageNominated, symBz, poolBz, &validatorIds)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("validatorId of symbol: %s, pool: %s not exist", symbol, hexutil.Encode(pool))
	}

	if len(validatorIds) == 0 {
		return nil, fmt.Errorf("no available validatorId, symbol: %s, pool: %s", symbol, hexutil.Encode(pool))
	}

	l.log.Info("get validatorId", "id", hexutil.Encode(validatorIds[0]), "symbol", symbol, "pool", hexutil.Encode(pool))

	switch symbol {
	case core.RMATIC:
		return big.NewInt(0).SetBytes(validatorIds[0]), nil
	case core.RBNB:
		addr, err := bncCmnTypes.ValAddressFromBech32(string(validatorIds[0]))
		if err != nil {
			return nil, err
		}
		return addr, nil
	default:
		return nil, fmt.Errorf("validatorId: symbol %s not supported", symbol)
	}
}

func (l *listener) leastBond(symbol core.RSymbol) (*big.Int, error) {
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return nil, err
	}

	var least types.U128
	exist, err := l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageLeastBond, symBz, nil, &least)
	if err != nil {
		return nil, err
	}

	if !exist {
		return big.NewInt(0), nil
	}

	return least.Int, nil
}

func (l *listener) eraSnapShots(symbol core.RSymbol, lastEra uint32) ([]types.Hash, error) {
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return nil, err
	}

	eraBz, err := types.EncodeToBytes(lastEra)
	if err != nil {
		return nil, err
	}

	ids := make([]types.Hash, 0)
	exist, err := l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageEraSnapShots, symBz, eraBz, &ids)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, eraSnapShotsNotExistError
	}

	return ids, nil

}
