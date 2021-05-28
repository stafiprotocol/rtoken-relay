package polkadot

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"math/big"
	"reflect"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	raw := "1054657374"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("String").(string)
	if r != "Test" {
		t.Errorf("Test String Process fail, decode return %s", r)
	}
}

func TestCompactU64(t *testing.T) {
	raw := "10"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("Compact<U64>").(uint64)
	if r != 4 {
		t.Errorf("Test TestCompactU64 Process fail, decode return %d", r)
	}
}

func TestU32(t *testing.T) {
	raw := "64000000"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("u32").(uint32)
	if r != 100 {
		t.Errorf("Test TestCompactU64 Process fail, expect return 100, decode return %d", r)
	}
}

func TestU16(t *testing.T) {
	raw := "0300"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("u16").(uint16)
	if r != 3 {
		t.Errorf("Test TestU16 Process fail, expect return 3, decode return %d", r)
	}
}

func TestRawBabePreDigest(t *testing.T) {
	raw := "0x02020000008b86750900000000"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	m.ProcessAndUpdateData("RawBabePreDigest")
}

func TestRawBabePreDigestVRF(t *testing.T) {
	raw := "0x030000000099decc0f0000000040a523a6fdd15ef7ffb2956689b828185b4d60cfac789f64d1b6f26257ebbe543349f8ceae602875c705a59b156af586c7cf907df5c8d5b541fa755638e32b07b02bfb5e7549fb88aa1f32da93519c67275e999da1cd58ec168c80b30e5b4d05"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	m.ProcessAndUpdateData("RawBabePreDigest")
}

func TestSet_Process(t *testing.T) {
	types.RuntimeType{}.Reg()
	types.RegCustomTypes(map[string]source.TypeStruct{
		"CustomSet": {
			Type:      "set",
			BitLength: 64,
			ValueList: []string{"Value1", "Value2", "Value3", "Value4", "Value5"},
		},
	})
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes("0x03000000")}, nil)
	r := m.ProcessAndUpdateData("CustomSet")
	if strings.Join(r.([]string), "") != "Value1Value2" {
		t.Errorf("Test TestSet_Process Process fail, decode return %v", r.([]string))
	}
}

// 0x025ed0b2 Compact<Balance>
func TestCompactBalance(t *testing.T) {
	raw := "0x025ed0b2"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	m.ProcessAndUpdateData("Compact<Balance>")
}

// 0xe52d2254c67c430a0000000000000000 Balance
func TestBalance(t *testing.T) {
	raw := "0xe52d2254c67c430a0000000000000000"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	c := m.ProcessAndUpdateData("Balance")
	fmt.Println(c)
}

//
func TestRegistration(t *testing.T) {
	raw := "0x04010000000200a0724e180900000000000000000000000d505552455354414b452d30310e507572655374616b65204c74641b68747470733a2f2f7777772e707572657374616b652e636f6d2f000000000d40707572657374616b65636f"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("Registration<BalanceOf>")
	rb, _ := json.Marshal(r)
	fmt.Println(string(rb))
}

func TestInt(t *testing.T) {
	raw := "0x2efb"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("i16")
	rb, _ := json.Marshal(r)
	fmt.Println(string(rb))
}

func TestBoolArray(t *testing.T) {
	raw := "0x00000100"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("Approvals")
	c := []interface{}{false, false, true, false}
	if !reflect.DeepEqual(c, r.([]interface{})) {
		t.Errorf("Test TestBoolArray Process fail, decode return %v", r)
	}
}

func TestReferendumInfo(t *testing.T) {
	raw := "0x00004e0c00295ce46278975a53b855188482af699f7726fbbeac89cf16a1741c4698dcdbc90080970600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("ReferendumInfo<BlockNumber, Hash, BalanceOf>")
	c := map[string]interface{}{
		"Ongoing": map[string]interface{}{
			"delay":        432000,
			"end":          806400,
			"proposalHash": "0x295ce46278975a53b855188482af699f7726fbbeac89cf16a1741c4698dcdbc9",
			"tally":        map[string]interface{}{"ayes": "0", "nays": "0", "turnout": "0"}, "threshold": "SuperMajorityApprove",
		}}
	if !reflect.DeepEqual(utiles.ToString(c), utiles.ToString(r)) {
		t.Errorf("Test TestReferendumInfo Process fail, decode return %v", r.(map[string]interface{}))
	}
}

func TestEthereumAccountId(t *testing.T) {
	raw := "0x4119b2e6c3cb618f4f0B93ac77f9Beec7ff02887"
	fmt.Println(len(utiles.HexToBytes(raw)))
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("EthereumAccountId")
	if r.(string) != "0x4119b2e6c3Cb618F4f0B93ac77f9BeeC7FF02887" {
		t.Errorf("Test TestEthereumAccountId Process fail, decode return %v", r)
	}
}

