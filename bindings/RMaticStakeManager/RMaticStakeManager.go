// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rmatic_stake_manager

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// RMaticStakeManagerMetaData contains all meta data concerning the RMaticStakeManager contract.
var RMaticStakeManagerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Delegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"era\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"ExecuteNewEra\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"govDelegated\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"localDelegated\",\"type\":\"uint256\"}],\"name\":\"RepairDelegated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unbondingDuration\",\"type\":\"uint256\"}],\"name\":\"SetUnbondingDuration\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"era\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"name\":\"Settle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"poolAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rTokenAmount\",\"type\":\"uint256\"}],\"name\":\"Stake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Undelegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"poolAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rTokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"burnAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unstakeIndex\",\"type\":\"uint256\"}],\"name\":\"Unstake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"poolAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"unstakeIndexList\",\"type\":\"uint256[]\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poolAddress\",\"type\":\"address\"}],\"name\":\"addStakePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poolAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentEra\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eraOffset\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"eraRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eraSeconds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"erc20TokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBondedPools\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"pools\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_staker\",\"type\":\"address\"}],\"name\":\"getUnstakeIndexListOf\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"unstakeIndexList\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poolAddress\",\"type\":\"address\"}],\"name\":\"getValidatorIdsOf\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"validatorIds\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20TokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_unbondingDuration\",\"type\":\"uint256\"}],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestEra\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poolAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_validatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_govDelegated\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_bond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_unbond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalRTokenSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalProtocolFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_era\",\"type\":\"uint256\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newEra\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextUnstakeIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"poolInfoOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"era\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"active\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"protocolFeeCommission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rateChangeLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poolAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_srcValidatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_dstValidatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"redelegate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poolAddress\",\"type\":\"address\"}],\"name\":\"rmStakePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_unstakeFeeCommission\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_protocolFeeCommission\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minStakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_unbondingDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rateChangeLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_eraSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_eraOffset\",\"type\":\"uint256\"}],\"name\":\"setParams\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_stakeAmount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poolAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_stakeAmount\",\"type\":\"uint256\"}],\"name\":\"stakeWithPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalProtocolFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRTokenSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unbondingDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rTokenAmount\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"unstakeAtIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"era\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unstakeFeeCommission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poolAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rTokenAmount\",\"type\":\"uint256\"}],\"name\":\"unstakeWithPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"withdrawProtocolFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poolAddress\",\"type\":\"address\"}],\"name\":\"withdrawWithPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RMaticStakeManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use RMaticStakeManagerMetaData.ABI instead.
var RMaticStakeManagerABI = RMaticStakeManagerMetaData.ABI

// RMaticStakeManager is an auto generated Go binding around an Ethereum contract.
type RMaticStakeManager struct {
	RMaticStakeManagerCaller     // Read-only binding to the contract
	RMaticStakeManagerTransactor // Write-only binding to the contract
	RMaticStakeManagerFilterer   // Log filterer for contract events
}

// RMaticStakeManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type RMaticStakeManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RMaticStakeManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RMaticStakeManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RMaticStakeManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RMaticStakeManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RMaticStakeManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RMaticStakeManagerSession struct {
	Contract     *RMaticStakeManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RMaticStakeManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RMaticStakeManagerCallerSession struct {
	Contract *RMaticStakeManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// RMaticStakeManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RMaticStakeManagerTransactorSession struct {
	Contract     *RMaticStakeManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// RMaticStakeManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type RMaticStakeManagerRaw struct {
	Contract *RMaticStakeManager // Generic contract binding to access the raw methods on
}

// RMaticStakeManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RMaticStakeManagerCallerRaw struct {
	Contract *RMaticStakeManagerCaller // Generic read-only contract binding to access the raw methods on
}

// RMaticStakeManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RMaticStakeManagerTransactorRaw struct {
	Contract *RMaticStakeManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRMaticStakeManager creates a new instance of RMaticStakeManager, bound to a specific deployed contract.
func NewRMaticStakeManager(address common.Address, backend bind.ContractBackend) (*RMaticStakeManager, error) {
	contract, err := bindRMaticStakeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManager{RMaticStakeManagerCaller: RMaticStakeManagerCaller{contract: contract}, RMaticStakeManagerTransactor: RMaticStakeManagerTransactor{contract: contract}, RMaticStakeManagerFilterer: RMaticStakeManagerFilterer{contract: contract}}, nil
}

// NewRMaticStakeManagerCaller creates a new read-only instance of RMaticStakeManager, bound to a specific deployed contract.
func NewRMaticStakeManagerCaller(address common.Address, caller bind.ContractCaller) (*RMaticStakeManagerCaller, error) {
	contract, err := bindRMaticStakeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManagerCaller{contract: contract}, nil
}

// NewRMaticStakeManagerTransactor creates a new write-only instance of RMaticStakeManager, bound to a specific deployed contract.
func NewRMaticStakeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*RMaticStakeManagerTransactor, error) {
	contract, err := bindRMaticStakeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManagerTransactor{contract: contract}, nil
}

// NewRMaticStakeManagerFilterer creates a new log filterer instance of RMaticStakeManager, bound to a specific deployed contract.
func NewRMaticStakeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*RMaticStakeManagerFilterer, error) {
	contract, err := bindRMaticStakeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManagerFilterer{contract: contract}, nil
}

// bindRMaticStakeManager binds a generic wrapper to an already deployed contract.
func bindRMaticStakeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RMaticStakeManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RMaticStakeManager *RMaticStakeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMaticStakeManager.Contract.RMaticStakeManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RMaticStakeManager *RMaticStakeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.RMaticStakeManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RMaticStakeManager *RMaticStakeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.RMaticStakeManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RMaticStakeManager *RMaticStakeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RMaticStakeManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RMaticStakeManager *RMaticStakeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RMaticStakeManager *RMaticStakeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_RMaticStakeManager *RMaticStakeManagerCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_RMaticStakeManager *RMaticStakeManagerSession) Admin() (common.Address, error) {
	return _RMaticStakeManager.Contract.Admin(&_RMaticStakeManager.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) Admin() (common.Address, error) {
	return _RMaticStakeManager.Contract.Admin(&_RMaticStakeManager.CallOpts)
}

// CurrentEra is a free data retrieval call binding the contract method 0x973628f6.
//
// Solidity: function currentEra() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) CurrentEra(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "currentEra")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentEra is a free data retrieval call binding the contract method 0x973628f6.
//
// Solidity: function currentEra() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) CurrentEra() (*big.Int, error) {
	return _RMaticStakeManager.Contract.CurrentEra(&_RMaticStakeManager.CallOpts)
}

// CurrentEra is a free data retrieval call binding the contract method 0x973628f6.
//
// Solidity: function currentEra() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) CurrentEra() (*big.Int, error) {
	return _RMaticStakeManager.Contract.CurrentEra(&_RMaticStakeManager.CallOpts)
}

// EraOffset is a free data retrieval call binding the contract method 0xc8c20263.
//
// Solidity: function eraOffset() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) EraOffset(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "eraOffset")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EraOffset is a free data retrieval call binding the contract method 0xc8c20263.
//
// Solidity: function eraOffset() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) EraOffset() (*big.Int, error) {
	return _RMaticStakeManager.Contract.EraOffset(&_RMaticStakeManager.CallOpts)
}

// EraOffset is a free data retrieval call binding the contract method 0xc8c20263.
//
// Solidity: function eraOffset() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) EraOffset() (*big.Int, error) {
	return _RMaticStakeManager.Contract.EraOffset(&_RMaticStakeManager.CallOpts)
}

// EraRate is a free data retrieval call binding the contract method 0xeb8ad76e.
//
// Solidity: function eraRate(uint256 ) view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) EraRate(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "eraRate", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EraRate is a free data retrieval call binding the contract method 0xeb8ad76e.
//
// Solidity: function eraRate(uint256 ) view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) EraRate(arg0 *big.Int) (*big.Int, error) {
	return _RMaticStakeManager.Contract.EraRate(&_RMaticStakeManager.CallOpts, arg0)
}

