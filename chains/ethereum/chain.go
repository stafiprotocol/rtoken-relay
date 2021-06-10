// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only
/*
The ethereum package contains the logic for interacting with ethereum chains.

There are 3 major components: the connection, the listener, and the writer.
The currently supported transfer types are Fungible (ERC20).

Client

The connection contains the ethereum RPC client and can be accessed by both the writer and listener.

Listener

The listener polls for each new block and looks for deposit events in the bridge contract. If a deposit occurs, the listener will fetch additional information from the handler before constructing a message and forwarding it to the router.

Writer

The writer recieves the message and creates a proposals on-chain. Once a proposal is made, the writer then watches for a finalization event and will attempt to execute the proposal if a matching event occurs. The writer skips over any proposals it has already seen.
*/
package ethereum

import (
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
)

type Chain struct {
	cfg      *core.ChainConfig // The config of the chain
	conn     *Connection       // THe chains connection
	listener *listener         // The listener of this chain
	writer   *writer           // The writer of the chain
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

	bs := new(blockstore.Blockstore)
	startBlk, err := chains.StartBlock(cfg, latestBlock, bs, conn.Address())
	if err != nil {
		return nil, err
	}

	// Setup listener & writer
	l := NewListener(cfg.Name, cfg.Symbol, cfg.Opts, startBlk, bs, conn, logger, stop, sysErr)
	w := NewWriter(cfg.Symbol, conn, logger, sysErr, stop)
	return &Chain{cfg: cfg, conn: conn, listener: l, writer: w, stop: stop}, nil
}

func (c *Chain) SetRouter(r *core.Router) {
	c.listener.setRouter(r)
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

	c.writer.log.Debug("Successfully started chain")
	return nil
}

func (c *Chain) Rsymbol() core.RSymbol {
	return c.cfg.Symbol
}

func (c *Chain) Name() string {
	return c.cfg.Name
}

// Stop signals to any running routines to exit
func (c *Chain) Stop() {
	close(c.stop)
}

func (c *Chain) InitBondedPools(symbols []core.RSymbol) error {
	return nil
}
