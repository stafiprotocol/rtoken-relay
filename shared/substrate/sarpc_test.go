package substrate

import (
	"encoding/json"
	"fmt"
	"github.com/stafiprotocol/rtoken-relay/core"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
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
			if p.Name == config.ParamDest && p.Type == config.ParamDestType {
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
			if p.Name == config.ParamDest && p.Type == config.ParamDestType {
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

func TestSarpcClient_GetExtrinsics2(t *testing.T) {
	sc, err := NewSarpcClient("ws://127.0.0.1:9944", stafiTypesFile, tlog)
	assert.NoError(t, err)

	/// stafi transfer_keep_alive
	exts, err := sc.GetExtrinsics("0xb39cc51509f579344d6e634ce885555871be4a5e4bccae129b3e7b348e5e55b9")

	assert.NoError(t, err)
	for _, ext := range exts {
		fmt.Println("exthash", ext.ExtrinsicHash)
		fmt.Println("moduleName", ext.CallModule.Name)
		fmt.Println("methodName", ext.Call.Name)
		fmt.Println("address", ext.Address)
		fmt.Println(ext.Params)
		for _, p := range ext.Params {
			if p.Name == config.ParamDest && p.Type == config.ParamDestType {
				dest, ok := p.Value.(string)
				fmt.Println("ok", ok)
				fmt.Println(dest)
			}

			fmt.Println("name", p.Name, "value", p.Value, "type", p.Type)
		}
	}
}

func TestEncode(t *testing.T) {
	s := "0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965"
	a, _ := hexutil.Decode(s)
	b := types.NewAccountID(a)
	c, err := types.EncodeToBytes(b)
	assert.NoError(t, err)
	assert.Equal(t, s, hexutil.Encode(c))

}

func TestEventNewMultisig(t *testing.T) {
	sc, err := NewSarpcClient("ws://127.0.0.1:9944", stafiTypesFile, tlog)
	assert.NoError(t, err)

	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AliceKey, tlog, stop)
	assert.NoError(t, err)

	evts, err := sc.GetEvents(906)
	assert.NoError(t, err)
	for _, evt := range evts {
		//fmt.Println(i)
		//fmt.Println(evt.ModuleId)
		//fmt.Println(evt.EventId)
		//fmt.Println(evt.Params)
		if evt.ModuleId != config.MultisigModuleId || evt.EventId != config.NewMultisigEventId {
			continue
		}

		d, err := EventNewMultisig(evt)
		assert.NoError(t, err)
		fmt.Println("who", hexutil.Encode(d.Who[:]))
		fmt.Println("id", hexutil.Encode(d.ID[:]))
		fmt.Println("hash", hexutil.Encode(d.CallHash[:]))

		mul := new(core.Multisig)
		exist, err := gc.QueryStorage(config.MultisigModuleId, config.StorageMultisigs, d.ID[:], d.CallHash[:], mul)
		assert.NoError(t, err)
		assert.True(t, exist)
		fmt.Println(mul)
	}
}

func TestEventMultisigExecuted(t *testing.T) {
	sc, err := NewSarpcClient("ws://127.0.0.1:9944", stafiTypesFile, tlog)
	assert.NoError(t, err)

	evts, err := sc.GetEvents(119)
	assert.NoError(t, err)
	for _, evt := range evts {
		if evt.ModuleId != config.MultisigModuleId || evt.EventId != config.MultisigExecutedEventId {
			continue
		}
		p := evt.Params[1].Value
		fmt.Println(p)
		bz, _ := json.Marshal(p)
		tp := new(types.TimePoint)
		json.Unmarshal(bz, tp)
		fmt.Println(tp)

		fmt.Println(evt.Params[4].Type)
		re := evt.Params[4].Value
		result, ok := re.(map[string]interface{})
		fmt.Println(ok)
		_, ok = result["Ok"]
		fmt.Println(ok)

		//
		//json.Unmarshal(bz, result)
		//fmt.Println(result)
	}
}

func TestEraPoolUpdated(t *testing.T) {
	sc, err := NewSarpcClient("ws://127.0.0.1:9944", stafiTypesFile, tlog)
	assert.NoError(t, err)

	evts, err := sc.GetEvents(113)
	assert.NoError(t, err)
	for _, evt := range evts {
		if evt.ModuleId != config.RTokenLedgerModuleId || evt.EventId != config.EraPoolUpdatedEventId {
			continue
		}
		fmt.Println(evt.ModuleId)
		fmt.Println(evt.EventId)
		fmt.Println(evt.Params)
		a, err := EraPoolUpdatedData(evt)
		assert.NoError(t, err)
		fmt.Println(a)
	}
}
