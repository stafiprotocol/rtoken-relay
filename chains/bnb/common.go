package bnb

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	staking "github.com/stafiprotocol/rtoken-relay/bindings/Staking"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
)

var (
	StakingAbi                      abi.ABI
	ErrNoAvailableValsForUnDelegate = errors.New("ErrNoAvailableValsForUnDelegate")
)

func init() {

	stakingAbi, err := abi.JSON(strings.NewReader(staking.StakingABI))
	if err != nil {
		panic(err)
	}
	StakingAbi = stakingAbi

}

func getProposalId(era uint32, event, step string, factor int) common.Hash {
	return crypto.Keccak256Hash([]byte(fmt.Sprintf("era-%d-%s-%s-%d", era, event, step, factor)))
}

func (w *writer) getClaimRewardProposal() ([]byte, error) {
	inputData, err := StakingAbi.Pack("claimReward")
	if err != nil {
		return nil, errors.Wrap(err, "staking abi pack failed")
	}

	bt := &ethmodel.BatchTransaction{
		Operation:  uint8(ethmodel.Call),
		To:         w.conn.GetStakingContractAddress(),
		Value:      big.NewInt(0),
		DataLength: big.NewInt(int64(len(inputData))),
		Data:       inputData,
	}

	return bt.Encode(), nil
}

func (w *writer) getClaimUndelegatedProposal() ([]byte, error) {
	inputData, err := StakingAbi.Pack("claimUndelegated")
	if err != nil {
		return nil, errors.Wrap(err, "staking abi pack failed")
	}

	bt := &ethmodel.BatchTransaction{
		Operation:  uint8(ethmodel.Call),
		To:         w.conn.GetStakingContractAddress(),
		Value:      big.NewInt(0),
		DataLength: big.NewInt(int64(len(inputData))),
		Data:       inputData,
	}

	return bt.Encode(), nil
}

func (w *writer) getDelegateProposal(totalAmount, relayerFee, leastBond decimal.Decimal, poolAddr common.Address, validators []common.Address, targetBlock int64) ([]byte, []common.Address, error) {
	if len(validators) == 0 {
		return nil, nil, fmt.Errorf("validators empty")
	}
	if totalAmount.LessThan(leastBond) {
		return nil, nil, fmt.Errorf("totalAmount %s less than leastBond %s", totalAmount, leastBond)
	}

	delegator := poolAddr

	validatorDelegatedAmount := make(map[common.Address]decimal.Decimal)
	for _, v := range validators {
		delegatedAmount, err := w.conn.stakingContract.GetDelegated(&bind.CallOpts{
			From:        poolAddr,
			BlockNumber: big.NewInt(targetBlock),
			Context:     context.Background(),
		}, delegator, v)
		if err != nil {
			return nil, nil, errors.Wrap(err, "stakingContract.GetDelegated")
		}
		validatorDelegatedAmount[v] = decimal.NewFromBigInt(delegatedAmount, 0)
	}
	// sort by delegated amount asc
	sort.SliceStable(validators, func(i, j int) bool {
		return validatorDelegatedAmount[validators[i]].LessThan(validatorDelegatedAmount[validators[j]])
	})

	distributedValidatorAmount := make(map[common.Address]decimal.Decimal)
out:
	for {
		for _, v := range validators {
			if totalAmount.GreaterThanOrEqual(leastBond) {
				distributedValidatorAmount[v] = distributedValidatorAmount[v].Add(leastBond)
				totalAmount = totalAmount.Sub(leastBond)
			} else {
				distributedValidatorAmount[v] = distributedValidatorAmount[v].Add(totalAmount)
				break out
			}
		}
	}

	distributedValidators := make([]common.Address, 0)
	for val := range distributedValidatorAmount {
		distributedValidators = append(distributedValidators, val)
	}

	// sort distributed validators by distributed amount desc
	sort.SliceStable(distributedValidators, func(i, j int) bool {
		return distributedValidatorAmount[validators[i]].GreaterThan(distributedValidatorAmount[validators[j]])
	})

	txs := make(ethmodel.BatchTransactions, 0)
	for _, val := range distributedValidators {
		amount := distributedValidatorAmount[val]
		inputData, err := StakingAbi.Pack("delegate", val, amount.BigInt())
		if err != nil {
			return nil, nil, errors.Wrap(err, "staking abi pack failed")
		}

		tx := &ethmodel.BatchTransaction{
			Operation:  uint8(ethmodel.Call),
			To:         w.conn.GetStakingContractAddress(),
			Value:      amount.Add(relayerFee).BigInt(),
			DataLength: big.NewInt(int64(len(inputData))),
			Data:       inputData,
		}
		txs = append(txs, tx)
	}

	return txs.Encode(), distributedValidators, nil
}

