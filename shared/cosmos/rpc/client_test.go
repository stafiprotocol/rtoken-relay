package rpc_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"

	// "strings"
	"sync"
	"testing"
	"time"

	"github.com/JFJun/go-substrate-crypto/ss58"
	// "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"github.com/stretchr/testify/assert"
)

var client *rpc.Client

// eda331e37bf66b2393c4c271e384dfaa2bfcdd35
var addrMultiSig1, _ = types.AccAddressFromBech32("cosmos1em384d8ek3y8nlugapz7p5k5skg58j66je3las")
var addrMultiSig2, _ = types.AccAddressFromBech32("cosmos1wmk9ys49zxgmx77pqs7cjnpamnnxuspqu2r87y")
var addrReceive, _ = types.AccAddressFromBech32("cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf")
var addrValidator, _ = types.ValAddressFromBech32("cosmosvaloper18pty5rk5fgsfnls5pq8ajxcwkf0we8l8q3g8pe")
var addrKey1, _ = types.AccAddressFromBech32("cosmos1a8mg9rj4nklhmwkf5vva8dvtgx4ucd9yjasret")
var addrValidatorTestnet2, _ = types.ValAddressFromBech32("cosmosvaloper19xczxvvdg8h67sk3cccrvxlj0ruyw3360rctfa")

var addrValidatorTestnet, _ = types.ValAddressFromBech32("cosmosvaloper17tpddyr578avyn95xngkjl8nl2l2tf6auh8kpc")
var addrValidatorTestnetStation, _ = types.ValAddressFromBech32("cosmosvaloper1x5wgh6vwye60wv3dtshs9dmqggwfx2ldk5cvqu")
var addrValidatorTestnetAteam, _ = types.ValAddressFromBech32("cosmosvaloper105gvcjgs6s4j5ws9srckx0drt4x8cwgywplh7p")

var adrValidatorTestnetTecos, _ = types.ValAddressFromBech32("cosmosvaloper1p7e37nztj62mmra8xhgqde7sql3llhhu6hvcx8")
var adrValidatorEverStake, _ = types.ValAddressFromBech32("cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3")
var adrValidatorForbole, _ = types.ValAddressFromBech32("cosmosvaloper1w96rrh9sx0h7n7qak00l90un0kx5wala2prmxt")

var addrReleaseAddress, _ = types.AccAddressFromBech32("cosmos1wmk9ys49zxgmx77pqs7cjnpamnnxuspqu2r87y")
var addrReleaseVal, _ = types.ValAddressFromBech32("cosmosvaloper14jlpmqquh0gste6nzf4dn43kc8h50l6fmx6dfs")
var addrReleaseValWetez, _ = types.ValAddressFromBech32("cosmosvaloper1s05va5d09xlq3et8mapsesqh6r5lqy7mkhwshm")

func TestGetAddrHex(t *testing.T) {
	t.Log("cosmosvaloper17tpddyr578avyn95xngkjl8nl2l2tf6auh8kpc", hexutil.Encode(addrValidatorTestnet.Bytes()))
	t.Log("cosmosvaloper1x5wgh6vwye60wv3dtshs9dmqggwfx2ldk5cvqu", hexutil.Encode(addrValidatorTestnetStation.Bytes()))
	t.Log("cosmosvaloper105gvcjgs6s4j5ws9srckx0drt4x8cwgywplh7p", hexutil.Encode(addrValidatorTestnetAteam.Bytes()))

	t.Log("cosmosvaloper1p7e37nztj62mmra8xhgqde7sql3llhhu6hvcx8", hexutil.Encode(adrValidatorTestnetTecos.Bytes()))
	t.Log("cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", hexutil.Encode(adrValidatorEverStake.Bytes()))
	t.Log("cosmosvaloper1w96rrh9sx0h7n7qak00l90un0kx5wala2prmxt", hexutil.Encode(adrValidatorForbole.Bytes()))
	//client_test.go:36: cosmosvaloper17tpddyr578avyn95xngkjl8nl2l2tf6auh8kpc 0xf2c2d69074f1fac24cb434d1697cf3fabea5a75d
	//client_test.go:38: cosmosvaloper1x5wgh6vwye60wv3dtshs9dmqggwfx2ldk5cvqu 0x351c8be98e2674f7322d5c2f02b760421c932bed
	//client_test.go:37: cosmosvaloper105gvcjgs6s4j5ws9srckx0drt4x8cwgywplh7p 0x7d10cc4910d42b2a3a0580f1633da35d4c7c3904

	//client_test.go:40: cosmosvaloper1p7e37nztj62mmra8xhgqde7sql3llhhu6hvcx8 0x0fb31f4c4b9695bd8fa735d006e7d007e3ffdefc
	//client_test.go:41: cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3 0x5a7f68bf60a3100937e42aad6bdc111f72c55e62
	//client_test.go:42: cosmosvaloper1w96rrh9sx0h7n7qak00l90un0kx5wala2prmxt 0x717431dcb033efe9f81db3dff2bf937d8d4777fd
}