// EraRate is a free data retrieval call binding the contract method 0xeb8ad76e.
//
// Solidity: function eraRate(uint256 ) view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) EraRate(arg0 *big.Int) (*big.Int, error) {
	return _RMaticStakeManager.Contract.EraRate(&_RMaticStakeManager.CallOpts, arg0)
}

// EraSeconds is a free data retrieval call binding the contract method 0xe81f1553.
//
// Solidity: function eraSeconds() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) EraSeconds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "eraSeconds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EraSeconds is a free data retrieval call binding the contract method 0xe81f1553.
//
// Solidity: function eraSeconds() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) EraSeconds() (*big.Int, error) {
	return _RMaticStakeManager.Contract.EraSeconds(&_RMaticStakeManager.CallOpts)
}

// EraSeconds is a free data retrieval call binding the contract method 0xe81f1553.
//
// Solidity: function eraSeconds() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) EraSeconds() (*big.Int, error) {
	return _RMaticStakeManager.Contract.EraSeconds(&_RMaticStakeManager.CallOpts)
}

// Erc20TokenAddress is a free data retrieval call binding the contract method 0xf835cd3c.
//
// Solidity: function erc20TokenAddress() view returns(address)
func (_RMaticStakeManager *RMaticStakeManagerCaller) Erc20TokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "erc20TokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Erc20TokenAddress is a free data retrieval call binding the contract method 0xf835cd3c.
//
// Solidity: function erc20TokenAddress() view returns(address)
func (_RMaticStakeManager *RMaticStakeManagerSession) Erc20TokenAddress() (common.Address, error) {
	return _RMaticStakeManager.Contract.Erc20TokenAddress(&_RMaticStakeManager.CallOpts)
}

// Erc20TokenAddress is a free data retrieval call binding the contract method 0xf835cd3c.
//
// Solidity: function erc20TokenAddress() view returns(address)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) Erc20TokenAddress() (common.Address, error) {
	return _RMaticStakeManager.Contract.Erc20TokenAddress(&_RMaticStakeManager.CallOpts)
}

// GetBondedPools is a free data retrieval call binding the contract method 0x58af8bf0.
//
// Solidity: function getBondedPools() view returns(address[] pools)
func (_RMaticStakeManager *RMaticStakeManagerCaller) GetBondedPools(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "getBondedPools")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetBondedPools is a free data retrieval call binding the contract method 0x58af8bf0.
//
// Solidity: function getBondedPools() view returns(address[] pools)
func (_RMaticStakeManager *RMaticStakeManagerSession) GetBondedPools() ([]common.Address, error) {
	return _RMaticStakeManager.Contract.GetBondedPools(&_RMaticStakeManager.CallOpts)
}

// GetBondedPools is a free data retrieval call binding the contract method 0x58af8bf0.
//
// Solidity: function getBondedPools() view returns(address[] pools)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) GetBondedPools() ([]common.Address, error) {
	return _RMaticStakeManager.Contract.GetBondedPools(&_RMaticStakeManager.CallOpts)
}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) GetRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "getRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) GetRate() (*big.Int, error) {
	return _RMaticStakeManager.Contract.GetRate(&_RMaticStakeManager.CallOpts)
}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) GetRate() (*big.Int, error) {
	return _RMaticStakeManager.Contract.GetRate(&_RMaticStakeManager.CallOpts)
}

// GetUnstakeIndexListOf is a free data retrieval call binding the contract method 0x615acc36.
//
// Solidity: function getUnstakeIndexListOf(address _staker) view returns(uint256[] unstakeIndexList)
func (_RMaticStakeManager *RMaticStakeManagerCaller) GetUnstakeIndexListOf(opts *bind.CallOpts, _staker common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "getUnstakeIndexListOf", _staker)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetUnstakeIndexListOf is a free data retrieval call binding the contract method 0x615acc36.
//
// Solidity: function getUnstakeIndexListOf(address _staker) view returns(uint256[] unstakeIndexList)
func (_RMaticStakeManager *RMaticStakeManagerSession) GetUnstakeIndexListOf(_staker common.Address) ([]*big.Int, error) {
	return _RMaticStakeManager.Contract.GetUnstakeIndexListOf(&_RMaticStakeManager.CallOpts, _staker)
}

// GetUnstakeIndexListOf is a free data retrieval call binding the contract method 0x615acc36.
//
// Solidity: function getUnstakeIndexListOf(address _staker) view returns(uint256[] unstakeIndexList)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) GetUnstakeIndexListOf(_staker common.Address) ([]*big.Int, error) {
	return _RMaticStakeManager.Contract.GetUnstakeIndexListOf(&_RMaticStakeManager.CallOpts, _staker)
}

// GetValidatorIdsOf is a free data retrieval call binding the contract method 0x12faf11a.
//
// Solidity: function getValidatorIdsOf(address _poolAddress) view returns(uint256[] validatorIds)
func (_RMaticStakeManager *RMaticStakeManagerCaller) GetValidatorIdsOf(opts *bind.CallOpts, _poolAddress common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "getValidatorIdsOf", _poolAddress)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetValidatorIdsOf is a free data retrieval call binding the contract method 0x12faf11a.
//
// Solidity: function getValidatorIdsOf(address _poolAddress) view returns(uint256[] validatorIds)
func (_RMaticStakeManager *RMaticStakeManagerSession) GetValidatorIdsOf(_poolAddress common.Address) ([]*big.Int, error) {
	return _RMaticStakeManager.Contract.GetValidatorIdsOf(&_RMaticStakeManager.CallOpts, _poolAddress)
}

// GetValidatorIdsOf is a free data retrieval call binding the contract method 0x12faf11a.
//
// Solidity: function getValidatorIdsOf(address _poolAddress) view returns(uint256[] validatorIds)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) GetValidatorIdsOf(_poolAddress common.Address) ([]*big.Int, error) {
	return _RMaticStakeManager.Contract.GetValidatorIdsOf(&_RMaticStakeManager.CallOpts, _poolAddress)
}

// LatestEra is a free data retrieval call binding the contract method 0x3f6f5f32.
//
// Solidity: function latestEra() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) LatestEra(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "latestEra")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LatestEra is a free data retrieval call binding the contract method 0x3f6f5f32.
//
// Solidity: function latestEra() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) LatestEra() (*big.Int, error) {
	return _RMaticStakeManager.Contract.LatestEra(&_RMaticStakeManager.CallOpts)
}

// LatestEra is a free data retrieval call binding the contract method 0x3f6f5f32.
//
// Solidity: function latestEra() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) LatestEra() (*big.Int, error) {
	return _RMaticStakeManager.Contract.LatestEra(&_RMaticStakeManager.CallOpts)
}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) MinStakeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "minStakeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) MinStakeAmount() (*big.Int, error) {
	return _RMaticStakeManager.Contract.MinStakeAmount(&_RMaticStakeManager.CallOpts)
}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) MinStakeAmount() (*big.Int, error) {
	return _RMaticStakeManager.Contract.MinStakeAmount(&_RMaticStakeManager.CallOpts)
}

// NextUnstakeIndex is a free data retrieval call binding the contract method 0x3bea9ee3.
//
// Solidity: function nextUnstakeIndex() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) NextUnstakeIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "nextUnstakeIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextUnstakeIndex is a free data retrieval call binding the contract method 0x3bea9ee3.
//
// Solidity: function nextUnstakeIndex() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) NextUnstakeIndex() (*big.Int, error) {
	return _RMaticStakeManager.Contract.NextUnstakeIndex(&_RMaticStakeManager.CallOpts)
}

// NextUnstakeIndex is a free data retrieval call binding the contract method 0x3bea9ee3.
//
// Solidity: function nextUnstakeIndex() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) NextUnstakeIndex() (*big.Int, error) {
	return _RMaticStakeManager.Contract.NextUnstakeIndex(&_RMaticStakeManager.CallOpts)
}

