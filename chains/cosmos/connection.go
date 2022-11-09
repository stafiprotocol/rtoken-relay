package cosmos

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/JFJun/go-substrate-crypto/ss58"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Connection struct {
	url                string
	symbol             core.RSymbol
	validatorTargets   []types.ValAddress
	currentHeight      int64
	increaseRewardInfo *IncreaseRewardInfo
	poolClients        map[string]*cosmos.PoolClient //map[addressHexStr]subClient
	log                log15.Logger
	stop               <-chan int
}

type IncreaseRewardInfo struct {
	Era    uint32 `json:"era"`
	Amount BigInt `json:"amount"`
}

type BigInt struct {
	big.Int
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	if string(p) == "null" {
		return nil
	}

	var val string
	err := json.Unmarshal(p, &val)
	if err != nil {
		return err
	}

	var z big.Int
	_, ok := z.SetString(val, 10)
	if !ok {
		return fmt.Errorf("not a valid big integer: %s", p)
	}
	b.Int = z
	return nil
}

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {
	subKeys := make(map[string]string)
	for _, account := range cfg.Accounts {
		subKey, ok := cfg.Opts[account].(string)
		if !ok || len(subKey) == 0 {
			return nil, fmt.Errorf("account %s has no subKeys", account)
		}
		subKeys[account] = subKey
	}

	targets := make([]types.ValAddress, 0)
	optTargets := cfg.Opts["validatorTargets"]
	log.Info("NewConnection", "targets", optTargets)
	if optTargets != nil {
		if tmpCares, ok := optTargets.([]interface{}); ok {
			for _, tc := range tmpCares {
				target, ok := tc.(string)
				if !ok {
					panic("target not string")
				}
				val, err := types.ValAddressFromBech32(target)
				if err != nil {
					panic(err)
				}
				targets = append(targets, val)
			}
		} else {
			panic("opt validatorTarget not string array")
		}
	}
	if len(targets) == 0 {
		panic("validatorTargets empty")
	}

	optIncreaseRewardInfo := cfg.Opts["increaseRewardInfo"]
	log.Info("NewConnection", "increaseRewardInfo", optIncreaseRewardInfo)

	var increaseRewardInfo *IncreaseRewardInfo
	if optIncreaseRewardInfo != nil {
		bts, err := json.Marshal(optIncreaseRewardInfo)
		if err != nil {
			return nil, err
		}

		increaseRewardInfo = &IncreaseRewardInfo{}
		err = json.Unmarshal(bts, increaseRewardInfo)
		if err != nil {
			return nil, err
		}
	}

	chainId, ok := cfg.Opts[config.ChainId].(string)
	if !ok || len(chainId) == 0 {
		return nil, errors.New("config must has chainId")
	}
	denom, ok := cfg.Opts[config.Denom].(string)
	if !ok || len(chainId) == 0 {
		return nil, errors.New("config must has denom")
	}
	gasPrice, ok := cfg.Opts[config.GasPrice].(string)
	if !ok || len(gasPrice) == 0 {
		return nil, errors.New("config must has gasPrice")
	}
	eraBlockNumberStr, ok := cfg.Opts[config.EraBlockNumber].(string)
	if !ok {
		return nil, errors.New("config must has eraBlockNumber")
	}
	eraBlockNumber, err := strconv.ParseInt(eraBlockNumberStr, 10, 64)
	if err != nil {
		return nil, errors.New("config parse eraBlockNumber failed")
	}
	if eraBlockNumber == 0 {
		return nil, errors.New("eraBlockNumber is zero")
	}

	subClients := make(map[string]*cosmos.PoolClient)
	keys := make([]string, 0)

	fmt.Printf("Will open cosmos wallet from <%s>. \nPlease ", cfg.KeystorePath)
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, cfg.KeystorePath, os.Stdin)
	if err != nil {
		return nil, err
	}

	for _, account := range cfg.Accounts {
		client, err := rpc.NewClient(key, chainId, account, gasPrice, denom, cfg.Endpoint)
		if err != nil {
			return nil, err
		}

		keyInfo, err := key.Key(account)
		if err != nil {
			return nil, err
		}
		addrHexStr := hex.EncodeToString(keyInfo.GetAddress().Bytes())
		subClients[addrHexStr] = cosmos.NewPoolClient(log, client, subKeys[account], eraBlockNumber)
		keys = append(keys, addrHexStr)
	}

	if len(keys) == 0 {
		return nil, errors.New("no poolKeys")
	}

	return &Connection{
		url:                cfg.Endpoint,
		symbol:             cfg.Symbol,
		log:                log,
		stop:               stop,
		poolClients:        subClients,
		validatorTargets:   targets,
		increaseRewardInfo: increaseRewardInfo,
	}, nil
}

