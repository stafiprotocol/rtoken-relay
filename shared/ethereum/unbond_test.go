package ethereum

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/rtoken-relay/bindings/OldMultisig"
	"github.com/stafiprotocol/rtoken-relay/bindings/StakeManager"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

var (
	unbondPool       = common.HexToAddress("0xb91f931ebeb626126b50ae2e9ce8cf7496497d98")
	unbondPoolOwners = []common.Address{common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6"), common.HexToAddress("0xBd39f5936969828eD9315220659cD11129071814")}

	unbondAmount, _ = utils.StringToBigint("")
	unbondTxHash    = common.HexToHash("0x9bd668ca5c97508167f046131a37b4ef10ccbd621dabf920eeeddbb62fe77e1d")
)

func TestMultisigUnbond(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)
	keys := make([]*secp256k1.Keypair, 0)

	for _, own := range unbondPoolOwners {
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

	cd, err := ValidatorShareAbi.Pack("sellVoucher_new", unbondAmount, unbondAmount)
	manager, err := StakeManager.NewStakeManager(goerliStakeManagerContract, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	share, err := manager.Validators(nil, big.NewInt(9))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("share", share.ContractAddress)

	mt := &ethmodel.MultiTransaction{
		To:        share.ContractAddress,
		Value:     big.NewInt(0),
		CallData:  cd,
		Operation: ethmodel.Call,
		SafeTxGas: big.NewInt(400000),
		TotalGas:  big.NewInt(400000),
	}

	multi, err := OldMultisig.NewMultisig(unbondPool, client.Client())
	if err != nil {
		t.Fatal(err)
	}
	_ = multi

	msg, err := multi.MessageToSign(nil, mt.To, mt.Value, mt.CallData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("msg", hexutil.Encode(msg[:]))

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

	tx, err := multi.ExecTransaction(
		client.Opts(),
		mt.To,
		mt.Value,
		mt.CallData,
		uint8(mt.Operation),
		mt.SafeTxGas,
		unbondTxHash,
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
