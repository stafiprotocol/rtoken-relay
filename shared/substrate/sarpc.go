package substrate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/gorilla/websocket"
	scale "github.com/itering/scale.go"
	"github.com/itering/scale.go/source"
	scaleTypes "github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/pkg/recws"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/util"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client"
	gsrpcConfig "github.com/stafiprotocol/go-substrate-rpc-client/config"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	wbskt "github.com/stafiprotocol/rtoken-relay/shared/substrate/websocket"
	"github.com/stafiprotocol/rtoken-relay/types/polkadot"
	"github.com/stafiprotocol/rtoken-relay/types/stafi"
)

const (
	wsId       = 1
	storageKey = "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7"
)

type SarpcClient struct {
	endpoint    string
	addressType string
	api         *gsrpc.SubstrateAPI
	key         *signature.KeyringPair
	genesisHash types.Hash
	stop        <-chan int

	wsPool    wbskt.Pool
	log       log15.Logger
	chainType string
	typesPath string

	stafiMetaDecoderMap map[int]*stafi.MetadataDecoder
	polkaMetaDecoderMap map[int]*scale.MetadataDecoder
	sync.RWMutex
}

var (
	metaDecoders = map[string]interface{}{
		ChainTypeStafi:    &stafi.MetadataDecoder{},
		ChainTypePolkadot: &polkadot.MetadataDecoder{},
	}
)

func NewSarpcClient(chainType, endpoint, typesPath, addressType string, key *signature.KeyringPair, log log15.Logger, stop <-chan int) (*SarpcClient, error) {
	log.Info("Connecting to substrate chain with sarpc", "endpoint", endpoint)

	if addressType != AddressTypeAccountId && addressType != AddressTypeMultiAddress {
		return nil, errors.New("addressType not supported")
	}

	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		return nil, err
	}

	gsrpcConfig.SetSubscribeTimeout(2 * time.Minute)
	latestHash, err := api.RPC.Chain.GetFinalizedHead()
	if err != nil {
		return nil, err
	}
	log.Info("NewSarpcClient", "latestHash", latestHash.Hex())

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}

	md := metaDecoders[chainType]
	if md == nil {
		return nil, errors.New("chainType not supported")
	}

	sc := &SarpcClient{
		endpoint:            endpoint,
		chainType:           chainType,
		addressType:         addressType,
		api:                 api,
		key:                 key,
		genesisHash:         genesisHash,
		stop:                stop,
		wsPool:              nil,
		log:                 log,
		typesPath:           typesPath,
		stafiMetaDecoderMap: make(map[int]*stafi.MetadataDecoder),
		polkaMetaDecoderMap: make(map[int]*scale.MetadataDecoder),
	}

	sc.regCustomTypes()

	if err != nil {
		return nil, err
	}

	return sc, nil
}

func (s *SarpcClient) getStafiMetaDecoder(blockHash string) (*stafi.MetadataDecoder, error) {
	v := &rpc.JsonRpcResult{}
	// runtime version
	if err := s.sendWsRequest(nil, v, rpc.ChainGetRuntimeVersion(wsId, blockHash)); err != nil {
		return nil, err
	}

	r := v.ToRuntimeVersion()
	if r == nil {
		return nil, fmt.Errorf("runtime version nil")
	}
	s.RLock()
	if decoder, exist := s.stafiMetaDecoderMap[r.SpecVersion]; exist {
		s.RUnlock()
		return decoder, nil
	}
	s.RUnlock()

	// check metadata need update, maybe  get ahead hash
	if err := s.sendWsRequest(nil, v, rpc.StateGetMetadata(wsId, blockHash)); err != nil {
		return nil, err
	}
	metaRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	md := &stafi.MetadataDecoder{}
	md.Init(utiles.HexToBytes(metaRaw))
	if err := md.Process(); err != nil {
		return nil, err
	}
	s.Lock()
	s.stafiMetaDecoderMap[r.SpecVersion] = md
	s.Unlock()
	return md, nil

}

