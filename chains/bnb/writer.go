// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package bnb

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	multisigOnchain "github.com/stafiprotocol/rtoken-relay/bindings/MultisigOnchain"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
	"github.com/stafiprotocol/rtoken-relay/utils"
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

func (w *writer) start() error {
	// check rate on polygon and evm
	go func() {
		retry := 0
		for {
			if retry > GetRetryLimit {
				w.log.Warn("check rate on bsc reach retry limit")
				return
			}
			latestBlockTimestamp, err := w.conn.LatestBlockTimestamp()
			if err != nil {
				w.log.Warn("Unable to get latest block", "err", err)
				time.Sleep(6 * time.Second)
				retry++
				continue
			}

			era := int64(latestBlockTimestamp/w.eraSeconds) + w.eraOffset
			if era <= 0 {
				w.log.Error("check rate on bsc", "err", fmt.Errorf("era must > 0: %d", era))
				time.Sleep(6 * time.Second)
				retry++
				continue
			}

			rate, err := w.mustGetEraRateFromStafi(core.RBNB, types.U32(era))
			if era <= 0 {
				w.log.Error("check rate on bsc", "err", err)
				time.Sleep(6 * time.Second)
				retry++
				continue
			}

			evmRate := new(big.Int).Mul(big.NewInt(int64(rate)), big.NewInt(1e6)) // decimals 12 on stafi, decimals 18 on evm
			proposalId := getRateProposalId(uint32(era), evmRate, 0)

			// evm rate
			if w.conn.bscStakePortalRateContract != nil {
				poolAddrStr := ""
				for pool := range w.bondedPools {
					poolAddrStr = pool
					break
				}

				poolAddr := common.HexToAddress(poolAddrStr)

				multisigOnchainContract, err := multisigOnchain.NewMultisigOnchain(poolAddr, w.conn.queryClient.Client())
				if err != nil {
					w.log.Error("multisigOnchain.NewMultisigOnchain", "err", err)
					return
				}

				localAccountClients, err := w.conn.GetAccountClients(multisigOnchainContract)
				if err != nil {
					w.log.Error("GetAccountClients", "err", err)
					return
				}
				if len(localAccountClients) == 0 {
					w.log.Error("subAccounts not exist", "err", err)
					return
				}

				for _, client := range localAccountClients {
					err := w.bscVoteRate(client, proposalId, evmRate)
					if err != nil {
						w.log.Error("check rate on bsc bscVoteRate failed", "err", err)
						return
					}
				}

				err = w.waitBscRateUpdated(proposalId)
				if err != nil {
					w.log.Error("check rate on bsc waitBscRateUpdated failed", "err", err)
					return
				}
			}

			break
		}
	}()

	return nil
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
	w.log.Info("processEraPoolUpdated start")
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

	w.log.Info("processEraPoolUpdated detail", "era", flow.Snap.Era, "pool", poolAddr.String(), "snapshot", fmt.Sprintf("%+v", snap))

	multisigOnchainContract, err := multisigOnchain.NewMultisigOnchain(poolAddr, w.conn.queryClient.Client())
	if err != nil {
		return errors.Wrap(err, "multisigOnchain.NewMultisigOnchain")
	}

	localAccountClients, err := w.conn.GetAccountClients(multisigOnchainContract)
	if err != nil {
		return errors.Wrap(err, "GetAccountClients")
	}
	if len(localAccountClients) == 0 {
		return fmt.Errorf("subAccounts not exist")
	}

	w.log.Debug("GetHeightByEra")
	targetHeight, targetTimestamp, err := w.conn.GetHeightTimestampByEra(snap.Era, int64(w.eraSeconds), w.eraOffset)
	if err != nil {
		return errors.Wrap(err, "GetHeightByEra")
	}
	preEra := snap.Era - 1

	// Should add the rewad of era 1368 1369 1370 during the era of 1371, which were skipped by wrong reward address on bc
	if snap.Era == 1371 {
		preEra = 1367
	}
	_, preEraTimestamp, err := w.conn.GetHeightTimestampByEra(preEra, int64(w.eraSeconds), w.eraOffset)
	if err != nil {
		return errors.Wrap(err, "GetHeightByEra")
	}

	minDelegation, err := w.conn.GetMinDelegation(poolAddr, targetHeight)
	if err != nil {
		return errors.Wrap(err, "GetMinDelegation")
	}

	// get reward form pre era to this era
	w.log.Debug("get reward fomr pre era to this era")
	newRewadOnBc, err := w.conn.RewardOnBcDuTimes(poolAddr, targetTimestamp, preEraTimestamp)
	if err != nil {
		return errors.Wrap(err, "RewardOnBc")
	}

	pendingStakeDeci := decimal.NewFromBigInt(flow.PendingStake, 0)   // decimals 18
	pendingRewardDeci := decimal.NewFromBigInt(flow.PendingReward, 0) // decimals 18
	leastBondDeci := decimal.NewFromBigInt(minDelegation, 0)          // decimals 18

	// decimals is 8 on bc and 18 on bsc
	bondDeci := decimal.NewFromBigInt(snap.Bond.Int, 0).Mul(tenDecimals)
	unbondDeci := decimal.NewFromBigInt(snap.Unbond.Int, 0).Mul(tenDecimals)
	newRewadOnBcDeci := decimal.NewFromBigInt(big.NewInt(newRewadOnBc), 0).Mul(tenDecimals)

	pendingRewardDeci = pendingRewardDeci.Add(newRewadOnBcDeci)

	//-------- claim reward on bsc
	w.log.Debug("claim reward on bsc", "staking-contract-minDelegate", leastBondDeci.StringFixed(0))
	rewardOnBsc, err := w.conn.GetDistributedReward(poolAddr, targetHeight)
	if err != nil {
		return errors.Wrap(err, "stakingContract.GetDistributedReward")
	}

	if rewardOnBsc.Sign() > 0 {
		proposalId := getProposalId(snap.Era, "processEraPoolUpdated", "claimReward", 0)
		for _, client := range localAccountClients {
			needSend, err := needSendProposal(client, multisigOnchainContract, proposalId)
			if err != nil {
				return errors.Wrap(err, "needSendProposal")
			}
			if needSend {
				proposalBts, err := w.getClaimRewardProposal()
				if err != nil {
					return errors.Wrap(err, "getClaimRewardProposal")
				}
				err = w.submitProposal(client, multisigOnchainContract, proposalId, proposalBts)
				if err != nil {
					if !strings.Contains(err.Error(), "proposal already executed") {
						return errors.Wrap(err, "submitProposal")
					}
				}
			}
		}

		err = w.waitProposalExecuted(multisigOnchainContract, proposalId)
		if err != nil {
			return errors.Wrap(err, "waitProposalExecuted")
		}
		realRewardClaimed, claimRewardTxHeight, err := w.findRealRewardAmountClaimed(multisigOnchainContract, proposalId, poolAddr, uint64(targetHeight))
		if err != nil {
			return errors.Wrap(err, "findRealRewardClaimed")
		}

		realRewardOnBscDeci := decimal.NewFromBigInt(realRewardClaimed, 0)
		pendingRewardDeci = pendingRewardDeci.Sub(realRewardOnBscDeci)
		pendingStakeDeci = pendingStakeDeci.Add(realRewardOnBscDeci)

		// pool balance change at this height
		targetHeight = int64(claimRewardTxHeight)
	}

	//---- claim undelegated
	w.log.Debug("claim undelegated")
	undelegatedAmount, err := w.conn.GetUndelegated(poolAddr, targetHeight)
	if err != nil {
		return errors.Wrap(err, "stakingContract.GetUndelegated")
	}
	if undelegatedAmount.Sign() > 0 {
		proposalId := getProposalId(snap.Era, "processEraPoolUpdated", "claimUndelegated", 0)

		for _, client := range localAccountClients {

			needSend, err := needSendProposal(client, multisigOnchainContract, proposalId)
			if err != nil {
				return errors.Wrap(err, "needSendProposal")
			}
			if needSend {
				proposalBts, err := w.getClaimUndelegatedProposal()
				if err != nil {
					return errors.Wrap(err, "getClaimUndelegatedProposal")
				}

				err = w.submitProposal(client, multisigOnchainContract, proposalId, proposalBts)
				if err != nil {
					if !strings.Contains(err.Error(), "proposal already executed") {
						return errors.Wrap(err, "submitProposal")
					}
				}
			}
		}

		err = w.waitProposalExecuted(multisigOnchainContract, proposalId)
		if err != nil {
			return errors.Wrap(err, "waitProposalExecuted")
		}
		_, undelegatedClaimTxHeight, err := w.findRealUndelegatedAmountClaimed(multisigOnchainContract, proposalId, poolAddr, uint64(targetHeight))
		if err != nil {
			return errors.Wrap(err, "findRealRewardClaimed")
		}

		// pool balance change at this height
		targetHeight = int64(undelegatedClaimTxHeight)
	}

	//------ balance of pool on target height
	w.log.Debug("balance of pool on target height")
	poolBalance, err := w.conn.BalanceAt(poolAddr, big.NewInt(targetHeight))
	if err != nil {
		return errors.Wrap(err, "pool balance get failed")
	}
	poolBalanceDeci := decimal.NewFromBigInt(poolBalance, 0)

	relayerFee, err := w.conn.GetRelayerFee(poolAddr, targetHeight)
	if err != nil {
		return errors.Wrap(err, "stakingContract.GetRelayerFee")
	}
	relayerFeeDeci := decimal.NewFromBigInt(relayerFee, 0)

	bscRelayerFee, err := w.conn.BSCRelayerFee(poolAddr, targetHeight)
	if err != nil {
		return errors.Wrap(err, "stakingContract.BSCRelayerFee")
	}
	bscRelayerFeeDeci := decimal.NewFromBigInt(bscRelayerFee, 0)

	//------- switch cal
	w.log.Debug("switch cal")
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
			unbondable, err := w.unbondable(tempUnDelegateAmountDeci, relayerFeeDeci, bscRelayerFeeDeci, leastBondDeci, mef.BnbValidators, poolAddr, targetHeight)
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
		w.log.Info("will delegate", "amount", willDelegateAmountDeci.StringFixed(0), "pool", poolAddr.String())
		proposalId := getProposalId(snap.Era, "processEraPoolUpdated", "delegate", 0)
		proposalBts, distributedValidators, err := w.getDelegateProposal(willDelegateAmountDeci, relayerFeeDeci, leastBondDeci, poolAddr, mef.BnbValidators, targetHeight)
		if err != nil {
			return errors.Wrap(err, "getDelegateProposal")
		}
		for _, client := range localAccountClients {

			needSend, err := needSendProposal(client, multisigOnchainContract, proposalId)
			if err != nil {
				return errors.Wrap(err, "needSendProposal")
			}
			if needSend {
				err = w.submitProposal(client, multisigOnchainContract, proposalId, proposalBts)
				if err != nil {
					if !strings.Contains(err.Error(), "proposal already executed") {
						return errors.Wrap(err, "submitProposal")
					}
				}
			}
		}

		err = w.waitProposalExecuted(multisigOnchainContract, proposalId)
		if err != nil {
			return errors.Wrap(err, "waitProposalExecuted")
		}
		err = w.waitDelegateCrossChainOk(poolAddr, proposalId, uint64(targetHeight), distributedValidators)
		if err != nil {
			return errors.Wrap(err, "waitDelegateCrossChainOk")
		}
	}

	// ----- unDelegate
	if willUnDelegateAmountDeci.IsPositive() {
		w.log.Info("will undelegate", "amount", willUnDelegateAmountDeci.StringFixed(0), "pool", poolAddr.String())
		proposalId := getProposalId(snap.Era, "processEraPoolUpdated", "unDelegate", 0)
		proposalBts, selectedValidator, err := w.getUnDelegateProposal(willUnDelegateAmountDeci, relayerFeeDeci, bscRelayerFeeDeci, leastBondDeci, mef.BnbValidators, poolAddr, targetHeight)
		if err != nil {
			return errors.Wrap(err, "getUnDelegateProposal")
		}

		for _, client := range localAccountClients {
			needSend, err := needSendProposal(client, multisigOnchainContract, proposalId)
			if err != nil {
				return errors.Wrap(err, "needSendProposal")
			}

			if needSend {
				err = w.submitProposal(client, multisigOnchainContract, proposalId, proposalBts)
				if err != nil {
					if !strings.Contains(err.Error(), "proposal already executed") {
						return errors.Wrap(err, "submitProposal")
					}
				}
			}
		}

		err = w.waitProposalExecuted(multisigOnchainContract, proposalId)
		if err != nil {
			return errors.Wrap(err, "waitProposalExecuted")
		}
		err = w.waitUnDelegateCrossChainOk(poolAddr, proposalId, uint64(targetHeight), selectedValidator)
		if err != nil {
			return errors.Wrap(err, "waitUnDelegateCrossChainOk")
		}
	}

	// ----- bond and active report with pending value
	w.log.Debug("bond and active report with pending value")
	// got current total staked amount
	staked, err := w.stakedAndCheck(poolAddr, willDelegateAmountDeci.BigInt(), willUnDelegateAmountDeci.BigInt(), targetHeight)
	if err != nil {
		return errors.Wrap(err, "stakedAndCheck failed")
	}
	stakedDeci := decimal.NewFromBigInt(staked, 0)
	if stakedDeci.IsNegative() {
		stakedDeci = decimal.Zero
	}
	if pendingStakeDeci.IsNegative() {
		pendingStakeDeci = decimal.Zero
	}
	if pendingRewardDeci.IsNegative() {
		pendingRewardDeci = decimal.Zero
	}

	activeDeci := stakedDeci.Add(pendingStakeDeci).Add(pendingRewardDeci)
	if bondAction == submodel.EitherBondUnbond {
		activeDeci = activeDeci.Sub(diffDeci)
	}
	if activeDeci.IsNegative() {
		activeDeci = decimal.Zero
	}

	flow.Active = activeDeci.Div(tenDecimals).BigInt() // decimals 8
	flow.PendingStake = pendingStakeDeci.BigInt()      // decimals 18
	flow.PendingReward = pendingRewardDeci.BigInt()    // decimals 18
	flow.BondCall = &submodel.BondCall{
		ReportType: submodel.BondAndReportActiveWithPendingValue,
		Action:     bondAction,
	}
	w.log.Info("will informChain", "reportType", "BondAndReportActiveWithPendingValue", "action", bondAction, "pool", poolAddr.String(),
		"active", flow.Active, "pendingStake", flow.PendingStake, "pendingReward", flow.PendingReward)

	err = w.informChain(m.Destination, m.Source, mef)
	if err != nil {
		return errors.Wrap(err, "informChain")
	}

	// report rate on bsc
	rate, err := w.mustGetEraRateFromStafi(snap.Symbol, types.U32(snap.Era))
	if err != nil {
		return fmt.Errorf("processSignatureEnough mustGetEraRateFromStafi error %s pool %s", err, poolAddr)
	}

	evmRate := new(big.Int).Mul(big.NewInt(int64(rate)), big.NewInt(1e6)) // decimals 12 on stafi, decimals 18 on evm
	proposalId := getRateProposalId(snap.Era, evmRate, 0)

	// evm rate
	if w.conn.bscStakePortalRateContract != nil {
		for _, client := range localAccountClients {
			err := w.bscVoteRate(client, proposalId, evmRate)
			if err != nil {
				return err
			}
		}

		err = w.waitBscRateUpdated(proposalId)
		if err != nil {
			return fmt.Errorf("waitRateUpdated error %s", err)
		}
	}

	return nil
}

