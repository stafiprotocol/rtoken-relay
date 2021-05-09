package cosmos

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/types"
	errType "github.com/cosmos/cosmos-sdk/types/errors"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/scale"
	substrateTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"github.com/stafiprotocol/rtoken-relay/utils"
	"math/big"
	"time"
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

func GetClaimRewardProposalId(shotId substrateTypes.Hash, height uint64) []byte {
	proposalId := make([]byte, 40)
	copy(proposalId, shotId[:])
	binary.BigEndian.PutUint64(proposalId[32:], height)
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

func GetTransferProposalId(shotId substrateTypes.Hash) []byte {
	proposalId := make([]byte, 32)
	copy(proposalId, shotId[:])
	return proposalId
}

func ParseTransferProposalId(content []byte) (shotId substrateTypes.Hash, err error) {
	if len(content) != 32 {
		err = errors.New("cont length is not right")
		return
	}
	shotId = substrateTypes.NewHash(content)
	return
}

func GetValidatorUpdateProposalId(content []byte) []byte {
	hash := utils.BlakeTwo256(content)
	return hash[:]
}
func GetBondUnbondUnsignedTx(client *rpc.Client, bond, unbond substrateTypes.U128,
	poolAddr types.AccAddress, height int64) (unSignedTx []byte, err error) {
	if bond.Int.Cmp(unbond.Int) == 0 {
		return nil, errors.New("bond equal to unbond")
	}

	deleRes, err := client.QueryDelegations(poolAddr, height)
	if err != nil {
		return nil, err
	}

	totalDelegateAmount := types.NewInt(0)
	valAddrs := make([]types.ValAddress, 0)
	for _, dele := range deleRes.GetDelegationResponses() {
		//filter old validator
		if dele.GetBalance().Amount.LT(types.NewInt(2)) {
			continue
		}
		valAddr, err := types.ValAddressFromBech32(dele.GetDelegation().ValidatorAddress)
		if err != nil {
			return nil, err
		}
		valAddrs = append(valAddrs, valAddr)
		totalDelegateAmount = totalDelegateAmount.Add(dele.GetBalance().Amount)
	}

	valAddrsLen := len(valAddrs)
	if valAddrsLen == 0 {
		return nil, fmt.Errorf("no valAddrs,pool: %s", poolAddr)
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
		val := unbond.Int.Sub(unbond.Int, bond.Int)
		//make val <= totalDelegateAmount
		if val.Cmp(totalDelegateAmount.BigInt()) > 0 {
			val = totalDelegateAmount.BigInt()
		}

		val = val.Div(val, big.NewInt(int64(valAddrsLen)))
		unSignedTx, err = client.GenMultiSigRawUnDelegateTx(
			poolAddr,
			valAddrs,
			types.NewCoin(client.GetDenom(), types.NewIntFromBigInt(val)))
	}

	return
}

//if bond > unbond only gen delegate tx  (so delegate tx)

//if bond < unbond  (so withdraw and delegate tx)
//	(1) unbond validator gen delegate tx
//	(2) orther validator gen withdraw and delegate tx
//if bond == unbond  (so withdraw and delegate tx)
//	(1)if balanceAmount > rewardAmount of era height ,gen withdraw and delegate tx
//	(2)if balanceAmount < rewardAmount of era height, gen withdraw tx
func GetClaimRewardUnsignedTx(client *rpc.Client, poolAddr types.AccAddress, height int64,
	bond substrateTypes.U128, unBond substrateTypes.U128) ([]byte, error) {
	//get reward of height
	rewardRes, err := client.QueryDelegationTotalRewards(poolAddr, height)
	if err != nil {
		return nil, err
	}
	rewardAmount := rewardRes.GetTotal().AmountOf(client.GetDenom()).TruncateInt()

	bondCmpUnbond := bond.Cmp(unBond.Int)

	//check when we behind several eras,only bond==unbond we can check this
	// if rewardAmount > rewardAmountNow no need claim and delegate just return ErrNoMsgs
	if bondCmpUnbond == 0 {
		//get reward of now
		rewardResNow, err := client.QueryDelegationTotalRewards(poolAddr, 0)
		if err != nil {
			return nil, err
		}
		rewardAmountNow := rewardResNow.GetTotal().AmountOf(client.GetDenom()).TruncateInt()
		if rewardAmount.GT(rewardAmountNow) {
			return nil, rpc.ErrNoMsgs
		}
	}

	//get balanceAmount of height
	balanceAmount, err := client.QueryBalance(poolAddr, client.GetDenom(), height)
	if err != nil {
		return nil, err
	}

	var unSignedTx []byte
	if bondCmpUnbond <= 0 {
		//check balanceAmount and rewardAmount
		//(1)if balanceAmount>rewardAmount gen withdraw and delegate tx
		//(2)if balanceAmount<rewardAmount gen withdraw tx
		if balanceAmount.Balance.Amount.GT(rewardAmount) {
			unSignedTx, err = client.GenMultiSigRawWithdrawAllRewardThenDeleTx(
				poolAddr,
				height)
		} else {
			unSignedTx, err = client.GenMultiSigRawWithdrawAllRewardTx(
				poolAddr,
				height)
		}
	} else {
		//check balanceAmount and rewardAmount
		//(1)if balanceAmount>rewardAmount gen delegate tx
		//(2)if balanceAmount<rewardAmount gen withdraw tx
		if balanceAmount.Balance.Amount.GT(rewardAmount) {
			unSignedTx, err = client.GenMultiSigRawDeleRewardTx(
				poolAddr,
				height)
		} else {
			unSignedTx, err = client.GenMultiSigRawWithdrawAllRewardTx(
				poolAddr,
				height)
		}
	}

	if err != nil {
		return nil, err
	}
	return unSignedTx, nil
}

func GetTransferUnsignedTx(client *rpc.Client, poolAddr types.AccAddress, receives []*submodel.Receive, logger log15.Logger) ([]byte, error) {
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
	if len(outPuts) == 0 {
		return nil, ErrNoOutPuts
	}

	return client.GenMultiSigRawBatchTransferTx(poolAddr, outPuts)
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
				height := poolClient.GetHeightByEra(wrappedUnSignedTx.Era)
				poolClient.RemoveUnsignedTx(wrappedUnSignedTx.Key)
				return w.ActiveReport(client, poolAddr, height, sigs.Symbol, sigs.Pool,
					wrappedUnSignedTx.SnapshotId, wrappedUnSignedTx.Era, wrappedUnSignedTx.Bond, wrappedUnSignedTx.Unbond)
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
			height := poolClient.GetHeightByEra(wrappedUnSignedTx.Era)
			poolClient.RemoveUnsignedTx(wrappedUnSignedTx.Key)
			return w.ActiveReport(client, poolAddr, height, sigs.Symbol, sigs.Pool,
				wrappedUnSignedTx.SnapshotId, wrappedUnSignedTx.Era, wrappedUnSignedTx.Bond, wrappedUnSignedTx.Unbond)
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

//get total delegation of now and report
func (w *writer) ActiveReport(client *rpc.Client, poolAddr types.AccAddress, height int64,
	symbol core.RSymbol, poolBts []byte, shotId substrateTypes.Hash, era uint32,
	bond substrateTypes.U128, unBond substrateTypes.U128) bool {

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

	// rewardRes, err := client.QueryDelegationTotalRewards(poolAddr, height)
	// if err != nil {
	// 	w.log.Error("checkAndSend QueryDelegationTotalRewards failed",
	// 		"pool", poolAddr,
	// 		"err", err)
	// 	return false
	// }
	// rewardTotal := big.NewInt(0)
	// rewardDe := rewardRes.GetTotal().AmountOf(client.GetDenom()).TruncateInt()
	// if !rewardDe.IsNil() {
	// 	rewardTotal = rewardTotal.Add(rewardTotal, rewardDe.BigInt())
	// }
	// //reward
	// total = total.Add(types.NewIntFromBigInt(rewardTotal))
	// //bond/unbond
	// if bond.Cmp(unBond.Int) > 0 {
	// 	willBond := bond.Sub(bond.Int, unBond.Int)
	// 	total = total.Add(types.NewIntFromBigInt(willBond))
	// } else {
	// 	willUnBond := unBond.Sub(unBond.Int, bond.Int)
	// 	total = total.Sub(types.NewIntFromBigInt(willUnBond))
	// }

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
