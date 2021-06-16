package solana

import (
	"fmt"

	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	solCommon "github.com/tpkeeper/solana-go-sdk/common"
)

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

func (w *writer) activeReport(source, dest core.RSymbol, flow *submodel.BondReportedFlow) bool {
	msg := &core.Message{Source: source, Destination: dest, Reason: core.ActiveReport, Content: flow}
	return w.submitMessage(msg)
}

type MultisigTxType string

var MultisigTxStakeType = MultisigTxType("stake")
var MultisigTxUnStakeType = MultisigTxType("unstake")
var MultisigTxWithdrawType = MultisigTxType("withdraw")
var MultisigTxTransferType = MultisigTxType("transfer")

func GetMultisigTxAccountPubkey(baseAccount, programID solCommon.PublicKey, txType MultisigTxType, era uint32) (solCommon.PublicKey, string) {
	seed := fmt.Sprintf("multisig2:%s:%d", txType, era)
	return solCommon.CreateWithSeed(baseAccount, seed, programID), seed
}

func GetStakeAccountPubkey(baseAccount solCommon.PublicKey, era uint32) (solCommon.PublicKey, string) {
	seed := fmt.Sprintf("stake:%d", era)
	return solCommon.CreateWithSeed(baseAccount, seed, solCommon.StakeProgramID), seed
}
