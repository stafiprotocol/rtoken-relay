// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package matic

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

type writer struct {
	symbol          core.RSymbol
	conn            *Connection
	router          chains.Router
	log             log15.Logger
	sysErr          chan<- error
	liquidityBonds  chan *core.Message
	currentChainEra uint32
	bondedPoolsMtx  sync.RWMutex
	bondedPools     map[common.Address]bool
	eventMtx        sync.RWMutex
	events          map[common.Hash]*submodel.MultiEventFlow
	bondReportedMtx sync.RWMutex
	bondReporteds   map[common.Hash]*submodel.BondReportedFlow
	propMtx         sync.RWMutex
	proposalIds     map[common.Hash][]byte
	stop            <-chan int
}

const (
	bondFlowLimit = 2048
)

func NewWriter(symbol core.RSymbol, conn *Connection, log log15.Logger, sysErr chan<- error, stop <-chan int) *writer {
	return &writer{
		symbol:          symbol,
		conn:            conn,
		log:             log,
		sysErr:          sysErr,
		liquidityBonds:  make(chan *core.Message, bondFlowLimit),
		currentChainEra: 0,
		bondedPools:     make(map[common.Address]bool),
		events:          make(map[common.Hash]*submodel.MultiEventFlow),
		bondReporteds:   make(map[common.Hash]*submodel.BondReportedFlow),
		proposalIds:     make(map[common.Hash][]byte),
		stop:            stop,
	}
}

func (w *writer) setRouter(r chains.Router) {
	w.router = r
}

func (w *writer) ResolveMessage(m *core.Message) (processOk bool) {
	defer func() {
		if !processOk {
			panic(fmt.Sprintf("resolveMessage process failed. %+v", m))
		}
	}()

	switch m.Reason {
	case core.LiquidityBond:
		return w.processLiquidityBond(m)
	case core.BondedPools:
		return w.processBondedPools(m)
	case core.EraPoolUpdated:
		return w.processEraPoolUpdated(m)
	case core.BondReportEvent:
		return w.processBondReported(m)
	case core.ActiveReportedEvent:
		return w.processActiveReported(m)
	case core.WithdrawReportedEvent:
		return w.processWithdrawReported(m)
	case core.SignatureEnough:
		return w.processSignatureEnough(m)
	//case core.ValidatorUpdatedEvent:
	//	return w.processValidatorUpdatedEvent(m)
	default:
		w.log.Warn("message reason unsupported", "reason", m.Reason)
		return true
	}
}