func initClient() {
	// key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, "/Users/tpkeeper/.gaia", strings.NewReader("tpkeeper\n"))
	// if err != nil {
	// 	panic(err)
	// }

	// client, err = rpc.NewClient(key, "stargate-final", "key0", "0.04umuon", "umuon", "https://testcosmosrpc.wetez.io:443")
	// client, err = rpc.NewClient(key, "stargate-final", "recipient", "0.04umuon", "umuon", "http://127.0.0.1:26657")
	// client, _ = rpc.NewClient(key, "cosmoshub-4", "self", "0.00001uatom", "uatom", "https://cosmos-rpc1.stafi.io:443")
	var err error
	// client, err = rpc.NewClient(nil, "cosmoshub-4", "", "0.00001uatom", "uatom", "https://cosmos-rpc1.stafi.io:443")
	client, err = rpc.NewClient(nil, "local-cosmos", "", "0.00001stake", "stake", "http://127.0.0.1:26657")
	if err != nil {
		panic(err)
	}
}

// {"height":"901192","txhash":"327DA2048B6D66BCB27C0F1A6D1E407D88FE719B95A30D108B5906FD6934F7B1","codespace":"","code":0,"data":"0A060A0473656E64","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos1ak3nrcmm7e4j8y7ycfc78pxl4g4lehf43vw6wu\"},{\"key\":\"sender\",\"value\":\"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf\"},{\"key\":\"amount\",\"value\":\"100umuon\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"send"},{"key":"sender","value":"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf"},{"key":"module","value":"bank"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"cosmos1ak3nrcmm7e4j8y7ycfc78pxl4g4lehf43vw6wu"},{"key":"sender","value":"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf"},{"key":"amount","value":"100umuon"}]}]}],"info":"","gas_wanted":"200000","gas_used":"51169","tx":null,"timestamp":""}
// {"height":"903451","txhash":"0E4F8F8FF7A3B67121711DA17FBE5AE8CB25DB272DDBF7DC0E02122947266604","codespace":"","code":0,"data":"0A060A0473656E64","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos1ak3nrcmm7e4j8y7ycfc78pxl4g4lehf43vw6wu\"},{\"key\":\"sender\",\"value\":\"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf\"},{\"key\":\"amount\",\"value\":\"10umuon\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"send"},{"key":"sender","value":"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf"},{"key":"module","value":"bank"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"cosmos1ak3nrcmm7e4j8y7ycfc78pxl4g4lehf43vw6wu"},{"key":"sender","value":"cosmos1cgs647rewxyzh5wu4e606kk7qyuj5f8hk20rgf"},{"key":"amount","value":"10umuon"}]}]}],"info":"","gas_wanted":"200000","gas_used":"51159","tx":null,"timestamp":""}
// block hash 0x16E8297663210ABF6937FE4C1C139D4BACD0D27A22EFD9E3FE06B1DA8E3F7BB3
func TestClient_SendTo(t *testing.T) {
	initClient()
	err := client.SingleTransferTo(addrMultiSig1, types.NewCoins(types.NewInt64Coin(client.GetDenom(), 50000)))
	assert.NoError(t, err)
}

func TestClient_ReDelegate(t *testing.T) {
	initClient()
	h, err := client.SingleReDelegate(addrValidatorTestnetAteam, addrValidatorTestnetStation,
		types.NewCoin(client.GetDenom(), types.NewInt(10)))
	assert.NoError(t, err)
	t.Log("hash", h)
}

func TestClient_GenRawTx(t *testing.T) {
	initClient()
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawTransferTx(addrReceive, types.NewCoins(types.NewInt64Coin(client.GetDenom(), 10)))
	assert.NoError(t, err)
	t.Log(string(rawTx))
}

