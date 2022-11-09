package cosmos

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/types"
	errType "github.com/cosmos/cosmos-sdk/types/errors"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	xStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/scale"
	substrateTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

var ErrNoOutPuts = errors.New("outputs length is zero")

func GetBondUnBondProposalId(shotId substrateTypes.Hash, bond, unbond substrateTypes.U128, seq uint64) []byte {
	proposalId := make([]byte, 72)
	copy(proposalId, shotId[:])

	var buffer = bytes.Buffer{}
	encoder := scale.NewEncoder(&buffer)
	encoder.Encode(bond)
	copy(proposalId[32:], buffer.Bytes())

	buffer.Reset()
	encoder.Encode(unbond)
	copy(proposalId[48:], buffer.Bytes())

	binary.BigEndian.PutUint64(proposalId[64:], seq)
	return proposalId
}

func ParseBondUnBondProposalId(content []byte) (shotId substrateTypes.Hash, bond, unbond substrateTypes.U128, seq uint64, err error) {
	if len(content) != 72 {
		err = errors.New("cont length is not right")
		return
	}
	shotId = substrateTypes.NewHash(content[:32])

	decoder := scale.NewDecoder(bytes.NewBuffer(content[32:48]))
	err = decoder.Decode(&bond)
	if err != nil {
		return
	}
	decoder2 := scale.NewDecoder(bytes.NewBuffer(content[48:64]))
	err = decoder2.Decode(&unbond)
	if err != nil {
		return
	}
	seq = binary.BigEndian.Uint64(content[64:])
	return
}

func GetClaimRewardProposalId(shotId substrateTypes.Hash, height uint64, factor uint32) []byte {
	proposalId := make([]byte, 40+4)
	copy(proposalId, shotId[:])
	binary.BigEndian.PutUint64(proposalId[32:], height)
	binary.BigEndian.PutUint32(proposalId[40:], factor)
	return proposalId
}

func ParseClaimRewardProposalId(content []byte) (shotId substrateTypes.Hash, height uint64, err error) {
	if len(content) != 40 {
		err = errors.New("cont length is not right")
		return
	}
	shotId = substrateTypes.NewHash(content[:32])
	height = binary.BigEndian.Uint64(content[32:])
	return
}

func GetTransferProposalId(txHash [32]byte, seq uint64) []byte {
	proposalId := make([]byte, 40)
	copy(proposalId, txHash[:])
	binary.BigEndian.PutUint64(proposalId[32:], seq)
	return proposalId
}

func GetValidatorUpdateProposalId(content []byte) []byte {
	hash := utils.BlakeTwo256(content)
	return hash[:]
}

