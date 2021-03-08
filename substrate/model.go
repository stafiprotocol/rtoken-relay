package substrate

import (
	"errors"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"

	scalecodec "github.com/itering/scale.go"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
)

const (
	TransferModuleName = "Balances"
	TransferKeepAlive  = "transfer_keep_alive"
	Transfer = "transfer"

	ParamDest     = "dest"
	ParamDestType = "Address"

	ParamValue     = "value"
	ParamValueType = "Compact<Balance>"
)

var (
	TerminatedError = errors.New("terminated")
)

type ChainEvent struct {
	ModuleId string                  `json:"module_id" `
	EventId  string                  `json:"event_id" `
	Params   []scalecodec.EventParam `json:"params"`
}

type FullSubClient struct {
	Sc         *SarpcClient
	Gc         *GsrpcClient
	Keys       []*signature.KeyringPair
	SubClients map[*signature.KeyringPair]*GsrpcClient
}

type StakingLedger struct {
	Stash  types.AccountID
	Total  types.UCompact
	Active types.UCompact
}

/// Hex accountIds
var (
	Nominators = []string{}
)
