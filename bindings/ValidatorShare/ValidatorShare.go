// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ValidatorShare

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ValidatorShareABI is the input ABI used to generate the binding from.
const ValidatorShareABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"activeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minSharesToMint\",\"type\":\"uint256\"}],\"name\":\"buyVoucher\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountToDeposit\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"commissionRate_deprecated\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"destination\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"drain\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"eventsHub\",\"outputs\":[{\"internalType\":\"contractEventsHub\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"exchangeRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getLiquidRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRewardPerShare\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getTotalStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"initalRewardPerShare\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_validatorId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_stakingLogger\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_stakeManager\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lastCommissionUpdate_deprecated\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"locked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"migrateIn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"migrateOut\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"restake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"rewardPerShare\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"claimAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maximumSharesToBurn\",\"type\":\"uint256\"}],\"name\":\"sellVoucher\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"claimAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maximumSharesToBurn\",\"type\":\"uint256\"}],\"name\":\"sellVoucher_new\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegatedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalAmountToSlash\",\"type\":\"uint256\"}],\"name\":\"slash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stakeManager\",\"outputs\":[{\"internalType\":\"contractIStakeManager\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stakingLogger\",\"outputs\":[{\"internalType\":\"contractStakingInfo\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalStake_deprecated\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"unbondNonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"unbonds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawEpoch\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"unbonds_new\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawEpoch\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unstakeClaimTokens\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"unbondNonce\",\"type\":\"uint256\"}],\"name\":\"unstakeClaimTokens_new\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_delegation\",\"type\":\"bool\"}],\"name\":\"updateDelegation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorRewards_deprecated\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawExchangeRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ValidatorShare is an auto generated Go binding around an Ethereum contract.
type ValidatorShare struct {
	ValidatorShareCaller     // Read-only binding to the contract
	ValidatorShareTransactor // Write-only binding to the contract
	ValidatorShareFilterer   // Log filterer for contract events
}

