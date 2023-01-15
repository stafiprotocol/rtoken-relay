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
	relay2       = "31iZqreRXv9Hc1CBSh9TnmJD8uW9DiaPUjeBcnofKNpeSyzv"
	KeystorePath = "/Users/tpkeeper/gowork/stafi/rtoken-relay/keys"
)

//func TestGsrpcClient(t *testing.T) {
//	password := "123456"
//	os.Setenv(keystore.EnvPassword, password)
//
//	kp, err := keystore.KeypairFromAddress(From, keystore.SubChain, KeystorePath, false)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	krp := kp.(*sr25519.Keypair).AsKeyringPair()
//
//	stop := make(chan int)
//	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, krp, tlog, stop)
//	assert.NoError(t, err)
//
//	bId, err := types.NewHashFromHexString("0xd2f787195c3498f941653fe542e62b397988eeb3e2b867cf39a93c1ae5127b41")
//	assert.NoError(t, err)
//
//	symbz, _ := types.EncodeToBytes(core.RDOT)
//
//	br := new(submodel.BondRecord)
//	exist, err := gc.QueryStorage(config.RTokenSeriesModuleId, config.StorageBondRecords, symbz, bId[:], br)
//	assert.NoError(t, err)
//	fmt.Println("exist:", exist)
//
//	meta, err := gc.GetLatestMetadata()
//	assert.NoError(t, err)
//
//	call, err := types.NewCall(
//		meta,
//		config.MethodExecuteBondRecord,
//		core.RDOT,
//		bId,
//		submodel.Pass,
//	)
//	assert.NoError(t, err)
//
//	//prop := &conn.Proposal{call, bondKey}
//	//symbol: RSymbol, prop_id: T::Hash, in_favour: bool
//	ext, err := gc.NewUnsignedExtrinsic(config.MethodRacknowledgeProposal, core.RDOT, bId, true, call)
//	err = gc.SignAndSubmitTx(ext)
//	assert.NoError(t, err)
//}
//
//func TestGsrpcClient1(t *testing.T) {
//	password := "123456"
//	os.Setenv(keystore.EnvPassword, password)
//
//	kp, err := keystore.KeypairFromAddress(From, keystore.SubChain, KeystorePath, false)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	krp := kp.(*sr25519.Keypair).AsKeyringPair()
//	stop := make(chan int)
//	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, krp, tlog, stop)
//	assert.NoError(t, err)
//
//	bId, err := types.NewHashFromHexString("0x7942f736614444159dba03bc5fef684dcebb7c5edb16d4e57315e148ea4525be")
//	assert.NoError(t, err)
//
//	//prop := &conn.Proposal{call, bondKey}
//	//symbol: RSymbol, prop_id: T::Hash, in_favour: bool
//	ext, err := gc.NewUnsignedExtrinsic("RTokenSeries.add_hashes", core.RDOT, bId, submodel.Pass)
//	err = gc.SignAndSubmitTx(ext)
//	assert.NoError(t, err)
//}
//
//func TestGsrpcClient_StakingActive(t *testing.T) {
//	password := "123456"
//	os.Setenv(keystore.EnvPassword, password)
//
//	kp, err := keystore.KeypairFromAddress(From, keystore.SubChain, KeystorePath, false)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	krp := kp.(*sr25519.Keypair).AsKeyringPair()
//	stop := make(chan int)
//	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, krp, tlog, stop)
//	assert.NoError(t, err)
//
//	b, _ := hexutil.Decode("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
//
//	jun := types.NewAddressFromAccountID(b)
//	s, err := gc.StakingLedger(jun.AsAccountID)
//	fmt.Println(s.Active)
//}
//
//func TestGsrpcClient_QueryStorage(t *testing.T) {
//	stop := make(chan int)
//	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, AliceKey, tlog, stop)
//	assert.NoError(t, err)
//
//	pool, err := hexutil.Decode("0xbeb93b63149fd6a1b69f2d50492fe01983be085cbf09f689219838f6322763d8")
//	assert.NoError(t, err)
//
//	symbz, _ := types.EncodeToBytes(core.RDOT)
//
//	subAcs := make([]types.Bytes, 0)
//	exist, err := gc.QueryStorage(config.RTokenLedgerModuleId, config.StorageSubAccounts, symbz, pool, &subAcs)
//	assert.NoError(t, err)
//	assert.True(t, exist)
//	for _, ac := range subAcs {
//		fmt.Println(hexutil.Encode(ac))
//	}
//}

