// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package matic

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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
	symbol          core.RSymbol
	conn            *Connection
	router          chains.Router
	log             log15.Logger
	sysErr          chan<- error
	liquidityBonds  chan *core.Message
	currentChainEra uint32
	bondedPoolsMtx  sync.RWMutex
	bondedPools     map[string]bool
	eventMtx        sync.RWMutex
	events          map[string]*submodel.MultiEventFlow
	bondReportedMtx sync.RWMutex
	bondReporteds   map[string]*submodel.BondReportedFlow
	stop            <-chan int
}

const (
	bondFlowLimit = 2048
)

var (
	UnclaimableHash = crypto.Keccak256Hash([]byte(`unclaimable`))
)

func NewWriter(symbol core.RSymbol, conn *Connection, log log15.Logger, sysErr chan<- error, stop <-chan int) *writer {
	return &writer{
		symbol:          symbol,
		conn:            conn,
		log:             log,
		sysErr:          sysErr,
		liquidityBonds:  make(chan *core.Message, bondFlowLimit),
		currentChainEra: 0,
		bondedPools:     make(map[string]bool),
		events:          make(map[string]*submodel.MultiEventFlow),
		bondReporteds:   make(map[string]*submodel.BondReportedFlow),
		stop:            stop,
	}
}

func (w *writer) setRouter(r chains.Router) {
	w.router = r
}

