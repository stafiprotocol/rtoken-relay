// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package matic

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ChainSafe/log15"
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
	name     string
	symbol   core.RSymbol
	eraBlock uint64
	conn     *Connection
	router   chains.Router
	log      log15.Logger
	stop     <-chan int
	sysErr   chan<- error // Reports fatal error to core
}

// NewListener creates and returns a listener
func NewListener(name string, symbol core.RSymbol, opts map[string]interface{}, conn *Connection, log log15.Logger, stop <-chan int, sysErr chan<- error) *listener {
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
			l.sysErr <- err
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
			l.log.Info("get stop signal, stop pool blocks")
			return nil
		default:
			// No more retries, goto next block
			if retry <= 0 {
				l.log.Error("Polling failed, retries exceeded")
				return ErrFatalPolling
			}

			latestBlock, err := l.conn.LatestBlock()
			if err != nil {
				if err.Error() == "client is closed" {
					err = l.conn.ReConnect()
					if err != nil {
						return err
					}
				}

				l.log.Error("Unable to get latest block", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			if latestBlock%100 == 0 {
				l.log.Debug("pollBlocks", "latest", latestBlock)
			}

			era := uint32(latestBlock / l.eraBlock)

			err = l.processEra(era)
			if err != nil {
				l.log.Error("processEra failed", "err", err)
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
