// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"math/big"
	"os"
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
	symbol                core.RSymbol
	conn                  *Connection
	router                chains.Router
	log                   log15.Logger
	sysErr                chan<- error
	eventMtx              sync.RWMutex
	newMulTicsMtx         sync.RWMutex
	bondedPoolsMtx        sync.RWMutex
	events                map[string]*submodel.MultiEventFlow
	newMultics            map[string]*submodel.EventNewMultisig
	bondedPools           map[string]bool
	liquidityBonds        chan *core.Message
	currentChainEra       uint32
	stop                  <-chan int
	transferRecord        string
	transferRecordHistory string
}

const (
	bondFlowLimit = 2048
)

func NewWriter(symbol core.RSymbol, opts map[string]interface{}, conn *Connection, log log15.Logger, sysErr chan<- error, stop <-chan int) *writer {
	transferRecord := ""
	transferRecordHistory := ""
	if symbol != core.RFIS {
		path, ok := opts["transferRecord"].(string)
		if !ok {
			panic("no filepath to save transferRecord")
		}

		historyPath, ok := opts["transferRecordHistory"].(string)
		if !ok {
			panic("no filepath to save transferRecord history")
		}
		transferRecord = path
		transferRecordHistory = historyPath
		if _, err := os.Stat(transferRecord); os.IsNotExist(err) {
			err = utils.WriteCSV(transferRecord, [][]string{})
			if err != nil {
				panic(err)
			}
		}

		if _, err := os.Stat(transferRecordHistory); os.IsNotExist(err) {
			err = utils.WriteCSV(transferRecordHistory, [][]string{})
			if err != nil {
				panic(err)
			}
		}
	}

	return &writer{
		symbol:                symbol,
		conn:                  conn,
		log:                   log,
		sysErr:                sysErr,
		events:                make(map[string]*submodel.MultiEventFlow),
		newMultics:            make(map[string]*submodel.EventNewMultisig),
		bondedPools:           make(map[string]bool),
		liquidityBonds:        make(chan *core.Message, bondFlowLimit),
		currentChainEra:       0,
		stop:                  stop,
		transferRecord:        transferRecord,
		transferRecordHistory: transferRecordHistory,
	}
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
	case core.TransferReportedEvent:
		return w.processTransferReportedEvent(m)
	case core.NominationUpdatedEvent:
		return w.processNominationUpdatedEvent(m)
	case core.GetEraNominated:
		return w.processGetEraNominated(m)
	case core.SubmitSignature:
		return w.processSubmitSignature(m)
	default:
		w.log.Warn("message reason unsupported", "reason", m.Reason)
		return false
	}
}

func (w *writer) processGetEraNominated(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.GetEraNominatedFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	validator, err := w.conn.GetEraNominated(flow.Symbol, flow.Pool, flow.Era)
	if err != nil {
		w.log.Warn("GetEraNominated failed", "err", err)
		if err.Error() == NotExistError.Error() {
			flow.NewValidators <- make([]types.AccountID, 0)
			return true
		}
	}
	flow.NewValidators <- validator
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
	w.log.Info("", "len(bondedPools)", len(bondedPools), "symbol", m.Source)
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

	if data, ok := flow.EventData.(*submodel.NominationUpdatedFlow); ok {
		w.log.Info("NominationUpdated", "symbol", data.Symbol, "era", data.Era)
		return true
	}

	var prop *submodel.Proposal
	var err error
	if data, ok := flow.EventData.(*submodel.EraPoolUpdatedFlow); ok {
		call := data.BondCall
		switch data.Symbol {
		case core.RMATIC:
			switch call.ReportType {
			case submodel.NewBondReport:
				prop, err = w.conn.NewBondReportProposal(data)
				if err != nil {
					w.log.Error("MethodNewBondReportProposal", "error", err)
					return false
				}
			case submodel.BondAndReportActive:
				prop, err = w.conn.BondAndReportActiveProposal(data)
				if err != nil {
					w.log.Error("MethodBondAndReportActiveProposal", "error", err)
					return false
				}
			default:
				w.log.Error("processInformChain: ReportType not supported", "ReportType", call.ReportType)
				return false
			}
		case core.RBNB:
			switch call.ReportType {
			case submodel.NewBondReport:
				prop, err = w.conn.NewBondReportProposal(data)
				if err != nil {
					w.log.Error("MethodNewBondReportProposal", "error", err)
					return false
				}
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

	if data, ok := flow.EventData.(*submodel.ActiveReportedFlow); ok {
		prop, err = w.conn.CommonReportProposal(config.MethodWithdrawReport, m.Source, bondId, data.ShotId)
		if err != nil {
			w.log.Error("MethodWithdrawReportProposal", "error", err)
			return false
		}
		result := w.conn.resolveProposal(prop, true)
		w.log.Info("MethodWithdrawReportProposal resolveProposal", "result", result)
		return result
	}

	if data, ok := flow.EventData.(*submodel.WithdrawReportedFlow); ok {
		prop, err = w.conn.CommonReportProposal(config.MethodTransferReport, m.Source, bondId, data.ShotId)
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
		if err.Error() == TargetNotExistError.Error() {
			w.sysErr <- fmt.Errorf("TargetNotExistError, pool: %s, symbol: %s, lastEra: %d", hexutil.Encode(flow.Snap.Pool), flow.Symbol, flow.LastEra)
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

func (w *writer) queryAndReportActive(m *core.Message) {
	flow, ok := m.Content.(*submodel.BondReportedFlow)
	if !ok {
		w.printContentError(m)
		return
	}

	ledger := new(submodel.StakingLedger)
	exist, err := w.conn.QueryStorage(config.StakingModuleId, config.StorageLedger, flow.Snap.Pool, nil, ledger)
	if err != nil {
		w.log.Error("queryAndReportActive get ledger error", "error", err, "pool", hexutil.Encode(flow.Snap.Pool))
		return
	}
	if !exist {
		w.log.Error("queryAndReportActive ledger not exist", "pool", hexutil.Encode(flow.Snap.Pool))
		return
	}

	flow.Snap.Active = types.NewU128(big.Int(ledger.Active))
	w.log.Info("queryAndReportActive", "pool", hexutil.Encode(flow.Snap.Pool), "active", flow.Snap.Active)
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

	var prop *submodel.Proposal
	var err error
	if !flow.NewActiveReportFlag {
		prop, err = w.conn.ActiveReportProposal(flow)
		if err != nil {
			w.log.Error("ActiveReportProposal", "error", err)
			return false
		}
	} else {
		prop, err = w.conn.NewActiveReportProposal(flow)
		if err != nil {
			w.log.Error("NewActiveReportProposal", "error", err)
			return false
		}
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

func (w *writer) start() error {
	if w.symbol == core.RFIS {
		return nil
	}

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
