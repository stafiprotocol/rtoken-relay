package cosmos

import (
	"encoding/hex"
	"errors"
	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	substrateTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

//write to cosmos
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
	case core.BondedPools:
		return w.processBondedPools(m)
	case core.EraPoolUpdated:
		return w.processEraPoolUpdatedEvt(m)
	case core.BondReportEvent:
		return w.processBondReportEvent(m)
	case core.ActiveReportedEvent:
		return w.processActiveReportedEvent(m)
	case core.SignatureEnough:
		return w.processSignatureEnoughEvt(m)
	default:
		w.log.Warn("message reason unsupported", "reason", m.Reason)
		return false
	}
}

//process LiquidityBond event from stafi
//1 check liquidityBond data  on cosmos chain
//2 return check result to stafi
func (w *writer) processLiquidityBond(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondFlow)
	if !ok {
		w.printContentError(m, errors.New("msg cast to BondFlow not ok"))
		return false
	}

	if flow.Reason != submodel.BondReasonDefault {
		w.log.Error("processLiquidityBond receive a message of which reason is not default",
			"bondId", flow.BondId.Hex(),
			"reason", flow.Reason)
		return false
	}

	bondReason, err := w.conn.TransferVerify(flow.Record)
	if err != nil {
		w.log.Error("TransferVerify error", "err", err, "bondId", flow.BondId.Hex())
		return false
	}

	flow.Reason = bondReason

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.LiquidityBondResult, Content: flow}
	return w.submitMessage(result)
}

//process eraPoolUpdate event
//1 gen bond/unbond multiSig unsigned tx and cache it
//2 sign it with subKey
//3 send signature to stafi
func (w *writer) processEraPoolUpdatedEvt(m *core.Message) bool {
	w.log.Trace("processEraPoolUpdate", "source", m.Source, "dest", m.Destination, "content", m.Content)
	mFlow, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m, errors.New("msg cast to MultisigFlow not ok"))
		return false
	}
	flow, ok := mFlow.EventData.(*submodel.EraPoolUpdatedFlow)
	if !ok {
		w.log.Error("processEraPoolUpdated HeadFlow is not EraPoolUpdatedFlow")
		return false
	}

	snap := flow.Snap

	//check bond/unbond is needed
	//bond report if no need
	if snap.Bond.Int.Cmp(snap.Unbond.Int) == 0 {
		w.log.Info("EvtEraPoolUpdated bond=unbond, no need to bond/unbond")
		callHash := utils.BlakeTwo256([]byte{})
		mFlow.OpaqueCalls = []*submodel.MultiOpaqueCall{
			&submodel.MultiOpaqueCall{
				CallHash: hexutil.Encode(callHash[:])}}
		return w.informChain(m.Destination, m.Source, mFlow)
	}

	//get poolClient of this pool address
	poolAddrHexStr := hex.EncodeToString(snap.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrHexStr)
	if err != nil {
		w.log.Error("EraPoolUpdated pool failed",
			"pool hex address", poolAddrHexStr,
			"err", err)
		return false
	}

	poolAddr, err := types.AccAddressFromHex(poolAddrHexStr)
	if err != nil {
		w.log.Error("hexPoolAddr cast to cosmos AccAddress failed",
			"pool hex address", poolAddrHexStr,
			"err", err)
		return false
	}

	//todo cosmos validator just for test,will got from stafi or cosmos
	var addrValidatorTestnetAteam, _ = types.ValAddressFromBech32("cosmosvaloper105gvcjgs6s4j5ws9srckx0drt4x8cwgywplh7p")
	client := poolClient.GetRpcClient()

	unSignedTx, err := GetBondUnbondUnsignedTx(client, snap.Bond, snap.Unbond, poolAddr, addrValidatorTestnetAteam)
	if err != nil {
		w.log.Error("GenMultiSigRawDelegateTx failed",
			"pool address", poolAddr.String(),
			"validator address", addrValidatorTestnetAteam.String(),
			"err", err)
		return false
	}

	//use current seq
	seq, err := client.GetSequence(0, poolAddr)
	if err != nil {
		w.log.Error("GetSequence failed",
			"pool address", poolAddr.String(),
			"err", err)
		return false
	}

	sigBts, err := client.SignMultiSigRawTxWithSeq(seq, unSignedTx, poolClient.GetSubKey())
	if err != nil {
		w.log.Error("SignMultiSigRawTx failed",
			"pool address", poolAddr.String(),
			"validator address", addrValidatorTestnetAteam.String(),
			"err", err)
		return false
	}

	//cache unSignedTx
	proposalId := GetBondUnBondProposalId(flow.ShotId, snap.Bond, snap.Unbond, seq)
	proposalIdHexStr := hex.EncodeToString(proposalId)
	wrapUnsignedTx := cosmos.WrapUnsignedTx{
		UnsignedTx: unSignedTx,
		SnapshotId: flow.ShotId,
		Key:        proposalIdHexStr,
		Type:       submodel.OriginalBond}

	poolClient.CacheUnsignedTx(proposalIdHexStr, &wrapUnsignedTx)

	param := submodel.SubmitSignatureParams{
		Symbol:     w.conn.symbol,
		Era:        substrateTypes.NewU32(snap.Era),
		Pool:       substrateTypes.NewBytes(snap.Pool),
		TxType:     submodel.OriginalBond,
		ProposalId: substrateTypes.NewBytes(proposalId),
		Signature:  substrateTypes.NewBytes(sigBts),
	}

	w.log.Info("processEraPoolUpdatedEvt gen unsigned bond/unbond Tx",
		"pool address", poolAddr.String(),
		"proposalId", hex.EncodeToString(proposalId))

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	return w.submitMessage(result)
}

