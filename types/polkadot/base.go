package polkadot

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/itering/scale.go/utiles"
)

var RuntimeCodecType []string

type ScaleDecoderOption struct {
	Spec        int
	SubType     string
	Module      string
	ValueList   []string
	Metadata    *MetadataStruct
	FixedLength int
}

type TypeMapping struct {
	Names []string
	Types []string
}

type IScaleDecoder interface {
	Init(data ScaleBytes, option *ScaleDecoderOption)
	Process()
	buildStruct()
	NextBytes(int) []byte
	GetNextU8() int
	reset()
}

type ScaleDecoder struct {
	Data        ScaleBytes      `json:"-"`
	TypeString  string          `json:"-"`
	SubType     string          `json:"-"`
	Value       interface{}     `json:"-"`
	RawValue    string          `json:"-"`
	TypeMapping *TypeMapping    `json:"-"`
	Metadata    *MetadataStruct `json:"-"`
	Spec        int             `json:"-"`
	Module      string          `json:"-"`
}

func (s *ScaleDecoder) Init(data ScaleBytes, option *ScaleDecoderOption) {
	if option != nil {
		if option.Metadata != nil {
			s.Metadata = option.Metadata
		}
		if option.SubType != "" {
			s.SubType = option.SubType
		}
		if option.Spec != 0 {
			s.Spec = option.Spec
		}
		if option.Module != "" {
			s.Module = option.Module
		}
	}
	s.Data = data
	s.RawValue = ""
	s.Value = nil
	if s.TypeMapping == nil && s.TypeString != "" {
		s.buildStruct()
	}
}

func (s *ScaleDecoder) Process() {}

func (s *ScaleDecoder) NextBytes(length int) []byte {
	data := s.Data.GetNextBytes(length)
	s.RawValue += utiles.BytesToHex(data)
	return data
}

func (s *ScaleDecoder) GetNextU8() int {
	b := s.NextBytes(1)
	if len(b) > 0 {
		return int(b[0])
	}
	return 0
}

func (s *ScaleDecoder) getNextBool() bool {
	data := s.NextBytes(1)
	return utiles.BytesToHex(data) == "01"
}

func (s *ScaleDecoder) reset() {
	s.Data.Data = []byte{}
	s.Data.Offset = 0
}

func (s *ScaleDecoder) buildStruct() {
	if s.TypeString != "" && string(s.TypeString[0]) == "(" && s.TypeString[len(s.TypeString)-1:] == ")" {

		var names, types []string
		reg := regexp.MustCompile(`\((.*?)\)`)
		typeString := s.TypeString[1 : len(s.TypeString)-1]
		typeParts := reg.FindAllString(typeString, -1)
		for _, part := range typeParts {
			typeString = strings.ReplaceAll(typeString, part, strings.ReplaceAll(part, ",", "#"))
		}

		for k, v := range strings.Split(typeString, ",") {
			types = append(types, strings.ReplaceAll(strings.TrimSpace(v), "#", ","))
			names = append(names, fmt.Sprintf("col%d", k+1))
		}

		s.TypeMapping = &TypeMapping{Names: names, Types: types}
	}
}

func (s *ScaleDecoder) ProcessAndUpdateData(typeString string) interface{} {
	r := RuntimeType{Module: s.Module}

	if TypeRegistry == nil {
		r.Reg()
	}

	class, value, subType := r.DecoderClass(typeString, s.Spec)
	if class == nil {
		panic(fmt.Sprintf("Not found decoder class %s", typeString))
	}

	offsetStart := s.Data.Offset

	// init
	method, exist := class.MethodByName("Init")
	if !exist {
		panic(fmt.Sprintf("%s not implement init function", typeString))
	}
	option := ScaleDecoderOption{SubType: subType, Spec: s.Spec, Metadata: s.Metadata, Module: s.Module}
	method.Func.Call([]reflect.Value{value, reflect.ValueOf(s.Data), reflect.ValueOf(&option)})

	// process
	value.MethodByName("Process").Call(nil)

	s.Data.Offset = int(value.Elem().FieldByName("Data").FieldByName("Offset").Int())
	s.Data.Data = value.Elem().FieldByName("Data").FieldByName("Data").Bytes()
	s.RawValue = utiles.BytesToHex(s.Data.Data[offsetStart:s.Data.Offset])

	return value.Elem().FieldByName("Value").Interface()
}