// if bond == unbond return err
// if bond > unbond gen delegate tx
// if bond < unbond gen undelegate tx
func GetBondUnbondUnsignedTx(client *rpc.Client, bond, unbond substrateTypes.U128,
	poolAddr types.AccAddress, height int64) (unSignedTx []byte, err error) {
	//check bond unbond
	if bond.Int.Cmp(unbond.Int) == 0 {
		return nil, errors.New("bond equal to unbond")
	}

	deleRes, err := client.QueryDelegations(poolAddr, height)
	if err != nil {
		return nil, err
	}

	totalDelegateAmount := types.NewInt(0)
	valAddrs := make([]types.ValAddress, 0)
	deleAmount := make(map[string]types.Int)
	//get validators amount>=3
	for _, dele := range deleRes.GetDelegationResponses() {
		//filter old validator,we say validator is old if amount < 3 uatom
		if dele.GetBalance().Amount.LT(types.NewInt(3)) {
			continue
		}

		valAddr, err := types.ValAddressFromBech32(dele.GetDelegation().ValidatorAddress)
		if err != nil {
			return nil, err
		}
		valAddrs = append(valAddrs, valAddr)
		totalDelegateAmount = totalDelegateAmount.Add(dele.GetBalance().Amount)
		deleAmount[valAddr.String()] = dele.GetBalance().Amount
	}

	valAddrsLen := len(valAddrs)
	//check valAddrs length
	if valAddrsLen == 0 {
		return nil, fmt.Errorf("no valAddrs,pool: %s", poolAddr)
	}
	//check totalDelegateAmount
	if totalDelegateAmount.LT(types.NewInt(3 * int64(valAddrsLen))) {
		return nil, fmt.Errorf("validators have no reserve value to unbond")
	}

	//bond or unbond to their validators average
	if bond.Int.Cmp(unbond.Int) > 0 {
		val := bond.Int.Sub(bond.Int, unbond.Int)
		val = val.Div(val, big.NewInt(int64(valAddrsLen)))
		unSignedTx, err = client.GenMultiSigRawDelegateTx(
			poolAddr,
			valAddrs,
			types.NewCoin(client.GetDenom(), types.NewIntFromBigInt(val)))
	} else {

		//make val <= totalDelegateAmount-3*len and we revserve 3 uatom
		val := unbond.Int.Sub(unbond.Int, bond.Int)
		willUsetotalDelegateAmount := totalDelegateAmount.Sub(types.NewInt(3 * int64(valAddrsLen)))
		if val.Cmp(willUsetotalDelegateAmount.BigInt()) >= 0 {
			val = willUsetotalDelegateAmount.BigInt()
		}
		willUseTotalVal := types.NewIntFromBigInt(val)

		//sort validators by delegate amount
		sort.Slice(valAddrs, func(i int, j int) bool {
			return deleAmount[valAddrs[i].String()].
				GT(deleAmount[valAddrs[j].String()])
		})

		//choose validators to be undelegated
		choosedVals := make([]types.ValAddress, 0)
		choosedAmount := make(map[string]types.Int)

		selectedAmount := types.NewInt(0)
		for _, validator := range valAddrs {
			nowValMaxUnDeleAmount := deleAmount[validator.String()].Sub(types.NewInt(3))
			if selectedAmount.Add(nowValMaxUnDeleAmount).GTE(willUseTotalVal) {
				willUseChoosedAmount := willUseTotalVal.Sub(selectedAmount)

				choosedVals = append(choosedVals, validator)
				choosedAmount[validator.String()] = willUseChoosedAmount
				selectedAmount = selectedAmount.Add(willUseChoosedAmount)
				break
			}

			choosedVals = append(choosedVals, validator)
			choosedAmount[validator.String()] = nowValMaxUnDeleAmount
			selectedAmount = selectedAmount.Add(nowValMaxUnDeleAmount)
		}

		unSignedTx, err = client.GenMultiSigRawUnDelegateTxV2(
			poolAddr,
			choosedVals,
			choosedAmount)
	}

	return
}

