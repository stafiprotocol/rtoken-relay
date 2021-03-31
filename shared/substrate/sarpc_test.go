package substrate

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
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
	//sc, err := NewSarpcClient("wss://stafi-seiya.stafi.io", stafiTypesFile, tlog)
	//sc, err := NewSarpcClient("wss://mainnet-rpc.stafi.io", stafiTypesFile, tlog)
	//sc, err := NewSarpcClient("wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, tlog)
	assert.NoError(t, err)

	for i := 1; i < 1000; i++ {
		evts, err := sc.GetEvents(uint64(i))
		assert.NoError(t, err)
		for _, evt := range evts {
			fmt.Println(evt.ModuleId)
			fmt.Println(evt.EventId)
			fmt.Println(evt.Params)
		}
	}
}

func TestSarpcClient_GetExtrinsics(t *testing.T) {
	//sc, err := NewSarpcClient("wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	sc, err := NewSarpcClient(ChainTypeStafi, "wss://stafi-seiya.stafi.io", stafiTypesFile, tlog)
	assert.NoError(t, err)

	/// polkadot-test
	//exts, err := sc.GetExtrinsics("0x21b81342a6e31d1bb9d9c74e11483fecc7e9374a0dc1b1f254978808025f477e")
	//exts, err := sc.GetExtrinsics("0x3d55fb40d3ac4f96373f5d2d9860154145c09df9b5b83a88062014cea0da5ad3")

	/// stafi transfer_keep_alive
	exts, err := sc.GetExtrinsics("0x8431e885f1e4b799cc2a86962e109bd8cc6d4070fc3ee1787562a9ba83ed5da4")

	assert.NoError(t, err)
	for _, ext := range exts {
		fmt.Println("exthash", ext.ExtrinsicHash)
		fmt.Println("moduleName", ext.CallModuleName)
		fmt.Println("methodName", ext.CallName)
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
	sc, err := NewSarpcClient(ChainTypePolkadot, "wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	//sc, err := NewSarpcClient("wss://stafi-seiya.stafi.io", stafiTypesFile, tlog)
	assert.NoError(t, err)

	/// polkadot-test
	exts, err := sc.GetExtrinsics("0x9e4810599ef1ad3a59c916c0bbb16e06252d868a30ee6b3707cb37e775b854ac")
	assert.NoError(t, err)
	//exts, err := sc.GetExtrinsics("0x3d55fb40d3ac4f96373f5d2d9860154145c09df9b5b83a88062014cea0da5ad3")

	/// stafi transfer_keep_alive
	//exts, err := sc.GetExtrinsics("0x8431e885f1e4b799cc2a86962e109bd8cc6d4070fc3ee1787562a9ba83ed5da4")

	for _, ext := range exts {
		fmt.Println("exthash", ext.ExtrinsicHash)
		fmt.Println("moduleName", ext.CallModuleName)
		fmt.Println("methodName", ext.CallName)
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
	sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, tlog)
	assert.NoError(t, err)

	exts, err := sc.GetExtrinsics("0x8431e885f1e4b799cc2a86962e109bd8cc6d4070fc3ee1787562a9ba83ed5da4")
	assert.NoError(t, err)
	for _, ext := range exts {
		fmt.Println("exthash", ext.ExtrinsicHash)
		fmt.Println("moduleName", ext.CallModuleName)
		fmt.Println("methodName", ext.CallName)
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
	sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, tlog)
	assert.NoError(t, err)

	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, AliceKey, tlog, stop)
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

		d, err := submodel.EventNewMultisigData(evt)
		assert.NoError(t, err)
		fmt.Println("who", hexutil.Encode(d.Who[:]))
		fmt.Println("id", hexutil.Encode(d.ID[:]))
		fmt.Println("hash", hexutil.Encode(d.CallHash[:]))

		mul := new(submodel.Multisig)
		exist, err := gc.QueryStorage(config.MultisigModuleId, config.StorageMultisigs, d.ID[:], d.CallHash[:], mul)
		assert.NoError(t, err)
		assert.True(t, exist)
		fmt.Println(mul)
	}
}

func TestEventMultisigExecuted(t *testing.T) {
	sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, tlog)
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
	sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, tlog)
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
		a, err := submodel.EraPoolUpdatedData(evt)
		assert.NoError(t, err)
		fmt.Println(a)
	}
}
