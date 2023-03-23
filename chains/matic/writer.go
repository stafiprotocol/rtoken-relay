// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package matic

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	log             core.Logger
	sysErr          chan<- error
	liquidityBonds  chan *core.Message
	currentChainEra uint32
	bondedPoolsMtx  sync.RWMutex
	bondedPools     map[common.Address]bool
	stop            <-chan int
}

const (
	bondFlowLimit = 2048
)

func NewWriter(symbol core.RSymbol, conn *Connection, log core.Logger, sysErr chan<- error, stop <-chan int) *writer {
	return &writer{
		symbol:          symbol,
		conn:            conn,
		log:             log,
		sysErr:          sysErr,
		liquidityBonds:  make(chan *core.Message, bondFlowLimit),
		currentChainEra: 0,
		bondedPools:     make(map[common.Address]bool),
		stop:            stop,
	}
}

func (w *writer) setRouter(r chains.Router) {
	w.router = r
}

func (w *writer) ResolveMessage(m *core.Message) bool {

	var err error
	switch m.Reason {
	case core.LiquidityBondEvent:
		err = w.processLiquidityBond(m)
	case core.BondedPools:
		err = w.processBondedPools(m)

	case core.EraPoolUpdatedEvent:
		err = w.processEraPoolUpdated(m)
	case core.BondReportedEvent:
		err = w.processBondReported(m)
	case core.ActiveReportedEvent:
		err = w.processActiveReported(m)
	case core.WithdrawReportedEvent:
		err = w.processWithdrawReported(m)
	default:
		err = fmt.Errorf("message reason unsupported, reason: %s", m.Reason)
		w.log.Warn("resolve message", "err", err)
		return true
	}
	if err != nil {
		w.log.Error("resolve message", "err", err)
		w.sysErr <- err
		return false
	}
	return true
}

func (w *writer) processLiquidityBond(m *core.Message) error {
	flow, ok := m.Content.(*submodel.BondFlow)
	if !ok {
		w.printContentError(m)
		return fmt.Errorf("cast err")
	}

	if flow.Reason != submodel.BondReasonDefault {
		w.log.Error("processLiquidityBond receive a message of which reason is not default", "bondId", flow.BondId.Hex(), "reason", flow.Reason, "symbol", flow.Symbol)
		return fmt.Errorf("flow rewason err")
	}

	var bondReason submodel.BondReason
	var err error
	if flow.VerifyTimes >= 10 {
		bondReason = submodel.BlockhashUnmatch
	} else {
		bondReason, err = w.conn.TransferVerify(flow.Record)
		if err != nil {
			w.log.Error("TransferVerify error", "err", err, "bondId", flow.BondId.Hex())
			flow.VerifyTimes += 1
			w.liquidityBonds <- m
			w.log.Info("processLiquidityBond", "size of liquidityBonds", len(w.liquidityBonds))
			return nil
		}
	}
	w.log.Info("processLiquidityBond", "bondId", flow.BondId.Hex(), "bondReason", bondReason, "VerifyTimes", flow.VerifyTimes)
	flow.Reason = bondReason

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.LiquidityBondResult, Content: flow}
	return w.submitMessage(result)
}

func (w *writer) processBondedPools(m *core.Message) error {
	pools, ok := m.Content.([]types.Bytes)
	if !ok {
		w.printContentError(m)
		return fmt.Errorf("cast err")
	}

	for _, p := range pools {
		addr := common.BytesToAddress(p)
		w.log.Info("processBondedPools", "pool", addr)
		w.setBondedPools(addr, true)
	}

	return nil
}

