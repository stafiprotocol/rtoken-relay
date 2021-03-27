package substrate

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client"
	"github.com/stafiprotocol/go-substrate-rpc-client/rpc/author"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

type GsrpcClient struct {
	endpoint    string
	api         *gsrpc.SubstrateAPI
	key         *signature.KeyringPair
	genesisHash types.Hash
	stop        <-chan int
	log         log15.Logger
}

func NewGsrpcClient(endpoint string, key *signature.KeyringPair, log log15.Logger, stop <-chan int) (*GsrpcClient, error) {
	log.Info("Connecting to substrate chain with Gsrpc", "endpoint", endpoint)

	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		return nil, err
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}

	return &GsrpcClient{
		endpoint:    endpoint,
		api:         api,
		key:         key,
		genesisHash: genesisHash,
		stop:        stop,
		log:         log,
	}, nil
}

func (gc *GsrpcClient) Address() string {
	return gc.key.Address
}

func (gc *GsrpcClient) GetLatestBlockNumber() (uint64, error) {
	h, err := gc.GetHeaderLatest()
	if err != nil {
		return 0, err
	}

	return uint64(h.Number), nil
}

func (gc *GsrpcClient) GetFinalizedBlockNumber() (uint64, error) {
	hash, err := gc.GetFinalizedHead()
	if err != nil {
		return 0, err
	}

	header, err := gc.GetHeader(hash)
	if err != nil {
		return 0, err
	}

	return uint64(header.Number), nil
}

func (gc *GsrpcClient) GetHeaderLatest() (*types.Header, error) {
	return gc.api.RPC.Chain.GetHeaderLatest()
}

func (gc *GsrpcClient) GetFinalizedHead() (types.Hash, error) {
	return gc.api.RPC.Chain.GetFinalizedHead()
}

func (gc *GsrpcClient) GetHeader(blockHash types.Hash) (*types.Header, error) {
	return gc.api.RPC.Chain.GetHeader(blockHash)
}

