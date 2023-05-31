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
	addr, err := BnbValAddressFromHex("0xaabfcf13c725951c373c937842a0e80183e37e2e")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(addr.String())
	assert.Equal(t, "bva142lu7y78yk23cdeujduy9g8gqxp7xl3wrz4ska", addr.String())

	hexStr, err := BnbValAddressFromBech32("bva142lu7y78yk23cdeujduy9g8gqxp7xl3wrz4ska")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hexutil.Encode(hexStr))
	assert.Equal(t, "0xaabfcf13c725951c373c937842a0e80183e37e2e", hexutil.Encode(hexStr))

	bts, err := hexutil.Decode("0x6276613134326c7537793738796b3233636465756a6475793967386771787037786c3377727a34736b61")
	if err != nil {
		t.Fatal(err)
	}

	addr2, err := BnbValAddressFromBech32(string(bts))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(addr2.String())
	assert.Equal(t, "bva142lu7y78yk23cdeujduy9g8gqxp7xl3wrz4ska", addr2.String())

	hexStr2, err := BnbValAddressFromBech32("bva18w9m90ksmcnsw6rd2kwpd74m4l5agkc68su3et")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hexutil.Encode(hexStr2))

	bts, _ = hexutil.Decode("0x992701c853301b120a6db0e4767aeecc737c515a")
	t.Log(GetBcRewardAddressFromBsc(bts).String())

	bncCmnTypes.Network = bncCmnTypes.TestNetwork
	bts, _ = hexutil.Decode("0x44f95eef755ed4fbdc19e3e8f617773d23e44a5b")
	t.Log(GetBcRewardAddressFromBsc(bts).String())

	hexStr, err = BnbValAddressFromBech32("bva15vvmregln3skagjcjxq9hshj452a2ppjkhperr")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hexutil.Encode(hexStr))
	
	hexStr, err = BnbValAddressFromBech32("bva15mgzha93ny878kuvjl0pnqmjygwccdad08uecu")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hexutil.Encode(hexStr))
}
