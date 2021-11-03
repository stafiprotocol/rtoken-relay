// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package bnb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ChainSafe/log15"
	goeth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethCoreTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	bncClient "github.com/stafiprotocol/go-sdk/client"
	bncRpc "github.com/stafiprotocol/go-sdk/client/rpc"
	"github.com/stafiprotocol/go-sdk/client/websocket"
	bncCmnTypes "github.com/stafiprotocol/go-sdk/common/types"
	bnckeys "github.com/stafiprotocol/go-sdk/keys"
	bncTypes "github.com/stafiprotocol/go-sdk/types"
	"github.com/stafiprotocol/go-sdk/types/msgtype"
	"github.com/stafiprotocol/rtoken-relay/bindings/MultiSendCallOnly"
	"github.com/stafiprotocol/rtoken-relay/bindings/TokenHub"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/bnc"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
)

var (
	DefaultValue        = big.NewInt(0)
	ErrNonceTooLow      = errors.New("nonce too low")
	NoBondError         = errors.New("Staked found no bond")
	CheckBcBalanceError = errors.New("check bc balance error")

	TxRetryInterval  = time.Second * 2
	ErrTxUnderpriced = errors.New("replacement transaction underpriced")
	ZeroAddress      = common.HexToAddress("0x0000000000000000000000000000000000000000")
	sideChainId      = bncTypes.ChapelNet
	apiUrl           = "https://testnet-api.binance.org"
	baseUrl          = "testnet-dex.binance.org:443"
	dexUrl           = "https://testnet-dex.binance.org/api/v1"
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
	url               string
	symbol            core.RSymbol
	bscClient         *ethereum.Client
	keyManager        bnckeys.KeyManager
	bcKeys            map[common.Address]bnckeys.KeyManager
	bscClients        map[common.Address]*ethereum.Client
	bcRpcClient       *bncRpc.HTTP
	rpcEndpoint       string
	tokenHubContract  common.Address
	multisendContract common.Address
	bcBlockHeight     int64
	eraBlock          uint64
	log               log15.Logger
	stop              <-chan int

	tokenHub *TokenHub.TokenHub
}

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {
	log.Info("NewClient", "name", cfg.Name, "KeystorePath", cfg.KeystorePath, "Endpoint", cfg.Endpoint)

	rpcEndpointCfg := cfg.Opts["rpcEndpoint"]
	rpcEndpoint, ok := rpcEndpointCfg.(string)
	if !ok {
		return nil, errors.New("rpcEndpoint not exist")
	}

	if strings.HasPrefix(cfg.Accounts[0], "tbnb") {
		bncCmnTypes.Network = bncCmnTypes.TestNetwork
		log.Info("bnc network is TestNetwork")
	} else if strings.HasPrefix(cfg.Accounts[0], "bnb") {
		sideChainId = "bsc"
		apiUrl = "https://api.binance.org"
		baseUrl = "dex.binance.org:443"
		dexUrl = "https://dex.binance.org/api/v1"
		log.Info("bnc network is ProdNetwork")
	} else {
		return nil, fmt.Errorf("unknown bnc network")
	}

	rpcClient := bncRpc.NewRPCClient(rpcEndpoint, bncCmnTypes.Network)
	if !rpcClient.IsActive() {
		panic("rpcClient is not active")
	}

	var bscClient *ethereum.Client
	var keyManager bnckeys.KeyManager
	bcKeys := make(map[common.Address]bnckeys.KeyManager)
	bscClients := make(map[common.Address]*ethereum.Client)
	acSize := len(cfg.Accounts)
	if acSize == 0 || acSize%2 != 0 {
		return nil, fmt.Errorf("account size not even")
	}

	for i := 0; i < acSize; i += 2 {
		file := fmt.Sprintf("%s/%s.key", cfg.KeystorePath, cfg.Accounts[i])
		pswd := keystore.GetPassword(fmt.Sprintf("Enter password for key %s:", file))
		km, err := bnckeys.NewKeyStoreKeyManager(file, string(pswd))
		if err != nil {
			return nil, err
		}

		kpI, err := keystore.KeypairFromAddress(cfg.Accounts[i+1], keystore.EthChain, cfg.KeystorePath, cfg.Insecure)
		if err != nil {
			return nil, err
		}
		kp, _ := kpI.(*secp256k1.Keypair)

		conn := ethereum.NewClient(cfg.Endpoint, kp, log, big.NewInt(0), big.NewInt(0))
		if err := conn.Connect(); err != nil {
			return nil, err
		}

		address := kp.CommonAddress()
		bcKeys[address] = km
		bscClients[address] = conn
		if i == 0 {
			bscClient = conn
			keyManager = km
		}
	}

	hub, err := initTokenhub(cfg.Opts["TokenHubContract"])
	if err != nil {
		return nil, err
	}

	multiSend, err := initMultiSend(cfg.Opts["MultiSendContract"])
	if err != nil {
		return nil, err
	}

	return &Connection{
		url:               cfg.Endpoint,
		symbol:            cfg.Symbol,
		bscClient:         bscClient,
		keyManager:        keyManager,
		bcKeys:            bcKeys,
		bscClients:        bscClients,
		bcRpcClient:       rpcClient,
		rpcEndpoint:       rpcEndpoint,
		tokenHubContract:  hub,
		multisendContract: multiSend,
		log:               log,
		stop:              stop,
	}, nil
}