//process bondReportEvent from stafi
//1 query reward on era height
//2 gen (claim reward && delegate) or (claim reward) unsigned tx and cache it
//3 sign it with subKey
//4 send signature to stafi
func (w *writer) processBondReportEvent(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondReportedFlow)
	if !ok {
		w.printContentError(m, errors.New("msg cast to BondReportFlow not ok"))
		return false
	}
	poolAddrHexStr := hex.EncodeToString(flow.Snap.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrHexStr)
	if err != nil {
		w.log.Error("processBondReportEvent failed",
			"pool hex address", poolAddrHexStr,
			"error", err)
		return false
	}

	poolAddr, err := types.AccAddressFromHex(poolAddrHexStr)
	if err != nil {
		w.log.Error("hexPoolAddr cast to cosmos AccAddress failed",
			"pool hex address", poolAddrHexStr,
			"err", err)
		return false
	}

	height := poolClient.GetHeightByEra(flow.Snap.Era)
	client := poolClient.GetRpcClient()
	unSignedTx, err := GetClaimRewardUnsignedTx(client, poolAddr, height)
	if err != nil {
		w.log.Error("GetClaimRewardUnsignedTx failed",
			"pool address", poolAddr.String(),
			"err", err)
		return false
	}

	//use current seq
	seq, err := client.GetSequence(0, poolAddr)
	if err != nil {
		w.log.Error("GetSequence failed",
			"pool address", poolAddr.String(),
			"err", err)
		return false
	}

	sigBts, err := client.SignMultiSigRawTxWithSeq(seq, unSignedTx, poolClient.GetSubKey())
	if err != nil {
		w.log.Error("SignMultiSigRawTx failed",
			"pool address", poolAddr.String(),
			"unsignedTx", string(unSignedTx),
			"err", err)
		return false
	}

	//cache unSignedTx
	proposalId := GetClaimRewardProposalId(flow.ShotId, uint64(height))
	proposalIdHexStr := hex.EncodeToString(proposalId)
	wrapUnsignedTx := cosmos.WrapUnsignedTx{
		UnsignedTx: unSignedTx,
		Key:        proposalIdHexStr,
		SnapshotId: flow.ShotId,
		Era:        flow.Snap.Era,
		Type:       submodel.OriginalClaimRewards}

	poolClient.CacheUnsignedTx(proposalIdHexStr, &wrapUnsignedTx)

	param := submodel.SubmitSignatureParams{
		Symbol:     w.conn.symbol,
		Era:        substrateTypes.NewU32(flow.Snap.Era),
		Pool:       substrateTypes.NewBytes(flow.Snap.Pool),
		TxType:     submodel.OriginalClaimRewards,
		ProposalId: substrateTypes.NewBytes(proposalId),
		Signature:  substrateTypes.NewBytes(sigBts),
	}

	w.log.Info("processBondReportEvent gen unsigned claim reward Tx",
		"pool address", poolAddr.String(),
		"wrapped unsigned tx key", proposalIdHexStr)

	//send signature to stafi
	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	return w.submitMessage(result)
}