func (w *writer) ResolveMessage(m *core.Message) bool {
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
		errMsg := "processEraPoolUpdated no keys"
		w.log.Error(errMsg)
		w.sysErr <- errors.New(errMsg)
		return false
	}

	shareAddr, err := w.conn.GetValidator(mef.ValidatorId)
	if err != nil {
		w.log.Error("processEraPoolUpdated get GetValidator error", "error", err)
		w.sysErr <- err
		return false
	}

	method, tx, err := w.conn.BondOrUnbondCall(shareAddr, snap.Bond.Int, snap.Unbond.Int)
	if err != nil {
		if err.Error() == substrate.BondEqualToUnbondError.Error() {
			w.log.Info("No need to send any call", "symbol", snap.Symbol, "era", snap.Era)
			return w.informChain(m.Destination, m.Source, mef)
		}
		w.log.Error("BondOrUnbondCall error", "err", err)
		return false
	}

	poolAddr := common.BytesToAddress(snap.Pool)
	msg, err := w.conn.MessageToSign(tx, poolAddr)
	if err != nil {
		w.log.Error("processEraPoolUpdated MessageToSign error", "err", err)
		return false
	}

	signature, err := crypto.Sign(msg[:], key.PrivateKey())
	if err != nil {
		w.log.Error("processEraPoolUpdated sign msg error", "error", err, "msg", hexutil.Encode(msg[:]))
		w.sysErr <- err
		return false
	}
	signature = append(msg[:], signature...)
	propId := append(shareAddr.Bytes(), tx.CallData...)
	param := submodel.SubmitSignatureParams{
		Symbol:     flow.Symbol,
		Era:        types.NewU32(snap.Era),
		Pool:       snap.Pool,
		TxType:     method,
		ProposalId: propId,
		Signature:  signature,
	}

	txHash, err := param.EncodeToHash()
	if err != nil {
		w.log.Error("processEraPoolUpdated EncodeToHash error", "error", err)
		w.sysErr <- err
		return false
	}
	hash := strings.ToLower(txHash.Hex())
	w.setEvents(hash, mef)

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
		errMsg := "processBondReported no keys"
		w.log.Error(errMsg)
		w.sysErr <- errors.New(errMsg)
		return false
	}

	shareAddr, err := w.conn.GetValidator(flow.ValidatorId)
	if err != nil {
		w.log.Error("processBondReported get GetValidator error", "error", err)
		w.sysErr <- err
		return false
	}

	claimable, err := w.conn.Claimable(shareAddr, common.BytesToAddress(snap.Pool))
	if err != nil {
		w.log.Error("Claimable error", "err", err)
		return false
	}

	if !claimable {
		w.setBondReported(UnclaimableHash.Hex(), flow)
		err = w.queryAndReportActive(UnclaimableHash, shareAddr, m)
		if err != nil {
			w.log.Error("queryAndReportActive error", "err", err)
			return false
		}
		return true
	}

	tx, err := w.conn.RestakeCall(shareAddr)
	if err != nil {
		w.log.Error("RestakeCall error", "err", err)
		return false
	}

	poolAddr := common.BytesToAddress(snap.Pool)
	msg, err := w.conn.MessageToSign(tx, poolAddr)
	if err != nil {
		w.log.Error("processBondReported MessageToSign error", "err", err)
		return false
	}

	signature, err := crypto.Sign(msg[:], key.PrivateKey())
	if err != nil {
		w.log.Error("processBondReported sign msg error", "error", err, "msg", hexutil.Encode(msg[:]))
		w.sysErr <- err
		return false
	}
	signature = append(msg[:], signature...)
	w.log.Info("processBondReported size of signature", "size", len(signature))
	propId := append(shareAddr.Bytes(), tx.CallData...)
	param := submodel.SubmitSignatureParams{
		Symbol:     flow.Symbol,
		Era:        types.NewU32(snap.Era),
		Pool:       snap.Pool,
		TxType:     submodel.OriginalClaimRewards,
		ProposalId: propId,
		Signature:  signature,
	}

	txHash, err := param.EncodeToHash()
	if err != nil {
		w.log.Error("processBondReported EncodeToHash error", "error", err)
		w.sysErr <- err
		return false
	}
	hash := strings.ToLower(txHash.Hex())
	w.setBondReported(hash, flow)
	w.log.Info("processBondReported BondReported set", "txHash", hash)

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
		errMsg := "processActiveReported no keys"
		w.log.Error(errMsg)
		w.sysErr <- errors.New(errMsg)
		return false
	}

	shareAddr, err := w.conn.GetValidator(mef.ValidatorId)
	if err != nil {
		w.log.Error("processActiveReported get GetValidator error", "error", err)
		w.sysErr <- err
		return false
	}

	poolAddr := common.BytesToAddress(snap.Pool)
	withdrawable, err := w.conn.Withdrawable(shareAddr, poolAddr)
	if err != nil {
		w.log.Error("Withdrawable error", "err", err, "shareAddr", shareAddr, "poolAddr", poolAddr)
		return false
	}

	if !withdrawable {
		w.log.Info("no need to withdraw")
		return w.informChain(m.Destination, m.Source, mef)
	}

	tx, err := w.conn.WithdrawCall(shareAddr, common.BytesToAddress(snap.Pool))
	if err != nil {
		w.log.Error("BondOrUnbondCall error", "err", err)
		return false
	}

	msg, err := w.conn.MessageToSign(tx, poolAddr)
	if err != nil {
		w.log.Error("processActiveReported MessageToSign error", "err", err)
		return false
	}

	signature, err := crypto.Sign(msg[:], key.PrivateKey())
	if err != nil {
		w.log.Error("processActiveReported sign msg error", "error", err, "msg", hexutil.Encode(msg[:]))
		w.sysErr <- err
		return false
	}
	signature = append(msg[:], signature...)
	propId := append(shareAddr.Bytes(), tx.CallData...)
	param := submodel.SubmitSignatureParams{
		Symbol:     flow.Symbol,
		Era:        types.NewU32(snap.Era),
		Pool:       snap.Pool,
		TxType:     submodel.OriginalWithdrawUnbond,
		ProposalId: propId,
		Signature:  signature,
	}

	txHash, err := param.EncodeToHash()
	if err != nil {
		w.log.Error("processActiveReported EncodeToHash error", "error", err)
		w.sysErr <- err
		return false
	}
	hash := strings.ToLower(txHash.Hex())
	w.setEvents(hash, mef)

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
		errMsg := "processWithdrawReported no keys"
		w.log.Error(errMsg)
		w.sysErr <- errors.New(errMsg)
		return false
	}

	if flow.TotalAmount.Uint64() == 0 {
		w.log.Info("processWithdrawReported: no need to do transfer call")
		return w.informChain(m.Destination, m.Source, mef)
	}

	poolAddr := common.BytesToAddress(snap.Pool)
	balance, err := w.conn.BalanceOf(poolAddr)
	if err != nil {
		w.log.Error("BalanceOf  error", "err", err, "Address", poolAddr.Hex())
		return false
	}

	if balance.Cmp(flow.TotalAmount.Int) < 0 {
		w.sysErr <- fmt.Errorf("free balance not enough for transfer back, symbol: %s, Address: %s, balance: %s, TotalAmount: %s", flow.Symbol, poolAddr.Hex(), balance.String(), flow.TotalAmount.Int.String())
		return false
	}

	tx, err := w.conn.TransferCall(flow.Receives)
	if err != nil {
		w.log.Error("TransferCall error", "err", err)
		return false
	}

	msg, err := w.conn.MessageToSign(tx, poolAddr)
	if err != nil {
		w.log.Error("processWithdrawReported MessageToSign error", "err", err)
		return false
	}

	signature, err := crypto.Sign(msg[:], key.PrivateKey())
	if err != nil {
		w.log.Error("processWithdrawReported sign msg error", "error", err, "msg", hexutil.Encode(msg[:]))
		w.sysErr <- err
		return false
	}
	signature = append(msg[:], signature...)
	propId := append(tx.To.Bytes(), tx.CallData...)
	param := submodel.SubmitSignatureParams{
		Symbol:     flow.Symbol,
		Era:        types.NewU32(snap.Era),
		Pool:       snap.Pool,
		TxType:     submodel.OriginalWithdrawUnbond,
		ProposalId: propId,
		Signature:  signature,
	}

	txHash, err := param.EncodeToHash()
	if err != nil {
		w.log.Error("processWithdrawReported EncodeToHash error", "error", err)
		w.sysErr <- err
		return false
	}
	hash := strings.ToLower(txHash.Hex())
	w.setEvents(hash, mef)

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

	txHash, err := sigs.EncodeToHash()
	if err != nil {
		w.log.Error("processSignatureEnough sigs EncodeToHash error", "error", err)
		w.sysErr <- err
		return false
	}
	hash := strings.ToLower(txHash.Hex())

	var mef *submodel.MultiEventFlow

	if sigs.TxType == submodel.OriginalClaimRewards {
		_, ok = w.getBondReported(hash)
		if !ok {
			w.log.Error("processSignatureEnough: no bond report for txHash", "txHash", hash)
			return false
		}
	} else {
		mef, ok = w.getEvents(hash)
		if !ok {
			w.log.Error("processSignatureEnough: no event for txHash", "txHash", hash)
			return false
		}
	}

	poolAddr := common.BytesToAddress(sigs.Pool)
	executed, err := w.conn.IsTxHashExecuted(txHash, poolAddr)
	if err != nil {
		w.log.Error("processSignatureEnough IsTxHashExecuted error", "error", err, "txHash", txHash)
		w.sysErr <- err
		return false
	}

	if executed {
		w.log.Warn("TxHash already executed", "txHash", hash)
		return true
	}

	if len(sigs.ProposalId) <= common.AddressLength {
		errMsg := "processSignatureEnough ProposalId too short"
		w.log.Error(errMsg)
		w.sysErr <- errors.New(errMsg)
		return false
	}

	signatures := make([][]byte, 0)
	for _, sig := range sigs.Signature {
		// 32 + 65 = 97
		if len(sig) != 97 {
			err := fmt.Errorf("processSignatureEnough: size of sig %s not 97", hexutil.Encode(sig))
			w.log.Error(err.Error())
			w.sysErr <- err
			return false
		}
		signatures = append(signatures, sig[32:])
	}

	to := common.BytesToAddress(sigs.ProposalId[:20])
	calldata := sigs.ProposalId[20:]
	msg := sigs.Signature[0][:32]
	if !w.conn.IsFirstSigner(msg, signatures[0]) {
		w.log.Info("processSignatureEnough: not first signer, will ignore")
	}

	vs, rs, ss := utils.DecomposeSignature(signatures)
	callType := config.Call
	var safeTxGas *big.Int
	switch sigs.TxType {
	case submodel.OriginalBond:
		safeTxGas = BuyVoucherSafeTxGas
	case submodel.OriginalUnbond:
		safeTxGas = SellVoucherNewSafeTxGas
	case submodel.OriginalClaimRewards:
		safeTxGas = RestakeSafeTxGas
	case submodel.OriginalWithdrawUnbond:
		safeTxGas = WithdrawTxGas
	case submodel.OriginalTransfer:
		safeTxGas = TransferTxGas
	default:
		errMsg := "processSignatureEnough TxType not supported"
		w.log.Error(errMsg)
		w.sysErr <- errors.New(errMsg)
		return false
	}

	err = w.conn.AsMulti(poolAddr, to, DefaultValue, calldata, callType, safeTxGas, txHash, vs, rs, ss)
	if err != nil {
		w.log.Error("AsMulti error", "err", err)
		return false
	}

	w.log.Info("AsMulti success", "txHash", txHash.Hex())
	err = w.checkAndSend(txHash, sigs, to, mef, m)
	if err != nil {
		w.log.Error("processSignatureEnough: checkAndSend error", "txHash", txHash.Hex(), "to", to.Hex())
		return false
	}
	return true
}

