package conn

import (
	"fmt"
	"github.com/stafiprotocol/go-substrate-rpc-client/scale"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

type BondVoteStatus string

const (
	VoteStatusInitiated = BondVoteStatus("Initiated")
	VoteStatusApproved  = BondVoteStatus("Approved")
	VoteStatusRejected  = BondVoteStatus("Rejected")
	VoteStatusExpired   = BondVoteStatus("Expired")
)

func (bvs *BondVoteStatus) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		*bvs = VoteStatusInitiated
	case 1:
		*bvs = VoteStatusApproved
	case 2:
		*bvs = VoteStatusRejected
	case 3:
		*bvs = VoteStatusExpired
	default:
		return fmt.Errorf("VoteStatus decode error: %d", b)
	}

	return nil
}

type RSymbol uint

const (
	RSymbolRdot = RSymbol(1)
)

//func (sym *RSymbol) Decode(decoder scale.Decoder) error {
//	b, err := decoder.ReadOneByte()
//	if err != nil {
//		return err
//	}
//
//	switch b {
//	case 1:
//		*sym = RSymbolRdot
//	default:
//		return fmt.Errorf("VoteStatus decode error: %d", b)
//	}
//
//	return nil
//}

type BondRecord struct {
	Bonder    types.AccountID `json:"bonder"`
	Rsymbol   RSymbol         `json:"rsymbol"`
	Pubkey    []byte          `json:"pubkey"`
	Signature []byte          `json:"signature"`
	Pool      []byte          `json:"pool"`
	Blockhash []byte          `json:"blockhash"`
	Txhash    []byte          `json:"txhash"`
	Amount    types.U128      `json:"amount"`
}

type BondVote struct {
	VotesFor     []types.AccountID
	VotesAgainst []types.AccountID
	Status       BondVoteStatus
}

type OpposeReason int

const (
	BLOCKHASH = OpposeReason(0)
	TXHASH    = OpposeReason(1)

	NoReason = OpposeReason(99)
)

type EvLiquidityBondEvt struct {
	Dest      types.AccountID
	Blockhash []byte
	Txhash    []byte
	Pubkey    []byte
}
