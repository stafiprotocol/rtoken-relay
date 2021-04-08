// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

type writer struct {
	conn            *Connection
	router          chains.Router
	log             log15.Logger
	sysErr          chan<- error
	eventMtx        sync.RWMutex
	newMulTicsMtx   sync.RWMutex
	bondedPoolsMtx  sync.RWMutex
	events          map[string]*submodel.MultiEventFlow
	newMultics      map[string]*submodel.EventNewMultisig
	multiExecuted   map[string]*submodel.EventMultisigExecuted
	bondedPools     map[string]bool
	currentChainEra uint32
}

var (
	waitBlockNum    = uint64(30)
	singleBlockTime = 6 * time.Second
)

func NewWriter(conn *Connection, log log15.Logger, sysErr chan<- error) *writer {
	return &writer{
		conn:            conn,
		log:             log,
		sysErr:          sysErr,
		events:          make(map[string]*submodel.MultiEventFlow),
		newMultics:      make(map[string]*submodel.EventNewMultisig),
		multiExecuted:   make(map[string]*submodel.EventMultisigExecuted),
		bondedPools:     make(map[string]bool),
		currentChainEra: 0,
	}
}

func (w *writer) start() error {
	return nil
}

func (w *writer) setRouter(r chains.Router) {
	w.router = r
}

func (w *writer) ResolveMessage(m *core.Message) bool {
	switch m.Reason {
	case core.LiquidityBond:
		return w.processLiquidityBond(m)
	case core.LiquidityBondResult:
		return w.processLiquidityBondResult(m)
	case core.NewEra:
		return w.processNewEra(m)
	case core.BondedPools:
		return w.processBondedPools(m)
	case core.EraPoolUpdated:
		return w.processEraPoolUpdated(m)
	case core.NewMultisig:
		return w.processNewMultisig(m)
	case core.MultisigExecuted:
		return w.processMultisigExecuted(m)
	case core.InformChain:
		return w.processInformChain(m)
	case core.BondReportEvent:
		return w.processBondReportEvent(m)
	case core.ActiveReport:
		return w.processActiveReport(m)
	case core.ActiveReportedEvent:
		return w.processActiveReported(m)
	case core.WithdrawReportedEvent:
		return w.processWithdrawReportedEvent(m)
	default:
		w.log.Warn("message reason unsupported", "reason", m.Reason)
		return false
	}
}

func (w *writer) processLiquidityBond(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	if flow.Reason != submodel.BondReasonDefault {
		w.log.Error("processLiquidityBond receive a message of which reason is not default", "bondId", flow.Key.BondId.Hex(), "reason", flow.Reason)
		return false
	}

	bondReason, err := w.conn.TransferVerify(flow.Record)
	if err != nil {
		w.log.Error("TransferVerify error", "err", err, "bondId", flow.Key.BondId.Hex())
		return false
	}
	w.log.Info("processLiquidityBond", "bondReason", bondReason)
	flow.Reason = bondReason

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.LiquidityBondResult, Content: flow}
	return w.submitMessage(result)
}

func (w *writer) processLiquidityBondResult(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	if flow.Reason == submodel.BondReasonDefault {
		w.log.Error("processLiquidityBondResult receive a message of which reason is default", "bondId", flow.Key.BondId.Hex(), "reason", flow.Reason)
		return false
	}

	prop, err := w.conn.LiquidityBondProposal(flow.Key, flow.Reason)
	if err != nil {
		w.log.Error("processLiquidityBondResult proposal", "error", err)
		w.sysErr <- err
		return false
	}

	result := w.conn.resolveProposal(prop, flow.Reason == submodel.Pass)
	w.log.Info("processLiquidityBondResult resolveProposal", "result", result)

	return result
}

func (w *writer) processNewEra(m *core.Message) bool {
	cur, ok := m.Content.(uint32)
	if !ok {
		w.printContentError(m)
		return false
	}

	if w.currentChainEra != 0 && cur <= w.currentChainEra {
		return true
	}

	continuable, err := w.conn.EraContinuable(m.Source)
	if err != nil {
		w.log.Error("EraContinuable error", "error", err)
		return false
	}
	if !continuable {
		return true
	}

	old, err := w.conn.CurrentChainEra(m.Source)
	if err != nil {
		if err.Error() != fmt.Sprintf("era of rsymbol %s not exist", m.Source) {
			w.log.Error("failed to get CurrentChainEra", "error", err)
			return false
		}
	}

	bondedPools, err := w.getLatestBondedPools(m.Source)
	if err != nil {
		w.log.Error("processNewEra error", "error", err)
		return false
	}
	w.log.Info("", "len(bondedPools)", len(bondedPools), "symbol", m.Source)
	msg := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.BondedPools, Content: bondedPools}
	if !w.submitMessage(msg) {
		w.log.Warn("failed to update bondedPools")
		return false
	} else {
		w.log.Info("successed to update bondedPools")
	}

	should := old + 1
	eraBz, _ := types.EncodeToBytes(should)
	bondId := types.Hash(utils.BlakeTwo256(eraBz))
	bk := &submodel.BondKey{Rsymbol: m.Source, BondId: bondId}
	prop, err := w.conn.SetChainEraProposal(bk, should)
	result := w.conn.resolveProposal(prop, true)
	w.log.Info("processNewEra", "symbol", m.Source, "uploadEra", should, "current", cur, "result", result)
	if result {
		w.currentChainEra = should
	}
	return result
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

