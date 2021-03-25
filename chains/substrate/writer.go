// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

type writer struct {
	conn          *Connection
	router        chains.Router
	log           log15.Logger
	sysErr        chan<- error
	multisigFlows map[string]*core.MultisigFlow // CallHash => flow
	bondedPools   map[string]bool
}

var (
	waitBlockNum    = uint64(50)
	singleBlockTime = 6 * time.Second
)

func NewWriter(conn *Connection, log log15.Logger, sysErr chan<- error) *writer {
	return &writer{
		conn:          conn,
		log:           log,
		sysErr:        sysErr,
		multisigFlows: make(map[string]*core.MultisigFlow),
		bondedPools:   make(map[string]bool),
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
	case core.InitLastVoter:
		return w.initLastVoter(m)
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
	case core.WithdrawUnbond:
		return w.processWithdrawUnbond(m)
	default:
		w.log.Warn("message reason unsupported", "reason", m.Reason)
		return false
	}
}

func (w *writer) initLastVoter(m *core.Message) bool {
	sym := m.Source

	bz, err := types.EncodeToBytes(sym)
	if err != nil {
		w.sysErr <- err
		return false
	}

	var voter types.AccountID
	exist, err := w.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageLastVoter, bz, nil, &voter)
	if err != nil {
		w.sysErr <- err
		return false
	}

	if exist {
		return true
	}

	bk := &core.BondKey{Rsymbol: sym, BondId: utils.BlakeTwo256(bz)}
	prop, err := w.conn.InitLastVoterProposal(bk)
	if err != nil {
		w.log.Error("InitLastVoterProposal", "error", err)
		w.sysErr <- err
		return false
	}

	result := w.conn.resolveProposal(prop, true)
	w.log.Info("InitLastVoterProposal resolveProposal", "result", result)

	return result
}

func (w *writer) processLiquidityBond(m *core.Message) bool {
	flow, ok := m.Content.(*core.BondFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	if flow.Reason != core.BondReasonDefault {
		w.log.Error("processLiquidityBond receive a message of which reason is not default", "bondId", flow.Key.BondId.Hex(), "reason", flow.Reason)
		return false
	}

	bondReason, err := w.conn.TransferVerify(flow.Record)
	if err != nil {
		w.log.Error("TransferVerify error", "err", err, "bondId", flow.Key.BondId.Hex())
		return false
	}

	flow.Reason = bondReason

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.LiquidityBondResult, Content: flow}
	return w.submitMessage(result)
}

func (w *writer) processLiquidityBondResult(m *core.Message) bool {
	flow, ok := m.Content.(*core.BondFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	if flow.Reason == core.BondReasonDefault {
		w.log.Error("processLiquidityBondResult receive a message of which reason is default", "bondId", flow.Key.BondId.Hex(), "reason", flow.Reason)
		return false
	}

	prop, err := w.conn.LiquidityBondProposal(flow.Key, flow.Reason)
	if err != nil {
		w.log.Error("processLiquidityBondResult proposal", "error", err)
		w.sysErr <- err
		return false
	}

	result := w.conn.resolveProposal(prop, flow.Reason == core.Pass)
	w.log.Info("processLiquidityBondResult resolveProposal", "result", result)

	return result
}

func (w *writer) processNewEra(m *core.Message) bool {
	neew, ok := m.Content.(uint32)
	if !ok {
		w.printContentError(m)
		return false
	}

	old, err := w.conn.CurrentRsymbolEra(m.Source)
	if err != nil {
		w.sysErr <- err
		return false
	}

	if neew <= old {
		w.log.Warn("rsymbol era is nog bigger than the storage one")
		return false
	}

	symbz, err := types.EncodeToBytes(m.Source)
	if err != nil {
		w.sysErr <- err
		return false
	}
	bondedPools := make([]types.Bytes, 0)
	exist, err := w.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageBondedPools, symbz, nil, &bondedPools)
	if err != nil {
		w.sysErr <- err
		return false
	}
	if !exist {
		w.log.Warn("processNewEra", "no bonded bondedPools for rsymbol", m.Source)
	}
	w.log.Info("", "len(bondedPools)", len(bondedPools))

	msg := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.BondedPools, Content: bondedPools}
	if !w.submitMessage(msg) {
		w.log.Warn("bondedPools failed")
	} else {
		w.log.Info("bondedPools successed")
	}

	newEra := types.U32(neew)
	eraBz, _ := types.EncodeToBytes(newEra)
	bondId := types.Hash(utils.BlakeTwo256(eraBz))
	bk := &core.BondKey{Rsymbol: m.Source, BondId: bondId}
	prop, err := w.conn.newUpdateEraProposal(bk, newEra)
	result := w.conn.resolveProposal(prop, true)
	w.log.Info("processNewEra", "rsymbol", m.Source, "era", newEra, "result", result)
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
		w.bondedPools[hexutil.Encode(p)] = true
	}

	return true
}

