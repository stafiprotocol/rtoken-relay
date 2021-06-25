// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"github.com/stafiprotocol/rtoken-relay/bindings/StakeManager"
	"github.com/stafiprotocol/rtoken-relay/bindings/ValidatorShare"
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
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/bindings/MaticToken"
	"github.com/stafiprotocol/rtoken-relay/bindings/MultiSend"
	"github.com/stafiprotocol/rtoken-relay/bindings/Multisig"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stretchr/testify/assert"
)

var (
	goerliEndPoint          = "wss://goerli.infura.io/ws/v3/86f8d5ba0d524274bce7780a83dbc0a4"
	goerliMultisendContract = common.HexToAddress("0x747e29a783a9EE438bD25ac32bB341f12c827217")
	goerliErc20Token        = common.HexToAddress("0x7c338c09fcdb43db9877032d06eea43a254c6a28")

	owner        = common.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6")
	receiver1    = common.HexToAddress("0xaD0bf51f7fc89e262edBbdF53C260088B024D857")
	receiver2    = common.HexToAddress("0x1Bf32E717FfeD95c5629bd9628e6F11E380e096B")
	defaultValue = big.NewInt(0)

	testLogger   = newTestLogger("test")
	keystorePath = "/Users/fwj/Go/stafi/rtoken-relay/keys/ethereum/"

	mabi, _    = abi.JSON(strings.NewReader(MaticToken.MaticTokenABI))
	sendAbi, _ = abi.JSON(strings.NewReader(MultiSend.MultiSendABI))
)

func TestTransferCallData(t *testing.T) {
	cd, err := mabi.Pack("transfer", receiver1, big.NewInt(0).Mul(big.NewInt(1000000000000000000), big.NewInt(10)))
	assert.NoError(t, err)
	t.Log(hexutil.Encode(cd))
	// 0xa9059cbb000000000000000000000000ad0bf51f7fc89e262edbbdf53c260088b024d8570000000000000000000000000000000000000000000000008ac7230489e80000
}

func TestMultisigTransfer(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner.Hex(), keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client := NewClient(goerliEndPoint, kp, testLogger, big.NewInt(0), big.NewInt(0))
	err = client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	multi, err := Multisig.NewMultisig(goerliMultisigContrat, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	cd, _ := mabi.Pack("transfer", receiver1, big.NewInt(0).Mul(big.NewInt(1000000000000000000), big.NewInt(10)))
	msg, err := multi.MessageToSign(nil, goerliErc20Token, defaultValue, cd)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hexutil.Encode(msg[:]))

	signature, err := crypto.Sign(msg[:], client.Keypair().PrivateKey())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hexutil.Encode(signature))
	t.Log(len(signature))

	var rs [32]byte
	copy(rs[:], signature[:32])
	var ss [32]byte
	copy(ss[:], signature[32:64])

	rses := [][32]byte{rs}
	sses := [][32]byte{ss}
	vses := []uint8{signature[64:][0]}

	err = client.LockAndUpdateOpts()
	if err != nil {
		t.Fatal(err)
	}

	txhash := common.HexToHash("0x8ed668ca5c97408167f046131a37b4ef10ccbd621d83f920eefddaa62fe77e0c")
	tx, err := multi.ExecTransaction(
		client.Opts(),
		goerliErc20Token,
		defaultValue,
		cd,
		config.Call,
		big.NewInt(100000),
		txhash,
		vses,
		rses,
		sses,
	)

	client.UnlockOpts()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(tx.Hash())
}

func TestMakeUpTransactions(t *testing.T) {
	amount := big.NewInt(0).Mul(big.NewInt(1000000000000000000), big.NewInt(500))

	x := common.LeftPadBytes(amount.Bytes(), 32)
	t.Log(hexutil.Encode(x))
}