/*func TestGsrpcClient_Multisig(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kp, err := keystore.KeypairFromAddress(From, keystore.SubChain, KeystorePath, false)
	if err != nil {
		t.Fatal(err)
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()
	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, krp, tlog, stop)
	assert.NoError(t, err)

	pool, err := hexutil.Decode("0xbeb93b63149fd6a1b69f2d50492fe01983be085cbf09f689219838f6322763d8")
	assert.NoError(t, err)

	symbz, _ := types.EncodeToBytes(core.RDOT)

	var threshold uint16
	exist, err := gc.QueryStorage(config.RTokenLedgerModuleId, config.StorageMultiThresholds, symbz, pool, &threshold)
	assert.NoError(t, err)
	assert.True(t, exist)
	fmt.Println(threshold)

	subs := make([]types.Bytes, 0)
	exist, err = gc.QueryStorage(config.RTokenLedgerModuleId, config.StorageSubAccounts, symbz, pool, &subs)
	assert.NoError(t, err)
	assert.True(t, exist)
	fmt.Println(subs)

	others := make([]types.AccountID, 0)
	for i, ac := range subs {
		if hexutil.Encode(gc.PublicKey()) == hexutil.Encode(ac) {
			bzs := append(subs[:i], subs[i+1:]...)
			for _, bz := range bzs {
				others = append(others, types.NewAccountID(bz))
			}
		}
	}

	for _, oth := range others {
		fmt.Println(hexutil.Encode(oth[:]))
	}

	bond, _ := utils.StringToBigint("10000000000000")
	unbond := big.NewInt(0)

	call, err := gc.BondOrUnbondCall(bond, unbond)
	assert.NoError(t, err)
	fmt.Println(call.Extrinsic)
	fmt.Println(hexutil.Encode(call.Opaque))
	h := utils.BlakeTwo256(call.Opaque)
	fmt.Println("callHash", hexutil.Encode(h[:]))

	sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, tlog)
	assert.NoError(t, err)

	info, err := sc.GetPaymentQueryInfo(call.Extrinsic)
	assert.NoError(t, err)
	fmt.Println("info", info.Class, info.PartialFee, info.Weight)

	tp := submodel.NewOptionTimePointEmpty()
	ext, err := gc.NewUnsignedExtrinsic(config.MethodAsMulti, threshold, others, tp, call.Opaque, false, info.Weight)
	err = gc.SignAndSubmitTx(ext)
	assert.NoError(t, err)
}*/

//	func TestStorages(t *testing.T) {
//		stop := make(chan int)
//		gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, AliceKey, tlog, stop)
//		assert.NoError(t, err)
//
//		addr, err := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
//		assert.NoError(t, err)
//
//		targets := make([]types.AccountID, 0)
//		exist, err := gc.QueryStorage(config.StakingModuleId, config.StorageNominators, addr.AsAccountID[:], nil, &targets)
//		assert.NoError(t, err)
//		fmt.Println(exist)
//		//fmt.Println(targets)
//		for _, t := range targets {
//			fmt.Println(hexutil.Encode(t[:]))
//		}
//
//		era := uint32(5)
//		bz, err := types.EncodeToBytes(era)
//		assert.NoError(t, err)
//
//		re := new(submodel.EraRewardPoints)
//		exist1, err := gc.QueryStorage(config.StakingModuleId, config.StorageErasRewardPoints, bz, nil, re)
//		assert.NoError(t, err)
//		fmt.Println(exist1)
//		fmt.Println(re)
//
//		ledger := new(submodel.StakingLedger)
//		exist, err = gc.QueryStorage(config.StakingModuleId, config.StorageLedger, addr.AsAccountID[:], nil, ledger)
//		assert.NoError(t, err)
//		fmt.Println(exist)
//		fmt.Println(ledger)
//	}
//
//	func TestBatch(t *testing.T) {
//		password := "123456"
//		os.Setenv(keystore.EnvPassword, password)
//
//		kp, err := keystore.KeypairFromAddress(From, keystore.SubChain, KeystorePath, false)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		krp := kp.(*sr25519.Keypair).AsKeyringPair()
//		stop := make(chan int)
//		gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, krp, tlog, stop)
//		assert.NoError(t, err)
//
//		addr, err := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
//		assert.NoError(t, err)
//
//		amount, _ := utils.StringToBigint("1000000000000")
//		value := types.NewUCompact(amount)
//
//		calls := make([]types.Call, 0)
//		meta, err := gc.GetLatestMetadata()
//		assert.NoError(t, err)
//
//		for i := 0; i <= 4000; i++ {
//			call, err := types.NewCall(
//				meta,
//				config.MethodTransferKeepAlive,
//				addr,
//				value,
//			)
//			assert.NoError(t, err)
//			calls = append(calls, call)
//		}
//
//		ext, err := gc.NewUnsignedExtrinsic(config.MethodBatch, calls)
//		assert.NoError(t, err)
//
//		err = gc.SignAndSubmitTx(ext)
//		assert.NoError(t, err)
//	}
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

