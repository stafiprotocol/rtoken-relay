package solana

import (
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
)

// process LiquidityBond event from stafi
// 1 check liquidityBond data  on solana chain
// 2 return check result to stafi
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

	retry := 0
	var err error
	var bondReason submodel.BondReason
	for {
		if retry > 600 {
			w.log.Error(fmt.Sprintf("TransferVerify reach retry limit, err: %s", err))
			return false
		}
		bondReason, err = w.conn.TransferVerify(flow.Record)
		if err != nil {
			w.log.Warn("TransferVerify error", "err", err, "bondId", flow.BondId.Hex())
			retry++
			time.Sleep(waitTime * 2)
			continue
		}
		break
	}

	flow.Reason = bondReason
	w.log.Info("processLiquidityBond", "bonder", hexutil.Encode(flow.Record.Bonder[:]),
		"bondReason", bondReason, "bondId", flow.BondId.Hex())
	result := &core.Message{Source: m.Destination, Destination: m.Source, Reason: core.LiquidityBondResult, Content: flow}
	return w.submitMessage(result)
}
