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

var multiEndError = errors.New("multiEnd")

func (l *listener) processLiquidityBondEvent(evt *submodel.ChainEvent) (*submodel.BondFlow, error) {
	evtData, err := submodel.LiquidityBondEventData(evt)
	if err != nil {
		return nil, err
	}

	bondKey := &submodel.BondKey{Rsymbol: evtData.Rsymbol, BondId: evtData.BondId}
	bk, err := types.EncodeToBytes(bondKey)
	if err != nil {
		return nil, err
	}

	br := new(submodel.BondRecord)
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

	return &submodel.BondFlow{Key: bondKey, Record: br, Reason: submodel.BondReasonDefault}, nil
}

func (l *listener) processEraPoolUpdatedEvt(evt *submodel.ChainEvent) (*submodel.MultiEventFlow, error) {
	data, err := submodel.EraPoolUpdatedData(evt)
	if err != nil {
		return nil, err
	}

	snap, err := l.snapshot(data.ShotId)
	if err != nil {
		return nil, err
	}

	th, sub, err := l.thresholdAndSubAccounts(snap.Rsymbol, snap.Pool)
	if err != nil {
		return nil, err
	}

	data.LastVoterFlag = l.conn.IsLastVoter(data.LastVoter)
	data.Snap = snap

	return &submodel.MultiEventFlow{
		EventId:     config.EraPoolUpdatedEventId,
		EventData:   data,
		Threshold:   th,
		SubAccounts: sub,
	}, nil
}

func (l *listener) processSignatureEnoughEvt(evt *submodel.ChainEvent) (*submodel.SubmitSignatures, error) {
	data, err := submodel.SignatureEnoughData(evt)
	if err != nil {
		return nil, err
	}

	sk := submodel.SignaturesKey{
		RSymbol:    data.RSymbol,
		Era:        data.Era,
		Pool:       data.Pool,
		TxType:     data.TxType,
		ProposalId: data.ProposalId,
	}
	pkBz, err := types.EncodeToBytes(sk)
	if err != nil {
		return nil, err
	}

	var sigs []types.Bytes
	exist, err := l.conn.QueryStorage(config.SubmitSignaturesModuleId, config.StorageSignatures, pkBz, nil, &sigs)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("unable to get signatures: signature key %+v ", sk)
	}

	return &submodel.SubmitSignatures{
		Symbol:     data.RSymbol,
		Era:        types.NewU32(data.Era),
		Pool:       data.Pool,
		TxType:     data.TxType,
		ProposalId: data.ProposalId,
		Signature:  sigs,
	}, nil
}


func (l *listener) processBondReportEvt(evt *submodel.ChainEvent) (*submodel.BondReportFlow, error) {
	flow, err := submodel.EventBondReport(evt)
	if err != nil {
		return nil, err
	}

	flow.LastVoterFlag = l.conn.IsLastVoter(flow.LastVoter)
	flow.LastEra = flow.Era - 1

	return flow, nil
}

func (l *listener) processWithdrawUnbondEvt(evt *submodel.ChainEvent) (*submodel.MultiEventFlow, error) {
	data, err := submodel.EventWithdrawUnbond(evt)
	if err != nil {
		return nil, err
	}

	//turn to process Transfer if rsymol==ratom
	//if data.Rsymbol == core.RATOM {
	//	return l.processTransferBackEvt(evt)
	//}

	th, sub, err := l.thresholdAndSubAccounts(data.Rsymbol, data.Pool)
	if err != nil {
		return nil, err
	}

	data.LastVoterFlag = l.conn.IsLastVoter(data.LastVoter)

	return &submodel.MultiEventFlow{
		EventId:     config.WithdrawUnbondEventId,
		EventData:   data,
		Threshold:   th,
		SubAccounts: sub,
	}, nil
}

func (l *listener) processTransferBackEvt(evt *submodel.ChainEvent) (*submodel.MultiEventFlow, error) {
	data, err := submodel.EventTransferBack(evt)
	if err != nil {
		return nil, err
	}

	receives, total, err := l.unbondings(data.Rsymbol, data.Pool, data.Era)
	if err != nil {
		return nil, err
	}

	th, sub, err := l.thresholdAndSubAccounts(data.Rsymbol, data.Pool)
	if err != nil {
		return nil, err
	}

	data.LastVoterFlag = l.conn.IsLastVoter(data.LastVoter)
	data.Receives = receives
	data.TotalAmount = total

	return &submodel.MultiEventFlow{
		EventId:     config.TransferBackEventId,
		EventData:   data,
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

func (l *listener) snapshot(shotId types.Hash) (*submodel.EraPoolSnapshot, error) {
	bz, err := types.EncodeToBytes(shotId)
	snap := new(submodel.EraPoolSnapshot)
	exist, err := l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageSnapshots, bz, nil, snap)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("snap of shotId: %s not exist", hexutil.Encode(shotId[:]))
	}

	return snap, nil
}

func (l *listener) thresholdAndSubAccounts(symbol core.RSymbol, pool []byte) (uint16, []types.Bytes, error) {
	pk := &submodel.PoolKey{Rsymbol: symbol, Pool: pool}
	pkBz, err := types.EncodeToBytes(pk)
	if err != nil {
		return 0, nil, err
	}

	var threshold uint16
	exist, err := l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageMultiThresholds, pkBz, nil, &threshold)
	if err != nil {
		return 0, nil, err
	}
	if !exist {
		return 0, nil, fmt.Errorf("threshold of pool: %s, rsymbol: %s not exist", symbol, hexutil.Encode(pool))
	}

	subs := make([]types.Bytes, 0)
	exist, err = l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageSubAccounts, pkBz, nil, &subs)
	if err != nil {
		return 0, nil, err
	}
	if !exist {
		return 0, nil, fmt.Errorf("subAccounts of pool: %s, rsymbol: %s not exist", symbol, hexutil.Encode(pool))
	}

	return threshold, subs, nil
}

func (l *listener) unbondings(symbol core.RSymbol, pool []byte, era uint32) ([]*submodel.Receive, types.U128, error) {
	type poolkey struct {
		Rsymbol core.RSymbol
		Pool    types.Bytes
		Era     types.U32
	}


	bz, err := types.EncodeToBytes(poolkey{symbol, pool, types.NewU32(era)})
	if err != nil {
		return nil, types.U128{}, err
	}

	type Unbonding struct {
		Who       types.AccountID
		Value     types.U128
		Recipient types.Bytes
	}
	unbonds := make([]Unbonding, 0)

	exist, err := l.conn.QueryStorage(config.RTokenLedgerModuleId, config.StoragePoolUnbonds, bz, nil, &unbonds)
	if err != nil {

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
		r, _ := hexutil.Decode(k)
		rec := &submodel.Receive{Recipient: types.NewAddressFromAccountID(r), Value: types.NewUCompact(v.Int)}
		receives = append(receives, rec)
		total = utils.AddU128(total, v)
	}

	return receives, total, nil
}