// ValidatorShareCaller is an auto generated read-only Go binding around an Ethereum contract.
type ValidatorShareCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorShareTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ValidatorShareTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorShareFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ValidatorShareFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorShareSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ValidatorShareSession struct {
	Contract     *ValidatorShare   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ValidatorShareCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ValidatorShareCallerSession struct {
	Contract *ValidatorShareCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ValidatorShareTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ValidatorShareTransactorSession struct {
	Contract     *ValidatorShareTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ValidatorShareRaw is an auto generated low-level Go binding around an Ethereum contract.
type ValidatorShareRaw struct {
	Contract *ValidatorShare // Generic contract binding to access the raw methods on
}

// ValidatorShareCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ValidatorShareCallerRaw struct {
	Contract *ValidatorShareCaller // Generic read-only contract binding to access the raw methods on
}

// ValidatorShareTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ValidatorShareTransactorRaw struct {
	Contract *ValidatorShareTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValidatorShare creates a new instance of ValidatorShare, bound to a specific deployed contract.
func NewValidatorShare(address common.Address, backend bind.ContractBackend) (*ValidatorShare, error) {
	contract, err := bindValidatorShare(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ValidatorShare{ValidatorShareCaller: ValidatorShareCaller{contract: contract}, ValidatorShareTransactor: ValidatorShareTransactor{contract: contract}, ValidatorShareFilterer: ValidatorShareFilterer{contract: contract}}, nil
}

// NewValidatorShareCaller creates a new read-only instance of ValidatorShare, bound to a specific deployed contract.
func NewValidatorShareCaller(address common.Address, caller bind.ContractCaller) (*ValidatorShareCaller, error) {
	contract, err := bindValidatorShare(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorShareCaller{contract: contract}, nil
}

// NewValidatorShareTransactor creates a new write-only instance of ValidatorShare, bound to a specific deployed contract.
func NewValidatorShareTransactor(address common.Address, transactor bind.ContractTransactor) (*ValidatorShareTransactor, error) {
	contract, err := bindValidatorShare(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorShareTransactor{contract: contract}, nil
}

// NewValidatorShareFilterer creates a new log filterer instance of ValidatorShare, bound to a specific deployed contract.
func NewValidatorShareFilterer(address common.Address, filterer bind.ContractFilterer) (*ValidatorShareFilterer, error) {
	contract, err := bindValidatorShare(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValidatorShareFilterer{contract: contract}, nil
}

// bindValidatorShare binds a generic wrapper to an already deployed contract.
func bindValidatorShare(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ValidatorShareABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ValidatorShare *ValidatorShareRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ValidatorShare.Contract.ValidatorShareCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ValidatorShare *ValidatorShareRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorShare.Contract.ValidatorShareTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ValidatorShare *ValidatorShareRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ValidatorShare.Contract.ValidatorShareTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ValidatorShare *ValidatorShareCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ValidatorShare.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ValidatorShare *ValidatorShareTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorShare.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ValidatorShare *ValidatorShareTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ValidatorShare.Contract.contract.Transact(opts, method, params...)
}

// ActiveAmount is a free data retrieval call binding the contract method 0x3a09bf44.
//
// Solidity: function activeAmount() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) ActiveAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "activeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveAmount is a free data retrieval call binding the contract method 0x3a09bf44.
//
// Solidity: function activeAmount() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) ActiveAmount() (*big.Int, error) {
	return _ValidatorShare.Contract.ActiveAmount(&_ValidatorShare.CallOpts)
}

// ActiveAmount is a free data retrieval call binding the contract method 0x3a09bf44.
//
// Solidity: function activeAmount() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) ActiveAmount() (*big.Int, error) {
	return _ValidatorShare.Contract.ActiveAmount(&_ValidatorShare.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ValidatorShare.Contract.Allowance(&_ValidatorShare.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ValidatorShare.Contract.Allowance(&_ValidatorShare.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ValidatorShare.Contract.BalanceOf(&_ValidatorShare.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ValidatorShare.Contract.BalanceOf(&_ValidatorShare.CallOpts, owner)
}

// CommissionRateDeprecated is a free data retrieval call binding the contract method 0x8ccdd289.
//
// Solidity: function commissionRate_deprecated() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) CommissionRateDeprecated(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "commissionRate_deprecated")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CommissionRateDeprecated is a free data retrieval call binding the contract method 0x8ccdd289.
//
// Solidity: function commissionRate_deprecated() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) CommissionRateDeprecated() (*big.Int, error) {
	return _ValidatorShare.Contract.CommissionRateDeprecated(&_ValidatorShare.CallOpts)
}

// CommissionRateDeprecated is a free data retrieval call binding the contract method 0x8ccdd289.
//
// Solidity: function commissionRate_deprecated() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) CommissionRateDeprecated() (*big.Int, error) {
	return _ValidatorShare.Contract.CommissionRateDeprecated(&_ValidatorShare.CallOpts)
}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(bool)
func (_ValidatorShare *ValidatorShareCaller) Delegation(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "delegation")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(bool)
func (_ValidatorShare *ValidatorShareSession) Delegation() (bool, error) {
	return _ValidatorShare.Contract.Delegation(&_ValidatorShare.CallOpts)
}

// Delegation is a free data retrieval call binding the contract method 0xdf5cf723.
//
// Solidity: function delegation() view returns(bool)
func (_ValidatorShare *ValidatorShareCallerSession) Delegation() (bool, error) {
	return _ValidatorShare.Contract.Delegation(&_ValidatorShare.CallOpts)
}

// EventsHub is a free data retrieval call binding the contract method 0x883b455f.
//
// Solidity: function eventsHub() view returns(address)
func (_ValidatorShare *ValidatorShareCaller) EventsHub(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "eventsHub")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EventsHub is a free data retrieval call binding the contract method 0x883b455f.
//
// Solidity: function eventsHub() view returns(address)
func (_ValidatorShare *ValidatorShareSession) EventsHub() (common.Address, error) {
	return _ValidatorShare.Contract.EventsHub(&_ValidatorShare.CallOpts)
}

// EventsHub is a free data retrieval call binding the contract method 0x883b455f.
//
// Solidity: function eventsHub() view returns(address)
func (_ValidatorShare *ValidatorShareCallerSession) EventsHub() (common.Address, error) {
	return _ValidatorShare.Contract.EventsHub(&_ValidatorShare.CallOpts)
}

// ExchangeRate is a free data retrieval call binding the contract method 0x3ba0b9a9.
//
// Solidity: function exchangeRate() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) ExchangeRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "exchangeRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExchangeRate is a free data retrieval call binding the contract method 0x3ba0b9a9.
//
// Solidity: function exchangeRate() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) ExchangeRate() (*big.Int, error) {
	return _ValidatorShare.Contract.ExchangeRate(&_ValidatorShare.CallOpts)
}

// ExchangeRate is a free data retrieval call binding the contract method 0x3ba0b9a9.
//
// Solidity: function exchangeRate() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) ExchangeRate() (*big.Int, error) {
	return _ValidatorShare.Contract.ExchangeRate(&_ValidatorShare.CallOpts)
}

// GetLiquidRewards is a free data retrieval call binding the contract method 0x676e5550.
//
// Solidity: function getLiquidRewards(address user) view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) GetLiquidRewards(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "getLiquidRewards", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLiquidRewards is a free data retrieval call binding the contract method 0x676e5550.
//
// Solidity: function getLiquidRewards(address user) view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) GetLiquidRewards(user common.Address) (*big.Int, error) {
	return _ValidatorShare.Contract.GetLiquidRewards(&_ValidatorShare.CallOpts, user)
}

// GetLiquidRewards is a free data retrieval call binding the contract method 0x676e5550.
//
// Solidity: function getLiquidRewards(address user) view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) GetLiquidRewards(user common.Address) (*big.Int, error) {
	return _ValidatorShare.Contract.GetLiquidRewards(&_ValidatorShare.CallOpts, user)
}

// GetRewardPerShare is a free data retrieval call binding the contract method 0x1bf494a7.
//
// Solidity: function getRewardPerShare() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) GetRewardPerShare(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "getRewardPerShare")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRewardPerShare is a free data retrieval call binding the contract method 0x1bf494a7.
//
// Solidity: function getRewardPerShare() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) GetRewardPerShare() (*big.Int, error) {
	return _ValidatorShare.Contract.GetRewardPerShare(&_ValidatorShare.CallOpts)
}

// GetRewardPerShare is a free data retrieval call binding the contract method 0x1bf494a7.
//
// Solidity: function getRewardPerShare() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) GetRewardPerShare() (*big.Int, error) {
	return _ValidatorShare.Contract.GetRewardPerShare(&_ValidatorShare.CallOpts)
}

// GetTotalStake is a free data retrieval call binding the contract method 0x1e7ff8f6.
//
// Solidity: function getTotalStake(address user) view returns(uint256, uint256)
func (_ValidatorShare *ValidatorShareCaller) GetTotalStake(opts *bind.CallOpts, user common.Address) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "getTotalStake", user)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetTotalStake is a free data retrieval call binding the contract method 0x1e7ff8f6.
//
// Solidity: function getTotalStake(address user) view returns(uint256, uint256)
func (_ValidatorShare *ValidatorShareSession) GetTotalStake(user common.Address) (*big.Int, *big.Int, error) {
	return _ValidatorShare.Contract.GetTotalStake(&_ValidatorShare.CallOpts, user)
}

// GetTotalStake is a free data retrieval call binding the contract method 0x1e7ff8f6.
//
// Solidity: function getTotalStake(address user) view returns(uint256, uint256)
func (_ValidatorShare *ValidatorShareCallerSession) GetTotalStake(user common.Address) (*big.Int, *big.Int, error) {
	return _ValidatorShare.Contract.GetTotalStake(&_ValidatorShare.CallOpts, user)
}

