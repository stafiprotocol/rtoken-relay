package cosmos

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/types"
	utils "github.com/stafiprotocol/chainbridge/shared/ethereum"
	substrateTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"time"
)

//write to cosmos
type writer struct {
	conn             *Connection
	router           chains.Router
	log              log15.Logger
	sysErr           chan<- error
	cachedUnsignedTx map[string][]byte
}

func NewWriter(conn *Connection, log log15.Logger, sysErr chan<- error) *writer {
	return &writer{
		conn:             conn,
		log:              log,
		sysErr:           sysErr,
		cachedUnsignedTx: make(map[string][]byte),
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
		return w.processSignatureEnough(m)
	default:
		w.log.Warn("message reason unsupported", "reason", m.Reason)
		return false
	}
}

//process LiquidityBond event from stafi
//1 check liquidityBond data  on cosmos chain
//2 return check result to stafi
func (w *writer) processLiquidityBond(m *core.Message) bool {
	flow, ok := m.Content.(*core.BondFlow)
	if !ok {
		w.printContentError(m, errors.New("msg cast to BondFlow not ok"))
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

//process eraPoolUpdate event
//1 gen bond/unbond multiSig unsigned tx and sign it with subKey
//2 send signature to stafi
func (w *writer) processEraPoolUpdated(m *core.Message) bool {
	w.log.Trace("processEraPoolUpdate", "source", m.Source, "dest", m.Destination, "content", m.Content)
	mFlow, ok := m.Content.(*core.MultisigFlow)
	if !ok {
		w.printContentError(m, errors.New("msg cast to MultisigFlow not ok"))
		return false
	}

	e := mFlow.EvtEraPoolUpdated

	//check bond/unbond is needed
	if e.Bond.Int.Cmp(e.Unbond.Int) == 0 {
		w.log.Info("EvtEraPoolUpdated bond=unbond, no need to bond/unbond")
		return true
	}

	//get subClient of this pool address
	poolAddrHexStr := hex.EncodeToString(e.Pool)
	subClient, err := w.conn.GetPoolClient(poolAddrHexStr)
	if err != nil {
		w.log.Error("EraPoolUpdated pool failed",
			"pool hex address", poolAddrHexStr,
			"err", err)
		return false
	}

	poolAddr, err := types.AccAddressFromHex(poolAddrHexStr)
	if err != nil {
		w.log.Error("hexPoolAddr cast to cosmos AccAddress failed", "pool hex address", poolAddrHexStr, "err", err)
		return false
	}

	//todo cosmos validator just for test,will got from stafi or cosmos
	var addrValidatorTestnetAteam, _ = types.ValAddressFromBech32("cosmosvaloper105gvcjgs6s4j5ws9srckx0drt4x8cwgywplh7p")
	client := subClient.GetRpcClient()
	//just for test
	coin := types.NewCoin(client.GetDenom(), types.NewInt(100))

	unSignedTx, err := client.GenMultiSigRawDelegateTx(
		poolAddr,
		addrValidatorTestnetAteam,
		coin)

	if err != nil {
		w.log.Error("GenMultiSigRawDelegateTx failed",
			"pool address", poolAddr.String(),
			"validator address", addrValidatorTestnetAteam.String(),
			"err", err)
		return false
	}

	sigBts, err := client.SignMultiSigRawTx(unSignedTx, subClient.GetSubKey())
	if err != nil {
		w.log.Error("SignMultiSigRawTx failed",
			"pool address", poolAddr.String(),
			"validator address", addrValidatorTestnetAteam.String(),
			"err", err)
		return false
	}

	//cache unSignedTx
	proposalId := utils.Hash(unSignedTx)
	w.cachedUnsignedTx[hex.EncodeToString(proposalId[:])] = unSignedTx
	param := core.SubmitSignatureParams{
		Symbol:     w.conn.symbol,
		Era:        substrateTypes.NewU32(e.NewEra),
		Pool:       substrateTypes.NewBytes(e.Pool),
		TxType:     core.OriginalTx(core.Bond),
		ProposalId: substrateTypes.NewBytes(proposalId[:]),
		Signature:  substrateTypes.NewBytes(sigBts),
	}

	w.log.Info("processEraPoolUpdated gen unsigned Tx",
		"pool address", poolAddr.String(),
		"tx hash", hex.EncodeToString(proposalId[:]))

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	return w.submitMessage(result)
}

//process SignatureEnough event
//1 assemble unsigned tx and signatures
//2 send tx to cosmos until it is confirmed or reach the retry limit
func (w *writer) processSignatureEnough(m *core.Message) bool {
	w.log.Trace("processSignatureEnough", "source", m.Source, "dest", m.Destination, "content", m.Content)
	sigs, ok := m.Content.(*core.SubmitSignatures)
	if !ok {
		w.printContentError(m, errors.New("msg cast to SubmitSignatures not ok"))
		return false
	}

	poolAddrHexStr := hex.EncodeToString(sigs.Pool)
	subClient, err := w.conn.GetPoolClient(poolAddrHexStr)
	if err != nil {
		w.log.Error("processSignatureEnough failed",
			"pool hex address", poolAddrHexStr,
			"error", err)
		return false
	}

	client := subClient.GetRpcClient()
	signatures := make([][]byte, 0)
	for _, sig := range sigs.Signature {
		signatures = append(signatures, sig)
	}

	unSignedTx, exist := w.cachedUnsignedTx[hex.EncodeToString(sigs.ProposalId)]
	if !exist {
		w.log.Error("processSignatureEnough failed",
			"proposalId", hex.EncodeToString(sigs.ProposalId),
			"err", "can`t find unsignedTx in cachedTx")
		return false
	}
	txHash, txBts, err := client.AssembleMultiSigTx(unSignedTx, signatures)
	if err != nil {
		w.log.Error("processSignatureEnough AssembleMultiSigTx failed",
			"pool hex address ", poolAddrHexStr,
			"err", err)
		return false
	}

	retry := BlockRetryLimit
	txHashHexStr := hex.EncodeToString(txHash)
	for {
		if retry <= 0 {
			w.log.Error("processSignatureEnough broadcast tx reach retry limit",
				"pool hex address", poolAddrHexStr)
			break
		}
		//check on chain
		res, err := client.QueryTxByHash(txHashHexStr)
		if err != nil || res.Empty() {
			w.log.Warn(fmt.Sprintf("processSignatureEnough QueryTxByHash failed will retry after %d second", BlockRetryInterval),
				"err or res.empty", err)
		} else {
			w.log.Info("processSignatureEnough success",
				"pool hex address", poolAddrHexStr,
				"txHash", txHashHexStr)
			//return true only check on chain
			return true
		}

		//broadcast if not on chain
		_, err = client.BroadcastTx(txBts)
		if err != nil {
			w.log.Warn(fmt.Sprintf("processSignatureEnough BroadcastTx failed, will retry after %d second", BlockRetryLimit),
				"err", err)
		}
		time.Sleep(BlockRetryInterval)
	}

	return false
}

func (w *writer) printContentError(m *core.Message, err error) {
	w.log.Error("msg resolve failed", "source", m.Source, "dest", m.Destination, "reason", m.Reason, "err", err)
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