//func TestBatchTransfer1(t *testing.T) {
//	stop := make(chan int)
//	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, AliceKey, tlog, stop)
//	assert.NoError(t, err)
//
//	less, _ := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
//	jun, _ := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
//	wen, _ := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
//	bao, _ := types.NewAddressFromHexAccountID("0x9c4189297ad2140c85861f64656d1d1318994599130d98b75ff094176d2ca31e")
//
//	addrs := []types.Address{less, jun, wen, bao}
//	receives := make([]*submodel.Receive, 0)
//	amount, _ := utils.StringToBigint("3000" + "000000000000")
//	value := types.NewUCompact(amount)
//	for _, addr := range addrs {
//		func(addr types.Address) {
//			rec := &submodel.Receive{Recipient: addr.AsAccountID[:], Value: value}
//			receives = append(receives, rec)
//		}(addr)
//	}
//	for _, rec := range receives {
//		t.Log(*rec)
//	}
//
//	err = gc.BatchTransfer(receives)
//	assert.NoError(t, err)
//}

//func TestMultiBatch(t *testing.T) {
//	password := "123456"
//	os.Setenv(keystore.EnvPassword, password)
//
//	kp, err := keystore.KeypairFromAddress(From, keystore.SubChain, KeystorePath, false)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	krp := kp.(*sr25519.Keypair).AsKeyringPair()
//	stop := make(chan int)
//	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, krp, tlog, stop)
//	assert.NoError(t, err)
//
//	//less, _ := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
//	jun, _ := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
//	wen, _ := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
//	bao, _ := types.NewAddressFromHexAccountID("0x9c4189297ad2140c85861f64656d1d1318994599130d98b75ff094176d2ca31e")
//	//mul, _ := types.NewAddressFromHexAccountID("0xbeb93b63149fd6a1b69f2d50492fe01983be085cbf09f689219838f6322763d8")
//	threshold := uint16(2)
//	others := []types.AccountID{wen.AsAccountID, jun.AsAccountID}
//	addrs := []types.Address{bao, jun, bao}
//
//	amount, _ := utils.StringToBigint("1000000000000")
//	value := types.NewUCompact(amount)
//
//	calls := make([]types.Call, 0)
//
//	var info *rpc.PaymentQueryInfo
//	tp := submodel.NewOptionTimePointEmpty()
//	for _, addr := range addrs {
//		c, err := gc.TransferCall(addr.AsAccountID[:], value)
//		fmt.Println(c.CallHash)
//		assert.NoError(t, err)
//		if info == nil {
//			sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, tlog)
//			assert.NoError(t, err)
//			info, err = sc.GetPaymentQueryInfo(c.Extrinsic)
//			assert.NoError(t, err)
//		}
//
//		ext, err := gc.NewUnsignedExtrinsic(config.MethodAsMulti, threshold, others, tp, c.Opaque, false, info.Weight)
//		xt, ok := ext.(*types.Extrinsic)
//		assert.True(t, ok)
//		calls = append(calls, xt.Method)
//	}
//
//	ext, err := gc.NewUnsignedExtrinsic(config.MethodBatch, calls)
//	assert.NoError(t, err)
//
//	err = gc.SignAndSubmitTx(ext)
//	assert.NoError(t, err)
//}

//	func TestBalance(t *testing.T) {
//		stop := make(chan int)
//		gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, AliceKey, tlog, stop)
//		assert.NoError(t, err)
//
//		less, _ := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
//		ac := new(types.AccountInfo)
//		exist, err := gc.QueryStorage(config.SystemModuleId, config.StorageAccount, less.AsAccountID[:], nil, ac)
//		assert.NoError(t, err)
//
//		fmt.Println(exist)
//		fmt.Println(ac.Data.Free)
//	}
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

