// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package bnb

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"math/big"
	"strings"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	bncCmnTypes "github.com/stafiprotocol/go-sdk/common/types"
	bnckeys "github.com/stafiprotocol/go-sdk/keys"
	bncTypes "github.com/stafiprotocol/go-sdk/types"
	"github.com/stafiprotocol/rtoken-relay/bindings/TokenHub"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
)

var (
	DefaultValue     = big.NewInt(0)
	ErrNonceTooLow   = errors.New("nonce too low")
	ErrTxUnderpriced = errors.New("replacement transaction underpriced")
	ZeroAddress      = common.HexToAddress("0x0000000000000000000000000000000000000000")
	sideChainId      = bncTypes.ChapelNet
)

const (
	TxRetryLimit = 10
)

type Connection struct {
	url              string
	symbol           core.RSymbol
	conn             *ethereum.Client
	bscKeys          []*secp256k1.Keypair
	bcKeys           map[*secp256k1.Keypair]bnckeys.KeyManager
	conns            map[*secp256k1.Keypair]*ethereum.Client
	tokenHubContract common.Address
	log              log15.Logger
	stop             <-chan int

	tokenHub *TokenHub.TokenHub
}

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {
	log.Info("NewClient", "name", cfg.Name, "KeystorePath", cfg.KeystorePath, "Endpoint", cfg.Endpoint)

	var key *secp256k1.Keypair
	bscKeys := make([]*secp256k1.Keypair, 0)
	bcKeys := make(map[*secp256k1.Keypair]bnckeys.KeyManager)
	conns := make(map[*secp256k1.Keypair]*ethereum.Client)
	acSize := len(cfg.Accounts)
	if acSize == 0 || acSize%2 != 0 {
		return nil, fmt.Errorf("account size not even")
	}

	if strings.HasPrefix(cfg.Accounts[0], "tbnb") {
		bncCmnTypes.Network = bncCmnTypes.TestNetwork
		log.Info("bnc networ is TestNetwork")
	} else if strings.HasPrefix(cfg.Accounts[0], "bnb") {
		sideChainId = "bsc"
		log.Info("bnc network is ProdNetwork")
	} else {
		return nil, fmt.Errorf("unknown bnc network")
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

		bscKeys = append(bscKeys, kp)
		bcKeys[kp] = km
		conns[kp] = conn
	}

	conn := conns[bscKeys[0]]
	hub, err := initTokenhub(cfg.Opts["TokenHubContract"])
	if err != nil {
		return nil, err
	}

	return &Connection{
		url:              cfg.Endpoint,
		symbol:           cfg.Symbol,
		conn:             conn,
		bscKeys:          bscKeys,
		bcKeys:           bcKeys,
		conns:            conns,
		tokenHubContract: hub,
		log:              log,
		stop:             stop,
	}, nil
}

func (c *Connection) ReConnect() error {
	return c.conn.Connect()
}

// LatestBlock returns the latest block from the current chain
func (c *Connection) LatestBlock() (uint64, error) {
	blk, err := c.conn.LatestBlock()
	if err != nil {
		return 0, err
	}
	return blk.Uint64(), nil
}

func (c *Connection) Address() string {
	return c.conn.Keypair().Address()
}

func (c *Connection) TransferVerify(r *submodel.BondRecord) (submodel.BondReason, error) {
	return c.conn.BnbTransferVerify(r)
}

func (c *Connection) BscTransferToBc(r *submodel.BondRecord) error {
	bscKey := c.FoundBscKey(r.Pool)
	if bscKey == nil {
		return fmt.Errorf("found no bscKey")
	}

	conn := c.conns[bscKey]
	hub, err := TokenHub.NewTokenHub(c.tokenHubContract, conn.Client())
	if err != nil {
		return err
	}

	fee, err := hub.RelayFee(nil)
	if err != nil {
		return err
	}

	bcKey := c.bcKeys[bscKey]
	receiver := common.HexToAddress(hexutil.Encode(bcKey.GetAddr()))
	value := big.NewInt(0).Add(r.Amount.Int, fee)

	err = conn.LockAndUpdateOpts(big.NewInt(0), value)
	if err != nil {
		return err
	}

	tx, err := hub.TransferOut(conn.Opts(), ZeroAddress, receiver, r.Amount.Int, 0x17b1f307761)
	conn.UnlockOpts()

	if err != nil {
		return err
	}
	c.log.Info("BscTransferToBc txHash", tx.Hash())
	return nil
}

func (c *Connection) FoundBscKey(pool []byte) *secp256k1.Keypair {
	for _, key := range c.bscKeys {
		if !bytes.Equal(key.CommonAddress().Bytes(), pool) {
			continue
		}

		return key
	}

	return nil
}

func (c *Connection) Close() {
	c.conn.Close()
}

func initTokenhub(TokenHubCfg interface{}) (common.Address, error) {
	tokenHubAddr, ok := TokenHubCfg.(string)
	if !ok {
		return ZeroAddress, errors.New("TokenHubContract not ok")
	}

	//hub, err := TokenHub.NewTokenHub(common.HexToAddress(tokenHubAddr), conn)
	//if err != nil {
	//	return nil, err
	//}

	return common.HexToAddress(tokenHubAddr), nil
}
