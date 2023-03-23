// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stake_erc20_portal_with_rate

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

// StakeERC20PortalWithRateMetaData contains all meta data concerning the StakeERC20PortalWithRate contract.
var StakeERC20PortalWithRateMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_stakePoolAddressList\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_initialSubAccounts\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"_erc20TokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_rTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_stakeUsePoolAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minStakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_stakeRelayFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_unstakeRelayFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_unstakeFeeCommission\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_initialThreshold\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"proposalId\",\"type\":\"bytes32\"}],\"name\":\"ProposalExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stakePool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rTokenAmount\",\"type\":\"uint256\"}],\"name\":\"Stake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stakePool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rTokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"burnAmount\",\"type\":\"uint256\"}],\"name\":\"Unstake\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_stakePoolAddressList\",\"type\":\"address[]\"}],\"name\":\"addStakePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"addSubAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newThreshold\",\"type\":\"uint256\"}],\"name\":\"changeThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"erc20TokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"getSubAccountIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_proposalId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"hasVoted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"proposals\",\"outputs\":[{\"internalType\":\"enumMultisig.ProposalStatus\",\"name\":\"_status\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"_yesVotes\",\"type\":\"uint16\"},{\"internalType\":\"uint8\",\"name\":\"_yesVotesTotal\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rateChangeLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"removeSubAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"}],\"name\":\"rmStakePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minStakeAmount\",\"type\":\"uint256\"}],\"name\":\"setMinStakeAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"setRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rateChangeLimit\",\"type\":\"uint256\"}],\"name\":\"setRateChangeLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_stakeRelayFee\",\"type\":\"uint256\"}],\"name\":\"setStakeRelayFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakeUsePoolAddress\",\"type\":\"address\"}],\"name\":\"setStakeUsePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_unstakeFeeCommission\",\"type\":\"uint256\"}],\"name\":\"setUnstakeFeeCommission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_unstakeRelayFee\",\"type\":\"uint256\"}],\"name\":\"setUnstakeRelayFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"stakePoolAddressExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeRelayFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeSwitch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeUsePoolAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"}],\"name\":\"stakeWithPool\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"toggleStakeSwitch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalUnstakeProtocolFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rTokenAmount\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unstakeFeeCommission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unstakeRelayFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"}],\"name\":\"unstakeWithPool\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_proposalId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"voteRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StakeERC20PortalWithRateABI is the input ABI used to generate the binding from.
// Deprecated: Use StakeERC20PortalWithRateMetaData.ABI instead.
var StakeERC20PortalWithRateABI = StakeERC20PortalWithRateMetaData.ABI

// StakeERC20PortalWithRate is an auto generated Go binding around an Ethereum contract.
type StakeERC20PortalWithRate struct {
	StakeERC20PortalWithRateCaller     // Read-only binding to the contract
	StakeERC20PortalWithRateTransactor // Write-only binding to the contract
	StakeERC20PortalWithRateFilterer   // Log filterer for contract events
}

// StakeERC20PortalWithRateCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakeERC20PortalWithRateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeERC20PortalWithRateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakeERC20PortalWithRateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeERC20PortalWithRateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakeERC20PortalWithRateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeERC20PortalWithRateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakeERC20PortalWithRateSession struct {
	Contract     *StakeERC20PortalWithRate // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// StakeERC20PortalWithRateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakeERC20PortalWithRateCallerSession struct {
	Contract *StakeERC20PortalWithRateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// StakeERC20PortalWithRateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakeERC20PortalWithRateTransactorSession struct {
	Contract     *StakeERC20PortalWithRateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// StakeERC20PortalWithRateRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakeERC20PortalWithRateRaw struct {
	Contract *StakeERC20PortalWithRate // Generic contract binding to access the raw methods on
}

// StakeERC20PortalWithRateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakeERC20PortalWithRateCallerRaw struct {
	Contract *StakeERC20PortalWithRateCaller // Generic read-only contract binding to access the raw methods on
}

// StakeERC20PortalWithRateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakeERC20PortalWithRateTransactorRaw struct {
	Contract *StakeERC20PortalWithRateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakeERC20PortalWithRate creates a new instance of StakeERC20PortalWithRate, bound to a specific deployed contract.
