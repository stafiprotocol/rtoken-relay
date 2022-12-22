// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package staking

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

// StakingMetaData contains all meta data concerning the Staking contract.
var StakingMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"eventType\",\"type\":\"uint8\"}],\"name\":\"crashResponse\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"errCode\",\"type\":\"uint8\"}],\"name\":\"delegateFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"relayerFee\",\"type\":\"uint256\"}],\"name\":\"delegateSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"delegateSuccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"errCode\",\"type\":\"uint256\"}],\"name\":\"failedSynPackage\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"paramChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"valSrc\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"valDst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"errCode\",\"type\":\"uint8\"}],\"name\":\"redelegateFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorSrc\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorDst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"relayerFee\",\"type\":\"uint256\"}],\"name\":\"redelegateSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"valSrc\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"valDst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"redelegateSuccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"rewardClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"rewardReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"errCode\",\"type\":\"uint8\"}],\"name\":\"undelegateFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"relayerFee\",\"type\":\"uint256\"}],\"name\":\"undelegateSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"undelegateSuccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"undelegatedClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"undelegatedReceived\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BIND_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CODE_FAILED\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CODE_OK\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CODE_SUCCESS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CROSS_CHAIN_CONTRACT_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CROSS_STAKE_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ERROR_FAIL_DECODE\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ERROR_WITHDRAW_BNB\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EVENT_DELEGATE\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EVENT_DISTRIBUTE_REWARD\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EVENT_DISTRIBUTE_UNDELEGATED\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EVENT_REDELEGATE\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EVENT_UNDELEGATE\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GOV_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GOV_HUB_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INCENTIVIZE_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INIT_BSC_RELAYER_FEE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INIT_MIN_DELEGATION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INIT_RELAYER_FEE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LIGHT_CLIENT_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"LOCK_TIME\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RELAYERHUB_CONTRACT_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SLASH_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SLASH_CONTRACT_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"STAKING_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"STAKING_CONTRACT_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SYSTEM_REWARD_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TEN_DECIMALS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TOKEN_HUB_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TOKEN_MANAGER_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TRANSFER_IN_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TRANSFER_OUT_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"VALIDATOR_CONTRACT_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"alreadyInit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bSCRelayerFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bscChainID\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minDelegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"relayerFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"msgBytes\",\"type\":\"bytes\"}],\"name\":\"handleSynPackage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"msgBytes\",\"type\":\"bytes\"}],\"name\":\"handleAckPackage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"msgBytes\",\"type\":\"bytes\"}],\"name\":\"handleFailAckPackage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"delegate\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorSrc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorDst\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"redelegate\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimUndelegated\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"getDelegated\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"}],\"name\":\"getTotalDelegated\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"}],\"name\":\"getDistributedReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"valSrc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"valDst\",\"type\":\"address\"}],\"name\":\"getPendingRedelegateTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"}],\"name\":\"getUndelegated\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"getPendingUndelegateTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRelayerFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinDelegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"}],\"name\":\"getRequestInFly\",\"outputs\":[{\"internalType\":\"uint256[3]\",\"name\":\"\",\"type\":\"uint256[3]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"updateParam\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StakingABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingMetaData.ABI instead.
var StakingABI = StakingMetaData.ABI

// Staking is an auto generated Go binding around an Ethereum contract.
type Staking struct {
	StakingCaller     // Read-only binding to the contract
	StakingTransactor // Write-only binding to the contract
	StakingFilterer   // Log filterer for contract events
}

// StakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingSession struct {
	Contract     *Staking          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingCallerSession struct {
	Contract *StakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingTransactorSession struct {
	Contract     *StakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingRaw struct {
	Contract *Staking // Generic contract binding to access the raw methods on
}

// StakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingCallerRaw struct {
	Contract *StakingCaller // Generic read-only contract binding to access the raw methods on
}

// StakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingTransactorRaw struct {
	Contract *StakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStaking creates a new instance of Staking, bound to a specific deployed contract.
func NewStaking(address common.Address, backend bind.ContractBackend) (*Staking, error) {
	contract, err := bindStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// NewStakingCaller creates a new read-only instance of Staking, bound to a specific deployed contract.
func NewStakingCaller(address common.Address, caller bind.ContractCaller) (*StakingCaller, error) {
	contract, err := bindStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingCaller{contract: contract}, nil
}

// NewStakingTransactor creates a new write-only instance of Staking, bound to a specific deployed contract.
func NewStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingTransactor, error) {
	contract, err := bindStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTransactor{contract: contract}, nil
}

// NewStakingFilterer creates a new log filterer instance of Staking, bound to a specific deployed contract.
func NewStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingFilterer, error) {
	contract, err := bindStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingFilterer{contract: contract}, nil
}

// bindStaking binds a generic wrapper to an already deployed contract.
func bindStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.StakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transact(opts, method, params...)
}

// BINDCHANNELID is a free data retrieval call binding the contract method 0x3dffc387.
//
// Solidity: function BIND_CHANNELID() view returns(uint8)
func (_Staking *StakingCaller) BINDCHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "BIND_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// BINDCHANNELID is a free data retrieval call binding the contract method 0x3dffc387.
//
// Solidity: function BIND_CHANNELID() view returns(uint8)
func (_Staking *StakingSession) BINDCHANNELID() (uint8, error) {
	return _Staking.Contract.BINDCHANNELID(&_Staking.CallOpts)
}

// BINDCHANNELID is a free data retrieval call binding the contract method 0x3dffc387.
//
// Solidity: function BIND_CHANNELID() view returns(uint8)
func (_Staking *StakingCallerSession) BINDCHANNELID() (uint8, error) {
	return _Staking.Contract.BINDCHANNELID(&_Staking.CallOpts)
}

// CODEFAILED is a free data retrieval call binding the contract method 0x552aaf93.
//
// Solidity: function CODE_FAILED() view returns(uint8)
func (_Staking *StakingCaller) CODEFAILED(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "CODE_FAILED")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CODEFAILED is a free data retrieval call binding the contract method 0x552aaf93.
//
// Solidity: function CODE_FAILED() view returns(uint8)
func (_Staking *StakingSession) CODEFAILED() (uint8, error) {
	return _Staking.Contract.CODEFAILED(&_Staking.CallOpts)
}

// CODEFAILED is a free data retrieval call binding the contract method 0x552aaf93.
//
// Solidity: function CODE_FAILED() view returns(uint8)
func (_Staking *StakingCallerSession) CODEFAILED() (uint8, error) {
	return _Staking.Contract.CODEFAILED(&_Staking.CallOpts)
}

// CODEOK is a free data retrieval call binding the contract method 0xab51bb96.
//
// Solidity: function CODE_OK() view returns(uint32)
func (_Staking *StakingCaller) CODEOK(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "CODE_OK")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// CODEOK is a free data retrieval call binding the contract method 0xab51bb96.
//
// Solidity: function CODE_OK() view returns(uint32)
func (_Staking *StakingSession) CODEOK() (uint32, error) {
	return _Staking.Contract.CODEOK(&_Staking.CallOpts)
}

// CODEOK is a free data retrieval call binding the contract method 0xab51bb96.
//
// Solidity: function CODE_OK() view returns(uint32)
func (_Staking *StakingCallerSession) CODEOK() (uint32, error) {
	return _Staking.Contract.CODEOK(&_Staking.CallOpts)
}

// CODESUCCESS is a free data retrieval call binding the contract method 0x17c9efb0.
//
// Solidity: function CODE_SUCCESS() view returns(uint8)
func (_Staking *StakingCaller) CODESUCCESS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "CODE_SUCCESS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CODESUCCESS is a free data retrieval call binding the contract method 0x17c9efb0.
//
// Solidity: function CODE_SUCCESS() view returns(uint8)
func (_Staking *StakingSession) CODESUCCESS() (uint8, error) {
	return _Staking.Contract.CODESUCCESS(&_Staking.CallOpts)
}

// CODESUCCESS is a free data retrieval call binding the contract method 0x17c9efb0.
//
// Solidity: function CODE_SUCCESS() view returns(uint8)
func (_Staking *StakingCallerSession) CODESUCCESS() (uint8, error) {
	return _Staking.Contract.CODESUCCESS(&_Staking.CallOpts)
}

// CROSSCHAINCONTRACTADDR is a free data retrieval call binding the contract method 0x51e80672.
//
// Solidity: function CROSS_CHAIN_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingCaller) CROSSCHAINCONTRACTADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "CROSS_CHAIN_CONTRACT_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CROSSCHAINCONTRACTADDR is a free data retrieval call binding the contract method 0x51e80672.
//
// Solidity: function CROSS_CHAIN_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingSession) CROSSCHAINCONTRACTADDR() (common.Address, error) {
	return _Staking.Contract.CROSSCHAINCONTRACTADDR(&_Staking.CallOpts)
}

// CROSSCHAINCONTRACTADDR is a free data retrieval call binding the contract method 0x51e80672.
//
// Solidity: function CROSS_CHAIN_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingCallerSession) CROSSCHAINCONTRACTADDR() (common.Address, error) {
	return _Staking.Contract.CROSSCHAINCONTRACTADDR(&_Staking.CallOpts)
}

// CROSSSTAKECHANNELID is a free data retrieval call binding the contract method 0x718a8aa8.
//
// Solidity: function CROSS_STAKE_CHANNELID() view returns(uint8)
func (_Staking *StakingCaller) CROSSSTAKECHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "CROSS_STAKE_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CROSSSTAKECHANNELID is a free data retrieval call binding the contract method 0x718a8aa8.
//
// Solidity: function CROSS_STAKE_CHANNELID() view returns(uint8)
func (_Staking *StakingSession) CROSSSTAKECHANNELID() (uint8, error) {
	return _Staking.Contract.CROSSSTAKECHANNELID(&_Staking.CallOpts)
}

// CROSSSTAKECHANNELID is a free data retrieval call binding the contract method 0x718a8aa8.
//
// Solidity: function CROSS_STAKE_CHANNELID() view returns(uint8)
func (_Staking *StakingCallerSession) CROSSSTAKECHANNELID() (uint8, error) {
	return _Staking.Contract.CROSSSTAKECHANNELID(&_Staking.CallOpts)
}

// ERRORFAILDECODE is a free data retrieval call binding the contract method 0x0bee7a67.
//
// Solidity: function ERROR_FAIL_DECODE() view returns(uint32)
func (_Staking *StakingCaller) ERRORFAILDECODE(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "ERROR_FAIL_DECODE")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// ERRORFAILDECODE is a free data retrieval call binding the contract method 0x0bee7a67.
//
// Solidity: function ERROR_FAIL_DECODE() view returns(uint32)
func (_Staking *StakingSession) ERRORFAILDECODE() (uint32, error) {
	return _Staking.Contract.ERRORFAILDECODE(&_Staking.CallOpts)
}

// ERRORFAILDECODE is a free data retrieval call binding the contract method 0x0bee7a67.
//
// Solidity: function ERROR_FAIL_DECODE() view returns(uint32)
func (_Staking *StakingCallerSession) ERRORFAILDECODE() (uint32, error) {
	return _Staking.Contract.ERRORFAILDECODE(&_Staking.CallOpts)
}

// ERRORWITHDRAWBNB is a free data retrieval call binding the contract method 0x333ad3e7.
//
// Solidity: function ERROR_WITHDRAW_BNB() view returns(uint32)
func (_Staking *StakingCaller) ERRORWITHDRAWBNB(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "ERROR_WITHDRAW_BNB")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// ERRORWITHDRAWBNB is a free data retrieval call binding the contract method 0x333ad3e7.
//
// Solidity: function ERROR_WITHDRAW_BNB() view returns(uint32)
func (_Staking *StakingSession) ERRORWITHDRAWBNB() (uint32, error) {
	return _Staking.Contract.ERRORWITHDRAWBNB(&_Staking.CallOpts)
}

// ERRORWITHDRAWBNB is a free data retrieval call binding the contract method 0x333ad3e7.
//
// Solidity: function ERROR_WITHDRAW_BNB() view returns(uint32)
func (_Staking *StakingCallerSession) ERRORWITHDRAWBNB() (uint32, error) {
	return _Staking.Contract.ERRORWITHDRAWBNB(&_Staking.CallOpts)
}

// EVENTDELEGATE is a free data retrieval call binding the contract method 0x92b888a4.
//
// Solidity: function EVENT_DELEGATE() view returns(uint8)
func (_Staking *StakingCaller) EVENTDELEGATE(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "EVENT_DELEGATE")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// EVENTDELEGATE is a free data retrieval call binding the contract method 0x92b888a4.
//
// Solidity: function EVENT_DELEGATE() view returns(uint8)
func (_Staking *StakingSession) EVENTDELEGATE() (uint8, error) {
	return _Staking.Contract.EVENTDELEGATE(&_Staking.CallOpts)
}

