// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package multisig_onchain

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

// MultisigOnchainMetaData contains all meta data concerning the MultisigOnchain contract.
var MultisigOnchainMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"proposalId\",\"type\":\"bytes32\"}],\"name\":\"ProposalExecuted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"addSubAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newThreshold\",\"type\":\"uint256\"}],\"name\":\"changeThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_proposalId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_transactions\",\"type\":\"bytes\"}],\"name\":\"execTransactions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"getSubAccountIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_proposalId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"hasVoted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_initialSubAccounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_initialThreshold\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"proposals\",\"outputs\":[{\"internalType\":\"enumMultisigOnchain.ProposalStatus\",\"name\":\"_status\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"_yesVotes\",\"type\":\"uint16\"},{\"internalType\":\"uint8\",\"name\":\"_yesVotesTotal\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"removeSubAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// MultisigOnchainABI is the input ABI used to generate the binding from.
// Deprecated: Use MultisigOnchainMetaData.ABI instead.
var MultisigOnchainABI = MultisigOnchainMetaData.ABI

// MultisigOnchain is an auto generated Go binding around an Ethereum contract.
type MultisigOnchain struct {
	MultisigOnchainCaller     // Read-only binding to the contract
	MultisigOnchainTransactor // Write-only binding to the contract
	MultisigOnchainFilterer   // Log filterer for contract events
}

// MultisigOnchainCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultisigOnchainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisigOnchainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultisigOnchainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisigOnchainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultisigOnchainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisigOnchainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultisigOnchainSession struct {
	Contract     *MultisigOnchain  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MultisigOnchainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultisigOnchainCallerSession struct {
	Contract *MultisigOnchainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// MultisigOnchainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultisigOnchainTransactorSession struct {
	Contract     *MultisigOnchainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// MultisigOnchainRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultisigOnchainRaw struct {
	Contract *MultisigOnchain // Generic contract binding to access the raw methods on
}

// MultisigOnchainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultisigOnchainCallerRaw struct {
	Contract *MultisigOnchainCaller // Generic read-only contract binding to access the raw methods on
}

// MultisigOnchainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultisigOnchainTransactorRaw struct {
	Contract *MultisigOnchainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultisigOnchain creates a new instance of MultisigOnchain, bound to a specific deployed contract.
func NewMultisigOnchain(address common.Address, backend bind.ContractBackend) (*MultisigOnchain, error) {
	contract, err := bindMultisigOnchain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultisigOnchain{MultisigOnchainCaller: MultisigOnchainCaller{contract: contract}, MultisigOnchainTransactor: MultisigOnchainTransactor{contract: contract}, MultisigOnchainFilterer: MultisigOnchainFilterer{contract: contract}}, nil
}

// NewMultisigOnchainCaller creates a new read-only instance of MultisigOnchain, bound to a specific deployed contract.
func NewMultisigOnchainCaller(address common.Address, caller bind.ContractCaller) (*MultisigOnchainCaller, error) {
	contract, err := bindMultisigOnchain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultisigOnchainCaller{contract: contract}, nil
}

// NewMultisigOnchainTransactor creates a new write-only instance of MultisigOnchain, bound to a specific deployed contract.
func NewMultisigOnchainTransactor(address common.Address, transactor bind.ContractTransactor) (*MultisigOnchainTransactor, error) {
	contract, err := bindMultisigOnchain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultisigOnchainTransactor{contract: contract}, nil
}

// NewMultisigOnchainFilterer creates a new log filterer instance of MultisigOnchain, bound to a specific deployed contract.
func NewMultisigOnchainFilterer(address common.Address, filterer bind.ContractFilterer) (*MultisigOnchainFilterer, error) {
	contract, err := bindMultisigOnchain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultisigOnchainFilterer{contract: contract}, nil
}

// bindMultisigOnchain binds a generic wrapper to an already deployed contract.
func bindMultisigOnchain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MultisigOnchainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultisigOnchain *MultisigOnchainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultisigOnchain.Contract.MultisigOnchainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultisigOnchain *MultisigOnchainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.MultisigOnchainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultisigOnchain *MultisigOnchainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.MultisigOnchainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultisigOnchain *MultisigOnchainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultisigOnchain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultisigOnchain *MultisigOnchainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultisigOnchain *MultisigOnchainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.contract.Transact(opts, method, params...)
}