// queryStorage performs a storage lookup. Arguments may be nil, result must be a pointer.
func (gc *GsrpcClient) QueryStorage(prefix, method string, arg1, arg2 []byte, result interface{}) (bool, error) {
	meta, err := gc.GetLatestMetadata()
	if err != nil {
		return false, err
	}

	key, err := types.CreateStorageKey(meta, prefix, method, arg1, arg2)
	if err != nil {
		return false, err
	}

	ok, err := gc.api.RPC.State.GetStorageLatest(key, result)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (gc *GsrpcClient) GetLatestMetadata() (*types.Metadata, error) {
	meta, err := gc.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (gc *GsrpcClient) GetLatestRuntimeVersion() (*types.RuntimeVersion, error) {
	rv, err := gc.api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (gc *GsrpcClient) GetLatestNonce() (types.U32, error) {
	ac, err := gc.GetAccountInfo()
	if err != nil {
		return 0, err
	}

	return ac.Nonce, nil
}

func (gc *GsrpcClient) GetAccountInfo() (*types.AccountInfo, error) {
	ac := new(types.AccountInfo)
	exist, err := gc.QueryStorage("System", "Account", gc.key.PublicKey, nil, &ac)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("account not exist")
	}

	return ac, nil
}

func (gc *GsrpcClient) NewUnsignedExtrinsic(callMethod string, args ...interface{}) (*types.Extrinsic, error) {
	gc.log.Debug("Submitting substrate call...", "callMethod", callMethod, "sender", gc.key.Address)
	meta, err := gc.GetLatestMetadata()
	if err != nil {
		return nil, err
	}

	call, err := types.NewCall(meta, callMethod, args...)
	if err != nil {
		return nil, err
	}

	unsignedExt := types.NewExtrinsic(call)
	return &unsignedExt, nil
}

func (gc *GsrpcClient) SignAndSubmitTx(ext *types.Extrinsic) error {
	err := gc.signExtrinsic(ext)
	if err != nil {
		return err
	}

	// Do the transfer and track the actual status
	sub, err := gc.api.RPC.Author.SubmitAndWatchExtrinsic(*ext)
	if err != nil {
		return err
	}
	gc.log.Trace("Extrinsic submission succeeded")
	defer sub.Unsubscribe()

	return gc.watchSubmission(sub)
}

func (gc *GsrpcClient) watchSubmission(sub *author.ExtrinsicStatusSubscription) error {
	for {
		select {
		case <-gc.stop:
			return TerminatedError
		case status := <-sub.Chan():
			switch {
			case status.IsInBlock:
				gc.log.Info("Extrinsic included in block", "block", status.AsInBlock.Hex())
				return nil
			case status.IsRetracted:
				return fmt.Errorf("extrinsic retracted: %s", status.AsRetracted.Hex())
			case status.IsDropped:
				return fmt.Errorf("extrinsic dropped from network")
			case status.IsInvalid:
				return fmt.Errorf("extrinsic invalid")
			}
		case err := <-sub.Err():
			gc.log.Trace("Extrinsic subscription error", "err", err)
			return err
		}
	}
}

func (gc *GsrpcClient) signExtrinsic(ext *types.Extrinsic) error {
	rv, err := gc.GetLatestRuntimeVersion()
	if err != nil {
		return err
	}

	nonce, err := gc.GetLatestNonce()
	if err != nil {
		return err
	}

	o := types.SignatureOptions{
		BlockHash:          gc.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        gc.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	err = ext.Sign(*gc.key, o)
	if err != nil {
		return err
	}

	return nil
}

func (gc *GsrpcClient) PublicKey() []byte {
	return gc.key.PublicKey
}

func (gc *GsrpcClient) StakingLedger(ac types.AccountID) (*StakingLedger, error) {
	s := new(StakingLedger)
	exist, err := gc.QueryStorage(config.StakingModuleId, config.StorageLedger, ac[:], nil, s)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("can not get active for account: %s", hexutil.Encode(ac[:]))
	}

	return s, nil
}

func (gc *GsrpcClient) BondOrUnbondCall(bond, unbond *big.Int) (*MultiOpaqueCall, error) {
	gc.log.Info("BondOrUnbondCall", "bond", bond, "unbond", unbond)
	var method string
	var val types.UCompact

	if bond.Cmp(unbond) < 0 {
		gc.log.Info("unbond larger than bond, UnbondCall")
		diff := big.NewInt(0).Sub(unbond, bond)
		method = config.MethodUnbond
		val = types.NewUCompact(diff)
	} else if bond.Cmp(unbond) > 0 {
		gc.log.Info("bond larger than unbond, BondCall")
		diff := big.NewInt(0).Sub(bond, unbond)
		method = config.MethodBondExtra
		val = types.NewUCompact(diff)
	} else {
		gc.log.Info("bond is equal to unbond, NoCall")
		return nil, BondEqualToUnbondError
	}

	ext, err := gc.NewUnsignedExtrinsic(method, val)
	if err != nil {
		return nil, err
	}

	return OpaqueCall(ext)
}

func (gc *GsrpcClient) WithdrawCall() (*MultiOpaqueCall, error) {
	ext, err := gc.NewUnsignedExtrinsic(config.MethodWithdrawUnbonded, uint32(0))
	if err != nil {
		return nil, err
	}

	return OpaqueCall(ext)
}

func (gc *GsrpcClient) TransferCall(dest types.Address, value types.UCompact) (*MultiOpaqueCall, error) {
	ext, err := gc.NewUnsignedExtrinsic(config.MethodTransferKeepAlive, dest, value)
	if err != nil {
		return nil, err
	}

	return OpaqueCall(ext)
}

func (gc *GsrpcClient) AccountInfo(who []byte) (*types.AccountInfo, error) {
	ac := new(types.AccountInfo)
	exist, err := gc.QueryStorage(config.SystemModuleId, config.StorageAccount, who, nil, ac)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("can not get accountInfo for account: %s", hexutil.Encode(who))
	}

	return ac, nil
}

func (gc *GsrpcClient) ExistentialDeposit() (types.U128, error) {
	var e types.U128
	err := gc.api.RPC.State.GetConst(config.BalancesModuleId, config.ConstExistentialDeposit, &e)
	if err != nil {
		return types.U128{}, err
	}
	return e, nil
}

func OpaqueCall(ext *types.Extrinsic) (*MultiOpaqueCall, error) {
	opaque, err := types.EncodeToBytes(ext.Method)
	if err != nil {
		return nil, err
	}

	bz, err := types.EncodeToBytes(ext)
	if err != nil {
		return nil, err
	}

	callhash := utils.BlakeTwo256(opaque)
	return &MultiOpaqueCall{
		Extrinsic: hexutil.Encode(bz),
		Opaque:    opaque,
		CallHash:  hexutil.Encode(callhash[:]),
	}, nil
}
