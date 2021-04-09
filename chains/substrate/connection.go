// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"errors"
	"fmt"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/stafiprotocol/chainbridge/utils/crypto/sr25519"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
)

type Connection struct {
	url     string
	rsymbol core.RSymbol
	sc      *substrate.SarpcClient
	gc      *substrate.GsrpcClient
	keys    []*signature.KeyringPair
	gcs     map[*signature.KeyringPair]*substrate.GsrpcClient
	log     log15.Logger
	stop    <-chan int
}

var (
	TargetNotExistError = errors.New("TargetNotExistError")
	BlockInterval       = 6 * time.Second
	WaitUntilFinalized  = 10 * BlockInterval
)

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {
	log.Info("NewConnection", "KeystorePath", cfg.KeystorePath, "Endpoint", cfg.Endpoint, "typesPath", cfg.Opts["typesPath"])

	typesPath := cfg.Opts["typesPath"]
	path, ok := typesPath.(string)
	if !ok {
		return nil, errors.New("no typesPath")
	}

	adType := cfg.Opts["addressType"]
	addressType, ok := adType.(string)
	if !ok {
		return nil, errors.New("addressType not ok")
	}

	ct := cfg.Opts["chainType"]
	chainType, ok := ct.(string)
	if !ok {
		return nil, errors.New("chainType not ok")
	}

	sc, err := substrate.NewSarpcClient(chainType, cfg.Endpoint, path, log)
	if err != nil {
		return nil, err
	}

	keys := make([]*signature.KeyringPair, 0)
	gcs := make(map[*signature.KeyringPair]*substrate.GsrpcClient)
	for _, account := range cfg.Accounts {
		kp, err := keystore.KeypairFromAddress(account, keystore.SubChain, cfg.KeystorePath, cfg.Insecure)
		if err != nil {
			return nil, err
		}
		krp := kp.(*sr25519.Keypair).AsKeyringPair()
		gc, err := substrate.NewGsrpcClient(cfg.Endpoint, addressType, krp, log, stop)
		if err != nil {
			return nil, err
		}
		keys = append(keys, krp)
		gcs[krp] = gc
	}

	if len(keys) == 0 {
		return nil, errors.New("no keys")
	}

	return &Connection{
		url:     cfg.Endpoint,
		rsymbol: cfg.Symbol,
		log:     log,
		stop:    stop,
		sc:      sc,
		gc:      gcs[keys[0]],
		keys:    keys,
		gcs:     gcs,
	}, nil
}

func (c *Connection) GetBlockNumber(hash types.Hash) (uint64, error) {
	return c.gc.GetBlockNumber(hash)
}

func (c *Connection) LatestBlockNumber() (uint64, error) {
	return c.gc.GetLatestBlockNumber()
}

func (c *Connection) FinalizedBlockNumber() (uint64, error) {
	return c.gc.GetFinalizedBlockNumber()
}

func (c *Connection) Address() string {
	return c.gc.Address()
}

func (c *Connection) IsConnected() bool {
	return c.sc.IsConnected()
}

func (c *Connection) Reconnect() error {
	return c.sc.WebsocketReconnect()
}

func (c *Connection) GetEvents(blockNum uint64) ([]*submodel.ChainEvent, error) {
	return c.sc.GetEvents(blockNum)
}

// queryStorage performs a storage lookup. Arguments may be nil, result must be a pointer.
func (c *Connection) QueryStorage(prefix, method string, arg1, arg2 []byte, result interface{}) (bool, error) {
	return c.gc.QueryStorage(prefix, method, arg1, arg2, result)
}

func (c *Connection) GetExtrinsics(blockhash string) ([]*submodel.Transaction, error) {
	return c.sc.GetExtrinsics(blockhash)
}

func (c *Connection) LatestMetadata() (*types.Metadata, error) {
	return c.gc.GetLatestMetadata()
}

func (c *Connection) FreeBalance(who []byte) (types.U128, error) {
	return c.gc.FreeBalance(who)
}

func (c *Connection) ExistentialDeposit() (types.U128, error) {
	return c.gc.ExistentialDeposit()
}

