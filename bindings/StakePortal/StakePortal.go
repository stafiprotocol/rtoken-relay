// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stake_portal

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

// StakePortalMetaData contains all meta data concerning the StakePortal contract.
var StakePortalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_stakePoolAddressList\",\"type\":\"address[]\"},{\"internalType\":\"uint8[]\",\"name\":\"_chainIdList\",\"type\":\"uint8[]\"},{\"internalType\":\"address\",\"name\":\"_erc20TokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_relayFee\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stafiRecipient\",\"type\":\"bytes32\"}],\"name\":\"RecoverStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stakePool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"chainId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stafiRecipient\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destRecipient\",\"type\":\"address\"}],\"name\":\"Stake\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint8[]\",\"name\":\"_chaindIdList\",\"type\":\"uint8[]\"}],\"name\":\"addChainId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_stakePoolAddressList\",\"type\":\"address[]\"}],\"name\":\"addStakePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"chainIdExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"erc20TokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_txHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_stafiRecipient\",\"type\":\"bytes32\"}],\"name\":\"recoverStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"relayFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_chaindId\",\"type\":\"uint8\"}],\"name\":\"rmChainId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"}],\"name\":\"rmStakePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minAmount\",\"type\":\"uint256\"}],\"name\":\"setMinAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_relayFee\",\"type\":\"uint256\"}],\"name\":\"setRelayFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"_destChainId\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_stafiRecipient\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_destRecipient\",\"type\":\"address\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"stakePoolAddressExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeSwitch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"toggleSwitch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StakePortalABI is the input ABI used to generate the binding from.
// Deprecated: Use StakePortalMetaData.ABI instead.
var StakePortalABI = StakePortalMetaData.ABI

// StakePortal is an auto generated Go binding around an Ethereum contract.
type StakePortal struct {
	StakePortalCaller     // Read-only binding to the contract
	StakePortalTransactor // Write-only binding to the contract
	StakePortalFilterer   // Log filterer for contract events
}

// StakePortalCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakePortalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakePortalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakePortalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakePortalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakePortalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakePortalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakePortalSession struct {
	Contract     *StakePortal      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakePortalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakePortalCallerSession struct {
	Contract *StakePortalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StakePortalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakePortalTransactorSession struct {
	Contract     *StakePortalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StakePortalRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakePortalRaw struct {
	Contract *StakePortal // Generic contract binding to access the raw methods on
}

// StakePortalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakePortalCallerRaw struct {
	Contract *StakePortalCaller // Generic read-only contract binding to access the raw methods on
}

// StakePortalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakePortalTransactorRaw struct {
	Contract *StakePortalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakePortal creates a new instance of StakePortal, bound to a specific deployed contract.
func NewStakePortal(address common.Address, backend bind.ContractBackend) (*StakePortal, error) {
	contract, err := bindStakePortal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakePortal{StakePortalCaller: StakePortalCaller{contract: contract}, StakePortalTransactor: StakePortalTransactor{contract: contract}, StakePortalFilterer: StakePortalFilterer{contract: contract}}, nil
}

// NewStakePortalCaller creates a new read-only instance of StakePortal, bound to a specific deployed contract.
func NewStakePortalCaller(address common.Address, caller bind.ContractCaller) (*StakePortalCaller, error) {
	contract, err := bindStakePortal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakePortalCaller{contract: contract}, nil
}

// NewStakePortalTransactor creates a new write-only instance of StakePortal, bound to a specific deployed contract.
func NewStakePortalTransactor(address common.Address, transactor bind.ContractTransactor) (*StakePortalTransactor, error) {
	contract, err := bindStakePortal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakePortalTransactor{contract: contract}, nil
}

// NewStakePortalFilterer creates a new log filterer instance of StakePortal, bound to a specific deployed contract.
func NewStakePortalFilterer(address common.Address, filterer bind.ContractFilterer) (*StakePortalFilterer, error) {
	contract, err := bindStakePortal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakePortalFilterer{contract: contract}, nil
}

// bindStakePortal binds a generic wrapper to an already deployed contract.
func bindStakePortal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakePortalABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakePortal *StakePortalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakePortal.Contract.StakePortalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakePortal *StakePortalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakePortal.Contract.StakePortalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakePortal *StakePortalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakePortal.Contract.StakePortalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakePortal *StakePortalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakePortal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakePortal *StakePortalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakePortal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakePortal *StakePortalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakePortal.Contract.contract.Transact(opts, method, params...)
}

