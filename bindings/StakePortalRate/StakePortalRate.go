// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stake_portal_rate

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

// StakePortalRateMetaData contains all meta data concerning the StakePortalRate contract.
var StakePortalRateMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_initialSubAccounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_initialThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"proposalId\",\"type\":\"bytes32\"}],\"name\":\"ProposalExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"SetRate\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"addSubAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newThreshold\",\"type\":\"uint256\"}],\"name\":\"changeThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"getSubAccountIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_proposalId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"hasVoted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"proposals\",\"outputs\":[{\"internalType\":\"enumMultisig.ProposalStatus\",\"name\":\"_status\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"_yesVotes\",\"type\":\"uint16\"},{\"internalType\":\"uint8\",\"name\":\"_yesVotesTotal\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rateChangeLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"removeSubAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"setRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rateChangeLimit\",\"type\":\"uint256\"}],\"name\":\"setRateChangeLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_proposalId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"voteRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StakePortalRateABI is the input ABI used to generate the binding from.
// Deprecated: Use StakePortalRateMetaData.ABI instead.
var StakePortalRateABI = StakePortalRateMetaData.ABI

// StakePortalRate is an auto generated Go binding around an Ethereum contract.
type StakePortalRate struct {
	StakePortalRateCaller     // Read-only binding to the contract
	StakePortalRateTransactor // Write-only binding to the contract
	StakePortalRateFilterer   // Log filterer for contract events
}

// StakePortalRateCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakePortalRateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakePortalRateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakePortalRateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakePortalRateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakePortalRateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakePortalRateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakePortalRateSession struct {
	Contract     *StakePortalRate  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakePortalRateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakePortalRateCallerSession struct {
	Contract *StakePortalRateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// StakePortalRateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakePortalRateTransactorSession struct {
	Contract     *StakePortalRateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// StakePortalRateRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakePortalRateRaw struct {
	Contract *StakePortalRate // Generic contract binding to access the raw methods on
}

// StakePortalRateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakePortalRateCallerRaw struct {
	Contract *StakePortalRateCaller // Generic read-only contract binding to access the raw methods on
}

// StakePortalRateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakePortalRateTransactorRaw struct {
	Contract *StakePortalRateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakePortalRate creates a new instance of StakePortalRate, bound to a specific deployed contract.
