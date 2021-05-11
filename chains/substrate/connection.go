// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"bytes"
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
	symbol  core.RSymbol
	sc      *substrate.SarpcClient
	gc      *substrate.GsrpcClient
	keys    []*signature.KeyringPair
	gcs     map[*signature.KeyringPair]*substrate.GsrpcClient
	log     log15.Logger
	stop    <-chan int
	lastKey *signature.KeyringPair
}

var (
	TargetNotExistError = errors.New("TargetNotExistError")
	NotExistError       = errors.New("not exist")
	BlockInterval       = 6 * time.Second
	WaitUntilFinalized  = 10 * BlockInterval

	WsRetryLimit    = 240
	WsRetryInterval = 500 * time.Millisecond
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

	acSize := len(cfg.Accounts)
	var lk *signature.KeyringPair
	for i := 0; i < acSize; i++ {
		kp, err := keystore.KeypairFromAddress(cfg.Accounts[i], keystore.SubChain, cfg.KeystorePath, cfg.Insecure)
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

		if cfg.Symbol != core.RFIS && i+1 == acSize {
			lk = krp
		}
	}

	if len(keys) == 0 {
		return nil, errors.New("no keys")
	}

	return &Connection{
		url:     cfg.Endpoint,
		symbol:  cfg.Symbol,
		log:     log,
		stop:    stop,
		sc:      sc,
		gc:      gcs[keys[0]],
		keys:    keys,
		gcs:     gcs,
		lastKey: lk,
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
		for i := 0; i < 10; i++ {
			time.Sleep(BlockInterval)
			blkNum, err = c.GetBlockNumber(hash)
			if err != nil {
				return submodel.BondReasonDefault, err
			}
			if blkNum != 0 {
				break
			}
		}
		if blkNum == 0 {
			return submodel.BondReasonDefault, errors.New("after waiting more than one minute, blkNum is still zero")
		}
	}

	final, err := c.FinalizedBlockNumber()
	if err != nil {
		return submodel.BondReasonDefault, err
	}

	if blkNum > final {
		c.log.Info("TransferVerify: block hash not finalized, waiting", "blockHash", bh, "symbol", r.Symbol)
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
		c.log.Warn("TransferVerify: get extrinsics error", "err", err, "blockHash", bh)
		return submodel.BondReasonDefault, err
	}

	th := hexutil.Encode(r.Txhash)
	for _, ext := range exts {
		c.log.Info("TransferVerify loop extrinsics", "ext", ext)
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
		return 0, fmt.Errorf("unable to get activeEraInfo for: %s", c.symbol)
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

func (c *Connection) CurrentEraSnapshots(symbol core.RSymbol) ([]types.Hash, error) {
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return nil, err
	}

	ids := make([]types.Hash, 0)
	exists, err := c.QueryStorage(config.RTokenLedgerModuleId, config.StorageCurrentEraSnapShots, symBz, nil, &ids)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("currentEraSnapshots not exist")
	}
	return ids, nil
}

func (c *Connection) GetEraNominated(symbol core.RSymbol, pool []byte, era uint32) ([]types.AccountID, error) {
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return nil, err
	}

	puk := &submodel.PoolUnbondKey{Pool: pool, Era: era}
	pkbz, err := types.EncodeToBytes(puk)
	if err != nil {
		return nil, err
	}

	validators := make([]types.AccountID, 0)
	exist, err := c.QueryStorage(config.RTokenSeriesModuleId, config.StorageEraNominated, symBz, pkbz, &validators)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, NotExistError
	}
	return validators, nil
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
		return 0, fmt.Errorf("era of symbol %s not exist", sym)
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

func (c *Connection) FoundIndexOfSubAccount(accounts []types.Bytes) (int, *substrate.GsrpcClient) {
	for i, ac := range accounts {
		for _, key := range c.keys {
			if hexutil.Encode(key.PublicKey) == hexutil.Encode(ac) {
				return i, c.gcs[key]
			}
		}
	}
	return -1, nil
}