func (c *Connection) ReConnect() error {
	return c.bscClient.Connect()
}

// LatestBlock returns the latest block from the current chain
func (c *Connection) LatestBlock() (int64, error) {
	bcClient, err := bncClient.NewDexClient(baseUrl, bncCmnTypes.Network, c.keyManager)
	if err != nil {
		return 0, err
	}

	quit := make(chan struct{})
	ch, err := bcClient.WsGet("$all@blockheight", func(bz []byte) (interface{}, error) {
		var event websocket.BlockHeightEvent
		err := json.Unmarshal(bz, &event)
		return event.BlockHeight, err
	}, quit)

	if err != nil {
		c.log.Error("LatestBlock error", "error", err)
		return 0, err
	}

	timer := time.NewTimer(time.Minute)
	defer timer.Stop()

	select {
	case <-timer.C:
		return 0, errors.New("LatestBlock timeout")
	case h := <-ch:
		blk, ok := h.(int64)
		if !ok {
			return 0, errors.New("LatestBlock height type not int64")
		}
		return blk, nil
	}
}

func (c *Connection) TransferVerify(r *submodel.BondRecord) (submodel.BondReason, error) {
	return c.bscClient.BnbTransferVerify(r)
}

func (c *Connection) IsPoolKeyExist(pool common.Address) bool {
	return c.bcKeys[pool] != nil && c.bscClients[pool] != nil
}

func (c *Connection) TransferFromBscToBc(pool common.Address, amount int64) (int64, error) {
	c.log.Info("TransferFromBscToBc", "pool", pool, "amount", amount)

	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return 0, fmt.Errorf("TransferFromBscToBc no bc key found: %s", pool.Hex())
	}

	bscClient := c.bscClients[pool]
	if bscClient == nil {
		return 0, fmt.Errorf("TransferFromBscToBc found no bsc client : %s", pool.Hex())
	}

	bscBalance, err := c.bscClient.Client().BalanceAt(context.Background(), pool, nil)
	if err != nil {
		return 0, err
	}

	hub, err := TokenHub.NewTokenHub(c.tokenHubContract, bscClient.Client())
	if err != nil {
		return 0, err
	}

	fee, err := hub.RelayFee(nil)
	if err != nil {
		return 0, err
	}

	amt := big.NewInt(0).Mul(big.NewInt(1e10), big.NewInt(amount))
	leastBal := big.NewInt(0).Add(fee, amt)
	if bscBalance.Cmp(leastBal) < 0 {
		return 0, fmt.Errorf("bscBalance: %v too small to transfer, total: %v", bscBalance, leastBal)
	}

	curBal, err := c.bcBalance(bcKey)
	if err != nil {
		return 0, err
	}
	futureBal := curBal + amount
	c.log.Info("TransferFromBscToBc bc balance", "curBal", curBal, "futureBal", futureBal)

	receiver := common.HexToAddress(hexutil.Encode(bcKey.GetAddr()))

	expireTime := time.Now().Add(time.Hour).Unix()

	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-c.stop:
			return 0, errors.New("TransferFromBscToBc stopped")
		default:
			err = bscClient.LockAndUpdateOpts(big.NewInt(0), leastBal)
			if err != nil {
				return 0, err
			}

			tx, err := hub.TransferOut(bscClient.Opts(), ZeroAddress, receiver, amt, uint64(expireTime))
			bscClient.UnlockOpts()

			if err == nil {
				c.log.Info("TransferFromBscToBc result", "tx", tx.Hash(), "gasPrice", tx.GasPrice())

				timer := time.NewTimer(5 * time.Minute)
				defer timer.Stop()

				for {
					select {
					case <-timer.C:
						return 0, errors.New("TransferFromBscToBc transaction status timeout")
					default:
						receipt, err := bscClient.TransactionReceipt(tx.Hash())
						if err != nil {
							if err.Error() == goeth.NotFound.Error() {
								time.Sleep(2 * time.Second)
								continue
							}
							return 0, fmt.Errorf("TransferFromBscToBc TransactionReceipt error: %s", err)
						}

						if receipt.Status == ethCoreTypes.ReceiptStatusSuccessful {
							return futureBal, nil
						} else {
							return 0, errors.New("TransferFromBscToBc TransactionReceipt status fail")
						}
					}
				}
			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				c.log.Debug("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				c.log.Warn("TransferFromBscToBc failed", "err", err)
				time.Sleep(TxRetryInterval)
			}
		}
	}
	return 0, fmt.Errorf("TransferFromBscToBc failed")
}