func (w *writer) processEraPoolUpdated(m *core.Message) bool {
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

	era, err := w.conn.CurrentEra()
	if err != nil {
		w.log.Error("CurrentEra error", "rsymbol", m.Source)
		return false
	}

	snap := flow.Snap
	if snap.Era != era {
		w.log.Warn("era_pool_updated_event of past era, ignored", "current", era, "eventEra", snap.Era, "rsymbol", snap.Rsymbol)
		return true
	}

	key, others := w.conn.FoundFirstSubAccount(mef.SubAccounts)
	if key == nil {
		if flow.LastVoterFlag {
			w.sysErr <- fmt.Errorf("the last voter relay does not have key for Multisig, pool: %s, rsymbol: %s", hexutil.Encode(snap.Pool), snap.Rsymbol)
			return false
		}

		w.log.Warn("EraPoolUpdated ignored for no key")
		return false
	}
	mef.Key, mef.Others = key, others

	call, err := w.conn.BondOrUnbondCall(snap)
	if err != nil {
		if err.Error() == substrate.BondEqualToUnbondError.Error() {
			w.log.Info("No need to send any call", "rsymbol", snap.Rsymbol, "era", snap.Era)
			return w.informChain(m.Destination, m.Source, mef)
		}
		w.log.Error("BondOrUnbondCall error", "err", err)
		return false
	}

	info, err := w.conn.PaymentQueryInfo(call.Extrinsic)
	if err != nil {
		w.log.Error("PaymentQueryInfo error", "err", err, "callhash", call.CallHash, "Extrinsic", call.Extrinsic)
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
			w.log.Error("AsMulti error", "err", err, "callhash", callhash)
			return false
		}
		w.log.Error("AsMulti success", "callhash", callhash)
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
		w.log.Error("AsMulti error", "err", err, "callhash", callhash)
		return false
	}

	w.log.Error("AsMulti success", "callhash", callhash)
	return true
}

func (w *writer) processActiveReported(m *core.Message) bool {
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
			w.sysErr <- fmt.Errorf("the last voter relay does not have key for Multisig, pool: %s, rsymbol: %s", hexutil.Encode(flow.Snap.Pool), flow.Snap.Rsymbol)
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
		w.log.Error("PaymentQueryInfo error", "err", err, "callhash", call.CallHash, "Extrinsic", call.Extrinsic)
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
			w.log.Error("AsMulti error", "err", err, "callhash", callhash)
			return false
		}
		w.log.Error("AsMulti success", "callhash", callhash)
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
		w.log.Error("AsMulti error", "err", err, "callhash", callhash)
		return false
	}

	w.log.Error("AsMulti success", "callhash", callhash)
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
	e, err := w.conn.ExistentialDeposit()
	least := utils.AddU128(flow.TotalAmount, e)
	if balance.Cmp(least.Int) < 0 {
		w.log.Error("free balance not enough for transfer back")
	}

	key, others := w.conn.FoundFirstSubAccount(mef.SubAccounts)
	if key == nil {
		if flow.LastVoterFlag {
			w.sysErr <- fmt.Errorf("the last voter relay does not have key for Multisig, pool: %s, rsymbol: %s", hexutil.Encode(flow.Snap.Pool), flow.Snap.Rsymbol)
			return false
		}

		w.log.Warn("WithdrawReported ignored for no key")
		return false
	}
	mef.Key, mef.Others = key, others

	calls, hashs1, hashs2, err := w.conn.TransferCalls(flow.Receives)
	if err != nil {
		w.log.Error("TransferCalls error", "rsymbol", m.Source)
		return false
	}

	info, err := w.conn.PaymentQueryInfo(calls[0].Extrinsic)
	if err != nil {
		w.log.Error("PaymentQueryInfo error", "err", err, "Extrinsic", calls[0].Extrinsic)
		return false
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

	if flow.LastVoterFlag {
		err = w.conn.AsMulti(mef)
		if err != nil {
			w.log.Error("AsMulti error", "err", err, "first callhash", mef.OpaqueCalls[0].CallHash)
			return false
		}
		w.log.Error("AsMulti success", "first callhash", mef.OpaqueCalls[0].CallHash)
		return true
	}

	for _, call := range calls {

		newMuls, ok := w.getNewMultics(call.CallHash)
		if ok {
			call.TimePoint = newMuls.TimePoint
			delete(mef.NewMulCallHashs, call.CallHash)
		}
	}

	if len(mef.NewMulCallHashs) != 0 {
		w.log.Info("not last voter, wait for NewMultisigEvent")
		return true
	}

	err = w.conn.AsMulti(mef)
	if err != nil {
		w.log.Error("AsMulti error", "err", err, "eventId", mef.EventId)
		return false
	}

	w.log.Error("AsMulti success", "eventId", mef.EventId)
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
		w.log.Info("receive a newMultisig, wait for more flow data")
		return true
	}

	for _, call := range evt.OpaqueCalls {
		if call.CallHash == flow.CallHashStr {
			call.TimePoint = flow.TimePoint
		}
		delete(evt.NewMulCallHashs, call.CallHash)
	}

	if len(evt.NewMulCallHashs) != 0 {
		w.log.Info("processNewMultisig wait for more callhash", "eventId", evt.EventId)
		return true
	}

	err := w.conn.AsMulti(evt)
	if err != nil {
		w.log.Error("AsMulti error", "err", err, "callhash", flow.CallHash)
		return false
	}

	w.log.Error("AsMulti success", "callhash", flow.CallHash)
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

