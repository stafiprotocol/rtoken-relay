// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

type writer struct {
	conn                *Connection
	router              chains.Router
	log                 log15.Logger
	sysErr              chan<- error
	eraPoolUpdatedFlows map[string]*core.EraPoolUpdatedFlow
}

func NewWriter(conn *Connection, log log15.Logger, sysErr chan<- error) *writer {
	return &writer{
		conn:   conn,
		log:    log,
		sysErr: sysErr,
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
		w.log.Warn("symbol era is smaller than the storage one")
		return false
	}

	newEra := types.U32(neew)
	eraBz, _ := types.EncodeToBytes(newEra)
	bondId := types.Hash(utils.BlakeTwo256(eraBz))
	bk := &core.BondKey{Rsymbol: m.Source, BondId: bondId}
	prop, err := w.conn.newUpdateEraProposal(bk, newEra)
	result := w.conn.resolveProposal(prop, true)
	w.log.Info("processNewEra", "symbol", m.Source, "era", newEra, "result", result)
	return result
}

func (w *writer) processEraPoolUpdated(m *core.Message) bool {
	flow, ok := m.Content.(*core.EraPoolUpdatedFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	era, err := w.conn.CurrentEra()
	if err != nil {
		w.log.Error("CurrentEra error", "rsymbol", m.Source)
		return false
	}

	if flow.Evt.NewEra != era {
		w.log.Warn("era_pool_updated_event of past era, ignored", "current", era, "eventEra", flow.Evt.NewEra, "rsymbol", flow.Evt.Rsymbol)
		return true
	}

	key, others := w.conn.FoundFirstSubAccount(flow.SubAccounts)
	if key == nil {
		if flow.LastVoterFlag {
			w.sysErr <- fmt.Errorf("the last voter relay does not have key for Multisig, pool: %s, rsymbol: %s", hexutil.Encode(flow.Evt.Pool), flow.Evt.Rsymbol)
			return false
		}

		w.log.Warn("EraPoolUpdated ignored for no key")
		return false
	}

	call, err := w.conn.OpaqueCall(flow)
	if err != nil {
		if err.Error() == substrate.BondEqualToUnbondError.Error() {
			// todo
			return true
		}
	}

	if flow.LastVoterFlag {
		err := w.conn.NewMultisig(key, flow.Threshold, others, nil, call)
		if err != nil {
			w.log.Error("NewMultisig error", "err", err)
			return false
		}
	}

	return false
}

func (w *writer) processNewMultisig(m *core.Message) bool {
	return false
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
