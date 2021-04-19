package submodel

import (
	"fmt"

	scalecodec "github.com/itering/scale.go"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/stafiprotocol/go-substrate-rpc-client/scale"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	chainTypes "github.com/stafiprotocol/rtoken-relay/types"
)

type EvtLiquidityBond struct {
	AccountId types.AccountID
	Symbol    core.RSymbol
	BondId    types.Hash
}

type BondFlow struct {
	Symbol core.RSymbol
	BondId types.Hash
	Record *BondRecord
	Reason BondReason
}

type BondRecord struct {
	Bonder    types.AccountID
	Symbol    core.RSymbol
	Pubkey    types.Bytes
	Pool      types.Bytes
	Blockhash types.Bytes
	Txhash    types.Bytes
	Amount    types.U128
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

type PoolBondState string

const (
	EraUpdated       = PoolBondState("EraUpdated")
	BondReported     = PoolBondState("BondReported")
	ActiveReported   = PoolBondState("ActiveReported")
	WithdrawSkipped  = PoolBondState("WithdrawSkipped")
	WithdrawReported = PoolBondState("WithdrawReported")
	TransferReported = PoolBondState("TransferReported")
)

func (p *PoolBondState) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		*p = EraUpdated
	case 1:
		*p = BondReported
	case 2:
		*p = ActiveReported
	case 3:
		*p = WithdrawSkipped
	case 4:
		*p = WithdrawReported
	case 5:
		*p = TransferReported
	default:
		return fmt.Errorf("PoolBondState decode error: %d", b)
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
	Symbol     core.RSymbol
	BondId     types.Hash
	MethodName string
}

// encode takes only nonce and call and encodes them for storage queries
func (p *Proposal) Encode() ([]byte, error) {
	return types.EncodeToBytes(struct {
		types.Hash
		types.Call
	}{p.BondId, p.Call})
}

type PoolSnapshot struct {
	Symbol    core.RSymbol
	Era       uint32
	Pool      []byte
	Bond      types.U128
	Unbond    types.U128
	Active    types.U128
	LastVoter types.AccountID
	BondState PoolBondState
}

type EraPoolUpdatedFlow struct {
	Symbol        core.RSymbol
	Era           uint32
	ShotId        types.Hash
	LastVoter     types.AccountID
	LastVoterFlag bool
	Snap          *PoolSnapshot
}

type BondReportedFlow struct {
	Symbol        core.RSymbol
	ShotId        types.Hash
	LastVoter     types.AccountID
	LastVoterFlag bool
	Snap          *PoolSnapshot
	LastEra       uint32
	SubAccounts   []types.Bytes
	Stashes       []types.AccountID
}

type ActiveReportedFlow struct {
	Symbol        core.RSymbol
	ShotId        types.Hash
	LastVoter     types.AccountID
	LastVoterFlag bool
	Snap          *PoolSnapshot
}

type WithdrawReportedFlow struct {
	Symbol        core.RSymbol
	ShotId        types.Hash
	LastVoter     types.AccountID
	LastVoterFlag bool
	Snap          *PoolSnapshot
	Receives      []*Receive
	TotalAmount   types.U128
}

type NominationUpdatedFlow struct {
	Symbol        core.RSymbol
	Pool          []byte
	NewValidators []types.Bytes
	Era           uint32
	LastVoter     types.AccountID
	LastVoterFlag bool
}

type GetEraNominatedFlow struct {
	Symbol        core.RSymbol
	Pool          []byte
	Era           uint32
	NewValidators chan []types.AccountID
}

type MultiEventFlow struct {
	EventId         string
	Symbol          core.RSymbol
	EventData       interface{}
	Threshold       uint16
	SubAccounts     []types.Bytes
	Key             *signature.KeyringPair
	Others          []types.AccountID
	OpaqueCalls     []*MultiOpaqueCall
	PaymentInfo     *rpc.PaymentQueryInfo
	NewMulCallHashs map[string]bool
	MulExeCallHashs map[string]bool
}

type EventNewMultisig struct {
	Who, ID     types.AccountID
	CallHash    types.Hash
	CallHashStr string
	TimePoint   *OptionTimePoint
	Approvals   []types.AccountID
}

