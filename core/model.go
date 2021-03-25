package core

import (
	"fmt"
	"github.com/stafiprotocol/go-substrate-rpc-client/scale"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

type EvtLiquidityBond struct {
	AccountId types.AccountID
	Rsymbol   RSymbol
	BondId    types.Hash
}

type BondFlow struct {
	Key    *BondKey
	Record *BondRecord
	Reason BondReason
}

type BondKey struct {
	Rsymbol RSymbol
	BondId  types.Hash
}

type BondRecord struct {
	Bonder    types.AccountID
	Rsymbol   RSymbol
	Pubkey    types.Bytes
	Pool      types.Bytes
	Blockhash types.Bytes
	Txhash    types.Bytes
	Amount    types.U128
}

type RSymbol string

const (
	RFIS  = RSymbol("RFIS")
	RDOT  = RSymbol("RDOT")
	RKSM  = RSymbol("RKSM")
	RATOM = RSymbol("RATOM")
)

func (r *RSymbol) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		//*r = RFIS
		return fmt.Errorf("RSymbol decode error: %d", b)
	case 1:
		*r = RDOT
	case 2:
		*r = RKSM
	case 3:
		*r = RATOM
	default:
		return fmt.Errorf("RSymbol decode error: %d", b)
	}

	return nil
}

func (r RSymbol) Encode(encoder scale.Encoder) error {
	switch r {
	case RFIS:
		//return encoder.PushByte(0)
		return fmt.Errorf("RFIS not supported")
	case RDOT:
		return encoder.PushByte(1)
	case RKSM:
		return encoder.PushByte(2)
	case RATOM:
		return encoder.PushByte(3)
	default:
		return fmt.Errorf("RSymbol %s not supported", r)
	}
}

type BondReason string

const (
	BondReasonDefault = BondReason("Default")
	Pass              = BondReason("Pass")
	BlockhashUnmatch  = BondReason("BlockhashUnmatch")
	TxhashUnmatch     = BondReason("TxhashUnmatch")
	PubkeyUnmatch     = BondReason("PubkeyUnmatch")
	PoolUnmatch       = BondReason("PoolUnmatch")
	AmountUnmatch     = BondReason("AmountUnmatch")
)

func (br BondReason) Encode(encoder scale.Encoder) error {
	switch br {
	case Pass:
		return encoder.PushByte(0)
	case BlockhashUnmatch:
		return encoder.PushByte(1)
	case TxhashUnmatch:
		return encoder.PushByte(2)
	case PubkeyUnmatch:
		return encoder.PushByte(3)
	case PoolUnmatch:
		return encoder.PushByte(4)
	case AmountUnmatch:
		return encoder.PushByte(5)
	default:
		return fmt.Errorf("BondReason %s not supported", br)
	}
}

func (br *BondReason) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		*br = Pass
	case 1:
		*br = BlockhashUnmatch
	case 2:
		*br = TxhashUnmatch
	case 3:
		*br = PubkeyUnmatch
	case 4:
		*br = PoolUnmatch
	case 5:
		*br = AmountUnmatch
	default:
		return fmt.Errorf("BondReason decode error: %d", b)
	}

	return nil
}

type RproposalStatus string

const (
	Initiated = RproposalStatus("Initiated")
	Approved  = RproposalStatus("Approved")
	Rejected  = RproposalStatus("Rejected")
	Expired   = RproposalStatus("Expired")
)

func (r *RproposalStatus) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		*r = Initiated
	case 1:
		*r = Approved
	case 2:
		*r = Rejected
	case 3:
		*r = Expired
	default:
		return fmt.Errorf("RproposalStatus decode error: %d", b)
	}

	return nil
}

type VoteState struct {
	VotesFor     []types.AccountID
	VotesAgainst []types.AccountID
	Status       RproposalStatus
}

type Proposal struct {
	Call       types.Call
	Key        *BondKey
	MethodName string
}

