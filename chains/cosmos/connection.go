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
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	rpcHttp "github.com/tendermint/tendermint/rpc/client/http"
	"os"
)

type Connection struct {
	url        string
	symbol     core.RSymbol
	poolKeys   []string //pool addr hex string
	subClients map[string]*cosmos.SubClient
	log        log15.Logger
	stop       <-chan int
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

	subClients := make(map[string]*cosmos.SubClient)
	keys := make([]string, 0)

	fmt.Printf("Will open cosmos wallet from <%s>. Please ", cfg.KeystorePath)
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, cfg.KeystorePath, os.Stdin)
	if err != nil {
		return nil, err
	}

	for _, account := range cfg.Accounts {
		rpcClient, err := rpcHttp.New(cfg.Endpoint, "/websocket")
		if err != nil {
			panic(err)
		}
		client, err := rpc.NewClient(rpcClient, key, "stargate-final", account)
		if err != nil {
			return nil, err
		}
		client.SetGasPrice("0.000001umuon")
		client.SetDenom("umuon")
		keyInfo, err := key.Key(account)
		if err != nil {
			return nil, err
		}
		addrHexStr := hex.EncodeToString(keyInfo.GetAddress().Bytes())
		subClients[addrHexStr] = cosmos.NewSubClient(log, client, subKeys[account])
		keys = append(keys, addrHexStr)
	}

	if len(keys) == 0 {
		return nil, errors.New("no poolKeys")
	}

	return &Connection{
		url:        cfg.Endpoint,
		symbol:     cfg.Symbol,
		log:        log,
		stop:       stop,
		poolKeys:   keys,
		subClients: subClients,
	}, nil
}

func (fc *Connection) TransferVerify(r *core.BondRecord) (core.BondReason, error) {
	return core.Pass, nil
	if len(fc.subClients) == 0 || len(fc.poolKeys) == 0 {
		return "", fmt.Errorf("no subClient")
	}
	hashStr := hex.EncodeToString(r.Txhash)
	client := fc.subClients[fc.poolKeys[0]]

	//check tx hash
	res, err := client.GetRpcClient().QueryTxByHash(hashStr)
	if err != nil {
		return core.TxhashUnmatch, err
	}
	if res.Empty() {
		return core.TxhashUnmatch, nil
	}

	//check block hash
	blockRes, err := client.GetRpcClient().QueryBlock(res.Height)
	if err != nil {
		return core.BlockhashUnmatch, err
	}
	if !bytes.Equal(blockRes.BlockID.Hash, r.Blockhash) {
		return core.BlockhashUnmatch, nil
	}

	//check amount and pool
	amountIsMatch := false
	poolIsMatch := false
	fromAddressStr := ""
	msgs := res.GetTx().GetMsgs()
	for i, _ := range msgs {
		if msgs[i].Type() == xBankTypes.TypeMsgSend {
			sendMsg, ok := msgs[i].(*xBankTypes.MsgSend)
			if ok {
				toAddr, err := types.AccAddressFromBech32(sendMsg.ToAddress)
				if err == nil {
					//amount and pool address must in one message
					if bytes.Equal(toAddr.Bytes(), r.Pool) &&
						sendMsg.Amount.AmountOf(client.GetRpcClient().GetDenom()).Equal(types.NewIntFromBigInt(r.Amount.Int)) {
						poolIsMatch = true
						amountIsMatch = true
						fromAddressStr = sendMsg.FromAddress
					}
				}

			}

		}
	}
	if !amountIsMatch {
		return core.AmountUnmatch, nil
	}
	if !poolIsMatch {
		return core.PoolUnmatch, nil
	}

	//check pubkey
	fromAddress, err := types.AccAddressFromBech32(fromAddressStr)
	if err != nil {
		return core.PubkeyUnmatch, err
	}
	accountRes, err := client.GetRpcClient().QueryAccount(fromAddress)
	if err != nil {
		return core.PubkeyUnmatch, err
	}

	if !bytes.Equal(accountRes.GetPubKey().Bytes(), r.Pubkey) {
		return core.PubkeyUnmatch, nil
	}
	return core.Pass, nil
}
