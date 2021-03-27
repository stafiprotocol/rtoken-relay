package substrate

import (
	"fmt"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
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
	TransferBack   = eventName(config.TransferBackEventId)

	NewMultisig      = eventName(config.NewMultisigEventId)
	MultisigExecuted = eventName(config.MultisigExecutedEventId)
)

var MainSubscriptions = []eventHandlerSubscriptions{
	{LiquidityBond, liquidityBondHandler},
	{EraPoolUpdated, eraPoolUpdatedHandler},
	{BondReport, bondReportHandler},
	{WithdrawUnbond, withdrawUnbondHandler},
	{TransferBack, transferBackHandler},
}

var OtherSubscriptions = []eventHandlerSubscriptions{
	{NewMultisig, newMultisigHandler},
	{MultisigExecuted, multisigExecutedHandler},
}

func liquidityBondHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*submodel.BondFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast bondflow")
	}

	return &core.Message{Destination: d.Key.Rsymbol, Reason: core.LiquidityBond, Content: d}, nil
}

func eraPoolUpdatedHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*submodel.MultiEventFlow)
	if !ok {
		return nil, fmt.Errorf("eraPoolUpdatedHandler: failed to cast MultisigFlow")
	}

	if d.EventId != config.EraPoolUpdatedEventId {
		return nil, fmt.Errorf("eventId not matched, expected: %s, got: %s", config.EraPoolUpdatedEventId, d.EventId)
	}

	flow, ok := d.EventData.(*submodel.EraPoolUpdatedFlow)
	if !ok {
		return nil, fmt.Errorf("eraPoolUpdatedHandler: failed to cast EraPoolUpdatedFlow")
	}

	return &core.Message{Destination: flow.Snap.Rsymbol, Reason: core.EraPoolUpdated, Content: d}, nil
}

func bondReportHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*submodel.BondReportFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast bond informChain")
	}

	return &core.Message{Destination: d.Rsymbol, Reason: core.BondReportEvent, Content: d}, nil
}

func withdrawUnbondHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*submodel.MultiEventFlow)
	if !ok {
		return nil, fmt.Errorf("withdrawUnbondHandler: failed to cast MultiEventFlow")
	}

	if d.EventId != config.WithdrawUnbondEventId {
		return nil, fmt.Errorf("eventId not matched, expected: %s, got: %s", config.WithdrawUnbondEventId, d.EventId)
	}

	flow, ok := d.EventData.(*submodel.WithdrawUnbondFlow)
	if !ok {
		return nil, fmt.Errorf("withdrawUnbondHandler: failed to cast WithdrawUnbondFlow")
	}

	return &core.Message{Destination: flow.Rsymbol, Reason: core.WithdrawUnbondEvent, Content: d}, nil
}

func transferBackHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*submodel.MultiEventFlow)
	if !ok {
		return nil, fmt.Errorf("transferBackHandler: failed to cast MultiEventFlow")
	}

	if d.EventId != config.TransferBackEventId {
		return nil, fmt.Errorf("eventId not matched, expected: %s, got: %s", config.TransferBackEventId, d.EventId)
	}

	flow, ok := d.EventData.(*submodel.TransferFlow)
	if !ok {
		return nil, fmt.Errorf("transferBackHandler: failed to cast WithdrawUnbondFlow")
	}

	return &core.Message{Destination: flow.Rsymbol, Reason: core.TransferBackEvent, Content: d}, nil
}

func newMultisigHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*submodel.EventNewMultisig)
	if !ok {
		return nil, fmt.Errorf("failed to cast newMultisig")
	}

	return &core.Message{Reason: core.NewMultisig, Content: d}, nil
}

func multisigExecutedHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*submodel.EventMultisigExecuted)
	if !ok {
		return nil, fmt.Errorf("failed to cast multisigExecuted")
	}

	return &core.Message{Reason: core.MultisigExecuted, Content: d}, nil
}
