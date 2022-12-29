package submodel

import (
	"fmt"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

type BondStateKey struct {
	BlockHash types.Bytes
	TxHash    types.Bytes
}

type BondFlow struct {
	Symbol      core.RSymbol
	BondId      types.Hash
	Record      *BondRecord
	Reason      BondReason
	State       BondState
	VerifyTimes int
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

type BondState string

const (
	Dealing = BondState("Dealing")
	Fail    = BondState("Fail")
	Success = BondState("Success")
	Default = BondState("Default")
)

func (bs BondState) Encode(encoder scale.Encoder) error {
	switch bs {
	case Dealing:
		return encoder.PushByte(0)
	case Fail:
		return encoder.PushByte(1)
	case Success:
		return encoder.PushByte(2)
	case Default:
		return encoder.PushByte(100)
	default:
		return fmt.Errorf("BondState %s not supported", bs)
	}
}

func (bs *BondState) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		*bs = Dealing
	case 1:
		*bs = Fail
	case 2:
		*bs = Success
	case 100:
		*bs = Default
	default:
		return fmt.Errorf("BondState decode error: %d", b)
	}

	return nil
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

type BondAction string

const (
	BondOnly         = BondAction("BondOnly")
	UnbondOnly       = BondAction("UnbondOnly")
	BothBondUnbond   = BondAction("BothBondUnbond")
	EitherBondUnbond = BondAction("EitherBondUnbond")
	InterDeduct      = BondAction("InterDeduct")
)

func (ba BondAction) Encode(encoder scale.Encoder) error {
	switch ba {
	case BondOnly:
		return encoder.PushByte(0)
	case UnbondOnly:
		return encoder.PushByte(1)
	case BothBondUnbond:
		return encoder.PushByte(2)
	case EitherBondUnbond:
		return encoder.PushByte(3)
	case InterDeduct:
		return encoder.PushByte(4)
	default:
		return fmt.Errorf("BondAction %s not supported", ba)
	}
}

func (ba *BondAction) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		*ba = BondOnly
	case 1:
		*ba = UnbondOnly
	case 2:
		*ba = BothBondUnbond
	case 3:
		*ba = EitherBondUnbond
	case 4:
		*ba = InterDeduct
	default:
		return fmt.Errorf("BondAction decode error: %d", b)
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

type BondReportType uint8

const (
	NewBondReport                       = BondReportType(0)
	BondAndReportActive                 = BondReportType(1)
	BondAndReportActiveWithPendingValue = BondReportType(2)
)

type BondCall struct {
	ReportType BondReportType
	Action     BondAction
}

type EraPoolUpdatedFlow struct {
	Symbol        core.RSymbol
	Era           uint32
	ShotId        types.Hash
	LastVoter     types.AccountID
	LastVoterFlag bool
	Snap          *PoolSnapshot
	LeastBond     *big.Int
	BondCall      *BondCall
	Active        *big.Int
	Reward        *big.Int
	PendingStake  *big.Int // rBNB use
	PendingReward *big.Int // rBNB use
}

type BondReportedFlow struct {
	Symbol              core.RSymbol
	ShotId              types.Hash
	LastVoter           types.AccountID
	LastVoterFlag       bool
	Snap                *PoolSnapshot
	LastEra             uint32
	EraBlock            uint64
	Unstaked            types.U128
	SubAccounts         []types.Bytes
	Threshold           uint32
	Stashes             []types.AccountID
	MaticValidatorId    *big.Int
	MultiTransaction    *ethmodel.MultiTransaction
	NewActiveReportFlag bool
	LeastBond           *big.Int
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

type TransferReportedFlow struct {
	Symbol        core.RSymbol
	ShotId        types.Hash
	Snap          *PoolSnapshot
	LastVoterFlag bool
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

type ValidatorUpdatedFlow struct {
	Symbol       core.RSymbol
	Pool         []byte
	OldValidator types.Bytes
	NewValidator types.Bytes
	Era          uint32
}

type GetEraNominatedFlow struct {
	Symbol        core.RSymbol
	Pool          []byte
	Era           uint32
	NewValidators chan []types.AccountID
}

type GetBondStateFlow struct {
	Symbol    core.RSymbol
	BlockHash types.Bytes
	TxHash    types.Bytes
	BondState chan BondState
}

type MultiEventFlow struct {
	EventId          string
	Symbol           core.RSymbol
	EventData        interface{}
	Threshold        uint16
	SubAccounts      []types.Bytes
	Key              *signature.KeyringPair
	Others           []types.AccountID
	OpaqueCalls      []*MultiOpaqueCall
	PaymentInfo      *rpc.PaymentQueryInfo
	NewMulCallHashs  map[string]bool
	MulExeCallHashs  map[string]bool
	MaticValidatorId *big.Int         // rmaitic use
	BnbValidators    []common.Address // rbnb use
	MultiTransaction *ethmodel.MultiTransaction
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
	Total  types.UCompact
	Own    types.UCompact
	Others []IndividualExposure
}

type IndividualExposure struct {
	Who   types.AccountID
	Value types.UCompact
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
	OriginalTxDefault      = OriginalTx("default")
	OriginalTransfer       = OriginalTx("Transfer") //transfer
	OriginalBond           = OriginalTx("Bond")     //bond or unbond
	OriginalUnbond         = OriginalTx("Unbond")
	OriginalWithdrawUnbond = OriginalTx("WithdrawUnbond") //redelegate: validator update
	OriginalClaimRewards   = OriginalTx("ClaimRewards")   // claim
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

func (ssp *SubmitSignatureParams) EncodeToHash() (common.Hash, error) {
	symBz, err := types.EncodeToBytes(ssp.Symbol)
	if err != nil {
		return [32]byte{}, err
	}

	eraBz, err := types.EncodeToBytes(ssp.Era)
	if err != nil {
		return [32]byte{}, err
	}

	txTypeBz, err := types.EncodeToBytes(ssp.TxType)
	if err != nil {
		return [32]byte{}, err
	}

	packed := make([]byte, 0)
	packed = append(packed, symBz...)
	packed = append(packed, eraBz...)
	packed = append(packed, ssp.Pool...)
	packed = append(packed, txTypeBz...)

	return crypto.Keccak256Hash(packed), nil
}

type GetReceiversParams struct {
	Symbol core.RSymbol
	Era    types.U32
	Pool   types.Bytes
}

type GetSubmitSignaturesFlow struct {
	Symbol     core.RSymbol
	Era        types.U32
	Pool       types.Bytes
	TxType     OriginalTx
	ProposalId types.Bytes
	Signatures chan []types.Bytes
}
type GetPoolThresholdFlow struct {
	Symbol    core.RSymbol
	Pool      types.Bytes
	Threshold chan uint32
}

type SubmitSignatures struct {
	Symbol     core.RSymbol
	Era        types.U32
	Pool       types.Bytes
	TxType     OriginalTx
	ProposalId types.Bytes
	Signature  []types.Bytes
	Threshold  uint32
}

func (ss *SubmitSignatures) EncodeToHash() (common.Hash, error) {
	symBz, err := types.EncodeToBytes(ss.Symbol)
	if err != nil {
		return [32]byte{}, err
	}

	eraBz, err := types.EncodeToBytes(ss.Era)
	if err != nil {
		return [32]byte{}, err
	}

	txTypeBz, err := types.EncodeToBytes(ss.TxType)
	if err != nil {
		return [32]byte{}, err
	}

	packed := make([]byte, 0)
	packed = append(packed, symBz...)
	packed = append(packed, eraBz...)
	packed = append(packed, ss.Pool...)
	packed = append(packed, txTypeBz...)

	return crypto.Keccak256Hash(packed), nil
}

type SignaturesKey struct {
	Era        uint32
	Pool       []byte
	TxType     OriginalTx
	ProposalId []byte
}

// SignaturesEnough(RSymbol, u32, Vec<u8>, OriginalTxType, Vec<u8>),
type EvtSignatureEnough struct {
	RSymbol    core.RSymbol
	Era        uint32
	Pool       []byte
	TxType     OriginalTx
	ProposalId []byte
}

// execute_bond_and_swap(origin, pool: Vec<u8>, blockhash: Vec<u8>, txhash: Vec<u8>, amount: u128, symbol: RSymbol, stafi_recipient: T::AccountId, dest_recipient: Vec<u8>, dest_id: ChainId, reason: BondReason)
type ExeLiquidityBondAndSwapFlow struct {
	Pool           types.Bytes
	Blockhash      types.Bytes
	Txhash         types.Bytes
	Amount         types.U128
	Symbol         core.RSymbol
	StafiRecipient types.AccountID
	DestRecipient  types.Bytes
	DestId         types.U8
	Reason         BondReason
}
