// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TokenHub

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

// TokenHubABI is the input ABI used to generate the binding from.
const TokenHubABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"paramChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"receiveDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bep2eAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"refundAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"status\",\"type\":\"uint32\"}],\"name\":\"refundFailure\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bep2eAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"refundAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"status\",\"type\":\"uint32\"}],\"name\":\"refundSuccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"rewardTo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bep2eAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"refundAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferInSuccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bep2eAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"senderAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"relayFee\",\"type\":\"uint256\"}],\"name\":\"transferOutSuccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"channelId\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"msgBytes\",\"type\":\"bytes\"}],\"name\":\"unexpectedPackage\",\"type\":\"event\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[],\"name\":\"BEP2_TOKEN_DECIMALS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"BEP2_TOKEN_SYMBOL_FOR_BNB\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"BIND_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CODE_OK\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CROSS_CHAIN_CONTRACT_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ERROR_FAIL_DECODE\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"GOV_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"GOV_HUB_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"INCENTIVIZE_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"INIT_MINIMUM_RELAY_FEE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"LIGHT_CLIENT_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAXIMUM_BEP2E_SYMBOL_LEN\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_BEP2_TOTAL_SUPPLY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_GAS_FOR_CALLING_BEP2E\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_GAS_FOR_TRANSFER_BNB\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MINIMUM_BEP2E_SYMBOL_LEN\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"RELAYERHUB_CONTRACT_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"REWARD_UPPER_LIMIT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"SLASH_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"SLASH_CONTRACT_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"STAKING_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"SYSTEM_REWARD_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TEN_DECIMALS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TOKEN_HUB_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TOKEN_MANAGER_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TRANSFER_IN_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TRANSFER_IN_FAILURE_INSUFFICIENT_BALANCE\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TRANSFER_IN_FAILURE_NON_PAYABLE_RECIPIENT\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TRANSFER_IN_FAILURE_TIMEOUT\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TRANSFER_IN_FAILURE_UNBOUND_TOKEN\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TRANSFER_IN_FAILURE_UNKNOWN\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TRANSFER_IN_SUCCESS\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TRANSFER_OUT_CHANNELID\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VALIDATOR_CONTRACT_ADDR\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"alreadyInit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"bep2eContractDecimals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bscChainID\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"relayFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"init\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"claimRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMiniRelayFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"channelId\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"msgBytes\",\"type\":\"bytes\"}],\"name\":\"handleSynPackage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"channelId\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"msgBytes\",\"type\":\"bytes\"}],\"name\":\"handleAckPackage\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"channelId\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"msgBytes\",\"type\":\"bytes\"}],\"name\":\"handleFailAckPackage\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"expireTime\",\"type\":\"uint64\"}],\"name\":\"transferOut\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipientAddrs\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address[]\",\"name\":\"refundAddrs\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"expireTime\",\"type\":\"uint64\"}],\"name\":\"batchTransferOutBNB\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"updateParam\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"bep2Symbol\",\"type\":\"bytes32\"}],\"name\":\"getContractAddrByBEP2Symbol\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddr\",\"type\":\"address\"}],\"name\":\"getBep2SymbolByContractAddr\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"bep2Symbol\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"contractAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"}],\"name\":\"bindToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"bep2Symbol\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"contractAddr\",\"type\":\"address\"}],\"name\":\"unbindToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"string\",\"name\":\"bep2Symbol\",\"type\":\"string\"}],\"name\":\"getBoundContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddr\",\"type\":\"address\"}],\"name\":\"getBoundBep2Symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// TokenHub is an auto generated Go binding around an Ethereum contract.
type TokenHub struct {
	TokenHubCaller     // Read-only binding to the contract
	TokenHubTransactor // Write-only binding to the contract
	TokenHubFilterer   // Log filterer for contract events
}

// TokenHubCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenHubCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenHubTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenHubTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenHubFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenHubFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenHubSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenHubSession struct {
	Contract     *TokenHub         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenHubCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenHubCallerSession struct {
	Contract *TokenHubCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TokenHubTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenHubTransactorSession struct {
	Contract     *TokenHubTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TokenHubRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenHubRaw struct {
	Contract *TokenHub // Generic contract binding to access the raw methods on
}

// TokenHubCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenHubCallerRaw struct {
	Contract *TokenHubCaller // Generic read-only contract binding to access the raw methods on
}

// TokenHubTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenHubTransactorRaw struct {
	Contract *TokenHubTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenHub creates a new instance of TokenHub, bound to a specific deployed contract.
func NewTokenHub(address common.Address, backend bind.ContractBackend) (*TokenHub, error) {
	contract, err := bindTokenHub(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenHub{TokenHubCaller: TokenHubCaller{contract: contract}, TokenHubTransactor: TokenHubTransactor{contract: contract}, TokenHubFilterer: TokenHubFilterer{contract: contract}}, nil
}

// NewTokenHubCaller creates a new read-only instance of TokenHub, bound to a specific deployed contract.
func NewTokenHubCaller(address common.Address, caller bind.ContractCaller) (*TokenHubCaller, error) {
	contract, err := bindTokenHub(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenHubCaller{contract: contract}, nil
}

// NewTokenHubTransactor creates a new write-only instance of TokenHub, bound to a specific deployed contract.
func NewTokenHubTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenHubTransactor, error) {
	contract, err := bindTokenHub(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenHubTransactor{contract: contract}, nil
}

// NewTokenHubFilterer creates a new log filterer instance of TokenHub, bound to a specific deployed contract.
func NewTokenHubFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenHubFilterer, error) {
	contract, err := bindTokenHub(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenHubFilterer{contract: contract}, nil
}

// bindTokenHub binds a generic wrapper to an already deployed contract.
func bindTokenHub(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenHubABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenHub *TokenHubRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenHub.Contract.TokenHubCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenHub *TokenHubRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenHub.Contract.TokenHubTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenHub *TokenHubRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenHub.Contract.TokenHubTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenHub *TokenHubCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenHub.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenHub *TokenHubTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenHub.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenHub *TokenHubTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenHub.Contract.contract.Transact(opts, method, params...)
}

// BEP2TOKENDECIMALS is a free data retrieval call binding the contract method 0x61368475.
//
// Solidity: function BEP2_TOKEN_DECIMALS() view returns(uint8)
func (_TokenHub *TokenHubCaller) BEP2TOKENDECIMALS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "BEP2_TOKEN_DECIMALS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// BEP2TOKENDECIMALS is a free data retrieval call binding the contract method 0x61368475.
//
// Solidity: function BEP2_TOKEN_DECIMALS() view returns(uint8)
func (_TokenHub *TokenHubSession) BEP2TOKENDECIMALS() (uint8, error) {
	return _TokenHub.Contract.BEP2TOKENDECIMALS(&_TokenHub.CallOpts)
}

// BEP2TOKENDECIMALS is a free data retrieval call binding the contract method 0x61368475.
//
// Solidity: function BEP2_TOKEN_DECIMALS() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) BEP2TOKENDECIMALS() (uint8, error) {
	return _TokenHub.Contract.BEP2TOKENDECIMALS(&_TokenHub.CallOpts)
}

// BEP2TOKENSYMBOLFORBNB is a free data retrieval call binding the contract method 0xb9fd21e3.
//
// Solidity: function BEP2_TOKEN_SYMBOL_FOR_BNB() view returns(bytes32)
func (_TokenHub *TokenHubCaller) BEP2TOKENSYMBOLFORBNB(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "BEP2_TOKEN_SYMBOL_FOR_BNB")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BEP2TOKENSYMBOLFORBNB is a free data retrieval call binding the contract method 0xb9fd21e3.
//
// Solidity: function BEP2_TOKEN_SYMBOL_FOR_BNB() view returns(bytes32)
func (_TokenHub *TokenHubSession) BEP2TOKENSYMBOLFORBNB() ([32]byte, error) {
	return _TokenHub.Contract.BEP2TOKENSYMBOLFORBNB(&_TokenHub.CallOpts)
}

// BEP2TOKENSYMBOLFORBNB is a free data retrieval call binding the contract method 0xb9fd21e3.
//
// Solidity: function BEP2_TOKEN_SYMBOL_FOR_BNB() view returns(bytes32)
func (_TokenHub *TokenHubCallerSession) BEP2TOKENSYMBOLFORBNB() ([32]byte, error) {
	return _TokenHub.Contract.BEP2TOKENSYMBOLFORBNB(&_TokenHub.CallOpts)
}

// BINDCHANNELID is a free data retrieval call binding the contract method 0x3dffc387.
//
// Solidity: function BIND_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubCaller) BINDCHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "BIND_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// BINDCHANNELID is a free data retrieval call binding the contract method 0x3dffc387.
//
// Solidity: function BIND_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubSession) BINDCHANNELID() (uint8, error) {
	return _TokenHub.Contract.BINDCHANNELID(&_TokenHub.CallOpts)
}

// BINDCHANNELID is a free data retrieval call binding the contract method 0x3dffc387.
//
// Solidity: function BIND_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) BINDCHANNELID() (uint8, error) {
	return _TokenHub.Contract.BINDCHANNELID(&_TokenHub.CallOpts)
}

// CODEOK is a free data retrieval call binding the contract method 0xab51bb96.
//
// Solidity: function CODE_OK() view returns(uint32)
func (_TokenHub *TokenHubCaller) CODEOK(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "CODE_OK")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// CODEOK is a free data retrieval call binding the contract method 0xab51bb96.
//
// Solidity: function CODE_OK() view returns(uint32)
func (_TokenHub *TokenHubSession) CODEOK() (uint32, error) {
	return _TokenHub.Contract.CODEOK(&_TokenHub.CallOpts)
}

// CODEOK is a free data retrieval call binding the contract method 0xab51bb96.
//
// Solidity: function CODE_OK() view returns(uint32)
func (_TokenHub *TokenHubCallerSession) CODEOK() (uint32, error) {
	return _TokenHub.Contract.CODEOK(&_TokenHub.CallOpts)
}

// CROSSCHAINCONTRACTADDR is a free data retrieval call binding the contract method 0x51e80672.
//
// Solidity: function CROSS_CHAIN_CONTRACT_ADDR() view returns(address)
func (_TokenHub *TokenHubCaller) CROSSCHAINCONTRACTADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "CROSS_CHAIN_CONTRACT_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CROSSCHAINCONTRACTADDR is a free data retrieval call binding the contract method 0x51e80672.
//
// Solidity: function CROSS_CHAIN_CONTRACT_ADDR() view returns(address)
func (_TokenHub *TokenHubSession) CROSSCHAINCONTRACTADDR() (common.Address, error) {
	return _TokenHub.Contract.CROSSCHAINCONTRACTADDR(&_TokenHub.CallOpts)
}

