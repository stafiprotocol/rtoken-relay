package substrate

import (
	"bytes"
	"fmt"
	"time"

	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
)

func (c *Connection) LiquidityBondProposal(key *core.BondKey, reason core.BondReason) (*core.Proposal, error) {
	meta, err := c.LatestMetadata()
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

	return &core.Proposal{call, key, method}, nil
}

func (c *Connection) resolveProposal(prop *core.Proposal, inFavour bool) bool {
	for i := 0; i < BlockRetryLimit; i++ {
		// Ensure we only submit a vote if status of the proposal is Initiated
		valid, reason, err := c.proposalValid(prop)
		c.log.Info("ResolveProposal proposalValid", "valid", valid, "reason", reason, "method", prop.MethodName)
		if err != nil {
			c.log.Error("Failed to assert proposal state", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}

		if !valid {
			c.log.Debug("Ignoring proposal", "reason", reason)
			return true
		}

		c.log.Info("Acknowledging proposal on chain...")
		//rsymbol: RSymbol, prop_id: T::Hash, in_favour: bool
		ext, err := c.gc.NewUnsignedExtrinsic(config.RacknowledgeProposal, prop.Key.Rsymbol, prop.Key.BondId, inFavour, prop.Call)
		err = c.gc.SignAndSubmitTx(ext)
		if err != nil {
			if err.Error() == TerminatedError.Error() {
				c.log.Error("Acknowledging proposal met TerminatedError")
				return false
			}
			c.log.Error("Acknowledging proposal error", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}
		return true
	}
	return true
}

func (c *Connection) proposalValid(prop *core.Proposal) (bool, string, error) {
	var state core.VoteState

	symBz, err := types.EncodeToBytes(prop.Key.Rsymbol)
	if err != nil {
		return false, "", err
	}

	propBz, err := prop.Encode()
	if err != nil {
		return false, "", err
	}

	exists, err := c.QueryStorage(config.RtokenVoteModuleId, config.StorageVotes, symBz, propBz, &state)
	if err != nil {
		return false, "", err
	}

	if !exists {
		return true, "", nil
	}

	if state.Status != core.Initiated {
		return false, fmt.Sprintf("CurrentVoteStatus: %s", state.Status), nil
	}

	if containsVote(state.VotesFor, types.NewAccountID(c.gc.PublicKey())) {
		return false, "already voted for", nil
	}

	if containsVote(state.VotesAgainst, types.NewAccountID(c.gc.PublicKey())) {
		return false, "already voted against", nil
	}

	return true, "", nil
}

func (c *Connection) newUpdateEraProposal(key *core.BondKey, newEra types.U32) (*core.Proposal, error) {
	meta, err := c.LatestMetadata()
	if err != nil {
		return nil, err
	}
	method := config.SetChainEra

	call, err := types.NewCall(
		meta,
		method,
		key.Rsymbol,
		newEra,
	)
	if err != nil {
		return nil, err
	}

	return &core.Proposal{call, key, method}, nil
}

//
//func (l *listener) newSetPoolActiveProposal(key *conn.BondKey, rsymbol conn.RSymbol, newEra types.U32,
//	pool types.Bytes, active types.U128) (*conn.Proposal, error) {
//	meta, err := l.gsrpc.GetLatestMetadata()
//	if err != nil {
//		return nil, err
//	}
//	method := config.SetPoolActive
//
//	call, err := types.NewCall(
//		meta,
//		method,
//		rsymbol,
//		newEra,
//		pool,
//		active,
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	return &conn.Proposal{call, key, method}, nil
//}

func containsVote(votes []types.AccountID, voter types.AccountID) bool {
	for _, v := range votes {
		if bytes.Equal(v[:], voter[:]) {
			return true
		}
	}
	return false
}
