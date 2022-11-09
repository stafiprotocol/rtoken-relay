// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

func (w *writer) processLiquidityBond(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	if flow.Reason != submodel.BondReasonDefault {
		w.log.Error("processLiquidityBond receive a message of which reason is not default", "bondId", flow.BondId.Hex(), "reason", flow.Reason)
		return false
	}

	var bondReason submodel.BondReason
	var err error
	if flow.VerifyTimes >= 5 {
		bondReason = submodel.BlockhashUnmatch
	} else {
		bondReason, err = w.conn.TransferVerify(flow.Record)
		if err != nil {
			w.log.Error("TransferVerify error", "err", err, "bondId", flow.BondId.Hex())
			flow.VerifyTimes += 1
			w.liquidityBonds <- m
			w.log.Info("processLiquidityBond", "size of liquidityBonds", len(w.liquidityBonds))
			return false
		}
	}
	w.log.Info("processLiquidityBond", "bondId", flow.BondId.Hex(), "bondReason", bondReason, "VerifyTimes", flow.VerifyTimes)
	flow.Reason = bondReason

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.LiquidityBondResult, Content: flow}
	return w.submitMessage(result)
}

func (w *writer) processBondedPools(m *core.Message) bool {
	pools, ok := m.Content.([]types.Bytes)
	if !ok {
		w.printContentError(m)
		return false
	}

	for _, p := range pools {
		w.log.Info("processBondedPools", "pool", utiles.AddHex(hexutil.Encode(p)))
		w.setBondedPools(hexutil.Encode(p), true)
	}

	return true
}

func (w *writer) processEraPoolUpdatedEvent(m *core.Message) bool {
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	flow, ok := mef.EventData.(*submodel.EraPoolUpdatedFlow)
	if !ok {
		w.log.Error("processEraPoolUpdated HeadFlow is not EraPoolUpdatedFlow")
		return false
	}

	snap := flow.Snap
	key, others := w.conn.FoundFirstSubAccount(mef.SubAccounts)
	if key == nil {
		if flow.LastVoterFlag {
			w.sysErr <- fmt.Errorf("the last voter relay does not have key for Multisig, pool: %s, symbol: %s", hexutil.Encode(snap.Pool), snap.Symbol)
			return false
		}

		w.log.Warn("EraPoolUpdated ignored for no key")
		return false
	}
	mef.Key, mef.Others = key, others

	call, err := w.conn.BondOrUnbondCall(snap)
	if err != nil {
		if err.Error() == substrate.BondEqualToUnbondError.Error() {
			w.log.Info("No need to send any call", "symbol", snap.Symbol, "era", snap.Era)
			return w.informChain(m.Destination, m.Source, mef)
		}
		w.log.Error("BondOrUnbondCall error", "err", err)
		return false
	}

	info, err := w.conn.PaymentQueryInfo(call.Extrinsic)
	if err != nil {
		w.log.Error("PaymentQueryInfo error", "err", err, "callHash", call.CallHash, "Extrinsic", call.Extrinsic)
		return false
	}
	mef.PaymentInfo = info
	mef.OpaqueCalls = []*submodel.MultiOpaqueCall{call}
	callhash := call.CallHash
	mef.NewMulCallHashs = map[string]bool{callhash: true}
	mef.MulExeCallHashs = map[string]bool{callhash: true}

	w.setEvents(call.CallHash, mef)

	if flow.LastVoterFlag {
		call.TimePoint = submodel.NewOptionTimePointEmpty()
		err = w.conn.AsMulti(mef)
		if err != nil {
			w.log.Error("AsMulti error", "err", err, "callHash", callhash)
			return false
		}
		w.log.Error("AsMulti success", "callHash", callhash)
		return true
	}

	newMuls, ok := w.getNewMultics(callhash)
	if !ok {
		w.log.Info("not last voter, wait for NewMultisigEvent")
		w.setEvents(call.CallHash, mef)
		return true
	}
	call.TimePoint = newMuls.TimePoint

	err = w.conn.AsMulti(mef)
	if err != nil {
		w.log.Error("AsMulti error", "err", err, "callHash", callhash)
		return false
	}

	w.log.Info("AsMulti success", "callHash", callhash)
	return true
}

