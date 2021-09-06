// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package bnb

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	DefaultValue     = big.NewInt(0)
	ErrNonceTooLow   = errors.New("nonce too low")
	TxRetryInterval  = time.Second * 2
	ErrTxUnderpriced = errors.New("replacement transaction underpriced")
	ZeroAddress      = common.HexToAddress("0x0000000000000000000000000000000000000000")
	sideChainId      = bncTypes.ChapelNet
	apiUrl           = "https://testnet-api.binance.org"
)

const (
	TxRetryLimit = 10

	BondFee     = int64(20000)
	UnbondFee   = int64(40000)
	TransferFee = int64(7500)
)

type Connection struct {
	url               string
	symbol            core.RSymbol
	bscClient         *ethereum.Client
	bcClient          bncClient.DexClient
	bcKeys            map[common.Address]bnckeys.KeyManager
	bscClients        map[common.Address]*ethereum.Client
	bcRpcClient       *bncRpc.HTTP
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
		log.Info("bnc networ is TestNetwork")
	} else if strings.HasPrefix(cfg.Accounts[0], "bnb") {
		sideChainId = "bsc"
		apiUrl = "https://api.binance.org"
		log.Info("bnc network is ProdNetwork")
	} else {
		return nil, fmt.Errorf("unknown bnc network")
	}

	rpcClient := bncRpc.NewRPCClient(rpcEndpoint, bncCmnTypes.Network)

	var key *secp256k1.Keypair
	var bscClient *ethereum.Client
	var bcClient bncClient.DexClient
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

		conn := ethereum.NewClient(cfg.Endpoint, key, log, big.NewInt(0), big.NewInt(0))
		if err := conn.Connect(); err != nil {
			return nil, err
		}

		address := kp.CommonAddress()
		bcKeys[address] = km
		bscClients[address] = conn
		if i == 0 {
			bscClient = conn

			bcClient, err = bncClient.NewDexClient("testnet-dex.binance.org:443", bncCmnTypes.Network, km)
			if err != nil {
				return nil, err
			}
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
		bcClient:          bcClient,
		bcKeys:            bcKeys,
		bscClients:        bscClients,
		bcRpcClient:       rpcClient,
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
func (c *Connection) LatestBlock() (latest int64, err error) {
	quit := make(chan struct{})
	defer close(quit)
	errCount := 0

	timer := time.NewTimer(time.Minute)
	defer timer.Stop()

	onReceive := func(event *websocket.BlockHeightEvent) {
		latest = event.BlockHeight
		quit <- struct{}{}
	}

	onError := func(inErr error) {
		errCount += 1
		if errCount > 10 {
			err = inErr
			quit <- struct{}{}
		}
	}
	err = c.bcClient.SubscribeBlockHeightEvent(quit, onReceive, onError, nil)
	if err != nil {
		return
	}

	select {
	case <-timer.C:
		err = errors.New("LatestBlock timeout")
		quit <- struct{}{}
		return
	}
}

func (c *Connection) TransferVerify(r *submodel.BondRecord) (submodel.BondReason, error) {
	return c.bscClient.BnbTransferVerify(r)
}

func (c *Connection) BckeyByPool(pool common.Address) bnckeys.KeyManager {
	return c.bcKeys[pool]
}

func (c *Connection) BscTransferToBc(r *submodel.BondRecord) error {
	addr := common.BytesToAddress(r.Pool)
	bscClient := c.bscClients[addr]
	if bscClient == nil {
		return fmt.Errorf("BscTransferToBc no bsc client found: %s", addr.Hex())
	}

	hub, err := TokenHub.NewTokenHub(c.tokenHubContract, bscClient.Client())
	if err != nil {
		return err
	}

	fee, err := hub.RelayFee(nil)
	if err != nil {
		return err
	}

	bcKey := c.bcKeys[addr]
	if bcKey == nil {
		return fmt.Errorf("BscTransferToBc no bc key found: %s", addr.Hex())
	}

	receiver := common.HexToAddress(hexutil.Encode(bcKey.GetAddr()))
	value := big.NewInt(0).Add(r.Amount.Int, fee)
	expireTime := time.Now().Add(time.Hour).Unix()

	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-c.stop:
			return errors.New("BscTransferToBc stopped")
		default:
			err = bscClient.LockAndUpdateOpts(big.NewInt(0), value)
			if err != nil {
				return err
			}

			tx, err := hub.TransferOut(bscClient.Opts(), ZeroAddress, receiver, r.Amount.Int, uint64(expireTime))
			bscClient.UnlockOpts()

			if err == nil {
				c.log.Info("BscTransferToBc result", "tx", tx.Hash(), "gasPrice", tx.GasPrice())
				return nil
			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				c.log.Debug("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				c.log.Warn("BscTransferToBc failed", "err", err)
				time.Sleep(TxRetryInterval)
			}
		}
	}
	return fmt.Errorf("BscTransferToBc failed")
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
	c.log.Info("BondOrUnbondCall", "bond", bond, "unbond", unbond)

	if bond < unbond {
		diff := unbond - bond
		if diff < leastBond {
			c.log.Info("bond is smaller than unbond while diff is smaller than leastBond, BondOnlyReport", "bond", bond, "unbond", unbond, "leastBond", leastBond)
			return submodel.EitherBondUnbond, 0
		}
		return submodel.UnbondOnly, diff
	} else if bond > unbond {
		diff := bond - unbond
		if diff < leastBond {
			c.log.Info("unbond is smaller than bond while diff is smaller than leastBond, BondOnlyReport", "bond", bond, "unbond", unbond, "leastBond", leastBond)
			return submodel.EitherBondUnbond, 0
		}
		return submodel.BondOnly, diff
	} else {
		c.log.Info("bond is equal to unbond, NoCall")
		return submodel.BothBondUnbond, 0
	}
}

