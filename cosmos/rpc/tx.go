package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xAuthClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	xDistriTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	xStakeTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (c *Client) SingleTransferTo(toAddr types.AccAddress, amount types.Coins) error {
	msg := xBankTypes.NewMsgSend(c.clientCtx.GetFromAddress(), toAddr, amount)
	cmd := cobra.Command{}
	return clientTx.GenerateOrBroadcastTxCLI(c.clientCtx, cmd.Flags(), msg)
}

func (c *Client) SingleReDelegate(srcValAddr, desValAddr types.ValAddress, amount types.Coin) (string, error) {
	msg := xStakeTypes.NewMsgBeginRedelegate(c.clientCtx.GetFromAddress(), srcValAddr, desValAddr, amount)
	account, err := c.clientCtx.AccountRetriever.GetAccount(c.clientCtx, c.clientCtx.GetFromAddress())
	if err != nil {
		return "", err
	}
	cmd := cobra.Command{}
	txf := clientTx.NewFactoryCLI(c.clientCtx, cmd.Flags())
	txf = txf.WithSequence(account.GetSequence()).
		WithAccountNumber(account.GetAccountNumber()).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT). //multi sig need this mod
		WithGasAdjustment(1.5).
		WithGas(0).
		WithGasPrices(c.gasPrice).
		WithSimulateAndExecute(true)

	//auto cal gas
	_, adjusted, err := clientTx.CalculateGas(c.clientCtx.QueryWithData, txf, msg)
	if err != nil {
		return "", err
	}
	txf = txf.WithGas(adjusted)

	txBuilderRaw, err := clientTx.BuildUnsignedTx(txf, msg)
	if err != nil {
		return "", err
	}
	xAuthClient.SignTx(txf, c.clientCtx, c.clientCtx.GetFromName(), txBuilderRaw, true, true)

	txBytes, err := c.clientCtx.TxConfig.TxEncoder()(txBuilderRaw.GetTx())
	if err != nil {
		return "", err
	}
	return c.BroadcastTx(txBytes)
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

func (c *Client) QueryDelegationRewards(delegatorAddr types.AccAddress, validatorAddr types.ValAddress) (*xDistriTypes.QueryDelegationRewardsResponse, error) {
	queryClient := xDistriTypes.NewQueryClient(c.clientCtx)
	return queryClient.DelegationRewards(
		context.Background(),
		&xDistriTypes.QueryDelegationRewardsRequest{DelegatorAddress: delegatorAddr.String(), ValidatorAddress: validatorAddr.String()},
	)
}

func (c *Client) QueryDelegationTotalRewards(delegatorAddr types.AccAddress) (*xDistriTypes.QueryDelegationTotalRewardsResponse, error) {
	queryClient := xDistriTypes.NewQueryClient(c.clientCtx)
	return queryClient.DelegationTotalRewards(
		context.Background(),
		&xDistriTypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: delegatorAddr.String()},
	)
}

func (c *Client) QueryBlock(height int64) (*ctypes.ResultBlock, error) {
	node, err := c.clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	return node.Block(context.Background(), &height)
}

func (c *Client) QueryAccount(addr types.AccAddress) (client.Account, error) {
	return c.clientCtx.AccountRetriever.GetAccount(c.clientCtx, addr)
}