func (w *writer) processEraPoolUpdated(m *core.Message) error {
	w.log.Info("processEraPoolUpdated")
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m)
		return fmt.Errorf("cast error")
	}

	flow, ok := mef.EventData.(*submodel.EraPoolUpdatedFlow)
	if !ok {
		return fmt.Errorf("processEraPoolUpdated HeadFlow is not EraPoolUpdatedFlow")
	}

	snap := flow.Snap
	key := w.conn.FoundKey(mef.SubAccounts)
	if key == nil {
		return fmt.Errorf("processEraPoolUpdated no keys symbol %s", core.RMATIC)
	}

	validatorId := mef.MaticValidatorId

	shareAddr, err := w.conn.GetValidator(validatorId)
	if err != nil {
		return fmt.Errorf("processEraPoolUpdated get GetValidator error: %s", err)
	}

	poolAddr := common.BytesToAddress(snap.Pool)
	txType, tx, err := w.conn.BondOrUnbondCall(shareAddr, snap.Bond.Int, snap.Unbond.Int, flow.LeastBond)
	if err != nil {
		if err.Error() == substrate.ErrorBondEqualToUnbond.Error() {
			w.log.Info("BondOrUnbondCall ErrorBondEqualToUnbond", "symbol", snap.Symbol, "era", snap.Era)
			flow.BondCall = &submodel.BondCall{
				ReportType: submodel.NewBondReport,
				Action:     submodel.BothBondUnbond,
			}
			return w.informChain(m.Destination, m.Source, mef)
		} else if err.Error() == substrate.ErrorDiffSmallerThanLeast.Error() {
			w.log.Info("BondOrUnbondCall ErrorDiffSmallerThanLeast", "symbol", snap.Symbol, "era", snap.Era)
			flow.BondCall = &submodel.BondCall{
				ReportType: submodel.NewBondReport,
				Action:     submodel.EitherBondUnbond,
			}
			return w.informChain(m.Destination, m.Source, mef)
		} else {
			return fmt.Errorf("BondOrUnbondCall error %s symbol %s era %d", err, snap.Symbol, snap.Era)
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
		TxType: txType,
	}

	txHash, err := param.EncodeToHash()
	if err != nil {
		return fmt.Errorf("processEraPoolUpdated EncodeToHash error %s", err)
	}

	msg := w.conn.MessageToSign(tx, poolAddr, txHash)
	signature, err := crypto.Sign(msg[:], key.PrivateKey())
	if err != nil {
		return fmt.Errorf("processEraPoolUpdated sign msg error %s, msg %s", err, hexutil.Encode(msg[:]))
	}
	shareAddrAndCallData := append(shareAddr.Bytes(), tx.CallData...)

	param.ProposalId = crypto.Keccak256Hash(shareAddrAndCallData).Bytes()
	param.Signature = append(msg[:], signature...)

	mef.MultiTransaction = tx
	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}

	err = w.submitMessage(result)
	if err != nil {
		return err
	}

	signatures, err := w.mustGetSignatureFromStafi(&param, uint32(mef.Threshold))
	if err != nil {
		return err
	}

	submitSignatures := submodel.SubmitSignatures{
		Symbol:     param.Symbol,
		Era:        param.Era,
		Pool:       param.Pool,
		TxType:     param.TxType,
		ProposalId: param.ProposalId,
		Signature:  signatures,
		Threshold:  uint32(mef.Threshold),
	}
	return w.processSignatureEnough(&submitSignatures, shareAddrAndCallData, mef, nil)
}

func (w *writer) processBondReported(m *core.Message) error {
	w.log.Info("processBondReported")
	flow, ok := m.Content.(*submodel.BondReportedFlow)
	if !ok {
		w.printContentError(m)
		return fmt.Errorf("cast error")
	}

	snap := flow.Snap
	key := w.conn.FoundKey(flow.SubAccounts)
	if key == nil {
		return fmt.Errorf("processBondReported no keys symbol: %s", m.Destination)
	}

	validatorId := flow.MaticValidatorId

	shareAddr, err := w.conn.GetValidator(validatorId)
	if err != nil {
		return fmt.Errorf("processBondReported get GetValidator error %s", err)
	}

	poolAddr := common.BytesToAddress(snap.Pool)
	tx, err := w.conn.RestakeCall(shareAddr)
	if err != nil {
		return fmt.Errorf("RestakeCall error %s", err)
	}

	param := submodel.SubmitSignatureParams{
		Symbol: flow.Symbol,
		Era:    types.NewU32(snap.Era),
		Pool:   snap.Pool,
		TxType: submodel.OriginalClaimRewards,
	}

	txHash, err := param.EncodeToHash()
	if err != nil {
		return fmt.Errorf("processBondReported EncodeToHash error %s", err)
	}

	msg := w.conn.MessageToSign(tx, poolAddr, txHash)
	signature, err := crypto.Sign(msg[:], key.PrivateKey())
	if err != nil {
		return fmt.Errorf("processBondReported sign msg error %s, msg %s", err, hexutil.Encode(msg[:]))
	}
	shareAddrAndCallData := append(shareAddr.Bytes(), tx.CallData...)

	param.ProposalId = crypto.Keccak256Hash(shareAddrAndCallData).Bytes()
	param.Signature = append(msg[:], signature...)

	flow.MultiTransaction = tx

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	err = w.submitMessage(result)
	if err != nil {
		return err
	}

	signatures, err := w.mustGetSignatureFromStafi(&param, flow.Threshold)
	if err != nil {
		return err
	}

	submitSignatures := submodel.SubmitSignatures{
		Symbol:     param.Symbol,
		Era:        param.Era,
		Pool:       param.Pool,
		TxType:     param.TxType,
		ProposalId: param.ProposalId,
		Signature:  signatures,
		Threshold:  flow.Threshold,
	}
	return w.processSignatureEnough(&submitSignatures, shareAddrAndCallData, nil, flow)

}

