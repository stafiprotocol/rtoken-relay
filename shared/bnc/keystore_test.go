package bnc

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stafiprotocol/go-sdk/keys"
	subtypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/bnc"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	bncRpc "github.com/stafiprotocol/go-sdk/client/rpc"
	bncCmnTypes "github.com/stafiprotocol/go-sdk/common/types"
	bncTypes "github.com/stafiprotocol/go-sdk/types"
	"github.com/stafiprotocol/rtoken-relay/bindings/TokenHub"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
)

var (
	tokenHubContract = common.HexToAddress("0x0000000000000000000000000000000000001004")
	bscTestEndpoint  = "https://data-seed-prebsc-2-s3.binance.org:8545"
	zeroAddress      = common.HexToAddress("0x0000000000000000000000000000000000000000")
	relayfee         = big.NewInt(1e16)

	keystorePath = "/Users/fwj/Go/stafi/rtoken-relay/keys/ethereum/"
	testLogger   = newTestLogger("test")
	AliceKp      = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
)

//func TestLoadKeyStore(t *testing.T) {
//	keystore := "/Users/fwj/Go/stafi/rtoken-relay/keys/"
//
//	addr := "tbnb12h62xhn5klt6xug4e9hjhz3dkwav78wz8u72ha"
//
//	file := keystore + addr + ".key"
//
//
//}

func TestVerifyBNBTransfer(t *testing.T) {
	client := newBscTestClient()

	block, err := client.Client().BlockByNumber(context.Background(), big.NewInt(10408048))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("blockNum", block.Number())

	t.Log("blockHash", block.Hash())

	txHash := common.HexToHash("0xe172e9cc8178d7ac3f9b216def7014df4d2c230d98aea61bf8a1b4c26b6b3954")
	for _, tx := range block.Transactions() {
		if !bytes.Equal(txHash.Bytes(), tx.Hash().Bytes()) {
			continue
		}

		t.Log("data", hexutil.Encode(tx.Data()))
		t.Log("to", tx.To())
		t.Log("type", tx.Type())
		t.Log("chainId", tx.ChainId())

		msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), big.NewInt(0))
		if err != nil {
			t.Fatal(err)
		}

		t.Log("from", msg.From())
		t.Log("value", msg.Value())
	}
}

func TestBnbTransferVerify(t *testing.T) {
	bh := common.HexToHash("0x4ffa036ca7a2308b3be5e0cf21c8a90a1ea9374e3f5ea2a444372efd7e11c908").Bytes()
	th := common.HexToHash("0xe172e9cc8178d7ac3f9b216def7014df4d2c230d98aea61bf8a1b4c26b6b3954").Bytes()
	pk := common.HexToAddress("0xad0bf51f7fc89e262edbbdf53c260088b024d857").Bytes()
	pool := common.HexToAddress("0xb7fda60ab83a35479e9ee373a43a861ba7aa0006").Bytes()

	r := &submodel.BondRecord{
		Pubkey:    pk,
		Pool:      pool,
		Blockhash: bh,
		Txhash:    th,
		Amount:    subtypes.NewU128(*big.NewInt(28e17)),
	}

	client := newBscTestClient()
	reason, err := client.TransferVerifyNative(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, submodel.Pass, reason)

	wrongBh, _ := hexutil.Decode("0x7144512dbd193241ca6c518526599c760df161a620a696116b01cde815cd8349")
	r.Blockhash = wrongBh

	reason, err = client.TransferVerifyNative(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, submodel.BlockhashUnmatch, reason)
	r.Blockhash = bh

	wrongTh, _ := hexutil.Decode("0x49c61f50df6cf784a5021d9e4b532a83c9e728e6c3824b08376a2a1a6941a602")
	r.Txhash = wrongTh
	reason, err = client.TransferVerifyNative(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, submodel.PubkeyUnmatch, reason)
	r.Txhash = th

	wrongFrom := common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6").Bytes()
	r.Pubkey = wrongFrom
	reason, err = client.TransferVerifyNative(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, submodel.PubkeyUnmatch, reason)
	r.Pubkey = pk

	wrongTo := common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6").Bytes()
	r.Pool = wrongTo
	reason, err = client.TransferVerifyNative(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, submodel.PoolUnmatch, reason)
	r.Pool = pool

	wrongAmt := subtypes.NewU128(*big.NewInt(28e16))
	r.Amount = wrongAmt
	reason, err = client.TransferVerifyNative(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, submodel.AmountUnmatch, reason)
}