func NewStakeERC20PortalWithRate(address common.Address, backend bind.ContractBackend) (*StakeERC20PortalWithRate, error) {
	contract, err := bindStakeERC20PortalWithRate(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalWithRate{StakeERC20PortalWithRateCaller: StakeERC20PortalWithRateCaller{contract: contract}, StakeERC20PortalWithRateTransactor: StakeERC20PortalWithRateTransactor{contract: contract}, StakeERC20PortalWithRateFilterer: StakeERC20PortalWithRateFilterer{contract: contract}}, nil
}

// NewStakeERC20PortalWithRateCaller creates a new read-only instance of StakeERC20PortalWithRate, bound to a specific deployed contract.
func NewStakeERC20PortalWithRateCaller(address common.Address, caller bind.ContractCaller) (*StakeERC20PortalWithRateCaller, error) {
	contract, err := bindStakeERC20PortalWithRate(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalWithRateCaller{contract: contract}, nil
}

// NewStakeERC20PortalWithRateTransactor creates a new write-only instance of StakeERC20PortalWithRate, bound to a specific deployed contract.
func NewStakeERC20PortalWithRateTransactor(address common.Address, transactor bind.ContractTransactor) (*StakeERC20PortalWithRateTransactor, error) {
	contract, err := bindStakeERC20PortalWithRate(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalWithRateTransactor{contract: contract}, nil
}

// NewStakeERC20PortalWithRateFilterer creates a new log filterer instance of StakeERC20PortalWithRate, bound to a specific deployed contract.
func NewStakeERC20PortalWithRateFilterer(address common.Address, filterer bind.ContractFilterer) (*StakeERC20PortalWithRateFilterer, error) {
	contract, err := bindStakeERC20PortalWithRate(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalWithRateFilterer{contract: contract}, nil
}

// bindStakeERC20PortalWithRate binds a generic wrapper to an already deployed contract.
func bindStakeERC20PortalWithRate(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakeERC20PortalWithRateABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakeERC20PortalWithRate.Contract.StakeERC20PortalWithRateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.StakeERC20PortalWithRateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.StakeERC20PortalWithRateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakeERC20PortalWithRate.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.contract.Transact(opts, method, params...)
}

// Erc20TokenAddress is a free data retrieval call binding the contract method 0xf835cd3c.
//
// Solidity: function erc20TokenAddress() view returns(address)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) Erc20TokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "erc20TokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Erc20TokenAddress is a free data retrieval call binding the contract method 0xf835cd3c.
//
// Solidity: function erc20TokenAddress() view returns(address)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) Erc20TokenAddress() (common.Address, error) {
	return _StakeERC20PortalWithRate.Contract.Erc20TokenAddress(&_StakeERC20PortalWithRate.CallOpts)
}

// Erc20TokenAddress is a free data retrieval call binding the contract method 0xf835cd3c.
//
// Solidity: function erc20TokenAddress() view returns(address)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) Erc20TokenAddress() (common.Address, error) {
	return _StakeERC20PortalWithRate.Contract.Erc20TokenAddress(&_StakeERC20PortalWithRate.CallOpts)
}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) GetRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "getRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) GetRate() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.GetRate(&_StakeERC20PortalWithRate.CallOpts)
}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) GetRate() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.GetRate(&_StakeERC20PortalWithRate.CallOpts)
}

// GetSubAccountIndex is a free data retrieval call binding the contract method 0x763f8680.
//
// Solidity: function getSubAccountIndex(address _subAccount) view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) GetSubAccountIndex(opts *bind.CallOpts, _subAccount common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "getSubAccountIndex", _subAccount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSubAccountIndex is a free data retrieval call binding the contract method 0x763f8680.
//
// Solidity: function getSubAccountIndex(address _subAccount) view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) GetSubAccountIndex(_subAccount common.Address) (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.GetSubAccountIndex(&_StakeERC20PortalWithRate.CallOpts, _subAccount)
}

// GetSubAccountIndex is a free data retrieval call binding the contract method 0x763f8680.
//
// Solidity: function getSubAccountIndex(address _subAccount) view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) GetSubAccountIndex(_subAccount common.Address) (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.GetSubAccountIndex(&_StakeERC20PortalWithRate.CallOpts, _subAccount)
}

// HasVoted is a free data retrieval call binding the contract method 0xaadc3b72.
//
// Solidity: function hasVoted(bytes32 _proposalId, address _subAccount) view returns(bool)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) HasVoted(opts *bind.CallOpts, _proposalId [32]byte, _subAccount common.Address) (bool, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "hasVoted", _proposalId, _subAccount)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasVoted is a free data retrieval call binding the contract method 0xaadc3b72.
//
// Solidity: function hasVoted(bytes32 _proposalId, address _subAccount) view returns(bool)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) HasVoted(_proposalId [32]byte, _subAccount common.Address) (bool, error) {
	return _StakeERC20PortalWithRate.Contract.HasVoted(&_StakeERC20PortalWithRate.CallOpts, _proposalId, _subAccount)
}

// HasVoted is a free data retrieval call binding the contract method 0xaadc3b72.
//
// Solidity: function hasVoted(bytes32 _proposalId, address _subAccount) view returns(bool)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) HasVoted(_proposalId [32]byte, _subAccount common.Address) (bool, error) {
	return _StakeERC20PortalWithRate.Contract.HasVoted(&_StakeERC20PortalWithRate.CallOpts, _proposalId, _subAccount)
}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) MinStakeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "minStakeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) MinStakeAmount() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.MinStakeAmount(&_StakeERC20PortalWithRate.CallOpts)
}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) MinStakeAmount() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.MinStakeAmount(&_StakeERC20PortalWithRate.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) Owner() (common.Address, error) {
	return _StakeERC20PortalWithRate.Contract.Owner(&_StakeERC20PortalWithRate.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) Owner() (common.Address, error) {
	return _StakeERC20PortalWithRate.Contract.Owner(&_StakeERC20PortalWithRate.CallOpts)
}

// Proposals is a free data retrieval call binding the contract method 0x32ed5b12.
//
// Solidity: function proposals(bytes32 ) view returns(uint8 _status, uint16 _yesVotes, uint8 _yesVotesTotal)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) Proposals(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Status        uint8
	YesVotes      uint16
	YesVotesTotal uint8
}, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "proposals", arg0)

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
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) Proposals(arg0 [32]byte) (struct {
	Status        uint8
	YesVotes      uint16
	YesVotesTotal uint8
}, error) {
	return _StakeERC20PortalWithRate.Contract.Proposals(&_StakeERC20PortalWithRate.CallOpts, arg0)
}

