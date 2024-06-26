package substrate

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gorilla/websocket"
	scale "github.com/itering/scale.go"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/source"
	scaleTypes "github.com/itering/scale.go/types"
	scaleBytes "github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/pkg/recws"
	"github.com/itering/substrate-api-rpc/rpc"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client"
	gsrpcConfig "github.com/stafiprotocol/go-substrate-rpc-client/config"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	wbskt "github.com/stafiprotocol/rtoken-relay/shared/substrate/websocket"
	chainTypes "github.com/stafiprotocol/rtoken-relay/types"
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
	log       core.Logger
	chainType string
	typesPath string

	currentSpecVersion int

	stafiMetaDecoderMap map[int]*stafi.MetadataDecoder
	polkaMetaDecoderMap map[int]*scale.MetadataDecoder
	sync.RWMutex

	metaDataVersion int
}

func NewSarpcClient(chainType, endpoint, typesPath, addressType string, key *signature.KeyringPair, log core.Logger, stop <-chan int) (*SarpcClient, error) {
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
		currentSpecVersion:  -1,
		stafiMetaDecoderMap: make(map[int]*stafi.MetadataDecoder),
		polkaMetaDecoderMap: make(map[int]*scale.MetadataDecoder),
	}

	sc.regCustomTypes()

	switch chainType {
	case ChainTypeStafi:
		_, err := sc.getStafiMetaDecoder(latestHash.Hex())
		if err != nil {
			return nil, err
		}
	case ChainTypePolkadot:
		_, err := sc.getPolkaMetaDecoder(latestHash.Hex())
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("chain type not support: %s", chainType)
	}

	return sc, nil
}

