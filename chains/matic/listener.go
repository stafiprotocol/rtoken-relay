// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package matic

import (
	"errors"
	"math/big"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
)

// Frequency of polling for a new block
var (
	BlockDelay         = uint64(10)
	BlockRetryInterval = time.Second * 10
	BlockRetryLimit    = 30
	ErrFatalPolling    = errors.New("listener block polling failed")
	BlockIntervalToProcessEra = uint64(10)
)

type listener struct {
	name         string
	symbol       core.RSymbol
	startBlock   uint64
	eraBlock uint64
	conn         *Connection
	router       chains.Router
	log          log15.Logger
	blockstore   blockstore.Blockstorer
	stop         <-chan int
	sysErr       chan<- error // Reports fatal error to core
	nextEra uint32
}

// NewListener creates and returns a listener
func NewListener(name string, symbol core.RSymbol, opts map[string]interface{}, startBlock uint64, bs blockstore.Blockstorer, conn *Connection, log log15.Logger, stop <-chan int, sysErr chan<- error) *listener {
	eraBlockCfg := opts["eraBlockCfg"]
	eraBlock, ok := eraBlockCfg.(int)
	if !ok {
		panic("eraBlockIntervalCfg not int")
	}

	return &listener{
		name:         name,
		symbol:       symbol,
		startBlock:   startBlock,
		eraBlock: uint64(eraBlock),
		conn:         conn,
		log:          log,
		blockstore:   bs,
		stop:         stop,
		sysErr:       sysErr,
	}
}

// setContracts sets the listener with the appropriate contracts
func (l *listener) setContracts() {
}

// sets the router
func (l *listener) setRouter(r chains.Router) {
	l.router = r
}

// start registers all subscriptions provided by the config
func (l *listener) start() error {
	l.log.Debug("Starting listener...")

	l.currentChainEra()

	go func() {
		err := l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

func (l *listener) currentChainEra() {
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	CurrentEra := make(chan uint32)
	msg := &core.Message{Destination: core.RFIS, Reason: core.CurrentChainEra, Content: CurrentEra}
	l.submitMessage(msg, nil)

	l.log.Debug("wait current era from stafi")
	select {
	case <-timer.C:
		panic("timeout to get current era")
	case era := <-CurrentEra:
		l.nextEra = era + 1
	}
}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.cfg.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before continuing to the next block.
func (l *listener) pollBlocks() error {
	l.log.Info("Polling Blocks...")
	var currentBlock = l.startBlock
	var retry = BlockRetryLimit
	for {
		select {
		case <-l.stop:
			return errors.New("polling terminated")
		default:
			// No more retries, goto next block
			if retry == 0 {
				l.log.Error("Polling failed, retries exceeded")
				l.sysErr <- ErrFatalPolling
				return nil
			}

			latestBlock, err := l.conn.LatestBlock()
			if err != nil {
				l.log.Error("Unable to get latest block", "block", currentBlock, "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Sleep if the difference is less than BlockDelay; (latest - current) < BlockDelay
			if currentBlock+BlockDelay > latestBlock {
				l.log.Debug("Block not ready, will retry", "target", currentBlock, "latest", latestBlock)
				time.Sleep(BlockRetryInterval)
				continue
			}

			if l.nextEra == 1 || (currentBlock - l.startBlock) >= l.eraBlock {
				l.log.Info("time to process era", "era", l.nextEra, "currentBlock", currentBlock, "eraBlock", l.eraBlock)
				l.startBlock = currentBlock
				l.processEra(l.nextEra)
				l.nextEra++
			}

			// Write to blockstore
			err = l.blockstore.StoreBlock(big.NewInt(0).SetUint64(currentBlock))
			if err != nil {
				l.log.Error("Failed to write to blockstore", "err", err)
			}
			currentBlock++
			retry = BlockRetryLimit
		}
	}
}

func (l *listener) processEra(era uint32) {
	msg := &core.Message{Destination: core.RFIS, Reason: core.NewEra, Content: era}
	l.submitMessage(msg, nil)
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (l *listener) submitMessage(m *core.Message, err error) {
	if err != nil {
		l.log.Error("Critical error before sending message", "err", err)
		return
	}
	m.Source = l.symbol
	err = l.router.Send(m)
	if err != nil {
		l.log.Error("failed to send message", "err", err, "msg", m)
	}
}
