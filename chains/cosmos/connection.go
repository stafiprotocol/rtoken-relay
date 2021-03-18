package cosmos

import (
	"errors"
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
	poolKeys   []keyring.Info
	subClients map[keyring.Info]*cosmos.SubClient
	log        log15.Logger
	stop       <-chan int
}

type subKeys map[string][]string

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {

	subKeyMap, ok := cfg.Opts["subKeys"].(subKeys)
	//todo check number
	if !ok {
		return nil, errors.New("no subKeys")
	}

	subClients := make(map[keyring.Info]*cosmos.SubClient)
	keys := make([]keyring.Info, 0)

	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, cfg.KeystorePath, os.Stdin)
	if err != nil {
		return nil, err
	}

	for _, account := range cfg.Accounts {
		rpcClient, err := rpcHttp.New(cfg.Endpoint, "/websocket")
		if err != nil {
			panic(err)
		}
		client := rpc.NewClient(rpcClient, key, "stargate-final", account)
		client.SetGasPrice("0.000001umuon")
		client.SetDenom("umuon")
		keyInfo, err := key.Key(account)
		if err != nil {
			return nil, err
		}
		subClients[keyInfo] = cosmos.NewSubClient(log, client, subKeyMap[account])
		keys = append(keys, keyInfo)
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
