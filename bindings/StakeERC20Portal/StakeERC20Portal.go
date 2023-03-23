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
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_stakePoolAddressList\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_initialSubAccounts\",\"type\":\"address[]\"},{\"internalType\":\"uint8[]\",\"name\":\"_chainIdList\",\"type\":\"uint8[]\"},{\"internalType\":\"address\",\"name\":\"_erc20TokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_rTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_stakeUsePoolAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minStakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_stakeRelayFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_unstakeRelayFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_unstakeFeeCommission\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_initialThreshold\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"proposalId\",\"type\":\"bytes32\"}],\"name\":\"ProposalExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stafiRecipient\",\"type\":\"bytes32\"}],\"name\":\"RecoverStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stakePool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rTokenAmount\",\"type\":\"uint256\"}],\"name\":\"Stake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stakePool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"chainId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stafiRecipient\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"destRecipient\",\"type\":\"address\"}],\"name\":\"StakeAndCross\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"stakePool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rTokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"burnAmount\",\"type\":\"uint256\"}],\"name\":\"Unstake\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint8[]\",\"name\":\"_chaindIdList\",\"type\":\"uint8[]\"}],\"name\":\"addChainId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_stakePoolAddressList\",\"type\":\"address[]\"}],\"name\":\"addStakePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"addSubAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"bridgeFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"name\":\"chainIdExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newThreshold\",\"type\":\"uint256\"}],\"name\":\"changeThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"erc20TokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"getSubAccountIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_proposalId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"hasVoted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"proposals\",\"outputs\":[{\"internalType\":\"enumMultisig.ProposalStatus\",\"name\":\"_status\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"_yesVotes\",\"type\":\"uint16\"},{\"internalType\":\"uint8\",\"name\":\"_yesVotesTotal\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rateChangeLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_txHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_stafiRecipient\",\"type\":\"bytes32\"}],\"name\":\"recoverStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subAccount\",\"type\":\"address\"}],\"name\":\"removeSubAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_chaindId\",\"type\":\"uint8\"}],\"name\":\"rmChainId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"}],\"name\":\"rmStakePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_chainId\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_bridgeFee\",\"type\":\"uint256\"}],\"name\":\"setBridgeFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minStakeAmount\",\"type\":\"uint256\"}],\"name\":\"setMinStakeAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"setRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rateChangeLimit\",\"type\":\"uint256\"}],\"name\":\"setRateChangeLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_stakeRelayFee\",\"type\":\"uint256\"}],\"name\":\"setStakeRelayFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakeUsePoolAddress\",\"type\":\"address\"}],\"name\":\"setStakeUsePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_unstakeFeeCommission\",\"type\":\"uint256\"}],\"name\":\"setUnstakeFeeCommission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_unstakeRelayFee\",\"type\":\"uint256\"}],\"name\":\"setUnstakeRelayFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"_destChainId\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_stafiRecipient\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_destRecipient\",\"type\":\"address\"}],\"name\":\"stakeAndCross\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeCrossSwitch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"stakePoolAddressExist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeRelayFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeSwitch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeUsePoolAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"}],\"name\":\"stakeWithPool\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"toggleStakeCrossSwitch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"toggleStakeSwitch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalUnstakeProtocolFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rTokenAmount\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unstakeFeeCommission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unstakeRelayFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_stakePoolAddress\",\"type\":\"address\"}],\"name\":\"unstakeWithPool\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_proposalId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"voteRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCaller) GetRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "getRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalSession) GetRate() (*big.Int, error) {
	return _StakeERC20Portal.Contract.GetRate(&_StakeERC20Portal.CallOpts)
}

// GetRate is a free data retrieval call binding the contract method 0x679aefce.
//
// Solidity: function getRate() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) GetRate() (*big.Int, error) {
	return _StakeERC20Portal.Contract.GetRate(&_StakeERC20Portal.CallOpts)
}

// GetSubAccountIndex is a free data retrieval call binding the contract method 0x763f8680.
//
// Solidity: function getSubAccountIndex(address _subAccount) view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCaller) GetSubAccountIndex(opts *bind.CallOpts, _subAccount common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "getSubAccountIndex", _subAccount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSubAccountIndex is a free data retrieval call binding the contract method 0x763f8680.
//
// Solidity: function getSubAccountIndex(address _subAccount) view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalSession) GetSubAccountIndex(_subAccount common.Address) (*big.Int, error) {
	return _StakeERC20Portal.Contract.GetSubAccountIndex(&_StakeERC20Portal.CallOpts, _subAccount)
}

// GetSubAccountIndex is a free data retrieval call binding the contract method 0x763f8680.
//
// Solidity: function getSubAccountIndex(address _subAccount) view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) GetSubAccountIndex(_subAccount common.Address) (*big.Int, error) {
	return _StakeERC20Portal.Contract.GetSubAccountIndex(&_StakeERC20Portal.CallOpts, _subAccount)
}