func (w *writer) processActiveReported(m *core.Message) error {
	w.log.Info("processActiveReported start")
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		return fmt.Errorf("content cast failed")
	}

	flow, ok := mef.EventData.(*submodel.WithdrawReportedFlow)
	if !ok {
		return fmt.Errorf("processActiveReported eventData is not ActiveReportedFlow")
	}

	snap := flow.Snap
	poolAddr := common.BytesToAddress(snap.Pool)

	multisigOnchainContract, err := multisigOnchain.NewMultisigOnchain(poolAddr, w.conn.queryClient.Client())
	if err != nil {
		return errors.Wrap(err, "multisigOnchain.NewMultisigOnchain")
	}
	localAccountClients, err := w.conn.GetAccountClients(multisigOnchainContract)
	if err != nil {
		return errors.Wrap(err, "GetAccountClients")
	}
	if len(localAccountClients) == 0 {
		return fmt.Errorf("subAccounts not exist")
	}

	proposalId := getProposalId(snap.Era, "processActiveReported", "transfer", 0)
	for _, client := range localAccountClients {

		needSend, err := needSendProposal(client, multisigOnchainContract, proposalId)
		if err != nil {
			return errors.Wrap(err, "needSendProposal")
		}
		proposalBts, totalAmountDeci, err := w.getTransferProposal(poolAddr, flow.Receives)
		if err != nil {
			return errors.Wrap(err, "getTransferProposal")
		}
		w.log.Info("processActiveReported detail", "poolAddr", poolAddr, "proposalId", hex.EncodeToString(proposalId[:]),
			"totalAmount", totalAmountDeci.StringFixed(0), "receives", utils.StrReceives(flow.Receives), "needSend", needSend)
		if needSend {
			for {
				poolBalance, err := w.conn.BalanceAt(poolAddr, nil)
				if err != nil {
					return errors.Wrap(err, "BalanceAt")
				}
				poolBalanceDeci := decimal.NewFromBigInt(poolBalance, 0)

				w.log.Debug("needSendProposal", "totalAmount", totalAmountDeci.StringFixed(0), "poolBalance", poolBalanceDeci.StringFixed(0))
				if poolBalanceDeci.LessThanOrEqual(totalAmountDeci) {
					// check again
					needSend, err = needSendProposal(client, multisigOnchainContract, proposalId)
					if err != nil {
						return errors.Wrap(err, "needSendProposal")
					}
					if needSend {
						time.Sleep(WaitInterval)
						w.log.Warn("pool balance not enough will wait",
							"pool", poolAddr.String(), "balance", poolBalanceDeci.String(), "totalTransferAmount", totalAmountDeci.StringFixed(0))
						continue
					}
					break
				}
				break
			}

			err = w.submitProposal(client, multisigOnchainContract, proposalId, proposalBts)
			if err != nil {
				return errors.Wrap(err, "submitProposal")
			}
		}
	}
	err = w.waitProposalExecuted(multisigOnchainContract, proposalId)
	if err != nil {
		return errors.Wrap(err, "waitProposalExecuted")
	}

	mef.RunTimeCalls = []*submodel.RunTimeCall{{CallHash: hexutil.Encode(proposalId[:])}}
	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.InformChain, Content: mef}
	return w.submitMessage(result)
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

