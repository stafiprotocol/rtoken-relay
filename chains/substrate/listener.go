// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
)

type listener struct {
	name          string
	symbol        core.RSymbol
	startBlock    uint64
	blockstore    blockstore.Blockstorer
	conn          *Connection
	subscriptions map[*eventId]eventHandler // Handlers for specific events
	router        chains.Router
	log           log15.Logger
	stop          <-chan int
	sysErr        chan<- error
	currentEra    uint32
}

// Frequency of polling for a new block
var (
	BlockRetryInterval = time.Second * 5
	BlockRetryLimit    = 5
)

func NewListener(name string, symbol core.RSymbol, startBlock uint64, bs blockstore.Blockstorer, conn *Connection, log log15.Logger, stop <-chan int, sysErr chan<- error) *listener {
	return &listener{
		name:          name,
		symbol:        symbol,
		startBlock:    startBlock,
		blockstore:    bs,
		conn:          conn,
		subscriptions: make(map[*eventId]eventHandler),
		log:           log,
		stop:          stop,
		sysErr:        sysErr,
		currentEra:    0,
	}
}

func (l *listener) setRouter(r chains.Router) {
	l.router = r
}

// Start creates the initial subscription for all events
func (l *listener) start() error {
	latestBlk, err := l.conn.LatestBlockNumber()
	if err != nil {
		return err
	}

	if latestBlk < l.startBlock {
		return fmt.Errorf("starting block (%d) is greater than latest known block (%d)", l.startBlock, latestBlk)
	}

	for _, sub := range Subscriptions {
		if sub.eId.symbol != l.symbol {
			continue
		}

		err := l.registerEventHandler(sub.eId, sub.handler)
		if err != nil {
			return err
		}
	}

	go func() {
		err = l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

// registerEventHandler enables a handler for a given event. This cannot be used after Start is called.
func (l *listener) registerEventHandler(id *eventId, handler eventHandler) error {
	if l.subscriptions[id] != nil {
		return fmt.Errorf("event %s, %s already registered", id.symbol, id.name)
	}
	l.subscriptions[id] = handler
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
		case <-l.stop:
			return TerminatedError
		default:
			// No more retries, goto next block
			if retry == 0 {
				if l.symbol == "Default" {
					l.sysErr <- fmt.Errorf("event polling retries exceeded: %s", l.symbol)
				} else {
					l.log.Error("pollBlocks error", "symbol", l.symbol)
				}

				return nil
			}

			finalBlk, err := l.conn.FinalizedBlockNumber()
			if err != nil {
				l.log.Error("Failed to fetch latest blockNumber", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Sleep if the block we want comes after the most recently finalized block
			if currentBlock+l.blockDelay() > finalBlk {
				l.log.Trace("Block not yet finalized", "target", currentBlock, "finalBlk", finalBlk)
				time.Sleep(BlockRetryInterval)
				continue
			}

			if l.symbol != core.RFIS {
				err = l.processEra()
			}

			err = l.processEvents(currentBlock)
			if err != nil {
				l.log.Error("Failed to process events in block", "block", currentBlock, "err", err)
				if strings.Contains(err.Error(), "close 1006") || strings.Contains(err.Error(), "websocket: not connected") {
					l.log.Info("listener", "is webscoket connected", l.conn.IsConnected())
					if err := l.conn.Reconnect(); err != nil {
						l.log.Error("listener", "websocket reconnect error", err)
					}
				}
				retry--
				continue
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

func (l *listener) processEra() error {
	era, err := l.conn.CurrentEra()
	if err != nil {
		return err
	}

	if era == l.currentEra {
		return nil
	}

	l.log.Info("get a new era, prepare to send message", "symbol", l.symbol, "currentEra", l.currentEra, "newEra", era)
	l.currentEra = era
	msg := &core.Message{Destination: core.RFIS, Reason: core.NewEra, Content: era}
	l.submitMessage(msg, nil)
	return nil
}

// processEvents fetches a block and parses out the events, calling Listener.handleEvents()
func (l *listener) processEvents(blockNum uint64) error {
	if blockNum%100 == 0 {
		l.log.Debug("processEvents", "blockNum", blockNum)
	}

	evts, err := l.conn.GetEvents(blockNum)
	if err != nil {
		return err
	}

	for _, evt := range evts {
		switch l.symbol {
		case core.RFIS:
			if evt.ModuleId == config.LiquidityBondModuleId && evt.EventId == config.LiquidityBondEventId {
				l.log.Trace("Handling LiquidityBondEvent event")
				flow, err := l.processLiquidityBondEvent(evt)
				if err != nil {
					return err
				}
				if l.subscriptions[LiquidityBond] != nil {
					l.submitMessage(l.subscriptions[LiquidityBond](flow, l.log))
				}
			} else if evt.ModuleId == config.RTokenLedgerModuleId && evt.EventId == config.EraPoolUpdatedEventId {
				l.log.Trace("Handling EraPoolUpdatedEvent event")
				flow, err := l.processEraPoolUpdatedEvt(evt)
				if err != nil {
					return err
				}
				if l.subscriptions[EraPoolUpdated] != nil {
					l.submitMessage(l.subscriptions[EraPoolUpdated](flow, l.log))
				}
			}
		case core.RDOT:
			if evt.ModuleId == config.MultisigModuleId && evt.EventId == config.NewMultisigEventId {
				l.log.Trace("Handling NewMultisigEvent event")

			} else if evt.ModuleId == config.MultisigModuleId && evt.EventId == config.MultisigExecutedEventId {

			}
		default:
			l.log.Error("process event unsupport rsymbol", "rsymbol", l.symbol)
		}
	}

	return nil
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
		l.log.Error("failed to send message", "err", err)
	}
}

func (l *listener) blockDelay() uint64 {
	switch l.symbol {
	case core.RFIS:
		return 5
	default:
		return 0
	}
}