func (c *Connection) CheckBcBalance(pool common.Address, future int64) error {
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return fmt.Errorf("CheckBcBalance no bc key found: %s", pool.Hex())
	}

	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return errors.New("CheckBcBalance timeout")
		default:
			realBal, err := c.bcBalance(bcKey)
			if err != nil {
				return err
			}

			c.log.Info("CheckBcBalance realBal", "realBal", realBal, "future", future)
			if realBal >= future {
				return nil
			}

			time.Sleep(30 * time.Second)
		}
	}
}

func (c *Connection) Unbondable(pool common.Address, validator bncCmnTypes.ValAddress) (bool, error) {
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return false, fmt.Errorf("Unbondable no bc key found: %s", pool.Hex())
	}

	unbonds, err := c.bcRpcClient.QuerySideChainUnbondingDelegations(sideChainId, bcKey.GetAddr())
	if err != nil {
		return false, err
	}

	for _, ub := range unbonds {
		if bytes.Equal(ub.ValidatorAddr, validator) {
			return false, nil
		}
	}

	return true, nil
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

func (c *Connection) ExecuteBond(pool common.Address, validator bncCmnTypes.ValAddress, amount int64) error {
	c.log.Info("ExecuteBond", "pool", pool, "validator", validator.String(), "amount", amount)
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return fmt.Errorf("ExecuteBond no bc key found: %s", pool.Hex())
	}

	enough, err := c.isBalanceEnough(bcKey, BcBondAction, amount)
	if err != nil {
		return err
	}

	if !enough {
		return fmt.Errorf("ExecuteBond free not enough")
	}

	c.bcRpcClient.SetKeyManager(bcKey)
	coin := bncCmnTypes.Coin{Denom: CoinSymbol, Amount: amount}

	res, err := c.bcRpcClient.SideChainDelegate(sideChainId, validator, coin, bncRpc.Sync)
	if err != nil {
		return err
	}

	c.log.Info("ExecuteBond", "txHash", res.Hash)

	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return errors.New("ExecuteBond txhash result timeout")
		default:
			err = c.txHashResult(res.Hash.String())
			if err == nil {
				return nil
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func (c *Connection) ExecuteUnbond(pool common.Address, validator bncCmnTypes.ValAddress, amount int64) error {
	c.log.Info("ExecuteUnbond", "pool", pool, "validator", validator.String(), "amount", amount)
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return fmt.Errorf("ExecuteUnbond no bc key found: %s", pool.Hex())
	}

	enough, err := c.isBalanceEnough(bcKey, BcUnbondAction, amount)
	if err != nil {
		return err
	}

	if !enough {
		return fmt.Errorf("ExecuteUnbond free not enough")
	}

	c.bcRpcClient.SetKeyManager(bcKey)
	coin := bncCmnTypes.Coin{Denom: CoinSymbol, Amount: amount}

	res, err := c.bcRpcClient.SideChainUnbond(sideChainId, validator, coin, bncRpc.Sync)
	if err != nil {
		return err
	}

	c.log.Info("ExecuteUnbond", "txHash", res.Hash)

	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return errors.New("ExecuteUnbond txhash result timeout")
		default:
			err = c.txHashResult(res.Hash.String())
			if err == nil {
				return nil
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func (c *Connection) Reward(pool common.Address, curHeight, lastHeight int64) (int64, error) {
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return 0, fmt.Errorf("Reward no bc key found: %s", pool.Hex())
	}
	delegator := bcKey.GetAddr().String()

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

	return 0, fmt.Errorf("Reward failed")
}

func (c *Connection) Staked(pool common.Address, validator bncCmnTypes.ValAddress) (int64, error) {
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return 0, fmt.Errorf("Staked no bc key found: %s", pool.Hex())
	}

	bonds, err := c.bcRpcClient.QuerySideChainDelegations(sideChainId, bcKey.GetAddr())
	if err != nil {
		return 0, err
	}

	for _, ub := range bonds {
		if bytes.Equal(ub.ValidatorAddr, validator) {
			return ub.Balance.Amount, nil
		}
	}

	return 0, NoBondError
}

