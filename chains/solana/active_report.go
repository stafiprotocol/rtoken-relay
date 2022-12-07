package solana

import (
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mr-tron/base58"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

// 1 merge and withdraw
// 3 withdraw report to stafi
func (w *writer) processActiveReportedEvent(m *core.Message) bool {
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m, errors.New("msg cast to MultiEventFlow not ok"))
		return false
	}

	flow, ok := mef.EventData.(*submodel.ActiveReportedFlow)
	if !ok {
		w.log.Error("processActiveReportedEvent eventData is not TransferFlow")
		return false
	}

	poolAddrBase58Str := base58.Encode(flow.Snap.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrBase58Str)
	if err != nil {
		w.log.Error("processActiveReportedEvent failed",
			"pool address", poolAddrBase58Str,
			"error", err)
		return false
	}

	ok = w.MergeAndWithdraw(poolClient, poolAddrBase58Str, flow.Snap.Era, flow.ShotId, flow.Snap.Pool)
	if !ok {
		return false
	}

	callHash := utils.BlakeTwo256(flow.Snap.Pool)
	mef.RunTimeCalls = []*submodel.RunTimeCall{
		{CallHash: hexutil.Encode(callHash[:])}}

	return w.informChain(m.Destination, m.Source, mef)
}