func (w *writer) processLiquidityBond(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	if flow.Reason != submodel.BondReasonDefault {
		w.log.Error("processLiquidityBond receive a message of which reason is not default", "bondId", flow.BondId.Hex(), "reason", flow.Reason, "symbol", flow.Symbol)
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
			return true
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
		addr := common.BytesToAddress(p)
		w.log.Info("processBondedPools", "pool", addr)
		w.setBondedPools(addr, true)
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

	snap := flow.Snap
	key := w.conn.FoundKey(mef.SubAccounts)
	if key == nil {
		w.log.Error("processEraPoolUpdated no keys", "symbol", m.Destination)
		return false
	}

	validatorId, ok := mef.ValidatorId.(*big.Int)
	if !ok {
		w.log.Error("processEraPoolUpdated validatorId not bigint")
		return false
	}

	shareAddr, err := w.conn.GetValidator(validatorId)
	if err != nil {
		w.log.Error("processEraPoolUpdated get GetValidator error", "error", err)
		return false
	}

	poolAddr := common.BytesToAddress(snap.Pool)
	method, tx, err := w.conn.BondOrUnbondCall(shareAddr, snap.Bond.Int, snap.Unbond.Int, flow.LeastBond)
	if err != nil {
		if err.Error() == substrate.BondEqualToUnbondError.Error() {
			w.log.Info("BondOrUnbondCall BondEqualToUnbondError", "symbol", snap.Symbol, "era", snap.Era)
			flow.BondCall = &submodel.BondCall{
				ReportType: submodel.NewBondReport,
				Action:     submodel.BothBondUnbond,
			}
			return w.informChain(m.Destination, m.Source, mef)
		} else if err.Error() == substrate.DiffSmallerThanLeastError.Error() {
			w.log.Info("BondOrUnbondCall DiffSmallerThanLeastError", "symbol", snap.Symbol, "era", snap.Era)
			flow.BondCall = &submodel.BondCall{
				ReportType: submodel.NewBondReport,
				Action:     submodel.EitherBondUnbond,
			}
			return w.informChain(m.Destination, m.Source, mef)
		} else {
			w.log.Error("BondOrUnbondCall error", "error", err, "symbol", snap.Symbol, "era", snap.Era)
			return false
		}
	}

	flow.BondCall = &submodel.BondCall{
		ReportType: submodel.BondAndReportActive,
		Action:     submodel.BothBondUnbond,
	}
	param := submodel.SubmitSignatureParams{
		Symbol: flow.Symbol,
		Era:    types.NewU32(snap.Era),
		Pool:   snap.Pool,
		TxType: method,
	}

	txHash, err := param.EncodeToHash()
	if err != nil {
		w.log.Error("processEraPoolUpdated EncodeToHash error", "error", err)
		return false
	}

	state, err := w.conn.TxHashState(txHash, poolAddr)
	if err != nil {
		w.log.Error("processEraPoolUpdated: TxHashState error", "error", err, "txHash", txHash, "poolAddr", poolAddr)
		return false
	}

	if state == ethmodel.HashStateSuccess {
		flow.Active, flow.Reward, err = w.conn.StakedAndReward(txHash, shareAddr, poolAddr)
		if err != nil {
			w.log.Error("processEraPoolUpdated: RewardByTxHash error", "error", err, "txHash", txHash, "pool", poolAddr)
			return false
		}
		w.log.Info("processEraPoolUpdated RewardByTxHash", "reward", flow.Active, "txHash", txHash, "pool", poolAddr)
		return w.informChain(m.Destination, m.Source, mef)
	}

	msg := w.conn.MessageToSign(tx, poolAddr, txHash)
	signature, err := crypto.Sign(msg[:], key.PrivateKey())
	if err != nil {
		w.log.Error("processEraPoolUpdated sign msg error", "error", err, "msg", hexutil.Encode(msg[:]))
		return false
	}
	param.ProposalId = append(shareAddr.Bytes(), tx.CallData...)
	propKey := crypto.Keccak256Hash(param.ProposalId)
	w.setProposalIds(propKey, param.ProposalId)
	param.ProposalId = propKey.Bytes()
	param.Signature = append(msg[:], signature...)

	mef.MultiTransaction = tx
	w.setEvents(txHash, mef)
	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	return w.submitMessage(result)
}

func (w *writer) processBondReported(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondReportedFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	snap := flow.Snap
	key := w.conn.FoundKey(flow.SubAccounts)
	if key == nil {
		w.log.Error("processBondReported no keys", "symbol", m.Destination)
		return false
	}

	validatorId, ok := flow.ValidatorId.(*big.Int)
	if !ok {
		w.log.Error("processBondReported validatorId not bigint")
		return false
	}

	shareAddr, err := w.conn.GetValidator(validatorId)
	if err != nil {
		w.log.Error("processBondReported get GetValidator error", "error", err)
		return false
	}

	poolAddr := common.BytesToAddress(snap.Pool)
	tx, err := w.conn.RestakeCall(shareAddr)
	if err != nil {
		w.log.Error("RestakeCall error", "err", err)
		return false
	}

	param := submodel.SubmitSignatureParams{
		Symbol: flow.Symbol,
		Era:    types.NewU32(snap.Era),
		Pool:   snap.Pool,
		TxType: submodel.OriginalClaimRewards,
	}

	txHash, err := param.EncodeToHash()
	if err != nil {
		w.log.Error("processBondReported EncodeToHash error", "error", err)
		return false
	}

	state, err := w.conn.TxHashState(txHash, poolAddr)
	if err != nil {
		w.log.Error("processBondReported: TxHashState error", "error", err, "txHash", txHash, "poolAddr", poolAddr)
		return false
	}

	if state == ethmodel.HashStateSuccess {
		active, err := w.conn.TotalStaked(shareAddr, poolAddr)
		if err != nil {
			w.log.Error("processBondReported TotalStaked error", "error", err, "share", shareAddr, "pool", poolAddr)
			return false
		}

		bond, unbond := snap.Bond.Int, snap.Unbond.Int
		if bond.Cmp(unbond) > 0 {
			diff := big.NewInt(0).Sub(bond, unbond)
			if diff.Cmp(flow.LeastBond) <= 0 {
				active = active.Add(active, diff)
			}
		}
		flow.Snap.Active = types.NewU128(*active)

		w.log.Info("processBondReported total active", "pool", poolAddr, "active", active)
		msg := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.ActiveReport, Content: flow}
		return w.submitMessage(msg)
	}

	msg := w.conn.MessageToSign(tx, poolAddr, txHash)
	signature, err := crypto.Sign(msg[:], key.PrivateKey())
	if err != nil {
		w.log.Error("processBondReported sign msg error", "error", err, "msg", hexutil.Encode(msg[:]))
		return false
	}
	param.ProposalId = append(shareAddr.Bytes(), tx.CallData...)
	propKey := crypto.Keccak256Hash(param.ProposalId)
	w.setProposalIds(propKey, param.ProposalId)
	param.ProposalId = propKey.Bytes()
	param.Signature = append(msg[:], signature...)

	flow.MultiTransaction = tx
	w.setBondReported(txHash, flow)
	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	return w.submitMessage(result)
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

	snap := flow.Snap
	key := w.conn.FoundKey(mef.SubAccounts)
	if key == nil {
		w.log.Error("processActiveReported no keys", "symbol", m.Destination)
		return false
	}

	validatorId, ok := mef.ValidatorId.(*big.Int)
	if !ok {
		w.log.Error("processActiveReported validatorId not bigint")
		return false
	}

	shareAddr, err := w.conn.GetValidator(validatorId)
	if err != nil {
		w.log.Error("processActiveReported get GetValidator error", "error", err)
		return false
	}

	poolAddr := common.BytesToAddress(snap.Pool)
	param := submodel.SubmitSignatureParams{
		Symbol: flow.Symbol,
		Era:    types.NewU32(snap.Era),
		Pool:   snap.Pool,
		TxType: submodel.OriginalWithdrawUnbond,
	}
	txHash, err := param.EncodeToHash()
	if err != nil {
		w.log.Error("processActiveReported EncodeToHash error", "error", err)
		return false
	}

	state, err := w.conn.TxHashState(txHash, poolAddr)
	if err != nil {
		w.log.Error("processActiveReported: TxHashState error", "error", err, "txHash", txHash, "poolAddr", poolAddr)
		return false
	}

	if state == ethmodel.HashStateSuccess {
		return w.informChain(m.Destination, m.Source, mef)
	}

	nonce, err := w.conn.WithdrawNonce(shareAddr, poolAddr)
	if err != nil {
		w.log.Error("WithdrawNonce error", "err", err, "shareAddr", shareAddr, "poolAddr", poolAddr)
		return false
	}

	if nonce.Uint64() == 0 {
		w.log.Info("withdrawn", "shareAddr", shareAddr, "poolAddr", poolAddr)
		mef.OpaqueCalls = []*submodel.MultiOpaqueCall{{CallHash: txHash.Hex()}}
		return w.informChain(m.Destination, m.Source, mef)
	}

	tx, err := w.conn.WithdrawCall(shareAddr, common.BytesToAddress(snap.Pool), nonce)
	if err != nil {
		w.log.Error("processActiveReported WithdrawCall error", "err", err)
		return false
	}

	msg := w.conn.MessageToSign(tx, poolAddr, txHash)
	signature, err := crypto.Sign(msg[:], key.PrivateKey())
	if err != nil {
		w.log.Error("processActiveReported sign msg error", "error", err, "msg", hexutil.Encode(msg[:]))
		return false
	}
	param.ProposalId = append(shareAddr.Bytes(), tx.CallData...)
	propKey := crypto.Keccak256Hash(param.ProposalId)
	w.setProposalIds(propKey, param.ProposalId)
	param.ProposalId = propKey.Bytes()
	param.Signature = append(msg[:], signature...)

	mef.MultiTransaction = tx
	w.setEvents(txHash, mef)
	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	return w.submitMessage(result)
}

