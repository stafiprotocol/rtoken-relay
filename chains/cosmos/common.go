package cosmos

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/cosmos/cosmos-sdk/types"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stafiprotocol/go-substrate-rpc-client/scale"
	substrateTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"math/big"
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