func (w *writer) processActiveReportedEvent(m *core.Message) bool {
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	flow, ok := mef.EventData.(*submodel.ActiveReportedFlow)
	if !ok {
		w.log.Error("processActiveReported eventData is not ActiveReportedFlow")
		return false
	}

	key, others := w.conn.FoundFirstSubAccount(mef.SubAccounts)
	if key == nil {
		if flow.LastVoterFlag {
			w.sysErr <- fmt.Errorf("the last voter relay does not have key for Multisig, pool: %s, symbol: %s", hexutil.Encode(flow.Snap.Pool), flow.Snap.Symbol)
			return false
		}

		w.log.Warn("ActiveReportedEvent ignored for no key")
		return false
	}
	mef.Key, mef.Others = key, others

	call, err := w.conn.WithdrawCall()
	if err != nil {
		w.log.Error("WithdrawCall error", "err", err)
		return false
	}

	info, err := w.conn.PaymentQueryInfo(call.Extrinsic)
	if err != nil {
		w.log.Error("PaymentQueryInfo error", "err", err, "callHash", call.CallHash, "Extrinsic", call.Extrinsic)
		return false
	}
	mef.PaymentInfo = info
	mef.OpaqueCalls = []*submodel.MultiOpaqueCall{call}
	callhash := call.CallHash
	mef.NewMulCallHashs = map[string]bool{callhash: true}
	mef.MulExeCallHashs = map[string]bool{callhash: true}
	w.setEvents(call.CallHash, mef)

	if flow.LastVoterFlag {
		call.TimePoint = submodel.NewOptionTimePointEmpty()
		err = w.conn.AsMulti(mef)
		if err != nil {
			w.log.Error("AsMulti error", "err", err, "callHash", callhash)
			return false
		}
		w.log.Error("AsMulti success", "callHash", callhash)
		return true
	}

	newMuls, ok := w.getNewMultics(callhash)
	if !ok {
		w.log.Info("not last voter, wait for NewMultisigEvent")
		w.setEvents(call.CallHash, mef)
		return true
	}
	call.TimePoint = newMuls.TimePoint

	err = w.conn.AsMulti(mef)
	if err != nil {
		w.log.Error("AsMulti error", "err", err, "callHash", callhash)
		return false
	}

	w.log.Info("AsMulti success", "callHash", callhash)
	return true
}

