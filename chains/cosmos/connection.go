package cosmos

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	rpcHttp "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"os"
	"time"
)

type Connection struct {
	url           string
	symbol        core.RSymbol
	currentHeight int64                         //todo enable automic
	poolClients   map[string]*cosmos.PoolClient //map[addressHexStr]subClient
	log           log15.Logger
	stop          <-chan int
}

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {
	subKeys := make(map[string]string)
	for _, account := range cfg.Accounts {
		subKey, ok := cfg.Opts[account].(string)
		if !ok || len(subKey) == 0 {
			return nil, errors.New(fmt.Sprintf("account %s has no subKeys", account))
		}
		subKeys[account] = subKey
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

	subClients := make(map[string]*cosmos.PoolClient)
	keys := make([]string, 0)

	fmt.Printf("Will open cosmos wallet from <%s>. \nPlease ", cfg.KeystorePath)
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, cfg.KeystorePath, os.Stdin)
	if err != nil {
		return nil, err
	}

	//todo some params just for test
	for _, account := range cfg.Accounts {
		rpcClient, err := rpcHttp.New(cfg.Endpoint, "/websocket")
		if err != nil {
			panic(err)
		}
		client, err := rpc.NewClient(rpcClient, key, chainId, account)
		if err != nil {
			return nil, err
		}
		client.SetGasPrice(gasPrice)
		client.SetDenom(denom)
		keyInfo, err := key.Key(account)
		if err != nil {
			return nil, err
		}
		addrHexStr := hex.EncodeToString(keyInfo.GetAddress().Bytes())
		subClients[addrHexStr] = cosmos.NewPoolClient(log, client, subKeys[account])
		keys = append(keys, addrHexStr)
	}

	if len(keys) == 0 {
		return nil, errors.New("no poolKeys")
	}

	return &Connection{
		url:         cfg.Endpoint,
		symbol:      cfg.Symbol,
		log:         log,
		stop:        stop,
		poolClients: subClients,
	}, nil
}

func (c *Connection) TransferVerify(r *submodel.BondRecord) (submodel.BondReason, error) {
	//todo test only,must rm on release
	//return submodel.Pass, nil

	hashStr := hex.EncodeToString(r.Txhash)

	poolClient, err := c.GetOnePoolClient()
	if err != nil {
		return submodel.BondReasonDefault, err
	}

	//check tx hash
	var txRes *types.TxResponse
	retryTx := BlockRetryLimit
	for {
		if retryTx <= 0 {
			return submodel.TxhashUnmatch, errors.New("QueryTxByHash reach retryTx limit")
		}
		txRes, err = poolClient.GetRpcClient().QueryTxByHash(hashStr)
		if err != nil || txRes == nil || txRes.Empty() {
			c.log.Warn(fmt.Sprintf("QueryTxByHash empty or err: %s ,will retryTx after %f second", err, BlockRetryInterval.Seconds()))
			time.Sleep(BlockRetryInterval)
			retryTx--
			continue
		}
		if txRes.Height+BlockConfirmNumber > c.currentHeight {
			c.log.Warn(fmt.Sprintf("confirm number is smaller than %d ,will retryTx after %f second", BlockConfirmNumber, BlockRetryInterval.Seconds()))
			time.Sleep(BlockRetryInterval)
			retryTx--
			continue
		} else {
			break
		}

	}

	//check code
	if txRes.Code != 0 {
		return submodel.TxhashUnmatch, nil
	}

	//check block hash
	var blockRes *ctypes.ResultBlock
	retryBlock := BlockRetryLimit
	for {
		if retryBlock <= 0 {
			return submodel.BlockhashUnmatch, errors.New("QueryBlock reach retryTx limit")
		}
		blockRes, err = poolClient.GetRpcClient().QueryBlock(txRes.Height)
		if err != nil || blockRes == nil {
			c.log.Warn(fmt.Sprintf("QueryBlock empty or err: %s ,will retryTx after %f second", err, BlockRetryInterval.Seconds()))
			time.Sleep(BlockRetryInterval)
			continue
		} else {
			break
		}
	}

	if !bytes.Equal(blockRes.BlockID.Hash, r.Blockhash) {
		return submodel.BlockhashUnmatch, nil
	}

	//check amount and pool
	amountIsMatch := false
	poolIsMatch := false
	fromAddressStr := ""
	msgs := txRes.GetTx().GetMsgs()
	for i, _ := range msgs {
		if msgs[i].Type() == xBankTypes.TypeMsgSend {
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
		return submodel.PubkeyUnmatch, err
	}
	accountRes, err := poolClient.GetRpcClient().QueryAccount(fromAddress)
	if err != nil {
		return submodel.PubkeyUnmatch, err
	}

	if !bytes.Equal(accountRes.GetPubKey().Bytes(), r.Pubkey) {
		return submodel.PubkeyUnmatch, nil
	}
	//check memo
	//tx.ConvertTxToStdTx(poolClient.GetRpcClient().GetLegacyAmino(), txRes.Data)

	return submodel.Pass, nil
}

//fetch one for query
func (c *Connection) GetOnePoolClient() (*cosmos.PoolClient, error) {
	//todo check connect state
	for _, sub := range c.poolClients {
		if sub != nil {
			return sub, nil
		}
	}
	return nil, errors.New("no subClient")
}

func (c *Connection) GetPoolClient(poolAddrHexStr string) (*cosmos.PoolClient, error) {
	//todo check connect state
	if sub, exist := c.poolClients[poolAddrHexStr]; exist {
		return sub, nil
	}
	return nil, errors.New("subClient of this pool not exist")
}

func (c *Connection) GetCurrentEra() (uint32, error) {
	client, err := c.GetOnePoolClient()
	if err != nil {
		return 0, err
	}

	status, err := client.GetRpcClient().GetStatus()
	if err != nil {
		return 0, err
	}
	era := uint32(status.SyncInfo.LatestBlockHeight / cosmos.EraBlockNumber)
	return era, nil
}
