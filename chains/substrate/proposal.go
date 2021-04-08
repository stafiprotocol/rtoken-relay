package substrate

import (
	"bytes"
	"fmt"
	"time"

	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
)

func (c *Connection) InitLastVoterProposal(key *submodel.BondKey) (*submodel.Proposal, error) {
	meta, err := c.LatestMetadata()
	if err != nil {
		return nil, err
	}
	method := config.MethodInitLastVoter

	call, err := types.NewCall(
		meta,
		method,
	)
	if err != nil {
		return nil, err
	}

	return &submodel.Proposal{call, key, method}, nil
}

func (c *Connection) LiquidityBondProposal(key *submodel.BondKey, reason submodel.BondReason) (*submodel.Proposal, error) {
	meta, err := c.LatestMetadata()
	if err != nil {
		return nil, err
	}
	method := config.MethodExecuteBondRecord

	call, err := types.NewCall(
		meta,
		method,
		key,
		reason,
	)
	if err != nil {
		return nil, err
	}

	return &submodel.Proposal{call, key, method}, nil
}

func (c *Connection) CommonReportProposal(method string, key *submodel.BondKey, shotId types.Hash) (*submodel.Proposal, error) {
	meta, err := c.LatestMetadata()
	if err != nil {
		return nil, err
	}
	//method := config.MethodBondReport

	call, err := types.NewCall(
		meta,
		method,
		shotId,
	)
	if err != nil {
		return nil, err
	}

	return &submodel.Proposal{call, key, method}, nil
}

func (c *Connection) ActiveReportProposal(key *submodel.BondKey, shotId types.Hash, active types.U128) (*submodel.Proposal, error) {
	meta, err := c.LatestMetadata()
	if err != nil {
		return nil, err
	}
	method := config.MethodActiveReport

	call, err := types.NewCall(
		meta,
		method,
		shotId,
		active,
	)
	if err != nil {
		return nil, err
	}

	return &submodel.Proposal{call, key, method}, nil
}

func (c *Connection) resolveProposal(prop *submodel.Proposal, inFavour bool) bool {
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
		ext, err := c.gc.NewUnsignedExtrinsic(config.MethodRacknowledgeProposal, prop.Key.Rsymbol, prop.Key.BondId, inFavour, prop.Call)
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

func (c *Connection) proposalValid(prop *submodel.Proposal) (bool, string, error) {
	var state submodel.VoteState

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

	if state.Status != submodel.Initiated {
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

func (c *Connection) SetChainEraProposal(key *submodel.BondKey, newEra uint32) (*submodel.Proposal, error) {
	meta, err := c.LatestMetadata()
	if err != nil {
		return nil, err
	}
	method := config.MethodSetChainEra

	call, err := types.NewCall(
		meta,
		method,
		key.Rsymbol,
		newEra,
	)
	if err != nil {
		return nil, err
	}

	return &submodel.Proposal{call, key, method}, nil
}

func containsVote(votes []types.AccountID, voter types.AccountID) bool {
	for _, v := range votes {
		if bytes.Equal(v[:], voter[:]) {
			return true
		}
	}
	return false
}