func (w *writer) processWithdrawReported(m *core.Message) bool {
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	flow, ok := mef.EventData.(*submodel.WithdrawReportedFlow)
	if !ok {
		w.log.Error("processWithdrawReported eventData is not WithdrawReportedFlow")
		return false
	}

	snap := flow.Snap
	key := w.conn.FoundKey(mef.SubAccounts)
	if key == nil {
		w.log.Error("processWithdrawReported no keys", "symbol", m.Destination)
		return false
	}

	poolAddr := common.BytesToAddress(snap.Pool)
	param := submodel.SubmitSignatureParams{
		Symbol: flow.Symbol,
		Era:    types.NewU32(snap.Era),
		Pool:   snap.Pool,
		TxType: submodel.OriginalTransfer,
	}
	txHash, err := param.EncodeToHash()
	if err != nil {
		w.log.Error("processWithdrawReported EncodeToHash error", "error", err)
		return false
	}

	state, err := w.conn.TxHashState(txHash, poolAddr)
	if err != nil {
		w.log.Error("processActiveReported: TxHashState error", "error", err, "txHash", txHash, "poolAddr", poolAddr)
		return false
	}

	if state == ethmodel.HashStateSuccess {
		return w.informChain(m.Destination, m.Source, mef)
	}

	if flow.TotalAmount.Uint64() == 0 {
		mef.OpaqueCalls = []*submodel.MultiOpaqueCall{{CallHash: txHash.Hex()}}
		w.log.Info("processWithdrawReported: no need to do transfer call")
		return w.informChain(m.Destination, m.Source, mef)
	}

	balance, err := w.conn.BalanceOf(poolAddr)
	if err != nil {
		w.log.Error("BalanceOf  error", "err", err, "Address", poolAddr)
		return false
	}

	if balance.Cmp(flow.TotalAmount.Int) < 0 {
		w.log.Error("free balance not enough for transfer back", "symbol", flow.Symbol, "Address", poolAddr, "balance", balance, "TotalAmount", flow.TotalAmount)
		return false
	}

	tx, err := w.conn.TransferCall(flow.Receives)
	if err != nil {
		w.log.Error("TransferCall error", "err", err)
		return false
	}

	msg := w.conn.MessageToSign(tx, poolAddr, txHash)
	signature, err := crypto.Sign(msg[:], key.PrivateKey())
	if err != nil {
		w.log.Error("processWithdrawReported sign msg error", "error", err, "msg", hexutil.Encode(msg[:]))
		return false
	}
	param.Signature = append(msg[:], signature...)
	param.ProposalId = append(tx.To.Bytes(), tx.CallData...)
	propKey := crypto.Keccak256Hash(param.ProposalId)
	w.log.Info("processWithdrawReported ProposalId", "to", tx.To, "calldata", hexutil.Encode(tx.CallData), "ProposalId", hexutil.Encode(param.ProposalId), "propKey", propKey)
	w.setProposalIds(propKey, param.ProposalId)
	param.ProposalId = propKey.Bytes()

	mef.MultiTransaction = tx
	w.setEvents(txHash, mef)
	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	return w.submitMessage(result)
}

