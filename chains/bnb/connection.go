// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package bnb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	multisig_onchain "github.com/stafiprotocol/rtoken-relay/bindings/MultisigOnchain"
	stake_portal "github.com/stafiprotocol/rtoken-relay/bindings/StakeNativePortal"
	staking "github.com/stafiprotocol/rtoken-relay/bindings/Staking"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/bnc"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
)

var (
	TxRetryInterval  = time.Second * 2
	ErrTxUnderpriced = errors.New("replacement transaction underpriced")
	ZeroAddress      = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

const (
	TxRetryLimit = 10
	CoinSymbol   = "BNB"
)

type Connection struct {
	symbol      core.RSymbol
	queryClient *ethereum.Client // first pool's bsc client used for query
	apiEndpint  string
	sideChainId string
	log         core.Logger
	stop        <-chan int

	pools               map[common.Address]*Pool // pool address => pool
	stakePortalContract *stake_portal.StakeNativePortal
	stakingContract     *staking.Staking
}

type Pool struct {
	poolAddress     common.Address
	bscClient       *ethereum.Client
	multisigOnchain *multisig_onchain.MultisigOnchain
}

func NewConnection(cfg *core.ChainConfig, log core.Logger, stop <-chan int) (*Connection, error) {
	log.Info("NewClient", "name", cfg.Name, "KeystorePath", cfg.KeystorePath, "Endpoint", cfg.Endpoint)

	rpcEndpointCfg := cfg.Opts["rpcEndpoint"]
	rpcEndpoint, ok := rpcEndpointCfg.(string)
	if !ok {
		return nil, errors.New("rpcEndpoint not exist")
	}

	ethClient, err := ethclient.Dial(rpcEndpoint)
	if err != nil {
		return nil, err
	}

	multisigContractsCfg := []string{}

	bts, err := json.Marshal(cfg.Opts["MultiSigContracts"])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bts, &multisigContractsCfg)
	if err != nil {
		return nil, err
	}

	var bscClient *ethereum.Client
	var pools = make(map[common.Address]*Pool)
	acSize := len(cfg.Accounts)
	if acSize == 0 {
		return nil, fmt.Errorf("account empty")
	}
	if acSize != len(multisigContractsCfg) {
		return nil, fmt.Errorf("account size should equal multisig contracts size")
	}
	multisigContracts := make([]common.Address, 0)
	for _, addr := range multisigContractsCfg {
		multisigContracts = append(multisigContracts, common.HexToAddress(addr))
	}

	for i := 0; i < acSize; i++ {
		kpI, err := keystore.KeypairFromAddress(cfg.Accounts[i], keystore.EthChain, cfg.KeystorePath, cfg.Insecure)
		if err != nil {
			return nil, err
		}
		kp, _ := kpI.(*secp256k1.Keypair)

		client := ethereum.NewClient(cfg.Endpoint, kp, log, big.NewInt(0), big.NewInt(0))
		if err := client.Connect(); err != nil {
			return nil, err
		}

		multisigOnchain, err := multisig_onchain.NewMultisigOnchain(multisigContracts[i], ethClient)
		if err != nil {
			return nil, err
		}
		pool := Pool{
			poolAddress:     multisigContracts[i],
			bscClient:       client,
			multisigOnchain: multisigOnchain,
		}
		pools[multisigContracts[i]] = &pool

		if i == 0 {
			bscClient = client
		}

	}

	stakePortal, err := initStakePortal(cfg.Opts["StakePortalContract"], ethClient)
	if err != nil {
		return nil, err
	}

	staking, err := initStaking(ethClient)
	if err != nil {
		return nil, err
	}

	return &Connection{
		symbol:              cfg.Symbol,
		queryClient:         bscClient,
		log:                 log,
		stop:                stop,
		pools:               pools,
		stakePortalContract: stakePortal,
		stakingContract:     staking,
	}, nil
}