func (w *writer) checkAndSend(txHash common.Hash, sigs *submodel.SubmitSignatures, to common.Address, mef *submodel.MultiEventFlow, m *core.Message) error {
	poolAddr := common.BytesToAddress(sigs.Pool)
	err := w.conn.CheckTxHash(txHash, poolAddr)
	if err != nil {
		return err
	}

	if sigs.TxType == submodel.OriginalClaimRewards {
		return w.queryAndReportActive(txHash, to, m)
	}

	mef.OpaqueCalls = []*submodel.MultiOpaqueCall{{CallHash: txHash.Hex()}}
	if w.informChain(m.Destination, m.Source, mef) {
		w.deleteEvents(txHash.Hex())
	}
	return nil
}

func (w *writer) queryAndReportActive(txHash common.Hash, share common.Address, m *core.Message) error {
	hash := strings.ToLower(txHash.Hex())
	flow, ok := w.getBondReported(hash)
	if !ok {
		err := fmt.Errorf("queryAndReportActive no bondReportedFlow, txHash: %s", txHash.Hex())
		w.log.Error(err.Error())
		w.sysErr <- err
		return err
	}

	active, err := w.conn.TotalStaked(share, common.BytesToAddress(flow.Snap.Pool))
	if err != nil {
		w.log.Error("queryAndReportActive TotalStaked error", "error", err)
		return err
	}

	flow.Snap.Active = types.NewU128(*active)
	w.log.Info("queryAndReportActive", "pool", hexutil.Encode(flow.Snap.Pool), "active", flow.Snap.Active)
	msg := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.ActiveReport, Content: flow}
	_ = w.submitMessage(msg)
	return nil
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

func (w *writer) getBondReported(key string) (*submodel.BondReportedFlow, bool) {
	w.bondReportedMtx.RLock()
	defer w.bondReportedMtx.RUnlock()
	value, exist := w.bondReporteds[key]
	return value, exist
}

func (w *writer) setBondReported(key string, value *submodel.BondReportedFlow) {
	w.bondReportedMtx.Lock()
	defer w.bondReportedMtx.Unlock()
	w.bondReporteds[key] = value
}

func (w *writer) deleteBondReported(key string) {
	w.bondReportedMtx.Lock()
	defer w.bondReportedMtx.Unlock()
	delete(w.bondReporteds, key)
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
