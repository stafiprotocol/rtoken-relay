package matic

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/rtoken-relay/config"
	"gotest.tools/assert"
	"math/big"
	"testing"
)

func TestBuyVoucher(t *testing.T) {
	calldata := "0x6ab150710000000000000000000000000000000000000000000000056bc75e2d631000000000000000000000000000000000000000000000000000000000000000000000"
	sig, _ := hexutil.Decode(calldata)
	m, err := ValidatorShareAbi.MethodById(sig)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, BuyVoucherMethodName, m.Name)

	// Test pack/unpack
	packed, err := ValidatorShareAbi.Pack(m.Name, big.NewInt(0).Mul(&config.AmountBase, big.NewInt(100)), big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, calldata, hexutil.Encode(packed))
}

func TestSellVoucherNew(t *testing.T) {
	sig, _ := hexutil.Decode("0xc83ec04d")
	m, err := ValidatorShareAbi.MethodById(sig)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, SellVoucherNewMethodName, m.Name)
}

func TestMaticTransfer(t *testing.T) {
	sig, _ := hexutil.Decode("0xa9059cbb")
	m, err := MaticTokenAbi.MethodById(sig)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, TransferMethodName, m.Name)
}