func (c *Connection) TransferVerify(r *submodel.BondRecord) (submodel.BondReason, error) {
	bh := hexutil.Encode(r.Blockhash)
	hash := types.NewHash(r.Blockhash)
	blkNum, err := c.GetBlockNumber(hash)
	if err != nil {
		return submodel.BondReasonDefault, err
	}
	if blkNum == 0 {
		return submodel.BlockhashUnmatch, nil
	}

	if !c.IsConnected() {
		if err := c.Reconnect(); err != nil {
			c.log.Error("Reconnect error", "err", err)
			return submodel.BondReasonDefault, err
		}
	}

	final, err := c.FinalizedBlockNumber()
	if err != nil {
		return submodel.BondReasonDefault, err
	}

	if blkNum > final {
		c.log.Info("TransferVerify: block hash not finalized, waiting", "blockHash", bh, "symbol", r.Rsymbol)
		time.Sleep(WaitUntilFinalized)
		final, err = c.FinalizedBlockNumber()
		if err != nil {
			return submodel.BondReasonDefault, err
		}
		if blkNum > final {
			return submodel.BondReasonDefault, errors.New("block number not finalized")
		}
	}

	exts, err := c.GetExtrinsics(bh)
	if err != nil {
		return submodel.BlockhashUnmatch, nil
	}

	th := hexutil.Encode(r.Txhash)
	for _, ext := range exts {
		txhash := utiles.AddHex(ext.ExtrinsicHash)
		if th != txhash {
			c.log.Info("txhash not equal", "expected", th, "got", txhash)
			continue
		}
		c.log.Info("txhash equal", "expected", th, "got", txhash)
		c.log.Info("TransferVerify", "CallModuleName", ext.CallModuleName, "CallName", ext.CallName, "ext.Params number", len(ext.Params))

		if ext.CallModuleName != config.BalancesModuleId || (ext.CallName != config.TransferKeepAlive && ext.CallName != config.Transfer) {
			return submodel.TxhashUnmatch, nil
		}

		addr, ok := ext.Address.(string)
		if !ok {
			c.log.Warn("TransferVerify: address not string", "address", ext.Address)
			return submodel.PubkeyUnmatch, nil
		}

		if hexutil.Encode(r.Pubkey) != utiles.AddHex(addr) {
			c.log.Warn("TransferVerify: pubkey", "addr", addr, "pubkey", hexutil.Encode(r.Pubkey))
			return submodel.PubkeyUnmatch, nil
		}

		for _, p := range ext.Params {
			c.log.Info("TransferVerify", "name", p.Name, "type", p.Type)
			if p.Name == config.ParamDest && p.Type == config.ParamDestType {
				c.log.Debug("cmp dest", "pool", hexutil.Encode(r.Pool), "dest", p.Value)

				dest, ok := p.Value.(string)
				if !ok {
					dest, ok := p.Value.(map[string]interface{})
					if !ok {
						return submodel.PoolUnmatch, nil
					}

					destId, ok := dest["Id"]
					if !ok {
						return submodel.PoolUnmatch, nil
					}

					d, ok := destId.(string)
					if !ok {
						return submodel.PoolUnmatch, nil
					}

					if hexutil.Encode(r.Pool) != utiles.AddHex(d) {
						return submodel.PoolUnmatch, nil
					}
				} else {
					if hexutil.Encode(r.Pool) != utiles.AddHex(dest) {
						return submodel.PoolUnmatch, nil
					}
				}
			} else if p.Name == config.ParamValue && p.Type == config.ParamValueType {
				c.log.Debug("cmp amount", "amount", r.Amount, "paramAmount", p.Value)
				if fmt.Sprint(r.Amount) != fmt.Sprint(p.Value) {
					return submodel.AmountUnmatch, nil
				}
			} else {
				c.log.Error("TransferVerify unexpected param", "name", p.Name, "value", p.Value, "type", p.Type)
				return submodel.TxhashUnmatch, nil
			}
		}

		return submodel.Pass, nil
	}

	return submodel.TxhashUnmatch, nil
}