func (c *Connection) ExecuteBond(pool common.Address, validator bncCmnTypes.ValAddress, amount int64) error {
	c.log.Info("ExecuteBond", "pool", pool, "validator", validator.String(), "amount", amount)
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return fmt.Errorf("ExecuteBond no bc key found: %s", pool.Hex())
	}

	enough, err := c.isBalanceEnough(bcKey, 0, amount)
	if err != nil {
		return err
	}

	if !enough {
		return fmt.Errorf("ExecuteBond free not enough")
	}

	c.bcRpcClient.SetKeyManager(bcKey)
	coin := bncCmnTypes.Coin{Denom: "BNB", Amount: amount}

	res, err := c.bcRpcClient.SideChainDelegate(sideChainId, validator, coin, bncRpc.Sync)
	if err != nil {
		return err
	}

	c.log.Info("ExecuteBond", "txHash", res.Hash)
	time.Sleep(time.Minute)
	return nil
}

func (c *Connection) ExecuteUnbond(pool common.Address, validator bncCmnTypes.ValAddress, amount int64) error {
	c.log.Info("ExecuteUnbond", "pool", pool, "validator", validator.String(), "amount", amount)
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return fmt.Errorf("ExecuteUnbond no bc key found: %s", pool.Hex())
	}

	enough, err := c.isBalanceEnough(bcKey, 1, amount)
	if err != nil {
		return err
	}

	if !enough {
		return fmt.Errorf("ExecuteUnbond free not enough")
	}

	c.bcRpcClient.SetKeyManager(bcKey)
	coin := bncCmnTypes.Coin{Denom: "BNB", Amount: amount}

	res, err := c.bcRpcClient.SideChainUnbond(sideChainId, validator, coin, bncRpc.Sync)
	if err != nil {
		return err
	}

	c.log.Info("ExecuteUnbond", "txHash", res.Hash)
	time.Sleep(time.Minute)
	return nil
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
			if i+1 == TxRetryLimit {
				return 0, err
			}
			continue
		}

		if height < lastHeight {
			return 0, nil
		}

		api := c.rewardApi(delegator, total, 0)
		sr, err := bnc.GetStakingReward(api)
		if err != nil {
			if i+1 == TxRetryLimit {
				return 0, err
			}
			continue
		}

		rewardSum := int64(0)
		for _, rd := range sr.RewardDetails {
			if rd.Height > curHeight {
				continue
			}

			if rd.Height < lastHeight {
				return rewardSum, nil
			}

			rewardSum += int64(rd.Reward * 1e8)
		}
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

	return 0, errors.New("Staked found no bond")
}