// EVENTDELEGATE is a free data retrieval call binding the contract method 0x92b888a4.
//
// Solidity: function EVENT_DELEGATE() view returns(uint8)
func (_Staking *StakingCallerSession) EVENTDELEGATE() (uint8, error) {
	return _Staking.Contract.EVENTDELEGATE(&_Staking.CallOpts)
}

// EVENTDISTRIBUTEREWARD is a free data retrieval call binding the contract method 0xb14315df.
//
// Solidity: function EVENT_DISTRIBUTE_REWARD() view returns(uint8)
func (_Staking *StakingCaller) EVENTDISTRIBUTEREWARD(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "EVENT_DISTRIBUTE_REWARD")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// EVENTDISTRIBUTEREWARD is a free data retrieval call binding the contract method 0xb14315df.
//
// Solidity: function EVENT_DISTRIBUTE_REWARD() view returns(uint8)
func (_Staking *StakingSession) EVENTDISTRIBUTEREWARD() (uint8, error) {
	return _Staking.Contract.EVENTDISTRIBUTEREWARD(&_Staking.CallOpts)
}

// EVENTDISTRIBUTEREWARD is a free data retrieval call binding the contract method 0xb14315df.
//
// Solidity: function EVENT_DISTRIBUTE_REWARD() view returns(uint8)
func (_Staking *StakingCallerSession) EVENTDISTRIBUTEREWARD() (uint8, error) {
	return _Staking.Contract.EVENTDISTRIBUTEREWARD(&_Staking.CallOpts)
}

// EVENTDISTRIBUTEUNDELEGATED is a free data retrieval call binding the contract method 0x151817e3.
//
// Solidity: function EVENT_DISTRIBUTE_UNDELEGATED() view returns(uint8)
func (_Staking *StakingCaller) EVENTDISTRIBUTEUNDELEGATED(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "EVENT_DISTRIBUTE_UNDELEGATED")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// EVENTDISTRIBUTEUNDELEGATED is a free data retrieval call binding the contract method 0x151817e3.
//
// Solidity: function EVENT_DISTRIBUTE_UNDELEGATED() view returns(uint8)
func (_Staking *StakingSession) EVENTDISTRIBUTEUNDELEGATED() (uint8, error) {
	return _Staking.Contract.EVENTDISTRIBUTEUNDELEGATED(&_Staking.CallOpts)
}

// EVENTDISTRIBUTEUNDELEGATED is a free data retrieval call binding the contract method 0x151817e3.
//
// Solidity: function EVENT_DISTRIBUTE_UNDELEGATED() view returns(uint8)
func (_Staking *StakingCallerSession) EVENTDISTRIBUTEUNDELEGATED() (uint8, error) {
	return _Staking.Contract.EVENTDISTRIBUTEUNDELEGATED(&_Staking.CallOpts)
}

// EVENTREDELEGATE is a free data retrieval call binding the contract method 0x3fdfa7e4.
//
// Solidity: function EVENT_REDELEGATE() view returns(uint8)
func (_Staking *StakingCaller) EVENTREDELEGATE(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "EVENT_REDELEGATE")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// EVENTREDELEGATE is a free data retrieval call binding the contract method 0x3fdfa7e4.
//
// Solidity: function EVENT_REDELEGATE() view returns(uint8)
func (_Staking *StakingSession) EVENTREDELEGATE() (uint8, error) {
	return _Staking.Contract.EVENTREDELEGATE(&_Staking.CallOpts)
}

// EVENTREDELEGATE is a free data retrieval call binding the contract method 0x3fdfa7e4.
//
// Solidity: function EVENT_REDELEGATE() view returns(uint8)
func (_Staking *StakingCallerSession) EVENTREDELEGATE() (uint8, error) {
	return _Staking.Contract.EVENTREDELEGATE(&_Staking.CallOpts)
}

// EVENTUNDELEGATE is a free data retrieval call binding the contract method 0xd7ecfcb6.
//
// Solidity: function EVENT_UNDELEGATE() view returns(uint8)
func (_Staking *StakingCaller) EVENTUNDELEGATE(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "EVENT_UNDELEGATE")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// EVENTUNDELEGATE is a free data retrieval call binding the contract method 0xd7ecfcb6.
//
// Solidity: function EVENT_UNDELEGATE() view returns(uint8)
func (_Staking *StakingSession) EVENTUNDELEGATE() (uint8, error) {
	return _Staking.Contract.EVENTUNDELEGATE(&_Staking.CallOpts)
}

// EVENTUNDELEGATE is a free data retrieval call binding the contract method 0xd7ecfcb6.
//
// Solidity: function EVENT_UNDELEGATE() view returns(uint8)
func (_Staking *StakingCallerSession) EVENTUNDELEGATE() (uint8, error) {
	return _Staking.Contract.EVENTUNDELEGATE(&_Staking.CallOpts)
}

// GOVCHANNELID is a free data retrieval call binding the contract method 0x96713da9.
//
// Solidity: function GOV_CHANNELID() view returns(uint8)
func (_Staking *StakingCaller) GOVCHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "GOV_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GOVCHANNELID is a free data retrieval call binding the contract method 0x96713da9.
//
// Solidity: function GOV_CHANNELID() view returns(uint8)
func (_Staking *StakingSession) GOVCHANNELID() (uint8, error) {
	return _Staking.Contract.GOVCHANNELID(&_Staking.CallOpts)
}

// GOVCHANNELID is a free data retrieval call binding the contract method 0x96713da9.
//
// Solidity: function GOV_CHANNELID() view returns(uint8)
func (_Staking *StakingCallerSession) GOVCHANNELID() (uint8, error) {
	return _Staking.Contract.GOVCHANNELID(&_Staking.CallOpts)
}

// GOVHUBADDR is a free data retrieval call binding the contract method 0x9dc09262.
//
// Solidity: function GOV_HUB_ADDR() view returns(address)
func (_Staking *StakingCaller) GOVHUBADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "GOV_HUB_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GOVHUBADDR is a free data retrieval call binding the contract method 0x9dc09262.
//
// Solidity: function GOV_HUB_ADDR() view returns(address)
func (_Staking *StakingSession) GOVHUBADDR() (common.Address, error) {
	return _Staking.Contract.GOVHUBADDR(&_Staking.CallOpts)
}

// GOVHUBADDR is a free data retrieval call binding the contract method 0x9dc09262.
//
// Solidity: function GOV_HUB_ADDR() view returns(address)
func (_Staking *StakingCallerSession) GOVHUBADDR() (common.Address, error) {
	return _Staking.Contract.GOVHUBADDR(&_Staking.CallOpts)
}

// INCENTIVIZEADDR is a free data retrieval call binding the contract method 0x6e47b482.
//
// Solidity: function INCENTIVIZE_ADDR() view returns(address)
func (_Staking *StakingCaller) INCENTIVIZEADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "INCENTIVIZE_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// INCENTIVIZEADDR is a free data retrieval call binding the contract method 0x6e47b482.
//
// Solidity: function INCENTIVIZE_ADDR() view returns(address)
func (_Staking *StakingSession) INCENTIVIZEADDR() (common.Address, error) {
	return _Staking.Contract.INCENTIVIZEADDR(&_Staking.CallOpts)
}

// INCENTIVIZEADDR is a free data retrieval call binding the contract method 0x6e47b482.
//
// Solidity: function INCENTIVIZE_ADDR() view returns(address)
func (_Staking *StakingCallerSession) INCENTIVIZEADDR() (common.Address, error) {
	return _Staking.Contract.INCENTIVIZEADDR(&_Staking.CallOpts)
}

// INITBSCRELAYERFEE is a free data retrieval call binding the contract method 0x34c43354.
//
// Solidity: function INIT_BSC_RELAYER_FEE() view returns(uint256)
func (_Staking *StakingCaller) INITBSCRELAYERFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "INIT_BSC_RELAYER_FEE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INITBSCRELAYERFEE is a free data retrieval call binding the contract method 0x34c43354.
//
// Solidity: function INIT_BSC_RELAYER_FEE() view returns(uint256)
func (_Staking *StakingSession) INITBSCRELAYERFEE() (*big.Int, error) {
	return _Staking.Contract.INITBSCRELAYERFEE(&_Staking.CallOpts)
}

// INITBSCRELAYERFEE is a free data retrieval call binding the contract method 0x34c43354.
//
// Solidity: function INIT_BSC_RELAYER_FEE() view returns(uint256)
func (_Staking *StakingCallerSession) INITBSCRELAYERFEE() (*big.Int, error) {
	return _Staking.Contract.INITBSCRELAYERFEE(&_Staking.CallOpts)
}

// INITMINDELEGATION is a free data retrieval call binding the contract method 0xedc1a5b0.
//
// Solidity: function INIT_MIN_DELEGATION() view returns(uint256)
func (_Staking *StakingCaller) INITMINDELEGATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "INIT_MIN_DELEGATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INITMINDELEGATION is a free data retrieval call binding the contract method 0xedc1a5b0.
//
// Solidity: function INIT_MIN_DELEGATION() view returns(uint256)
func (_Staking *StakingSession) INITMINDELEGATION() (*big.Int, error) {
	return _Staking.Contract.INITMINDELEGATION(&_Staking.CallOpts)
}

// INITMINDELEGATION is a free data retrieval call binding the contract method 0xedc1a5b0.
//
// Solidity: function INIT_MIN_DELEGATION() view returns(uint256)
func (_Staking *StakingCallerSession) INITMINDELEGATION() (*big.Int, error) {
	return _Staking.Contract.INITMINDELEGATION(&_Staking.CallOpts)
}

// INITRELAYERFEE is a free data retrieval call binding the contract method 0xbaaafd3b.
//
// Solidity: function INIT_RELAYER_FEE() view returns(uint256)
func (_Staking *StakingCaller) INITRELAYERFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "INIT_RELAYER_FEE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INITRELAYERFEE is a free data retrieval call binding the contract method 0xbaaafd3b.
//
// Solidity: function INIT_RELAYER_FEE() view returns(uint256)
func (_Staking *StakingSession) INITRELAYERFEE() (*big.Int, error) {
	return _Staking.Contract.INITRELAYERFEE(&_Staking.CallOpts)
}

// INITRELAYERFEE is a free data retrieval call binding the contract method 0xbaaafd3b.
//
// Solidity: function INIT_RELAYER_FEE() view returns(uint256)
func (_Staking *StakingCallerSession) INITRELAYERFEE() (*big.Int, error) {
	return _Staking.Contract.INITRELAYERFEE(&_Staking.CallOpts)
}

// LIGHTCLIENTADDR is a free data retrieval call binding the contract method 0xdc927faf.
//
// Solidity: function LIGHT_CLIENT_ADDR() view returns(address)
func (_Staking *StakingCaller) LIGHTCLIENTADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "LIGHT_CLIENT_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LIGHTCLIENTADDR is a free data retrieval call binding the contract method 0xdc927faf.
//
// Solidity: function LIGHT_CLIENT_ADDR() view returns(address)
func (_Staking *StakingSession) LIGHTCLIENTADDR() (common.Address, error) {
	return _Staking.Contract.LIGHTCLIENTADDR(&_Staking.CallOpts)
}

// LIGHTCLIENTADDR is a free data retrieval call binding the contract method 0xdc927faf.
//
// Solidity: function LIGHT_CLIENT_ADDR() view returns(address)
func (_Staking *StakingCallerSession) LIGHTCLIENTADDR() (common.Address, error) {
	return _Staking.Contract.LIGHTCLIENTADDR(&_Staking.CallOpts)
}

// LOCKTIME is a free data retrieval call binding the contract method 0x413d9c3a.
//
// Solidity: function LOCK_TIME() view returns(uint256)
func (_Staking *StakingCaller) LOCKTIME(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "LOCK_TIME")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LOCKTIME is a free data retrieval call binding the contract method 0x413d9c3a.
//
// Solidity: function LOCK_TIME() view returns(uint256)
func (_Staking *StakingSession) LOCKTIME() (*big.Int, error) {
	return _Staking.Contract.LOCKTIME(&_Staking.CallOpts)
}

// LOCKTIME is a free data retrieval call binding the contract method 0x413d9c3a.
//
// Solidity: function LOCK_TIME() view returns(uint256)
func (_Staking *StakingCallerSession) LOCKTIME() (*big.Int, error) {
	return _Staking.Contract.LOCKTIME(&_Staking.CallOpts)
}

