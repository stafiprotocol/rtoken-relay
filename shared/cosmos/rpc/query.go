package rpc

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	xAuthClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	xDistriTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	xStakeTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"net"
	"net/url"
	"syscall"
	"time"
)

const retryLimit = 10
const waitTime = time.Millisecond * 500

//no 0x prefix
func (c *Client) QueryTxByHash(hashHexStr string) (*types.TxResponse, error) {
	cc, err := retry(func() (interface{}, error) {
		return xAuthClient.QueryTx(c.clientCtx, hashHexStr)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*types.TxResponse), nil
}

func (c *Client) QueryDelegation(delegatorAddr types.AccAddress, validatorAddr types.ValAddress, height int64) (*xStakeTypes.QueryDelegationResponse, error) {
	client := c.clientCtx.WithHeight(height)
	queryClient := xStakeTypes.NewQueryClient(client)
	params := &xStakeTypes.QueryDelegationRequest{
		DelegatorAddr: delegatorAddr.String(),
		ValidatorAddr: validatorAddr.String(),
	}

	cc, err := retry(func() (interface{}, error) {
		return queryClient.Delegation(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*xStakeTypes.QueryDelegationResponse), nil
}

func (c *Client) QueryDelegations(delegatorAddr types.AccAddress, height int64) (*xStakeTypes.QueryDelegatorDelegationsResponse, error) {
	client := c.clientCtx.WithHeight(height)
	queryClient := xStakeTypes.NewQueryClient(client)
	params := &xStakeTypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: delegatorAddr.String(),
		Pagination:    &query.PageRequest{},
	}
	cc, err := retry(func() (interface{}, error) {
		return queryClient.DelegatorDelegations(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*xStakeTypes.QueryDelegatorDelegationsResponse), nil
}

func (c *Client) QueryDelegationRewards(delegatorAddr types.AccAddress, validatorAddr types.ValAddress, height int64) (*xDistriTypes.QueryDelegationRewardsResponse, error) {
	client := c.clientCtx.WithHeight(height)
	queryClient := xDistriTypes.NewQueryClient(client)
	cc, err := retry(func() (interface{}, error) {
		return queryClient.DelegationRewards(
			context.Background(),
			&xDistriTypes.QueryDelegationRewardsRequest{DelegatorAddress: delegatorAddr.String(), ValidatorAddress: validatorAddr.String()},
		)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*xDistriTypes.QueryDelegationRewardsResponse), nil
}

func (c *Client) QueryDelegationTotalRewards(delegatorAddr types.AccAddress, height int64) (*xDistriTypes.QueryDelegationTotalRewardsResponse, error) {
	client := c.clientCtx.WithHeight(height)
	queryClient := xDistriTypes.NewQueryClient(client)

	cc, err := retry(func() (interface{}, error) {
		return queryClient.DelegationTotalRewards(
			context.Background(),
			&xDistriTypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: delegatorAddr.String()},
		)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*xDistriTypes.QueryDelegationTotalRewardsResponse), nil
}

func (c *Client) QueryBlock(height int64) (*ctypes.ResultBlock, error) {
	node, err := c.clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	cc, err := retry(func() (interface{}, error) {
		return node.Block(context.Background(), &height)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*ctypes.ResultBlock), nil
}

func (c *Client) QueryAccount(addr types.AccAddress) (client.Account, error) {
	return c.getAccount(0, addr)
}

func (c *Client) GetSequence(height int64, addr types.AccAddress) (uint64, error) {
	account, err := c.getAccount(height, addr)
	if err != nil {
		return 0, err
	}
	return account.GetSequence(), nil
}

func (c *Client) QueryBalance(addr types.AccAddress, denom string, height int64) (*xBankTypes.QueryBalanceResponse, error) {
	client := c.clientCtx.WithHeight(height)
	queryClient := xBankTypes.NewQueryClient(client)
	params := xBankTypes.NewQueryBalanceRequest(addr, denom)

	cc, err := retry(func() (interface{}, error) {
		return queryClient.Balance(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*xBankTypes.QueryBalanceResponse), nil
}

func (c *Client) GetCurrentBLockHeight() (int64, error) {
	status, err := c.GetStatus()
	if err != nil {
		return 0, err
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

func (c *Client) GetStatus() (*ctypes.ResultStatus, error) {
	cc, err := retry(func() (interface{}, error) {
		return c.clientCtx.Client.Status(context.Background())
	})
	if err != nil {
		return nil, err
	}
	return cc.(*ctypes.ResultStatus), nil
}

func (c *Client) GetAccount() (client.Account, error) {
	return c.getAccount(0, c.clientCtx.FromAddress)
}

func (c *Client) getAccount(height int64, addr types.AccAddress) (client.Account, error) {
	cc, err := retry(func() (interface{}, error) {
		client := c.clientCtx.WithHeight(height)
		return client.AccountRetriever.GetAccount(c.clientCtx, addr)
	})
	if err != nil {
		return nil, err
	}
	return cc.(client.Account), nil
}

//only retry func when return connection err here
func retry(f func() (interface{}, error)) (interface{}, error) {
	var err error
	var result interface{}
	for i := 0; i < retryLimit; i++ {
		result, err = f()
		if err != nil && isConnectionError(err) {
			time.Sleep(waitTime)
			continue
		}
		return result, err
	}
	panic(fmt.Sprintf("reach retry limit. err: %s", err))
}

func isConnectionError(err error) bool {
	switch t := err.(type) {
	case *url.Error:
		if t.Timeout() || t.Temporary() {
			return true
		}
		return isConnectionError(t.Err)
	}

	switch t := err.(type) {
	case *net.OpError:
		fmt.Println("operror", t.Op)
		if t.Op == "dial" || t.Op == "read" {
			return true
		}
		return isConnectionError(t.Err)

	case syscall.Errno:
		if t == syscall.ECONNREFUSED {
			return true
		}
	}

	switch t := err.(type) {
	case wrapError:
		fmt.Println("wrapError")
		newErr := t.Unwrap()
		return isConnectionError(newErr)
	}

	return false
}

type wrapError interface {
	Unwrap() error
}
