// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package bnb

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
)

var (
	tenDecimals = decimal.NewFromInt(1e10)
)

type writer struct {
	symbol         core.RSymbol
	conn           *Connection
	router         chains.Router
	eraSeconds     uint64
	eraOffset      int64
	log            core.Logger
	sysErr         chan<- error
	bondedPoolsMtx sync.RWMutex
	bondedPools    map[string]bool
	stop           <-chan int
}

func NewWriter(symbol core.RSymbol, eraSeconds uint64, eraOffset int64, conn *Connection, log core.Logger, sysErr chan<- error, stop <-chan int) *writer {

	return &writer{
		symbol:      symbol,
		conn:        conn,
		eraSeconds:  eraSeconds,
		eraOffset:   eraOffset,
		log:         log,
		sysErr:      sysErr,
		bondedPools: make(map[string]bool),
		stop:        stop,
	}
}

func (w *writer) setRouter(r chains.Router) {
	w.router = r
}

func (w *writer) ResolveMessage(m *core.Message) (processOk bool) {

	var err error
	switch m.Reason {
	case core.LiquidityBondEvent:
		err = w.processLiquidityBond(m)
	case core.BondedPools:
		err = w.processBondedPools(m)
	case core.EraPoolUpdatedEvent:
		err = w.processEraPoolUpdated(m)
	case core.ActiveReportedEvent:
		err = w.processActiveReported(m)
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
		return fmt.Errorf("content cast failed")
	}

	if flow.Reason != submodel.BondReasonDefault {
		return fmt.Errorf("processLiquidityBond receive a message of which reason is not default bondId:%s reason:%s symbol: %s", flow.BondId.Hex(), flow.Reason, flow.Symbol)
	}

	bondReason, err := w.conn.TransferVerify(flow.Record)
	if err != nil {
		return errors.Wrap(err, "TransferVerify")
	}

	w.log.Info("processLiquidityBond", "bondId", flow.BondId.Hex(), "bondReason", bondReason, "VerifyTimes", flow.VerifyTimes)
	flow.Reason = bondReason

	msg := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.LiquidityBondResult, Content: flow}
	return w.submitMessage(msg)
}

