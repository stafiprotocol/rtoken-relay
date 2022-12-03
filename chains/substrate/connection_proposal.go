package substrate

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/stafiprotocol/rtoken-relay/utils"

	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
)

var ErrExpired = errors.New("proposal expired")
var ErrAlreadyVoted = errors.New("proposal already voted")

func (c *Connection) LiquidityBondProposal(flow *submodel.BondFlow) (*submodel.Proposal, error) {
	method := config.MethodExecuteBondRecord
	ci, err := c.sc.FindCallIndex(method)
	if err != nil {
		return nil, err
	}

	call, err := types.NewCallWithCallIndex(
		ci,
		method,
		flow.Symbol,
		flow.BondId,
		flow.Reason,
	)
	if err != nil {
		return nil, err
	}

	return &submodel.Proposal{Call: call, Symbol: flow.Symbol, BondId: flow.BondId, MethodName: method}, nil
}

// execute_bond_and_swap(origin, pool: Vec<u8>, blockhash: Vec<u8>, txhash: Vec<u8>, amount: u128, symbol: RSymbol, stafi_recipient: T::AccountId, dest_recipient: Vec<u8>, dest_id: ChainId, reason: BondReason)
func (c *Connection) ExeLiquidityBondAndSwapProposal(flow *submodel.ExeLiquidityBondAndSwapFlow) (*submodel.Proposal, error) {
	method := config.MethodExecuteBondAndSwap
	ci, err := c.sc.FindCallIndex(method)
	if err != nil {
		return nil, err
	}

	call, err := types.NewCallWithCallIndex(
		ci,
		method,
		flow.Pool,
		flow.Blockhash,
		flow.Txhash,
		flow.Amount,
		flow.Symbol,
		flow.StafiRecipient,
		flow.DestRecipient,
		flow.DestId,
		flow.Reason,
	)
	if err != nil {
		return nil, err
	}

	bondId := types.NewHash(flow.Txhash)

	return &submodel.Proposal{Call: call, Symbol: flow.Symbol, BondId: bondId, MethodName: method}, nil
}

func (c *Connection) CommonReportProposal(method string, symbol core.RSymbol, bondId, shotId types.Hash) (*submodel.Proposal, error) {
	ci, err := c.sc.FindCallIndex(method)
	if err != nil {
		return nil, err
	}

	call, err := types.NewCallWithCallIndex(
		ci,
		method,
		symbol,
		shotId,
	)
	if err != nil {
		return nil, err
	}

	return &submodel.Proposal{Call: call, Symbol: symbol, BondId: bondId, MethodName: method}, nil
}

func (c *Connection) NewBondReportProposal(flow *submodel.EraPoolUpdatedFlow) (*submodel.Proposal, error) {
	method := config.MethodNewBondReport
	ci, err := c.sc.FindCallIndex(method)
	if err != nil {
		return nil, err
	}

	call, err := types.NewCallWithCallIndex(
		ci,
		method,
		flow.Symbol,
		flow.ShotId,
		flow.BondCall.Action,
	)
	if err != nil {
		return nil, err
	}

	return &submodel.Proposal{Call: call, Symbol: flow.Symbol, BondId: flow.ShotId, MethodName: method}, nil
}

func (c *Connection) BondAndReportActiveProposal(flow *submodel.EraPoolUpdatedFlow) (*submodel.Proposal, error) {
	method := config.MethodBondAndReportActive
	ci, err := c.sc.FindCallIndex(method)
	if err != nil {
		return nil, err
	}

	call, err := types.NewCallWithCallIndex(
		ci,
		method,
		flow.Symbol,
		flow.ShotId,
		flow.BondCall.Action,
		types.NewU128(*flow.Active),
		types.NewU128(*flow.Reward),
	)
	if err != nil {
		return nil, err
	}

	return &submodel.Proposal{Call: call, Symbol: flow.Symbol, BondId: flow.ShotId, MethodName: method}, nil
}

func (c *Connection) ActiveReportProposal(flow *submodel.BondReportedFlow) (*submodel.Proposal, error) {
	method := config.MethodActiveReport
	ci, err := c.sc.FindCallIndex(method)
	if err != nil {
		return nil, err
	}

	call, err := types.NewCallWithCallIndex(
		ci,
		method,
		flow.Symbol,
		flow.ShotId,
		flow.Snap.Active,
	)
	if err != nil {
		return nil, err
	}

	return &submodel.Proposal{Call: call, Symbol: flow.Symbol, BondId: flow.ShotId, MethodName: method}, nil
}

func (c *Connection) NewActiveReportProposal(flow *submodel.BondReportedFlow) (*submodel.Proposal, error) {
	method := config.MethodNewActiveReport
	ci, err := c.sc.FindCallIndex(method)
	if err != nil {
		return nil, err
	}

	call, err := types.NewCallWithCallIndex(
		ci,
		method,
		flow.Symbol,
		flow.ShotId,
		flow.Snap.Active,
		flow.Unstaked,
	)
	if err != nil {
		return nil, err
	}

	return &submodel.Proposal{Call: call, Symbol: flow.Symbol, BondId: flow.ShotId, MethodName: method}, nil
}

