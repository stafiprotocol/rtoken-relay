package ethmodel

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type MultiTransaction struct {
	To        common.Address
	Value     *big.Int
	CallData  []byte
	CallType  uint8
	SafeTxGas *big.Int
}

type BatchTransaction struct {
	Operation  uint8
	To         common.Address
	Value      *big.Int
	DataLength *big.Int
	Data       []byte
}

type BatchTransactions []*BatchTransaction

func (bts BatchTransactions) Encode() []byte {
	packed := make([]byte, 0)
	for _, bt := range bts {
		packed = append(packed, bt.Encode()...)
	}
	return packed
}

func (bt *BatchTransaction) Encode() []byte {
	packed := []byte{bt.Operation}
	packed = append(packed, bt.To.Bytes()...)
	packed = append(packed, common.LeftPadBytes(bt.Value.Bytes(), 32)...)
	packed = append(packed, common.LeftPadBytes(bt.DataLength.Bytes(), 32)...)
	packed = append(packed, bt.Data...)
	return packed
}

type EventSig string

const (
	TransferEvent         = EventSig("Transfer(address,address,uint256)")
	TransferEventTopicLen = 3
)

func (es EventSig) GetTopic() common.Hash {
	return crypto.Keccak256Hash([]byte(es))
}

type TxHashState uint8

var (
	HashStateUnsubmit = TxHashState(0)
	HashStateFail     = TxHashState(1)
	HashStateSuccess  = TxHashState(2)
)