func (w *writer) processWithdrawReportedEvent(m *core.Message) bool {
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	flow, ok := mef.EventData.(*submodel.WithdrawReportedFlow)
	if !ok {
		w.log.Error("processWithdrawReportedEvent eventData is not WithdrawReportedFlow")
		return false
	}

	balance, err := w.conn.FreeBalance(flow.Snap.Pool)
	if err != nil {
		w.log.Error("FreeBalance error", "err", err, "pool", hexutil.Encode(flow.Snap.Pool))
		return false
	}
	e, err := w.conn.ExistentialDeposit()
	if err != nil {
		w.log.Error("ExistentialDeposit error", "err", err, "pool", hexutil.Encode(flow.Snap.Pool))
		return false
	}
	least := utils.AddU128(flow.TotalAmount, e)
	if balance.Cmp(least.Int) < 0 {
		w.sysErr <- fmt.Errorf("free balance not enough for transfer back, symbol: %s, pool: %s, least: %s", flow.Symbol, hexutil.Encode(flow.Snap.Pool), least.Int.String())
		return false
	}

	key, others := w.conn.FoundFirstSubAccount(mef.SubAccounts)
	if key == nil {
		if flow.LastVoterFlag {
			w.sysErr <- fmt.Errorf("the last voter relay does not have key for Multisig, pool: %s, symbol: %s", hexutil.Encode(flow.Snap.Pool), flow.Snap.Symbol)
			return false
		}

		w.log.Warn("WithdrawReported ignored for no key")
		return false
	}
	mef.Key, mef.Others = key, others

	lastKey := w.conn.LastKey()
	if lastKey == nil {
		w.sysErr <- fmt.Errorf("processWithdrawReportedEvent: last key is nil, pool: %s, symbol: %s", hexutil.Encode(flow.Snap.Pool), flow.Snap.Symbol)
		return false
	}
	call, err := w.conn.TransferCall(lastKey.PublicKey, types.NewUCompact(flow.TotalAmount.Int))
	if err != nil {
		w.log.Error("TransferCall error", "symbol", m.Source)
		return false
	}

	info, err := w.conn.PaymentQueryInfo(call.Extrinsic)
	if err != nil {
		w.log.Error("PaymentQueryInfo error", "err", err, "callHash", call.CallHash, "Extrinsic", call.Extrinsic)
		return false
	}
	mef.PaymentInfo = info
	mef.OpaqueCalls = []*submodel.MultiOpaqueCall{call}
	callhash := call.CallHash
	mef.NewMulCallHashs = map[string]bool{callhash: true}
	mef.MulExeCallHashs = map[string]bool{callhash: true}
	w.setEvents(callhash, mef)
	w.log.Info("processWithdrawReportedEvent: event set", "callHash", callhash)

	if flow.LastVoterFlag {
		tb := &TransferBack{Symbol: string(flow.Symbol), Pool: hexutil.Encode(flow.Snap.Pool), Address: lastKey.Address, Era: fmt.Sprint(flow.Snap.Era)}
		w.log.Info("processWithdrawReportedEvent: lastVoter prepare to create transfer back", "TransferBack", *tb)
		err := CreateTransferback(w.transferRecord, tb)
		if err != nil {
			w.sysErr <- fmt.Errorf("processInformChain: CreateTransferback error: %s, TransferRecord: %+v", err, *tb)
			return false
		}

		err = CreateTransferback(w.transferRecordHistory, tb)
		if err != nil {
			w.sysErr <- fmt.Errorf("processInformChain: CreateTransferback history error: %s, TransferRecord: %+v", err, *tb)
			return false
		}

		w.log.Info("processWithdrawReportedEvent: create transfer back succeed")

		call.TimePoint = submodel.NewOptionTimePointEmpty()
		err = w.conn.AsMulti(mef)
		if err != nil {
			w.log.Error("AsMulti error", "err", err, "callHash", callhash)
			return false
		}
		w.log.Info("AsMulti success", "callHash", callhash)

		return true
	}

	newMuls, ok := w.getNewMultics(callhash)
	if !ok {
		w.log.Info("not last voter, wait for NewMultisigEvent", "callHash", callhash)
		w.setEvents(call.CallHash, mef)
		return true
	}
	call.TimePoint = newMuls.TimePoint

	err = w.conn.AsMulti(mef)
	if err != nil {
		w.log.Error("AsMulti error", "err", err, "callHash", callhash)
		return false
	}

	w.log.Info("AsMulti success", "callHash", callhash)
	return true
}

func (w *writer) processTransferReportedEvent(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.TransferReportedFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	key := w.conn.LastKey()
	if key == nil {
		w.sysErr <- fmt.Errorf("processTransferReportedEvent: last key is nil, pool: %s, symbol: %s", hexutil.Encode(flow.Snap.Pool), flow.Snap.Symbol)
		return false
	}

	tb := &TransferBack{Symbol: string(flow.Symbol), Pool: hexutil.Encode(flow.Snap.Pool), Address: key.Address, Era: fmt.Sprint(flow.Snap.Era)}
	w.log.Info("processTransferReportedEvent", "TransferBack", *tb)
	exist := IsTransferbackExist(w.transferRecord, tb)
	if !exist {
		w.log.Info("processTransferReportedEvent: transfer back not exist, will ignore", "TransferBack", tb)
		return true
	}

	client := w.conn.KeyIndex(key)
	if client == nil {
		w.sysErr <- fmt.Errorf("found transfer back record but do not have key of lastVoter, symbol: %s, Address: %s", flow.Symbol, key.Address)
		return false
	}

	balance, err := w.conn.FreeBalance(key.PublicKey)
	if err != nil {
		w.log.Error("FreeBalance error", "err", err, "Address", key.Address)
		return false
	}
	e, err := w.conn.ExistentialDeposit()
	if err != nil {
		w.log.Error("ExistentialDeposit error", "err", err, "Address", key.Address)
		return false
	}
	least := utils.AddU128(flow.TotalAmount, e)
	if balance.Cmp(least.Int) < 0 {
		w.sysErr <- fmt.Errorf("free balance not enough for transfer back, symbol: %s, Address: %s, balance: %s, least: %s", flow.Symbol, key.Address, balance.String(), least.Int.String())
		return false
	}

	err = client.BatchTransfer(flow.Receives)
	if err != nil {
		w.sysErr <- fmt.Errorf("TransferBack error: %s, symbol: %s, Address: %s", err, flow.Symbol, key.Address)
		return false
	}

	err = DeleteTransferback(w.transferRecord, tb)
	if err != nil {
		w.sysErr <- fmt.Errorf("TransferBack succeed but failed to delete Transferback: %s, %+v", err, *tb)
		return false
	}

	w.log.Info("TransferBack succeed", "TransferRecord", *tb)
	return true
}

