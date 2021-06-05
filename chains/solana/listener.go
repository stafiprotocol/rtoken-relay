package solana

import (
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
)

//listen event or block update from solana
type listener struct {
	name       string
	symbol     core.RSymbol
	startBlock uint64
	conn       *Connection
	//subscriptions map[*eventId]eventHandler // Handlers for specific events
	router chains.Router
	log    log15.Logger
	stop   <-chan int
	sysErr chan<- error
}
