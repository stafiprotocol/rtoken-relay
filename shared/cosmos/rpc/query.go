package rpc

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	xAuthClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	xDistriTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	xStakeTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

//no 0x prefix
func (c *Client) QueryTxByHash(hashHexStr string) (*types.TxResponse, error) {
	return xAuthClient.QueryTx(c.clientCtx, hashHexStr)
}

func (c *Client) QueryDelegation(delegatorAddr types.AccAddress, validatorAddr types.ValAddress) (*xStakeTypes.QueryDelegationResponse, error) {
	client := c.clientCtx
	queryClient := xStakeTypes.NewQueryClient(client)
	params := &xStakeTypes.QueryDelegationRequest{
		DelegatorAddr: delegatorAddr.String(),
		ValidatorAddr: validatorAddr.String(),
	}

	return queryClient.Delegation(context.Background(), params)
}

func (c *Client) QueryDelegations(delegatorAddr types.AccAddress, height int64) (*xStakeTypes.QueryDelegatorDelegationsResponse, error) {
	client := c.clientCtx.WithHeight(height)
	queryClient := xStakeTypes.NewQueryClient(client)
	params := &xStakeTypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: delegatorAddr.String(),
		Pagination:    &query.PageRequest{},
	}

	return queryClient.DelegatorDelegations(context.Background(), params)
}

func (c *Client) QueryDelegationRewards(delegatorAddr types.AccAddress, validatorAddr types.ValAddress, height int64) (*xDistriTypes.QueryDelegationRewardsResponse, error) {
	client := c.clientCtx.WithHeight(height)
	queryClient := xDistriTypes.NewQueryClient(client)
	return queryClient.DelegationRewards(
		context.Background(),
		&xDistriTypes.QueryDelegationRewardsRequest{DelegatorAddress: delegatorAddr.String(), ValidatorAddress: validatorAddr.String()},
	)
}

func (c *Client) QueryDelegationTotalRewards(delegatorAddr types.AccAddress, height int64) (*xDistriTypes.QueryDelegationTotalRewardsResponse, error) {
	client := c.clientCtx.WithHeight(height)
	queryClient := xDistriTypes.NewQueryClient(client)
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

func (c *Client) GetSequence(height int64, addr types.AccAddress) (uint64, error) {
	client := c.clientCtx.WithHeight(height)
	account, err := client.AccountRetriever.GetAccount(client, addr)
	if err != nil {
		return 0, err
	}

	return account.GetSequence(), nil
}

func (c *Client) QueryBalance(addr types.AccAddress, denom string, height int64) (*xBankTypes.QueryBalanceResponse, error) {
	client := c.clientCtx.WithHeight(height)
	queryClient := xBankTypes.NewQueryClient(client)
	params := xBankTypes.NewQueryBalanceRequest(addr, denom)
	return queryClient.Balance(context.Background(), params)
}

func (c *Client) GetCurrentBLockHeight() (int64, error) {
	status, err := c.GetStatus()
	if err != nil {
		return 0, err
	}
	return status.SyncInfo.LatestBlockHeight, nil
}
