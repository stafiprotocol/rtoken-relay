package cosmos

import (
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
)

//write to cosmos
type writer struct {
	conn                *Connection
	router              chains.Router
	log                 log15.Logger
	sysErr              chan<- error
	eraPoolUpdatedFlows map[string]*core.EraPoolUpdatedFlow
}

func NewWriter(conn *Connection, log log15.Logger, sysErr chan<- error) *writer {
	return &writer{
		conn:   conn,
		log:    log,
		sysErr: sysErr,
	}
}
func (w *writer) start() error {
	return nil
}

func (w *writer) setRouter(r chains.Router) {
	w.router = r
}

//resolve msg from other chains
func (w *writer) ResolveMessage(m *core.Message) bool {
	switch m.Reason {
	case core.LiquidityBond:
	case core.LiquidityBondResult:
	case core.NewEra:
	case core.EraPoolUpdated:
	case core.NewMultisig:
	case core.MultisigExecuted:
		return true
	default:
		w.log.Warn("message reason unsupported", "reason", m.Reason)
		return false
	}
}
