// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"errors"
	"strconv"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
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
	logger.Info("InitializeChain", "symbol", cfg.Symbol)

	stop := make(chan int)
	conn, err := NewConnection(cfg, logger, stop)
	if err != nil {
		return nil, err
	}

	blk := parseStartBlock(cfg)
	bs := new(blockstore.Blockstore)
	if cfg.LatestBlockFlag {
		blk, err = conn.LatestBlockNumber()
		if err != nil {
			return nil, err
		}
	} else {
		bp := cfg.Opts["blockstorePath"]
		if bp == nil {
			return nil, errors.New("blockstorePath nil")
		}

		bsPath, ok := bp.(string)
		if !ok {
			return nil, errors.New("blockstorePath not string")
		}
		bs, err = blockstore.NewBlockstore(bsPath, 100, conn.Address())
		if err != nil {
			return nil, err
		}

		blk, err = checkBlockstore(bs, blk)
		if err != nil {
			return nil, err
		}
	}

	// Setup listener & writer
	l := NewListener(cfg.Name, cfg.Symbol, cfg.Opts, blk, bs, conn, logger, stop, sysErr)
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

// checkBlockstore queries the blockstore for the latest known block. If the latest block is
// greater than startBlock, then the latest block is returned, otherwise startBlock is.
func checkBlockstore(bs *blockstore.Blockstore, startBlock uint64) (uint64, error) {
	latestBlock, err := bs.TryLoadLatestBlock()
	if err != nil {
		return 0, err
	}

	if latestBlock.Uint64() > startBlock {
		return latestBlock.Uint64(), nil
	} else {
		return startBlock, nil
	}
}

func parseStartBlock(cfg *core.ChainConfig) uint64 {
	if blk, ok := cfg.Opts["startBlock"]; ok {
		blkStr, ok := blk.(string)
		if !ok {
			panic("block not string")
		}
		res, err := strconv.ParseUint(blkStr, 10, 32)
		if err != nil {
			panic(err)
		}
		return res
	}
	return 0
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
