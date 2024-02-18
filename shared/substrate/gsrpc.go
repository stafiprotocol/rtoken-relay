package substrate

import (
	"fmt"
	"github.com/pkg/errors"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client"
	"github.com/stafiprotocol/go-substrate-rpc-client/rpc/author"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

func (sc *SarpcClient) FlashApi() (*gsrpc.SubstrateAPI, error) {
	_, err := sc.api.RPC.Chain.GetBlockHashLatest()
	if err != nil {
		var api *gsrpc.SubstrateAPI
		for i := 0; i < 3; i++ {
			api, err = gsrpc.NewSubstrateAPI(sc.endpoint)
			if err == nil {
				break
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
		if api != nil {
			sc.api = api
		}
	}
	return sc.api, nil
}

func (sc *SarpcClient) Address() string {
	return sc.key.Address
}

func (sc *SarpcClient) GetLatestBlockNumber() (uint64, error) {
	h, err := sc.GetHeaderLatest()
	if err != nil {
		return 0, err
	}

	return uint64(h.Number), nil
}

func (sc *SarpcClient) GetFinalizedBlockNumber() (uint64, error) {
	hash, err := sc.GetFinalizedHead()
	if err != nil {
		return 0, err
	}

	header, err := sc.GetHeader(hash)
	if err != nil {
		return 0, err
	}

	return uint64(header.Number), nil
}

func (sc *SarpcClient) GetHeaderLatest() (*types.Header, error) {
	api, err := sc.FlashApi()
	if err != nil {
		return nil, err
	}
	return api.RPC.Chain.GetHeaderLatest()
}

func (sc *SarpcClient) GetFinalizedHead() (types.Hash, error) {
	api, err := sc.FlashApi()
	if err != nil {
		return types.NewHash([]byte{}), err
	}
	return api.RPC.Chain.GetFinalizedHead()
}

func (sc *SarpcClient) GetHeader(blockHash types.Hash) (*types.Header, error) {
	api, err := sc.FlashApi()
	if err != nil {
		return nil, err
	}
	return api.RPC.Chain.GetHeader(blockHash)
}

func (sc *SarpcClient) GetBlockNumber(blockHash types.Hash) (uint64, error) {
	head, err := sc.GetHeader(blockHash)
	if err != nil {
		return 0, err
	}

	return uint64(head.Number), nil
}

// queryStorage performs a storage lookup. Arguments may be nil, result must be a pointer.
func (sc *SarpcClient) QueryStorage(prefix, method string, arg1, arg2 []byte, result interface{}) (bool, error) {
	entry, err := sc.FindStorageEntryMetadata(prefix, method)
	if err != nil {
		return false, err
	}

	var key types.StorageKey
	keySeted := false
	if entry.IsNMap() {
		hashers, err := entry.Hashers()
		if err != nil {
			return false, err
		}

		if len(hashers) == 1 {
			key, err = types.CreateStorageKeyWithEntryMeta(uint8(sc.metaDataVersion), entry, prefix, method, arg1)
			if err != nil {
				return false, err
			}
			keySeted = true
		}
	}

	if !keySeted {
		key, err = types.CreateStorageKeyWithEntryMeta(uint8(sc.metaDataVersion), entry, prefix, method, arg1, arg2)
		if err != nil {
			return false, err
		}
	}

	api, err := sc.FlashApi()
	if err != nil {
		return false, err
	}

	ok, err := api.RPC.State.GetStorageLatest(key, result)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (sc *SarpcClient) GetLatestRuntimeVersion() (*types.RuntimeVersion, error) {
	api, err := sc.FlashApi()
	if err != nil {
		return nil, err
	}
	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (sc *SarpcClient) GetLatestNonce() (types.U32, error) {
	ac, err := sc.GetAccountInfo()
	if err != nil {
		return 0, err
	}

	return ac.Nonce, nil
}

func (sc *SarpcClient) GetAccountInfo() (*types.AccountInfo, error) {
	ac := new(types.AccountInfo)
	exist, err := sc.QueryStorage("System", "Account", sc.key.PublicKey, nil, &ac)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("account not exist")
	}

	return ac, nil
}

func (sc *SarpcClient) NewUnsignedExtrinsic(callMethod string, args ...interface{}) (interface{}, error) {
	sc.log.Debug("NewUnsignedExtrinsic", "callMethod", callMethod, "addressType", sc.addressType, "sender", sc.key.Address)

	ci, err := sc.FindCallIndex(callMethod)
	if err != nil {
		return nil, errors.Wrap(err, "FindCallIndex")
	}
	sc.log.Debug("callIndex", "methodIndex", ci.MethodIndex, "sectionIndex", ci.SectionIndex)

	call, err := types.NewCallWithCallIndex(ci, callMethod, args...)
	if err != nil {
		return nil, errors.Wrap(err, "NewCallWithCallIndex")
	}

	switch sc.addressType {
	case AddressTypeAccountId:
		unsignedExt := types.NewExtrinsic(call)
		return &unsignedExt, nil
	case AddressTypeMultiAddress:
		unsignedExt := types.NewExtrinsicMulti(call)
		return &unsignedExt, nil
	default:
		return nil, fmt.Errorf("address type is not supported: %s", sc.addressType)
	}
}

func (sc *SarpcClient) SignAndSubmitTx(ext interface{}) error {
	err := sc.signExtrinsic(ext)
	if err != nil {
		return errors.Wrap(err, "signExtrinsic")
	}

	extHexStr, err := types.EncodeToHexString(ext)
	if err != nil {
		return errors.Wrap(err, "EncodeToHexString")
	}
	sc.log.Trace("SignAndSubmitTx", "extrinsic", extHexStr)

	api, err := sc.FlashApi()
	if err != nil {
		return errors.Wrap(err, "FlashApi")
	}
	// Do the transfer and track the actual status
	sub, err := api.RPC.Author.SubmitAndWatch(ext)
	if err != nil {
		return errors.Wrap(err, "SubmitAndWatch")
	}
	defer sub.Unsubscribe()

	return sc.watchSubmission(sub)
}

func (sc *SarpcClient) watchSubmission(sub *author.ExtrinsicStatusSubscription) error {
	for {
		select {
		case <-sc.stop:
			return ErrorTerminated
		case status := <-sub.Chan():
			switch {
			case status.IsInBlock:
				sc.log.Info("Extrinsic included in block", "block", status.AsInBlock.Hex())
				return nil
			case status.IsRetracted:
				return fmt.Errorf("extrinsic retracted: %s", status.AsRetracted.Hex())
			case status.IsDropped:
				return fmt.Errorf("extrinsic dropped from network")
			case status.IsInvalid:
				return fmt.Errorf("extrinsic invalid")
			}
		case err := <-sub.Err():
			sc.log.Trace("Extrinsic subscription error", "err", err)
			return err
		}
	}
}

func (sc *SarpcClient) signExtrinsic(xt interface{}) error {
	sc.log.Debug("signExtrinsic", "addressType", sc.addressType)
	rv, err := sc.GetLatestRuntimeVersion()
	if err != nil {
		return err
	}

	nonce, err := sc.GetLatestNonce()
	if err != nil {
		return err
	}

	o := types.SignatureOptions{
		BlockHash:          sc.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        sc.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	switch ext := xt.(type) {
	case *types.Extrinsic:
		return ext.Sign(*sc.key, o)
	case *types.ExtrinsicMulti:
		return ext.Sign(*sc.key, o)
	default:
		return errors.New("extrinsic unsupport")
	}
}

func (sc *SarpcClient) PublicKey() []byte {
	return sc.key.PublicKey
}

func (sc *SarpcClient) StakingLedger(ac types.AccountID) (*submodel.StakingLedger, error) {
	s := new(submodel.StakingLedger)
	exist, err := sc.QueryStorage(config.StakingModuleId, config.StorageLedger, ac[:], nil, s)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("can not get active for account: %s", hexutil.Encode(ac[:]))
	}

	return s, nil
}

func (sc *SarpcClient) BondOrUnbondCall(bond, unbond *big.Int) (*submodel.RunTimeCall, error) {
	sc.log.Info("BondOrUnbondCall", "bond", bond, "unbond", unbond)
	var method string
	var val types.UCompact

	if bond.Cmp(unbond) < 0 {
		sc.log.Info("unbond larger than bond, UnbondCall")
		diff := big.NewInt(0).Sub(unbond, bond)
		method = config.MethodUnbond
		val = types.NewUCompact(diff)
	} else if bond.Cmp(unbond) > 0 {
		sc.log.Info("bond larger than unbond, BondCall")
		diff := big.NewInt(0).Sub(bond, unbond)
		method = config.MethodBondExtra
		val = types.NewUCompact(diff)
	} else {
		sc.log.Info("bond is equal to unbond, NoCall")
		return nil, ErrorBondEqualToUnbond
	}

	ext, err := sc.NewUnsignedExtrinsic(method, val)
	if err != nil {
		return nil, err
	}

	return CreateRunTimeCall(ext)
}

func (sc *SarpcClient) BondOrUnbondExtrinsic(bond, unbond *big.Int) (interface{}, error) {
	sc.log.Info("BondOrUnbondCall", "bond", bond, "unbond", unbond)
	var method string
	var val types.UCompact

	if bond.Cmp(unbond) < 0 {
		sc.log.Info("unbond larger than bond, UnbondCall")
		diff := big.NewInt(0).Sub(unbond, bond)
		method = config.MethodUnbond
		val = types.NewUCompact(diff)
	} else if bond.Cmp(unbond) > 0 {
		sc.log.Info("bond larger than unbond, BondCall")
		diff := big.NewInt(0).Sub(bond, unbond)
		method = config.MethodBondExtra
		val = types.NewUCompact(diff)
	} else {
		sc.log.Info("bond is equal to unbond, NoCall")
		return nil, ErrorBondEqualToUnbond
	}

	return sc.NewUnsignedExtrinsic(method, val)
}

func (sc *SarpcClient) WithdrawCall() (*submodel.RunTimeCall, error) {
	ext, err := sc.NewUnsignedExtrinsic(config.MethodWithdrawUnbonded, uint32(0))
	if err != nil {
		return nil, err
	}

	return CreateRunTimeCall(ext)
}

func (sc *SarpcClient) TransferCall(accountId []byte, value types.UCompact) (*submodel.RunTimeCall, error) {
	var addr interface{}
	switch sc.addressType {
	case AddressTypeAccountId:
		addr = types.NewAddressFromAccountID(accountId)
	case AddressTypeMultiAddress:
		addr = types.NewMultiAddressFromAccountID(accountId)
	default:
		return nil, fmt.Errorf("addressType not supported: %s", sc.addressType)
	}

	ext, err := sc.NewUnsignedExtrinsic(config.MethodTransferKeepAlive, addr, value)
	if err != nil {
		return nil, err
	}

	return CreateRunTimeCall(ext)
}
func (sc *SarpcClient) TransferExtrinsic(accountId []byte, value types.UCompact) (interface{}, error) {
	var addr interface{}
	switch sc.addressType {
	case AddressTypeAccountId:
		addr = types.NewAddressFromAccountID(accountId)
	case AddressTypeMultiAddress:
		addr = types.NewMultiAddressFromAccountID(accountId)
	default:
		return nil, fmt.Errorf("addressType not supported: %s", sc.addressType)
	}

	return sc.NewUnsignedExtrinsic(config.MethodTransferKeepAlive, addr, value)
}

func (sc *SarpcClient) BatchTransfer(receives []*submodel.Receive) error {
	calls := make([]types.Call, 0)

	ci, err := sc.FindCallIndex(config.MethodTransferKeepAlive)
	if err != nil {
		return err
	}

	for _, rec := range receives {
		var addr interface{}
		switch sc.addressType {
		case AddressTypeAccountId:
			addr = types.NewAddressFromAccountID(rec.Recipient)
		case AddressTypeMultiAddress:
			addr = types.NewMultiAddressFromAccountID(rec.Recipient)
		default:
			return fmt.Errorf("addressType not supported: %s", sc.addressType)
		}

		call, err := types.NewCallWithCallIndex(
			ci,
			config.MethodTransferKeepAlive,
			addr,
			rec.Value,
		)
		if err != nil {
			return err
		}
		calls = append(calls, call)
	}

	ext, err := sc.NewUnsignedExtrinsic(config.MethodBatch, calls)
	if err != nil {
		return err
	}

	return sc.SignAndSubmitTx(ext)
}

func (sc *SarpcClient) NominateCall(validators []types.Bytes) (*submodel.RunTimeCall, error) {
	targets := make([]interface{}, 0)
	switch sc.addressType {
	case AddressTypeAccountId:
		for _, val := range validators {
			targets = append(targets, types.NewAddressFromAccountID(val))
		}
	case AddressTypeMultiAddress:
		for _, val := range validators {
			targets = append(targets, types.NewMultiAddressFromAccountID(val))
		}
	default:
		return nil, fmt.Errorf("addressType not supported: %s", sc.addressType)
	}

	ext, err := sc.NewUnsignedExtrinsic(config.MethodNominate, targets)
	if err != nil {
		return nil, err
	}

	return CreateRunTimeCall(ext)
}

func (sc *SarpcClient) FreeBalance(who []byte) (types.U128, error) {
	if sc.addressType == AddressTypeMultiAddress {
		info, err := sc.NewVersionAccountInfo(who)
		if err != nil {
			return types.U128{}, err
		}
		return info.Data.Free, nil
	}

	info, err := sc.AccountInfo(who)
	if err != nil {
		return types.U128{}, err
	}

	return info.Data.Free, nil
}

func (sc *SarpcClient) AccountInfo(who []byte) (*types.AccountInfo, error) {
	ac := new(types.AccountInfo)
	exist, err := sc.QueryStorage(config.SystemModuleId, config.StorageAccount, who, nil, ac)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("can not get accountInfo for account: %s", hexutil.Encode(who))
	}

	return ac, nil
}

func (sc *SarpcClient) NewVersionAccountInfo(who []byte) (*submodel.AccountInfo, error) {
	ac := new(submodel.AccountInfo)
	exist, err := sc.QueryStorage(config.SystemModuleId, config.StorageAccount, who, nil, ac)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("can not get accountInfo for account: %s", hexutil.Encode(who))
	}

	return ac, nil
}

func (sc *SarpcClient) ExistentialDeposit() (types.U128, error) {
	_, err := sc.FlashApi()
	if err != nil {
		return types.U128{}, err
	}
	var e types.U128
	err = sc.GetConst(config.BalancesModuleId, config.ConstExistentialDeposit, &e)
	if err != nil {
		return types.U128{}, err
	}
	return e, nil
}

func (sc *SarpcClient) GetConst(prefix, name string, res interface{}) error {

	switch sc.chainType {
	case ChainTypeStafi:
		return sc.api.RPC.State.GetConst(prefix, name, &res)
	case ChainTypePolkadot:
		blockHash, err := sc.GetFinalizedHead()
		if err != nil {
			return err
		}
		md, err := sc.getPolkaMetaDecoder(blockHash.Hex())
		if err != nil {
			return err
		}

		for _, mod := range md.Metadata.Metadata.Modules {
			if string(mod.Prefix) == prefix {
				for _, cons := range mod.Constants {
					if cons.Name == name {

						return types.DecodeFromHexString(cons.ConstantsValue, res)
					}
				}
			}
		}
		return fmt.Errorf("could not find constant %s.%s", prefix, name)
	default:
		return errors.New("GetConst chainType not supported")
	}
}

func CreateRunTimeCall(ext interface{}) (*submodel.RunTimeCall, error) {
	var call types.Call

	switch xt := ext.(type) {
	case *types.Extrinsic:
		call = xt.Method
	case *types.ExtrinsicMulti:
		call = xt.Method
	default:
		return nil, errors.New("not support extrinsic type")
	}

	opaque, err := types.EncodeToBytes(call)
	if err != nil {
		return nil, err
	}
	callhash := utils.BlakeTwo256(opaque)

	bz, err := types.EncodeToBytes(ext)
	if err != nil {
		return nil, err
	}

	return &submodel.RunTimeCall{
		Extrinsic: hexutil.Encode(bz),
		Call:      call,
		CallHash:  hexutil.Encode(callhash[:]),
	}, nil
}

func (sc *SarpcClient) FindStorageEntryMetadata(module string, fn string) (types.StorageEntryMetadata, error) {
	switch sc.chainType {
	case ChainTypeStafi:
		meta, err := sc.api.RPC.State.GetMetadataLatest()
		if err != nil {
			return nil, err
		}

		return meta.FindStorageEntryMetadata(module, fn)
	case ChainTypePolkadot:
		blockHash, err := sc.GetFinalizedHead()
		if err != nil {
			return nil, err
		}
		md, err := sc.getPolkaMetaDecoder(blockHash.Hex())
		if err != nil {
			return nil, err
		}

		for _, mod := range md.Metadata.Metadata.Modules {
			if string(mod.Prefix) != module {
				continue
			}
			for _, s := range mod.Storage {
				if string(s.Name) != fn {
					continue
				}

				sfm := types.StorageFunctionMetadataV13{
					Name: types.Text(s.Name),
				}

				if s.Type.PlainType != nil {
					sfm.Type = types.StorageFunctionTypeV13{
						IsType: true,
						AsType: types.Type(*s.Type.PlainType),
					}
				}

				if s.Type.DoubleMapType != nil {
					dmt := types.DoubleMapTypeV10{
						Key1:       types.Type(s.Type.DoubleMapType.Key),
						Key2:       types.Type(s.Type.DoubleMapType.Key2),
						Value:      types.Type(s.Type.DoubleMapType.Value),
						Hasher:     TransformHasher(s.Type.DoubleMapType.Hasher),
						Key2Hasher: TransformHasher(s.Type.DoubleMapType.Key2Hasher),
					}

					sfm.Type = types.StorageFunctionTypeV13{
						IsDoubleMap: true,
						AsDoubleMap: dmt,
					}
				}

				if s.Type.MapType != nil {
					mt := types.MapTypeV10{
						Key:    types.Type(s.Type.MapType.Key),
						Value:  types.Type(s.Type.MapType.Value),
						Linked: s.Type.MapType.IsLinked,
						Hasher: TransformHasher(s.Type.MapType.Hasher),
					}

					sfm.Type = types.StorageFunctionTypeV13{
						IsMap: true,
						AsMap: mt,
					}
				}

				if s.Type.NMapType != nil {
					keys := make([]types.Type, 0)
					for _, key := range s.Type.NMapType.KeyVec {
						keys = append(keys, types.Type(key))
					}

					hashers := make([]types.StorageHasherV10, 0)
					for _, hasher := range s.Type.NMapType.Hashers {
						hashers = append(hashers, TransformHasher(hasher))
					}

					nmt := types.NMapTypeV13{
						Keys:    keys,
						Hashers: hashers,
						Value:   types.Type(s.Type.NMapType.Value),
					}

					sfm.Type = types.StorageFunctionTypeV13{
						IsNMap: true,
						AsNMap: nmt,
					}
				}

				return sfm, nil
			}
			return nil, fmt.Errorf("storage %v not found within module %v", fn, module)
		}
		return nil, fmt.Errorf("module %v not found in metadata", module)
	default:
		return nil, errors.New("chainType not supported")
	}
}

func (sc *SarpcClient) FindCallIndex(call string) (types.CallIndex, error) {
	switch sc.chainType {
	case ChainTypeStafi:
		meta, err := sc.api.RPC.State.GetMetadataLatest()
		if err != nil {
			return types.CallIndex{}, err
		}

		return meta.FindCallIndex(call)
	case ChainTypePolkadot:
		blockHash, err := sc.GetFinalizedHead()
		if err != nil {
			return types.CallIndex{}, err
		}

		md, err := sc.getPolkaMetaDecoder(blockHash.Hex())
		if err != nil {
			return types.CallIndex{}, err
		}
		s := strings.Split(call, ".")

		for _, mod := range md.Metadata.Metadata.Modules {
			if string(mod.Name) != s[0] {
				continue
			}
			for ci, f := range mod.Calls {
				if string(f.Name) == s[1] {
					mIndex := uint8(ci)
					if strings.EqualFold(mod.Name, "Balances") {
						if !strings.EqualFold(f.Name, "transfer_allow_death") && !strings.EqualFold(f.Name, "force_transfer") {
							mIndex++
						}
					}
					return types.CallIndex{SectionIndex: uint8(mod.Index), MethodIndex: mIndex}, nil
				}
			}
			return types.CallIndex{}, fmt.Errorf("method %v not found within module %v for call %v", s[1], mod.Name, call)
		}
		return types.CallIndex{}, fmt.Errorf("module %v not found in metadata for call %v", s[0], call)

	default:
		return types.CallIndex{}, errors.New("FindCallIndex chainType not supported")
	}
}

func TransformHasher(Hasher string) types.StorageHasherV10 {
	if Hasher == "Blake2_128" {
		return types.StorageHasherV10{IsBlake2_128: true}
	}

	if Hasher == "Blake2_256" {
		return types.StorageHasherV10{IsBlake2_256: true}
	}

	if Hasher == "Blake2_128Concat" {
		return types.StorageHasherV10{IsBlake2_128Concat: true}
	}

	if Hasher == "Twox128" {
		return types.StorageHasherV10{IsTwox128: true}
	}

	if Hasher == "Twox256" {
		return types.StorageHasherV10{IsTwox256: true}
	}

	if Hasher == "Twox64Concat" {
		return types.StorageHasherV10{IsTwox64Concat: true}
	}

	return types.StorageHasherV10{IsIdentity: true}
}