func (c *Connection) resolveProposal(prop *submodel.Proposal, inFavour bool) bool {
	for i := 0; i < BlockRetryLimit; i++ {
		// Ensure we only submit a vote if status of the proposal is Initiated
		valid, reason, err := c.proposalValid(prop)
		if err != nil {
			switch err {
			case ErrExpired:
				newBondId := types.Hash(utils.BlakeTwo256(prop.BondId[:]))
				prop.BondId = newBondId
				continue
			case ErrAlreadyVoted:
				return true
			default:
				c.log.Error("Failed to assert proposal state", "err", err)
				time.Sleep(BlockRetryInterval)
				continue
			}
		}
		c.log.Info("ResolveProposal", "valid", valid, "reason", reason, "method", prop.MethodName)
		if !valid {
			c.log.Debug("Ignoring proposal", "reason", reason)
			return true
		}

		c.log.Debug("Acknowledging proposal on chain...")
		//symbol: RSymbol, prop_id: T::Hash, in_favour: bool
		ext, err := c.sc.NewUnsignedExtrinsic(config.MethodRacknowledgeProposal, prop.Symbol, prop.BondId, inFavour, prop.Call)
		if err != nil {
			c.log.Error("NewUnsignedExtrinsic error", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}
		err = c.sc.SignAndSubmitTx(ext)
		if err != nil {
			if err.Error() == ErrorTerminated.Error() {
				c.log.Error("Acknowledging proposal met ErrorTerminated")
				return false
			}
			c.log.Error("Acknowledging proposal error", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}
		return true
	}
	return false
}

func (c *Connection) proposalValid(prop *submodel.Proposal) (bool, string, error) {
	var state submodel.VoteState

	symBz, err := types.EncodeToBytes(prop.Symbol)
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

	if state.Status == submodel.Expired {
		return false, "", ErrExpired
	}

	if state.Status != submodel.Initiated {
		return false, fmt.Sprintf("CurrentVoteStatus: %s", state.Status), nil
	}

	if containsVote(state.VotesFor, types.NewAccountID(c.sc.PublicKey())) {
		return false, "already voted for", ErrAlreadyVoted
	}

	if containsVote(state.VotesAgainst, types.NewAccountID(c.sc.PublicKey())) {
		return false, "already voted against", ErrAlreadyVoted
	}

	return true, "", nil
}

func (c *Connection) SetChainEraProposal(symbol core.RSymbol, bondId types.Hash, newEra uint32) (*submodel.Proposal, error) {
	method := config.MethodSetChainEra
	ci, err := c.sc.FindCallIndex(method)
	if err != nil {
		return nil, err
	}

	call, err := types.NewCallWithCallIndex(
		ci,
		method,
		symbol,
		newEra,
	)
	if err != nil {
		return nil, err
	}

	return &submodel.Proposal{Call: call, Symbol: symbol, BondId: bondId, MethodName: method}, nil
}

// used by stafi chain
func (c *Connection) GetSelectedVoters(symbol core.RSymbol, newEra uint32) (types.AccountID, error) {
	th, err := c.RelayerThreshold(symbol)
	if err != nil {
		return types.AccountID{}, err
	}
	// select voters
	voters, err := c.GetNewChainEraProposalVoters(symbol, newEra)
	if err != nil {
		return types.AccountID{}, err
	}
	if len(voters) < int(th) {
		return types.AccountID{}, fmt.Errorf("newChainEra voters not enough")
	}

	usedVoters := voters[:th]
	sort.SliceStable(usedVoters, func(i, j int) bool {
		return bytes.Compare(usedVoters[i][:], usedVoters[j][:]) < 0
	})

	return usedVoters[0], nil
}

// used by stafi chain
func (c *Connection) GetNewChainEraProposalVoters(symbol core.RSymbol, newEra uint32) ([]types.AccountID, error) {
	eraBz, err := types.EncodeToBytes(newEra)
	if err != nil {
		c.log.Error("processNewEra EncodeToBytes error", "error", err, "newEra", newEra)
		return nil, err
	}
	bondId := types.Hash(utils.BlakeTwo256(eraBz))
	prop, err := c.SetChainEraProposal(symbol, bondId, newEra)
	if err != nil {
		return nil, err
	}

	var state submodel.VoteState

	symBz, err := types.EncodeToBytes(prop.Symbol)
	if err != nil {
		return nil, err
	}

	propBz, err := prop.Encode()
	if err != nil {
		return nil, err
	}

	exists, err := c.QueryStorage(config.RtokenVoteModuleId, config.StorageVotes, symBz, propBz, &state)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("not exists")
	}

	if len(state.VotesAgainst) != 0 {
		state.VotesFor = append(state.VotesFor, state.VotesAgainst...)
	}
	return state.VotesFor, nil
}

func containsVote(votes []types.AccountID, voter types.AccountID) bool {
	for _, v := range votes {
		if bytes.Equal(v[:], voter[:]) {
			return true
		}
	}
	return false
}