// HasVoted is a free data retrieval call binding the contract method 0xaadc3b72.
//
// Solidity: function hasVoted(bytes32 _proposalId, address _subAccount) view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalCaller) HasVoted(opts *bind.CallOpts, _proposalId [32]byte, _subAccount common.Address) (bool, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "hasVoted", _proposalId, _subAccount)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasVoted is a free data retrieval call binding the contract method 0xaadc3b72.
//
// Solidity: function hasVoted(bytes32 _proposalId, address _subAccount) view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalSession) HasVoted(_proposalId [32]byte, _subAccount common.Address) (bool, error) {
	return _StakeERC20Portal.Contract.HasVoted(&_StakeERC20Portal.CallOpts, _proposalId, _subAccount)
}

// HasVoted is a free data retrieval call binding the contract method 0xaadc3b72.
//
// Solidity: function hasVoted(bytes32 _proposalId, address _subAccount) view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) HasVoted(_proposalId [32]byte, _subAccount common.Address) (bool, error) {
	return _StakeERC20Portal.Contract.HasVoted(&_StakeERC20Portal.CallOpts, _proposalId, _subAccount)
}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCaller) MinStakeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "minStakeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalSession) MinStakeAmount() (*big.Int, error) {
	return _StakeERC20Portal.Contract.MinStakeAmount(&_StakeERC20Portal.CallOpts)
}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) MinStakeAmount() (*big.Int, error) {
	return _StakeERC20Portal.Contract.MinStakeAmount(&_StakeERC20Portal.CallOpts)
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

// Proposals is a free data retrieval call binding the contract method 0x32ed5b12.
//
// Solidity: function proposals(bytes32 ) view returns(uint8 _status, uint16 _yesVotes, uint8 _yesVotesTotal)
func (_StakeERC20Portal *StakeERC20PortalCaller) Proposals(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Status        uint8
	YesVotes      uint16
	YesVotesTotal uint8
}, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "proposals", arg0)

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
func (_StakeERC20Portal *StakeERC20PortalSession) Proposals(arg0 [32]byte) (struct {
	Status        uint8
	YesVotes      uint16
	YesVotesTotal uint8
}, error) {
	return _StakeERC20Portal.Contract.Proposals(&_StakeERC20Portal.CallOpts, arg0)
}

// Proposals is a free data retrieval call binding the contract method 0x32ed5b12.
//
// Solidity: function proposals(bytes32 ) view returns(uint8 _status, uint16 _yesVotes, uint8 _yesVotesTotal)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) Proposals(arg0 [32]byte) (struct {
	Status        uint8
	YesVotes      uint16
	YesVotesTotal uint8
}, error) {
	return _StakeERC20Portal.Contract.Proposals(&_StakeERC20Portal.CallOpts, arg0)
}

// RTokenAddress is a free data retrieval call binding the contract method 0x35381358.
//
// Solidity: function rTokenAddress() view returns(address)
func (_StakeERC20Portal *StakeERC20PortalCaller) RTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "rTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RTokenAddress is a free data retrieval call binding the contract method 0x35381358.
//
// Solidity: function rTokenAddress() view returns(address)
func (_StakeERC20Portal *StakeERC20PortalSession) RTokenAddress() (common.Address, error) {
	return _StakeERC20Portal.Contract.RTokenAddress(&_StakeERC20Portal.CallOpts)
}

// RTokenAddress is a free data retrieval call binding the contract method 0x35381358.
//
// Solidity: function rTokenAddress() view returns(address)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) RTokenAddress() (common.Address, error) {
	return _StakeERC20Portal.Contract.RTokenAddress(&_StakeERC20Portal.CallOpts)
}

// RateChangeLimit is a free data retrieval call binding the contract method 0xc0152b71.
//
// Solidity: function rateChangeLimit() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCaller) RateChangeLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "rateChangeLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RateChangeLimit is a free data retrieval call binding the contract method 0xc0152b71.
//
// Solidity: function rateChangeLimit() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalSession) RateChangeLimit() (*big.Int, error) {
	return _StakeERC20Portal.Contract.RateChangeLimit(&_StakeERC20Portal.CallOpts)
}

// RateChangeLimit is a free data retrieval call binding the contract method 0xc0152b71.
//
// Solidity: function rateChangeLimit() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) RateChangeLimit() (*big.Int, error) {
	return _StakeERC20Portal.Contract.RateChangeLimit(&_StakeERC20Portal.CallOpts)
}

// StakeCrossSwitch is a free data retrieval call binding the contract method 0x8638bfe6.
//
// Solidity: function stakeCrossSwitch() view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalCaller) StakeCrossSwitch(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "stakeCrossSwitch")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// StakeCrossSwitch is a free data retrieval call binding the contract method 0x8638bfe6.
//
// Solidity: function stakeCrossSwitch() view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalSession) StakeCrossSwitch() (bool, error) {
	return _StakeERC20Portal.Contract.StakeCrossSwitch(&_StakeERC20Portal.CallOpts)
}

// StakeCrossSwitch is a free data retrieval call binding the contract method 0x8638bfe6.
//
// Solidity: function stakeCrossSwitch() view returns(bool)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) StakeCrossSwitch() (bool, error) {
	return _StakeERC20Portal.Contract.StakeCrossSwitch(&_StakeERC20Portal.CallOpts)
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

// StakeRelayFee is a free data retrieval call binding the contract method 0xfb2bc539.
//
// Solidity: function stakeRelayFee() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCaller) StakeRelayFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "stakeRelayFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakeRelayFee is a free data retrieval call binding the contract method 0xfb2bc539.
//
// Solidity: function stakeRelayFee() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalSession) StakeRelayFee() (*big.Int, error) {
	return _StakeERC20Portal.Contract.StakeRelayFee(&_StakeERC20Portal.CallOpts)
}