// if bond == unbond return err
// if bond > unbond gen delegate tx
// if bond < unbond gen undelegate tx
func GetBondUnbondUnsignedTxWithTargets(client *rpc.Client, bond, unbond substrateTypes.U128,
	poolAddr types.AccAddress, height int64, targets []types.ValAddress) (unSignedTx []byte, err error) {
	//check bond unbond
	if bond.Int.Cmp(unbond.Int) == 0 {
		return nil, errors.New("bond equal to unbond")
	}

	deleRes, err := client.QueryDelegations(poolAddr, height)
	if err != nil {
		return nil, err
	}

	totalDelegateAmount := types.NewInt(0)
	valAddrs := make([]types.ValAddress, 0)
	deleAmount := make(map[string]types.Int)
	//get validators amount>=3
	for _, dele := range deleRes.GetDelegationResponses() {
		//filter old validator,we say validator is old if amount < 3 uatom
		if dele.GetBalance().Amount.LT(types.NewInt(3)) {
			continue
		}

		valAddr, err := types.ValAddressFromBech32(dele.GetDelegation().ValidatorAddress)
		if err != nil {
			return nil, err
		}

		valAddrs = append(valAddrs, valAddr)
		totalDelegateAmount = totalDelegateAmount.Add(dele.GetBalance().Amount)
		deleAmount[valAddr.String()] = dele.GetBalance().Amount
	}

	valAddrsLen := len(valAddrs)
	//check valAddrs length
	if valAddrsLen == 0 {
		return nil, fmt.Errorf("no valAddrs,pool: %s", poolAddr)
	}
	//check totalDelegateAmount
	if totalDelegateAmount.LT(types.NewInt(3 * int64(valAddrsLen))) {
		return nil, fmt.Errorf("validators have no reserve value to unbond")
	}

	//bond or unbond to their validators average
	if bond.Int.Cmp(unbond.Int) > 0 {
		valAddrs = targets
		valAddrsLen = len(valAddrs)
		//check valAddrs length
		if valAddrsLen == 0 {
			return nil, fmt.Errorf("no target valAddrs,pool: %s", poolAddr)
		}

		val := bond.Int.Sub(bond.Int, unbond.Int)
		val = val.Div(val, big.NewInt(int64(valAddrsLen)))
		unSignedTx, err = client.GenMultiSigRawDelegateTx(
			poolAddr,
			valAddrs,
			types.NewCoin(client.GetDenom(), types.NewIntFromBigInt(val)))
	} else {

		//make val <= totalDelegateAmount-3*len and we revserve 3 uatom
		val := unbond.Int.Sub(unbond.Int, bond.Int)
		willUsetotalDelegateAmount := totalDelegateAmount.Sub(types.NewInt(3 * int64(valAddrsLen)))
		if val.Cmp(willUsetotalDelegateAmount.BigInt()) > 0 {
			return nil, fmt.Errorf("no enough value can be used to unbond, pool: %s", poolAddr)
		}
		willUseTotalVal := types.NewIntFromBigInt(val)

		//remove if unbonding >= 7
		canUseValAddrs := make([]types.ValAddress, 0)
		for _, val := range valAddrs {
			res, err := client.QueryUnbondingDelegation(poolAddr, val, 0)
			if err != nil {
				// unbonding empty case
				if strings.Contains(err.Error(), "NotFound") {
					canUseValAddrs = append(canUseValAddrs, val)
					continue
				}
				return nil, err
			}
			if len(res.GetUnbond().Entries) < 7 {
				canUseValAddrs = append(canUseValAddrs, val)
			}
		}
		valAddrs = canUseValAddrs
		if len(valAddrs) == 0 {
			return nil, fmt.Errorf("no valAddrs can be used to unbond, pool: %s", poolAddr)
		}

		//sort validators by delegate amount
		sort.Slice(valAddrs, func(i int, j int) bool {
			return deleAmount[valAddrs[i].String()].
				GT(deleAmount[valAddrs[j].String()])
		})

		//choose validators to be undelegated
		choosedVals := make([]types.ValAddress, 0)
		choosedAmount := make(map[string]types.Int)

		selectedAmount := types.NewInt(0)
		enough := false
		for _, validator := range valAddrs {
			nowValMaxUnDeleAmount := deleAmount[validator.String()].Sub(types.NewInt(3))
			//if we find all validators needed
			if selectedAmount.Add(nowValMaxUnDeleAmount).GTE(willUseTotalVal) {
				willUseChoosedAmount := willUseTotalVal.Sub(selectedAmount)

				choosedVals = append(choosedVals, validator)
				choosedAmount[validator.String()] = willUseChoosedAmount
				selectedAmount = selectedAmount.Add(willUseChoosedAmount)

				enough = true
				break
			}

			choosedVals = append(choosedVals, validator)
			choosedAmount[validator.String()] = nowValMaxUnDeleAmount
			selectedAmount = selectedAmount.Add(nowValMaxUnDeleAmount)
		}

		if !enough {
			return nil, fmt.Errorf("can't find enough valAddrs to unbond, pool: %s", poolAddr)
		}

		unSignedTx, err = client.GenMultiSigRawUnDelegateTxV2(
			poolAddr,
			choosedVals,
			choosedAmount)
	}

	return
}