// PoolInfoOf is a free data retrieval call binding the contract method 0x008b5dd2.
//
// Solidity: function poolInfoOf(address ) view returns(uint256 era, uint256 bond, uint256 unbond, uint256 active)
func (_RMaticStakeManager *RMaticStakeManagerCaller) PoolInfoOf(opts *bind.CallOpts, arg0 common.Address) (struct {
	Era    *big.Int
	Bond   *big.Int
	Unbond *big.Int
	Active *big.Int
}, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "poolInfoOf", arg0)

	outstruct := new(struct {
		Era    *big.Int
		Bond   *big.Int
		Unbond *big.Int
		Active *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Era = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Bond = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Unbond = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Active = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// PoolInfoOf is a free data retrieval call binding the contract method 0x008b5dd2.
//
// Solidity: function poolInfoOf(address ) view returns(uint256 era, uint256 bond, uint256 unbond, uint256 active)
func (_RMaticStakeManager *RMaticStakeManagerSession) PoolInfoOf(arg0 common.Address) (struct {
	Era    *big.Int
	Bond   *big.Int
	Unbond *big.Int
	Active *big.Int
}, error) {
	return _RMaticStakeManager.Contract.PoolInfoOf(&_RMaticStakeManager.CallOpts, arg0)
}

// PoolInfoOf is a free data retrieval call binding the contract method 0x008b5dd2.
//
// Solidity: function poolInfoOf(address ) view returns(uint256 era, uint256 bond, uint256 unbond, uint256 active)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) PoolInfoOf(arg0 common.Address) (struct {
	Era    *big.Int
	Bond   *big.Int
	Unbond *big.Int
	Active *big.Int
}, error) {
	return _RMaticStakeManager.Contract.PoolInfoOf(&_RMaticStakeManager.CallOpts, arg0)
}

// ProtocolFeeCommission is a free data retrieval call binding the contract method 0x19301c26.
//
// Solidity: function protocolFeeCommission() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) ProtocolFeeCommission(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "protocolFeeCommission")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProtocolFeeCommission is a free data retrieval call binding the contract method 0x19301c26.
//
// Solidity: function protocolFeeCommission() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) ProtocolFeeCommission() (*big.Int, error) {
	return _RMaticStakeManager.Contract.ProtocolFeeCommission(&_RMaticStakeManager.CallOpts)
}

// ProtocolFeeCommission is a free data retrieval call binding the contract method 0x19301c26.
//
// Solidity: function protocolFeeCommission() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) ProtocolFeeCommission() (*big.Int, error) {
	return _RMaticStakeManager.Contract.ProtocolFeeCommission(&_RMaticStakeManager.CallOpts)
}

// RTokenAddress is a free data retrieval call binding the contract method 0x35381358.
//
// Solidity: function rTokenAddress() view returns(address)
func (_RMaticStakeManager *RMaticStakeManagerCaller) RTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "rTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RTokenAddress is a free data retrieval call binding the contract method 0x35381358.
//
// Solidity: function rTokenAddress() view returns(address)
func (_RMaticStakeManager *RMaticStakeManagerSession) RTokenAddress() (common.Address, error) {
	return _RMaticStakeManager.Contract.RTokenAddress(&_RMaticStakeManager.CallOpts)
}

// RTokenAddress is a free data retrieval call binding the contract method 0x35381358.
//
// Solidity: function rTokenAddress() view returns(address)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) RTokenAddress() (common.Address, error) {
	return _RMaticStakeManager.Contract.RTokenAddress(&_RMaticStakeManager.CallOpts)
}

// RateChangeLimit is a free data retrieval call binding the contract method 0xc0152b71.
//
// Solidity: function rateChangeLimit() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) RateChangeLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "rateChangeLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RateChangeLimit is a free data retrieval call binding the contract method 0xc0152b71.
//
// Solidity: function rateChangeLimit() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) RateChangeLimit() (*big.Int, error) {
	return _RMaticStakeManager.Contract.RateChangeLimit(&_RMaticStakeManager.CallOpts)
}

// RateChangeLimit is a free data retrieval call binding the contract method 0xc0152b71.
//
// Solidity: function rateChangeLimit() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) RateChangeLimit() (*big.Int, error) {
	return _RMaticStakeManager.Contract.RateChangeLimit(&_RMaticStakeManager.CallOpts)
}

// TotalProtocolFee is a free data retrieval call binding the contract method 0x88611f35.
//
// Solidity: function totalProtocolFee() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) TotalProtocolFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "totalProtocolFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalProtocolFee is a free data retrieval call binding the contract method 0x88611f35.
//
// Solidity: function totalProtocolFee() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) TotalProtocolFee() (*big.Int, error) {
	return _RMaticStakeManager.Contract.TotalProtocolFee(&_RMaticStakeManager.CallOpts)
}

// TotalProtocolFee is a free data retrieval call binding the contract method 0x88611f35.
//
// Solidity: function totalProtocolFee() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) TotalProtocolFee() (*big.Int, error) {
	return _RMaticStakeManager.Contract.TotalProtocolFee(&_RMaticStakeManager.CallOpts)
}

// TotalRTokenSupply is a free data retrieval call binding the contract method 0x7a7c27c0.
//
// Solidity: function totalRTokenSupply() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) TotalRTokenSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "totalRTokenSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalRTokenSupply is a free data retrieval call binding the contract method 0x7a7c27c0.
//
// Solidity: function totalRTokenSupply() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) TotalRTokenSupply() (*big.Int, error) {
	return _RMaticStakeManager.Contract.TotalRTokenSupply(&_RMaticStakeManager.CallOpts)
}

// TotalRTokenSupply is a free data retrieval call binding the contract method 0x7a7c27c0.
//
// Solidity: function totalRTokenSupply() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) TotalRTokenSupply() (*big.Int, error) {
	return _RMaticStakeManager.Contract.TotalRTokenSupply(&_RMaticStakeManager.CallOpts)
}

// UnbondingDuration is a free data retrieval call binding the contract method 0xccf6802a.
//
// Solidity: function unbondingDuration() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) UnbondingDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "unbondingDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnbondingDuration is a free data retrieval call binding the contract method 0xccf6802a.
//
// Solidity: function unbondingDuration() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) UnbondingDuration() (*big.Int, error) {
	return _RMaticStakeManager.Contract.UnbondingDuration(&_RMaticStakeManager.CallOpts)
}

// UnbondingDuration is a free data retrieval call binding the contract method 0xccf6802a.
//
// Solidity: function unbondingDuration() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) UnbondingDuration() (*big.Int, error) {
	return _RMaticStakeManager.Contract.UnbondingDuration(&_RMaticStakeManager.CallOpts)
}

// UnstakeAtIndex is a free data retrieval call binding the contract method 0x6e436c6e.
//
// Solidity: function unstakeAtIndex(uint256 ) view returns(uint256 era, address pool, address receiver, uint256 amount)
func (_RMaticStakeManager *RMaticStakeManagerCaller) UnstakeAtIndex(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Era      *big.Int
	Pool     common.Address
	Receiver common.Address
	Amount   *big.Int
}, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "unstakeAtIndex", arg0)

	outstruct := new(struct {
		Era      *big.Int
		Pool     common.Address
		Receiver common.Address
		Amount   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Era = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Pool = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Receiver = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// UnstakeAtIndex is a free data retrieval call binding the contract method 0x6e436c6e.
//
// Solidity: function unstakeAtIndex(uint256 ) view returns(uint256 era, address pool, address receiver, uint256 amount)
func (_RMaticStakeManager *RMaticStakeManagerSession) UnstakeAtIndex(arg0 *big.Int) (struct {
	Era      *big.Int
	Pool     common.Address
	Receiver common.Address
	Amount   *big.Int
}, error) {
	return _RMaticStakeManager.Contract.UnstakeAtIndex(&_RMaticStakeManager.CallOpts, arg0)
}

// UnstakeAtIndex is a free data retrieval call binding the contract method 0x6e436c6e.
//
// Solidity: function unstakeAtIndex(uint256 ) view returns(uint256 era, address pool, address receiver, uint256 amount)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) UnstakeAtIndex(arg0 *big.Int) (struct {
	Era      *big.Int
	Pool     common.Address
	Receiver common.Address
	Amount   *big.Int
}, error) {
	return _RMaticStakeManager.Contract.UnstakeAtIndex(&_RMaticStakeManager.CallOpts, arg0)
}