func (c *Connection) KeyIndex(key *signature.KeyringPair) *substrate.GsrpcClient {
	return c.gcs[key]
}

func (c *Connection) BondOrUnbondCall(snap *submodel.PoolSnapshot) (*submodel.MultiOpaqueCall, error) {
	return c.gc.BondOrUnbondCall(snap.Bond.Int, snap.Unbond.Int)
}

func (c *Connection) WithdrawCall() (*submodel.MultiOpaqueCall, error) {
	return c.gc.WithdrawCall()
}

func (c *Connection) TransferCall(recipient []byte, value types.UCompact) (*submodel.MultiOpaqueCall, error) {
	return c.gc.TransferCall(recipient, value)
}

func (c *Connection) LastKey() *signature.KeyringPair {
	return c.lastKey
}

func (c *Connection) NominateCall(validators []types.Bytes) (*submodel.MultiOpaqueCall, error) {
	return c.gc.NominateCall(validators)
}

func (c *Connection) PaymentQueryInfo(ext string) (info *rpc.PaymentQueryInfo, err error) {
	for i := 0; i < WsRetryLimit; i++ {
		info, err = c.sc.GetPaymentQueryInfo(ext)
		if err == nil {
			return
		}

		time.Sleep(WsRetryInterval)
	}

	return
}

func (c *Connection) AsMulti(flow *submodel.MultiEventFlow) error {
	for i := 0; i < BlockRetryLimit; i++ {
		err := c.asMulti(flow)
		if err != nil {
			c.log.Warn("asmulti err will retry after 10 s", "err", err)
			time.Sleep(BlockInterval)
			continue
		} else {
			return nil
		}
	}

	return errors.New(fmt.Sprintf("asmulti reach limit symbol %s", flow.Symbol))
}

