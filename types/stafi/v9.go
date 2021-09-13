package stafi

import (
	"encoding/json"

	"github.com/huandu/xstrings"
	"github.com/itering/scale.go/utiles"
)

type MetadataV9Decoder struct {
	ScaleDecoder
}

func (m *MetadataV9Decoder) Init(data ScaleBytes, option *ScaleDecoderOption) {
	m.ScaleDecoder.Init(data, option)
}

func (m *MetadataV9Decoder) Process() {
	result := MetadataStruct{
		Metadata: MetadataTag{
			Modules: nil,
		},
	}
	MetadataV9ModuleCall := m.ProcessAndUpdateData("Vec<MetadataV8Module>").([]interface{})
	callModuleIndex := 0
	eventModuleIndex := 0
	result.CallIndex = make(map[string]CallIndex)
	result.EventIndex = make(map[string]EventIndex)
	bm, _ := json.Marshal(MetadataV9ModuleCall)
	var modulesType []MetadataModules
	_ = json.Unmarshal(bm, &modulesType)
	for k, module := range modulesType {
		if module.Calls != nil {
			for callIndex, call := range module.Calls {
				modulesType[k].Calls[callIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(callModuleIndex), 2, "0") + xstrings.RightJustify(utiles.IntToHex(callIndex), 2, "0")
				result.CallIndex[modulesType[k].Calls[callIndex].Lookup] = CallIndex{Module: module, Call: call}
			}
			callModuleIndex++
		}
		if module.Events != nil {
			for eventIndex, event := range module.Events {
				modulesType[k].Events[eventIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(eventModuleIndex), 2, "0") + xstrings.RightJustify(utiles.IntToHex(eventIndex), 2, "0")
				result.EventIndex[modulesType[k].Events[eventIndex].Lookup] = EventIndex{Module: module, Call: event}
			}
			eventModuleIndex++
		}
	}

	result.Metadata.Modules = modulesType
	m.Value = result
}