// CROSSCHAINCONTRACTADDR is a free data retrieval call binding the contract method 0x51e80672.
//
// Solidity: function CROSS_CHAIN_CONTRACT_ADDR() view returns(address)
func (_TokenHub *TokenHubCallerSession) CROSSCHAINCONTRACTADDR() (common.Address, error) {
	return _TokenHub.Contract.CROSSCHAINCONTRACTADDR(&_TokenHub.CallOpts)
}

// ERRORFAILDECODE is a free data retrieval call binding the contract method 0x0bee7a67.
//
// Solidity: function ERROR_FAIL_DECODE() view returns(uint32)
func (_TokenHub *TokenHubCaller) ERRORFAILDECODE(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "ERROR_FAIL_DECODE")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// ERRORFAILDECODE is a free data retrieval call binding the contract method 0x0bee7a67.
//
// Solidity: function ERROR_FAIL_DECODE() view returns(uint32)
func (_TokenHub *TokenHubSession) ERRORFAILDECODE() (uint32, error) {
	return _TokenHub.Contract.ERRORFAILDECODE(&_TokenHub.CallOpts)
}

// ERRORFAILDECODE is a free data retrieval call binding the contract method 0x0bee7a67.
//
// Solidity: function ERROR_FAIL_DECODE() view returns(uint32)
func (_TokenHub *TokenHubCallerSession) ERRORFAILDECODE() (uint32, error) {
	return _TokenHub.Contract.ERRORFAILDECODE(&_TokenHub.CallOpts)
}

// GOVCHANNELID is a free data retrieval call binding the contract method 0x96713da9.
//
// Solidity: function GOV_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubCaller) GOVCHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "GOV_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GOVCHANNELID is a free data retrieval call binding the contract method 0x96713da9.
//
// Solidity: function GOV_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubSession) GOVCHANNELID() (uint8, error) {
	return _TokenHub.Contract.GOVCHANNELID(&_TokenHub.CallOpts)
}

// GOVCHANNELID is a free data retrieval call binding the contract method 0x96713da9.
//
// Solidity: function GOV_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) GOVCHANNELID() (uint8, error) {
	return _TokenHub.Contract.GOVCHANNELID(&_TokenHub.CallOpts)
}

// GOVHUBADDR is a free data retrieval call binding the contract method 0x9dc09262.
//
// Solidity: function GOV_HUB_ADDR() view returns(address)
func (_TokenHub *TokenHubCaller) GOVHUBADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "GOV_HUB_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GOVHUBADDR is a free data retrieval call binding the contract method 0x9dc09262.
//
// Solidity: function GOV_HUB_ADDR() view returns(address)
func (_TokenHub *TokenHubSession) GOVHUBADDR() (common.Address, error) {
	return _TokenHub.Contract.GOVHUBADDR(&_TokenHub.CallOpts)
}

// GOVHUBADDR is a free data retrieval call binding the contract method 0x9dc09262.
//
// Solidity: function GOV_HUB_ADDR() view returns(address)
func (_TokenHub *TokenHubCallerSession) GOVHUBADDR() (common.Address, error) {
	return _TokenHub.Contract.GOVHUBADDR(&_TokenHub.CallOpts)
}

// INCENTIVIZEADDR is a free data retrieval call binding the contract method 0x6e47b482.
//
// Solidity: function INCENTIVIZE_ADDR() view returns(address)
func (_TokenHub *TokenHubCaller) INCENTIVIZEADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "INCENTIVIZE_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// INCENTIVIZEADDR is a free data retrieval call binding the contract method 0x6e47b482.
//
// Solidity: function INCENTIVIZE_ADDR() view returns(address)
func (_TokenHub *TokenHubSession) INCENTIVIZEADDR() (common.Address, error) {
	return _TokenHub.Contract.INCENTIVIZEADDR(&_TokenHub.CallOpts)
}

// INCENTIVIZEADDR is a free data retrieval call binding the contract method 0x6e47b482.
//
// Solidity: function INCENTIVIZE_ADDR() view returns(address)
func (_TokenHub *TokenHubCallerSession) INCENTIVIZEADDR() (common.Address, error) {
	return _TokenHub.Contract.INCENTIVIZEADDR(&_TokenHub.CallOpts)
}

// INITMINIMUMRELAYFEE is a free data retrieval call binding the contract method 0x50432d32.
//
// Solidity: function INIT_MINIMUM_RELAY_FEE() view returns(uint256)
func (_TokenHub *TokenHubCaller) INITMINIMUMRELAYFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "INIT_MINIMUM_RELAY_FEE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INITMINIMUMRELAYFEE is a free data retrieval call binding the contract method 0x50432d32.
//
// Solidity: function INIT_MINIMUM_RELAY_FEE() view returns(uint256)
func (_TokenHub *TokenHubSession) INITMINIMUMRELAYFEE() (*big.Int, error) {
	return _TokenHub.Contract.INITMINIMUMRELAYFEE(&_TokenHub.CallOpts)
}

// INITMINIMUMRELAYFEE is a free data retrieval call binding the contract method 0x50432d32.
//
// Solidity: function INIT_MINIMUM_RELAY_FEE() view returns(uint256)
func (_TokenHub *TokenHubCallerSession) INITMINIMUMRELAYFEE() (*big.Int, error) {
	return _TokenHub.Contract.INITMINIMUMRELAYFEE(&_TokenHub.CallOpts)
}

// LIGHTCLIENTADDR is a free data retrieval call binding the contract method 0xdc927faf.
//
// Solidity: function LIGHT_CLIENT_ADDR() view returns(address)
func (_TokenHub *TokenHubCaller) LIGHTCLIENTADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "LIGHT_CLIENT_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LIGHTCLIENTADDR is a free data retrieval call binding the contract method 0xdc927faf.
//
// Solidity: function LIGHT_CLIENT_ADDR() view returns(address)
func (_TokenHub *TokenHubSession) LIGHTCLIENTADDR() (common.Address, error) {
	return _TokenHub.Contract.LIGHTCLIENTADDR(&_TokenHub.CallOpts)
}

// LIGHTCLIENTADDR is a free data retrieval call binding the contract method 0xdc927faf.
//
// Solidity: function LIGHT_CLIENT_ADDR() view returns(address)
func (_TokenHub *TokenHubCallerSession) LIGHTCLIENTADDR() (common.Address, error) {
	return _TokenHub.Contract.LIGHTCLIENTADDR(&_TokenHub.CallOpts)
}

// MAXIMUMBEP2ESYMBOLLEN is a free data retrieval call binding the contract method 0x077b8f35.
//
// Solidity: function MAXIMUM_BEP2E_SYMBOL_LEN() view returns(uint8)
func (_TokenHub *TokenHubCaller) MAXIMUMBEP2ESYMBOLLEN(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "MAXIMUM_BEP2E_SYMBOL_LEN")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// MAXIMUMBEP2ESYMBOLLEN is a free data retrieval call binding the contract method 0x077b8f35.
//
// Solidity: function MAXIMUM_BEP2E_SYMBOL_LEN() view returns(uint8)
func (_TokenHub *TokenHubSession) MAXIMUMBEP2ESYMBOLLEN() (uint8, error) {
	return _TokenHub.Contract.MAXIMUMBEP2ESYMBOLLEN(&_TokenHub.CallOpts)
}

// MAXIMUMBEP2ESYMBOLLEN is a free data retrieval call binding the contract method 0x077b8f35.
//
// Solidity: function MAXIMUM_BEP2E_SYMBOL_LEN() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) MAXIMUMBEP2ESYMBOLLEN() (uint8, error) {
	return _TokenHub.Contract.MAXIMUMBEP2ESYMBOLLEN(&_TokenHub.CallOpts)
}

// MAXBEP2TOTALSUPPLY is a free data retrieval call binding the contract method 0x9a854bbd.
//
// Solidity: function MAX_BEP2_TOTAL_SUPPLY() view returns(uint256)
func (_TokenHub *TokenHubCaller) MAXBEP2TOTALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "MAX_BEP2_TOTAL_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXBEP2TOTALSUPPLY is a free data retrieval call binding the contract method 0x9a854bbd.
//
// Solidity: function MAX_BEP2_TOTAL_SUPPLY() view returns(uint256)
func (_TokenHub *TokenHubSession) MAXBEP2TOTALSUPPLY() (*big.Int, error) {
	return _TokenHub.Contract.MAXBEP2TOTALSUPPLY(&_TokenHub.CallOpts)
}

// MAXBEP2TOTALSUPPLY is a free data retrieval call binding the contract method 0x9a854bbd.
//
// Solidity: function MAX_BEP2_TOTAL_SUPPLY() view returns(uint256)
func (_TokenHub *TokenHubCallerSession) MAXBEP2TOTALSUPPLY() (*big.Int, error) {
	return _TokenHub.Contract.MAXBEP2TOTALSUPPLY(&_TokenHub.CallOpts)
}

// MAXGASFORCALLINGBEP2E is a free data retrieval call binding the contract method 0xb7701861.
//
// Solidity: function MAX_GAS_FOR_CALLING_BEP2E() view returns(uint256)
func (_TokenHub *TokenHubCaller) MAXGASFORCALLINGBEP2E(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "MAX_GAS_FOR_CALLING_BEP2E")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXGASFORCALLINGBEP2E is a free data retrieval call binding the contract method 0xb7701861.
//
// Solidity: function MAX_GAS_FOR_CALLING_BEP2E() view returns(uint256)
func (_TokenHub *TokenHubSession) MAXGASFORCALLINGBEP2E() (*big.Int, error) {
	return _TokenHub.Contract.MAXGASFORCALLINGBEP2E(&_TokenHub.CallOpts)
}

// MAXGASFORCALLINGBEP2E is a free data retrieval call binding the contract method 0xb7701861.
//
// Solidity: function MAX_GAS_FOR_CALLING_BEP2E() view returns(uint256)
func (_TokenHub *TokenHubCallerSession) MAXGASFORCALLINGBEP2E() (*big.Int, error) {
	return _TokenHub.Contract.MAXGASFORCALLINGBEP2E(&_TokenHub.CallOpts)
}

// MAXGASFORTRANSFERBNB is a free data retrieval call binding the contract method 0xfa9e9159.
//
// Solidity: function MAX_GAS_FOR_TRANSFER_BNB() view returns(uint256)
func (_TokenHub *TokenHubCaller) MAXGASFORTRANSFERBNB(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "MAX_GAS_FOR_TRANSFER_BNB")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXGASFORTRANSFERBNB is a free data retrieval call binding the contract method 0xfa9e9159.
//
// Solidity: function MAX_GAS_FOR_TRANSFER_BNB() view returns(uint256)
func (_TokenHub *TokenHubSession) MAXGASFORTRANSFERBNB() (*big.Int, error) {
	return _TokenHub.Contract.MAXGASFORTRANSFERBNB(&_TokenHub.CallOpts)
}

