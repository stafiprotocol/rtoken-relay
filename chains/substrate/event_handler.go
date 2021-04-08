package substrate

import (
	"fmt"

	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
)

type eventName string
type eventHandler func(interface{}) (*core.Message, error)

type eventHandlerSubscriptions struct {
	name    eventName
	handler eventHandler
}

const (
	LiquidityBond    = eventName(config.LiquidityBondEventId)
	EraPoolUpdated   = eventName(config.EraPoolUpdatedEventId)
	BondReported     = eventName(config.BondReportedEventId)
	ActiveReported   = eventName(config.ActiveReportedEventId)
	WithdrawReported = eventName(config.WithdrawReportedEventId)

	NewMultisig      = eventName(config.NewMultisigEventId)
	MultisigExecuted = eventName(config.MultisigExecutedEventId)
)

var MainSubscriptions = []eventHandlerSubscriptions{
	{LiquidityBond, liquidityBondHandler},
	{EraPoolUpdated, eraPoolUpdatedHandler},
	{BondReported, bondReportedHandler},
	{ActiveReported, activeReportedHandler},
	{WithdrawReported, withdrawReportedHandler},
}

var OtherSubscriptions = []eventHandlerSubscriptions{
	{NewMultisig, newMultisigHandler},
	{MultisigExecuted, multisigExecutedHandler},
}

func liquidityBondHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.BondFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast bondflow")
	}

	return &core.Message{Destination: d.Key.Rsymbol, Reason: core.LiquidityBond, Content: d}, nil
}

func eraPoolUpdatedHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.MultiEventFlow)
	if !ok {
		return nil, fmt.Errorf("eraPoolUpdatedHandler: failed to cast MultisigFlow")
	}

	if d.EventId != config.EraPoolUpdatedEventId {
		return nil, fmt.Errorf("eventId not matched, expected: %s, got: %s", config.EraPoolUpdatedEventId, d.EventId)
	}

	return &core.Message{Destination: d.Rsymbol, Reason: core.EraPoolUpdated, Content: d}, nil
}

func bondReportedHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.BondReportedFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast bond informChain")
	}

	return &core.Message{Destination: d.Snap.Rsymbol, Reason: core.BondReportEvent, Content: d}, nil
}

func activeReportedHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.MultiEventFlow)
	if !ok {
		return nil, fmt.Errorf("activeReportedHandler: failed to cast MultiEventFlow")
	}

	if d.EventId != config.ActiveReportedEventId {
		return nil, fmt.Errorf("eventId not matched, expected: %s, got: %s", config.ActiveReportedEventId, d.EventId)
	}

	return &core.Message{Destination: d.Rsymbol, Reason: core.ActiveReportedEvent, Content: d}, nil
}

func withdrawReportedHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.MultiEventFlow)
	if !ok {
		return nil, fmt.Errorf("withdrawReportedHandler: failed to cast MultiEventFlow")
	}

	if d.EventId != config.WithdrawReportedEventId {
		return nil, fmt.Errorf("eventId not matched, expected: %s, got: %s", config.WithdrawReportedEventId, d.EventId)
	}

	return &core.Message{Destination: d.Rsymbol, Reason: core.WithdrawReportedEvent, Content: d}, nil
}

func newMultisigHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.EventNewMultisig)
	if !ok {
		return nil, fmt.Errorf("failed to cast newMultisig")
	}

	return &core.Message{Reason: core.NewMultisig, Content: d}, nil
}

func multisigExecutedHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.EventMultisigExecuted)
	if !ok {
		return nil, fmt.Errorf("failed to cast multisigExecuted")
	}

	return &core.Message{Reason: core.MultisigExecuted, Content: d}, nil
}