func (s *SarpcClient) getPolkaMetaDecoder(blockHash string) (*scale.MetadataDecoder, error) {
	v := &rpc.JsonRpcResult{}
	// runtime version
	if err := s.sendWsRequest(nil, v, rpc.ChainGetRuntimeVersion(wsId, blockHash)); err != nil {
		return nil, err
	}

	r := v.ToRuntimeVersion()
	if r == nil {
		return nil, fmt.Errorf("runtime version nil")
	}
	s.RLock()
	if decoder, exist := s.polkaMetaDecoderMap[r.SpecVersion]; exist {
		s.RUnlock()
		return decoder, nil
	}
	s.RUnlock()

	// check metadata need update, maybe  get ahead hash
	if err := s.sendWsRequest(nil, v, rpc.StateGetMetadata(wsId, blockHash)); err != nil {
		return nil, err
	}
	metaRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	md := &scale.MetadataDecoder{}
	md.Init(utiles.HexToBytes(metaRaw))
	if err := md.Process(); err != nil {
		return nil, err
	}
	s.Lock()
	s.polkaMetaDecoderMap[r.SpecVersion] = md
	s.Unlock()
	return md, nil

}

func (sc *SarpcClient) regCustomTypes() {
	content, err := ioutil.ReadFile(sc.typesPath)
	if err != nil {
		panic(err)
	}

	switch sc.chainType {
	case ChainTypeStafi:
		stafi.RuntimeType{}.Reg()
		stafi.RegCustomTypes(source.LoadTypeRegistry(content))
	case ChainTypePolkadot:
		polkadot.RuntimeType{}.Reg()
		polkadot.RegCustomTypes(source.LoadTypeRegistry(content))
	default:
		panic("chainType not supported")
	}
}

func (sc *SarpcClient) initial() (*wbskt.PoolConn, error) {
	var err error
	if sc.wsPool == nil {
		factory := func() (*recws.RecConn, error) {
			SubscribeConn := &recws.RecConn{KeepAliveTimeout: 10 * time.Second}
			SubscribeConn.Dial(sc.endpoint, nil)
			return SubscribeConn, err
		}
		if sc.wsPool, err = wbskt.NewChannelPool(1, 25, factory); err != nil {
			fmt.Println("NewChannelPool", err)
		}
	}
	if err != nil {
		return nil, err
	}
	conn, err := sc.wsPool.Get()
	return conn, err
}

func (sc *SarpcClient) sendWsRequest(p wbskt.WsConn, v interface{}, action []byte) (err error) {
	if p == nil {
		var pool *wbskt.PoolConn
		if pool, err = sc.initial(); err == nil {
			defer pool.Close()
			p = pool.Conn
		} else {
			return
		}

	}

	if err = p.WriteMessage(websocket.TextMessage, action); err != nil {
		if p != nil {
			p.MarkUnusable()
		}
		return fmt.Errorf("websocket send error: %v", err)
	}
	if err = p.ReadJSON(v); err != nil {
		if p != nil {
			p.MarkUnusable()
		}
		return
	}
	return nil
}


func (sc *SarpcClient) GetBlock(blockHash string) (*rpc.Block, error) {
	v := &rpc.JsonRpcResult{}
	if err := sc.sendWsRequest(nil, v, rpc.ChainGetBlock(wsId, blockHash)); err != nil {
		return nil, err
	}
	rpcBlock := v.ToBlock()
	return &rpcBlock.Block, nil
}