// MAXGASFORTRANSFERBNB is a free data retrieval call binding the contract method 0xfa9e9159.
//
// Solidity: function MAX_GAS_FOR_TRANSFER_BNB() view returns(uint256)
func (_TokenHub *TokenHubCallerSession) MAXGASFORTRANSFERBNB() (*big.Int, error) {
	return _TokenHub.Contract.MAXGASFORTRANSFERBNB(&_TokenHub.CallOpts)
}

// MINIMUMBEP2ESYMBOLLEN is a free data retrieval call binding the contract method 0xdc6f5e90.
//
// Solidity: function MINIMUM_BEP2E_SYMBOL_LEN() view returns(uint8)
func (_TokenHub *TokenHubCaller) MINIMUMBEP2ESYMBOLLEN(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "MINIMUM_BEP2E_SYMBOL_LEN")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// MINIMUMBEP2ESYMBOLLEN is a free data retrieval call binding the contract method 0xdc6f5e90.
//
// Solidity: function MINIMUM_BEP2E_SYMBOL_LEN() view returns(uint8)
func (_TokenHub *TokenHubSession) MINIMUMBEP2ESYMBOLLEN() (uint8, error) {
	return _TokenHub.Contract.MINIMUMBEP2ESYMBOLLEN(&_TokenHub.CallOpts)
}

// MINIMUMBEP2ESYMBOLLEN is a free data retrieval call binding the contract method 0xdc6f5e90.
//
// Solidity: function MINIMUM_BEP2E_SYMBOL_LEN() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) MINIMUMBEP2ESYMBOLLEN() (uint8, error) {
	return _TokenHub.Contract.MINIMUMBEP2ESYMBOLLEN(&_TokenHub.CallOpts)
}

// RELAYERHUBCONTRACTADDR is a free data retrieval call binding the contract method 0xa1a11bf5.
//
// Solidity: function RELAYERHUB_CONTRACT_ADDR() view returns(address)
func (_TokenHub *TokenHubCaller) RELAYERHUBCONTRACTADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "RELAYERHUB_CONTRACT_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RELAYERHUBCONTRACTADDR is a free data retrieval call binding the contract method 0xa1a11bf5.
//
// Solidity: function RELAYERHUB_CONTRACT_ADDR() view returns(address)
func (_TokenHub *TokenHubSession) RELAYERHUBCONTRACTADDR() (common.Address, error) {
	return _TokenHub.Contract.RELAYERHUBCONTRACTADDR(&_TokenHub.CallOpts)
}

// RELAYERHUBCONTRACTADDR is a free data retrieval call binding the contract method 0xa1a11bf5.
//
// Solidity: function RELAYERHUB_CONTRACT_ADDR() view returns(address)
func (_TokenHub *TokenHubCallerSession) RELAYERHUBCONTRACTADDR() (common.Address, error) {
	return _TokenHub.Contract.RELAYERHUBCONTRACTADDR(&_TokenHub.CallOpts)
}

// REWARDUPPERLIMIT is a free data retrieval call binding the contract method 0x43a368b9.
//
// Solidity: function REWARD_UPPER_LIMIT() view returns(uint256)
func (_TokenHub *TokenHubCaller) REWARDUPPERLIMIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "REWARD_UPPER_LIMIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// REWARDUPPERLIMIT is a free data retrieval call binding the contract method 0x43a368b9.
//
// Solidity: function REWARD_UPPER_LIMIT() view returns(uint256)
func (_TokenHub *TokenHubSession) REWARDUPPERLIMIT() (*big.Int, error) {
	return _TokenHub.Contract.REWARDUPPERLIMIT(&_TokenHub.CallOpts)
}

// REWARDUPPERLIMIT is a free data retrieval call binding the contract method 0x43a368b9.
//
// Solidity: function REWARD_UPPER_LIMIT() view returns(uint256)
func (_TokenHub *TokenHubCallerSession) REWARDUPPERLIMIT() (*big.Int, error) {
	return _TokenHub.Contract.REWARDUPPERLIMIT(&_TokenHub.CallOpts)
}

// SLASHCHANNELID is a free data retrieval call binding the contract method 0x7942fd05.
//
// Solidity: function SLASH_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubCaller) SLASHCHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "SLASH_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// SLASHCHANNELID is a free data retrieval call binding the contract method 0x7942fd05.
//
// Solidity: function SLASH_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubSession) SLASHCHANNELID() (uint8, error) {
	return _TokenHub.Contract.SLASHCHANNELID(&_TokenHub.CallOpts)
}

// SLASHCHANNELID is a free data retrieval call binding the contract method 0x7942fd05.
//
// Solidity: function SLASH_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) SLASHCHANNELID() (uint8, error) {
	return _TokenHub.Contract.SLASHCHANNELID(&_TokenHub.CallOpts)
}

// SLASHCONTRACTADDR is a free data retrieval call binding the contract method 0x43756e5c.
//
// Solidity: function SLASH_CONTRACT_ADDR() view returns(address)
func (_TokenHub *TokenHubCaller) SLASHCONTRACTADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "SLASH_CONTRACT_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SLASHCONTRACTADDR is a free data retrieval call binding the contract method 0x43756e5c.
//
// Solidity: function SLASH_CONTRACT_ADDR() view returns(address)
func (_TokenHub *TokenHubSession) SLASHCONTRACTADDR() (common.Address, error) {
	return _TokenHub.Contract.SLASHCONTRACTADDR(&_TokenHub.CallOpts)
}

// SLASHCONTRACTADDR is a free data retrieval call binding the contract method 0x43756e5c.
//
// Solidity: function SLASH_CONTRACT_ADDR() view returns(address)
func (_TokenHub *TokenHubCallerSession) SLASHCONTRACTADDR() (common.Address, error) {
	return _TokenHub.Contract.SLASHCONTRACTADDR(&_TokenHub.CallOpts)
}

// STAKINGCHANNELID is a free data retrieval call binding the contract method 0x4bf6c882.
//
// Solidity: function STAKING_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubCaller) STAKINGCHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "STAKING_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// STAKINGCHANNELID is a free data retrieval call binding the contract method 0x4bf6c882.
//
// Solidity: function STAKING_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubSession) STAKINGCHANNELID() (uint8, error) {
	return _TokenHub.Contract.STAKINGCHANNELID(&_TokenHub.CallOpts)
}

// STAKINGCHANNELID is a free data retrieval call binding the contract method 0x4bf6c882.
//
// Solidity: function STAKING_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) STAKINGCHANNELID() (uint8, error) {
	return _TokenHub.Contract.STAKINGCHANNELID(&_TokenHub.CallOpts)
}

// SYSTEMREWARDADDR is a free data retrieval call binding the contract method 0xc81b1662.
//
// Solidity: function SYSTEM_REWARD_ADDR() view returns(address)
func (_TokenHub *TokenHubCaller) SYSTEMREWARDADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "SYSTEM_REWARD_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SYSTEMREWARDADDR is a free data retrieval call binding the contract method 0xc81b1662.
//
// Solidity: function SYSTEM_REWARD_ADDR() view returns(address)
func (_TokenHub *TokenHubSession) SYSTEMREWARDADDR() (common.Address, error) {
	return _TokenHub.Contract.SYSTEMREWARDADDR(&_TokenHub.CallOpts)
}

// SYSTEMREWARDADDR is a free data retrieval call binding the contract method 0xc81b1662.
//
// Solidity: function SYSTEM_REWARD_ADDR() view returns(address)
func (_TokenHub *TokenHubCallerSession) SYSTEMREWARDADDR() (common.Address, error) {
	return _TokenHub.Contract.SYSTEMREWARDADDR(&_TokenHub.CallOpts)
}

// TENDECIMALS is a free data retrieval call binding the contract method 0x5d499b1b.
//
// Solidity: function TEN_DECIMALS() view returns(uint256)
func (_TokenHub *TokenHubCaller) TENDECIMALS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "TEN_DECIMALS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TENDECIMALS is a free data retrieval call binding the contract method 0x5d499b1b.
//
// Solidity: function TEN_DECIMALS() view returns(uint256)
func (_TokenHub *TokenHubSession) TENDECIMALS() (*big.Int, error) {
	return _TokenHub.Contract.TENDECIMALS(&_TokenHub.CallOpts)
}

// TENDECIMALS is a free data retrieval call binding the contract method 0x5d499b1b.
//
// Solidity: function TEN_DECIMALS() view returns(uint256)
func (_TokenHub *TokenHubCallerSession) TENDECIMALS() (*big.Int, error) {
	return _TokenHub.Contract.TENDECIMALS(&_TokenHub.CallOpts)
}

// TOKENHUBADDR is a free data retrieval call binding the contract method 0xfd6a6879.
//
// Solidity: function TOKEN_HUB_ADDR() view returns(address)
func (_TokenHub *TokenHubCaller) TOKENHUBADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "TOKEN_HUB_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TOKENHUBADDR is a free data retrieval call binding the contract method 0xfd6a6879.
//
// Solidity: function TOKEN_HUB_ADDR() view returns(address)
func (_TokenHub *TokenHubSession) TOKENHUBADDR() (common.Address, error) {
	return _TokenHub.Contract.TOKENHUBADDR(&_TokenHub.CallOpts)
}

// TOKENHUBADDR is a free data retrieval call binding the contract method 0xfd6a6879.
//
// Solidity: function TOKEN_HUB_ADDR() view returns(address)
func (_TokenHub *TokenHubCallerSession) TOKENHUBADDR() (common.Address, error) {
	return _TokenHub.Contract.TOKENHUBADDR(&_TokenHub.CallOpts)
}

// TOKENMANAGERADDR is a free data retrieval call binding the contract method 0x75d47a0a.
//
// Solidity: function TOKEN_MANAGER_ADDR() view returns(address)
func (_TokenHub *TokenHubCaller) TOKENMANAGERADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "TOKEN_MANAGER_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TOKENMANAGERADDR is a free data retrieval call binding the contract method 0x75d47a0a.
//
// Solidity: function TOKEN_MANAGER_ADDR() view returns(address)
func (_TokenHub *TokenHubSession) TOKENMANAGERADDR() (common.Address, error) {
	return _TokenHub.Contract.TOKENMANAGERADDR(&_TokenHub.CallOpts)
}

// TOKENMANAGERADDR is a free data retrieval call binding the contract method 0x75d47a0a.
//
// Solidity: function TOKEN_MANAGER_ADDR() view returns(address)
func (_TokenHub *TokenHubCallerSession) TOKENMANAGERADDR() (common.Address, error) {
	return _TokenHub.Contract.TOKENMANAGERADDR(&_TokenHub.CallOpts)
}