// Proposals is a free data retrieval call binding the contract method 0x32ed5b12.
//
// Solidity: function proposals(bytes32 ) view returns(uint8 _status, uint16 _yesVotes, uint8 _yesVotesTotal)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) Proposals(arg0 [32]byte) (struct {
	Status        uint8
	YesVotes      uint16
	YesVotesTotal uint8
}, error) {
	return _StakeERC20PortalWithRate.Contract.Proposals(&_StakeERC20PortalWithRate.CallOpts, arg0)
}

// RTokenAddress is a free data retrieval call binding the contract method 0x35381358.
//
// Solidity: function rTokenAddress() view returns(address)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) RTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "rTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RTokenAddress is a free data retrieval call binding the contract method 0x35381358.
//
// Solidity: function rTokenAddress() view returns(address)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) RTokenAddress() (common.Address, error) {
	return _StakeERC20PortalWithRate.Contract.RTokenAddress(&_StakeERC20PortalWithRate.CallOpts)
}

// RTokenAddress is a free data retrieval call binding the contract method 0x35381358.
//
// Solidity: function rTokenAddress() view returns(address)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) RTokenAddress() (common.Address, error) {
	return _StakeERC20PortalWithRate.Contract.RTokenAddress(&_StakeERC20PortalWithRate.CallOpts)
}

// RateChangeLimit is a free data retrieval call binding the contract method 0xc0152b71.
//
// Solidity: function rateChangeLimit() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) RateChangeLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "rateChangeLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RateChangeLimit is a free data retrieval call binding the contract method 0xc0152b71.
//
// Solidity: function rateChangeLimit() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) RateChangeLimit() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.RateChangeLimit(&_StakeERC20PortalWithRate.CallOpts)
}

// RateChangeLimit is a free data retrieval call binding the contract method 0xc0152b71.
//
// Solidity: function rateChangeLimit() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) RateChangeLimit() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.RateChangeLimit(&_StakeERC20PortalWithRate.CallOpts)
}

// StakePoolAddressExist is a free data retrieval call binding the contract method 0x3a68e4b2.
//
// Solidity: function stakePoolAddressExist(address ) view returns(bool)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) StakePoolAddressExist(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "stakePoolAddressExist", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StakePoolAddressExist is a free data retrieval call binding the contract method 0x3a68e4b2.
//
// Solidity: function stakePoolAddressExist(address ) view returns(bool)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) StakePoolAddressExist(arg0 common.Address) (bool, error) {
	return _StakeERC20PortalWithRate.Contract.StakePoolAddressExist(&_StakeERC20PortalWithRate.CallOpts, arg0)
}

// StakePoolAddressExist is a free data retrieval call binding the contract method 0x3a68e4b2.
//
// Solidity: function stakePoolAddressExist(address ) view returns(bool)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) StakePoolAddressExist(arg0 common.Address) (bool, error) {
	return _StakeERC20PortalWithRate.Contract.StakePoolAddressExist(&_StakeERC20PortalWithRate.CallOpts, arg0)
}

// StakeRelayFee is a free data retrieval call binding the contract method 0xfb2bc539.
//
// Solidity: function stakeRelayFee() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) StakeRelayFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "stakeRelayFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakeRelayFee is a free data retrieval call binding the contract method 0xfb2bc539.
//
// Solidity: function stakeRelayFee() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) StakeRelayFee() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.StakeRelayFee(&_StakeERC20PortalWithRate.CallOpts)
}

// StakeRelayFee is a free data retrieval call binding the contract method 0xfb2bc539.
//
// Solidity: function stakeRelayFee() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) StakeRelayFee() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.StakeRelayFee(&_StakeERC20PortalWithRate.CallOpts)
}

// StakeSwitch is a free data retrieval call binding the contract method 0x24387cdc.
//
// Solidity: function stakeSwitch() view returns(bool)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) StakeSwitch(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "stakeSwitch")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StakeSwitch is a free data retrieval call binding the contract method 0x24387cdc.
//
// Solidity: function stakeSwitch() view returns(bool)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) StakeSwitch() (bool, error) {
	return _StakeERC20PortalWithRate.Contract.StakeSwitch(&_StakeERC20PortalWithRate.CallOpts)
}

// StakeSwitch is a free data retrieval call binding the contract method 0x24387cdc.
//
// Solidity: function stakeSwitch() view returns(bool)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) StakeSwitch() (bool, error) {
	return _StakeERC20PortalWithRate.Contract.StakeSwitch(&_StakeERC20PortalWithRate.CallOpts)
}

// StakeUsePoolAddress is a free data retrieval call binding the contract method 0xed98f4bf.
//
// Solidity: function stakeUsePoolAddress() view returns(address)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) StakeUsePoolAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "stakeUsePoolAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeUsePoolAddress is a free data retrieval call binding the contract method 0xed98f4bf.
//
// Solidity: function stakeUsePoolAddress() view returns(address)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) StakeUsePoolAddress() (common.Address, error) {
	return _StakeERC20PortalWithRate.Contract.StakeUsePoolAddress(&_StakeERC20PortalWithRate.CallOpts)
}