// InitalRewardPerShare is a free data retrieval call binding the contract method 0xe4cd1aec.
//
// Solidity: function initalRewardPerShare(address ) view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) InitalRewardPerShare(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "initalRewardPerShare", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitalRewardPerShare is a free data retrieval call binding the contract method 0xe4cd1aec.
//
// Solidity: function initalRewardPerShare(address ) view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) InitalRewardPerShare(arg0 common.Address) (*big.Int, error) {
	return _ValidatorShare.Contract.InitalRewardPerShare(&_ValidatorShare.CallOpts, arg0)
}

// InitalRewardPerShare is a free data retrieval call binding the contract method 0xe4cd1aec.
//
// Solidity: function initalRewardPerShare(address ) view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) InitalRewardPerShare(arg0 common.Address) (*big.Int, error) {
	return _ValidatorShare.Contract.InitalRewardPerShare(&_ValidatorShare.CallOpts, arg0)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_ValidatorShare *ValidatorShareCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_ValidatorShare *ValidatorShareSession) IsOwner() (bool, error) {
	return _ValidatorShare.Contract.IsOwner(&_ValidatorShare.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_ValidatorShare *ValidatorShareCallerSession) IsOwner() (bool, error) {
	return _ValidatorShare.Contract.IsOwner(&_ValidatorShare.CallOpts)
}

// LastCommissionUpdateDeprecated is a free data retrieval call binding the contract method 0x23440679.
//
// Solidity: function lastCommissionUpdate_deprecated() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) LastCommissionUpdateDeprecated(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "lastCommissionUpdate_deprecated")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastCommissionUpdateDeprecated is a free data retrieval call binding the contract method 0x23440679.
//
// Solidity: function lastCommissionUpdate_deprecated() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) LastCommissionUpdateDeprecated() (*big.Int, error) {
	return _ValidatorShare.Contract.LastCommissionUpdateDeprecated(&_ValidatorShare.CallOpts)
}

// LastCommissionUpdateDeprecated is a free data retrieval call binding the contract method 0x23440679.
//
// Solidity: function lastCommissionUpdate_deprecated() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) LastCommissionUpdateDeprecated() (*big.Int, error) {
	return _ValidatorShare.Contract.LastCommissionUpdateDeprecated(&_ValidatorShare.CallOpts)
}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() view returns(bool)
func (_ValidatorShare *ValidatorShareCaller) Locked(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "locked")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() view returns(bool)
func (_ValidatorShare *ValidatorShareSession) Locked() (bool, error) {
	return _ValidatorShare.Contract.Locked(&_ValidatorShare.CallOpts)
}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() view returns(bool)
func (_ValidatorShare *ValidatorShareCallerSession) Locked() (bool, error) {
	return _ValidatorShare.Contract.Locked(&_ValidatorShare.CallOpts)
}

// MinAmount is a free data retrieval call binding the contract method 0x9b2cb5d8.
//
// Solidity: function minAmount() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) MinAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "minAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinAmount is a free data retrieval call binding the contract method 0x9b2cb5d8.
//
// Solidity: function minAmount() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) MinAmount() (*big.Int, error) {
	return _ValidatorShare.Contract.MinAmount(&_ValidatorShare.CallOpts)
}

// MinAmount is a free data retrieval call binding the contract method 0x9b2cb5d8.
//
// Solidity: function minAmount() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) MinAmount() (*big.Int, error) {
	return _ValidatorShare.Contract.MinAmount(&_ValidatorShare.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ValidatorShare *ValidatorShareCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ValidatorShare *ValidatorShareSession) Owner() (common.Address, error) {
	return _ValidatorShare.Contract.Owner(&_ValidatorShare.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ValidatorShare *ValidatorShareCallerSession) Owner() (common.Address, error) {
	return _ValidatorShare.Contract.Owner(&_ValidatorShare.CallOpts)
}

// RewardPerShare is a free data retrieval call binding the contract method 0x446a2ec8.
//
// Solidity: function rewardPerShare() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) RewardPerShare(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "rewardPerShare")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardPerShare is a free data retrieval call binding the contract method 0x446a2ec8.
//
// Solidity: function rewardPerShare() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) RewardPerShare() (*big.Int, error) {
	return _ValidatorShare.Contract.RewardPerShare(&_ValidatorShare.CallOpts)
}

// RewardPerShare is a free data retrieval call binding the contract method 0x446a2ec8.
//
// Solidity: function rewardPerShare() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) RewardPerShare() (*big.Int, error) {
	return _ValidatorShare.Contract.RewardPerShare(&_ValidatorShare.CallOpts)
}

// StakeManager is a free data retrieval call binding the contract method 0x7542ff95.
//
// Solidity: function stakeManager() view returns(address)
func (_ValidatorShare *ValidatorShareCaller) StakeManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "stakeManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeManager is a free data retrieval call binding the contract method 0x7542ff95.
//
// Solidity: function stakeManager() view returns(address)
func (_ValidatorShare *ValidatorShareSession) StakeManager() (common.Address, error) {
	return _ValidatorShare.Contract.StakeManager(&_ValidatorShare.CallOpts)
}

// StakeManager is a free data retrieval call binding the contract method 0x7542ff95.
//
// Solidity: function stakeManager() view returns(address)
func (_ValidatorShare *ValidatorShareCallerSession) StakeManager() (common.Address, error) {
	return _ValidatorShare.Contract.StakeManager(&_ValidatorShare.CallOpts)
}

// StakingLogger is a free data retrieval call binding the contract method 0x3d94eb05.
//
// Solidity: function stakingLogger() view returns(address)
func (_ValidatorShare *ValidatorShareCaller) StakingLogger(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "stakingLogger")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakingLogger is a free data retrieval call binding the contract method 0x3d94eb05.
//
// Solidity: function stakingLogger() view returns(address)
func (_ValidatorShare *ValidatorShareSession) StakingLogger() (common.Address, error) {
	return _ValidatorShare.Contract.StakingLogger(&_ValidatorShare.CallOpts)
}

// StakingLogger is a free data retrieval call binding the contract method 0x3d94eb05.
//
// Solidity: function stakingLogger() view returns(address)
func (_ValidatorShare *ValidatorShareCallerSession) StakingLogger() (common.Address, error) {
	return _ValidatorShare.Contract.StakingLogger(&_ValidatorShare.CallOpts)
}

