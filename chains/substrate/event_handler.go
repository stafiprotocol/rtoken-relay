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

	NewMultisig      = eventName(config.NewMultisigEventId)
	MultisigExecuted = eventName(config.MultisigExecutedEventId)
)

var MainSubscriptions = []eventHandlerSubscriptions{
	{LiquidityBond, liquidityBondHandler},
	{EraPoolUpdated, eraPoolUpdatedHandler},
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
		return nil, fmt.Errorf("failed to cast era pool updated")
	}

	return &core.Message{Destination: d.EvtEraPoolUpdated.Rsymbol, Reason: core.EraPoolUpdated, Content: d}, nil
}

func newMultisigHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*core.MultisigFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast bondflow")
	}

	return &core.Message{Reason: core.NewMultisig, Content: d}, nil
}

func multisigExecutedHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*core.MultisigFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast bondflow")
	}

	return &core.Message{Reason: core.NewMultisig, Content: d}, nil
}
