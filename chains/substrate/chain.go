// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"errors"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
)

var TerminatedError = errors.New("terminated")

type Chain struct {
	cfg      *core.ChainConfig // The config of the chain
	conn     *Connection
	listener *listener // The listener of this chain
	writer   *writer   // The writer of the chain
	stop     chan<- int
}

func InitializeChain(cfg *core.ChainConfig, logger log15.Logger, sysErr chan<- error) (*Chain, error) {
	logger.Info("InitializeChain", "type", "substrate", "symbol", cfg.Symbol)

	stop := make(chan int)
	conn, err := NewConnection(cfg, logger, stop)
	if err != nil {
		return nil, err
	}

	latestBlock, err := conn.LatestBlockNumber()
	if err != nil {
		return nil, err
	}

	bs, err := chains.NewBlockstore(cfg.Opts["blockstorePath"], conn.Address())
	if err != nil {
		return nil, err
	}

	var startBlk uint64
	if !cfg.LatestBlockFlag {
		startBlk, err = chains.StartBlock(bs, cfg.Opts["startBlock"])
		if err != nil {
			return nil, err
		}
	} else {
		startBlk = latestBlock
	}

	// Setup listener & writer
	l := NewListener(cfg.Name, cfg.Symbol, cfg.Opts, startBlk, bs, conn, logger, stop, sysErr)
	w := NewWriter(cfg.Symbol, cfg.Opts, conn, logger, sysErr, stop)
	return &Chain{cfg: cfg, conn: conn, listener: l, writer: w, stop: stop}, nil
}

func (c *Chain) Start() error {
	err := c.listener.start()
	if err != nil {
		return err
	}

	err = c.writer.start()
	if err != nil {
		return err
	}
	return nil
}

func (c *Chain) SetRouter(r *core.Router) {
	r.Listen(c.Rsymbol(), c.writer)
	c.listener.setRouter(r)
	c.writer.setRouter(r)
}

func (c *Chain) Rsymbol() core.RSymbol {
	return c.cfg.Symbol
}

func (c *Chain) Name() string {
	return c.cfg.Name
}

func (c *Chain) Stop() {
	close(c.stop)
}

func (c *Chain) InitBondedPools(symbols []core.RSymbol) error {
	//only stafi need
	if c.Rsymbol() != core.RFIS {
		return nil
	}
	for _, symbol := range symbols {
		bondedPools, err := c.writer.getLatestBondedPools(symbol)
		if err != nil {
			return err
		}
		msg := &core.Message{Source: core.RFIS, Destination: symbol, Reason: core.BondedPools, Content: bondedPools}
		if !c.writer.submitMessage(msg) {
			return errors.New("init bondedPools submitMessage failed")
		}
	}
	return nil
}
