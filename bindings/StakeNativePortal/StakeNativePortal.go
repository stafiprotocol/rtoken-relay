// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stake_native_portal

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

// StakeNativePortalMetaData contains all meta data concerning the StakeNativePortal contract.
var StakeNativePortalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_stakePoolAddressList\",\"type\":\"address[]\"},{\"internalType\":\"uint8[]\",\"name\":\"_chainIdList\",\"type\":\"uint8[]\"},{\"internalType\":\"uint256\",\"name\":\"_minAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_relayFee\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stafiRecipient\",\"type\":\"bytes32\"}],\"name\":\"RecoverStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stakePool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"chainId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stafiRecipient\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destRecipient\",\"type\":\"address\"}],\"name\":\"Stake\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint8[]\",\"name\":\"_chaindIdList\",\"type\":\"uint8[]\"}],\"name\":\"addChainId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_stakePoolAddressList\",\"type\":\"address[]\"}],\"name\":\"addStakePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"bridgeFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"chainIdExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_txHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_stafiRecipient\",\"type\":\"bytes32\"}],\"name\":\"recoverStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"relayFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_chaindId\",\"type\":\"uint8\"}],\"name\":\"rmChainId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"}],\"name\":\"rmStakePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_chainId\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_bridgeFee\",\"type\":\"uint256\"}],\"name\":\"setBridgeFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minAmount\",\"type\":\"uint256\"}],\"name\":\"setMinAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_relayFee\",\"type\":\"uint256\"}],\"name\":\"setRelayFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"_destChainId\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_stafiRecipient\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_destRecipient\",\"type\":\"address\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"stakePoolAddressExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeSwitch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"toggleSwitch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StakeNativePortalABI is the input ABI used to generate the binding from.
// Deprecated: Use StakeNativePortalMetaData.ABI instead.
var StakeNativePortalABI = StakeNativePortalMetaData.ABI

// StakeNativePortal is an auto generated Go binding around an Ethereum contract.
type StakeNativePortal struct {
	StakeNativePortalCaller     // Read-only binding to the contract
	StakeNativePortalTransactor // Write-only binding to the contract
	StakeNativePortalFilterer   // Log filterer for contract events
}

// StakeNativePortalCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakeNativePortalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeNativePortalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakeNativePortalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeNativePortalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakeNativePortalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeNativePortalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakeNativePortalSession struct {
	Contract     *StakeNativePortal // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StakeNativePortalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakeNativePortalCallerSession struct {
	Contract *StakeNativePortalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// StakeNativePortalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakeNativePortalTransactorSession struct {
	Contract     *StakeNativePortalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// StakeNativePortalRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakeNativePortalRaw struct {
	Contract *StakeNativePortal // Generic contract binding to access the raw methods on
}

// StakeNativePortalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakeNativePortalCallerRaw struct {
	Contract *StakeNativePortalCaller // Generic read-only contract binding to access the raw methods on
}

// StakeNativePortalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakeNativePortalTransactorRaw struct {
	Contract *StakeNativePortalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakeNativePortal creates a new instance of StakeNativePortal, bound to a specific deployed contract.
func NewStakeNativePortal(address common.Address, backend bind.ContractBackend) (*StakeNativePortal, error) {
	contract, err := bindStakeNativePortal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakeNativePortal{StakeNativePortalCaller: StakeNativePortalCaller{contract: contract}, StakeNativePortalTransactor: StakeNativePortalTransactor{contract: contract}, StakeNativePortalFilterer: StakeNativePortalFilterer{contract: contract}}, nil
}

// NewStakeNativePortalCaller creates a new read-only instance of StakeNativePortal, bound to a specific deployed contract.
func NewStakeNativePortalCaller(address common.Address, caller bind.ContractCaller) (*StakeNativePortalCaller, error) {
	contract, err := bindStakeNativePortal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakeNativePortalCaller{contract: contract}, nil
}

// NewStakeNativePortalTransactor creates a new write-only instance of StakeNativePortal, bound to a specific deployed contract.
func NewStakeNativePortalTransactor(address common.Address, transactor bind.ContractTransactor) (*StakeNativePortalTransactor, error) {
	contract, err := bindStakeNativePortal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakeNativePortalTransactor{contract: contract}, nil
}