// ChainIdExist is a free data retrieval call binding the contract method 0x0133228d.
//
// Solidity: function chainIdExist(uint8 ) view returns(bool)
func (_StakePortal *StakePortalCaller) ChainIdExist(opts *bind.CallOpts, arg0 uint8) (bool, error) {
	var out []interface{}
	err := _StakePortal.contract.Call(opts, &out, "chainIdExist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ChainIdExist is a free data retrieval call binding the contract method 0x0133228d.
//
// Solidity: function chainIdExist(uint8 ) view returns(bool)
func (_StakePortal *StakePortalSession) ChainIdExist(arg0 uint8) (bool, error) {
	return _StakePortal.Contract.ChainIdExist(&_StakePortal.CallOpts, arg0)
}

// ChainIdExist is a free data retrieval call binding the contract method 0x0133228d.
//
// Solidity: function chainIdExist(uint8 ) view returns(bool)
func (_StakePortal *StakePortalCallerSession) ChainIdExist(arg0 uint8) (bool, error) {
	return _StakePortal.Contract.ChainIdExist(&_StakePortal.CallOpts, arg0)
}

// Erc20TokenAddress is a free data retrieval call binding the contract method 0xf835cd3c.
//
// Solidity: function erc20TokenAddress() view returns(address)
func (_StakePortal *StakePortalCaller) Erc20TokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakePortal.contract.Call(opts, &out, "erc20TokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Erc20TokenAddress is a free data retrieval call binding the contract method 0xf835cd3c.
//
// Solidity: function erc20TokenAddress() view returns(address)
func (_StakePortal *StakePortalSession) Erc20TokenAddress() (common.Address, error) {
	return _StakePortal.Contract.Erc20TokenAddress(&_StakePortal.CallOpts)
}

// Erc20TokenAddress is a free data retrieval call binding the contract method 0xf835cd3c.
//
// Solidity: function erc20TokenAddress() view returns(address)
func (_StakePortal *StakePortalCallerSession) Erc20TokenAddress() (common.Address, error) {
	return _StakePortal.Contract.Erc20TokenAddress(&_StakePortal.CallOpts)
}

// MinAmount is a free data retrieval call binding the contract method 0x9b2cb5d8.
//
// Solidity: function minAmount() view returns(uint256)
func (_StakePortal *StakePortalCaller) MinAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakePortal.contract.Call(opts, &out, "minAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinAmount is a free data retrieval call binding the contract method 0x9b2cb5d8.
//
// Solidity: function minAmount() view returns(uint256)
func (_StakePortal *StakePortalSession) MinAmount() (*big.Int, error) {
	return _StakePortal.Contract.MinAmount(&_StakePortal.CallOpts)
}

// MinAmount is a free data retrieval call binding the contract method 0x9b2cb5d8.
//
// Solidity: function minAmount() view returns(uint256)
func (_StakePortal *StakePortalCallerSession) MinAmount() (*big.Int, error) {
	return _StakePortal.Contract.MinAmount(&_StakePortal.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakePortal *StakePortalCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakePortal.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakePortal *StakePortalSession) Owner() (common.Address, error) {
	return _StakePortal.Contract.Owner(&_StakePortal.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakePortal *StakePortalCallerSession) Owner() (common.Address, error) {
	return _StakePortal.Contract.Owner(&_StakePortal.CallOpts)
}

// RelayFee is a free data retrieval call binding the contract method 0x71d30863.
//
// Solidity: function relayFee() view returns(uint256)
func (_StakePortal *StakePortalCaller) RelayFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakePortal.contract.Call(opts, &out, "relayFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RelayFee is a free data retrieval call binding the contract method 0x71d30863.
//
// Solidity: function relayFee() view returns(uint256)
func (_StakePortal *StakePortalSession) RelayFee() (*big.Int, error) {
	return _StakePortal.Contract.RelayFee(&_StakePortal.CallOpts)
}

// RelayFee is a free data retrieval call binding the contract method 0x71d30863.
//
// Solidity: function relayFee() view returns(uint256)
func (_StakePortal *StakePortalCallerSession) RelayFee() (*big.Int, error) {
	return _StakePortal.Contract.RelayFee(&_StakePortal.CallOpts)
}

// StakePoolAddressExist is a free data retrieval call binding the contract method 0x3a68e4b2.
//
// Solidity: function stakePoolAddressExist(address ) view returns(bool)
func (_StakePortal *StakePortalCaller) StakePoolAddressExist(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _StakePortal.contract.Call(opts, &out, "stakePoolAddressExist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StakePoolAddressExist is a free data retrieval call binding the contract method 0x3a68e4b2.
//
// Solidity: function stakePoolAddressExist(address ) view returns(bool)
func (_StakePortal *StakePortalSession) StakePoolAddressExist(arg0 common.Address) (bool, error) {
	return _StakePortal.Contract.StakePoolAddressExist(&_StakePortal.CallOpts, arg0)
}

// StakePoolAddressExist is a free data retrieval call binding the contract method 0x3a68e4b2.
//
// Solidity: function stakePoolAddressExist(address ) view returns(bool)
func (_StakePortal *StakePortalCallerSession) StakePoolAddressExist(arg0 common.Address) (bool, error) {
	return _StakePortal.Contract.StakePoolAddressExist(&_StakePortal.CallOpts, arg0)
}

// StakeSwitch is a free data retrieval call binding the contract method 0x24387cdc.
//
// Solidity: function stakeSwitch() view returns(bool)
func (_StakePortal *StakePortalCaller) StakeSwitch(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StakePortal.contract.Call(opts, &out, "stakeSwitch")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StakeSwitch is a free data retrieval call binding the contract method 0x24387cdc.
//
// Solidity: function stakeSwitch() view returns(bool)
func (_StakePortal *StakePortalSession) StakeSwitch() (bool, error) {
	return _StakePortal.Contract.StakeSwitch(&_StakePortal.CallOpts)
}

// StakeSwitch is a free data retrieval call binding the contract method 0x24387cdc.
//
// Solidity: function stakeSwitch() view returns(bool)
func (_StakePortal *StakePortalCallerSession) StakeSwitch() (bool, error) {
	return _StakePortal.Contract.StakeSwitch(&_StakePortal.CallOpts)
}

// AddChainId is a paid mutator transaction binding the contract method 0x3fa2569b.
//
// Solidity: function addChainId(uint8[] _chaindIdList) returns()
func (_StakePortal *StakePortalTransactor) AddChainId(opts *bind.TransactOpts, _chaindIdList []uint8) (*types.Transaction, error) {
	return _StakePortal.contract.Transact(opts, "addChainId", _chaindIdList)
}

// AddChainId is a paid mutator transaction binding the contract method 0x3fa2569b.
//
// Solidity: function addChainId(uint8[] _chaindIdList) returns()
func (_StakePortal *StakePortalSession) AddChainId(_chaindIdList []uint8) (*types.Transaction, error) {
	return _StakePortal.Contract.AddChainId(&_StakePortal.TransactOpts, _chaindIdList)
}

// AddChainId is a paid mutator transaction binding the contract method 0x3fa2569b.
//
// Solidity: function addChainId(uint8[] _chaindIdList) returns()
func (_StakePortal *StakePortalTransactorSession) AddChainId(_chaindIdList []uint8) (*types.Transaction, error) {
	return _StakePortal.Contract.AddChainId(&_StakePortal.TransactOpts, _chaindIdList)
}

// AddStakePool is a paid mutator transaction binding the contract method 0xcb08cbae.
//
// Solidity: function addStakePool(address[] _stakePoolAddressList) returns()
func (_StakePortal *StakePortalTransactor) AddStakePool(opts *bind.TransactOpts, _stakePoolAddressList []common.Address) (*types.Transaction, error) {
	return _StakePortal.contract.Transact(opts, "addStakePool", _stakePoolAddressList)
}

// AddStakePool is a paid mutator transaction binding the contract method 0xcb08cbae.
//
// Solidity: function addStakePool(address[] _stakePoolAddressList) returns()
func (_StakePortal *StakePortalSession) AddStakePool(_stakePoolAddressList []common.Address) (*types.Transaction, error) {
	return _StakePortal.Contract.AddStakePool(&_StakePortal.TransactOpts, _stakePoolAddressList)
}

// AddStakePool is a paid mutator transaction binding the contract method 0xcb08cbae.
//
// Solidity: function addStakePool(address[] _stakePoolAddressList) returns()
func (_StakePortal *StakePortalTransactorSession) AddStakePool(_stakePoolAddressList []common.Address) (*types.Transaction, error) {
	return _StakePortal.Contract.AddStakePool(&_StakePortal.TransactOpts, _stakePoolAddressList)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x919a2f83.
//
// Solidity: function recoverStake(bytes32 _txHash, bytes32 _stafiRecipient) returns()
func (_StakePortal *StakePortalTransactor) RecoverStake(opts *bind.TransactOpts, _txHash [32]byte, _stafiRecipient [32]byte) (*types.Transaction, error) {
	return _StakePortal.contract.Transact(opts, "recoverStake", _txHash, _stafiRecipient)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x919a2f83.
//
// Solidity: function recoverStake(bytes32 _txHash, bytes32 _stafiRecipient) returns()
func (_StakePortal *StakePortalSession) RecoverStake(_txHash [32]byte, _stafiRecipient [32]byte) (*types.Transaction, error) {
	return _StakePortal.Contract.RecoverStake(&_StakePortal.TransactOpts, _txHash, _stafiRecipient)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x919a2f83.
//
// Solidity: function recoverStake(bytes32 _txHash, bytes32 _stafiRecipient) returns()
func (_StakePortal *StakePortalTransactorSession) RecoverStake(_txHash [32]byte, _stafiRecipient [32]byte) (*types.Transaction, error) {
	return _StakePortal.Contract.RecoverStake(&_StakePortal.TransactOpts, _txHash, _stafiRecipient)
}

// RmChainId is a paid mutator transaction binding the contract method 0xf8a9cf65.
//
// Solidity: function rmChainId(uint8 _chaindId) returns()
func (_StakePortal *StakePortalTransactor) RmChainId(opts *bind.TransactOpts, _chaindId uint8) (*types.Transaction, error) {
	return _StakePortal.contract.Transact(opts, "rmChainId", _chaindId)
}

// RmChainId is a paid mutator transaction binding the contract method 0xf8a9cf65.
//
// Solidity: function rmChainId(uint8 _chaindId) returns()
func (_StakePortal *StakePortalSession) RmChainId(_chaindId uint8) (*types.Transaction, error) {
	return _StakePortal.Contract.RmChainId(&_StakePortal.TransactOpts, _chaindId)
}

// RmChainId is a paid mutator transaction binding the contract method 0xf8a9cf65.
//
// Solidity: function rmChainId(uint8 _chaindId) returns()
func (_StakePortal *StakePortalTransactorSession) RmChainId(_chaindId uint8) (*types.Transaction, error) {
	return _StakePortal.Contract.RmChainId(&_StakePortal.TransactOpts, _chaindId)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _stakePoolAddress) returns()
func (_StakePortal *StakePortalTransactor) RmStakePool(opts *bind.TransactOpts, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakePortal.contract.Transact(opts, "rmStakePool", _stakePoolAddress)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _stakePoolAddress) returns()
func (_StakePortal *StakePortalSession) RmStakePool(_stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakePortal.Contract.RmStakePool(&_StakePortal.TransactOpts, _stakePoolAddress)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _stakePoolAddress) returns()
func (_StakePortal *StakePortalTransactorSession) RmStakePool(_stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakePortal.Contract.RmStakePool(&_StakePortal.TransactOpts, _stakePoolAddress)
}

// SetMinAmount is a paid mutator transaction binding the contract method 0x897b0637.
//
// Solidity: function setMinAmount(uint256 _minAmount) returns()
func (_StakePortal *StakePortalTransactor) SetMinAmount(opts *bind.TransactOpts, _minAmount *big.Int) (*types.Transaction, error) {
	return _StakePortal.contract.Transact(opts, "setMinAmount", _minAmount)
}

// SetMinAmount is a paid mutator transaction binding the contract method 0x897b0637.
//
// Solidity: function setMinAmount(uint256 _minAmount) returns()
func (_StakePortal *StakePortalSession) SetMinAmount(_minAmount *big.Int) (*types.Transaction, error) {
	return _StakePortal.Contract.SetMinAmount(&_StakePortal.TransactOpts, _minAmount)
}

// SetMinAmount is a paid mutator transaction binding the contract method 0x897b0637.
//
// Solidity: function setMinAmount(uint256 _minAmount) returns()
func (_StakePortal *StakePortalTransactorSession) SetMinAmount(_minAmount *big.Int) (*types.Transaction, error) {
	return _StakePortal.Contract.SetMinAmount(&_StakePortal.TransactOpts, _minAmount)
}

// SetRelayFee is a paid mutator transaction binding the contract method 0x98385109.
//
// Solidity: function setRelayFee(uint256 _relayFee) returns()
func (_StakePortal *StakePortalTransactor) SetRelayFee(opts *bind.TransactOpts, _relayFee *big.Int) (*types.Transaction, error) {
	return _StakePortal.contract.Transact(opts, "setRelayFee", _relayFee)
}

// SetRelayFee is a paid mutator transaction binding the contract method 0x98385109.
//
// Solidity: function setRelayFee(uint256 _relayFee) returns()
func (_StakePortal *StakePortalSession) SetRelayFee(_relayFee *big.Int) (*types.Transaction, error) {
	return _StakePortal.Contract.SetRelayFee(&_StakePortal.TransactOpts, _relayFee)
}

// SetRelayFee is a paid mutator transaction binding the contract method 0x98385109.
//
// Solidity: function setRelayFee(uint256 _relayFee) returns()
func (_StakePortal *StakePortalTransactorSession) SetRelayFee(_relayFee *big.Int) (*types.Transaction, error) {
	return _StakePortal.Contract.SetRelayFee(&_StakePortal.TransactOpts, _relayFee)
}

// Stake is a paid mutator transaction binding the contract method 0x0f52a4ec.
//
// Solidity: function stake(address _stakePoolAddress, uint256 _amount, uint8 _destChainId, bytes32 _stafiRecipient, address _destRecipient) payable returns()
func (_StakePortal *StakePortalTransactor) Stake(opts *bind.TransactOpts, _stakePoolAddress common.Address, _amount *big.Int, _destChainId uint8, _stafiRecipient [32]byte, _destRecipient common.Address) (*types.Transaction, error) {
	return _StakePortal.contract.Transact(opts, "stake", _stakePoolAddress, _amount, _destChainId, _stafiRecipient, _destRecipient)
}

// Stake is a paid mutator transaction binding the contract method 0x0f52a4ec.
//
// Solidity: function stake(address _stakePoolAddress, uint256 _amount, uint8 _destChainId, bytes32 _stafiRecipient, address _destRecipient) payable returns()
func (_StakePortal *StakePortalSession) Stake(_stakePoolAddress common.Address, _amount *big.Int, _destChainId uint8, _stafiRecipient [32]byte, _destRecipient common.Address) (*types.Transaction, error) {
	return _StakePortal.Contract.Stake(&_StakePortal.TransactOpts, _stakePoolAddress, _amount, _destChainId, _stafiRecipient, _destRecipient)
}

// Stake is a paid mutator transaction binding the contract method 0x0f52a4ec.
//
// Solidity: function stake(address _stakePoolAddress, uint256 _amount, uint8 _destChainId, bytes32 _stafiRecipient, address _destRecipient) payable returns()
func (_StakePortal *StakePortalTransactorSession) Stake(_stakePoolAddress common.Address, _amount *big.Int, _destChainId uint8, _stafiRecipient [32]byte, _destRecipient common.Address) (*types.Transaction, error) {
	return _StakePortal.Contract.Stake(&_StakePortal.TransactOpts, _stakePoolAddress, _amount, _destChainId, _stafiRecipient, _destRecipient)
}

// ToggleSwitch is a paid mutator transaction binding the contract method 0xbfe4d261.
//
// Solidity: function toggleSwitch() returns()
func (_StakePortal *StakePortalTransactor) ToggleSwitch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakePortal.contract.Transact(opts, "toggleSwitch")
}

// ToggleSwitch is a paid mutator transaction binding the contract method 0xbfe4d261.
//
// Solidity: function toggleSwitch() returns()
func (_StakePortal *StakePortalSession) ToggleSwitch() (*types.Transaction, error) {
	return _StakePortal.Contract.ToggleSwitch(&_StakePortal.TransactOpts)
}

// ToggleSwitch is a paid mutator transaction binding the contract method 0xbfe4d261.
//
// Solidity: function toggleSwitch() returns()
func (_StakePortal *StakePortalTransactorSession) ToggleSwitch() (*types.Transaction, error) {
	return _StakePortal.Contract.ToggleSwitch(&_StakePortal.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakePortal *StakePortalTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _StakePortal.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakePortal *StakePortalSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _StakePortal.Contract.TransferOwnership(&_StakePortal.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakePortal *StakePortalTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _StakePortal.Contract.TransferOwnership(&_StakePortal.TransactOpts, _newOwner)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns()
func (_StakePortal *StakePortalTransactor) WithdrawFee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakePortal.contract.Transact(opts, "withdrawFee")
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns()
func (_StakePortal *StakePortalSession) WithdrawFee() (*types.Transaction, error) {
	return _StakePortal.Contract.WithdrawFee(&_StakePortal.TransactOpts)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns()
func (_StakePortal *StakePortalTransactorSession) WithdrawFee() (*types.Transaction, error) {
	return _StakePortal.Contract.WithdrawFee(&_StakePortal.TransactOpts)
}

// StakePortalRecoverStakeIterator is returned from FilterRecoverStake and is used to iterate over the raw logs and unpacked data for RecoverStake events raised by the StakePortal contract.
type StakePortalRecoverStakeIterator struct {
	Event *StakePortalRecoverStake // Event containing the contract specifics and raw log

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
func (it *StakePortalRecoverStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakePortalRecoverStake)
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
		it.Event = new(StakePortalRecoverStake)
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
func (it *StakePortalRecoverStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakePortalRecoverStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakePortalRecoverStake represents a RecoverStake event raised by the StakePortal contract.
type StakePortalRecoverStake struct {
	TxHash         [32]byte
	StafiRecipient [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRecoverStake is a free log retrieval operation binding the contract event 0xe6824256d477a01bdb00762b9b44dd3e583772d99e373d00d1722223336586c7.
//
// Solidity: event RecoverStake(bytes32 txHash, bytes32 stafiRecipient)
func (_StakePortal *StakePortalFilterer) FilterRecoverStake(opts *bind.FilterOpts) (*StakePortalRecoverStakeIterator, error) {

	logs, sub, err := _StakePortal.contract.FilterLogs(opts, "RecoverStake")
	if err != nil {
		return nil, err
	}
	return &StakePortalRecoverStakeIterator{contract: _StakePortal.contract, event: "RecoverStake", logs: logs, sub: sub}, nil
}

// WatchRecoverStake is a free log subscription operation binding the contract event 0xe6824256d477a01bdb00762b9b44dd3e583772d99e373d00d1722223336586c7.
//
// Solidity: event RecoverStake(bytes32 txHash, bytes32 stafiRecipient)
func (_StakePortal *StakePortalFilterer) WatchRecoverStake(opts *bind.WatchOpts, sink chan<- *StakePortalRecoverStake) (event.Subscription, error) {

	logs, sub, err := _StakePortal.contract.WatchLogs(opts, "RecoverStake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakePortalRecoverStake)
				if err := _StakePortal.contract.UnpackLog(event, "RecoverStake", log); err != nil {
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

// ParseRecoverStake is a log parse operation binding the contract event 0xe6824256d477a01bdb00762b9b44dd3e583772d99e373d00d1722223336586c7.
//
// Solidity: event RecoverStake(bytes32 txHash, bytes32 stafiRecipient)
func (_StakePortal *StakePortalFilterer) ParseRecoverStake(log types.Log) (*StakePortalRecoverStake, error) {
	event := new(StakePortalRecoverStake)
	if err := _StakePortal.contract.UnpackLog(event, "RecoverStake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakePortalStakeIterator is returned from FilterStake and is used to iterate over the raw logs and unpacked data for Stake events raised by the StakePortal contract.
type StakePortalStakeIterator struct {
	Event *StakePortalStake // Event containing the contract specifics and raw log

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
func (it *StakePortalStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakePortalStake)
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
		it.Event = new(StakePortalStake)
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
func (it *StakePortalStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakePortalStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakePortalStake represents a Stake event raised by the StakePortal contract.
type StakePortalStake struct {
	Staker         common.Address
	StakePool      common.Address
	Amount         *big.Int
	ChainId        uint8
	StafiRecipient [32]byte
	DestRecipient  common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterStake is a free log retrieval operation binding the contract event 0xa3442c24a17de01319e72fe10476c331914acbf277e5e27004c94169bb16e883.
//
// Solidity: event Stake(address staker, address stakePool, uint256 amount, uint8 chainId, bytes32 stafiRecipient, address destRecipient)
func (_StakePortal *StakePortalFilterer) FilterStake(opts *bind.FilterOpts) (*StakePortalStakeIterator, error) {

	logs, sub, err := _StakePortal.contract.FilterLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return &StakePortalStakeIterator{contract: _StakePortal.contract, event: "Stake", logs: logs, sub: sub}, nil
}

// WatchStake is a free log subscription operation binding the contract event 0xa3442c24a17de01319e72fe10476c331914acbf277e5e27004c94169bb16e883.
//
// Solidity: event Stake(address staker, address stakePool, uint256 amount, uint8 chainId, bytes32 stafiRecipient, address destRecipient)
func (_StakePortal *StakePortalFilterer) WatchStake(opts *bind.WatchOpts, sink chan<- *StakePortalStake) (event.Subscription, error) {

	logs, sub, err := _StakePortal.contract.WatchLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakePortalStake)
				if err := _StakePortal.contract.UnpackLog(event, "Stake", log); err != nil {
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

// ParseStake is a log parse operation binding the contract event 0xa3442c24a17de01319e72fe10476c331914acbf277e5e27004c94169bb16e883.
//
// Solidity: event Stake(address staker, address stakePool, uint256 amount, uint8 chainId, bytes32 stafiRecipient, address destRecipient)
func (_StakePortal *StakePortalFilterer) ParseStake(log types.Log) (*StakePortalStake, error) {
	event := new(StakePortalStake)
	if err := _StakePortal.contract.UnpackLog(event, "Stake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
