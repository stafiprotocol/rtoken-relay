// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stake_erc20_portal

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

// StakeERC20PortalMetaData contains all meta data concerning the StakeERC20Portal contract.
var StakeERC20PortalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_stakePoolAddressList\",\"type\":\"address[]\"},{\"internalType\":\"uint8[]\",\"name\":\"_chainIdList\",\"type\":\"uint8[]\"},{\"internalType\":\"address\",\"name\":\"_erc20TokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_relayFee\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stafiRecipient\",\"type\":\"bytes32\"}],\"name\":\"RecoverStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stakePool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"chainId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stafiRecipient\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destRecipient\",\"type\":\"address\"}],\"name\":\"Stake\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint8[]\",\"name\":\"_chaindIdList\",\"type\":\"uint8[]\"}],\"name\":\"addChainId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_stakePoolAddressList\",\"type\":\"address[]\"}],\"name\":\"addStakePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"bridgeFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"chainIdExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"erc20TokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_txHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_stafiRecipient\",\"type\":\"bytes32\"}],\"name\":\"recoverStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"relayFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_chaindId\",\"type\":\"uint8\"}],\"name\":\"rmChainId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"}],\"name\":\"rmStakePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_chainId\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_bridgeFee\",\"type\":\"uint256\"}],\"name\":\"setBridgeFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minAmount\",\"type\":\"uint256\"}],\"name\":\"setMinAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_relayFee\",\"type\":\"uint256\"}],\"name\":\"setRelayFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"_destChainId\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_stafiRecipient\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_destRecipient\",\"type\":\"address\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"stakePoolAddressExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeSwitch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"toggleSwitch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StakeERC20PortalABI is the input ABI used to generate the binding from.
// Deprecated: Use StakeERC20PortalMetaData.ABI instead.
var StakeERC20PortalABI = StakeERC20PortalMetaData.ABI

// StakeERC20Portal is an auto generated Go binding around an Ethereum contract.
type StakeERC20Portal struct {
	StakeERC20PortalCaller     // Read-only binding to the contract
	StakeERC20PortalTransactor // Write-only binding to the contract
	StakeERC20PortalFilterer   // Log filterer for contract events
}

// StakeERC20PortalCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakeERC20PortalCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeERC20PortalTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakeERC20PortalTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeERC20PortalFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakeERC20PortalFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeERC20PortalSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakeERC20PortalSession struct {
	Contract     *StakeERC20Portal // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakeERC20PortalCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakeERC20PortalCallerSession struct {
	Contract *StakeERC20PortalCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// StakeERC20PortalTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakeERC20PortalTransactorSession struct {
	Contract     *StakeERC20PortalTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// StakeERC20PortalRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakeERC20PortalRaw struct {
	Contract *StakeERC20Portal // Generic contract binding to access the raw methods on
}

// StakeERC20PortalCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakeERC20PortalCallerRaw struct {
	Contract *StakeERC20PortalCaller // Generic read-only contract binding to access the raw methods on
}

// StakeERC20PortalTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakeERC20PortalTransactorRaw struct {
	Contract *StakeERC20PortalTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakeERC20Portal creates a new instance of StakeERC20Portal, bound to a specific deployed contract.
func NewStakeERC20Portal(address common.Address, backend bind.ContractBackend) (*StakeERC20Portal, error) {
	contract, err := bindStakeERC20Portal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakeERC20Portal{StakeERC20PortalCaller: StakeERC20PortalCaller{contract: contract}, StakeERC20PortalTransactor: StakeERC20PortalTransactor{contract: contract}, StakeERC20PortalFilterer: StakeERC20PortalFilterer{contract: contract}}, nil
}

// NewStakeERC20PortalCaller creates a new read-only instance of StakeERC20Portal, bound to a specific deployed contract.
func NewStakeERC20PortalCaller(address common.Address, caller bind.ContractCaller) (*StakeERC20PortalCaller, error) {
	contract, err := bindStakeERC20Portal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalCaller{contract: contract}, nil
}

// NewStakeERC20PortalTransactor creates a new write-only instance of StakeERC20Portal, bound to a specific deployed contract.
func NewStakeERC20PortalTransactor(address common.Address, transactor bind.ContractTransactor) (*StakeERC20PortalTransactor, error) {
	contract, err := bindStakeERC20Portal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalTransactor{contract: contract}, nil
}

