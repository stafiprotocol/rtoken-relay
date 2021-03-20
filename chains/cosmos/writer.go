package cosmos

import (
	"encoding/hex"
	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/types"
	substrateTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"time"
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
	case core.SignatureEnough:

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
	mFlow, ok := m.Content.(*core.MultisigFlow)
	if !ok {
		w.log.Debug("EvtEraPoolUpdated cast err", "msg", m)
		w.printContentError(m)
		return false
	}
	e := mFlow.EvtEraPoolUpdated

	//todo for test
	if e.Bond.Int.Cmp(e.Unbond.Int) == 0 {
		w.log.Debug("EvtEraPoolUpdated bond=unbond")
		return true
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

	sigBts, err := client.SignMultiSigRawTx(unSignedTx, client.GetFromName())
	if err != nil {
		w.log.Debug("SignMultiSigRawTx", "err", err)
		w.printContentError(m)
		return false
	}

	w.log.Info("processEraPoolUpdated gen unsigned Tx", "tx", string(unSignedTx))
	param := core.SubmitSignatureParams{
		Symbol:     w.conn.symbol,
		Era:        substrateTypes.NewU32(e.NewEra),
		Pool:       substrateTypes.NewBytes(e.Pool),
		TxType:     core.OriginalTx(core.Bond),
		ProposalId: substrateTypes.NewBytes(unSignedTx),
		Signature:  substrateTypes.NewBytes(sigBts),
	}

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	return w.submitMessage(result)
}

func (w *writer) processSignatureEnough(m *core.Message) bool {
	w.log.Trace("processSignatureEnough", "source", m.Source, "dest", m.Destination, "content", m.Content)
	sigs, ok := m.Content.(*core.SubmitSignatures)
	if !ok {
		w.log.Debug("SubmitSignatures cast err", "msg", m)
		w.printContentError(m)
		return false
	}

	poolAddrHexStr := hex.EncodeToString(sigs.Pool)
	subClient, exist := w.conn.subClients[poolAddrHexStr]
	if !exist {
		w.log.Debug("processSignatureEnough pool not exist")
		w.printContentError(m)
		return false
	}

	client := subClient.GetRpcClient()
	signatures := make([][]byte, 0)
	for _, sig := range sigs.Signature {
		signatures = append(signatures, sig)
	}

	txHash, txBts, err := client.AssembleMultiSigTx(sigs.ProposalId, signatures)
	if err != nil {
		w.log.Debug("processSignatureEnough AssembleMultiSigTx", "err", err)
		w.printContentError(m)
		return false
	}

	retry := BlockRetryLimit
	txHashHexStr := hex.EncodeToString(txHash)
	for {
		if retry <= 0 {
			w.log.Error("processSignatureEnough broadcast tx reach retry limit")
			break
		}
		//check on chain
		res, err := client.QueryTxByHash(txHashHexStr)
		if err != nil || res.Empty() {
			w.log.Debug("processSignatureEnough QueryTxByHash", "err or res.empty", err)
		} else {
			w.log.Info("processSignatureEnough success", "txHash", txHashHexStr)
			//return true only check on chain
			return true
		}

		//broadcast if not on chain
		_, err = client.BroadcastTx(txBts)
		if err != nil {
			w.log.Debug("processSignatureEnough BroadcastTx", "err", err)
		}

		time.Sleep(BlockRetryInterval)
	}

	return false
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