// UnstakeFeeCommission is a free data retrieval call binding the contract method 0xe1bb5897.
//
// Solidity: function unstakeFeeCommission() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCaller) UnstakeFeeCommission(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RMaticStakeManager.contract.Call(opts, &out, "unstakeFeeCommission")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakeFeeCommission is a free data retrieval call binding the contract method 0xe1bb5897.
//
// Solidity: function unstakeFeeCommission() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerSession) UnstakeFeeCommission() (*big.Int, error) {
	return _RMaticStakeManager.Contract.UnstakeFeeCommission(&_RMaticStakeManager.CallOpts)
}

// UnstakeFeeCommission is a free data retrieval call binding the contract method 0xe1bb5897.
//
// Solidity: function unstakeFeeCommission() view returns(uint256)
func (_RMaticStakeManager *RMaticStakeManagerCallerSession) UnstakeFeeCommission() (*big.Int, error) {
	return _RMaticStakeManager.Contract.UnstakeFeeCommission(&_RMaticStakeManager.CallOpts)
}

// AddStakePool is a paid mutator transaction binding the contract method 0x0f772a1d.
//
// Solidity: function addStakePool(address _poolAddress) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) AddStakePool(opts *bind.TransactOpts, _poolAddress common.Address) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "addStakePool", _poolAddress)
}

// AddStakePool is a paid mutator transaction binding the contract method 0x0f772a1d.
//
// Solidity: function addStakePool(address _poolAddress) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) AddStakePool(_poolAddress common.Address) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.AddStakePool(&_RMaticStakeManager.TransactOpts, _poolAddress)
}

// AddStakePool is a paid mutator transaction binding the contract method 0x0f772a1d.
//
// Solidity: function addStakePool(address _poolAddress) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) AddStakePool(_poolAddress common.Address) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.AddStakePool(&_RMaticStakeManager.TransactOpts, _poolAddress)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _poolAddress, uint256 _amount) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) Approve(opts *bind.TransactOpts, _poolAddress common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "approve", _poolAddress, _amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _poolAddress, uint256 _amount) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) Approve(_poolAddress common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Approve(&_RMaticStakeManager.TransactOpts, _poolAddress, _amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _poolAddress, uint256 _amount) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) Approve(_poolAddress common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Approve(&_RMaticStakeManager.TransactOpts, _poolAddress, _amount)
}

// Init is a paid mutator transaction binding the contract method 0x86863ec6.
//
// Solidity: function init(address _rTokenAddress, address _erc20TokenAddress, uint256 _unbondingDuration) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) Init(opts *bind.TransactOpts, _rTokenAddress common.Address, _erc20TokenAddress common.Address, _unbondingDuration *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "init", _rTokenAddress, _erc20TokenAddress, _unbondingDuration)
}

// Init is a paid mutator transaction binding the contract method 0x86863ec6.
//
// Solidity: function init(address _rTokenAddress, address _erc20TokenAddress, uint256 _unbondingDuration) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) Init(_rTokenAddress common.Address, _erc20TokenAddress common.Address, _unbondingDuration *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Init(&_RMaticStakeManager.TransactOpts, _rTokenAddress, _erc20TokenAddress, _unbondingDuration)
}

// Init is a paid mutator transaction binding the contract method 0x86863ec6.
//
// Solidity: function init(address _rTokenAddress, address _erc20TokenAddress, uint256 _unbondingDuration) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) Init(_rTokenAddress common.Address, _erc20TokenAddress common.Address, _unbondingDuration *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Init(&_RMaticStakeManager.TransactOpts, _rTokenAddress, _erc20TokenAddress, _unbondingDuration)
}

// Migrate is a paid mutator transaction binding the contract method 0x3b634e07.
//
// Solidity: function migrate(address _poolAddress, uint256 _validatorId, uint256 _govDelegated, uint256 _bond, uint256 _unbond, uint256 _rate, uint256 _totalRTokenSupply, uint256 _totalProtocolFee, uint256 _era) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) Migrate(opts *bind.TransactOpts, _poolAddress common.Address, _validatorId *big.Int, _govDelegated *big.Int, _bond *big.Int, _unbond *big.Int, _rate *big.Int, _totalRTokenSupply *big.Int, _totalProtocolFee *big.Int, _era *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "migrate", _poolAddress, _validatorId, _govDelegated, _bond, _unbond, _rate, _totalRTokenSupply, _totalProtocolFee, _era)
}

// Migrate is a paid mutator transaction binding the contract method 0x3b634e07.
//
// Solidity: function migrate(address _poolAddress, uint256 _validatorId, uint256 _govDelegated, uint256 _bond, uint256 _unbond, uint256 _rate, uint256 _totalRTokenSupply, uint256 _totalProtocolFee, uint256 _era) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) Migrate(_poolAddress common.Address, _validatorId *big.Int, _govDelegated *big.Int, _bond *big.Int, _unbond *big.Int, _rate *big.Int, _totalRTokenSupply *big.Int, _totalProtocolFee *big.Int, _era *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Migrate(&_RMaticStakeManager.TransactOpts, _poolAddress, _validatorId, _govDelegated, _bond, _unbond, _rate, _totalRTokenSupply, _totalProtocolFee, _era)
}

// Migrate is a paid mutator transaction binding the contract method 0x3b634e07.
//
// Solidity: function migrate(address _poolAddress, uint256 _validatorId, uint256 _govDelegated, uint256 _bond, uint256 _unbond, uint256 _rate, uint256 _totalRTokenSupply, uint256 _totalProtocolFee, uint256 _era) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) Migrate(_poolAddress common.Address, _validatorId *big.Int, _govDelegated *big.Int, _bond *big.Int, _unbond *big.Int, _rate *big.Int, _totalRTokenSupply *big.Int, _totalProtocolFee *big.Int, _era *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Migrate(&_RMaticStakeManager.TransactOpts, _poolAddress, _validatorId, _govDelegated, _bond, _unbond, _rate, _totalRTokenSupply, _totalProtocolFee, _era)
}

// NewEra is a paid mutator transaction binding the contract method 0x7b207727.
//
// Solidity: function newEra() returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) NewEra(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "newEra")
}

// NewEra is a paid mutator transaction binding the contract method 0x7b207727.
//
// Solidity: function newEra() returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) NewEra() (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.NewEra(&_RMaticStakeManager.TransactOpts)
}

// NewEra is a paid mutator transaction binding the contract method 0x7b207727.
//
// Solidity: function newEra() returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) NewEra() (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.NewEra(&_RMaticStakeManager.TransactOpts)
}

// Redelegate is a paid mutator transaction binding the contract method 0x53c91acc.
//
// Solidity: function redelegate(address _poolAddress, uint256 _srcValidatorId, uint256 _dstValidatorId, uint256 _amount) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) Redelegate(opts *bind.TransactOpts, _poolAddress common.Address, _srcValidatorId *big.Int, _dstValidatorId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "redelegate", _poolAddress, _srcValidatorId, _dstValidatorId, _amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x53c91acc.
//
// Solidity: function redelegate(address _poolAddress, uint256 _srcValidatorId, uint256 _dstValidatorId, uint256 _amount) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) Redelegate(_poolAddress common.Address, _srcValidatorId *big.Int, _dstValidatorId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Redelegate(&_RMaticStakeManager.TransactOpts, _poolAddress, _srcValidatorId, _dstValidatorId, _amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x53c91acc.
//
// Solidity: function redelegate(address _poolAddress, uint256 _srcValidatorId, uint256 _dstValidatorId, uint256 _amount) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) Redelegate(_poolAddress common.Address, _srcValidatorId *big.Int, _dstValidatorId *big.Int, _amount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Redelegate(&_RMaticStakeManager.TransactOpts, _poolAddress, _srcValidatorId, _dstValidatorId, _amount)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _poolAddress) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) RmStakePool(opts *bind.TransactOpts, _poolAddress common.Address) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "rmStakePool", _poolAddress)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _poolAddress) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) RmStakePool(_poolAddress common.Address) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.RmStakePool(&_RMaticStakeManager.TransactOpts, _poolAddress)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _poolAddress) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) RmStakePool(_poolAddress common.Address) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.RmStakePool(&_RMaticStakeManager.TransactOpts, _poolAddress)
}

