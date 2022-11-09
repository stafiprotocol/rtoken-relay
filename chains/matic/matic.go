package matic

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stafiprotocol/rtoken-relay/bindings/MaticToken"
	"github.com/stafiprotocol/rtoken-relay/bindings/MultiSend"
	"github.com/stafiprotocol/rtoken-relay/bindings/ValidatorShare"
)

// / ValidatorShare
var (
	ValidatorShareAbi, _ = abi.JSON(strings.NewReader(ValidatorShare.ValidatorShareABI))
	MultiSendAbi, _      = abi.JSON(strings.NewReader(MultiSend.MultiSendABI))

	BuyVoucherSafeTxGas     = big.NewInt(450000)
	SellVoucherNewSafeTxGas = big.NewInt(450000)
	RestakeSafeTxGas        = big.NewInt(200000)
	WithdrawTxGas           = big.NewInt(300000)
	TransferTxGas           = big.NewInt(100000)
)

const (
	BuyVoucherMethodName = "buyVoucher" // 0x6ab15071
	//BuyVoucherMethodId = "0x6ab15071"

	SellVoucherNewMethodName = "sellVoucher_new" //0xc83ec04d
	//SellVoucherNewMethodId = "0xc83ec04d"

	RestakeMethodName     = "restake"                // claim/payout
	UnstakeClaimTokensNew = "unstakeClaimTokens_new" // withdraw

	MultiSendMethodName = "multiSend"
)

// / MaticToken
var (
	MaticTokenAbi, _ = abi.JSON(strings.NewReader(MaticToken.MaticTokenABI))
)

const (
	TransferMethodName = "transfer" // 0xa9059cbb
	//TransferMethodId = "0xa9059cbb"
)
