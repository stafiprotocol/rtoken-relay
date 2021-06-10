// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
)

var (
	BlockRetryInterval = time.Second * 5
	ExtraGasPrice      = big.NewInt(10000000000)

	DefaultGasLimit = big.NewInt(1000000)
	DefaultGasPrice = big.NewInt(20000000000)
)

type Client struct {
	endpoint    string
	chainId     *big.Int
	http        bool
	kp          *secp256k1.Keypair
	gasLimit    *big.Int
	maxGasPrice *big.Int
	conn        *ethclient.Client
	opts        *bind.TransactOpts
	callOpts    *bind.CallOpts
	nonce       uint64
	optsLock    sync.Mutex
	log         log15.Logger
	stop        chan int // All routines should exit when this channel is closed
}

// NewClient returns an uninitialized connection, must call Client.Connect() before using.
func NewClient(endpoint string, http bool, kp *secp256k1.Keypair, log log15.Logger) *Client {
	return &Client{
		endpoint:    endpoint,
		http:        http,
		kp:          kp,
		gasLimit:    DefaultGasLimit,
		maxGasPrice: DefaultGasPrice,
		log:         log,
		stop:        make(chan int),
	}
}

// Connect starts the ethereum WS connection
func (c *Client) Connect() error {
	c.log.Info("Connecting to ethereum chain...", "url", c.endpoint)
	var rpcClient *rpc.Client
	var err error
	// Start http or ws client
	if c.http {
		rpcClient, err = rpc.DialHTTP(c.endpoint)
	} else {
		rpcClient, err = rpc.DialWebsocket(context.Background(), c.endpoint, "/ws")
	}
	if err != nil {
		return err
	}
	c.conn = ethclient.NewClient(rpcClient)

	// Construct tx opts, call opts, and nonce mechanism
	opts, _, err := c.newTransactOpts(big.NewInt(0), c.gasLimit, c.maxGasPrice, c.chainId)
	if err != nil {
		return err
	}
	c.opts = opts
	c.nonce = 0
	c.callOpts = &bind.CallOpts{From: c.kp.CommonAddress()}
	return nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (c *Client) newTransactOpts(value, gasLimit, gasPrice, chainId *big.Int) (*bind.TransactOpts, uint64, error) {
	privateKey := c.kp.PrivateKey()
	address := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := c.conn.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, 0, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, 0, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = uint64(gasLimit.Int64())
	auth.GasPrice = gasPrice
	auth.Context = context.Background()

	return auth, nonce, nil
}

func (c *Client) Keypair() *secp256k1.Keypair {
	return c.kp
}

func (c *Client) Client() *ethclient.Client {
	return c.conn
}

func (c *Client) Opts() *bind.TransactOpts {
	return c.opts
}

func (c *Client) CallOpts() *bind.CallOpts {
	return c.callOpts
}

func (c *Client) SafeEstimateGas(ctx context.Context) (*big.Int, error) {
	gasPrice, err := c.conn.SuggestGasPrice(context.TODO())
	if err != nil {
		return nil, err
	}

	return gasPrice.Add(gasPrice, ExtraGasPrice), nil
}

// LockAndUpdateOpts acquires a lock on the opts before updating the nonce
// and gas price.
func (c *Client) LockAndUpdateOpts() error {
	c.optsLock.Lock()

	gasPrice, err := c.SafeEstimateGas(context.TODO())
	if err != nil {
		return err
	}
	c.opts.GasPrice = gasPrice

	nonce, err := c.conn.PendingNonceAt(context.Background(), c.opts.From)
	if err != nil {
		c.optsLock.Unlock()
		return err
	}
	c.opts.Nonce.SetUint64(nonce)
	return nil
}

func (c *Client) UnlockOpts() {
	c.optsLock.Unlock()
}

// LatestBlock returns the latest block from the current chain
func (c *Client) LatestBlock() (*big.Int, error) {
	header, err := c.conn.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

// EnsureHasBytecode asserts if contract code exists at the specified address
func (c *Client) EnsureHasBytecode(addr ethcommon.Address) error {
	code, err := c.conn.CodeAt(context.Background(), addr, nil)
	if err != nil {
		return err
	}

	if len(code) == 0 {
		return fmt.Errorf("no bytecode found at %s", addr.Hex())
	}
	return nil
}

// WaitForBlock will poll for the block number until the current block is equal or greater than
func (c *Client) WaitForBlock(block *big.Int) error {
	blk := big.NewInt(3)
	blk = blk.Add(blk, block)
	for {
		select {
		case <-c.stop:
			return errors.New("connection terminated")
		default:
			currBlock, err := c.LatestBlock()
			if err != nil {
				return err
			}

			// Greater than target
			if currBlock.Cmp(blk) >= 0 {
				return nil
			}
			c.log.Trace("Block not ready, waiting", "target", block, "current", currBlock)
			time.Sleep(BlockRetryInterval)
			continue
		}
	}
}

// Close terminates the client connection and stops any running routines
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
	close(c.stop)
}