type Multisig struct {
	When      types.TimePoint
	Deposit   types.U128
	Depositor types.AccountID
	Approvals []types.AccountID
}

type EventMultisigExecuted struct {
	Who, ID     types.AccountID
	TimePoint   types.TimePoint
	CallHash    types.Hash
	CallHashStr string
	Result      bool
}

type MultiCallParam struct {
	TimePoint *OptionTimePoint
	Opaque    []byte
	Extrinsic string
	CallHash  string
}

type PoolUnbondKey struct {
	Pool []byte
	Era  uint32
}

type Unbonding struct {
	Who       types.AccountID
	Value     types.U128
	Recipient []byte
}

type Receive struct {
	Recipient []byte
	Value     types.UCompact
}

type Era struct {
	Type  string `json:"type"`
	Value uint32 `json:"value"`
}

type ChainEvent struct {
	ModuleId string                  `json:"module_id" `
	EventId  string                  `json:"event_id" `
	Params   []scalecodec.EventParam `json:"params"`
}

type EraRewardPoints struct {
	Total       uint32
	Individuals []Individual
}

type Individual struct {
	Validator   types.AccountID
	RewardPoint uint32
}

type Exposure struct {
	Total  types.U128
	Own    types.U128
	Others []*IndividualExposure
}

type IndividualExposure struct {
	Who   types.AccountID
	Value types.U128
}

type StakingLedger struct {
	Stash          types.AccountID
	Total          types.UCompact
	Active         types.UCompact
	Unlocking      []UnlockChunk
	ClaimedRewards []uint32
}

type UnlockChunk struct {
	Value types.UCompact
	Era   types.UCompact
}

type MultiOpaqueCall struct {
	Extrinsic string
	Opaque    []byte
	CallHash  string
	TimePoint *OptionTimePoint
}

type Transaction struct {
	ExtrinsicHash  string
	CallModuleName string
	CallName       string
	Address        interface{}
	Params         []chainTypes.ExtrinsicParam
}

type OriginalTx string

const (
	OriginalTransfer       = OriginalTx("Transfer")
	OriginalBond           = OriginalTx("Bond")
	OriginalUnbond         = OriginalTx("Unbond")
	OriginalWithdrawUnbond = OriginalTx("WithdrawUnbond")
	OriginalClaimRewards   = OriginalTx("ClaimRewards")
)

func (r *OriginalTx) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		*r = OriginalTransfer
	case 1:
		*r = OriginalBond
	case 2:
		*r = OriginalUnbond
	case 3:
		*r = OriginalWithdrawUnbond
	case 4:
		*r = OriginalClaimRewards
	default:
		return fmt.Errorf("OriginalTx decode error: %d", b)
	}

	return nil
}

func (r OriginalTx) Encode(encoder scale.Encoder) error {
	switch r {
	case OriginalTransfer:
		return encoder.PushByte(0)
	case OriginalBond:
		return encoder.PushByte(1)
	case OriginalUnbond:
		return encoder.PushByte(2)
	case OriginalWithdrawUnbond:
		return encoder.PushByte(3)
	case OriginalClaimRewards:
		return encoder.PushByte(4)
	default:
		return fmt.Errorf("OriginalTx %s not supported", r)
	}
}

type SubmitSignatureParams struct {
	Symbol     core.RSymbol
	Era        types.U32
	Pool       types.Bytes
	TxType     OriginalTx
	ProposalId types.Bytes
	Signature  types.Bytes
}

type GetReceiversParams struct {
	Symbol core.RSymbol
	Era    types.U32
	Pool   types.Bytes
}

type SubmitSignatures struct {
	Symbol     core.RSymbol
	Era        types.U32
	Pool       types.Bytes
	TxType     OriginalTx
	ProposalId types.Bytes
	Signature  []types.Bytes
}

type SignaturesKey struct {
	RSymbol    core.RSymbol
	Era        uint32
	Pool       []byte
	TxType     OriginalTx
	ProposalId []byte
}

//SignaturesEnough(RSymbol, u32, Vec<u8>, OriginalTxType, Vec<u8>),
type EvtSignatureEnough struct {
	RSymbol    core.RSymbol
	Era        uint32
	Pool       []byte
	TxType     OriginalTx
	ProposalId []byte
}
