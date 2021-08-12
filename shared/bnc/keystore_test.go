package bnc

import (
	"bytes"
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	subtypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
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
	reason, err := client.BnbTransferVerify(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, submodel.Pass, reason)

	wrongBh, _ := hexutil.Decode("0x7144512dbd193241ca6c518526599c760df161a620a696116b01cde815cd8349")
	r.Blockhash = wrongBh

	reason, err = client.BnbTransferVerify(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, submodel.BlockhashUnmatch, reason)
	r.Blockhash = bh

	wrongTh, _ := hexutil.Decode("0x49c61f50df6cf784a5021d9e4b532a83c9e728e6c3824b08376a2a1a6941a602")
	r.Txhash = wrongTh
	reason, err = client.BnbTransferVerify(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, submodel.PubkeyUnmatch, reason)
	r.Txhash = th

	wrongFrom := common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6").Bytes()
	r.Pubkey = wrongFrom
	reason, err = client.BnbTransferVerify(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, submodel.PubkeyUnmatch, reason)
	r.Pubkey = pk

	wrongTo := common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6").Bytes()
	r.Pool = wrongTo
	reason, err = client.BnbTransferVerify(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, submodel.PoolUnmatch, reason)
	r.Pool = pool

	wrongAmt := subtypes.NewU128(*big.NewInt(28e16))
	r.Amount = wrongAmt
	reason, err = client.BnbTransferVerify(r)
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

func newTestLogger(name string) log15.Logger {
	tLog := log15.New("chain", name)
	tLog.SetHandler(log15.LvlFilterHandler(log15.LvlError, tLog.GetHandler()))
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
