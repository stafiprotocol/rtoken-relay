package cosmos

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/types"
	errType "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	chainBridgeUtils "github.com/stafiprotocol/chainbridge/shared/ethereum"
	substrateTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
	"github.com/stafiprotocol/rtoken-relay/utils"
	"math/big"
	"time"
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
	case core.EraPoolUpdated:
		return w.processEraPoolUpdatedEvt(m)
	case core.SignatureEnough:
		return w.processSignatureEnoughEvt(m)
	case core.BondReportEvent:
		return w.processBondReportEvent(m)
	case core.BondedPools:
		return w.processBondedPools(m)
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

	era, err := w.conn.GetCurrentEra()
	if err != nil {
		w.log.Error("CurrentEra error", "rsymbol", m.Source)
		return false
	}
	snap := flow.Snap

	if snap.Era != era {
		w.log.Warn("era_pool_updated_event of past era, ignored", "current", era, "eventEra", snap.Era, "rsymbol", snap.Rsymbol)
		return true
	}

	//check bond/unbond is needed
	//bond report if no need
	if snap.Bond.Int.Cmp(snap.Unbond.Int) == 0 {
		w.log.Info("EvtEraPoolUpdated bond=unbond, no need to bond/unbond")
		callHash := utils.BlakeTwo256([]byte{})
		mFlow.OpaqueCalls = []*substrate.MultiOpaqueCall{
			&substrate.MultiOpaqueCall{
				CallHash: hexutil.Encode(callHash[:])}}
		return w.informChain(m.Destination, m.Source, mFlow)
	}

	//get subClient of this pool address
	poolAddrHexStr := hex.EncodeToString(snap.Pool)
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
	coin := types.NewCoin(client.GetDenom(), types.NewInt(10))

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
	proposalId := chainBridgeUtils.Hash(unSignedTx)
	proposalIdHexStr := hex.EncodeToString(proposalId[:])
	wrapUnsignedTx := cosmos.WrapUnsignedTx{
		UnsignedTx: unSignedTx,
		SnapshotId: flow.ShotId,
		Hash:       proposalIdHexStr,
		Type:       submodel.OriginalBond}
	subClient.CacheUnsignedTx(proposalIdHexStr, &wrapUnsignedTx)

	param := submodel.SubmitSignatureParams{
		Symbol:     w.conn.symbol,
		Era:        substrateTypes.NewU32(snap.Era),
		Pool:       substrateTypes.NewBytes(snap.Pool),
		TxType:     submodel.OriginalBond,
		ProposalId: substrateTypes.NewBytes(proposalId[:]),
		Signature:  substrateTypes.NewBytes(sigBts),
	}

	w.log.Info("processEraPoolUpdatedEvt gen unsigned Tx",
		"pool address", poolAddr.String(),
		"tx hash", hex.EncodeToString(proposalId[:]))

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

	poolAddr, err := types.AccAddressFromHex(poolAddrHexStr)
	if err != nil {
		w.log.Error("hexPoolAddr cast to cosmos AccAddress failed", "pool hex address", poolAddrHexStr, "err", err)
		return false
	}

	client := poolClient.GetRpcClient()
	signatures := make([][]byte, 0)
	for _, sig := range sigs.Signature {
		signatures = append(signatures, sig)
	}
	proposalIdHexStr := hex.EncodeToString(sigs.ProposalId)
	wrappedUnSignedTx, err := poolClient.GetWrappedUnsignedTx(proposalIdHexStr)
	if err != nil {
		w.log.Error("processSignatureEnoughEvt failed",
			"proposalId", hex.EncodeToString(sigs.ProposalId),
			"err", err)
		return false
	}
	if wrappedUnSignedTx.Type != submodel.OriginalBond && wrappedUnSignedTx.Type != submodel.OriginalClaimRewards {
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

	retry := BlockRetryLimit
	txHashHexStr := hex.EncodeToString(txHash)
	for {
		if retry <= 0 {
			w.log.Error("processSignatureEnoughEvt broadcast tx reach retry limit",
				"pool hex address", poolAddrHexStr)
			break
		}
		//check on chain
		res, err := client.QueryTxByHash(txHashHexStr)
		if err != nil || res.Empty() || res.Code != 0 {
			w.log.Warn(fmt.Sprintf("processSignatureEnoughEvt QueryTxByHash failed. will rebroadcast after %f second",
				BlockRetryInterval.Seconds()),
				"err or res.empty", err)
			retry--
		} else {
			w.log.Info("processSignatureEnoughEvt success",
				"pool hex address", poolAddrHexStr,
				"txHash", txHashHexStr)
			//return true only check on chain

			switch wrappedUnSignedTx.Type {
			case submodel.OriginalBond:
				callHash := utils.BlakeTwo256(sigs.Pool)
				mflow := submodel.MultiEventFlow{
					EventData: &submodel.EraPoolUpdatedFlow{
						ShotId: wrappedUnSignedTx.SnapshotId},
					OpaqueCalls: []*substrate.MultiOpaqueCall{
						&substrate.MultiOpaqueCall{
							CallHash: hexutil.Encode(callHash[:])}},
				}

				poolClient.RemoveUnsignedTx(proposalIdHexStr)
				return w.informChain(m.Destination, m.Source, &mflow)
			case submodel.OriginalClaimRewards:
				height := poolClient.GetHeightByEra(wrappedUnSignedTx.Era)
				delegationsRes, err := client.QueryDelegations(poolAddr, height)
				if err != nil {
					w.log.Error("processSignatureEnoughEvt QueryDelegations failed",
						"pool hex address", poolAddrHexStr,
						"err", err)
					return false
				}
				total := types.NewInt(0)
				for _, dele := range delegationsRes.GetDelegationResponses() {
					total = total.Add(dele.Balance.Amount)
				}

				rewardRes, err := client.QueryDelegationTotalRewards(poolAddr, height)
				if err != nil {
					w.log.Error("processSignatureEnoughEvt QueryDelegationTotalRewards failed",
						"pool hex address", poolAddrHexStr,
						"err", err)
					return false
				}
				rewardTotal := big.NewInt(0)
				rewardDe := rewardRes.Total.AmountOf(client.GetDenom())
				if !rewardDe.IsNil() {
					rewardTotal = rewardTotal.Add(rewardTotal, rewardDe.BigInt())
				}

				total.Add(types.NewIntFromBigInt(rewardTotal))
				f := submodel.BondReportFlow{
					Era:     wrappedUnSignedTx.Era,
					Rsymbol: sigs.Symbol,
					Pool:    sigs.Pool,
					ShotId:  wrappedUnSignedTx.SnapshotId,
					Active:  substrateTypes.NewU128(*total.BigInt())}

				return w.activeReport(m.Destination, m.Source, &f)

			default:
				return true
			}
		}

		//broadcast if not on chain
		_, err = client.BroadcastTx(txBts)
		if err != nil && err != errType.ErrTxInMempoolCache {
			w.log.Warn("processSignatureEnoughEvt BroadcastTx failed",
				"err", err)
		}
		time.Sleep(BlockRetryInterval)
	}

	return false
}