func (w *writer) bscVoteRate(client *ethereum.Client, proposalId [32]byte, evmRate *big.Int) error {
	proposal, err := w.conn.bscStakePortalRateContract.Proposals(&bind.CallOpts{}, proposalId)
	if err != nil {
		return fmt.Errorf("processSignatureEnough Proposals error %s ", err)
	}
	if proposal.Status == 2 { // success status
		return nil
	}
	hasVoted, err := w.conn.bscStakePortalRateContract.HasVoted(&bind.CallOpts{}, proposalId, client.Opts().From)
	if err != nil {
		return fmt.Errorf("processSignatureEnough HasVoted error %s", err)
	}
	if hasVoted {
		return nil
	}

	// send tx
	err = client.LockAndUpdateOpts(big.NewInt(0), big.NewInt(0))
	if err != nil {
		return fmt.Errorf("processSignatureEnough LockAndUpdateOpts error %s", err)
	}

	voteTx, err := w.conn.bscStakePortalRateContract.VoteRate(client.Opts(), proposalId, evmRate)
	if err != nil {
		client.UnlockOpts()
		return fmt.Errorf("processSignatureEnough VoteRate error %s", err)
	}
	client.UnlockOpts()

	err = w.waitBscTxOk(voteTx.Hash())
	if err != nil {
		return fmt.Errorf("processSignatureEnough waitTxOk error %s", err)
	}

	return nil
}

func (task *writer) waitBscTxOk(txHash common.Hash) error {
	retry := 0
	for {
		if retry > BlockRetryLimit*3 {
			return fmt.Errorf("waitPolygonTxOk tx reach retry limit")
		}
		_, pending, err := task.conn.queryClient.TransactionByHash(context.Background(), txHash)
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

func (task *writer) waitBscRateUpdated(proposalId [32]byte) error {
	retry := 0
	for {
		if retry > BlockRetryLimit*3 {
			return fmt.Errorf("waitPolygonRateUpdated tx reach retry limit")
		}

		proposal, err := task.conn.bscStakePortalRateContract.Proposals(&bind.CallOpts{}, proposalId)
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
