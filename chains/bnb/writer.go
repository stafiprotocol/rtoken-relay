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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/itering/scale.go/utiles"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
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

	pendingStakeDeci := decimal.NewFromBigInt(flow.PendingStake, 0)
	pendingRewardDeci := decimal.NewFromBigInt(flow.PendingStake, 0)
	bondDeci := decimal.NewFromBigInt(snap.Bond.Int, 0)
	unbondDeci := decimal.NewFromBigInt(snap.Unbond.Int, 0)
	leastBondDeci := decimal.NewFromBigInt(flow.LeastBond, 0)

	newRewadOnBc, err := w.conn.RewardOnBc(poolAddr, targetHeight, lastEraHeight)
	if err != nil {
		return errors.Wrap(err, "RewardOnBc")
	}
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
		proposalId := crypto.Keccak256Hash([]byte(fmt.Sprintf("era-%d-processEraPoolUpdated-claimReward", snap.Era)))

		needSend, err := needSendProposal(pool, proposalId)
		if err != nil {
			return errors.Wrap(err, "needSendProposal")
		}
		if needSend {
			proposalBts := []byte{}
			err := w.submitProposal(pool, proposalId, proposalBts)
			if err != nil {
				return errors.Wrap(err, "submitProposal")
			}
			err = w.waitProposalExecuted(pool, proposalId)
			if err != nil {
				return errors.Wrap(err, "waitProposalExecuted")
			}
			realRewardClaimed, claimRewardTxHeight, err := w.findRealRewardClaimed(pool, proposalId, poolAddr, uint64(targetHeight))
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
		proposalId := crypto.Keccak256Hash([]byte(fmt.Sprintf("era-%d-processEraPoolUpdated-claimUndelegated", snap.Era)))

		needSend, err := needSendProposal(pool, proposalId)
		if err != nil {
			return errors.Wrap(err, "needSendProposal")
		}
		if needSend {
			proposalBts := []byte{}
			err := w.submitProposal(pool, proposalId, proposalBts)
			if err != nil {
				return errors.Wrap(err, "submitProposal")
			}
			err = w.waitProposalExecuted(pool, proposalId)
			if err != nil {
				return errors.Wrap(err, "waitProposalExecuted")
			}
			_, undelegatedClaimTxHeight, err := w.findRealUndelegatedClaimed(pool, proposalId, poolAddr, uint64(targetHeight))
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
	//------- switch
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
			unbondable, err := w.unbondable(pool)
			if err != nil {
				return errors.Wrap(err, "unbondable")
			}
			if unbondable {
				willUnDelegateAmountDeci = diffDeci.Sub(pendingStakeDeci).Div(leastBondDeci).Ceil().Mul(leastBondDeci)
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
		proposalId := crypto.Keccak256Hash([]byte(fmt.Sprintf("era-%d-processEraPoolUpdated-delegate", snap.Era)))
		proposalBts := []byte{}
		err := w.submitProposal(pool, proposalId, proposalBts)
		if err != nil {
			return errors.Wrap(err, "submitProposal")
		}
		err = w.waitProposalExecuted(pool, proposalId)
		if err != nil {
			return errors.Wrap(err, "waitProposalExecuted")
		}
	}
	// ----- unDelegate
	if willUnDelegateAmountDeci.IsPositive() {
		proposalId := crypto.Keccak256Hash([]byte(fmt.Sprintf("era-%d-processEraPoolUpdated-unDelegate", snap.Era)))
		proposalBts := []byte{}
		err := w.submitProposal(pool, proposalId, proposalBts)
		if err != nil {
			return errors.Wrap(err, "submitProposal")
		}
		err = w.waitProposalExecuted(pool, proposalId)
		if err != nil {
			return errors.Wrap(err, "waitProposalExecuted")
		}
	}

	// ----- bond and active report with pending value
	staked, err := w.staked(pool)
	if err != nil {
		return errors.Wrap(err, "get total staked failed")
	}
	stakedDeci := decimal.NewFromBigInt(staked, 0)
	activeDeci := stakedDeci.Add(pendingStakeDeci).Add(pendingRewardDeci)
	if bondAction == submodel.EitherBondUnbond {
		activeDeci = activeDeci.Sub(diffDeci)
	}
	flow.Active = activeDeci.BigInt()
	flow.PendingStake = pendingStakeDeci.BigInt()
	flow.PendingReward = pendingRewardDeci.BigInt()
	flow.BondCall = &submodel.BondCall{
		ReportType: submodel.BondAndReportActiveWithPendingValue,
		Action:     bondAction,
	}

	return w.informChain(m.Destination, m.Source, mef)
}

func needSendProposal(pool *Pool, proposalId [32]byte) (bool, error) {
	proposal, err := pool.multisigOnchain.Proposals(&bind.CallOpts{}, proposalId)
	if err != nil {
		return false, errors.Wrap(err, "multisigOnchain.Proposals")
	}
	needSend := false
	switch proposal.Status {
	case 0:
		needSend = true
	case 1:
		voted, err := pool.multisigOnchain.HasVoted(&bind.CallOpts{}, proposalId, pool.bscClient.Opts().From)
		if err != nil {
			return false, errors.Wrap(err, "multisigOnchain.HasVoted")
		}
		if !voted {
			needSend = true
		}
	case 2:
	default:
		return false, fmt.Errorf("unknown proposal status: %d", proposal.Status)
	}
	return needSend, nil
}

func (w *writer) unbondable(pool *Pool) (bool, error) {
	// w.conn.stakingContract.GetPendingUndelegateTime()
	return false, nil
}

func (w *writer) staked(pool *Pool) (*big.Int, error) {
	return nil, nil
}

func (w *writer) findRealRewardClaimed(pool *Pool, proposalId [32]byte, poolAddr common.Address, targetHeight uint64) (*big.Int, uint64, error) {
	proposalExectedIterator, err := pool.multisigOnchain.FilterProposalExecuted(&bind.FilterOpts{
		Start:   targetHeight,
		Context: context.Background(),
	})
	if err != nil {
		return nil, 0, errors.Wrap(err, "multisigOnchain.FilterProposalExecuted")
	}
	for proposalExectedIterator.Next() {
		if proposalExectedIterator.Event.ProposalId == proposalId {
			rewardClaimedIterator, err := w.conn.stakingContract.FilterRewardClaimed(&bind.FilterOpts{
				Start:   proposalExectedIterator.Event.Raw.BlockNumber,
				End:     &proposalExectedIterator.Event.Raw.BlockNumber,
				Context: context.Background(),
			}, []common.Address{poolAddr})
			if err != nil {
				return nil, 0, errors.Wrap(err, "stakingContract.FilterRewardClaimed")
			}
			for rewardClaimedIterator.Next() {
				if rewardClaimedIterator.Event.Raw.TxHash == proposalExectedIterator.Event.Raw.TxHash {
					return rewardClaimedIterator.Event.Amount, rewardClaimedIterator.Event.Raw.BlockNumber, nil
				}
			}
		}
	}
	return nil, 0, fmt.Errorf("not find reward claim event")
}

func (w *writer) findRealUndelegatedClaimed(pool *Pool, proposalId [32]byte, poolAddr common.Address, targetHeight uint64) (*big.Int, uint64, error) {
	proposalExectedIterator, err := pool.multisigOnchain.FilterProposalExecuted(&bind.FilterOpts{
		Start:   targetHeight,
		Context: context.Background(),
	})
	if err != nil {
		return nil, 0, errors.Wrap(err, "multisigOnchain.FilterProposalExecuted")
	}
	for proposalExectedIterator.Next() {
		if proposalExectedIterator.Event.ProposalId == proposalId {
			rewardClaimedIterator, err := w.conn.stakingContract.FilterUndelegatedClaimed(&bind.FilterOpts{
				Start:   proposalExectedIterator.Event.Raw.BlockNumber,
				End:     &proposalExectedIterator.Event.Raw.BlockNumber,
				Context: context.Background(),
			}, []common.Address{poolAddr})
			if err != nil {
				return nil, 0, errors.Wrap(err, "stakingContract.FilterUndelegatedClaimed")
			}
			for rewardClaimedIterator.Next() {
				if rewardClaimedIterator.Event.Raw.TxHash == proposalExectedIterator.Event.Raw.TxHash {
					return rewardClaimedIterator.Event.Amount, rewardClaimedIterator.Event.Raw.BlockNumber, nil
				}
			}
		}
	}
	return nil, 0, fmt.Errorf("not find undelegated claim event")
}

func (w *writer) submitProposal(pool *Pool, proposalId [32]byte, proposalBts []byte) error {
	err := pool.bscClient.LockAndUpdateOpts(big.NewInt(0), big.NewInt(0))
	if err != nil {
		return errors.Wrap(err, "LockAndUpdateOpts")
	}
	defer pool.bscClient.UnlockOpts()

	tx, err := pool.multisigOnchain.ExecTransactions(pool.bscClient.Opts(), proposalId, proposalBts)
	if err != nil {
		return errors.Wrap(err, "multisigOnchain.ExecTransactions")
	}
	retry := 0
	for {
		if retry > GetRetryLimit*2 {
			return fmt.Errorf("multisigOnchain.ExecTransactions tx reach retry limit")
		}
		_, pending, err := pool.bscClient.Client().TransactionByHash(context.Background(), tx.Hash())
		if err == nil && !pending {
			break
		} else {
			if err != nil {
				w.log.Warn("tx status", "hash", tx.Hash(), "err", err.Error())
			} else {
				w.log.Warn("tx status", "hash", tx.Hash(), "status", "pending")
			}
			time.Sleep(WaitInterval)
			retry++
			continue
		}
	}
	return nil
}

func (w *writer) waitProposalExecuted(pool *Pool, proposalId [32]byte) error {
	retry := 0
	for {
		if retry > GetRetryLimit*6 {
			return fmt.Errorf("networkBalancesContract.SubmitBalances tx reach retry limit")
		}

		proposal, err := pool.multisigOnchain.Proposals(&bind.CallOpts{}, proposalId)
		if err != nil {
			w.log.Warn("get proposal failed, will retry", "err", err.Error(), "proposalId", proposalId)
			time.Sleep(WaitInterval)
			retry++
			continue
		}
		if proposal.Status != 2 {
			w.log.Warn("proposals not exexted, will wait", "proposalId", proposalId)
			time.Sleep(WaitInterval)
			retry++
			continue
		}
	}
	return nil
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
	_ = pool

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
