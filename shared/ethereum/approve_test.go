package ethereum

import (
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/rtoken-relay/bindings/Multisig"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

var (
	owners = []common.Address{common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6")}

	goerliMultisigContract  = common.HexToAddress("0x1Cb8b55cB11152E74D34Be1961E4FFe169F5B99A")
	goerliMultisigContract1 = common.HexToAddress("0xfc42de640aa9759d460e1a11416eca95d25c5908")
)

func TestMultisigApprove(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	keys := make([]*secp256k1.Keypair, 0)
	for _, owner := range owners {
		kpI, err := keystore.KeypairFromAddress(owner.Hex(), keystore.EthChain, keystorePath, false)
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

	multi, err := Multisig.NewMultisig(goerliMultisigContract, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	cd, _ := mabi.Pack("approve", goerliStakeManagerContract, big.NewInt(0).Mul(&config.AmountBase, big.NewInt(100000000000000000)))
	msg, err := multi.MessageToSign(nil, goerliMaticToken, defaultValue, cd)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hexutil.Encode(msg[:]))

	signatures := make([][]byte, 0)
	for _, key := range keys {
		signature, err := crypto.Sign(msg[:], key.PrivateKey())
		if err != nil {
			t.Fatal(err)
		}
		t.Log(hexutil.Encode(signature))
		t.Log(len(signature))

		signatures = append(signatures, signature)
	}

	vs, rs, ss := utils.DecomposeSignature(signatures)
	err = client.LockAndUpdateOpts()
	if err != nil {
		t.Fatal(err)
	}

	txhash := common.HexToHash("0x11d668ca5c97408167f046131a37b4ef10ccbd621d83f920eefddaa62fe77e0c")
	tx, err := multi.ExecTransaction(
		client.Opts(),
		goerliMaticToken,
		defaultValue,
		cd,
		uint8(config.Call),
		big.NewInt(100000),
		txhash,
		vs,
		rs,
		ss,
	)

	client.UnlockOpts()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(tx.Hash())
}
