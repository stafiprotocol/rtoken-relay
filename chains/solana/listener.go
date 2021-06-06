package solana

import (
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
)

//listen event or block update from solana
type listener struct {
	name   string
	symbol core.RSymbol
	conn   *Connection
	//subscriptions map[*eventId]eventHandler // Handlers for specific events
	router chains.Router
	log    log15.Logger
	stop   <-chan int
	sysErr chan<- error
}

func NewListener(name string, symbol core.RSymbol, conn *Connection, log log15.Logger, stop <-chan int, sysErr chan<- error) *listener {
	return &listener{
		name:   name,
		symbol: symbol,
		conn:   conn,
		//subscriptions: make(map[*eventId]eventHandler),
		log:    log,
		stop:   stop,
		sysErr: sysErr,
	}
}

func (l *listener) setRouter(r chains.Router) {
	l.router = r
}

func (l *listener) start() error {

	go func() {
		
	}()

	return nil
}