// StakeUsePoolAddress is a free data retrieval call binding the contract method 0xed98f4bf.
//
// Solidity: function stakeUsePoolAddress() view returns(address)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) StakeUsePoolAddress() (common.Address, error) {
	return _StakeERC20PortalWithRate.Contract.StakeUsePoolAddress(&_StakeERC20PortalWithRate.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) Threshold(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "threshold")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) Threshold() (uint8, error) {
	return _StakeERC20PortalWithRate.Contract.Threshold(&_StakeERC20PortalWithRate.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) Threshold() (uint8, error) {
	return _StakeERC20PortalWithRate.Contract.Threshold(&_StakeERC20PortalWithRate.CallOpts)
}

// TotalUnstakeProtocolFee is a free data retrieval call binding the contract method 0x050f541a.
//
// Solidity: function totalUnstakeProtocolFee() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) TotalUnstakeProtocolFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "totalUnstakeProtocolFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalUnstakeProtocolFee is a free data retrieval call binding the contract method 0x050f541a.
//
// Solidity: function totalUnstakeProtocolFee() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) TotalUnstakeProtocolFee() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.TotalUnstakeProtocolFee(&_StakeERC20PortalWithRate.CallOpts)
}

// TotalUnstakeProtocolFee is a free data retrieval call binding the contract method 0x050f541a.
//
// Solidity: function totalUnstakeProtocolFee() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) TotalUnstakeProtocolFee() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.TotalUnstakeProtocolFee(&_StakeERC20PortalWithRate.CallOpts)
}

// UnstakeFeeCommission is a free data retrieval call binding the contract method 0xe1bb5897.
//
// Solidity: function unstakeFeeCommission() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) UnstakeFeeCommission(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "unstakeFeeCommission")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakeFeeCommission is a free data retrieval call binding the contract method 0xe1bb5897.
//
// Solidity: function unstakeFeeCommission() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) UnstakeFeeCommission() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.UnstakeFeeCommission(&_StakeERC20PortalWithRate.CallOpts)
}

// UnstakeFeeCommission is a free data retrieval call binding the contract method 0xe1bb5897.
//
// Solidity: function unstakeFeeCommission() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) UnstakeFeeCommission() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.UnstakeFeeCommission(&_StakeERC20PortalWithRate.CallOpts)
}

// UnstakeRelayFee is a free data retrieval call binding the contract method 0xf94bf9c4.
//
// Solidity: function unstakeRelayFee() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCaller) UnstakeRelayFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20PortalWithRate.contract.Call(opts, &out, "unstakeRelayFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakeRelayFee is a free data retrieval call binding the contract method 0xf94bf9c4.
//
// Solidity: function unstakeRelayFee() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) UnstakeRelayFee() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.UnstakeRelayFee(&_StakeERC20PortalWithRate.CallOpts)
}

// UnstakeRelayFee is a free data retrieval call binding the contract method 0xf94bf9c4.
//
// Solidity: function unstakeRelayFee() view returns(uint256)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateCallerSession) UnstakeRelayFee() (*big.Int, error) {
	return _StakeERC20PortalWithRate.Contract.UnstakeRelayFee(&_StakeERC20PortalWithRate.CallOpts)
}

// AddStakePool is a paid mutator transaction binding the contract method 0xcb08cbae.
//
// Solidity: function addStakePool(address[] _stakePoolAddressList) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) AddStakePool(opts *bind.TransactOpts, _stakePoolAddressList []common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "addStakePool", _stakePoolAddressList)
}

// AddStakePool is a paid mutator transaction binding the contract method 0xcb08cbae.
//
// Solidity: function addStakePool(address[] _stakePoolAddressList) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) AddStakePool(_stakePoolAddressList []common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.AddStakePool(&_StakeERC20PortalWithRate.TransactOpts, _stakePoolAddressList)
}

// AddStakePool is a paid mutator transaction binding the contract method 0xcb08cbae.
//
// Solidity: function addStakePool(address[] _stakePoolAddressList) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) AddStakePool(_stakePoolAddressList []common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.AddStakePool(&_StakeERC20PortalWithRate.TransactOpts, _stakePoolAddressList)
}

// AddSubAccount is a paid mutator transaction binding the contract method 0xb52e1eff.
//
// Solidity: function addSubAccount(address _subAccount) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) AddSubAccount(opts *bind.TransactOpts, _subAccount common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "addSubAccount", _subAccount)
}

// AddSubAccount is a paid mutator transaction binding the contract method 0xb52e1eff.
//
// Solidity: function addSubAccount(address _subAccount) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) AddSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.AddSubAccount(&_StakeERC20PortalWithRate.TransactOpts, _subAccount)
}

// AddSubAccount is a paid mutator transaction binding the contract method 0xb52e1eff.
//
// Solidity: function addSubAccount(address _subAccount) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) AddSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.AddSubAccount(&_StakeERC20PortalWithRate.TransactOpts, _subAccount)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _newThreshold) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) ChangeThreshold(opts *bind.TransactOpts, _newThreshold *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "changeThreshold", _newThreshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _newThreshold) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) ChangeThreshold(_newThreshold *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.ChangeThreshold(&_StakeERC20PortalWithRate.TransactOpts, _newThreshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _newThreshold) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) ChangeThreshold(_newThreshold *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.ChangeThreshold(&_StakeERC20PortalWithRate.TransactOpts, _newThreshold)
}

// RemoveSubAccount is a paid mutator transaction binding the contract method 0x0ba3aae3.
//
// Solidity: function removeSubAccount(address _subAccount) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) RemoveSubAccount(opts *bind.TransactOpts, _subAccount common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "removeSubAccount", _subAccount)
}

