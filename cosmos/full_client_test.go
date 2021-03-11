package cosmos_test

import (
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	subClientTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/cosmos"
	"github.com/stafiprotocol/rtoken-relay/cosmos/rpc"
	"github.com/stretchr/testify/assert"
	rpcHttp "github.com/tendermint/tendermint/rpc/client/http"
	"math/big"
	"strings"
	"testing"
)

var client *rpc.Client
var fullClient cosmos.FullClient
var addrMultiSig1, _ = types.AccAddressFromBech32("cosmos1ak3nrcmm7e4j8y7ycfc78pxl4g4lehf43vw6wu")

func init() {
	rpcClient, err := rpcHttp.New("http://127.0.0.1:26657", "/websocket")
	if err != nil {
		panic(err)
	}
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, "/Users/tpkeeper/.gaia", strings.NewReader("tpkeeper\n"))
	if err != nil {
		panic(err)
	}

	client = rpc.NewClient(rpcClient, key, "stargate-final", "recipient")
	client.SetGasPrice("0.000001umuon")
	client.SetDenom("umuon")
	keyInfo, err := key.Key("recipient")
	if err != nil {
		panic(err)
	}

	fullClient = cosmos.FullClient{Keys: []keyring.Info{keyInfo}, SubClients: map[keyring.Info]*rpc.Client{keyInfo: client}}
}

func TestFullClient_TransferVerify(t *testing.T) {

	pubkey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeAccPub, "cosmospub1addwnpepqtnhzz60w9ruzzkepksxk6wj66u073ncdm6450c73zwr8h3t7z6pvxhuhtr")
	assert.NoError(t, err)
	txHashBts, _ := hex.DecodeString("CED193699DDC061A186BAD6C19BCDC4C9EB1C723605CD5A32E51A46E78421753")
	blockHashBts, _ := hex.DecodeString("597A5614245909C700083A3EBAF26A54A6D2A8EBF75DD9AE83477DCCF57483D0")

	bondRecord := conn.BondRecord{
		Txhash:    txHashBts,
		Blockhash: blockHashBts,
		Pubkey:    pubkey.Bytes(),
		Amount:    subClientTypes.NewU128(*big.NewInt(1000)),
		Pool:      addrMultiSig1.Bytes()}
	reason, err := fullClient.TransferVerify(&bondRecord)
	assert.NoError(t, err)
	t.Log(reason)
}

func TestFullClient_CurrentEra(t *testing.T) {
	era,err:=fullClient.CurrentEra()
	assert.NoError(t,err)
	t.Log(era)
}