// StakeRelayFee is a free data retrieval call binding the contract method 0xfb2bc539.
//
// Solidity: function stakeRelayFee() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) StakeRelayFee() (*big.Int, error) {
	return _StakeERC20Portal.Contract.StakeRelayFee(&_StakeERC20Portal.CallOpts)
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

// StakeUsePoolAddress is a free data retrieval call binding the contract method 0xed98f4bf.
//
// Solidity: function stakeUsePoolAddress() view returns(address)
func (_StakeERC20Portal *StakeERC20PortalCaller) StakeUsePoolAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "stakeUsePoolAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeUsePoolAddress is a free data retrieval call binding the contract method 0xed98f4bf.
//
// Solidity: function stakeUsePoolAddress() view returns(address)
func (_StakeERC20Portal *StakeERC20PortalSession) StakeUsePoolAddress() (common.Address, error) {
	return _StakeERC20Portal.Contract.StakeUsePoolAddress(&_StakeERC20Portal.CallOpts)
}

// StakeUsePoolAddress is a free data retrieval call binding the contract method 0xed98f4bf.
//
// Solidity: function stakeUsePoolAddress() view returns(address)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) StakeUsePoolAddress() (common.Address, error) {
	return _StakeERC20Portal.Contract.StakeUsePoolAddress(&_StakeERC20Portal.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_StakeERC20Portal *StakeERC20PortalCaller) Threshold(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "threshold")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_StakeERC20Portal *StakeERC20PortalSession) Threshold() (uint8, error) {
	return _StakeERC20Portal.Contract.Threshold(&_StakeERC20Portal.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint8)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) Threshold() (uint8, error) {
	return _StakeERC20Portal.Contract.Threshold(&_StakeERC20Portal.CallOpts)
}

// TotalUnstakeProtocolFee is a free data retrieval call binding the contract method 0x050f541a.
//
// Solidity: function totalUnstakeProtocolFee() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCaller) TotalUnstakeProtocolFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "totalUnstakeProtocolFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalUnstakeProtocolFee is a free data retrieval call binding the contract method 0x050f541a.
//
// Solidity: function totalUnstakeProtocolFee() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalSession) TotalUnstakeProtocolFee() (*big.Int, error) {
	return _StakeERC20Portal.Contract.TotalUnstakeProtocolFee(&_StakeERC20Portal.CallOpts)
}

// TotalUnstakeProtocolFee is a free data retrieval call binding the contract method 0x050f541a.
//
// Solidity: function totalUnstakeProtocolFee() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) TotalUnstakeProtocolFee() (*big.Int, error) {
	return _StakeERC20Portal.Contract.TotalUnstakeProtocolFee(&_StakeERC20Portal.CallOpts)
}

// UnstakeFeeCommission is a free data retrieval call binding the contract method 0xe1bb5897.
//
// Solidity: function unstakeFeeCommission() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCaller) UnstakeFeeCommission(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "unstakeFeeCommission")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakeFeeCommission is a free data retrieval call binding the contract method 0xe1bb5897.
//
// Solidity: function unstakeFeeCommission() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalSession) UnstakeFeeCommission() (*big.Int, error) {
	return _StakeERC20Portal.Contract.UnstakeFeeCommission(&_StakeERC20Portal.CallOpts)
}

// UnstakeFeeCommission is a free data retrieval call binding the contract method 0xe1bb5897.
//
// Solidity: function unstakeFeeCommission() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) UnstakeFeeCommission() (*big.Int, error) {
	return _StakeERC20Portal.Contract.UnstakeFeeCommission(&_StakeERC20Portal.CallOpts)
}

// UnstakeRelayFee is a free data retrieval call binding the contract method 0xf94bf9c4.
//
// Solidity: function unstakeRelayFee() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCaller) UnstakeRelayFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeERC20Portal.contract.Call(opts, &out, "unstakeRelayFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakeRelayFee is a free data retrieval call binding the contract method 0xf94bf9c4.
//
// Solidity: function unstakeRelayFee() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalSession) UnstakeRelayFee() (*big.Int, error) {
	return _StakeERC20Portal.Contract.UnstakeRelayFee(&_StakeERC20Portal.CallOpts)
}

// UnstakeRelayFee is a free data retrieval call binding the contract method 0xf94bf9c4.
//
// Solidity: function unstakeRelayFee() view returns(uint256)
func (_StakeERC20Portal *StakeERC20PortalCallerSession) UnstakeRelayFee() (*big.Int, error) {
	return _StakeERC20Portal.Contract.UnstakeRelayFee(&_StakeERC20Portal.CallOpts)
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

// AddSubAccount is a paid mutator transaction binding the contract method 0xb52e1eff.
//
// Solidity: function addSubAccount(address _subAccount) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) AddSubAccount(opts *bind.TransactOpts, _subAccount common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "addSubAccount", _subAccount)
}

// AddSubAccount is a paid mutator transaction binding the contract method 0xb52e1eff.
//
// Solidity: function addSubAccount(address _subAccount) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) AddSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.AddSubAccount(&_StakeERC20Portal.TransactOpts, _subAccount)
}