func (c *Connection) CurrentEra() (uint32, error) {
	var index uint32
	exist, err := c.QueryStorage(config.StakingModuleId, config.StorageActiveEra, nil, nil, &index)
	if err != nil {
		return 0, err
	}

	if !exist {
		return 0, fmt.Errorf("unable to get activeEraInfo for: %s", c.rsymbol)
	}

	return index, nil
}

func (c *Connection) EraContinuable(sym core.RSymbol) (bool, error) {
	symBz, err := types.EncodeToBytes(sym)
	if err != nil {
		return false, err
	}

	ids := make([]types.Hash, 0)
	exists, err := c.QueryStorage(config.RTokenLedgerModuleId, config.StorageCurrentEraSnapShots, symBz, nil, &ids)
	if err != nil {
		return false, err
	}

	if !exists {
		return true, nil
	}

	return len(ids) == 0, nil
}

func (c *Connection) CurrentChainEra(sym core.RSymbol) (uint32, error) {
	symBz, err := types.EncodeToBytes(sym)
	if err != nil {
		return 0, err
	}

	var era uint32
	exists, err := c.QueryStorage(config.RTokenLedgerModuleId, config.StorageChainEras, symBz, nil, &era)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, fmt.Errorf("era of rsymbol %s not exist", sym)
	}

	return era, nil
}

func (c *Connection) IsLastVoter(voter types.AccountID) bool {
	return hexutil.Encode(c.gc.PublicKey()) == hexutil.Encode(voter[:])
}

func (c *Connection) FoundFirstSubAccount(accounts []types.Bytes) (*signature.KeyringPair, []types.AccountID) {
	others := make([]types.AccountID, 0)
	for i, ac := range accounts {
		for _, key := range c.keys {
			if hexutil.Encode(key.PublicKey) == hexutil.Encode(ac) {
				bzs := append(accounts[:i], accounts[i+1:]...)
				for _, bz := range bzs {
					others = append(others, types.NewAccountID(bz))
				}
				return key, others
			}
		}
	}
	return nil, nil
}

func (c *Connection) BondOrUnbondCall(snap *submodel.PoolSnapshot) (*submodel.MultiOpaqueCall, error) {
	return c.gc.BondOrUnbondCall(snap.Bond.Int, snap.Unbond.Int)
}

func (c *Connection) WithdrawCall() (*submodel.MultiOpaqueCall, error) {
	return c.gc.WithdrawCall()
}

func (c *Connection) TransferCalls(receives []*submodel.Receive) ([]*submodel.MultiOpaqueCall, map[string]bool, map[string]bool, error) {
	calls := make([]*submodel.MultiOpaqueCall, 0)
	hashs1 := make(map[string]bool)
	hashs2 := make(map[string]bool)
	for _, rec := range receives {
		call, err := c.gc.TransferCall(rec.Recipient, rec.Value)
		if err != nil {
			return nil, nil, nil, err
		}

		hashs1[call.CallHash] = true
		hashs2[call.CallHash] = true
		calls = append(calls, call)
	}

	return calls, hashs1, hashs2, nil
}

func (c *Connection) PaymentQueryInfo(ext string) (*rpc.PaymentQueryInfo, error) {
	return c.sc.GetPaymentQueryInfo(ext)
}

func (c *Connection) AsMulti(flow *submodel.MultiEventFlow) error {
	gc := c.gcs[flow.Key]
	if gc == nil {
		panic(fmt.Sprintf("key disappear: %s, rsymbol: %s", hexutil.Encode(flow.Key.PublicKey), c.rsymbol))
	}

	l := len(flow.OpaqueCalls)
	if l == 1 {
		moc := flow.OpaqueCalls[0]
		ext, err := gc.NewUnsignedExtrinsic(config.MethodAsMulti, flow.Threshold, flow.Others, moc.TimePoint, moc.Opaque, false, flow.PaymentInfo.Weight)
		if err != nil {
			return err
		}

		return gc.SignAndSubmitTx(ext)
	}

	calls := make([]types.Call, 0)
	for _, oc := range flow.OpaqueCalls {
		ext, err := c.gc.NewUnsignedExtrinsic(config.MethodAsMulti, flow.Threshold, flow.Others, oc.TimePoint, oc.Opaque, false, flow.PaymentInfo.Weight)
		if err != nil {
			return err
		}

		if xt, ok := ext.(*types.Extrinsic); ok {
			calls = append(calls, xt.Method)
		} else if xt, ok := ext.(*types.ExtrinsicMulti); ok {
			calls = append(calls, xt.Method)
		}
	}

	ext, err := gc.NewUnsignedExtrinsic(config.MethodBatch, calls)
	if err != nil {
		return err
	}

	return gc.SignAndSubmitTx(ext)
}