func (w *writer) processActiveReported(m *core.Message) error {
	w.log.Info("processActiveReported")
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m)
		return fmt.Errorf("cast err")
	}

	flow, ok := mef.EventData.(*submodel.ActiveReportedFlow)
	if !ok {
		return fmt.Errorf("processActiveReported eventData is not ActiveReportedFlow")
	}

	snap := flow.Snap
	key := w.conn.FoundKey(mef.SubAccounts)
	if key == nil {
		return fmt.Errorf("processActiveReported no keys symbol %s", m.Destination)
	}

	validatorId := mef.MaticValidatorId

	shareAddr, err := w.conn.GetValidator(validatorId)
	if err != nil {
		return fmt.Errorf("processActiveReported get GetValidator error %s", err)
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
		return fmt.Errorf("processActiveReported EncodeToHash error %s", err)
	}

	var nonce *big.Int
	for {
		nonce, err = w.conn.WithdrawNonce(shareAddr, poolAddr)
		if err != nil {
			if err == ErrWithdrawEpochNotMatch {

				// check balance on test net, if balance > total unbond, will skip withdraw call
				if !w.conn.IsMainnet() {
					balance, err := w.conn.BalanceOf(poolAddr)
					if err != nil {
						return err
					}
					if balance.Cmp(mef.TotalUnbondAmount) > 0 {
						w.log.Info("withdrawn", "shareAddr", shareAddr, "poolAddr", poolAddr)
						mef.RunTimeCalls = []*submodel.RunTimeCall{{CallHash: txHash.Hex()}}
						return w.informChain(m.Destination, m.Source, mef)
					}
				}

				// will wait
				w.log.Warn(fmt.Sprintf("WithdrawNonce failed:%s, will wait. shareAddr %s poolAddr %s", err, shareAddr, poolAddr))
				time.Sleep(time.Minute * 2)
				continue
			} else {
				return fmt.Errorf("WithdrawNonce failed:%s, shareAddr %s poolAddr %s", err, shareAddr, poolAddr)
			}
		}
		break
	}

	if nonce.Uint64() == 0 {
		w.log.Info("withdrawn", "shareAddr", shareAddr, "poolAddr", poolAddr)
		mef.RunTimeCalls = []*submodel.RunTimeCall{{CallHash: txHash.Hex()}}
		return w.informChain(m.Destination, m.Source, mef)
	}

	tx, err := w.conn.WithdrawCall(shareAddr, common.BytesToAddress(snap.Pool), nonce)
	if err != nil {
		return fmt.Errorf("processActiveReported WithdrawCall error %s", err)
	}

	msg := w.conn.MessageToSign(tx, poolAddr, txHash)
	signature, err := crypto.Sign(msg[:], key.PrivateKey())
	if err != nil {
		return fmt.Errorf("processActiveReported sign msg error %s msg %s", err, hexutil.Encode(msg[:]))
	}
	shareAddrAndCallData := append(shareAddr.Bytes(), tx.CallData...)

	param.ProposalId = crypto.Keccak256Hash(shareAddrAndCallData).Bytes()
	param.Signature = append(msg[:], signature...)

	mef.MultiTransaction = tx
	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	err = w.submitMessage(result)
	if err != nil {
		return err
	}

	signatures, err := w.mustGetSignatureFromStafi(&param, uint32(mef.Threshold))
	if err != nil {
		return err
	}

	submitSignatures := submodel.SubmitSignatures{
		Symbol:     param.Symbol,
		Era:        param.Era,
		Pool:       param.Pool,
		TxType:     param.TxType,
		ProposalId: param.ProposalId,
		Signature:  signatures,
		Threshold:  uint32(mef.Threshold),
	}
	return w.processSignatureEnough(&submitSignatures, shareAddrAndCallData, mef, nil)
}

