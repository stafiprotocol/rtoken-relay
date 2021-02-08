package substrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ChainSafe/log15"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/websocket"
)

const (
	wsId       = 1
	storageKey = "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7"
)

type SarpcClient struct {
	endpoint           string
	wsconn             websocket.WsConn
	log                log15.Logger
	metaRaw            string
	currentSpecVersion int
	metaDecoder        scalecodec.MetadataDecoder
}

const DefaultTypeFilePath = "../network/stafi.json"

func NewSarpcClient(endpoint string, log log15.Logger) (*SarpcClient, error) {

	log.Info("Connecting to substrate chain with Sarpc", "Endpoint", endpoint)

	sc := &SarpcClient{
		endpoint:           endpoint,
		log:                log,
		metaRaw:            "",
		currentSpecVersion: 0,
		metaDecoder:        scalecodec.MetadataDecoder{},
	}
	websocket.SetEndpoint(sc.endpoint)
	types.RuntimeType{}.Reg()
	content, err := ioutil.ReadFile(DefaultTypeFilePath)
	if err != nil {
		panic(err)
	}
	types.RegCustomTypes(source.LoadTypeRegistry(content))

	var pool *websocket.PoolConn
	if pool, err = websocket.Init(); err == nil {
		defer pool.Close()
		sc.wsconn = pool.Conn
	} else {
		return nil, fmt.Errorf("websocket init error: %s", err)
	}

	return sc, nil
}

func (sc *SarpcClient) UpdateMeta(blockHash string) error {
	v := &rpc.JsonRpcResult{}
	// runtime version
	if err := websocket.SendWsRequest(sc.wsconn, v, rpc.ChainGetRuntimeVersion(wsId, blockHash)); err != nil {
		return fmt.Errorf("websocket get runtime version error: %v", err)
	}
	r := v.ToRuntimeVersion()
	if r == nil {
		return fmt.Errorf("runtime version nil")
	}

	// metadata raw
	if sc.metaRaw == "" || r.SpecVersion > sc.currentSpecVersion {
		if err := websocket.SendWsRequest(sc.wsconn, v, rpc.StateGetMetadata(wsId, blockHash)); err != nil {
			return fmt.Errorf("websocket get metadata error: %v", err)
		}
		metaRaw, err := v.ToString()
		if err != nil {
			return err
		}
		sc.metaRaw = metaRaw
		sc.metaDecoder.Init(utiles.HexToBytes(metaRaw))
		err = sc.metaDecoder.Process()
		if err != nil {
			return err
		}
		sc.currentSpecVersion = r.SpecVersion
	}

	return nil
}

func (sc *SarpcClient) GetBlock(blockHash string) (*rpc.Block, error) {
	v := &rpc.JsonRpcResult{}
	if err := websocket.SendWsRequest(sc.wsconn, v, rpc.ChainGetBlock(wsId, blockHash)); err != nil {
		return nil, err
	}
	rpcBlock := v.ToBlock()
	return &rpcBlock.Block, nil
}

func (sc *SarpcClient) GetExtrinsicHashs(blockHash string) ([]string, error) {
	blk, err := sc.GetBlock(blockHash)
	if err != nil {
		return nil, err
	}

	ehs := make([]string, 0)
	e := scalecodec.ExtrinsicDecoder{}
	option := types.ScaleDecoderOption{Metadata: &sc.metaDecoder.Metadata, Spec: sc.currentSpecVersion}
	for _, raw := range blk.Extrinsics {
		fmt.Println("Raw:", raw)
		e.Init(types.ScaleBytes{Data: util.HexToBytes(raw)}, &option)
		e.Process()
		if e.ExtrinsicHash != "" {
			ehs = append(ehs, e.ExtrinsicHash)
			fmt.Println(e.Address)
			fmt.Println(e.ContainsTransaction)
			fmt.Println(e.CallIndex)
			fmt.Println(e.Params)
			//fmt.Println(e.CallModule)
		}
	}

	return ehs, nil
}

func (sc *SarpcClient) GetEventsByModuleIdAndEventId(blockNum uint64, moduleId, eventId string) ([]*ChainEvent, error) {
	blockHash, err := sc.GetBlockHash(blockNum)
	if err != nil {
		return nil, nil
	}

	evts, err := sc.GetChainEvents(blockHash)
	if err != nil {
		return nil, nil
	}

	wanted := make([]*ChainEvent, 0)

	for _, evt := range evts {
		if evt.ModuleId != moduleId || evt.EventId != eventId {
			continue
		}

		wanted = append(wanted, evt)
	}

	return wanted, nil
}

func (sc *SarpcClient) IsConnected() bool {
	return sc.wsconn.IsConnected()
}

func (sc *SarpcClient) WebsocketReconnect() error {
	if _, _, err := sc.wsconn.ReadMessage(); err != nil {
		return fmt.Errorf("websocket reconnect error: %s", err)
	}
	return nil
}

func (sc *SarpcClient) GetBlockHash(blockNum uint64) (string, error) {
	v := &rpc.JsonRpcResult{}
	if err := websocket.SendWsRequest(sc.wsconn, v, rpc.ChainGetBlockHash(wsId, int(blockNum))); err != nil {
		return "", fmt.Errorf("websocket get block hash error: %v", err)
	}

	blockHash, err := v.ToString()
	if err != nil {
		return "", err
	}

	return blockHash, nil
}

func (sc *SarpcClient) GetChainEvents(blockHash string) ([]*ChainEvent, error) {
	v := &rpc.JsonRpcResult{}
	if err := websocket.SendWsRequest(sc.wsconn, v, rpc.StateGetStorage(wsId, storageKey, blockHash)); err != nil {
		return nil, fmt.Errorf("websocket get event raw error: %v", err)
	}
	eventRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	err = sc.UpdateMeta(blockHash)
	if err != nil {
		return nil, nil
	}

	// parse event raw into []ChainEvent
	e := scalecodec.EventsDecoder{}
	option := types.ScaleDecoderOption{Metadata: &sc.metaDecoder.Metadata}
	e.Init(types.ScaleBytes{Data: util.HexToBytes(eventRaw)}, &option)
	e.Process()

	var events []*ChainEvent
	b, err := json.Marshal(e.Value)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}
