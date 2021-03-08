// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/substrate"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

type listener struct {
	gsrpc      *substrate.GsrpcClient
	sarpc      *substrate.SarpcClient
	bs         blockstore.Blockstorer
	startBlock uint64
	chains     map[conn.RSymbol]conn.Chain
	chainEras  map[conn.RSymbol]types.U32
	sysErr     chan<- error
	ctx        context.Context
	log        log15.Logger
}

// Frequency of polling for a new block
var (
	BlockRetryInterval = time.Second * 5
	BlockRetryLimit = 5
	TerminatedError = errors.New("terminated")
	BlockDelay uint64 = 5
)

func NewListener(ctx context.Context, sarpc *substrate.SarpcClient, gsrpc *substrate.GsrpcClient, bs blockstore.Blockstorer,
	startBlock uint64, chainEras map[conn.RSymbol]types.U32,
	chains map[conn.RSymbol]conn.Chain, sysErr chan<- error, log log15.Logger) *listener {
	return &listener{
		sarpc:      sarpc,
		gsrpc:      gsrpc,
		bs:         bs,
		startBlock: startBlock,
		chains:     chains,
		chainEras:  chainEras,
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

	for sym, chain := range l.chains {
		go l.updateEra(sym, chain)
	}

	go func() {
		err := l.pollBlocks()
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
			return TerminatedError
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
			if currentBlock + BlockDelay > uint64(finalizedHeader.Number) {
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

	allEvts, err := l.sarpc.GetEvents(blockNum)
	if err != nil {
		return err
	}

	eraUpdatedEvts := make([]*substrate.ChainEvent, 0)
	bondEvts := make([]*substrate.ChainEvent, 0)

	for _, evt := range allEvts {
		if evt.ModuleId == config.RTokenLedgerModuleId && evt.EventId == config.EraPoolUpdatedEventId {
			eraUpdatedEvts = append(eraUpdatedEvts, evt)
			continue
		}

		if evt.ModuleId == config.LiquidityBondModuleId && evt.EventId == config.LiquidityBondEventId {
			bondEvts = append(bondEvts, evt)
		}
	}

	if len(eraUpdatedEvts) > 0 {
		l.log.Trace("processEvents", "EraPoolUpdatedEventNum", len(eraUpdatedEvts), "blockNum", blockNum)
		err = l.processEraPoolUpdatedEvts(eraUpdatedEvts)
		if err != nil {
			return err
		}
	}

	if len(bondEvts) > 0 {
		l.log.Trace("processEvents", "LiquidityBondEventNum", len(bondEvts), "blockNum", blockNum)
		if err := l.processLiquidityBondEvents(bondEvts); err != nil {
			l.log.Error("processLiquidityBondEvent", "error", err)
			return err
		}
	}

	return nil
}

func (l *listener) updateEra(sym conn.RSymbol, chain conn.Chain) {
	for {
		select {
		case <-l.ctx.Done():
			return
		default:
			local, ok := l.chainEras[sym]
			if !ok {
				l.sysErr <- fmt.Errorf("chainEra not found for symbol: %s", sym)
				return
			}
			l.log.Info("updateEra")

			foreign, err := chain.CurrentEra()
			if !ok {
				l.sysErr <- fmt.Errorf("get current era error %s for %s", err, sym)
				return
			}
			l.log.Info("updateEra", "local", local, "foreign", foreign)

			if foreign < local {
				//l.sysErr <- fmt.Errorf("new era %d is smaller than old %d", foreign, local)
				time.Sleep(1 * time.Minute)
				continue
			} else if foreign > local {
				h, err := utils.Blake2Hash(foreign)
				if err != nil {
					l.sysErr <- fmt.Errorf("Blake2Hash error %s for %d, %s", err, foreign, sym)
					return
				}

				bk := &conn.BondKey{sym, h}
				prop, err := l.newUpdateEraProposal(bk, foreign)
				if err != nil {
					l.sysErr <- fmt.Errorf("new era %d is smaller than old %d", foreign, local)
					return
				}

				result := l.resolveProposal(prop, true)
				l.log.Info("updateEra", "symbol", sym, "era", foreign, "result", result)
				if result {
					l.chainEras[sym] = foreign
					time.Sleep(1 * time.Minute)
				}
			} else {
				time.Sleep(1 * time.Minute)
			}

		}
	}
}
