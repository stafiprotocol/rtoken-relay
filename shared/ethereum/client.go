// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
)

var (
	BlockRetryInterval   = time.Second * 5
	BlockToFinalize      = big.NewInt(3)
	DefaultGasLimit      = big.NewInt(1000000)
	DefaultGasPrice      = big.NewInt(20000000000)
	DefaultExtraGasLimit = big.NewInt(150000)
	ExtraGasPrice        = big.NewInt(10000000000)
)

type Client struct {
	endpoint    string
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
func NewClient(endpoint string, kp *secp256k1.Keypair, log log15.Logger, gasLimit, gasPrice *big.Int) *Client {
	client := &Client{
		endpoint:    endpoint,
		kp:          kp,
		gasLimit:    gasLimit,
		maxGasPrice: gasPrice,
		log:         log,
		stop:        make(chan int),
	}

	if client.gasLimit.Uint64() == 0 {
		client.gasLimit = DefaultGasLimit
	}

	if client.maxGasPrice.Uint64() == 0 {
		client.maxGasPrice = DefaultGasPrice
	}

	return client
}

// Connect starts the ethereum WS connection
func (c *Client) Connect() error {
	c.log.Info("Connecting to ethereum chain...", "url", c.endpoint)
	var rpcClient *rpc.Client
	var err error
	if strings.HasPrefix(c.endpoint, "http") {
		rpcClient, err = rpc.DialHTTP(c.endpoint)
	} else {
		rpcClient, err = rpc.DialWebsocket(context.Background(), c.endpoint, "/ws")
	}
	if err != nil {
		return err
	}
	c.conn = ethclient.NewClient(rpcClient)

	// Construct tx opts, call opts, and nonce mechanism
	opts, _, err := c.newTransactOpts(big.NewInt(0), c.gasLimit, c.maxGasPrice)
	if err != nil {
		return err
	}
	c.opts = opts
	c.nonce = 0
	c.callOpts = &bind.CallOpts{From: c.kp.CommonAddress()}
	return nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (c *Client) newTransactOpts(value, gasLimit, gasPrice *big.Int) (*bind.TransactOpts, uint64, error) {
	privateKey := c.kp.PrivateKey()
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := c.conn.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, 0, err
	}

	auth := bind.NewKeyedTransactor(privateKey)
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
func (c *Client) LockAndUpdateOpts(limit, value *big.Int) error {
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

	if limit.Uint64() == 0 {
		c.opts.GasLimit = DefaultGasLimit.Uint64()
	} else {
		c.opts.GasLimit = big.NewInt(0).Add(limit, DefaultExtraGasLimit).Uint64()
	}

	c.opts.Value = value
	return nil
}

func (c *Client) UnlockOpts() {
	c.opts.Value = big.NewInt(0)
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
func (c *Client) EnsureHasBytecode(addr common.Address) error {
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
func (c *Client) WaitForBlock(block, extra *big.Int) error {
	extra = extra.Add(extra, block)
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
			if currBlock.Cmp(extra) >= 0 {
				return nil
			}
			c.log.Trace("Block not ready, waiting", "target", block, "current", currBlock)
			time.Sleep(BlockRetryInterval)
			continue
		}
	}
}

func (c *Client) TransferVerify(r *submodel.BondRecord, token common.Address) (submodel.BondReason, error) {
	blkHash := common.BytesToHash(r.Blockhash)
	txHash := common.BytesToHash(r.Txhash)

	block, err := c.conn.BlockByHash(context.Background(), blkHash)
	if err != nil {
		return submodel.BondReasonDefault, err
	}

	blk, err := c.conn.BlockByNumber(context.Background(), block.Number())
	if err != nil {
		return submodel.BondReasonDefault, err
	}

	if !bytes.Equal(blk.Hash().Bytes(), r.Blockhash) {
		return submodel.BlockhashUnmatch, nil
	}

	latestBlk, err := c.LatestBlock()
	if err != nil {
		return submodel.BondReasonDefault, err
	}

	num := big.NewInt(0).Add(block.Number(), BlockToFinalize)
	if num.Cmp(latestBlk) > 0 {
		err := c.WaitForBlock(latestBlk, BlockToFinalize)
		if err != nil {
			return submodel.BondReasonDefault, err
		}
	}

	for _, tx := range block.Transactions() {
		if !bytes.Equal(r.Txhash, tx.Hash().Bytes()) {
			continue
		}

		receipt, err := c.conn.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return submodel.BondReasonDefault, err
		}

		for _, elog := range receipt.Logs {
			if !bytes.Equal(elog.Address[:], token.Bytes()) {
				continue
			}

			if len(elog.Topics) != ethmodel.TransferEventTopicLen {
				c.log.Warn("TransferVerify: size of topics not right")
				return submodel.TxhashUnmatch, nil
			}

			if !bytes.Equal(elog.Topics[0].Bytes(), ethmodel.TransferEvent.GetTopic().Bytes()) {
				c.log.Warn("TransferVerify: first topic not TransferEvent")
				return submodel.TxhashUnmatch, nil
			}

			from := common.LeftPadBytes(r.Pubkey, 32)
			if !bytes.Equal(elog.Topics[1].Bytes(), from) {
				c.log.Warn("TransferVerify: second topic not pubkey")
				return submodel.PubkeyUnmatch, nil
			}

			to := common.LeftPadBytes(r.Pool, 32)
			if !bytes.Equal(elog.Topics[2].Bytes(), to) {
				c.log.Warn("TransferVerify: last topic not pool")
				return submodel.PoolUnmatch, nil
			}

			amount := common.LeftPadBytes(r.Amount.Int.Bytes(), 32)
			if !bytes.Equal(elog.Data, amount) {
				c.log.Warn("TransferVerify: data not amount")
				return submodel.AmountUnmatch, nil
			}

			return submodel.Pass, nil
		}
		c.log.Warn("TransferVerify: no address of log equal to the given token contract", "contract", token.Hex())
		return submodel.TxhashUnmatch, nil
	}
	c.log.Warn("TransferVerify: txhash not found", "blockhash", blkHash.Hex(), "txhash", txHash.Hex())
	return submodel.BlockhashUnmatch, nil
}

func (c *Client) BnbTransferVerify(r *submodel.BondRecord) (submodel.BondReason, error) {
	blkHash := common.BytesToHash(r.Blockhash)
	txHash := common.BytesToHash(r.Txhash)

	block, err := c.conn.BlockByHash(context.Background(), blkHash)
	if err != nil {
		return submodel.BondReasonDefault, err
	}

	blk, err := c.conn.BlockByNumber(context.Background(), block.Number())
	if err != nil {
		return submodel.BondReasonDefault, err
	}

	if !bytes.Equal(blk.Hash().Bytes(), r.Blockhash) {
		return submodel.BlockhashUnmatch, nil
	}

	latestBlk, err := c.LatestBlock()
	if err != nil {
		return submodel.BondReasonDefault, err
	}

	num := big.NewInt(0).Add(block.Number(), BlockToFinalize)
	if num.Cmp(latestBlk) > 0 {
		err := c.WaitForBlock(latestBlk, BlockToFinalize)
		if err != nil {
			return submodel.BondReasonDefault, err
		}
	}

	for _, tx := range block.Transactions() {
		if !bytes.Equal(r.Txhash, tx.Hash().Bytes()) {
			continue
		}

		msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), big.NewInt(0))
		if err != nil {
			c.log.Warn("BnbTransferVerify: As message error", "ChainId", tx.ChainId())
			return submodel.TxhashUnmatch, nil
		}

		if !bytes.Equal(r.Pubkey, msg.From().Bytes()) {
			return submodel.PubkeyUnmatch, nil
		}

		if !bytes.Equal(r.Pool, tx.To().Bytes()) {
			return submodel.PoolUnmatch, nil
		}

		if tx.Value().Cmp(r.Amount.Int) != 0 {
			c.log.Warn("BnbTransferVerify: amount not equal", "value", tx.Value(), "amount", r.Amount.Int)
			return submodel.AmountUnmatch, nil
		}

		return submodel.Pass, nil
	}
	c.log.Warn("TransferVerify: txhash not found", "blockhash", blkHash.Hex(), "txhash", txHash.Hex())
	return submodel.BlockhashUnmatch, nil
}

func (c *Client) TransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	return c.conn.TransactionReceipt(context.Background(), hash)
}

// Close terminates the client connection and stops any running routines
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
	close(c.stop)
}