// RemoveSubAccount is a paid mutator transaction binding the contract method 0x0ba3aae3.
//
// Solidity: function removeSubAccount(address _subAccount) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) RemoveSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.RemoveSubAccount(&_StakeERC20PortalWithRate.TransactOpts, _subAccount)
}

// RemoveSubAccount is a paid mutator transaction binding the contract method 0x0ba3aae3.
//
// Solidity: function removeSubAccount(address _subAccount) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) RemoveSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.RemoveSubAccount(&_StakeERC20PortalWithRate.TransactOpts, _subAccount)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _stakePoolAddress) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) RmStakePool(opts *bind.TransactOpts, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "rmStakePool", _stakePoolAddress)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _stakePoolAddress) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) RmStakePool(_stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.RmStakePool(&_StakeERC20PortalWithRate.TransactOpts, _stakePoolAddress)
}

// RmStakePool is a paid mutator transaction binding the contract method 0x58c76ca8.
//
// Solidity: function rmStakePool(address _stakePoolAddress) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) RmStakePool(_stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.RmStakePool(&_StakeERC20PortalWithRate.TransactOpts, _stakePoolAddress)
}

// SetMinStakeAmount is a paid mutator transaction binding the contract method 0xeb4af045.
//
// Solidity: function setMinStakeAmount(uint256 _minStakeAmount) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) SetMinStakeAmount(opts *bind.TransactOpts, _minStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "setMinStakeAmount", _minStakeAmount)
}

// SetMinStakeAmount is a paid mutator transaction binding the contract method 0xeb4af045.
//
// Solidity: function setMinStakeAmount(uint256 _minStakeAmount) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) SetMinStakeAmount(_minStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetMinStakeAmount(&_StakeERC20PortalWithRate.TransactOpts, _minStakeAmount)
}

// SetMinStakeAmount is a paid mutator transaction binding the contract method 0xeb4af045.
//
// Solidity: function setMinStakeAmount(uint256 _minStakeAmount) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) SetMinStakeAmount(_minStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetMinStakeAmount(&_StakeERC20PortalWithRate.TransactOpts, _minStakeAmount)
}

// SetRate is a paid mutator transaction binding the contract method 0x34fcf437.
//
// Solidity: function setRate(uint256 _rate) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) SetRate(opts *bind.TransactOpts, _rate *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "setRate", _rate)
}

// SetRate is a paid mutator transaction binding the contract method 0x34fcf437.
//
// Solidity: function setRate(uint256 _rate) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) SetRate(_rate *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetRate(&_StakeERC20PortalWithRate.TransactOpts, _rate)
}

// SetRate is a paid mutator transaction binding the contract method 0x34fcf437.
//
// Solidity: function setRate(uint256 _rate) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) SetRate(_rate *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetRate(&_StakeERC20PortalWithRate.TransactOpts, _rate)
}

// SetRateChangeLimit is a paid mutator transaction binding the contract method 0x19826e71.
//
// Solidity: function setRateChangeLimit(uint256 _rateChangeLimit) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) SetRateChangeLimit(opts *bind.TransactOpts, _rateChangeLimit *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "setRateChangeLimit", _rateChangeLimit)
}

// SetRateChangeLimit is a paid mutator transaction binding the contract method 0x19826e71.
//
// Solidity: function setRateChangeLimit(uint256 _rateChangeLimit) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) SetRateChangeLimit(_rateChangeLimit *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetRateChangeLimit(&_StakeERC20PortalWithRate.TransactOpts, _rateChangeLimit)
}

// SetRateChangeLimit is a paid mutator transaction binding the contract method 0x19826e71.
//
// Solidity: function setRateChangeLimit(uint256 _rateChangeLimit) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) SetRateChangeLimit(_rateChangeLimit *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetRateChangeLimit(&_StakeERC20PortalWithRate.TransactOpts, _rateChangeLimit)
}

// SetStakeRelayFee is a paid mutator transaction binding the contract method 0x31db4715.
//
// Solidity: function setStakeRelayFee(uint256 _stakeRelayFee) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) SetStakeRelayFee(opts *bind.TransactOpts, _stakeRelayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "setStakeRelayFee", _stakeRelayFee)
}

// SetStakeRelayFee is a paid mutator transaction binding the contract method 0x31db4715.
//
// Solidity: function setStakeRelayFee(uint256 _stakeRelayFee) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) SetStakeRelayFee(_stakeRelayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetStakeRelayFee(&_StakeERC20PortalWithRate.TransactOpts, _stakeRelayFee)
}

// SetStakeRelayFee is a paid mutator transaction binding the contract method 0x31db4715.
//
// Solidity: function setStakeRelayFee(uint256 _stakeRelayFee) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) SetStakeRelayFee(_stakeRelayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetStakeRelayFee(&_StakeERC20PortalWithRate.TransactOpts, _stakeRelayFee)
}

// SetStakeUsePool is a paid mutator transaction binding the contract method 0xe52a0d92.
//
// Solidity: function setStakeUsePool(address _stakeUsePoolAddress) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) SetStakeUsePool(opts *bind.TransactOpts, _stakeUsePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "setStakeUsePool", _stakeUsePoolAddress)
}

// SetStakeUsePool is a paid mutator transaction binding the contract method 0xe52a0d92.
//
// Solidity: function setStakeUsePool(address _stakeUsePoolAddress) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) SetStakeUsePool(_stakeUsePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetStakeUsePool(&_StakeERC20PortalWithRate.TransactOpts, _stakeUsePoolAddress)
}