func (w *writer) processWithdrawReported(m *core.Message) error {
	w.log.Info("processWithdrawReported")
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m)
		return fmt.Errorf("cast err")
	}

	flow, ok := mef.EventData.(*submodel.WithdrawReportedFlow)
	if !ok {
		return fmt.Errorf("processWithdrawReported eventData is not WithdrawReportedFlow")
	}

	snap := flow.Snap
	key := w.conn.FoundKey(mef.SubAccounts)
	if key == nil {
		return fmt.Errorf("processWithdrawReported no keys symbol %s", m.Destination)
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
		return fmt.Errorf("processWithdrawReported EncodeToHash error %s", err)
	}

	if flow.TotalAmount.Uint64() == 0 {
		mef.RunTimeCalls = []*submodel.RunTimeCall{{CallHash: txHash.Hex()}}
		w.log.Info("processWithdrawReported: no need to do transfer call")
		return w.informChain(m.Destination, m.Source, mef)
	}

	balance, err := w.conn.BalanceOf(poolAddr)
	if err != nil {
		return fmt.Errorf("BalanceOf  error %s address %s", err, poolAddr)
	}

	// check it is already dealed or not
	if balance.Cmp(flow.TotalAmount.Int) < 0 {
		state, err := w.conn.TxHashState(txHash, poolAddr)
		if err != nil {
			return fmt.Errorf("TxHashState error %s, txHash %s poolAddr %s", err, txHash, poolAddr)
		}

		if state != ethmodel.HashStateSuccess {
			return fmt.Errorf("free balance not enough for transfer back symbol %s pool %s balance %s, totalAmount %s", flow.Symbol, poolAddr, balance, flow.TotalAmount)
		}
	}

	tx, err := w.conn.TransferCall(flow.Receives)
	if err != nil {
		return fmt.Errorf("TransferCall error %s", err)
	}

	msg := w.conn.MessageToSign(tx, poolAddr, txHash)
	signature, err := crypto.Sign(msg[:], key.PrivateKey())
	if err != nil {
		return fmt.Errorf("processWithdrawReported sign msg error %s, msg %s", err, hexutil.Encode(msg[:]))
	}
	param.Signature = append(msg[:], signature...)
	shareAddrAndCallData := append(tx.To.Bytes(), tx.CallData...)

	param.ProposalId = crypto.Keccak256Hash(shareAddrAndCallData).Bytes()

	mef.MultiTransaction = tx
	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	err = w.submitMessage(result)
	if err != nil {
		return err
	}

	signatures, err := w.mustGetSignatureFromStafi(&param, uint32(mef.Threshold))
	if err != nil {
		return err
	}

	submitSignatures := submodel.SubmitSignatures{
		Symbol:     param.Symbol,
		Era:        param.Era,
		Pool:       param.Pool,
		TxType:     param.TxType,
		ProposalId: param.ProposalId,
		Signature:  signatures,
		Threshold:  uint32(mef.Threshold),
	}
	return w.processSignatureEnough(&submitSignatures, shareAddrAndCallData, mef, nil)
}

