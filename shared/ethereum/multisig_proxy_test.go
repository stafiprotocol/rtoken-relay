package ethereum

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/rtoken-relay/bindings/MaticToken"
	"github.com/stafiprotocol/rtoken-relay/bindings/Multisig"
	"github.com/stafiprotocol/rtoken-relay/bindings/ProxyAdmin"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stretchr/testify/assert"
)

var (
	goerliMultisigProxy = common.HexToAddress("0x5e0f61b16b1a5F7f557FCde20a0075bb52FFCd33")
	goerliProxyAdmin    = common.HexToAddress("0xBf996BFe7a62ab39130281eC1062eDbEC88B708d")

	goerliRFis = common.HexToAddress("0xb89f77f28c2d84d32ed8bb2311afa29a6089559c")
)

func TestMultisigProxyAdmin(t *testing.T) {
	client := NewGoerliClient()
	proxyAdmin, err := ProxyAdmin.NewProxyAdmin(goerliProxyAdmin, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	owner, err := proxyAdmin.Owner(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("owner", owner)

	admin, err := proxyAdmin.GetProxyAdmin(nil, goerliMultisigProxy)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("admin", admin)

	impl, err := proxyAdmin.GetProxyImplementation(nil, goerliMultisigProxy)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("impl", impl)
}

func TestMultisigProxyStorage(t *testing.T) {
	txhash := common.HexToHash("0x8ed668ca5c97408167f046131a37b4ef10ccbd621dabf920eefddaa62fe77e1b")
	client := NewGoerliClient()

	multisig, err := Multisig.NewMultisig(goerliMultisigProxy, client.Client())
	if err != nil {
		t.Fatal(err)
	}
	state, err := multisig.TxHashs(nil, txhash)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, state, uint8(ethmodel.HashStateSuccess))
}

//func TestMultisigProxyUpgrade(t *testing.T) {
//	password := "123456"
//	os.Setenv(keystore.EnvPassword, password)
//
//	admin := common.HexToAddress("0xBd39f5936969828eD9315220659cD11129071814")
//	kpI, err := keystore.KeypairFromAddress(admin.Hex(), keystore.EthChain, keystorePath, false)
//	if err != nil {
//		panic(err)
//	}
//	kp, _ := kpI.(*secp256k1.Keypair)
//
//	client := NewClient(goerliEndPoint, kp, testLogger, big.NewInt(0), big.NewInt(0))
//	err = client.Connect()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	proxyAdmin, err := ProxyAdmin.NewProxyAdmin(goerliProxyAdmin, client.Client())
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	newImpl := common.HexToAddress("0x33b520a2eB16e08C4110564b70D0F5e57a1Bcf04")
//
//	err = client.LockAndUpdateOpts(big.NewInt(0))
//	if err != nil {
//		t.Fatal(err)
//	}
//	tx, err := proxyAdmin.Upgrade(client.Opts(), goerliMultisigProxy, newImpl)
//	client.UnlockOpts()
//
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("Upgrade txHash", tx.Hash())
//}

//func TestMultisigProxyInitialize(t *testing.T) {
//	password := "123456"
//	os.Setenv(keystore.EnvPassword, password)
//
//	kpI, err := keystore.KeypairFromAddress(owner.Hex(), keystore.EthChain, keystorePath, false)
//	if err != nil {
//		panic(err)
//	}
//	kp, _ := kpI.(*secp256k1.Keypair)
//
//	client := NewClient(goerliEndPoint, kp, testLogger, big.NewInt(0), big.NewInt(0))
//	err = client.Connect()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	multisig, err := Multisig.NewMultisig(goerliMultisigProxy, client.Client())
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	th, err := multisig.GetThreshold(nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("threshold", th)
//
//	//owners := []common.Address{common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6"), common.HexToAddress("0xBd39f5936969828eD9315220659cD11129071814")}
//	//
//	//err = client.LockAndUpdateOpts()
//	//if err != nil {
//	//	t.Fatal(err)
//	//}
//	//tx, err := multisig.Initialize(client.opts, owners, big.NewInt(2))
//	//client.UnlockOpts()
//	//if err != nil {
//	//	t.Fatal(err)
//	//}
//	//t.Log("initialize txHash", tx.Hash())
//	//
//	//th, err = multisig.GetThreshold(nil)
//	//if err != nil {
//	//	t.Fatal(err)
//	//}
//	//t.Log("threshold after initialize", th)
//}

func TestErc20Balance(t *testing.T) {
	client := NewGoerliClient()
	rfis, err := MaticToken.NewMaticToken(goerliRFis, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	bal, err := rfis.BalanceOf(nil, goerliMultisigProxy)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("bal", bal)
}

//func TestProxyTransferErc20(t *testing.T) {
//	password := "123456"
//	os.Setenv(keystore.EnvPassword, password)
//
//	owners := []common.Address{common.HexToAddress("0xBd39f5936969828eD9315220659cD11129071814"), common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6")}
//	keys := make([]*secp256k1.Keypair, 0)
//
//	for _, own := range owners {
//		kpI, err := keystore.KeypairFromAddress(own.Hex(), keystore.EthChain, keystorePath, false)
//		if err != nil {
//			panic(err)
//		}
//		kp, _ := kpI.(*secp256k1.Keypair)
//
//		keys = append(keys, kp)
//	}
//
//	client := NewClient(goerliEndPoint, keys[0], testLogger, big.NewInt(0), big.NewInt(0))
//	err := client.Connect()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	cd, _ := maticTokenAbi.Pack("transfer", receiver1, big.NewInt(0).Mul(big.NewInt(10000000000000000), big.NewInt(10)))
//	mt := &ethmodel.MultiTransaction{
//		To:        goerliRFis,
//		Value:     big.NewInt(0),
//		CallData:  cd,
//		Operation: ethmodel.Call,
//		SafeTxGas: big.NewInt(100000),
//		TotalGas: big.NewInt(100000),
//	}
//	txhash := common.HexToHash("0x1ed668ca5c97408167f046131a37b4ef10ccbd621dabf920eefddaa62fe77e1b")
//	msg := mt.MessageToSign(txhash, goerliMultisigProxy)
//	t.Log("msg", hexutil.Encode(msg[:]))
//
//	multi, err := Multisig.NewMultisig(goerliMultisigProxy, client.Client())
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	sigs := make([][]byte, 0)
//	for _, key := range keys {
//		signature, err := crypto.Sign(msg[:], key.PrivateKey())
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		sigs = append(sigs, signature)
//	}
//
//	vs, rs, ss := utils.DecomposeSignature(sigs)
//	err = client.LockAndUpdateOpts(mt.TotalGas)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	tx, err := multi.ExecTransaction(
//		client.Opts(),
//		mt.To,
//		mt.Value,
//		mt.CallData,
//		uint8(mt.Operation),
//		mt.SafeTxGas,
//		txhash,
//		vs,
//		rs,
//		ss,
//	)
//
//	client.UnlockOpts()
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	t.Log("txHash", tx.Hash())
//}