func (w *writer) processSignatureEnough(m *core.Message) bool {
	sigs, ok := m.Content.(*submodel.SubmitSignatures)
	if !ok {
		w.log.Error("msg cast to SubmitSignatures not ok")
		w.printContentError(m)
		return false
	}
	w.log.Info("processSignatureEnough", "source", m.Source,
		"dest", m.Destination, "pool", hexutil.Encode(sigs.Pool), "txType", sigs.TxType)

	propKey := common.BytesToHash(sigs.ProposalId)
	proposalId, ok := w.getProposalIds(propKey)
	if !ok {
		w.log.Warn("processSignatureEnough unable to get proposalId", "symbol", m.Destination)
		return true
	}

	if len(proposalId) <= common.AddressLength {
		w.log.Error("processSignatureEnough ProposalId too short", "symbol", m.Destination)
		return false
	}

	signatures := make([][]byte, 0)
	for _, sig := range sigs.Signature {
		// 32 + 65 = 97
		if len(sig) != 97 {
			w.log.Error("processSignatureEnough: size of sig not 97", "sig", hexutil.Encode(sig))
			return false
		}
		signatures = append(signatures, sig[32:])
	}
	vs, rs, ss := utils.DecomposeSignature(signatures)
	to := common.BytesToAddress(proposalId[:20])
	calldata := proposalId[20:]
	poolAddr := common.BytesToAddress(sigs.Pool)

	txHash, err := sigs.EncodeToHash()
	if err != nil {
		w.log.Error("processSignatureEnough sigs EncodeToHash error", "error", err)
		return false
	}

	var mef *submodel.MultiEventFlow
	var bondFlow *submodel.BondReportedFlow
	var operation ethmodel.CallEnum
	var value, safeGas, totalGas *big.Int
	var subAccounts []types.Bytes
	var era uint32

	txTypeErr := fmt.Errorf("processSignatureEnough TxType %s not supported", sigs.TxType)
	if sigs.TxType != submodel.OriginalClaimRewards {
		mef, ok = w.getEvents(txHash)
		if !ok {
			w.log.Error("processSignatureEnough: no event for txHash", "txHash", txHash)
			return false
		}
		mef.OpaqueCalls = []*submodel.MultiOpaqueCall{{CallHash: txHash.Hex()}}
		operation = mef.MultiTransaction.Operation
		value = mef.MultiTransaction.Value
		safeGas = mef.MultiTransaction.SafeTxGas
		totalGas = mef.MultiTransaction.TotalGas
		subAccounts = mef.SubAccounts
		switch sigs.TxType {
		case submodel.OriginalBond, submodel.OriginalUnbond:
			flow, _ := mef.EventData.(*submodel.EraPoolUpdatedFlow)
			era = flow.Era
		case submodel.OriginalWithdrawUnbond:
			flow, _ := mef.EventData.(*submodel.ActiveReportedFlow)
			era = flow.Snap.Era
		case submodel.OriginalTransfer:
			flow, _ := mef.EventData.(*submodel.WithdrawReportedFlow)
			era = flow.Snap.Era
		default:
			w.log.Error(txTypeErr.Error())
			return false
		}
	} else {
		bondFlow, ok = w.getBondReported(txHash)
		if !ok {
			w.log.Error("processSignatureEnough: no bond bondReport for txHash", "txHash", txHash)
			return false
		}

		operation = bondFlow.MultiTransaction.Operation
		safeGas = bondFlow.MultiTransaction.SafeTxGas
		value = bondFlow.MultiTransaction.Value
		totalGas = bondFlow.MultiTransaction.TotalGas
		subAccounts = bondFlow.SubAccounts
		era = bondFlow.Snap.Era
	}

	report := func() bool {
		switch sigs.TxType {
		case submodel.OriginalBond, submodel.OriginalUnbond:
			flow, _ := mef.EventData.(*submodel.EraPoolUpdatedFlow)
			flow.Active, flow.Reward, err = w.conn.StakedAndReward(txHash, to, poolAddr)
			if err != nil {
				w.log.Error("processSignatureEnough: RewardByTxHash error", "error", err, "txHash", txHash, "pool", poolAddr)
				return false
			}
			w.log.Info("processSignatureEnough RewardByTxHash", "reward", flow.Active, "txHash", txHash, "pool", poolAddr)
			w.deleteProposalIds(propKey)
			return w.reportMultiEventResult(txHash, mef, m)
		case submodel.OriginalClaimRewards:
			active, err := w.conn.TotalStaked(to, poolAddr)
			if err != nil {
				w.log.Error("processSignatureEnough TotalStaked error", "error", err, "share", to, "pool", poolAddr)
				return false
			}

			snap := bondFlow.Snap
			bond, unbond := snap.Bond.Int, snap.Unbond.Int
			if bond.Cmp(unbond) > 0 {
				diff := big.NewInt(0).Sub(bond, unbond)
				if diff.Cmp(bondFlow.LeastBond) <= 0 {
					active = active.Add(active, diff)
				}
			}
			bondFlow.Snap.Active = types.NewU128(*active)
			w.log.Info("processSignatureEnough total active", "pool", poolAddr, "active", active)
			w.deleteProposalIds(propKey)
			return w.reportBondReportedResult(txHash, bondFlow, m)
		case submodel.OriginalWithdrawUnbond, submodel.OriginalTransfer:
			w.deleteProposalIds(propKey)
			return w.reportMultiEventResult(txHash, mef, m)
		default:
			w.log.Error(txTypeErr.Error())
			return false
		}
	}

	state, err := w.conn.TxHashState(txHash, poolAddr)
	if err != nil {
		w.log.Error("processSignatureEnough: TxHashState error", "error", err, "txHash", txHash, "poolAddr", poolAddr)
		return false
	}

	if state == ethmodel.HashStateSuccess {
		return report()
	}

	eraSignerFlag := w.conn.IsEraSigner(era, subAccounts)
	if !eraSignerFlag {
		w.log.Info("processSignatureEnough", "eraSignerFlag", eraSignerFlag, "txHash", txHash)
		err = w.conn.WaitTxHashSuccess(txHash, poolAddr, sigs.TxType)
		if err != nil {
			w.log.Error("processSignatureEnough: WaitTxHashSuccess error", "error", err, "txHash", txHash, "poolAddr", poolAddr)
			return false
		}
		return report()
	}

	err = w.conn.AsMulti(poolAddr, to, value, calldata, uint8(operation), safeGas, totalGas, txHash, vs, rs, ss)
	if err != nil {
		w.log.Error("AsMulti error", "err", err)
		return false
	}
	w.log.Info("AsMulti success", "txHash", txHash)

	err = w.conn.WaitTxHashSuccess(txHash, poolAddr, sigs.TxType)
	if err != nil {
		w.log.Error("processSignatureEnough: WaitTxHashSuccess error", "error", err, "txHash", txHash, "poolAddr", poolAddr)
		return false
	}

	return report()
}