func TestTransferOut(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	owner := common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6")
	kpI, err := keystore.KeypairFromAddress(owner.Hex(), keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client := ethereum.NewClient(bscTestEndpoint, kp, testLogger, big.NewInt(0), big.NewInt(0))
	err = client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	bal, err := client.Client().BalanceAt(context.Background(), owner, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("bal", bal)

	hub, err := TokenHub.NewTokenHub(tokenHubContract, client.Client())
	if err != nil {
		t.Fatal(err)
	}
	_ = hub

	receiver := common.HexToAddress("0x5acf525eccbe80a8dad05ad208e9bc94c89bab1f")
	amount := big.NewInt(0).Add(big.NewInt(1e10), big.NewInt(1e8))
	value := big.NewInt(0).Add(amount, relayfee)

	err = client.LockAndUpdateOpts(big.NewInt(0), value)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := hub.TransferOut(client.Opts(), zeroAddress, receiver, amount, 0x17b1f307761)
	client.UnlockOpts()

	if err != nil {
		t.Fatal(err)
	}
	t.Log("txHash", tx.Hash())

	time.Sleep(2 * time.Minute)

	bal, err = client.Client().BalanceAt(context.Background(), owner, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("bal-after", bal)
}

func TestBatchTransferOut(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	owner := common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6")
	kpI, err := keystore.KeypairFromAddress(owner.Hex(), keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client := ethereum.NewClient(bscTestEndpoint, kp, testLogger, big.NewInt(0), big.NewInt(0))
	err = client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	hub, err := TokenHub.NewTokenHub(tokenHubContract, client.Client())
	if err != nil {
		t.Fatal(err)
	}
	_ = hub

	receiver := common.HexToAddress("0x5acf525eccbe80a8dad05ad208e9bc94c89bab1f")
	amount := big.NewInt(5e17)
	value := big.NewInt(0).Add(amount, relayfee)

	err = client.LockAndUpdateOpts(big.NewInt(0), value)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := hub.TransferOut(client.Opts(), zeroAddress, receiver, amount, 0x17b1f307761)
	client.UnlockOpts()

	if err != nil {
		t.Fatal(err)
	}
	t.Log("txHash", tx.Hash())
}

func TestAmountAndExpireTime(t *testing.T) {
	a, _ := hexutil.Decode("0x120b3d02e0f68000")
	b := big.NewInt(0).SetBytes(a)
	t.Log(b) // 1300200000000000000

	c, _ := hexutil.Decode("0x610e1984")
	d := big.NewInt(0).SetBytes(c)
	t.Log(d) //
}

func newTestLogger(name string) core.Logger {
	tLog := core.NewLog("chain", name)
	return tLog
}

func newBscTestClient() *ethereum.Client {
	client := ethereum.NewClient(bscTestEndpoint, AliceKp, testLogger, big.NewInt(0), big.NewInt(0))
	err := client.Connect()
	if err != nil {
		panic(err)
	}

	return client
}

func TestStakingReward(t *testing.T) {
	url := "https://testnet-api.binance.org/v1/staking/chains/chapel/delegators/tbnb1tt84yhkvh6q23kksttfq36dujnyfh2cldrzux5/rewards?limit=10&offset=0"

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	sr := new(bnc.StakingReward)
	if err := json.Unmarshal(body, sr); err != nil {
		t.Fatal(err)
	}

	t.Log(sr)

	url1 := "https://testnet-api.binance.org/v1/staking/chains/chapel/delegators/tbnb1ufrtxk7f5w0skl5evusrmsd6cundpxvmpz4n4n/rewards?limit=10&offset=0"
	sr1, err := bnc.GetStakingReward(url1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sr1.Total)
}

func TestValidatorId(t *testing.T) {
	addr, err := bncCmnTypes.ValAddressFromBech32("bva16kujlngdxq4pvyf87gpzx2x7ya4lgsz96j0aqt")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(addr)

	a := []byte(`bva15mgzha93ny878kuvjl0pnqmjygwccdad08uecu`)
	b := hexutil.Encode(a)
	t.Log(b) // 0x62766131366b756a6c6e676478713470767966383767707a783278377961346c67737a39366a30617174
}

func TestRpcClient(t *testing.T) {
	rpcEndpoint := "tcp://data-seed-pre-1-s3.binance.org:80"

	keyManager, err := keys.NewPrivateKeyManager("64967ded205b00b1f872f59242031d4cc02a1bcca47017361d7f3854e86c545e")
	if err != nil {
		t.Fatal(err)
	}

	client := bncRpc.NewRPCClient(rpcEndpoint, bncCmnTypes.TestNetwork)
	t.Log("IsActive", client.IsActive())
	bal, err := client.GetBalance(keyManager.GetAddr(), "BNB")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("bal", bal.Free)

	dels, err := client.QuerySideChainDelegations(bncTypes.ChapelNet, keyManager.GetAddr())
	if err != nil {
		t.Fatal(err)
	}

	for _, del := range dels {
		t.Log(del)
	}

	addr, err := bncCmnTypes.AccAddressFromBech32("tbnb1ufrtxk7f5w0skl5evusrmsd6cundpxvmpz4n4n")
	if err != nil {
		t.Fatal(err)
	}

	unbonds, err := client.QuerySideChainUnbondingDelegations(bncTypes.ChapelNet, addr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(unbonds)
}

func TestMainnetRpcClient(t *testing.T) {
	rpcEndpoint := "tcp://dataseed3.ninicoin.io:80"

	addr, err := bncCmnTypes.AccAddressFromBech32("bnb17eak3g00hn5mvvlz6f0dv09p83z7cu5ym7cesl")
	if err != nil {
		t.Fatal(err)
	}

	client := bncRpc.NewRPCClient(rpcEndpoint, bncCmnTypes.ProdNetwork)
	t.Log("IsActive", client.IsActive())

	bal, err := client.GetBalance(addr, "BNB")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("bal", bal.Free)
}

func TestIsactive(t *testing.T) {
	client := bncRpc.NewRPCClient("tcp://data-seed-pre-1-s3.binance.org:80", bncCmnTypes.TestNetwork)
	t.Log("IsActive", client.IsActive())

	addr, err := bncCmnTypes.AccAddressFromBech32("tbnb1tt84yhkvh6q23kksttfq36dujnyfh2cldrzux5")
	if err != nil {
		t.Fatal(err)
	}

	bal, err := client.GetBalance(addr, "BNB")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("bal", bal.Free)
}

func TestTxHashStatus(t *testing.T) {
	txHash := common.HexToHash("0x35976b21bfca498cc14241fff3ebdbaea4565d216f8778a0f0a229c091bce871")

	client := newBscTestClient()
	receipt, err := client.TransactionReceipt(txHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("status", receipt.Status)

	anotherTxHash := common.HexToHash("0x65765caa0099feba25d973e47ec4376aa6367d401a94d1419907128669467b0b")
	receipt1, err := client.TransactionReceipt(anotherTxHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("status1", receipt1.Status)
}

func TestBcTxhash(t *testing.T) {
	//url := "https://testnet-dex.binance.org/api/v1/tx/AF1C536217B7B797622B9414B373A0E217EB3665F429720DF10A6A545BECA718"
	url1 := "https://testnet-dex.binance.org/api/v1/tx/A95B078EDE38E89FAC1819A019B7B135D3E0350A25F840636E9C73304EF8CE89"

	resp, err := http.Get(url1)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(body))

	thr := new(bnc.TxHashResult)
	if err := json.Unmarshal(body, thr); err != nil {
		t.Fatal(err)
	}

	t.Log(thr)
}

func TestTxHashStatus1(t *testing.T) {
	txHash := common.HexToHash("0x4c112c9fb5b1861579840b6eaf9229955c9de2737990addf586e779be3ebbf2e")

	client := newBscTestClient()
	receipt, err := client.TransactionReceipt(txHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("status", receipt.Status)
}

func TestBalance(t *testing.T) {
	bncCmnTypes.Network = bncCmnTypes.TestNetwork
	addr, err := bncCmnTypes.AccAddressFromBech32("tbnb1v5umyy5snft4xep474gt55k9gvan7dvwt4m90n")
	if err != nil {
		t.Fatal(err)
	}

	client := bncRpc.NewRPCClient("tcp://data-seed-pre-1-s3.binance.org:80", bncCmnTypes.TestNetwork)
	bal, err := client.GetBalance(addr, "BNB")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("bal", bal.Free)
}