// TRANSFERINCHANNELID is a free data retrieval call binding the contract method 0x70fd5bad.
//
// Solidity: function TRANSFER_IN_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubCaller) TRANSFERINCHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "TRANSFER_IN_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// TRANSFERINCHANNELID is a free data retrieval call binding the contract method 0x70fd5bad.
//
// Solidity: function TRANSFER_IN_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubSession) TRANSFERINCHANNELID() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINCHANNELID(&_TokenHub.CallOpts)
}

// TRANSFERINCHANNELID is a free data retrieval call binding the contract method 0x70fd5bad.
//
// Solidity: function TRANSFER_IN_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) TRANSFERINCHANNELID() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINCHANNELID(&_TokenHub.CallOpts)
}

// TRANSFERINFAILUREINSUFFICIENTBALANCE is a free data retrieval call binding the contract method 0xa7c9f02d.
//
// Solidity: function TRANSFER_IN_FAILURE_INSUFFICIENT_BALANCE() view returns(uint8)
func (_TokenHub *TokenHubCaller) TRANSFERINFAILUREINSUFFICIENTBALANCE(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "TRANSFER_IN_FAILURE_INSUFFICIENT_BALANCE")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// TRANSFERINFAILUREINSUFFICIENTBALANCE is a free data retrieval call binding the contract method 0xa7c9f02d.
//
// Solidity: function TRANSFER_IN_FAILURE_INSUFFICIENT_BALANCE() view returns(uint8)
func (_TokenHub *TokenHubSession) TRANSFERINFAILUREINSUFFICIENTBALANCE() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINFAILUREINSUFFICIENTBALANCE(&_TokenHub.CallOpts)
}

// TRANSFERINFAILUREINSUFFICIENTBALANCE is a free data retrieval call binding the contract method 0xa7c9f02d.
//
// Solidity: function TRANSFER_IN_FAILURE_INSUFFICIENT_BALANCE() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) TRANSFERINFAILUREINSUFFICIENTBALANCE() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINFAILUREINSUFFICIENTBALANCE(&_TokenHub.CallOpts)
}

// TRANSFERINFAILURENONPAYABLERECIPIENT is a free data retrieval call binding the contract method 0xebf71d53.
//
// Solidity: function TRANSFER_IN_FAILURE_NON_PAYABLE_RECIPIENT() view returns(uint8)
func (_TokenHub *TokenHubCaller) TRANSFERINFAILURENONPAYABLERECIPIENT(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "TRANSFER_IN_FAILURE_NON_PAYABLE_RECIPIENT")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// TRANSFERINFAILURENONPAYABLERECIPIENT is a free data retrieval call binding the contract method 0xebf71d53.
//
// Solidity: function TRANSFER_IN_FAILURE_NON_PAYABLE_RECIPIENT() view returns(uint8)
func (_TokenHub *TokenHubSession) TRANSFERINFAILURENONPAYABLERECIPIENT() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINFAILURENONPAYABLERECIPIENT(&_TokenHub.CallOpts)
}

// TRANSFERINFAILURENONPAYABLERECIPIENT is a free data retrieval call binding the contract method 0xebf71d53.
//
// Solidity: function TRANSFER_IN_FAILURE_NON_PAYABLE_RECIPIENT() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) TRANSFERINFAILURENONPAYABLERECIPIENT() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINFAILURENONPAYABLERECIPIENT(&_TokenHub.CallOpts)
}

// TRANSFERINFAILURETIMEOUT is a free data retrieval call binding the contract method 0x8b87b21f.
//
// Solidity: function TRANSFER_IN_FAILURE_TIMEOUT() view returns(uint8)
func (_TokenHub *TokenHubCaller) TRANSFERINFAILURETIMEOUT(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "TRANSFER_IN_FAILURE_TIMEOUT")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// TRANSFERINFAILURETIMEOUT is a free data retrieval call binding the contract method 0x8b87b21f.
//
// Solidity: function TRANSFER_IN_FAILURE_TIMEOUT() view returns(uint8)
func (_TokenHub *TokenHubSession) TRANSFERINFAILURETIMEOUT() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINFAILURETIMEOUT(&_TokenHub.CallOpts)
}

// TRANSFERINFAILURETIMEOUT is a free data retrieval call binding the contract method 0x8b87b21f.
//
// Solidity: function TRANSFER_IN_FAILURE_TIMEOUT() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) TRANSFERINFAILURETIMEOUT() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINFAILURETIMEOUT(&_TokenHub.CallOpts)
}

// TRANSFERINFAILUREUNBOUNDTOKEN is a free data retrieval call binding the contract method 0xff9c0027.
//
// Solidity: function TRANSFER_IN_FAILURE_UNBOUND_TOKEN() view returns(uint8)
func (_TokenHub *TokenHubCaller) TRANSFERINFAILUREUNBOUNDTOKEN(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "TRANSFER_IN_FAILURE_UNBOUND_TOKEN")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// TRANSFERINFAILUREUNBOUNDTOKEN is a free data retrieval call binding the contract method 0xff9c0027.
//
// Solidity: function TRANSFER_IN_FAILURE_UNBOUND_TOKEN() view returns(uint8)
func (_TokenHub *TokenHubSession) TRANSFERINFAILUREUNBOUNDTOKEN() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINFAILUREUNBOUNDTOKEN(&_TokenHub.CallOpts)
}

// TRANSFERINFAILUREUNBOUNDTOKEN is a free data retrieval call binding the contract method 0xff9c0027.
//
// Solidity: function TRANSFER_IN_FAILURE_UNBOUND_TOKEN() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) TRANSFERINFAILUREUNBOUNDTOKEN() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINFAILUREUNBOUNDTOKEN(&_TokenHub.CallOpts)
}

// TRANSFERINFAILUREUNKNOWN is a free data retrieval call binding the contract method 0xf0148472.
//
// Solidity: function TRANSFER_IN_FAILURE_UNKNOWN() view returns(uint8)
func (_TokenHub *TokenHubCaller) TRANSFERINFAILUREUNKNOWN(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "TRANSFER_IN_FAILURE_UNKNOWN")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// TRANSFERINFAILUREUNKNOWN is a free data retrieval call binding the contract method 0xf0148472.
//
// Solidity: function TRANSFER_IN_FAILURE_UNKNOWN() view returns(uint8)
func (_TokenHub *TokenHubSession) TRANSFERINFAILUREUNKNOWN() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINFAILUREUNKNOWN(&_TokenHub.CallOpts)
}

// TRANSFERINFAILUREUNKNOWN is a free data retrieval call binding the contract method 0xf0148472.
//
// Solidity: function TRANSFER_IN_FAILURE_UNKNOWN() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) TRANSFERINFAILUREUNKNOWN() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINFAILUREUNKNOWN(&_TokenHub.CallOpts)
}

// TRANSFERINSUCCESS is a free data retrieval call binding the contract method 0xa496fba2.
//
// Solidity: function TRANSFER_IN_SUCCESS() view returns(uint8)
func (_TokenHub *TokenHubCaller) TRANSFERINSUCCESS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "TRANSFER_IN_SUCCESS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// TRANSFERINSUCCESS is a free data retrieval call binding the contract method 0xa496fba2.
//
// Solidity: function TRANSFER_IN_SUCCESS() view returns(uint8)
func (_TokenHub *TokenHubSession) TRANSFERINSUCCESS() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINSUCCESS(&_TokenHub.CallOpts)
}

// TRANSFERINSUCCESS is a free data retrieval call binding the contract method 0xa496fba2.
//
// Solidity: function TRANSFER_IN_SUCCESS() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) TRANSFERINSUCCESS() (uint8, error) {
	return _TokenHub.Contract.TRANSFERINSUCCESS(&_TokenHub.CallOpts)
}

// TRANSFEROUTCHANNELID is a free data retrieval call binding the contract method 0xfc3e5908.
//
// Solidity: function TRANSFER_OUT_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubCaller) TRANSFEROUTCHANNELID(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "TRANSFER_OUT_CHANNELID")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// TRANSFEROUTCHANNELID is a free data retrieval call binding the contract method 0xfc3e5908.
//
// Solidity: function TRANSFER_OUT_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubSession) TRANSFEROUTCHANNELID() (uint8, error) {
	return _TokenHub.Contract.TRANSFEROUTCHANNELID(&_TokenHub.CallOpts)
}

// TRANSFEROUTCHANNELID is a free data retrieval call binding the contract method 0xfc3e5908.
//
// Solidity: function TRANSFER_OUT_CHANNELID() view returns(uint8)
func (_TokenHub *TokenHubCallerSession) TRANSFEROUTCHANNELID() (uint8, error) {
	return _TokenHub.Contract.TRANSFEROUTCHANNELID(&_TokenHub.CallOpts)
}

// VALIDATORCONTRACTADDR is a free data retrieval call binding the contract method 0xf9a2bbc7.
//
// Solidity: function VALIDATOR_CONTRACT_ADDR() view returns(address)
func (_TokenHub *TokenHubCaller) VALIDATORCONTRACTADDR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "VALIDATOR_CONTRACT_ADDR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VALIDATORCONTRACTADDR is a free data retrieval call binding the contract method 0xf9a2bbc7.
//
// Solidity: function VALIDATOR_CONTRACT_ADDR() view returns(address)
func (_TokenHub *TokenHubSession) VALIDATORCONTRACTADDR() (common.Address, error) {
	return _TokenHub.Contract.VALIDATORCONTRACTADDR(&_TokenHub.CallOpts)
}

// VALIDATORCONTRACTADDR is a free data retrieval call binding the contract method 0xf9a2bbc7.
//
// Solidity: function VALIDATOR_CONTRACT_ADDR() view returns(address)
func (_TokenHub *TokenHubCallerSession) VALIDATORCONTRACTADDR() (common.Address, error) {
	return _TokenHub.Contract.VALIDATORCONTRACTADDR(&_TokenHub.CallOpts)
}

// AlreadyInit is a free data retrieval call binding the contract method 0xa78abc16.
//
// Solidity: function alreadyInit() view returns(bool)
func (_TokenHub *TokenHubCaller) AlreadyInit(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "alreadyInit")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AlreadyInit is a free data retrieval call binding the contract method 0xa78abc16.
//
// Solidity: function alreadyInit() view returns(bool)
func (_TokenHub *TokenHubSession) AlreadyInit() (bool, error) {
	return _TokenHub.Contract.AlreadyInit(&_TokenHub.CallOpts)
}