// NewStakeNativePortalFilterer creates a new log filterer instance of StakeNativePortal, bound to a specific deployed contract.
func NewStakeNativePortalFilterer(address common.Address, filterer bind.ContractFilterer) (*StakeNativePortalFilterer, error) {
	contract, err := bindStakeNativePortal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakeNativePortalFilterer{contract: contract}, nil
}

// bindStakeNativePortal binds a generic wrapper to an already deployed contract.
func bindStakeNativePortal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakeNativePortalABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakeNativePortal *StakeNativePortalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakeNativePortal.Contract.StakeNativePortalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakeNativePortal *StakeNativePortalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.StakeNativePortalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakeNativePortal *StakeNativePortalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.StakeNativePortalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakeNativePortal *StakeNativePortalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakeNativePortal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakeNativePortal *StakeNativePortalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakeNativePortal *StakeNativePortalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.contract.Transact(opts, method, params...)
}

// BridgeFee is a free data retrieval call binding the contract method 0x046853ec.
//
// Solidity: function bridgeFee(uint8 ) view returns(uint256)
func (_StakeNativePortal *StakeNativePortalCaller) BridgeFee(opts *bind.CallOpts, arg0 uint8) (*big.Int, error) {
	var out []interface{}
	err := _StakeNativePortal.contract.Call(opts, &out, "bridgeFee", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BridgeFee is a free data retrieval call binding the contract method 0x046853ec.
//
// Solidity: function bridgeFee(uint8 ) view returns(uint256)
func (_StakeNativePortal *StakeNativePortalSession) BridgeFee(arg0 uint8) (*big.Int, error) {
	return _StakeNativePortal.Contract.BridgeFee(&_StakeNativePortal.CallOpts, arg0)
}

// BridgeFee is a free data retrieval call binding the contract method 0x046853ec.
//
// Solidity: function bridgeFee(uint8 ) view returns(uint256)
func (_StakeNativePortal *StakeNativePortalCallerSession) BridgeFee(arg0 uint8) (*big.Int, error) {
	return _StakeNativePortal.Contract.BridgeFee(&_StakeNativePortal.CallOpts, arg0)
}

// ChainIdExist is a free data retrieval call binding the contract method 0x0133228d.
//
// Solidity: function chainIdExist(uint8 ) view returns(bool)
func (_StakeNativePortal *StakeNativePortalCaller) ChainIdExist(opts *bind.CallOpts, arg0 uint8) (bool, error) {
	var out []interface{}
	err := _StakeNativePortal.contract.Call(opts, &out, "chainIdExist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ChainIdExist is a free data retrieval call binding the contract method 0x0133228d.
//
// Solidity: function chainIdExist(uint8 ) view returns(bool)
func (_StakeNativePortal *StakeNativePortalSession) ChainIdExist(arg0 uint8) (bool, error) {
	return _StakeNativePortal.Contract.ChainIdExist(&_StakeNativePortal.CallOpts, arg0)
}

// ChainIdExist is a free data retrieval call binding the contract method 0x0133228d.
//
// Solidity: function chainIdExist(uint8 ) view returns(bool)
func (_StakeNativePortal *StakeNativePortalCallerSession) ChainIdExist(arg0 uint8) (bool, error) {
	return _StakeNativePortal.Contract.ChainIdExist(&_StakeNativePortal.CallOpts, arg0)
}

// MinAmount is a free data retrieval call binding the contract method 0x9b2cb5d8.
//
// Solidity: function minAmount() view returns(uint256)
func (_StakeNativePortal *StakeNativePortalCaller) MinAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeNativePortal.contract.Call(opts, &out, "minAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinAmount is a free data retrieval call binding the contract method 0x9b2cb5d8.
//
// Solidity: function minAmount() view returns(uint256)
func (_StakeNativePortal *StakeNativePortalSession) MinAmount() (*big.Int, error) {
	return _StakeNativePortal.Contract.MinAmount(&_StakeNativePortal.CallOpts)
}

// MinAmount is a free data retrieval call binding the contract method 0x9b2cb5d8.
//
// Solidity: function minAmount() view returns(uint256)
func (_StakeNativePortal *StakeNativePortalCallerSession) MinAmount() (*big.Int, error) {
	return _StakeNativePortal.Contract.MinAmount(&_StakeNativePortal.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakeNativePortal *StakeNativePortalCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeNativePortal.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakeNativePortal *StakeNativePortalSession) Owner() (common.Address, error) {
	return _StakeNativePortal.Contract.Owner(&_StakeNativePortal.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakeNativePortal *StakeNativePortalCallerSession) Owner() (common.Address, error) {
	return _StakeNativePortal.Contract.Owner(&_StakeNativePortal.CallOpts)
}

// RelayFee is a free data retrieval call binding the contract method 0x71d30863.
//
// Solidity: function relayFee() view returns(uint256)
func (_StakeNativePortal *StakeNativePortalCaller) RelayFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeNativePortal.contract.Call(opts, &out, "relayFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RelayFee is a free data retrieval call binding the contract method 0x71d30863.
//
// Solidity: function relayFee() view returns(uint256)
func (_StakeNativePortal *StakeNativePortalSession) RelayFee() (*big.Int, error) {
	return _StakeNativePortal.Contract.RelayFee(&_StakeNativePortal.CallOpts)
}

// RelayFee is a free data retrieval call binding the contract method 0x71d30863.
//
// Solidity: function relayFee() view returns(uint256)
func (_StakeNativePortal *StakeNativePortalCallerSession) RelayFee() (*big.Int, error) {
	return _StakeNativePortal.Contract.RelayFee(&_StakeNativePortal.CallOpts)
}

// StakePoolAddressExist is a free data retrieval call binding the contract method 0x3a68e4b2.
//
// Solidity: function stakePoolAddressExist(address ) view returns(bool)
func (_StakeNativePortal *StakeNativePortalCaller) StakePoolAddressExist(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _StakeNativePortal.contract.Call(opts, &out, "stakePoolAddressExist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StakePoolAddressExist is a free data retrieval call binding the contract method 0x3a68e4b2.
//
// Solidity: function stakePoolAddressExist(address ) view returns(bool)
func (_StakeNativePortal *StakeNativePortalSession) StakePoolAddressExist(arg0 common.Address) (bool, error) {
	return _StakeNativePortal.Contract.StakePoolAddressExist(&_StakeNativePortal.CallOpts, arg0)
}

// StakePoolAddressExist is a free data retrieval call binding the contract method 0x3a68e4b2.
//
// Solidity: function stakePoolAddressExist(address ) view returns(bool)
func (_StakeNativePortal *StakeNativePortalCallerSession) StakePoolAddressExist(arg0 common.Address) (bool, error) {
	return _StakeNativePortal.Contract.StakePoolAddressExist(&_StakeNativePortal.CallOpts, arg0)
}

// StakeSwitch is a free data retrieval call binding the contract method 0x24387cdc.
//
// Solidity: function stakeSwitch() view returns(bool)
func (_StakeNativePortal *StakeNativePortalCaller) StakeSwitch(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StakeNativePortal.contract.Call(opts, &out, "stakeSwitch")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StakeSwitch is a free data retrieval call binding the contract method 0x24387cdc.
//
// Solidity: function stakeSwitch() view returns(bool)
func (_StakeNativePortal *StakeNativePortalSession) StakeSwitch() (bool, error) {
	return _StakeNativePortal.Contract.StakeSwitch(&_StakeNativePortal.CallOpts)
}

// StakeSwitch is a free data retrieval call binding the contract method 0x24387cdc.
//
// Solidity: function stakeSwitch() view returns(bool)
func (_StakeNativePortal *StakeNativePortalCallerSession) StakeSwitch() (bool, error) {
	return _StakeNativePortal.Contract.StakeSwitch(&_StakeNativePortal.CallOpts)
}

// AddChainId is a paid mutator transaction binding the contract method 0x3fa2569b.
//
// Solidity: function addChainId(uint8[] _chaindIdList) returns()
func (_StakeNativePortal *StakeNativePortalTransactor) AddChainId(opts *bind.TransactOpts, _chaindIdList []uint8) (*types.Transaction, error) {
	return _StakeNativePortal.contract.Transact(opts, "addChainId", _chaindIdList)
}

// AddChainId is a paid mutator transaction binding the contract method 0x3fa2569b.
//
// Solidity: function addChainId(uint8[] _chaindIdList) returns()
func (_StakeNativePortal *StakeNativePortalSession) AddChainId(_chaindIdList []uint8) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.AddChainId(&_StakeNativePortal.TransactOpts, _chaindIdList)
}

// AddChainId is a paid mutator transaction binding the contract method 0x3fa2569b.
//
// Solidity: function addChainId(uint8[] _chaindIdList) returns()
func (_StakeNativePortal *StakeNativePortalTransactorSession) AddChainId(_chaindIdList []uint8) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.AddChainId(&_StakeNativePortal.TransactOpts, _chaindIdList)
}

// AddStakePool is a paid mutator transaction binding the contract method 0xcb08cbae.
//
// Solidity: function addStakePool(address[] _stakePoolAddressList) returns()
func (_StakeNativePortal *StakeNativePortalTransactor) AddStakePool(opts *bind.TransactOpts, _stakePoolAddressList []common.Address) (*types.Transaction, error) {
	return _StakeNativePortal.contract.Transact(opts, "addStakePool", _stakePoolAddressList)
}

// AddStakePool is a paid mutator transaction binding the contract method 0xcb08cbae.
//
// Solidity: function addStakePool(address[] _stakePoolAddressList) returns()
func (_StakeNativePortal *StakeNativePortalSession) AddStakePool(_stakePoolAddressList []common.Address) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.AddStakePool(&_StakeNativePortal.TransactOpts, _stakePoolAddressList)
}

// AddStakePool is a paid mutator transaction binding the contract method 0xcb08cbae.
//
// Solidity: function addStakePool(address[] _stakePoolAddressList) returns()
func (_StakeNativePortal *StakeNativePortalTransactorSession) AddStakePool(_stakePoolAddressList []common.Address) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.AddStakePool(&_StakeNativePortal.TransactOpts, _stakePoolAddressList)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x919a2f83.
//
// Solidity: function recoverStake(bytes32 _txHash, bytes32 _stafiRecipient) returns()
func (_StakeNativePortal *StakeNativePortalTransactor) RecoverStake(opts *bind.TransactOpts, _txHash [32]byte, _stafiRecipient [32]byte) (*types.Transaction, error) {
	return _StakeNativePortal.contract.Transact(opts, "recoverStake", _txHash, _stafiRecipient)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x919a2f83.
//
// Solidity: function recoverStake(bytes32 _txHash, bytes32 _stafiRecipient) returns()
func (_StakeNativePortal *StakeNativePortalSession) RecoverStake(_txHash [32]byte, _stafiRecipient [32]byte) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.RecoverStake(&_StakeNativePortal.TransactOpts, _txHash, _stafiRecipient)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x919a2f83.
//
// Solidity: function recoverStake(bytes32 _txHash, bytes32 _stafiRecipient) returns()
func (_StakeNativePortal *StakeNativePortalTransactorSession) RecoverStake(_txHash [32]byte, _stafiRecipient [32]byte) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.RecoverStake(&_StakeNativePortal.TransactOpts, _txHash, _stafiRecipient)
}

// RmChainId is a paid mutator transaction binding the contract method 0xf8a9cf65.
//
// Solidity: function rmChainId(uint8 _chaindId) returns()
func (_StakeNativePortal *StakeNativePortalTransactor) RmChainId(opts *bind.TransactOpts, _chaindId uint8) (*types.Transaction, error) {
	return _StakeNativePortal.contract.Transact(opts, "rmChainId", _chaindId)
}

// RmChainId is a paid mutator transaction binding the contract method 0xf8a9cf65.
//
// Solidity: function rmChainId(uint8 _chaindId) returns()
func (_StakeNativePortal *StakeNativePortalSession) RmChainId(_chaindId uint8) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.RmChainId(&_StakeNativePortal.TransactOpts, _chaindId)
}

// RmChainId is a paid mutator transaction binding the contract method 0xf8a9cf65.
//
// Solidity: function rmChainId(uint8 _chaindId) returns()
func (_StakeNativePortal *StakeNativePortalTransactorSession) RmChainId(_chaindId uint8) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.RmChainId(&_StakeNativePortal.TransactOpts, _chaindId)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _stakePoolAddress) returns()
func (_StakeNativePortal *StakeNativePortalTransactor) RmStakePool(opts *bind.TransactOpts, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeNativePortal.contract.Transact(opts, "rmStakePool", _stakePoolAddress)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _stakePoolAddress) returns()
func (_StakeNativePortal *StakeNativePortalSession) RmStakePool(_stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.RmStakePool(&_StakeNativePortal.TransactOpts, _stakePoolAddress)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _stakePoolAddress) returns()
func (_StakeNativePortal *StakeNativePortalTransactorSession) RmStakePool(_stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.RmStakePool(&_StakeNativePortal.TransactOpts, _stakePoolAddress)
}

// SetBridgeFee is a paid mutator transaction binding the contract method 0xaddc3519.
//
// Solidity: function setBridgeFee(uint8 _chainId, uint256 _bridgeFee) returns()
func (_StakeNativePortal *StakeNativePortalTransactor) SetBridgeFee(opts *bind.TransactOpts, _chainId uint8, _bridgeFee *big.Int) (*types.Transaction, error) {
	return _StakeNativePortal.contract.Transact(opts, "setBridgeFee", _chainId, _bridgeFee)
}

// SetBridgeFee is a paid mutator transaction binding the contract method 0xaddc3519.
//
// Solidity: function setBridgeFee(uint8 _chainId, uint256 _bridgeFee) returns()
func (_StakeNativePortal *StakeNativePortalSession) SetBridgeFee(_chainId uint8, _bridgeFee *big.Int) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.SetBridgeFee(&_StakeNativePortal.TransactOpts, _chainId, _bridgeFee)
}

// SetBridgeFee is a paid mutator transaction binding the contract method 0xaddc3519.
//
// Solidity: function setBridgeFee(uint8 _chainId, uint256 _bridgeFee) returns()
func (_StakeNativePortal *StakeNativePortalTransactorSession) SetBridgeFee(_chainId uint8, _bridgeFee *big.Int) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.SetBridgeFee(&_StakeNativePortal.TransactOpts, _chainId, _bridgeFee)
}

// SetMinAmount is a paid mutator transaction binding the contract method 0x897b0637.
//
// Solidity: function setMinAmount(uint256 _minAmount) returns()
func (_StakeNativePortal *StakeNativePortalTransactor) SetMinAmount(opts *bind.TransactOpts, _minAmount *big.Int) (*types.Transaction, error) {
	return _StakeNativePortal.contract.Transact(opts, "setMinAmount", _minAmount)
}

// SetMinAmount is a paid mutator transaction binding the contract method 0x897b0637.
//
// Solidity: function setMinAmount(uint256 _minAmount) returns()
func (_StakeNativePortal *StakeNativePortalSession) SetMinAmount(_minAmount *big.Int) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.SetMinAmount(&_StakeNativePortal.TransactOpts, _minAmount)
}

// SetMinAmount is a paid mutator transaction binding the contract method 0x897b0637.
//
// Solidity: function setMinAmount(uint256 _minAmount) returns()
func (_StakeNativePortal *StakeNativePortalTransactorSession) SetMinAmount(_minAmount *big.Int) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.SetMinAmount(&_StakeNativePortal.TransactOpts, _minAmount)
}

// SetRelayFee is a paid mutator transaction binding the contract method 0x98385109.
//
// Solidity: function setRelayFee(uint256 _relayFee) returns()
func (_StakeNativePortal *StakeNativePortalTransactor) SetRelayFee(opts *bind.TransactOpts, _relayFee *big.Int) (*types.Transaction, error) {
	return _StakeNativePortal.contract.Transact(opts, "setRelayFee", _relayFee)
}

// SetRelayFee is a paid mutator transaction binding the contract method 0x98385109.
//
// Solidity: function setRelayFee(uint256 _relayFee) returns()
func (_StakeNativePortal *StakeNativePortalSession) SetRelayFee(_relayFee *big.Int) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.SetRelayFee(&_StakeNativePortal.TransactOpts, _relayFee)
}

// SetRelayFee is a paid mutator transaction binding the contract method 0x98385109.
//
// Solidity: function setRelayFee(uint256 _relayFee) returns()
func (_StakeNativePortal *StakeNativePortalTransactorSession) SetRelayFee(_relayFee *big.Int) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.SetRelayFee(&_StakeNativePortal.TransactOpts, _relayFee)
}

