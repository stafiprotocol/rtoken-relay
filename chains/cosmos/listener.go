package cosmos

import (
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"time"
)

// Frequency of polling for a new block
var (
	BlockRetryInterval = time.Second * 6
	BlockRetryLimit    = 10
	BlockConfirmNumber = int64(6)
)

//listen event from cosmos
type listener struct {
	name       string
	symbol     core.RSymbol
	startBlock uint64
	blockstore blockstore.Blockstorer
	conn       *Connection
	//subscriptions map[*eventId]eventHandler // Handlers for specific events
	router     chains.Router
	log        log15.Logger
	stop       <-chan int
	sysErr     chan<- error
	currentEra uint32
}

func NewListener(name string, symbol core.RSymbol, startBlock uint64, bs blockstore.Blockstorer, conn *Connection, log log15.Logger, stop <-chan int, sysErr chan<- error) *listener {
	return &listener{
		name:       name,
		symbol:     symbol,
		startBlock: startBlock,
		blockstore: bs,
		conn:       conn,
		//subscriptions: make(map[*eventId]eventHandler),
		log:        log,
		stop:       stop,
		sysErr:     sysErr,
		currentEra: 0,
	}
}

func (l *listener) setRouter(r chains.Router) {
	l.router = r
}

func (l *listener) start() error {

	go func() {
		err := l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

func (l *listener) pollBlocks() error {
	var retry = BlockRetryLimit
	ticker := time.NewTicker(BlockRetryInterval)
	defer ticker.Stop()
	for {
		select {
		case <-l.stop:
			return TerminatedError
		case <-ticker.C:
			if retry <= 0 {
				return fmt.Errorf("poolBlocks reach retry limit ,symbol: %s", l.symbol)
			}

			err := l.updateEra()
			if err != nil {
				retry--
				continue
			}
			retry = BlockRetryLimit
		}
	}
}

func (l *listener) updateEra() error {
	client, err := l.conn.GetOnePoolClient()
	if err != nil {
		return err
	}

	height, era, err := client.GetCurrentEra()
	if err != nil {
		return err
	}
	//update height era
	l.conn.currentHeight = height
	l.currentEra = era

	msg := &core.Message{Destination: core.RFIS, Reason: core.NewEra, Content: era}
	l.submitMessage(msg, nil)
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
		l.log.Error("failed to send message", "err", err, "msg", m)
	}
}