// TotalStakeDeprecated is a free data retrieval call binding the contract method 0x5f0c80cc.
//
// Solidity: function totalStake_deprecated() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) TotalStakeDeprecated(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "totalStake_deprecated")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalStakeDeprecated is a free data retrieval call binding the contract method 0x5f0c80cc.
//
// Solidity: function totalStake_deprecated() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) TotalStakeDeprecated() (*big.Int, error) {
	return _ValidatorShare.Contract.TotalStakeDeprecated(&_ValidatorShare.CallOpts)
}

// TotalStakeDeprecated is a free data retrieval call binding the contract method 0x5f0c80cc.
//
// Solidity: function totalStake_deprecated() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) TotalStakeDeprecated() (*big.Int, error) {
	return _ValidatorShare.Contract.TotalStakeDeprecated(&_ValidatorShare.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) TotalSupply() (*big.Int, error) {
	return _ValidatorShare.Contract.TotalSupply(&_ValidatorShare.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) TotalSupply() (*big.Int, error) {
	return _ValidatorShare.Contract.TotalSupply(&_ValidatorShare.CallOpts)
}

// UnbondNonces is a free data retrieval call binding the contract method 0x3046c204.
//
// Solidity: function unbondNonces(address ) view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) UnbondNonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "unbondNonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnbondNonces is a free data retrieval call binding the contract method 0x3046c204.
//
// Solidity: function unbondNonces(address ) view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) UnbondNonces(arg0 common.Address) (*big.Int, error) {
	return _ValidatorShare.Contract.UnbondNonces(&_ValidatorShare.CallOpts, arg0)
}

// UnbondNonces is a free data retrieval call binding the contract method 0x3046c204.
//
// Solidity: function unbondNonces(address ) view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) UnbondNonces(arg0 common.Address) (*big.Int, error) {
	return _ValidatorShare.Contract.UnbondNonces(&_ValidatorShare.CallOpts, arg0)
}

// Unbonds is a free data retrieval call binding the contract method 0x653ec134.
//
// Solidity: function unbonds(address ) view returns(uint256 shares, uint256 withdrawEpoch)
func (_ValidatorShare *ValidatorShareCaller) Unbonds(opts *bind.CallOpts, arg0 common.Address) (struct {
	Shares        *big.Int
	WithdrawEpoch *big.Int
}, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "unbonds", arg0)

	outstruct := new(struct {
		Shares        *big.Int
		WithdrawEpoch *big.Int
	})

	outstruct.Shares = out[0].(*big.Int)
	outstruct.WithdrawEpoch = out[1].(*big.Int)

	return *outstruct, err

}

// Unbonds is a free data retrieval call binding the contract method 0x653ec134.
//
// Solidity: function unbonds(address ) view returns(uint256 shares, uint256 withdrawEpoch)
func (_ValidatorShare *ValidatorShareSession) Unbonds(arg0 common.Address) (struct {
	Shares        *big.Int
	WithdrawEpoch *big.Int
}, error) {
	return _ValidatorShare.Contract.Unbonds(&_ValidatorShare.CallOpts, arg0)
}

// Unbonds is a free data retrieval call binding the contract method 0x653ec134.
//
// Solidity: function unbonds(address ) view returns(uint256 shares, uint256 withdrawEpoch)
func (_ValidatorShare *ValidatorShareCallerSession) Unbonds(arg0 common.Address) (struct {
	Shares        *big.Int
	WithdrawEpoch *big.Int
}, error) {
	return _ValidatorShare.Contract.Unbonds(&_ValidatorShare.CallOpts, arg0)
}

// UnbondsNew is a free data retrieval call binding the contract method 0x795be587.
//
// Solidity: function unbonds_new(address , uint256 ) view returns(uint256 shares, uint256 withdrawEpoch)
func (_ValidatorShare *ValidatorShareCaller) UnbondsNew(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	Shares        *big.Int
	WithdrawEpoch *big.Int
}, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "unbonds_new", arg0, arg1)

	outstruct := new(struct {
		Shares        *big.Int
		WithdrawEpoch *big.Int
	})

	outstruct.Shares = out[0].(*big.Int)
	outstruct.WithdrawEpoch = out[1].(*big.Int)

	return *outstruct, err

}

// UnbondsNew is a free data retrieval call binding the contract method 0x795be587.
//
// Solidity: function unbonds_new(address , uint256 ) view returns(uint256 shares, uint256 withdrawEpoch)
func (_ValidatorShare *ValidatorShareSession) UnbondsNew(arg0 common.Address, arg1 *big.Int) (struct {
	Shares        *big.Int
	WithdrawEpoch *big.Int
}, error) {
	return _ValidatorShare.Contract.UnbondsNew(&_ValidatorShare.CallOpts, arg0, arg1)
}

// UnbondsNew is a free data retrieval call binding the contract method 0x795be587.
//
// Solidity: function unbonds_new(address , uint256 ) view returns(uint256 shares, uint256 withdrawEpoch)
func (_ValidatorShare *ValidatorShareCallerSession) UnbondsNew(arg0 common.Address, arg1 *big.Int) (struct {
	Shares        *big.Int
	WithdrawEpoch *big.Int
}, error) {
	return _ValidatorShare.Contract.UnbondsNew(&_ValidatorShare.CallOpts, arg0, arg1)
}

// ValidatorId is a free data retrieval call binding the contract method 0x5c5f7dae.
//
// Solidity: function validatorId() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) ValidatorId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "validatorId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorId is a free data retrieval call binding the contract method 0x5c5f7dae.
//
// Solidity: function validatorId() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) ValidatorId() (*big.Int, error) {
	return _ValidatorShare.Contract.ValidatorId(&_ValidatorShare.CallOpts)
}