//func TestGsrpcClient_Multisig1(t *testing.T) {
//	password := "123456"
//	os.Setenv(keystore.EnvPassword, password)
//
//	kp, err := keystore.KeypairFromAddress(Wen, keystore.SubChain, KeystorePath, false)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	krp := kp.(*sr25519.Keypair).AsKeyringPair()
//	stop := make(chan int)
//	gc, err := NewGsrpcClient("wss://kusama-test-rpc.stafi.io", AddressTypeMultiAddress, krp, tlog, stop)
//	assert.NoError(t, err)
//
//	//pool, err := hexutil.Decode("0xac0df419ce0dc61b092a5cfa06a28e40cd82bc9de7e8c1e5591169360d66ba3c")
//	//assert.NoError(t, err)
//
//	wen, _ := types.NewMultiAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
//	//less, _ := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
//
//	//bond, _ := utils.StringToBigint("10000000000000")
//	//unbond := big.NewInt(0)
//	//
//	//call, err := gc.BondOrUnbondCall(bond, unbond)
//
//	amount, _ := utils.StringToBigint("10000000000000")
//	call, err := gc.TransferCall(wen.AsID[:], types.NewUCompact(amount))
//	assert.NoError(t, err)
//	fmt.Println(call)
//
//	fmt.Println("callHash", call.Extrinsic)
//
//	sc, err := NewSarpcClient(ChainTypePolkadot, "wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
//	assert.NoError(t, err)
//	info, err := sc.GetPaymentQueryInfo(call.Extrinsic)
//	assert.NoError(t, err)
//	fmt.Println(info)
//}

