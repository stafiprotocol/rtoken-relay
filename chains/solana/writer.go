package solana

import (
	"fmt"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
)

//write to solana
type writer struct {
	conn   *Connection
	router chains.Router
	log    log15.Logger
	sysErr chan<- error
}

func NewWriter(conn *Connection, log log15.Logger, sysErr chan<- error) *writer {
	return &writer{
		conn:   conn,
		log:    log,
		sysErr: sysErr,
	}
}

func (w *writer) setRouter(r chains.Router) {
	w.router = r
}


//resolve msg from other chains
func (w *writer) ResolveMessage(m *core.Message) (processOk bool) {
	defer func() {
		if !processOk {
			panic(fmt.Sprintf("resolveMessage process failed. %+v", m))
		}
	}()

	switch m.Reason {
	case core.LiquidityBond:
		return w.processLiquidityBond(m)
	case core.BondedPools:
		return w.processBondedPools(m)
	case core.EraPoolUpdated:
		return w.processEraPoolUpdatedEvt(m)
	case core.BondReportEvent:
		return w.processBondReportEvent(m)
	case core.ActiveReportedEvent:
		return w.processActiveReportedEvent(m)
	default:
		w.log.Warn("message reason unsupported", "reason", m.Reason)
		return true
	}
}

func (w *writer) processBondedPools(m *core.Message) bool {
	return true
}