// Stake is a paid mutator transaction binding the contract method 0x0f52a4ec.
//
// Solidity: function stake(address _stakePoolAddress, uint256 _amount, uint8 _destChainId, bytes32 _stafiRecipient, address _destRecipient) payable returns()
func (_StakeNativePortal *StakeNativePortalTransactor) Stake(opts *bind.TransactOpts, _stakePoolAddress common.Address, _amount *big.Int, _destChainId uint8, _stafiRecipient [32]byte, _destRecipient common.Address) (*types.Transaction, error) {
	return _StakeNativePortal.contract.Transact(opts, "stake", _stakePoolAddress, _amount, _destChainId, _stafiRecipient, _destRecipient)
}

// Stake is a paid mutator transaction binding the contract method 0x0f52a4ec.
//
// Solidity: function stake(address _stakePoolAddress, uint256 _amount, uint8 _destChainId, bytes32 _stafiRecipient, address _destRecipient) payable returns()
func (_StakeNativePortal *StakeNativePortalSession) Stake(_stakePoolAddress common.Address, _amount *big.Int, _destChainId uint8, _stafiRecipient [32]byte, _destRecipient common.Address) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.Stake(&_StakeNativePortal.TransactOpts, _stakePoolAddress, _amount, _destChainId, _stafiRecipient, _destRecipient)
}