func (c *Connection) asMulti(flow *submodel.MultiEventFlow) error {
	gc := c.gcs[flow.Key]
	if gc == nil {
		panic(fmt.Sprintf("key disappear: %s, symbol: %s", hexutil.Encode(flow.Key.PublicKey), c.symbol))
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

func (c *Connection) SetToPayoutStashes(flow *submodel.BondReportedFlow, validatorFromStafi []types.AccountID) error {
	fullTargets := make([]types.AccountID, 0)
	if len(validatorFromStafi) != 0 {
		fullTargets = validatorFromStafi
	} else {
		exist, err := c.QueryStorage(config.StakingModuleId, config.StorageNominators, flow.Snap.Pool, nil, &fullTargets)
		if err != nil {
			return err
		}
		if !exist || len(fullTargets) == 0 {
			return TargetNotExistError
		}
	}

	eraBz, err := types.EncodeToBytes(flow.LastEra)
	if err != nil {
		return err
	}

	points := new(submodel.EraRewardPoints)
	exist, err := c.QueryStorage(config.StakingModuleId, config.StorageErasRewardPoints, eraBz, nil, points)
	if err != nil {
		return err
	}
	if !exist {
		c.log.Warn("earRewardPoints not exits", "LastEra", flow.LastEra, "symbol", flow.Symbol)
		flow.Stashes = nil
		return nil
	}

	targets := make([]types.AccountID, 0)
	for _, tgt := range fullTargets {
		pointedFlag := false
		for _, idv := range points.Individuals {
			if hexutil.Encode(tgt[:]) == hexutil.Encode(idv.Validator[:]) {
				pointedFlag = true
			}
		}

		if !pointedFlag {
			continue
		}

		ep := new(submodel.Exposure)
		exist, err := c.QueryStorage(config.StakingModuleId, config.StorageErasStakersClipped, eraBz, tgt[:], ep)
		if err != nil {
			return err
		}

		if !exist {
			c.log.Info("ErasStakersClipped not exits", "LastEra", flow.LastEra, "symbol", flow.Symbol, "Validator", hexutil.Encode(tgt[:]))
			continue
		}

		for _, other := range ep.Others {
			if bytes.Equal(other.Who[:], flow.Snap.Pool[:]) {
				targets = append(targets, tgt)
			}
		}
	}

	flow.Stashes = targets
	return nil
}

func (c *Connection) TryPayout(flow *submodel.BondReportedFlow) error {
	controllers := make([]types.AccountID, 0)

	idx, client := c.FoundIndexOfSubAccount(flow.SubAccounts)
	if idx == -1 || client == nil {
		return errors.New("not a sub account")
	}

	stashes := flow.Stashes
	if idx%2 != 0 {
		for i, j := 0, len(stashes)-1; i < j; i, j = i+1, j-1 {
			stashes[i], stashes[j] = stashes[j], stashes[i]
		}
	}

	method := config.MethodPayoutStakers
	for _, stash := range stashes {
		stashStr := hexutil.Encode(stash[:])

		var controller types.AccountID
		exist, err := c.QueryStorage(config.StakingModuleId, config.StorageBonded, stash[:], nil, &controller)
		if err != nil {
			return fmt.Errorf("get controller error: %s, stash: %s", err, stashStr)
		}
		if !exist {
			return fmt.Errorf("get controller not exist, stash: %s", stashStr)
		}
		controllerStr := hexutil.Encode(controller[:])

		ledger := new(submodel.StakingLedger)
		exist, err = c.QueryStorage(config.StakingModuleId, config.StorageLedger, controller[:], nil, ledger)
		if err != nil {
			return fmt.Errorf("get ledger error: %s, stash: %s", err, stashStr)
		}
		if !exist {
			return fmt.Errorf("ledger not exist, stash: %s, controller: %s", stashStr, controllerStr)
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
		} else {
			controllers = append(controllers, controller)
			if idx >= 4 {
				continue
			}

			ext, err := client.NewUnsignedExtrinsic(method, stash, flow.LastEra)
			if err != nil {
				return fmt.Errorf("NewUnsignedExtrinsic error: %s", err)
			}

			err = client.SignAndSubmitTx(ext)
			c.log.Info("SignAndSubmitTx result", "err", err)
		}
	}
	for _, con := range controllers {
		conStr := hexutil.Encode(con[:])
		retry := 0
		for {
			if retry >= BlockRetryLimit*10 {
				//reach limit
				return errors.New("query controller reach limit")
			}
			claimed, err := c.IsClaimed(con[:], flow.LastEra)
			if err != nil {
				c.log.Debug("query controller claimed err will retry", "controller", conStr, "retry", retry, "err", err)
				retry++
				time.Sleep(BlockInterval)
				continue
			}
			if !claimed {
				c.log.Debug("not claimed will retry query controller", "controller", conStr, "retry", retry)
				retry++
				time.Sleep(BlockInterval)
				continue
			}
			//has claimed
			break
		}
	}
	return nil
}

func (c *Connection) IsClaimed(controller []byte, lastEra uint32) (bool, error) {
	ledger := new(submodel.StakingLedger)
	exist, err := c.QueryStorage(config.StakingModuleId, config.StorageLedger, controller, nil, ledger)
	if err != nil {
		return false, fmt.Errorf("get ledger error: %s, controller: %s", err, hexutil.Encode(controller))
	}
	if !exist {
		return false, fmt.Errorf("ledger not exist, controller: %s", hexutil.Encode(controller))
	}

	claimed := false
	for _, claimedEra := range ledger.ClaimedRewards {
		if lastEra == claimedEra {
			claimed = true
			break
		}
	}
	return claimed, nil
}

func (c *Connection) submitSignature(param *submodel.SubmitSignatureParams) bool {
	for i := 0; i < BlockRetryLimit; i++ {
		c.log.Info("submitSignature on chain...")
		ext, err := c.gc.NewUnsignedExtrinsic(config.SubmitSignatures, param.Symbol, param.Era, param.Pool,
			param.TxType, param.ProposalId, param.Signature)
		err = c.gc.SignAndSubmitTx(ext)
		if err != nil {
			if err.Error() == TerminatedError.Error() {
				c.log.Error("submitSignature  met TerminatedError")
				return false
			}
			c.log.Warn("submitSignature error will retry", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}
		return true
	}
	return true
}