// AddSubAccount is a paid mutator transaction binding the contract method 0xb52e1eff.
//
// Solidity: function addSubAccount(address _subAccount) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) AddSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.AddSubAccount(&_StakeERC20Portal.TransactOpts, _subAccount)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _newThreshold) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) ChangeThreshold(opts *bind.TransactOpts, _newThreshold *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "changeThreshold", _newThreshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _newThreshold) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) ChangeThreshold(_newThreshold *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.ChangeThreshold(&_StakeERC20Portal.TransactOpts, _newThreshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 _newThreshold) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) ChangeThreshold(_newThreshold *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.ChangeThreshold(&_StakeERC20Portal.TransactOpts, _newThreshold)
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

// RemoveSubAccount is a paid mutator transaction binding the contract method 0x0ba3aae3.
//
// Solidity: function removeSubAccount(address _subAccount) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) RemoveSubAccount(opts *bind.TransactOpts, _subAccount common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "removeSubAccount", _subAccount)
}

// RemoveSubAccount is a paid mutator transaction binding the contract method 0x0ba3aae3.
//
// Solidity: function removeSubAccount(address _subAccount) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) RemoveSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.RemoveSubAccount(&_StakeERC20Portal.TransactOpts, _subAccount)
}

// RemoveSubAccount is a paid mutator transaction binding the contract method 0x0ba3aae3.
//
// Solidity: function removeSubAccount(address _subAccount) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) RemoveSubAccount(_subAccount common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.RemoveSubAccount(&_StakeERC20Portal.TransactOpts, _subAccount)
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

// SetMinStakeAmount is a paid mutator transaction binding the contract method 0xeb4af045.
//
// Solidity: function setMinStakeAmount(uint256 _minStakeAmount) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) SetMinStakeAmount(opts *bind.TransactOpts, _minStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "setMinStakeAmount", _minStakeAmount)
}

// SetMinStakeAmount is a paid mutator transaction binding the contract method 0xeb4af045.
//
// Solidity: function setMinStakeAmount(uint256 _minStakeAmount) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) SetMinStakeAmount(_minStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetMinStakeAmount(&_StakeERC20Portal.TransactOpts, _minStakeAmount)
}

// SetMinStakeAmount is a paid mutator transaction binding the contract method 0xeb4af045.
//
// Solidity: function setMinStakeAmount(uint256 _minStakeAmount) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) SetMinStakeAmount(_minStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetMinStakeAmount(&_StakeERC20Portal.TransactOpts, _minStakeAmount)
}

// SetRate is a paid mutator transaction binding the contract method 0x34fcf437.
//
// Solidity: function setRate(uint256 _rate) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) SetRate(opts *bind.TransactOpts, _rate *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "setRate", _rate)
}

// SetRate is a paid mutator transaction binding the contract method 0x34fcf437.
//
// Solidity: function setRate(uint256 _rate) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) SetRate(_rate *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetRate(&_StakeERC20Portal.TransactOpts, _rate)
}

// SetRate is a paid mutator transaction binding the contract method 0x34fcf437.
//
// Solidity: function setRate(uint256 _rate) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) SetRate(_rate *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetRate(&_StakeERC20Portal.TransactOpts, _rate)
}

// SetRateChangeLimit is a paid mutator transaction binding the contract method 0x19826e71.
//
// Solidity: function setRateChangeLimit(uint256 _rateChangeLimit) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) SetRateChangeLimit(opts *bind.TransactOpts, _rateChangeLimit *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "setRateChangeLimit", _rateChangeLimit)
}

// SetRateChangeLimit is a paid mutator transaction binding the contract method 0x19826e71.
//
// Solidity: function setRateChangeLimit(uint256 _rateChangeLimit) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) SetRateChangeLimit(_rateChangeLimit *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetRateChangeLimit(&_StakeERC20Portal.TransactOpts, _rateChangeLimit)
}

// SetRateChangeLimit is a paid mutator transaction binding the contract method 0x19826e71.
//
// Solidity: function setRateChangeLimit(uint256 _rateChangeLimit) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) SetRateChangeLimit(_rateChangeLimit *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetRateChangeLimit(&_StakeERC20Portal.TransactOpts, _rateChangeLimit)
}

// SetStakeRelayFee is a paid mutator transaction binding the contract method 0x31db4715.
//
// Solidity: function setStakeRelayFee(uint256 _stakeRelayFee) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) SetStakeRelayFee(opts *bind.TransactOpts, _stakeRelayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "setStakeRelayFee", _stakeRelayFee)
}

// SetStakeRelayFee is a paid mutator transaction binding the contract method 0x31db4715.
//
// Solidity: function setStakeRelayFee(uint256 _stakeRelayFee) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) SetStakeRelayFee(_stakeRelayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetStakeRelayFee(&_StakeERC20Portal.TransactOpts, _stakeRelayFee)
}

// SetStakeRelayFee is a paid mutator transaction binding the contract method 0x31db4715.
//
// Solidity: function setStakeRelayFee(uint256 _stakeRelayFee) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) SetStakeRelayFee(_stakeRelayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetStakeRelayFee(&_StakeERC20Portal.TransactOpts, _stakeRelayFee)
}

