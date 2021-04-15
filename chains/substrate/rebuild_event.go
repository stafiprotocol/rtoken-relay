package substrate

import (
	"errors"
	"fmt"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
)

func (w *writer) rebuildEraUpdateEventUseStafiConn(shotId types.Hash, snap *submodel.PoolSnapshot) error {
	th, sub, err := w.stafiConn.thresholdAndSubAccounts(snap.Symbol, snap.Pool)
	if err != nil {
		return err
	}
	key, others := w.conn.FoundFirstSubAccount(sub)
	if key == nil {
		return errors.New("FoundFirstSubAccount have no sub key")
	}

	flow := &submodel.EraPoolUpdatedFlow{
		Symbol:    snap.Symbol,
		Era:       snap.Era,
		ShotId:    shotId,
		LastVoter: snap.LastVoter,
	}
	flow.LastVoterFlag = w.conn.IsLastVoter(flow.LastVoter)
	flow.Snap = snap

	mef := &submodel.MultiEventFlow{
		EventId:     config.EraPoolUpdatedEventId,
		Symbol:      snap.Symbol,
		EventData:   flow,
		Threshold:   th,
		SubAccounts: sub,
	}

	mef.Key, mef.Others = key, others
	call, err := w.conn.BondOrUnbondCall(snap)
	if err != nil {
		//no tx no need to rebuild
		if err.Error() == substrate.BondEqualToUnbondError.Error() {
			return err
		}
		w.log.Error("BondOrUnbondCall error", "err", err)
		return err
	}

	info, err := w.conn.PaymentQueryInfo(call.Extrinsic)
	if err != nil {
		w.log.Error("PaymentQueryInfo error", "err", err, "callHash", call.CallHash, "Extrinsic", call.Extrinsic)
		return err
	}
	mef.PaymentInfo = info
	mef.OpaqueCalls = []*submodel.MultiOpaqueCall{call}
	callhash := call.CallHash
	mef.NewMulCallHashs = map[string]bool{callhash: true}
	mef.MulExeCallHashs = map[string]bool{callhash: true}

	w.setEvents(call.CallHash, mef)
	return nil
}

func (w *writer) rebuildActiveReportedEventUseStafiConn(shotId types.Hash, snap *submodel.PoolSnapshot) error {
	th, sub, err := w.stafiConn.thresholdAndSubAccounts(snap.Symbol, snap.Pool)
	if err != nil {
		return err
	}
	key, others := w.conn.FoundFirstSubAccount(sub)
	if key == nil {
		return errors.New("EraPoolUpdated FoundFirstSubAccount have no sub key")
	}

	flow := &submodel.ActiveReportedFlow{
		Symbol:    snap.Symbol,
		ShotId:    shotId,
		LastVoter: snap.LastVoter,
	}
	flow.LastVoterFlag = w.conn.IsLastVoter(flow.LastVoter)
	flow.Snap = snap

	mef := &submodel.MultiEventFlow{
		EventId:     config.ActiveReportedEventId,
		Symbol:      snap.Symbol,
		EventData:   flow,
		Threshold:   th,
		SubAccounts: sub,
	}
	mef.Key, mef.Others = key, others

	call, err := w.conn.WithdrawCall()
	if err != nil {
		w.log.Error("WithdrawCall error", "err", err)
		return err
	}

	info, err := w.conn.PaymentQueryInfo(call.Extrinsic)
	if err != nil {
		w.log.Error("PaymentQueryInfo error", "err", err, "callHash", call.CallHash, "Extrinsic", call.Extrinsic)
		return err
	}
	mef.PaymentInfo = info
	mef.OpaqueCalls = []*submodel.MultiOpaqueCall{call}
	callhash := call.CallHash
	mef.NewMulCallHashs = map[string]bool{callhash: true}
	mef.MulExeCallHashs = map[string]bool{callhash: true}
	w.setEvents(call.CallHash, mef)
	return nil
}

func (w *writer) rebuildWithdrawReportedEventUseStafiConn(shotId types.Hash, snap *submodel.PoolSnapshot) error {
	th, sub, err := w.stafiConn.thresholdAndSubAccounts(snap.Symbol, snap.Pool)
	if err != nil {
		return err
	}
	key, others := w.conn.FoundFirstSubAccount(sub)
	if key == nil {
		return errors.New("EraPoolUpdated FoundFirstSubAccount have no sub key")
	}

	receives, total, err := w.stafiConn.unbondings(snap.Symbol, snap.Pool, snap.Era)
	if err != nil {
		return err
	}

	flow := &submodel.WithdrawReportedFlow{
		Symbol:    snap.Symbol,
		ShotId:    shotId,
		LastVoter: snap.LastVoter,
	}

	flow.LastVoterFlag = w.conn.IsLastVoter(flow.LastVoter)
	flow.Snap = snap
	flow.Receives = receives
	flow.TotalAmount = total

	mef := &submodel.MultiEventFlow{
		EventId:     config.WithdrawReportedEventId,
		Symbol:      snap.Symbol,
		EventData:   flow,
		Threshold:   th,
		SubAccounts: sub,
	}
	mef.Key, mef.Others = key, others

	calls, hashs1, hashs2, err := w.conn.TransferCalls(flow.Receives)
	if err != nil {
		w.log.Error("TransferCalls error", "symbol", snap.Symbol)
		return err
	}

	info, err := w.conn.PaymentQueryInfo(calls[0].Extrinsic)
	if err != nil {
		w.log.Error("PaymentQueryInfo error", "err", err, "Extrinsic", calls[0].Extrinsic)
		return err
	}

	mef.PaymentInfo = info
	mef.OpaqueCalls = calls
	mef.NewMulCallHashs = hashs1
	mef.MulExeCallHashs = hashs2
	for _, call := range calls {
		if flow.LastVoterFlag {
			call.TimePoint = submodel.NewOptionTimePointEmpty()
		}
		w.setEvents(call.CallHash, mef)
	}

	return nil
}

func (w *writer) getLatestBondedPoolsUseStafiConn(symbol core.RSymbol) ([]types.Bytes, error) {
	symbz, err := types.EncodeToBytes(symbol)
	if err != nil {
		w.sysErr <- err
		return nil, err
	}
	bondedPools := make([]types.Bytes, 0)
	exist, err := w.stafiConn.QueryStorage(config.RTokenLedgerModuleId, config.StorageBondedPools, symbz, nil, &bondedPools)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("bonded pools not extis: %s", symbol)
	}

	return bondedPools, nil
}
