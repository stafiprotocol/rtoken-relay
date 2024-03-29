// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package matic

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"

	stake_portal "github.com/stafiprotocol/rtoken-relay/bindings/StakeERC20Portal"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

// Frequency of polling for a new block
var (
	BlockDelay         = uint64(2)
	BlockRetryInterval = time.Second * 10
	EraInterval        = time.Minute * 2
	BlockRetryLimit    = 30
	GetRetryLimit      = 50
	ErrFatalPolling    = errors.New("listener block polling failed")
)

type listener struct {
	name       string
	symbol     core.RSymbol
	eraSeconds uint64
	eraOffset  int64
	startBlock uint64
	conn       *Connection
	router     chains.Router
	blockstore blockstore.Blockstorer
	log        core.Logger
	stop       <-chan int
	sysErr     chan<- error // Reports fatal error to core
}

// NewListener creates and returns a listener
func NewListener(name string, symbol core.RSymbol, opts map[string]interface{}, conn *Connection, startBlock uint64, bs blockstore.Blockstorer, log core.Logger, stop <-chan int, sysErr chan<- error) *listener {
	eraSeconds := opts["eraSeconds"]
	eraSecondsStr, ok := eraSeconds.(string)
	if !ok {
		panic("eraSeconds not string")
	}
	eraSecondsBig, ok := utils.StringToBigint(eraSecondsStr)
	if !ok {
		panic("eraSeconds is not digital string")
	}
	if eraSecondsBig.Sign() <= 0 {
		panic(fmt.Sprintf("wrong erablock: %s", eraSecondsBig))
	}

	eraOffset := opts["eraOffset"]
	eraOffsetStr, ok := eraOffset.(string)
	if !ok {
		panic("eraOffset not string")
	}
	eraOffsetBig, ok := utils.StringToBigint(eraOffsetStr)
	if !ok {
		panic("eraOffset is not digital string")
	}

	return &listener{
		name:       name,
		symbol:     symbol,
		eraSeconds: eraSecondsBig.Uint64(),
		eraOffset:  eraOffsetBig.Int64(),
		conn:       conn,
		startBlock: startBlock,
		blockstore: bs,
		log:        log,
		stop:       stop,
		sysErr:     sysErr,
	}
}

// sets the router
func (l *listener) setRouter(r chains.Router) {
	l.router = r
}

// start registers all subscriptions provided by the config
func (l *listener) start() error {
	l.log.Debug("Starting listener...")

	go func() {
		err := l.pollEras()
		if err != nil {
			l.log.Error("Polling eras failed", "err", err)
			l.sysErr <- err
		}
	}()

	go func() {
		err := l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
			l.sysErr <- err
		}
	}()

	return nil
}