func (sc *SarpcClient) GetExtrinsics(blockHash string) ([]*submodel.Transaction, error) {
	blk, err := sc.GetBlock(blockHash)
	if err != nil {
		return nil, err
	}

	exts := make([]*submodel.Transaction, 0)
	switch sc.chainType {
	case ChainTypeStafi:
		e := new(stafi.ExtrinsicDecoder)
		md, err := sc.getStafiMetaDecoder(blockHash)
		if err != nil {
			return nil, err
		}
		option := stafi.ScaleDecoderOption{Metadata: &md.Metadata, Spec: md.Spec}
		for _, raw := range blk.Extrinsics {
			e.Init(stafi.ScaleBytes{Data: util.HexToBytes(raw)}, &option)
			e.Process()
			if e.ExtrinsicHash != "" && e.ContainsTransaction {
				ext := &submodel.Transaction{
					ExtrinsicHash:  e.ExtrinsicHash,
					CallModuleName: e.CallModule.Name,
					CallName:       e.Call.Name,
					Address:        e.Address,
					Params:         e.Params,
				}
				exts = append(exts, ext)
			}
		}
		return exts, nil
	case ChainTypePolkadot:
		e := new(scale.ExtrinsicDecoder)
		// e := new(polkadot.ExtrinsicDecoder)
		md, err := sc.getPolkaMetaDecoder(blockHash)
		if err != nil {
			return nil, err
		}
		option := scaleTypes.ScaleDecoderOption{Metadata: &md.Metadata, Spec: md.Spec}
		for _, raw := range blk.Extrinsics {
			e.Init(scaleTypes.ScaleBytes{Data: util.HexToBytes(raw)}, &option)
			e.Process()
			decodeExtrinsic := e.Value.(map[string]interface{})
			var ce ChainExtrinsic
			eb, _ := json.Marshal(decodeExtrinsic)
			_ = json.Unmarshal(eb, &ce)
			if e.ExtrinsicHash != "" && e.ContainsTransaction {
				ext := &submodel.Transaction{
					ExtrinsicHash:  e.ExtrinsicHash,
					CallModuleName: ce.CallModule,
					CallName:       ce.CallModuleFunction,
					Address:        e.Address,
					Params:         e.Params,
				}
				exts = append(exts, ext)
			}
		}
		return exts, nil
	default:
		return nil, errors.New("chainType not supported")
	}
}

func (sc *SarpcClient) GetBlockHash(blockNum uint64) (string, error) {
	v := &rpc.JsonRpcResult{}
	if err := sc.sendWsRequest(nil, v, rpc.ChainGetBlockHash(wsId, int(blockNum))); err != nil {
		return "", fmt.Errorf("websocket get block hash error: %v", err)
	}

	blockHash, err := v.ToString()
	if err != nil {
		return "", err
	}

	return blockHash, nil
}

func (sc *SarpcClient) GetChainEvents(blockHash string) ([]*submodel.ChainEvent, error) {
	v := &rpc.JsonRpcResult{}
	if err := sc.sendWsRequest(nil, v, rpc.StateGetStorage(wsId, storageKey, blockHash)); err != nil {
		return nil, fmt.Errorf("websocket get event raw error: %v", err)
	}
	eventRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	var events []*submodel.ChainEvent
	switch sc.chainType {
	case ChainTypeStafi:
		e := stafi.EventsDecoder{}
		md, err := sc.getStafiMetaDecoder(blockHash)
		if err != nil {
			return nil, err
		}

		option := stafi.ScaleDecoderOption{Metadata: &md.Metadata}
		e.Init(stafi.ScaleBytes{Data: util.HexToBytes(eventRaw)}, &option)
		e.Process()
		b, err := json.Marshal(e.Value)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &events)
		if err != nil {
			return nil, err
		}
	case ChainTypePolkadot:
		e := scale.EventsDecoder{}
		md, err := sc.getPolkaMetaDecoder(blockHash)
		if err != nil {
			return nil, err
		}
		option := scaleTypes.ScaleDecoderOption{Metadata: &md.Metadata, Spec: md.Spec}
		e.Init(scaleTypes.ScaleBytes{Data: util.HexToBytes(eventRaw)}, &option)
		e.Process()
		b, err := json.Marshal(e.Value)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &events)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("chainType not supported")
	}

	return events, nil
}

func (sc *SarpcClient) GetEvents(blockNum uint64) ([]*submodel.ChainEvent, error) {
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

func (sc *SarpcClient) GetPaymentQueryInfo(encodedExtrinsic string) (paymentInfo *rpc.PaymentQueryInfo, err error) {
	v := &rpc.JsonRpcResult{}
	if err = sc.sendWsRequest(nil, v, rpc.SystemPaymentQueryInfo(wsId, encodedExtrinsic)); err != nil {
		return
	}

	paymentInfo = v.ToPaymentQueryInfo()
	if paymentInfo == nil {
		return nil, fmt.Errorf("get PaymentQueryInfo error")
	}
	return
}