func (w *writer) reportMultiEventResult(txHash common.Hash, mef *submodel.MultiEventFlow, m *core.Message) bool {
	result := w.informChain(m.Destination, m.Source, mef)
	if result {
		w.deleteEvents(txHash)
	}
	return result
}

func (w *writer) reportBondReportedResult(txHash common.Hash, flow *submodel.BondReportedFlow, m *core.Message) bool {
	msg := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.ActiveReport, Content: flow}
	result := w.submitMessage(msg)
	if result {
		w.deleteBondReported(txHash)
	}
	return result
}

func (w *writer) printContentError(m *core.Message) {
	w.log.Error("msg resolve failed", "source", m.Source, "dest", m.Destination, "reason", m.Reason)
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

func (w *writer) informChain(source, dest core.RSymbol, flow *submodel.MultiEventFlow) bool {
	msg := &core.Message{Source: source, Destination: dest, Reason: core.InformChain, Content: flow}
	return w.submitMessage(msg)
}

func (w *writer) getBondedPools(key common.Address) (bool, bool) {
	w.bondedPoolsMtx.RLock()
	defer w.bondedPoolsMtx.RUnlock()
	value, exist := w.bondedPools[key]
	return value, exist
}

func (w *writer) setBondedPools(key common.Address, value bool) {
	w.bondedPoolsMtx.Lock()
	defer w.bondedPoolsMtx.Unlock()
	w.bondedPools[key] = value
}

func (w *writer) getEvents(key common.Hash) (*submodel.MultiEventFlow, bool) {
	w.eventMtx.RLock()
	defer w.eventMtx.RUnlock()
	value, exist := w.events[key]
	return value, exist
}

func (w *writer) setEvents(key common.Hash, value *submodel.MultiEventFlow) {
	w.eventMtx.Lock()
	defer w.eventMtx.Unlock()
	w.events[key] = value
}

func (w *writer) deleteEvents(key common.Hash) {
	w.eventMtx.Lock()
	defer w.eventMtx.Unlock()
	delete(w.events, key)
}

func (w *writer) getBondReported(key common.Hash) (*submodel.BondReportedFlow, bool) {
	w.bondReportedMtx.RLock()
	defer w.bondReportedMtx.RUnlock()
	value, exist := w.bondReporteds[key]
	return value, exist
}

func (w *writer) setBondReported(key common.Hash, value *submodel.BondReportedFlow) {
	w.bondReportedMtx.Lock()
	defer w.bondReportedMtx.Unlock()
	w.bondReporteds[key] = value
}

func (w *writer) deleteBondReported(key common.Hash) {
	w.bondReportedMtx.Lock()
	defer w.bondReportedMtx.Unlock()
	delete(w.bondReporteds, key)
}

func (w *writer) getProposalIds(key common.Hash) ([]byte, bool) {
	w.propMtx.RLock()
	defer w.propMtx.RUnlock()
	value, exist := w.proposalIds[key]
	return value, exist
}

func (w *writer) setProposalIds(key common.Hash, value []byte) {
	w.propMtx.Lock()
	defer w.propMtx.Unlock()
	w.proposalIds[key] = value
}

func (w *writer) deleteProposalIds(key common.Hash) {
	w.propMtx.Lock()
	defer w.propMtx.Unlock()
	delete(w.proposalIds, key)
}

func (w *writer) start() error {
	go func() {
		for {
			select {
			case <-w.stop:
				close(w.liquidityBonds)
				w.log.Info("writer stopped")
				return
			case msg := <-w.liquidityBonds:
				result := w.processLiquidityBond(msg)
				w.log.Info("retry processLiquidityBond", "result", result)
			}
		}
	}()

	return nil
}
