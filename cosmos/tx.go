package cosmos

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xAuthClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
)

func (c *Client) TransferTo(toAddr types.AccAddress, amount types.Coins) error {
	msg := xBankTypes.NewMsgSend(c.clientCtx.GetFromAddress(), toAddr, amount)
	cmd := cobra.Command{}
	return clientTx.GenerateOrBroadcastTxCLI(c.clientCtx, cmd.Flags(), msg)
}

func (c *Client) GenRawTransferTx(toAddr types.AccAddress, amount types.Coins) ([]byte, error) {
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

func (c *Client) SignRawTx(rawTx []byte, fromKey string) (signature []byte, err error) {
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
	xAuthClient.SignTxWithSignerAddress(txf, c.clientCtx, c.clientCtx.GetFromAddress(), fromKey, txBuilder, true, false)

	return marshalSignatureJSON(c.clientCtx.TxConfig, txBuilder, true)
}

func (c *Client) QueryTxByHash(hashHexStr string) (*types.TxResponse, error) {
	return xAuthClient.QueryTx(c.clientCtx, hashHexStr)
}
func (c *Client) BroadcastTx(tx []byte) error {
	res, err := c.clientCtx.BroadcastTx(tx)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return errors.New(fmt.Sprintf("Boradcast err with res.code: %d", res.Code))
	}
	return nil
}

func marshalSignatureJSON(txConfig client.TxConfig, txBldr client.TxBuilder, signatureOnly bool) ([]byte, error) {
	parsedTx := txBldr.GetTx()
	if signatureOnly {
		sigs, err := parsedTx.GetSignaturesV2()
		if err != nil {
			return nil, err
		}
		return txConfig.MarshalSignatureJSON(sigs)
	}

	return txConfig.TxJSONEncoder()(parsedTx)
}
