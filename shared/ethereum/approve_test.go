package ethereum

import (
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/rtoken-relay/bindings/MaticToken"
	"github.com/stafiprotocol/rtoken-relay/bindings/MultiSend"
	"github.com/stafiprotocol/rtoken-relay/bindings/Multisig"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

var (
	owners = []common.Address{common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6"), common.HexToAddress("0xBd39f5936969828eD9315220659cD11129071814")}

	/// update these vars before running TestMultisigApprove
	goerliMultisigProxyContract = common.HexToAddress("")
	goerliEndPoint              = "wss://goerli.infura.io/ws/v3/86f8d5ba0d524274bce7780a83dbc0a4"
	keystorePath                = "/Users/fwj/Go/stafi/rtoken-relay/keys/ethereum/"
	proxyOwners                 = []common.Address{common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6")}
	txhash                      = common.HexToHash("0x8bd668ca5c97508167f046131a37b4ef10ccbd621dabf920eefddaa62fe77e1d")

	testLogger                 = newTestLogger("test")
	goerliStakeManagerContract = common.HexToAddress("0x00200eA4Ee292E253E6Ca07dBA5EdC07c8Aa37A3")
	goerliMaticToken           = common.HexToAddress("0x499d11e0b6eac7c0593d8fb292dcbbf815fb29ae")
	maticTokenAbi, _           = abi.JSON(strings.NewReader(MaticToken.MaticTokenABI))
	sendAbi, _                 = abi.JSON(strings.NewReader(MultiSend.MultiSendABI))
	multisigAbi, _             = abi.JSON(strings.NewReader(Multisig.MultisigABI))
	AliceKp                    = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
	AmountBase                 = big.NewInt(1e18)
)

func TestMultisigApprove(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)
	keys := make([]*secp256k1.Keypair, 0)

	for _, own := range proxyOwners {
		kpI, err := keystore.KeypairFromAddress(own.Hex(), keystore.EthChain, keystorePath, false)
		if err != nil {
			panic(err)
		}
		kp, _ := kpI.(*secp256k1.Keypair)

		keys = append(keys, kp)
	}

	client := NewClient(goerliEndPoint, keys[0], testLogger, big.NewInt(0), big.NewInt(0))
	err := client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	cd, _ := maticTokenAbi.Pack("approve", goerliStakeManagerContract, big.NewInt(0).Mul(AmountBase, big.NewInt(1e17)))
	mt := &ethmodel.MultiTransaction{
		To:        goerliMaticToken,
		Value:     big.NewInt(0),
		CallData:  cd,
		Operation: ethmodel.Call,
		SafeTxGas: big.NewInt(100000),
		TotalGas:  big.NewInt(100000),
	}

	msg := mt.MessageToSign(txhash, goerliMultisigProxyContract)
	t.Log("msg", hexutil.Encode(msg[:]))

	multi, err := Multisig.NewMultisig(goerliMultisigProxyContract, client.Client())
	if err != nil {
		t.Fatal(err)
	}
	_ = multi

	sigs := make([][]byte, 0)
	for _, key := range keys {
		signature, err := crypto.Sign(msg[:], key.PrivateKey())
		if err != nil {
			t.Fatal(err)
		}

		sigs = append(sigs, signature)
	}
	vs, rs, ss := utils.DecomposeSignature(sigs)

	err = client.LockAndUpdateOpts(mt.TotalGas)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := multi.ExecTransaction(
		client.Opts(),
		mt.To,
		mt.Value,
		mt.CallData,
		uint8(mt.Operation),
		mt.SafeTxGas,
		txhash,
		vs,
		rs,
		ss,
	)

	client.UnlockOpts()

	if err != nil {
		t.Fatal(err)
	}

	t.Log("txHash", tx.Hash())
}

func TestMultisigProxyInitial(t *testing.T) {
	threshold := big.NewInt(1)
	cd, err := multisigAbi.Pack("initialize", proxyOwners, threshold)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("initializeData", hexutil.Encode(cd))
}

func newTestLogger(name string) log15.Logger {
	tLog := log15.New("chain", name)
	tLog.SetHandler(log15.LvlFilterHandler(log15.LvlError, tLog.GetHandler()))
	return tLog
}