func (s *SarpcClient) getStafiMetaDecoder(blockHash string) (*stafi.MetadataDecoder, error) {
	v := &rpc.JsonRpcResult{}
	// runtime version
	if err := s.sendWsRequest(v, rpc.ChainGetRuntimeVersion(wsId, blockHash)); err != nil {
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
	if err := s.sendWsRequest(v, rpc.StateGetMetadata(wsId, blockHash)); err != nil {
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

	if r.SpecVersion > s.currentSpecVersion {
		s.currentSpecVersion = r.SpecVersion
		s.metaDataVersion = md.Metadata.MetadataVersion
	}
	return md, nil

}

func (s *SarpcClient) getPolkaMetaDecoder(blockHash string) (*scale.MetadataDecoder, error) {
	v := &rpc.JsonRpcResult{}
	// runtime version
	if err := s.sendWsRequest(v, rpc.ChainGetRuntimeVersion(wsId, blockHash)); err != nil {
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
	if err := s.sendWsRequest(v, rpc.StateGetMetadata(wsId, blockHash)); err != nil {
		return nil, err
	}
	metaRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	md := scale.MetadataDecoder{}
	md.Init(utiles.HexToBytes(metaRaw))
	if err := md.Process(); err != nil {
		return nil, err
	}
	s.Lock()
	s.polkaMetaDecoderMap[r.SpecVersion] = &md
	s.Unlock()

	if r.SpecVersion > s.currentSpecVersion {
		s.currentSpecVersion = r.SpecVersion
		s.metaDataVersion = md.Metadata.MetadataVersion
	}

	return &md, nil

}

func (sc *SarpcClient) regCustomTypes() {
	content, err := os.ReadFile(sc.typesPath)
	if err != nil {
		panic(err)
	}

	switch sc.chainType {
	case ChainTypeStafi:
		stafi.RuntimeType{}.Reg()
		stafi.RegCustomTypes(source.LoadTypeRegistry(content))
	case ChainTypePolkadot:
		scaleTypes.RegCustomTypes(source.LoadTypeRegistry(content))
	default:
		panic("chainType not supported")
	}
}

func (sc *SarpcClient) initial() (*wbskt.PoolConn, error) {
	var err error
	if sc.wsPool == nil {
		factory := func() (*recws.RecConn, error) {
			SubscribeConn := &recws.RecConn{KeepAliveTimeout: 2 * time.Minute}
			SubscribeConn.Dial(sc.endpoint, nil)
			sc.log.Debug("conn factory create new conn", "endpoint", sc.endpoint)
			return SubscribeConn, err
		}
		if sc.wsPool, err = wbskt.NewChannelPool(1, 25, factory); err != nil {
			sc.log.Error("wbskt.NewChannelPool", "err", err)
		}
	}
	if err != nil {
		return nil, err
	}
	conn, err := sc.wsPool.Get()
	return conn, err
}

func (sc *SarpcClient) sendWsRequest(v interface{}, action []byte) error {

	retry := 0
	for {
		if retry >= 100 {
			return fmt.Errorf("sendWsRequest reach retry limit")
		}

		if poolConn, err := sc.initial(); err != nil {
			return err
		} else {
			p := poolConn.Conn

			if err = p.WriteMessage(websocket.TextMessage, action); err != nil {
				poolConn.MarkUnusable()
				poolConn.Close()

				sc.log.Warn("websocket send error", "err", err)
				time.Sleep(time.Millisecond * 100)
				retry++
				continue
			}

			if err = p.ReadJSON(v); err != nil {
				poolConn.MarkUnusable()
				poolConn.Close()

				sc.log.Warn("websocket read error", "err", err)
				time.Sleep(time.Millisecond * 100)
				retry++
				continue
			}
			poolConn.Close()
			return nil
		}
	}
}

func (sc *SarpcClient) GetBlock(blockHash string) (*rpc.Block, error) {
	v := &rpc.JsonRpcResult{}
	if err := sc.sendWsRequest(v, rpc.ChainGetBlock(wsId, blockHash)); err != nil {
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
		md, err := sc.getStafiMetaDecoder(blockHash)
		if err != nil {
			return nil, err
		}
		e := new(stafi.ExtrinsicDecoder)
		option := stafi.ScaleDecoderOption{Metadata: &md.Metadata}
		for _, raw := range blk.Extrinsics {
			e.Init(stafi.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
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
		md, err := sc.getPolkaMetaDecoder(blockHash)
		if err != nil {
			return nil, err
		}
		e := new(scale.ExtrinsicDecoder)
		option := scaleTypes.ScaleDecoderOption{Metadata: &md.Metadata}
		for _, raw := range blk.Extrinsics {
			e.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(raw)}, &option)
			e.Process()
			decodeExtrinsic := e.Value.(*scalecodec.GenericExtrinsic)
			var ce ChainExtrinsic
			eb, _ := json.Marshal(decodeExtrinsic)
			_ = json.Unmarshal(eb, &ce)
			if e.ExtrinsicHash != "" && e.ContainsTransaction {
				params := make([]chainTypes.ExtrinsicParam, 0)
				for _, p := range e.Params {
					params = append(params, chainTypes.ExtrinsicParam{
						Name:  p.Name,
						Type:  p.Type,
						Value: p.Value,
					})
				}

				ext := &submodel.Transaction{
					ExtrinsicHash:  e.ExtrinsicHash,
					CallModuleName: ce.CallModule,
					CallName:       ce.CallModuleFunction,
					Address:        e.Address,
					Params:         params,
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
	if err := sc.sendWsRequest(v, rpc.ChainGetBlockHash(wsId, int(blockNum))); err != nil {
		return "", fmt.Errorf("websocket get block hash error: %v", err)
	}

	blockHash, err := v.ToString()
	if err != nil {
		return "", fmt.Errorf("ChainGetBlockHash get error %v", err)
	}
	if blockHash == "" {
		return "", fmt.Errorf("ChainGetBlockHash error, blockHash empty")
	}

	return blockHash, nil
}

func (sc *SarpcClient) GetChainEvents(blockHash string) ([]*submodel.ChainEvent, error) {
	v := &rpc.JsonRpcResult{}
	if err := sc.sendWsRequest(v, rpc.StateGetStorage(wsId, storageKey, blockHash)); err != nil {
		return nil, fmt.Errorf("websocket get event raw error: %v", err)
	}
	eventRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	var events []*submodel.ChainEvent
	switch sc.chainType {
	case ChainTypeStafi:
		md, err := sc.getStafiMetaDecoder(blockHash)
		if err != nil {
			return nil, err
		}
		e := stafi.EventsDecoder{}
		option := stafi.ScaleDecoderOption{Metadata: &md.Metadata}
		e.Init(stafi.ScaleBytes{Data: utiles.HexToBytes(eventRaw)}, &option)
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
		md, err := sc.getPolkaMetaDecoder(blockHash)
		if err != nil {
			return nil, err
		}
		option := scaleTypes.ScaleDecoderOption{Metadata: &md.Metadata}

		e := scale.EventsDecoder{}
		e.Init(scaleBytes.ScaleBytes{Data: utiles.HexToBytes(eventRaw)}, &option)

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
	if err = sc.sendWsRequest(v, rpc.SystemPaymentQueryInfo(wsId, encodedExtrinsic)); err != nil {
		return
	}
	sc.log.Trace("GetPaymentQueryInfo", "json", fmt.Sprintf("%+v", v))

	paymentInfo = v.ToPaymentQueryInfo()
	if paymentInfo == nil {
		return nil, fmt.Errorf("get PaymentQueryInfo error")
	}
	return
}

func (sc *SarpcClient) GetPaymentQueryInfoV2(encodedExtrinsic string) (paymentInfo *submodel.PaymentQueryInfoV2, err error) {
	v := &rpc.JsonRpcResult{}
	if err = sc.sendWsRequest(v, rpc.SystemPaymentQueryInfo(wsId, encodedExtrinsic)); err != nil {
		return
	}
	sc.log.Trace("GetPaymentQueryInfo", "json", fmt.Sprintf("%+v", v))

	if v.Error != nil {
		return nil, errors.New(v.Error.Message)
	}

	if (v).Result == nil {
		return nil, fmt.Errorf("get PaymentQueryInfo error")
	}
	result := (v).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, fmt.Errorf("get PaymentQueryInfo error")
	}
	value := &submodel.PaymentQueryInfoV2{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("get PaymentQueryInfo error")
	}
	_ = json.Unmarshal([]byte(marshal), value)

	return value, nil
}

func (sc *SarpcClient) IsClaimed(controller []byte, validator []byte, lastEra uint32) (bool, error) {
	ledger := new(submodel.StakingLedger)
	exist, err := sc.QueryStorage(config.StakingModuleId, config.StorageLedger, controller, nil, ledger)
	if err != nil {
		return false, fmt.Errorf("get ledger error: %s, controller: %s, era:%d", err, hexutil.Encode(controller), lastEra)
	}
	if !exist {
		return false, fmt.Errorf("get StorageLedger not exist, controller: %s, era:%d", hexutil.Encode(controller), lastEra)
	}

	claimed := false
	for _, claimedEra := range ledger.ClaimedRewards {
		if lastEra == claimedEra {
			return true, nil
		}
	}

	if !claimed {
		pageCount := uint32(1)
		overview, exist, err := sc.ErasStakersOverview(lastEra, validator)
		if err != nil {
			return false, fmt.Errorf("ErasStakersOverview failed: %s, validator: %s", err.Error(), hex.EncodeToString(validator))
		}
		if exist {
			if overview.PageCount > 1 {
				pageCount = overview.PageCount
			}
		}
		rewards, err := sc.ClaimedRewards(lastEra, validator)
		if err != nil {
			return false, fmt.Errorf("ClaimedRewards failed: %s, validator: %s", err.Error(), hex.EncodeToString(validator))
		}
		rewardsPageExist := make(map[uint32]bool)
		for _, r := range rewards {
			rewardsPageExist[r] = true
		}

		for i := uint32(0); i < pageCount; i++ {
			if !rewardsPageExist[i] {
				return false, nil
			}
		}
		return true, nil

	}

	return claimed, nil
}