func (w *writer) processNominationUpdatedEvent(m *core.Message) bool {
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	flow, ok := mef.EventData.(*submodel.NominationUpdatedFlow)
	if !ok {
		w.log.Error("processNominationUpdatedEvent eventData is not NominationUpdatedFlow")
		return false
	}

	era, err := w.conn.CurrentEra()
	if err != nil {
		w.log.Error("processNominationUpdatedEvent: CurrentEra error", "error", err, "symbol", flow.Symbol)
		return false
	}

	if flow.Era != era {
		w.log.Warn("processNominationUpdatedEvent: event era is not current era", "flow.era", flow.Era, "current era", era, "symbol", flow.Symbol)
		return true
	}

	key, others := w.conn.FoundFirstSubAccount(mef.SubAccounts)
	if key == nil {
		if flow.LastVoterFlag {
			w.sysErr <- fmt.Errorf("the last voter relay does not have key for Multisig, pool: %s, symbol: %s", hexutil.Encode(flow.Pool), flow.Symbol)
			return false
		}

		w.log.Warn("NominationUpdated ignored for no key")
		return false
	}
	mef.Key, mef.Others = key, others

	call, err := w.conn.NominateCall(flow.NewValidators)
	if err != nil {
		w.log.Error("NominateCall error", "err", err)
		return false
	}

	info, err := w.conn.PaymentQueryInfo(call.Extrinsic)
	if err != nil {
		w.log.Error("PaymentQueryInfo error", "err", err, "callHash", call.CallHash, "Extrinsic", call.Extrinsic)
		return false
	}
	mef.PaymentInfo = info
	mef.OpaqueCalls = []*submodel.MultiOpaqueCall{call}
	callhash := call.CallHash
	mef.NewMulCallHashs = map[string]bool{callhash: true}
	mef.MulExeCallHashs = map[string]bool{callhash: true}

	w.setEvents(call.CallHash, mef)

	if flow.LastVoterFlag {
		call.TimePoint = submodel.NewOptionTimePointEmpty()
		err = w.conn.AsMulti(mef)
		if err != nil {
			w.log.Error("AsMulti error", "err", err, "callHash", callhash)
			return false
		}
		w.log.Error("AsMulti success", "callHash", callhash)
		return true
	}

	newMuls, ok := w.getNewMultics(callhash)
	if !ok {
		w.log.Info("not last voter, wait for NewMultisigEvent")
		w.setEvents(call.CallHash, mef)
		return true
	}
	call.TimePoint = newMuls.TimePoint

	err = w.conn.AsMulti(mef)
	if err != nil {
		w.log.Error("AsMulti error", "err", err, "callHash", callhash)
		return false
	}

	w.log.Info("AsMulti success", "callHash", callhash)
	return true
}

