// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Multisig

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

// MultisigABI is the input ABI used to generate the binding from.
const MultisigABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"AddedOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"ChangedThreshold\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumEnum.HashState\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"ExecutionResult\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"RemovedOwner\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"TxHashs\",\"outputs\":[{\"internalType\":\"enumEnum.HashState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_threshold\",\"type\":\"uint256\"}],\"name\":\"addOwnerWithThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_threshold\",\"type\":\"uint256\"}],\"name\":\"changeThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"enumEnum.Operation\",\"name\":\"operation\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"safeTxGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8[]\",\"name\":\"vs\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"}],\"name\":\"execTransaction\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwners\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_owners\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_threshold\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"prevOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_threshold\",\"type\":\"uint256\"}],\"name\":\"removeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"prevOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"swapOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Multisig is an auto generated Go binding around an Ethereum contract.
type Multisig struct {
	MultisigCaller     // Read-only binding to the contract
	MultisigTransactor // Write-only binding to the contract
	MultisigFilterer   // Log filterer for contract events
}

// MultisigCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultisigCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisigTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultisigTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisigFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultisigFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisigSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultisigSession struct {
	Contract     *Multisig         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MultisigCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultisigCallerSession struct {
	Contract *MultisigCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// MultisigTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultisigTransactorSession struct {
	Contract     *MultisigTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// MultisigRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultisigRaw struct {
	Contract *Multisig // Generic contract binding to access the raw methods on
}

// MultisigCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultisigCallerRaw struct {
	Contract *MultisigCaller // Generic read-only contract binding to access the raw methods on
}

// MultisigTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultisigTransactorRaw struct {
	Contract *MultisigTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultisig creates a new instance of Multisig, bound to a specific deployed contract.
