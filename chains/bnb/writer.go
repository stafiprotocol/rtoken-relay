// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package bnb

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
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
	bondedPools     map[string]bool
	stop            <-chan int
}

const (
	bondFlowLimit = 2048
)

func NewWriter(symbol core.RSymbol, opts map[string]interface{}, conn *Connection, log core.Logger, sysErr chan<- error, stop <-chan int) *writer {

	return &writer{
		symbol:          symbol,
		conn:            conn,
		log:             log,
		sysErr:          sysErr,
		liquidityBonds:  make(chan *core.Message, bondFlowLimit),
		currentChainEra: 0,
		bondedPools:     make(map[string]bool),
		stop:            stop,
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
	case core.LiquidityBondEvent:
		return w.processLiquidityBond(m)
	case core.BondedPools:
		return w.processBondedPools(m)
	case core.EraPoolUpdatedEvent:
		return w.processEraPoolUpdated(m)
	case core.ActiveReportedEvent:
		return w.processActiveReported(m)
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

/*
targetHeight = heightOfEra(currentEra)
pendingReward += newRewardOnBc
if rewardOnBsc > 0:

	reward = claimReward()
	pendingReward -= reward
	pendingStake += reward
	targetHeight = claimRewardTxHeight

if undelegatedAmountOnBsc > 0:

	claimUndelegated()
	targetHeight = claimUndelegatedTxHeight

poolBalance = balanceOfHeight(targetHeight)

willDelegateAmount = 0
willUnDelegateAmount = 0
bondAction = bondBondUnbond
switch {
case bond > unbond:

	diff = bond-unbond
	pendingStake += diff
	if (pendingStake > leastBond) && (poolBalance > leastBond):
		willDelegateAmount = min(pendingStake, poolBalance)
		pendingStake -= willDelegateAmount

case bond < unbond:

	diff = unbond - bond
	if pendingStake >= diff:
		pendingStake -= diff
		if (pendingStake > leastBond) && (poolBalance > leastBond):
			willDelegateAmount = min(pendingStake, poolBalance)
			pendingStake -= willDelegateAmount
	else:
		if unBondable:
			willUnDelegateAmount = ceil((diff - pendingStake)/leastBond) * leastBond
			pendingStake = pendingStake + willUnDelegateAmount - diff
		else:
			bondAction = eitherBondUnbond

case bond == unbond:

	if (pendingStake > leastBond) && (poolBalance > leastBond):
		willDelegateAmount = min(pendingStake, poolBalance)
		pendingStake -= willDelegateAmount

}

if willDelegateAmount > 0:

	delegate(willDelegateAmount)

if willUnDelegateAmount > 0:

	unDelegate(willUnDelegateAmount)

active = staked + pendingStake + pendingReward
if action == eitherBondUnbond:

	active -= diff

return bond_and_active_report_with_pending_value(action, active, pendingStake, pendingReward)
*/
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

	return w.informChain(m.Destination, m.Source, mef)
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
	poolAddr := common.BytesToAddress(snap.Pool)
	if !w.conn.IsPoolKeyExist(poolAddr) {
		w.log.Info("has no pool key, will ignore", "pool", poolAddr)
		return true
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
				time.Sleep(2 * time.Second)
				result := w.processLiquidityBond(msg)
				w.log.Info("retry processLiquidityBond", "result", result)
			}
		}
	}()

	return nil
}
