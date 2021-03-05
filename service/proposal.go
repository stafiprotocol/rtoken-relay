package service

import (
	"fmt"
	"time"

	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/substrate"
)

func (l *listener) resolveProposal(prop *conn.Proposal, inFavour bool) bool {
	for i := 0; i < BlockRetryLimit; i++ {
		// Ensure we only submit a vote if status of the proposal is Initiated
		valid, reason, err := l.proposalValid(prop)
		l.log.Info("ResolveProposal proposalValid", "valid", valid, "reason", reason, "method", prop.MethodName)
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
		ext, err := l.gsrpc.NewUnsignedExtrinsic(config.RacknowledgeProposal, prop.Key.Symbol, prop.Key.BondId, inFavour, prop.Call)
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

func (l *listener) newLiquidityBondProposal(key *conn.BondKey, reason conn.BondReason) (*conn.Proposal, error) {
	meta, err := l.gsrpc.GetLatestMetadata()
	if err != nil {
		return nil, err
	}
	method := config.ExecuteBondRecord

	call, err := types.NewCall(
		meta,
		method,
		key,
		reason,
	)
	if err != nil {
		return nil, err
	}

	return &conn.Proposal{call, key, method}, nil
}

func (l *listener) newUpdateEraProposal(key *conn.BondKey, newEra types.U32) (*conn.Proposal, error) {
	meta, err := l.gsrpc.GetLatestMetadata()
	if err != nil {
		return nil, err
	}
	method := config.SetChainEra

	call, err := types.NewCall(
		meta,
		method,
		key.Symbol,
		newEra,
	)
	if err != nil {
		return nil, err
	}

	return &conn.Proposal{call, key, method}, nil
}

func (l *listener) newSetPoolActiveProposal(key *conn.BondKey, symbol conn.RSymbol, newEra types.U32,
	pool types.Bytes, active types.U128) (*conn.Proposal, error) {
	meta, err := l.gsrpc.GetLatestMetadata()
	if err != nil {
		return nil, err
	}
	method := config.SetPoolActive

	call, err := types.NewCall(
		meta,
		method,
		symbol,
		newEra,
		pool,
		active,
	)
	if err != nil {
		return nil, err
	}

	return &conn.Proposal{call, key, method}, nil
}