// ValidatorId is a free data retrieval call binding the contract method 0x5c5f7dae.
//
// Solidity: function validatorId() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) ValidatorId() (*big.Int, error) {
	return _ValidatorShare.Contract.ValidatorId(&_ValidatorShare.CallOpts)
}

// ValidatorRewardsDeprecated is a free data retrieval call binding the contract method 0x39c31e93.
//
// Solidity: function validatorRewards_deprecated() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) ValidatorRewardsDeprecated(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "validatorRewards_deprecated")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorRewardsDeprecated is a free data retrieval call binding the contract method 0x39c31e93.
//
// Solidity: function validatorRewards_deprecated() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) ValidatorRewardsDeprecated() (*big.Int, error) {
	return _ValidatorShare.Contract.ValidatorRewardsDeprecated(&_ValidatorShare.CallOpts)
}

// ValidatorRewardsDeprecated is a free data retrieval call binding the contract method 0x39c31e93.
//
// Solidity: function validatorRewards_deprecated() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) ValidatorRewardsDeprecated() (*big.Int, error) {
	return _ValidatorShare.Contract.ValidatorRewardsDeprecated(&_ValidatorShare.CallOpts)
}

// WithdrawExchangeRate is a free data retrieval call binding the contract method 0xbfb18f29.
//
// Solidity: function withdrawExchangeRate() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) WithdrawExchangeRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "withdrawExchangeRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawExchangeRate is a free data retrieval call binding the contract method 0xbfb18f29.
//
// Solidity: function withdrawExchangeRate() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) WithdrawExchangeRate() (*big.Int, error) {
	return _ValidatorShare.Contract.WithdrawExchangeRate(&_ValidatorShare.CallOpts)
}

// WithdrawExchangeRate is a free data retrieval call binding the contract method 0xbfb18f29.
//
// Solidity: function withdrawExchangeRate() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) WithdrawExchangeRate() (*big.Int, error) {
	return _ValidatorShare.Contract.WithdrawExchangeRate(&_ValidatorShare.CallOpts)
}

// WithdrawPool is a free data retrieval call binding the contract method 0x5c42c733.
//
// Solidity: function withdrawPool() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) WithdrawPool(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "withdrawPool")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawPool is a free data retrieval call binding the contract method 0x5c42c733.
//
// Solidity: function withdrawPool() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) WithdrawPool() (*big.Int, error) {
	return _ValidatorShare.Contract.WithdrawPool(&_ValidatorShare.CallOpts)
}

// WithdrawPool is a free data retrieval call binding the contract method 0x5c42c733.
//
// Solidity: function withdrawPool() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) WithdrawPool() (*big.Int, error) {
	return _ValidatorShare.Contract.WithdrawPool(&_ValidatorShare.CallOpts)
}

// WithdrawShares is a free data retrieval call binding the contract method 0x8d086da4.
//
// Solidity: function withdrawShares() view returns(uint256)
func (_ValidatorShare *ValidatorShareCaller) WithdrawShares(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorShare.contract.Call(opts, &out, "withdrawShares")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawShares is a free data retrieval call binding the contract method 0x8d086da4.
//
// Solidity: function withdrawShares() view returns(uint256)
func (_ValidatorShare *ValidatorShareSession) WithdrawShares() (*big.Int, error) {
	return _ValidatorShare.Contract.WithdrawShares(&_ValidatorShare.CallOpts)
}

// WithdrawShares is a free data retrieval call binding the contract method 0x8d086da4.
//
// Solidity: function withdrawShares() view returns(uint256)
func (_ValidatorShare *ValidatorShareCallerSession) WithdrawShares() (*big.Int, error) {
	return _ValidatorShare.Contract.WithdrawShares(&_ValidatorShare.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ValidatorShare *ValidatorShareTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ValidatorShare *ValidatorShareSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.Approve(&_ValidatorShare.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ValidatorShare *ValidatorShareTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.Approve(&_ValidatorShare.TransactOpts, spender, value)
}

// BuyVoucher is a paid mutator transaction binding the contract method 0x6ab15071.
//
// Solidity: function buyVoucher(uint256 _amount, uint256 _minSharesToMint) returns(uint256 amountToDeposit)
func (_ValidatorShare *ValidatorShareTransactor) BuyVoucher(opts *bind.TransactOpts, _amount *big.Int, _minSharesToMint *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "buyVoucher", _amount, _minSharesToMint)
}

// BuyVoucher is a paid mutator transaction binding the contract method 0x6ab15071.
//
// Solidity: function buyVoucher(uint256 _amount, uint256 _minSharesToMint) returns(uint256 amountToDeposit)
func (_ValidatorShare *ValidatorShareSession) BuyVoucher(_amount *big.Int, _minSharesToMint *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.BuyVoucher(&_ValidatorShare.TransactOpts, _amount, _minSharesToMint)
}

// BuyVoucher is a paid mutator transaction binding the contract method 0x6ab15071.
//
// Solidity: function buyVoucher(uint256 _amount, uint256 _minSharesToMint) returns(uint256 amountToDeposit)
func (_ValidatorShare *ValidatorShareTransactorSession) BuyVoucher(_amount *big.Int, _minSharesToMint *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.BuyVoucher(&_ValidatorShare.TransactOpts, _amount, _minSharesToMint)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ValidatorShare *ValidatorShareTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ValidatorShare *ValidatorShareSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.DecreaseAllowance(&_ValidatorShare.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ValidatorShare *ValidatorShareTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.DecreaseAllowance(&_ValidatorShare.TransactOpts, spender, subtractedValue)
}

// Drain is a paid mutator transaction binding the contract method 0xabf59fc9.
//
// Solidity: function drain(address token, address destination, uint256 amount) returns()
func (_ValidatorShare *ValidatorShareTransactor) Drain(opts *bind.TransactOpts, token common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "drain", token, destination, amount)
}

// Drain is a paid mutator transaction binding the contract method 0xabf59fc9.
//
// Solidity: function drain(address token, address destination, uint256 amount) returns()
func (_ValidatorShare *ValidatorShareSession) Drain(token common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.Drain(&_ValidatorShare.TransactOpts, token, destination, amount)
}

// Drain is a paid mutator transaction binding the contract method 0xabf59fc9.
//
// Solidity: function drain(address token, address destination, uint256 amount) returns()
func (_ValidatorShare *ValidatorShareTransactorSession) Drain(token common.Address, destination common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.Drain(&_ValidatorShare.TransactOpts, token, destination, amount)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ValidatorShare *ValidatorShareTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ValidatorShare *ValidatorShareSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.IncreaseAllowance(&_ValidatorShare.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ValidatorShare *ValidatorShareTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.IncreaseAllowance(&_ValidatorShare.TransactOpts, spender, addedValue)
}

// Initialize is a paid mutator transaction binding the contract method 0xb4988fd0.
//
// Solidity: function initialize(uint256 _validatorId, address _stakingLogger, address _stakeManager) returns()
func (_ValidatorShare *ValidatorShareTransactor) Initialize(opts *bind.TransactOpts, _validatorId *big.Int, _stakingLogger common.Address, _stakeManager common.Address) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "initialize", _validatorId, _stakingLogger, _stakeManager)
}

// Initialize is a paid mutator transaction binding the contract method 0xb4988fd0.
//
// Solidity: function initialize(uint256 _validatorId, address _stakingLogger, address _stakeManager) returns()
func (_ValidatorShare *ValidatorShareSession) Initialize(_validatorId *big.Int, _stakingLogger common.Address, _stakeManager common.Address) (*types.Transaction, error) {
	return _ValidatorShare.Contract.Initialize(&_ValidatorShare.TransactOpts, _validatorId, _stakingLogger, _stakeManager)
}

// Initialize is a paid mutator transaction binding the contract method 0xb4988fd0.
//
// Solidity: function initialize(uint256 _validatorId, address _stakingLogger, address _stakeManager) returns()
func (_ValidatorShare *ValidatorShareTransactorSession) Initialize(_validatorId *big.Int, _stakingLogger common.Address, _stakeManager common.Address) (*types.Transaction, error) {
	return _ValidatorShare.Contract.Initialize(&_ValidatorShare.TransactOpts, _validatorId, _stakingLogger, _stakeManager)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_ValidatorShare *ValidatorShareTransactor) Lock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "lock")
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_ValidatorShare *ValidatorShareSession) Lock() (*types.Transaction, error) {
	return _ValidatorShare.Contract.Lock(&_ValidatorShare.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_ValidatorShare *ValidatorShareTransactorSession) Lock() (*types.Transaction, error) {
	return _ValidatorShare.Contract.Lock(&_ValidatorShare.TransactOpts)
}

// MigrateIn is a paid mutator transaction binding the contract method 0xa0c1ca34.
//
// Solidity: function migrateIn(address user, uint256 amount) returns()
func (_ValidatorShare *ValidatorShareTransactor) MigrateIn(opts *bind.TransactOpts, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "migrateIn", user, amount)
}