func TestHash(t *testing.T) {
	h, _ := types.NewHashFromHexString("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	a, _ := types.EncodeToBytes(h)

	fmt.Println(hexutil.Encode(a))
}

//func TestGsrpcClient_Multisig2(t *testing.T) {
//	password := "123456"
//	os.Setenv(keystore.EnvPassword, password)
//
//	kp, err := keystore.KeypairFromAddress(Jun, keystore.SubChain, KeystorePath, false)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	krp := kp.(*sr25519.Keypair).AsKeyringPair()
//	stop := make(chan int)
//	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, krp, tlog, stop)
//	assert.NoError(t, err)
//
//	//pool, err := hexutil.Decode("0xac0df419ce0dc61b092a5cfa06a28e40cd82bc9de7e8c1e5591169360d66ba3c")
//	//assert.NoError(t, err)
//
//	wen, _ := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
//	//jun, _ := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
//	bao, _ := types.NewAddressFromHexAccountID("0x9c4189297ad2140c85861f64656d1d1318994599130d98b75ff094176d2ca31e")
//
//	others := []types.AccountID{wen.AsAccountID, bao.AsAccountID}
//	threshold := uint16(2)
//	targets := []types.Bytes{wen.AsAccountID[:]}
//
//	call, err := gc.NominateCall(targets)
//	assert.NoError(t, err)
//	fmt.Println(call)
//
//	fmt.Println("callHash", call.Extrinsic)
//
//	sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, tlog)
//	assert.NoError(t, err)
//	info, err := sc.GetPaymentQueryInfo(call.Extrinsic)
//	assert.NoError(t, err)
//	fmt.Println(info)
//
//	tp := types.TimePoint{Height: types.NewU32(1850), Index: types.NewU32(1)}
//	opTp := submodel.NewOptionTimePoint(tp)
//	ext, err := gc.NewUnsignedExtrinsic(config.MethodAsMulti, threshold, others, opTp, call.Opaque, false, info.Weight)
//	err = gc.SignAndSubmitTx(ext)
//	assert.NoError(t, err)
//}

func TestPool(t *testing.T) {
	p := "0x782a467d4ff23b660ca5f1ecf47f8537d4c35049541b6ebbf5381c00c4c158f7"
	pool, _ := hexutil.Decode(p)
	pbz, _ := types.EncodeToBytes(pool)
	fmt.Println(pool)
	fmt.Println(pbz)

}

func Test_KSM_GsrpcClient_Multisig_bond(t *testing.T) {

	logrus.SetLevel(logrus.TraceLevel)

	password := "tpkeeper"
	os.Setenv(keystore.EnvPassword, password)

	kp, err := keystore.KeypairFromAddress(relay1, keystore.SubChain, KeystorePath, false)
	if err != nil {
		t.Fatal(err)
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()
	stop := make(chan int)
	// sc, err := NewSarpcClient(ChainTypePolkadot, "wss://polkadot-rpc2.stafi.io", polkaTypesFile, AddressTypeMultiAddress, krp, tlog, stop)
	// sc, err := NewSarpcClient(ChainTypePolkadot, "wss://rpc.polkadot.io", polkaTypesFile, AddressTypeMultiAddress, krp, tlog, stop)
	// sc, err := NewSarpcClient(ChainTypePolkadot, "wss://public-rpc.pinknode.io/polkadot", polkaTypesFile, AddressTypeMultiAddress, krp, tlog, stop)
	// sc, err := NewSarpcClient(ChainTypePolkadot, "wss://1rpc.io/dot", polkaTypesFile, AddressTypeMultiAddress, krp, tlog, stop)
	sc, err := NewSarpcClient(ChainTypePolkadot, "wss://rpc.dotters.network/polkadot", polkaTypesFile, AddressTypeMultiAddress, krp, tlog, stop)
	// sc, err := NewSarpcClient(ChainTypePolkadot, "ws://127.0.0.1:9944", kusamaTypesFile, AddressTypeMultiAddress, krp, tlog, stop)
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

	unbond, _ := utils.StringToBigint("1000000000000")
	bond := big.NewInt(0)

	ext, err := sc.BondOrUnbondExtrinsic(bond, unbond)
	if err != nil {
		t.Fatal(err)
	}

	var cal types.Call
	switch ext := ext.(type) {
	case *types.Extrinsic:
		cal = ext.Method
	case *types.ExtrinsicMulti:
		cal = ext.Method
	default:
		t.Fatal("ext unsupported")
	}
	_ = cal
	// h := utils.BlakeTwo256(call.Opaque)

	// t.Log("Extrinsic", call.Extrinsic)
	// t.Log("Opaque", hexutil.Encode(call.Opaque))
	// t.Log("callHash", hexutil.Encode(h[:]))
	// calBts,err:=types.EncodeToBytes(cal)
	// if err!=nil{
	// 	t.Fatal(err)
	// }
	// t.Log(hexutil.Encode(calBts))

	extBz, err := types.EncodeToBytes(ext)
	if err != nil {
		t.Fatal(err)
	}
	info, err := sc.GetPaymentQueryInfo(hexutil.Encode(extBz))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("info", info.Class, info.PartialFee, info.Weight)

	//optp := types.TimePoint{Height: 1964877, Index: 1}
	//tp := submodel.NewOptionTimePoint(optp)
	// opaque, err := types.EncodeToBytes(cal)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	tp := submodel.NewOptionTimePointEmpty()
	multiExt, err := sc.NewUnsignedExtrinsic(config.MethodAsMulti, threshold, others, tp, cal, info.Weight)
	if err != nil {
		t.Fatal(err)
	}

	err = sc.SignAndSubmitTx(multiExt)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_KSM_GsrpcClient_Multisig_transfer(t *testing.T) {

	logrus.SetLevel(logrus.TraceLevel)

	password := "tpkeeper"
	os.Setenv(keystore.EnvPassword, password)

	kp, err := keystore.KeypairFromAddress(relay1, keystore.SubChain, KeystorePath, false)
	if err != nil {
		t.Fatal(err)
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()
	stop := make(chan int)
	// sc, err := NewSarpcClient(ChainTypePolkadot, "wss://kusama-test-rpc.stafi.io", kusamaTypesFile, AddressTypeMultiAddress, krp, tlog, stop)
	sc, err := NewSarpcClient(ChainTypePolkadot, "ws://127.0.0.1:9944", polkaTypesFile, AddressTypeMultiAddress, krp, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	threshold := uint16(2)
	relay2, _ := types.NewMultiAddressFromHexAccountID("0x2afeb305f32a12507a6b211d818218577b0e425692766b08b8bc5d714fccac3b")

	others := []types.AccountID{
		relay2.AsID,
	}

	ext, err := sc.TransferExtrinsic(relay2.AsID[:], types.NewUCompact(big.NewInt(1000000000000)))
	if err != nil {
		t.Fatal(err)
	}

	var cal types.Call
	switch ext := ext.(type) {
	case *types.Extrinsic:
		cal = ext.Method
	case *types.ExtrinsicMulti:
		cal = ext.Method
	default:
		t.Fatal("ext unsupported")
	}
	_ = cal
	// h := utils.BlakeTwo256(call.Opaque)

	// t.Log("Extrinsic", call.Extrinsic)
	// t.Log("Opaque", hexutil.Encode(call.Opaque))
	// t.Log("callHash", hexutil.Encode(h[:]))
	// calBts,err:=types.EncodeToBytes(cal)
	// if err!=nil{
	// 	t.Fatal(err)
	// }
	// t.Log(hexutil.Encode(calBts))

	extBz, err := types.EncodeToBytes(ext)
	if err != nil {
		t.Fatal(err)
	}
	info, err := sc.GetPaymentQueryInfo(hexutil.Encode(extBz))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("info", info.Class, info.PartialFee, info.Weight)
	weight := submodel.WeightV2{
		RefTime:   types.NewUCompact(big.NewInt(info.Weight)),
		ProofSize: types.NewUCompact(big.NewInt(0)),
	}
	//optp := types.TimePoint{Height: 1964877, Index: 1}
	//tp := submodel.NewOptionTimePoint(optp)
	// opaque, err := types.EncodeToBytes(cal)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	tp := submodel.NewOptionTimePointEmpty()
	multiExt, err := sc.NewUnsignedExtrinsic(config.MethodAsMulti, threshold, others, tp, cal, weight)
	if err != nil {
		t.Fatal(err)
	}

	err = sc.SignAndSubmitTx(multiExt)
	if err != nil {
		t.Fatal(err)
	}

}

func Test_AsMulti_transfer(t *testing.T) {

	logrus.SetLevel(logrus.TraceLevel)

	password := "tpkeeper"
	os.Setenv(keystore.EnvPassword, password)

	kp, err := keystore.KeypairFromAddress(relay2, keystore.SubChain, KeystorePath, false)
	if err != nil {
		t.Fatal(err)
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()
	stop := make(chan int)
	// sc, err := NewSarpcClient(ChainTypePolkadot, "wss://kusama-test-rpc.stafi.io", kusamaTypesFile, AddressTypeMultiAddress, krp, tlog, stop)
	sc, err := NewSarpcClient(ChainTypePolkadot, "ws://127.0.0.1:9944", polkaTypesFile, AddressTypeMultiAddress, krp, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}

	threshold := uint16(2)
	relay1, _ := types.NewMultiAddressFromHexAccountID("0x8e7750f4276116f8f089a5a4b24ca6577a13c7a1bcfe15868291b563336a7729")
	relay2, _ := types.NewMultiAddressFromHexAccountID("0x2afeb305f32a12507a6b211d818218577b0e425692766b08b8bc5d714fccac3b")

	others := []types.AccountID{
		relay1.AsID,
	}

	ext, err := sc.TransferExtrinsic(relay2.AsID[:], types.NewUCompact(big.NewInt(1000000000000)))
	if err != nil {
		t.Fatal(err)
	}

	var cal types.Call
	switch ext := ext.(type) {
	case *types.Extrinsic:
		cal = ext.Method
	case *types.ExtrinsicMulti:
		cal = ext.Method
	default:
		t.Fatal("ext unsupported")
	}
	_ = cal

	extBz, err := types.EncodeToBytes(ext)
	if err != nil {
		t.Fatal(err)
	}
	info, err := sc.GetPaymentQueryInfo(hexutil.Encode(extBz))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("info", info.Class, info.PartialFee, info.Weight)
	weight := submodel.WeightV2{
		RefTime:   types.NewUCompact(big.NewInt(info.Weight)),
		ProofSize: types.NewUCompact(big.NewInt(0)),
	}
	optp := types.TimePoint{Height: types.NewU32(1122), Index: 2}
	tp := submodel.NewOptionTimePoint(optp)

	multiExt, err := sc.NewUnsignedExtrinsic(config.MethodAsMulti, threshold, others, tp, cal, weight)
	if err != nil {
		t.Fatal(err)
	}

	err = sc.SignAndSubmitTx(multiExt)
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
	relay2, _ := types.NewMultiAddressFromHexAccountID("0x2afeb305f32a12507a6b211d818218577b0e425692766b08b8bc5d714fccac3b")

	ext, err := sc.TransferExtrinsic(relay2.AsID[:], types.NewUCompact(big.NewInt(1e10)))
	// ext, err := sc.NewUnsignedExtrinsic(config.MethodTransfer, relay2, types.NewUCompact(big.NewInt(1e10)))
	if err != nil {
		t.Fatal(err)
	}

	err = sc.SignAndSubmitTx(ext)
	if err != nil {
		t.Fatal(err)
	}
}