func TestClient_SignRawTx(t *testing.T) {
	initClient()
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawTransferTx(addrReceive, types.NewCoins(types.NewInt64Coin(client.GetDenom(), 10)))
	assert.NoError(t, err)

	signature, err := client.SignMultiSigRawTx(rawTx, "key1")
	assert.NoError(t, err)
	t.Log(string(signature))
}

func TestClient_CreateMultiSigTx(t *testing.T) {
	initClient()
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
	hash, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature2, signature1}, 2)
	assert.NoError(t, err)
	t.Log(hex.EncodeToString(hash))
	t.Log(string(tx))

	txHash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(txHash)
}

func TestClient_BroadcastTx(t *testing.T) {
	initClient()
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawTransferTx(addrReceive, types.NewCoins(types.NewInt64Coin(client.GetDenom(), 10)))
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key1")
	assert.NoError(t, err)

	_, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1}, 2)
	assert.NoError(t, err)

	_, err = client.BroadcastTx(tx)
	assert.ErrorIs(t, err, errors.New(fmt.Sprintf("Boradcast err with res.code: %d", 4)))
}

func TestClient_QueryTxByHash(t *testing.T) {
	initClient()

	for {

		res, err := client.QueryTxByHash("cb21ff0aae81b98f39e1488bb26c098597fab336db45521be7d8aad057a973e3")
		assert.NoError(t, err)
		t.Log(res.Height)
		// for _, msg := range res.GetTx().GetMsgs() {

		// 	t.Log(msg.String())
		// 	t.Log(msg.Type())
		// 	t.Log(msg.Route())
		// }

		curBlock, err := client.GetCurrentBLockHeight()
		assert.NoError(t, err)
		t.Log(curBlock)
		time.Sleep(1 * time.Second)
		t.Log("\n")
	}
}

func TestClient_QueryDelegationRewards(t *testing.T) {
	initClient()
	// r, err := client.QueryDelegations(addrMultiSig1, 0)
	// assert.NoError(t, err)
	// t.Log(r.String())

	rewards, err := client.QueryDelegationTotalRewards(addrMultiSig1, 0)
	assert.NoError(t, err)
	t.Log(rewards.String())

	// res, err := client.QueryDelegationRewards(addrMultiSig1, addrValidator, 0)
	// assert.NoError(t, err)
	// t.Log(res.GetRewards().AmountOf(client.GetDenom()))
}

func TestClient_GenMultiSigRawDelegateTx(t *testing.T) {
	initClient()
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawDelegateTx(addrMultiSig1, []types.ValAddress{adrValidatorEverStake}, types.NewCoin(client.GetDenom(), types.NewInt(1)))
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key1")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature3, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	_, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2, signature3}, 2)
	assert.NoError(t, err)

	hash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log("hash", hash)
}

func TestClient_GenMultiSigRawReDelegateTx(t *testing.T) {
	initClient()
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

	_, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2}, 2)
	assert.NoError(t, err)

	txHash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(txHash)
}

func TestClient_GenMultiSigRawWithdrawDeleRewardTx(t *testing.T) {
	initClient()
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawWithdrawDeleRewardTx(addrMultiSig1, addrValidatorTestnetAteam)
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	h, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2}, 2)
	assert.NoError(t, err)
	t.Log(hex.EncodeToString(h))
	hash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(hash)
}

func TestClient_GenMultiSigRawWithdrawDeleRewardAndDelegratTx(t *testing.T) {
	initClient()
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawWithdrawAllRewardThenDeleTx(addrMultiSig1, 0)
	assert.NoError(t, err)
	t.Log(string(rawTx))

	signature1, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	h, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2}, 2)
	assert.NoError(t, err)
	t.Log(string(tx))
	t.Log(hex.EncodeToString(h))
	hash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(hash)
}

func TestClient_GenMultiSigRawWithdrawAllRewardTx(t *testing.T) {
	initClient()
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawWithdrawAllRewardTx(addrMultiSig1, 0)
	assert.NoError(t, err)
	t.Log(string(rawTx))

	signature1, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	h, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2}, 2)
	assert.NoError(t, err)
	t.Log(string(tx))
	t.Log(hex.EncodeToString(h))
	hash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(hash)
}