// SetParams is a paid mutator transaction binding the contract method 0xce1c4556.
//
// Solidity: function setParams(uint256 _unstakeFeeCommission, uint256 _protocolFeeCommission, uint256 _minStakeAmount, uint256 _unbondingDuration, uint256 _rateChangeLimit, uint256 _eraSeconds, uint256 _eraOffset) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) SetParams(opts *bind.TransactOpts, _unstakeFeeCommission *big.Int, _protocolFeeCommission *big.Int, _minStakeAmount *big.Int, _unbondingDuration *big.Int, _rateChangeLimit *big.Int, _eraSeconds *big.Int, _eraOffset *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "setParams", _unstakeFeeCommission, _protocolFeeCommission, _minStakeAmount, _unbondingDuration, _rateChangeLimit, _eraSeconds, _eraOffset)
}

// SetParams is a paid mutator transaction binding the contract method 0xce1c4556.
//
// Solidity: function setParams(uint256 _unstakeFeeCommission, uint256 _protocolFeeCommission, uint256 _minStakeAmount, uint256 _unbondingDuration, uint256 _rateChangeLimit, uint256 _eraSeconds, uint256 _eraOffset) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) SetParams(_unstakeFeeCommission *big.Int, _protocolFeeCommission *big.Int, _minStakeAmount *big.Int, _unbondingDuration *big.Int, _rateChangeLimit *big.Int, _eraSeconds *big.Int, _eraOffset *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.SetParams(&_RMaticStakeManager.TransactOpts, _unstakeFeeCommission, _protocolFeeCommission, _minStakeAmount, _unbondingDuration, _rateChangeLimit, _eraSeconds, _eraOffset)
}

// SetParams is a paid mutator transaction binding the contract method 0xce1c4556.
//
// Solidity: function setParams(uint256 _unstakeFeeCommission, uint256 _protocolFeeCommission, uint256 _minStakeAmount, uint256 _unbondingDuration, uint256 _rateChangeLimit, uint256 _eraSeconds, uint256 _eraOffset) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) SetParams(_unstakeFeeCommission *big.Int, _protocolFeeCommission *big.Int, _minStakeAmount *big.Int, _unbondingDuration *big.Int, _rateChangeLimit *big.Int, _eraSeconds *big.Int, _eraOffset *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.SetParams(&_RMaticStakeManager.TransactOpts, _unstakeFeeCommission, _protocolFeeCommission, _minStakeAmount, _unbondingDuration, _rateChangeLimit, _eraSeconds, _eraOffset)
}

// Stake is a paid mutator transaction binding the contract method 0xa694fc3a.
//
// Solidity: function stake(uint256 _stakeAmount) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) Stake(opts *bind.TransactOpts, _stakeAmount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "stake", _stakeAmount)
}

// Stake is a paid mutator transaction binding the contract method 0xa694fc3a.
//
// Solidity: function stake(uint256 _stakeAmount) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) Stake(_stakeAmount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Stake(&_RMaticStakeManager.TransactOpts, _stakeAmount)
}

// Stake is a paid mutator transaction binding the contract method 0xa694fc3a.
//
// Solidity: function stake(uint256 _stakeAmount) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) Stake(_stakeAmount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Stake(&_RMaticStakeManager.TransactOpts, _stakeAmount)
}

// StakeWithPool is a paid mutator transaction binding the contract method 0x1525be32.
//
// Solidity: function stakeWithPool(address _poolAddress, uint256 _stakeAmount) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) StakeWithPool(opts *bind.TransactOpts, _poolAddress common.Address, _stakeAmount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "stakeWithPool", _poolAddress, _stakeAmount)
}

// StakeWithPool is a paid mutator transaction binding the contract method 0x1525be32.
//
// Solidity: function stakeWithPool(address _poolAddress, uint256 _stakeAmount) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) StakeWithPool(_poolAddress common.Address, _stakeAmount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.StakeWithPool(&_RMaticStakeManager.TransactOpts, _poolAddress, _stakeAmount)
}

// StakeWithPool is a paid mutator transaction binding the contract method 0x1525be32.
//
// Solidity: function stakeWithPool(address _poolAddress, uint256 _stakeAmount) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) StakeWithPool(_poolAddress common.Address, _stakeAmount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.StakeWithPool(&_RMaticStakeManager.TransactOpts, _poolAddress, _stakeAmount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 _rTokenAmount) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) Unstake(opts *bind.TransactOpts, _rTokenAmount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "unstake", _rTokenAmount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 _rTokenAmount) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) Unstake(_rTokenAmount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Unstake(&_RMaticStakeManager.TransactOpts, _rTokenAmount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 _rTokenAmount) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) Unstake(_rTokenAmount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Unstake(&_RMaticStakeManager.TransactOpts, _rTokenAmount)
}

// UnstakeWithPool is a paid mutator transaction binding the contract method 0xb608b458.
//
// Solidity: function unstakeWithPool(address _poolAddress, uint256 _rTokenAmount) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) UnstakeWithPool(opts *bind.TransactOpts, _poolAddress common.Address, _rTokenAmount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "unstakeWithPool", _poolAddress, _rTokenAmount)
}

// UnstakeWithPool is a paid mutator transaction binding the contract method 0xb608b458.
//
// Solidity: function unstakeWithPool(address _poolAddress, uint256 _rTokenAmount) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) UnstakeWithPool(_poolAddress common.Address, _rTokenAmount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.UnstakeWithPool(&_RMaticStakeManager.TransactOpts, _poolAddress, _rTokenAmount)
}

// UnstakeWithPool is a paid mutator transaction binding the contract method 0xb608b458.
//
// Solidity: function unstakeWithPool(address _poolAddress, uint256 _rTokenAmount) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) UnstakeWithPool(_poolAddress common.Address, _rTokenAmount *big.Int) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.UnstakeWithPool(&_RMaticStakeManager.TransactOpts, _poolAddress, _rTokenAmount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) Withdraw() (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Withdraw(&_RMaticStakeManager.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) Withdraw() (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.Withdraw(&_RMaticStakeManager.TransactOpts)
}

// WithdrawProtocolFee is a paid mutator transaction binding the contract method 0x668fb6dc.
//
// Solidity: function withdrawProtocolFee(address _to) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) WithdrawProtocolFee(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "withdrawProtocolFee", _to)
}

// WithdrawProtocolFee is a paid mutator transaction binding the contract method 0x668fb6dc.
//
// Solidity: function withdrawProtocolFee(address _to) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) WithdrawProtocolFee(_to common.Address) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.WithdrawProtocolFee(&_RMaticStakeManager.TransactOpts, _to)
}

// WithdrawProtocolFee is a paid mutator transaction binding the contract method 0x668fb6dc.
//
// Solidity: function withdrawProtocolFee(address _to) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) WithdrawProtocolFee(_to common.Address) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.WithdrawProtocolFee(&_RMaticStakeManager.TransactOpts, _to)
}

// WithdrawWithPool is a paid mutator transaction binding the contract method 0xf737abac.
//
// Solidity: function withdrawWithPool(address _poolAddress) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactor) WithdrawWithPool(opts *bind.TransactOpts, _poolAddress common.Address) (*types.Transaction, error) {
	return _RMaticStakeManager.contract.Transact(opts, "withdrawWithPool", _poolAddress)
}

// WithdrawWithPool is a paid mutator transaction binding the contract method 0xf737abac.
//
// Solidity: function withdrawWithPool(address _poolAddress) returns()
func (_RMaticStakeManager *RMaticStakeManagerSession) WithdrawWithPool(_poolAddress common.Address) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.WithdrawWithPool(&_RMaticStakeManager.TransactOpts, _poolAddress)
}

