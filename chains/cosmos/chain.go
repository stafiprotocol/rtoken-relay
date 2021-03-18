package cosmos

import (
	"errors"
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
	stop := make(chan int)
	conn, err := NewConnection(cfg, logger, stop)
	if err != nil {
		return nil, err
	}

	bs := new(blockstore.Blockstore)
	startBlock := uint64(0)
	// Setup listener & writer
	l := NewListener(cfg.Name, cfg.Symbol, startBlock, bs, conn, logger, stop, sysErr)
	w := NewWriter(conn, logger, sysErr)
	return &Chain{cfg: cfg, conn: conn, listener: l, writer: w, stop: stop}, nil
}

func (c *Chain) Start() error {
	err := c.listener.start()
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