// encode takes only nonce and call and encodes them for storage queries
func (p *Proposal) Encode() ([]byte, error) {
	return types.EncodeToBytes(struct {
		types.Hash
		types.Call
	}{p.Key.BondId, p.Call})
}

type PoolKey struct {
	Rsymbol RSymbol
	Pool    types.Bytes
}

type SignaturesKey struct {
	RSymbol    RSymbol
	Era        uint32
	Pool       []byte
	TxType     OriginalTx
	ProposalId []byte
}

type MultisigFlow struct {
	UpdatedData     *PoolUpdatedData
	LastVoterFlag   bool
	Threshold       uint16
	SubAccounts     []types.Bytes
	Key             *signature.KeyringPair
	Others          []types.AccountID
	TimePoint       *OptionTimePoint
	Opaque          []byte
	EncodeExtrinsic string
	CallHash        string
	NewMul          *EventNewMultisig
	Multisig        *Multisig
	MulExecuted     *EventMultisigExecuted
}

type BondReportFlow struct {
	ShotId        types.Hash
	Rsymbol       RSymbol
	Pool          []byte
	Era           uint32
	LastVoter     types.AccountID
	LastEra       uint32
	LastVoterFlag bool
	Active        types.UCompact
	Stashes       []types.AccountID
}

//type ToPayout struct {
//	Stash      types.AccountID
//	Era        uint32
//}

type PoolUpdatedData struct {
	Evt  *EvtEraPoolUpdated
	Snap *BondSnapshot
}

//EraPoolUpdated(RSymbol, u32, Vec<u8>, u128, u128, AccountId),
type EvtEraPoolUpdated struct {
	ShotId    types.Hash
	LastVoter types.AccountID
}

type BondSnapshot struct {
	Rsymbol   RSymbol
	Era       uint32
	Pool      []byte
	Bond      types.U128
	Unbond    types.U128
	Active    types.U128
	LastVoter types.AccountID
}

//SignaturesEnough(RSymbol, u32, Vec<u8>, OriginalTxType, Vec<u8>),
type EvtSignatureEnough struct {
	RSymbol    RSymbol
	Era        uint32
	Pool       []byte
	TxType     OriginalTx
	ProposalId []byte
}

type EventNewMultisig struct {
	Who, ID  types.AccountID
	CallHash types.Hash
}

type Multisig struct {
	When      types.TimePoint
	Deposit   types.U128
	Depositor types.AccountID
	Approvals []types.AccountID
}

type EventMultisigExecuted struct {
	Who       types.AccountID
	TimePoint types.TimePoint
	ID        types.AccountID
	CallHash  types.Hash
	Result    bool
}

type OriginalTx string

const (
	Transfer       = "Transfer"
	Bond           = "Bond"
	Unbond         = "Unbond"
	WithdrawUnbond = "WithdrawUnbond"
	ClaimRewards   = "ClaimRewards"
)

func (r *OriginalTx) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		*r = Transfer
	case 1:
		*r = Bond
	case 2:
		*r = Unbond
	case 3:
		*r = WithdrawUnbond
	case 4:
		*r = ClaimRewards
	default:
		return fmt.Errorf("OriginalTx decode error: %d", b)
	}

	return nil
}

func (r OriginalTx) Encode(encoder scale.Encoder) error {
	switch r {
	case Transfer:
		return encoder.PushByte(0)
	case Bond:
		return encoder.PushByte(1)
	case Unbond:
		return encoder.PushByte(2)
	case WithdrawUnbond:
		return encoder.PushByte(3)
	case ClaimRewards:
		return encoder.PushByte(4)
	default:
		return fmt.Errorf("OriginalTx %s not supported", r)
	}
}

type SubmitSignatureParams struct {
	Symbol     RSymbol
	Era        types.U32
	Pool       types.Bytes
	TxType     OriginalTx
	ProposalId types.Bytes
	Signature  types.Bytes
}

type SubmitSignatures struct {
	Symbol     RSymbol
	Era        types.U32
	Pool       types.Bytes
	TxType     OriginalTx
	ProposalId types.Bytes
	Signature  []types.Bytes
}