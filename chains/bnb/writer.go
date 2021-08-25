// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package bnb

import (
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	bncCmnTypes "github.com/stafiprotocol/go-sdk/common/types"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
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
	stop            <-chan int
	rewardsRecord   string
}

const (
	bondFlowLimit = 2048
)

func NewWriter(symbol core.RSymbol, opts map[string]interface{}, conn *Connection, log log15.Logger, sysErr chan<- error, stop <-chan int) *writer {
	rewardsRecord, ok := opts["rewardsRecord"].(string)
	if !ok {
		panic("no filepath to save rewardsRecord")
	}

	if _, err := os.Stat(rewardsRecord); os.IsNotExist(err) {
		err = utils.WriteCSV(rewardsRecord, [][]string{})
		if err != nil {
			panic(err)
		}
	}

	return &writer{
		symbol:          symbol,
		conn:            conn,
		log:             log,
		sysErr:          sysErr,
		liquidityBonds:  make(chan *core.Message, bondFlowLimit),
		currentChainEra: 0,
		bondedPools:     make(map[string]bool),
		stop:            stop,
		rewardsRecord:   rewardsRecord,
	}
}

func (w *writer) setRouter(r chains.Router) {
	w.router = r
}

func (w *writer) ResolveMessage(m *core.Message) (processOk bool) {
	defer func() {
		if !processOk {
			w.sysErr <- fmt.Errorf("resolveMessage process failed. %+v", m)
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

	if w.conn.BscTransferToBc(flow.Record) != nil {
		w.log.Error("BscTransferToBc error", "error", err)
		return false
	}

	msg := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.LiquidityBondResult, Content: flow}
	return w.submitMessage(msg)
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

	validatorId, ok := mef.ValidatorId.(bncCmnTypes.ValAddress)
	if !ok {
		w.log.Error("processEraPoolUpdated validatorId not ValAddress")
		return false
	}

	snap := flow.Snap
	poolAddr := common.BytesToAddress(snap.Pool)
	unbondable, err := w.conn.Unbondable(poolAddr, validatorId)
	if err != nil {
		w.log.Error("processEraPoolUpdated Unbondable error", "pool", poolAddr, "validator", validatorId.String())
		return false
	}

	bond := snap.Bond.Int64()
	unbond := snap.Unbond.Int64()
	least := flow.LeastBond.Int64()
	flow.BondCall = &submodel.BondCall{ReportType: submodel.NewBondReport}
	if unbondable {
		flow.BondCall.Action = submodel.EitherBondUnbond
		return w.informChain(m.Destination, m.Source, mef)
	}

	action, amount := w.conn.BondOrUnbondCall(bond, unbond, least)
	w.log.Info("processEraPoolUpdated", "action", action, "symbol", snap.Symbol, "era", snap.Era)
	switch action {
	case submodel.BondOnly:
		err := w.conn.ExecuteBond(poolAddr, validatorId, amount)
		if err != nil {
			w.log.Error("ExecuteBond error", "error", err)
			return false
		}
		flow.BondCall.Action = submodel.BothBondUnbond
	case submodel.UnbondOnly:
		err := w.conn.ExecuteUnbond(poolAddr, validatorId, amount)
		if err != nil {
			w.log.Error("ExecuteUnbond error", "error", err)
			return false
		}
		flow.BondCall.Action = submodel.BothBondUnbond
	case submodel.BothBondUnbond, submodel.EitherBondUnbond:
		w.log.Info("processEraPoolUpdated: no need to bond or unbond")
		flow.BondCall.Action = action
	}

	return w.informChain(m.Destination, m.Source, mef)
}

func (w *writer) processBondReported(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondReportedFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	snap := flow.Snap
	poolAddr := common.BytesToAddress(snap.Pool)
	validatorId, ok := flow.ValidatorId.(bncCmnTypes.ValAddress)
	if !ok {
		w.log.Error("processBondReported validatorId not ValAddress")
		return false
	}

	eraBlock := w.conn.EraBlock()
	lastHeight := int64(0)
	if flow.LastEra != 0 {
		lastHeight = int64(flow.LastEra) * int64(eraBlock)
	}
	curHeight := int64(snap.Era) * int64(eraBlock)
	w.log.Info("processBondReported reward", "pool", poolAddr, "validator", validatorId.String(), "curHeight", curHeight, "lastHeight", lastHeight)

	reward, err := w.conn.Reward(poolAddr, curHeight, lastHeight)
	if err != nil {
		w.log.Error("processBondReported reward error", "err", err)
		return false
	}

	staked, err := w.conn.Staked(poolAddr, validatorId)
	if err != nil {
		w.log.Error("processBondReported reward error", "err", err)
		return false
	}

	flow.Snap.Active = types.NewU128(*big.NewInt(staked))
	flow.Unstaked = types.NewU128(*big.NewInt(reward))
	flow.NewActiveReportFlag = true
	w.log.Info("queryAndReportActive", "pool", poolAddr, "active", staked, "unstaked", reward)

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.ActiveReport, Content: flow}
	return w.submitMessage(result)
}

func (w *writer) processActiveReported(m *core.Message) bool {
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	flow, ok := mef.EventData.(*submodel.WithdrawReportedFlow)
	if !ok {
		w.log.Error("processActiveReported eventData is not ActiveReportedFlow")
		return false
	}

	snap := flow.Snap
	total := flow.TotalAmount.Int64()
	poolAddr := common.BytesToAddress(snap.Pool)
	w.log.Info("processActiveReported prepare to transfer", "pool", poolAddr, "total", total)
	err := w.conn.CheckTransfer(poolAddr, total)
	if err != nil {
		w.log.Error("processActiveReported unable to transfer", "error", err)
		return false
	}

	err = w.conn.TransferFromBcToBsc(poolAddr, total)
	if err != nil {
		w.log.Error("processActiveReported TransferFromBcToBsc error", "error", err)
		return false
	}

	err = w.conn.BatchTransfer(poolAddr, flow.Receives, total)
	if err != nil {
		w.log.Error("processActiveReported BatchTransfer error", "error", err)
		return false
	}

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.InformChain, Content: mef}
	return w.submitMessage(result)
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
