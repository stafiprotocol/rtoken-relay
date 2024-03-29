// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"math/big"

	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

func (w *writer) processGetEraNominated(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.GetEraNominatedFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	validator, err := w.conn.GetEraNominated(flow.Symbol, flow.Pool, flow.Era)
	if err != nil {
		w.log.Warn("GetEraNominated failed", "err", err)
		if err.Error() == ErrorNotExist.Error() {
			flow.NewValidators <- make([]types.AccountID, 0)
			return true
		}
	}
	flow.NewValidators <- validator
	return true
}

func (w *writer) processGetBondState(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.GetBondStateFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	bondState, err := w.conn.bondState(flow.Symbol, flow.BlockHash, flow.TxHash)
	if err != nil {
		w.log.Warn("bondState failed", "err", err)
		if err.Error() == ErrorNotExist.Error() {
			flow.BondState <- submodel.Default
			return true
		}
		flow.BondState <- submodel.Default
		return false
	}

	flow.BondState <- bondState
	return true
}

func (w *writer) processGetSubmitSignatures(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.GetSubmitSignaturesFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	submitSigs, err := w.conn.getSubmitSignatures(flow.Symbol, uint32(flow.Era), flow.Pool, flow.TxType, flow.ProposalId)
	if err != nil {
		if err.Error() == ErrorNotExist.Error() {
			flow.Signatures <- []types.Bytes{}
			return true
		}
		flow.Signatures <- []types.Bytes{}
		return false
	}

	flow.Signatures <- submitSigs
	return true
}

func (w *writer) processGetPoolThreshold(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.GetPoolThresholdFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	threshold, err := w.conn.poolThreshold(flow.Symbol, flow.Pool)
	if err != nil {
		if err.Error() == ErrorNotExist.Error() {
			flow.Threshold <- 0
			return true
		}
		flow.Threshold <- 0
		return false
	}

	flow.Threshold <- uint32(threshold)
	return true
}

func (w *writer) processGetEraRate(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.GetEraRateFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	rate, err := w.conn.GetEraRate(flow.Symbol, flow.Era)
	if err != nil {
		if err.Error() == ErrorNotExist.Error() {
			flow.Rate <- 0
			return true
		}
		flow.Rate <- 0
		return false
	}

	flow.Rate <- rate
	return true
}

func (w *writer) processSubmitSignature(m *core.Message) bool {
	param, ok := m.Content.(submodel.SubmitSignatureParams)
	if !ok {
		w.printContentError(m)
		return false
	}

	need, err := w.conn.NeedMoreSignature(&param)
	if err != nil {
		w.log.Error("NeedMoreSignature error", "error", err)
		return false
	}

	if need {
		result := w.conn.submitSignature(&param)
		w.log.Info("submitSignature", "symbol", m.Source, "result", result)
		return result
	}

	w.log.Info("processSubmitSignature: signature already enough")
	return true
}

func (w *writer) processLiquidityBondResult(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	if flow.Reason == submodel.BondReasonDefault {
		w.log.Error("processLiquidityBondResult receive a message of which reason is default", "bondId", flow.BondId.Hex(), "reason", flow.Reason)
		return false
	}

	prop, err := w.conn.LiquidityBondProposal(flow)
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

	if w.currentChainEra == 0 {
		currentEra, err := w.conn.CurrentChainEra(m.Source)
		if err != nil {
			if err.Error() != fmt.Sprintf("era of symbol %s not exist", m.Source) {
				w.log.Error("failed to get CurrentChainEra", "error", err)
				return false
			}
		}
		w.currentChainEra = currentEra
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
		if err.Error() != fmt.Sprintf("era of symbol %s not exist", m.Source) {
			w.log.Error("failed to get CurrentChainEra", "error", err)
			return false
		}
	}

	bondedPools, err := w.getLatestBondedPools(m.Source)
	if err != nil {
		w.log.Error("processNewEra error", "error", err)
		return false
	}
	w.log.Debug("processNewEra", "bondedPools len", len(bondedPools), "symbol", m.Source)
	msg := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.BondedPools, Content: bondedPools}
	if !w.submitMessage(msg) {
		w.log.Warn("failed to update bondedPools")
		return false
	} else {
		w.log.Info("successed to update bondedPools")
	}

	should := cur
	if old != 0 && cur > old {
		should = old + 1
	}
	eraBz, err := types.EncodeToBytes(should)
	if err != nil {
		w.log.Error("processNewEra EncodeToBytes error", "error", err, "should", should)
		return false
	}
	bondId := types.Hash(utils.BlakeTwo256(eraBz))
	prop, err := w.conn.SetChainEraProposal(m.Source, bondId, should)
	if err != nil {
		w.log.Error("processNewEra SetChainEraProposal error", "error", err)
		return false
	}
	result := w.conn.resolveProposal(prop, true)
	w.log.Info("processNewEra", "symbol", m.Source, "uploadEra", should, "current", cur, "result", result)
	if result {
		w.currentChainEra = should
	}
	return result
}