// RELAYERHUBCONTRACTADDR is a free data retrieval call binding the contract method 0xa1a11bf5.
//
// Solidity: function RELAYERHUB_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingCaller) RELAYERHUBCONTRACTADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "RELAYERHUB_CONTRACT_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RELAYERHUBCONTRACTADDR is a free data retrieval call binding the contract method 0xa1a11bf5.
//
// Solidity: function RELAYERHUB_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingSession) RELAYERHUBCONTRACTADDR() (common.Address, error) {
	return _Staking.Contract.RELAYERHUBCONTRACTADDR(&_Staking.CallOpts)
}

// RELAYERHUBCONTRACTADDR is a free data retrieval call binding the contract method 0xa1a11bf5.
//
// Solidity: function RELAYERHUB_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingCallerSession) RELAYERHUBCONTRACTADDR() (common.Address, error) {
	return _Staking.Contract.RELAYERHUBCONTRACTADDR(&_Staking.CallOpts)
}

// SLASHCHANNELID is a free data retrieval call binding the contract method 0x7942fd05.
//
// Solidity: function SLASH_CHANNELID() view returns(uint8)
func (_Staking *StakingCaller) SLASHCHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "SLASH_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// SLASHCHANNELID is a free data retrieval call binding the contract method 0x7942fd05.
//
// Solidity: function SLASH_CHANNELID() view returns(uint8)
func (_Staking *StakingSession) SLASHCHANNELID() (uint8, error) {
	return _Staking.Contract.SLASHCHANNELID(&_Staking.CallOpts)
}

// SLASHCHANNELID is a free data retrieval call binding the contract method 0x7942fd05.
//
// Solidity: function SLASH_CHANNELID() view returns(uint8)
func (_Staking *StakingCallerSession) SLASHCHANNELID() (uint8, error) {
	return _Staking.Contract.SLASHCHANNELID(&_Staking.CallOpts)
}

// SLASHCONTRACTADDR is a free data retrieval call binding the contract method 0x43756e5c.
//
// Solidity: function SLASH_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingCaller) SLASHCONTRACTADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "SLASH_CONTRACT_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SLASHCONTRACTADDR is a free data retrieval call binding the contract method 0x43756e5c.
//
// Solidity: function SLASH_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingSession) SLASHCONTRACTADDR() (common.Address, error) {
	return _Staking.Contract.SLASHCONTRACTADDR(&_Staking.CallOpts)
}

// SLASHCONTRACTADDR is a free data retrieval call binding the contract method 0x43756e5c.
//
// Solidity: function SLASH_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingCallerSession) SLASHCONTRACTADDR() (common.Address, error) {
	return _Staking.Contract.SLASHCONTRACTADDR(&_Staking.CallOpts)
}

// STAKINGCHANNELID is a free data retrieval call binding the contract method 0x4bf6c882.
//
// Solidity: function STAKING_CHANNELID() view returns(uint8)
func (_Staking *StakingCaller) STAKINGCHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "STAKING_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// STAKINGCHANNELID is a free data retrieval call binding the contract method 0x4bf6c882.
//
// Solidity: function STAKING_CHANNELID() view returns(uint8)
func (_Staking *StakingSession) STAKINGCHANNELID() (uint8, error) {
	return _Staking.Contract.STAKINGCHANNELID(&_Staking.CallOpts)
}

// STAKINGCHANNELID is a free data retrieval call binding the contract method 0x4bf6c882.
//
// Solidity: function STAKING_CHANNELID() view returns(uint8)
func (_Staking *StakingCallerSession) STAKINGCHANNELID() (uint8, error) {
	return _Staking.Contract.STAKINGCHANNELID(&_Staking.CallOpts)
}

// STAKINGCONTRACTADDR is a free data retrieval call binding the contract method 0x0e2374a5.
//
// Solidity: function STAKING_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingCaller) STAKINGCONTRACTADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "STAKING_CONTRACT_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// STAKINGCONTRACTADDR is a free data retrieval call binding the contract method 0x0e2374a5.
//
// Solidity: function STAKING_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingSession) STAKINGCONTRACTADDR() (common.Address, error) {
	return _Staking.Contract.STAKINGCONTRACTADDR(&_Staking.CallOpts)
}

// STAKINGCONTRACTADDR is a free data retrieval call binding the contract method 0x0e2374a5.
//
// Solidity: function STAKING_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingCallerSession) STAKINGCONTRACTADDR() (common.Address, error) {
	return _Staking.Contract.STAKINGCONTRACTADDR(&_Staking.CallOpts)
}

// SYSTEMREWARDADDR is a free data retrieval call binding the contract method 0xc81b1662.
//
// Solidity: function SYSTEM_REWARD_ADDR() view returns(address)
func (_Staking *StakingCaller) SYSTEMREWARDADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "SYSTEM_REWARD_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SYSTEMREWARDADDR is a free data retrieval call binding the contract method 0xc81b1662.
//
// Solidity: function SYSTEM_REWARD_ADDR() view returns(address)
func (_Staking *StakingSession) SYSTEMREWARDADDR() (common.Address, error) {
	return _Staking.Contract.SYSTEMREWARDADDR(&_Staking.CallOpts)
}

// SYSTEMREWARDADDR is a free data retrieval call binding the contract method 0xc81b1662.
//
// Solidity: function SYSTEM_REWARD_ADDR() view returns(address)
func (_Staking *StakingCallerSession) SYSTEMREWARDADDR() (common.Address, error) {
	return _Staking.Contract.SYSTEMREWARDADDR(&_Staking.CallOpts)
}

// TENDECIMALS is a free data retrieval call binding the contract method 0x5d499b1b.
//
// Solidity: function TEN_DECIMALS() view returns(uint256)
func (_Staking *StakingCaller) TENDECIMALS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "TEN_DECIMALS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TENDECIMALS is a free data retrieval call binding the contract method 0x5d499b1b.
//
// Solidity: function TEN_DECIMALS() view returns(uint256)
func (_Staking *StakingSession) TENDECIMALS() (*big.Int, error) {
	return _Staking.Contract.TENDECIMALS(&_Staking.CallOpts)
}

// TENDECIMALS is a free data retrieval call binding the contract method 0x5d499b1b.
//
// Solidity: function TEN_DECIMALS() view returns(uint256)
func (_Staking *StakingCallerSession) TENDECIMALS() (*big.Int, error) {
	return _Staking.Contract.TENDECIMALS(&_Staking.CallOpts)
}

// TOKENHUBADDR is a free data retrieval call binding the contract method 0xfd6a6879.
//
// Solidity: function TOKEN_HUB_ADDR() view returns(address)
func (_Staking *StakingCaller) TOKENHUBADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "TOKEN_HUB_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TOKENHUBADDR is a free data retrieval call binding the contract method 0xfd6a6879.
//
// Solidity: function TOKEN_HUB_ADDR() view returns(address)
func (_Staking *StakingSession) TOKENHUBADDR() (common.Address, error) {
	return _Staking.Contract.TOKENHUBADDR(&_Staking.CallOpts)
}

// TOKENHUBADDR is a free data retrieval call binding the contract method 0xfd6a6879.
//
// Solidity: function TOKEN_HUB_ADDR() view returns(address)
func (_Staking *StakingCallerSession) TOKENHUBADDR() (common.Address, error) {
	return _Staking.Contract.TOKENHUBADDR(&_Staking.CallOpts)
}

// TOKENMANAGERADDR is a free data retrieval call binding the contract method 0x75d47a0a.
//
// Solidity: function TOKEN_MANAGER_ADDR() view returns(address)
func (_Staking *StakingCaller) TOKENMANAGERADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "TOKEN_MANAGER_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TOKENMANAGERADDR is a free data retrieval call binding the contract method 0x75d47a0a.
//
// Solidity: function TOKEN_MANAGER_ADDR() view returns(address)
func (_Staking *StakingSession) TOKENMANAGERADDR() (common.Address, error) {
	return _Staking.Contract.TOKENMANAGERADDR(&_Staking.CallOpts)
}

// TOKENMANAGERADDR is a free data retrieval call binding the contract method 0x75d47a0a.
//
// Solidity: function TOKEN_MANAGER_ADDR() view returns(address)
func (_Staking *StakingCallerSession) TOKENMANAGERADDR() (common.Address, error) {
	return _Staking.Contract.TOKENMANAGERADDR(&_Staking.CallOpts)
}

// TRANSFERINCHANNELID is a free data retrieval call binding the contract method 0x70fd5bad.
//
// Solidity: function TRANSFER_IN_CHANNELID() view returns(uint8)
func (_Staking *StakingCaller) TRANSFERINCHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "TRANSFER_IN_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// TRANSFERINCHANNELID is a free data retrieval call binding the contract method 0x70fd5bad.
//
// Solidity: function TRANSFER_IN_CHANNELID() view returns(uint8)
func (_Staking *StakingSession) TRANSFERINCHANNELID() (uint8, error) {
	return _Staking.Contract.TRANSFERINCHANNELID(&_Staking.CallOpts)
}

// TRANSFERINCHANNELID is a free data retrieval call binding the contract method 0x70fd5bad.
//
// Solidity: function TRANSFER_IN_CHANNELID() view returns(uint8)
func (_Staking *StakingCallerSession) TRANSFERINCHANNELID() (uint8, error) {
	return _Staking.Contract.TRANSFERINCHANNELID(&_Staking.CallOpts)
}

// TRANSFEROUTCHANNELID is a free data retrieval call binding the contract method 0xfc3e5908.
//
// Solidity: function TRANSFER_OUT_CHANNELID() view returns(uint8)
func (_Staking *StakingCaller) TRANSFEROUTCHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "TRANSFER_OUT_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// TRANSFEROUTCHANNELID is a free data retrieval call binding the contract method 0xfc3e5908.
//
// Solidity: function TRANSFER_OUT_CHANNELID() view returns(uint8)
func (_Staking *StakingSession) TRANSFEROUTCHANNELID() (uint8, error) {
	return _Staking.Contract.TRANSFEROUTCHANNELID(&_Staking.CallOpts)
}

// TRANSFEROUTCHANNELID is a free data retrieval call binding the contract method 0xfc3e5908.
//
// Solidity: function TRANSFER_OUT_CHANNELID() view returns(uint8)
func (_Staking *StakingCallerSession) TRANSFEROUTCHANNELID() (uint8, error) {
	return _Staking.Contract.TRANSFEROUTCHANNELID(&_Staking.CallOpts)
}

// VALIDATORCONTRACTADDR is a free data retrieval call binding the contract method 0xf9a2bbc7.
//
// Solidity: function VALIDATOR_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingCaller) VALIDATORCONTRACTADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "VALIDATOR_CONTRACT_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VALIDATORCONTRACTADDR is a free data retrieval call binding the contract method 0xf9a2bbc7.
//
// Solidity: function VALIDATOR_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingSession) VALIDATORCONTRACTADDR() (common.Address, error) {
	return _Staking.Contract.VALIDATORCONTRACTADDR(&_Staking.CallOpts)
}

// VALIDATORCONTRACTADDR is a free data retrieval call binding the contract method 0xf9a2bbc7.
//
// Solidity: function VALIDATOR_CONTRACT_ADDR() view returns(address)
func (_Staking *StakingCallerSession) VALIDATORCONTRACTADDR() (common.Address, error) {
	return _Staking.Contract.VALIDATORCONTRACTADDR(&_Staking.CallOpts)
}

// AlreadyInit is a free data retrieval call binding the contract method 0xa78abc16.
//
// Solidity: function alreadyInit() view returns(bool)
func (_Staking *StakingCaller) AlreadyInit(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "alreadyInit")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AlreadyInit is a free data retrieval call binding the contract method 0xa78abc16.
//
// Solidity: function alreadyInit() view returns(bool)
func (_Staking *StakingSession) AlreadyInit() (bool, error) {
	return _Staking.Contract.AlreadyInit(&_Staking.CallOpts)
}

// AlreadyInit is a free data retrieval call binding the contract method 0xa78abc16.
//
// Solidity: function alreadyInit() view returns(bool)
func (_Staking *StakingCallerSession) AlreadyInit() (bool, error) {
	return _Staking.Contract.AlreadyInit(&_Staking.CallOpts)
}

// BSCRelayerFee is a free data retrieval call binding the contract method 0x5d17c8bd.
//
// Solidity: function bSCRelayerFee() view returns(uint256)
func (_Staking *StakingCaller) BSCRelayerFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "bSCRelayerFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BSCRelayerFee is a free data retrieval call binding the contract method 0x5d17c8bd.
//
// Solidity: function bSCRelayerFee() view returns(uint256)
func (_Staking *StakingSession) BSCRelayerFee() (*big.Int, error) {
	return _Staking.Contract.BSCRelayerFee(&_Staking.CallOpts)
}

