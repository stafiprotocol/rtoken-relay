package cosmos

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
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

func GetBondUnbondUnsignedTx(client *rpc.Client, bond, unbond substrateTypes.U128,
	poolAddr types.AccAddress, valAddr types.ValAddress) (unSignedTx []byte, err error) {
	if bond.Int.Cmp(unbond.Int) > 0 {
		val := bond.Int.Sub(bond.Int, unbond.Int)
		unSignedTx, err = client.GenMultiSigRawDelegateTx(
			poolAddr,
			valAddr,
			types.NewCoin(client.GetDenom(), types.NewIntFromBigInt(val)))
	} else {
		val := unbond.Int.Sub(unbond.Int, bond.Int)
		unSignedTx, err = client.GenMultiSigRawUnDelegateTx(
			poolAddr,
			valAddr,
			types.NewCoin(client.GetDenom(), types.NewIntFromBigInt(val)))
	}

	return
}

func GetClaimRewardUnsignedTx(client *rpc.Client, poolAddr types.AccAddress, height int64) ([]byte, error) {
	//get reward
	rewardRes, err := client.QueryDelegationTotalRewards(poolAddr, height)
	if err != nil {
		return nil, err
	}
	rewardAmount := rewardRes.GetTotal().AmountOf(client.GetDenom()).TruncateInt()

	//get balanceAmount
	balanceAmount, err := client.QueryBalance(poolAddr, client.GetDenom(), height)
	if err != nil {
		return nil, err
	}

	//check balanceAmount and rewardAmount
	//(1)if balanceAmount>rewardAmount gen withdraw and delegate tx
	//(2)if balanceAmount<rewardAmount gen withdraw tx
	var unSignedTx []byte
	if balanceAmount.Balance.Amount.GT(rewardAmount) {
		unSignedTx, err = client.GenMultiSigRawWithdrawAllRewardThenDeleTx(
			poolAddr,
			height)
	} else {
		unSignedTx, err = client.GenMultiSigRawWithdrawAllRewardTx(
			poolAddr,
			height)
	}
	if err != nil {
		return nil, err
	}
	return unSignedTx, nil
}

