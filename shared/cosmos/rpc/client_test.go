package rpc_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"github.com/stretchr/testify/assert"
	rpcHttp "github.com/tendermint/tendermint/rpc/client/http"
	"strings"
	"testing"
)

var client *rpc.Client

//eda331e37bf66b2393c4c271e384dfaa2bfcdd35
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

	client, _ = rpc.NewClient(rpcClient, key, "stargate-final", "recipient")
	client.SetGasPrice("0.000001umuon")
	client.SetDenom("umuon")
}

//{"height":"901192","txhash":"327DA2048B6D66BCB27C0F1A6D1E407D88FE719B95A30D108B5906FD6934F7B1","codespace":"","code":0,"data":"0A060A0473656E64","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos1ak3nrcmm7e4j8y7ycfc78pxl4g4lehf43vw6wu\"},{\"key\":\"sender\",\"value\":\"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf\"},{\"key\":\"amount\",\"value\":\"100umuon\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"send"},{"key":"sender","value":"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf"},{"key":"module","value":"bank"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"cosmos1ak3nrcmm7e4j8y7ycfc78pxl4g4lehf43vw6wu"},{"key":"sender","value":"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf"},{"key":"amount","value":"100umuon"}]}]}],"info":"","gas_wanted":"200000","gas_used":"51169","tx":null,"timestamp":""}
//{"height":"903451","txhash":"0E4F8F8FF7A3B67121711DA17FBE5AE8CB25DB272DDBF7DC0E02122947266604","codespace":"","code":0,"data":"0A060A0473656E64","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos1ak3nrcmm7e4j8y7ycfc78pxl4g4lehf43vw6wu\"},{\"key\":\"sender\",\"value\":\"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf\"},{\"key\":\"amount\",\"value\":\"10umuon\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"send"},{"key":"sender","value":"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf"},{"key":"module","value":"bank"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"cosmos1ak3nrcmm7e4j8y7ycfc78pxl4g4lehf43vw6wu"},{"key":"sender","value":"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf"},{"key":"amount","value":"10umuon"}]}]}],"info":"","gas_wanted":"200000","gas_used":"51159","tx":null,"timestamp":""}
//block hash 0x16E8297663210ABF6937FE4C1C139D4BACD0D27A22EFD9E3FE06B1DA8E3F7BB3
func TestClient_SendTo(t *testing.T) {
	err := client.SingleTransferTo(addrMultiSig1, types.NewCoins(types.NewInt64Coin(client.GetDenom(), 10)))
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

	signature1, err := client.SignMultiSigRawTxWithSeq(56, rawTx, "key1")
	assert.NoError(t, err)
	t.Log(string(signature1))

	signature2, err := client.SignMultiSigRawTxWithSeq(56, rawTx, "key3")
	assert.NoError(t, err)
	t.Log(string(signature2))
	//signature3, err := client.SignMultiSigRawTxWithSeq(56, rawTx, "key2")
	//assert.NoError(t, err)
	//t.Log(string(signature2))

	hash, tx, err := client.AssembleMultiSigTxWithSeq(50, rawTx, [][]byte{ signature2, signature1})
	assert.NoError(t, err)
	t.Log(hex.EncodeToString(hash))
	t.Log(string(tx))

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

	_, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1})
	assert.NoError(t, err)

	_, err = client.BroadcastTx(tx)
	assert.ErrorIs(t, err, errors.New(fmt.Sprintf("Boradcast err with res.code: %d", 4)))
}

func TestClient_QueryTxByHash(t *testing.T) {
	res, err := client.QueryTxByHash("6C017062FD3F48F13B640E5FEDD59EB050B148E67EF12EC0A511442D32BD4C88")
	t.Log(err)
	assert.NoError(t, err)
	for _, msg := range res.GetTx().GetMsgs() {

		t.Log(msg.String())
		t.Log(msg.Type())
		t.Log(msg.Route())
	}
}

func TestClient_QueryDelegationRewards(t *testing.T) {
	res, err := client.QueryDelegationRewards(addrMultiSig1, addrValidator, 0)
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

	_, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2, signature3})
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

	_, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2})
	assert.NoError(t, err)

	txHash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(txHash)
}

func TestClient_GenMultiSigRawWithdrawDeleRewardTx(t *testing.T) {
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawWithdrawDeleRewardTx(addrMultiSig1, addrValidatorTestnetAteam)
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	h, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2})
	assert.NoError(t, err)
	t.Log(hex.EncodeToString(h))
	hash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(hash)
}

