// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"math/big"
	"strings"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/substrate"
)

type listener struct {
	gsrpc      *substrate.GsrpcClient
	sarpc      *substrate.SarpcClient
	bs         blockstore.Blockstorer
	startBlock uint64
	validators map[conn.RSymbol]conn.Validator
	sysErr     chan<- error
	ctx        context.Context
	log        log15.Logger
}

// Frequency of polling for a new block
var BlockRetryInterval = time.Second * 5
var BlockRetryLimit = 5

func NewListener(ctx context.Context, sarpc *substrate.SarpcClient, gsrpc *substrate.GsrpcClient, bs blockstore.Blockstorer,
	startBlock uint64, validators map[conn.RSymbol]conn.Validator, sysErr chan<- error, log log15.Logger) *listener {
	return &listener{
		sarpc:      sarpc,
		gsrpc:      gsrpc,
		bs:         bs,
		startBlock: startBlock,
		validators: validators,
		sysErr:     sysErr,
		ctx:        ctx,
		log:        log,
	}
}

// Start creates the initial subscription for all events
func (l *listener) Start() error {
	// Check whether latest is less than starting block
	header, err := l.gsrpc.GetHeaderLatest()
	if err != nil {
		return err
	}

	if uint64(header.Number) < l.startBlock {
		return fmt.Errorf("starting block (%d) is greater than latest known block (%d)", l.startBlock, header.Number)
	}

	go func() {
		err = l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before returning with an error.
func (l *listener) pollBlocks() error {
	var currentBlock = l.startBlock
	var retry = BlockRetryLimit
	for {
		select {
		case <-l.ctx.Done():
			return errors.New("terminated")
		default:
			// No more retries, goto next block
			if retry == 0 {
				l.sysErr <- errors.New("event polling retries exceeded")
				return nil
			}

			// Get finalized block hash
			finalizedHash, err := l.gsrpc.GetFinalizedHead()
			if err != nil {
				l.log.Error("Failed to fetch finalized hash", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Get finalized block header
			finalizedHeader, err := l.gsrpc.GetHeader(finalizedHash)
			if err != nil {
				l.log.Error("Failed to fetch finalized header", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Sleep if the block we want comes after the most recently finalized block
			if currentBlock > uint64(finalizedHeader.Number) {
				l.log.Trace("Block not yet finalized", "target", currentBlock, "latest", finalizedHeader.Number)
				time.Sleep(BlockRetryInterval)
				continue
			}

			err = l.processEvents(currentBlock)
			if err != nil {
				l.log.Error("Failed to process events in block", "block", currentBlock, "err", err)
				if strings.Contains(err.Error(), "close 1006") || strings.Contains(err.Error(), "websocket: not connected") {
					l.log.Info("listener", "is webscoket connected", l.sarpc.IsConnected())
					l.sarpc.WebsocketReconnect()
				}
				retry--
				continue
			}

			// Write to blockstore
			err = l.bs.StoreBlock(big.NewInt(0).SetUint64(currentBlock))
			if err != nil {
				l.log.Error("Failed to write to blockstore", "err", err)
			}

			currentBlock++
			retry = BlockRetryLimit
		}
	}
}

func (l *listener) processEvents(blockNum uint64) error {
	if blockNum%100 == 0 {
		l.log.Debug("processEvents", "blockNum", blockNum)
	}

	allEvts, err := l.sarpc.GetEvents(blockNum)
	if err != nil {
		return err
	}

	lbes := filterEvts(allEvts, config.LiquidityBondModuleId, config.LiquidityBondEventId)
	l.log.Trace("processEvents", "LiquidityBondEventNum", len(lbes), "blockNum", blockNum)
	if err := l.processLiquidityBondEvents(lbes); err != nil {
		l.log.Error("processLiquidityBondEvent", "error", err)
		return err
	}

	return nil
}

func filterEvts(evts []*substrate.ChainEvent, moduleId, eventId string) []*substrate.ChainEvent {
	wanted := make([]*substrate.ChainEvent, 0)

	for _, evt := range evts {
		if evt.ModuleId != moduleId || evt.EventId != eventId {
			continue
		}

		wanted = append(wanted, evt)
	}
	return wanted
}

func (l *listener) processLiquidityBondEvents(evts []*substrate.ChainEvent) error {
	for _, evt := range evts {
		err := l.processLiquidityBondEvent(evt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *listener) processLiquidityBondEvent(evt *substrate.ChainEvent) error {
	lb, err := liquidityBondEventData(evt)
	if err != nil {
		return err
	}

	bondKey := lb.bondKey()
	bk, err := types.EncodeToBytes(bondKey)
	if err != nil {
		return err
	}

	br := new(conn.BondRecord)
	re, err := l.gsrpc.QueryStorage(config.LiquidityBondModuleId, config.StorageBondRecords, bk, nil, br)
	if err != nil {
		return err
	}

	if !re {
		return fmt.Errorf("unable to get bondrecord by bondkey: %+v", lb)
	}

	if br.Bonder != lb.accountId {
		return fmt.Errorf("bonder not matched: %s, %s", hexutil.Encode(br.Bonder[:]), hexutil.Encode(lb.accountId[:]))
	}

	val, ok := l.validators[br.Symbol]
	if !ok {
		return fmt.Errorf("no validator for symbol: %s", br.Symbol)
	}

	reason, err := val.TransferVerify(br)
	if err != nil {
		return err
	}

	brp, err := l.bondProposal(bondKey, reason)
	if err != nil {
		return err
	}

	result := l.resolveBondProposal(brp)
	l.log.Info("processLiquidityBondEvent", "result", result)

	return nil
}

func (l *listener) bondProposal(key *conn.BondKey, reason conn.BondReason) (*conn.BondRecordProposal, error) {
	meta, err := l.gsrpc.GetLatestMetadata()
	if err != nil {
		return nil, err
	}

	call, err := types.NewCall(
		meta,
		config.ExecuteBondRecord,
		key,
		reason,
	)
	if err != nil {
		return nil, err
	}

	prop := &conn.Proposal{call, key}
	return &conn.BondRecordProposal{prop, reason}, nil
}

func (l *listener) resolveBondProposal(p *conn.BondRecordProposal) bool {
	for i := 0; i < BlockRetryLimit; i++ {
		// Ensure we only submit a vote if status of the proposal is Initiated
		valid, reason, err := l.proposalValid(p.Prop)
		l.log.Info("ResolveBondProposal proposalValid", "valid", valid, "reason", reason)
		if err != nil {
			l.log.Error("Failed to assert proposal state", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}

		if !valid {
			l.log.Debug("Ignoring proposal", "reason", reason)
			return true
		}

		l.log.Info("Acknowledging proposal on chain...")
		//symbol: RSymbol, prop_id: T::Hash, in_favour: bool
		inFavour := p.Reason == conn.Pass
		ext, err := l.gsrpc.NewUnsignedExtrinsic(config.RacknowledgeProposal, p.Prop.Key.Symbol, p.Prop.Key.BondId, inFavour, p.Prop.Call)
		err = l.gsrpc.SignAndSubmitTx(ext)
		if err != nil {
			if err.Error() == substrate.TerminatedError.Error() {
				l.log.Error("Acknowledging proposal met TerminatedError")
				return false
			}
			l.log.Error("Acknowledging proposal error", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}
		return true
	}
	return true
}

func (l *listener) proposalValid(prop *conn.Proposal) (bool, string, error) {
	var state conn.VoteState

	symBz, err := types.EncodeToBytes(prop.Key.Symbol)
	if err != nil {
		return false, "", err
	}

	propBz, err := prop.Encode()
	if err != nil {
		return false, "", err
	}

	exists, err := l.gsrpc.QueryStorage(config.RtokenVoteModuleId, config.StorageVotes, symBz, propBz, &state)
	if err != nil {
		return false, "", err
	}

	if !exists {
		return true, "", nil
	}

	if state.Status != conn.Initiated {
		return false, fmt.Sprintf("CurrentVoteStatus: %s", state.Status), nil
	}

	if containsVote(state.VotesFor, types.NewAccountID(l.gsrpc.PublicKey())) {
		return false, "already voted for", nil
	}

	if containsVote(state.VotesAgainst, types.NewAccountID(l.gsrpc.PublicKey())) {
		return false, "already voted against", nil
	}

	return true, "", nil
}
