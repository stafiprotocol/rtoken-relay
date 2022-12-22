// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package bnb

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	multisig_onchain "github.com/stafiprotocol/rtoken-relay/bindings/MultisigOnchain"
	stake_portal "github.com/stafiprotocol/rtoken-relay/bindings/StakeNativePortal"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
)

var (
	DefaultValue   = big.NewInt(0)
	ErrNonceTooLow = errors.New("nonce too low")
	NoBondError    = errors.New("Staked found no bond")

	TxRetryInterval  = time.Second * 2
	ErrTxUnderpriced = errors.New("replacement transaction underpriced")
	ZeroAddress      = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

const (
	TxRetryLimit = 10

	bondFee     = int64(20000)
	unbondFee   = int64(40000)
	transferFee = int64(7500)

	CoinSymbol = "BNB"
)

type BcActionType int

const (
	BcBondAction   = BcActionType(0)
	BcUnbondAction = BcActionType(1)
	BcSwapAction   = BcActionType(2)
)

type Connection struct {
	url         string
	symbol      core.RSymbol
	bscClient   *ethereum.Client                    // for query
	bscClients  map[common.Address]*ethereum.Client // pool -> clients
	rpcEndpoint string
	log         core.Logger
	stop        <-chan int

	multisigOnchains    map[common.Address]*multisig_onchain.MultisigOnchain // pool -> multisigOnchain
	stakePortalContract *stake_portal.StakeNativePortal
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
	bscClients := make(map[common.Address]*ethereum.Client)
	multisigs := make(map[common.Address]*multisig_onchain.MultisigOnchain)
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

		conn := ethereum.NewClient(cfg.Endpoint, kp, log, big.NewInt(0), big.NewInt(0))
		if err := conn.Connect(); err != nil {
			return nil, err
		}
		bscClients[multisigContracts[i]] = conn

		multisigs[multisigContracts[i]], err = multisig_onchain.NewMultisigOnchain(multisigContracts[i], ethClient)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			bscClient = conn
		}
	}

	stakePortal, err := initStakePortal(cfg.Opts["StakePortalContract"], bscClient.Client())
	if err != nil {
		return nil, err
	}

	return &Connection{
		url:                 cfg.Endpoint,
		symbol:              cfg.Symbol,
		bscClient:           bscClient,
		rpcEndpoint:         rpcEndpoint,
		log:                 log,
		stop:                stop,
		multisigOnchains:    multisigs,
		stakePortalContract: stakePortal,
	}, nil
}

func (c *Connection) ReConnect() error {
	return c.bscClient.Connect()
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

func (c *Connection) TransferVerify(r *submodel.BondRecord) (result submodel.BondReason, err error) {
	for i := 0; i < 5; i++ {
		result, err = c.bscClient.BnbTransferVerify(r)
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

func (c *Connection) IsPoolKeyExist(pool common.Address) bool {
	return c.multisigOnchains[pool] != nil
}

func (c *Connection) BondOrUnbondCall(bond, unbond, leastBond int64) (submodel.BondAction, int64) {
	c.log.Info("BondOrUnbondCall", "bond", bond, "unbond", unbond, "leastBond", leastBond)

	if bond < unbond {
		diff := unbond - bond
		if diff < leastBond {
			c.log.Info("bond is smaller than unbond while diff is smaller than leastBond, InterDeduct")
			return submodel.InterDeduct, 0
		}
		return submodel.UnbondOnly, diff
	} else if bond > unbond {
		diff := bond - unbond
		if diff < leastBond {
			c.log.Info("unbond is smaller than bond while diff is smaller than leastBond, InterDeduct")
			return submodel.InterDeduct, 0
		}
		return submodel.BondOnly, diff
	} else {
		c.log.Info("bond is equal to unbond, NoCall")
		if bond == 0 {
			return submodel.EitherBondUnbond, 0
		}
		return submodel.BothBondUnbond, 0
	}
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
	for _, c := range c.bscClients {
		c.Close()
	}
}

func (c *Connection) LatestBlockTimestamp() (uint64, error) {
	blkTime, err := c.bscClient.LatestBlockTimestamp()
	if err != nil {
		return 0, err
	}

	return blkTime, nil
}

func (c *Connection) LatestBlock() (uint64, error) {
	blk, err := c.bscClient.LatestBlock()
	if err != nil {
		return 0, err
	}

	return blk.Uint64(), nil
}
