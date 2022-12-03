package ethereum

import (
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/rtoken-relay/bindings/MaticToken"
	"github.com/stafiprotocol/rtoken-relay/bindings/Multisig"
	"github.com/stafiprotocol/rtoken-relay/bindings/ValidatorShare"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

var (
	// update these vars before running TestMultisigApprove
	sender                      = common.HexToAddress("0x76397f6e8899793D48EC316d2492De6919a4a32E")
	goerliMultisigProxyContract = common.HexToAddress("0xe74340EAfD516e8DD2e0B528dA7D73089e156c1E")
	keystorePath                = "/Users/tpkeeper/gowork/stafi/rtoken-relay/keys"
	proxyOwners                 = []common.Address{
		common.HexToAddress("0xA96577dA157b173618bd7420005a14F73cb0e294"),
		common.HexToAddress("0x4D50094712505B484e5115F0987281b82969b9F5"),
	}
	txhash = common.HexToHash("0x8bd668ca5c97508167f046131a37b4ef10ccbd621dabf920eefddaa62fe77e1d")

	goerliEndPoint             = "wss://goerli.infura.io/ws/v3/86f8d5ba0d524274bce7780a83dbc0a4"
	testLogger                 = newTestLogger("test")
	goerliStakeManagerContract = common.HexToAddress("0x00200eA4Ee292E253E6Ca07dBA5EdC07c8Aa37A3")
	goerliMaticToken           = common.HexToAddress("0x499d11e0b6eac7c0593d8fb292dcbbf815fb29ae")
	maticTokenAbi, _           = abi.JSON(strings.NewReader(MaticToken.MaticTokenABI))
	multisigAbi, _             = abi.JSON(strings.NewReader(Multisig.MultisigABI))
	ValidatorShareAbi, _       = abi.JSON(strings.NewReader(ValidatorShare.ValidatorShareABI))
	AliceKp                    = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
	AmountBase                 = big.NewInt(1e18)
)

func TestMultisigApprove(t *testing.T) {
	password := "tpkeeper"
	os.Setenv(keystore.EnvPassword, password)
	keys := make([]*secp256k1.Keypair, 0)

	for _, own := range proxyOwners {
		kpI, err := keystore.KeypairFromAddress(own.Hex(), keystore.EthChain, keystorePath, false)
		if err != nil {
			t.Fatal(err)
		}
		kp, _ := kpI.(*secp256k1.Keypair)

		keys = append(keys, kp)
	}
	kpI, err := keystore.KeypairFromAddress(sender.Hex(), keystore.EthChain, keystorePath, false)
	if err != nil {
		t.Fatal(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)
	sender := kp

	client := NewClient(goerliEndPoint, sender, testLogger, big.NewInt(0), big.NewInt(0))
	err = client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	// approve
	cd, _ := maticTokenAbi.Pack("approve", goerliStakeManagerContract, big.NewInt(0).Mul(AmountBase, big.NewInt(1e17)))
	mt := &ethmodel.MultiTransaction{
		To:        goerliMaticToken,
		Value:     big.NewInt(0),
		CallData:  cd,
		Operation: ethmodel.Call,
		SafeTxGas: big.NewInt(100000),
		TotalGas:  big.NewInt(100000),
	}

	// sigs
	msg := mt.MessageToSign(txhash, goerliMultisigProxyContract)
	t.Log("msg", hexutil.Encode(msg[:]))

	multi, err := Multisig.NewMultisig(goerliMultisigProxyContract, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	sigs := make([][]byte, 0)
	for _, key := range keys {
		signature, err := crypto.Sign(msg[:], key.PrivateKey())
		if err != nil {
			t.Fatal(err)
		}

		sigs = append(sigs, signature)
	}
	vs, rs, ss := utils.DecomposeSignature(sigs)

	err = client.LockAndUpdateOpts(mt.TotalGas, big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	// send
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

func newTestLogger(name string) core.Logger {
	tLog := core.NewLog("chain", name)
	return tLog
}
