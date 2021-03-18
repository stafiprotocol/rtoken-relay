package cosmos

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	rpcHttp "github.com/tendermint/tendermint/rpc/client/http"
	"os"
)

type Connection struct {
	url        string
	symbol     core.RSymbol
	poolKeys   []string //hex string
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