// MigrateIn is a paid mutator transaction binding the contract method 0xa0c1ca34.
//
// Solidity: function migrateIn(address user, uint256 amount) returns()
func (_ValidatorShare *ValidatorShareSession) MigrateIn(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.MigrateIn(&_ValidatorShare.TransactOpts, user, amount)
}

// MigrateIn is a paid mutator transaction binding the contract method 0xa0c1ca34.
//
// Solidity: function migrateIn(address user, uint256 amount) returns()
func (_ValidatorShare *ValidatorShareTransactorSession) MigrateIn(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.MigrateIn(&_ValidatorShare.TransactOpts, user, amount)
}

// MigrateOut is a paid mutator transaction binding the contract method 0x6e7ce591.
//
// Solidity: function migrateOut(address user, uint256 amount) returns()
func (_ValidatorShare *ValidatorShareTransactor) MigrateOut(opts *bind.TransactOpts, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "migrateOut", user, amount)
}

// MigrateOut is a paid mutator transaction binding the contract method 0x6e7ce591.
//
// Solidity: function migrateOut(address user, uint256 amount) returns()
func (_ValidatorShare *ValidatorShareSession) MigrateOut(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.MigrateOut(&_ValidatorShare.TransactOpts, user, amount)
}

// MigrateOut is a paid mutator transaction binding the contract method 0x6e7ce591.
//
// Solidity: function migrateOut(address user, uint256 amount) returns()
func (_ValidatorShare *ValidatorShareTransactorSession) MigrateOut(user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.MigrateOut(&_ValidatorShare.TransactOpts, user, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ValidatorShare *ValidatorShareTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ValidatorShare *ValidatorShareSession) RenounceOwnership() (*types.Transaction, error) {
	return _ValidatorShare.Contract.RenounceOwnership(&_ValidatorShare.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ValidatorShare *ValidatorShareTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ValidatorShare.Contract.RenounceOwnership(&_ValidatorShare.TransactOpts)
}

// Restake is a paid mutator transaction binding the contract method 0x4f91440d.
//
// Solidity: function restake() returns(uint256, uint256)
func (_ValidatorShare *ValidatorShareTransactor) Restake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "restake")
}

// Restake is a paid mutator transaction binding the contract method 0x4f91440d.
//
// Solidity: function restake() returns(uint256, uint256)
func (_ValidatorShare *ValidatorShareSession) Restake() (*types.Transaction, error) {
	return _ValidatorShare.Contract.Restake(&_ValidatorShare.TransactOpts)
}

// Restake is a paid mutator transaction binding the contract method 0x4f91440d.
//
// Solidity: function restake() returns(uint256, uint256)
func (_ValidatorShare *ValidatorShareTransactorSession) Restake() (*types.Transaction, error) {
	return _ValidatorShare.Contract.Restake(&_ValidatorShare.TransactOpts)
}

// SellVoucher is a paid mutator transaction binding the contract method 0x029d3040.
//
// Solidity: function sellVoucher(uint256 claimAmount, uint256 maximumSharesToBurn) returns()
func (_ValidatorShare *ValidatorShareTransactor) SellVoucher(opts *bind.TransactOpts, claimAmount *big.Int, maximumSharesToBurn *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "sellVoucher", claimAmount, maximumSharesToBurn)
}

