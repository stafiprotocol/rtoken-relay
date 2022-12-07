package substrate

import (
	"sync"
	"testing"
	"time"

	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
)

var (
	tlog = core.NewLog()
)

const (
	stafiTypesFile  = "../../network/stafi.json"
	polkaTypesFile  = "../../network/polkadot.json"
	kusamaTypesFile = "../../network/kusama.json"
)

func TestSarpcClient_GetChainEvents(t *testing.T) {
	//sc, err := NewSarpcClient("wss://mainnet-rpc.stafi.io", stafiTypesFile, tlog)
	//sc, err := NewSarpcClient("wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	//sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, tlog)
	stop := make(chan int)
	sc, err := NewSarpcClient(ChainTypeStafi, "wss://stafi-seiya.stafi.io", stafiTypesFile, AddressTypeAccountId, AliceKey, tlog, stop)
	// sc, err := NewSarpcClient(ChainTypePolkadot, "wss://kusama-rpc.polkadot.io", polkaTypesFile, AddressTypeMultiAddress, AliceKey, tlog, stop)
	// sc, err := NewSarpcClient(ChainTypePolkadot, "wss://kusama-rpc.stafi.io", kusamaTypesFile, AddressTypeMultiAddress, AliceKey, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}
	//evt, err := sc.GetEvents(7112781)
	//assert.NoError(t, err)
	//for _, e := range evt {
	//	t.Log(e.EventId)
	//}

	wg := sync.WaitGroup{}
	for i := 746800; i <= 746824; i++ {
		t.Log("i", i)
		go func() {
			wg.Add(1)
			_, err := sc.GetEvents(uint64(i))
			if err != nil {
				time.Sleep(time.Second)
				t.Log(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

}

func TestSarpcClient_GetChainEventNominationUpdated(t *testing.T) {
	stop := make(chan int)
	// sc, err := NewSarpcClient(ChainTypeStafi, "wss://stafi-seiya.stafi.io", stafiTypesFile, AddressTypeAccountId, AliceKey, tlog, stop)
	sc, err := NewSarpcClient(ChainTypeStafi, "wss://mainnet-rpc.stafi.io", stafiTypesFile, AddressTypeAccountId, AliceKey, tlog, stop)
	// sc, err := NewSarpcClient(ChainTypePolkadot,"wss://polkadot-test-rpc.stafi.io", polkaTypesFile, AddressTypeAccountId, AliceKey, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	symbz, err := types.EncodeToBytes(core.RKSM)
	if err != nil {
		t.Fatal(err)
	}
	bondedPools := make([]types.Bytes, 0)
	exist, err := sc.QueryStorage(config.RTokenLedgerModuleId, config.StorageBondedPools, symbz, nil, &bondedPools)
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		t.Fatal("bonded pools not extis")
	}

	t.Log(bondedPools)

	evts, err := sc.GetEvents(11561482)
	if err != nil {
		t.Fatal(err)
	}
	for _, evt := range evts {
		if evt.EventId != config.NominationUpdatedEventId {
			continue
		}
		flow, err := submodel.EventNominationUpdated(evt)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(flow)
		//t.Log(evt.ModuleId)
		//t.Log(evt.EventId)
		//t.Log(evt.Params)
	}
}

func TestSarpcClient_GetExtrinsics1(t *testing.T) {
	//sc, err := NewSarpcClient(ChainTypePolkadot, "wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	//sc, err := NewSarpcClient("wss://stafi-seiya.stafi.io", stafiTypesFile, tlog)
	stop := make(chan int)
	sc, err := NewSarpcClient(ChainTypePolkadot, "wss://rpc.polkadot.io", polkaTypesFile, AddressTypeMultiAddress, AliceKey, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	for i := 7411010; i >= 7311010; i-- {
		if i%10 == 0 {
			t.Log("i", i)
		}

		bh, err := sc.GetBlockHash(uint64(i))
		exts, err := sc.GetExtrinsics(bh)
		if err != nil {
			t.Fatal(err)
		}

		for _, ext := range exts {
			t.Log("exthash", ext.ExtrinsicHash)
			t.Log("moduleName", ext.CallModuleName)
			t.Log("methodName", ext.CallName)
			t.Log("address", ext.Address)
			t.Log(ext.Params)
		}
	}
}

func TestSarpcClient_GetExtrinsics2(t *testing.T) {
	stop := make(chan int)
	sc, err := NewSarpcClient(ChainTypePolkadot, "wss://kusama-rpc.polkadot.io", polkaTypesFile, AddressTypeMultiAddress, AliceKey, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	exts, err := sc.GetExtrinsics("0x6157da60a188b3f31d250afe5acb2da786417fec00973f1c7f863504fbca4642")
	if err != nil {
		t.Fatal(err)
	}

	for _, ext := range exts {
		t.Log("exthash", ext.ExtrinsicHash)
		t.Log("moduleName", ext.CallModuleName)
		t.Log("methodName", ext.CallName)
		t.Log("address", ext.Address)
		t.Log(ext.Params)
		for _, p := range ext.Params {
			if p.Name == config.ParamDest && p.Type == config.ParamDestType {
				//dest, ok := p.Value.(string)
				//fmt.Println("ok", ok)
				//fmt.Println(dest)

				// polkadot-test
				dest, ok := p.Value.(map[string]interface{})
				t.Log("ok", ok)
				v, ok := dest["Id"]
				t.Log("ok1", ok)
				val, ok := v.(string)
				t.Log("ok2", ok)
				t.Log(val)
			}

			t.Log("name", p.Name, "value", p.Value, "type", p.Type)
		}
	}
}
