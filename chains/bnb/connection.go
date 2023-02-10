// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package bnb

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	multisig_onchain "github.com/stafiprotocol/rtoken-relay/bindings/MultisigOnchain"
	stake_portal "github.com/stafiprotocol/rtoken-relay/bindings/StakeNativePortal"
	staking "github.com/stafiprotocol/rtoken-relay/bindings/Staking"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/bnc"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

var (
	TxRetryInterval  = time.Second * 2
	ErrTxUnderpriced = errors.New("replacement transaction underpriced")
	ZeroAddress      = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

const (
	TxRetryLimit = 10
)

type Connection struct {
	symbol         core.RSymbol
	queryClient    *ethereum.Client // first pool's bsc client used for query
	bcApiEndpoint  string
	bscSideChainId string
	log            core.Logger
	stop           <-chan int

	accountClients         map[common.Address]*ethereum.Client // account address => client
	stakePortalContract    *stake_portal.StakeNativePortal
	stakingContract        *staking.Staking
	stakingContractAddress common.Address
}

func NewConnection(cfg *core.ChainConfig, log core.Logger, stop <-chan int, forTest bool) (*Connection, error) {
	log.Info("NewClient", "name", cfg.Name, "KeystorePath", cfg.KeystorePath, "Endpoint", cfg.Endpoint)

	bcEndpoint := cfg.Opts["bcEndpoint"]
	bcEndpointStr, ok := bcEndpoint.(string)
	if !ok {
		return nil, errors.New("bcEndpoint not exist")
	}

	bscSideChainId := cfg.Opts["bscSideChainId"]
	bscSideChainIdStr, ok := bscSideChainId.(string)
	if !ok {
		return nil, errors.New("bscSideChainId not exist")
	}

	ethClient, err := ethclient.Dial(cfg.Endpoint)
	if err != nil {
		return nil, err
	}
	bscClient := ethereum.NewClient(cfg.Endpoint, nil, log, big.NewInt(0), big.NewInt(0))
	if err := bscClient.Connect(); err != nil {
		return nil, err
	}

	acSize := len(cfg.Accounts)
	if acSize == 0 && !forTest {
		return nil, fmt.Errorf("account empty")
	}
	accounts := make(map[common.Address]*ethereum.Client)
	for i := 0; i < acSize; i++ {
		if !common.IsHexAddress(cfg.Accounts[i]) {
			return nil, fmt.Errorf("account not hex address, index: %d", i)
		}
		accountAddress := common.HexToAddress(cfg.Accounts[i])

		kpI, err := keystore.KeypairFromAddress(cfg.Accounts[i], keystore.EthChain, cfg.KeystorePath, cfg.Insecure)
		if err != nil {
			return nil, err
		}
		kp, ok := kpI.(*secp256k1.Keypair)
		if !ok {
			return nil, fmt.Errorf("secp256k1.Keypair cast failed")
		}

		client := ethereum.NewClient(cfg.Endpoint, kp, log, big.NewInt(0), big.NewInt(0))
		if err := client.Connect(); err != nil {
			return nil, err
		}
		accounts[accountAddress] = client
		if i == 0 {
			bscClient = client
		}
	}

	stakePortal, err := initStakePortal(cfg.Opts["StakePortalContract"], ethClient)
	if err != nil {
		return nil, err
	}

	staking, stakingAddr, err := initStaking(ethClient, cfg.Opts["StakingContract"])
	if err != nil {
		return nil, err
	}

	return &Connection{
		symbol:                 cfg.Symbol,
		queryClient:            bscClient,
		bcApiEndpoint:          bcEndpointStr,
		bscSideChainId:         bscSideChainIdStr,
		log:                    log,
		stop:                   stop,
		stakePortalContract:    stakePortal,
		stakingContract:        staking,
		stakingContractAddress: stakingAddr,
		accountClients:         accounts,
	}, nil
}

func initStakePortal(stakePortalCfg interface{}, conn *ethclient.Client) (*stake_portal.StakeNativePortal, error) {
	stakePortalAddr, ok := stakePortalCfg.(string)
	if !ok {
		return nil, errors.New("stakePortalCfg not ok")
	}
	if !common.IsHexAddress(stakePortalAddr) {
		return nil, errors.New("stakePortalCfg not hex string")
	}
	stakePortal, err := stake_portal.NewStakeNativePortal(common.HexToAddress(stakePortalAddr), conn)
	if err != nil {
		return nil, err
	}

	return stakePortal, nil
}

func initStaking(conn *ethclient.Client, stakingContractCfg interface{}) (*staking.Staking, common.Address, error) {
	stakingContractStr, ok := stakingContractCfg.(string)
	if !ok {
		return nil, common.Address{}, errors.New("stakingContractCfg not ok")
	}
	if !common.IsHexAddress(stakingContractStr) {
		return nil, common.Address{}, errors.New("stakingContractCfg not hex string")
	}

	addr := common.HexToAddress(stakingContractStr)
	staking, err := staking.NewStaking(addr, conn)
	if err != nil {
		return nil, common.Address{}, err
	}

	return staking, addr, nil
}

func (c *Connection) TransferVerify(r *submodel.BondRecord) (result submodel.BondReason, err error) {
	for i := 0; i < 5; i++ {
		result, err = c.queryClient.TransferVerifyBNB(r)
		if err != nil {
			return
		}

		if result == submodel.BlockhashUnmatch {
			time.Sleep(2 * time.Second)
			continue
		}
	}

	return
}

func (c *Connection) GetAccountClients(multisigContract *multisig_onchain.MultisigOnchain) ([]*ethereum.Client, error) {
	ret := make([]*ethereum.Client, 0)
	for account, client := range c.accountClients {
		index, err := multisigContract.GetSubAccountIndex(&bind.CallOpts{
			Context: context.Background(),
		}, account)
		if err != nil {
			return nil, err
		}
		if index.Int64() > 0 {
			ret = append(ret, client)
		}
	}

	return ret, nil
}

// use conn address as blockstore use address
func (c *Connection) BlockStoreUseAddress() string {
	return c.queryClient.Keypair().Address()
}

func (c *Connection) Close() {
	for _, c := range c.accountClients {
		c.Close()
	}
}

func (c *Connection) LatestBlockTimestamp() (uint64, error) {
	blkTime, err := c.queryClient.LatestBlockTimestamp()
	if err != nil {
		return 0, err
	}

	return blkTime, nil
}

func (c *Connection) GetStakingContractAddress() common.Address {
	return c.stakingContractAddress
}

func (c *Connection) LatestBlock() (uint64, error) {
	blk, err := c.queryClient.LatestBlock()
	if err != nil {
		return 0, err
	}

	return blk.Uint64(), nil
}

func (c *Connection) QueryBlock(number int64) (*types.Block, error) {
	blk, err := c.queryClient.Client().BlockByNumber(context.Background(), big.NewInt(number))
	if err != nil {
		return nil, err
	}

	return blk, nil
}

func (c *Connection) GetHeightTimestampByEra(era uint32, eraSeconds, offset int64) (int64, int64, error) {
	if int64(era) < offset {
		return 0, 0, fmt.Errorf("era: %d is less than offset: %d", era, offset)
	}
	targetTimestamp := (int64(era) - offset) * eraSeconds
	height, err := c.GetHeightByTimestamp(targetTimestamp)
	if err != nil {
		return 0, 0, err
	}
	return height, targetTimestamp, nil
}

func (c *Connection) GetHeightByTimestamp(targetTimestamp int64) (int64, error) {
	c.log.Trace("GetHeightByTimestamp", "targetTimestamp", targetTimestamp)

	blockNumber, timestamp, err := c.queryClient.LatestBlockAndTimestamp()
	if err != nil {
		return 0, err
	}
	seconds := int64(timestamp) - targetTimestamp
	if seconds < 0 {
		// return if over 20 minutes
		if seconds < -60*20 {
			return 0, fmt.Errorf("latest block timestamp: %d is less than targetTimestamp: %d", timestamp, targetTimestamp)
		}

		retry := 0
		for {
			if retry > BlockRetryLimit {
				return 0, fmt.Errorf("latest block timestamp: %d is less than targetTimestamp: %d", timestamp, targetTimestamp)
			}

			blockNumber, timestamp, err = c.queryClient.LatestBlockAndTimestamp()
			if err != nil {
				return 0, err
			}
			if int64(timestamp) < targetTimestamp {
				c.log.Warn(fmt.Sprintf("latest block timestamp: %d is less than targetTimestamp: %d, will wait...", timestamp, targetTimestamp))

				time.Sleep(WaitInterval)
				retry++

				continue
			}

			seconds = int64(timestamp) - targetTimestamp
			break
		}
	}

	tmpTargetBlock := int64(blockNumber) - seconds/3
	if tmpTargetBlock <= 0 {
		tmpTargetBlock = 1
	}

	var block *types.Block
	for {
		block, err = c.QueryBlock(tmpTargetBlock)
		if err != nil {
			return 0, errors.Wrap(err, "loop queryBlock")
		}
		du := int64(block.Time()) - targetTimestamp
		tmpTargetBlock = block.Number().Int64() - du/3
		if du < 20 && du > -20 {
			break
		}
		c.log.Trace("loop block", "block", block.Number().Int64(), "time", block.Time(), "tempBlock", tmpTargetBlock)
	}

	// return after blocknumber
	var afterBlockNumber int64
	var preBlockNumber int64
	if int64(block.Time()) > targetTimestamp {
		afterBlockNumber = block.Number().Int64()
		for {
			if afterBlockNumber <= 2 {
				return 1, nil
			}
			block, err := c.QueryBlock(afterBlockNumber - 1)
			if err != nil {
				return 0, err
			}
			c.log.Trace("afterBlock", "block", block.Number().Int64(), "time", block.Time())
			if int64(block.Time()) > targetTimestamp {
				afterBlockNumber = block.Number().Int64()
				continue
			}

			break
		}

	} else {
		preBlockNumber = block.Number().Int64()
		for {
			block, err := c.QueryBlock(preBlockNumber + 1)
			if err != nil {
				return 0, err
			}
			c.log.Trace("preBlock", "block", block.Number().Int64(), "time", block.Time())
			if int64(block.Time()) > targetTimestamp {
				afterBlockNumber = block.Number().Int64()
				break
			} else {
				preBlockNumber = block.Number().Int64()
			}
		}
	}

	return afterBlockNumber, nil
}

// return reward decimals 8
func (c *Connection) RewardOnBcDuTimes(pool common.Address, targetTimestamp, lastEraTimestamp int64) (int64, error) {
	delegator := utils.GetDelegaterAddressOnBc(pool[:]).String()

	for i := 0; i < TxRetryLimit; i++ {
		total, lastRewardTimestamp, err := c.RewardTotalTimesAndLastRewardTimestamp(delegator)
		if err != nil {
			c.log.Warn("totalAndLastHeight error", "err", err)
			if i+1 == TxRetryLimit {
				return 0, err
			}
			continue
		}
		c.log.Trace("RewardOnBcDuTimes", "total", total, "lastRewardTimestamp", lastRewardTimestamp, "delegator", delegator)
		if lastRewardTimestamp < lastEraTimestamp {
			return 0, nil
		}

		rewardSum, err := c.stakingReward(delegator, total, targetTimestamp, lastEraTimestamp)
		if err != nil {
			c.log.Error("stakingReward error", "err", err)
			if i+1 == TxRetryLimit {
				return 0, err
			}
			continue
		}

		return rewardSum, nil
	}

	return 0, fmt.Errorf("get reward failed")
}

func (c *Connection) RewardTotalTimesAndLastRewardTimestamp(delegator string) (int64, int64, error) {
	api := c.rewardApi(delegator, 1, 0)
	c.log.Trace("totalAndLastHeight rewardApi", "rewardApi", api)
	sr, err := bnc.GetStakingReward(api)
	if err != nil {
		return 0, 0, err
	}

	if len(sr.RewardDetails) == 0 {
		return 0, 0, nil
	}

	rewardTime, err := time.Parse(time.RFC3339, sr.RewardDetails[0].RewardTime)
	if err != nil {
		return 0, 0, err
	}
	return sr.Total, rewardTime.Unix(), nil
}

func (c *Connection) stakingReward(delegator string, total, targetTimestamp, lastEraTimestamp int64) (int64, error) {
	c.log.Trace("stakingReward", "delegator", delegator, "total", total, "targetTimestamp", targetTimestamp, "lastEraTimestamp", lastEraTimestamp)
	offset := int64(0)
	rewardSum := int64(0)

OUT:
	for i := total; i > 0; i -= 100 {
		api := c.rewardApi(delegator, 100, offset)

		sr, err := bnc.GetStakingReward(api)
		if err != nil {
			return 0, errors.Wrap(err, "bnc.GetStakingReward")
		}

		for _, rd := range sr.RewardDetails {
			rewardTime, err := time.Parse(time.RFC3339, rd.RewardTime)
			if err != nil {
				return 0, err
			}
			if rewardTime.Unix() > targetTimestamp {
				continue
			}

			if rewardTime.Unix() <= lastEraTimestamp {
				break OUT
			}

			rewardSum += int64(rd.Reward * 1e8)
			c.log.Trace("stakingReward", "add", rd.Reward, "height", rd.Height)
		}

		offset += 100
	}

	return rewardSum, nil
}

func (c *Connection) rewardApi(delegator string, limit, offset int64) string {
	return fmt.Sprintf("%s/v1/staking/chains/%s/delegators/%s/rewards?limit=%d&offset=%d", c.bcApiEndpoint, c.bscSideChainId, delegator, limit, offset)
}

func (c *Connection) GetStakingContract() *staking.Staking {
	return c.stakingContract
}

func (c *Connection) GetQueryClient() *ethereum.Client {
	return c.queryClient
}
