package utils

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	bncCmnTypes "github.com/stafiprotocol/go-sdk/common/types"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stretchr/testify/assert"
)

var (
	csvFile = "../datas/test.csv"
)

func TestWriteAndReadeCSV(t *testing.T) {
	contents := [][]string{
		{"id1", "name1", "score60"},
		{"id2", "name2", "score62"},
	}

	err := WriteCSV(csvFile, contents)
	assert.NoError(t, err)

	lines := ReadCSV(csvFile)
	t.Log(lines)
}

func TestWriteAndReadeCSV1(t *testing.T) {
	contents := [][]string{
		{"id1", "name1", "score63"},
		{"id2", "name2", "score64"},
	}

	err := WriteCSV(csvFile, contents)
	assert.NoError(t, err)

	lines := ReadCSV(csvFile)
	t.Log(lines)
}

func TestGetStakeCAoB(t *testing.T) {
	exp, err := bncCmnTypes.AccAddressFromHex("91D7deA99716Cbb247E81F1cfB692009164a967E")
	if err != nil {
		t.Fatal(err)
	}
	stakeCAoB := GetStakeCAoB(exp.Bytes(), DelegateCAoBSalt)
	fmt.Println(stakeCAoB.String())
	if delAddr := GetStakeCAoB(stakeCAoB.Bytes(), DelegateCAoBSalt); delAddr.String() != exp.String() {
		t.Fatal()
	}
}

func TestStrRecieves(t *testing.T) {
	receives := []*submodel.Receive{{
		Recipient: []byte{9, 3, 4, 4, 3},
		Value:     types.NewUCompact(big.NewInt(1)),
	}, {
		Recipient: []byte{13, 3, 4, 5},
		Value:     types.NewUCompact(big.NewInt(1)),
	}}

	t.Log(StrReceives(receives))
}

func TestDecimal(t *testing.T) {
	t.Log(decimal.NewFromInt(5).Div(decimal.NewFromInt(100)).String())
	t.Log(decimal.NewFromInt(5).Div(decimal.NewFromInt(100)).Ceil().BigInt().Int64())
	t.Log(decimal.NewFromInt(0).Div(decimal.NewFromInt(100)).Ceil().BigInt().Int64())
}

func TestBnbVal(t *testing.T) {
	addr, err := BnbValAddressFromHex("0x13D0fC20edb23021e8675ead7676301ce9B748Ee")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(addr.String())

	hexStr, err := BnbValAddressFromBech32("bva142lu7y78yk23cdeujduy9g8gqxp7xl3wrz4ska")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hexutil.Encode(hexStr))

	bts, err := hexutil.Decode("0x627661317a306730636738646b67637a7236723874366b6876613373726e356d776a387735746c753768")
	if err != nil {
		t.Fatal(err)
	}

	addr2, err := BnbValAddressFromBech32(string(bts))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(addr2.String())

}