func (c *Connection) CheckTransfer(pool common.Address, amount int64) error {
	c.log.Info("Transferable", "pool", pool, "amount", amount)
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return fmt.Errorf("CheckTransfer no bc key found: %s", pool.Hex())
	}

	enough, err := c.isBalanceEnough(bcKey, BcSwapAction, amount)
	if err != nil {
		return err
	}

	if !enough {
		return fmt.Errorf("CheckTransfer free not enough")
	}

	return nil
}

func (c *Connection) TransferFromBcToBsc(pool common.Address, amount int64) error {
	c.log.Info("TransferFromBcToBsc", "pool", pool, "amount", amount)
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return fmt.Errorf("TransferFromBcToBsc no bc key found: %s", pool.Hex())
	}

	bscBalance, err := c.bscClient.Client().BalanceAt(context.Background(), pool, nil)
	if err != nil {
		return err
	}

	bscAmount := big.NewInt(0).Mul(big.NewInt(1e10), big.NewInt(amount))
	newBalance := bscAmount.Add(bscAmount, bscBalance)

	c.bcRpcClient.SetKeyManager(bcKey)
	coin := bncCmnTypes.Coin{Denom: CoinSymbol, Amount: amount}
	expireTime := time.Now().Add(time.Hour).Unix()

	tx, err := c.bcRpcClient.TransferOut(msgtype.SmartChainAddress(pool), coin, expireTime, bncRpc.Sync)
	if err != nil {
		return err
	}

	c.log.Info("TransferFromBcToBsc", "txHash", tx.Hash)

	timer := time.NewTimer(10 * time.Minute)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return errors.New("ExecuteBond txhash result timeout")
		default:
			err = c.txHashResult(tx.Hash.String())
			if err == nil {
				bscBalance, err := c.bscClient.Client().BalanceAt(context.Background(), pool, nil)
				if err != nil {
					return err
				}

				if bscBalance.Cmp(newBalance) >= 0 {
					return nil
				}
			}

			time.Sleep(10 * time.Second)
		}
	}
}