func (w *writer) informChain(source, dest core.RSymbol, flow *submodel.MultiEventFlow) bool {
	msg := &core.Message{Source: source, Destination: dest, Reason: core.InformChain, Content: flow}
	return w.submitMessage(msg)
}

func (w *writer) processInformChain(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	if flow.EventData == nil {
		w.sysErr <- fmt.Errorf("the headflow is nil: %s", flow.EventId)
		return false
	}

	if data, ok := flow.EventData.(*submodel.EraPoolUpdatedFlow); ok {
		bk := &submodel.BondKey{Rsymbol: m.Source, BondId: data.ShotId}
		prop, err := w.conn.CommonReportProposal(config.MethodBondReport, bk, data.ShotId)
		if err != nil {
			w.log.Error("MethodBondReportProposal", "error", err)
			return false
		}
		result := w.conn.resolveProposal(prop, true)
		w.log.Info("MethodBondReportProposal resolveProposal", "result", result)
		return result
	}

	callhash := flow.OpaqueCalls[0].CallHash
	bondId, err := types.NewHashFromHexString(utiles.AddHex(callhash))
	if err != nil {
		w.sysErr <- fmt.Errorf("processInformChain: callhash %s decode error: %s", bondId, err)
		return false
	}

	bk := &submodel.BondKey{Rsymbol: m.Source, BondId: bondId}
	if data, ok := flow.EventData.(*submodel.ActiveReportedFlow); ok {
		prop, err := w.conn.CommonReportProposal(config.MethodWithdrawReport, bk, data.ShotId)
		if err != nil {
			w.log.Error("MethodWithdrawReportProposal", "error", err)
			return false
		}
		result := w.conn.resolveProposal(prop, true)
		w.log.Info("MethodWithdrawReportProposal resolveProposal", "result", result)
		return result
	}

	if data, ok := flow.EventData.(*submodel.WithdrawReportedFlow); ok {
		prop, err := w.conn.CommonReportProposal(config.MethodTransferReport, bk, data.ShotId)
		if err != nil {
			w.log.Error("MethodTransferReportProposal", "error", err)
			return false
		}
		result := w.conn.resolveProposal(prop, true)
		w.log.Info("MethodTransferReportProposal resolveProposal", "result", result)
		return result
	}

	return false
}

func (w *writer) processBondReportEvent(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondReportedFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	err := w.conn.SetToPayoutStashes(flow)
	if err != nil {
		if err.Error() == TargetNotExistError.Error() {
			w.sysErr <- fmt.Errorf("TargetNotExistError, pool: %s, rsymbol: %s, lastEra: %d", hexutil.Encode(flow.Snap.Pool), flow.Snap.Rsymbol, flow.LastEra)
			return false
		}
		w.log.Error("SetToPayoutStashes error", "error", err, "pool", hexutil.Encode(flow.Snap.Pool), "rsymbol", flow.Snap.Rsymbol, "lastEra", flow.LastEra)
		return false
	}

	m.Source, m.Destination = m.Destination, m.Source
	waitFlag := true
	if flow.LastVoterFlag {
		w.conn.TryPayout(flow)
		waitFlag = false
	}
	go w.waitPayout(m, waitFlag)

	return true
}