func (c *Connection) TransferVerify(r *submodel.BondRecord) (submodel.BondReason, error) {
	hashStr := hex.EncodeToString(r.Txhash)
	poolClient, err := c.GetOnePoolClient()
	if err != nil {
		return submodel.BondReasonDefault, err
	}

	//check tx hash
	txRes, err := c.GetTx(poolClient, hashStr)
	if err != nil {
		return submodel.TxhashUnmatch, nil
	}

	if txRes.Empty() {
		return submodel.TxhashUnmatch, nil
	}

	if txRes.Code != 0 {
		return submodel.TxhashUnmatch, nil
	}

	//check block hash
	blockRes, err := c.GetBlock(poolClient, txRes.Height)
	if err != nil {
		return submodel.BondReasonDefault, err
	}
	if !bytes.Equal(blockRes.BlockID.Hash, r.Blockhash) {
		return submodel.BlockhashUnmatch, nil
	}

	//check amount and pool
	amountIsMatch := false
	poolIsMatch := false
	var fromAddressStr string

	msgs := txRes.GetTx().GetMsgs()
	for i, _ := range msgs {
		if sendMsg, ok := msgs[i].(*xBankTypes.MsgSend); ok {
			toAddr, err := types.AccAddressFromBech32(sendMsg.ToAddress)
			if err == nil {
				//amount and pool address must in one message
				if bytes.Equal(toAddr.Bytes(), r.Pool) &&
					sendMsg.Amount.AmountOf(poolClient.GetRpcClient().GetDenom()).
						Equal(types.NewIntFromBigInt(r.Amount.Int)) {
					poolIsMatch = true
					amountIsMatch = true
					fromAddressStr = sendMsg.FromAddress
				}
			}

		}
	}
	if !amountIsMatch {
		return submodel.AmountUnmatch, nil
	}
	if !poolIsMatch {
		return submodel.PoolUnmatch, nil
	}

	//check pubkey
	fromAddress, err := types.AccAddressFromBech32(fromAddressStr)
	if err != nil {
		return submodel.BondReasonDefault, err
	}
	accountRes, err := poolClient.GetRpcClient().QueryAccount(fromAddress)
	if err != nil {
		return submodel.BondReasonDefault, err
	}

	if !bytes.Equal(accountRes.GetPubKey().Bytes(), r.Pubkey) {
		return submodel.PubkeyUnmatch, nil
	}
	//check memo
	var memoInTx string
	tx, err := poolClient.GetRpcClient().GetTxConfig().TxDecoder()(txRes.Tx.GetValue())
	if err == nil {
		memoTx, ok := tx.(types.TxWithMemo)
		if ok {
			memoInTx = memoTx.GetMemo()
		}
	}

	bonderAddr, err := ss58.Encode(r.Bonder[:], ss58.StafiPrefix)
	if err != nil {
		return submodel.PubkeyUnmatch, nil // memo unmatch
	}

	if memoInTx != bonderAddr {
		return submodel.PubkeyUnmatch, nil // memo unmatch
	}

	return submodel.Pass, nil
}

// fetch one for query
func (c *Connection) GetOnePoolClient() (*cosmos.PoolClient, error) {
	for _, sub := range c.poolClients {
		if sub != nil {
			return sub, nil
		}
	}
	return nil, errors.New("no subClient")
}

func (c *Connection) GetPoolClient(poolAddrHexStr string) (*cosmos.PoolClient, error) {
	if sub, exist := c.poolClients[poolAddrHexStr]; exist {
		return sub, nil
	}
	return nil, errors.New("subClient of this pool not exist")
}

func (c *Connection) GetTx(poolClient *cosmos.PoolClient, txHash string) (*types.TxResponse, error) {
	var txRes *types.TxResponse
	var err error
	retryTx := 0
	for {
		if retryTx >= BlockRetryLimit {
			return nil, errors.New("QueryTxByHash reach retry limit")
		}
		txRes, err = poolClient.GetRpcClient().QueryTxByHash(txHash)
		if err != nil {
			c.log.Warn(fmt.Sprintf("QueryTxByHash err: %s ,will retry queryTx after %f second", err, BlockRetryInterval.Seconds()))
			time.Sleep(BlockRetryInterval)
			retryTx++
			continue
		}
		currentHeight, err := poolClient.GetRpcClient().GetCurrentBLockHeight()
		if err != nil {
			c.log.Warn(fmt.Sprintf("GetCurrentBLockHeight err: %s ,will retry queryTx after %f second", err, BlockRetryInterval.Seconds()))
			time.Sleep(BlockRetryInterval)
			retryTx++
			continue
		}

		if txRes.Height+BlockConfirmNumber > currentHeight {
			c.log.Warn(fmt.Sprintf("confirm number is smaller than %d ,will retry queryTx after %f second", BlockConfirmNumber, BlockRetryInterval.Seconds()))
			time.Sleep(BlockRetryInterval)
			retryTx++
			continue
		} else {
			break
		}

	}
	return txRes, nil
}

func (c *Connection) GetBlock(poolClient *cosmos.PoolClient, height int64) (*ctypes.ResultBlock, error) {
	var blockRes *ctypes.ResultBlock
	var err error
	retryTx := 0
	for {
		if retryTx >= BlockRetryLimit {
			return nil, errors.New("QueryBlock reach retry limit")
		}
		blockRes, err = poolClient.GetRpcClient().QueryBlock(height)
		if err != nil {
			c.log.Warn(fmt.Sprintf("QueryBlock err: %s ,will retry queryBlock after %f second", err, BlockRetryInterval.Seconds()))
			time.Sleep(BlockRetryInterval)
			retryTx++
			continue
		}
		break
	}
	return blockRes, nil
}