// SetStakeUsePool is a paid mutator transaction binding the contract method 0xe52a0d92.
//
// Solidity: function setStakeUsePool(address _stakeUsePoolAddress) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) SetStakeUsePool(opts *bind.TransactOpts, _stakeUsePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "setStakeUsePool", _stakeUsePoolAddress)
}

// SetStakeUsePool is a paid mutator transaction binding the contract method 0xe52a0d92.
//
// Solidity: function setStakeUsePool(address _stakeUsePoolAddress) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) SetStakeUsePool(_stakeUsePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetStakeUsePool(&_StakeERC20Portal.TransactOpts, _stakeUsePoolAddress)
}

// SetStakeUsePool is a paid mutator transaction binding the contract method 0xe52a0d92.
//
// Solidity: function setStakeUsePool(address _stakeUsePoolAddress) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) SetStakeUsePool(_stakeUsePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetStakeUsePool(&_StakeERC20Portal.TransactOpts, _stakeUsePoolAddress)
}

// SetUnstakeFeeCommission is a paid mutator transaction binding the contract method 0x3fbd062b.
//
// Solidity: function setUnstakeFeeCommission(uint256 _unstakeFeeCommission) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) SetUnstakeFeeCommission(opts *bind.TransactOpts, _unstakeFeeCommission *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "setUnstakeFeeCommission", _unstakeFeeCommission)
}

// SetUnstakeFeeCommission is a paid mutator transaction binding the contract method 0x3fbd062b.
//
// Solidity: function setUnstakeFeeCommission(uint256 _unstakeFeeCommission) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) SetUnstakeFeeCommission(_unstakeFeeCommission *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetUnstakeFeeCommission(&_StakeERC20Portal.TransactOpts, _unstakeFeeCommission)
}

// SetUnstakeFeeCommission is a paid mutator transaction binding the contract method 0x3fbd062b.
//
// Solidity: function setUnstakeFeeCommission(uint256 _unstakeFeeCommission) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) SetUnstakeFeeCommission(_unstakeFeeCommission *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetUnstakeFeeCommission(&_StakeERC20Portal.TransactOpts, _unstakeFeeCommission)
}

// SetUnstakeRelayFee is a paid mutator transaction binding the contract method 0xeb7f83d9.
//
// Solidity: function setUnstakeRelayFee(uint256 _unstakeRelayFee) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) SetUnstakeRelayFee(opts *bind.TransactOpts, _unstakeRelayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "setUnstakeRelayFee", _unstakeRelayFee)
}

// SetUnstakeRelayFee is a paid mutator transaction binding the contract method 0xeb7f83d9.
//
// Solidity: function setUnstakeRelayFee(uint256 _unstakeRelayFee) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) SetUnstakeRelayFee(_unstakeRelayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetUnstakeRelayFee(&_StakeERC20Portal.TransactOpts, _unstakeRelayFee)
}

// SetUnstakeRelayFee is a paid mutator transaction binding the contract method 0xeb7f83d9.
//
// Solidity: function setUnstakeRelayFee(uint256 _unstakeRelayFee) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) SetUnstakeRelayFee(_unstakeRelayFee *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.SetUnstakeRelayFee(&_StakeERC20Portal.TransactOpts, _unstakeRelayFee)
}

// Stake is a paid mutator transaction binding the contract method 0xa694fc3a.
//
// Solidity: function stake(uint256 _amount) payable returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) Stake(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "stake", _amount)
}

// Stake is a paid mutator transaction binding the contract method 0xa694fc3a.
//
// Solidity: function stake(uint256 _amount) payable returns()
func (_StakeERC20Portal *StakeERC20PortalSession) Stake(_amount *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.Stake(&_StakeERC20Portal.TransactOpts, _amount)
}

// Stake is a paid mutator transaction binding the contract method 0xa694fc3a.
//
// Solidity: function stake(uint256 _amount) payable returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) Stake(_amount *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.Stake(&_StakeERC20Portal.TransactOpts, _amount)
}

// StakeAndCross is a paid mutator transaction binding the contract method 0x5e55e287.
//
// Solidity: function stakeAndCross(address _stakePoolAddress, uint256 _amount, uint8 _destChainId, bytes32 _stafiRecipient, address _destRecipient) payable returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) StakeAndCross(opts *bind.TransactOpts, _stakePoolAddress common.Address, _amount *big.Int, _destChainId uint8, _stafiRecipient [32]byte, _destRecipient common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "stakeAndCross", _stakePoolAddress, _amount, _destChainId, _stafiRecipient, _destRecipient)
}

// StakeAndCross is a paid mutator transaction binding the contract method 0x5e55e287.
//
// Solidity: function stakeAndCross(address _stakePoolAddress, uint256 _amount, uint8 _destChainId, bytes32 _stafiRecipient, address _destRecipient) payable returns()
func (_StakeERC20Portal *StakeERC20PortalSession) StakeAndCross(_stakePoolAddress common.Address, _amount *big.Int, _destChainId uint8, _stafiRecipient [32]byte, _destRecipient common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.StakeAndCross(&_StakeERC20Portal.TransactOpts, _stakePoolAddress, _amount, _destChainId, _stafiRecipient, _destRecipient)
}