func (w *writer) processSignatureEnough(sigs *submodel.SubmitSignatures, shareAddrAndCallData []byte, mef *submodel.MultiEventFlow, bondFlow *submodel.BondReportedFlow) error {
	w.log.Info("processSignatureEnough", "pool", hexutil.Encode(sigs.Pool), "txType", sigs.TxType)
	if len(shareAddrAndCallData) <= common.AddressLength {
		return fmt.Errorf("processSignatureEnough shareAddrAndCallData too short, symbol %s", sigs.Symbol)
	}

	signatures := make([][]byte, 0)
	for _, sig := range sigs.Signature {
		// 32 + 65 = 97
		if len(sig) != 97 {
			return fmt.Errorf("processSignatureEnough: size of sig not 97 sig: %s", hexutil.Encode(sig))
		}
		signatures = append(signatures, sig[32:])
	}

	vs, rs, ss := utils.DecomposeSignature(signatures)
	to := common.BytesToAddress(shareAddrAndCallData[:20])
	calldata := shareAddrAndCallData[20:]
	poolAddr := common.BytesToAddress(sigs.Pool)

	txHash, err := sigs.EncodeToHash()
	if err != nil {
		return fmt.Errorf("processSignatureEnough sigs EncodeToHash error %s", err)
	}

	var operation ethmodel.CallEnum
	var value, safeGas, totalGas *big.Int
	var subAccounts []types.Bytes
	var era uint32

	txTypeErr := fmt.Errorf("processSignatureEnough TxType %s not supported", sigs.TxType)
	switch sigs.TxType {
	case submodel.OriginalClaimRewards:
		if bondFlow == nil {
			return fmt.Errorf("bondFlow shouldn't be nil")
		}
		operation = bondFlow.MultiTransaction.Operation
		safeGas = bondFlow.MultiTransaction.SafeTxGas
		value = bondFlow.MultiTransaction.Value
		totalGas = bondFlow.MultiTransaction.TotalGas
		subAccounts = bondFlow.SubAccounts
		era = bondFlow.Snap.Era

	case submodel.OriginalBond, submodel.OriginalUnbond, submodel.OriginalWithdrawUnbond, submodel.OriginalTransfer:
		if mef == nil {
			return fmt.Errorf("mef shouldn't be nil")
		}
		mef.RunTimeCalls = []*submodel.RunTimeCall{{CallHash: txHash.Hex()}}
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
			return fmt.Errorf(txTypeErr.Error())
		}

	default:
		return fmt.Errorf(txTypeErr.Error())
	}

	report := func() error {
		switch sigs.TxType {
		case submodel.OriginalBond, submodel.OriginalUnbond:
			flow, _ := mef.EventData.(*submodel.EraPoolUpdatedFlow)
			flow.Active, flow.Reward, err = w.conn.StakedAndReward(txHash, to, poolAddr)
			if err != nil {
				return fmt.Errorf("processSignatureEnough: RewardByTxHash error %s txHash %s pool %s", err, txHash, poolAddr)
			}
			w.log.Info("processSignatureEnough ok", "active", flow.Active, "txHash", txHash, "pool", poolAddr)
			return w.reportMultiEventResult(txHash, mef)
		case submodel.OriginalClaimRewards:
			active, err := w.conn.TotalStaked(to, poolAddr)
			if err != nil {
				return fmt.Errorf("processSignatureEnough TotalStaked error %s shre %s pool %s", err, to, poolAddr)
			}

			snap := bondFlow.Snap
			bond, unbond := snap.Bond.Int, snap.Unbond.Int
			if bond.Cmp(unbond) > 0 {
				diff := big.NewInt(0).Sub(bond, unbond)
				if diff.Cmp(bondFlow.LeastBond) <= 0 {
					active = active.Add(active, diff)
				}
			}
			bondFlow.ReportActive = types.NewU128(*active)
			w.log.Info("processSignatureEnough ok", "pool", poolAddr, "active", active, "txHash", txHash)
			err = w.activeReport(txHash, bondFlow)
			if err != nil {
				return fmt.Errorf("processSignatureEnough activeReport error %s shre %s pool %s", err, to, poolAddr)
			}

			// report rate on evm
			rate, err := w.mustGetEraRateFromStafi(snap.Symbol, types.U32(snap.Era))
			if err != nil {
				return fmt.Errorf("processSignatureEnough mustGetEraRateFromStafi error %s shre %s pool %s", err, to, poolAddr)
			}

			evmRate := new(big.Int).Mul(big.NewInt(int64(rate)), big.NewInt(1e6)) // decimals 12 on stafi, decimals 18 on evm
			proposalId := getProposalId(snap.Era, evmRate, 0)
			proposal, err := w.conn.stakePortalWithRateContract.Proposals(&bind.CallOpts{}, proposalId)
			if err != nil {
				return fmt.Errorf("processSignatureEnough Proposals error %s shre %s pool %s", err, to, poolAddr)
			}
			if proposal.Status == 2 { // success status
				return nil
			}
			hasVoted, err := w.conn.stakePortalWithRateContract.HasVoted(&bind.CallOpts{}, proposalId, w.conn.conn.Opts().From)
			if err != nil {
				return fmt.Errorf("processSignatureEnough HasVoted error %s shre %s pool %s", err, to, poolAddr)
			}
			if hasVoted {
				return nil
			}

			// send tx
			err = w.conn.conn.LockAndUpdateOpts(totalGas, big.NewInt(0))
			if err != nil {
				return fmt.Errorf("processSignatureEnough LockAndUpdateOpts error %s shre %s pool %s", err, to, poolAddr)
			}
			voteTx, err := w.conn.stakePortalWithRateContract.VoteRate(w.conn.conn.Opts(), proposalId, evmRate)
			w.conn.conn.UnlockOpts()

			if err != nil {
				return fmt.Errorf("processSignatureEnough VoteRate error %s shre %s pool %s", err, to, poolAddr)
			}

			err = w.waitTxOk(voteTx.Hash())
			if err != nil {
				return fmt.Errorf("processSignatureEnough waitTxOk error %s shre %s pool %s", err, to, poolAddr)
			}

			return w.waitRateUpdated(proposalId)

		case submodel.OriginalWithdrawUnbond, submodel.OriginalTransfer:
			w.log.Info("processSignatureEnough ok", "pool", poolAddr, "txHash", txHash)
			return w.reportMultiEventResult(txHash, mef)
		default:
			return fmt.Errorf(txTypeErr.Error())
		}
	}

	state, err := w.conn.TxHashState(txHash, poolAddr)
	if err != nil {
		return fmt.Errorf("processSignatureEnough: TxHashState error %s, txHash %s poolAddr %s", err, txHash, poolAddr)
	}

	if state == ethmodel.HashStateSuccess {
		return report()
	}

	eraSignerFlag := w.conn.IsSubmitterOfEra(era, subAccounts)
	if !eraSignerFlag {
		w.log.Info("processSignatureEnough", "eraSignerFlag", eraSignerFlag, "era", era, "len(subAccounts)", len(subAccounts), "txHash", txHash)
		err = w.conn.WaitTxHashSuccess(txHash, poolAddr, sigs.TxType)
		if err != nil {
			return fmt.Errorf("processSignatureEnough: WaitTxHashSuccess error %s txHash %s, pool %s", err, txHash, poolAddr)
		}
		return report()
	}

	err = w.conn.AsMulti(poolAddr, to, value, calldata, uint8(operation), safeGas, totalGas, txHash, vs, rs, ss)
	if err != nil {
		return fmt.Errorf("AsMulti error %s", err)
	}
	w.log.Info("AsMulti success", "txHash", txHash)

	err = w.conn.WaitTxHashSuccess(txHash, poolAddr, sigs.TxType)
	if err != nil {
		return fmt.Errorf("processSignatureEnough: WaitTxHashSuccess error %s txhash %s pool %s", err, txHash, poolAddr)
	}

	return report()
}