func (w *writer) processNewMultisig(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.EventNewMultisig)
	if !ok {
		w.printContentError(m)
		return false
	}

	_, ok = w.getBondedPools(hexutil.Encode(flow.ID[:]))
	if !ok {
		w.log.Info("received a newMultisig event which the ID is not in the bondedPools, ignored")
		return true
	}

	w.setNewMultics(flow.CallHashStr, flow)

	evt, ok := w.getEvents(flow.CallHashStr)
	if !ok {
		w.log.Info("receive a newMultisig, wait for more flow data", "callHash", flow.CallHashStr)
		return true
	}

	identify := hexutil.Encode(evt.Key.PublicKey)
	for _, apv := range flow.Approvals {
		if identify == hexutil.Encode(apv[:]) {
			w.log.Info("receive a newMultisig which has already approved, will ignore", "callHash", flow.CallHashStr)
			return true
		}
	}

	for _, call := range evt.OpaqueCalls {
		if call.CallHash == flow.CallHashStr {
			call.TimePoint = flow.TimePoint
			delete(evt.NewMulCallHashs, flow.CallHashStr)
		}
	}

	if len(evt.NewMulCallHashs) != 0 {
		w.log.Info("processNewMultisig wait for more callhash", "eventId", evt.EventId)
		return true
	}

	err := w.conn.AsMulti(evt)
	if err != nil {
		w.log.Error("AsMulti error", "err", err, "callHash", flow.CallHashStr)
		return false
	}

	w.log.Info("AsMulti success", "callHash", flow.CallHashStr)
	return true
}

func (w *writer) processMultisigExecuted(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.EventMultisigExecuted)
	if !ok {
		w.printContentError(m)
		return false
	}

	_, ok = w.getBondedPools(hexutil.Encode(flow.ID[:]))
	if !ok {
		w.log.Info("received a multisigExecuted event which the ID is not in the bondedPools, ignored")
		return true
	}

	evt, ok := w.getEvents(flow.CallHashStr)
	if !ok {
		w.log.Info("receive a multisigExecuted but no evt found")
		return true
	}

	delete(evt.MulExeCallHashs, flow.CallHashStr)
	if len(evt.MulExeCallHashs) != 0 {
		w.log.Info("processMultisigExecuted wait for more callhash", "eventId", evt.EventId)
		return true
	}
	w.deleteEvents(flow.CallHashStr)
	w.deleteNewMultics(flow.CallHashStr)
	return w.informChain(m.Source, "", evt)
}

func (w *writer) processBondReportEvent(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondReportedFlow)
	if !ok {
		w.printContentError(m)
		return false
	}
	//get targets from stafi
	eraNominatedFlow := &submodel.GetEraNominatedFlow{
		Symbol:        flow.Symbol,
		Pool:          flow.Snap.Pool,
		Era:           flow.LastEra,
		NewValidators: make(chan []types.AccountID, 1),
	}

	validatorsFromStafi := make([]types.AccountID, 0)
	msg := &core.Message{
		Source: m.Destination, Destination: core.RFIS,
		Reason: core.GetEraNominated, Content: eraNominatedFlow}
	w.submitMessage(msg)

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	w.log.Debug("wait validator from stafi", "pool", eraNominatedFlow)
	//wait for validators
	select {
	case <-timer.C:
	case validatorsFromStafi = <-eraNominatedFlow.NewValidators:
	}
	w.log.Debug("validatorsFromStafi", "validator", validatorsFromStafi)
	err := w.conn.SetToPayoutStashes(flow, validatorsFromStafi)
	if err != nil {
		if err.Error() == ErrorTargetNotExist.Error() {
			w.sysErr <- fmt.Errorf("ErrorTargetNotExist, pool: %s, symbol: %s, lastEra: %d", hexutil.Encode(flow.Snap.Pool), flow.Symbol, flow.LastEra)
			return false
		}
		w.log.Error("SetToPayoutStashes error", "error", err, "pool", hexutil.Encode(flow.Snap.Pool), "symbol", flow.Symbol, "lastEra", flow.LastEra)
		return false
	}

	m.Source, m.Destination = m.Destination, m.Source
	if flow.Stashes == nil || len(flow.Stashes) == 0 {
		w.queryAndReportActive(m)
	} else {
		err := w.conn.TryPayout(flow)
		if err != nil {
			w.log.Error("TryPayout error", "error", err)
			return false
		}
		w.queryAndReportActive(m)
	}

	return true
}
