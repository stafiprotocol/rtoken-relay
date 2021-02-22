package substrate

import (
	"errors"
	scalecodec "github.com/itering/scale.go"
)

const (
	TransferModuleName = "Balances"
	TransferKeepAlive  = "transfer_keep_alive"

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