// Stake is a paid mutator transaction binding the contract method 0x0f52a4ec.
//
// Solidity: function stake(address _stakePoolAddress, uint256 _amount, uint8 _destChainId, bytes32 _stafiRecipient, address _destRecipient) payable returns()
func (_StakeNativePortal *StakeNativePortalTransactorSession) Stake(_stakePoolAddress common.Address, _amount *big.Int, _destChainId uint8, _stafiRecipient [32]byte, _destRecipient common.Address) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.Stake(&_StakeNativePortal.TransactOpts, _stakePoolAddress, _amount, _destChainId, _stafiRecipient, _destRecipient)
}

// ToggleSwitch is a paid mutator transaction binding the contract method 0xbfe4d261.
//
// Solidity: function toggleSwitch() returns()
func (_StakeNativePortal *StakeNativePortalTransactor) ToggleSwitch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeNativePortal.contract.Transact(opts, "toggleSwitch")
}

// ToggleSwitch is a paid mutator transaction binding the contract method 0xbfe4d261.
//
// Solidity: function toggleSwitch() returns()
func (_StakeNativePortal *StakeNativePortalSession) ToggleSwitch() (*types.Transaction, error) {
	return _StakeNativePortal.Contract.ToggleSwitch(&_StakeNativePortal.TransactOpts)
}