// SellVoucher is a paid mutator transaction binding the contract method 0x029d3040.
//
// Solidity: function sellVoucher(uint256 claimAmount, uint256 maximumSharesToBurn) returns()
func (_ValidatorShare *ValidatorShareSession) SellVoucher(claimAmount *big.Int, maximumSharesToBurn *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.SellVoucher(&_ValidatorShare.TransactOpts, claimAmount, maximumSharesToBurn)
}

// SellVoucher is a paid mutator transaction binding the contract method 0x029d3040.
//
// Solidity: function sellVoucher(uint256 claimAmount, uint256 maximumSharesToBurn) returns()
func (_ValidatorShare *ValidatorShareTransactorSession) SellVoucher(claimAmount *big.Int, maximumSharesToBurn *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.SellVoucher(&_ValidatorShare.TransactOpts, claimAmount, maximumSharesToBurn)
}

// SellVoucherNew is a paid mutator transaction binding the contract method 0xc83ec04d.
//
// Solidity: function sellVoucher_new(uint256 claimAmount, uint256 maximumSharesToBurn) returns()
func (_ValidatorShare *ValidatorShareTransactor) SellVoucherNew(opts *bind.TransactOpts, claimAmount *big.Int, maximumSharesToBurn *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "sellVoucher_new", claimAmount, maximumSharesToBurn)
}

// SellVoucherNew is a paid mutator transaction binding the contract method 0xc83ec04d.
//
// Solidity: function sellVoucher_new(uint256 claimAmount, uint256 maximumSharesToBurn) returns()
func (_ValidatorShare *ValidatorShareSession) SellVoucherNew(claimAmount *big.Int, maximumSharesToBurn *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.SellVoucherNew(&_ValidatorShare.TransactOpts, claimAmount, maximumSharesToBurn)
}

// SellVoucherNew is a paid mutator transaction binding the contract method 0xc83ec04d.
//
// Solidity: function sellVoucher_new(uint256 claimAmount, uint256 maximumSharesToBurn) returns()
func (_ValidatorShare *ValidatorShareTransactorSession) SellVoucherNew(claimAmount *big.Int, maximumSharesToBurn *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.SellVoucherNew(&_ValidatorShare.TransactOpts, claimAmount, maximumSharesToBurn)
}

// Slash is a paid mutator transaction binding the contract method 0x6cbb6050.
//
// Solidity: function slash(uint256 validatorStake, uint256 delegatedAmount, uint256 totalAmountToSlash) returns(uint256)
func (_ValidatorShare *ValidatorShareTransactor) Slash(opts *bind.TransactOpts, validatorStake *big.Int, delegatedAmount *big.Int, totalAmountToSlash *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "slash", validatorStake, delegatedAmount, totalAmountToSlash)
}

// Slash is a paid mutator transaction binding the contract method 0x6cbb6050.
//
// Solidity: function slash(uint256 validatorStake, uint256 delegatedAmount, uint256 totalAmountToSlash) returns(uint256)
func (_ValidatorShare *ValidatorShareSession) Slash(validatorStake *big.Int, delegatedAmount *big.Int, totalAmountToSlash *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.Slash(&_ValidatorShare.TransactOpts, validatorStake, delegatedAmount, totalAmountToSlash)
}

// Slash is a paid mutator transaction binding the contract method 0x6cbb6050.
//
// Solidity: function slash(uint256 validatorStake, uint256 delegatedAmount, uint256 totalAmountToSlash) returns(uint256)
func (_ValidatorShare *ValidatorShareTransactorSession) Slash(validatorStake *big.Int, delegatedAmount *big.Int, totalAmountToSlash *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.Slash(&_ValidatorShare.TransactOpts, validatorStake, delegatedAmount, totalAmountToSlash)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ValidatorShare *ValidatorShareTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ValidatorShare *ValidatorShareSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.Transfer(&_ValidatorShare.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ValidatorShare *ValidatorShareTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.Transfer(&_ValidatorShare.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ValidatorShare *ValidatorShareTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ValidatorShare *ValidatorShareSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.TransferFrom(&_ValidatorShare.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ValidatorShare *ValidatorShareTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.TransferFrom(&_ValidatorShare.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ValidatorShare *ValidatorShareTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ValidatorShare *ValidatorShareSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ValidatorShare.Contract.TransferOwnership(&_ValidatorShare.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ValidatorShare *ValidatorShareTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ValidatorShare.Contract.TransferOwnership(&_ValidatorShare.TransactOpts, newOwner)
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_ValidatorShare *ValidatorShareTransactor) Unlock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "unlock")
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_ValidatorShare *ValidatorShareSession) Unlock() (*types.Transaction, error) {
	return _ValidatorShare.Contract.Unlock(&_ValidatorShare.TransactOpts)
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_ValidatorShare *ValidatorShareTransactorSession) Unlock() (*types.Transaction, error) {
	return _ValidatorShare.Contract.Unlock(&_ValidatorShare.TransactOpts)
}

// UnstakeClaimTokens is a paid mutator transaction binding the contract method 0x8d16a14a.
//
// Solidity: function unstakeClaimTokens() returns()
func (_ValidatorShare *ValidatorShareTransactor) UnstakeClaimTokens(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "unstakeClaimTokens")
}

// UnstakeClaimTokens is a paid mutator transaction binding the contract method 0x8d16a14a.
//
// Solidity: function unstakeClaimTokens() returns()
func (_ValidatorShare *ValidatorShareSession) UnstakeClaimTokens() (*types.Transaction, error) {
	return _ValidatorShare.Contract.UnstakeClaimTokens(&_ValidatorShare.TransactOpts)
}

// UnstakeClaimTokens is a paid mutator transaction binding the contract method 0x8d16a14a.
//
// Solidity: function unstakeClaimTokens() returns()
func (_ValidatorShare *ValidatorShareTransactorSession) UnstakeClaimTokens() (*types.Transaction, error) {
	return _ValidatorShare.Contract.UnstakeClaimTokens(&_ValidatorShare.TransactOpts)
}

// UnstakeClaimTokensNew is a paid mutator transaction binding the contract method 0xe97fddc2.
//
// Solidity: function unstakeClaimTokens_new(uint256 unbondNonce) returns()
func (_ValidatorShare *ValidatorShareTransactor) UnstakeClaimTokensNew(opts *bind.TransactOpts, unbondNonce *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "unstakeClaimTokens_new", unbondNonce)
}

// UnstakeClaimTokensNew is a paid mutator transaction binding the contract method 0xe97fddc2.
//
// Solidity: function unstakeClaimTokens_new(uint256 unbondNonce) returns()
func (_ValidatorShare *ValidatorShareSession) UnstakeClaimTokensNew(unbondNonce *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.UnstakeClaimTokensNew(&_ValidatorShare.TransactOpts, unbondNonce)
}

// UnstakeClaimTokensNew is a paid mutator transaction binding the contract method 0xe97fddc2.
//
// Solidity: function unstakeClaimTokens_new(uint256 unbondNonce) returns()
func (_ValidatorShare *ValidatorShareTransactorSession) UnstakeClaimTokensNew(unbondNonce *big.Int) (*types.Transaction, error) {
	return _ValidatorShare.Contract.UnstakeClaimTokensNew(&_ValidatorShare.TransactOpts, unbondNonce)
}

// UpdateDelegation is a paid mutator transaction binding the contract method 0x7ba8c820.
//
// Solidity: function updateDelegation(bool _delegation) returns()
func (_ValidatorShare *ValidatorShareTransactor) UpdateDelegation(opts *bind.TransactOpts, _delegation bool) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "updateDelegation", _delegation)
}