// NewStakeERC20PortalFilterer creates a new log filterer instance of StakeERC20Portal, bound to a specific deployed contract.
func NewStakeERC20PortalFilterer(address common.Address, filterer bind.ContractFilterer) (*StakeERC20PortalFilterer, error) {
	contract, err := bindStakeERC20Portal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalFilterer{contract: contract}, nil
}

// bindStakeERC20Portal binds a generic wrapper to an already deployed contract.
func bindStakeERC20Portal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakeERC20PortalABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakeERC20Portal *StakeERC20PortalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakeERC20Portal.Contract.StakeERC20PortalCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakeERC20Portal *StakeERC20PortalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.StakeERC20PortalTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakeERC20Portal *StakeERC20PortalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.StakeERC20PortalTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakeERC20Portal *StakeERC20PortalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakeERC20Portal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakeERC20Portal *StakeERC20PortalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakeERC20Portal *StakeERC20PortalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.contract.Transact(opts, method, params...)
}

// BridgeFee is a free data retrieval call binding the contract method 0x046853ec.
//
// Solidity: function bridgeFee(uint8 ) view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCaller) BridgeFee(opts *bind.CallOpts, arg0 uint8) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "bridgeFee", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BridgeFee is a free data retrieval call binding the contract method 0x046853ec.
//
// Solidity: function bridgeFee(uint8 ) view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalSession) BridgeFee(arg0 uint8) (*big.Int, error) {
	return _StakeERC20Portal.Contract.BridgeFee(&_StakeERC20Portal.CallOpts, arg0)
}

// BridgeFee is a free data retrieval call binding the contract method 0x046853ec.
//
// Solidity: function bridgeFee(uint8 ) view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) BridgeFee(arg0 uint8) (*big.Int, error) {
	return _StakeERC20Portal.Contract.BridgeFee(&_StakeERC20Portal.CallOpts, arg0)
}

// ChainIdExist is a free data retrieval call binding the contract method 0x0133228d.
//
// Solidity: function chainIdExist(uint8 ) view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalCaller) ChainIdExist(opts *bind.CallOpts, arg0 uint8) (bool, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "chainIdExist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ChainIdExist is a free data retrieval call binding the contract method 0x0133228d.
//
// Solidity: function chainIdExist(uint8 ) view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalSession) ChainIdExist(arg0 uint8) (bool, error) {
	return _StakeERC20Portal.Contract.ChainIdExist(&_StakeERC20Portal.CallOpts, arg0)
}

// ChainIdExist is a free data retrieval call binding the contract method 0x0133228d.
//
// Solidity: function chainIdExist(uint8 ) view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) ChainIdExist(arg0 uint8) (bool, error) {
	return _StakeERC20Portal.Contract.ChainIdExist(&_StakeERC20Portal.CallOpts, arg0)
}

// Erc20TokenAddress is a free data retrieval call binding the contract method 0xf835cd3c.
//
// Solidity: function erc20TokenAddress() view returns(address)
func (_StakeERC20Portal *StakeERC20PortalCaller) Erc20TokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "erc20TokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Erc20TokenAddress is a free data retrieval call binding the contract method 0xf835cd3c.
//
// Solidity: function erc20TokenAddress() view returns(address)
func (_StakeERC20Portal *StakeERC20PortalSession) Erc20TokenAddress() (common.Address, error) {
	return _StakeERC20Portal.Contract.Erc20TokenAddress(&_StakeERC20Portal.CallOpts)
}

// Erc20TokenAddress is a free data retrieval call binding the contract method 0xf835cd3c.
//
// Solidity: function erc20TokenAddress() view returns(address)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) Erc20TokenAddress() (common.Address, error) {
	return _StakeERC20Portal.Contract.Erc20TokenAddress(&_StakeERC20Portal.CallOpts)
}

// MinAmount is a free data retrieval call binding the contract method 0x9b2cb5d8.
//
// Solidity: function minAmount() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCaller) MinAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "minAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinAmount is a free data retrieval call binding the contract method 0x9b2cb5d8.
//
// Solidity: function minAmount() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalSession) MinAmount() (*big.Int, error) {
	return _StakeERC20Portal.Contract.MinAmount(&_StakeERC20Portal.CallOpts)
}