// SetStakeUsePool is a paid mutator transaction binding the contract method 0xe52a0d92.
//
// Solidity: function setStakeUsePool(address _stakeUsePoolAddress) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) SetStakeUsePool(_stakeUsePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetStakeUsePool(&_StakeERC20PortalWithRate.TransactOpts, _stakeUsePoolAddress)
}

// SetUnstakeFeeCommission is a paid mutator transaction binding the contract method 0x3fbd062b.
//
// Solidity: function setUnstakeFeeCommission(uint256 _unstakeFeeCommission) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) SetUnstakeFeeCommission(opts *bind.TransactOpts, _unstakeFeeCommission *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "setUnstakeFeeCommission", _unstakeFeeCommission)
}

// SetUnstakeFeeCommission is a paid mutator transaction binding the contract method 0x3fbd062b.
//
// Solidity: function setUnstakeFeeCommission(uint256 _unstakeFeeCommission) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) SetUnstakeFeeCommission(_unstakeFeeCommission *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetUnstakeFeeCommission(&_StakeERC20PortalWithRate.TransactOpts, _unstakeFeeCommission)
}

// SetUnstakeFeeCommission is a paid mutator transaction binding the contract method 0x3fbd062b.
//
// Solidity: function setUnstakeFeeCommission(uint256 _unstakeFeeCommission) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) SetUnstakeFeeCommission(_unstakeFeeCommission *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetUnstakeFeeCommission(&_StakeERC20PortalWithRate.TransactOpts, _unstakeFeeCommission)
}

// SetUnstakeRelayFee is a paid mutator transaction binding the contract method 0xeb7f83d9.
//
// Solidity: function setUnstakeRelayFee(uint256 _unstakeRelayFee) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) SetUnstakeRelayFee(opts *bind.TransactOpts, _unstakeRelayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "setUnstakeRelayFee", _unstakeRelayFee)
}

// SetUnstakeRelayFee is a paid mutator transaction binding the contract method 0xeb7f83d9.
//
// Solidity: function setUnstakeRelayFee(uint256 _unstakeRelayFee) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) SetUnstakeRelayFee(_unstakeRelayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetUnstakeRelayFee(&_StakeERC20PortalWithRate.TransactOpts, _unstakeRelayFee)
}

// SetUnstakeRelayFee is a paid mutator transaction binding the contract method 0xeb7f83d9.
//
// Solidity: function setUnstakeRelayFee(uint256 _unstakeRelayFee) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) SetUnstakeRelayFee(_unstakeRelayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.SetUnstakeRelayFee(&_StakeERC20PortalWithRate.TransactOpts, _unstakeRelayFee)
}

// Stake is a paid mutator transaction binding the contract method 0xa694fc3a.
//
// Solidity: function stake(uint256 _amount) payable returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) Stake(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "stake", _amount)
}

// Stake is a paid mutator transaction binding the contract method 0xa694fc3a.
//
// Solidity: function stake(uint256 _amount) payable returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) Stake(_amount *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.Stake(&_StakeERC20PortalWithRate.TransactOpts, _amount)
}

// Stake is a paid mutator transaction binding the contract method 0xa694fc3a.
//
// Solidity: function stake(uint256 _amount) payable returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) Stake(_amount *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.Stake(&_StakeERC20PortalWithRate.TransactOpts, _amount)
}

// StakeWithPool is a paid mutator transaction binding the contract method 0x227cdfb7.
//
// Solidity: function stakeWithPool(uint256 _amount, address _stakePoolAddress) payable returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) StakeWithPool(opts *bind.TransactOpts, _amount *big.Int, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "stakeWithPool", _amount, _stakePoolAddress)
}

// StakeWithPool is a paid mutator transaction binding the contract method 0x227cdfb7.
//
// Solidity: function stakeWithPool(uint256 _amount, address _stakePoolAddress) payable returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) StakeWithPool(_amount *big.Int, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.StakeWithPool(&_StakeERC20PortalWithRate.TransactOpts, _amount, _stakePoolAddress)
}

// StakeWithPool is a paid mutator transaction binding the contract method 0x227cdfb7.
//
// Solidity: function stakeWithPool(uint256 _amount, address _stakePoolAddress) payable returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) StakeWithPool(_amount *big.Int, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.StakeWithPool(&_StakeERC20PortalWithRate.TransactOpts, _amount, _stakePoolAddress)
}

// ToggleStakeSwitch is a paid mutator transaction binding the contract method 0x995bb87d.
//
// Solidity: function toggleStakeSwitch() returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) ToggleStakeSwitch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "toggleStakeSwitch")
}

// ToggleStakeSwitch is a paid mutator transaction binding the contract method 0x995bb87d.
//
// Solidity: function toggleStakeSwitch() returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) ToggleStakeSwitch() (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.ToggleStakeSwitch(&_StakeERC20PortalWithRate.TransactOpts)
}