// 0xf954ad81a546df9de6c79051ca67f7d0d08e9d861604c76fb8767dce2ce8d4f8
func TestClient_GenMultiSigRawUnDelegateTx(t *testing.T) {
	initClient()
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawUnDelegateTx(addrMultiSig1, []types.ValAddress{adrValidatorEverStake},
		types.NewCoin(client.GetDenom(), types.NewInt(10)))

	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	hash, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2}, 2)
	assert.NoError(t, err)
	t.Log("hash", hex.EncodeToString(hash))
	hash2, err := client.BroadcastTx(tx)
	t.Log("hash2", hash2)
	assert.NoError(t, err)
}

func TestClient_GenMultiSigRawBatchTransferTx(t *testing.T) {
	initClient()
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

	rawTx, err := client.GenMultiSigRawBatchTransferTx(addrMultiSig1, []xBankTypes.Output{out1, out2})
	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "key3")
	assert.NoError(t, err)

	_, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2}, 2)
	assert.NoError(t, err)

	txHash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(txHash)
}

func TestGetPubKey(t *testing.T) {
	initClient()
	test, err := types.AccAddressFromBech32("cosmos1k4mrzs67tnd7fl6ax3l440wpa3amat7teurg9a")
	if err != nil {
		t.Fatal(err)
	}
	account, err := client.QueryAccount(test)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(account)
	t.Log(hex.EncodeToString(account.GetPubKey().Bytes()))

	// res, err := client.QueryTxByHash("327DA2048B6D66BCB27C0F1A6D1E407D88FE719B95A30D108B5906FD6934F7B1")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// msgs := res.GetTx().GetMsgs()
	// for i, _ := range msgs {
	// 	if msgs[i].Type() == xBankTypes.TypeMsgSend {
	// 		msg, _ := msgs[i].(*xBankTypes.MsgSend)
	// 		t.Log(msg.Amount.AmountOf("umuon").Uint64())
	// 	}
	// }

}

func TestClient_Sign(t *testing.T) {
	initClient()
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
	// addrKey1, _ := types.AccAddressFromBech32("cosmos1a8mg9rj4nklhmwkf5vva8dvtgx4ucd9yjasret")
	// addrKey2, _ := types.AccAddressFromBech32("cosmos1ztquzhpkve7szl99jkugq4l8jtpnhln76aetam")
	// addrKey3, _ := types.AccAddressFromBech32("cosmos12zz2hm02sxe9f4pwt7y5q9wjhcu98vnuwmjz4x")
	// addrKey4, _ := types.AccAddressFromBech32("cosmos12yprrdprzat35zhqxe2fcnn3u26gwlt6xcq0pj")
	// addrKey5, _ := types.AccAddressFromBech32("cosmos1em384d8ek3y8nlugapz7p5k5skg58j66je3las")
	// t.Log(hex.EncodeToString(addrKey1.Bytes()))
	// t.Log(hex.EncodeToString(addrKey2.Bytes()))
	// t.Log(hex.EncodeToString(addrKey3.Bytes()))
	// t.Log(hex.EncodeToString(addrKey4.Bytes()))
	// t.Log(hex.EncodeToString(addrKey5.Bytes()))

	addrKey6, _ := types.AccAddressFromBech32("cosmos13jd2vn5wt8h6slj0gcv05lasgpkwpm26n04y75")
	t.Log(hex.EncodeToString(addrKey6.Bytes()))
	addrKey7, _ := types.AccAddressFromBech32("cosmos1cad0efr25faywnjp8qp36l8zlqa2sgz0jwn0hl")
	t.Log(hex.EncodeToString(addrKey7.Bytes()))
	addrKey8, _ := types.AccAddressFromBech32("cosmos13mwxtgrljf9d5r72sc28496ua4lsga0jvmqz8x")
	t.Log(hex.EncodeToString(addrKey8.Bytes()))
	// client_test.go:410: 8c9aa64e8e59efa87e4f4618fa7fb0406ce0ed5a
	// client_test.go:412: c75afca46aa27a474e4138031d7ce2f83aa8204f
	// client_test.go:414: 8edc65a07f924ada0fca86147a975ced7f0475f2
}

func TestClient_QueryDelegations(t *testing.T) {
	initClient()
	res, err := client.QueryDelegations(addrReleaseAddress, 0)
	assert.NoError(t, err)
	t.Log(res.String())
	for i, d := range res.GetDelegationResponses() {
		t.Log(i, d.Balance.Amount.IsZero())
	}
}

