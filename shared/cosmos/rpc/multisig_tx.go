package rpc

import (
	"errors"
	"fmt"
	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	kMultiSig "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xAuthClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	xAuthSigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	xDistriTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	xStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
	tendermintTypes "github.com/tendermint/tendermint/types"
)

var ErrNoMsgs = errors.New("no tx msgs")

//c.clientCtx.FromAddress must be multi sig address
func (c *Client) GenMultiSigRawTransferTx(toAddr types.AccAddress, amount types.Coins) ([]byte, error) {
	msg := xBankTypes.NewMsgSend(c.clientCtx.GetFromAddress(), toAddr, amount)
	return c.GenMultiSigRawTx(msg)
}

//only support one type coin
func (c *Client) GenMultiSigRawBatchTransferTx(poolAddr types.AccAddress, outs []xBankTypes.Output) ([]byte, error) {
	totalAmount := types.NewInt(0)
	for _, out := range outs {
		for _, coin := range out.Coins {
			totalAmount = totalAmount.Add(coin.Amount)
		}
	}
	input := xBankTypes.Input{
		Address: poolAddr.String(),
		Coins:   types.NewCoins(types.NewCoin(c.denom, totalAmount))}

	msg := xBankTypes.NewMsgMultiSend([]xBankTypes.Input{input}, outs)
	return c.GenMultiSigRawTx(msg)
}

//generate unsigned delegate tx
func (c *Client) GenMultiSigRawDelegateTx(delAddr types.AccAddress, valAddrs []types.ValAddress, amount types.Coin) ([]byte, error) {
	if len(valAddrs) == 0 {
		return nil, errors.New("no valAddrs")
	}
	if amount.IsZero() {
		return nil, errors.New("amount is zero")
	}

	msgs := make([]types.Msg, 0)
	for _, valAddr := range valAddrs {
		msg := xStakingTypes.NewMsgDelegate(delAddr, valAddr, amount)
		msgs = append(msgs, msg)
	}

	return c.GenMultiSigRawTx(msgs...)
}

//generate unsigned unDelegate tx
func (c *Client) GenMultiSigRawUnDelegateTx(delAddr types.AccAddress, valAddrs []types.ValAddress, amount types.Coin) ([]byte, error) {
	if len(valAddrs) == 0 {
		return nil, errors.New("no valAddrs")
	}
	if amount.IsZero() {
		return nil, errors.New("amount is zero")
	}
	msgs := make([]types.Msg, 0)
	for _, valAddr := range valAddrs {
		msg := xStakingTypes.NewMsgUndelegate(delAddr, valAddr, amount)
		msgs = append(msgs, msg)
	}
	return c.GenMultiSigRawTx(msgs...)
}

//generate unsigned unDelegate tx
func (c *Client) GenMultiSigRawUnDelegateTxV2(delAddr types.AccAddress, valAddrs []types.ValAddress,
	amounts map[string]types.Int) ([]byte, error) {

	if len(valAddrs) == 0 {
		return nil, errors.New("no valAddrs")
	}
	msgs := make([]types.Msg, 0)
	for _, valAddr := range valAddrs {
		amount := types.NewCoin(c.GetDenom(), amounts[valAddr.String()])
		if amount.IsZero() {
			return nil, errors.New("amount is zero")
		}
		msg := xStakingTypes.NewMsgUndelegate(delAddr, valAddr, amount)
		msgs = append(msgs, msg)
	}
	return c.GenMultiSigRawTx(msgs...)
}

//generate unsigned reDelegate tx
func (c *Client) GenMultiSigRawReDelegateTx(delAddr types.AccAddress, valSrcAddr, valDstAddr types.ValAddress, amount types.Coin) ([]byte, error) {
	msg := xStakingTypes.NewMsgBeginRedelegate(delAddr, valSrcAddr, valDstAddr, amount)
	return c.GenMultiSigRawTx(msg)
}

//generate unsigned withdraw delegate reward tx
func (c *Client) GenMultiSigRawWithdrawDeleRewardTx(delAddr types.AccAddress, valAddr types.ValAddress) ([]byte, error) {
	msg := xDistriTypes.NewMsgWithdrawDelegatorReward(delAddr, valAddr)
	return c.GenMultiSigRawTx(msg)
}

//generate unsigned withdraw reward then delegate reward tx
func (c *Client) GenMultiSigRawWithdrawRewardThenDeleTx(delAddr types.AccAddress, valAddr types.ValAddress, amount types.Coin) ([]byte, error) {
	msg := xDistriTypes.NewMsgWithdrawDelegatorReward(delAddr, valAddr)
	msg2 := xStakingTypes.NewMsgDelegate(delAddr, valAddr, amount)
	return c.GenMultiSigRawTx(msg, msg2)
}

