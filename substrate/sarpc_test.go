package substrate

import (
	"fmt"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/stretchr/testify/assert"
)

var (
	tlog = log15.Root()
)

const (
	stafiTypesFile = "/Users/fwj/Go/stafi/rtoken-relay/network/stafi.json"
	polkaTypesFile = "/Users/fwj/Go/stafi/rtoken-relay/network/polkadot.json"
)

func TestSarpcClient_GetChainEvents(t *testing.T) {
	sc, err := NewSarpcClient("wss://stafi-seiya.stafi.io", stafiTypesFile, tlog)
	//sc, err := NewSarpcClient("wss://mainnet-rpc.stafi.io", stafiTypesFile, tlog)
	//sc, err := NewSarpcClient("wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	assert.NoError(t, err)

	evts, err := sc.GetEvents(1006395)
	assert.NoError(t, err)
	for _, evt := range evts {
		fmt.Println(evt.ModuleId)
		fmt.Println(evt.EventId)
		fmt.Println(evt.Params)
	}
}

func TestSarpcClient_GetExtrinsics(t *testing.T) {
	//sc, err := NewSarpcClient("wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	sc, err := NewSarpcClient("wss://stafi-seiya.stafi.io", stafiTypesFile, tlog)
	assert.NoError(t, err)

	/// polkadot-test
	//exts, err := sc.GetExtrinsics("0x21b81342a6e31d1bb9d9c74e11483fecc7e9374a0dc1b1f254978808025f477e")
	//exts, err := sc.GetExtrinsics("0x3d55fb40d3ac4f96373f5d2d9860154145c09df9b5b83a88062014cea0da5ad3")

	/// stafi transfer_keep_alive
	exts, err := sc.GetExtrinsics("0x8431e885f1e4b799cc2a86962e109bd8cc6d4070fc3ee1787562a9ba83ed5da4")

	assert.NoError(t, err)
	for _, ext := range exts {
		fmt.Println("exthash", ext.ExtrinsicHash)
		fmt.Println("moduleName", ext.CallModule.Name)
		fmt.Println("methodName", ext.Call.Name)
		fmt.Println("address", ext.Address)
		fmt.Println(ext.Params)
		for _, p := range ext.Params {
			if p.Name == ParamDest && p.Type == ParamDestType {
				dest, ok := p.Value.(string)
				fmt.Println("ok", ok)
				fmt.Println(dest)

				/// polkadot-test
				//dest, ok := p.Value.(map[string]interface{})
				//fmt.Println("ok", ok)
				//v, ok := dest["Id"]
				//fmt.Println("ok1", ok)
				//val, ok := v.(string)
				//fmt.Println("ok2", ok)
				//fmt.Println(val)
			}

			fmt.Println("name", p.Name, "value", p.Value, "type", p.Type)
		}
	}
}

func TestSarpcClient_GetExtrinsics1(t *testing.T) {
	sc, err := NewSarpcClient("wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	//sc, err := NewSarpcClient("wss://stafi-seiya.stafi.io", stafiTypesFile, tlog)
	assert.NoError(t, err)

	call, ok := sc.metaDecoder.Metadata.CallIndex["1459"]
	fmt.Println("ok", ok)
	fmt.Println("call", call)

	/// polkadot-test
	exts, err := sc.GetExtrinsics("0x04b257be925ad91ef2655643ad1dc1456283d9dc8e4f252601eece246bae670a")
	//exts, err := sc.GetExtrinsics("0x3d55fb40d3ac4f96373f5d2d9860154145c09df9b5b83a88062014cea0da5ad3")

	/// stafi transfer_keep_alive
	//exts, err := sc.GetExtrinsics("0x8431e885f1e4b799cc2a86962e109bd8cc6d4070fc3ee1787562a9ba83ed5da4")

	assert.NoError(t, err)
	for _, ext := range exts {
		fmt.Println("exthash", ext.ExtrinsicHash)
		fmt.Println("moduleName", ext.CallModule.Name)
		fmt.Println("methodName", ext.Call.Name)
		fmt.Println("address", ext.Address)
		fmt.Println(ext.Params)
		for _, p := range ext.Params {
			if p.Name == ParamDest && p.Type == ParamDestType {
				//dest, ok := p.Value.(string)
				//fmt.Println("ok", ok)
				//fmt.Println(dest)

				// polkadot-test
				dest, ok := p.Value.(map[string]interface{})
				fmt.Println("ok", ok)
				v, ok := dest["Id"]
				fmt.Println("ok1", ok)
				val, ok := v.(string)
				fmt.Println("ok2", ok)
				fmt.Println(val)
			}

			fmt.Println("name", p.Name, "value", p.Value, "type", p.Type)
		}
	}
}
