package solana

import (
	"errors"

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