func GetTransferUnsignedTx(client *rpc.Client, poolAddr types.AccAddress, receives []*submodel.Receive) ([]byte, error) {
	outPuts := make([]xBankTypes.Output, 0)
	for _, receive := range receives {
		hexAccountStr := hex.EncodeToString(receive.Recipient.AsAccountID[:])
		addr, err := types.AccAddressFromHex(hexAccountStr)
		if err != nil {
			//skp is err
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

func (w *writer) RebuildUnsignedTxFromSigs(client *rpc.Client, sigs *submodel.SubmitSignatures) (*cosmos.WrapUnsignedTx, error) {
	poolAddrHexStr := hex.EncodeToString(sigs.Pool)
	proposalIdHexStr := hex.EncodeToString(sigs.ProposalId)
	poolAddr, err := types.AccAddressFromHex(poolAddrHexStr)
	if err != nil {
		return nil, err
	}

	var wrappedUnSignedTx *cosmos.WrapUnsignedTx

	switch sigs.TxType {
	case submodel.OriginalBond:
		shotId, bond, unbond, _, err := ParseBondUnBondProposalId(sigs.ProposalId)
		if err != nil {
			return nil, err
		}
		//todo cosmos validator just for test,will got from stafi or cosmos
		var addrValidatorTestnetAteam, _ = types.ValAddressFromBech32("cosmosvaloper105gvcjgs6s4j5ws9srckx0drt4x8cwgywplh7p")
		unsignedTx, err := GetBondUnbondUnsignedTx(client, bond, unbond, poolAddr, addrValidatorTestnetAteam)
		if err != nil {
			return nil, err
		}
		wrappedUnSignedTx = &cosmos.WrapUnsignedTx{
			UnsignedTx: unsignedTx,
			SnapshotId: shotId,
			Era:        uint32(sigs.Era),
			Type:       submodel.OriginalBond,
			Key:        proposalIdHexStr,
		}

	case submodel.OriginalClaimRewards:
		shotId, height, err := ParseClaimRewardProposalId(sigs.ProposalId)
		if err != nil {
			return nil, err
		}
		unsignedTx, err := GetClaimRewardUnsignedTx(client, poolAddr, int64(height))
		if err != nil {
			return nil, err
		}
		wrappedUnSignedTx = &cosmos.WrapUnsignedTx{
			UnsignedTx: unsignedTx,
			SnapshotId: shotId,
			Era:        uint32(sigs.Era),
			Type:       submodel.OriginalClaimRewards,
			Key:        proposalIdHexStr,
		}

	case submodel.OriginalTransfer:
		showtId, err := ParseTransferProposalId(sigs.ProposalId)
		if err != nil {
			return nil, err
		}
		//get receivers from stafi
		msg := &core.Message{Source: sigs.Symbol, Destination: core.RFIS,
			Reason: core.GetReceivers, Content: &submodel.GetReceiversParams{
				Symbol: sigs.Symbol, Era: sigs.Era, Pool: sigs.Pool,
			}}
		ok := w.submitMessage(msg)
		if !ok {
			return nil, fmt.Errorf("get receiver from stafi failed ")
		}

		receiveList, ok := msg.Content.([]*submodel.Receive)
		if !ok {
			return nil, fmt.Errorf("cat msg to receive failed ")
		}
		unsignedTx, err := GetTransferUnsignedTx(client, poolAddr, receiveList)
		if err != nil {
			return nil, err
		}
		wrappedUnSignedTx = &cosmos.WrapUnsignedTx{
			UnsignedTx: unsignedTx,
			SnapshotId: showtId,
			Era:        uint32(sigs.Era),
			Type:       submodel.OriginalTransfer,
			Key:        proposalIdHexStr,
		}

	default:
		return nil, fmt.Errorf("rebuild from proposalId failed unknown tx type:%s", sigs.TxType)
	}

	return wrappedUnSignedTx, nil
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

func (w *writer) activeReport(source, dest core.RSymbol, flow *submodel.BondReportFlow) bool {
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
		return false
	}

	for {
		if retry <= 0 {
			w.log.Error("processSignatureEnoughEvt broadcast tx reach retry limit",
				"pool hex address", poolAddrHexStr)
			break
		}
		//check on chain
		res, err := client.QueryTxByHash(txHashHexStr)
		if err != nil || res.Empty() || res.Code != 0 {
			w.log.Warn(fmt.Sprintf("processSignatureEnoughEvt QueryTxByHash failed. will rebroadcast after %f second",
				BlockRetryInterval.Seconds()),
				"err or res.empty", err)
			retry--
		} else {
			w.log.Info("processSignatureEnoughEvt success",
				"pool hex address", poolAddrHexStr,
				"txHash", txHashHexStr)
			//return true only check on chain

			switch wrappedUnSignedTx.Type {
			case submodel.OriginalBond:
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
				delegationsRes, err := client.QueryDelegations(poolAddr, height)
				if err != nil {
					w.log.Error("processSignatureEnoughEvt QueryDelegations failed",
						"pool hex address", poolAddrHexStr,
						"err", err)
					return false
				}
				total := types.NewInt(0)
				for _, dele := range delegationsRes.GetDelegationResponses() {
					total = total.Add(dele.Balance.Amount)
				}

				rewardRes, err := client.QueryDelegationTotalRewards(poolAddr, height)
				if err != nil {
					w.log.Error("processSignatureEnoughEvt QueryDelegationTotalRewards failed",
						"pool hex address", poolAddrHexStr,
						"err", err)
					return false
				}
				rewardTotal := big.NewInt(0)
				rewardDe := rewardRes.Total.AmountOf(client.GetDenom())
				if !rewardDe.IsNil() {
					rewardTotal = rewardTotal.Add(rewardTotal, rewardDe.BigInt())
				}

				total.Add(types.NewIntFromBigInt(rewardTotal))
				f := submodel.BondReportFlow{
					Era:     wrappedUnSignedTx.Era,
					Rsymbol: sigs.Symbol,
					Pool:    sigs.Pool,
					ShotId:  wrappedUnSignedTx.SnapshotId,
					Active:  substrateTypes.NewU128(*total.BigInt())}

				poolClient.RemoveUnsignedTx(wrappedUnSignedTx.Key)
				return w.activeReport(m.Destination, m.Source, &f)
			case submodel.OriginalTransfer:
				callHash := utils.BlakeTwo256(sigs.Pool)
				mflow := submodel.MultiEventFlow{
					EventData: &submodel.TransferFlow{
						ShotId: wrappedUnSignedTx.SnapshotId},
					OpaqueCalls: []*submodel.MultiOpaqueCall{
						&submodel.MultiOpaqueCall{
							CallHash: hexutil.Encode(callHash[:])}},
				}

				poolClient.RemoveUnsignedTx(wrappedUnSignedTx.Key)
				return w.informChain(m.Destination, m.Source, &mflow)
			default:
				return true
			}
		}

		//broadcast if not on chain
		_, err = client.BroadcastTx(txBts)
		if err != nil && err != errType.ErrTxInMempoolCache {
			w.log.Warn("processSignatureEnoughEvt BroadcastTx failed",
				"err", err)
		}
		time.Sleep(BlockRetryInterval)
	}
	return false
}