// if bond > unbond only gen delegate tx  (txType: 2)
// if bond <= unbond
//
//	(1)if balanceAmount > rewardAmount of era height ,gen withdraw and delegate tx  (txType: 3)
//	(2)if balanceAmount < rewardAmount of era height, gen withdraw tx  (txTpye: 1)
func GetClaimRewardUnsignedTx(client *rpc.Client, poolAddr types.AccAddress, height int64,
	bond substrateTypes.U128, unBond substrateTypes.U128) ([]byte, int, *types.Int, error) {
	//get reward of height
	rewardRes, err := client.QueryDelegationTotalRewards(poolAddr, height)
	if err != nil {
		return nil, 0, nil, err
	}
	rewardAmount := rewardRes.GetTotal().AmountOf(client.GetDenom()).TruncateInt()

	bondCmpUnbond := bond.Cmp(unBond.Int)

	//check when we behind several eras,only bond==unbond we can check this
	// if rewardAmount > rewardAmountNow no need claim and delegate just return ErrNoMsgs
	if bondCmpUnbond == 0 {
		//get reward of now
		rewardResNow, err := client.QueryDelegationTotalRewards(poolAddr, 0)
		if err != nil {
			return nil, 0, nil, err
		}
		rewardAmountNow := rewardResNow.GetTotal().AmountOf(client.GetDenom()).TruncateInt()
		if rewardAmount.GT(rewardAmountNow) {
			return nil, 0, nil, rpc.ErrNoMsgs
		}
	}

	//get balanceAmount of height
	balanceAmount, err := client.QueryBalance(poolAddr, client.GetDenom(), height)
	if err != nil {
		return nil, 0, nil, err
	}

	txType := 0
	var unSignedTx []byte
	if bondCmpUnbond <= 0 {
		//check balanceAmount and rewardAmount
		//(1)if balanceAmount>rewardAmount gen withdraw and delegate tx
		//(2)if balanceAmount<rewardAmount gen withdraw tx
		if balanceAmount.Balance.Amount.GT(rewardAmount) {
			unSignedTx, err = client.GenMultiSigRawWithdrawAllRewardThenDeleTx(
				poolAddr,
				height)
			txType = 3
		} else {
			unSignedTx, err = client.GenMultiSigRawWithdrawAllRewardTx(
				poolAddr,
				height)
			txType = 1
		}
	} else {
		//check balanceAmount and rewardAmount
		//(1)if balanceAmount>rewardAmount gen delegate tx
		//(2)if balanceAmount<rewardAmount gen withdraw tx
		if balanceAmount.Balance.Amount.GT(rewardAmount) {
			unSignedTx, err = client.GenMultiSigRawDeleRewardTx(
				poolAddr,
				height)
			txType = 2
		} else {
			unSignedTx, err = client.GenMultiSigRawWithdrawAllRewardTx(
				poolAddr,
				height)
			txType = 1
		}
	}

	if err != nil {
		return nil, 0, nil, err
	}

	decodedTx, err := client.GetTxConfig().TxJSONDecoder()(unSignedTx)
	if err != nil {
		return nil, 0, nil, err
	}
	totalAmountRet := types.NewInt(0)
	for _, msg := range decodedTx.GetMsgs() {
		if m, ok := msg.(*xStakingTypes.MsgDelegate); ok {
			totalAmountRet = totalAmountRet.Add(m.Amount.Amount)
		}

	}

	return unSignedTx, txType, &totalAmountRet, nil
}

func GetTransferUnsignedTx(client *rpc.Client, poolAddr types.AccAddress, receives []*submodel.Receive,
	logger log15.Logger) ([]byte, []xBankTypes.Output, error) {

	outPuts := make([]xBankTypes.Output, 0)
	for _, receive := range receives {
		hexAccountStr := hex.EncodeToString(receive.Recipient[:20])
		addr, err := types.AccAddressFromHex(hexAccountStr)
		if err != nil {
			logger.Error("GetTransferUnsignedTx AccAddressFromHex failed", "hexAccount", hexAccountStr, "err", err)
			continue
		}
		valueBigInt := big.Int(receive.Value)
		out := xBankTypes.Output{
			Address: addr.String(),
			Coins:   types.NewCoins(types.NewCoin(client.GetDenom(), types.NewIntFromBigInt(&valueBigInt))),
		}
		outPuts = append(outPuts, out)
	}

	//len should not be 0
	if len(outPuts) == 0 {
		return nil, nil, ErrNoOutPuts
	}

	//sort outPuts for the same rawTx from different relayer
	sort.SliceStable(outPuts, func(i, j int) bool {
		return bytes.Compare([]byte(outPuts[i].Address), []byte(outPuts[j].Address)) < 0
	})

	txBts, err := client.GenMultiSigRawBatchTransferTx(poolAddr, outPuts)
	if err != nil {
		return nil, nil, ErrNoOutPuts
	}
	return txBts, outPuts, nil
}