// GetSubAccountIndex is a free data retrieval call binding the contract method 0x763f8680.
//
// Solidity: function getSubAccountIndex(address _subAccount) view returns(uint256)
func (_MultisigOnchain *MultisigOnchainCaller) GetSubAccountIndex(opts *bind.CallOpts, _subAccount common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MultisigOnchain.contract.Call(opts, &out, "getSubAccountIndex", _subAccount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSubAccountIndex is a free data retrieval call binding the contract method 0x763f8680.
//
// Solidity: function getSubAccountIndex(address _subAccount) view returns(uint256)
func (_MultisigOnchain *MultisigOnchainSession) GetSubAccountIndex(_subAccount common.Address) (*big.Int, error) {
	return _MultisigOnchain.Contract.GetSubAccountIndex(&_MultisigOnchain.CallOpts, _subAccount)
}

// GetSubAccountIndex is a free data retrieval call binding the contract method 0x763f8680.
//
// Solidity: function getSubAccountIndex(address _subAccount) view returns(uint256)
func (_MultisigOnchain *MultisigOnchainCallerSession) GetSubAccountIndex(_subAccount common.Address) (*big.Int, error) {
	return _MultisigOnchain.Contract.GetSubAccountIndex(&_MultisigOnchain.CallOpts, _subAccount)
}

// HasVoted is a free data retrieval call binding the contract method 0xaadc3b72.
//
// Solidity: function hasVoted(bytes32 _proposalId, address _subAccount) view returns(bool)
func (_MultisigOnchain *MultisigOnchainCaller) HasVoted(opts *bind.CallOpts, _proposalId [32]byte, _subAccount common.Address) (bool, error) {
	var out []interface{}
	err := _MultisigOnchain.contract.Call(opts, &out, "hasVoted", _proposalId, _subAccount)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasVoted is a free data retrieval call binding the contract method 0xaadc3b72.
//
// Solidity: function hasVoted(bytes32 _proposalId, address _subAccount) view returns(bool)
func (_MultisigOnchain *MultisigOnchainSession) HasVoted(_proposalId [32]byte, _subAccount common.Address) (bool, error) {
	return _MultisigOnchain.Contract.HasVoted(&_MultisigOnchain.CallOpts, _proposalId, _subAccount)
}

// HasVoted is a free data retrieval call binding the contract method 0xaadc3b72.
//
// Solidity: function hasVoted(bytes32 _proposalId, address _subAccount) view returns(bool)
func (_MultisigOnchain *MultisigOnchainCallerSession) HasVoted(_proposalId [32]byte, _subAccount common.Address) (bool, error) {
	return _MultisigOnchain.Contract.HasVoted(&_MultisigOnchain.CallOpts, _proposalId, _subAccount)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MultisigOnchain *MultisigOnchainCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultisigOnchain.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MultisigOnchain *MultisigOnchainSession) Owner() (common.Address, error) {
	return _MultisigOnchain.Contract.Owner(&_MultisigOnchain.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MultisigOnchain *MultisigOnchainCallerSession) Owner() (common.Address, error) {
	return _MultisigOnchain.Contract.Owner(&_MultisigOnchain.CallOpts)
}

// Proposals is a free data retrieval call binding the contract method 0x32ed5b12.
//
// Solidity: function proposals(bytes32 ) view returns(uint8 _status, uint16 _yesVotes, uint8 _yesVotesTotal)
func (_MultisigOnchain *MultisigOnchainCaller) Proposals(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Status        uint8
	YesVotes      uint16
	YesVotesTotal uint8
}, error) {
	var out []interface{}
	err := _MultisigOnchain.contract.Call(opts, &out, "proposals", arg0)

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
func (_MultisigOnchain *MultisigOnchainSession) Proposals(arg0 [32]byte) (struct {
	Status        uint8
	YesVotes      uint16
	YesVotesTotal uint8
}, error) {
	return _MultisigOnchain.Contract.Proposals(&_MultisigOnchain.CallOpts, arg0)
}

// Proposals is a free data retrieval call binding the contract method 0x32ed5b12.
//
// Solidity: function proposals(bytes32 ) view returns(uint8 _status, uint16 _yesVotes, uint8 _yesVotesTotal)
func (_MultisigOnchain *MultisigOnchainCallerSession) Proposals(arg0 [32]byte) (struct {
	Status        uint8
	YesVotes      uint16
	YesVotesTotal uint8
}, error) {
	return _MultisigOnchain.Contract.Proposals(&_MultisigOnchain.CallOpts, arg0)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_MultisigOnchain *MultisigOnchainCaller) Threshold(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MultisigOnchain.contract.Call(opts, &out, "threshold")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_MultisigOnchain *MultisigOnchainSession) Threshold() (uint8, error) {
	return _MultisigOnchain.Contract.Threshold(&_MultisigOnchain.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_MultisigOnchain *MultisigOnchainCallerSession) Threshold() (uint8, error) {
	return _MultisigOnchain.Contract.Threshold(&_MultisigOnchain.CallOpts)
}

// AddSubAccount is a paid mutator transaction binding the contract method 0xb52e1eff.
//
// Solidity: function addSubAccount(address _subAccount) returns()
func (_MultisigOnchain *MultisigOnchainTransactor) AddSubAccount(opts *bind.TransactOpts, _subAccount common.Address) (*types.Transaction, error) {
	return _MultisigOnchain.contract.Transact(opts, "addSubAccount", _subAccount)
}

// AddSubAccount is a paid mutator transaction binding the contract method 0xb52e1eff.
//
// Solidity: function addSubAccount(address _subAccount) returns()
func (_MultisigOnchain *MultisigOnchainSession) AddSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.AddSubAccount(&_MultisigOnchain.TransactOpts, _subAccount)
}

// AddSubAccount is a paid mutator transaction binding the contract method 0xb52e1eff.
//
// Solidity: function addSubAccount(address _subAccount) returns()
func (_MultisigOnchain *MultisigOnchainTransactorSession) AddSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.AddSubAccount(&_MultisigOnchain.TransactOpts, _subAccount)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _newThreshold) returns()
func (_MultisigOnchain *MultisigOnchainTransactor) ChangeThreshold(opts *bind.TransactOpts, _newThreshold *big.Int) (*types.Transaction, error) {
	return _MultisigOnchain.contract.Transact(opts, "changeThreshold", _newThreshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _newThreshold) returns()
func (_MultisigOnchain *MultisigOnchainSession) ChangeThreshold(_newThreshold *big.Int) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.ChangeThreshold(&_MultisigOnchain.TransactOpts, _newThreshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _newThreshold) returns()
func (_MultisigOnchain *MultisigOnchainTransactorSession) ChangeThreshold(_newThreshold *big.Int) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.ChangeThreshold(&_MultisigOnchain.TransactOpts, _newThreshold)
}

// ExecTransactions is a paid mutator transaction binding the contract method 0x3a9a317b.
//
// Solidity: function execTransactions(bytes32 _proposalId, bytes _transactions) returns()
func (_MultisigOnchain *MultisigOnchainTransactor) ExecTransactions(opts *bind.TransactOpts, _proposalId [32]byte, _transactions []byte) (*types.Transaction, error) {
	return _MultisigOnchain.contract.Transact(opts, "execTransactions", _proposalId, _transactions)
}

// ExecTransactions is a paid mutator transaction binding the contract method 0x3a9a317b.
//
// Solidity: function execTransactions(bytes32 _proposalId, bytes _transactions) returns()
func (_MultisigOnchain *MultisigOnchainSession) ExecTransactions(_proposalId [32]byte, _transactions []byte) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.ExecTransactions(&_MultisigOnchain.TransactOpts, _proposalId, _transactions)
}

// ExecTransactions is a paid mutator transaction binding the contract method 0x3a9a317b.
//
// Solidity: function execTransactions(bytes32 _proposalId, bytes _transactions) returns()
func (_MultisigOnchain *MultisigOnchainTransactorSession) ExecTransactions(_proposalId [32]byte, _transactions []byte) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.ExecTransactions(&_MultisigOnchain.TransactOpts, _proposalId, _transactions)
}

// Initialize is a paid mutator transaction binding the contract method 0x60b5bb3f.
//
// Solidity: function initialize(address[] _initialSubAccounts, uint256 _initialThreshold) returns()
func (_MultisigOnchain *MultisigOnchainTransactor) Initialize(opts *bind.TransactOpts, _initialSubAccounts []common.Address, _initialThreshold *big.Int) (*types.Transaction, error) {
	return _MultisigOnchain.contract.Transact(opts, "initialize", _initialSubAccounts, _initialThreshold)
}

// Initialize is a paid mutator transaction binding the contract method 0x60b5bb3f.
//
// Solidity: function initialize(address[] _initialSubAccounts, uint256 _initialThreshold) returns()
func (_MultisigOnchain *MultisigOnchainSession) Initialize(_initialSubAccounts []common.Address, _initialThreshold *big.Int) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.Initialize(&_MultisigOnchain.TransactOpts, _initialSubAccounts, _initialThreshold)
}