//generate unsigned withdraw all reward then delegate reward tx
func (c *Client) GenMultiSigRawWithdrawAllRewardTx(delAddr types.AccAddress, height int64) ([]byte, error) {
	delValsRes, err := c.QueryDelegations(delAddr, height)
	if err != nil {
		return nil, err
	}

	delegations := delValsRes.GetDelegationResponses()
	// build multi-message transaction
	msgs := make([]types.Msg, 0)
	for _, delegation := range delegations {
		valAddr := delegation.Delegation.ValidatorAddress
		val, err := types.ValAddressFromBech32(valAddr)
		if err != nil {
			return nil, err
		}
		//gen withdraw
		msg := xDistriTypes.NewMsgWithdrawDelegatorReward(delAddr, val)
		if err := msg.ValidateBasic(); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)

	}
	return c.GenMultiSigRawTx(msgs...)
}

//generate unsigned withdraw all reward then delegate reward tx
func (c *Client) GenMultiSigRawWithdrawAllRewardThenDeleTx(delAddr types.AccAddress, height int64) ([]byte, error) {
	delValsRes, err := c.QueryDelegations(delAddr, height)
	if err != nil {
		return nil, err
	}
	totalReward, err := c.QueryDelegationTotalRewards(delAddr, height)
	if err != nil {
		return nil, err
	}
	rewards := make(map[string]types.Coin)
	for _, r := range totalReward.Rewards {
		rewards[r.ValidatorAddress] = types.NewCoin(c.GetDenom(), r.Reward.AmountOf(c.GetDenom()).TruncateInt())
	}

	delegations := delValsRes.GetDelegationResponses()
	// build multi-message transaction
	msgs := make([]types.Msg, 0)
	for _, delegation := range delegations {
		valAddr := delegation.Delegation.ValidatorAddress
		//must filter zero value or tx will failure
		if rewards[valAddr].IsZero() {
			continue
		}

		val, err := types.ValAddressFromBech32(valAddr)
		if err != nil {
			return nil, err
		}
		//gen withdraw
		msg := xDistriTypes.NewMsgWithdrawDelegatorReward(delAddr, val)
		if err := msg.ValidateBasic(); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)

		//gen delegate
		msg2 := xStakingTypes.NewMsgDelegate(delAddr, val, rewards[valAddr])
		if err := msg2.ValidateBasic(); err != nil {
			return nil, err
		}

		msgs = append(msgs, msg2)
	}

	if len(msgs) == 0 {
		return nil, ErrNoMsgs
	}

	return c.GenMultiSigRawTx(msgs...)
}

//generate unsigned delegate reward tx
func (c *Client) GenMultiSigRawDeleRewardTx(delAddr types.AccAddress, height int64) ([]byte, error) {
	delValsRes, err := c.QueryDelegations(delAddr, height)
	if err != nil {
		return nil, err
	}
	totalReward, err := c.QueryDelegationTotalRewards(delAddr, height)
	if err != nil {
		return nil, err
	}
	rewards := make(map[string]types.Coin)
	for _, r := range totalReward.Rewards {
		rewards[r.ValidatorAddress] = types.NewCoin(c.GetDenom(), r.Reward.AmountOf(c.GetDenom()).TruncateInt())
	}

	delegations := delValsRes.GetDelegationResponses()
	// build multi-message transaction
	msgs := make([]types.Msg, 0)
	for _, delegation := range delegations {
		valAddr := delegation.Delegation.ValidatorAddress
		//must filter zero value or tx will failure
		if rewards[valAddr].IsZero() {
			continue
		}

		val, err := types.ValAddressFromBech32(valAddr)
		if err != nil {
			return nil, err
		}

		//gen delegate
		msg2 := xStakingTypes.NewMsgDelegate(delAddr, val, rewards[valAddr])
		if err := msg2.ValidateBasic(); err != nil {
			return nil, err
		}

		msgs = append(msgs, msg2)
	}

	if len(msgs) == 0 {
		return nil, ErrNoMsgs
	}

	return c.GenMultiSigRawTx(msgs...)
}

//c.clientCtx.FromAddress must be multi sig address,no need sequence
func (c *Client) GenMultiSigRawTx(msgs ...types.Msg) ([]byte, error) {
	cmd := cobra.Command{}
	txf := clientTx.NewFactoryCLI(c.clientCtx, cmd.Flags())
	txf = txf.WithAccountNumber(c.accountNumber).
		WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON). //multi sig need this mod
		WithGasAdjustment(1.5).
		WithGasPrices(c.gasPrice).
		WithGas(1500000).
		WithSimulateAndExecute(true)

	txBuilderRaw, err := clientTx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return nil, err
	}
	return c.clientCtx.TxConfig.TxJSONEncoder()(txBuilderRaw.GetTx())
}