func (w *writer) printContentError(m *core.Message, err error) {
	w.log.Error("msg resolve failed", "source", m.Source, "dest", m.Destination, "reason", m.Reason, "err", err)
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (w *writer) submitMessage(m *core.Message) bool {
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

func (w *writer) activeReport(source, dest core.RSymbol, flow *submodel.BondReportedFlow) bool {
	msg := &core.Message{Source: source, Destination: dest, Reason: core.ActiveReport, Content: flow}
	return w.submitMessage(msg)
}

func (w *writer) processBondedPools(m *core.Message) bool {
	return true
}

func (w *writer) checkAndSend(poolClient *cosmos.PoolClient, wrappedUnSignedTx *cosmos.WrapUnsignedTx,
	sigs *submodel.SubmitSignatures, m *core.Message, txHash, txBts []byte) bool {
	retry := BlockRetryLimit
	txHashHexStr := hex.EncodeToString(txHash)
	client := poolClient.GetRpcClient()
	poolAddrHexStr := hex.EncodeToString(sigs.Pool)
	poolAddr, err := types.AccAddressFromHex(poolAddrHexStr)
	if err != nil {
		w.log.Error("checkAndSend AccAddressFromHex failed", "err", err)
		return false
	}

	for {
		if retry <= 0 {
			w.log.Error("checkAndSend broadcast tx reach retry limit",
				"pool hex address", poolAddrHexStr)

			if wrappedUnSignedTx.Type == submodel.OriginalClaimRewards {
				w.log.Info("claimRewards failed we still active report")
				poolClient.RemoveUnsignedTx(wrappedUnSignedTx.Key)
				return w.ActiveReport(client, poolAddr, sigs.Symbol, sigs.Pool,
					wrappedUnSignedTx.SnapshotId, wrappedUnSignedTx.Era)
			}
			break
		}
		//check on chain
		res, err := client.QueryTxByHash(txHashHexStr)
		if err != nil || res.Empty() || res.Code != 0 {
			w.log.Warn(fmt.Sprintf(
				"checkAndSend QueryTxByHash failed. will rebroadcast after %f second",
				BlockRetryInterval.Seconds()),
				"tx hash", txHashHexStr,
				"err or res.empty", err)

			//broadcast if not on chain
			_, err = client.BroadcastTx(txBts)
			if err != nil && err != errType.ErrTxInMempoolCache {
				w.log.Warn("checkAndSend BroadcastTx failed  will retry",
					"err", err)
			}
			time.Sleep(BlockRetryInterval)
			retry--
			continue
		}

		w.log.Info("checkAndSend success",
			"pool hex address", poolAddrHexStr,
			"tx type", wrappedUnSignedTx.Type,
			"era", wrappedUnSignedTx.Era,
			"txHash", txHashHexStr)

		//inform stafi
		switch wrappedUnSignedTx.Type {
		case submodel.OriginalBond: //bond or unbond
			callHash := utils.BlakeTwo256(sigs.Pool)
			mflow := submodel.MultiEventFlow{
				EventData: &submodel.EraPoolUpdatedFlow{
					ShotId: wrappedUnSignedTx.SnapshotId},
				OpaqueCalls: []*submodel.MultiOpaqueCall{
					&submodel.MultiOpaqueCall{
						CallHash: hexutil.Encode(callHash[:])}},
			}

			poolClient.RemoveUnsignedTx(wrappedUnSignedTx.Key)
			return w.informChain(m.Destination, m.Source, &mflow)
		case submodel.OriginalClaimRewards:
			poolClient.RemoveUnsignedTx(wrappedUnSignedTx.Key)
			return w.ActiveReport(client, poolAddr, sigs.Symbol, sigs.Pool,
				wrappedUnSignedTx.SnapshotId, wrappedUnSignedTx.Era)
		case submodel.OriginalTransfer:
			callHash := utils.BlakeTwo256(sigs.Pool)
			mflow := submodel.MultiEventFlow{
				EventData: &submodel.WithdrawReportedFlow{
					ShotId: wrappedUnSignedTx.SnapshotId},
				OpaqueCalls: []*submodel.MultiOpaqueCall{
					&submodel.MultiOpaqueCall{
						CallHash: hexutil.Encode(callHash[:])}},
			}

			poolClient.RemoveUnsignedTx(wrappedUnSignedTx.Key)
			return w.informChain(m.Destination, m.Source, &mflow)
		case submodel.OriginalWithdrawUnbond: //update validator
			poolClient.RemoveUnsignedTx(wrappedUnSignedTx.Key)
			return true
		default:
			w.log.Error("checkAndSend failed,unknown unsigned tx type",
				"pool", hex.EncodeToString(sigs.Pool),
				"type", wrappedUnSignedTx.Type)
			return false
		}

	}
	return false
}

// get total delegation of now and report
func (w *writer) ActiveReport(client *rpc.Client, poolAddr types.AccAddress,
	symbol core.RSymbol, poolBts []byte, shotId substrateTypes.Hash, era uint32) bool {

	delegationsRes, err := client.QueryDelegations(poolAddr, 0)
	if err != nil {
		w.log.Error("activeReport failed",
			"pool", poolAddr,
			"err", err)
		return false
	}
	total := types.NewInt(0)
	for _, dele := range delegationsRes.GetDelegationResponses() {
		total = total.Add(dele.Balance.Amount)
	}

	f := submodel.BondReportedFlow{
		Symbol: symbol,
		ShotId: shotId,
		Snap: &submodel.PoolSnapshot{
			Era:    era,
			Symbol: symbol,
			Pool:   poolBts,
			Active: substrateTypes.NewU128(*total.BigInt())},
	}

	w.log.Info("active report", "pool", poolAddr,
		"era", era, "active", total.String(), "symbol", symbol)
	return w.activeReport(symbol, core.RFIS, &f)
}

// ensue every validator claim reward
// if bond == unbond: if no delegation before, return errNoMsgs, else gen withdraw tx
// if bond > unbond: gen delegate tx
// if bond < unbond: gen undelegate+withdraw tx
func GetBondUnbondWithdrawUnsignedTxWithTargets(client *rpc.Client, bond, unbond *big.Int,
	poolAddr types.AccAddress, height int64, targets []types.ValAddress, memo string) (unSignedTx []byte, unSignedType int, err error) {

	poolAddrStr := poolAddr.String()

	switch bond.Cmp(unbond) {
	case 0:
		// return errnoMsgs if no delegation before
		var delegationRes *xStakingTypes.QueryDelegatorDelegationsResponse
		delegationRes, err = client.QueryDelegations(poolAddr, height)
		if err != nil {
			if strings.Contains(err.Error(), "unable to find delegations for address") {
				err = rpc.ErrNoMsgs
				return
			} else {
				return
			}
		}
		if len(delegationRes.DelegationResponses) == 0 {
			err = rpc.ErrNoMsgs
			return
		}

		unSignedTx, err = client.GenMultiSigRawWithdrawAllRewardTxWithMemo(
			poolAddr,
			height,
			memo)
		unSignedType = 0
		return
	case 1:
		valAddrs := targets
		valAddrsLen := len(valAddrs)
		//check valAddrs length
		if valAddrsLen == 0 {
			return nil, 0, fmt.Errorf("no target valAddrs, pool: %s", poolAddrStr)
		}

		//collect withdraw validators
		targetMap := make(map[string]bool)
		for _, target := range targets {
			targetMap[target.String()] = true
		}

		var deleRes *xStakingTypes.QueryDelegatorDelegationsResponse
		deleRes, err = client.QueryDelegations(poolAddr, height)
		if err != nil {
			return nil, 0, fmt.Errorf("QueryDelegations failed: %s", err)
		}

		withdrawValAddrs := make([]types.ValAddress, 0)
		//get validators amount>=3
		for _, dele := range deleRes.GetDelegationResponses() {
			//filter old validator,we say validator is old if amount < 3 uatom
			if dele.GetBalance().Amount.LT(types.NewInt(3)) {
				continue
			}

			valAddr, err := types.ValAddressFromBech32(dele.GetDelegation().ValidatorAddress)
			if err != nil {
				return nil, 0, err
			}
			if !targetMap[valAddr.String()] {
				withdrawValAddrs = append(withdrawValAddrs, valAddr)
			}
		}

		val := bond.Sub(bond, unbond)
		val = val.Div(val, big.NewInt(int64(valAddrsLen)))
		unSignedTx, err = client.GenMultiSigRawDelegateWithdrawTxWithMemo(
			poolAddr,
			valAddrs,
			types.NewCoin(client.GetDenom(), types.NewIntFromBigInt(val)),
			withdrawValAddrs,
			memo)
		unSignedType = 1
		return
	case -1:
		var deleRes *xStakingTypes.QueryDelegatorDelegationsResponse
		deleRes, err = client.QueryDelegations(poolAddr, height)
		if err != nil {
			return nil, 0, fmt.Errorf("QueryDelegations failed: %s", err)
		}

		totalDelegateAmount := types.NewInt(0)
		valAddrs := make([]types.ValAddress, 0)
		deleAmount := make(map[string]types.Int)

		//get validators amount>=3
		for _, dele := range deleRes.GetDelegationResponses() {
			//filter old validator,we say validator is old if amount < 3 uatom
			if dele.GetBalance().Amount.LT(types.NewInt(3)) {
				continue
			}

			valAddr, err := types.ValAddressFromBech32(dele.GetDelegation().ValidatorAddress)
			if err != nil {
				return nil, 0, err
			}

			valAddrs = append(valAddrs, valAddr)
			totalDelegateAmount = totalDelegateAmount.Add(dele.GetBalance().Amount)
			deleAmount[valAddr.String()] = dele.GetBalance().Amount
		}

		//check valAddrs length
		valAddrsLen := len(valAddrs)
		if valAddrsLen == 0 {
			return nil, 0, fmt.Errorf("no valAddrs, pool: %s", poolAddr)
		}

		//check totalDelegateAmount
		if totalDelegateAmount.LT(types.NewInt(3 * int64(valAddrsLen))) {
			return nil, 0, fmt.Errorf("validators have no reserve value to unbond")
		}

		//make val <= totalDelegateAmount-3*len and we reserve 3 uatom
		val := new(big.Int).Sub(unbond, bond)
		willUsetotalDelegateAmount := totalDelegateAmount.Sub(types.NewInt(3 * int64(valAddrsLen)))
		if val.Cmp(willUsetotalDelegateAmount.BigInt()) > 0 {
			return nil, 0, fmt.Errorf("no enough value can be used to unbond, pool: %s", poolAddr)
		}
		willUseTotalVal := types.NewIntFromBigInt(val)

		//remove validator who's unbonding >= 7
		canUseValAddrs := make([]types.ValAddress, 0)
		for _, val := range valAddrs {
			res, err := client.QueryUnbondingDelegation(poolAddr, val, height)
			if err != nil {
				// unbonding empty case
				if strings.Contains(err.Error(), "NotFound") {
					canUseValAddrs = append(canUseValAddrs, val)
					continue
				}
				return nil, 0, err
			}
			if len(res.GetUnbond().Entries) < 7 {
				canUseValAddrs = append(canUseValAddrs, val)
			}
		}
		valAddrs = canUseValAddrs
		if len(valAddrs) == 0 {
			return nil, 0, fmt.Errorf("no valAddrs can be used to unbond, pool: %s", poolAddr)
		}

		//sort validators by delegate amount
		sort.Slice(valAddrs, func(i int, j int) bool {
			return deleAmount[valAddrs[i].String()].
				GT(deleAmount[valAddrs[j].String()])
		})

		//choose validators to be undelegated
		choosedVals := make([]types.ValAddress, 0)
		choosedAmount := make(map[string]types.Int)

		selectedAmount := types.NewInt(0)
		enough := false
		for _, validator := range valAddrs {
			nowValMaxUnDeleAmount := deleAmount[validator.String()].Sub(types.NewInt(3))
			//if we find all validators needed
			if selectedAmount.Add(nowValMaxUnDeleAmount).GTE(willUseTotalVal) {
				willUseChoosedAmount := willUseTotalVal.Sub(selectedAmount)

				choosedVals = append(choosedVals, validator)
				choosedAmount[validator.String()] = willUseChoosedAmount
				selectedAmount = selectedAmount.Add(willUseChoosedAmount)

				enough = true
				break
			}

			choosedVals = append(choosedVals, validator)
			choosedAmount[validator.String()] = nowValMaxUnDeleAmount
			selectedAmount = selectedAmount.Add(nowValMaxUnDeleAmount)
		}

		if !enough {
			return nil, 0, fmt.Errorf("can't find enough valAddrs to unbond, pool: %s", poolAddrStr)
		}

		// filter withdraw validators
		withdrawVals := make([]types.ValAddress, 0)
		for valAddressStr := range deleAmount {
			if _, exist := choosedAmount[valAddressStr]; !exist {
				valAddr, err := types.ValAddressFromBech32(valAddressStr)
				if err != nil {
					return nil, 0, err
				}
				withdrawVals = append(withdrawVals, valAddr)
			}
		}

		unSignedTx, err = client.GenMultiSigRawUnDelegateWithdrawTxWithMemo(
			poolAddr,
			choosedVals,
			choosedAmount,
			withdrawVals,
			memo)
		unSignedType = -1
		return
	default:
		return nil, 0, fmt.Errorf("unreached case err")
	}
}

// Notice: delegate/undelegate/withdraw operates will withdraw all reward
// all delegations had withdraw all reward in eraUpdatedEvent handler
// (0)if rewardAmount of  height == 0, hubClient.ErrNoMsgs
// (1)else gen delegate tx
func GetDelegateRewardUnsignedTxWithReward(client *rpc.Client, poolAddr types.AccAddress, height int64,
	rewards map[string]types.Coin, memo string) ([]byte, *types.Int, error) {

	unSignedTx, err := client.GenMultiSigRawDeleRewardTxWithRewardsWithMemo(
		poolAddr,
		height,
		rewards,
		memo)
	if err != nil {
		return nil, nil, err
	}

	decodedTx, err := client.GetTxConfig().TxJSONDecoder()(unSignedTx)
	if err != nil {
		return nil, nil, fmt.Errorf("GetTxConfig().TxDecoder() failed: %s, unSignedTx: %s", err, string(unSignedTx))
	}
	totalAmountRet := types.NewInt(0)
	for _, msg := range decodedTx.GetMsgs() {
		if m, ok := msg.(*xStakingTypes.MsgDelegate); ok {
			totalAmountRet = totalAmountRet.Add(m.Amount.Amount)
		}
	}
	return unSignedTx, &totalAmountRet, nil
}

func GetTransferUnsignedTxWithMemo(client *rpc.Client, poolAddr types.AccAddress, receives []*submodel.Receive, memo string,
	logger log15.Logger) ([]byte, []xBankTypes.Output, error) {

	outPuts := make([]xBankTypes.Output, 0)
	for _, receive := range receives {
		hexAccountStr := hex.EncodeToString(receive.Recipient[:20])
		addr, err := types.AccAddressFromHex(hexAccountStr)
		if err != nil {
			logger.Error("GetTransferUnsignedTx AccAddressFromHex failed", "hexAccount", hexAccountStr, "err", err)
			continue
		}
		valueBigInt := big.Int(receive.Value)
		out := xBankTypes.Output{
			Address: addr.String(),
			Coins:   types.NewCoins(types.NewCoin(client.GetDenom(), types.NewIntFromBigInt(&valueBigInt))),
		}
		outPuts = append(outPuts, out)
	}

	//len should not be 0
	if len(outPuts) == 0 {
		return nil, nil, ErrNoOutPuts
	}

	//sort outPuts for the same rawTx from different relayer
	sort.SliceStable(outPuts, func(i, j int) bool {
		return bytes.Compare([]byte(outPuts[i].Address), []byte(outPuts[j].Address)) < 0
	})

	txBts, err := client.GenMultiSigRawBatchTransferTxWithMemo(poolAddr, outPuts, memo)
	if err != nil {
		return nil, nil, ErrNoOutPuts
	}
	return txBts, outPuts, nil
}