func (l *listener) pollBlocks() error {
	l.log.Info("Polling blocks...")
	var willDealBlock = l.startBlock
	var retry = BlockRetryLimit
	for {
		select {
		case <-l.stop:
			l.log.Info("pollBlocks receive stop chan, will stop")
			return nil
		default:
			if retry <= 0 {
				return fmt.Errorf("pollBlocks reach retry limit ,symbol: %s", l.symbol)
			}

			latestBlk, err := l.conn.LatestBlock()
			if err != nil {
				l.log.Error("Failed to fetch latest blockNumber", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}
			// Sleep if the block we want comes after the most recently finalized block
			if willDealBlock+BlockDelay > latestBlk {
				if willDealBlock%100 == 0 {
					l.log.Trace("Block not yet finalized", "target", willDealBlock, "finalBlk", latestBlk)
				}
				time.Sleep(BlockRetryInterval)
				continue
			}

			err = l.processBlockEvents(willDealBlock)
			if err != nil {
				l.log.Error("Failed to process events in block", "block", willDealBlock, "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Write to blockstore
			err = l.blockstore.StoreBlock(new(big.Int).SetUint64(willDealBlock))
			if err != nil {
				l.log.Error("Failed to write to blockstore", "err", err)
				return err
			}
			willDealBlock++

			retry = BlockRetryLimit
		}
	}
}

func (l *listener) processBlockEvents(currentBlock uint64) error {
	if currentBlock%100 == 0 {
		l.log.Debug("processBlockEvents", "blockNum", currentBlock)
	}
	// stake event
	stakeIterator, err := l.conn.stakePortalContract.FilterStake(&bind.FilterOpts{
		Start:   currentBlock,
		End:     &currentBlock,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}

	err = l.processStakeEvent(stakeIterator, false, [32]byte{}, common.Hash{}, common.Address{})
	if err != nil {
		return err
	}

	// recover event
	recoverStakeIterator, err := l.conn.stakePortalContract.FilterRecoverStake(&bind.FilterOpts{
		Start:   currentBlock,
		End:     &currentBlock,
		Context: context.Background(),
	})
	if err != nil {
		return err
	}

	for recoverStakeIterator.Next() {
		oldTx, err := l.conn.conn.TransactionReceipt(recoverStakeIterator.Event.TxHash)
		if err != nil {
			return err
		}
		oldBlockNumber := oldTx.BlockNumber.Uint64()

		oldStakeIterator, err := l.conn.stakePortalContract.FilterStake(&bind.FilterOpts{
			Start:   oldBlockNumber,
			End:     &oldBlockNumber,
			Context: context.Background(),
		})
		if err != nil {
			return err
		}

		recoverBlockHash := recoverStakeIterator.Event.Raw.BlockHash
		recoverTxhash := recoverStakeIterator.Event.Raw.TxHash
		recoverTxIndex := recoverStakeIterator.Event.Raw.TxIndex
		recoverTx, _, err := l.conn.conn.TransactionByHash(context.Background(), recoverTxhash)
		if err != nil {
			return err
		}

		recoverTxSender, err := l.conn.conn.TransactionSender(context.Background(), recoverTx, recoverBlockHash, recoverTxIndex)
		if err != nil {
			return err
		}

		err = l.processStakeEvent(oldStakeIterator, true, recoverStakeIterator.Event.StafiRecipient, recoverStakeIterator.Event.TxHash, recoverTxSender)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *listener) processStakeEvent(stakeIterator *stake_portal.StakeERC20PortalStakeIterator, isRecover bool, stafiRecipient [32]byte, oldTxHash common.Hash, newTxSender common.Address) error {

	for stakeIterator.Next() {
		willUseStafiRecipient := stakeIterator.Event.StafiRecipient

		if isRecover {
			willUseStafiRecipient = stafiRecipient
			if stakeIterator.Event.Raw.TxHash != oldTxHash {
				continue
			}

			// check tx sender
			if stakeIterator.Event.Staker != newTxSender {
				continue
			}
		}

		l.log.Info("find stake event", "stafiRecipient", hex.EncodeToString(willUseStafiRecipient[:]), "pool", stakeIterator.Event.StakePool.String(), "staker", stakeIterator.Event.Staker.String())

		// check stafi recipient ?
		flow := submodel.ExeLiquidityBondAndSwapFlow{
			Pool:           types.NewBytes(stakeIterator.Event.StakePool.Bytes()),
			Blockhash:      types.NewBytes(stakeIterator.Event.Raw.BlockHash.Bytes()),
			Txhash:         types.NewBytes(stakeIterator.Event.Raw.TxHash.Bytes()),
			Amount:         types.NewU128(*stakeIterator.Event.Amount),
			Symbol:         core.RMATIC,
			StafiRecipient: types.NewAccountID(willUseStafiRecipient[:]),
			DestRecipient:  types.NewBytes(stakeIterator.Event.DestRecipient[:]),
			DestId:         types.U8(stakeIterator.Event.ChainId),
			Reason:         submodel.Pass,
		}

		err := l.processExeLiquidityBondAndSwap(&flow)
		if err != nil {
			return err
		}

		// wait until it is dealed
		retry := 0
		for {
			if retry > GetRetryLimit {
				return fmt.Errorf("GetBondStateFlow reach retry limit")
			}

			bondState, err := l.getbondStateFromStafi(core.RMATIC, flow.Blockhash, flow.Txhash)
			if err != nil {
				l.log.Warn("getbondStateFromStafi", "err", err)
				retry++
				time.Sleep(WaitInterval)
				continue
			}

			if bondState == submodel.Default {
				retry++
				time.Sleep(WaitInterval)
				continue
			}
			break
		}
	}
	return nil
}

func (l *listener) getbondStateFromStafi(symbol core.RSymbol, blockHash, txHash types.Bytes) (submodel.BondState, error) {
	getBondStateFlow := submodel.GetBondStateFlow{
		Symbol:    core.RMATIC,
		BlockHash: blockHash,
		TxHash:    txHash,
		BondState: make(chan submodel.BondState, 1),
	}
	msg := &core.Message{
		Source: core.RMATIC, Destination: core.RFIS,
		Reason: core.GetBondState, Content: &getBondStateFlow}
	err := l.submitMessage(msg)
	if err != nil {
		return submodel.Default, err
	}

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		return submodel.Default, fmt.Errorf("get bond state from stafi timeout")
	case bs := <-getBondStateFlow.BondState:
		return bs, nil
	}
}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.cfg.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before continuing to the next block.
func (l *listener) pollEras() error {
	l.log.Info("Polling eras...")
	var retry = BlockRetryLimit
	for {
		select {
		case <-l.stop:
			l.log.Info("get stop signal, stop pool blocks")
			return nil
		default:
			// No more retries, goto next block
			if retry <= 0 {
				l.log.Error("Polling eras failed, retries exceeded")
				return ErrFatalPolling
			}

			latestBlockTimestamp, err := l.conn.LatestBlockTimestamp()
			if err != nil {
				l.log.Warn("Unable to get latest block", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			era := int64(latestBlockTimestamp/l.eraSeconds) + l.eraOffset
			if era <= 0 {
				return fmt.Errorf("era must > 0: %d", era)
			}

			err = l.processEra(uint32(era))
			if err != nil {
				l.log.Warn("processEra failed", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}
			time.Sleep(EraInterval)
			retry = BlockRetryLimit
		}
	}
}

func (l *listener) processEra(era uint32) error {
	msg := &core.Message{Destination: core.RFIS, Reason: core.NewEra, Content: era}
	return l.submitMessage(msg)
}

func (l *listener) processExeLiquidityBondAndSwap(flow *submodel.ExeLiquidityBondAndSwapFlow) error {
	msg := &core.Message{Destination: core.RFIS, Reason: core.ExeLiquidityBondAndSwap, Content: flow}
	return l.submitMessage(msg)
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (l *listener) submitMessage(m *core.Message) error {
	m.Source = l.symbol
	return l.router.Send(m)
}
