package rpc

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	xAuthTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	xAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	xDistriTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	xStakeTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"net"
	"net/url"
	"syscall"
	"time"
)

const retryLimit = 60
const waitTime = time.Millisecond * 500

var ErrNoTxIncludeWithdraw = fmt.Errorf("no tx include withdraw")
var ErrNoRewardNeedDelegate = fmt.Errorf("no tx reward need delegate")

func GetMemo(era uint32, txType string) string {
	return fmt.Sprintf("%d:%s", era, txType)
}

//no 0x prefix
func (c *Client) QueryTxByHash(hashHexStr string) (*types.TxResponse, error) {
	cc, err := retry(func() (interface{}, error) {
		return xAuthTx.QueryTx(c.clientCtx, hashHexStr)
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

func (c *Client) QueryUnbondingDelegation(delegatorAddr types.AccAddress, validatorAddr types.ValAddress, height int64) (*xStakeTypes.QueryUnbondingDelegationResponse, error) {
	client := c.clientCtx.WithHeight(height)
	queryClient := xStakeTypes.NewQueryClient(client)
	params := &xStakeTypes.QueryUnbondingDelegationRequest{
		DelegatorAddr: delegatorAddr.String(),
		ValidatorAddr: validatorAddr.String(),
	}

	cc, err := retry(func() (interface{}, error) {
		return queryClient.UnbondingDelegation(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*xStakeTypes.QueryUnbondingDelegationResponse), nil
}

func (c *Client) QueryDelegations(delegatorAddr types.AccAddress, height int64) (*xStakeTypes.QueryDelegatorDelegationsResponse, error) {
	retry := 0
	var err error
	var delegations *xStakeTypes.QueryDelegatorDelegationsResponse
	for {
		if retry > retryLimit {
			return nil, fmt.Errorf("QueryDelegationsWithRetry reach retry: %s", err)
		}
		delegations, err = c.queryDelegations(delegatorAddr, height)
		if err != nil || len(delegations.DelegationResponses) == 0 {
			time.Sleep(waitTime)
			retry++
			continue
		}
		break
	}
	return delegations, nil
}

func (c *Client) queryDelegations(delegatorAddr types.AccAddress, height int64) (*xStakeTypes.QueryDelegatorDelegationsResponse, error) {
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

func (c *Client) GetTxs(events []string, page, limit int, orderBy string) (*types.SearchTxsResult, error) {
	cc, err := retry(func() (interface{}, error) {
		return xAuthTx.QueryTxsByEvents(c.clientCtx, events, page, limit, orderBy)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*types.SearchTxsResult), nil
}

func (c *Client) GetRewardToBeDelegated(delegatorAddr string, era uint32) (map[string]types.Coin, int64, error) {
	moduleAddressStr := xAuthTypes.NewModuleAddress(xDistriTypes.ModuleName).String()
	delAddress, err := types.AccAddressFromBech32(delegatorAddr)
	if err != nil {
		return nil, 0, err
	}

	txs, err := c.GetTxs(
		[]string{
			fmt.Sprintf("transfer.recipient='%s'", delegatorAddr),
			fmt.Sprintf("transfer.sender='%s'", moduleAddressStr),
		}, 1, 4, "desc")
	if err != nil {
		return nil, 0, err
	}

	if len(txs.Txs) == 0 {
		return nil, 0, ErrNoRewardNeedDelegate
	}

	valRewards := make(map[string]types.Coin)
	retHeight := int64(0)
	for _, tx := range txs.Txs {
		txValue := tx.Tx.Value

		decodeTx, err := c.GetTxConfig().TxDecoder()(txValue)
		if err != nil {
			return nil, 0, err
		}
		memoTx, ok := decodeTx.(types.TxWithMemo)
		if !ok {
			return nil, 0, fmt.Errorf("tx is not type TxWithMemo, txhash: %s", txs.Txs[0].TxHash)
		}
		memoInTx := memoTx.GetMemo()

		switch memoInTx {
		case GetMemo(era, TxTypeHandleEraPoolUpdatedEvent):
			//return tx handleEraPoolUpdatedEvent height
			retHeight = tx.Height - 1
			fallthrough
		case GetMemo(era-1, TxTypeHandleBondReportedEvent):
			height := tx.Height - 1
			totalReward, err := c.QueryDelegationTotalRewards(delAddress, height)
			if err != nil {
				return nil, 0, err
			}

			for _, r := range totalReward.Rewards {
				rewardCoin := types.NewCoin(c.GetDenom(), r.Reward.AmountOf(c.GetDenom()).TruncateInt())
				if rewardCoin.IsZero() {
					continue
				}
				if _, exist := valRewards[r.ValidatorAddress]; !exist {
					valRewards[r.ValidatorAddress] = rewardCoin
				} else {
					valRewards[r.ValidatorAddress] = valRewards[r.ValidatorAddress].Add(rewardCoin)
				}

			}
		default:
		}
	}

	if len(valRewards) == 0 {
		return nil, 0, ErrNoRewardNeedDelegate
	}

	return valRewards, retHeight, nil
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
		newErr := t.Unwrap()
		return isConnectionError(newErr)
	}

	return false
}

type wrapError interface {
	Unwrap() error
}
