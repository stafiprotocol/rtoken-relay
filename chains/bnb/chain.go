package bnb

import (
	"github.com/ChainSafe/log15"
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

	latestBlock, err := conn.LatestBlock()
	if err != nil {
		return nil, err
	}

	bs, err := chains.NewBlockstore(cfg.Opts["blockstorePath"], conn.Address())
	if err != nil {
		return nil, err
	}

	var startBlk uint64
	if cfg.LatestBlockFlag {
		startBlk = latestBlock
	} else {
		startBlk, err = chains.StartBlock(bs, cfg.Opts["startBlock"])
		if err != nil {
			return nil, err
		}
	}

	// Setup listener & writer
	l := NewListener(cfg.Name, cfg.Symbol, cfg.Opts, startBlk, bs, conn, logger, stop, sysErr)
	w := NewWriter(cfg.Symbol, conn, logger, sysErr, stop)
	return &Chain{cfg: cfg, conn: conn, listener: l, writer: w, stop: stop}, nil
}

func (c *Chain) SetRouter(r *core.Router) {
	r.Listen(c.Rsymbol(), c.writer)
	c.listener.setRouter(r)
	c.writer.setRouter(r)
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
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Chain) InitBondedPools(symbols []core.RSymbol) error {
	return nil
}
