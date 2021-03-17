package substrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/gorilla/websocket"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/pkg/recws"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/util"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client"
)

const (
	wsId       = 1
	storageKey = "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7"
)

type SarpcClient struct {
	endpoint           string
	rec                *recws.RecConn
	log                log15.Logger
	metaRaw            string
	typesPath          string
	currentSpecVersion int
	metaDecoder        *scalecodec.MetadataDecoder
}

func NewSarpcClient(endpoint, typesPath string, log log15.Logger) (*SarpcClient, error) {
	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		return nil, err
	}

	latestHash, err := api.RPC.Chain.GetFinalizedHead()
	if err != nil {
		return nil, err
	}
	log.Info("NewSarpcClient", "latestHash", latestHash.Hex())

	rec := &recws.RecConn{KeepAliveTimeout: 10 * time.Second}
	rec.Dial(endpoint, nil)
	sc := &SarpcClient{
		endpoint:           endpoint,
		rec:                rec,
		log:                log,
		metaRaw:            "",
		typesPath:          typesPath,
		currentSpecVersion: 0,
		metaDecoder:        &scalecodec.MetadataDecoder{},
	}

	err = sc.UpdateMeta(latestHash.Hex())
	if err != nil {
		return nil, err
	}

	return sc, nil
}

func (sc *SarpcClient) SendWsRequest(v interface{}, action []byte) (err error) {
	if err = sc.rec.WriteMessage(websocket.TextMessage, action); err != nil {
		if sc.rec != nil {
			sc.rec.MarkUnusable()
		}
		return fmt.Errorf("websocket send error: %v", err)
	}
	if err = sc.rec.ReadJSON(v); err != nil {
		if sc.rec != nil {
			sc.rec.MarkUnusable()
		}
		return
	}
	return nil
}

func (sc *SarpcClient) UpdateMeta(blockHash string) error {
	v := &rpc.JsonRpcResult{}
	// runtime version
	if err := sc.SendWsRequest(v, rpc.ChainGetRuntimeVersion(wsId, blockHash)); err != nil {
		return err
	}

	r := v.ToRuntimeVersion()
	if r == nil {
		return fmt.Errorf("runtime version nil")
	}

	// metadata raw
	if sc.metaRaw == "" || r.SpecVersion > sc.currentSpecVersion {
		if err := sc.SendWsRequest(v, rpc.StateGetMetadata(wsId, blockHash)); err != nil {
			return err
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
	if err := sc.SendWsRequest(v, rpc.ChainGetBlock(wsId, blockHash)); err != nil {
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
	return sc.rec.IsConnected()
}

func (sc *SarpcClient) WebsocketReconnect() error {
	if _, _, err := sc.rec.ReadMessage(); err != nil {
		return fmt.Errorf("websocket reconnect error: %s", err)
	}
	return nil
}

func (sc *SarpcClient) GetBlockHash(blockNum uint64) (string, error) {
	v := &rpc.JsonRpcResult{}
	if err := sc.SendWsRequest(v, rpc.ChainGetBlockHash(wsId, int(blockNum))); err != nil {
		return "", fmt.Errorf("websocket get block hash error: %v", err)
	}

	blockHash, err := v.ToString()
	if err != nil {
		return "", err
	}

	return blockHash, nil
}

func (sc *SarpcClient) GetChainEvents(blockHash string) ([]*ChainEvent, error) {
	err := sc.UpdateMeta(blockHash)
	if err != nil {
		return nil, err
	}

	v := &rpc.JsonRpcResult{}
	if err := sc.SendWsRequest(v, rpc.StateGetStorage(wsId, storageKey, blockHash)); err != nil {
		return nil, fmt.Errorf("websocket get event raw error: %v", err)
	}
	eventRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	types.RuntimeType{}.Reg()
	content, err := ioutil.ReadFile(sc.typesPath)
	if err != nil {
		panic(err)
	}
	types.RegCustomTypes(source.LoadTypeRegistry(content))

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