// MinAmount is a free data retrieval call binding the contract method 0x9b2cb5d8.
//
// Solidity: function minAmount() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) MinAmount() (*big.Int, error) {
	return _StakeERC20Portal.Contract.MinAmount(&_StakeERC20Portal.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakeERC20Portal *StakeERC20PortalCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakeERC20Portal *StakeERC20PortalSession) Owner() (common.Address, error) {
	return _StakeERC20Portal.Contract.Owner(&_StakeERC20Portal.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) Owner() (common.Address, error) {
	return _StakeERC20Portal.Contract.Owner(&_StakeERC20Portal.CallOpts)
}

// RelayFee is a free data retrieval call binding the contract method 0x71d30863.
//
// Solidity: function relayFee() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCaller) RelayFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "relayFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RelayFee is a free data retrieval call binding the contract method 0x71d30863.
//
// Solidity: function relayFee() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalSession) RelayFee() (*big.Int, error) {
	return _StakeERC20Portal.Contract.RelayFee(&_StakeERC20Portal.CallOpts)
}

// RelayFee is a free data retrieval call binding the contract method 0x71d30863.
//
// Solidity: function relayFee() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) RelayFee() (*big.Int, error) {
	return _StakeERC20Portal.Contract.RelayFee(&_StakeERC20Portal.CallOpts)
}

// StakePoolAddressExist is a free data retrieval call binding the contract method 0x3a68e4b2.
//
// Solidity: function stakePoolAddressExist(address ) view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalCaller) StakePoolAddressExist(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "stakePoolAddressExist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StakePoolAddressExist is a free data retrieval call binding the contract method 0x3a68e4b2.
//
// Solidity: function stakePoolAddressExist(address ) view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalSession) StakePoolAddressExist(arg0 common.Address) (bool, error) {
	return _StakeERC20Portal.Contract.StakePoolAddressExist(&_StakeERC20Portal.CallOpts, arg0)
}

// StakePoolAddressExist is a free data retrieval call binding the contract method 0x3a68e4b2.
//
// Solidity: function stakePoolAddressExist(address ) view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) StakePoolAddressExist(arg0 common.Address) (bool, error) {
	return _StakeERC20Portal.Contract.StakePoolAddressExist(&_StakeERC20Portal.CallOpts, arg0)
}

// StakeSwitch is a free data retrieval call binding the contract method 0x24387cdc.
//
// Solidity: function stakeSwitch() view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalCaller) StakeSwitch(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "stakeSwitch")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StakeSwitch is a free data retrieval call binding the contract method 0x24387cdc.
//
// Solidity: function stakeSwitch() view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalSession) StakeSwitch() (bool, error) {
	return _StakeERC20Portal.Contract.StakeSwitch(&_StakeERC20Portal.CallOpts)
}

// StakeSwitch is a free data retrieval call binding the contract method 0x24387cdc.
//
// Solidity: function stakeSwitch() view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) StakeSwitch() (bool, error) {
	return _StakeERC20Portal.Contract.StakeSwitch(&_StakeERC20Portal.CallOpts)
}

// AddChainId is a paid mutator transaction binding the contract method 0x3fa2569b.
//
// Solidity: function addChainId(uint8[] _chaindIdList) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) AddChainId(opts *bind.TransactOpts, _chaindIdList []uint8) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "addChainId", _chaindIdList)
}

// AddChainId is a paid mutator transaction binding the contract method 0x3fa2569b.
//
// Solidity: function addChainId(uint8[] _chaindIdList) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) AddChainId(_chaindIdList []uint8) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.AddChainId(&_StakeERC20Portal.TransactOpts, _chaindIdList)
}

// AddChainId is a paid mutator transaction binding the contract method 0x3fa2569b.
//
// Solidity: function addChainId(uint8[] _chaindIdList) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) AddChainId(_chaindIdList []uint8) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.AddChainId(&_StakeERC20Portal.TransactOpts, _chaindIdList)
}

// AddStakePool is a paid mutator transaction binding the contract method 0xcb08cbae.
//
// Solidity: function addStakePool(address[] _stakePoolAddressList) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) AddStakePool(opts *bind.TransactOpts, _stakePoolAddressList []common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "addStakePool", _stakePoolAddressList)
}

