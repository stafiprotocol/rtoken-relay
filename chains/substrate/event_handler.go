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
	LiquidityBond     = eventName(config.LiquidityBondEventId)
	EraPoolUpdated    = eventName(config.EraPoolUpdatedEventId)
	BondReported      = eventName(config.BondReportedEventId)
	ActiveReported    = eventName(config.ActiveReportedEventId)
	WithdrawReported  = eventName(config.WithdrawReportedEventId)
	TransferReported  = eventName(config.TransferReportedEventId)
	NominationUpdated = eventName(config.NominationUpdatedEventId)

	NewMultisig      = eventName(config.NewMultisigEventId)
	MultisigExecuted = eventName(config.MultisigExecutedEventId)

	SignatureEnough  = eventName(config.SignaturesEnoughEventId)
	ValidaterUpdated = eventName(config.ValidatorUpdatedEventId)
)

var StafiSubscriptions = []eventHandlerSubscriptions{
	{LiquidityBond, liquidityBondHandler},
	{EraPoolUpdated, eraPoolUpdatedHandler},
	{BondReported, bondReportedHandler},
	{ActiveReported, activeReportedHandler},
	{WithdrawReported, withdrawReportedHandler},
	{TransferReported, transferReportedHandler},
	{NominationUpdated, nominationUpdatedHandler},
	{SignatureEnough, signatureEnoughHandler},
	{ValidaterUpdated, validatorUpdatedHandler},
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

	return &core.Message{Destination: d.Symbol, Reason: core.LiquidityBondEvent, Content: d}, nil
}

func eraPoolUpdatedHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.MultiEventFlow)
	if !ok {
		return nil, fmt.Errorf("eraPoolUpdatedHandler: failed to cast MultisigFlow")
	}

	if d.EventId != config.EraPoolUpdatedEventId {
		return nil, fmt.Errorf("eventId not matched, expected: %s, got: %s", config.EraPoolUpdatedEventId, d.EventId)
	}

	return &core.Message{Destination: d.Symbol, Reason: core.EraPoolUpdatedEvent, Content: d}, nil
}

func bondReportedHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.BondReportedFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast bond informChain")
	}

	return &core.Message{Destination: d.Snap.Symbol, Reason: core.BondReportedEvent, Content: d}, nil
}

func activeReportedHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.MultiEventFlow)
	if !ok {
		return nil, fmt.Errorf("activeReportedHandler: failed to cast MultiEventFlow")
	}

	if d.EventId != config.ActiveReportedEventId {
		return nil, fmt.Errorf("eventId not matched, expected: %s, got: %s", config.ActiveReportedEventId, d.EventId)
	}

	return &core.Message{Destination: d.Symbol, Reason: core.ActiveReportedEvent, Content: d}, nil
}

func withdrawReportedHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.MultiEventFlow)
	if !ok {
		return nil, fmt.Errorf("withdrawReportedHandler: failed to cast MultiEventFlow")
	}

	if d.EventId != config.WithdrawReportedEventId {
		return nil, fmt.Errorf("eventId not matched, expected: %s, got: %s", config.WithdrawReportedEventId, d.EventId)
	}

	return &core.Message{Destination: d.Symbol, Reason: core.WithdrawReportedEvent, Content: d}, nil
}

func transferReportedHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.TransferReportedFlow)
	if !ok {
		return nil, fmt.Errorf("transferReportedHandler: failed to cast TransferReportedFlow")
	}

	return &core.Message{Destination: d.Symbol, Reason: core.TransferReportedEvent, Content: d}, nil
}

func nominationUpdatedHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.MultiEventFlow)
	if !ok {
		return nil, fmt.Errorf("nominationUpdatedHandler: failed to cast MultiEventFlow")
	}

	if d.EventId != config.NominationUpdatedEventId {
		return nil, fmt.Errorf("eventId not matched, expected: %s, got: %s", config.NominationUpdatedEventId, d.EventId)
	}

	return &core.Message{Destination: d.Symbol, Reason: core.NominationUpdatedEvent, Content: d}, nil
}

func validatorUpdatedHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.MultiEventFlow)
	if !ok {
		return nil, fmt.Errorf("validatorUpdatedHandler: failed to cast MultiEventFlow")
	}

	if d.EventId != config.ValidatorUpdatedEventId {
		return nil, fmt.Errorf("eventId not matched, expected: %s, got: %s", config.ValidatorUpdatedEventId, d.EventId)
	}

	return &core.Message{Destination: d.Symbol, Reason: core.ValidatorUpdatedEvent, Content: d}, nil
}

func signatureEnoughHandler(data interface{}) (*core.Message, error) {
	d, ok := data.(*submodel.SubmitSignatures)
	if !ok {
		return nil, fmt.Errorf("failed to cast submitSignatures")
	}

	return &core.Message{Destination: d.Symbol, Reason: core.SignatureEnoughEvent, Content: d}, nil
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
