package cosmos

import (
	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	kMultiSig "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xAuthClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
)

func (c *Client) GenMultiSigRawTransferTx(toAddr types.AccAddress, amount types.Coins) ([]byte, error) {
	account, err := c.clientCtx.AccountRetriever.GetAccount(c.clientCtx, c.clientCtx.GetFromAddress())
	if err != nil {
		return nil, err
	}

	msg := xBankTypes.NewMsgSend(c.clientCtx.GetFromAddress(), toAddr, amount)
	cmd := cobra.Command{}
	txf := clientTx.NewFactoryCLI(c.clientCtx, cmd.Flags())
	txf = txf.WithSequence(account.GetSequence()).
		WithAccountNumber(account.GetAccountNumber()).
		WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON) //multi sig need this mod

	txBuilderRaw, err := clientTx.BuildUnsignedTx(txf, msg)
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

func (c *Client) CreateMultiSigTx(rawTx []byte, signatures [][]byte) (txBts []byte, err error) {
	accountMultiSign, err := c.clientCtx.AccountRetriever.GetAccount(c.clientCtx, c.clientCtx.GetFromAddress())
	if err != nil {
		return nil, err
	}
	legacyPubkey := accountMultiSign.GetPubKey().(*kMultiSig.LegacyAminoPubKey)

	willUseSigs := make([]signing.SignatureV2, 0)
	for _, s := range signatures {
		ss, err := c.clientCtx.TxConfig.UnmarshalSignatureJSON(s)
		if err != nil {
			return nil, err
		}
		willUseSigs = append(willUseSigs, ss...)
	}

	multiSigData := multisig.NewMultisig(len(legacyPubkey.PubKeys))
	for _, sig := range willUseSigs {
		if err := multisig.AddSignatureV2(multiSigData, sig, legacyPubkey.GetPubKeys()); err != nil {
			return nil, err
		}
	}

	sigV2 := signing.SignatureV2{
		PubKey:   legacyPubkey,
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