// StakeAndCross is a paid mutator transaction binding the contract method 0x5e55e287.
//
// Solidity: function stakeAndCross(address _stakePoolAddress, uint256 _amount, uint8 _destChainId, bytes32 _stafiRecipient, address _destRecipient) payable returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) StakeAndCross(_stakePoolAddress common.Address, _amount *big.Int, _destChainId uint8, _stafiRecipient [32]byte, _destRecipient common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.StakeAndCross(&_StakeERC20Portal.TransactOpts, _stakePoolAddress, _amount, _destChainId, _stafiRecipient, _destRecipient)
}

// StakeWithPool is a paid mutator transaction binding the contract method 0x227cdfb7.
//
// Solidity: function stakeWithPool(uint256 _amount, address _stakePoolAddress) payable returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) StakeWithPool(opts *bind.TransactOpts, _amount *big.Int, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "stakeWithPool", _amount, _stakePoolAddress)
}

// StakeWithPool is a paid mutator transaction binding the contract method 0x227cdfb7.
//
// Solidity: function stakeWithPool(uint256 _amount, address _stakePoolAddress) payable returns()
func (_StakeERC20Portal *StakeERC20PortalSession) StakeWithPool(_amount *big.Int, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.StakeWithPool(&_StakeERC20Portal.TransactOpts, _amount, _stakePoolAddress)
}

// StakeWithPool is a paid mutator transaction binding the contract method 0x227cdfb7.
//
// Solidity: function stakeWithPool(uint256 _amount, address _stakePoolAddress) payable returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) StakeWithPool(_amount *big.Int, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.StakeWithPool(&_StakeERC20Portal.TransactOpts, _amount, _stakePoolAddress)
}

// ToggleStakeCrossSwitch is a paid mutator transaction binding the contract method 0x6759d907.
//
// Solidity: function toggleStakeCrossSwitch() returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) ToggleStakeCrossSwitch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "toggleStakeCrossSwitch")
}

// ToggleStakeCrossSwitch is a paid mutator transaction binding the contract method 0x6759d907.
//
// Solidity: function toggleStakeCrossSwitch() returns()
func (_StakeERC20Portal *StakeERC20PortalSession) ToggleStakeCrossSwitch() (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.ToggleStakeCrossSwitch(&_StakeERC20Portal.TransactOpts)
}

// ToggleStakeCrossSwitch is a paid mutator transaction binding the contract method 0x6759d907.
//
// Solidity: function toggleStakeCrossSwitch() returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) ToggleStakeCrossSwitch() (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.ToggleStakeCrossSwitch(&_StakeERC20Portal.TransactOpts)
}

// ToggleStakeSwitch is a paid mutator transaction binding the contract method 0x995bb87d.
//
// Solidity: function toggleStakeSwitch() returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) ToggleStakeSwitch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "toggleStakeSwitch")
}

// ToggleStakeSwitch is a paid mutator transaction binding the contract method 0x995bb87d.
//
// Solidity: function toggleStakeSwitch() returns()
func (_StakeERC20Portal *StakeERC20PortalSession) ToggleStakeSwitch() (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.ToggleStakeSwitch(&_StakeERC20Portal.TransactOpts)
}

// ToggleStakeSwitch is a paid mutator transaction binding the contract method 0x995bb87d.
//
// Solidity: function toggleStakeSwitch() returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) ToggleStakeSwitch() (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.ToggleStakeSwitch(&_StakeERC20Portal.TransactOpts)
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

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 _rTokenAmount) payable returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) Unstake(opts *bind.TransactOpts, _rTokenAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "unstake", _rTokenAmount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 _rTokenAmount) payable returns()
func (_StakeERC20Portal *StakeERC20PortalSession) Unstake(_rTokenAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.Unstake(&_StakeERC20Portal.TransactOpts, _rTokenAmount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 _rTokenAmount) payable returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) Unstake(_rTokenAmount *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.Unstake(&_StakeERC20Portal.TransactOpts, _rTokenAmount)
}

// UnstakeWithPool is a paid mutator transaction binding the contract method 0xac8ef2db.
//
// Solidity: function unstakeWithPool(uint256 _rTokenAmount, address _stakePoolAddress) payable returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) UnstakeWithPool(opts *bind.TransactOpts, _rTokenAmount *big.Int, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "unstakeWithPool", _rTokenAmount, _stakePoolAddress)
}

// UnstakeWithPool is a paid mutator transaction binding the contract method 0xac8ef2db.
//
// Solidity: function unstakeWithPool(uint256 _rTokenAmount, address _stakePoolAddress) payable returns()
func (_StakeERC20Portal *StakeERC20PortalSession) UnstakeWithPool(_rTokenAmount *big.Int, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.UnstakeWithPool(&_StakeERC20Portal.TransactOpts, _rTokenAmount, _stakePoolAddress)
}

// UnstakeWithPool is a paid mutator transaction binding the contract method 0xac8ef2db.
//
// Solidity: function unstakeWithPool(uint256 _rTokenAmount, address _stakePoolAddress) payable returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) UnstakeWithPool(_rTokenAmount *big.Int, _stakePoolAddress common.Address) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.UnstakeWithPool(&_StakeERC20Portal.TransactOpts, _rTokenAmount, _stakePoolAddress)
}

// VoteRate is a paid mutator transaction binding the contract method 0x113e3709.
//
// Solidity: function voteRate(bytes32 _proposalId, uint256 _rate) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactor) VoteRate(opts *bind.TransactOpts, _proposalId [32]byte, _rate *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.contract.Transact(opts, "voteRate", _proposalId, _rate)
}

