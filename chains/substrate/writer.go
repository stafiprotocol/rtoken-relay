// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

type writer struct {
	conn          *Connection
	router        chains.Router
	log           log15.Logger
	sysErr        chan<- error
	MultisigFlows map[string]*core.MultisigFlow // CallHash => flow
}

func NewWriter(conn *Connection, log log15.Logger, sysErr chan<- error) *writer {
	return &writer{
		conn:          conn,
		log:           log,
		sysErr:        sysErr,
		MultisigFlows: make(map[string]*core.MultisigFlow),
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
	case core.EraPoolUpdated:
		return w.processEraPoolUpdated(m)
	case core.NewMultisig:
		return w.processNewMultisig(m)
	case core.MultisigExecuted:
		// todo
		return true
	case core.SubmitSignature:
		return w.processSubmitSignature(m)
	default:
		w.log.Warn("message reason unsupported", "reason", m.Reason)
		return false
	}
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
		w.log.Warn("rsymbol era is smaller than the storage one")
		return false
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

func (w *writer) processEraPoolUpdated(m *core.Message) bool {
	flow, ok := m.Content.(*core.MultisigFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	era, err := w.conn.CurrentEra()
	if err != nil {
		w.log.Error("CurrentEra error", "rsymbol", m.Source)
		return false
	}

	evt := flow.EvtEraPoolUpdated
	if evt.NewEra != era {
		w.log.Warn("era_pool_updated_event of past era, ignored", "current", era, "eventEra", evt.NewEra, "rsymbol", evt.Rsymbol)
		return true
	}

	key, others := w.conn.FoundFirstSubAccount(flow.SubAccounts)
	if key == nil {
		if flow.LastVoterFlag {
			w.sysErr <- fmt.Errorf("the last voter relay does not have key for Multisig, pool: %s, rsymbol: %s", hexutil.Encode(evt.Pool), evt.Rsymbol)
			return false
		}

		w.log.Warn("EraPoolUpdated ignored for no key")
		return false
	}

	flow.Key = key
	flow.Others = others
	err = w.conn.SetCallHash(flow)
	if err != nil {
		if err.Error() == substrate.BondEqualToUnbondError.Error() {
			w.log.Info("No need to send any call", "callhash", flow.CallHash)
			return true
		}
		w.log.Error("SetCallHash error", "err", err)
		return false
	}

	oldFlow := w.MultisigFlows[flow.CallHash]
	if oldFlow != nil {
		if oldFlow.MulExecute != nil {
			w.log.Info("already executed", "callhash", flow.CallHash)
			delete(w.MultisigFlows, flow.CallHash)
			return true
		}

		if oldFlow.NewMul == nil {
			w.sysErr <- fmt.Errorf("found old flow, but its NewMul is nil, callhash: %s", flow.CallHash)
			return false
		}

		approvals := oldFlow.Multisig.Approvals
		ac := hexutil.Encode(key.PublicKey)
		for _, apv := range approvals {
			if ac == hexutil.Encode(apv[:]) {
				w.log.Info("already approved", "approver", ac, "callhash", flow.CallHash)
				delete(w.MultisigFlows, flow.CallHash)
				return true
			}
		}

		flow.NewMul = oldFlow.NewMul
		flow.TimePoint = oldFlow.TimePoint
		flow.Multisig = oldFlow.Multisig
		err = w.conn.AsMulti(flow)
		if err != nil {
			w.log.Error("AsMulti error", "err", err, "callhash", flow.CallHash)
			return false
		}

		w.log.Error("AsMulti success", "callhash", flow.CallHash)
		return true
	}

	w.MultisigFlows[flow.CallHash] = flow

	if !flow.LastVoterFlag {
		w.log.Info("not last voter, wait for NewMultisigEvent")
		return true
	}

	err = w.conn.AsMulti(flow)
	if err != nil {
		w.log.Error("AsMulti error", "err", err, "callhash", flow.CallHash)
		return false
	}

	w.log.Error("AsMulti success", "callhash", flow.CallHash)
	return true
}

func (w *writer) processNewMultisig(m *core.Message) bool {
	flow, ok := m.Content.(*core.MultisigFlow)
	if !ok {
		w.printContentError(m)
		return false
	}
	oldFlow := w.MultisigFlows[flow.CallHash]
	if oldFlow == nil {
		w.log.Info("receive a newMultisig, wait for more flow data")
		w.MultisigFlows[flow.CallHash] = flow
		return true
	}

	if flow.EvtEraPoolUpdated == nil {
		w.sysErr <- fmt.Errorf("found old flow, but its eraPoolUpdated is nil, callhash: %s", flow.CallHash)
		return false
	}

	oldFlow.NewMul = flow.NewMul
	oldFlow.TimePoint = flow.TimePoint
	oldFlow.Multisig = flow.Multisig
	err := w.conn.AsMulti(flow)
	if err != nil {
		w.log.Error("AsMulti error", "err", err, "callhash", flow.CallHash)
		return false
	}

	w.log.Error("AsMulti success", "callhash", flow.CallHash)
	return true
}

func (w *writer) processSubmitSignature(m *core.Message) bool {
	param, ok := m.Content.(core.SubmitSignatureParams)
	if !ok {
		w.printContentError(m)
		return false
	}

	//old, err := w.conn.CurrentRsymbolEra(m.Source)
	//if err != nil {
	//	w.sysErr <- err
	//	return false
	//}
	//
	//if param.Era <= old {
	//	w.log.Warn("rsymbol era is smaller than the storage one")
	//	return false
	//}

	//newEra := types.U32(neew)
	//eraBz, _ := types.EncodeToBytes(newEra)
	//bondId := types.Hash(utils.BlakeTwo256(eraBz))
	//bk := &core.BondKey{Rsymbol: m.Source, BondId: bondId}
	//prop, err := w.conn.newUpdateEraProposal(bk, newEra)
	//
	//param:=core.SubmitSignatureParams{}
	result := w.conn.submitSignature(&param)
	w.log.Info("submitSignature", "rsymbol", m.Source, "result", result)
	return result
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