func TestClient_QueryBalance(t *testing.T) {
	initClient()
	res, err := client.QueryBalance(addrMultiSig1, "umuon", 440000)
	assert.NoError(t, err)
	t.Log(res.Balance.Amount)
}

func TestClient_QueryDelegationTotalRewards(t *testing.T) {
	initClient()
	res, err := client.QueryDelegationTotalRewards(addrMultiSig2, 0)
	assert.NoError(t, err)
	for i, _ := range res.Rewards {
		t.Log(i, res.Rewards[i].Reward.AmountOf(client.GetDenom()))
		t.Log(i, res.Rewards[i].Reward.AmountOf(client.GetDenom()).TruncateInt())

	}
	t.Log("total ", res.GetTotal().AmountOf(client.GetDenom()).TruncateInt())
}

func TestClient_GetSequence(t *testing.T) {
	initClient()
	seq, err := client.GetSequence(0, addrReleaseAddress)
	assert.NoError(t, err)
	t.Log(seq)
	// txRes,err:=client.QueryTxByHash("FBD05BD4B9DB0386B16E679184EAC88D444B38DA992F8AFC35B5A580B3FC6AA4")
	// assert.NoError(t,err)
	// t.Log(txRes.String())

	res, err := client.QueryUnbondingDelegation(addrReleaseAddress, addrReleaseValWetez, 0)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			t.Log(err)
		}

		t.Fatal(err)
	}
	t.Log(res.String())
	t.Log(len(res.GetUnbond().Entries))
}

func TestMaxTransfer(t *testing.T) {
	initClient()
	err := client.SetFromName("multiSign1")
	assert.NoError(t, err)
	outputs := make([]xBankTypes.Output, 0)
	for i := 0; i < 3; i++ {
		out1 := xBankTypes.Output{
			Address: addrReceive.String(),
			Coins:   types.NewCoins(types.NewCoin(client.GetDenom(), types.NewInt(1))),
		}
		outputs = append(outputs, out1)
	}

	rawTx, err := client.GenMultiSigRawBatchTransferTx(addrMultiSig1, outputs)
	assert.NoError(t, err)

	sequence, err := client.GetSequence(0, addrMultiSig1)
	assert.NoError(t, err)
	signature1, err := client.SignMultiSigRawTxWithSeq(sequence, rawTx, "key2")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTxWithSeq(sequence, rawTx, "key3")
	assert.NoError(t, err)
	t.Log(string(signature1))
	t.Log(string(signature2))

	hash, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2}, 2)
	assert.NoError(t, err)
	t.Log(hex.EncodeToString(hash))
	t.Log(len(tx))
	txHash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log(txHash)
}

func TestMemo(t *testing.T) {
	initClient()
	res, err := client.QueryTxByHash("c7e3f7baf5a5f1d8cbc112080f32070dddd7cca5fe4272e06f8d42c17b25193f")
	assert.NoError(t, err)
	tx, err := client.GetTxConfig().TxDecoder()(res.Tx.GetValue())
	//tx, err := client.GetTxConfig().TxJSONDecoder()(res.Tx.Value)
	assert.NoError(t, err)
	memoTx, ok := tx.(types.TxWithMemo)
	assert.Equal(t, true, ok)
	t.Log(memoTx.GetMemo())
	hb, _ := hexutil.Decode("0xbebd0355ae360c8e6a7ed940a819838c66ca7b8f581f9c0e81dbb5faff346a30")
	//t.Log(string(hb))
	bonderAddr, _ := ss58.Encode(hb, ss58.StafiPrefix)
	t.Log(bonderAddr)
}