// AddStakePool is a paid mutator transaction binding the contract method 0xcb08cbae.
//
// Solidity: function addStakePool(address[] _stakePoolAddressList) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) AddStakePool(_stakePoolAddressList []common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.AddStakePool(&_StakeERC20Portal.TransactOpts, _stakePoolAddressList)
}

// AddStakePool is a paid mutator transaction binding the contract method 0xcb08cbae.
//
// Solidity: function addStakePool(address[] _stakePoolAddressList) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) AddStakePool(_stakePoolAddressList []common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.AddStakePool(&_StakeERC20Portal.TransactOpts, _stakePoolAddressList)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x919a2f83.
//
// Solidity: function recoverStake(bytes32 _txHash, bytes32 _stafiRecipient) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) RecoverStake(opts *bind.TransactOpts, _txHash [32]byte, _stafiRecipient [32]byte) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "recoverStake", _txHash, _stafiRecipient)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x919a2f83.
//
// Solidity: function recoverStake(bytes32 _txHash, bytes32 _stafiRecipient) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) RecoverStake(_txHash [32]byte, _stafiRecipient [32]byte) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.RecoverStake(&_StakeERC20Portal.TransactOpts, _txHash, _stafiRecipient)
}

// RecoverStake is a paid mutator transaction binding the contract method 0x919a2f83.
//
// Solidity: function recoverStake(bytes32 _txHash, bytes32 _stafiRecipient) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) RecoverStake(_txHash [32]byte, _stafiRecipient [32]byte) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.RecoverStake(&_StakeERC20Portal.TransactOpts, _txHash, _stafiRecipient)
}

// RmChainId is a paid mutator transaction binding the contract method 0xf8a9cf65.
//
// Solidity: function rmChainId(uint8 _chaindId) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) RmChainId(opts *bind.TransactOpts, _chaindId uint8) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "rmChainId", _chaindId)
}

// RmChainId is a paid mutator transaction binding the contract method 0xf8a9cf65.
//
// Solidity: function rmChainId(uint8 _chaindId) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) RmChainId(_chaindId uint8) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.RmChainId(&_StakeERC20Portal.TransactOpts, _chaindId)
}

// RmChainId is a paid mutator transaction binding the contract method 0xf8a9cf65.
//
// Solidity: function rmChainId(uint8 _chaindId) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) RmChainId(_chaindId uint8) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.RmChainId(&_StakeERC20Portal.TransactOpts, _chaindId)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _stakePoolAddress) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) RmStakePool(opts *bind.TransactOpts, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "rmStakePool", _stakePoolAddress)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _stakePoolAddress) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) RmStakePool(_stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.RmStakePool(&_StakeERC20Portal.TransactOpts, _stakePoolAddress)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _stakePoolAddress) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) RmStakePool(_stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.RmStakePool(&_StakeERC20Portal.TransactOpts, _stakePoolAddress)
}

// SetBridgeFee is a paid mutator transaction binding the contract method 0xaddc3519.
//
// Solidity: function setBridgeFee(uint8 _chainId, uint256 _bridgeFee) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) SetBridgeFee(opts *bind.TransactOpts, _chainId uint8, _bridgeFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "setBridgeFee", _chainId, _bridgeFee)
}

// SetBridgeFee is a paid mutator transaction binding the contract method 0xaddc3519.
//
// Solidity: function setBridgeFee(uint8 _chainId, uint256 _bridgeFee) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) SetBridgeFee(_chainId uint8, _bridgeFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetBridgeFee(&_StakeERC20Portal.TransactOpts, _chainId, _bridgeFee)
}

// SetBridgeFee is a paid mutator transaction binding the contract method 0xaddc3519.
//
// Solidity: function setBridgeFee(uint8 _chainId, uint256 _bridgeFee) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) SetBridgeFee(_chainId uint8, _bridgeFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetBridgeFee(&_StakeERC20Portal.TransactOpts, _chainId, _bridgeFee)
}

// SetMinAmount is a paid mutator transaction binding the contract method 0x897b0637.
//
// Solidity: function setMinAmount(uint256 _minAmount) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) SetMinAmount(opts *bind.TransactOpts, _minAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "setMinAmount", _minAmount)
}

// SetMinAmount is a paid mutator transaction binding the contract method 0x897b0637.
//
// Solidity: function setMinAmount(uint256 _minAmount) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) SetMinAmount(_minAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetMinAmount(&_StakeERC20Portal.TransactOpts, _minAmount)
}

