package substrate

import scalecodec "github.com/itering/scale.go"

type ChainEvent struct {
	ModuleId string                  `json:"module_id" `
	EventId  string                  `json:"event_id" `
	Params   []scalecodec.EventParam `json:"params"`
}

//type EventLiquidityBond struct {
//
//}
