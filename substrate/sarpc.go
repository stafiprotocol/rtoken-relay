package substrate

import (
	"encoding/json"
	"fmt"
	"github.com/ChainSafe/log15"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/websocket"
	"io/ioutil"
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
	typesPath          string
	currentSpecVersion int
	metaDecoder        scalecodec.MetadataDecoder
}

func NewSarpcClient(endpoint, typesPath string, log log15.Logger) (*SarpcClient, error) {
	log.Info("Connecting to substrate chain with Sarpc", "Endpoint", endpoint)

	sc := &SarpcClient{
		endpoint:           endpoint,
		log:                log,
		metaRaw:            "",
		typesPath:          typesPath,
		currentSpecVersion: 0,
		metaDecoder:        scalecodec.MetadataDecoder{},
	}
	websocket.SetEndpoint(sc.endpoint)
	var err error
	types.RuntimeType{}.Reg()
	content, err := ioutil.ReadFile(typesPath)
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
		return err
	}
	r := v.ToRuntimeVersion()
	if r == nil {
		return fmt.Errorf("runtime version nil")
	}

	// metadata raw
	if sc.metaRaw == "" || r.SpecVersion > sc.currentSpecVersion {
		if err := websocket.SendWsRequest(sc.wsconn, v, rpc.StateGetMetadata(wsId, blockHash)); err != nil {
			return err
		}
		metaRaw, err := v.ToString()
		if err != nil {
			return err
		}
		sc.metaRaw = metaRaw
		sc.metaDecoder.Init(utiles.HexToBytes(metaRaw))
		err = sc.metaDecoder.Process()

		types.RuntimeType{}.Reg()
		content, err := ioutil.ReadFile(sc.typesPath)
		if err != nil {
			panic(err)
		}
		types.RegCustomTypes(source.LoadTypeRegistry(content))

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

func (sc *SarpcClient) GetExtrinsics(blockHash string) ([]*scalecodec.ExtrinsicDecoder, error) {
	err := sc.UpdateMeta(blockHash)
	if err != nil {
		return nil, err
	}

	blk, err := sc.GetBlock(blockHash)
	if err != nil {
		return nil, err
	}

	exts := make([]*scalecodec.ExtrinsicDecoder, 0)
	e := new(scalecodec.ExtrinsicDecoder)
	option := types.ScaleDecoderOption{Metadata: &sc.metaDecoder.Metadata, Spec: sc.currentSpecVersion}
	for _, raw := range blk.Extrinsics {
		e.Init(types.ScaleBytes{Data: util.HexToBytes(raw)}, &option)
		e.Process()
		if e.ExtrinsicHash != "" && e.ContainsTransaction {
			exts = append(exts, e)
		}
	}

	return exts, nil
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

func (sc *SarpcClient) GetEvents(blockNum uint64) ([]*ChainEvent, error) {
	blockHash, err := sc.GetBlockHash(blockNum)
	if err != nil {
		return nil, err
	}

	evts, err := sc.GetChainEvents(blockHash)
	if err != nil {
		return nil, err
	}

	return evts, nil
}