func (c *Connection) CheckTransfer(pool common.Address, amount int64) error {
	c.log.Info("Transferable", "pool", pool, "amount", amount)
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return fmt.Errorf("Transferable no bc key found: %s", pool.Hex())
	}

	enough, err := c.isBalanceEnough(bcKey, 2, amount)
	if err != nil {
		return err
	}

	if !enough {
		return fmt.Errorf("Transferable free not enough")
	}

	return nil
}

func (c *Connection) TransferFromBcToBsc(pool common.Address, amount int64) error {
	c.log.Info("TransferFromBcToBsc", "pool", pool, "amount", amount)
	bcKey := c.bcKeys[pool]
	if bcKey == nil {
		return fmt.Errorf("TransferFromBcToBsc no bc key found: %s", pool.Hex())
	}

	c.bcRpcClient.SetKeyManager(bcKey)
	coin := bncCmnTypes.Coin{Denom: "BNB", Amount: amount}
	expireTime := time.Now().Add(time.Hour).Unix()

	tx, err := c.bcRpcClient.TransferOut(msgtype.SmartChainAddress(pool), coin, expireTime, bncRpc.Sync)
	if err != nil {
		return err
	}

	c.log.Info("TransferFromBcToBsc", "txHash", tx.Hash)
	return nil
}

func (c *Connection) BatchTransfer(pool common.Address, receives []*submodel.Receive, amount int64) error {
	bscClient := c.bscClients[pool]
	if bscClient == nil {
		return fmt.Errorf("BatchTransfer no bsc client found: %s", pool.Hex())
	}

	sender, err := MultiSendCallOnly.NewMultiSendCallOnly(c.multisendContract, bscClient.Client())
	if err != nil {
		return err
	}

	bts := make(ethmodel.BatchTransactions, 0)
	totalGas := big.NewInt(0)
	for _, rec := range receives {
		addr := common.BytesToAddress(rec.Recipient)
		value := big.Int(rec.Value)

		bt := &ethmodel.BatchTransaction{
			Operation:  uint8(ethmodel.Call),
			To:         addr,
			Value:      &value,
			DataLength: big.NewInt(0),
			Data:       nil,
		}
		totalGas.Add(totalGas, big.NewInt(1e5))
		bts = append(bts, bt)
	}

	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-c.stop:
			return errors.New("BatchTransfer stopped")
		default:
			err = bscClient.LockAndUpdateOpts(totalGas, big.NewInt(amount))
			if err != nil {
				return err
			}
			tx, err := sender.MultiSend(bscClient.Opts(), bts.Encode())
			bscClient.UnlockOpts()

			if err == nil {
				c.log.Info("BatchTransfer result", "tx", tx.Hash(), "gasPrice", tx.GasPrice())
				return nil
			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				c.log.Debug("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				c.log.Warn("BatchTransfer failed", "err", err)
				time.Sleep(TxRetryInterval)
			}
		}
	}
	return fmt.Errorf("BatchTransfer failed")
}

func (c *Connection) isBalanceEnough(key bnckeys.KeyManager, action int, amount int64) (bool, error) {
	free, err := c.bcBanlance(key)
	if err != nil {
		return false, err
	}

	switch action {
	case 0:
		return amount+BondFee < free, nil
	case 1:
		return amount+UnbondFee < free, nil
	case 2:
		return amount+TransferFee < free, nil
	default:
		return false, fmt.Errorf("action not supported")
	}
}

func (c *Connection) bcBanlance(key bnckeys.KeyManager) (int64, error) {
	addr := key.GetAddr()
	bal, err := c.bcRpcClient.GetBalance(addr, "BNB")
	if err != nil {
		return 0, err
	}
	c.log.Info("current Balance", "bal", bal.Free.ToInt64())
	return bal.Free.ToInt64(), nil
}

func (c *Connection) totalAndLastHeight(delegator string) (int64, int64, error) {
	api := c.rewardApi(delegator, 1, 0)
	sr, err := bnc.GetStakingReward(api)
	if err != nil {
		return 0, 0, err
	}

	if len(sr.RewardDetails) == 0 {
		return 0, 0, nil
	}

	return sr.Total, sr.RewardDetails[0].Height, nil
}

func (c *Connection) rewardApi(delegator string, limit, offset int64) string {
	return fmt.Sprintf("%s/v1/staking/chains/%s/delegators/%s/rewards?limit=%d&offset=%d", apiUrl, sideChainId, delegator, limit, offset)
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