func getProposalId(era uint32, rate *big.Int, factor int) common.Hash {
	return crypto.Keccak256Hash([]byte(fmt.Sprintf("era-%d-%s-%s-%d", era, "voteRate", rate.String(), factor)))
}

func (w *writer) reportMultiEventResult(txHash common.Hash, mef *submodel.MultiEventFlow) error {
	return w.informChain(core.RMATIC, core.RFIS, mef)
}

func (w *writer) activeReport(txHash common.Hash, flow *submodel.BondReportedFlow) error {
	msg := &core.Message{Source: core.RMATIC, Destination: core.RFIS, Reason: core.ActiveReport, Content: flow}
	return w.submitMessage(msg)
}

func (w *writer) printContentError(m *core.Message) {
	w.log.Error("msg resolve failed", "source", m.Source, "dest", m.Destination, "reason", m.Reason)
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (w *writer) submitMessage(m *core.Message) error {
	if m.Destination == "" {
		m.Destination = core.RFIS
	}
	return w.router.Send(m)
}

func (w *writer) informChain(source, dest core.RSymbol, flow *submodel.MultiEventFlow) error {
	msg := &core.Message{Source: source, Destination: dest, Reason: core.InformChain, Content: flow}
	return w.submitMessage(msg)
}

func (w *writer) setBondedPools(key common.Address, value bool) {
	w.bondedPoolsMtx.Lock()
	defer w.bondedPoolsMtx.Unlock()
	w.bondedPools[key] = value
}

func (w *writer) start() error {
	go func() {
		for {
			select {
			case <-w.stop:
				close(w.liquidityBonds)
				w.log.Info("get stop signal, stop liquidityBonds handler")
				return
			case msg := <-w.liquidityBonds:
				time.Sleep(5 * time.Second)
				result := w.processLiquidityBond(msg)
				w.log.Info("retry processLiquidityBond", "result", result)
			}
		}
	}()

	return nil
}

func (h *writer) mustGetSignatureFromStafi(param *submodel.SubmitSignatureParams, threshold uint32) (signatures []types.Bytes, err error) {
	flow := submodel.GetSubmitSignaturesFlow{
		Symbol:     param.Symbol,
		Era:        param.Era,
		Pool:       param.Pool,
		TxType:     param.TxType,
		ProposalId: param.ProposalId,
		Signatures: make(chan []types.Bytes, 1),
	}

	for {
		sigs, err := h.getSignatureFromStafi(&flow)
		if err != nil {
			h.log.Debug("getSignatureFromStafiHub failed, will retry.", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}
		if len(sigs) < int(threshold) {
			h.log.Debug("getSignatureFromStafiHub sigs not enough, will retry.", "sigs len", len(sigs), "threshold", threshold)
			time.Sleep(BlockRetryInterval)
			continue
		}
		return sigs, nil
	}
}

func (h *writer) getSignatureFromStafi(param *submodel.GetSubmitSignaturesFlow) (signatures []types.Bytes, err error) {
	msg := core.Message{
		Source:      h.conn.symbol,
		Destination: core.RFIS,
		Reason:      core.GetSubmitSignatures,
		Content:     param,
	}
	err = h.router.Send(&msg)
	if err != nil {
		return nil, err
	}

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	h.log.Debug("wait getSignature from stafihub", "rSymbol", h.conn.symbol)
	select {
	case <-timer.C:
		return nil, fmt.Errorf("get signatures from stafihub timeout")
	case sigs := <-param.Signatures:
		return sigs, nil
	}
}

func (h *writer) mustGetEraRateFromStafi(symbol core.RSymbol, era types.U32) (rate uint64, err error) {
	flow := submodel.GetEraRateFlow{
		Symbol: symbol,
		Era:    era,
		Rate:   make(chan uint64, 1),
	}

	for {
		sigs, err := h.getEraRateFromStafi(&flow)
		if err != nil {
			h.log.Debug("mustGetEraRateFromStafi failed, will retry.", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}
		if sigs == 0 {
			h.log.Debug("mustGetEraRateFromStafi rate zero, will retry.")
			time.Sleep(BlockRetryInterval)
			continue
		}
		return sigs, nil
	}
}

func (h *writer) getEraRateFromStafi(param *submodel.GetEraRateFlow) (rate uint64, err error) {
	msg := core.Message{
		Source:      h.conn.symbol,
		Destination: core.RFIS,
		Reason:      core.GetEraRate,
		Content:     param,
	}
	err = h.router.Send(&msg)
	if err != nil {
		return 0, err
	}

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	h.log.Debug("wait getEraRateFromStafi from stafi", "rSymbol", h.conn.symbol)
	select {
	case <-timer.C:
		return 0, fmt.Errorf("get getEraRateFromStafi from stafi timeout")
	case sigs := <-param.Rate:
		return sigs, nil
	}
}

func (task *writer) waitTxOk(txHash common.Hash) error {
	retry := 0
	for {
		if retry > BlockRetryLimit*3 {
			return fmt.Errorf("networkBalancesContract.SubmitBalances tx reach retry limit")
		}
		_, pending, err := task.conn.conn.TransactionByHash(context.Background(), txHash)
		if err == nil && !pending {
			break
		} else {
			if err != nil {
				task.log.Warn("tx status", "hash", txHash, "err", err.Error())
			} else {
				task.log.Warn("tx status", "hash", txHash, "status", "pending")
			}
			time.Sleep(BlockRetryInterval)
			retry++
			continue
		}

	}
	task.log.Info("tx send ok", "tx", txHash.String())
	return nil
}

func (task *writer) waitRateUpdated(proposalId [32]byte) error {
	retry := 0
	for {
		if retry > BlockRetryLimit*3 {
			return fmt.Errorf("networkBalancesContract.SubmitBalances tx reach retry limit")
		}

		proposal, err := task.conn.stakePortalWithRateContract.Proposals(&bind.CallOpts{}, proposalId)
		if err != nil {
			time.Sleep(BlockRetryInterval)
			retry++
			continue
		}
		if proposal.Status != 2 {
			time.Sleep(BlockRetryInterval)
			retry++
			continue
		}
		break
	}
	return nil
}