// Initialize is a paid mutator transaction binding the contract method 0x60b5bb3f.
//
// Solidity: function initialize(address[] _initialSubAccounts, uint256 _initialThreshold) returns()
func (_MultisigOnchain *MultisigOnchainTransactorSession) Initialize(_initialSubAccounts []common.Address, _initialThreshold *big.Int) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.Initialize(&_MultisigOnchain.TransactOpts, _initialSubAccounts, _initialThreshold)
}

// RemoveSubAccount is a paid mutator transaction binding the contract method 0x0ba3aae3.
//
// Solidity: function removeSubAccount(address _subAccount) returns()
func (_MultisigOnchain *MultisigOnchainTransactor) RemoveSubAccount(opts *bind.TransactOpts, _subAccount common.Address) (*types.Transaction, error) {
	return _MultisigOnchain.contract.Transact(opts, "removeSubAccount", _subAccount)
}

// RemoveSubAccount is a paid mutator transaction binding the contract method 0x0ba3aae3.
//
// Solidity: function removeSubAccount(address _subAccount) returns()
func (_MultisigOnchain *MultisigOnchainSession) RemoveSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.RemoveSubAccount(&_MultisigOnchain.TransactOpts, _subAccount)
}

// RemoveSubAccount is a paid mutator transaction binding the contract method 0x0ba3aae3.
//
// Solidity: function removeSubAccount(address _subAccount) returns()
func (_MultisigOnchain *MultisigOnchainTransactorSession) RemoveSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.RemoveSubAccount(&_MultisigOnchain.TransactOpts, _subAccount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_MultisigOnchain *MultisigOnchainTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _MultisigOnchain.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_MultisigOnchain *MultisigOnchainSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.TransferOwnership(&_MultisigOnchain.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_MultisigOnchain *MultisigOnchainTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _MultisigOnchain.Contract.TransferOwnership(&_MultisigOnchain.TransactOpts, _newOwner)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MultisigOnchain *MultisigOnchainTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultisigOnchain.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MultisigOnchain *MultisigOnchainSession) Receive() (*types.Transaction, error) {
	return _MultisigOnchain.Contract.Receive(&_MultisigOnchain.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_MultisigOnchain *MultisigOnchainTransactorSession) Receive() (*types.Transaction, error) {
	return _MultisigOnchain.Contract.Receive(&_MultisigOnchain.TransactOpts)
}

// MultisigOnchainProposalExecutedIterator is returned from FilterProposalExecuted and is used to iterate over the raw logs and unpacked data for ProposalExecuted events raised by the MultisigOnchain contract.
type MultisigOnchainProposalExecutedIterator struct {
	Event *MultisigOnchainProposalExecuted // Event containing the contract specifics and raw log

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
func (it *MultisigOnchainProposalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigOnchainProposalExecuted)
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
		it.Event = new(MultisigOnchainProposalExecuted)
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
func (it *MultisigOnchainProposalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigOnchainProposalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigOnchainProposalExecuted represents a ProposalExecuted event raised by the MultisigOnchain contract.
type MultisigOnchainProposalExecuted struct {
	ProposalId [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalExecuted is a free log retrieval operation binding the contract event 0x7b1bcf1ccf901a11589afff5504d59fd0a53780eed2a952adade0348985139e0.
//
// Solidity: event ProposalExecuted(bytes32 indexed proposalId)
func (_MultisigOnchain *MultisigOnchainFilterer) FilterProposalExecuted(opts *bind.FilterOpts, proposalId [][32]byte) (*MultisigOnchainProposalExecutedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _MultisigOnchain.contract.FilterLogs(opts, "ProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &MultisigOnchainProposalExecutedIterator{contract: _MultisigOnchain.contract, event: "ProposalExecuted", logs: logs, sub: sub}, nil
}

// WatchProposalExecuted is a free log subscription operation binding the contract event 0x7b1bcf1ccf901a11589afff5504d59fd0a53780eed2a952adade0348985139e0.
//
// Solidity: event ProposalExecuted(bytes32 indexed proposalId)
func (_MultisigOnchain *MultisigOnchainFilterer) WatchProposalExecuted(opts *bind.WatchOpts, sink chan<- *MultisigOnchainProposalExecuted, proposalId [][32]byte) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _MultisigOnchain.contract.WatchLogs(opts, "ProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigOnchainProposalExecuted)
				if err := _MultisigOnchain.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
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
func (_MultisigOnchain *MultisigOnchainFilterer) ParseProposalExecuted(log types.Log) (*MultisigOnchainProposalExecuted, error) {
	event := new(MultisigOnchainProposalExecuted)
	if err := _MultisigOnchain.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