func (w *writer) processEraPoolUpdated(m *core.Message) bool {
	mf, ok := m.Content.(*core.MultisigFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	flow, ok := mf.HeadFlow.(*core.EraPoolUpdatedFlow)
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

	key, others := w.conn.FoundFirstSubAccount(mf.MulCall.SubAccounts)
	if key == nil {
		if flow.LastVoterFlag {
			w.sysErr <- fmt.Errorf("the last voter relay does not have key for Multisig, pool: %s, rsymbol: %s", hexutil.Encode(snap.Pool), snap.Rsymbol)
			return false
		}

		w.log.Warn("EraPoolUpdated ignored for no key")
		return false
	}

	mf.MulCall.Key = key
	mf.MulCall.Others = others
	call, err := w.conn.BondOrUnbondCall(snap)
	if err != nil {
		if err.Error() == substrate.BondEqualToUnbondError.Error() {
			w.log.Info("No need to send any call", "callhash", mf.MulCall.CallHash)
			return w.informChain(m.Destination, m.Source, mf)
		}
		w.log.Error("BondOrUnbondCall error", "err", err)
		return false
	}
	mf.MulCall.Extrinsic, mf.MulCall.Opaque, mf.MulCall.CallHash = call.Extrinsic, call.Opaque, call.CallHash

	callhash := mf.MulCall.CallHash
	oldFlow := w.multisigFlows[callhash]
	if oldFlow != nil {
		if oldFlow.MulExecuted != nil {
			w.log.Info("already executed", "callhash", callhash)
			delete(w.multisigFlows, callhash)
			return true
		}

		if oldFlow.NewMul == nil {
			w.sysErr <- fmt.Errorf("found old flow, but its NewMul is nil, callhash: %s", callhash)
			return false
		}

		approvals := oldFlow.Multisig.Approvals
		ac := hexutil.Encode(key.PublicKey)
		for _, apv := range approvals {
			if ac == hexutil.Encode(apv[:]) {
				w.log.Info("already approved", "approver", ac, "callhash", callhash)
				delete(w.multisigFlows, callhash)
				return true
			}
		}

		mf.MulCall.TimePoint = oldFlow.MulCall.TimePoint
		mf.NewMul = oldFlow.NewMul
		mf.Multisig = oldFlow.Multisig
		w.multisigFlows[callhash] = mf

		err = w.conn.AsMulti(mf.MulCall)
		if err != nil {
			w.log.Error("AsMulti error", "err", err, "callhash", callhash)
			return false
		}

		w.log.Error("AsMulti success", "callhash", callhash)
		return true
	}

	w.multisigFlows[callhash] = mf
	if !flow.LastVoterFlag {
		w.log.Info("not last voter, wait for NewMultisigEvent")
		return true
	}

	mf.MulCall.TimePoint = core.NewOptionTimePointEmpty()
	err = w.conn.AsMulti(mf.MulCall)
	if err != nil {
		w.log.Error("AsMulti error", "err", err, "callhash", callhash)
		return false
	}

	w.log.Error("AsMulti success", "callhash", callhash)
	return true
}

func (w *writer) processNewMultisig(m *core.Message) bool {
	flow, ok := m.Content.(*core.MultisigFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	_, ok = w.bondedPools[hexutil.Encode(flow.NewMul.ID[:])]
	if !ok {
		w.log.Info("received a newMultisig event which the ID is not in the bondedPools, ignored")
		return true
	}

	oldFlow := w.multisigFlows[flow.CallHash]
	if oldFlow == nil {
		w.log.Info("receive a newMultisig, wait for more flow data")
		w.multisigFlows[flow.CallHash] = flow
		return true
	}

	/// received multiExecuted first
	if oldFlow.HeadFlow == nil || oldFlow.MulCall == nil {
		w.log.Warn("NewMultisig found old flow, but its HeadFlow/mulcall is nil", "callhash", flow.CallHash)
		return false
	}

	oldFlow.NewMul = flow.NewMul
	oldFlow.MulCall.TimePoint = flow.MulCall.TimePoint
	oldFlow.Multisig = flow.Multisig
	err := w.conn.AsMulti(oldFlow.MulCall)
	if err != nil {
		w.log.Error("AsMulti error", "err", err, "callhash", flow.CallHash)
		return false
	}

	w.log.Error("AsMulti success", "callhash", flow.CallHash)
	return true
}

func (w *writer) processMultisigExecuted(m *core.Message) bool {
	flow, ok := m.Content.(*core.MultisigFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	_, ok = w.bondedPools[hexutil.Encode(flow.MulExecuted.ID[:])]
	if !ok {
		w.log.Info("received a multisigExecuted event which the ID is not in the bondedPools, ignored")
		return true
	}

	oldFlow := w.multisigFlows[flow.CallHash]
	if oldFlow == nil {
		w.log.Warn("received a multisigExecuted event, but found no oldFlow")
		return true
	}

	/// received multiExecuted first
	if oldFlow.HeadFlow == nil || oldFlow.MulCall == nil {
		w.sysErr <- fmt.Errorf("MultisigExecuted found old flow, but its HeadFlow/MulCall is nil, callhash: %s", flow.CallHash)
		return false
	}

	result := w.informChain(m.Destination, m.Source, oldFlow)
	if result {
		delete(w.multisigFlows, flow.CallHash)
	}

	return result
}

func (w *writer) informChain(source, dest core.RSymbol, flow *core.MultisigFlow) bool {
	msg := &core.Message{Source: source, Destination: dest, Reason: core.InformChain, Content: flow}
	return w.submitMessage(msg)
}

func (w *writer) processInformChain(m *core.Message) bool {
	flow, ok := m.Content.(*core.MultisigFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	if flow.HeadFlow == nil {
		w.sysErr <- fmt.Errorf("the headflow is nil: %s", flow.CallHash)
		return false
	}

	bondId, err := types.NewHashFromHexString(utiles.AddHex(flow.CallHash))
	if err != nil {
		w.sysErr <- fmt.Errorf("processBondReport: callhash %s decode error: %s", flow.CallHash, err)
		return false
	}

	hf := flow.HeadFlow
	bk := &core.BondKey{Rsymbol: m.Source, BondId: bondId}

	if data, ok := hf.(*core.EraPoolUpdatedFlow); ok {
		prop, err := w.conn.CommonReportProposal(config.MethodBondReport, bk, data.ShotId)
		if err != nil {
			w.log.Error("MethodBondReportProposal", "error", err)
			return false
		}
		result := w.conn.resolveProposal(prop, true)
		w.log.Info("MethodBondReportProposal resolveProposal", "result", result)
		return result
	}

	if data, ok := hf.(*core.WithdrawUnbondFlow); ok {
		prop, err := w.conn.CommonReportProposal(config.MethodWithdrawReport, bk, data.ShotId)
		if err != nil {
			w.log.Error("MethodWithdrawReportProposal", "error", err)
			return false
		}
		result := w.conn.resolveProposal(prop, true)
		w.log.Info("MethodWithdrawReportProposal resolveProposal", "result", result)
		return result
	}

	if data, ok := hf.(*core.TransferFlow); ok {
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
	flow, ok := m.Content.(*core.BondReportFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	err := w.conn.SetToPayoutStashes(flow)
	if err != nil {
		if err.Error() == TargetNotExistError.Error() {
			w.sysErr <- fmt.Errorf("TargetNotExistError, pool: %s, rsymbol: %s, lastEra: %d", hexutil.Encode(flow.Pool), flow.Rsymbol, flow.LastEra)
			return false
		}
		w.log.Error("SetToPayoutStashes error", "error", err, "pool", hexutil.Encode(flow.Pool), "rsymbol", flow.Rsymbol, "lastEra", flow.LastEra)
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
	flow, ok := m.Content.(*core.BondReportFlow)
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
					break
				}
			}
		}
	}

	ledger := new(substrate.StakingLedger)
	exist, err := w.conn.QueryStorage(config.StakingModuleId, config.StorageLedger, flow.Pool, nil, ledger)
	if err != nil {
		w.log.Error("waitPayout get ledger error", "error", err, "pool", hexutil.Encode(flow.Pool))
		return
	}
	if !exist {
		w.log.Error("waitPayout ledger not exist", "pool", hexutil.Encode(flow.Pool))
		return
	}

	flow.Active = ledger.Active
	m.Content = flow
	m.Reason = core.ActiveReport

	w.submitMessage(m)
}

func (w *writer) processActiveReport(m *core.Message) bool {
	flow, ok := m.Content.(*core.BondReportFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	bk := &core.BondKey{Rsymbol: m.Source, BondId: flow.ShotId}
	prop, err := w.conn.ActiveReportProposal(bk, flow.ShotId, flow.Active)
	if err != nil {
		w.log.Error("ActiveReportProposal", "error", err)
		return false
	}

	result := w.conn.resolveProposal(prop, true)
	w.log.Info("ActiveReportProposal resolveProposal", "result", result)

	return result
}

func (w *writer) processWithdrawUnbond(m *core.Message) bool {
	mf, ok := m.Content.(*core.MultisigFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	flow, ok := mf.HeadFlow.(*core.WithdrawUnbondFlow)
	if !ok {
		w.log.Error("processWithdrawUnbond HeadFlow is not WithdrawUnbondFlow")
		return false
	}

	era, err := w.conn.CurrentEra()
	if err != nil {
		w.log.Error("CurrentEra error", "rsymbol", m.Source)
		return false
	}

	if flow.Era != era {
		w.log.Warn("era_pool_updated_event of past era, ignored", "current", era, "eventEra", flow.Era, "rsymbol", flow.Rsymbol)
		return true
	}

	key, others := w.conn.FoundFirstSubAccount(mf.MulCall.SubAccounts)
	if key == nil {
		if flow.LastVoterFlag {
			w.sysErr <- fmt.Errorf("the last voter relay does not have key for Multisig, pool: %s, rsymbol: %s", hexutil.Encode(flow.Pool), flow.Rsymbol)
			return false
		}

		w.log.Warn("WithdrawUnbond ignored for no key")
		return false
	}

	mf.MulCall.Key = key
	mf.MulCall.Others = others
	call, err := w.conn.WithdrawCall()
	if err != nil {
		w.log.Error("WithdrawCall error", "err", err)
		return false
	}
	mf.MulCall.Extrinsic, mf.MulCall.Opaque, mf.MulCall.CallHash = call.Extrinsic, call.Opaque, call.CallHash


	callhash := mf.MulCall.CallHash
	oldFlow := w.multisigFlows[callhash]
	if oldFlow != nil {
		if oldFlow.MulExecuted != nil {
			w.log.Info("already executed", "callhash", callhash)
			delete(w.multisigFlows, callhash)
			return true
		}

		if oldFlow.NewMul == nil {
			w.sysErr <- fmt.Errorf("found old flow, but its NewMul is nil, callhash: %s", callhash)
			return false
		}

		approvals := oldFlow.Multisig.Approvals
		ac := hexutil.Encode(key.PublicKey)
		for _, apv := range approvals {
			if ac == hexutil.Encode(apv[:]) {
				w.log.Info("already approved", "approver", ac, "callhash", callhash)
				delete(w.multisigFlows, callhash)
				return true
			}
		}

		mf.MulCall.TimePoint = oldFlow.MulCall.TimePoint
		mf.NewMul = oldFlow.NewMul
		mf.Multisig = oldFlow.Multisig
		w.multisigFlows[callhash] = mf

		err = w.conn.AsMulti(mf.MulCall)
		if err != nil {
			w.log.Error("AsMulti error", "err", err, "callhash", callhash)
			return false
		}

		w.log.Error("AsMulti success", "callhash", callhash)
		return true
	}

	w.multisigFlows[callhash] = mf
	if !flow.LastVoterFlag {
		w.log.Info("not last voter, wait for NewMultisigEvent")
		return true
	}

	mf.MulCall.TimePoint = core.NewOptionTimePointEmpty()
	err = w.conn.AsMulti(mf.MulCall)
	if err != nil {
		w.log.Error("AsMulti error", "err", err, "callhash", callhash)
		return false
	}

	w.log.Error("AsMulti success", "callhash", callhash)
	return true
}

func (w *writer) printContentError(m *core.Message) {
	w.log.Error("msg resolve failed", "source", m.Source, "dest", m.Destination, "reason", m.Reason)
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (w *writer) submitMessage(m *core.Message) bool {
	err := w.router.Send(m)
	if err != nil {
		w.log.Error("failed to process event", "err", err)
		return false
	}

	return true
}