// Will deal:
// - bondReport/newBondReport/bondAndReportActive/bondAndReportActiveWithPendingValue
// - withdrawReport
// - transferReport
func (w *writer) processInformChain(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	if flow.EventData == nil {
		w.log.Error("flow.EventDate is nil")
		return false
	}

	var prop *submodel.Proposal
	var err error

	switch data := flow.EventData.(type) {
	case *submodel.NominationUpdatedFlow:
		w.log.Info("NominationUpdated", "symbol", data.Symbol, "era", data.Era)
		return true
	case *submodel.EraPoolUpdatedFlow:
		call := data.BondCall // total 4 case(bond_report/new_bond_report/bond_and_report_active/bond_and_report_active_with_pending_value)
		method := ""
		switch data.Symbol {
		case core.RMATIC:
			switch call.ReportType {
			case submodel.NewBondReport:
				prop, err = w.conn.NewBondReportProposal(data)
				if err != nil {
					w.log.Error("MethodNewBondReportProposal", "error", err)
					return false
				}
				method = "NewBondReportProposal"
			case submodel.BondAndReportActive:
				reportActive := new(big.Int).Add(data.Active, data.Reward)
				ok, err := w.checkActive(data.Snap.Symbol, data.Snap.Active, types.NewU128(*reportActive))
				if err != nil {
					w.log.Error("checkActive", "error", err)
					return false
				}
				if !ok {
					w.log.Error("checkActive failed", "snapActive", data.Snap.Active, "reportActive", reportActive)
					return false
				}

				prop, err = w.conn.BondAndReportActiveProposal(data)
				if err != nil {
					w.log.Error("MethodBondAndReportActiveProposal", "error", err)
					return false
				}
				method = "BondAndReportActiveProposal"
			default:
				w.log.Error("processInformChain: ReportType not supported", "ReportType", call.ReportType)
				return false
			}
		case core.RBNB:
			switch call.ReportType {
			case submodel.BondAndReportActiveWithPendingValue:
				ok, err := w.checkActive(data.Snap.Symbol, data.Snap.Active, types.NewU128(*data.Active))
				if err != nil {
					w.log.Error("checkActive", "error", err)
					return false
				}
				if !ok {
					w.log.Error("checkActive failed", "snapActive", data.Snap.Active, "reportActive", data.Active)
					return false
				}

				prop, err = w.conn.BondAndReportActiveWithPendingValueProposal(data)
				if err != nil {
					w.log.Error("BondAndReportActiveWithPendingValueProposal", "error", err)
					return false
				}
				method = "BondAndReportActiveWithPendingValueProposal"
			default:
				w.log.Error("processInformChain: ReportType not supported", "ReportType", call.ReportType)
				return false
			}
		default:
			prop, err = w.conn.CommonReportProposal(config.MethodBondReport, m.Source, data.ShotId, data.ShotId)
			if err != nil {
				w.log.Error("MethodBondReportProposal", "error", err)
				return false
			}
			method = "CommonReportProposal"
		}

		result := w.conn.resolveProposal(prop, true)
		w.log.Info(fmt.Sprintf("%s resolveProposal", method), "result", result)
		return result
	case *submodel.ActiveReportedFlow:
		callhash := flow.RunTimeCalls[0].CallHash
		bondId, err := types.NewHashFromHexString(utiles.AddHex(callhash))
		if err != nil {
			w.sysErr <- fmt.Errorf("processInformChain: callhash %s decode error: %s", bondId, err)
			return false
		}

		prop, err = w.conn.CommonReportProposal(config.MethodWithdrawReport, m.Source, bondId, data.ShotId)
		if err != nil {
			w.log.Error("MethodWithdrawReportProposal", "error", err)
			return false
		}
		result := w.conn.resolveProposal(prop, true)
		w.log.Info("MethodWithdrawReportProposal resolveProposal", "result", result)
		return result

	case *submodel.WithdrawReportedFlow:
		callhash := flow.RunTimeCalls[0].CallHash
		bondId, err := types.NewHashFromHexString(utiles.AddHex(callhash))
		if err != nil {
			w.sysErr <- fmt.Errorf("processInformChain: callhash %s decode error: %s", bondId, err)
			return false
		}
		prop, err = w.conn.CommonReportProposal(config.MethodTransferReport, m.Source, bondId, data.ShotId)
		if err != nil {
			w.log.Error("MethodTransferReportProposal", "error", err)
			return false
		}
		result := w.conn.resolveProposal(prop, true)
		w.log.Info("MethodTransferReportProposal resolveProposal", "result", result)
		return result

	default:
		return false
	}
}