func TestMultiThread(t *testing.T) {
	initClient()
	wg := sync.WaitGroup{}
	wg.Add(50)

	for i := 0; i < 50; i++ {
		go func(i int) {
			t.Log(i)
			time.Sleep(5 * time.Second)
			height, err := client.GetAccount()
			if err != nil {
				t.Log("fail", i, err)
			} else {
				t.Log("success", i, height.GetSequence())
			}
			time.Sleep(15 * time.Second)
			height, err = client.GetAccount()
			if err != nil {
				t.Log("fail", i, err)
			} else {
				t.Log("success", i, height.GetSequence())
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestSort(t *testing.T) {
	a := []string{"cosmos1kuyde8vpt8c0ty4pxqgxw3makse7md80umthvg"}
	t.Log(a)
	sort.SliceStable(a, func(i, j int) bool {
		return bytes.Compare([]byte(a[i]), []byte(a[j])) < 0
	})
	t.Log(a)
	// rawTx := "7b22626f6479223a7b226d65737361676573223a5b7b224074797065223a222f636f736d6f732e62616e6b2e763162657461312e4d73674d756c746953656e64222c22696e70757473223a5b7b2261646472657373223a22636f736d6f7331776d6b39797334397a78676d78373770717337636a6e70616d6e6e7875737071753272383779222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a22313134373730227d5d7d5d2c226f757470757473223a5b7b2261646472657373223a22636f736d6f733135366b6b326b71747777776670733836673534377377646c7263326377367163746d36633877222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2231393936227d5d7d2c7b2261646472657373223a22636f736d6f73316b7579646538767074386330747934707871677877336d616b7365376d643830756d74687667222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a223939383030227d5d7d2c7b2261646472657373223a22636f736d6f73316a6b6b68666c753871656471743463796173643674673730676a7778346a6b6872736536727a222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a223132393734227d5d7d5d7d5d2c226d656d6f223a22222c2274696d656f75745f686569676874223a2230222c22657874656e73696f6e5f6f7074696f6e73223a5b5d2c226e6f6e5f637269746963616c5f657874656e73696f6e5f6f7074696f6e73223a5b5d7d2c22617574685f696e666f223a7b227369676e65725f696e666f73223a5b5d2c22666565223a7b22616d6f756e74223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2237353030227d5d2c226761735f6c696d6974223a2231353030303030222c227061796572223a22222c226772616e746572223a22227d7d2c227369676e617475726573223a5b5d7d"
	rawTx := "7b22626f6479223a7b226d65737361676573223a5b7b224074797065223a222f636f736d6f732e62616e6b2e763162657461312e4d73674d756c746953656e64222c22696e70757473223a5b7b2261646472657373223a22636f736d6f7331776d6b39797334397a78676d78373770717337636a6e70616d6e6e7875737071753272383779222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2231303539363338303631227d5d7d5d2c226f757470757473223a5b7b2261646472657373223a22636f736d6f733135366b6b326b71747777776670733836673534377377646c7263326377367163746d36633877222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2231303539363338303631227d5d7d5d7d5d2c226d656d6f223a22222c2274696d656f75745f686569676874223a2230222c22657874656e73696f6e5f6f7074696f6e73223a5b5d2c226e6f6e5f637269746963616c5f657874656e73696f6e5f6f7074696f6e73223a5b5d7d2c22617574685f696e666f223a7b227369676e65725f696e666f73223a5b5d2c22666565223a7b22616d6f756e74223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2237353030227d5d2c226761735f6c696d6974223a2231353030303030222c227061796572223a22222c226772616e746572223a22227d7d2c227369676e617475726573223a5b5d7d"
	// rawTx:="7b22626f6479223a7b226d65737361676573223a5b7b224074797065223a222f636f736d6f732e62616e6b2e763162657461312e4d73674d756c746953656e64222c22696e70757473223a5b7b2261646472657373223a22636f736d6f7331776d6b39797334397a78676d78373770717337636a6e70616d6e6e7875737071753272383779222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a22313134373730227d5d7d5d2c226f757470757473223a5b7b2261646472657373223a22636f736d6f73316a6b6b68666c753871656471743463796173643674673730676a7778346a6b6872736536727a222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a223132393734227d5d7d2c7b2261646472657373223a22636f736d6f733135366b6b326b71747777776670733836673534377377646c7263326377367163746d36633877222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2231393936227d5d7d2c7b2261646472657373223a22636f736d6f73316b7579646538767074386330747934707871677877336d616b7365376d643830756d74687667222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a223939383030227d5d7d5d7d5d2c226d656d6f223a22222c2274696d656f75745f686569676874223a2230222c22657874656e73696f6e5f6f7074696f6e73223a5b5d2c226e6f6e5f637269746963616c5f657874656e73696f6e5f6f7074696f6e73223a5b5d7d2c22617574685f696e666f223a7b227369676e65725f696e666f73223a5b5d2c22666565223a7b22616d6f756e74223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2237353030227d5d2c226761735f6c696d6974223a2231353030303030222c227061796572223a22222c226772616e746572223a22227d7d2c227369676e617475726573223a5b5d7d"
	txBts, err := hex.DecodeString(rawTx)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(txBts))
}