// ToggleStakeSwitch is a paid mutator transaction binding the contract method 0x995bb87d.
//
// Solidity: function toggleStakeSwitch() returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) ToggleStakeSwitch() (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.ToggleStakeSwitch(&_StakeERC20PortalWithRate.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.TransferOwnership(&_StakeERC20PortalWithRate.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address _newOwner) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.TransferOwnership(&_StakeERC20PortalWithRate.TransactOpts, _newOwner)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 _rTokenAmount) payable returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) Unstake(opts *bind.TransactOpts, _rTokenAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "unstake", _rTokenAmount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 _rTokenAmount) payable returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) Unstake(_rTokenAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.Unstake(&_StakeERC20PortalWithRate.TransactOpts, _rTokenAmount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 _rTokenAmount) payable returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) Unstake(_rTokenAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.Unstake(&_StakeERC20PortalWithRate.TransactOpts, _rTokenAmount)
}

// UnstakeWithPool is a paid mutator transaction binding the contract method 0xac8ef2db.
//
// Solidity: function unstakeWithPool(uint256 _rTokenAmount, address _stakePoolAddress) payable returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) UnstakeWithPool(opts *bind.TransactOpts, _rTokenAmount *big.Int, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "unstakeWithPool", _rTokenAmount, _stakePoolAddress)
}

// UnstakeWithPool is a paid mutator transaction binding the contract method 0xac8ef2db.
//
// Solidity: function unstakeWithPool(uint256 _rTokenAmount, address _stakePoolAddress) payable returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) UnstakeWithPool(_rTokenAmount *big.Int, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.UnstakeWithPool(&_StakeERC20PortalWithRate.TransactOpts, _rTokenAmount, _stakePoolAddress)
}

// UnstakeWithPool is a paid mutator transaction binding the contract method 0xac8ef2db.
//
// Solidity: function unstakeWithPool(uint256 _rTokenAmount, address _stakePoolAddress) payable returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) UnstakeWithPool(_rTokenAmount *big.Int, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.UnstakeWithPool(&_StakeERC20PortalWithRate.TransactOpts, _rTokenAmount, _stakePoolAddress)
}

// VoteRate is a paid mutator transaction binding the contract method 0x113e3709.
//
// Solidity: function voteRate(bytes32 _proposalId, uint256 _rate) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactor) VoteRate(opts *bind.TransactOpts, _proposalId [32]byte, _rate *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.contract.Transact(opts, "voteRate", _proposalId, _rate)
}

// VoteRate is a paid mutator transaction binding the contract method 0x113e3709.
//
// Solidity: function voteRate(bytes32 _proposalId, uint256 _rate) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateSession) VoteRate(_proposalId [32]byte, _rate *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.VoteRate(&_StakeERC20PortalWithRate.TransactOpts, _proposalId, _rate)
}

// VoteRate is a paid mutator transaction binding the contract method 0x113e3709.
//
// Solidity: function voteRate(bytes32 _proposalId, uint256 _rate) returns()
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateTransactorSession) VoteRate(_proposalId [32]byte, _rate *big.Int) (*types.Transaction, error) {
	return _StakeERC20PortalWithRate.Contract.VoteRate(&_StakeERC20PortalWithRate.TransactOpts, _proposalId, _rate)
}

