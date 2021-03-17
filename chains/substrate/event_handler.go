package substrate

import (
	"fmt"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
)

type eventId struct {
	symbol core.RSymbol
	name   eventName
}

type eventName string
type eventHandler func(interface{}, log15.Logger) (*core.Message, error)

var (
	LiquidityBond = &eventId{
		symbol: core.RFIS,
		name:   eventName(config.LiquidityBondEventId),
	}

	EraPoolUpdated = &eventId{
		symbol: core.RFIS,
		name:   eventName(config.EraPoolUpdatedEventId),
	}

	NewMultisig = &eventId{
		symbol: core.RDOT,
		name:   eventName(config.NewMultisigEventId),
	}

	MultisigExecuted = &eventId{
		symbol: core.RFIS,
		name:   eventName(config.EraPoolUpdatedEventId),
	}
)

var Subscriptions = []struct {
	eId     *eventId
	handler eventHandler
}{
	{LiquidityBond, liquidityBondHandler},
	{EraPoolUpdated, eraPoolUpdatedHandler},
}

func liquidityBondHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*core.BondFlow)
	if !ok {
		return nil, fmt.Errorf("failed to cast bondflow")
	}

	return &core.Message{Destination: d.Key.Rsymbol, Reason: core.LiquidityBond, Content: d}, nil
}

func eraPoolUpdatedHandler(data interface{}, log log15.Logger) (*core.Message, error) {
	d, ok := data.(*core.EvtEraPoolUpdated)
	if !ok {
		return nil, fmt.Errorf("failed to cast era pool updated")
	}

	return &core.Message{Destination: d.Rsymbol, Reason: core.EraPoolUpdated, Content: d}, nil
}