//process bondReportEvent from stafi
//1 query reward on era height
//2 gen (claim reward && delegate) or (claim reward) unsigned tx and cached
//3 send unsigned tx to stafi
func (w *writer) processBondReportEvent(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondReportFlow)
	if !ok {
		w.printContentError(m, errors.New("msg cast to BondReportFlow not ok"))
		return false
	}
	poolAddrHexStr := hex.EncodeToString(flow.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrHexStr)
	if err != nil {
		w.log.Error("processBondReportEvent failed",
			"pool hex address", poolAddrHexStr,
			"error", err)
		return false
	}

	poolAddr, err := types.AccAddressFromHex(poolAddrHexStr)
	if err != nil {
		w.log.Error("hexPoolAddr cast to cosmos AccAddress failed", "pool hex address", poolAddrHexStr, "err", err)
		return false
	}

	height := poolClient.GetHeightByEra(flow.Era)
	client := poolClient.GetRpcClient()

	//get reward
	rewardRes, err := client.QueryDelegationTotalRewards(poolAddr, height)
	if err != nil {
		w.log.Error("QueryDelegationTotalRewards failed",
			"pool hex address", poolAddrHexStr,
			"error", err)
		return false
	}
	rewardAmount := rewardRes.GetTotal().AmountOf(client.GetDenom()).TruncateInt()

	//get balanceAmount
	balanceAmount, err := client.QueryBalance(poolAddr, client.GetDenom(), height)
	if err != nil {
		w.log.Error("QueryBalance failed",
			"pool hex address", poolAddrHexStr,
			"error", err)
		return false
	}

	//check balanceAmount and rewardAmount
	//(1)if balanceAmount>rewardAmount gen withdraw and delegate tx
	//(2)if balanceAmount<rewardAmount gen withdraw tx
	var unSignedTx []byte
	if balanceAmount.Balance.Amount.GT(rewardAmount) {
		unSignedTx, err = client.GenMultiSigRawWithdrawAllRewardThenDeleTx(
			poolAddr,
			height)
	} else {
		unSignedTx, err = client.GenMultiSigRawWithdrawAllRewardTx(
			poolAddr,
			height)
	}
	if err != nil {
		w.log.Error("GenMultiSigRawWithdrawAllRewardThenDeleTx failed",
			"pool address", poolAddr.String(),
			"height", height,
			"err", err)
		return false
	}

	sigBts, err := client.SignMultiSigRawTx(unSignedTx, poolClient.GetSubKey())
	if err != nil {
		w.log.Error("SignMultiSigRawTx failed",
			"pool address", poolAddr.String(),
			"unsignedTx", string(unSignedTx),
			"err", err)
		return false
	}

	//cache unSignedTx
	proposalId := chainBridgeUtils.Hash(unSignedTx)
	proposalIdHexStr := hex.EncodeToString(proposalId[:])
	wrapUnsignedTx := cosmos.WrapUnsignedTx{
		UnsignedTx: unSignedTx,
		Hash:       proposalIdHexStr,
		SnapshotId: flow.ShotId,
		Era:        flow.Era,
		Type:       submodel.OriginalClaimRewards}

	poolClient.CacheUnsignedTx(proposalIdHexStr, &wrapUnsignedTx)

	param := submodel.SubmitSignatureParams{
		Symbol:     w.conn.symbol,
		Era:        substrateTypes.NewU32(flow.Era),
		Pool:       substrateTypes.NewBytes(flow.Pool),
		TxType:     submodel.OriginalClaimRewards,
		ProposalId: substrateTypes.NewBytes(proposalId[:]),
		Signature:  substrateTypes.NewBytes(sigBts),
	}

	w.log.Info("processBondReportEvent gen unsigned claim reward Tx",
		"pool address", poolAddr.String(),
		"unsigned tx hash", hex.EncodeToString(proposalId[:]))

	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.SubmitSignature, Content: param}
	return w.submitMessage(result)
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

func (w *writer) informChain(source, dest core.RSymbol, flow *submodel.MultiEventFlow) bool {
	msg := &core.Message{Source: source, Destination: dest, Reason: core.InformChain, Content: flow}
	return w.submitMessage(msg)
}

func (w *writer) activeReport(source, dest core.RSymbol, flow *submodel.BondReportFlow) bool {
	msg := &core.Message{Source: source, Destination: dest, Reason: core.ActiveReport, Content: flow}
	return w.submitMessage(msg)
}

func (w *writer) processBondedPools(m *core.Message) bool {
	return true
}