func TestMultisigSend(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner.Hex(), keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client := NewClient(goerliEndPoint, kp, testLogger, big.NewInt(0), big.NewInt(0))
	err = client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	multi, err := Multisig.NewMultisig(goerliMultisigContrat, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	cd1, _ := mabi.Pack("transfer", receiver1, big.NewInt(0).Mul(big.NewInt(1000000000000000000), big.NewInt(10)))
	cd2, _ := mabi.Pack("transfer", receiver2, big.NewInt(0).Mul(big.NewInt(1000000000000000000), big.NewInt(500)))

	bts := ethmodel.BatchTransactions{
		&ethmodel.BatchTransaction{
			Operation:  config.Call,
			To:         goerliErc20Token,
			Value:      defaultValue,
			DataLength: big.NewInt(int64(len(cd1))),
			Data:       cd1,
		},
		&ethmodel.BatchTransaction{
			Operation:  config.Call,
			To:         goerliErc20Token,
			Value:      defaultValue,
			DataLength: big.NewInt(int64(len(cd2))),
			Data:       cd2,
		},
	}

	txs := bts.Encode()
	cd3, err := sendAbi.Pack("multiSend", txs)
	assert.NoError(t, err)

	msg, err := multi.MessageToSign(nil, goerliMultisendContract, defaultValue, cd3)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hexutil.Encode(msg[:]))

	signature, err := crypto.Sign(msg[:], client.Keypair().PrivateKey())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hexutil.Encode(signature))
	t.Log(len(signature))

	var rs [32]byte
	copy(rs[:], signature[:32])
	var ss [32]byte
	copy(ss[:], signature[32:64])

	rses := [][32]byte{rs}
	sses := [][32]byte{ss}
	vses := []uint8{signature[64:][0]}

	err = client.LockAndUpdateOpts()
	if err != nil {
		t.Fatal(err)
	}

	txhash := common.HexToHash("0x8ed668ca5c97408167f046131a37b4ef10ccbd621d83f920eefddaa62fe77e0c")
	tx, err := multi.ExecTransaction(
		client.Opts(),
		goerliMultisendContract,
		defaultValue,
		cd3,
		uint8(1),
		big.NewInt(200000),
		txhash,
		vses,
		rses,
		sses,
	)

	client.UnlockOpts()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(tx.Hash())
}

func newTestLogger(name string) log15.Logger {
	tLog := log15.New("chain", name)
	tLog.SetHandler(log15.LvlFilterHandler(log15.LvlError, tLog.GetHandler()))
	return tLog
}

func TestTransferPack(t *testing.T) {
	value := big.NewInt(0).Mul(big.NewInt(1000000000000000000), big.NewInt(10))
	a := types.NewUCompact(value)
	b := big.Int(a)
	cd1, err := mabi.Pack("transfer", receiver1, &b)

	assert.NoError(t, err)
	t.Log(hexutil.Encode(cd1))
	t.Log(crypto.Keccak256Hash(cd1).Hex())
}

func TestVerify(t *testing.T) {
	bh, _ := hexutil.Decode("0x2e61dec15e7b3fcd19af31603de13e44f6fa8a4df311c981408a24e5e0bf02b4")
	th, _ := hexutil.Decode("0xa7ab84dac6dae5700a3ddcfb6fda62d5117f39809e99b3737cc01250d75c0660")

	pk, _ := hexutil.Decode("0xbd39f5936969828ed9315220659cd11129071814")
	pool, _ := hexutil.Decode("0x37c9c42eedbc72842cc48f0e51006ac804987e38")
	amt := big.NewInt(0).Mul(big.NewInt(1000000000000000000), big.NewInt(500))

	a := &submodel.BondRecord{
		Pubkey:    pk,
		Pool:      pool,
		Blockhash: bh,
		Txhash:    th,
		Amount:    types.NewU128(*amt),
	}

	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner.Hex(), keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client := NewClient(goerliEndPoint, kp, testLogger, big.NewInt(0), big.NewInt(0))
	err = client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	reason, err := client.TransferVerify(a, goerliErc20Token)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(reason)

	wrongBh, _ := hexutil.Decode("0x64c375983dbf3f1680c252684695a17a8f58d7f84ce71e406bebd1d9de67304f")
	a.Blockhash = wrongBh

	reason, err = client.TransferVerify(a, goerliErc20Token)
	assert.NoError(t, err)
	assert.Equal(t, submodel.BlockhashUnmatch, reason)
	a.Blockhash = bh

	wrongTh, _ := hexutil.Decode("0x165bc1fc1cea7d0f6df6fc33fa0e838a5dc15bb460603f170de384b85afc878a")
	a.Txhash = wrongTh
	reason, err = client.TransferVerify(a, goerliErc20Token)
	assert.NoError(t, err)
	assert.Equal(t, submodel.BlockhashUnmatch, reason)
	a.Txhash = th

	a.Pubkey = pool
	reason, err = client.TransferVerify(a, goerliErc20Token)
	assert.NoError(t, err)
	assert.Equal(t, submodel.PubkeyUnmatch, reason)
	a.Pubkey = pk

	a.Pool = pk
	reason, err = client.TransferVerify(a, goerliErc20Token)
	assert.NoError(t, err)
	assert.Equal(t, submodel.PoolUnmatch, reason)
	a.Pool = pool

	wrongAmt := big.NewInt(0).Mul(big.NewInt(1000000000000000000), big.NewInt(50))
	a.Amount = types.NewU128(*wrongAmt)
	reason, err = client.TransferVerify(a, goerliErc20Token)
	assert.NoError(t, err)
	assert.Equal(t, submodel.AmountUnmatch, reason)
}