func (w *writer) processBondedPools(m *core.Message) error {
	pools, ok := m.Content.([]types.Bytes)
	if !ok {
		w.printContentError(m)
		return fmt.Errorf("content cast failed")
	}

	for _, p := range pools {
		w.log.Info("processBondedPools", "pool", utiles.AddHex(hexutil.Encode(p)))
		w.setBondedPools(hexutil.Encode(p), true)
	}

	return nil
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
func (w *writer) processEraPoolUpdated(m *core.Message) error {
	w.log.Info("processEraPoolUpdated")
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		return fmt.Errorf("content cast failed")
	}

	flow, ok := mef.EventData.(*submodel.EraPoolUpdatedFlow)
	if !ok {
		return fmt.Errorf("processEraPoolUpdated HeadFlow is not EraPoolUpdatedFlow")
	}

	snap := flow.Snap
	poolAddr := common.BytesToAddress(snap.Pool)

	var pool *Pool
	var exist bool
	if pool, exist = w.conn.GetPool(poolAddr); !exist {
		return fmt.Errorf("has no pool key, will ignore")
	}

	targetHeight, err := w.conn.GetHeightByEra(snap.Era, int64(w.eraSeconds), w.eraOffset)
	if err != nil {
		return errors.Wrap(err, "GetHeightByEra")
	}
	lastEraHeight, err := w.conn.GetHeightByEra(snap.Era-1, int64(w.eraSeconds), w.eraOffset)
	if err != nil {
		return errors.Wrap(err, "GetHeightByEra")
	}

	minDelegation, err := w.conn.stakingContract.GetMinDelegation(&bind.CallOpts{
		From:        poolAddr,
		BlockNumber: big.NewInt(targetHeight),
		Context:     context.Background(),
	})
	if err != nil {
		return errors.Wrap(err, "GetMinDelegation")
	}

	newRewadOnBc, err := w.conn.RewardOnBc(poolAddr, targetHeight, lastEraHeight)
	if err != nil {
		return errors.Wrap(err, "RewardOnBc")
	}

	pendingStakeDeci := decimal.NewFromBigInt(flow.PendingStake, 0)
	pendingRewardDeci := decimal.NewFromBigInt(flow.PendingStake, 0)
	bondDeci := decimal.NewFromBigInt(snap.Bond.Int, 0).Mul(tenDecimals) // decimals is 8 on bc and 18 on bsc
	unbondDeci := decimal.NewFromBigInt(snap.Unbond.Int, 0).Mul(tenDecimals)
	leastBondDeci := decimal.NewFromBigInt(minDelegation, 0)
	pendingRewardDeci = pendingRewardDeci.Add(decimal.NewFromInt(newRewadOnBc))

	//-------- claim reward on bsc
	rewardOnBsc, err := w.conn.stakingContract.GetDistributedReward(&bind.CallOpts{
		From:        poolAddr,
		BlockNumber: big.NewInt(targetHeight),
		Context:     context.Background(),
	}, poolAddr)
	if err != nil {
		return errors.Wrap(err, "stakingContract.GetDistributedReward")
	}

	if rewardOnBsc.Sign() > 0 {

		proposalId := getProposalId(snap.Era, "processEraPoolUpdated", "claimReward", 0)

		needSend, err := needSendProposal(pool, proposalId)
		if err != nil {
			return errors.Wrap(err, "needSendProposal")
		}
		if needSend {
			proposalBts, err := w.getClaimRewardProposal()
			if err != nil {
				return errors.Wrap(err, "getClaimRewardProposal")
			}
			err = w.submitProposal(pool, proposalId, proposalBts)
			if err != nil {
				return errors.Wrap(err, "submitProposal")
			}
			err = w.waitProposalExecuted(pool, proposalId)
			if err != nil {
				return errors.Wrap(err, "waitProposalExecuted")
			}
			realRewardClaimed, claimRewardTxHeight, err := w.findRealRewardAmountClaimed(pool, proposalId, poolAddr, uint64(targetHeight))
			if err != nil {
				return errors.Wrap(err, "findRealRewardClaimed")
			}

			realRewardOnBscDeci := decimal.NewFromBigInt(realRewardClaimed, 0)
			pendingRewardDeci = pendingRewardDeci.Sub(realRewardOnBscDeci)
			pendingStakeDeci = pendingStakeDeci.Add(realRewardOnBscDeci)

			targetHeight = int64(claimRewardTxHeight)
		}
	}

	//---- claim undelegated
	undelegatedAmount, err := w.conn.stakingContract.GetUndelegated(&bind.CallOpts{
		From:        poolAddr,
		BlockNumber: big.NewInt(targetHeight),
		Context:     context.Background(),
	}, poolAddr)
	if err != nil {
		return errors.Wrap(err, "stakingContract.GetUndelegated")
	}
	if undelegatedAmount.Sign() > 0 {
		proposalId := getProposalId(snap.Era, "processEraPoolUpdated", "claimUndelegated", 0)

		needSend, err := needSendProposal(pool, proposalId)
		if err != nil {
			return errors.Wrap(err, "needSendProposal")
		}
		if needSend {
			proposalBts, err := w.getClaimUndelegatedProposal()
			if err != nil {
				return errors.Wrap(err, "getClaimUndelegatedProposal")
			}

			err = w.submitProposal(pool, proposalId, proposalBts)
			if err != nil {
				return errors.Wrap(err, "submitProposal")
			}
			err = w.waitProposalExecuted(pool, proposalId)
			if err != nil {
				return errors.Wrap(err, "waitProposalExecuted")
			}
			_, undelegatedClaimTxHeight, err := w.findRealUndelegatedAmountClaimed(pool, proposalId, poolAddr, uint64(targetHeight))
			if err != nil {
				return errors.Wrap(err, "findRealRewardClaimed")
			}

			targetHeight = int64(undelegatedClaimTxHeight)
		}
	}

	//------ balance of pool on target height
	poolBalance, err := pool.bscClient.Client().BalanceAt(context.Background(), poolAddr, big.NewInt(targetHeight))
	if err != nil {
		return errors.Wrap(err, "pool balance get failed")
	}
	poolBalanceDeci := decimal.NewFromBigInt(poolBalance, 0)

	relayerFee, err := w.conn.stakingContract.GetRelayerFee(&bind.CallOpts{
		From:        poolAddr,
		BlockNumber: big.NewInt(targetHeight),
		Context:     context.Background(),
	})
	if err != nil {
		return errors.Wrap(err, "stakingContract.GetRelayerFee")
	}
	relayerFeeDeci := decimal.NewFromBigInt(relayerFee, 0)

	//------- switch cal
	willDelegateAmountDeci := decimal.Zero
	willUnDelegateAmountDeci := decimal.Zero
	diffDeci := decimal.Zero
	bondAction := submodel.BothBondUnbond
	switch bondDeci.Cmp(unbondDeci) {
	case 1:
		diffDeci = bondDeci.Sub(unbondDeci)
		pendingStakeDeci = pendingStakeDeci.Add(diffDeci)
		if pendingStakeDeci.GreaterThan(leastBondDeci) && poolBalanceDeci.GreaterThan(leastBondDeci) {
			willDelegateAmountDeci = decimal.Min(pendingStakeDeci, poolBalanceDeci)
			pendingStakeDeci = pendingStakeDeci.Sub(willDelegateAmountDeci)
		}
	case -1:
		diffDeci = unbondDeci.Sub(bondDeci)
		if pendingStakeDeci.GreaterThanOrEqual(diffDeci) {
			pendingStakeDeci = pendingStakeDeci.Sub(diffDeci)
			if pendingStakeDeci.GreaterThan(leastBondDeci) && poolBalanceDeci.GreaterThan(leastBondDeci) {
				willDelegateAmountDeci = decimal.Min(pendingStakeDeci, poolBalanceDeci)
				pendingStakeDeci = pendingStakeDeci.Sub(willDelegateAmountDeci)
			}
		} else {
			tempUnDelegateAmountDeci := diffDeci.Sub(pendingStakeDeci).Div(leastBondDeci).Ceil().Mul(leastBondDeci)
			unbondable, err := w.unbondable(tempUnDelegateAmountDeci, relayerFeeDeci, leastBondDeci, mef.BnbValidators, poolAddr, targetHeight)
			if err != nil {
				return errors.Wrap(err, "unbondable")
			}
			if unbondable {
				willUnDelegateAmountDeci = tempUnDelegateAmountDeci
				pendingStakeDeci = pendingStakeDeci.Add(willUnDelegateAmountDeci).Sub(diffDeci)
			} else {
				bondAction = submodel.EitherBondUnbond
			}
		}

	case 0:
		if pendingStakeDeci.GreaterThan(leastBondDeci) && poolBalanceDeci.GreaterThan(leastBondDeci) {
			willDelegateAmountDeci = decimal.Min(pendingStakeDeci, poolBalanceDeci)
			pendingStakeDeci = pendingStakeDeci.Sub(willDelegateAmountDeci)
		}

	default:
		return fmt.Errorf("unknown cmp result")
	}

	// ----- delegate
	if willDelegateAmountDeci.IsPositive() {
		proposalId := getProposalId(snap.Era, "processEraPoolUpdated", "delegate", 0)
		needSend, err := needSendProposal(pool, proposalId)
		if err != nil {
			return errors.Wrap(err, "needSendProposal")
		}
		if needSend {
			proposalBts, err := w.getDelegateProposal(willDelegateAmountDeci, relayerFeeDeci, leastBondDeci, poolAddr, mef.BnbValidators, targetHeight)
			if err != nil {
				return errors.Wrap(err, "getDelegateProposal")
			}

			err = w.submitProposal(pool, proposalId, proposalBts)
			if err != nil {
				return errors.Wrap(err, "submitProposal")
			}
			err = w.waitProposalExecuted(pool, proposalId)
			if err != nil {
				return errors.Wrap(err, "waitProposalExecuted")
			}
		}
	}

	// ----- unDelegate
	if willUnDelegateAmountDeci.IsPositive() {
		proposalId := getProposalId(snap.Era, "processEraPoolUpdated", "unDelegate", 0)
		needSend, err := needSendProposal(pool, proposalId)
		if err != nil {
			return errors.Wrap(err, "needSendProposal")
		}
		if needSend {
			proposalBts, err := w.getUnDelegateProposal(willUnDelegateAmountDeci, relayerFeeDeci, leastBondDeci, mef.BnbValidators, poolAddr, targetHeight)
			if err != nil {
				return errors.Wrap(err, "getUnDelegateProposal")
			}

			err = w.submitProposal(pool, proposalId, proposalBts)
			if err != nil {
				return errors.Wrap(err, "submitProposal")
			}
			err = w.waitProposalExecuted(pool, proposalId)
			if err != nil {
				return errors.Wrap(err, "waitProposalExecuted")
			}
		}
	}

	// ----- bond and active report with pending value
	staked, err := w.staked(pool, poolAddr)
	if err != nil {
		return errors.Wrap(err, "get total staked failed")
	}
	stakedDeci := decimal.NewFromBigInt(staked, 0)
	activeDeci := stakedDeci.Add(pendingStakeDeci).Add(pendingRewardDeci)
	if bondAction == submodel.EitherBondUnbond {
		activeDeci = activeDeci.Sub(diffDeci)
	}
	flow.Active = activeDeci.Div(tenDecimals).BigInt() // decimals 8
	flow.PendingStake = pendingStakeDeci.BigInt()      // decimals 18
	flow.PendingReward = pendingRewardDeci.BigInt()    // decimals 18
	flow.BondCall = &submodel.BondCall{
		ReportType: submodel.BondAndReportActiveWithPendingValue,
		Action:     bondAction,
	}

	return w.informChain(m.Destination, m.Source, mef)
}

