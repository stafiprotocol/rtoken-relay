package polkadot

import (
	"fmt"

	"github.com/itering/scale.go/utiles"
)

type EventsDecoder struct {
	Vec
	Metadata *MetadataStruct
}

func (e *EventsDecoder) Init(data ScaleBytes, option *ScaleDecoderOption) {
	if option.Metadata == nil {
		panic("ExtrinsicDecoder option metadata required")
	}
	e.TypeString = "Vec<EventRecord>"
	e.Metadata = option.Metadata
	e.Vec.Init(data, option)
}

type EventParam struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func (e *EventsDecoder) Process() {
	elementCount := e.ProcessAndUpdateData("Compact<u32>").(int)

	er := EventRecord{Metadata: e.Metadata}
	er.Data = e.Data
	er.Spec = e.Spec
	var result []interface{}
	for i := 0; i < elementCount; i++ {
		element := er.Process()
		element["event_idx"] = i
		result = append(result, element)
	}
	e.Value = result
}

type EventRecord struct {
	ScaleDecoder
	Metadata     *MetadataStruct
	Phase        int             `json:"phase"`
	ExtrinsicIdx int             `json:"extrinsic_idx"`
	Type         string          `json:"type"`
	Params       []EventParam    `json:"params"`
	Event        MetadataEvents  `json:"event"`
	EventModule  MetadataModules `json:"event_module"`
	Topic        []string        `json:"topic"`
}

func (e *EventRecord) Process() map[string]interface{} {
	e.Params = []EventParam{}
	e.Topic = []string{}

	e.Phase = e.GetNextU8()

	if e.Phase == 0 {
		e.ExtrinsicIdx = int(e.ProcessAndUpdateData("U32").(uint32))
	}
	e.Type = utiles.BytesToHex(e.NextBytes(2))

	call, ok := e.Metadata.EventIndex[e.Type]
	if !ok {
		panic(fmt.Sprintf("Not find Extrinsic Lookup %s, please check metadata info", e.Type))
	}

	e.Event = call.Call
	e.EventModule = call.Module

	for _, argType := range e.Event.Args {
		e.Module = e.EventModule.Name
		e.Params = append(e.Params, EventParam{Type: argType, Value: e.ProcessAndUpdateData(argType)})
	}

	if e.Metadata.MetadataVersion >= 5 {
		if topic := e.ProcessAndUpdateData("Vec<Hash>"); topic != nil {
			topicValue := topic.([]interface{})
			for _, v := range topicValue {
				e.Topic = append(e.Topic, v.(string))
			}
		}
	}

	return map[string]interface{}{
		"phase":         e.Phase,
		"extrinsic_idx": e.ExtrinsicIdx,
		"type":          e.Type,
		"module_id":     e.EventModule.Name,
		"event_id":      e.Event.Name,
		"params":        e.Params,
		"topic":         e.Topic,
	}

}