// StakeERC20PortalWithRateProposalExecutedIterator is returned from FilterProposalExecuted and is used to iterate over the raw logs and unpacked data for ProposalExecuted events raised by the StakeERC20PortalWithRate contract.
type StakeERC20PortalWithRateProposalExecutedIterator struct {
	Event *StakeERC20PortalWithRateProposalExecuted // Event containing the contract specifics and raw log

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
func (it *StakeERC20PortalWithRateProposalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeERC20PortalWithRateProposalExecuted)
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
		it.Event = new(StakeERC20PortalWithRateProposalExecuted)
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
func (it *StakeERC20PortalWithRateProposalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeERC20PortalWithRateProposalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeERC20PortalWithRateProposalExecuted represents a ProposalExecuted event raised by the StakeERC20PortalWithRate contract.
type StakeERC20PortalWithRateProposalExecuted struct {
	ProposalId [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalExecuted is a free log retrieval operation binding the contract event 0x7b1bcf1ccf901a11589afff5504d59fd0a53780eed2a952adade0348985139e0.
//
// Solidity: event ProposalExecuted(bytes32 indexed proposalId)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateFilterer) FilterProposalExecuted(opts *bind.FilterOpts, proposalId [][32]byte) (*StakeERC20PortalWithRateProposalExecutedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _StakeERC20PortalWithRate.contract.FilterLogs(opts, "ProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalWithRateProposalExecutedIterator{contract: _StakeERC20PortalWithRate.contract, event: "ProposalExecuted", logs: logs, sub: sub}, nil
}

// WatchProposalExecuted is a free log subscription operation binding the contract event 0x7b1bcf1ccf901a11589afff5504d59fd0a53780eed2a952adade0348985139e0.
//
// Solidity: event ProposalExecuted(bytes32 indexed proposalId)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateFilterer) WatchProposalExecuted(opts *bind.WatchOpts, sink chan<- *StakeERC20PortalWithRateProposalExecuted, proposalId [][32]byte) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _StakeERC20PortalWithRate.contract.WatchLogs(opts, "ProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeERC20PortalWithRateProposalExecuted)
				if err := _StakeERC20PortalWithRate.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
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
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateFilterer) ParseProposalExecuted(log types.Log) (*StakeERC20PortalWithRateProposalExecuted, error) {
	event := new(StakeERC20PortalWithRateProposalExecuted)
	if err := _StakeERC20PortalWithRate.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeERC20PortalWithRateStakeIterator is returned from FilterStake and is used to iterate over the raw logs and unpacked data for Stake events raised by the StakeERC20PortalWithRate contract.
type StakeERC20PortalWithRateStakeIterator struct {
	Event *StakeERC20PortalWithRateStake // Event containing the contract specifics and raw log

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
func (it *StakeERC20PortalWithRateStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeERC20PortalWithRateStake)
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
		it.Event = new(StakeERC20PortalWithRateStake)
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
func (it *StakeERC20PortalWithRateStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeERC20PortalWithRateStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeERC20PortalWithRateStake represents a Stake event raised by the StakeERC20PortalWithRate contract.
type StakeERC20PortalWithRateStake struct {
	Staker       common.Address
	StakePool    common.Address
	TokenAmount  *big.Int
	RTokenAmount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStake is a free log retrieval operation binding the contract event 0x63602d0ecc7b3a0ef7ff1a116e23056662d64280355ba8031b6d0d767c4b4458.
//
// Solidity: event Stake(address staker, address stakePool, uint256 tokenAmount, uint256 rTokenAmount)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateFilterer) FilterStake(opts *bind.FilterOpts) (*StakeERC20PortalWithRateStakeIterator, error) {

	logs, sub, err := _StakeERC20PortalWithRate.contract.FilterLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalWithRateStakeIterator{contract: _StakeERC20PortalWithRate.contract, event: "Stake", logs: logs, sub: sub}, nil
}

// WatchStake is a free log subscription operation binding the contract event 0x63602d0ecc7b3a0ef7ff1a116e23056662d64280355ba8031b6d0d767c4b4458.
//
// Solidity: event Stake(address staker, address stakePool, uint256 tokenAmount, uint256 rTokenAmount)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateFilterer) WatchStake(opts *bind.WatchOpts, sink chan<- *StakeERC20PortalWithRateStake) (event.Subscription, error) {

	logs, sub, err := _StakeERC20PortalWithRate.contract.WatchLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeERC20PortalWithRateStake)
				if err := _StakeERC20PortalWithRate.contract.UnpackLog(event, "Stake", log); err != nil {
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
// Solidity: event Stake(address staker, address stakePool, uint256 tokenAmount, uint256 rTokenAmount)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateFilterer) ParseStake(log types.Log) (*StakeERC20PortalWithRateStake, error) {
	event := new(StakeERC20PortalWithRateStake)
	if err := _StakeERC20PortalWithRate.contract.UnpackLog(event, "Stake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeERC20PortalWithRateUnstakeIterator is returned from FilterUnstake and is used to iterate over the raw logs and unpacked data for Unstake events raised by the StakeERC20PortalWithRate contract.
type StakeERC20PortalWithRateUnstakeIterator struct {
	Event *StakeERC20PortalWithRateUnstake // Event containing the contract specifics and raw log

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
func (it *StakeERC20PortalWithRateUnstakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeERC20PortalWithRateUnstake)
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
		it.Event = new(StakeERC20PortalWithRateUnstake)
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
func (it *StakeERC20PortalWithRateUnstakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeERC20PortalWithRateUnstakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeERC20PortalWithRateUnstake represents a Unstake event raised by the StakeERC20PortalWithRate contract.
type StakeERC20PortalWithRateUnstake struct {
	Staker       common.Address
	StakePool    common.Address
	TokenAmount  *big.Int
	RTokenAmount *big.Int
	BurnAmount   *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUnstake is a free log retrieval operation binding the contract event 0xfe7007b2e89d80edda76299251df08366480cac22e5e260f5e662e850b1f7a32.
//
// Solidity: event Unstake(address staker, address stakePool, uint256 tokenAmount, uint256 rTokenAmount, uint256 burnAmount)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateFilterer) FilterUnstake(opts *bind.FilterOpts) (*StakeERC20PortalWithRateUnstakeIterator, error) {

	logs, sub, err := _StakeERC20PortalWithRate.contract.FilterLogs(opts, "Unstake")
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalWithRateUnstakeIterator{contract: _StakeERC20PortalWithRate.contract, event: "Unstake", logs: logs, sub: sub}, nil
}

// WatchUnstake is a free log subscription operation binding the contract event 0xfe7007b2e89d80edda76299251df08366480cac22e5e260f5e662e850b1f7a32.
//
// Solidity: event Unstake(address staker, address stakePool, uint256 tokenAmount, uint256 rTokenAmount, uint256 burnAmount)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateFilterer) WatchUnstake(opts *bind.WatchOpts, sink chan<- *StakeERC20PortalWithRateUnstake) (event.Subscription, error) {

	logs, sub, err := _StakeERC20PortalWithRate.contract.WatchLogs(opts, "Unstake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeERC20PortalWithRateUnstake)
				if err := _StakeERC20PortalWithRate.contract.UnpackLog(event, "Unstake", log); err != nil {
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

// ParseUnstake is a log parse operation binding the contract event 0xfe7007b2e89d80edda76299251df08366480cac22e5e260f5e662e850b1f7a32.
//
// Solidity: event Unstake(address staker, address stakePool, uint256 tokenAmount, uint256 rTokenAmount, uint256 burnAmount)
func (_StakeERC20PortalWithRate *StakeERC20PortalWithRateFilterer) ParseUnstake(log types.Log) (*StakeERC20PortalWithRateUnstake, error) {
	event := new(StakeERC20PortalWithRateUnstake)
	if err := _StakeERC20PortalWithRate.contract.UnpackLog(event, "Unstake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