// ToggleSwitch is a paid mutator transaction binding the contract method 0xbfe4d261.
//
// Solidity: function toggleSwitch() returns()
func (_StakeNativePortal *StakeNativePortalTransactorSession) ToggleSwitch() (*types.Transaction, error) {
	return _StakeNativePortal.Contract.ToggleSwitch(&_StakeNativePortal.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakeNativePortal *StakeNativePortalTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _StakeNativePortal.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakeNativePortal *StakeNativePortalSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.TransferOwnership(&_StakeNativePortal.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakeNativePortal *StakeNativePortalTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _StakeNativePortal.Contract.TransferOwnership(&_StakeNativePortal.TransactOpts, _newOwner)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns()
func (_StakeNativePortal *StakeNativePortalTransactor) WithdrawFee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeNativePortal.contract.Transact(opts, "withdrawFee")
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns()
func (_StakeNativePortal *StakeNativePortalSession) WithdrawFee() (*types.Transaction, error) {
	return _StakeNativePortal.Contract.WithdrawFee(&_StakeNativePortal.TransactOpts)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns()
func (_StakeNativePortal *StakeNativePortalTransactorSession) WithdrawFee() (*types.Transaction, error) {
	return _StakeNativePortal.Contract.WithdrawFee(&_StakeNativePortal.TransactOpts)
}

// StakeNativePortalRecoverStakeIterator is returned from FilterRecoverStake and is used to iterate over the raw logs and unpacked data for RecoverStake events raised by the StakeNativePortal contract.
type StakeNativePortalRecoverStakeIterator struct {
	Event *StakeNativePortalRecoverStake // Event containing the contract specifics and raw log

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
func (it *StakeNativePortalRecoverStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeNativePortalRecoverStake)
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
		it.Event = new(StakeNativePortalRecoverStake)
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
func (it *StakeNativePortalRecoverStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeNativePortalRecoverStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeNativePortalRecoverStake represents a RecoverStake event raised by the StakeNativePortal contract.
type StakeNativePortalRecoverStake struct {
	TxHash         [32]byte
	StafiRecipient [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRecoverStake is a free log retrieval operation binding the contract event 0xe6824256d477a01bdb00762b9b44dd3e583772d99e373d00d1722223336586c7.
//
// Solidity: event RecoverStake(bytes32 txHash, bytes32 stafiRecipient)
func (_StakeNativePortal *StakeNativePortalFilterer) FilterRecoverStake(opts *bind.FilterOpts) (*StakeNativePortalRecoverStakeIterator, error) {

	logs, sub, err := _StakeNativePortal.contract.FilterLogs(opts, "RecoverStake")
	if err != nil {
		return nil, err
	}
	return &StakeNativePortalRecoverStakeIterator{contract: _StakeNativePortal.contract, event: "RecoverStake", logs: logs, sub: sub}, nil
}

// WatchRecoverStake is a free log subscription operation binding the contract event 0xe6824256d477a01bdb00762b9b44dd3e583772d99e373d00d1722223336586c7.
//
// Solidity: event RecoverStake(bytes32 txHash, bytes32 stafiRecipient)
func (_StakeNativePortal *StakeNativePortalFilterer) WatchRecoverStake(opts *bind.WatchOpts, sink chan<- *StakeNativePortalRecoverStake) (event.Subscription, error) {

	logs, sub, err := _StakeNativePortal.contract.WatchLogs(opts, "RecoverStake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeNativePortalRecoverStake)
				if err := _StakeNativePortal.contract.UnpackLog(event, "RecoverStake", log); err != nil {
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
func (_StakeNativePortal *StakeNativePortalFilterer) ParseRecoverStake(log types.Log) (*StakeNativePortalRecoverStake, error) {
	event := new(StakeNativePortalRecoverStake)
	if err := _StakeNativePortal.contract.UnpackLog(event, "RecoverStake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeNativePortalStakeIterator is returned from FilterStake and is used to iterate over the raw logs and unpacked data for Stake events raised by the StakeNativePortal contract.
type StakeNativePortalStakeIterator struct {
	Event *StakeNativePortalStake // Event containing the contract specifics and raw log

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
func (it *StakeNativePortalStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeNativePortalStake)
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
		it.Event = new(StakeNativePortalStake)
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
func (it *StakeNativePortalStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeNativePortalStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeNativePortalStake represents a Stake event raised by the StakeNativePortal contract.
type StakeNativePortalStake struct {
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
func (_StakeNativePortal *StakeNativePortalFilterer) FilterStake(opts *bind.FilterOpts) (*StakeNativePortalStakeIterator, error) {

	logs, sub, err := _StakeNativePortal.contract.FilterLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return &StakeNativePortalStakeIterator{contract: _StakeNativePortal.contract, event: "Stake", logs: logs, sub: sub}, nil
}

// WatchStake is a free log subscription operation binding the contract event 0xa3442c24a17de01319e72fe10476c331914acbf277e5e27004c94169bb16e883.
//
// Solidity: event Stake(address staker, address stakePool, uint256 amount, uint8 chainId, bytes32 stafiRecipient, address destRecipient)
func (_StakeNativePortal *StakeNativePortalFilterer) WatchStake(opts *bind.WatchOpts, sink chan<- *StakeNativePortalStake) (event.Subscription, error) {

	logs, sub, err := _StakeNativePortal.contract.WatchLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeNativePortalStake)
				if err := _StakeNativePortal.contract.UnpackLog(event, "Stake", log); err != nil {
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
func (_StakeNativePortal *StakeNativePortalFilterer) ParseStake(log types.Log) (*StakeNativePortalStake, error) {
	event := new(StakeNativePortalStake)
	if err := _StakeNativePortal.contract.UnpackLog(event, "Stake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
