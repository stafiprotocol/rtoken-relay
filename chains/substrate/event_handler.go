package substrate

import (
	"fmt"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
)

type eventName string
type eventHandler func(interface{}, log15.Logger) (*core.Message, error)

type eventHandlerSubscriptions struct {
	name    eventName
	handler eventHandler
}

const (
	LiquidityBond  = eventName(config.LiquidityBondEventId)
	EraPoolUpdated = eventName(config.EraPoolUpdatedEventId)
	BondReport     = eventName(config.BondReportEventId)
	WithdrawUnbond = eventName(config.WithdrawUnbondEventId)

	NewMultisig      = eventName(config.NewMultisigEventId)
	MultisigExecuted = eventName(config.MultisigExecutedEventId)
)

var MainSubscriptions = []eventHandlerSubscriptions{
	{LiquidityBond, liquidityBondHandler},
	{EraPoolUpdated, eraPoolUpdatedHandler},
	{BondReport, bondReportHandler},
	{WithdrawUnbond, withdrawUnbondHandler},
}

var OtherSubscriptions = []eventHandlerSubscriptions{
	{NewMultisig, newMultisigHandler},
	{MultisigExecuted, multisigExecutedHandler},
}

func liquidityBondHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*core.BondFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast bondflow")
	}

	return &core.Message{Destination: d.Key.Rsymbol, Reason: core.LiquidityBond, Content: d}, nil
}

func eraPoolUpdatedHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*core.MultisigFlow)
	if !ok {
		return nil, fmt.Errorf("eraPoolUpdatedHandler: failed to cast MultisigFlow")
	}

	flow, ok := d.HeadFlow.(*core.EraPoolUpdatedFlow)
	if !ok {
		return nil, fmt.Errorf("eraPoolUpdatedHandler: failed to cast EraPoolUpdatedFlow")
	}

	return &core.Message{Destination: flow.Snap.Rsymbol, Reason: core.EraPoolUpdated, Content: d}, nil
}

func bondReportHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*core.BondReportFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast bond informChain")
	}

	return &core.Message{Destination: d.Rsymbol, Reason: core.BondReportEvent, Content: d}, nil
}

func withdrawUnbondHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*core.WithdrawUnbondFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast bond informChain")
	}

	return &core.Message{Destination: d.Rsymbol, Reason: core.BondReportEvent, Content: d}, nil
}

func newMultisigHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*core.MultisigFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast newMultisig")
	}

	return &core.Message{Reason: core.NewMultisig, Content: d}, nil
}

func multisigExecutedHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*core.MultisigFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast multisigExecuted")
	}

	return &core.Message{Reason: core.MultisigExecuted, Content: d}, nil
}
