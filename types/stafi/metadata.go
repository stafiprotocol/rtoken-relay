package stafi

type MetadataModules struct {
	Name      string                `json:"name"`
	Prefix    string                `json:"prefix"`
	Storage   []MetadataStorage     `json:"storage"`
	Calls     []MetadataCalls       `json:"calls"`
	Events    []MetadataEvents      `json:"events"`
	Constants []MetadataConstants   `json:"constants,omitempty"`
	Errors    []MetadataModuleError `json:"errors"`
	Index     int                   `json:"index"`
}

type MetadataStorage struct {
	Name     string      `json:"name"`
	Modifier string      `json:"modifier"`
	Type     StorageType `json:"type"`
	Fallback string      `json:"fallback"`
	Docs     []string    `json:"docs"`
	Hasher   string      `json:"hasher,omitempty"`
}

type StorageType struct {
	Origin        string   `json:"origin"`
	PlainType     *string  `json:"plain_type,omitempty"`
	MapType       *MapType `json:"map_type,omitempty"`
	DoubleMapType *MapType `json:"double_map_type,omitempty"`
}

type MapType struct {
	Hasher     string `json:"hasher"`
	Key        string `json:"key"`
	Key2       string `json:"key2,omitempty"`
	Key2Hasher string `json:"key2Hasher,omitempty"`
	Value      string `json:"value"`
	IsLinked   bool   `json:"isLinked"`
}

type MetadataCalls struct {
	Lookup string                       `json:"lookup"`
	Name   string                       `json:"name"`
	Docs   []string                     `json:"docs"`
	Args   []MetadataModuleCallArgument `json:"args"`
}

type MetadataEvents struct {
	Lookup string   `json:"lookup"`
	Name   string   `json:"name"`
	Docs   []string `json:"docs"`
	Args   []string `json:"args"`
}

type MetadataStruct struct {
	MetadataVersion int                   `json:"metadata_version"`
	Metadata        MetadataTag           `json:"metadata"`
	CallIndex       map[string]CallIndex  `json:"call_index"`
	EventIndex      map[string]EventIndex `json:"event_index"`
	Extrinsic       *ExtrinsicMetadata    `json:"extrinsic"`
}

type CallIndex struct {
	Module MetadataModules `json:"module"`
	Call   MetadataCalls   `json:"call"`
}

type EventIndex struct {
	Module MetadataModules `json:"module"`
	Call   MetadataEvents  `json:"call"`
}

type MetadataTag struct {
	Modules []MetadataModules `json:"modules"`
}

type MetadataConstants struct {
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	ConstantsValue string   `json:"constants_value"`
	Docs           []string `json:"docs"`
}
