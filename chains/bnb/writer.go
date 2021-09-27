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
	swapRecord      string
	swapHistory     string
}

const (
	bondFlowLimit = 2048
)

func NewWriter(symbol core.RSymbol, opts map[string]interface{}, conn *Connection, log log15.Logger, sysErr chan<- error, stop <-chan int) *writer {
	record, ok := opts["SwapRecord"].(string)
	if !ok {
		panic("no filepath to save SwapRecord")
	}

	history, ok := opts["SwapHistory"].(string)
	if !ok {
		panic("no filepath to save SwapRecord history")
	}

	if _, err := os.Stat(record); os.IsNotExist(err) {
		err = utils.WriteCSV(record, [][]string{})
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(history); os.IsNotExist(err) {
		err = utils.WriteCSV(history, [][]string{})
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
		swapRecord:      record,
		swapHistory:     history,
	}
}

func (w *writer) setRouter(r chains.Router) {
	w.router = r
}

func (w *writer) ResolveMessage(m *core.Message) (processOk bool) {
	defer func() {
		if !processOk {
			panic(fmt.Errorf("resolveMessage process failed. %+v", m))
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

	snap := flow.Snap
	poolAddr := common.BytesToAddress(snap.Pool)
	if !w.conn.IsPoolKeyExist(poolAddr) {
		w.log.Info("has no pool key, will ignore")
		return true
	}

	validatorId, ok := mef.ValidatorId.(bncCmnTypes.ValAddress)
	if !ok {
		w.log.Error("processEraPoolUpdated validatorId not ValAddress")
		return false
	}

	w.log.Info("processEraPoolUpdated infos", "pool", poolAddr, "validator", validatorId)
	unbondable, err := w.conn.Unbondable(poolAddr, validatorId)
	if err != nil {
		w.log.Error("processEraPoolUpdated Unbondable error", "error", err)
		return false
	}

	bond := snap.Bond.Int64()
	unbond := snap.Unbond.Int64()
	diff := bond - unbond
	least := flow.LeastBond.Int64()
	flow.BondCall = &submodel.BondCall{ReportType: submodel.NewBondReport}
	action, amount := w.conn.BondOrUnbondCall(bond, unbond, least)

	lastEra := int64(snap.Era - 1)
	eraBeforeLast := lastEra - 1
	eraBlock := int64(w.conn.EraBlock())
	lastHeight := eraBeforeLast * eraBlock
	curHeight := lastEra * eraBlock
	reward, err := w.conn.Reward(poolAddr, curHeight, lastHeight)
	if err != nil {
		w.log.Error("processEraPoolUpdated last era reward error", "err", err)
		return false
	}

	w.log.Info("processEraPoolUpdated", "action", action, "symbol", snap.Symbol, "era", snap.Era, "reward", reward)

	if bond > 0 && bond > reward && (diff <= 0 || diff >= least) {
		swap := &Swap{Symbol: string(flow.Symbol), Pool: poolAddr.Hex(), Era: fmt.Sprint(flow.Snap.Era), From: FromBsc}
		historied := IsSwapExist(w.swapHistory, swap)
		recorded := IsSwapExist(w.swapRecord, swap)
		w.log.Info("processEraPoolUpdated", "historied", historied, "recorded", recorded)

		swapFun := func() bool {
			futureBal, err := w.conn.TransferFromBscToBc(poolAddr, bond-reward)
			if err != nil {
				w.log.Error("processEraPoolUpdated swap error", "error", err)
				return false
			}

			if err := DeleteSwap(w.swapRecord, swap); err != nil {
				w.log.Error("processEraPoolUpdated delete swap error", "error", err)
				return false
			}

			if err := w.conn.CheckBcBalance(poolAddr, futureBal); err != nil {
				w.log.Info("CheckBcBalance error", "err", err)
				return false
			}

			return true
		}

		if !historied {
			err = CreateSwap(w.swapRecord, swap)
			if err != nil {
				w.log.Error("processEraPoolUpdated: create swap record err", "err", err, "swap", *swap)
				return false
			}

			err = CreateSwap(w.swapHistory, swap)
			if err != nil {
				w.log.Error("processEraPoolUpdated: create swap history err", "err", err, "swap", *swap)
				return false
			}

			if !swapFun() {
				return false
			}
		} else if recorded {
			if !swapFun() {
				return false
			}
		}
	}

	switch action {
	case submodel.BondOnly:
		if err = w.conn.ExecuteBond(poolAddr, validatorId, amount); err != nil {
			w.log.Error("ExecuteBond error", "error", err)
			return false
		}

		flow.BondCall.Action = submodel.BothBondUnbond
	case submodel.UnbondOnly:
		if unbondable {
			err = w.conn.ExecuteUnbond(poolAddr, validatorId, amount)
			if err != nil {
				w.log.Error("ExecuteUnbond error", "error", err)
				return false
			}

			flow.BondCall.Action = submodel.BothBondUnbond
		} else {
			flow.BondCall.Action = submodel.InterDeduct
		}
	case submodel.BothBondUnbond, submodel.EitherBondUnbond, submodel.InterDeduct:
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
	if !w.conn.IsPoolKeyExist(poolAddr) {
		w.log.Info("has no pool key, will ignore", "pool", poolAddr)
		return true
	}

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

	var reward int64
	staked, err := w.conn.Staked(poolAddr, validatorId)
	if err != nil {
		if err.Error() != NoBondError.Error() {
			w.log.Error("processBondReported Staked error", "err", err)
			return false
		}
		reward = 0
	} else {
		reward, err = w.conn.Reward(poolAddr, curHeight, lastHeight)
		if err != nil {
			w.log.Error("processBondReported reward error", "err", err)
			return false
		}
	}

	bond := snap.Bond.Int64()
	unbond := snap.Unbond.Int64()
	if bond > unbond {
		diff := bond - unbond
		if diff < flow.LeastBond.Int64() {
			staked += diff
		}
	} else if unbond > bond {
		diff := unbond - bond
		if diff < flow.LeastBond.Int64() {
			staked -= diff
		}
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
	if !w.conn.IsPoolKeyExist(poolAddr) {
		w.log.Info("has no pool key, will ignore", "pool", poolAddr)
		return true
	}

	swap := &Swap{Symbol: string(snap.Symbol), Pool: poolAddr.Hex(), Era: fmt.Sprint(snap.Era), From: FromBc}
	historied := IsSwapExist(w.swapHistory, swap)
	recorded := IsSwapExist(w.swapRecord, swap)
	w.log.Info("processActiveReported prepare to swap", "pool", poolAddr, "total", total, "historied", historied, "recorded", recorded)

	var err error
	swapFun := func() bool {
		err = w.conn.CheckTransfer(poolAddr, total)
		if err != nil {
			w.log.Error("processActiveReported unable to transfer", "error", err)
			return false
		}

		err = w.conn.TransferFromBcToBsc(poolAddr, total)
		if err != nil {
			w.log.Error("processActiveReported TransferFromBcToBsc error", "error", err)
			return false
		}

		err = DeleteSwap(w.swapRecord, swap)
		if err != nil {
			w.log.Error("processActiveReported delete swap error", "error", err)
			return false
		}

		return true
	}

	if !historied {
		err = CreateSwap(w.swapRecord, swap)
		if err != nil {
			w.log.Error("processActiveReported: create swap record err", "err", err, "swap", *swap)
			return false
		}

		err = CreateSwap(w.swapHistory, swap)
		if err != nil {
			w.log.Error("processActiveReported: create swap history err", "err", err, "swap", *swap)
			return false
		}

		if !swapFun() {
			return false
		}
	} else if recorded {
		if !swapFun() {
			return false
		}
	}

	transformedTotal := big.NewInt(0).Mul(flow.TotalAmount.Int, big.NewInt(1e10))
	txHash, err := w.conn.BatchTransfer(poolAddr, flow.Receives, transformedTotal)
	if err != nil {
		w.log.Error("processActiveReported BatchTransfer error", "error", err)
		return false
	}

	mef.OpaqueCalls = []*submodel.MultiOpaqueCall{{CallHash: txHash.Hex()}}
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
