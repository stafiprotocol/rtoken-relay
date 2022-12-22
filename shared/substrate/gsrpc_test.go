package substrate

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/sr25519"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

var (
	AliceKey     = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
	From         = "31yavGB5CVb8EwpqKQaS9XY7JZcfbK6QpWPn5kkweHVpqcov"
	LessPolka    = "1334v66HrtqQndbugYxX9m56V6222m97LbavB4KAMmqgjsas"
	From1        = "31d96Cq9idWQqPq3Ch5BFY84zrThVE3r98M7vG4xYaSWHwsX"
	From2        = "1TgYb5x8xjsZRyL5bwvxUoAWBn36psr4viSMHbRXA8bkB2h"
	Wen          = "1swvN162p1siDjm63UhhWoa59bpPZTSNKGVcbCwHUYkfRRW"
	Jun          = "33RQ73d9XfPTaE2SV7dzdhQQ17YaeMQ4yzhzAQhhFVenxMuJ"
	relay1       = "33xzQzUk75cAxt7i3hHb2XWwJNFqzcSULfoCRsAkiCG4jh5d"
	KeystorePath = "/Users/tpkeeper/gowork/stafi/rtoken-relay/keys"
)

func TestBatchTransfer(t *testing.T) {
	stop := make(chan int)
	sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, AddressTypeAccountId, AliceKey, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	less, _ := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
	jun, _ := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
	wen, _ := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	bao, _ := types.NewAddressFromHexAccountID("0x9c4189297ad2140c85861f64656d1d1318994599130d98b75ff094176d2ca31e")

	addrs := []types.Address{less, jun, wen, bao}

	amount, _ := utils.StringToBigint("3000" + "000000000000")
	value := types.NewUCompact(amount)
	calls := make([]types.Call, 0)

	ci, err := sc.FindCallIndex(config.MethodTransferKeepAlive)
	if err != nil {
		t.Fatal(err)
	}

	for _, addr := range addrs {
		call, err := types.NewCallWithCallIndex(
			ci,
			config.MethodTransferKeepAlive,
			addr,
			value,
		)
		if err != nil {
			t.Fatal(err)
		}
		calls = append(calls, call)
	}

	ext, err := sc.NewUnsignedExtrinsic(config.MethodBatch, calls)
	if err != nil {
		t.Fatal(err)
	}

	err = sc.SignAndSubmitTx(ext)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetConst(t *testing.T) {
	stop := make(chan int)
	sc, err := NewSarpcClient(ChainTypePolkadot, "wss://kusama-rpc.polkadot.io", polkaTypesFile, AddressTypeMultiAddress, AliceKey, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	e, err := sc.ExistentialDeposit()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(e)
}

func TestPolkaQueryStorage(t *testing.T) {
	stop := make(chan int)
	sc, err := NewSarpcClient(ChainTypePolkadot, "wss://kusama-rpc.polkadot.io", polkaTypesFile, AddressTypeMultiAddress, AliceKey, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	var index uint32
	exist, err := sc.QueryStorage(config.StakingModuleId, config.StorageActiveEra, nil, nil, &index)
	if err != nil {
		panic(err)
	}

	if !exist {
		panic("not exist")
	}

	t.Log(index)
}

func TestStafiLocalQueryActiveEra(t *testing.T) {
	stop := make(chan int)
	sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, AddressTypeAccountId, AliceKey, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	var index uint32
	exist, err := sc.QueryStorage(config.StakingModuleId, config.StorageActiveEra, nil, nil, &index)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
	t.Log("activeEra", index)
}

func TestActive(t *testing.T) {
	stop := make(chan int)
	//sc, err := NewSarpcClient(ChainTypePolkadot, "wss://kusama-test-rpc.stafi.io", polkaTypesFile, AddressTypeMultiAddress, AliceKey, tlog, stop)
	sc, err := NewSarpcClient(ChainTypePolkadot, "wss://kusama-test-rpc.stafi.io", polkaTypesFile, AddressTypeMultiAddress, AliceKey, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	a := "0xac0df419ce0dc61b092a5cfa06a28e40cd82bc9de7e8c1e5591169360d66ba3c"
	mac, err := types.NewMultiAddressFromHexAccountID(a)
	if err != nil {
		t.Fatal(err)
	}
	ledger := new(submodel.StakingLedger)
	exist, err := sc.QueryStorage(config.StakingModuleId, config.StorageLedger, mac.AsID[:], nil, ledger)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
	t.Log("ledger", ledger)

	t.Log(types.NewU128(big.Int(ledger.Active)))
}

func TestActive1(t *testing.T) {
	stop := make(chan int)
	sc, err := NewSarpcClient(ChainTypePolkadot, "wss://polkadot-test-rpc.stafi.io", polkaTypesFile, AddressTypeMultiAddress, AliceKey, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	a := "0x782a467d4ff23b660ca5f1ecf47f8537d4c35049541b6ebbf5381c00c4c158f7"
	b, _ := hexutil.Decode(a) // work
	//mac, err := types.NewAddressFromHexAccountID(a) // work
	//mac, err := types.NewMultiAddressFromHexAccountID(a) // work
	ledger := new(submodel.StakingLedger)
	exist, err := sc.QueryStorage(config.StakingModuleId, config.StorageLedger, b, nil, ledger)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
	t.Log(types.NewU128(big.Int(ledger.Active)))
}

func TestHash(t *testing.T) {
	h, _ := types.NewHashFromHexString("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	a, _ := types.EncodeToBytes(h)

	fmt.Println(hexutil.Encode(a))
}

func TestPool(t *testing.T) {
	p := "0x782a467d4ff23b660ca5f1ecf47f8537d4c35049541b6ebbf5381c00c4c158f7"
	pool, _ := hexutil.Decode(p)
	pbz, _ := types.EncodeToBytes(pool)
	fmt.Println(pool)
	fmt.Println(pbz)

	//stop := make(chan int)
	//gc, err := NewGsrpcClient("wss://stafi-seiya.stafi.io", AddressTypeAccountId, AliceKey, tlog, stop)
	//assert.NoError(t, err)
	//
	//
	////poolBz, err := types.EncodeToBytes(pool)
	//symBz, err := types.EncodeToBytes(core.RKSM)
	//assert.NoError(t, err)
	//
	//var threshold uint16
	//exist, err := gc.QueryStorage(config.RTokenLedgerModuleId, config.StorageMultiThresholds, symBz, pbz, &threshold)
	//assert.NoError(t, err)
	//fmt.Println(exist)
	//fmt.Println()

}

func Test_KSM_GsrpcClient_Multisig(t *testing.T) {

	logrus.SetLevel(logrus.TraceLevel)

	password := "tpkeeper"
	os.Setenv(keystore.EnvPassword, password)

	kp, err := keystore.KeypairFromAddress(relay1, keystore.SubChain, KeystorePath, false)
	if err != nil {
		t.Fatal(err)
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()
	stop := make(chan int)
	sc, err := NewSarpcClient(ChainTypePolkadot, "wss://kusama-test-rpc.stafi.io", kusamaTypesFile, AddressTypeMultiAddress, krp, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	//pool, err := hexutil.Decode("ac0df419ce0dc61b092a5cfa06a28e40cd82bc9de7e8c1e5591169360d66ba3c")
	//assert.NoError(t, err)

	threshold := uint16(2)
	//wen, _ := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	// jun, _ := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
	relay2, _ := types.NewMultiAddressFromHexAccountID("0x2afeb305f32a12507a6b211d818218577b0e425692766b08b8bc5d714fccac3b")

	others := []types.AccountID{
		relay2.AsID,
	}

	//for _, oth := range others {
	//	fmt.Println(hexutil.Encode(oth[:]))
	//}

	bond, _ := utils.StringToBigint("1000000000000")
	unbond := big.NewInt(0)

	call, err := sc.BondOrUnbondCall(bond, unbond)
	if err != nil {
		t.Fatal(err)
	}

	h := utils.BlakeTwo256(call.Opaque)
	t.Log("Extrinsic", call.Extrinsic)
	t.Log("Opaque", hexutil.Encode(call.Opaque))
	t.Log("callHash", hexutil.Encode(h[:]))

	info, err := sc.GetPaymentQueryInfo(call.Extrinsic)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("info", info.Class, info.PartialFee, info.Weight)

	//optp := types.TimePoint{Height: 1964877, Index: 1}
	//tp := submodel.NewOptionTimePoint(optp)

	tp := submodel.NewOptionTimePointEmpty()
	ext, err := sc.NewUnsignedExtrinsic(config.MethodAsMulti, threshold, others, tp, call.Opaque, false, info.Weight)
	if err != nil {
		t.Fatal(err)
	}

	err = sc.SignAndSubmitTx(ext)
	if err != nil {
		t.Fatal(err)
	}
}
func Test_KSM_GsrpcClient_transfer(t *testing.T) {

	logrus.SetLevel(logrus.TraceLevel)

	password := "tpkeeper"
	os.Setenv(keystore.EnvPassword, password)

	kp, err := keystore.KeypairFromAddress(relay1, keystore.SubChain, KeystorePath, false)
	if err != nil {
		t.Fatal(err)
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()
	stop := make(chan int)
	// sc, err := NewSarpcClient(ChainTypePolkadot, "wss://kusama-test-rpc.stafi.io", kusamaTypesFile, AddressTypeAccountId, krp, tlog, stop)
	sc, err := NewSarpcClient(ChainTypePolkadot, "ws://127.0.0.1:9944", kusamaTypesFile, AddressTypeMultiAddress, krp, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	//pool, err := hexutil.Decode("ac0df419ce0dc61b092a5cfa06a28e40cd82bc9de7e8c1e5591169360d66ba3c")
	//assert.NoError(t, err)

	// threshold := uint16(2)
	//wen, _ := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	// jun, _ := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
	relay2, _ := types.NewMultiAddressFromHexAccountID("0x2afeb305f32a12507a6b211d818218577b0e425692766b08b8bc5d714fccac3b")

	// others := []types.AccountID{
	// 	relay2.AsAccountID,
	// }

	//for _, oth := range others {
	//	fmt.Println(hexutil.Encode(oth[:]))
	//}

	// bond, _ := utils.StringToBigint("1000000000000")
	// unbond := big.NewInt(0)

	// call,err:=sc.TransferCall(relay2.AsAccountID[:],types.NewUCompact(big.NewInt(1000000)))
	// if err!=nil{
	// 	t.Fatal(err)
	// }
	ext, err := sc.NewUnsignedExtrinsic(config.MethodTransfer, relay2, types.NewUCompact(big.NewInt(1e10)))
	if err != nil {
		t.Fatal(err)
	}

	// call, err := sc.BondOrUnbondCall(bond, unbond)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// h := utils.BlakeTwo256(call.Opaque)
	// t.Log("Extrinsic", call.Extrinsic)
	// t.Log("Opaque", hexutil.Encode(call.Opaque))
	// t.Log("callHash", hexutil.Encode(h[:]))

	// info, err := sc.GetPaymentQueryInfo(call.Extrinsic)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log("info", info.Class, info.PartialFee, info.Weight)

	//optp := types.TimePoint{Height: 1964877, Index: 1}
	//tp := submodel.NewOptionTimePoint(optp)

	// tp := submodel.NewOptionTimePointEmpty()
	// ext, err := sc.NewUnsignedExtrinsic(config.MethodAsMulti, threshold, others, tp, call.Opaque, false, info.Weight)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	err = sc.SignAndSubmitTx(ext)
	if err != nil {
		t.Fatal(err)
	}
}