//process TransferBackEvent
//1 gen transfer  unsigned tx and cache it
//2 sign it with subKey
//3 send signature to stafi
func (w *writer) processActiveReportedEvent(m *core.Message) bool {
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m, errors.New("msg cast to MultiEventFlow not ok"))
		return false
	}

	flow, ok := mef.EventData.(*submodel.WithdrawReportedFlow)
	if !ok {
		w.log.Error("processActiveReportedEvent eventData is not TransferFlow")
		return false
	}

	poolAddrHexStr := hex.EncodeToString(flow.Snap.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrHexStr)
	if err != nil {
		w.log.Error("processBondReportEvent failed",
			"pool hex address", poolAddrHexStr,
			"error", err)
		return false
	}

	poolAddr, err := types.AccAddressFromHex(poolAddrHexStr)
	if err != nil {
		w.log.Error("hexPoolAddr cast to cosmos AccAddress failed",
			"pool hex address", poolAddrHexStr,
			"err", err)
		return false
	}
	client := poolClient.GetRpcClient()

	unSignedTx, err := GetTransferUnsignedTx(client, poolAddr, flow.Receives)
	if err != nil && err != ErrNoOutPuts {
		w.log.Error("GetTransferUnsignedTx failed", "pool hex address", poolAddrHexStr, "err", err)
		return false
	}
	if err == ErrNoOutPuts {
		callHash := utils.BlakeTwo256(flow.Snap.Pool)
		mflow := submodel.MultiEventFlow{
			EventData: &submodel.WithdrawReportedFlow{
				ShotId: flow.ShotId},
			OpaqueCalls: []*submodel.MultiOpaqueCall{
				&submodel.MultiOpaqueCall{
					CallHash: hexutil.Encode(callHash[:])}},
		}
		return w.informChain(m.Destination, m.Source, &mflow)
	}

	//use current seq
	seq, err := client.GetSequence(0, poolAddr)
	if err != nil {
		w.log.Error("GetSequence failed",
			"pool address", poolAddr.String(),
			"err", err)
		return false
	}

	sigBts, err := client.SignMultiSigRawTxWithSeq(seq, unSignedTx, poolClient.GetSubKey())
	if err != nil {
		w.log.Error("processActiveReportedEvent SignMultiSigRawTx failed",
			"pool address", poolAddr.String(),
			"unsignedTx", string(unSignedTx),
			"err", err)
		return false
	}

	//cache unSignedTx
	proposalId := GetTransferProposalId(flow.ShotId)
	proposalIdHexStr := hex.EncodeToString(proposalId)
	wrapUnsignedTx := cosmos.WrapUnsignedTx{
		UnsignedTx: unSignedTx,
		Key:        proposalIdHexStr,
		SnapshotId: flow.ShotId,
		Era:        flow.Snap.Era,
		Type:       submodel.OriginalTransfer}

	poolClient.CacheUnsignedTx(proposalIdHexStr, &wrapUnsignedTx)

	param := submodel.SubmitSignatureParams{
		Symbol:     w.conn.symbol,
		Era:        substrateTypes.NewU32(flow.Snap.Era),
		Pool:       substrateTypes.NewBytes(flow.Snap.Pool),
		TxType:     submodel.OriginalTransfer,
		ProposalId: substrateTypes.NewBytes(proposalId),
		Signature:  substrateTypes.NewBytes(sigBts),
	}

	w.log.Info("processActiveReportedEvent gen unsigned transfer Tx",
		"pool address", poolAddr.String(),
		"unsigned tx hash", hex.EncodeToString(proposalId))

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	return w.submitMessage(result)
}

//process SignatureEnough event
//1 assemble unsigned tx and signatures
//2 send tx to cosmos until it is confirmed or reach the retry limit
//3 (1)bondUnbond type: bond report bond result to stafi
//	(2)claimThenDelegate type:report active to stafi
func (w *writer) processSignatureEnoughEvt(m *core.Message) bool {
	w.log.Trace("processSignatureEnoughEvt", "source", m.Source, "dest", m.Destination, "content", m.Content)
	sigs, ok := m.Content.(*submodel.SubmitSignatures)
	if !ok {
		w.printContentError(m, errors.New("msg cast to SubmitSignatures not ok"))
		return false
	}

	era, err := w.conn.GetCurrentEra()
	if err != nil {
		w.log.Error("CurrentEra error", "rsymbol", m.Source)
		return false
	}

	if uint32(sigs.Era) != era {
		w.log.Warn("processSignatureEnoughEvt of past era, ignored", "current", era, "eventEra", sigs.Era,
			"rsymbol", sigs.Symbol)
		return true
	}

	poolAddrHexStr := hex.EncodeToString(sigs.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrHexStr)
	if err != nil {
		w.log.Error("processSignatureEnoughEvt failed",
			"pool hex address", poolAddrHexStr,
			"error", err)
		return false
	}

	client := poolClient.GetRpcClient()
	signatures := make([][]byte, 0)
	for _, sig := range sigs.Signature {
		signatures = append(signatures, sig)
	}
	proposalIdHexStr := hex.EncodeToString(sigs.ProposalId)

	//if cached tx not exist,return false,not rebuild from proposalId
	wrappedUnSignedTx, err := poolClient.GetWrappedUnsignedTx(proposalIdHexStr)
	if err != nil {
		w.log.Warn("processSignatureEnoughEvt GetWrappedUnsignedTx,failed",
			"proposalId", proposalIdHexStr,
			"err", err)
		return false
	}

	if wrappedUnSignedTx.Type != submodel.OriginalBond &&
		wrappedUnSignedTx.Type != submodel.OriginalClaimRewards &&
		wrappedUnSignedTx.Type != submodel.OriginalTransfer {
		w.log.Error("processSignatureEnoughEvt failed,unknown unsigned tx type",
			"proposalId", hex.EncodeToString(sigs.ProposalId),
			"type", wrappedUnSignedTx.Type)
		return false
	}

	txHash, txBts, err := client.AssembleMultiSigTx(wrappedUnSignedTx.UnsignedTx, signatures)
	if err != nil {
		w.log.Error("processSignatureEnoughEvt AssembleMultiSigTx failed",
			"pool hex address ", poolAddrHexStr,
			"err", err)
		return false
	}

	return w.checkAndSend(poolClient, wrappedUnSignedTx, sigs, m, txHash, txBts)
}