// VoteRate is a paid mutator transaction binding the contract method 0x113e3709.
//
// Solidity: function voteRate(bytes32 _proposalId, uint256 _rate) returns()
func (_StakeERC20Portal *StakeERC20PortalSession) VoteRate(_proposalId [32]byte, _rate *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.VoteRate(&_StakeERC20Portal.TransactOpts, _proposalId, _rate)
}

// VoteRate is a paid mutator transaction binding the contract method 0x113e3709.
//
// Solidity: function voteRate(bytes32 _proposalId, uint256 _rate) returns()
func (_StakeERC20Portal *StakeERC20PortalTransactorSession) VoteRate(_proposalId [32]byte, _rate *big.Int) (*types.Transaction, error) {
	return _StakeERC20Portal.Contract.VoteRate(&_StakeERC20Portal.TransactOpts, _proposalId, _rate)
}

// StakeERC20PortalProposalExecutedIterator is returned from FilterProposalExecuted and is used to iterate over the raw logs and unpacked data for ProposalExecuted events raised by the StakeERC20Portal contract.
type StakeERC20PortalProposalExecutedIterator struct {
	Event *StakeERC20PortalProposalExecuted // Event containing the contract specifics and raw log

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
func (it *StakeERC20PortalProposalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeERC20PortalProposalExecuted)
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
		it.Event = new(StakeERC20PortalProposalExecuted)
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
func (it *StakeERC20PortalProposalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeERC20PortalProposalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeERC20PortalProposalExecuted represents a ProposalExecuted event raised by the StakeERC20Portal contract.
type StakeERC20PortalProposalExecuted struct {
	ProposalId [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalExecuted is a free log retrieval operation binding the contract event 0x7b1bcf1ccf901a11589afff5504d59fd0a53780eed2a952adade0348985139e0.
//
// Solidity: event ProposalExecuted(bytes32 indexed proposalId)
func (_StakeERC20Portal *StakeERC20PortalFilterer) FilterProposalExecuted(opts *bind.FilterOpts, proposalId [][32]byte) (*StakeERC20PortalProposalExecutedIterator, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _StakeERC20Portal.contract.FilterLogs(opts, "ProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalProposalExecutedIterator{contract: _StakeERC20Portal.contract, event: "ProposalExecuted", logs: logs, sub: sub}, nil
}

// WatchProposalExecuted is a free log subscription operation binding the contract event 0x7b1bcf1ccf901a11589afff5504d59fd0a53780eed2a952adade0348985139e0.
//
// Solidity: event ProposalExecuted(bytes32 indexed proposalId)
func (_StakeERC20Portal *StakeERC20PortalFilterer) WatchProposalExecuted(opts *bind.WatchOpts, sink chan<- *StakeERC20PortalProposalExecuted, proposalId [][32]byte) (event.Subscription, error) {

	var proposalIdRule []interface{}
	for _, proposalIdItem := range proposalId {
		proposalIdRule = append(proposalIdRule, proposalIdItem)
	}

	logs, sub, err := _StakeERC20Portal.contract.WatchLogs(opts, "ProposalExecuted", proposalIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeERC20PortalProposalExecuted)
				if err := _StakeERC20Portal.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
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
func (_StakeERC20Portal *StakeERC20PortalFilterer) ParseProposalExecuted(log types.Log) (*StakeERC20PortalProposalExecuted, error) {
	event := new(StakeERC20PortalProposalExecuted)
	if err := _StakeERC20Portal.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	Staker       common.Address
	StakePool    common.Address
	TokenAmount  *big.Int
	RTokenAmount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStake is a free log retrieval operation binding the contract event 0x63602d0ecc7b3a0ef7ff1a116e23056662d64280355ba8031b6d0d767c4b4458.
//
// Solidity: event Stake(address staker, address stakePool, uint256 tokenAmount, uint256 rTokenAmount)
func (_StakeERC20Portal *StakeERC20PortalFilterer) FilterStake(opts *bind.FilterOpts) (*StakeERC20PortalStakeIterator, error) {

	logs, sub, err := _StakeERC20Portal.contract.FilterLogs(opts, "Stake")
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalStakeIterator{contract: _StakeERC20Portal.contract, event: "Stake", logs: logs, sub: sub}, nil
}

// WatchStake is a free log subscription operation binding the contract event 0x63602d0ecc7b3a0ef7ff1a116e23056662d64280355ba8031b6d0d767c4b4458.
//
// Solidity: event Stake(address staker, address stakePool, uint256 tokenAmount, uint256 rTokenAmount)
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

// ParseStake is a log parse operation binding the contract event 0x63602d0ecc7b3a0ef7ff1a116e23056662d64280355ba8031b6d0d767c4b4458.
//
// Solidity: event Stake(address staker, address stakePool, uint256 tokenAmount, uint256 rTokenAmount)
func (_StakeERC20Portal *StakeERC20PortalFilterer) ParseStake(log types.Log) (*StakeERC20PortalStake, error) {
	event := new(StakeERC20PortalStake)
	if err := _StakeERC20Portal.contract.UnpackLog(event, "Stake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeERC20PortalStakeAndCrossIterator is returned from FilterStakeAndCross and is used to iterate over the raw logs and unpacked data for StakeAndCross events raised by the StakeERC20Portal contract.
type StakeERC20PortalStakeAndCrossIterator struct {
	Event *StakeERC20PortalStakeAndCross // Event containing the contract specifics and raw log

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
func (it *StakeERC20PortalStakeAndCrossIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeERC20PortalStakeAndCross)
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
		it.Event = new(StakeERC20PortalStakeAndCross)
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
func (it *StakeERC20PortalStakeAndCrossIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeERC20PortalStakeAndCrossIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeERC20PortalStakeAndCross represents a StakeAndCross event raised by the StakeERC20Portal contract.
type StakeERC20PortalStakeAndCross struct {
	Staker         common.Address
	StakePool      common.Address
	Amount         *big.Int
	ChainId        uint8
	StafiRecipient [32]byte
	DestRecipient  common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterStakeAndCross is a free log retrieval operation binding the contract event 0xfb3b5205262b9f062f710d9de902868cf013909615405f915d7f472e0ac152ff.
//
// Solidity: event StakeAndCross(address staker, address stakePool, uint256 amount, uint8 chainId, bytes32 stafiRecipient, address destRecipient)
func (_StakeERC20Portal *StakeERC20PortalFilterer) FilterStakeAndCross(opts *bind.FilterOpts) (*StakeERC20PortalStakeAndCrossIterator, error) {

	logs, sub, err := _StakeERC20Portal.contract.FilterLogs(opts, "StakeAndCross")
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalStakeAndCrossIterator{contract: _StakeERC20Portal.contract, event: "StakeAndCross", logs: logs, sub: sub}, nil
}

// WatchStakeAndCross is a free log subscription operation binding the contract event 0xfb3b5205262b9f062f710d9de902868cf013909615405f915d7f472e0ac152ff.
//
// Solidity: event StakeAndCross(address staker, address stakePool, uint256 amount, uint8 chainId, bytes32 stafiRecipient, address destRecipient)
func (_StakeERC20Portal *StakeERC20PortalFilterer) WatchStakeAndCross(opts *bind.WatchOpts, sink chan<- *StakeERC20PortalStakeAndCross) (event.Subscription, error) {

	logs, sub, err := _StakeERC20Portal.contract.WatchLogs(opts, "StakeAndCross")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeERC20PortalStakeAndCross)
				if err := _StakeERC20Portal.contract.UnpackLog(event, "StakeAndCross", log); err != nil {
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

// ParseStakeAndCross is a log parse operation binding the contract event 0xfb3b5205262b9f062f710d9de902868cf013909615405f915d7f472e0ac152ff.
//
// Solidity: event StakeAndCross(address staker, address stakePool, uint256 amount, uint8 chainId, bytes32 stafiRecipient, address destRecipient)
func (_StakeERC20Portal *StakeERC20PortalFilterer) ParseStakeAndCross(log types.Log) (*StakeERC20PortalStakeAndCross, error) {
	event := new(StakeERC20PortalStakeAndCross)
	if err := _StakeERC20Portal.contract.UnpackLog(event, "StakeAndCross", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeERC20PortalUnstakeIterator is returned from FilterUnstake and is used to iterate over the raw logs and unpacked data for Unstake events raised by the StakeERC20Portal contract.
type StakeERC20PortalUnstakeIterator struct {
	Event *StakeERC20PortalUnstake // Event containing the contract specifics and raw log

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
func (it *StakeERC20PortalUnstakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeERC20PortalUnstake)
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
		it.Event = new(StakeERC20PortalUnstake)
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
func (it *StakeERC20PortalUnstakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeERC20PortalUnstakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeERC20PortalUnstake represents a Unstake event raised by the StakeERC20Portal contract.
type StakeERC20PortalUnstake struct {
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
func (_StakeERC20Portal *StakeERC20PortalFilterer) FilterUnstake(opts *bind.FilterOpts) (*StakeERC20PortalUnstakeIterator, error) {

	logs, sub, err := _StakeERC20Portal.contract.FilterLogs(opts, "Unstake")
	if err != nil {
		return nil, err
	}
	return &StakeERC20PortalUnstakeIterator{contract: _StakeERC20Portal.contract, event: "Unstake", logs: logs, sub: sub}, nil
}

// WatchUnstake is a free log subscription operation binding the contract event 0xfe7007b2e89d80edda76299251df08366480cac22e5e260f5e662e850b1f7a32.
//
// Solidity: event Unstake(address staker, address stakePool, uint256 tokenAmount, uint256 rTokenAmount, uint256 burnAmount)
func (_StakeERC20Portal *StakeERC20PortalFilterer) WatchUnstake(opts *bind.WatchOpts, sink chan<- *StakeERC20PortalUnstake) (event.Subscription, error) {

	logs, sub, err := _StakeERC20Portal.contract.WatchLogs(opts, "Unstake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeERC20PortalUnstake)
				if err := _StakeERC20Portal.contract.UnpackLog(event, "Unstake", log); err != nil {
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
func (_StakeERC20Portal *StakeERC20PortalFilterer) ParseUnstake(log types.Log) (*StakeERC20PortalUnstake, error) {
	event := new(StakeERC20PortalUnstake)
	if err := _StakeERC20Portal.contract.UnpackLog(event, "Unstake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