func TestApprove(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner.Hex(), keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client := NewClient(goerliEndPoint, kp, testLogger, big.NewInt(0), big.NewInt(0))
	err = client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	token, err := MaticToken.NewMaticToken(goerliMaticToken, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	err = client.LockAndUpdateOpts()
	if err != nil {
		t.Fatal(err)
	}

	approveTx, err := token.Approve(client.Opts(), goerliStakeManagerContract, big.NewInt(0).Mul(&config.AmountBase, big.NewInt(1000000000000)))
	if err != nil {
		t.Fatal(err)
	}

	client.UnlockOpts()

	if err != nil {
		t.Fatal(err)
	}

	t.Log("approveTxHahs", approveTx.Hash())
}

func TestBond(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner.Hex(), keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client := NewClient(goerliEndPoint, kp, testLogger, big.NewInt(0), big.NewInt(0))
	err = client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	manager, err := StakeManager.NewStakeManager(goerliStakeManagerContract, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	id := big.NewInt(9)
	shareData, err := manager.Validators(nil, id)
	if err != nil {
		t.Fatal(err)
	}

	shareContract := shareData.ContractAddress
	t.Log("ShareContract", shareContract)

	valFalg, err := manager.IsValidator(nil, id)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("valFalg", valFalg)

	share, err := ValidatorShare.NewValidatorShare(shareData.ContractAddress, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	err = client.LockAndUpdateOpts()
	if err != nil {
		t.Fatal(err)
	}

	tx, err := share.BuyVoucher(client.Opts(), &config.AmountBase, big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}

	client.UnlockOpts()

	if err != nil {
		t.Fatal(err)
	}

	t.Log("txHash", tx.Hash())
}

func TestRestake(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner.Hex(), keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client := NewClient(goerliEndPoint, kp, testLogger, big.NewInt(0), big.NewInt(0))
	err = client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	manager, err := StakeManager.NewStakeManager(goerliStakeManagerContract, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	id := big.NewInt(9)
	shareData, err := manager.Validators(nil, id)
	if err != nil {
		t.Fatal(err)
	}

	shareContract := shareData.ContractAddress
	t.Log("ShareContract", shareContract)

	valFalg, err := manager.IsValidator(nil, id)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("valFalg", valFalg)

	share, err := ValidatorShare.NewValidatorShare(shareData.ContractAddress, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	stake, rate, err := share.GetTotalStake(nil, owner)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("rate", rate)
	t.Log("stake", stake)

	reward, err := share.GetLiquidRewards(nil, owner)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("reward", reward)

	nonce, err := share.UnbondNonces(nil, owner)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("nonce", nonce)

	unbond, err := share.UnbondsNew(nil, owner, nonce)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("unbond", unbond)

	//err = client.LockAndUpdateOpts()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	////tx, err := share.BuyVoucher(client.Opts(), &config.AmountBase, big.NewInt(0))
	//tx, err := share.Restake(client.Opts())
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//client.UnlockOpts()
	//
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Log("txHash", tx.Hash())
}

func TestMultisigApprove(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner.Hex(), keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client := NewClient(goerliEndPoint, kp, testLogger, big.NewInt(0), big.NewInt(0))
	err = client.Connect()
	if err != nil {
		t.Fatal(err)
	}

	multi, err := Multisig.NewMultisig(goerliMultisigContrat, client.Client())
	if err != nil {
		t.Fatal(err)
	}

	cd, _ := mabi.Pack("approve", goerliStakeManagerContract, big.NewInt(0).Mul(&config.AmountBase, big.NewInt(100000000000000000)))
	msg, err := multi.MessageToSign(nil, goerliMaticToken, defaultValue, cd)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hexutil.Encode(msg[:]))

	signature, err := crypto.Sign(msg[:], client.Keypair().PrivateKey())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hexutil.Encode(signature))
	t.Log(len(signature))

	var rs [32]byte
	copy(rs[:], signature[:32])
	var ss [32]byte
	copy(ss[:], signature[32:64])

	rses := [][32]byte{rs}
	sses := [][32]byte{ss}
	vses := []uint8{signature[64:][0]}

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
		config.Call,
		big.NewInt(100000),
		txhash,
		vses,
		rses,
		sses,
	)

	client.UnlockOpts()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(tx.Hash())
}

func TestValidatorIdToHex(t *testing.T) {
	id := big.NewInt(9)
	t.Log("idHex", hexutil.Encode(id.Bytes())) //0x09
}

func TestKecca(t *testing.T) {
	a := []byte(`unclaimable`)

	t.Log(crypto.Keccak256Hash(a))
}