// WithdrawWithPool is a paid mutator transaction binding the contract method 0xf737abac.
//
// Solidity: function withdrawWithPool(address _poolAddress) returns()
func (_RMaticStakeManager *RMaticStakeManagerTransactorSession) WithdrawWithPool(_poolAddress common.Address) (*types.Transaction, error) {
	return _RMaticStakeManager.Contract.WithdrawWithPool(&_RMaticStakeManager.TransactOpts, _poolAddress)
}

// RMaticStakeManagerDelegateIterator is returned from FilterDelegate and is used to iterate over the raw logs and unpacked data for Delegate events raised by the RMaticStakeManager contract.
type RMaticStakeManagerDelegateIterator struct {
	Event *RMaticStakeManagerDelegate // Event containing the contract specifics and raw log

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
func (it *RMaticStakeManagerDelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMaticStakeManagerDelegate)
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
		it.Event = new(RMaticStakeManagerDelegate)
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
func (it *RMaticStakeManagerDelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RMaticStakeManagerDelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RMaticStakeManagerDelegate represents a Delegate event raised by the RMaticStakeManager contract.
type RMaticStakeManagerDelegate struct {
	Pool      common.Address
	Validator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDelegate is a free log retrieval operation binding the contract event 0x510b11bb3f3c799b11307c01ab7db0d335683ef5b2da98f7697de744f465eacc.
//
// Solidity: event Delegate(address pool, address validator, uint256 amount)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) FilterDelegate(opts *bind.FilterOpts) (*RMaticStakeManagerDelegateIterator, error) {

	logs, sub, err := _RMaticStakeManager.contract.FilterLogs(opts, "Delegate")
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManagerDelegateIterator{contract: _RMaticStakeManager.contract, event: "Delegate", logs: logs, sub: sub}, nil
}

// WatchDelegate is a free log subscription operation binding the contract event 0x510b11bb3f3c799b11307c01ab7db0d335683ef5b2da98f7697de744f465eacc.
//
// Solidity: event Delegate(address pool, address validator, uint256 amount)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) WatchDelegate(opts *bind.WatchOpts, sink chan<- *RMaticStakeManagerDelegate) (event.Subscription, error) {

	logs, sub, err := _RMaticStakeManager.contract.WatchLogs(opts, "Delegate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RMaticStakeManagerDelegate)
				if err := _RMaticStakeManager.contract.UnpackLog(event, "Delegate", log); err != nil {
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

// ParseDelegate is a log parse operation binding the contract event 0x510b11bb3f3c799b11307c01ab7db0d335683ef5b2da98f7697de744f465eacc.
//
// Solidity: event Delegate(address pool, address validator, uint256 amount)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) ParseDelegate(log types.Log) (*RMaticStakeManagerDelegate, error) {
	event := new(RMaticStakeManagerDelegate)
	if err := _RMaticStakeManager.contract.UnpackLog(event, "Delegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RMaticStakeManagerExecuteNewEraIterator is returned from FilterExecuteNewEra and is used to iterate over the raw logs and unpacked data for ExecuteNewEra events raised by the RMaticStakeManager contract.
type RMaticStakeManagerExecuteNewEraIterator struct {
	Event *RMaticStakeManagerExecuteNewEra // Event containing the contract specifics and raw log

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
func (it *RMaticStakeManagerExecuteNewEraIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMaticStakeManagerExecuteNewEra)
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
		it.Event = new(RMaticStakeManagerExecuteNewEra)
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
func (it *RMaticStakeManagerExecuteNewEraIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RMaticStakeManagerExecuteNewEraIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RMaticStakeManagerExecuteNewEra represents a ExecuteNewEra event raised by the RMaticStakeManager contract.
type RMaticStakeManagerExecuteNewEra struct {
	Era  *big.Int
	Rate *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterExecuteNewEra is a free log retrieval operation binding the contract event 0x02105621fc31aa3ac04a9845beacd54c700e2ab23ff8acdd755dfd878ae61f02.
//
// Solidity: event ExecuteNewEra(uint256 indexed era, uint256 rate)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) FilterExecuteNewEra(opts *bind.FilterOpts, era []*big.Int) (*RMaticStakeManagerExecuteNewEraIterator, error) {

	var eraRule []interface{}
	for _, eraItem := range era {
		eraRule = append(eraRule, eraItem)
	}

	logs, sub, err := _RMaticStakeManager.contract.FilterLogs(opts, "ExecuteNewEra", eraRule)
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManagerExecuteNewEraIterator{contract: _RMaticStakeManager.contract, event: "ExecuteNewEra", logs: logs, sub: sub}, nil
}

// WatchExecuteNewEra is a free log subscription operation binding the contract event 0x02105621fc31aa3ac04a9845beacd54c700e2ab23ff8acdd755dfd878ae61f02.
//
// Solidity: event ExecuteNewEra(uint256 indexed era, uint256 rate)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) WatchExecuteNewEra(opts *bind.WatchOpts, sink chan<- *RMaticStakeManagerExecuteNewEra, era []*big.Int) (event.Subscription, error) {

	var eraRule []interface{}
	for _, eraItem := range era {
		eraRule = append(eraRule, eraItem)
	}

	logs, sub, err := _RMaticStakeManager.contract.WatchLogs(opts, "ExecuteNewEra", eraRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RMaticStakeManagerExecuteNewEra)
				if err := _RMaticStakeManager.contract.UnpackLog(event, "ExecuteNewEra", log); err != nil {
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

// ParseExecuteNewEra is a log parse operation binding the contract event 0x02105621fc31aa3ac04a9845beacd54c700e2ab23ff8acdd755dfd878ae61f02.
//
// Solidity: event ExecuteNewEra(uint256 indexed era, uint256 rate)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) ParseExecuteNewEra(log types.Log) (*RMaticStakeManagerExecuteNewEra, error) {
	event := new(RMaticStakeManagerExecuteNewEra)
	if err := _RMaticStakeManager.contract.UnpackLog(event, "ExecuteNewEra", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RMaticStakeManagerRepairDelegatedIterator is returned from FilterRepairDelegated and is used to iterate over the raw logs and unpacked data for RepairDelegated events raised by the RMaticStakeManager contract.
type RMaticStakeManagerRepairDelegatedIterator struct {
	Event *RMaticStakeManagerRepairDelegated // Event containing the contract specifics and raw log

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
func (it *RMaticStakeManagerRepairDelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMaticStakeManagerRepairDelegated)
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
		it.Event = new(RMaticStakeManagerRepairDelegated)
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
func (it *RMaticStakeManagerRepairDelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RMaticStakeManagerRepairDelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RMaticStakeManagerRepairDelegated represents a RepairDelegated event raised by the RMaticStakeManager contract.
type RMaticStakeManagerRepairDelegated struct {
	Pool           common.Address
	Validator      common.Address
	GovDelegated   *big.Int
	LocalDelegated *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRepairDelegated is a free log retrieval operation binding the contract event 0x373c6bebde9765bcea50da315cc4bb34983406b93947382285efab03bc4990ef.
//
// Solidity: event RepairDelegated(address pool, address validator, uint256 govDelegated, uint256 localDelegated)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) FilterRepairDelegated(opts *bind.FilterOpts) (*RMaticStakeManagerRepairDelegatedIterator, error) {

	logs, sub, err := _RMaticStakeManager.contract.FilterLogs(opts, "RepairDelegated")
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManagerRepairDelegatedIterator{contract: _RMaticStakeManager.contract, event: "RepairDelegated", logs: logs, sub: sub}, nil
}

// WatchRepairDelegated is a free log subscription operation binding the contract event 0x373c6bebde9765bcea50da315cc4bb34983406b93947382285efab03bc4990ef.
//
// Solidity: event RepairDelegated(address pool, address validator, uint256 govDelegated, uint256 localDelegated)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) WatchRepairDelegated(opts *bind.WatchOpts, sink chan<- *RMaticStakeManagerRepairDelegated) (event.Subscription, error) {

	logs, sub, err := _RMaticStakeManager.contract.WatchLogs(opts, "RepairDelegated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RMaticStakeManagerRepairDelegated)
				if err := _RMaticStakeManager.contract.UnpackLog(event, "RepairDelegated", log); err != nil {
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

// ParseRepairDelegated is a log parse operation binding the contract event 0x373c6bebde9765bcea50da315cc4bb34983406b93947382285efab03bc4990ef.
//
// Solidity: event RepairDelegated(address pool, address validator, uint256 govDelegated, uint256 localDelegated)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) ParseRepairDelegated(log types.Log) (*RMaticStakeManagerRepairDelegated, error) {
	event := new(RMaticStakeManagerRepairDelegated)
	if err := _RMaticStakeManager.contract.UnpackLog(event, "RepairDelegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RMaticStakeManagerSetUnbondingDurationIterator is returned from FilterSetUnbondingDuration and is used to iterate over the raw logs and unpacked data for SetUnbondingDuration events raised by the RMaticStakeManager contract.
type RMaticStakeManagerSetUnbondingDurationIterator struct {
	Event *RMaticStakeManagerSetUnbondingDuration // Event containing the contract specifics and raw log

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
func (it *RMaticStakeManagerSetUnbondingDurationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMaticStakeManagerSetUnbondingDuration)
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
		it.Event = new(RMaticStakeManagerSetUnbondingDuration)
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
func (it *RMaticStakeManagerSetUnbondingDurationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RMaticStakeManagerSetUnbondingDurationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RMaticStakeManagerSetUnbondingDuration represents a SetUnbondingDuration event raised by the RMaticStakeManager contract.
type RMaticStakeManagerSetUnbondingDuration struct {
	UnbondingDuration *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterSetUnbondingDuration is a free log retrieval operation binding the contract event 0x279066b062766bf26597e98ef1d6fb6ec39502061f95f271089f727c875414d0.
//
// Solidity: event SetUnbondingDuration(uint256 unbondingDuration)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) FilterSetUnbondingDuration(opts *bind.FilterOpts) (*RMaticStakeManagerSetUnbondingDurationIterator, error) {

	logs, sub, err := _RMaticStakeManager.contract.FilterLogs(opts, "SetUnbondingDuration")
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManagerSetUnbondingDurationIterator{contract: _RMaticStakeManager.contract, event: "SetUnbondingDuration", logs: logs, sub: sub}, nil
}

// WatchSetUnbondingDuration is a free log subscription operation binding the contract event 0x279066b062766bf26597e98ef1d6fb6ec39502061f95f271089f727c875414d0.
//
// Solidity: event SetUnbondingDuration(uint256 unbondingDuration)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) WatchSetUnbondingDuration(opts *bind.WatchOpts, sink chan<- *RMaticStakeManagerSetUnbondingDuration) (event.Subscription, error) {

	logs, sub, err := _RMaticStakeManager.contract.WatchLogs(opts, "SetUnbondingDuration")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RMaticStakeManagerSetUnbondingDuration)
				if err := _RMaticStakeManager.contract.UnpackLog(event, "SetUnbondingDuration", log); err != nil {
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

// ParseSetUnbondingDuration is a log parse operation binding the contract event 0x279066b062766bf26597e98ef1d6fb6ec39502061f95f271089f727c875414d0.
//
// Solidity: event SetUnbondingDuration(uint256 unbondingDuration)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) ParseSetUnbondingDuration(log types.Log) (*RMaticStakeManagerSetUnbondingDuration, error) {
	event := new(RMaticStakeManagerSetUnbondingDuration)
	if err := _RMaticStakeManager.contract.UnpackLog(event, "SetUnbondingDuration", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RMaticStakeManagerSettleIterator is returned from FilterSettle and is used to iterate over the raw logs and unpacked data for Settle events raised by the RMaticStakeManager contract.
type RMaticStakeManagerSettleIterator struct {
	Event *RMaticStakeManagerSettle // Event containing the contract specifics and raw log

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
func (it *RMaticStakeManagerSettleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMaticStakeManagerSettle)
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
		it.Event = new(RMaticStakeManagerSettle)
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
func (it *RMaticStakeManagerSettleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RMaticStakeManagerSettleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RMaticStakeManagerSettle represents a Settle event raised by the RMaticStakeManager contract.
type RMaticStakeManagerSettle struct {
	Era  *big.Int
	Pool common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSettle is a free log retrieval operation binding the contract event 0x5a208f0e8f17e3d4544100afcf65662217cc8fb06e47d9016a83e21a1f96a791.
//
// Solidity: event Settle(uint256 indexed era, address indexed pool)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) FilterSettle(opts *bind.FilterOpts, era []*big.Int, pool []common.Address) (*RMaticStakeManagerSettleIterator, error) {

	var eraRule []interface{}
	for _, eraItem := range era {
		eraRule = append(eraRule, eraItem)
	}
	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _RMaticStakeManager.contract.FilterLogs(opts, "Settle", eraRule, poolRule)
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManagerSettleIterator{contract: _RMaticStakeManager.contract, event: "Settle", logs: logs, sub: sub}, nil
}

// WatchSettle is a free log subscription operation binding the contract event 0x5a208f0e8f17e3d4544100afcf65662217cc8fb06e47d9016a83e21a1f96a791.
//
// Solidity: event Settle(uint256 indexed era, address indexed pool)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) WatchSettle(opts *bind.WatchOpts, sink chan<- *RMaticStakeManagerSettle, era []*big.Int, pool []common.Address) (event.Subscription, error) {

	var eraRule []interface{}
	for _, eraItem := range era {
		eraRule = append(eraRule, eraItem)
	}
	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _RMaticStakeManager.contract.WatchLogs(opts, "Settle", eraRule, poolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RMaticStakeManagerSettle)
				if err := _RMaticStakeManager.contract.UnpackLog(event, "Settle", log); err != nil {
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

// ParseSettle is a log parse operation binding the contract event 0x5a208f0e8f17e3d4544100afcf65662217cc8fb06e47d9016a83e21a1f96a791.
//
// Solidity: event Settle(uint256 indexed era, address indexed pool)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) ParseSettle(log types.Log) (*RMaticStakeManagerSettle, error) {
	event := new(RMaticStakeManagerSettle)
	if err := _RMaticStakeManager.contract.UnpackLog(event, "Settle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RMaticStakeManagerStakeIterator is returned from FilterStake and is used to iterate over the raw logs and unpacked data for Stake events raised by the RMaticStakeManager contract.
type RMaticStakeManagerStakeIterator struct {
	Event *RMaticStakeManagerStake // Event containing the contract specifics and raw log

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
func (it *RMaticStakeManagerStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMaticStakeManagerStake)
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
		it.Event = new(RMaticStakeManagerStake)
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
func (it *RMaticStakeManagerStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RMaticStakeManagerStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RMaticStakeManagerStake represents a Stake event raised by the RMaticStakeManager contract.
type RMaticStakeManagerStake struct {
	Staker       common.Address
	PoolAddress  common.Address
	TokenAmount  *big.Int
	RTokenAmount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStake is a free log retrieval operation binding the contract event 0x63602d0ecc7b3a0ef7ff1a116e23056662d64280355ba8031b6d0d767c4b4458.
//
// Solidity: event Stake(address staker, address poolAddress, uint256 tokenAmount, uint256 rTokenAmount)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) FilterStake(opts *bind.FilterOpts) (*RMaticStakeManagerStakeIterator, error) {

	logs, sub, err := _RMaticStakeManager.contract.FilterLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManagerStakeIterator{contract: _RMaticStakeManager.contract, event: "Stake", logs: logs, sub: sub}, nil
}

// WatchStake is a free log subscription operation binding the contract event 0x63602d0ecc7b3a0ef7ff1a116e23056662d64280355ba8031b6d0d767c4b4458.
//
// Solidity: event Stake(address staker, address poolAddress, uint256 tokenAmount, uint256 rTokenAmount)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) WatchStake(opts *bind.WatchOpts, sink chan<- *RMaticStakeManagerStake) (event.Subscription, error) {

	logs, sub, err := _RMaticStakeManager.contract.WatchLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RMaticStakeManagerStake)
				if err := _RMaticStakeManager.contract.UnpackLog(event, "Stake", log); err != nil {
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

// ParseStake is a log parse operation binding the contract event 0x63602d0ecc7b3a0ef7ff1a116e23056662d64280355ba8031b6d0d767c4b4458.
//
// Solidity: event Stake(address staker, address poolAddress, uint256 tokenAmount, uint256 rTokenAmount)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) ParseStake(log types.Log) (*RMaticStakeManagerStake, error) {
	event := new(RMaticStakeManagerStake)
	if err := _RMaticStakeManager.contract.UnpackLog(event, "Stake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RMaticStakeManagerUndelegateIterator is returned from FilterUndelegate and is used to iterate over the raw logs and unpacked data for Undelegate events raised by the RMaticStakeManager contract.
type RMaticStakeManagerUndelegateIterator struct {
	Event *RMaticStakeManagerUndelegate // Event containing the contract specifics and raw log

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
func (it *RMaticStakeManagerUndelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMaticStakeManagerUndelegate)
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
		it.Event = new(RMaticStakeManagerUndelegate)
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
func (it *RMaticStakeManagerUndelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RMaticStakeManagerUndelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RMaticStakeManagerUndelegate represents a Undelegate event raised by the RMaticStakeManager contract.
type RMaticStakeManagerUndelegate struct {
	Pool      common.Address
	Validator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUndelegate is a free log retrieval operation binding the contract event 0xbda8c0e95802a0e6788c3e9027292382d5a41b86556015f846b03a9874b2b827.
//
// Solidity: event Undelegate(address pool, address validator, uint256 amount)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) FilterUndelegate(opts *bind.FilterOpts) (*RMaticStakeManagerUndelegateIterator, error) {

	logs, sub, err := _RMaticStakeManager.contract.FilterLogs(opts, "Undelegate")
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManagerUndelegateIterator{contract: _RMaticStakeManager.contract, event: "Undelegate", logs: logs, sub: sub}, nil
}

// WatchUndelegate is a free log subscription operation binding the contract event 0xbda8c0e95802a0e6788c3e9027292382d5a41b86556015f846b03a9874b2b827.
//
// Solidity: event Undelegate(address pool, address validator, uint256 amount)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) WatchUndelegate(opts *bind.WatchOpts, sink chan<- *RMaticStakeManagerUndelegate) (event.Subscription, error) {

	logs, sub, err := _RMaticStakeManager.contract.WatchLogs(opts, "Undelegate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RMaticStakeManagerUndelegate)
				if err := _RMaticStakeManager.contract.UnpackLog(event, "Undelegate", log); err != nil {
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

// ParseUndelegate is a log parse operation binding the contract event 0xbda8c0e95802a0e6788c3e9027292382d5a41b86556015f846b03a9874b2b827.
//
// Solidity: event Undelegate(address pool, address validator, uint256 amount)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) ParseUndelegate(log types.Log) (*RMaticStakeManagerUndelegate, error) {
	event := new(RMaticStakeManagerUndelegate)
	if err := _RMaticStakeManager.contract.UnpackLog(event, "Undelegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RMaticStakeManagerUnstakeIterator is returned from FilterUnstake and is used to iterate over the raw logs and unpacked data for Unstake events raised by the RMaticStakeManager contract.
type RMaticStakeManagerUnstakeIterator struct {
	Event *RMaticStakeManagerUnstake // Event containing the contract specifics and raw log

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
func (it *RMaticStakeManagerUnstakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMaticStakeManagerUnstake)
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
		it.Event = new(RMaticStakeManagerUnstake)
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
func (it *RMaticStakeManagerUnstakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RMaticStakeManagerUnstakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RMaticStakeManagerUnstake represents a Unstake event raised by the RMaticStakeManager contract.
type RMaticStakeManagerUnstake struct {
	Staker       common.Address
	PoolAddress  common.Address
	TokenAmount  *big.Int
	RTokenAmount *big.Int
	BurnAmount   *big.Int
	UnstakeIndex *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUnstake is a free log retrieval operation binding the contract event 0x4e5916c8cf4042e78e6a44e60bf6b48e8e969ce5abf8d4227613ddc3454015ea.
//
// Solidity: event Unstake(address staker, address poolAddress, uint256 tokenAmount, uint256 rTokenAmount, uint256 burnAmount, uint256 unstakeIndex)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) FilterUnstake(opts *bind.FilterOpts) (*RMaticStakeManagerUnstakeIterator, error) {

	logs, sub, err := _RMaticStakeManager.contract.FilterLogs(opts, "Unstake")
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManagerUnstakeIterator{contract: _RMaticStakeManager.contract, event: "Unstake", logs: logs, sub: sub}, nil
}

// WatchUnstake is a free log subscription operation binding the contract event 0x4e5916c8cf4042e78e6a44e60bf6b48e8e969ce5abf8d4227613ddc3454015ea.
//
// Solidity: event Unstake(address staker, address poolAddress, uint256 tokenAmount, uint256 rTokenAmount, uint256 burnAmount, uint256 unstakeIndex)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) WatchUnstake(opts *bind.WatchOpts, sink chan<- *RMaticStakeManagerUnstake) (event.Subscription, error) {

	logs, sub, err := _RMaticStakeManager.contract.WatchLogs(opts, "Unstake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RMaticStakeManagerUnstake)
				if err := _RMaticStakeManager.contract.UnpackLog(event, "Unstake", log); err != nil {
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

// ParseUnstake is a log parse operation binding the contract event 0x4e5916c8cf4042e78e6a44e60bf6b48e8e969ce5abf8d4227613ddc3454015ea.
//
// Solidity: event Unstake(address staker, address poolAddress, uint256 tokenAmount, uint256 rTokenAmount, uint256 burnAmount, uint256 unstakeIndex)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) ParseUnstake(log types.Log) (*RMaticStakeManagerUnstake, error) {
	event := new(RMaticStakeManagerUnstake)
	if err := _RMaticStakeManager.contract.UnpackLog(event, "Unstake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RMaticStakeManagerWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the RMaticStakeManager contract.
type RMaticStakeManagerWithdrawIterator struct {
	Event *RMaticStakeManagerWithdraw // Event containing the contract specifics and raw log

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
func (it *RMaticStakeManagerWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RMaticStakeManagerWithdraw)
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
		it.Event = new(RMaticStakeManagerWithdraw)
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
func (it *RMaticStakeManagerWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RMaticStakeManagerWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RMaticStakeManagerWithdraw represents a Withdraw event raised by the RMaticStakeManager contract.
type RMaticStakeManagerWithdraw struct {
	Staker           common.Address
	PoolAddress      common.Address
	TokenAmount      *big.Int
	UnstakeIndexList []*big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xe4b7499d334dcb3a4338114f8df473bb4444d9cace993f8d2eb779921f074dd3.
//
// Solidity: event Withdraw(address staker, address poolAddress, uint256 tokenAmount, uint256[] unstakeIndexList)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) FilterWithdraw(opts *bind.FilterOpts) (*RMaticStakeManagerWithdrawIterator, error) {

	logs, sub, err := _RMaticStakeManager.contract.FilterLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return &RMaticStakeManagerWithdrawIterator{contract: _RMaticStakeManager.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xe4b7499d334dcb3a4338114f8df473bb4444d9cace993f8d2eb779921f074dd3.
//
// Solidity: event Withdraw(address staker, address poolAddress, uint256 tokenAmount, uint256[] unstakeIndexList)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *RMaticStakeManagerWithdraw) (event.Subscription, error) {

	logs, sub, err := _RMaticStakeManager.contract.WatchLogs(opts, "Withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RMaticStakeManagerWithdraw)
				if err := _RMaticStakeManager.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0xe4b7499d334dcb3a4338114f8df473bb4444d9cace993f8d2eb779921f074dd3.
//
// Solidity: event Withdraw(address staker, address poolAddress, uint256 tokenAmount, uint256[] unstakeIndexList)
func (_RMaticStakeManager *RMaticStakeManagerFilterer) ParseWithdraw(log types.Log) (*RMaticStakeManagerWithdraw, error) {
	event := new(RMaticStakeManagerWithdraw)
	if err := _RMaticStakeManager.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