//c.clientCtx.FromAddress  must be multi sig address
func (c *Client) SignMultiSigRawTx(rawTx []byte, fromSubKey string) (signature []byte, err error) {
	account, err := c.GetAccount()
	if err != nil {
		return nil, err
	}

	cmd := cobra.Command{}
	txf := clientTx.NewFactoryCLI(c.clientCtx, cmd.Flags())
	txf = txf.WithSequence(account.GetSequence()).
		WithAccountNumber(account.GetAccountNumber()).
		WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON) //multi sig need this mod

	tx, err := c.clientCtx.TxConfig.TxJSONDecoder()(rawTx)
	if err != nil {
		return nil, err
	}
	txBuilder, err := c.clientCtx.TxConfig.WrapTxBuilder(tx)
	if err != nil {
		return nil, err
	}
	err = xAuthClient.SignTxWithSignerAddress(txf, c.clientCtx, c.clientCtx.GetFromAddress(), fromSubKey, txBuilder, true, true)
	if err != nil {
		return nil, err
	}
	return marshalSignatureJSON(c.clientCtx.TxConfig, txBuilder, true)
}

//c.clientCtx.FromAddress  must be multi sig address
func (c *Client) SignMultiSigRawTxWithSeq(sequence uint64, rawTx []byte, fromSubKey string) (signature []byte, err error) {
	cmd := cobra.Command{}
	txf := clientTx.NewFactoryCLI(c.clientCtx, cmd.Flags())
	txf = txf.WithSequence(sequence).
		WithAccountNumber(c.accountNumber).
		WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON) //multi sig need this mod

	tx, err := c.clientCtx.TxConfig.TxJSONDecoder()(rawTx)
	if err != nil {
		return nil, err
	}
	txBuilder, err := c.clientCtx.TxConfig.WrapTxBuilder(tx)
	if err != nil {
		return nil, err
	}
	err = xAuthClient.SignTxWithSignerAddress(txf, c.clientCtx, c.clientCtx.GetFromAddress(), fromSubKey, txBuilder, true, true)
	if err != nil {
		return nil, err
	}
	return marshalSignatureJSON(c.clientCtx.TxConfig, txBuilder, true)
}

//assemble multiSig tx bytes for broadcast
func (c *Client) AssembleMultiSigTx(rawTx []byte, signatures [][]byte, threshold uint32) (txHash, txBts []byte, err error) {

	multisigInfo, err := c.clientCtx.Keyring.Key(c.clientCtx.FromName)
	if err != nil {
		return
	}
	if multisigInfo.GetType() != keyring.TypeMulti {
		return nil, nil, fmt.Errorf("%q must be of type %s: %s",
			c.clientCtx.FromName, keyring.TypeMulti, multisigInfo.GetType())
	}
	multiSigPub := multisigInfo.GetPubKey().(*kMultiSig.LegacyAminoPubKey)

	tx, err := c.clientCtx.TxConfig.TxJSONDecoder()(rawTx)
	if err != nil {
		return nil, nil, err
	}
	txBuilder, err := c.clientCtx.TxConfig.WrapTxBuilder(tx)
	if err != nil {
		return nil, nil, err
	}

	willUseSigs := make([]signing.SignatureV2, 0)
	for _, s := range signatures {
		ss, err := c.clientCtx.TxConfig.UnmarshalSignatureJSON(s)
		if err != nil {
			return nil, nil, err
		}
		willUseSigs = append(willUseSigs, ss...)
	}

	multiSigData := multisig.NewMultisig(len(multiSigPub.PubKeys))
	var useSequence uint64

	correntSigNumber := uint32(0)
	for i, sig := range willUseSigs {
		if correntSigNumber == threshold {
			break
		}
		//check sequence
		if i == 0 {
			useSequence = sig.Sequence
		} else {
			if useSequence != sig.Sequence {
				continue
			}
		}
		//check sig
		signingData := xAuthSigning.SignerData{
			ChainID:       c.clientCtx.ChainID,
			AccountNumber: c.accountNumber,
			Sequence:      useSequence,
		}

		err = xAuthSigning.VerifySignature(sig.PubKey, signingData, sig.Data, c.clientCtx.TxConfig.SignModeHandler(), txBuilder.GetTx())
		if err != nil {
			continue
		}

		if err := multisig.AddSignatureV2(multiSigData, sig, multiSigPub.GetPubKeys()); err != nil {
			continue
		}
		correntSigNumber++
	}

	if correntSigNumber != threshold {
		return nil, nil, fmt.Errorf("correct sig number:%d  threshold %d", correntSigNumber, threshold)
	}

	sigV2 := signing.SignatureV2{
		PubKey:   multiSigPub,
		Data:     multiSigData,
		Sequence: useSequence,
	}

	txBuilder.SetSignatures(sigV2)
	txBytes, err := c.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, nil, err
	}
	tendermintTx := tendermintTypes.Tx(txBytes)
	return tendermintTx.Hash(), tendermintTx, nil
}
