// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package matic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/bindings/MaticToken"
	"github.com/stafiprotocol/rtoken-relay/bindings/Multisig"
	stake_portal "github.com/stafiprotocol/rtoken-relay/bindings/StakeERC20Portal"
	stake_portal_with_rate "github.com/stafiprotocol/rtoken-relay/bindings/StakeERC20PortalWithRate"
	"github.com/stafiprotocol/rtoken-relay/bindings/StakeManager"
	stake_portal_rate "github.com/stafiprotocol/rtoken-relay/bindings/StakePortalRate"
	"github.com/stafiprotocol/rtoken-relay/bindings/ValidatorShare"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
)

var (
	DefaultValue             = big.NewInt(0)
	TxRetryInterval          = time.Second * 2
	WaitInterval             = time.Second * 6
	TxHashLogInterval        = time.Second * 10
	ErrNonceTooLow           = errors.New("nonce too low")
	ErrWithdrawEpochNotMatch = errors.New("found nonce but unable to withdraw")
	ErrTxUnderpriced         = errors.New("replacement transaction underpriced")
	ZeroAddress              = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

const (
	TxRetryLimit        = 10
	TxConfirmLimit      = 50
	TxHashLogRetryLimit = 60
)

type Connection struct {
	url    string
	symbol core.RSymbol
	conn   *ethereum.Client     // submitter's client
	keys   []*secp256k1.Keypair // key of signer or submitter
	addrs  [][]byte             // address of signer or submitter
	log    core.Logger
	stop   <-chan int

	polygonConn *ethereum.Client

	stateManagerContract common.Address
	maticTokenContract   common.Address
	multiSendContract    common.Address

	stakeManager                *StakeManager.StakeManager
	maticToken                  *MaticToken.MaticToken
	stakePortalContract         *stake_portal.StakeERC20Portal
	stakePortalWithRateContract *stake_portal_with_rate.StakeERC20PortalWithRate

	polygonStakePortalRateContract *stake_portal_rate.StakePortalRate
	isMainnet                      bool
}

func NewConnection(cfg *core.ChainConfig, log core.Logger, stop <-chan int) (*Connection, error) {
	log.Info("NewClient", "name", cfg.Name, "KeystorePath", cfg.KeystorePath, "Endpoint", cfg.Endpoint)

	var key *secp256k1.Keypair
	keys := make([]*secp256k1.Keypair, 0)
	addrs := make([][]byte, 0)
	acSize := len(cfg.Accounts)
	for i := 0; i < acSize; i++ {
		kpI, err := keystore.KeypairFromAddress(cfg.Accounts[i], keystore.EthChain, cfg.KeystorePath, cfg.Insecure)
		if err != nil {
			return nil, err
		}
		kp, ok := kpI.(*secp256k1.Keypair)
		if !ok {
			return nil, fmt.Errorf("keypair cast failed")
		}

		keys = append(keys, kp)
		addrs = append(addrs, kp.CommonAddress().Bytes())

		// use the last one as submitter
		if i == acSize-1 {
			key = kp
		}
	}

	if len(keys) == 0 {
		return nil, errors.New("no keys")
	}

	client, err := ethclient.Dial(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	conn := ethereum.NewClient(cfg.Endpoint, key, log, big.NewInt(0), big.NewInt(0))
	if err := conn.Connect(); err != nil {
		return nil, err
	}

	stakeManager, stateManagerContract, err := initStakeManager(cfg.Opts["StakeManagerContract"], client)
	if err != nil {
		return nil, err
	}
	matic, maticAddr, err := initMaticToken(cfg.Opts["MaticTokenContract"], client)
	if err != nil {
		return nil, err
	}
	multiSendAddr, err := initMultisend(cfg.Opts["MultiSendContract"], client)
	if err != nil {
		return nil, err
	}
	stakePortal, err := initStakePortal(cfg.Opts["StakePortalContract"], client)
	if err != nil {
		return nil, err
	}

	var stakePortalWithRate *stake_portal_with_rate.StakeERC20PortalWithRate
	if cfg.Opts["StakePortalWithRateContract"] != nil {
		stakePortalWithRate, err = initStakePortalWithRate(cfg.Opts["StakePortalWithRateContract"], client)
		if err != nil {
			return nil, err
		}
	}

	chainId, err := conn.Client().ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	var isMainnet bool
	if chainId.Int64() == 1 {
		isMainnet = true
	} else {
		isMainnet = false
	}

	// polygon
	polygonEndpointStr, ok := cfg.Opts["polygonEndpoint"].(string)
	if !ok {
		return nil, errors.New("polygonEndpoint not ok")
	}

	polygonConn := ethereum.NewClient(polygonEndpointStr, key, log, big.NewInt(0), big.NewInt(0))
	if err := polygonConn.Connect(); err != nil {
		return nil, err
	}
	polygonClient, err := ethclient.Dial(polygonEndpointStr)
	if err != nil {
		return nil, err
	}

	var stakePortalRate *stake_portal_rate.StakePortalRate
	if cfg.Opts["polygonStakePortalRateContract"] != nil {
		stakePortalRate, err = initPolygonStakePortalRate(cfg.Opts["polygonStakePortalRateContract"], polygonClient)
		if err != nil {
			return nil, err
		}
	}

	return &Connection{
		url:                            cfg.Endpoint,
		symbol:                         cfg.Symbol,
		conn:                           conn,
		keys:                           keys,
		addrs:                          addrs,
		log:                            log,
		stop:                           stop,
		stakeManager:                   stakeManager,
		stateManagerContract:           stateManagerContract,
		maticTokenContract:             maticAddr,
		maticToken:                     matic,
		multiSendContract:              multiSendAddr,
		stakePortalContract:            stakePortal,
		stakePortalWithRateContract:    stakePortalWithRate,
		isMainnet:                      isMainnet,
		polygonConn:                    polygonConn,
		polygonStakePortalRateContract: stakePortalRate,
	}, nil
}

// LatestBlock returns the latest block from the current chain
func (c *Connection) LatestBlock() (uint64, error) {
	blk, err := c.conn.LatestBlock()
	if err != nil {
		return 0, err
	}

	return blk.Uint64(), nil
}

func (c *Connection) LatestBlockTimestamp() (uint64, error) {
	blkTime, err := c.conn.LatestBlockTimestamp()
	if err != nil {
		return 0, err
	}

	return blkTime, nil
}

// use conn address as blockstore use address
func (c *Connection) BlockStoreUseAddress() string {
	return c.conn.Keypair().Address()
}

func (c *Connection) TransferVerify(r *submodel.BondRecord) (result submodel.BondReason, err error) {
	for i := 0; i < 10; i++ {
		result, err = c.conn.TransferVerifyERC20(r, c.maticTokenContract)
		if err != nil {
			return
		}

		if result == submodel.BlockhashUnmatch {
			time.Sleep(5 * time.Second)
			continue
		}
	}

	return
}

func (c *Connection) FoundKey(accounts []types.Bytes) *secp256k1.Keypair {
	for _, ac := range accounts {
		for _, key := range c.keys {
			//c.log.Info("FoundKey", "keyAddress", key.Address(), "Account", hexutil.Encode(ac))
			// don't use address string as case problem exist
			if bytes.Equal(key.CommonAddress().Bytes(), ac) {
				return key
			}
		}
	}
	return nil
}

func (c *Connection) GetValidator(validatorId *big.Int) (common.Address, error) {
	c.log.Info("GetValidator", "validatorId", validatorId, "stateManagerContract", c.stateManagerContract)

	valFlag, err := c.stakeManager.IsValidator(nil, validatorId)
	if err != nil {
		return ZeroAddress, err
	}

	if !valFlag {
		return ZeroAddress, fmt.Errorf("validatorId: %s is not validator", validatorId.String())
	}

	validator, err := c.stakeManager.Validators(nil, validatorId)
	if err != nil {
		return ZeroAddress, err
	}

	return validator.ContractAddress, nil
}

func (c *Connection) UnbondNonce(shareAddress, user common.Address) (*big.Int, error) {
	share, err := ValidatorShare.NewValidatorShare(shareAddress, c.conn.Client())
	if err != nil {
		return nil, err
	}

	return share.UnbondNonces(nil, user)
}

func (c *Connection) Unbond(shareAddress, user common.Address, nonce *big.Int) (*big.Int, *big.Int, error) {
	share, err := ValidatorShare.NewValidatorShare(shareAddress, c.conn.Client())
	if err != nil {
		return nil, nil, err
	}

	unbond, err := share.UnbondsNew(nil, user, nonce)
	if err != nil {
		return nil, nil, err
	}

	return unbond.Shares, unbond.WithdrawEpoch, nil
}

func (c *Connection) BondOrUnbondCall(share common.Address, bond, unbond, leastBond *big.Int) (submodel.OriginalTx, *ethmodel.MultiTransaction, error) {
	c.log.Info("BondOrUnbondCall", "bond", bond, "unbond", unbond)
	tx := &ethmodel.MultiTransaction{To: share, Value: DefaultValue, Operation: ethmodel.Call}
	var err error

	if bond.Cmp(unbond) < 0 {
		c.log.Info("unbond larger than bond, UnbondCall")
		diff := big.NewInt(0).Sub(unbond, bond)
		tx.CallData, err = ValidatorShareAbi.Pack(SellVoucherNewMethodName, diff, diff)
		if err != nil {
			return submodel.OriginalTxDefault, nil, err
		}
		tx.SafeTxGas = BuyVoucherSafeTxGas
		tx.TotalGas = BuyVoucherSafeTxGas
		return submodel.OriginalUnbond, tx, nil
	} else if bond.Cmp(unbond) > 0 {
		diff := big.NewInt(0).Sub(bond, unbond)
		if diff.Cmp(leastBond) <= 0 {
			c.log.Info("diff is smaller than leastBond, NoCall", "bond", bond, "leastBond", leastBond)
			return submodel.OriginalTxDefault, nil, substrate.ErrorDiffSmallerThanLeast
		}

		c.log.Info("bond larger than unbond, BondCall")
		tx.CallData, err = ValidatorShareAbi.Pack(BuyVoucherMethodName, diff, big.NewInt(0))
		if err != nil {
			return submodel.OriginalTxDefault, nil, err
		}
		tx.SafeTxGas = SellVoucherNewSafeTxGas
		tx.TotalGas = SellVoucherNewSafeTxGas
		return submodel.OriginalBond, tx, nil
	} else {
		c.log.Info("bond is equal to unbond, NoCall")
		return submodel.OriginalTxDefault, nil, substrate.ErrorBondEqualToUnbond
	}
}

func (c *Connection) Claimable(shareAddress, user common.Address) (bool, error) {
	share, err := ValidatorShare.NewValidatorShare(shareAddress, c.conn.Client())
	if err != nil {
		return false, err
	}

	min, err := share.MinAmount(nil)
	if err != nil {
		return false, err
	}

	reward, err := share.GetLiquidRewards(nil, user)
	if err != nil {
		return false, err
	}

	return reward.Cmp(min) >= 0, nil
}

func (c *Connection) RestakeCall(share common.Address) (*ethmodel.MultiTransaction, error) {
	packed, err := ValidatorShareAbi.Pack(RestakeMethodName)
	if err != nil {
		return nil, err
	}

	return &ethmodel.MultiTransaction{
		To:        share,
		Value:     DefaultValue,
		CallData:  packed,
		Operation: ethmodel.Call,
		SafeTxGas: RestakeSafeTxGas,
		TotalGas:  RestakeSafeTxGas,
	}, nil
}

func (c *Connection) WithdrawNonce(share, pool common.Address) (*big.Int, error) {
	c.log.Debug("WithdrawNonce", "share", share, "pool", pool)

	latestNonce, err := c.UnbondNonce(share, pool)
	if err != nil {
		return nil, err
	}

	if latestNonce.Uint64() == 0 {
		return big.NewInt(0), nil
	}

	currentEpoch, err := c.stakeManager.Epoch(nil)
	if err != nil {
		return nil, err
	}

	delay, err := c.stakeManager.WITHDRAWALDELAY(nil)
	if err != nil {
		return nil, err
	}
	c.log.Debug("WithdrawNonce", "currentEpoch", currentEpoch, "delay", delay)

	ableIdx := uint64(0)
	unableIdx := uint64(0)

	for i := latestNonce.Uint64(); i >= 1; i-- {
		nonce := big.NewInt(int64(i))
		shares, withdrawEpoch, err := c.Unbond(share, pool, nonce)
		if err != nil {
			return nil, err
		}

		if shares.Uint64() == 0 {
			break
		}

		withdrawEpoch.Add(withdrawEpoch, delay)
		if withdrawEpoch.Cmp(currentEpoch) <= 0 {
			ableIdx = i
		} else {
			unableIdx = i
		}
	}

	if ableIdx == 0 && unableIdx != 0 {
		c.log.Debug(fmt.Sprintf("found nonce but unable to withdraw, ableIdx: %d, unableIdx: %d", ableIdx, unableIdx))
		return nil, ErrWithdrawEpochNotMatch
	}

	return big.NewInt(int64(ableIdx)), nil
}

func (c *Connection) WithdrawCall(share, pool common.Address, nonce *big.Int) (*ethmodel.MultiTransaction, error) {
	packed, err := ValidatorShareAbi.Pack(UnstakeClaimTokensNew, nonce)
	if err != nil {
		return nil, err
	}

	return &ethmodel.MultiTransaction{
		To:        share,
		Value:     DefaultValue,
		CallData:  packed,
		Operation: ethmodel.Call,
		SafeTxGas: WithdrawTxGas,
		TotalGas:  WithdrawTxGas,
	}, nil
}

func (c *Connection) MessageToSign(tx *ethmodel.MultiTransaction, pool common.Address, txHash common.Hash) common.Hash {
	return tx.MessageToSign(txHash, pool)
}

// selector one to submit signature
func (c *Connection) IsSubmitterOfEra(era uint32, subAccounts []types.Bytes) bool {
	len := uint32(len(subAccounts))
	index := era % len
	acc := subAccounts[index]

	for _, addr := range c.addrs {
		if bytes.Equal(acc, addr) {
			return true
		}
	}

	return false
}

func (c *Connection) VerifySigs(msg []byte, sigs [][]byte, pool common.Address) error {
	multi, err := Multisig.NewMultisig(pool, c.conn.Client())
	if err != nil {
		return err
	}

	for _, sig := range sigs {
		sigPublicKey, err := crypto.Ecrecover(msg, sig)
		if err != nil {
			panic(err)
		}

		pk, err := crypto.UnmarshalPubkey(sigPublicKey)
		if err != nil {
			return err
		}
		owner := crypto.PubkeyToAddress(*pk)
		isOwner, err := multi.IsOwner(nil, owner)
		if err != nil {
			return err
		}

		if !isOwner {
			return fmt.Errorf("VerifySig: not owner, address: %s, msg: %s, sig: %s, pool: %s", owner.Hex(), hexutil.Encode(msg), hexutil.Encode(sig), pool.Hex())
		}
	}

	return nil
}

func (c *Connection) BalanceOf(owner common.Address) (*big.Int, error) {
	return c.maticToken.BalanceOf(&bind.CallOpts{}, owner)
}

func (c *Connection) TotalStaked(share, staker common.Address) (*big.Int, error) {
	shr, err := ValidatorShare.NewValidatorShare(share, c.conn.Client())
	if err != nil {
		return nil, err
	}

	total, _, err := shr.GetTotalStake(nil, staker)
	if err != nil {
		return nil, err
	}

	return total, nil
}

func (c *Connection) TransferCall(receives []*submodel.Receive) (*ethmodel.MultiTransaction, error) {
	bts := make(ethmodel.BatchTransactions, 0)
	totalGas := big.NewInt(0)
	for _, rec := range receives {
		addr := common.BytesToAddress(rec.Recipient)
		value := big.Int(rec.Value)
		calldata, err := MaticTokenAbi.Pack(TransferMethodName, addr, &value)
		if err != nil {
			return nil, err
		}

		bt := &ethmodel.BatchTransaction{
			Operation:  uint8(ethmodel.Call),
			To:         c.maticTokenContract,
			Value:      DefaultValue,
			DataLength: big.NewInt(int64(len(calldata))),
			Data:       calldata,
		}
		totalGas.Add(totalGas, TransferTxGas)
		bts = append(bts, bt)
	}

	cd, err := MultiSendAbi.Pack(MultiSendMethodName, bts.Encode())
	if err != nil {
		return nil, err
	}

	return &ethmodel.MultiTransaction{
		To:        c.multiSendContract,
		Value:     DefaultValue,
		CallData:  cd,
		Operation: ethmodel.DelegateCall,
		SafeTxGas: TransferTxGas,
		TotalGas:  totalGas,
	}, nil
}

func (c *Connection) AsMulti(
	pool, to common.Address,
	value *big.Int,
	calldata []byte,
	operation uint8,
	safeTxGas *big.Int,
	totalGas *big.Int,
	txHash [32]byte,
	vs []uint8, rs [][32]byte, ss [][32]byte) error {
	multisig, err := Multisig.NewMultisig(pool, c.conn.Client())
	if err != nil {
		return err
	}
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-c.stop:
			return errors.New("AsMulting stopped")
		default:
			err := c.conn.LockAndUpdateOpts(totalGas, big.NewInt(0))
			if err != nil {
				c.log.Error("Failed to update tx opts", "err", err)
				continue
			}

			tx, err := multisig.ExecTransaction(
				c.conn.Opts(),
				to,
				value,
				calldata,
				operation,
				safeTxGas,
				txHash,
				vs,
				rs,
				ss,
			)
			c.conn.UnlockOpts()

			if err == nil {
				c.log.Info("multisig ExecTransaction result", "tx", tx.Hash(), "gasPrice", tx.GasPrice())
				return nil
			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				c.log.Debug("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				c.log.Warn("ExecTransaction failed", "err", err, "to", to, "calldata", hexutil.Encode(calldata), "safeTxGas", safeTxGas)
				time.Sleep(TxRetryInterval)
			}
		}
	}
	return fmt.Errorf("multisig ExecTransaction failed, to: %s, calldata: %s, safeTxGas: %s", to.Hex(), hexutil.Encode(calldata), safeTxGas.String())
}

func (c *Connection) TxHashState(hash common.Hash, pool common.Address) (ethmodel.TxHashState, error) {
	multisig, err := Multisig.NewMultisig(pool, c.conn.Client())
	if err != nil {
		return 0, err
	}

	state, err := multisig.TxHashs(nil, hash)
	return ethmodel.TxHashState(state), err
}

func (c *Connection) WaitTxHashSuccess(hash common.Hash, pool common.Address, txType submodel.OriginalTx) error {
	wait := func() error {
		latest, err := c.conn.LatestBlock()
		if err != nil {
			return err
		}
		return c.conn.WaitForBlock(latest, big.NewInt(4))
	}

	for i := 0; i < TxConfirmLimit; i++ {
		state, err := c.TxHashState(hash, pool)
		if err != nil {
			return err
		}

		switch state {
		case ethmodel.HashStateUnsubmit, ethmodel.HashStateFail:
			if txType == submodel.OriginalClaimRewards && state == ethmodel.HashStateFail {
				return nil
			}
			err = wait()
			if err != nil {
				return err
			}
		case ethmodel.HashStateSuccess:
			return nil
		default:
			return fmt.Errorf("unsupported HashState: %d", state)
		}
	}

	return errors.New("tx not executed")
}

func (c *Connection) StakedAndReward(txHash common.Hash, share, pool common.Address) (*big.Int, *big.Int, error) {
	reward, err := c.RewardByTxHash(txHash, pool)
	if err != nil {
		return nil, nil, err
	}

	staked, err := c.TotalStaked(share, pool)
	if err != nil {
		return nil, nil, err
	}

	return staked, reward, nil
}

// / txhash is not transaction hash but a param of multi.execTransaction
func (c *Connection) RewardByTxHash(txHash common.Hash, pool common.Address) (*big.Int, error) {
	for i := 0; i <= TxHashLogRetryLimit; i++ {
		multisig, err := Multisig.NewMultisig(pool, c.conn.Client())
		if err != nil {
			return nil, err
		}

		iter, err := multisig.FilterExecutionResult(nil)
		if err != nil {
			return nil, err
		}

		for {
			if !iter.Next() {
				break
			}
			evt := iter.Event
			if !bytes.Equal(evt.TxHash[:], txHash.Bytes()) || evt.Arg1 != uint8(ethmodel.HashStateSuccess) {
				continue
			}

			return c.RewardByTransactionHash(evt.Raw.TxHash, pool)
		}

		time.Sleep(TxHashLogInterval)
	}
	return nil, fmt.Errorf("RewardByTxHash: no log, txHash: %s, pool: %s", txHash.Hex(), pool.Hex())
}

func (c *Connection) RewardByTransactionHash(hash common.Hash, pool common.Address) (*big.Int, error) {
	c.log.Debug("RewardByTransactionHash", "hash", hash, "pool", pool)
	receipt, err := c.conn.TransactionReceipt(hash)
	if err != nil {
		return nil, err
	}

	for _, elog := range receipt.Logs {
		if !bytes.Equal(elog.Address.Bytes(), c.maticTokenContract.Bytes()) {
			continue
		}

		transfer, err := c.maticToken.ParseTransfer(*elog)
		if err != nil {
			continue
		}

		if !bytes.Equal(transfer.From.Bytes(), c.stateManagerContract.Bytes()) || !bytes.Equal(transfer.To.Bytes(), pool.Bytes()) {
			continue
		}

		return transfer.Value, nil
	}

	return big.NewInt(0), nil
}

func (c *Connection) IsMainnet() bool {
	return c.isMainnet
}

func (c *Connection) Close() {
	c.conn.Close()
}
