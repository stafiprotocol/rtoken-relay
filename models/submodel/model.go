package submodel

import (
	"fmt"
	"github.com/itering/substrate-api-rpc/rpc"

	"github.com/stafiprotocol/go-substrate-rpc-client/scale"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
)

type EvtLiquidityBond struct {
	AccountId types.AccountID
	Rsymbol   core.RSymbol
	BondId    types.Hash
}

type BondFlow struct {
	Key    *BondKey
	Record *BondRecord
	Reason BondReason
}

type BondKey struct {
	Rsymbol core.RSymbol
	BondId  types.Hash
}

type BondRecord struct {
	Bonder    types.AccountID
	Rsymbol   core.RSymbol
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
	Rsymbol core.RSymbol
	Pool    []byte
}

type MultiEventFlow struct {
	EventId         string
	EventData       interface{}
	Threshold       uint16
	SubAccounts     []types.Bytes
	Key             *signature.KeyringPair
	Others          []types.AccountID
	OpaqueCalls     []*substrate.MultiOpaqueCall
	PaymentInfo     *rpc.PaymentQueryInfo
	NewMulCallHashs map[string]bool
	MulExeCallHashs map[string]bool
}

type BondReportFlow struct {
	ShotId        types.Hash
	Rsymbol       core.RSymbol
	Pool          []byte
	Era           uint32
	LastVoter     types.AccountID
	LastEra       uint32
	LastVoterFlag bool
	Active        types.U128
	Stashes       []types.AccountID
}

type WithdrawUnbondFlow struct {
	ShotId        types.Hash
	Rsymbol       core.RSymbol
	Pool          []byte
	Era           uint32
	LastVoter     types.AccountID
	LastVoterFlag bool
}

type TransferFlow struct {
	ShotId        types.Hash
	Rsymbol       core.RSymbol
	Pool          []byte
	Era           uint32
	LastVoter     types.AccountID
	LastVoterFlag bool
	Receives      []*Receive
	TotalAmount   types.U128
}

type EraPoolUpdatedFlow struct {
	ShotId        types.Hash
	LastVoter     types.AccountID
	LastVoterFlag bool
	Snap          *EraPoolSnapshot
}

type EraPoolSnapshot struct {
	Rsymbol   core.RSymbol
	Era       uint32
	Pool      []byte
	Bond      types.U128
	Unbond    types.U128
	Active    types.U128
	LastVoter types.AccountID
}

type EventNewMultisig struct {
	Who, ID     types.AccountID
	CallHash    types.Hash
	CallHashStr string
	TimePoint   *substrate.OptionTimePoint
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

//type MultisigFlow struct {
//	HeadFlow    interface{}
//	MulCall     *MultisigCall
//	CallHash    string
//	NewMul      *EventNewMultisig
//	Multisig    *Multisig
//	MulExecuted *EventMultisigExecuted
//	Params      []*MultiBatchCallParam
//}

type MultisigCall struct {
	//Params      []*MultiCallParam
	//TimePoint   *OptionTimePoint
	//Opaque      []byte
	//Extrinsic   string
	//CallHash    string
}

type MultiCallParam struct {
	TimePoint *substrate.OptionTimePoint
	Opaque    []byte
	Extrinsic string
	CallHash  string
}

type Unbonding struct {
	Who        types.AccountID
	Symbol     core.RSymbol
	Pool       []byte
	Rvalue     types.U128
	Value      types.U128
	CurrentEra uint32
	UnlockEra  uint32
	Recipient  []byte
}

type Receive struct {
	Recipient types.Address
	Value     types.UCompact
}

type Era struct {
	Type  string `json:"type"`
	Value uint32 `json:"value"`
}