// Will deal:
// - activeReport/newActiveReport
func (w *writer) processActiveReport(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondReportedFlow)
	if !ok {
		w.printContentError(m)
		return false
	}
	ok, err := w.checkActive(flow.Snap.Symbol, flow.Snap.Active, flow.ReportActive)
	if err != nil {
		w.log.Error("checkActive", "error", err)
		return false
	}
	if !ok {
		w.log.Error("checkActive failed", "snapActive", flow.Snap.Active, "reportActive", flow.ReportActive)
		return false
	}

	var prop *submodel.Proposal
	if flow.NewActiveReportFlag {
		prop, err = w.conn.NewActiveReportProposal(flow)
		if err != nil {
			w.log.Error("NewActiveReportProposal", "error", err)
			return false
		}
	} else {
		prop, err = w.conn.ActiveReportProposal(flow)
		if err != nil {
			w.log.Error("ActiveReportProposal", "error", err)
			return false
		}
	}

	result := w.conn.resolveProposal(prop, true)
	w.log.Info("ActiveReportProposal resolveProposal", "result", result)

	return result
}

var perBillBase = big.NewInt(1e9)

func (w *writer) checkActive(symbol core.RSymbol, snapActive, reportActive types.U128) (bool, error) {
	if snapActive.Cmp(big.NewInt(0)) == 0 {
		return true, nil
	}
	changeRateLimit, err := w.conn.ActiveChangeRateLimit(symbol)
	if err != nil {
		if err != ErrorNotExist {
			return false, err
		}
		changeRateLimit = DefaultActiveChangeRateLimit
	}

	if changeRateLimit == 0 {
		return true, nil
	}

	var change *big.Int
	if snapActive.Cmp(reportActive.Int) > 0 {
		change = new(big.Int).Sub(snapActive.Int, reportActive.Int)
	} else {
		change = new(big.Int).Sub(reportActive.Int, snapActive.Int)
	}
	tmp := new(big.Int).Mul(snapActive.Int, big.NewInt(int64(changeRateLimit)))
	changeLimit := new(big.Int).Div(tmp, perBillBase)

	if change.Cmp(changeLimit) > 0 {
		return false, nil
	}

	return true, nil
}

func (w *writer) processExeLiquidityBondAndSwap(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.ExeLiquidityBondAndSwapFlow)
	if !ok {
		w.printContentError(m)
		return false
	}
	if flow.Reason == submodel.BondReasonDefault {
		w.log.Error("processExeLiquidityBondAndSwap receive a message of which reason is default", "txHash", flow.Txhash, "reason", flow.Reason)
		return false
	}
	// should exe when state not exist or failed
	bondState, err := w.conn.bondState(flow.Symbol, flow.Blockhash, flow.Txhash)
	if err != nil {
		if err != ErrorNotExist {
			w.log.Error("processExeLiquidityBondAndSwap get bond state failed", "txHash", flow.Txhash, "reason", flow.Reason, "err", err)
			return false
		}
	} else {
		if bondState == submodel.Success {
			return true
		}
	}

	prop, err := w.conn.ExeLiquidityBondAndSwapProposal(flow)
	if err != nil {
		w.log.Error("processExeLiquidityBondAndSwap proposal", "error", err)
		w.sysErr <- err
		return false
	}

	result := w.conn.resolveProposal(prop, flow.Reason == submodel.Pass)
	w.log.Info("processExeLiquidityBondAndSwap resolveProposal", "result", result)

	return result
}