// AlreadyInit is a free data retrieval call binding the contract method 0xa78abc16.
//
// Solidity: function alreadyInit() view returns(bool)
func (_TokenHub *TokenHubCallerSession) AlreadyInit() (bool, error) {
	return _TokenHub.Contract.AlreadyInit(&_TokenHub.CallOpts)
}

// Bep2eContractDecimals is a free data retrieval call binding the contract method 0xa5cd588b.
//
// Solidity: function bep2eContractDecimals(address ) view returns(uint256)
func (_TokenHub *TokenHubCaller) Bep2eContractDecimals(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "bep2eContractDecimals", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Bep2eContractDecimals is a free data retrieval call binding the contract method 0xa5cd588b.
//
// Solidity: function bep2eContractDecimals(address ) view returns(uint256)
func (_TokenHub *TokenHubSession) Bep2eContractDecimals(arg0 common.Address) (*big.Int, error) {
	return _TokenHub.Contract.Bep2eContractDecimals(&_TokenHub.CallOpts, arg0)
}

// Bep2eContractDecimals is a free data retrieval call binding the contract method 0xa5cd588b.
//
// Solidity: function bep2eContractDecimals(address ) view returns(uint256)
func (_TokenHub *TokenHubCallerSession) Bep2eContractDecimals(arg0 common.Address) (*big.Int, error) {
	return _TokenHub.Contract.Bep2eContractDecimals(&_TokenHub.CallOpts, arg0)
}

// BscChainID is a free data retrieval call binding the contract method 0x493279b1.
//
// Solidity: function bscChainID() view returns(uint16)
func (_TokenHub *TokenHubCaller) BscChainID(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "bscChainID")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// BscChainID is a free data retrieval call binding the contract method 0x493279b1.
//
// Solidity: function bscChainID() view returns(uint16)
func (_TokenHub *TokenHubSession) BscChainID() (uint16, error) {
	return _TokenHub.Contract.BscChainID(&_TokenHub.CallOpts)
}

// BscChainID is a free data retrieval call binding the contract method 0x493279b1.
//
// Solidity: function bscChainID() view returns(uint16)
func (_TokenHub *TokenHubCallerSession) BscChainID() (uint16, error) {
	return _TokenHub.Contract.BscChainID(&_TokenHub.CallOpts)
}

// GetBep2SymbolByContractAddr is a free data retrieval call binding the contract method 0xbd466461.
//
// Solidity: function getBep2SymbolByContractAddr(address contractAddr) view returns(bytes32)
func (_TokenHub *TokenHubCaller) GetBep2SymbolByContractAddr(opts *bind.CallOpts, contractAddr common.Address) ([32]byte, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "getBep2SymbolByContractAddr", contractAddr)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBep2SymbolByContractAddr is a free data retrieval call binding the contract method 0xbd466461.
//
// Solidity: function getBep2SymbolByContractAddr(address contractAddr) view returns(bytes32)
func (_TokenHub *TokenHubSession) GetBep2SymbolByContractAddr(contractAddr common.Address) ([32]byte, error) {
	return _TokenHub.Contract.GetBep2SymbolByContractAddr(&_TokenHub.CallOpts, contractAddr)
}

// GetBep2SymbolByContractAddr is a free data retrieval call binding the contract method 0xbd466461.
//
// Solidity: function getBep2SymbolByContractAddr(address contractAddr) view returns(bytes32)
func (_TokenHub *TokenHubCallerSession) GetBep2SymbolByContractAddr(contractAddr common.Address) ([32]byte, error) {
	return _TokenHub.Contract.GetBep2SymbolByContractAddr(&_TokenHub.CallOpts, contractAddr)
}

// GetBoundBep2Symbol is a free data retrieval call binding the contract method 0xfc1a598f.
//
// Solidity: function getBoundBep2Symbol(address contractAddr) view returns(string)
func (_TokenHub *TokenHubCaller) GetBoundBep2Symbol(opts *bind.CallOpts, contractAddr common.Address) (string, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "getBoundBep2Symbol", contractAddr)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetBoundBep2Symbol is a free data retrieval call binding the contract method 0xfc1a598f.
//
// Solidity: function getBoundBep2Symbol(address contractAddr) view returns(string)
func (_TokenHub *TokenHubSession) GetBoundBep2Symbol(contractAddr common.Address) (string, error) {
	return _TokenHub.Contract.GetBoundBep2Symbol(&_TokenHub.CallOpts, contractAddr)
}

// GetBoundBep2Symbol is a free data retrieval call binding the contract method 0xfc1a598f.
//
// Solidity: function getBoundBep2Symbol(address contractAddr) view returns(string)
func (_TokenHub *TokenHubCallerSession) GetBoundBep2Symbol(contractAddr common.Address) (string, error) {
	return _TokenHub.Contract.GetBoundBep2Symbol(&_TokenHub.CallOpts, contractAddr)
}

// GetBoundContract is a free data retrieval call binding the contract method 0x3d713223.
//
// Solidity: function getBoundContract(string bep2Symbol) view returns(address)
func (_TokenHub *TokenHubCaller) GetBoundContract(opts *bind.CallOpts, bep2Symbol string) (common.Address, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "getBoundContract", bep2Symbol)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBoundContract is a free data retrieval call binding the contract method 0x3d713223.
//
// Solidity: function getBoundContract(string bep2Symbol) view returns(address)
func (_TokenHub *TokenHubSession) GetBoundContract(bep2Symbol string) (common.Address, error) {
	return _TokenHub.Contract.GetBoundContract(&_TokenHub.CallOpts, bep2Symbol)
}

// GetBoundContract is a free data retrieval call binding the contract method 0x3d713223.
//
// Solidity: function getBoundContract(string bep2Symbol) view returns(address)
func (_TokenHub *TokenHubCallerSession) GetBoundContract(bep2Symbol string) (common.Address, error) {
	return _TokenHub.Contract.GetBoundContract(&_TokenHub.CallOpts, bep2Symbol)
}

// GetContractAddrByBEP2Symbol is a free data retrieval call binding the contract method 0x59b92789.
//
// Solidity: function getContractAddrByBEP2Symbol(bytes32 bep2Symbol) view returns(address)
func (_TokenHub *TokenHubCaller) GetContractAddrByBEP2Symbol(opts *bind.CallOpts, bep2Symbol [32]byte) (common.Address, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "getContractAddrByBEP2Symbol", bep2Symbol)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetContractAddrByBEP2Symbol is a free data retrieval call binding the contract method 0x59b92789.
//
// Solidity: function getContractAddrByBEP2Symbol(bytes32 bep2Symbol) view returns(address)
func (_TokenHub *TokenHubSession) GetContractAddrByBEP2Symbol(bep2Symbol [32]byte) (common.Address, error) {
	return _TokenHub.Contract.GetContractAddrByBEP2Symbol(&_TokenHub.CallOpts, bep2Symbol)
}

// GetContractAddrByBEP2Symbol is a free data retrieval call binding the contract method 0x59b92789.
//
// Solidity: function getContractAddrByBEP2Symbol(bytes32 bep2Symbol) view returns(address)
func (_TokenHub *TokenHubCallerSession) GetContractAddrByBEP2Symbol(bep2Symbol [32]byte) (common.Address, error) {
	return _TokenHub.Contract.GetContractAddrByBEP2Symbol(&_TokenHub.CallOpts, bep2Symbol)
}

// GetMiniRelayFee is a free data retrieval call binding the contract method 0x149d14d9.
//
// Solidity: function getMiniRelayFee() view returns(uint256)
func (_TokenHub *TokenHubCaller) GetMiniRelayFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "getMiniRelayFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMiniRelayFee is a free data retrieval call binding the contract method 0x149d14d9.
//
// Solidity: function getMiniRelayFee() view returns(uint256)
func (_TokenHub *TokenHubSession) GetMiniRelayFee() (*big.Int, error) {
	return _TokenHub.Contract.GetMiniRelayFee(&_TokenHub.CallOpts)
}

// GetMiniRelayFee is a free data retrieval call binding the contract method 0x149d14d9.
//
// Solidity: function getMiniRelayFee() view returns(uint256)
func (_TokenHub *TokenHubCallerSession) GetMiniRelayFee() (*big.Int, error) {
	return _TokenHub.Contract.GetMiniRelayFee(&_TokenHub.CallOpts)
}

// RelayFee is a free data retrieval call binding the contract method 0x71d30863.
//
// Solidity: function relayFee() view returns(uint256)
func (_TokenHub *TokenHubCaller) RelayFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenHub.contract.Call(opts, &out, "relayFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RelayFee is a free data retrieval call binding the contract method 0x71d30863.
//
// Solidity: function relayFee() view returns(uint256)
func (_TokenHub *TokenHubSession) RelayFee() (*big.Int, error) {
	return _TokenHub.Contract.RelayFee(&_TokenHub.CallOpts)
}

// RelayFee is a free data retrieval call binding the contract method 0x71d30863.
//
// Solidity: function relayFee() view returns(uint256)
func (_TokenHub *TokenHubCallerSession) RelayFee() (*big.Int, error) {
	return _TokenHub.Contract.RelayFee(&_TokenHub.CallOpts)
}

// BatchTransferOutBNB is a paid mutator transaction binding the contract method 0x6e056520.
//
// Solidity: function batchTransferOutBNB(address[] recipientAddrs, uint256[] amounts, address[] refundAddrs, uint64 expireTime) payable returns(bool)
func (_TokenHub *TokenHubTransactor) BatchTransferOutBNB(opts *bind.TransactOpts, recipientAddrs []common.Address, amounts []*big.Int, refundAddrs []common.Address, expireTime uint64) (*types.Transaction, error) {
	return _TokenHub.contract.Transact(opts, "batchTransferOutBNB", recipientAddrs, amounts, refundAddrs, expireTime)
}

// BatchTransferOutBNB is a paid mutator transaction binding the contract method 0x6e056520.
//
// Solidity: function batchTransferOutBNB(address[] recipientAddrs, uint256[] amounts, address[] refundAddrs, uint64 expireTime) payable returns(bool)
func (_TokenHub *TokenHubSession) BatchTransferOutBNB(recipientAddrs []common.Address, amounts []*big.Int, refundAddrs []common.Address, expireTime uint64) (*types.Transaction, error) {
	return _TokenHub.Contract.BatchTransferOutBNB(&_TokenHub.TransactOpts, recipientAddrs, amounts, refundAddrs, expireTime)
}

// BatchTransferOutBNB is a paid mutator transaction binding the contract method 0x6e056520.
//
// Solidity: function batchTransferOutBNB(address[] recipientAddrs, uint256[] amounts, address[] refundAddrs, uint64 expireTime) payable returns(bool)
func (_TokenHub *TokenHubTransactorSession) BatchTransferOutBNB(recipientAddrs []common.Address, amounts []*big.Int, refundAddrs []common.Address, expireTime uint64) (*types.Transaction, error) {
	return _TokenHub.Contract.BatchTransferOutBNB(&_TokenHub.TransactOpts, recipientAddrs, amounts, refundAddrs, expireTime)
}

// BindToken is a paid mutator transaction binding the contract method 0x8eff336c.
//
// Solidity: function bindToken(bytes32 bep2Symbol, address contractAddr, uint256 decimals) returns()
func (_TokenHub *TokenHubTransactor) BindToken(opts *bind.TransactOpts, bep2Symbol [32]byte, contractAddr common.Address, decimals *big.Int) (*types.Transaction, error) {
	return _TokenHub.contract.Transact(opts, "bindToken", bep2Symbol, contractAddr, decimals)
}

// BindToken is a paid mutator transaction binding the contract method 0x8eff336c.
//
// Solidity: function bindToken(bytes32 bep2Symbol, address contractAddr, uint256 decimals) returns()
func (_TokenHub *TokenHubSession) BindToken(bep2Symbol [32]byte, contractAddr common.Address, decimals *big.Int) (*types.Transaction, error) {
	return _TokenHub.Contract.BindToken(&_TokenHub.TransactOpts, bep2Symbol, contractAddr, decimals)
}

// BindToken is a paid mutator transaction binding the contract method 0x8eff336c.
//
// Solidity: function bindToken(bytes32 bep2Symbol, address contractAddr, uint256 decimals) returns()
func (_TokenHub *TokenHubTransactorSession) BindToken(bep2Symbol [32]byte, contractAddr common.Address, decimals *big.Int) (*types.Transaction, error) {
	return _TokenHub.Contract.BindToken(&_TokenHub.TransactOpts, bep2Symbol, contractAddr, decimals)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x9a99b4f0.
//
// Solidity: function claimRewards(address to, uint256 amount) returns(uint256)
func (_TokenHub *TokenHubTransactor) ClaimRewards(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenHub.contract.Transact(opts, "claimRewards", to, amount)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x9a99b4f0.
//
// Solidity: function claimRewards(address to, uint256 amount) returns(uint256)
func (_TokenHub *TokenHubSession) ClaimRewards(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenHub.Contract.ClaimRewards(&_TokenHub.TransactOpts, to, amount)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x9a99b4f0.
//
// Solidity: function claimRewards(address to, uint256 amount) returns(uint256)
func (_TokenHub *TokenHubTransactorSession) ClaimRewards(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenHub.Contract.ClaimRewards(&_TokenHub.TransactOpts, to, amount)
}

// HandleAckPackage is a paid mutator transaction binding the contract method 0x831d65d1.
//
// Solidity: function handleAckPackage(uint8 channelId, bytes msgBytes) returns()
func (_TokenHub *TokenHubTransactor) HandleAckPackage(opts *bind.TransactOpts, channelId uint8, msgBytes []byte) (*types.Transaction, error) {
	return _TokenHub.contract.Transact(opts, "handleAckPackage", channelId, msgBytes)
}

// HandleAckPackage is a paid mutator transaction binding the contract method 0x831d65d1.
//
// Solidity: function handleAckPackage(uint8 channelId, bytes msgBytes) returns()
func (_TokenHub *TokenHubSession) HandleAckPackage(channelId uint8, msgBytes []byte) (*types.Transaction, error) {
	return _TokenHub.Contract.HandleAckPackage(&_TokenHub.TransactOpts, channelId, msgBytes)
}

// HandleAckPackage is a paid mutator transaction binding the contract method 0x831d65d1.
//
// Solidity: function handleAckPackage(uint8 channelId, bytes msgBytes) returns()
func (_TokenHub *TokenHubTransactorSession) HandleAckPackage(channelId uint8, msgBytes []byte) (*types.Transaction, error) {
	return _TokenHub.Contract.HandleAckPackage(&_TokenHub.TransactOpts, channelId, msgBytes)
}

// HandleFailAckPackage is a paid mutator transaction binding the contract method 0xc8509d81.
//
// Solidity: function handleFailAckPackage(uint8 channelId, bytes msgBytes) returns()
func (_TokenHub *TokenHubTransactor) HandleFailAckPackage(opts *bind.TransactOpts, channelId uint8, msgBytes []byte) (*types.Transaction, error) {
	return _TokenHub.contract.Transact(opts, "handleFailAckPackage", channelId, msgBytes)
}

// HandleFailAckPackage is a paid mutator transaction binding the contract method 0xc8509d81.
//
// Solidity: function handleFailAckPackage(uint8 channelId, bytes msgBytes) returns()
func (_TokenHub *TokenHubSession) HandleFailAckPackage(channelId uint8, msgBytes []byte) (*types.Transaction, error) {
	return _TokenHub.Contract.HandleFailAckPackage(&_TokenHub.TransactOpts, channelId, msgBytes)
}

// HandleFailAckPackage is a paid mutator transaction binding the contract method 0xc8509d81.
//
// Solidity: function handleFailAckPackage(uint8 channelId, bytes msgBytes) returns()
func (_TokenHub *TokenHubTransactorSession) HandleFailAckPackage(channelId uint8, msgBytes []byte) (*types.Transaction, error) {
	return _TokenHub.Contract.HandleFailAckPackage(&_TokenHub.TransactOpts, channelId, msgBytes)
}

// HandleSynPackage is a paid mutator transaction binding the contract method 0x1182b875.
//
// Solidity: function handleSynPackage(uint8 channelId, bytes msgBytes) returns(bytes)
func (_TokenHub *TokenHubTransactor) HandleSynPackage(opts *bind.TransactOpts, channelId uint8, msgBytes []byte) (*types.Transaction, error) {
	return _TokenHub.contract.Transact(opts, "handleSynPackage", channelId, msgBytes)
}

// HandleSynPackage is a paid mutator transaction binding the contract method 0x1182b875.
//
// Solidity: function handleSynPackage(uint8 channelId, bytes msgBytes) returns(bytes)
func (_TokenHub *TokenHubSession) HandleSynPackage(channelId uint8, msgBytes []byte) (*types.Transaction, error) {
	return _TokenHub.Contract.HandleSynPackage(&_TokenHub.TransactOpts, channelId, msgBytes)
}

// HandleSynPackage is a paid mutator transaction binding the contract method 0x1182b875.
//
// Solidity: function handleSynPackage(uint8 channelId, bytes msgBytes) returns(bytes)
func (_TokenHub *TokenHubTransactorSession) HandleSynPackage(channelId uint8, msgBytes []byte) (*types.Transaction, error) {
	return _TokenHub.Contract.HandleSynPackage(&_TokenHub.TransactOpts, channelId, msgBytes)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_TokenHub *TokenHubTransactor) Init(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenHub.contract.Transact(opts, "init")
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_TokenHub *TokenHubSession) Init() (*types.Transaction, error) {
	return _TokenHub.Contract.Init(&_TokenHub.TransactOpts)
}

// Init is a paid mutator transaction binding the contract method 0xe1c7392a.
//
// Solidity: function init() returns()
func (_TokenHub *TokenHubTransactorSession) Init() (*types.Transaction, error) {
	return _TokenHub.Contract.Init(&_TokenHub.TransactOpts)
}

// TransferOut is a paid mutator transaction binding the contract method 0xaa7415f5.
//
// Solidity: function transferOut(address contractAddr, address recipient, uint256 amount, uint64 expireTime) payable returns(bool)
func (_TokenHub *TokenHubTransactor) TransferOut(opts *bind.TransactOpts, contractAddr common.Address, recipient common.Address, amount *big.Int, expireTime uint64) (*types.Transaction, error) {
	return _TokenHub.contract.Transact(opts, "transferOut", contractAddr, recipient, amount, expireTime)
}

// TransferOut is a paid mutator transaction binding the contract method 0xaa7415f5.
//
// Solidity: function transferOut(address contractAddr, address recipient, uint256 amount, uint64 expireTime) payable returns(bool)
func (_TokenHub *TokenHubSession) TransferOut(contractAddr common.Address, recipient common.Address, amount *big.Int, expireTime uint64) (*types.Transaction, error) {
	return _TokenHub.Contract.TransferOut(&_TokenHub.TransactOpts, contractAddr, recipient, amount, expireTime)
}

// TransferOut is a paid mutator transaction binding the contract method 0xaa7415f5.
//
// Solidity: function transferOut(address contractAddr, address recipient, uint256 amount, uint64 expireTime) payable returns(bool)
func (_TokenHub *TokenHubTransactorSession) TransferOut(contractAddr common.Address, recipient common.Address, amount *big.Int, expireTime uint64) (*types.Transaction, error) {
	return _TokenHub.Contract.TransferOut(&_TokenHub.TransactOpts, contractAddr, recipient, amount, expireTime)
}

// UnbindToken is a paid mutator transaction binding the contract method 0xb99328c5.
//
// Solidity: function unbindToken(bytes32 bep2Symbol, address contractAddr) returns()
func (_TokenHub *TokenHubTransactor) UnbindToken(opts *bind.TransactOpts, bep2Symbol [32]byte, contractAddr common.Address) (*types.Transaction, error) {
	return _TokenHub.contract.Transact(opts, "unbindToken", bep2Symbol, contractAddr)
}

// UnbindToken is a paid mutator transaction binding the contract method 0xb99328c5.
//
// Solidity: function unbindToken(bytes32 bep2Symbol, address contractAddr) returns()
func (_TokenHub *TokenHubSession) UnbindToken(bep2Symbol [32]byte, contractAddr common.Address) (*types.Transaction, error) {
	return _TokenHub.Contract.UnbindToken(&_TokenHub.TransactOpts, bep2Symbol, contractAddr)
}

// UnbindToken is a paid mutator transaction binding the contract method 0xb99328c5.
//
// Solidity: function unbindToken(bytes32 bep2Symbol, address contractAddr) returns()
func (_TokenHub *TokenHubTransactorSession) UnbindToken(bep2Symbol [32]byte, contractAddr common.Address) (*types.Transaction, error) {
	return _TokenHub.Contract.UnbindToken(&_TokenHub.TransactOpts, bep2Symbol, contractAddr)
}

// UpdateParam is a paid mutator transaction binding the contract method 0xac431751.
//
// Solidity: function updateParam(string key, bytes value) returns()
func (_TokenHub *TokenHubTransactor) UpdateParam(opts *bind.TransactOpts, key string, value []byte) (*types.Transaction, error) {
	return _TokenHub.contract.Transact(opts, "updateParam", key, value)
}

// UpdateParam is a paid mutator transaction binding the contract method 0xac431751.
//
// Solidity: function updateParam(string key, bytes value) returns()
func (_TokenHub *TokenHubSession) UpdateParam(key string, value []byte) (*types.Transaction, error) {
	return _TokenHub.Contract.UpdateParam(&_TokenHub.TransactOpts, key, value)
}

// UpdateParam is a paid mutator transaction binding the contract method 0xac431751.
//
// Solidity: function updateParam(string key, bytes value) returns()
func (_TokenHub *TokenHubTransactorSession) UpdateParam(key string, value []byte) (*types.Transaction, error) {
	return _TokenHub.Contract.UpdateParam(&_TokenHub.TransactOpts, key, value)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TokenHub *TokenHubTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TokenHub.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TokenHub *TokenHubSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenHub.Contract.Fallback(&_TokenHub.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TokenHub *TokenHubTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TokenHub.Contract.Fallback(&_TokenHub.TransactOpts, calldata)
}

// TokenHubParamChangeIterator is returned from FilterParamChange and is used to iterate over the raw logs and unpacked data for ParamChange events raised by the TokenHub contract.
type TokenHubParamChangeIterator struct {
	Event *TokenHubParamChange // Event containing the contract specifics and raw log

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
func (it *TokenHubParamChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHubParamChange)
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
		it.Event = new(TokenHubParamChange)
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
func (it *TokenHubParamChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHubParamChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHubParamChange represents a ParamChange event raised by the TokenHub contract.
type TokenHubParamChange struct {
	Key   string
	Value []byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterParamChange is a free log retrieval operation binding the contract event 0x6cdb0ac70ab7f2e2d035cca5be60d89906f2dede7648ddbd7402189c1eeed17a.
//
// Solidity: event paramChange(string key, bytes value)
func (_TokenHub *TokenHubFilterer) FilterParamChange(opts *bind.FilterOpts) (*TokenHubParamChangeIterator, error) {

	logs, sub, err := _TokenHub.contract.FilterLogs(opts, "paramChange")
	if err != nil {
		return nil, err
	}
	return &TokenHubParamChangeIterator{contract: _TokenHub.contract, event: "paramChange", logs: logs, sub: sub}, nil
}

// WatchParamChange is a free log subscription operation binding the contract event 0x6cdb0ac70ab7f2e2d035cca5be60d89906f2dede7648ddbd7402189c1eeed17a.
//
// Solidity: event paramChange(string key, bytes value)
func (_TokenHub *TokenHubFilterer) WatchParamChange(opts *bind.WatchOpts, sink chan<- *TokenHubParamChange) (event.Subscription, error) {

	logs, sub, err := _TokenHub.contract.WatchLogs(opts, "paramChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHubParamChange)
				if err := _TokenHub.contract.UnpackLog(event, "paramChange", log); err != nil {
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
func (_TokenHub *TokenHubFilterer) ParseParamChange(log types.Log) (*TokenHubParamChange, error) {
	event := new(TokenHubParamChange)
	if err := _TokenHub.contract.UnpackLog(event, "paramChange", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenHubReceiveDepositIterator is returned from FilterReceiveDeposit and is used to iterate over the raw logs and unpacked data for ReceiveDeposit events raised by the TokenHub contract.
type TokenHubReceiveDepositIterator struct {
	Event *TokenHubReceiveDeposit // Event containing the contract specifics and raw log

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
func (it *TokenHubReceiveDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHubReceiveDeposit)
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
		it.Event = new(TokenHubReceiveDeposit)
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
func (it *TokenHubReceiveDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHubReceiveDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHubReceiveDeposit represents a ReceiveDeposit event raised by the TokenHub contract.
type TokenHubReceiveDeposit struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReceiveDeposit is a free log retrieval operation binding the contract event 0x6c98249d85d88c3753a04a22230f595e4dc8d3dc86c34af35deeeedc861b89db.
//
// Solidity: event receiveDeposit(address from, uint256 amount)
func (_TokenHub *TokenHubFilterer) FilterReceiveDeposit(opts *bind.FilterOpts) (*TokenHubReceiveDepositIterator, error) {

	logs, sub, err := _TokenHub.contract.FilterLogs(opts, "receiveDeposit")
	if err != nil {
		return nil, err
	}
	return &TokenHubReceiveDepositIterator{contract: _TokenHub.contract, event: "receiveDeposit", logs: logs, sub: sub}, nil
}

// WatchReceiveDeposit is a free log subscription operation binding the contract event 0x6c98249d85d88c3753a04a22230f595e4dc8d3dc86c34af35deeeedc861b89db.
//
// Solidity: event receiveDeposit(address from, uint256 amount)
func (_TokenHub *TokenHubFilterer) WatchReceiveDeposit(opts *bind.WatchOpts, sink chan<- *TokenHubReceiveDeposit) (event.Subscription, error) {

	logs, sub, err := _TokenHub.contract.WatchLogs(opts, "receiveDeposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHubReceiveDeposit)
				if err := _TokenHub.contract.UnpackLog(event, "receiveDeposit", log); err != nil {
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

// ParseReceiveDeposit is a log parse operation binding the contract event 0x6c98249d85d88c3753a04a22230f595e4dc8d3dc86c34af35deeeedc861b89db.
//
// Solidity: event receiveDeposit(address from, uint256 amount)
func (_TokenHub *TokenHubFilterer) ParseReceiveDeposit(log types.Log) (*TokenHubReceiveDeposit, error) {
	event := new(TokenHubReceiveDeposit)
	if err := _TokenHub.contract.UnpackLog(event, "receiveDeposit", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenHubRefundFailureIterator is returned from FilterRefundFailure and is used to iterate over the raw logs and unpacked data for RefundFailure events raised by the TokenHub contract.
type TokenHubRefundFailureIterator struct {
	Event *TokenHubRefundFailure // Event containing the contract specifics and raw log

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
func (it *TokenHubRefundFailureIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHubRefundFailure)
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
		it.Event = new(TokenHubRefundFailure)
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
func (it *TokenHubRefundFailureIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHubRefundFailureIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHubRefundFailure represents a RefundFailure event raised by the TokenHub contract.
type TokenHubRefundFailure struct {
	Bep2eAddr  common.Address
	RefundAddr common.Address
	Amount     *big.Int
	Status     uint32
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRefundFailure is a free log retrieval operation binding the contract event 0x203f9f67a785f4f81be4d48b109aa0c498d1bc8097ecc2627063f480cc5fe73e.
//
// Solidity: event refundFailure(address bep2eAddr, address refundAddr, uint256 amount, uint32 status)
func (_TokenHub *TokenHubFilterer) FilterRefundFailure(opts *bind.FilterOpts) (*TokenHubRefundFailureIterator, error) {

	logs, sub, err := _TokenHub.contract.FilterLogs(opts, "refundFailure")
	if err != nil {
		return nil, err
	}
	return &TokenHubRefundFailureIterator{contract: _TokenHub.contract, event: "refundFailure", logs: logs, sub: sub}, nil
}

// WatchRefundFailure is a free log subscription operation binding the contract event 0x203f9f67a785f4f81be4d48b109aa0c498d1bc8097ecc2627063f480cc5fe73e.
//
// Solidity: event refundFailure(address bep2eAddr, address refundAddr, uint256 amount, uint32 status)
func (_TokenHub *TokenHubFilterer) WatchRefundFailure(opts *bind.WatchOpts, sink chan<- *TokenHubRefundFailure) (event.Subscription, error) {

	logs, sub, err := _TokenHub.contract.WatchLogs(opts, "refundFailure")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHubRefundFailure)
				if err := _TokenHub.contract.UnpackLog(event, "refundFailure", log); err != nil {
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

// ParseRefundFailure is a log parse operation binding the contract event 0x203f9f67a785f4f81be4d48b109aa0c498d1bc8097ecc2627063f480cc5fe73e.
//
// Solidity: event refundFailure(address bep2eAddr, address refundAddr, uint256 amount, uint32 status)
func (_TokenHub *TokenHubFilterer) ParseRefundFailure(log types.Log) (*TokenHubRefundFailure, error) {
	event := new(TokenHubRefundFailure)
	if err := _TokenHub.contract.UnpackLog(event, "refundFailure", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenHubRefundSuccessIterator is returned from FilterRefundSuccess and is used to iterate over the raw logs and unpacked data for RefundSuccess events raised by the TokenHub contract.
type TokenHubRefundSuccessIterator struct {
	Event *TokenHubRefundSuccess // Event containing the contract specifics and raw log

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
func (it *TokenHubRefundSuccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHubRefundSuccess)
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
		it.Event = new(TokenHubRefundSuccess)
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
func (it *TokenHubRefundSuccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHubRefundSuccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHubRefundSuccess represents a RefundSuccess event raised by the TokenHub contract.
type TokenHubRefundSuccess struct {
	Bep2eAddr  common.Address
	RefundAddr common.Address
	Amount     *big.Int
	Status     uint32
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRefundSuccess is a free log retrieval operation binding the contract event 0xd468d4fa5e8fb4adc119b29a983fd0785e04af5cb8b7a3a69a47270c54b6901a.
//
// Solidity: event refundSuccess(address bep2eAddr, address refundAddr, uint256 amount, uint32 status)
func (_TokenHub *TokenHubFilterer) FilterRefundSuccess(opts *bind.FilterOpts) (*TokenHubRefundSuccessIterator, error) {

	logs, sub, err := _TokenHub.contract.FilterLogs(opts, "refundSuccess")
	if err != nil {
		return nil, err
	}
	return &TokenHubRefundSuccessIterator{contract: _TokenHub.contract, event: "refundSuccess", logs: logs, sub: sub}, nil
}

// WatchRefundSuccess is a free log subscription operation binding the contract event 0xd468d4fa5e8fb4adc119b29a983fd0785e04af5cb8b7a3a69a47270c54b6901a.
//
// Solidity: event refundSuccess(address bep2eAddr, address refundAddr, uint256 amount, uint32 status)
func (_TokenHub *TokenHubFilterer) WatchRefundSuccess(opts *bind.WatchOpts, sink chan<- *TokenHubRefundSuccess) (event.Subscription, error) {

	logs, sub, err := _TokenHub.contract.WatchLogs(opts, "refundSuccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHubRefundSuccess)
				if err := _TokenHub.contract.UnpackLog(event, "refundSuccess", log); err != nil {
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

// ParseRefundSuccess is a log parse operation binding the contract event 0xd468d4fa5e8fb4adc119b29a983fd0785e04af5cb8b7a3a69a47270c54b6901a.
//
// Solidity: event refundSuccess(address bep2eAddr, address refundAddr, uint256 amount, uint32 status)
func (_TokenHub *TokenHubFilterer) ParseRefundSuccess(log types.Log) (*TokenHubRefundSuccess, error) {
	event := new(TokenHubRefundSuccess)
	if err := _TokenHub.contract.UnpackLog(event, "refundSuccess", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenHubRewardToIterator is returned from FilterRewardTo and is used to iterate over the raw logs and unpacked data for RewardTo events raised by the TokenHub contract.
type TokenHubRewardToIterator struct {
	Event *TokenHubRewardTo // Event containing the contract specifics and raw log

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
func (it *TokenHubRewardToIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHubRewardTo)
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
		it.Event = new(TokenHubRewardTo)
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
func (it *TokenHubRewardToIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHubRewardToIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHubRewardTo represents a RewardTo event raised by the TokenHub contract.
type TokenHubRewardTo struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardTo is a free log retrieval operation binding the contract event 0xf8b71c64315fc33b2ead2adfa487955065152a8ac33d9d5193aafd7f45dc15a0.
//
// Solidity: event rewardTo(address to, uint256 amount)
func (_TokenHub *TokenHubFilterer) FilterRewardTo(opts *bind.FilterOpts) (*TokenHubRewardToIterator, error) {

	logs, sub, err := _TokenHub.contract.FilterLogs(opts, "rewardTo")
	if err != nil {
		return nil, err
	}
	return &TokenHubRewardToIterator{contract: _TokenHub.contract, event: "rewardTo", logs: logs, sub: sub}, nil
}

// WatchRewardTo is a free log subscription operation binding the contract event 0xf8b71c64315fc33b2ead2adfa487955065152a8ac33d9d5193aafd7f45dc15a0.
//
// Solidity: event rewardTo(address to, uint256 amount)
func (_TokenHub *TokenHubFilterer) WatchRewardTo(opts *bind.WatchOpts, sink chan<- *TokenHubRewardTo) (event.Subscription, error) {

	logs, sub, err := _TokenHub.contract.WatchLogs(opts, "rewardTo")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHubRewardTo)
				if err := _TokenHub.contract.UnpackLog(event, "rewardTo", log); err != nil {
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

// ParseRewardTo is a log parse operation binding the contract event 0xf8b71c64315fc33b2ead2adfa487955065152a8ac33d9d5193aafd7f45dc15a0.
//
// Solidity: event rewardTo(address to, uint256 amount)
func (_TokenHub *TokenHubFilterer) ParseRewardTo(log types.Log) (*TokenHubRewardTo, error) {
	event := new(TokenHubRewardTo)
	if err := _TokenHub.contract.UnpackLog(event, "rewardTo", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenHubTransferInSuccessIterator is returned from FilterTransferInSuccess and is used to iterate over the raw logs and unpacked data for TransferInSuccess events raised by the TokenHub contract.
type TokenHubTransferInSuccessIterator struct {
	Event *TokenHubTransferInSuccess // Event containing the contract specifics and raw log

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
func (it *TokenHubTransferInSuccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHubTransferInSuccess)
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
		it.Event = new(TokenHubTransferInSuccess)
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
func (it *TokenHubTransferInSuccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHubTransferInSuccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHubTransferInSuccess represents a TransferInSuccess event raised by the TokenHub contract.
type TokenHubTransferInSuccess struct {
	Bep2eAddr  common.Address
	RefundAddr common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTransferInSuccess is a free log retrieval operation binding the contract event 0x471eb9cc1ffe55ffadf15b32595415eb9d80f22e761d24bd6dffc607e1284d59.
//
// Solidity: event transferInSuccess(address bep2eAddr, address refundAddr, uint256 amount)
func (_TokenHub *TokenHubFilterer) FilterTransferInSuccess(opts *bind.FilterOpts) (*TokenHubTransferInSuccessIterator, error) {

	logs, sub, err := _TokenHub.contract.FilterLogs(opts, "transferInSuccess")
	if err != nil {
		return nil, err
	}
	return &TokenHubTransferInSuccessIterator{contract: _TokenHub.contract, event: "transferInSuccess", logs: logs, sub: sub}, nil
}

// WatchTransferInSuccess is a free log subscription operation binding the contract event 0x471eb9cc1ffe55ffadf15b32595415eb9d80f22e761d24bd6dffc607e1284d59.
//
// Solidity: event transferInSuccess(address bep2eAddr, address refundAddr, uint256 amount)
func (_TokenHub *TokenHubFilterer) WatchTransferInSuccess(opts *bind.WatchOpts, sink chan<- *TokenHubTransferInSuccess) (event.Subscription, error) {

	logs, sub, err := _TokenHub.contract.WatchLogs(opts, "transferInSuccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHubTransferInSuccess)
				if err := _TokenHub.contract.UnpackLog(event, "transferInSuccess", log); err != nil {
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

// ParseTransferInSuccess is a log parse operation binding the contract event 0x471eb9cc1ffe55ffadf15b32595415eb9d80f22e761d24bd6dffc607e1284d59.
//
// Solidity: event transferInSuccess(address bep2eAddr, address refundAddr, uint256 amount)
func (_TokenHub *TokenHubFilterer) ParseTransferInSuccess(log types.Log) (*TokenHubTransferInSuccess, error) {
	event := new(TokenHubTransferInSuccess)
	if err := _TokenHub.contract.UnpackLog(event, "transferInSuccess", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenHubTransferOutSuccessIterator is returned from FilterTransferOutSuccess and is used to iterate over the raw logs and unpacked data for TransferOutSuccess events raised by the TokenHub contract.
type TokenHubTransferOutSuccessIterator struct {
	Event *TokenHubTransferOutSuccess // Event containing the contract specifics and raw log

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
func (it *TokenHubTransferOutSuccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHubTransferOutSuccess)
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
		it.Event = new(TokenHubTransferOutSuccess)
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
func (it *TokenHubTransferOutSuccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHubTransferOutSuccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHubTransferOutSuccess represents a TransferOutSuccess event raised by the TokenHub contract.
type TokenHubTransferOutSuccess struct {
	Bep2eAddr  common.Address
	SenderAddr common.Address
	Amount     *big.Int
	RelayFee   *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTransferOutSuccess is a free log retrieval operation binding the contract event 0x74eab09b0e53aefc23f2e1b16da593f95c2dd49c6f5a23720463d10d9c330b2a.
//
// Solidity: event transferOutSuccess(address bep2eAddr, address senderAddr, uint256 amount, uint256 relayFee)
func (_TokenHub *TokenHubFilterer) FilterTransferOutSuccess(opts *bind.FilterOpts) (*TokenHubTransferOutSuccessIterator, error) {

	logs, sub, err := _TokenHub.contract.FilterLogs(opts, "transferOutSuccess")
	if err != nil {
		return nil, err
	}
	return &TokenHubTransferOutSuccessIterator{contract: _TokenHub.contract, event: "transferOutSuccess", logs: logs, sub: sub}, nil
}

// WatchTransferOutSuccess is a free log subscription operation binding the contract event 0x74eab09b0e53aefc23f2e1b16da593f95c2dd49c6f5a23720463d10d9c330b2a.
//
// Solidity: event transferOutSuccess(address bep2eAddr, address senderAddr, uint256 amount, uint256 relayFee)
func (_TokenHub *TokenHubFilterer) WatchTransferOutSuccess(opts *bind.WatchOpts, sink chan<- *TokenHubTransferOutSuccess) (event.Subscription, error) {

	logs, sub, err := _TokenHub.contract.WatchLogs(opts, "transferOutSuccess")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHubTransferOutSuccess)
				if err := _TokenHub.contract.UnpackLog(event, "transferOutSuccess", log); err != nil {
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

// ParseTransferOutSuccess is a log parse operation binding the contract event 0x74eab09b0e53aefc23f2e1b16da593f95c2dd49c6f5a23720463d10d9c330b2a.
//
// Solidity: event transferOutSuccess(address bep2eAddr, address senderAddr, uint256 amount, uint256 relayFee)
func (_TokenHub *TokenHubFilterer) ParseTransferOutSuccess(log types.Log) (*TokenHubTransferOutSuccess, error) {
	event := new(TokenHubTransferOutSuccess)
	if err := _TokenHub.contract.UnpackLog(event, "transferOutSuccess", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TokenHubUnexpectedPackageIterator is returned from FilterUnexpectedPackage and is used to iterate over the raw logs and unpacked data for UnexpectedPackage events raised by the TokenHub contract.
type TokenHubUnexpectedPackageIterator struct {
	Event *TokenHubUnexpectedPackage // Event containing the contract specifics and raw log

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
func (it *TokenHubUnexpectedPackageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenHubUnexpectedPackage)
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
		it.Event = new(TokenHubUnexpectedPackage)
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
func (it *TokenHubUnexpectedPackageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenHubUnexpectedPackageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenHubUnexpectedPackage represents a UnexpectedPackage event raised by the TokenHub contract.
type TokenHubUnexpectedPackage struct {
	ChannelId uint8
	MsgBytes  []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUnexpectedPackage is a free log retrieval operation binding the contract event 0x41ce201247b6ceb957dcdb217d0b8acb50b9ea0e12af9af4f5e7f38902101605.
//
// Solidity: event unexpectedPackage(uint8 channelId, bytes msgBytes)
func (_TokenHub *TokenHubFilterer) FilterUnexpectedPackage(opts *bind.FilterOpts) (*TokenHubUnexpectedPackageIterator, error) {

	logs, sub, err := _TokenHub.contract.FilterLogs(opts, "unexpectedPackage")
	if err != nil {
		return nil, err
	}
	return &TokenHubUnexpectedPackageIterator{contract: _TokenHub.contract, event: "unexpectedPackage", logs: logs, sub: sub}, nil
}

// WatchUnexpectedPackage is a free log subscription operation binding the contract event 0x41ce201247b6ceb957dcdb217d0b8acb50b9ea0e12af9af4f5e7f38902101605.
//
// Solidity: event unexpectedPackage(uint8 channelId, bytes msgBytes)
func (_TokenHub *TokenHubFilterer) WatchUnexpectedPackage(opts *bind.WatchOpts, sink chan<- *TokenHubUnexpectedPackage) (event.Subscription, error) {

	logs, sub, err := _TokenHub.contract.WatchLogs(opts, "unexpectedPackage")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenHubUnexpectedPackage)
				if err := _TokenHub.contract.UnpackLog(event, "unexpectedPackage", log); err != nil {
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

// ParseUnexpectedPackage is a log parse operation binding the contract event 0x41ce201247b6ceb957dcdb217d0b8acb50b9ea0e12af9af4f5e7f38902101605.
//
// Solidity: event unexpectedPackage(uint8 channelId, bytes msgBytes)
func (_TokenHub *TokenHubFilterer) ParseUnexpectedPackage(log types.Log) (*TokenHubUnexpectedPackage, error) {
	event := new(TokenHubUnexpectedPackage)
	if err := _TokenHub.contract.UnpackLog(event, "unexpectedPackage", log); err != nil {
		return nil, err
	}
	return event, nil
}
