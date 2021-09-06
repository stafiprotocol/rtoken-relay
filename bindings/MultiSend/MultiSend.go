// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package MultiSend

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

// MultiSendABI is the input ABI used to generate the binding from.
const MultiSendABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"}],\"name\":\"multiSend\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]"

// MultiSend is an auto generated Go binding around an Ethereum contract.
type MultiSend struct {
	MultiSendCaller     // Read-only binding to the contract
	MultiSendTransactor // Write-only binding to the contract
	MultiSendFilterer   // Log filterer for contract events
}

// MultiSendCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultiSendCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiSendTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultiSendTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiSendFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultiSendFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiSendSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultiSendSession struct {
	Contract     *MultiSend        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MultiSendCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultiSendCallerSession struct {
	Contract *MultiSendCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// MultiSendTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultiSendTransactorSession struct {
	Contract     *MultiSendTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// MultiSendRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultiSendRaw struct {
	Contract *MultiSend // Generic contract binding to access the raw methods on
}

// MultiSendCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultiSendCallerRaw struct {
	Contract *MultiSendCaller // Generic read-only contract binding to access the raw methods on
}

// MultiSendTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultiSendTransactorRaw struct {
	Contract *MultiSendTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultiSend creates a new instance of MultiSend, bound to a specific deployed contract.
func NewMultiSend(address common.Address, backend bind.ContractBackend) (*MultiSend, error) {
	contract, err := bindMultiSend(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiSend{MultiSendCaller: MultiSendCaller{contract: contract}, MultiSendTransactor: MultiSendTransactor{contract: contract}, MultiSendFilterer: MultiSendFilterer{contract: contract}}, nil
}

// NewMultiSendCaller creates a new read-only instance of MultiSend, bound to a specific deployed contract.
func NewMultiSendCaller(address common.Address, caller bind.ContractCaller) (*MultiSendCaller, error) {
	contract, err := bindMultiSend(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiSendCaller{contract: contract}, nil
}

// NewMultiSendTransactor creates a new write-only instance of MultiSend, bound to a specific deployed contract.
func NewMultiSendTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiSendTransactor, error) {
	contract, err := bindMultiSend(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiSendTransactor{contract: contract}, nil
}

// NewMultiSendFilterer creates a new log filterer instance of MultiSend, bound to a specific deployed contract.
func NewMultiSendFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiSendFilterer, error) {
	contract, err := bindMultiSend(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiSendFilterer{contract: contract}, nil
}

// bindMultiSend binds a generic wrapper to an already deployed contract.
func bindMultiSend(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MultiSendABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiSend *MultiSendRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiSend.Contract.MultiSendCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiSend *MultiSendRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiSend.Contract.MultiSendTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiSend *MultiSendRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiSend.Contract.MultiSendTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiSend *MultiSendCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiSend.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiSend *MultiSendTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiSend.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiSend *MultiSendTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiSend.Contract.contract.Transact(opts, method, params...)
}

// MultiSend is a paid mutator transaction binding the contract method 0x8d80ff0a.
//
// Solidity: function multiSend(bytes transactions) payable returns()
func (_MultiSend *MultiSendTransactor) MultiSend(opts *bind.TransactOpts, transactions []byte) (*types.Transaction, error) {
	return _MultiSend.contract.Transact(opts, "multiSend", transactions)
}

// MultiSend is a paid mutator transaction binding the contract method 0x8d80ff0a.
//
// Solidity: function multiSend(bytes transactions) payable returns()
func (_MultiSend *MultiSendSession) MultiSend(transactions []byte) (*types.Transaction, error) {
	return _MultiSend.Contract.MultiSend(&_MultiSend.TransactOpts, transactions)
}

// MultiSend is a paid mutator transaction binding the contract method 0x8d80ff0a.
//
// Solidity: function multiSend(bytes transactions) payable returns()
func (_MultiSend *MultiSendTransactorSession) MultiSend(transactions []byte) (*types.Transaction, error) {
	return _MultiSend.Contract.MultiSend(&_MultiSend.TransactOpts, transactions)
}