func initStakePortal(stakeManagerCfg interface{}, conn *ethclient.Client) (*stake_portal.StakeNativePortal, error) {
	stakePortalAddr, ok := stakeManagerCfg.(string)
	if !ok {
		return nil, errors.New("StakeManagerContract not ok")
	}
	stakePortal, err := stake_portal.NewStakeNativePortal(common.HexToAddress(stakePortalAddr), conn)
	if err != nil {
		return nil, err
	}

	return stakePortal, nil
}

func initStaking(conn *ethclient.Client) (*staking.Staking, error) {
	staking, err := staking.NewStaking(StakingContractAddr, conn)
	if err != nil {
		return nil, err
	}

	return staking, nil
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

func (c *Connection) GetPool(pool common.Address) (*Pool, bool) {
	p, exist := c.pools[pool]
	return p, exist
}

// use conn address as blockstore use address
func (c *Connection) BlockStoreUseAddress() string {
	return c.queryClient.Keypair().Address()
}

// func (c *Connection) BatchTransfer(pool common.Address, receives []*submodel.Receive, total *big.Int) (common.Hash, error) {
// 	bscClient := c.bscClients[pool]
// 	if bscClient == nil {
// 		return common.Hash{}, fmt.Errorf("BatchTransfer no bsc client found: %s", pool.Hex())
// 	}

// 	bscBalance, err := bscClient.Client().BalanceAt(context.Background(), pool, nil)
// 	if err != nil {
// 		return common.Hash{}, err
// 	}

// 	if bscBalance.Cmp(total) <= 0 {
// 		return common.Hash{}, fmt.Errorf("bscBalance: %v too small to transfer, total: %v", bscBalance, total)
// 	}

// 	sender, err := MultiSendCallOnly.NewMultiSendCallOnly(c.multisendContract, bscClient.Client())
// 	if err != nil {
// 		return common.Hash{}, err
// 	}

// 	bts := make(ethmodel.BatchTransactions, 0)
// 	totalGas := big.NewInt(0)
// 	for _, rec := range receives {
// 		addr := common.BytesToAddress(rec.Recipient)
// 		val := big.Int(rec.Value)
// 		value := big.NewInt(0).Mul(&val, big.NewInt(1e10))

// 		bt := &ethmodel.BatchTransaction{
// 			Operation:  uint8(ethmodel.Call),
// 			To:         addr,
// 			Value:      value,
// 			DataLength: big.NewInt(0),
// 			Data:       nil,
// 		}
// 		totalGas.Add(totalGas, big.NewInt(1e5))
// 		bts = append(bts, bt)
// 	}

// 	for i := 0; i < TxRetryLimit; i++ {
// 		select {
// 		case <-c.stop:
// 			return common.Hash{}, errors.New("BatchTransfer stopped")
// 		default:
// 			err = bscClient.LockAndUpdateOpts(totalGas, total)
// 			if err != nil {
// 				return common.Hash{}, err
// 			}
// 			tx, err := sender.MultiSend(bscClient.Opts(), bts.Encode())
// 			bscClient.UnlockOpts()

// 			if err == nil {
// 				c.log.Info("BatchTransfer result", "tx", tx.Hash(), "gasPrice", tx.GasPrice())

// 				timer := time.NewTimer(5 * time.Minute)
// 				defer timer.Stop()

// 				for {
// 					select {
// 					case <-timer.C:
// 						return common.Hash{}, errors.New("BatchTransfer transaction status timeout")
// 					default:
// 						receipt, err := bscClient.TransactionReceipt(tx.Hash())
// 						if err != nil {
// 							if err.Error() == goeth.NotFound.Error() {
// 								time.Sleep(2 * time.Second)
// 								continue
// 							}
// 							return common.Hash{}, fmt.Errorf("BatchTransfer TransactionReceipt error: %s", err)
// 						}

// 						if receipt.Status == ethCoreTypes.ReceiptStatusSuccessful {
// 							return tx.Hash(), nil
// 						} else {
// 							return common.Hash{}, errors.New("BatchTransfer TransactionReceipt status fail")
// 						}
// 					}
// 				}
// 			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
// 				c.log.Debug("Nonce too low, will retry")
// 				time.Sleep(TxRetryInterval)
// 			} else {
// 				c.log.Warn("BatchTransfer failed", "err", err)
// 				time.Sleep(TxRetryInterval)
// 			}
// 		}
// 	}
// 	return common.Hash{}, fmt.Errorf("BatchTransfer failed")
// }

func (c *Connection) Close() {
	for _, c := range c.pools {
		c.bscClient.Close()
	}
}

func (c *Connection) LatestBlockTimestamp() (uint64, error) {
	blkTime, err := c.queryClient.LatestBlockTimestamp()
	if err != nil {
		return 0, err
	}

	return blkTime, nil
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

func (c *Connection) GetHeightByEra(era uint32, eraSeconds, offset int64) (int64, error) {
	if int64(era) < offset {
		return 0, fmt.Errorf("era: %d is less than offset: %d", era, offset)
	}
	targetTimestamp := (int64(era) - offset) * eraSeconds
	return c.GetHeightByTimestamp(targetTimestamp)
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

	tmpTargetBlock := int64(blockNumber) - seconds/7
	if tmpTargetBlock <= 0 {
		tmpTargetBlock = 1
	}

	block, err := c.QueryBlock(tmpTargetBlock)
	if err != nil {
		return 0, err
	}

	// return after blocknumber
	var afterBlockNumber int64
	var preBlockNumber int64
	if int64(block.Time()) > targetTimestamp {
		afterBlockNumber = block.Number().Int64()
		for {
			c.log.Trace("afterBlock", "block", afterBlockNumber)
			if afterBlockNumber <= 2 {
				return 1, nil
			}
			block, err := c.QueryBlock(afterBlockNumber - 1)
			if err != nil {
				return 0, err
			}
			if int64(block.Time()) > targetTimestamp {
				afterBlockNumber = block.Number().Int64()
				continue
			}

			break
		}

	} else {
		preBlockNumber = block.Number().Int64()
		for {
			c.log.Trace("preBlock", "block", preBlockNumber)
			block, err := c.QueryBlock(preBlockNumber + 1)
			if err != nil {
				return 0, err
			}
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

func (c *Connection) RewardOnBc(pool common.Address, curHeight, lastHeight int64) (int64, error) {
	delegator := ""

	for i := 0; i < TxRetryLimit; i++ {
		total, height, err := c.totalAndLastHeight(delegator)
		if err != nil {
			c.log.Error("totalAndLastHeight error", "err", err)
			if i+1 == TxRetryLimit {
				return 0, err
			}
			continue
		}

		if height < lastHeight {
			return 0, nil
		}

		rewardSum, err := c.stakingReward(delegator, total, curHeight, lastHeight)
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

func (c *Connection) totalAndLastHeight(delegator string) (int64, int64, error) {
	api := c.rewardApi(delegator, 1, 0)
	c.log.Info("totalAndLastHeight rewardApi", "rewardApi", api)
	sr, err := bnc.GetStakingReward(api)
	if err != nil {
		return 0, 0, err
	}

	if len(sr.RewardDetails) == 0 {
		return 0, 0, nil
	}

	return sr.Total, sr.RewardDetails[0].Height, nil
}

func (c *Connection) stakingReward(delegator string, total, curHeight, lastHeight int64) (int64, error) {
	offset := int64(0)
	rewardSum := int64(0)

OUT:
	for i := total; i > 0; i -= 100 {
		api := c.rewardApi(delegator, 100, offset)

		sr, err := bnc.GetStakingReward(api)
		if err != nil {
			c.log.Info("stakingReward GetStakingReward error", "err", err)
			return 0, err
		}

		for _, rd := range sr.RewardDetails {
			if rd.Height > curHeight {
				continue
			}

			if rd.Height <= lastHeight {
				break OUT
			}

			rewardSum += int64(rd.Reward * 1e8)
		}

		offset += 100
	}

	return rewardSum, nil
}

func (c *Connection) rewardApi(delegator string, limit, offset int64) string {
	return fmt.Sprintf("%s/v1/staking/chains/%s/delegators/%s/rewards?limit=%d&offset=%d", c.apiEndpint, c.sideChainId, delegator, limit, offset)
}
