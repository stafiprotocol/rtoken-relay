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
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

var (
	transferPool       = common.HexToAddress("0xb91f931ebeb626126b50ae2e9ce8cf7496497d98")
	transferPoolOwners = []common.Address{common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6"), common.HexToAddress("0xBd39f5936969828eD9315220659cD11129071814")}
	transferAmount     = big.NewInt(0).Mul(AmountBase, big.NewInt(42)) // 表示42个matic
	transferDst        = common.HexToHash("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6")
	transferTxHash     = common.HexToHash("0x9bd668ca5c97508178f123131a37b4ab10ccbd621dabf920eeeddbb62fe77e1d")
)

func TestMultisigTransfer(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)
	keys := make([]*secp256k1.Keypair, 0)

	for _, own := range transferPoolOwners {
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

	cd, err := maticTokenAbi.Pack("transfer", transferDst, transferAmount)
	mt := &ethmodel.MultiTransaction{
		To:        goerliMaticToken,
		Value:     big.NewInt(0),
		CallData:  cd,
		Operation: ethmodel.Call,
		SafeTxGas: big.NewInt(400000),
		TotalGas:  big.NewInt(400000),
	}

	multi, err := OldMultisig.NewMultisig(transferPool, client.Client())
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
		transferTxHash,
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