func (w *writer) waitPayout(m *core.Message, waitFlag bool) {
	flow, ok := m.Content.(*submodel.BondReportedFlow)
	if !ok {
		w.printContentError(m)
		return
	}

	if waitFlag {
		startBlk, err := w.conn.LatestBlockNumber()
		if err != nil {
			w.log.Error("waitPayout latest block error", "error", err)
			return
		}

		endBlk := startBlk + waitBlockNum

	LOOP:
		for {
			select {
			default:
				blk, err := w.conn.LatestBlockNumber()
				if err != nil {
					return
				}
				if blk < endBlk {
					time.Sleep(singleBlockTime)
				} else {
					break LOOP
				}
			}
		}
	}

	ledger := new(submodel.StakingLedger)
	exist, err := w.conn.QueryStorage(config.StakingModuleId, config.StorageLedger, flow.Snap.Pool, nil, ledger)
	if err != nil {
		w.log.Error("waitPayout get ledger error", "error", err, "pool", hexutil.Encode(flow.Snap.Pool))
		return
	}
	if !exist {
		w.log.Error("waitPayout ledger not exist", "pool", hexutil.Encode(flow.Snap.Pool))
		return
	}

	flow.Snap.Active = types.NewU128(big.Int(ledger.Active))
	w.log.Info("waitPayout", "waitFlag", waitFlag, "pool", hexutil.Encode(flow.Snap.Pool), "active", flow.Snap.Active)
	m.Content = flow
	m.Reason = core.ActiveReport

	w.submitMessage(m)
}

func (w *writer) processActiveReport(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondReportedFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	bk := &submodel.BondKey{Rsymbol: m.Source, BondId: flow.ShotId}
	prop, err := w.conn.ActiveReportProposal(bk, flow.ShotId, flow.Snap.Active)
	if err != nil {
		w.log.Error("ActiveReportProposal", "error", err)
		return false
	}

	result := w.conn.resolveProposal(prop, true)
	w.log.Info("ActiveReportProposal resolveProposal", "result", result)

	return result
}

func (w *writer) printContentError(m *core.Message) {
	w.log.Error("msg resolve failed", "source", m.Source, "dest", m.Destination, "reason", m.Reason)
}

func (w *writer) getLatestBondedPools(symbol core.RSymbol) ([]types.Bytes, error) {
	symbz, err := types.EncodeToBytes(symbol)
	if err != nil {
		w.sysErr <- err
		return nil, err
	}
	bondedPools := make([]types.Bytes, 0)
	exist, err := w.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageBondedPools, symbz, nil, &bondedPools)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("bonded pools not extis: %s", symbol)
	}

	return bondedPools, nil
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (w *writer) submitMessage(m *core.Message) bool {
	if m.Destination == "" {
		m.Destination = core.RFIS
	}
	err := w.router.Send(m)
	if err != nil {
		w.log.Error("failed to process event", "err", err)
		return false
	}

	return true
}

func (w *writer) getEvents(key string) (*submodel.MultiEventFlow, bool) {
	w.eventMtx.RLock()
	defer w.eventMtx.RUnlock()
	value, exist := w.events[key]
	return value, exist
}

func (w *writer) setEvents(key string, value *submodel.MultiEventFlow) {
	w.eventMtx.Lock()
	defer w.eventMtx.Unlock()
	w.events[key] = value
}

func (w *writer) deleteEvents(key string) {
	w.eventMtx.Lock()
	defer w.eventMtx.Unlock()
	delete(w.events, key)
}

func (w *writer) getNewMultics(key string) (*submodel.EventNewMultisig, bool) {
	w.newMulTicsMtx.RLock()
	defer w.newMulTicsMtx.RUnlock()
	value, exist := w.newMultics[key]
	return value, exist
}

func (w *writer) setNewMultics(key string, value *submodel.EventNewMultisig) {
	w.newMulTicsMtx.Lock()
	defer w.newMulTicsMtx.Unlock()
	w.newMultics[key] = value
}

func (w *writer) deleteNewMultics(key string) {
	w.newMulTicsMtx.Lock()
	defer w.newMulTicsMtx.Unlock()
	delete(w.newMultics, key)
}

func (w *writer) getBondedPools(key string) (bool, bool) {
	w.bondedPoolsMtx.RLock()
	defer w.bondedPoolsMtx.RUnlock()
	value, exist := w.bondedPools[key]
	return value, exist
}

func (w *writer) setBondedPools(key string, value bool) {
	w.bondedPoolsMtx.Lock()
	defer w.bondedPoolsMtx.Unlock()
	w.bondedPools[key] = value
}
