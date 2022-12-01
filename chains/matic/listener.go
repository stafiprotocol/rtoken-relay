// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package matic

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"

	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

// Frequency of polling for a new block
var (
	BlockDelay         = uint64(10)
	BlockRetryInterval = time.Second * 10
	EraInterval        = time.Minute * 2
	BlockRetryLimit    = 30
	ErrFatalPolling    = errors.New("listener block polling failed")
)

type listener struct {
	name       string
	symbol     core.RSymbol
	eraBlock   uint64
	startBlock uint64
	conn       *Connection
	router     chains.Router
	blockstore blockstore.Blockstorer
	log        core.Logger
	stop       <-chan int
	sysErr     chan<- error // Reports fatal error to core
}

// NewListener creates and returns a listener
func NewListener(name string, symbol core.RSymbol, opts map[string]interface{}, conn *Connection, startBlock uint64, bs blockstore.Blockstorer, log core.Logger, stop <-chan int, sysErr chan<- error) *listener {
	eraBlockCfg := opts["eraBlockCfg"]
	eraBlockStr, ok := eraBlockCfg.(string)
	if !ok {
		panic("eraBlockCfg not string")
	}
	eraBlock, ok := utils.StringToBigint(eraBlockStr)
	if !ok {
		panic("eraBlockCfg not digital string")
	}
	if eraBlock.Cmp(big.NewInt(0)) <= 0 {
		panic(fmt.Sprintf("wrong erablock: %s", eraBlock))
	}

	return &listener{
		name:       name,
		symbol:     symbol,
		eraBlock:   eraBlock.Uint64(),
		conn:       conn,
		startBlock: startBlock,
		blockstore: bs,
		log:        log,
		stop:       stop,
		sysErr:     sysErr,
	}
}

// sets the router
func (l *listener) setRouter(r chains.Router) {
	l.router = r
}

// start registers all subscriptions provided by the config
func (l *listener) start() error {
	l.log.Debug("Starting listener...")

	go func() {
		err := l.pollEras()
		if err != nil {
			l.log.Error("Polling eras failed", "err", err)
			l.sysErr <- err
		}
	}()

	go func() {
		err := l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
			l.sysErr <- err
		}
	}()

	return nil
}

func (l *listener) pollBlocks() error {
	var willDealBlock = l.startBlock
	var retry = BlockRetryLimit
	for {
		select {
		case <-l.stop:
			l.log.Info("pollBlocks receive stop chan, will stop")
			return nil
		default:
			if retry <= 0 {
				return fmt.Errorf("pollBlocks reach retry limit ,symbol: %s", l.symbol)
			}

			latestBlk, err := l.conn.LatestBlock()
			if err != nil {
				l.log.Error("Failed to fetch latest blockNumber", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}
			// Sleep if the block we want comes after the most recently finalized block
			if willDealBlock+BlockDelay > latestBlk {
				if willDealBlock%100 == 0 {
					l.log.Trace("Block not yet finalized", "target", willDealBlock, "finalBlk", latestBlk)
				}
				time.Sleep(BlockRetryInterval)
				continue
			}

			err = l.processBlockEvents(willDealBlock)
			if err != nil {
				l.log.Error("Failed to process events in block", "block", willDealBlock, "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Write to blockstore
			err = l.blockstore.StoreBlock(new(big.Int).SetUint64(willDealBlock))
			if err != nil {
				l.log.Error("Failed to write to blockstore", "err", err)
				return err
			}
			willDealBlock++

			retry = BlockRetryLimit
		}
	}
}

func (l *listener) processBlockEvents(currentBlock uint64) error {
	if currentBlock%100 == 0 {
		l.log.Debug("processBlockEvents", "blockNum", currentBlock)
	}
	stakeIterator, err := l.conn.stakePortalContract.FilterStake(&bind.FilterOpts{
		Start:   currentBlock,
		End:     &currentBlock,
		Context: context.Background(),
	})

	if err != nil {
		return err
	}
	for stakeIterator.Next() {

	}

	return nil
}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.cfg.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before continuing to the next block.
func (l *listener) pollEras() error {
	l.log.Info("Polling Eras...")
	var retry = BlockRetryLimit
	for {
		select {
		case <-l.stop:
			l.log.Info("get stop signal, stop pool blocks")
			return nil
		default:
			// No more retries, goto next block
			if retry <= 0 {
				l.log.Error("Polling eras failed, retries exceeded")
				return ErrFatalPolling
			}

			latestBlock, err := l.conn.LatestBlock()
			if err != nil {
				l.log.Warn("Unable to get latest block", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			era := uint32(latestBlock / l.eraBlock)

			err = l.processEra(era)
			if err != nil {
				l.log.Warn("processEra failed", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}
			time.Sleep(EraInterval)
			retry = BlockRetryLimit
		}
	}
}

func (l *listener) processEra(era uint32) error {
	msg := &core.Message{Destination: core.RFIS, Reason: core.NewEra, Content: era}
	return l.submitMessage(msg)
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (l *listener) submitMessage(m *core.Message) error {
	m.Source = l.symbol
	return l.router.Send(m)
}