func (c *Connection) BatchTransfer(pool common.Address, receives []*submodel.Receive, total *big.Int) (common.Hash, error) {
	bscClient := c.bscClients[pool]
	if bscClient == nil {
		return common.Hash{}, fmt.Errorf("BatchTransfer no bsc client found: %s", pool.Hex())
	}

	bscBalance, err := bscClient.Client().BalanceAt(context.Background(), pool, nil)
	if err != nil {
		return common.Hash{}, err
	}

	if bscBalance.Cmp(total) <= 0 {
		return common.Hash{}, fmt.Errorf("bscBalance: %v too small to transfer, total: %v", bscBalance, total)
	}

	sender, err := MultiSendCallOnly.NewMultiSendCallOnly(c.multisendContract, bscClient.Client())
	if err != nil {
		return common.Hash{}, err
	}

	bts := make(ethmodel.BatchTransactions, 0)
	totalGas := big.NewInt(0)
	for _, rec := range receives {
		addr := common.BytesToAddress(rec.Recipient)
		val := big.Int(rec.Value)
		value := big.NewInt(0).Mul(&val, big.NewInt(1e10))

		bt := &ethmodel.BatchTransaction{
			Operation:  uint8(ethmodel.Call),
			To:         addr,
			Value:      value,
			DataLength: big.NewInt(0),
			Data:       nil,
		}
		totalGas.Add(totalGas, big.NewInt(1e5))
		bts = append(bts, bt)
	}

	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-c.stop:
			return common.Hash{}, errors.New("BatchTransfer stopped")
		default:
			err = bscClient.LockAndUpdateOpts(totalGas, total)
			if err != nil {
				return common.Hash{}, err
			}
			tx, err := sender.MultiSend(bscClient.Opts(), bts.Encode())
			bscClient.UnlockOpts()

			if err == nil {
				c.log.Info("BatchTransfer result", "tx", tx.Hash(), "gasPrice", tx.GasPrice())

				timer := time.NewTimer(5 * time.Minute)
				defer timer.Stop()

				for {
					select {
					case <-timer.C:
						return common.Hash{}, errors.New("BatchTransfer transaction status timeout")
					default:
						receipt, err := bscClient.TransactionReceipt(tx.Hash())
						if err != nil {
							if err.Error() == goeth.NotFound.Error() {
								time.Sleep(2 * time.Second)
								continue
							}
							return common.Hash{}, fmt.Errorf("BatchTransfer TransactionReceipt error: %s", err)
						}

						if receipt.Status == ethCoreTypes.ReceiptStatusSuccessful {
							return tx.Hash(), nil
						} else {
							return common.Hash{}, errors.New("BatchTransfer TransactionReceipt status fail")
						}
					}
				}
			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				c.log.Debug("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				c.log.Warn("BatchTransfer failed", "err", err)
				time.Sleep(TxRetryInterval)
			}
		}
	}
	return common.Hash{}, fmt.Errorf("BatchTransfer failed")
}

func (c *Connection) FreeBalance(pool common.Address) (int64, error) {
	c.log.Info("FreeBalance", "pool", pool)
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return 0, fmt.Errorf("FreeBalance no bc key found: %s", pool.Hex())
	}

	bal, err := c.bcBalance(bcKey)
	if err != nil {
		return 0, err
	}

	return bal, nil
}

func (c *Connection) isBalanceEnough(key bnckeys.KeyManager, action BcActionType, amount int64) (bool, error) {
	free, err := c.bcBalance(key)
	if err != nil {
		return false, err
	}

	switch action {
	case BcBondAction:
		return amount+bondFee < free, nil
	case BcUnbondAction:
		return unbondFee < free, nil
	case BcSwapAction:
		return amount+transferFee < free, nil
	default:
		return false, fmt.Errorf("action not supported")
	}
}

func (c *Connection) bcBalance(key bnckeys.KeyManager) (int64, error) {
	addr := key.GetAddr()
	bal, err := c.bcRpcClient.GetBalance(addr, CoinSymbol)
	if err != nil {
		return 0, err
	}
	c.log.Info("current Balance", "bal", bal.Free.ToInt64())
	return bal.Free.ToInt64(), nil
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
	return fmt.Sprintf("%s/v1/staking/chains/%s/delegators/%s/rewards?limit=%d&offset=%d", apiUrl, sideChainId, delegator, limit, offset)
}

func (c *Connection) txHashResult(txHash string) error {
	url := fmt.Sprintf("%s/tx/%s", dexUrl, txHash)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	thr := new(bnc.TxHashResult)
	if err := json.Unmarshal(body, thr); err != nil {
		return err
	}

	c.log.Info("txHashResult", "txHash", txHash, "result", thr)
	return nil
}

func (c *Connection) SetEraBlock(eraBlock uint64) {
	c.eraBlock = eraBlock
}

func (c *Connection) EraBlock() uint64 {
	return c.eraBlock
}

func (c *Connection) Close() {
	c.bscClient.Close()
}

func initTokenhub(tokenHubCfg interface{}) (common.Address, error) {
	tokenHubAddr, ok := tokenHubCfg.(string)
	if !ok {
		return ZeroAddress, errors.New("TokenHubContract not ok")
	}

	return common.HexToAddress(tokenHubAddr), nil
}

func initMultiSend(multiSendCfg interface{}) (common.Address, error) {
	multiSendAddr, ok := multiSendCfg.(string)
	if !ok {
		return ZeroAddress, errors.New("MultiSendContract not ok")
	}

	return common.HexToAddress(multiSendAddr), nil
}
