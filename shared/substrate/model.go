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

type EraRewardPoints struct {
	Total       uint32
	Individuals []Individual
}

type Individual struct {
	Validator   types.AccountID
	RewardPoint uint32
}

type Exposure struct {
	Total  types.U128
	Own    types.U128
	Others []*IndividualExposure
}

type IndividualExposure struct {
	Who   types.AccountID
	Value types.U128
}

type StakingLedger struct {
	Stash          types.AccountID
	Total          types.UCompact
	Active         types.UCompact
	Unlocking      []*UnlockChunk
	ClaimedRewards []uint32
}

type UnlockChunk struct {
	Value types.U128
	Era   uint32
}

type MultiOpaqueCall struct {
	Extrinsic string
	Opaque    []byte
	CallHash  string
	TimePoint *OptionTimePoint
}
