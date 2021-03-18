package rpc

import (
	"fmt"
	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	kMultiSig "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xAuthClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	xDistriTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	xStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

//c.clientCtx.FromAddress must be multi sig address
func (c *Client) GenMultiSigRawTransferTx(toAddr types.AccAddress, amount types.Coins) ([]byte, error) {
	msg := xBankTypes.NewMsgSend(c.clientCtx.GetFromAddress(), toAddr, amount)
	return c.GenMultiSigRawTx(msg)
}

//only support one type coin
func (c *Client) GenMultiSigRawBatchTransferTx(outs []xBankTypes.Output) ([]byte, error) {
	totalAmount := types.NewInt(0)
	for _, out := range outs {
		for _, coin := range out.Coins {
			totalAmount = totalAmount.Add(coin.Amount)
		}
	}
	input := xBankTypes.Input{
		Address: c.clientCtx.GetFromAddress().String(),
		Coins:   types.NewCoins(types.NewCoin(c.denom, totalAmount))}

	msg := xBankTypes.NewMsgMultiSend([]xBankTypes.Input{input}, outs)
	return c.GenMultiSigRawTx(msg)
}

//generate unsigned delegate tx
func (c *Client) GenMultiSigRawDelegateTx(delAddr types.AccAddress, valAddr types.ValAddress, amount types.Coin) ([]byte, error) {
	msg := xStakingTypes.NewMsgDelegate(delAddr, valAddr, amount)
	return c.GenMultiSigRawTx(msg)
}

//generate unsigned reDelegate tx
func (c *Client) GenMultiSigRawReDelegateTx(delAddr types.AccAddress, valSrcAddr, valDstAddr types.ValAddress, amount types.Coin) ([]byte, error) {
	msg := xStakingTypes.NewMsgBeginRedelegate(delAddr, valSrcAddr, valDstAddr, amount)
	return c.GenMultiSigRawTx(msg)
}

//generate unsigned unDelegate tx
func (c *Client) GenMultiSigRawUnDelegateTx(delAddr types.AccAddress, valAddr types.ValAddress, amount types.Coin) ([]byte, error) {
	msg := xStakingTypes.NewMsgUndelegate(delAddr, valAddr, amount)
	return c.GenMultiSigRawTx(msg)
}

////generate unsigned withdraw delegate reward tx
func (c *Client) GenMultiSigRawWithdrawDeleRewardTx(delAddr types.AccAddress, valAddr types.ValAddress) ([]byte, error) {
	msg := xDistriTypes.NewMsgWithdrawDelegatorReward(delAddr, valAddr)
	return c.GenMultiSigRawTx(msg)
}

//c.clientCtx.FromAddress must be multi sig address
func (c *Client) GenMultiSigRawTx(msgs ...types.Msg) ([]byte, error) {
	account, err := c.clientCtx.AccountRetriever.GetAccount(c.clientCtx, c.clientCtx.GetFromAddress())
	if err != nil {
		return nil, err
	}
	cmd := cobra.Command{}
	txf := clientTx.NewFactoryCLI(c.clientCtx, cmd.Flags())
	txf = txf.WithSequence(account.GetSequence()).
		WithAccountNumber(account.GetAccountNumber()).
		WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON). //multi sig need this mod
		WithGasAdjustment(1.5).
		WithGasPrices(c.gasPrice).
		WithGas(300000).
		WithSimulateAndExecute(true)

	//todo fix auto cal gas
	//_, adjusted, err := clientTx.CalculateGas(c.clientCtx.QueryWithData, txf, msgs...)
	//if err != nil {
	//	return nil, err
	//}
	//txf = txf.WithGas(adjusted)

	txBuilderRaw, err := clientTx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return nil, err
	}
	return c.clientCtx.TxConfig.TxJSONEncoder()(txBuilderRaw.GetTx())
}

//c.clientCtx.FromAddress  must be multi sig address
func (c *Client) SignMultiSigRawTx(rawTx []byte, fromKey string) (signature []byte, err error) {
	account, err := c.clientCtx.AccountRetriever.GetAccount(c.clientCtx, c.clientCtx.GetFromAddress())
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
	xAuthClient.SignTxWithSignerAddress(txf, c.clientCtx, c.clientCtx.GetFromAddress(), fromKey, txBuilder, true, true)
	if err != nil {
		return nil, err
	}
	return marshalSignatureJSON(c.clientCtx.TxConfig, txBuilder, true)
}

//assemble multiSig tx bytes for broadcast
func (c *Client) AssembleMultiSigTx(rawTx []byte, signatures [][]byte) (txBts []byte, err error) {
	accountMultiSign, err := c.clientCtx.AccountRetriever.GetAccount(c.clientCtx, c.clientCtx.GetFromAddress())
	if err != nil {
		return nil, err
	}

	multisigInfo, err := c.clientCtx.Keyring.Key(c.clientCtx.FromName)
	if err != nil {
		return
	}
	if multisigInfo.GetType() != keyring.TypeMulti {
		return nil, fmt.Errorf("%q must be of type %s: %s", c.clientCtx.FromName, keyring.TypeMulti, multisigInfo.GetType())
	}

	multiSigPub := multisigInfo.GetPubKey().(*kMultiSig.LegacyAminoPubKey)

	willUseSigs := make([]signing.SignatureV2, 0)
	for _, s := range signatures {
		ss, err := c.clientCtx.TxConfig.UnmarshalSignatureJSON(s)
		if err != nil {
			return nil, err
		}
		willUseSigs = append(willUseSigs, ss...)
	}

	multiSigData := multisig.NewMultisig(len(multiSigPub.PubKeys))
	//todo check sig
	for _, sig := range willUseSigs {
		if err := multisig.AddSignatureV2(multiSigData, sig, multiSigPub.GetPubKeys()); err != nil {
			return nil, err
		}
	}

	sigV2 := signing.SignatureV2{
		PubKey:   multiSigPub,
		Data:     multiSigData,
		Sequence: accountMultiSign.GetSequence(),
	}

	tx, err := c.clientCtx.TxConfig.TxJSONDecoder()(rawTx)
	if err != nil {
		return nil, err
	}
	txBuilder, err := c.clientCtx.TxConfig.WrapTxBuilder(tx)
	if err != nil {
		return nil, err
	}
	txBuilder.SetSignatures(sigV2)
	return c.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
}
