package cosmos_test

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stafiprotocol/rtoken-relay/cosmos"
	"github.com/stretchr/testify/assert"
	rpcHttp "github.com/tendermint/tendermint/rpc/client/http"
	"strings"
	"testing"
)

var client *cosmos.Client
var addrMultiSig1, _ = types.AccAddressFromBech32("cosmos1ak3nrcmm7e4j8y7ycfc78pxl4g4lehf43vw6wu")
var addrReceive, _ = types.AccAddressFromBech32("cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf")
var addrValidator, _ = types.ValAddressFromBech32("cosmosvaloper1y6zkfcvwkpqz89z7rwu9kcdm4kc7uc4e5y5a2r")
var addrKey1, _ = types.AccAddressFromBech32("cosmos1a8mg9rj4nklhmwkf5vva8dvtgx4ucd9yjasret")
var addrValidatorTestnet, _ = types.ValAddressFromBech32("cosmosvaloper17tpddyr578avyn95xngkjl8nl2l2tf6auh8kpc")
var addrValidatorTestnet2, _ = types.ValAddressFromBech32("cosmosvaloper19xczxvvdg8h67sk3cccrvxlj0ruyw3360rctfa")
var addrValidatorTestnetAteam, _ = types.ValAddressFromBech32("cosmosvaloper105gvcjgs6s4j5ws9srckx0drt4x8cwgywplh7p")
var addrValidatorTestnetStation, _ = types.ValAddressFromBech32("cosmosvaloper1x5wgh6vwye60wv3dtshs9dmqggwfx2ldk5cvqu")

func init() {
	rpcClient, err := rpcHttp.New("http://127.0.0.1:26657", "/websocket")
	if err != nil {
		panic(err)
	}
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, "/Users/tpkeeper/.gaia", strings.NewReader("tpkeeper\n"))
	if err != nil {
		panic(err)
	}

	client = cosmos.NewClient(rpcClient, key, "stargate-final", "recipient")
	client.SetGasPrice("0.000001umuon")
	client.SetDenom("umuon")
}

func TestClient_SendTo(t *testing.T) {
	err := client.SingleTransferTo(addrMultiSig1, types.NewCoins(types.NewInt64Coin(client.GetDenom(), 1000)))
	assert.NoError(t, err)
}

func TestClient_ReDelegate(t *testing.T) {
	h, err := client.SingleReDelegate(addrValidatorTestnetAteam, addrValidatorTestnetStation,
		types.NewCoin(client.GetDenom(), types.NewInt(10)))
	assert.NoError(t, err)
	t.Log("hash", h)
}

func TestClient_GenRawTx(t *testing.T) {
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawTransferTx(addrReceive, types.NewCoins(types.NewInt64Coin(client.GetDenom(), 10)))
	assert.NoError(t, err)
	t.Log(string(rawTx))
}

func TestClient_SignRawTx(t *testing.T) {
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawTransferTx(addrReceive, types.NewCoins(types.NewInt64Coin(client.GetDenom(), 10)))
	assert.NoError(t, err)

	signature, err := client.SignMultiSigRawTx(rawTx, "key1")
	assert.NoError(t, err)
	t.Log(string(signature))
}

func TestClient_CreateMultiSigTx(t *testing.T) {
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawTransferTx(addrReceive, types.NewCoins(types.NewInt64Coin(client.GetDenom(), 10)))
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key1")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	tx, err := client.CreateMultiSigTx(rawTx, [][]byte{signature1, signature2})
	assert.NoError(t, err)

	txHash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(txHash)
}

func TestClient_BroadcastTx(t *testing.T) {
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawTransferTx(addrReceive, types.NewCoins(types.NewInt64Coin(client.GetDenom(), 10)))
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key1")
	assert.NoError(t, err)

	tx, err := client.CreateMultiSigTx(rawTx, [][]byte{signature1})
	assert.NoError(t, err)

	_, err = client.BroadcastTx(tx)
	assert.ErrorIs(t, err, errors.New(fmt.Sprintf("Boradcast err with res.code: %d", 4)))
}

func TestClient_QueryTxByHash(t *testing.T) {
	res, err := client.QueryTxByHash("6C017062FD3F48F13B640E5FEDD59EB050B148E67EF12EC0A511442D32BD4C88")
	assert.NoError(t, err)
	for _, msg := range res.GetTx().GetMsgs() {

		t.Log(msg.String())
		t.Log(msg.Type())
		t.Log(msg.Route())
	}
}

func TestClient_QueryDelegationRewards(t *testing.T) {
	res, err := client.QueryDelegationRewards(addrMultiSig1, addrValidator)
	assert.NoError(t, err)
	t.Log(res.GetRewards().AmountOf(client.GetDenom()))
}

func TestClient_GenMultiSigRawDelegateTx(t *testing.T) {
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawDelegateTx(addrMultiSig1, addrValidatorTestnetAteam, types.NewCoin(client.GetDenom(), types.NewInt(100)))
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key1")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature3, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	tx, err := client.CreateMultiSigTx(rawTx, [][]byte{signature1, signature2, signature3})
	assert.NoError(t, err)

	_, err = client.BroadcastTx(tx)
	assert.NoError(t, err)
}

func TestClient_GenMultiSigRawReDelegateTx(t *testing.T) {
	err := client.SetFromName("multiSign1")

	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawReDelegateTx(addrMultiSig1, addrValidatorTestnetAteam, addrValidatorTestnetStation,
		types.NewCoin(client.GetDenom(), types.NewInt(10)))
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key1")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	//signature3, err := client.SignMultiSigRawTx(rawTx, "key3")
	//assert.NoError(t, err)

	tx, err := client.CreateMultiSigTx(rawTx, [][]byte{signature1, signature2})
	assert.NoError(t, err)

	txHash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(txHash)
}

func TestClient_GenMultiSigRawWithdrawDeleRewardTx(t *testing.T) {
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawWithdrawDeleRewardTx(addrMultiSig1, addrValidator)
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	tx, err := client.CreateMultiSigTx(rawTx, [][]byte{signature1, signature2})
	assert.NoError(t, err)

	_, err = client.BroadcastTx(tx)
	assert.NoError(t, err)
}

func TestClient_GenMultiSigRawUnDelegateTx(t *testing.T) {
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawUnDelegateTx(addrMultiSig1, addrValidator, types.NewCoin(client.GetDenom(), types.NewInt(100)))
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	tx, err := client.CreateMultiSigTx(rawTx, [][]byte{signature1, signature2})
	assert.NoError(t, err)

	_, err = client.BroadcastTx(tx)
	assert.NoError(t, err)
}

func TestClient_GenMultiSigRawBatchTransferTx(t *testing.T) {
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	out1 := xBankTypes.Output{
		Address: addrReceive.String(),
		Coins:   types.NewCoins(types.NewCoin(client.GetDenom(), types.NewInt(10))),
	}
	out2 := xBankTypes.Output{
		Address: addrKey1.String(),
		Coins:   types.NewCoins(types.NewCoin(client.GetDenom(), types.NewInt(10))),
	}

	rawTx, err := client.GenMultiSigRawBatchTransferTx([]xBankTypes.Output{out1, out2})
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	tx, err := client.CreateMultiSigTx(rawTx, [][]byte{signature1, signature2})
	assert.NoError(t, err)

	txHash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(txHash)
}
