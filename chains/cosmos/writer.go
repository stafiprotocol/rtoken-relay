package cosmos

import (
	"encoding/hex"
	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
)

//write to cosmos
type writer struct {
	conn          *Connection
	router        chains.Router
	log           log15.Logger
	sysErr        chan<- error
	multisigFlows map[string]*core.MultisigFlow
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
		return w.processLiquidityBond(m)
	case core.EraPoolUpdated:
		return w.processEraPoolUpdated(m)
	case core.NewMultisig:
	case core.MultisigExecuted:
		return true
	default:
		w.log.Warn("message reason unsupported", "reason", m.Reason)
		return false
	}
	return true
}

func (w *writer) processLiquidityBond(m *core.Message) bool {
	flow, ok := m.Content.(*core.BondFlow)
	if !ok {
		w.printContentError(m)
		return false
	}

	if flow.Reason != core.BondReasonDefault {
		w.log.Error("processLiquidityBond receive a message of which reason is not default", "bondId", flow.Key.BondId.Hex(), "reason", flow.Reason)
		return false
	}

	bondReason, err := w.conn.TransferVerify(flow.Record)
	if err != nil {
		w.log.Error("TransferVerify error", "err", err, "bondId", flow.Key.BondId.Hex())
		return false
	}

	flow.Reason = bondReason

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.LiquidityBondResult, Content: flow}
	return w.submitMessage(result)
}

func (w *writer) processEraPoolUpdated(m *core.Message) bool {
	w.log.Trace("processEraPoolUpdate", "source", m.Source, "dest", m.Destination, "content", m.Content)
	e, ok := m.Content.(core.EvtEraPoolUpdated)
	if !ok {
		w.log.Debug("EvtEraPoolUpdated cast err", "msg", m)
		w.printContentError(m)
		return false
	}
	poolAddrHexStr := hex.EncodeToString(e.Pool)
	subClient, exist := w.conn.subClients[poolAddrHexStr]
	if !exist {
		w.log.Debug("processEraPoolUpdated pool not exist")
		w.printContentError(m)
		return false
	}
	var addrValidatorTestnetAteam, _ = types.ValAddressFromBech32("cosmosvaloper105gvcjgs6s4j5ws9srckx0drt4x8cwgywplh7p")
	poolAddr, err := types.AccAddressFromHex(poolAddrHexStr)
	if err != nil {
		w.log.Debug("accAddressFromHex", "err", err)
		w.printContentError(m)
		return false
	}
	client := subClient.GetRpcClient()

	unSignedTx, err := client.GenMultiSigRawDelegateTx(
		poolAddr,
		addrValidatorTestnetAteam,
		types.NewCoin(client.GetDenom(), types.NewInt(100)))

	if err != nil {
		w.log.Debug("GenMultiSigRawDelegateTx", "err", err)
		w.printContentError(m)
		return false
	}

	return true
}

func (w *writer) printContentError(m *core.Message) {
	w.log.Error("msg resolve failed", "source", m.Source, "dest", m.Destination, "reason", m.Reason)
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (w *writer) submitMessage(m *core.Message) bool {
	err := w.router.Send(m)
	if err != nil {
		w.log.Error("failed to process event", "err", err)
		return false
	}

	return true
}
