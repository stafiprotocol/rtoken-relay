package stafi

import (
	"errors"
	"github.com/itering/scale.go/utiles"
)

type MetadataDecoder struct {
	ScaleDecoder
	Version       string         `json:"version"`
	VersionNumber int            `json:"version_number"`
	Metadata      MetadataStruct `json:"metadata"`
	CodecTypes    []string       `json:"codec_types"`
}

func (m *MetadataDecoder) Init(data []byte) {
	sData := ScaleBytes{Data: data}
	m.ScaleDecoder.Init(sData, nil)
}

func (m *MetadataDecoder) Process() error {
	magicBytes := m.NextBytes(4)
	if string(magicBytes) == "meta" {
		metadataVersion := utiles.U256(utiles.BytesToHex(m.Data.Data[m.Data.Offset : m.Data.Offset+1]))
		m.Version = m.ProcessAndUpdateData("MetadataVersion").(string)
		m.Metadata = m.ProcessAndUpdateData(m.Version).(MetadataStruct)
		m.Metadata.MetadataVersion = int(metadataVersion.Int64())
		m.CodecTypes = utiles.UniqueSlice(RuntimeCodecType)
		return nil
	}
	return errors.New("not metadata")

}
