package stafi

import (
	"fmt"

	scale "github.com/itering/scale.go"
	"github.com/itering/scale.go/utiles"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/blake2b"
)

type ExtrinsicDecoder struct {
	ScaleDecoder
	ExtrinsicLength     int                    `json:"extrinsic_length"`
	ExtrinsicHash       string                 `json:"extrinsic_hash"`
	VersionInfo         string                 `json:"version_info"`
	ContainsTransaction bool                   `json:"contains_transaction"`
	Address             interface{}            `json:"address"`
	Signature           string                 `json:"signature"`
	SignatureVersion    int                    `json:"signature_version"`
	Nonce               int                    `json:"nonce"`
	Era                 string                 `json:"era"`
	CallIndex           string                 `json:"call_index"`
	Tip                 interface{}            `json:"tip"`
	CallModule          MetadataModules        `json:"call_module"`
	Call                MetadataCalls          `json:"call"`
	Params              []scale.ExtrinsicParam `json:"params"`
	Metadata            *MetadataStruct
}

func (e *ExtrinsicDecoder) Init(data ScaleBytes, option *ScaleDecoderOption) {
	if option == nil || option.Metadata == nil {
		panic("ExtrinsicDecoder option metadata required")
	}
	e.Params = []scale.ExtrinsicParam{}
	e.Metadata = option.Metadata
	e.ScaleDecoder.Init(data, option)
}

func (e *ExtrinsicDecoder) generateHash() string {
	if !e.ContainsTransaction {
		return ""
	}
	var extrinsicData []byte
	if e.ExtrinsicLength > 0 {
		extrinsicData = e.Data.Data
	} else {
		extrinsicLengthType := CompactU32{}
		extrinsicLengthType.Encode(len(e.Data.Data))
		extrinsicData = append(extrinsicLengthType.Data.Data[:], e.Data.Data[:]...)
	}
	checksum, _ := blake2b.New(32, []byte{})
	_, _ = checksum.Write(extrinsicData)
	h := checksum.Sum(nil)
	return utiles.BytesToHex(h)
}

func (e *ExtrinsicDecoder) Process() {
	e.ExtrinsicLength = e.ProcessAndUpdateData("Compact<u32>").(int)
	if e.ExtrinsicLength != e.Data.GetRemainingLength() {
		e.ExtrinsicLength = 0
		e.Data.Reset()
	}

	e.VersionInfo = utiles.BytesToHex(e.NextBytes(1))

	e.ContainsTransaction = utiles.U256(e.VersionInfo).Int64() >= 80

	result := map[string]interface{}{
		"extrinsic_length": e.ExtrinsicLength,
		"version_info":     e.VersionInfo,
	}

	if e.VersionInfo == "01" || e.VersionInfo == "81" {

		if e.ContainsTransaction {
			e.Address = e.ProcessAndUpdateData("Address").(string)
			e.Signature = e.ProcessAndUpdateData("Signature").(string)
			e.Nonce = e.ProcessAndUpdateData("Compact<u32>").(int)
			e.Era = e.ProcessAndUpdateData("Era").(string)
			e.ExtrinsicHash = e.generateHash()
		}
		e.CallIndex = utiles.BytesToHex(e.NextBytes(2))

	} else if e.VersionInfo == "02" || e.VersionInfo == "82" {

		if e.ContainsTransaction {
			e.Address = e.ProcessAndUpdateData("Address").(string)
			e.Signature = e.ProcessAndUpdateData("Signature").(string)
			e.Era = e.ProcessAndUpdateData("Era").(string)
			e.Nonce = int(e.ProcessAndUpdateData("Compact<U64>").(uint64))
			e.Tip = e.ProcessAndUpdateData("Compact<Balance>").(decimal.Decimal)
			e.ExtrinsicHash = e.generateHash()
		}
		e.CallIndex = utiles.BytesToHex(e.NextBytes(2))

	} else if e.VersionInfo == "03" || e.VersionInfo == "83" {

		if e.ContainsTransaction {
			e.Address = e.ProcessAndUpdateData("Address").(string)
			e.Signature = e.ProcessAndUpdateData("Signature").(string)
			e.Era = e.ProcessAndUpdateData("Era").(string)
			e.Nonce = int(e.ProcessAndUpdateData("Compact<U64>").(uint64))
			e.Tip = e.ProcessAndUpdateData("Compact<Balance>").(decimal.Decimal)
			e.ExtrinsicHash = e.generateHash()
		}
		e.CallIndex = utiles.BytesToHex(e.NextBytes(2))

	} else if e.VersionInfo == "04" || e.VersionInfo == "84" {

		if e.ContainsTransaction {

			address := e.ProcessAndUpdateData("Address")
			switch v := address.(type) {
			case string:
				e.Address = v
				result["address_type"] = "AccountId"
			case map[string]interface{}:
				for name, value := range v {
					result["address_type"] = name
					e.Address = value
				}
			}
			e.SignatureVersion = e.ProcessAndUpdateData("U8").(int)
			if e.SignatureVersion == 2 {
				e.Signature = e.ProcessAndUpdateData("EcdsaSignature").(string)
			} else {
				e.Signature = e.ProcessAndUpdateData("Signature").(string)
			}
			e.Era = e.ProcessAndUpdateData("EraExtrinsic").(string)
			e.Nonce = int(e.ProcessAndUpdateData("Compact<U64>").(uint64))
			if e.Metadata.Extrinsic != nil {
				if utiles.SliceIndex("ChargeTransactionPayment", e.Metadata.Extrinsic.SignedExtensions) != -1 {
					e.Tip = e.ProcessAndUpdateData("Compact<Balance>")
				}
			} else {
				e.Tip = e.ProcessAndUpdateData("Compact<Balance>")
			}
			e.ExtrinsicHash = e.generateHash()
		}
		e.CallIndex = utiles.BytesToHex(e.NextBytes(2))
	} else {
		panic(fmt.Sprintf("Extrinsics version %s is not support", e.VersionInfo))
	}
	if e.CallIndex == "" {
		panic("Not find Extrinsic Lookup, please check type registry")
	}

	call, ok := e.Metadata.CallIndex[e.CallIndex]
	if !ok {
		panic(fmt.Sprintf("Not find Extrinsic Lookup %s, please check metadata info", e.CallIndex))
	}
	e.Call = call.Call
	e.CallModule = call.Module

	for _, arg := range e.Call.Args {
		e.Params = append(e.Params, scale.ExtrinsicParam{
			Name:  arg.Name,
			Type:  arg.Type,
			Value: e.ProcessAndUpdateData(arg.Type)})
	}

	if e.ContainsTransaction {
		result["account_id"] = e.Address
		result["signature"] = e.Signature
		result["nonce"] = e.Nonce
		result["era"] = e.Era
		result["extrinsic_hash"] = e.ExtrinsicHash
	}

	if e.CallIndex != "" {
		result["call_code"] = e.CallIndex
		result["call_module_function"] = e.Call.Name
		result["call_module"] = e.CallModule.Name
	}

	result["nonce"] = e.Nonce
	result["era"] = e.Era
	result["tip"] = e.Tip
	result["params"] = e.Params
	e.Value = result
}