// SetMinAmount is a paid mutator transaction binding the contract method 0x897b0637.
//
// Solidity: function setMinAmount(uint256 _minAmount) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) SetMinAmount(_minAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetMinAmount(&_StakeERC20Portal.TransactOpts, _minAmount)
}

// SetRelayFee is a paid mutator transaction binding the contract method 0x98385109.
//
// Solidity: function setRelayFee(uint256 _relayFee) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) SetRelayFee(opts *bind.TransactOpts, _relayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "setRelayFee", _relayFee)
}

// SetRelayFee is a paid mutator transaction binding the contract method 0x98385109.
//
// Solidity: function setRelayFee(uint256 _relayFee) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) SetRelayFee(_relayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetRelayFee(&_StakeERC20Portal.TransactOpts, _relayFee)
}

// SetRelayFee is a paid mutator transaction binding the contract method 0x98385109.
//
// Solidity: function setRelayFee(uint256 _relayFee) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) SetRelayFee(_relayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetRelayFee(&_StakeERC20Portal.TransactOpts, _relayFee)
}

// Stake is a paid mutator transaction binding the contract method 0x0f52a4ec.
//
// Solidity: function stake(address _stakePoolAddress, uint256 _amount, uint8 _destChainId, bytes32 _stafiRecipient, address _destRecipient) payable returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) Stake(opts *bind.TransactOpts, _stakePoolAddress common.Address, _amount *big.Int, _destChainId uint8, _stafiRecipient [32]byte, _destRecipient common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "stake", _stakePoolAddress, _amount, _destChainId, _stafiRecipient, _destRecipient)
}

// Stake is a paid mutator transaction binding the contract method 0x0f52a4ec.
//
// Solidity: function stake(address _stakePoolAddress, uint256 _amount, uint8 _destChainId, bytes32 _stafiRecipient, address _destRecipient) payable returns()
func (_StakeERC20Portal *StakeERC20PortalSession) Stake(_stakePoolAddress common.Address, _amount *big.Int, _destChainId uint8, _stafiRecipient [32]byte, _destRecipient common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.Stake(&_StakeERC20Portal.TransactOpts, _stakePoolAddress, _amount, _destChainId, _stafiRecipient, _destRecipient)
}

// Stake is a paid mutator transaction binding the contract method 0x0f52a4ec.
//
// Solidity: function stake(address _stakePoolAddress, uint256 _amount, uint8 _destChainId, bytes32 _stafiRecipient, address _destRecipient) payable returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) Stake(_stakePoolAddress common.Address, _amount *big.Int, _destChainId uint8, _stafiRecipient [32]byte, _destRecipient common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.Stake(&_StakeERC20Portal.TransactOpts, _stakePoolAddress, _amount, _destChainId, _stafiRecipient, _destRecipient)
}

// ToggleSwitch is a paid mutator transaction binding the contract method 0xbfe4d261.
//
// Solidity: function toggleSwitch() returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) ToggleSwitch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "toggleSwitch")
}

// ToggleSwitch is a paid mutator transaction binding the contract method 0xbfe4d261.
//
// Solidity: function toggleSwitch() returns()
func (_StakeERC20Portal *StakeERC20PortalSession) ToggleSwitch() (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.ToggleSwitch(&_StakeERC20Portal.TransactOpts)
}

// ToggleSwitch is a paid mutator transaction binding the contract method 0xbfe4d261.
//
// Solidity: function toggleSwitch() returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) ToggleSwitch() (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.ToggleSwitch(&_StakeERC20Portal.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.TransferOwnership(&_StakeERC20Portal.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.TransferOwnership(&_StakeERC20Portal.TransactOpts, _newOwner)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) WithdrawFee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "withdrawFee")
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns()
func (_StakeERC20Portal *StakeERC20PortalSession) WithdrawFee() (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.WithdrawFee(&_StakeERC20Portal.TransactOpts)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0xe941fa78.
//
// Solidity: function withdrawFee() returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) WithdrawFee() (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.WithdrawFee(&_StakeERC20Portal.TransactOpts)
}