func (c *Connection) SetToPayoutStashes(flow *submodel.BondReportedFlow) error {
	fullTargets := make([]types.AccountID, 0)
	exist, err := c.QueryStorage(config.StakingModuleId, config.StorageNominators, flow.Snap.Pool, nil, &fullTargets)
	if err != nil {
		return err
	}
	if !exist || len(fullTargets) == 0 {
		return TargetNotExistError
	}

	bz, err := types.EncodeToBytes(flow.LastEra)
	if err != nil {
		return err
	}

	points := new(submodel.EraRewardPoints)
	exist, err = c.QueryStorage(config.StakingModuleId, config.StorageErasRewardPoints, bz, nil, points)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("earRewardPoints not exits for era: %d, rsymbol: %s", flow.LastEra, c.rsymbol)
	}

	pointedTargets := make([]types.AccountID, 0)
	for _, tgt := range fullTargets {
		for _, idv := range points.Individuals {
			if hexutil.Encode(tgt[:]) == hexutil.Encode(idv.Validator[:]) {
				pointedTargets = append(pointedTargets, tgt)
			}
		}
	}

	flow.Stashes = pointedTargets
	return nil
}

func (c *Connection) TryPayout(flow *submodel.BondReportedFlow) {
	calls := make([]types.Call, 0)
	meta, err := c.LatestMetadata()
	if err != nil {
		c.log.Error("TryPayout LatestMetadata error", "error", err)
	}
	method := config.MethodPayoutStakers

	for _, stash := range flow.Stashes {
		stashStr := hexutil.Encode(stash[:])

		var controller types.AccountID
		exist, err := c.QueryStorage(config.StakingModuleId, config.StorageBonded, stash[:], nil, &controller)
		if err != nil {
			c.log.Error("TryPayout get controller error", "error", err, "stash", stashStr)
			continue
		}
		if !exist {
			c.log.Error("TryPayout get controller not exist", "stash", stashStr)
			continue
		}
		controllerStr := hexutil.Encode(controller[:])

		ledger := new(submodel.StakingLedger)
		exist, err = c.QueryStorage(config.StakingModuleId, config.StorageLedger, controller[:], nil, ledger)
		if err != nil {
			c.log.Error("TryPayout get ledger error", "error", err, "stash", stashStr)
			continue
		}
		if !exist {
			c.log.Error("TryPayout ledger not exist", "stash", stashStr, "controller", controllerStr)
			continue
		}

		claimed := false
		for _, claimedEra := range ledger.ClaimedRewards {
			if flow.LastEra == claimedEra {
				claimed = true
				break
			}
		}
		if claimed {
			c.log.Info("TryPayout already claimed", "stash", stashStr)
			continue
		}

		call, err := types.NewCall(
			meta,
			method,
			stash,
			flow.LastEra,
		)

		if err != nil {
			c.log.Error("TryPayout NewCall error", "error", err, "stash", stashStr)
			continue
		}

		calls = append(calls, call)
	}

	if len(calls) == 0 {
		return
	}

	ext, err := c.gc.NewUnsignedExtrinsic(config.MethodBatch, calls)
	if err != nil {
		c.log.Error("TryPayout NewUnsignedExtrinsic error", "error", err)
		return
	}

	err = c.gc.SignAndSubmitTx(ext)
	c.log.Info("TryPayout SignAndSubmitTx result", "err", err)
}
