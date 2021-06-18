// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package matic

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/bindings/MaticToken"
	"github.com/stafiprotocol/rtoken-relay/bindings/Multisig"
	"github.com/stafiprotocol/rtoken-relay/bindings/StakeManager"
	"github.com/stafiprotocol/rtoken-relay/bindings/ValidatorShare"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

var (
	DefaultValue = big.NewInt(0)
	TxConfirmLimit   = 200
	TxRetryInterval  = time.Second * 2
	ErrNonceTooLow   = errors.New("nonce too low")
	ErrTxUnderpriced = errors.New("replacement transaction underpriced")
	ZeroAddress      = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

const (
	TxRetryLimit = 10
)

type Connection struct {
	url          string
	symbol       core.RSymbol
	conn         *ethereum.Client
	keys         []*secp256k1.Keypair
	publicKeys   [][]byte
	log          log15.Logger
	stop         <-chan int

	stakeManager *StakeManager.StakeManager
	multisig     *Multisig.Multisig
	maticTokenContract common.Address
	maticToken *MaticToken.MaticToken
	multiSendContract common.Address
}

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {
	log.Info("NewClient", "name", cfg.Name, "KeystorePath", cfg.KeystorePath, "Endpoint", cfg.Endpoint)

	var key *secp256k1.Keypair
	keys := make([]*secp256k1.Keypair, 0)
	publicKeys := make([][]byte, 0)
	acSize := len(cfg.Accounts)
	for i := 0; i < acSize; i++ {
		kpI, err := keystore.KeypairFromAddress(cfg.Accounts[i], keystore.EthChain, cfg.KeystorePath, cfg.Insecure)
		if err != nil {
			return nil, err
		}
		kp, _ := kpI.(*secp256k1.Keypair)
		pk := utils.PublicKeyFromKeypair(kp)

		if i == acSize - 1 {
			key = kp
		} else {
			keys = append(keys, kp)
			publicKeys = append(publicKeys, pk)
		}
	}

	if len(keys) == 0 {
		return nil, errors.New("no keys")
	}

	conn := ethereum.NewClient(cfg.Endpoint, key, log, big.NewInt(0), big.NewInt(0))
	if err := conn.Connect(); err != nil {
		return nil, err
	}

	stakeManager, err := initStakeManager(cfg.Opts["StakeManagerContract"], conn.Client())
	if err != nil {
		return nil, err
	}

	multisig, err := initMultisig(cfg.Opts["MultisigContract"], conn.Client())
	if err != nil {
		return nil, err
	}

	matic, maticAddr, err := initMaticToken(cfg.Opts["MaticTokenContract"], conn.Client())
	if err != nil {
		return nil, err
	}

	multiSendAddr, err := initMultisend(cfg.Opts["MultiSendContract"], conn.Client())
	if err != nil {
		return nil, err
	}

	return &Connection{
		url:          cfg.Endpoint,
		symbol:       cfg.Symbol,
		conn:         conn,
		keys:         keys,
		publicKeys:   publicKeys,
		log:          log,
		stop:         stop,
		stakeManager: stakeManager,
		multisig:     multisig,
		maticTokenContract: maticAddr,
		maticToken: matic,
		multiSendContract: multiSendAddr,
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

func (c *Connection) Address() string {
	return c.conn.Keypair().Address()
}

func (c *Connection) TransferVerify(r *submodel.BondRecord) (submodel.BondReason, error) {
	return c.conn.TransferVerify(r, c.maticTokenContract)
}

func (c *Connection) FoundKey(accounts []types.Bytes) *secp256k1.Keypair {
	for _, ac := range accounts {
		for _, key := range c.keys {
			if key.Address() == hexutil.Encode(ac) {
				return key
			}
		}
	}
	return nil
}

func (c *Connection) Validator(validatorId *big.Int) (common.Address, error) {
	validator, err := c.stakeManager.Validators(nil, validatorId)
	if err != nil {
		return ZeroAddress, err
	}

	// todo check status
	return validator.ContractAddress, nil
}

func (c *Connection) UnbondNonce(shareAddress, user common.Address) (*big.Int, error) {
	share, err := ValidatorShare.NewValidatorShare(shareAddress, c.conn.Client())
	if err != nil {
		return nil, err
	}

	return share.UnbondNonces(nil, user)
}

func (c *Connection) BondOrUnbondCall(share common.Address, bond, unbond *big.Int) (submodel.OriginalTx, *ethmodel.MultiTransaction, error) {
	c.log.Info("BondOrUnbondCall", "bond", bond, "unbond", unbond)
	tx := &ethmodel.MultiTransaction{To: share, Value: DefaultValue}
	var err error

	if bond.Cmp(unbond) < 0 {
		c.log.Info("unbond larger than bond, UnbondCall")
		diff := big.NewInt(0).Sub(unbond, bond)
		tx.CallData, err = ValidatorShareAbi.Pack(SellVoucherNewMethodName, diff, diff)
		if err != nil {
			return submodel.OriginalTxDefault, nil, err
		}
		return submodel.OriginalUnbond, tx, nil
	} else if bond.Cmp(unbond) > 0 {
		c.log.Info("bond larger than unbond, BondCall")
		diff := big.NewInt(0).Sub(bond, unbond)
		tx.CallData, err = ValidatorShareAbi.Pack(BuyVoucherMethodName, diff, big.NewInt(0))
		if err != nil {
			return submodel.OriginalTxDefault, nil, err
		}
		return submodel.OriginalBond, tx, nil
	} else {
		c.log.Info("bond is equal to unbond, NoCall")
		return submodel.OriginalTxDefault, nil, substrate.BondEqualToUnbondError
	}
}

func (c *Connection) RestakeCall(share common.Address) (*ethmodel.MultiTransaction, error) {
	packed, err := ValidatorShareAbi.Pack(RestakeMethodName)
	if err != nil {
		return nil, err
	}

	return &ethmodel.MultiTransaction{
		To: share,
		Value: DefaultValue,
		CallData: packed,
	}, nil
}

func (c *Connection) WithdrawCall(share, pool common.Address) (*ethmodel.MultiTransaction, error) {
	UnbondNonce, err := c.UnbondNonce(share, pool)
	if err != nil {
		return nil, err
	}

	packed, err := ValidatorShareAbi.Pack(UnstakeClaimTokensNew, UnbondNonce)
	if err != nil {
		return nil, err
	}

	return &ethmodel.MultiTransaction{
		To: share,
		Value: DefaultValue,
		CallData: packed,
	}, nil
}

func (c *Connection) MessageToSign(tx *ethmodel.MultiTransaction) ([32]byte, error) {
	return c.multisig.MessageToSign(nil, tx.To, tx.Value, tx.CallData)
}

func (c *Connection) IsFirstSigner(msg, sig []byte) bool {
	sigPublicKey, err := crypto.Ecrecover(msg, sig)
	if err != nil {
		panic(err)
	}

	for _, key := range c.publicKeys {
		if bytes.Equal(sigPublicKey, key) {
			return true
		}
	}

	return false
}

func (c *Connection) BalanceOf(owner common.Address) (*big.Int, error) {
	return c.maticToken.BalanceOf(c.conn.CallOpts(), owner)
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
	for _, rec := range receives {
		addr := common.BytesToAddress(rec.Recipient)
		value := big.Int(rec.Value)
		calldata, err := MaticTokenAbi.Pack(TransferMethodName, addr, &value)
		if err != nil {
			return nil, err
		}

		bt := &ethmodel.BatchTransaction{
			Operation: config.Call,
			To: c.maticTokenContract,
			Value: DefaultValue,
			DataLength: big.NewInt(int64(len(calldata))),
			Data: calldata,
		}

		bts = append(bts, bt)
	}

	cd, err := MultiSendAbi.Pack(MultiSendMethodName, bts.Encode())
	if err != nil {
		return nil, err
	}

	return &ethmodel.MultiTransaction{
		To: c.multiSendContract,
		Value: DefaultValue,
		CallData: cd,
	}, nil
}

func (c *Connection) AsMulti(
	to common.Address,
	value *big.Int,
	calldata []byte,
	operation uint8,
	safeTxGas *big.Int,
	txHash [32]byte,
	vs []uint8, rs [][32]byte, ss [][32]byte) error {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-c.stop:
			return errors.New("AsMulting stopped")
		default:
			err := c.conn.LockAndUpdateOpts()
			if err != nil {
				c.log.Error("Failed to update tx opts", "err", err)
				continue
			}

			tx, err := c.multisig.ExecTransaction(
				c.conn.Opts(),
				to,
				value,
				calldata,
				config.Call,
				safeTxGas,
				txHash,
				vs,
				rs,
				ss,
			)
			c.conn.UnlockOpts()

			if err == nil {
				c.log.Info("multisig ExecTransaction result", "tx", tx.Hash(), "gasPrice", tx.GasPrice().String())
				return nil
			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				c.log.Debug("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				c.log.Warn("ExecTransaction failed", "err", err, "to", to.Hex(), "calldata", hexutil.Encode(calldata), "safeTxGas", safeTxGas)
				time.Sleep(TxRetryInterval)
			}
		}
	}
	return fmt.Errorf("multisig ExecTransaction failed, to: %s, calldata: %s, safeTxGas: %s", to.Hex(), hexutil.Encode(calldata), safeTxGas.String())
}


func (c *Connection) IsTxHashExecuted(hash common.Hash) (bool, error) {
	return c.multisig.ExecutedTxHashs(nil, hash)
}

func (c *Connection) CheckTxHash(hash common.Hash) error {
	latest, err := c.conn.LatestBlock()
	if err != nil {
		return err
	}

	for i := 0; i < TxConfirmLimit; i++  {
		executed, err := c.IsTxHashExecuted(hash)
		if err != nil {
			return err
		}

		if executed {
			return nil
		}

		err = c.conn.WaitForBlock(latest, big.NewInt(1))
		if err != nil {
			return err
		}
	}

	return errors.New("tx not executed")
}

func (c *Connection) Close() {
	c.conn.Close()
}