func TestRegistrarInfo(t *testing.T) {
	raw := "0x08014c4bf7f93d0a5ed801ef778f8e7ef58201bdd7e33e167faf42a01d439283cb430000000000000000000000000000000000000000000000000112ccb53338ac0da571d3697548346fb5f0b637ac9412f8abbf6d13588be7563200d8c3795800000000000000000000000000000000000000"
	fmt.Println(len(utiles.HexToBytes(raw)))
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("Vec<Option<RegistrarInfo<BalanceOf, AccountId>>>")
	fmt.Println(r)
}

func TestRewardDestinationLatest(t *testing.T) {
	raw := "0x03f8764d575b96b30e095a201d90b6ddaf944d042811846f7a3fe5ffda2a01c045"
	fmt.Println(len(utiles.HexToBytes(raw)))
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("RewardDestination")
	fmt.Println(r)
}

func TestGenericLookupSource(t *testing.T) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(256))
	fmt.Println(bs)
	c := []byte{255, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}
	raw := utiles.BytesToHex(c)
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("GenericLookupSource")
	if r.(string) != "0102030405060708010203040506070801020304050607080102030405060708" {
		t.Errorf("Test TestGenericLookupSource Process fail, decode return %d", r)
	}
	m.Init(types.ScaleBytes{Data: []byte{0xfc, 0, 1}}, nil)
	r = m.ProcessAndUpdateData("GenericLookupSource")
}

func TestBTreeMap(t *testing.T) {
	raw := "0x041c62617a7a696e6745000000"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("BTreeMap<Text,u32>")
	if utiles.ToString(r) != `[{"bazzing":69}]` {
		t.Errorf("Test TestBTreeMap Process fail, decode return %v", utiles.ToString(r))
	}
}

func TestClikeEnum(t *testing.T) {
	raw := "0x45"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	types.RuntimeType{}.Reg()
	types.RegCustomTypes(source.LoadTypeRegistry([]byte(`{"t": {"type": "enum","type_mapping": [["A","42"],["B","69"],["C","255"]]}}`)))
	r := m.ProcessAndUpdateData("t")
	if utiles.ToString(r) != `B` {
		t.Errorf("Test TestClikeEnum Process fail, decode return %v", utiles.ToString(r))
	}
}

func TestNamespaceInt(t *testing.T) {
	raw := "0x2efb"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("Eth::i16")
	if r.(*big.Int).String() != `-1234` {
		t.Errorf("Test TestNamespaceInt Process fail, decode return %v", utiles.ToString(r))
	}
}

func TestBoundedVec(t *testing.T) {
	raw := "0x08014c4bf7f93d0a5ed801ef778f8e7ef58201bdd7e33e167faf42a01d439283cb430000000000000000000000000000000000000000000000000112ccb53338ac0da571d3697548346fb5f0b637ac9412f8abbf6d13588be7563200d8c3795800000000000000000000000000000000000000"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	m.ProcessAndUpdateData("BoundedVec<Option<RegistrarInfo<BalanceOf, AccountId>>,5>")
}

func TestBTreeSet(t *testing.T) {
	raw := "0x1002000000180000001e00000050000000"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	m.ProcessAndUpdateData("BTreeSet<U32>")
}

func TestModuleTypeOverride(t *testing.T) {
	raw := "0x21033e010000830100004403000025010000ec020000bb010000050100002a010000b8000000680300008f010000d702000040010000f00000006b02000014000000bf010000460100006200000044000000b40200004603000022000000760000004303000021030000a3000000cc000000030300002e030000ac020000c5020000f701000029010000ab010000280300001b00000091000000e7010000f8000000930200005e010000f4020000240000004b00000043020000c70100007e020000da0200000d0100000d0200008b0100000d000000d10200007d030000ef0000005b0200005c0100008a020000e3000000b1010000ad000000210200003d0300007b020000e9010000b8020000c30000008b0200003f000000330000005a03000037010000a20200004901000096020000340100007a0000004e0100004a020000d10000003f0200009e0200002e020000e0010000c101000087000000110100009e000000bd000000dc000000a20100004d000000b102000053010000040300006b01000027020000b5010000fd0000001c0000001c020000da00000048010000300300002700000080000000890100005000000017030000bd010000e20100006b03000009030000bb0000007e00000020030000e0020000ee020000730000008f0200000a000000570100000f0100001a000000e801000074020000be01000081010000a6020000f201000058010000570300002c000000b0000000a6010000420000007d0200006802000055020000010300009f020000a901000066000000b6020000a900000077000000c40000002f000000b10000001b010000c60100006c030000d0010000300200003c000000cd020000000100006f01000063020000d60000005e030000b90100005401000095010000010000000b000000140300002d02000071010000fc0200005303000005020000350300005d0300009c000000c902000064030000ea0100006202000010010000090100004e00000018020000dc02000023010000250000004f0100003b010000fb000000f501000020020000730300009202000059030000c1000000fb020000fc010000e20200004c020000"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, &types.ScaleDecoderOption{Module: "parasShared"})
	m.ProcessAndUpdateData("Vec<ValidatorIndex>")
}

func TestAccountInfo(t *testing.T) {
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(`0x0100000001000000`)}, nil)
	utiles.Debug(m.ProcessAndUpdateData("AccountInfo<Index, AccountData>"))
}