// UpdateDelegation is a paid mutator transaction binding the contract method 0x7ba8c820.
//
// Solidity: function updateDelegation(bool _delegation) returns()
func (_ValidatorShare *ValidatorShareSession) UpdateDelegation(_delegation bool) (*types.Transaction, error) {
	return _ValidatorShare.Contract.UpdateDelegation(&_ValidatorShare.TransactOpts, _delegation)
}

// UpdateDelegation is a paid mutator transaction binding the contract method 0x7ba8c820.
//
// Solidity: function updateDelegation(bool _delegation) returns()
func (_ValidatorShare *ValidatorShareTransactorSession) UpdateDelegation(_delegation bool) (*types.Transaction, error) {
	return _ValidatorShare.Contract.UpdateDelegation(&_ValidatorShare.TransactOpts, _delegation)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0xc7b8981c.
//
// Solidity: function withdrawRewards() returns()
func (_ValidatorShare *ValidatorShareTransactor) WithdrawRewards(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorShare.contract.Transact(opts, "withdrawRewards")
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0xc7b8981c.
//
// Solidity: function withdrawRewards() returns()
func (_ValidatorShare *ValidatorShareSession) WithdrawRewards() (*types.Transaction, error) {
	return _ValidatorShare.Contract.WithdrawRewards(&_ValidatorShare.TransactOpts)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0xc7b8981c.
//
// Solidity: function withdrawRewards() returns()
func (_ValidatorShare *ValidatorShareTransactorSession) WithdrawRewards() (*types.Transaction, error) {
	return _ValidatorShare.Contract.WithdrawRewards(&_ValidatorShare.TransactOpts)
}

// ValidatorShareApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ValidatorShare contract.
type ValidatorShareApprovalIterator struct {
	Event *ValidatorShareApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ValidatorShareApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorShareApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ValidatorShareApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ValidatorShareApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorShareApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorShareApproval represents a Approval event raised by the ValidatorShare contract.
type ValidatorShareApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ValidatorShare *ValidatorShareFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ValidatorShareApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ValidatorShare.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorShareApprovalIterator{contract: _ValidatorShare.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ValidatorShare *ValidatorShareFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ValidatorShareApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ValidatorShare.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorShareApproval)
				if err := _ValidatorShare.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ValidatorShare *ValidatorShareFilterer) ParseApproval(log types.Log) (*ValidatorShareApproval, error) {
	event := new(ValidatorShareApproval)
	if err := _ValidatorShare.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ValidatorShareOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ValidatorShare contract.
type ValidatorShareOwnershipTransferredIterator struct {
	Event *ValidatorShareOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ValidatorShareOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorShareOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ValidatorShareOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ValidatorShareOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorShareOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorShareOwnershipTransferred represents a OwnershipTransferred event raised by the ValidatorShare contract.
type ValidatorShareOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ValidatorShare *ValidatorShareFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ValidatorShareOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ValidatorShare.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorShareOwnershipTransferredIterator{contract: _ValidatorShare.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ValidatorShare *ValidatorShareFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ValidatorShareOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ValidatorShare.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorShareOwnershipTransferred)
				if err := _ValidatorShare.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ValidatorShare *ValidatorShareFilterer) ParseOwnershipTransferred(log types.Log) (*ValidatorShareOwnershipTransferred, error) {
	event := new(ValidatorShareOwnershipTransferred)
	if err := _ValidatorShare.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ValidatorShareTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ValidatorShare contract.
type ValidatorShareTransferIterator struct {
	Event *ValidatorShareTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ValidatorShareTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorShareTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ValidatorShareTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ValidatorShareTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorShareTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorShareTransfer represents a Transfer event raised by the ValidatorShare contract.
type ValidatorShareTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ValidatorShare *ValidatorShareFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ValidatorShareTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ValidatorShare.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorShareTransferIterator{contract: _ValidatorShare.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ValidatorShare *ValidatorShareFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ValidatorShareTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ValidatorShare.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorShareTransfer)
				if err := _ValidatorShare.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ValidatorShare *ValidatorShareFilterer) ParseTransfer(log types.Log) (*ValidatorShareTransfer, error) {
	event := new(ValidatorShareTransfer)
	if err := _ValidatorShare.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}
