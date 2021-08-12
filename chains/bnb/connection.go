// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package bnb

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	bncRpc "github.com/stafiprotocol/go-sdk/client/rpc"
	bncCmnTypes "github.com/stafiprotocol/go-sdk/common/types"
	bnckeys "github.com/stafiprotocol/go-sdk/keys"
	bncTypes "github.com/stafiprotocol/go-sdk/types"
	"github.com/stafiprotocol/rtoken-relay/bindings/TokenHub"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
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
	bscClient        *ethereum.Client
	bcKeys           map[common.Address]bnckeys.KeyManager
	bscClients       map[common.Address]*ethereum.Client
	bcRpcClient      *bncRpc.HTTP
	tokenHubContract common.Address
	log              log15.Logger
	stop             <-chan int

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
		log.Info("bnc network is ProdNetwork")
	} else {
		return nil, fmt.Errorf("unknown bnc network")
	}

	rpcClient := bncRpc.NewRPCClient(rpcEndpoint, bncCmnTypes.Network)

	var key *secp256k1.Keypair
	//bscKeys := make([]common.Address, 0)
	var bscClient *ethereum.Client
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
		//bscKeys = append(bscKeys, address)
		bcKeys[address] = km
		bscClients[address] = conn
		if i == 0 {
			bscClient = conn
		}
	}

	hub, err := initTokenhub(cfg.Opts["TokenHubContract"])
	if err != nil {
		return nil, err
	}

	return &Connection{
		url:              cfg.Endpoint,
		symbol:           cfg.Symbol,
		bscClient:        bscClient,
		bcKeys:           bcKeys,
		bscClients:       bscClients,
		bcRpcClient:      rpcClient,
		tokenHubContract: hub,
		log:              log,
		stop:             stop,
	}, nil
}

func (c *Connection) ReConnect() error {
	return c.bscClient.Connect()
}

// LatestBlock returns the latest block from the current chain
func (c *Connection) LatestBlock() (uint64, error) {
	blk, err := c.bscClient.LatestBlock()
	if err != nil {
		return 0, err
	}
	return blk.Uint64(), nil
}

func (c *Connection) Address() string {
	return c.bscClient.Keypair().Address()
}

func (c *Connection) TransferVerify(r *submodel.BondRecord) (submodel.BondReason, error) {
	return c.bscClient.BnbTransferVerify(r)
}

func (c *Connection) BscTransferToBc(r *submodel.BondRecord) error {
	addr := common.BytesToAddress(r.Pool)
	bscClient := c.bscClients[addr]
	if bscClient == nil {
		return fmt.Errorf("no bsc client found: %s", addr.Hex())
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
		return fmt.Errorf("no bc key found: %s", addr.Hex())
	}

	receiver := common.HexToAddress(hexutil.Encode(bcKey.GetAddr()))
	value := big.NewInt(0).Add(r.Amount.Int, fee)

	err = bscClient.LockAndUpdateOpts(big.NewInt(0), value)
	if err != nil {
		return err
	}

	tx, err := hub.TransferOut(bscClient.Opts(), ZeroAddress, receiver, r.Amount.Int, 0x17b1f307761)
	bscClient.UnlockOpts()

	if err != nil {
		return err
	}
	c.log.Info("BscTransferToBc txHash", tx.Hash())
	return nil
}

func (c *Connection) BondOrUnbondCall(bond, unbond, leastBond int64) (submodel.BondReportType, int64) {
	c.log.Info("BondOrUnbondCall", "bond", bond, "unbond", unbond)

	if bond < unbond {
		diff := unbond - bond
		if diff < leastBond {
			c.log.Info("bond is smaller than unbond while diff is smaller than leastBond, PureBondReport", "bond", bond, "unbond", unbond, "leastBond", leastBond)
			return submodel.PureBondReport, 0
		}
		return submodel.UnBondReport, diff
	} else if bond > unbond {
		diff := bond - unbond
		if diff < leastBond {
			c.log.Info("unbond is smaller than bond while diff is smaller than leastBond, PureBondReport", "bond", bond, "unbond", unbond, "leastBond", leastBond)
			return submodel.PureBondReport, 0
		}
		return submodel.BondReport, diff
	} else {
		c.log.Info("bond is equal to unbond, NoCall")
		return submodel.BondReport, 0
	}
}

func (c *Connection) Close() {
	c.bscClient.Close()
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