func (w *writer) getUnDelegateProposal(totalAmount, relayerFee, leastBond decimal.Decimal, validators []common.Address, poolAddr common.Address, targetBlock int64) ([]byte, []common.Address, error) {
	if len(validators) == 0 {
		return nil, nil, fmt.Errorf("validators empty")
	}
	if totalAmount.LessThan(leastBond) {
		return nil, nil, fmt.Errorf("totalAmount %s less than leastBond %s", totalAmount, leastBond)
	}

	delegator := poolAddr

	block, err := w.conn.QueryBlock(targetBlock)
	if err != nil {
		return nil, nil, errors.Wrap(err, "QueryBlock")
	}

	validatorDelegated := make(map[common.Address]decimal.Decimal)
	for _, v := range validators {
		undelegateTime, err := w.conn.stakingContract.GetPendingUndelegateTime(&bind.CallOpts{
			From:        poolAddr,
			BlockNumber: big.NewInt(targetBlock),
			Context:     context.Background(),
		}, delegator, v)
		if err != nil {
			return nil, nil, errors.Wrap(err, "stakingContract.GetPendingUndelegateTime")
		}

		if block.Time() < undelegateTime.Uint64() {
			continue
		}

		delegatedAmount, err := w.conn.stakingContract.GetDelegated(&bind.CallOpts{
			From:        poolAddr,
			BlockNumber: big.NewInt(targetBlock),
			Context:     context.Background(),
		}, delegator, v)
		if err != nil {
			return nil, nil, errors.Wrap(err, "stakingContract.GetDelegated")
		}
		validatorDelegated[v] = decimal.NewFromBigInt(delegatedAmount, 0)
	}
	sort.SliceStable(validators, func(i, j int) bool {
		return validatorDelegated[validators[i]].GreaterThan(validatorDelegated[validators[j]])
	})

	selectedValidatorsAmount := make(map[common.Address]decimal.Decimal)
	selectedValidators := make([]common.Address, 0)
	selectedAmount := decimal.Zero

	for _, val := range validators {
		for {
			if validatorDelegated[val].GreaterThanOrEqual(leastBond) && selectedAmount.LessThan(totalAmount) {
				selectedValidatorsAmount[val] = selectedValidatorsAmount[val].Add(leastBond)
				validatorDelegated[val] = validatorDelegated[val].Sub(leastBond)
				selectedAmount = selectedAmount.Add(leastBond)
				continue
			}
			break
		}
	}
	if !selectedAmount.Equal(totalAmount) {
		return nil, nil, ErrNoAvailableValsForUnDelegate
	}

	for v := range selectedValidatorsAmount {
		selectedValidators = append(selectedValidators, v)
	}

	// sort by amount asc
	sort.Slice(selectedValidators, func(i, j int) bool {
		return selectedValidatorsAmount[selectedValidators[i]].GreaterThan(selectedValidatorsAmount[selectedValidators[j]])
	})

	txs := make(ethmodel.BatchTransactions, 0)
	for _, val := range selectedValidators {
		amount := selectedValidatorsAmount[val]
		inputData, err := StakingAbi.Pack("undelegate", val, amount.BigInt())
		if err != nil {
			return nil, nil, errors.Wrap(err, "staking abi pack failed")
		}

		tx := &ethmodel.BatchTransaction{
			Operation:  uint8(ethmodel.Call),
			To:         w.conn.GetStakingContractAddress(),
			Value:      relayerFee.BigInt(),
			DataLength: big.NewInt(int64(len(inputData))),
			Data:       inputData,
		}
		txs = append(txs, tx)
	}

	return txs.Encode(), selectedValidators, nil
}

