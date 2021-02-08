// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/substrate"
	"math/big"
	"strings"
	"time"

	"github.com/stafiprotocol/chainbridge/utils/blockstore"
)

type listener struct {
	gsrpc      *substrate.GsrpcClient
	sarpc      *substrate.SarpcClient
	bs         blockstore.Blockstorer
	startBlock uint64
	validators map[conn.RSymbol]conn.Validator
	sysErr     chan<- error
	ctx        context.Context
	log        log15.Logger
}

// Frequency of polling for a new block
var BlockRetryInterval = time.Second * 5
var BlockRetryLimit = 5

func NewListener(ctx context.Context, sarpc *substrate.SarpcClient, gsrpc *substrate.GsrpcClient, bs blockstore.Blockstorer,
	startBlock uint64, validators map[conn.RSymbol]conn.Validator, sysErr chan<- error, log log15.Logger) *listener {
	return &listener{
		sarpc:      sarpc,
		gsrpc:      gsrpc,
		bs:         bs,
		startBlock: startBlock,
		validators: validators,
		sysErr:     sysErr,
		ctx:        ctx,
		log:        log,
	}
}

// Start creates the initial subscription for all events
func (l *listener) Start() error {
	// Check whether latest is less than starting block
	header, err := l.gsrpc.GetHeaderLatest()
	if err != nil {
		return err
	}

	if uint64(header.Number) < l.startBlock {
		return fmt.Errorf("starting block (%d) is greater than latest known block (%d)", l.startBlock, header.Number)
	}

	go func() {
		err = l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before returning with an error.
func (l *listener) pollBlocks() error {
	var currentBlock = l.startBlock
	var retry = BlockRetryLimit
	for {
		select {
		case <-l.ctx.Done():
			return errors.New("terminated")
		default:
			// No more retries, goto next block
			if retry == 0 {
				l.sysErr <- errors.New("event polling retries exceeded")
				return nil
			}

			// Get finalized block hash
			finalizedHash, err := l.gsrpc.GetFinalizedHead()
			if err != nil {
				l.log.Error("Failed to fetch finalized hash", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Get finalized block header
			finalizedHeader, err := l.gsrpc.GetHeader(finalizedHash)
			if err != nil {
				l.log.Error("Failed to fetch finalized header", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Sleep if the block we want comes after the most recently finalized block
			if currentBlock > uint64(finalizedHeader.Number) {
				l.log.Trace("Block not yet finalized", "target", currentBlock, "latest", finalizedHeader.Number)
				time.Sleep(BlockRetryInterval)
				continue
			}

			err = l.processEvents(currentBlock)
			if err != nil {
				l.log.Error("Failed to process events in block", "block", currentBlock, "err", err)
				if strings.Contains(err.Error(), "close 1006") || strings.Contains(err.Error(), "websocket: not connected") {
					l.log.Info("listener", "is webscoket connected", l.sarpc.IsConnected())
					l.sarpc.WebsocketReconnect()
				}
				retry--
				continue
			}

			// Write to blockstore
			err = l.bs.StoreBlock(big.NewInt(0).SetUint64(currentBlock))
			if err != nil {
				l.log.Error("Failed to write to blockstore", "err", err)
			}

			currentBlock++
			retry = BlockRetryLimit
		}
	}
}

func (l *listener) processEvents(blockNum uint64) error {
	if blockNum%100 == 0 {
		l.log.Debug("processEvents", "blockNum", blockNum)
	}

	evts, err := l.sarpc.GetEventsByModuleIdAndEventId(blockNum, config.LiquidityBondModuleId, config.LiquidityBondEventId)
	if err != nil {
		return err
	}

	l.log.Trace("block", "LiquidityEventNum", len(evts), "blockNum", blockNum)

	for _, evt := range evts {
		err := l.processLiquidityBondEvent(evt)
		if err != nil {
			l.log.Error("processLiquidityBondEvent", "error", err)
		}
	}

	for _, evt := range evts {
		l.log.Info("LiquidityBondEvents", "evt", evt)
	}

	return nil
}

func (l *listener) processLiquidityBondEvent(evt *substrate.ChainEvent) error {
	fmt.Println(evt)
	return nil
}
