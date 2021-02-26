package substrate

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ChainSafe/log15"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client"
	"github.com/stafiprotocol/go-substrate-rpc-client/rpc/author"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/conn"
)

type GsrpcClient struct {
	endpoint    string
	api         *gsrpc.SubstrateAPI
	key         *signature.KeyringPair
	genesisHash types.Hash
	ctx         context.Context
	log         log15.Logger
}

func NewGsrpcClient(ctx context.Context, endpoint string, key *signature.KeyringPair, log log15.Logger) (*GsrpcClient, error) {
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
		ctx:         ctx,
		log:         log,
	}, nil
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
		case <-gc.ctx.Done():
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

func (gc *GsrpcClient) TryToBondOrUnbond(lc *conn.LinkChunk) error {
	bond := lc.Bond
	unbond := lc.Unbond
	if bond.Cmp(unbond.Int) < 0 {
		diff := big.NewInt(0).Sub(unbond.Int, bond.Int)
		realUnbond := types.NewU128(*diff)
		err := gc.TryToUnbond(realUnbond)
		if err != nil {
			return err
		}

	} else if bond.Cmp(unbond.Int) > 0 {
		diff := big.NewInt(0).Sub(bond.Int, unbond.Int)
		realBond := types.NewU128(*diff)
		err := gc.TryToBondExtra(realBond)
		if err != nil {
			return err
		}
	} else {
		gc.log.Info("bond is equal to unbond, nothing need to do")
	}

	return nil
}

func (gc *GsrpcClient) TryToUnbond(value types.U128) error {
	ext, err := gc.NewUnsignedExtrinsic(UnBondMethod, value)
	err = gc.SignAndSubmitTx(ext)
	if err != nil {
		gc.log.Error("TryToBondExtra error", "err", err)
		return err
	}
	return nil
}

func (gc *GsrpcClient) TryToBondExtra(value types.U128) error {
	ext, err := gc.NewUnsignedExtrinsic(BondExtraMethod, value)
	err = gc.SignAndSubmitTx(ext)
	if err != nil {
		gc.log.Error("TryToBondExtra error", "err", err)
		return err
	}
	return nil
}

func (gc *GsrpcClient) TryToClaim(lc *conn.LinkChunk) error {
	//for _, nom := range Nominators {
	//}
	//return nil
	return nil
}
