package solana

import (
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/chains"
)

//write to solana
type writer struct {
	conn   *Connection
	router chains.Router
	log    log15.Logger
	sysErr chan<- error
}
