package substrate

import (
	"errors"
	scalecodec "github.com/itering/scale.go"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

var (
	TerminatedError        = errors.New("terminated")
	BondEqualToUnbondError = errors.New("BondEqualToUnbondError")
)

type ChainEvent struct {
	ModuleId string                  `json:"module_id" `
	EventId  string                  `json:"event_id" `
	Params   []scalecodec.EventParam `json:"params"`
}

type StakingLedger struct {
	Stash  types.AccountID
	Total  types.UCompact
	Active types.UCompact
}