func NewMultisig(address common.Address, backend bind.ContractBackend) (*Multisig, error) {
	contract, err := bindMultisig(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Multisig{MultisigCaller: MultisigCaller{contract: contract}, MultisigTransactor: MultisigTransactor{contract: contract}, MultisigFilterer: MultisigFilterer{contract: contract}}, nil
}

// NewMultisigCaller creates a new read-only instance of Multisig, bound to a specific deployed contract.
func NewMultisigCaller(address common.Address, caller bind.ContractCaller) (*MultisigCaller, error) {
	contract, err := bindMultisig(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultisigCaller{contract: contract}, nil
}

// NewMultisigTransactor creates a new write-only instance of Multisig, bound to a specific deployed contract.
func NewMultisigTransactor(address common.Address, transactor bind.ContractTransactor) (*MultisigTransactor, error) {
	contract, err := bindMultisig(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultisigTransactor{contract: contract}, nil
}

// NewMultisigFilterer creates a new log filterer instance of Multisig, bound to a specific deployed contract.
func NewMultisigFilterer(address common.Address, filterer bind.ContractFilterer) (*MultisigFilterer, error) {
	contract, err := bindMultisig(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultisigFilterer{contract: contract}, nil
}

// bindMultisig binds a generic wrapper to an already deployed contract.
func bindMultisig(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MultisigABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Multisig *MultisigRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multisig.Contract.MultisigCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Multisig *MultisigRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multisig.Contract.MultisigTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Multisig *MultisigRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Multisig.Contract.MultisigTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Multisig *MultisigCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multisig.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Multisig *MultisigTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multisig.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Multisig *MultisigTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Multisig.Contract.contract.Transact(opts, method, params...)
}

// TxHashs is a free data retrieval call binding the contract method 0x4e3b0b48.
//
// Solidity: function TxHashs(bytes32 ) view returns(uint8)
func (_Multisig *MultisigCaller) TxHashs(opts *bind.CallOpts, arg0 [32]byte) (uint8, error) {
	var out []interface{}
	err := _Multisig.contract.Call(opts, &out, "TxHashs", arg0)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// TxHashs is a free data retrieval call binding the contract method 0x4e3b0b48.
//
// Solidity: function TxHashs(bytes32 ) view returns(uint8)
func (_Multisig *MultisigSession) TxHashs(arg0 [32]byte) (uint8, error) {
	return _Multisig.Contract.TxHashs(&_Multisig.CallOpts, arg0)
}

// TxHashs is a free data retrieval call binding the contract method 0x4e3b0b48.
//
// Solidity: function TxHashs(bytes32 ) view returns(uint8)
func (_Multisig *MultisigCallerSession) TxHashs(arg0 [32]byte) (uint8, error) {
	return _Multisig.Contract.TxHashs(&_Multisig.CallOpts, arg0)
}

// GetChainId is a free data retrieval call binding the contract method 0x3408e470.
//
// Solidity: function getChainId() pure returns(uint256)
func (_Multisig *MultisigCaller) GetChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multisig.contract.Call(opts, &out, "getChainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetChainId is a free data retrieval call binding the contract method 0x3408e470.
//
// Solidity: function getChainId() pure returns(uint256)
func (_Multisig *MultisigSession) GetChainId() (*big.Int, error) {
	return _Multisig.Contract.GetChainId(&_Multisig.CallOpts)
}

// GetChainId is a free data retrieval call binding the contract method 0x3408e470.
//
// Solidity: function getChainId() pure returns(uint256)
func (_Multisig *MultisigCallerSession) GetChainId() (*big.Int, error) {
	return _Multisig.Contract.GetChainId(&_Multisig.CallOpts)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_Multisig *MultisigCaller) GetOwners(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Multisig.contract.Call(opts, &out, "getOwners")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_Multisig *MultisigSession) GetOwners() ([]common.Address, error) {
	return _Multisig.Contract.GetOwners(&_Multisig.CallOpts)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_Multisig *MultisigCallerSession) GetOwners() ([]common.Address, error) {
	return _Multisig.Contract.GetOwners(&_Multisig.CallOpts)
}

// GetThreshold is a free data retrieval call binding the contract method 0xe75235b8.
//
// Solidity: function getThreshold() view returns(uint256)
func (_Multisig *MultisigCaller) GetThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multisig.contract.Call(opts, &out, "getThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetThreshold is a free data retrieval call binding the contract method 0xe75235b8.
//
// Solidity: function getThreshold() view returns(uint256)
func (_Multisig *MultisigSession) GetThreshold() (*big.Int, error) {
	return _Multisig.Contract.GetThreshold(&_Multisig.CallOpts)
}

// GetThreshold is a free data retrieval call binding the contract method 0xe75235b8.
//
// Solidity: function getThreshold() view returns(uint256)
func (_Multisig *MultisigCallerSession) GetThreshold() (*big.Int, error) {
	return _Multisig.Contract.GetThreshold(&_Multisig.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address owner) view returns(bool)
func (_Multisig *MultisigCaller) IsOwner(opts *bind.CallOpts, owner common.Address) (bool, error) {
	var out []interface{}
	err := _Multisig.contract.Call(opts, &out, "isOwner", owner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address owner) view returns(bool)
func (_Multisig *MultisigSession) IsOwner(owner common.Address) (bool, error) {
	return _Multisig.Contract.IsOwner(&_Multisig.CallOpts, owner)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address owner) view returns(bool)
func (_Multisig *MultisigCallerSession) IsOwner(owner common.Address) (bool, error) {
	return _Multisig.Contract.IsOwner(&_Multisig.CallOpts, owner)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_Multisig *MultisigCaller) Nonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Multisig.contract.Call(opts, &out, "nonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_Multisig *MultisigSession) Nonce() (*big.Int, error) {
	return _Multisig.Contract.Nonce(&_Multisig.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_Multisig *MultisigCallerSession) Nonce() (*big.Int, error) {
	return _Multisig.Contract.Nonce(&_Multisig.CallOpts)
}

// AddOwnerWithThreshold is a paid mutator transaction binding the contract method 0x0d582f13.
//
// Solidity: function addOwnerWithThreshold(address owner, uint256 _threshold) returns()
func (_Multisig *MultisigTransactor) AddOwnerWithThreshold(opts *bind.TransactOpts, owner common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _Multisig.contract.Transact(opts, "addOwnerWithThreshold", owner, _threshold)
}

// AddOwnerWithThreshold is a paid mutator transaction binding the contract method 0x0d582f13.
//
// Solidity: function addOwnerWithThreshold(address owner, uint256 _threshold) returns()
func (_Multisig *MultisigSession) AddOwnerWithThreshold(owner common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _Multisig.Contract.AddOwnerWithThreshold(&_Multisig.TransactOpts, owner, _threshold)
}

// AddOwnerWithThreshold is a paid mutator transaction binding the contract method 0x0d582f13.
//
// Solidity: function addOwnerWithThreshold(address owner, uint256 _threshold) returns()
func (_Multisig *MultisigTransactorSession) AddOwnerWithThreshold(owner common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _Multisig.Contract.AddOwnerWithThreshold(&_Multisig.TransactOpts, owner, _threshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _threshold) returns()
func (_Multisig *MultisigTransactor) ChangeThreshold(opts *bind.TransactOpts, _threshold *big.Int) (*types.Transaction, error) {
	return _Multisig.contract.Transact(opts, "changeThreshold", _threshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _threshold) returns()
func (_Multisig *MultisigSession) ChangeThreshold(_threshold *big.Int) (*types.Transaction, error) {
	return _Multisig.Contract.ChangeThreshold(&_Multisig.TransactOpts, _threshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _threshold) returns()
func (_Multisig *MultisigTransactorSession) ChangeThreshold(_threshold *big.Int) (*types.Transaction, error) {
	return _Multisig.Contract.ChangeThreshold(&_Multisig.TransactOpts, _threshold)
}

// ExecTransaction is a paid mutator transaction binding the contract method 0xd84edcb9.
//
// Solidity: function execTransaction(address to, uint256 value, bytes data, uint8 operation, uint256 safeTxGas, bytes32 txHash, uint8[] vs, bytes32[] rs, bytes32[] ss) payable returns(bool success)
func (_Multisig *MultisigTransactor) ExecTransaction(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, txHash [32]byte, vs []uint8, rs [][32]byte, ss [][32]byte) (*types.Transaction, error) {
	return _Multisig.contract.Transact(opts, "execTransaction", to, value, data, operation, safeTxGas, txHash, vs, rs, ss)
}

// ExecTransaction is a paid mutator transaction binding the contract method 0xd84edcb9.
//
// Solidity: function execTransaction(address to, uint256 value, bytes data, uint8 operation, uint256 safeTxGas, bytes32 txHash, uint8[] vs, bytes32[] rs, bytes32[] ss) payable returns(bool success)
func (_Multisig *MultisigSession) ExecTransaction(to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, txHash [32]byte, vs []uint8, rs [][32]byte, ss [][32]byte) (*types.Transaction, error) {
	return _Multisig.Contract.ExecTransaction(&_Multisig.TransactOpts, to, value, data, operation, safeTxGas, txHash, vs, rs, ss)
}

// ExecTransaction is a paid mutator transaction binding the contract method 0xd84edcb9.
//
// Solidity: function execTransaction(address to, uint256 value, bytes data, uint8 operation, uint256 safeTxGas, bytes32 txHash, uint8[] vs, bytes32[] rs, bytes32[] ss) payable returns(bool success)
func (_Multisig *MultisigTransactorSession) ExecTransaction(to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, txHash [32]byte, vs []uint8, rs [][32]byte, ss [][32]byte) (*types.Transaction, error) {
	return _Multisig.Contract.ExecTransaction(&_Multisig.TransactOpts, to, value, data, operation, safeTxGas, txHash, vs, rs, ss)
}

// Initialize is a paid mutator transaction binding the contract method 0x60b5bb3f.
//
// Solidity: function initialize(address[] _owners, uint256 _threshold) returns()
func (_Multisig *MultisigTransactor) Initialize(opts *bind.TransactOpts, _owners []common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _Multisig.contract.Transact(opts, "initialize", _owners, _threshold)
}

// Initialize is a paid mutator transaction binding the contract method 0x60b5bb3f.
//
// Solidity: function initialize(address[] _owners, uint256 _threshold) returns()
func (_Multisig *MultisigSession) Initialize(_owners []common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _Multisig.Contract.Initialize(&_Multisig.TransactOpts, _owners, _threshold)
}

// Initialize is a paid mutator transaction binding the contract method 0x60b5bb3f.
//
// Solidity: function initialize(address[] _owners, uint256 _threshold) returns()
func (_Multisig *MultisigTransactorSession) Initialize(_owners []common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _Multisig.Contract.Initialize(&_Multisig.TransactOpts, _owners, _threshold)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0xf8dc5dd9.
//
// Solidity: function removeOwner(address prevOwner, address owner, uint256 _threshold) returns()
func (_Multisig *MultisigTransactor) RemoveOwner(opts *bind.TransactOpts, prevOwner common.Address, owner common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _Multisig.contract.Transact(opts, "removeOwner", prevOwner, owner, _threshold)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0xf8dc5dd9.
//
// Solidity: function removeOwner(address prevOwner, address owner, uint256 _threshold) returns()
func (_Multisig *MultisigSession) RemoveOwner(prevOwner common.Address, owner common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _Multisig.Contract.RemoveOwner(&_Multisig.TransactOpts, prevOwner, owner, _threshold)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0xf8dc5dd9.
//
// Solidity: function removeOwner(address prevOwner, address owner, uint256 _threshold) returns()
func (_Multisig *MultisigTransactorSession) RemoveOwner(prevOwner common.Address, owner common.Address, _threshold *big.Int) (*types.Transaction, error) {
	return _Multisig.Contract.RemoveOwner(&_Multisig.TransactOpts, prevOwner, owner, _threshold)
}

// SwapOwner is a paid mutator transaction binding the contract method 0xe318b52b.
//
// Solidity: function swapOwner(address prevOwner, address oldOwner, address newOwner) returns()
func (_Multisig *MultisigTransactor) SwapOwner(opts *bind.TransactOpts, prevOwner common.Address, oldOwner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _Multisig.contract.Transact(opts, "swapOwner", prevOwner, oldOwner, newOwner)
}

// SwapOwner is a paid mutator transaction binding the contract method 0xe318b52b.
//
// Solidity: function swapOwner(address prevOwner, address oldOwner, address newOwner) returns()
func (_Multisig *MultisigSession) SwapOwner(prevOwner common.Address, oldOwner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _Multisig.Contract.SwapOwner(&_Multisig.TransactOpts, prevOwner, oldOwner, newOwner)
}

// SwapOwner is a paid mutator transaction binding the contract method 0xe318b52b.
//
// Solidity: function swapOwner(address prevOwner, address oldOwner, address newOwner) returns()
func (_Multisig *MultisigTransactorSession) SwapOwner(prevOwner common.Address, oldOwner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _Multisig.Contract.SwapOwner(&_Multisig.TransactOpts, prevOwner, oldOwner, newOwner)
}

// MultisigAddedOwnerIterator is returned from FilterAddedOwner and is used to iterate over the raw logs and unpacked data for AddedOwner events raised by the Multisig contract.
type MultisigAddedOwnerIterator struct {
	Event *MultisigAddedOwner // Event containing the contract specifics and raw log

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
func (it *MultisigAddedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigAddedOwner)
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
		it.Event = new(MultisigAddedOwner)
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
func (it *MultisigAddedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigAddedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigAddedOwner represents a AddedOwner event raised by the Multisig contract.
type MultisigAddedOwner struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAddedOwner is a free log retrieval operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: event AddedOwner(address owner)
func (_Multisig *MultisigFilterer) FilterAddedOwner(opts *bind.FilterOpts) (*MultisigAddedOwnerIterator, error) {

	logs, sub, err := _Multisig.contract.FilterLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return &MultisigAddedOwnerIterator{contract: _Multisig.contract, event: "AddedOwner", logs: logs, sub: sub}, nil
}

// WatchAddedOwner is a free log subscription operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: event AddedOwner(address owner)
func (_Multisig *MultisigFilterer) WatchAddedOwner(opts *bind.WatchOpts, sink chan<- *MultisigAddedOwner) (event.Subscription, error) {

	logs, sub, err := _Multisig.contract.WatchLogs(opts, "AddedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigAddedOwner)
				if err := _Multisig.contract.UnpackLog(event, "AddedOwner", log); err != nil {
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

// ParseAddedOwner is a log parse operation binding the contract event 0x9465fa0c962cc76958e6373a993326400c1c94f8be2fe3a952adfa7f60b2ea26.
//
// Solidity: event AddedOwner(address owner)
func (_Multisig *MultisigFilterer) ParseAddedOwner(log types.Log) (*MultisigAddedOwner, error) {
	event := new(MultisigAddedOwner)
	if err := _Multisig.contract.UnpackLog(event, "AddedOwner", log); err != nil {
		return nil, err
	}
	return event, nil
}

// MultisigChangedThresholdIterator is returned from FilterChangedThreshold and is used to iterate over the raw logs and unpacked data for ChangedThreshold events raised by the Multisig contract.
type MultisigChangedThresholdIterator struct {
	Event *MultisigChangedThreshold // Event containing the contract specifics and raw log

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
func (it *MultisigChangedThresholdIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigChangedThreshold)
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
		it.Event = new(MultisigChangedThreshold)
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
func (it *MultisigChangedThresholdIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigChangedThresholdIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigChangedThreshold represents a ChangedThreshold event raised by the Multisig contract.
type MultisigChangedThreshold struct {
	Threshold *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChangedThreshold is a free log retrieval operation binding the contract event 0x610f7ff2b304ae8903c3de74c60c6ab1f7d6226b3f52c5161905bb5ad4039c93.
//
// Solidity: event ChangedThreshold(uint256 threshold)
func (_Multisig *MultisigFilterer) FilterChangedThreshold(opts *bind.FilterOpts) (*MultisigChangedThresholdIterator, error) {

	logs, sub, err := _Multisig.contract.FilterLogs(opts, "ChangedThreshold")
	if err != nil {
		return nil, err
	}
	return &MultisigChangedThresholdIterator{contract: _Multisig.contract, event: "ChangedThreshold", logs: logs, sub: sub}, nil
}

// WatchChangedThreshold is a free log subscription operation binding the contract event 0x610f7ff2b304ae8903c3de74c60c6ab1f7d6226b3f52c5161905bb5ad4039c93.
//
// Solidity: event ChangedThreshold(uint256 threshold)
func (_Multisig *MultisigFilterer) WatchChangedThreshold(opts *bind.WatchOpts, sink chan<- *MultisigChangedThreshold) (event.Subscription, error) {

	logs, sub, err := _Multisig.contract.WatchLogs(opts, "ChangedThreshold")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigChangedThreshold)
				if err := _Multisig.contract.UnpackLog(event, "ChangedThreshold", log); err != nil {
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

// ParseChangedThreshold is a log parse operation binding the contract event 0x610f7ff2b304ae8903c3de74c60c6ab1f7d6226b3f52c5161905bb5ad4039c93.
//
// Solidity: event ChangedThreshold(uint256 threshold)
func (_Multisig *MultisigFilterer) ParseChangedThreshold(log types.Log) (*MultisigChangedThreshold, error) {
	event := new(MultisigChangedThreshold)
	if err := _Multisig.contract.UnpackLog(event, "ChangedThreshold", log); err != nil {
		return nil, err
	}
	return event, nil
}

// MultisigExecutionResultIterator is returned from FilterExecutionResult and is used to iterate over the raw logs and unpacked data for ExecutionResult events raised by the Multisig contract.
type MultisigExecutionResultIterator struct {
	Event *MultisigExecutionResult // Event containing the contract specifics and raw log

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
func (it *MultisigExecutionResultIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigExecutionResult)
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
		it.Event = new(MultisigExecutionResult)
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
func (it *MultisigExecutionResultIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigExecutionResultIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigExecutionResult represents a ExecutionResult event raised by the Multisig contract.
type MultisigExecutionResult struct {
	TxHash [32]byte
	Arg1   uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExecutionResult is a free log retrieval operation binding the contract event 0x033f1a25a5e9b1378e89053b7b1e656cffc91fc1e5ecb5b1ff4882814c50b41a.
//
// Solidity: event ExecutionResult(bytes32 txHash, uint8 arg1)
func (_Multisig *MultisigFilterer) FilterExecutionResult(opts *bind.FilterOpts) (*MultisigExecutionResultIterator, error) {

	logs, sub, err := _Multisig.contract.FilterLogs(opts, "ExecutionResult")
	if err != nil {
		return nil, err
	}
	return &MultisigExecutionResultIterator{contract: _Multisig.contract, event: "ExecutionResult", logs: logs, sub: sub}, nil
}

// WatchExecutionResult is a free log subscription operation binding the contract event 0x033f1a25a5e9b1378e89053b7b1e656cffc91fc1e5ecb5b1ff4882814c50b41a.
//
// Solidity: event ExecutionResult(bytes32 txHash, uint8 arg1)
func (_Multisig *MultisigFilterer) WatchExecutionResult(opts *bind.WatchOpts, sink chan<- *MultisigExecutionResult) (event.Subscription, error) {

	logs, sub, err := _Multisig.contract.WatchLogs(opts, "ExecutionResult")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigExecutionResult)
				if err := _Multisig.contract.UnpackLog(event, "ExecutionResult", log); err != nil {
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

// ParseExecutionResult is a log parse operation binding the contract event 0x033f1a25a5e9b1378e89053b7b1e656cffc91fc1e5ecb5b1ff4882814c50b41a.
//
// Solidity: event ExecutionResult(bytes32 txHash, uint8 arg1)
func (_Multisig *MultisigFilterer) ParseExecutionResult(log types.Log) (*MultisigExecutionResult, error) {
	event := new(MultisigExecutionResult)
	if err := _Multisig.contract.UnpackLog(event, "ExecutionResult", log); err != nil {
		return nil, err
	}
	return event, nil
}

// MultisigRemovedOwnerIterator is returned from FilterRemovedOwner and is used to iterate over the raw logs and unpacked data for RemovedOwner events raised by the Multisig contract.
type MultisigRemovedOwnerIterator struct {
	Event *MultisigRemovedOwner // Event containing the contract specifics and raw log

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
func (it *MultisigRemovedOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigRemovedOwner)
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
		it.Event = new(MultisigRemovedOwner)
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
func (it *MultisigRemovedOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigRemovedOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigRemovedOwner represents a RemovedOwner event raised by the Multisig contract.
type MultisigRemovedOwner struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterRemovedOwner is a free log retrieval operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: event RemovedOwner(address owner)
func (_Multisig *MultisigFilterer) FilterRemovedOwner(opts *bind.FilterOpts) (*MultisigRemovedOwnerIterator, error) {

	logs, sub, err := _Multisig.contract.FilterLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return &MultisigRemovedOwnerIterator{contract: _Multisig.contract, event: "RemovedOwner", logs: logs, sub: sub}, nil
}

// WatchRemovedOwner is a free log subscription operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: event RemovedOwner(address owner)
func (_Multisig *MultisigFilterer) WatchRemovedOwner(opts *bind.WatchOpts, sink chan<- *MultisigRemovedOwner) (event.Subscription, error) {

	logs, sub, err := _Multisig.contract.WatchLogs(opts, "RemovedOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigRemovedOwner)
				if err := _Multisig.contract.UnpackLog(event, "RemovedOwner", log); err != nil {
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

// ParseRemovedOwner is a log parse operation binding the contract event 0xf8d49fc529812e9a7c5c50e69c20f0dccc0db8fa95c98bc58cc9a4f1c1299eaf.
//
// Solidity: event RemovedOwner(address owner)
func (_Multisig *MultisigFilterer) ParseRemovedOwner(log types.Log) (*MultisigRemovedOwner, error) {
	event := new(MultisigRemovedOwner)
	if err := _Multisig.contract.UnpackLog(event, "RemovedOwner", log); err != nil {
		return nil, err
	}
	return event, nil
}
