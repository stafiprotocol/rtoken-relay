package cosmos

import (
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"time"
)

// Frequency of polling for a new block
var (
	BlockRetryInterval = time.Second * 5
	BlockRetryLimit    = 5
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
	for {
		select {
		case <-l.stop:
			return TerminatedError
		default:
			if retry == 0 {
				l.log.Error("pollBlocks error", "symbol", l.symbol)
				return nil
			}

			err := l.updateEra()
			if err != nil {
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}
			time.Sleep(BlockRetryInterval)
			retry = BlockRetryLimit
		}
	}
}

func (l *listener) updateEra() error {

	client, err := l.conn.GetOnePoolClient()
	if err != nil {
		return err
	}

	era, err := client.GetCurrentEra()
	if err != nil {
		return err
	}
	if era <= l.currentEra {
		return nil
	}

	l.log.Info("get a new era, prepare to send message", "symbol", l.symbol, "currentEra", l.currentEra, "newEra", era)
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