func (w *writer) getTransferProposal(poolAddr common.Address, receives []*submodel.Receive) ([]byte, decimal.Decimal, error) {

	totalAmount := decimal.Zero
	txs := make(ethmodel.BatchTransactions, 0)
	for _, r := range receives {
		amount := decimal.NewFromBigInt((*big.Int)(&r.Value), 0).Mul(tenDecimals) // decimals: bc 8 bsc 18
		to := common.BytesToAddress(r.Recipient)
		tx := &ethmodel.BatchTransaction{
			Operation:  uint8(ethmodel.Call),
			To:         to,
			Value:      amount.BigInt(),
			DataLength: big.NewInt(0),
			Data:       []byte{},
		}
		txs = append(txs, tx)
		totalAmount = totalAmount.Add(amount)
	}

	return txs.Encode(), totalAmount, nil
}

func (w *writer) unbondable(totalAmount, relayerFee, leastBond decimal.Decimal, validators []common.Address, poolAddr common.Address, targetBlock int64) (bool, error) {
	_, _, err := w.getUnDelegateProposal(totalAmount, relayerFee, leastBond, validators, poolAddr, targetBlock)
	if err != nil {
		if err == ErrNoAvailableValsForUnDelegate {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (w *writer) staked(pool *Pool) (*big.Int, error) {
	delegator := pool.poolAddress
	return w.conn.stakingContract.GetTotalDelegated(&bind.CallOpts{
		From:    delegator,
		Context: context.Background(),
	}, delegator)
}

func (w *writer) findRealRewardAmountClaimed(pool *Pool, proposalId [32]byte, poolAddr common.Address, targetHeight uint64) (*big.Int, uint64, error) {
	proposalExectedIterator, err := pool.multisigOnchain.FilterProposalExecuted(&bind.FilterOpts{
		Start:   targetHeight,
		Context: context.Background(),
	}, [][32]byte{proposalId})
	if err != nil {
		return nil, 0, errors.Wrap(err, "multisigOnchain.FilterProposalExecuted")
	}
	for proposalExectedIterator.Next() {
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
	return nil, 0, fmt.Errorf("not find reward claim event")
}

func (w *writer) findRealUndelegatedAmountClaimed(pool *Pool, proposalId [32]byte, poolAddr common.Address, targetHeight uint64) (*big.Int, uint64, error) {
	proposalExectedIterator, err := pool.multisigOnchain.FilterProposalExecuted(&bind.FilterOpts{
		Start:   targetHeight,
		Context: context.Background(),
	}, [][32]byte{proposalId})
	if err != nil {
		return nil, 0, errors.Wrap(err, "multisigOnchain.FilterProposalExecuted")
	}
	for proposalExectedIterator.Next() {
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
	w.log.Info("submitProposal ok", "pool", pool.poolAddress, "proposalId", hexutil.Encode(proposalId[:]), "txHash", tx.Hash())
	return nil
}

func needSendProposal(pool *Pool, proposalId [32]byte) (bool, error) {
	proposal, err := pool.multisigOnchain.Proposals(&bind.CallOpts{
		Context: context.Background(),
	}, proposalId)
	if err != nil {
		return false, errors.Wrap(err, "multisigOnchain.Proposals")
	}
	needSend := false
	switch proposal.Status {
	case 0:
		needSend = true
	case 1:
		voted, err := pool.multisigOnchain.HasVoted(&bind.CallOpts{
			Context: context.Background(),
		}, proposalId, pool.bscClient.Opts().From)
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

func (w *writer) waitProposalExecuted(pool *Pool, proposalId [32]byte) error {
	retry := 0
	for {
		if retry > GetRetryLimit*6 {
			return fmt.Errorf("waitProposalExecuted reach retry limit")
		}

		proposal, err := pool.multisigOnchain.Proposals(&bind.CallOpts{}, proposalId)
		if err != nil {
			w.log.Warn("get proposal failed, will retry", "err", err.Error(), "proposalId", hexutil.Encode(proposalId[:]))
			time.Sleep(WaitInterval)
			retry++
			continue
		}
		if proposal.Status != 2 {
			w.log.Warn("proposals not executed, will wait", "proposalId", hexutil.Encode(proposalId[:]))
			time.Sleep(WaitInterval)
			retry++
			continue
		}
		return nil
	}
}

// bsc staking contract:
//
//	function getRequestInFly(address delegator) override external view returns(uint256[3] memory) {
//	    uint256[3] memory request;
//	    request[0] = delegateInFly[delegator];
//	    request[1] = undelegateInFly[delegator];
//	    request[2] = redelegateInFly[delegator];
//	    return request;
//	}
func (w *writer) waitDelegateCrossChainOk(pool *Pool, proposalId [32]byte, targetHeight uint64, validators []common.Address) error {

	delegator := pool.poolAddress
	retry := 0
	for {
		if retry > GetRetryLimit*6 {
			return fmt.Errorf("waitDelegateCrossChainOk reach retry limit")
		}
		inFlys, err := w.conn.stakingContract.GetRequestInFly(&bind.CallOpts{
			From:    delegator,
			Context: context.Background(),
		}, delegator)

		if err != nil {
			w.log.Warn("GetRequestInFly failed, will retry", "err", err.Error(), "proposalId", hexutil.Encode(proposalId[:]))
			time.Sleep(WaitInterval)
			retry++
			continue
		}
		if inFlys[0].Sign() != 0 {
			w.log.Warn("delegate is in fly, will retry", "proposalId", hexutil.Encode(proposalId[:]))
			time.Sleep(WaitInterval)
			retry++
			continue
		}

		delegateSucessIterator, err := w.conn.stakingContract.FilterDelegateSuccess(&bind.FilterOpts{
			Start:   targetHeight,
			Context: context.Background(),
		}, []common.Address{delegator}, validators)

		if err != nil {
			w.log.Warn("FilterDelegateSuccess failed, will retry", "err", err.Error(), "proposalId", hexutil.Encode(proposalId[:]))
			time.Sleep(WaitInterval)
			retry++
			continue
		}
		successCount := 0
		for delegateSucessIterator.Next() {
			successCount++
		}
		if successCount != len(validators) {
			return fmt.Errorf("some validators delegate failed, pool: %s", pool.poolAddress.String())
		}

		return nil
	}
}

func (w *writer) waitUnDelegateCrossChainOk(pool *Pool, proposalId [32]byte, targetHeight uint64, validators []common.Address) error {

	delegator := pool.poolAddress
	retry := 0
	for {
		if retry > GetRetryLimit*6 {
			return fmt.Errorf("waitDelegateCrossChainOk reach retry limit")
		}
		inFlys, err := w.conn.stakingContract.GetRequestInFly(&bind.CallOpts{
			From:    delegator,
			Context: context.Background(),
		}, delegator)

		if err != nil {
			w.log.Warn("GetRequestInFly failed, will retry", "err", err.Error(), "proposalId", hexutil.Encode(proposalId[:]))
			time.Sleep(WaitInterval)
			retry++
			continue
		}
		if inFlys[1].Sign() != 0 {
			w.log.Warn("undelegate is in fly, will retry", "proposalId", hexutil.Encode(proposalId[:]))
			time.Sleep(WaitInterval)
			retry++
			continue
		}

		undelegateSucessIterator, err := w.conn.stakingContract.FilterUndelegateSuccess(&bind.FilterOpts{
			Start:   targetHeight,
			Context: context.Background(),
		}, []common.Address{delegator}, validators)

		if err != nil {
			w.log.Warn("FilterDelegateSuccess failed, will retry", "err", err.Error(), "proposalId", hexutil.Encode(proposalId[:]))
			time.Sleep(WaitInterval)
			retry++
			continue
		}
		successCount := 0
		for undelegateSucessIterator.Next() {
			successCount++
		}
		if successCount != len(validators) {
			return fmt.Errorf("some validators undelegate failed, pool: %s", pool.poolAddress.String())
		}

		return nil
	}
}
