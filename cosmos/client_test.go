package cosmos_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stafiprotocol/rtoken-relay/cosmos"
	"github.com/stretchr/testify/assert"
	rpcHttp "github.com/tendermint/tendermint/rpc/client/http"
	"strings"
	"testing"
)

var client *cosmos.Client
var addrMultiSig1, _ = types.AccAddressFromBech32("cosmos1ak3nrcmm7e4j8y7ycfc78pxl4g4lehf43vw6wu")
var addrReceive, _ = types.AccAddressFromBech32("cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf")

func init() {
	rpcClient, err := rpcHttp.New("http://127.0.0.1:26657", "/websocket")
	if err != nil {
		panic(err)
	}
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, "/Users/tpkeeper/.gaia", strings.NewReader("tpkeeper\n"))
	if err != nil {
		panic(err)
	}

	client = cosmos.NewClient(rpcClient, key, "my-test-chain", "validator")
}

func TestClient_SendTo(t *testing.T) {
	err := client.SendTo(addrMultiSig1, types.NewCoins(types.NewInt64Coin("stake", 100)))
	assert.NoError(t, err)
}

func TestClient_GenRawTx(t *testing.T) {
	err := client.SetFromKey("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenRawTx(addrReceive, types.NewCoins(types.NewInt64Coin("stake", 10)))
	assert.NoError(t, err)
	t.Log(string(rawTx))
}

func TestClient_SignRawTx(t *testing.T) {
	err := client.SetFromKey("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenRawTx(addrReceive, types.NewCoins(types.NewInt64Coin("stake", 10)))
	assert.NoError(t, err)

	signature, err := client.SignRawTx(rawTx, "key1")
	assert.NoError(t, err)
	t.Log(string(signature))
}

func TestClient_CreateMultiSigTx(t *testing.T) {
	err := client.SetFromKey("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenRawTx(addrReceive, types.NewCoins(types.NewInt64Coin("stake", 10)))
	assert.NoError(t, err)

	signature1, err := client.SignRawTx(rawTx, "key1")
	assert.NoError(t, err)
	signature2, err := client.SignRawTx(rawTx, "key2")
	assert.NoError(t, err)

	tx, err := client.CreateMultiSigTx(rawTx, [][]byte{signature1, signature2})
	assert.NoError(t, err)

	err = client.BroadcastTx(tx)
	assert.NoError(t, err)
}