func NewStakePortalRate(address common.Address, backend bind.ContractBackend) (*StakePortalRate, error) {
	contract, err := bindStakePortalRate(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakePortalRate{StakePortalRateCaller: StakePortalRateCaller{contract: contract}, StakePortalRateTransactor: StakePortalRateTransactor{contract: contract}, StakePortalRateFilterer: StakePortalRateFilterer{contract: contract}}, nil
}

// NewStakePortalRateCaller creates a new read-only instance of StakePortalRate, bound to a specific deployed contract.
func NewStakePortalRateCaller(address common.Address, caller bind.ContractCaller) (*StakePortalRateCaller, error) {
	contract, err := bindStakePortalRate(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakePortalRateCaller{contract: contract}, nil
}

// NewStakePortalRateTransactor creates a new write-only instance of StakePortalRate, bound to a specific deployed contract.
func NewStakePortalRateTransactor(address common.Address, transactor bind.ContractTransactor) (*StakePortalRateTransactor, error) {
	contract, err := bindStakePortalRate(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakePortalRateTransactor{contract: contract}, nil
}

// NewStakePortalRateFilterer creates a new log filterer instance of StakePortalRate, bound to a specific deployed contract.
func NewStakePortalRateFilterer(address common.Address, filterer bind.ContractFilterer) (*StakePortalRateFilterer, error) {
	contract, err := bindStakePortalRate(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakePortalRateFilterer{contract: contract}, nil
}

// bindStakePortalRate binds a generic wrapper to an already deployed contract.
func bindStakePortalRate(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakePortalRateABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakePortalRate *StakePortalRateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakePortalRate.Contract.StakePortalRateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakePortalRate *StakePortalRateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakePortalRate.Contract.StakePortalRateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakePortalRate *StakePortalRateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakePortalRate.Contract.StakePortalRateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakePortalRate *StakePortalRateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakePortalRate.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakePortalRate *StakePortalRateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakePortalRate.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakePortalRate *StakePortalRateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakePortalRate.Contract.contract.Transact(opts, method, params...)
}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() view returns(uint256)
func (_StakePortalRate *StakePortalRateCaller) GetRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakePortalRate.contract.Call(opts, &out, "getRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() view returns(uint256)
func (_StakePortalRate *StakePortalRateSession) GetRate() (*big.Int, error) {
	return _StakePortalRate.Contract.GetRate(&_StakePortalRate.CallOpts)
}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() view returns(uint256)
func (_StakePortalRate *StakePortalRateCallerSession) GetRate() (*big.Int, error) {
	return _StakePortalRate.Contract.GetRate(&_StakePortalRate.CallOpts)
}

// GetSubAccountIndex is a free data retrieval call binding the contract method 0x763f8680.
//
// Solidity: function getSubAccountIndex(address _subAccount) view returns(uint256)
func (_StakePortalRate *StakePortalRateCaller) GetSubAccountIndex(opts *bind.CallOpts, _subAccount common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakePortalRate.contract.Call(opts, &out, "getSubAccountIndex", _subAccount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSubAccountIndex is a free data retrieval call binding the contract method 0x763f8680.
//
// Solidity: function getSubAccountIndex(address _subAccount) view returns(uint256)
func (_StakePortalRate *StakePortalRateSession) GetSubAccountIndex(_subAccount common.Address) (*big.Int, error) {
	return _StakePortalRate.Contract.GetSubAccountIndex(&_StakePortalRate.CallOpts, _subAccount)
}

// GetSubAccountIndex is a free data retrieval call binding the contract method 0x763f8680.
//
// Solidity: function getSubAccountIndex(address _subAccount) view returns(uint256)
func (_StakePortalRate *StakePortalRateCallerSession) GetSubAccountIndex(_subAccount common.Address) (*big.Int, error) {
	return _StakePortalRate.Contract.GetSubAccountIndex(&_StakePortalRate.CallOpts, _subAccount)
}

// HasVoted is a free data retrieval call binding the contract method 0xaadc3b72.
//
// Solidity: function hasVoted(bytes32 _proposalId, address _subAccount) view returns(bool)
func (_StakePortalRate *StakePortalRateCaller) HasVoted(opts *bind.CallOpts, _proposalId [32]byte, _subAccount common.Address) (bool, error) {
	var out []interface{}
	err := _StakePortalRate.contract.Call(opts, &out, "hasVoted", _proposalId, _subAccount)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasVoted is a free data retrieval call binding the contract method 0xaadc3b72.
//
// Solidity: function hasVoted(bytes32 _proposalId, address _subAccount) view returns(bool)
func (_StakePortalRate *StakePortalRateSession) HasVoted(_proposalId [32]byte, _subAccount common.Address) (bool, error) {
	return _StakePortalRate.Contract.HasVoted(&_StakePortalRate.CallOpts, _proposalId, _subAccount)
}

// HasVoted is a free data retrieval call binding the contract method 0xaadc3b72.
//
// Solidity: function hasVoted(bytes32 _proposalId, address _subAccount) view returns(bool)
func (_StakePortalRate *StakePortalRateCallerSession) HasVoted(_proposalId [32]byte, _subAccount common.Address) (bool, error) {
	return _StakePortalRate.Contract.HasVoted(&_StakePortalRate.CallOpts, _proposalId, _subAccount)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakePortalRate *StakePortalRateCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakePortalRate.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakePortalRate *StakePortalRateSession) Owner() (common.Address, error) {
	return _StakePortalRate.Contract.Owner(&_StakePortalRate.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakePortalRate *StakePortalRateCallerSession) Owner() (common.Address, error) {
	return _StakePortalRate.Contract.Owner(&_StakePortalRate.CallOpts)
}

// Proposals is a free data retrieval call binding the contract method 0x32ed5b12.
//
// Solidity: function proposals(bytes32 ) view returns(uint8 _status, uint16 _yesVotes, uint8 _yesVotesTotal)
func (_StakePortalRate *StakePortalRateCaller) Proposals(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Status        uint8
	YesVotes      uint16
	YesVotesTotal uint8
}, error) {
	var out []interface{}
	err := _StakePortalRate.contract.Call(opts, &out, "proposals", arg0)

	outstruct := new(struct {
		Status        uint8
		YesVotes      uint16
		YesVotesTotal uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Status = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.YesVotes = *abi.ConvertType(out[1], new(uint16)).(*uint16)
	outstruct.YesVotesTotal = *abi.ConvertType(out[2], new(uint8)).(*uint8)

	return *outstruct, err

}

// Proposals is a free data retrieval call binding the contract method 0x32ed5b12.
//
// Solidity: function proposals(bytes32 ) view returns(uint8 _status, uint16 _yesVotes, uint8 _yesVotesTotal)
func (_StakePortalRate *StakePortalRateSession) Proposals(arg0 [32]byte) (struct {
	Status        uint8
	YesVotes      uint16
	YesVotesTotal uint8
}, error) {
	return _StakePortalRate.Contract.Proposals(&_StakePortalRate.CallOpts, arg0)
}

// Proposals is a free data retrieval call binding the contract method 0x32ed5b12.
//
// Solidity: function proposals(bytes32 ) view returns(uint8 _status, uint16 _yesVotes, uint8 _yesVotesTotal)
func (_StakePortalRate *StakePortalRateCallerSession) Proposals(arg0 [32]byte) (struct {
	Status        uint8
	YesVotes      uint16
	YesVotesTotal uint8
}, error) {
	return _StakePortalRate.Contract.Proposals(&_StakePortalRate.CallOpts, arg0)
}

// RateChangeLimit is a free data retrieval call binding the contract method 0xc0152b71.
//
// Solidity: function rateChangeLimit() view returns(uint256)
func (_StakePortalRate *StakePortalRateCaller) RateChangeLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakePortalRate.contract.Call(opts, &out, "rateChangeLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RateChangeLimit is a free data retrieval call binding the contract method 0xc0152b71.
//
// Solidity: function rateChangeLimit() view returns(uint256)
func (_StakePortalRate *StakePortalRateSession) RateChangeLimit() (*big.Int, error) {
	return _StakePortalRate.Contract.RateChangeLimit(&_StakePortalRate.CallOpts)
}

// RateChangeLimit is a free data retrieval call binding the contract method 0xc0152b71.
//
// Solidity: function rateChangeLimit() view returns(uint256)
func (_StakePortalRate *StakePortalRateCallerSession) RateChangeLimit() (*big.Int, error) {
	return _StakePortalRate.Contract.RateChangeLimit(&_StakePortalRate.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_StakePortalRate *StakePortalRateCaller) Threshold(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _StakePortalRate.contract.Call(opts, &out, "threshold")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_StakePortalRate *StakePortalRateSession) Threshold() (uint8, error) {
	return _StakePortalRate.Contract.Threshold(&_StakePortalRate.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_StakePortalRate *StakePortalRateCallerSession) Threshold() (uint8, error) {
	return _StakePortalRate.Contract.Threshold(&_StakePortalRate.CallOpts)
}

// AddSubAccount is a paid mutator transaction binding the contract method 0xb52e1eff.
//
// Solidity: function addSubAccount(address _subAccount) returns()
func (_StakePortalRate *StakePortalRateTransactor) AddSubAccount(opts *bind.TransactOpts, _subAccount common.Address) (*types.Transaction, error) {
	return _StakePortalRate.contract.Transact(opts, "addSubAccount", _subAccount)
}

// AddSubAccount is a paid mutator transaction binding the contract method 0xb52e1eff.
//
// Solidity: function addSubAccount(address _subAccount) returns()
func (_StakePortalRate *StakePortalRateSession) AddSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _StakePortalRate.Contract.AddSubAccount(&_StakePortalRate.TransactOpts, _subAccount)
}

// AddSubAccount is a paid mutator transaction binding the contract method 0xb52e1eff.
//
// Solidity: function addSubAccount(address _subAccount) returns()
func (_StakePortalRate *StakePortalRateTransactorSession) AddSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _StakePortalRate.Contract.AddSubAccount(&_StakePortalRate.TransactOpts, _subAccount)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _newThreshold) returns()
func (_StakePortalRate *StakePortalRateTransactor) ChangeThreshold(opts *bind.TransactOpts, _newThreshold *big.Int) (*types.Transaction, error) {
	return _StakePortalRate.contract.Transact(opts, "changeThreshold", _newThreshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _newThreshold) returns()
func (_StakePortalRate *StakePortalRateSession) ChangeThreshold(_newThreshold *big.Int) (*types.Transaction, error) {
	return _StakePortalRate.Contract.ChangeThreshold(&_StakePortalRate.TransactOpts, _newThreshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _newThreshold) returns()
func (_StakePortalRate *StakePortalRateTransactorSession) ChangeThreshold(_newThreshold *big.Int) (*types.Transaction, error) {
	return _StakePortalRate.Contract.ChangeThreshold(&_StakePortalRate.TransactOpts, _newThreshold)
}

// RemoveSubAccount is a paid mutator transaction binding the contract method 0x0ba3aae3.
//
// Solidity: function removeSubAccount(address _subAccount) returns()
func (_StakePortalRate *StakePortalRateTransactor) RemoveSubAccount(opts *bind.TransactOpts, _subAccount common.Address) (*types.Transaction, error) {
	return _StakePortalRate.contract.Transact(opts, "removeSubAccount", _subAccount)
}

// RemoveSubAccount is a paid mutator transaction binding the contract method 0x0ba3aae3.
//
// Solidity: function removeSubAccount(address _subAccount) returns()
func (_StakePortalRate *StakePortalRateSession) RemoveSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _StakePortalRate.Contract.RemoveSubAccount(&_StakePortalRate.TransactOpts, _subAccount)
}

// RemoveSubAccount is a paid mutator transaction binding the contract method 0x0ba3aae3.
//
// Solidity: function removeSubAccount(address _subAccount) returns()
func (_StakePortalRate *StakePortalRateTransactorSession) RemoveSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _StakePortalRate.Contract.RemoveSubAccount(&_StakePortalRate.TransactOpts, _subAccount)
}

// SetRate is a paid mutator transaction binding the contract method 0x34fcf437.
//
// Solidity: function setRate(uint256 _rate) returns()
func (_StakePortalRate *StakePortalRateTransactor) SetRate(opts *bind.TransactOpts, _rate *big.Int) (*types.Transaction, error) {
	return _StakePortalRate.contract.Transact(opts, "setRate", _rate)
}

// SetRate is a paid mutator transaction binding the contract method 0x34fcf437.
//
// Solidity: function setRate(uint256 _rate) returns()
func (_StakePortalRate *StakePortalRateSession) SetRate(_rate *big.Int) (*types.Transaction, error) {
	return _StakePortalRate.Contract.SetRate(&_StakePortalRate.TransactOpts, _rate)
}

// SetRate is a paid mutator transaction binding the contract method 0x34fcf437.
//
// Solidity: function setRate(uint256 _rate) returns()
func (_StakePortalRate *StakePortalRateTransactorSession) SetRate(_rate *big.Int) (*types.Transaction, error) {
	return _StakePortalRate.Contract.SetRate(&_StakePortalRate.TransactOpts, _rate)
}

// SetRateChangeLimit is a paid mutator transaction binding the contract method 0x19826e71.
//
// Solidity: function setRateChangeLimit(uint256 _rateChangeLimit) returns()
func (_StakePortalRate *StakePortalRateTransactor) SetRateChangeLimit(opts *bind.TransactOpts, _rateChangeLimit *big.Int) (*types.Transaction, error) {
	return _StakePortalRate.contract.Transact(opts, "setRateChangeLimit", _rateChangeLimit)
}

// SetRateChangeLimit is a paid mutator transaction binding the contract method 0x19826e71.
//
// Solidity: function setRateChangeLimit(uint256 _rateChangeLimit) returns()
func (_StakePortalRate *StakePortalRateSession) SetRateChangeLimit(_rateChangeLimit *big.Int) (*types.Transaction, error) {
	return _StakePortalRate.Contract.SetRateChangeLimit(&_StakePortalRate.TransactOpts, _rateChangeLimit)
}

// SetRateChangeLimit is a paid mutator transaction binding the contract method 0x19826e71.
//
// Solidity: function setRateChangeLimit(uint256 _rateChangeLimit) returns()
func (_StakePortalRate *StakePortalRateTransactorSession) SetRateChangeLimit(_rateChangeLimit *big.Int) (*types.Transaction, error) {
	return _StakePortalRate.Contract.SetRateChangeLimit(&_StakePortalRate.TransactOpts, _rateChangeLimit)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakePortalRate *StakePortalRateTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _StakePortalRate.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakePortalRate *StakePortalRateSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _StakePortalRate.Contract.TransferOwnership(&_StakePortalRate.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakePortalRate *StakePortalRateTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _StakePortalRate.Contract.TransferOwnership(&_StakePortalRate.TransactOpts, _newOwner)
}

// VoteRate is a paid mutator transaction binding the contract method 0x113e3709.
//
// Solidity: function voteRate(bytes32 _proposalId, uint256 _rate) returns()
func (_StakePortalRate *StakePortalRateTransactor) VoteRate(opts *bind.TransactOpts, _proposalId [32]byte, _rate *big.Int) (*types.Transaction, error) {
	return _StakePortalRate.contract.Transact(opts, "voteRate", _proposalId, _rate)
}

// VoteRate is a paid mutator transaction binding the contract method 0x113e3709.
//
// Solidity: function voteRate(bytes32 _proposalId, uint256 _rate) returns()
func (_StakePortalRate *StakePortalRateSession) VoteRate(_proposalId [32]byte, _rate *big.Int) (*types.Transaction, error) {
	return _StakePortalRate.Contract.VoteRate(&_StakePortalRate.TransactOpts, _proposalId, _rate)
}

// VoteRate is a paid mutator transaction binding the contract method 0x113e3709.
//
// Solidity: function voteRate(bytes32 _proposalId, uint256 _rate) returns()
func (_StakePortalRate *StakePortalRateTransactorSession) VoteRate(_proposalId [32]byte, _rate *big.Int) (*types.Transaction, error) {
	return _StakePortalRate.Contract.VoteRate(&_StakePortalRate.TransactOpts, _proposalId, _rate)
}

// StakePortalRateProposalExecutedIterator is returned from FilterProposalExecuted and is used to iterate over the raw logs and unpacked data for ProposalExecuted events raised by the StakePortalRate contract.
type StakePortalRateProposalExecutedIterator struct {
	Event *StakePortalRateProposalExecuted // Event containing the contract specifics and raw log

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
func (it *StakePortalRateProposalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakePortalRateProposalExecuted)
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
		it.Event = new(StakePortalRateProposalExecuted)
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
func (it *StakePortalRateProposalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakePortalRateProposalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakePortalRateProposalExecuted represents a ProposalExecuted event raised by the StakePortalRate contract.
type StakePortalRateProposalExecuted struct {
	ProposalId [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalExecuted is a free log retrieval operation binding the contract event 0x7b1bcf1ccf901a11589afff5504d59fd0a53780eed2a952adade0348985139e0.
//
// Solidity: event ProposalExecuted(bytes32 indexed proposalId)
func (_StakePortalRate *StakePortalRateFilterer) FilterProposalExecuted(opts *bind.FilterOpts, proposalId [][32]byte) (*StakePortalRateProposalExecutedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _StakePortalRate.contract.FilterLogs(opts, "ProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &StakePortalRateProposalExecutedIterator{contract: _StakePortalRate.contract, event: "ProposalExecuted", logs: logs, sub: sub}, nil
}

// WatchProposalExecuted is a free log subscription operation binding the contract event 0x7b1bcf1ccf901a11589afff5504d59fd0a53780eed2a952adade0348985139e0.
//
// Solidity: event ProposalExecuted(bytes32 indexed proposalId)
func (_StakePortalRate *StakePortalRateFilterer) WatchProposalExecuted(opts *bind.WatchOpts, sink chan<- *StakePortalRateProposalExecuted, proposalId [][32]byte) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _StakePortalRate.contract.WatchLogs(opts, "ProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakePortalRateProposalExecuted)
				if err := _StakePortalRate.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
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

// ParseProposalExecuted is a log parse operation binding the contract event 0x7b1bcf1ccf901a11589afff5504d59fd0a53780eed2a952adade0348985139e0.
//
// Solidity: event ProposalExecuted(bytes32 indexed proposalId)
func (_StakePortalRate *StakePortalRateFilterer) ParseProposalExecuted(log types.Log) (*StakePortalRateProposalExecuted, error) {
	event := new(StakePortalRateProposalExecuted)
	if err := _StakePortalRate.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakePortalRateSetRateIterator is returned from FilterSetRate and is used to iterate over the raw logs and unpacked data for SetRate events raised by the StakePortalRate contract.
type StakePortalRateSetRateIterator struct {
	Event *StakePortalRateSetRate // Event containing the contract specifics and raw log

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
func (it *StakePortalRateSetRateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakePortalRateSetRate)
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
		it.Event = new(StakePortalRateSetRate)
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
func (it *StakePortalRateSetRateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakePortalRateSetRateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakePortalRateSetRate represents a SetRate event raised by the StakePortalRate contract.
type StakePortalRateSetRate struct {
	Rate *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSetRate is a free log retrieval operation binding the contract event 0x2640b4015d3473fd09bf2b30939e17deb4068cdacf3892136e737e166ceb3210.
//
// Solidity: event SetRate(uint256 rate)
func (_StakePortalRate *StakePortalRateFilterer) FilterSetRate(opts *bind.FilterOpts) (*StakePortalRateSetRateIterator, error) {

	logs, sub, err := _StakePortalRate.contract.FilterLogs(opts, "SetRate")
	if err != nil {
		return nil, err
	}
	return &StakePortalRateSetRateIterator{contract: _StakePortalRate.contract, event: "SetRate", logs: logs, sub: sub}, nil
}

// WatchSetRate is a free log subscription operation binding the contract event 0x2640b4015d3473fd09bf2b30939e17deb4068cdacf3892136e737e166ceb3210.
//
// Solidity: event SetRate(uint256 rate)
func (_StakePortalRate *StakePortalRateFilterer) WatchSetRate(opts *bind.WatchOpts, sink chan<- *StakePortalRateSetRate) (event.Subscription, error) {

	logs, sub, err := _StakePortalRate.contract.WatchLogs(opts, "SetRate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakePortalRateSetRate)
				if err := _StakePortalRate.contract.UnpackLog(event, "SetRate", log); err != nil {
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

// ParseSetRate is a log parse operation binding the contract event 0x2640b4015d3473fd09bf2b30939e17deb4068cdacf3892136e737e166ceb3210.
//
// Solidity: event SetRate(uint256 rate)
func (_StakePortalRate *StakePortalRateFilterer) ParseSetRate(log types.Log) (*StakePortalRateSetRate, error) {
	event := new(StakePortalRateSetRate)
	if err := _StakePortalRate.contract.UnpackLog(event, "SetRate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
