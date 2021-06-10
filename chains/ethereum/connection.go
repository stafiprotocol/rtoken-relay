// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"errors"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
)

type Connection struct {
	url     string
	symbol  core.RSymbol
	key     *secp256k1.Keypair
	keys    []*secp256k1.Keypair
	client  *ethereum.Client
	clients map[*secp256k1.Keypair]*ethereum.Client
	log     log15.Logger
	stop    <-chan int
	lastKey *signature.KeyringPair
}

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {
	log.Info("NewConnection", "name", cfg.Name, "KeystorePath", cfg.KeystorePath, "Endpoint", cfg.Endpoint)

	httpCfg := cfg.Opts["httpFlag"]
	httpFlag, ok := httpCfg.(bool)
	if !ok {
		return nil, errors.New("httpFlag config not ok")
	}

	var key *secp256k1.Keypair
	var clt *ethereum.Client
	keys := make([]*secp256k1.Keypair, 0)
	clients := make(map[*secp256k1.Keypair]*ethereum.Client)

	acSize := len(cfg.Accounts)
	for i := 0; i < acSize; i++ {
		kpI, err := keystore.KeypairFromAddress(cfg.Accounts[i], keystore.EthChain, cfg.KeystorePath, cfg.Insecure)
		if err != nil {
			return nil, err
		}
		kp, _ := kpI.(*secp256k1.Keypair)

		client := ethereum.NewClient(cfg.Endpoint, httpFlag, kp, log)
		if err != nil {
			return nil, err
		}
		keys = append(keys, kp)
		clients[kp] = client
		if i == 0 {
			key = kp
			clt = client
		}
	}

	return &Connection{
		url:     cfg.Endpoint,
		symbol:  cfg.Symbol,
		key:     key,
		keys:    keys,
		client:  clt,
		clients: clients,
		log:     log,
		stop:    stop,
	}, nil
}

func (c *Connection) LatestBlockNumber() (uint64, error) {
	blk, err := c.client.LatestBlock()
	if err != nil {
		return 0, err
	}

	return blk.Uint64(), nil
}

func (c *Connection) Address() string {
	return c.key.Address()
}

func (c *Connection) TransferVerify(r *submodel.BondRecord) (submodel.BondReason, error) {
	return submodel.Pass, nil
}