func (w *writer) processActiveReported(m *core.Message) error {
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m)
		return fmt.Errorf("content cast failed")
	}

	flow, ok := mef.EventData.(*submodel.WithdrawReportedFlow)
	if !ok {
		return fmt.Errorf("processActiveReported eventData is not ActiveReportedFlow")
	}

	snap := flow.Snap
	poolAddr := common.BytesToAddress(snap.Pool)
	var pool *Pool
	var exist bool
	if pool, exist = w.conn.GetPool(poolAddr); !exist {
		return fmt.Errorf("has no pool key, will ignore pool: %s", poolAddr)
	}
	proposalId := getProposalId(snap.Era, "processActiveReported", "transfer", 0)
	needSend, err := needSendProposal(pool, proposalId)
	if err != nil {
		return errors.Wrap(err, "needSendProposal")
	}
	if needSend {
		proposalBts, totalAmount, err := w.getTransferProposal(poolAddr, flow.Receives)
		if err != nil {
			return errors.Wrap(err, "getTransferProposal")
		}

		for {
			poolBalance, err := w.conn.queryClient.Client().BalanceAt(context.Background(), poolAddr, nil)
			if err != nil {
				return errors.Wrap(err, "BalanceAt")
			}
			poolBalanceDeci := decimal.NewFromBigInt(poolBalance, 0)

			if poolBalanceDeci.LessThanOrEqual(totalAmount) {
				needSend, err = needSendProposal(pool, proposalId)
				if err != nil {
					return errors.Wrap(err, "needSendProposal")
				}
				if needSend {
					time.Sleep(WaitInterval)
					w.log.Warn("pool balance not enough will wait", "pool", poolAddr.String(), "balance", poolBalanceDeci.String(), "totalTransferAmount", totalAmount)
					continue
				}
				break
			}
			break
		}
		err = w.submitProposal(pool, proposalId, proposalBts)
		if err != nil {
			return errors.Wrap(err, "submitProposal")
		}
		err = w.waitProposalExecuted(pool, proposalId)
		if err != nil {
			return errors.Wrap(err, "waitProposalExecuted")
		}

	}

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.InformChain, Content: mef}
	return w.submitMessage(result)
}

func (w *writer) printContentError(m *core.Message) {
	w.log.Error("msg resolve failed", "source", m.Source, "dest", m.Destination, "reason", m.Reason)
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (w *writer) submitMessage(m *core.Message) error {
	if len(m.Destination) == 0 {
		m.Destination = core.RFIS
	}
	return w.router.Send(m)
}

func (w *writer) informChain(source, dest core.RSymbol, flow *submodel.MultiEventFlow) error {
	msg := &core.Message{Source: source, Destination: dest, Reason: core.InformChain, Content: flow}
	return w.submitMessage(msg)
}

func (w *writer) setBondedPools(key string, value bool) {
	w.bondedPoolsMtx.Lock()
	defer w.bondedPoolsMtx.Unlock()
	w.bondedPools[key] = value
}
