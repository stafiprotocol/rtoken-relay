package bnb

import (
	"fmt"

	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

type Chain struct {
	cfg      *core.ChainConfig // The config of the chain
	conn     *Connection       // THe chains connection
	listener *listener         // The listener of this chain
	writer   *writer           // The writer of the chain
	stop     chan<- int
}

func InitializeChain(cfg *core.ChainConfig, logger core.Logger, sysErr chan<- error) (*Chain, error) {
	logger.Info("InitializeChain", "type", "bnb", "symbol", cfg.Symbol)

	stop := make(chan int)
	conn, err := NewConnection(cfg, logger, stop)
	if err != nil {
		return nil, err
	}
	bs, err := chains.NewBlockstore(cfg.Opts["blockstorePath"], conn.BlockStoreUseAddress())
	if err != nil {
		return nil, err
	}
	startBlk, err := chains.StartBlock(bs, cfg.Opts["startBlock"])
	if err != nil {
		return nil, err
	}

	eraSeconds := cfg.Opts["eraSeconds"]
	eraSecondsStr, ok := eraSeconds.(string)
	if !ok {
		panic("eraSeconds not string")
	}
	eraSecondsBig, ok := utils.StringToBigint(eraSecondsStr)
	if !ok {
		panic("eraSeconds is not digital string")
	}
	if eraSecondsBig.Sign() <= 0 {
		panic(fmt.Sprintf("wrong erablock: %s", eraSecondsBig))
	}

	eraOffset := cfg.Opts["eraOffset"]
	eraOffsetStr, ok := eraOffset.(string)
	if !ok {
		panic("eraOffset not string")
	}
	eraOffsetBig, ok := utils.StringToBigint(eraOffsetStr)
	if !ok {
		panic("eraOffset is not digital string")
	}

	// Setup listener & writer
	l := NewListener(cfg.Name, cfg.Symbol, conn, bs, startBlk, eraSecondsBig.Uint64(), eraOffsetBig.Int64(), logger, stop, sysErr)
	w := NewWriter(cfg.Symbol, eraSecondsBig.Uint64(), eraOffsetBig.Int64(), conn, logger, sysErr, stop)
	return &Chain{cfg: cfg, conn: conn, listener: l, writer: w, stop: stop}, nil
}

func (c *Chain) SetRouter(r *core.Router) {
	r.Listen(c.Rsymbol(), c.writer)
	c.listener.setRouter(r)
	c.writer.setRouter(r)
}

func (c *Chain) Start() error {
	return c.listener.start()
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