// BSCRelayerFee is a free data retrieval call binding the contract method 0x5d17c8bd.
//
// Solidity: function bSCRelayerFee() view returns(uint256)
func (_Staking *StakingCallerSession) BSCRelayerFee() (*big.Int, error) {
	return _Staking.Contract.BSCRelayerFee(&_Staking.CallOpts)
}

// BscChainID is a free data retrieval call binding the contract method 0x493279b1.
//
// Solidity: function bscChainID() view returns(uint16)
func (_Staking *StakingCaller) BscChainID(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "bscChainID")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// BscChainID is a free data retrieval call binding the contract method 0x493279b1.
//
// Solidity: function bscChainID() view returns(uint16)
func (_Staking *StakingSession) BscChainID() (uint16, error) {
	return _Staking.Contract.BscChainID(&_Staking.CallOpts)
}

// BscChainID is a free data retrieval call binding the contract method 0x493279b1.
//
// Solidity: function bscChainID() view returns(uint16)
func (_Staking *StakingCallerSession) BscChainID() (uint16, error) {
	return _Staking.Contract.BscChainID(&_Staking.CallOpts)
}

// GetDelegated is a free data retrieval call binding the contract method 0xd61b9b93.
//
// Solidity: function getDelegated(address delegator, address validator) view returns(uint256)
func (_Staking *StakingCaller) GetDelegated(opts *bind.CallOpts, delegator common.Address, validator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getDelegated", delegator, validator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDelegated is a free data retrieval call binding the contract method 0xd61b9b93.
//
// Solidity: function getDelegated(address delegator, address validator) view returns(uint256)
func (_Staking *StakingSession) GetDelegated(delegator common.Address, validator common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDelegated(&_Staking.CallOpts, delegator, validator)
}

// GetDelegated is a free data retrieval call binding the contract method 0xd61b9b93.
//
// Solidity: function getDelegated(address delegator, address validator) view returns(uint256)
func (_Staking *StakingCallerSession) GetDelegated(delegator common.Address, validator common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDelegated(&_Staking.CallOpts, delegator, validator)
}

// GetDistributedReward is a free data retrieval call binding the contract method 0x11fe9ec6.
//
// Solidity: function getDistributedReward(address delegator) view returns(uint256)
func (_Staking *StakingCaller) GetDistributedReward(opts *bind.CallOpts, delegator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getDistributedReward", delegator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDistributedReward is a free data retrieval call binding the contract method 0x11fe9ec6.
//
// Solidity: function getDistributedReward(address delegator) view returns(uint256)
func (_Staking *StakingSession) GetDistributedReward(delegator common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDistributedReward(&_Staking.CallOpts, delegator)
}

// GetDistributedReward is a free data retrieval call binding the contract method 0x11fe9ec6.
//
// Solidity: function getDistributedReward(address delegator) view returns(uint256)
func (_Staking *StakingCallerSession) GetDistributedReward(delegator common.Address) (*big.Int, error) {
	return _Staking.Contract.GetDistributedReward(&_Staking.CallOpts, delegator)
}

// GetMinDelegation is a free data retrieval call binding the contract method 0x69b635b6.
//
// Solidity: function getMinDelegation() view returns(uint256)
func (_Staking *StakingCaller) GetMinDelegation(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getMinDelegation")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinDelegation is a free data retrieval call binding the contract method 0x69b635b6.
//
// Solidity: function getMinDelegation() view returns(uint256)
func (_Staking *StakingSession) GetMinDelegation() (*big.Int, error) {
	return _Staking.Contract.GetMinDelegation(&_Staking.CallOpts)
}

// GetMinDelegation is a free data retrieval call binding the contract method 0x69b635b6.
//
// Solidity: function getMinDelegation() view returns(uint256)
func (_Staking *StakingCallerSession) GetMinDelegation() (*big.Int, error) {
	return _Staking.Contract.GetMinDelegation(&_Staking.CallOpts)
}

// GetPendingRedelegateTime is a free data retrieval call binding the contract method 0xf45fd80b.
//
// Solidity: function getPendingRedelegateTime(address delegator, address valSrc, address valDst) view returns(uint256)
func (_Staking *StakingCaller) GetPendingRedelegateTime(opts *bind.CallOpts, delegator common.Address, valSrc common.Address, valDst common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getPendingRedelegateTime", delegator, valSrc, valDst)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPendingRedelegateTime is a free data retrieval call binding the contract method 0xf45fd80b.
//
// Solidity: function getPendingRedelegateTime(address delegator, address valSrc, address valDst) view returns(uint256)
func (_Staking *StakingSession) GetPendingRedelegateTime(delegator common.Address, valSrc common.Address, valDst common.Address) (*big.Int, error) {
	return _Staking.Contract.GetPendingRedelegateTime(&_Staking.CallOpts, delegator, valSrc, valDst)
}

// GetPendingRedelegateTime is a free data retrieval call binding the contract method 0xf45fd80b.
//
// Solidity: function getPendingRedelegateTime(address delegator, address valSrc, address valDst) view returns(uint256)
func (_Staking *StakingCallerSession) GetPendingRedelegateTime(delegator common.Address, valSrc common.Address, valDst common.Address) (*big.Int, error) {
	return _Staking.Contract.GetPendingRedelegateTime(&_Staking.CallOpts, delegator, valSrc, valDst)
}

// GetPendingUndelegateTime is a free data retrieval call binding the contract method 0xbf8546ca.
//
// Solidity: function getPendingUndelegateTime(address delegator, address validator) view returns(uint256)
func (_Staking *StakingCaller) GetPendingUndelegateTime(opts *bind.CallOpts, delegator common.Address, validator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getPendingUndelegateTime", delegator, validator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPendingUndelegateTime is a free data retrieval call binding the contract method 0xbf8546ca.
//
// Solidity: function getPendingUndelegateTime(address delegator, address validator) view returns(uint256)
func (_Staking *StakingSession) GetPendingUndelegateTime(delegator common.Address, validator common.Address) (*big.Int, error) {
	return _Staking.Contract.GetPendingUndelegateTime(&_Staking.CallOpts, delegator, validator)
}

// GetPendingUndelegateTime is a free data retrieval call binding the contract method 0xbf8546ca.
//
// Solidity: function getPendingUndelegateTime(address delegator, address validator) view returns(uint256)
func (_Staking *StakingCallerSession) GetPendingUndelegateTime(delegator common.Address, validator common.Address) (*big.Int, error) {
	return _Staking.Contract.GetPendingUndelegateTime(&_Staking.CallOpts, delegator, validator)
}

// GetRelayerFee is a free data retrieval call binding the contract method 0xc2117d82.
//
// Solidity: function getRelayerFee() view returns(uint256)
func (_Staking *StakingCaller) GetRelayerFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getRelayerFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRelayerFee is a free data retrieval call binding the contract method 0xc2117d82.
//
// Solidity: function getRelayerFee() view returns(uint256)
func (_Staking *StakingSession) GetRelayerFee() (*big.Int, error) {
	return _Staking.Contract.GetRelayerFee(&_Staking.CallOpts)
}

// GetRelayerFee is a free data retrieval call binding the contract method 0xc2117d82.
//
// Solidity: function getRelayerFee() view returns(uint256)
func (_Staking *StakingCallerSession) GetRelayerFee() (*big.Int, error) {
	return _Staking.Contract.GetRelayerFee(&_Staking.CallOpts)
}

// GetRequestInFly is a free data retrieval call binding the contract method 0x047636d1.
//
// Solidity: function getRequestInFly(address delegator) view returns(uint256[3])
func (_Staking *StakingCaller) GetRequestInFly(opts *bind.CallOpts, delegator common.Address) ([3]*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getRequestInFly", delegator)

	if err != nil {
		return *new([3]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([3]*big.Int)).(*[3]*big.Int)

	return out0, err

}

// GetRequestInFly is a free data retrieval call binding the contract method 0x047636d1.
//
// Solidity: function getRequestInFly(address delegator) view returns(uint256[3])
func (_Staking *StakingSession) GetRequestInFly(delegator common.Address) ([3]*big.Int, error) {
	return _Staking.Contract.GetRequestInFly(&_Staking.CallOpts, delegator)
}

// GetRequestInFly is a free data retrieval call binding the contract method 0x047636d1.
//
// Solidity: function getRequestInFly(address delegator) view returns(uint256[3])
func (_Staking *StakingCallerSession) GetRequestInFly(delegator common.Address) ([3]*big.Int, error) {
	return _Staking.Contract.GetRequestInFly(&_Staking.CallOpts, delegator)
}

// GetTotalDelegated is a free data retrieval call binding the contract method 0x6fb7f7eb.
//
// Solidity: function getTotalDelegated(address delegator) view returns(uint256)
func (_Staking *StakingCaller) GetTotalDelegated(opts *bind.CallOpts, delegator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getTotalDelegated", delegator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalDelegated is a free data retrieval call binding the contract method 0x6fb7f7eb.
//
// Solidity: function getTotalDelegated(address delegator) view returns(uint256)
func (_Staking *StakingSession) GetTotalDelegated(delegator common.Address) (*big.Int, error) {
	return _Staking.Contract.GetTotalDelegated(&_Staking.CallOpts, delegator)
}

// GetTotalDelegated is a free data retrieval call binding the contract method 0x6fb7f7eb.
//
// Solidity: function getTotalDelegated(address delegator) view returns(uint256)
func (_Staking *StakingCallerSession) GetTotalDelegated(delegator common.Address) (*big.Int, error) {
	return _Staking.Contract.GetTotalDelegated(&_Staking.CallOpts, delegator)
}

// GetUndelegated is a free data retrieval call binding the contract method 0x75aca593.
//
// Solidity: function getUndelegated(address delegator) view returns(uint256)
func (_Staking *StakingCaller) GetUndelegated(opts *bind.CallOpts, delegator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "getUndelegated", delegator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUndelegated is a free data retrieval call binding the contract method 0x75aca593.
//
// Solidity: function getUndelegated(address delegator) view returns(uint256)
func (_Staking *StakingSession) GetUndelegated(delegator common.Address) (*big.Int, error) {
	return _Staking.Contract.GetUndelegated(&_Staking.CallOpts, delegator)
}

// GetUndelegated is a free data retrieval call binding the contract method 0x75aca593.
//
// Solidity: function getUndelegated(address delegator) view returns(uint256)
func (_Staking *StakingCallerSession) GetUndelegated(delegator common.Address) (*big.Int, error) {
	return _Staking.Contract.GetUndelegated(&_Staking.CallOpts, delegator)
}

// MinDelegation is a free data retrieval call binding the contract method 0x02985992.
//
// Solidity: function minDelegation() view returns(uint256)
func (_Staking *StakingCaller) MinDelegation(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "minDelegation")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinDelegation is a free data retrieval call binding the contract method 0x02985992.
//
// Solidity: function minDelegation() view returns(uint256)
func (_Staking *StakingSession) MinDelegation() (*big.Int, error) {
	return _Staking.Contract.MinDelegation(&_Staking.CallOpts)
}

// MinDelegation is a free data retrieval call binding the contract method 0x02985992.
//
// Solidity: function minDelegation() view returns(uint256)
func (_Staking *StakingCallerSession) MinDelegation() (*big.Int, error) {
	return _Staking.Contract.MinDelegation(&_Staking.CallOpts)
}

// RelayerFee is a free data retrieval call binding the contract method 0x2fdeb111.
//
// Solidity: function relayerFee() view returns(uint256)
func (_Staking *StakingCaller) RelayerFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "relayerFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RelayerFee is a free data retrieval call binding the contract method 0x2fdeb111.
//
// Solidity: function relayerFee() view returns(uint256)
func (_Staking *StakingSession) RelayerFee() (*big.Int, error) {
	return _Staking.Contract.RelayerFee(&_Staking.CallOpts)
}

// RelayerFee is a free data retrieval call binding the contract method 0x2fdeb111.
//
// Solidity: function relayerFee() view returns(uint256)
func (_Staking *StakingCallerSession) RelayerFee() (*big.Int, error) {
	return _Staking.Contract.RelayerFee(&_Staking.CallOpts)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xb88a802f.
//
// Solidity: function claimReward() returns(uint256 amount)
func (_Staking *StakingTransactor) ClaimReward(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "claimReward")
}

// ClaimReward is a paid mutator transaction binding the contract method 0xb88a802f.
//
// Solidity: function claimReward() returns(uint256 amount)
func (_Staking *StakingSession) ClaimReward() (*types.Transaction, error) {
	return _Staking.Contract.ClaimReward(&_Staking.TransactOpts)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xb88a802f.
//
// Solidity: function claimReward() returns(uint256 amount)
func (_Staking *StakingTransactorSession) ClaimReward() (*types.Transaction, error) {
	return _Staking.Contract.ClaimReward(&_Staking.TransactOpts)
}

// ClaimUndelegated is a paid mutator transaction binding the contract method 0x62b171d2.
//
// Solidity: function claimUndelegated() returns(uint256 amount)
func (_Staking *StakingTransactor) ClaimUndelegated(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "claimUndelegated")
}

// ClaimUndelegated is a paid mutator transaction binding the contract method 0x62b171d2.
//
// Solidity: function claimUndelegated() returns(uint256 amount)
func (_Staking *StakingSession) ClaimUndelegated() (*types.Transaction, error) {
	return _Staking.Contract.ClaimUndelegated(&_Staking.TransactOpts)
}

// ClaimUndelegated is a paid mutator transaction binding the contract method 0x62b171d2.
//
// Solidity: function claimUndelegated() returns(uint256 amount)
func (_Staking *StakingTransactorSession) ClaimUndelegated() (*types.Transaction, error) {
	return _Staking.Contract.ClaimUndelegated(&_Staking.TransactOpts)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address validator, uint256 amount) payable returns()
func (_Staking *StakingTransactor) Delegate(opts *bind.TransactOpts, validator common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "delegate", validator, amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address validator, uint256 amount) payable returns()
func (_Staking *StakingSession) Delegate(validator common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Delegate(&_Staking.TransactOpts, validator, amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address validator, uint256 amount) payable returns()
func (_Staking *StakingTransactorSession) Delegate(validator common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Delegate(&_Staking.TransactOpts, validator, amount)
}

// HandleAckPackage is a paid mutator transaction binding the contract method 0x831d65d1.
//
// Solidity: function handleAckPackage(uint8 , bytes msgBytes) returns()
func (_Staking *StakingTransactor) HandleAckPackage(opts *bind.TransactOpts, arg0 uint8, msgBytes []byte) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "handleAckPackage", arg0, msgBytes)
}

// HandleAckPackage is a paid mutator transaction binding the contract method 0x831d65d1.
//
// Solidity: function handleAckPackage(uint8 , bytes msgBytes) returns()
func (_Staking *StakingSession) HandleAckPackage(arg0 uint8, msgBytes []byte) (*types.Transaction, error) {
	return _Staking.Contract.HandleAckPackage(&_Staking.TransactOpts, arg0, msgBytes)
}

// HandleAckPackage is a paid mutator transaction binding the contract method 0x831d65d1.
//
// Solidity: function handleAckPackage(uint8 , bytes msgBytes) returns()
func (_Staking *StakingTransactorSession) HandleAckPackage(arg0 uint8, msgBytes []byte) (*types.Transaction, error) {
	return _Staking.Contract.HandleAckPackage(&_Staking.TransactOpts, arg0, msgBytes)
}

// HandleFailAckPackage is a paid mutator transaction binding the contract method 0xc8509d81.
//
// Solidity: function handleFailAckPackage(uint8 , bytes msgBytes) returns()
func (_Staking *StakingTransactor) HandleFailAckPackage(opts *bind.TransactOpts, arg0 uint8, msgBytes []byte) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "handleFailAckPackage", arg0, msgBytes)
}

// HandleFailAckPackage is a paid mutator transaction binding the contract method 0xc8509d81.
//
// Solidity: function handleFailAckPackage(uint8 , bytes msgBytes) returns()
func (_Staking *StakingSession) HandleFailAckPackage(arg0 uint8, msgBytes []byte) (*types.Transaction, error) {
	return _Staking.Contract.HandleFailAckPackage(&_Staking.TransactOpts, arg0, msgBytes)
}

// HandleFailAckPackage is a paid mutator transaction binding the contract method 0xc8509d81.
//
// Solidity: function handleFailAckPackage(uint8 , bytes msgBytes) returns()
func (_Staking *StakingTransactorSession) HandleFailAckPackage(arg0 uint8, msgBytes []byte) (*types.Transaction, error) {
	return _Staking.Contract.HandleFailAckPackage(&_Staking.TransactOpts, arg0, msgBytes)
}

// HandleSynPackage is a paid mutator transaction binding the contract method 0x1182b875.
//
// Solidity: function handleSynPackage(uint8 , bytes msgBytes) returns(bytes)
func (_Staking *StakingTransactor) HandleSynPackage(opts *bind.TransactOpts, arg0 uint8, msgBytes []byte) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "handleSynPackage", arg0, msgBytes)
}

// HandleSynPackage is a paid mutator transaction binding the contract method 0x1182b875.
//
// Solidity: function handleSynPackage(uint8 , bytes msgBytes) returns(bytes)
func (_Staking *StakingSession) HandleSynPackage(arg0 uint8, msgBytes []byte) (*types.Transaction, error) {
	return _Staking.Contract.HandleSynPackage(&_Staking.TransactOpts, arg0, msgBytes)
}

// HandleSynPackage is a paid mutator transaction binding the contract method 0x1182b875.
//
// Solidity: function handleSynPackage(uint8 , bytes msgBytes) returns(bytes)
func (_Staking *StakingTransactorSession) HandleSynPackage(arg0 uint8, msgBytes []byte) (*types.Transaction, error) {
	return _Staking.Contract.HandleSynPackage(&_Staking.TransactOpts, arg0, msgBytes)
}

// Redelegate is a paid mutator transaction binding the contract method 0x6bd8f804.
//
// Solidity: function redelegate(address validatorSrc, address validatorDst, uint256 amount) payable returns()
func (_Staking *StakingTransactor) Redelegate(opts *bind.TransactOpts, validatorSrc common.Address, validatorDst common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "redelegate", validatorSrc, validatorDst, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x6bd8f804.
//
// Solidity: function redelegate(address validatorSrc, address validatorDst, uint256 amount) payable returns()
func (_Staking *StakingSession) Redelegate(validatorSrc common.Address, validatorDst common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Redelegate(&_Staking.TransactOpts, validatorSrc, validatorDst, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x6bd8f804.
//
// Solidity: function redelegate(address validatorSrc, address validatorDst, uint256 amount) payable returns()
func (_Staking *StakingTransactorSession) Redelegate(validatorSrc common.Address, validatorDst common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Redelegate(&_Staking.TransactOpts, validatorSrc, validatorDst, amount)
}

// Undelegate is a paid mutator transaction binding the contract method 0x4d99dd16.
//
// Solidity: function undelegate(address validator, uint256 amount) payable returns()
func (_Staking *StakingTransactor) Undelegate(opts *bind.TransactOpts, validator common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "undelegate", validator, amount)
}

// Undelegate is a paid mutator transaction binding the contract method 0x4d99dd16.
//
// Solidity: function undelegate(address validator, uint256 amount) payable returns()
func (_Staking *StakingSession) Undelegate(validator common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Undelegate(&_Staking.TransactOpts, validator, amount)
}

// Undelegate is a paid mutator transaction binding the contract method 0x4d99dd16.
//
// Solidity: function undelegate(address validator, uint256 amount) payable returns()
func (_Staking *StakingTransactorSession) Undelegate(validator common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Undelegate(&_Staking.TransactOpts, validator, amount)
}

// UpdateParam is a paid mutator transaction binding the contract method 0xac431751.
//
// Solidity: function updateParam(string key, bytes value) returns()
func (_Staking *StakingTransactor) UpdateParam(opts *bind.TransactOpts, key string, value []byte) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "updateParam", key, value)
}

// UpdateParam is a paid mutator transaction binding the contract method 0xac431751.
//
// Solidity: function updateParam(string key, bytes value) returns()
func (_Staking *StakingSession) UpdateParam(key string, value []byte) (*types.Transaction, error) {
	return _Staking.Contract.UpdateParam(&_Staking.TransactOpts, key, value)
}

// UpdateParam is a paid mutator transaction binding the contract method 0xac431751.
//
// Solidity: function updateParam(string key, bytes value) returns()
func (_Staking *StakingTransactorSession) UpdateParam(key string, value []byte) (*types.Transaction, error) {
	return _Staking.Contract.UpdateParam(&_Staking.TransactOpts, key, value)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Staking *StakingTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Staking *StakingSession) Receive() (*types.Transaction, error) {
	return _Staking.Contract.Receive(&_Staking.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Staking *StakingTransactorSession) Receive() (*types.Transaction, error) {
	return _Staking.Contract.Receive(&_Staking.TransactOpts)
}

// StakingCrashResponseIterator is returned from FilterCrashResponse and is used to iterate over the raw logs and unpacked data for CrashResponse events raised by the Staking contract.
type StakingCrashResponseIterator struct {
	Event *StakingCrashResponse // Event containing the contract specifics and raw log

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
func (it *StakingCrashResponseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingCrashResponse)
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
		it.Event = new(StakingCrashResponse)
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
func (it *StakingCrashResponseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingCrashResponseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingCrashResponse represents a CrashResponse event raised by the Staking contract.
type StakingCrashResponse struct {
	EventType uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCrashResponse is a free log retrieval operation binding the contract event 0xf83de021914a4585482db5ca47d520a5657165b443fa2c7ef8ed4635f054da9b.
//
// Solidity: event crashResponse(uint8 indexed eventType)
func (_Staking *StakingFilterer) FilterCrashResponse(opts *bind.FilterOpts, eventType []uint8) (*StakingCrashResponseIterator, error) {

	var eventTypeRule []interface{}
	for _, eventTypeItem := range eventType {
		eventTypeRule = append(eventTypeRule, eventTypeItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "crashResponse", eventTypeRule)
	if err != nil {
		return nil, err
	}
	return &StakingCrashResponseIterator{contract: _Staking.contract, event: "crashResponse", logs: logs, sub: sub}, nil
}

// WatchCrashResponse is a free log subscription operation binding the contract event 0xf83de021914a4585482db5ca47d520a5657165b443fa2c7ef8ed4635f054da9b.
//
// Solidity: event crashResponse(uint8 indexed eventType)
func (_Staking *StakingFilterer) WatchCrashResponse(opts *bind.WatchOpts, sink chan<- *StakingCrashResponse, eventType []uint8) (event.Subscription, error) {

	var eventTypeRule []interface{}
	for _, eventTypeItem := range eventType {
		eventTypeRule = append(eventTypeRule, eventTypeItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "crashResponse", eventTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingCrashResponse)
				if err := _Staking.contract.UnpackLog(event, "crashResponse", log); err != nil {
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

// ParseCrashResponse is a log parse operation binding the contract event 0xf83de021914a4585482db5ca47d520a5657165b443fa2c7ef8ed4635f054da9b.
//
// Solidity: event crashResponse(uint8 indexed eventType)
func (_Staking *StakingFilterer) ParseCrashResponse(log types.Log) (*StakingCrashResponse, error) {
	event := new(StakingCrashResponse)
	if err := _Staking.contract.UnpackLog(event, "crashResponse", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingDelegateFailedIterator is returned from FilterDelegateFailed and is used to iterate over the raw logs and unpacked data for DelegateFailed events raised by the Staking contract.
type StakingDelegateFailedIterator struct {
	Event *StakingDelegateFailed // Event containing the contract specifics and raw log

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
func (it *StakingDelegateFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingDelegateFailed)
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
		it.Event = new(StakingDelegateFailed)
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
func (it *StakingDelegateFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingDelegateFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingDelegateFailed represents a DelegateFailed event raised by the Staking contract.
type StakingDelegateFailed struct {
	Delegator common.Address
	Validator common.Address
	Amount    *big.Int
	ErrCode   uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDelegateFailed is a free log retrieval operation binding the contract event 0xcbd481ae600289fad8c0484d07ce0ffe4f010d7c844ecfdeaf2a13fead52886e.
//
// Solidity: event delegateFailed(address indexed delegator, address indexed validator, uint256 amount, uint8 errCode)
func (_Staking *StakingFilterer) FilterDelegateFailed(opts *bind.FilterOpts, delegator []common.Address, validator []common.Address) (*StakingDelegateFailedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "delegateFailed", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingDelegateFailedIterator{contract: _Staking.contract, event: "delegateFailed", logs: logs, sub: sub}, nil
}

// WatchDelegateFailed is a free log subscription operation binding the contract event 0xcbd481ae600289fad8c0484d07ce0ffe4f010d7c844ecfdeaf2a13fead52886e.
//
// Solidity: event delegateFailed(address indexed delegator, address indexed validator, uint256 amount, uint8 errCode)
func (_Staking *StakingFilterer) WatchDelegateFailed(opts *bind.WatchOpts, sink chan<- *StakingDelegateFailed, delegator []common.Address, validator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "delegateFailed", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingDelegateFailed)
				if err := _Staking.contract.UnpackLog(event, "delegateFailed", log); err != nil {
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

// ParseDelegateFailed is a log parse operation binding the contract event 0xcbd481ae600289fad8c0484d07ce0ffe4f010d7c844ecfdeaf2a13fead52886e.
//
// Solidity: event delegateFailed(address indexed delegator, address indexed validator, uint256 amount, uint8 errCode)
func (_Staking *StakingFilterer) ParseDelegateFailed(log types.Log) (*StakingDelegateFailed, error) {
	event := new(StakingDelegateFailed)
	if err := _Staking.contract.UnpackLog(event, "delegateFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingDelegateSubmittedIterator is returned from FilterDelegateSubmitted and is used to iterate over the raw logs and unpacked data for DelegateSubmitted events raised by the Staking contract.
type StakingDelegateSubmittedIterator struct {
	Event *StakingDelegateSubmitted // Event containing the contract specifics and raw log

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
func (it *StakingDelegateSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingDelegateSubmitted)
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
		it.Event = new(StakingDelegateSubmitted)
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
func (it *StakingDelegateSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingDelegateSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingDelegateSubmitted represents a DelegateSubmitted event raised by the Staking contract.
type StakingDelegateSubmitted struct {
	Delegator  common.Address
	Validator  common.Address
	Amount     *big.Int
	RelayerFee *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDelegateSubmitted is a free log retrieval operation binding the contract event 0x5f32ed2794e2e72d19e3cb2320e8820a499c4204887372beba51f5e61c040867.
//
// Solidity: event delegateSubmitted(address indexed delegator, address indexed validator, uint256 amount, uint256 relayerFee)
func (_Staking *StakingFilterer) FilterDelegateSubmitted(opts *bind.FilterOpts, delegator []common.Address, validator []common.Address) (*StakingDelegateSubmittedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "delegateSubmitted", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingDelegateSubmittedIterator{contract: _Staking.contract, event: "delegateSubmitted", logs: logs, sub: sub}, nil
}

// WatchDelegateSubmitted is a free log subscription operation binding the contract event 0x5f32ed2794e2e72d19e3cb2320e8820a499c4204887372beba51f5e61c040867.
//
// Solidity: event delegateSubmitted(address indexed delegator, address indexed validator, uint256 amount, uint256 relayerFee)
func (_Staking *StakingFilterer) WatchDelegateSubmitted(opts *bind.WatchOpts, sink chan<- *StakingDelegateSubmitted, delegator []common.Address, validator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "delegateSubmitted", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingDelegateSubmitted)
				if err := _Staking.contract.UnpackLog(event, "delegateSubmitted", log); err != nil {
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

// ParseDelegateSubmitted is a log parse operation binding the contract event 0x5f32ed2794e2e72d19e3cb2320e8820a499c4204887372beba51f5e61c040867.
//
// Solidity: event delegateSubmitted(address indexed delegator, address indexed validator, uint256 amount, uint256 relayerFee)
func (_Staking *StakingFilterer) ParseDelegateSubmitted(log types.Log) (*StakingDelegateSubmitted, error) {
	event := new(StakingDelegateSubmitted)
	if err := _Staking.contract.UnpackLog(event, "delegateSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingDelegateSuccessIterator is returned from FilterDelegateSuccess and is used to iterate over the raw logs and unpacked data for DelegateSuccess events raised by the Staking contract.
type StakingDelegateSuccessIterator struct {
	Event *StakingDelegateSuccess // Event containing the contract specifics and raw log

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
func (it *StakingDelegateSuccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingDelegateSuccess)
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
		it.Event = new(StakingDelegateSuccess)
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
func (it *StakingDelegateSuccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingDelegateSuccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingDelegateSuccess represents a DelegateSuccess event raised by the Staking contract.
type StakingDelegateSuccess struct {
	Delegator common.Address
	Validator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDelegateSuccess is a free log retrieval operation binding the contract event 0x9a57c81564ab02642f34fd87e41baa9b074c18342cec3b7268b62bf752018fd1.
//
// Solidity: event delegateSuccess(address indexed delegator, address indexed validator, uint256 amount)
func (_Staking *StakingFilterer) FilterDelegateSuccess(opts *bind.FilterOpts, delegator []common.Address, validator []common.Address) (*StakingDelegateSuccessIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "delegateSuccess", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingDelegateSuccessIterator{contract: _Staking.contract, event: "delegateSuccess", logs: logs, sub: sub}, nil
}

// WatchDelegateSuccess is a free log subscription operation binding the contract event 0x9a57c81564ab02642f34fd87e41baa9b074c18342cec3b7268b62bf752018fd1.
//
// Solidity: event delegateSuccess(address indexed delegator, address indexed validator, uint256 amount)
func (_Staking *StakingFilterer) WatchDelegateSuccess(opts *bind.WatchOpts, sink chan<- *StakingDelegateSuccess, delegator []common.Address, validator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "delegateSuccess", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingDelegateSuccess)
				if err := _Staking.contract.UnpackLog(event, "delegateSuccess", log); err != nil {
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

// ParseDelegateSuccess is a log parse operation binding the contract event 0x9a57c81564ab02642f34fd87e41baa9b074c18342cec3b7268b62bf752018fd1.
//
// Solidity: event delegateSuccess(address indexed delegator, address indexed validator, uint256 amount)
func (_Staking *StakingFilterer) ParseDelegateSuccess(log types.Log) (*StakingDelegateSuccess, error) {
	event := new(StakingDelegateSuccess)
	if err := _Staking.contract.UnpackLog(event, "delegateSuccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingFailedSynPackageIterator is returned from FilterFailedSynPackage and is used to iterate over the raw logs and unpacked data for FailedSynPackage events raised by the Staking contract.
type StakingFailedSynPackageIterator struct {
	Event *StakingFailedSynPackage // Event containing the contract specifics and raw log

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
func (it *StakingFailedSynPackageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingFailedSynPackage)
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
		it.Event = new(StakingFailedSynPackage)
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
func (it *StakingFailedSynPackageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingFailedSynPackageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingFailedSynPackage represents a FailedSynPackage event raised by the Staking contract.
type StakingFailedSynPackage struct {
	EventType uint8
	ErrCode   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFailedSynPackage is a free log retrieval operation binding the contract event 0x391d6e5ea6ab6c49b9a0abb1782cae5def8d711f973b00c729658c0b2a80b31b.
//
// Solidity: event failedSynPackage(uint8 indexed eventType, uint256 errCode)
func (_Staking *StakingFilterer) FilterFailedSynPackage(opts *bind.FilterOpts, eventType []uint8) (*StakingFailedSynPackageIterator, error) {

	var eventTypeRule []interface{}
	for _, eventTypeItem := range eventType {
		eventTypeRule = append(eventTypeRule, eventTypeItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "failedSynPackage", eventTypeRule)
	if err != nil {
		return nil, err
	}
	return &StakingFailedSynPackageIterator{contract: _Staking.contract, event: "failedSynPackage", logs: logs, sub: sub}, nil
}

// WatchFailedSynPackage is a free log subscription operation binding the contract event 0x391d6e5ea6ab6c49b9a0abb1782cae5def8d711f973b00c729658c0b2a80b31b.
//
// Solidity: event failedSynPackage(uint8 indexed eventType, uint256 errCode)
func (_Staking *StakingFilterer) WatchFailedSynPackage(opts *bind.WatchOpts, sink chan<- *StakingFailedSynPackage, eventType []uint8) (event.Subscription, error) {

	var eventTypeRule []interface{}
	for _, eventTypeItem := range eventType {
		eventTypeRule = append(eventTypeRule, eventTypeItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "failedSynPackage", eventTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingFailedSynPackage)
				if err := _Staking.contract.UnpackLog(event, "failedSynPackage", log); err != nil {
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

// ParseFailedSynPackage is a log parse operation binding the contract event 0x391d6e5ea6ab6c49b9a0abb1782cae5def8d711f973b00c729658c0b2a80b31b.
//
// Solidity: event failedSynPackage(uint8 indexed eventType, uint256 errCode)
func (_Staking *StakingFilterer) ParseFailedSynPackage(log types.Log) (*StakingFailedSynPackage, error) {
	event := new(StakingFailedSynPackage)
	if err := _Staking.contract.UnpackLog(event, "failedSynPackage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingParamChangeIterator is returned from FilterParamChange and is used to iterate over the raw logs and unpacked data for ParamChange events raised by the Staking contract.
type StakingParamChangeIterator struct {
	Event *StakingParamChange // Event containing the contract specifics and raw log

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
func (it *StakingParamChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingParamChange)
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
		it.Event = new(StakingParamChange)
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
func (it *StakingParamChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingParamChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingParamChange represents a ParamChange event raised by the Staking contract.
type StakingParamChange struct {
	Key   string
	Value []byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterParamChange is a free log retrieval operation binding the contract event 0x6cdb0ac70ab7f2e2d035cca5be60d89906f2dede7648ddbd7402189c1eeed17a.
//
// Solidity: event paramChange(string key, bytes value)
func (_Staking *StakingFilterer) FilterParamChange(opts *bind.FilterOpts) (*StakingParamChangeIterator, error) {

	logs, sub, err := _Staking.contract.FilterLogs(opts, "paramChange")
	if err != nil {
		return nil, err
	}
	return &StakingParamChangeIterator{contract: _Staking.contract, event: "paramChange", logs: logs, sub: sub}, nil
}

// WatchParamChange is a free log subscription operation binding the contract event 0x6cdb0ac70ab7f2e2d035cca5be60d89906f2dede7648ddbd7402189c1eeed17a.
//
// Solidity: event paramChange(string key, bytes value)
func (_Staking *StakingFilterer) WatchParamChange(opts *bind.WatchOpts, sink chan<- *StakingParamChange) (event.Subscription, error) {

	logs, sub, err := _Staking.contract.WatchLogs(opts, "paramChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingParamChange)
				if err := _Staking.contract.UnpackLog(event, "paramChange", log); err != nil {
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

// ParseParamChange is a log parse operation binding the contract event 0x6cdb0ac70ab7f2e2d035cca5be60d89906f2dede7648ddbd7402189c1eeed17a.
//
// Solidity: event paramChange(string key, bytes value)
func (_Staking *StakingFilterer) ParseParamChange(log types.Log) (*StakingParamChange, error) {
	event := new(StakingParamChange)
	if err := _Staking.contract.UnpackLog(event, "paramChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingRedelegateFailedIterator is returned from FilterRedelegateFailed and is used to iterate over the raw logs and unpacked data for RedelegateFailed events raised by the Staking contract.
type StakingRedelegateFailedIterator struct {
	Event *StakingRedelegateFailed // Event containing the contract specifics and raw log

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
func (it *StakingRedelegateFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingRedelegateFailed)
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
		it.Event = new(StakingRedelegateFailed)
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
func (it *StakingRedelegateFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingRedelegateFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingRedelegateFailed represents a RedelegateFailed event raised by the Staking contract.
type StakingRedelegateFailed struct {
	Delegator common.Address
	ValSrc    common.Address
	ValDst    common.Address
	Amount    *big.Int
	ErrCode   uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRedelegateFailed is a free log retrieval operation binding the contract event 0xb93bee5c59f85ede6b074a99f4ffcd3e3fc0d5c3d8156de331de89a49e0ce77c.
//
// Solidity: event redelegateFailed(address indexed delegator, address indexed valSrc, address indexed valDst, uint256 amount, uint8 errCode)
func (_Staking *StakingFilterer) FilterRedelegateFailed(opts *bind.FilterOpts, delegator []common.Address, valSrc []common.Address, valDst []common.Address) (*StakingRedelegateFailedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var valSrcRule []interface{}
	for _, valSrcItem := range valSrc {
		valSrcRule = append(valSrcRule, valSrcItem)
	}
	var valDstRule []interface{}
	for _, valDstItem := range valDst {
		valDstRule = append(valDstRule, valDstItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "redelegateFailed", delegatorRule, valSrcRule, valDstRule)
	if err != nil {
		return nil, err
	}
	return &StakingRedelegateFailedIterator{contract: _Staking.contract, event: "redelegateFailed", logs: logs, sub: sub}, nil
}

// WatchRedelegateFailed is a free log subscription operation binding the contract event 0xb93bee5c59f85ede6b074a99f4ffcd3e3fc0d5c3d8156de331de89a49e0ce77c.
//
// Solidity: event redelegateFailed(address indexed delegator, address indexed valSrc, address indexed valDst, uint256 amount, uint8 errCode)
func (_Staking *StakingFilterer) WatchRedelegateFailed(opts *bind.WatchOpts, sink chan<- *StakingRedelegateFailed, delegator []common.Address, valSrc []common.Address, valDst []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var valSrcRule []interface{}
	for _, valSrcItem := range valSrc {
		valSrcRule = append(valSrcRule, valSrcItem)
	}
	var valDstRule []interface{}
	for _, valDstItem := range valDst {
		valDstRule = append(valDstRule, valDstItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "redelegateFailed", delegatorRule, valSrcRule, valDstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingRedelegateFailed)
				if err := _Staking.contract.UnpackLog(event, "redelegateFailed", log); err != nil {
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

// ParseRedelegateFailed is a log parse operation binding the contract event 0xb93bee5c59f85ede6b074a99f4ffcd3e3fc0d5c3d8156de331de89a49e0ce77c.
//
// Solidity: event redelegateFailed(address indexed delegator, address indexed valSrc, address indexed valDst, uint256 amount, uint8 errCode)
func (_Staking *StakingFilterer) ParseRedelegateFailed(log types.Log) (*StakingRedelegateFailed, error) {
	event := new(StakingRedelegateFailed)
	if err := _Staking.contract.UnpackLog(event, "redelegateFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingRedelegateSubmittedIterator is returned from FilterRedelegateSubmitted and is used to iterate over the raw logs and unpacked data for RedelegateSubmitted events raised by the Staking contract.
type StakingRedelegateSubmittedIterator struct {
	Event *StakingRedelegateSubmitted // Event containing the contract specifics and raw log

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
func (it *StakingRedelegateSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingRedelegateSubmitted)
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
		it.Event = new(StakingRedelegateSubmitted)
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
func (it *StakingRedelegateSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingRedelegateSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingRedelegateSubmitted represents a RedelegateSubmitted event raised by the Staking contract.
type StakingRedelegateSubmitted struct {
	Delegator    common.Address
	ValidatorSrc common.Address
	ValidatorDst common.Address
	Amount       *big.Int
	RelayerFee   *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRedelegateSubmitted is a free log retrieval operation binding the contract event 0xdb0d03fdfcb145c486c442659e6a341a8828985505097cb5190afcf541e84015.
//
// Solidity: event redelegateSubmitted(address indexed delegator, address indexed validatorSrc, address indexed validatorDst, uint256 amount, uint256 relayerFee)
func (_Staking *StakingFilterer) FilterRedelegateSubmitted(opts *bind.FilterOpts, delegator []common.Address, validatorSrc []common.Address, validatorDst []common.Address) (*StakingRedelegateSubmittedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorSrcRule []interface{}
	for _, validatorSrcItem := range validatorSrc {
		validatorSrcRule = append(validatorSrcRule, validatorSrcItem)
	}
	var validatorDstRule []interface{}
	for _, validatorDstItem := range validatorDst {
		validatorDstRule = append(validatorDstRule, validatorDstItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "redelegateSubmitted", delegatorRule, validatorSrcRule, validatorDstRule)
	if err != nil {
		return nil, err
	}
	return &StakingRedelegateSubmittedIterator{contract: _Staking.contract, event: "redelegateSubmitted", logs: logs, sub: sub}, nil
}

// WatchRedelegateSubmitted is a free log subscription operation binding the contract event 0xdb0d03fdfcb145c486c442659e6a341a8828985505097cb5190afcf541e84015.
//
// Solidity: event redelegateSubmitted(address indexed delegator, address indexed validatorSrc, address indexed validatorDst, uint256 amount, uint256 relayerFee)
func (_Staking *StakingFilterer) WatchRedelegateSubmitted(opts *bind.WatchOpts, sink chan<- *StakingRedelegateSubmitted, delegator []common.Address, validatorSrc []common.Address, validatorDst []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorSrcRule []interface{}
	for _, validatorSrcItem := range validatorSrc {
		validatorSrcRule = append(validatorSrcRule, validatorSrcItem)
	}
	var validatorDstRule []interface{}
	for _, validatorDstItem := range validatorDst {
		validatorDstRule = append(validatorDstRule, validatorDstItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "redelegateSubmitted", delegatorRule, validatorSrcRule, validatorDstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingRedelegateSubmitted)
				if err := _Staking.contract.UnpackLog(event, "redelegateSubmitted", log); err != nil {
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

// ParseRedelegateSubmitted is a log parse operation binding the contract event 0xdb0d03fdfcb145c486c442659e6a341a8828985505097cb5190afcf541e84015.
//
// Solidity: event redelegateSubmitted(address indexed delegator, address indexed validatorSrc, address indexed validatorDst, uint256 amount, uint256 relayerFee)
func (_Staking *StakingFilterer) ParseRedelegateSubmitted(log types.Log) (*StakingRedelegateSubmitted, error) {
	event := new(StakingRedelegateSubmitted)
	if err := _Staking.contract.UnpackLog(event, "redelegateSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingRedelegateSuccessIterator is returned from FilterRedelegateSuccess and is used to iterate over the raw logs and unpacked data for RedelegateSuccess events raised by the Staking contract.
type StakingRedelegateSuccessIterator struct {
	Event *StakingRedelegateSuccess // Event containing the contract specifics and raw log

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
func (it *StakingRedelegateSuccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingRedelegateSuccess)
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
		it.Event = new(StakingRedelegateSuccess)
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
func (it *StakingRedelegateSuccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingRedelegateSuccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingRedelegateSuccess represents a RedelegateSuccess event raised by the Staking contract.
type StakingRedelegateSuccess struct {
	Delegator common.Address
	ValSrc    common.Address
	ValDst    common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRedelegateSuccess is a free log retrieval operation binding the contract event 0x78bffae3f8c6691ac7fc1a3bff800cb2d612f5ad9ae5b0444cfe2eb15c189e18.
//
// Solidity: event redelegateSuccess(address indexed delegator, address indexed valSrc, address indexed valDst, uint256 amount)
func (_Staking *StakingFilterer) FilterRedelegateSuccess(opts *bind.FilterOpts, delegator []common.Address, valSrc []common.Address, valDst []common.Address) (*StakingRedelegateSuccessIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var valSrcRule []interface{}
	for _, valSrcItem := range valSrc {
		valSrcRule = append(valSrcRule, valSrcItem)
	}
	var valDstRule []interface{}
	for _, valDstItem := range valDst {
		valDstRule = append(valDstRule, valDstItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "redelegateSuccess", delegatorRule, valSrcRule, valDstRule)
	if err != nil {
		return nil, err
	}
	return &StakingRedelegateSuccessIterator{contract: _Staking.contract, event: "redelegateSuccess", logs: logs, sub: sub}, nil
}

// WatchRedelegateSuccess is a free log subscription operation binding the contract event 0x78bffae3f8c6691ac7fc1a3bff800cb2d612f5ad9ae5b0444cfe2eb15c189e18.
//
// Solidity: event redelegateSuccess(address indexed delegator, address indexed valSrc, address indexed valDst, uint256 amount)
func (_Staking *StakingFilterer) WatchRedelegateSuccess(opts *bind.WatchOpts, sink chan<- *StakingRedelegateSuccess, delegator []common.Address, valSrc []common.Address, valDst []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var valSrcRule []interface{}
	for _, valSrcItem := range valSrc {
		valSrcRule = append(valSrcRule, valSrcItem)
	}
	var valDstRule []interface{}
	for _, valDstItem := range valDst {
		valDstRule = append(valDstRule, valDstItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "redelegateSuccess", delegatorRule, valSrcRule, valDstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingRedelegateSuccess)
				if err := _Staking.contract.UnpackLog(event, "redelegateSuccess", log); err != nil {
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

// ParseRedelegateSuccess is a log parse operation binding the contract event 0x78bffae3f8c6691ac7fc1a3bff800cb2d612f5ad9ae5b0444cfe2eb15c189e18.
//
// Solidity: event redelegateSuccess(address indexed delegator, address indexed valSrc, address indexed valDst, uint256 amount)
func (_Staking *StakingFilterer) ParseRedelegateSuccess(log types.Log) (*StakingRedelegateSuccess, error) {
	event := new(StakingRedelegateSuccess)
	if err := _Staking.contract.UnpackLog(event, "redelegateSuccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingRewardClaimedIterator is returned from FilterRewardClaimed and is used to iterate over the raw logs and unpacked data for RewardClaimed events raised by the Staking contract.
type StakingRewardClaimedIterator struct {
	Event *StakingRewardClaimed // Event containing the contract specifics and raw log

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
func (it *StakingRewardClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingRewardClaimed)
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
		it.Event = new(StakingRewardClaimed)
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
func (it *StakingRewardClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingRewardClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingRewardClaimed represents a RewardClaimed event raised by the Staking contract.
type StakingRewardClaimed struct {
	Delegator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRewardClaimed is a free log retrieval operation binding the contract event 0x83b78188b13346b2ffb484da70d42ee27de7fbf9f2bd8045269e10ed643ccd76.
//
// Solidity: event rewardClaimed(address indexed delegator, uint256 amount)
func (_Staking *StakingFilterer) FilterRewardClaimed(opts *bind.FilterOpts, delegator []common.Address) (*StakingRewardClaimedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "rewardClaimed", delegatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingRewardClaimedIterator{contract: _Staking.contract, event: "rewardClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardClaimed is a free log subscription operation binding the contract event 0x83b78188b13346b2ffb484da70d42ee27de7fbf9f2bd8045269e10ed643ccd76.
//
// Solidity: event rewardClaimed(address indexed delegator, uint256 amount)
func (_Staking *StakingFilterer) WatchRewardClaimed(opts *bind.WatchOpts, sink chan<- *StakingRewardClaimed, delegator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "rewardClaimed", delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingRewardClaimed)
				if err := _Staking.contract.UnpackLog(event, "rewardClaimed", log); err != nil {
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

// ParseRewardClaimed is a log parse operation binding the contract event 0x83b78188b13346b2ffb484da70d42ee27de7fbf9f2bd8045269e10ed643ccd76.
//
// Solidity: event rewardClaimed(address indexed delegator, uint256 amount)
func (_Staking *StakingFilterer) ParseRewardClaimed(log types.Log) (*StakingRewardClaimed, error) {
	event := new(StakingRewardClaimed)
	if err := _Staking.contract.UnpackLog(event, "rewardClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingRewardReceivedIterator is returned from FilterRewardReceived and is used to iterate over the raw logs and unpacked data for RewardReceived events raised by the Staking contract.
type StakingRewardReceivedIterator struct {
	Event *StakingRewardReceived // Event containing the contract specifics and raw log

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
func (it *StakingRewardReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingRewardReceived)
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
		it.Event = new(StakingRewardReceived)
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
func (it *StakingRewardReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingRewardReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingRewardReceived represents a RewardReceived event raised by the Staking contract.
type StakingRewardReceived struct {
	Delegator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRewardReceived is a free log retrieval operation binding the contract event 0x7cc266c7b444f808013fa187f7b904d470a051a6564e78f482aa496581ba4bf8.
//
// Solidity: event rewardReceived(address indexed delegator, uint256 amount)
func (_Staking *StakingFilterer) FilterRewardReceived(opts *bind.FilterOpts, delegator []common.Address) (*StakingRewardReceivedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "rewardReceived", delegatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingRewardReceivedIterator{contract: _Staking.contract, event: "rewardReceived", logs: logs, sub: sub}, nil
}

// WatchRewardReceived is a free log subscription operation binding the contract event 0x7cc266c7b444f808013fa187f7b904d470a051a6564e78f482aa496581ba4bf8.
//
// Solidity: event rewardReceived(address indexed delegator, uint256 amount)
func (_Staking *StakingFilterer) WatchRewardReceived(opts *bind.WatchOpts, sink chan<- *StakingRewardReceived, delegator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "rewardReceived", delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingRewardReceived)
				if err := _Staking.contract.UnpackLog(event, "rewardReceived", log); err != nil {
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

// ParseRewardReceived is a log parse operation binding the contract event 0x7cc266c7b444f808013fa187f7b904d470a051a6564e78f482aa496581ba4bf8.
//
// Solidity: event rewardReceived(address indexed delegator, uint256 amount)
func (_Staking *StakingFilterer) ParseRewardReceived(log types.Log) (*StakingRewardReceived, error) {
	event := new(StakingRewardReceived)
	if err := _Staking.contract.UnpackLog(event, "rewardReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingUndelegateFailedIterator is returned from FilterUndelegateFailed and is used to iterate over the raw logs and unpacked data for UndelegateFailed events raised by the Staking contract.
type StakingUndelegateFailedIterator struct {
	Event *StakingUndelegateFailed // Event containing the contract specifics and raw log

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
func (it *StakingUndelegateFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUndelegateFailed)
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
		it.Event = new(StakingUndelegateFailed)
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
func (it *StakingUndelegateFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUndelegateFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUndelegateFailed represents a UndelegateFailed event raised by the Staking contract.
type StakingUndelegateFailed struct {
	Delegator common.Address
	Validator common.Address
	Amount    *big.Int
	ErrCode   uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUndelegateFailed is a free log retrieval operation binding the contract event 0x4417d10c1e33efa83a770b8d4f47176e78c08c1298d534901ad3b16bb585fa2e.
//
// Solidity: event undelegateFailed(address indexed delegator, address indexed validator, uint256 amount, uint8 errCode)
func (_Staking *StakingFilterer) FilterUndelegateFailed(opts *bind.FilterOpts, delegator []common.Address, validator []common.Address) (*StakingUndelegateFailedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "undelegateFailed", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingUndelegateFailedIterator{contract: _Staking.contract, event: "undelegateFailed", logs: logs, sub: sub}, nil
}

// WatchUndelegateFailed is a free log subscription operation binding the contract event 0x4417d10c1e33efa83a770b8d4f47176e78c08c1298d534901ad3b16bb585fa2e.
//
// Solidity: event undelegateFailed(address indexed delegator, address indexed validator, uint256 amount, uint8 errCode)
func (_Staking *StakingFilterer) WatchUndelegateFailed(opts *bind.WatchOpts, sink chan<- *StakingUndelegateFailed, delegator []common.Address, validator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "undelegateFailed", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUndelegateFailed)
				if err := _Staking.contract.UnpackLog(event, "undelegateFailed", log); err != nil {
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

// ParseUndelegateFailed is a log parse operation binding the contract event 0x4417d10c1e33efa83a770b8d4f47176e78c08c1298d534901ad3b16bb585fa2e.
//
// Solidity: event undelegateFailed(address indexed delegator, address indexed validator, uint256 amount, uint8 errCode)
func (_Staking *StakingFilterer) ParseUndelegateFailed(log types.Log) (*StakingUndelegateFailed, error) {
	event := new(StakingUndelegateFailed)
	if err := _Staking.contract.UnpackLog(event, "undelegateFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingUndelegateSubmittedIterator is returned from FilterUndelegateSubmitted and is used to iterate over the raw logs and unpacked data for UndelegateSubmitted events raised by the Staking contract.
type StakingUndelegateSubmittedIterator struct {
	Event *StakingUndelegateSubmitted // Event containing the contract specifics and raw log

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
func (it *StakingUndelegateSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUndelegateSubmitted)
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
		it.Event = new(StakingUndelegateSubmitted)
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
func (it *StakingUndelegateSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUndelegateSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUndelegateSubmitted represents a UndelegateSubmitted event raised by the Staking contract.
type StakingUndelegateSubmitted struct {
	Delegator  common.Address
	Validator  common.Address
	Amount     *big.Int
	RelayerFee *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUndelegateSubmitted is a free log retrieval operation binding the contract event 0xdf0b6ac27f3f3bb31cee3dab0f4fe40cc19c6a3f8daaec52e06b261e58a12519.
//
// Solidity: event undelegateSubmitted(address indexed delegator, address indexed validator, uint256 amount, uint256 relayerFee)
func (_Staking *StakingFilterer) FilterUndelegateSubmitted(opts *bind.FilterOpts, delegator []common.Address, validator []common.Address) (*StakingUndelegateSubmittedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "undelegateSubmitted", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingUndelegateSubmittedIterator{contract: _Staking.contract, event: "undelegateSubmitted", logs: logs, sub: sub}, nil
}

// WatchUndelegateSubmitted is a free log subscription operation binding the contract event 0xdf0b6ac27f3f3bb31cee3dab0f4fe40cc19c6a3f8daaec52e06b261e58a12519.
//
// Solidity: event undelegateSubmitted(address indexed delegator, address indexed validator, uint256 amount, uint256 relayerFee)
func (_Staking *StakingFilterer) WatchUndelegateSubmitted(opts *bind.WatchOpts, sink chan<- *StakingUndelegateSubmitted, delegator []common.Address, validator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "undelegateSubmitted", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUndelegateSubmitted)
				if err := _Staking.contract.UnpackLog(event, "undelegateSubmitted", log); err != nil {
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

// ParseUndelegateSubmitted is a log parse operation binding the contract event 0xdf0b6ac27f3f3bb31cee3dab0f4fe40cc19c6a3f8daaec52e06b261e58a12519.
//
// Solidity: event undelegateSubmitted(address indexed delegator, address indexed validator, uint256 amount, uint256 relayerFee)
func (_Staking *StakingFilterer) ParseUndelegateSubmitted(log types.Log) (*StakingUndelegateSubmitted, error) {
	event := new(StakingUndelegateSubmitted)
	if err := _Staking.contract.UnpackLog(event, "undelegateSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingUndelegateSuccessIterator is returned from FilterUndelegateSuccess and is used to iterate over the raw logs and unpacked data for UndelegateSuccess events raised by the Staking contract.
type StakingUndelegateSuccessIterator struct {
	Event *StakingUndelegateSuccess // Event containing the contract specifics and raw log

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
func (it *StakingUndelegateSuccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUndelegateSuccess)
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
		it.Event = new(StakingUndelegateSuccess)
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
func (it *StakingUndelegateSuccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUndelegateSuccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUndelegateSuccess represents a UndelegateSuccess event raised by the Staking contract.
type StakingUndelegateSuccess struct {
	Delegator common.Address
	Validator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUndelegateSuccess is a free log retrieval operation binding the contract event 0xd6f878a5bcbbe79a64e6418bb0d56aaa20b9a60587d45749819df88dfc7c3c44.
//
// Solidity: event undelegateSuccess(address indexed delegator, address indexed validator, uint256 amount)
func (_Staking *StakingFilterer) FilterUndelegateSuccess(opts *bind.FilterOpts, delegator []common.Address, validator []common.Address) (*StakingUndelegateSuccessIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "undelegateSuccess", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingUndelegateSuccessIterator{contract: _Staking.contract, event: "undelegateSuccess", logs: logs, sub: sub}, nil
}

// WatchUndelegateSuccess is a free log subscription operation binding the contract event 0xd6f878a5bcbbe79a64e6418bb0d56aaa20b9a60587d45749819df88dfc7c3c44.
//
// Solidity: event undelegateSuccess(address indexed delegator, address indexed validator, uint256 amount)
func (_Staking *StakingFilterer) WatchUndelegateSuccess(opts *bind.WatchOpts, sink chan<- *StakingUndelegateSuccess, delegator []common.Address, validator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "undelegateSuccess", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUndelegateSuccess)
				if err := _Staking.contract.UnpackLog(event, "undelegateSuccess", log); err != nil {
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

// ParseUndelegateSuccess is a log parse operation binding the contract event 0xd6f878a5bcbbe79a64e6418bb0d56aaa20b9a60587d45749819df88dfc7c3c44.
//
// Solidity: event undelegateSuccess(address indexed delegator, address indexed validator, uint256 amount)
func (_Staking *StakingFilterer) ParseUndelegateSuccess(log types.Log) (*StakingUndelegateSuccess, error) {
	event := new(StakingUndelegateSuccess)
	if err := _Staking.contract.UnpackLog(event, "undelegateSuccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingUndelegatedClaimedIterator is returned from FilterUndelegatedClaimed and is used to iterate over the raw logs and unpacked data for UndelegatedClaimed events raised by the Staking contract.
type StakingUndelegatedClaimedIterator struct {
	Event *StakingUndelegatedClaimed // Event containing the contract specifics and raw log

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
func (it *StakingUndelegatedClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUndelegatedClaimed)
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
		it.Event = new(StakingUndelegatedClaimed)
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
func (it *StakingUndelegatedClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUndelegatedClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUndelegatedClaimed represents a UndelegatedClaimed event raised by the Staking contract.
type StakingUndelegatedClaimed struct {
	Delegator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUndelegatedClaimed is a free log retrieval operation binding the contract event 0xc712d133b8d448221aaed2198ed1f0db6dfc860fb01bc3a630916fe6cbef946f.
//
// Solidity: event undelegatedClaimed(address indexed delegator, uint256 amount)
func (_Staking *StakingFilterer) FilterUndelegatedClaimed(opts *bind.FilterOpts, delegator []common.Address) (*StakingUndelegatedClaimedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "undelegatedClaimed", delegatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingUndelegatedClaimedIterator{contract: _Staking.contract, event: "undelegatedClaimed", logs: logs, sub: sub}, nil
}

// WatchUndelegatedClaimed is a free log subscription operation binding the contract event 0xc712d133b8d448221aaed2198ed1f0db6dfc860fb01bc3a630916fe6cbef946f.
//
// Solidity: event undelegatedClaimed(address indexed delegator, uint256 amount)
func (_Staking *StakingFilterer) WatchUndelegatedClaimed(opts *bind.WatchOpts, sink chan<- *StakingUndelegatedClaimed, delegator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "undelegatedClaimed", delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUndelegatedClaimed)
				if err := _Staking.contract.UnpackLog(event, "undelegatedClaimed", log); err != nil {
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

// ParseUndelegatedClaimed is a log parse operation binding the contract event 0xc712d133b8d448221aaed2198ed1f0db6dfc860fb01bc3a630916fe6cbef946f.
//
// Solidity: event undelegatedClaimed(address indexed delegator, uint256 amount)
func (_Staking *StakingFilterer) ParseUndelegatedClaimed(log types.Log) (*StakingUndelegatedClaimed, error) {
	event := new(StakingUndelegatedClaimed)
	if err := _Staking.contract.UnpackLog(event, "undelegatedClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingUndelegatedReceivedIterator is returned from FilterUndelegatedReceived and is used to iterate over the raw logs and unpacked data for UndelegatedReceived events raised by the Staking contract.
type StakingUndelegatedReceivedIterator struct {
	Event *StakingUndelegatedReceived // Event containing the contract specifics and raw log

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
func (it *StakingUndelegatedReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUndelegatedReceived)
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
		it.Event = new(StakingUndelegatedReceived)
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
func (it *StakingUndelegatedReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUndelegatedReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUndelegatedReceived represents a UndelegatedReceived event raised by the Staking contract.
type StakingUndelegatedReceived struct {
	Delegator common.Address
	Validator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUndelegatedReceived is a free log retrieval operation binding the contract event 0x35a799836f74fac7eccf5c73902823b970543d2274d3b93d8da3d37a255772a2.
//
// Solidity: event undelegatedReceived(address indexed delegator, address indexed validator, uint256 amount)
func (_Staking *StakingFilterer) FilterUndelegatedReceived(opts *bind.FilterOpts, delegator []common.Address, validator []common.Address) (*StakingUndelegatedReceivedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "undelegatedReceived", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingUndelegatedReceivedIterator{contract: _Staking.contract, event: "undelegatedReceived", logs: logs, sub: sub}, nil
}

// WatchUndelegatedReceived is a free log subscription operation binding the contract event 0x35a799836f74fac7eccf5c73902823b970543d2274d3b93d8da3d37a255772a2.
//
// Solidity: event undelegatedReceived(address indexed delegator, address indexed validator, uint256 amount)
func (_Staking *StakingFilterer) WatchUndelegatedReceived(opts *bind.WatchOpts, sink chan<- *StakingUndelegatedReceived, delegator []common.Address, validator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "undelegatedReceived", delegatorRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUndelegatedReceived)
				if err := _Staking.contract.UnpackLog(event, "undelegatedReceived", log); err != nil {
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

// ParseUndelegatedReceived is a log parse operation binding the contract event 0x35a799836f74fac7eccf5c73902823b970543d2274d3b93d8da3d37a255772a2.
//
// Solidity: event undelegatedReceived(address indexed delegator, address indexed validator, uint256 amount)
func (_Staking *StakingFilterer) ParseUndelegatedReceived(log types.Log) (*StakingUndelegatedReceived, error) {
	event := new(StakingUndelegatedReceived)
	if err := _Staking.contract.UnpackLog(event, "undelegatedReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