// StakeERC20PortalRecoverStakeIterator is returned from FilterRecoverStake and is used to iterate over the raw logs and unpacked data for RecoverStake events raised by the StakeERC20Portal contract.
type StakeERC20PortalRecoverStakeIterator struct {
	Event *StakeERC20PortalRecoverStake // Event containing the contract specifics and raw log

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
func (it *StakeERC20PortalRecoverStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeERC20PortalRecoverStake)
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
		it.Event = new(StakeERC20PortalRecoverStake)
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
func (it *StakeERC20PortalRecoverStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeERC20PortalRecoverStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeERC20PortalRecoverStake represents a RecoverStake event raised by the StakeERC20Portal contract.
type StakeERC20PortalRecoverStake struct {
	TxHash         [32]byte
	StafiRecipient [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRecoverStake is a free log retrieval operation binding the contract event 0xe6824256d477a01bdb00762b9b44dd3e583772d99e373d00d1722223336586c7.
//
// Solidity: event RecoverStake(bytes32 txHash, bytes32 stafiRecipient)
func (_StakeERC20Portal *StakeERC20PortalFilterer) FilterRecoverStake(opts *bind.FilterOpts) (*StakeERC20PortalRecoverStakeIterator, error) {

	logs, sub, err := _StakeERC20Portal.contract.FilterLogs(opts, "RecoverStake")
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalRecoverStakeIterator{contract: _StakeERC20Portal.contract, event: "RecoverStake", logs: logs, sub: sub}, nil
}

// WatchRecoverStake is a free log subscription operation binding the contract event 0xe6824256d477a01bdb00762b9b44dd3e583772d99e373d00d1722223336586c7.
//
// Solidity: event RecoverStake(bytes32 txHash, bytes32 stafiRecipient)
func (_StakeERC20Portal *StakeERC20PortalFilterer) WatchRecoverStake(opts *bind.WatchOpts, sink chan<- *StakeERC20PortalRecoverStake) (event.Subscription, error) {

	logs, sub, err := _StakeERC20Portal.contract.WatchLogs(opts, "RecoverStake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeERC20PortalRecoverStake)
				if err := _StakeERC20Portal.contract.UnpackLog(event, "RecoverStake", log); err != nil {
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
func (_StakeERC20Portal *StakeERC20PortalFilterer) ParseRecoverStake(log types.Log) (*StakeERC20PortalRecoverStake, error) {
	event := new(StakeERC20PortalRecoverStake)
	if err := _StakeERC20Portal.contract.UnpackLog(event, "RecoverStake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeERC20PortalStakeIterator is returned from FilterStake and is used to iterate over the raw logs and unpacked data for Stake events raised by the StakeERC20Portal contract.
type StakeERC20PortalStakeIterator struct {
	Event *StakeERC20PortalStake // Event containing the contract specifics and raw log

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
func (it *StakeERC20PortalStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeERC20PortalStake)
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
		it.Event = new(StakeERC20PortalStake)
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
func (it *StakeERC20PortalStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeERC20PortalStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeERC20PortalStake represents a Stake event raised by the StakeERC20Portal contract.
type StakeERC20PortalStake struct {
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
func (_StakeERC20Portal *StakeERC20PortalFilterer) FilterStake(opts *bind.FilterOpts) (*StakeERC20PortalStakeIterator, error) {

	logs, sub, err := _StakeERC20Portal.contract.FilterLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalStakeIterator{contract: _StakeERC20Portal.contract, event: "Stake", logs: logs, sub: sub}, nil
}

// WatchStake is a free log subscription operation binding the contract event 0xa3442c24a17de01319e72fe10476c331914acbf277e5e27004c94169bb16e883.
//
// Solidity: event Stake(address staker, address stakePool, uint256 amount, uint8 chainId, bytes32 stafiRecipient, address destRecipient)
func (_StakeERC20Portal *StakeERC20PortalFilterer) WatchStake(opts *bind.WatchOpts, sink chan<- *StakeERC20PortalStake) (event.Subscription, error) {

	logs, sub, err := _StakeERC20Portal.contract.WatchLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeERC20PortalStake)
				if err := _StakeERC20Portal.contract.UnpackLog(event, "Stake", log); err != nil {
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
func (_StakeERC20Portal *StakeERC20PortalFilterer) ParseStake(log types.Log) (*StakeERC20PortalStake, error) {
	event := new(StakeERC20PortalStake)
	if err := _StakeERC20Portal.contract.UnpackLog(event, "Stake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