func TestClient_GenMultiSigRawWithdrawDeleRewardAndDelegratTx(t *testing.T) {
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawWithdrawAllRewardThenDeleTx(addrMultiSig1, 1035322)
	assert.NoError(t, err)
	t.Log(string(rawTx))

	signature1, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	h, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2})
	assert.NoError(t, err)
	t.Log(string(tx))
	t.Log(hex.EncodeToString(h))
	hash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(hash)
}

func TestClient_GenMultiSigRawWithdrawAllRewardTx(t *testing.T) {
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawWithdrawAllRewardTx(addrMultiSig1, 1037755)
	assert.NoError(t, err)
	t.Log(string(rawTx))

	signature1, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	h, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2})
	assert.NoError(t, err)
	t.Log(string(tx))
	t.Log(hex.EncodeToString(h))
	hash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(hash)
}

//0xf954ad81a546df9de6c79051ca67f7d0d08e9d861604c76fb8767dce2ce8d4f8
func TestClient_GenMultiSigRawUnDelegateTx(t *testing.T) {
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawUnDelegateTx(addrMultiSig1, addrValidatorTestnetAteam, types.NewCoin(client.GetDenom(), types.NewInt(15)))
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	hash, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2})
	assert.NoError(t, err)
	t.Log("hash", hex.EncodeToString(hash))
	hash2, err := client.BroadcastTx(tx)
	t.Log("hash2", hash2)
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

	_, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2})
	assert.NoError(t, err)

	txHash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(txHash)
}

func TestGetPubKey(t *testing.T) {
	account, _ := client.QueryAccount(addrReceive)
	t.Log(hex.EncodeToString(account.GetPubKey().Bytes()))

	res, err := client.QueryTxByHash("327DA2048B6D66BCB27C0F1A6D1E407D88FE719B95A30D108B5906FD6934F7B1")
	if err != nil {
		t.Fatal(err)
	}
	msgs := res.GetTx().GetMsgs()
	for i, _ := range msgs {
		if msgs[i].Type() == xBankTypes.TypeMsgSend {
			msg, _ := msgs[i].(*xBankTypes.MsgSend)
			t.Log(msg.Amount.AmountOf("umuon").Uint64())
		}
	}

}

func TestClient_Sign(t *testing.T) {
	bts, err := hex.DecodeString("0E4F8F8FF7A3B67121711DA17FBE5AE8CB25DB272DDBF7DC0E02122947266604")
	assert.NoError(t, err)
	sigs, pubkey, err := client.Sign("recipient", bts)
	assert.NoError(t, err)
	t.Log(hex.EncodeToString(sigs))
	//4c6902bda88424923c62f95b3e3ead40769edab4ec794108d1c18994fac90d490087815823bd1a8af3d6a0271538cef4622b4b500a6253d2bd4c80d38e95aa6d
	t.Log(hex.EncodeToString(pubkey.Bytes()))
	//02e7710b4f7147c10ad90da06b69d2d6b8ff46786ef55a3f1e889c33de2bf0b416
}

func TestAddress(t *testing.T) {
	addrKey1, _ := types.AccAddressFromBech32("cosmos1a8mg9rj4nklhmwkf5vva8dvtgx4ucd9yjasret")
	addrKey2, _ := types.AccAddressFromBech32("cosmos1ztquzhpkve7szl99jkugq4l8jtpnhln76aetam")
	addrKey3, _ := types.AccAddressFromBech32("cosmos12zz2hm02sxe9f4pwt7y5q9wjhcu98vnuwmjz4x")
	t.Log(hex.EncodeToString(addrKey1.Bytes()))
	t.Log(hex.EncodeToString(addrKey2.Bytes()))
	t.Log(hex.EncodeToString(addrKey3.Bytes()))
}

func TestClient_QueryDelegations(t *testing.T) {
	res, err := client.QueryDelegations(addrMultiSig1, 0)
	assert.NoError(t, err)
	t.Log(res.String())
}

func TestClient_QueryBalance(t *testing.T) {
	res, err := client.QueryBalance(addrMultiSig1, "umuon", 0)
	assert.NoError(t, err)
	t.Log(res.Balance.Amount)
}

func TestClient_QueryDelegationTotalRewards(t *testing.T) {
	res, err := client.QueryDelegationTotalRewards(addrMultiSig1, 0)
	assert.NoError(t, err)
	t.Log(res.Total.AmountOf(client.GetDenom()).TruncateInt())
}

func TestClient_GetSequence(t *testing.T) {
	seq,err:=client.GetSequence(0,addrMultiSig1)
	assert.NoError(t,err)
	t.Log(seq)
	t.Log(hex.EncodeToString(addrValidatorTestnetAteam.Bytes()))
}