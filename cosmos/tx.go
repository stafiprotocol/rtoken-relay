package cosmos

import (
	"errors"
	"fmt"
	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types"
	xAuthClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
)

func (c *Client) TransferTo(toAddr types.AccAddress, amount types.Coins) error {
	msg := xBankTypes.NewMsgSend(c.clientCtx.GetFromAddress(), toAddr, amount)
	cmd := cobra.Command{}
	return clientTx.GenerateOrBroadcastTxCLI(c.clientCtx, cmd.Flags(), msg)
}

func (c *Client) QueryTxByHash(hashHexStr string) (*types.TxResponse, error) {
	return xAuthClient.QueryTx(c.clientCtx, hashHexStr)
}

func (c *Client) BroadcastTx(tx []byte) (string, error) {
	res, err := c.clientCtx.BroadcastTx(tx)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", errors.New(fmt.Sprintf("Boradcast err with res.code: %d", res.Code))
	}
	return res.TxHash, nil
}
