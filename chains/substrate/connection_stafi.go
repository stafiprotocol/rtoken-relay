// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
)

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

func (c *Connection) RelayerThreshold(symbol core.RSymbol) (uint32, error) {
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return 0, err
	}

	var th types.U32
	exists, err := c.QueryStorage(config.RTokenRelayersModuleId, config.StorageRelayerThreshold, symBz, nil, &th)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, errors.New("RelayerThreshold not exist")
	}
	return uint32(th), nil
}

var DefaultActiveChangeRateLimit = uint32(1e7)

func (c *Connection) ActiveChangeRateLimit(sym core.RSymbol) (uint32, error) {
	symBz, err := types.EncodeToBytes(sym)
	if err != nil {
		return 0, err
	}

	var PerBill types.U32
	exists, err := c.QueryStorage(config.RTokenLedgerModuleId, config.StorageActiveChangeRateLimit, symBz, nil, &PerBill)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, ErrorNotExist
	}

	return uint32(PerBill), nil
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
		return nil, ErrorNotExist
	}
	return validators, nil
}

func (sc *Connection) GetEraRate(symbol core.RSymbol, era types.U32) (rate uint64, err error) {
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return 0, err
	}
	eraIndex, err := types.EncodeToBytes(era)
	if err != nil {
		return 0, err
	}
	exists, err := sc.QueryStorage(config.RTokenRateModuleId, config.StorageEraRate, symBz, eraIndex, &rate)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrorNotExist
	}
	return rate, nil
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
	return hexutil.Encode(c.sc.PublicKey()) == hexutil.Encode(voter[:])
}

func (c *Connection) submitSignature(param *submodel.SubmitSignatureParams) bool {
	for i := 0; i < BlockRetryLimit; i++ {
		c.log.Info("submitSignature on chain...")
		ext, err := c.sc.NewUnsignedExtrinsic(config.SubmitSignatures, param.Symbol, param.Era, param.Pool,
			param.TxType, param.ProposalId, param.Signature)
		if err != nil {
			c.log.Warn("submitSignature error will retry", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}
		err = c.sc.SignAndSubmitTx(ext)
		if err != nil {
			if err.Error() == ErrorTerminated.Error() {
				c.log.Error("submitSignature  met ErrorTerminated")
				return false
			}
			c.log.Warn("submitSignature error will retry", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}
		return true
	}
	return false
}

func (c *Connection) poolThreshold(symbol core.RSymbol, pool []byte) (uint16, error) {
	poolBz, err := types.EncodeToBytes(pool)
	if err != nil {
		return 0, err
	}

	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return 0, err
	}

	var threshold uint16
	exist, err := c.QueryStorage(config.RTokenLedgerModuleId, config.StorageMultiThresholds, symBz, poolBz, &threshold)
	if err != nil {
		return 0, err
	}
	if !exist {
		return 0, ErrorNotExist
	}
	return threshold, nil
}

func (c *Connection) bondState(symbol core.RSymbol, blockHash, txHash []byte) (submodel.BondState, error) {
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return submodel.Default, err
	}

	bsk := submodel.BondStateKey{BlockHash: blockHash, TxHash: txHash}
	bskBz, err := types.EncodeToBytes(bsk)
	if err != nil {
		return submodel.Default, err
	}

	var bs submodel.BondState
	exist, err := c.QueryStorage(config.RTokenSeriesModuleId, config.StorageBondStates, symBz, bskBz, &bs)
	if err != nil {
		return submodel.Default, err
	}

	if !exist {
		return submodel.Default, ErrorNotExist
	}
	return bs, nil
}

func (c *Connection) getSubmitSignatures(symbol core.RSymbol, era uint32, pool []byte, txType submodel.OriginalTx, proposalId []byte) ([]types.Bytes, error) {
	symBz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return nil, err
	}
	sk := submodel.SignaturesKey{
		Era:        era,
		Pool:       pool,
		TxType:     txType,
		ProposalId: proposalId,
	}

	skBz, err := types.EncodeToBytes(sk)
	if err != nil {
		return nil, err
	}

	var sigs []types.Bytes
	exist, err := c.QueryStorage(config.RTokenSeriesModuleId, config.StorageSignatures, symBz, skBz, &sigs)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, ErrorNotExist
	}

	return sigs, nil
}
