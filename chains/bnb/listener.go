// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package bnb

import (
	"errors"
	"time"

	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

// Frequency of polling for a new block
var (
	BlockDelay         = uint64(10)
	BlockRetryInterval = time.Second * 10
	EraInterval        = time.Minute * 5
	BlockRetryLimit    = 30
	ErrFatalPolling    = errors.New("listener block polling failed")
)

type listener struct {
	name     string
	symbol   core.RSymbol
	eraBlock uint64
	conn     *Connection
	router   chains.Router
	log      core.Logger
	stop     <-chan int
	sysErr   chan<- error // Reports fatal error to core
}

// NewListener creates and returns a listener
func NewListener(name string, symbol core.RSymbol, opts map[string]interface{}, conn *Connection, log core.Logger, stop <-chan int, sysErr chan<- error) *listener {
	eraBlockCfg := opts["eraBlockCfg"]
	eraBlockStr, ok := eraBlockCfg.(string)
	if !ok {
		panic("eraBlockCfg not string")
	}
	eraBlock, ok := utils.StringToBigint(eraBlockStr)
	if !ok {
		panic("eraBlockCfg not digital string")
	}

	conn.SetEraBlock(eraBlock.Uint64())
	return &listener{
		name:     name,
		symbol:   symbol,
		eraBlock: eraBlock.Uint64(),
		conn:     conn,
		log:      log,
		stop:     stop,
		sysErr:   sysErr,
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
		err := l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.cfg.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before continuing to the next block.
func (l *listener) pollBlocks() error {
	l.log.Info("Polling Blocks...")
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

			latestBlock, err := l.conn.LatestBlock2()
			if err != nil {
				l.log.Error("Unable to get latest block", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			if latestBlock%1000 == 0 {
				l.log.Debug("pollBlocks", "latest", latestBlock)
			}

			era := uint32(uint64(latestBlock) / l.eraBlock)
			l.processEra(era)
			time.Sleep(EraInterval)
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
