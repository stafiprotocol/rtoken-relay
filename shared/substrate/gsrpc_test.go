package substrate

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/stafiprotocol/chainbridge/utils/crypto/sr25519"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
	"github.com/stretchr/testify/assert"
)

var (
	AliceKey     = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
	From         = "31yavGB5CVb8EwpqKQaS9XY7JZcfbK6QpWPn5kkweHVpqcov"
	LessPolka    = "1334v66HrtqQndbugYxX9m56V6222m97LbavB4KAMmqgjsas"
	From1        = "31d96Cq9idWQqPq3Ch5BFY84zrThVE3r98M7vG4xYaSWHwsX"
	From2        = "1TgYb5x8xjsZRyL5bwvxUoAWBn36psr4viSMHbRXA8bkB2h"
	KeystorePath = "/Users/fwj/Go/stafi/rtoken-relay/keys"
)

func TestGsrpcClient(t *testing.T) {
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

	bId, err := types.NewHashFromHexString("0xd2f787195c3498f941653fe542e62b397988eeb3e2b867cf39a93c1ae5127b41")
	assert.NoError(t, err)

	bondKey := &submodel.BondKey{
		Rsymbol: core.RDOT,
		BondId:  bId,
	}

	fmt.Println(bondKey)

	bk, err := types.EncodeToBytes(bondKey)
	assert.NoError(t, err)

	br := new(submodel.BondRecord)
	exist, err := gc.QueryStorage(config.LiquidityBondModuleId, config.StorageBondRecords, bk, nil, br)
	assert.NoError(t, err)
	fmt.Println("exist:", exist)

	meta, err := gc.GetLatestMetadata()
	assert.NoError(t, err)

	call, err := types.NewCall(
		meta,
		config.MethodExecuteBondRecord,
		bondKey,
		submodel.Pass,
	)
	assert.NoError(t, err)

	//prop := &conn.Proposal{call, bondKey}
	//symbol: RSymbol, prop_id: T::Key, in_favour: bool
	ext, err := gc.NewUnsignedExtrinsic(config.MethodRacknowledgeProposal, bondKey.Rsymbol, bondKey.BondId, true, call)
	err = gc.SignAndSubmitTx(ext)
	assert.NoError(t, err)
}

func TestGsrpcClient1(t *testing.T) {
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

	bId, err := types.NewHashFromHexString("0x7942f736614444159dba03bc5fef684dcebb7c5edb16d4e57315e148ea4525be")
	assert.NoError(t, err)

	bondKey := &submodel.BondKey{
		Rsymbol: core.RDOT,
		BondId:  bId,
	}

	//prop := &conn.Proposal{call, bondKey}
	//symbol: RSymbol, prop_id: T::Key, in_favour: bool
	ext, err := gc.NewUnsignedExtrinsic("RTokenSeries.add_hashes", bondKey, submodel.Pass)
	err = gc.SignAndSubmitTx(ext)
	assert.NoError(t, err)
}

func TestGsrpcClient_StakingActive(t *testing.T) {
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

	b, _ := hexutil.Decode("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")

	jun := types.NewAddressFromAccountID(b)
	s, err := gc.StakingLedger(jun.AsAccountID)
	fmt.Println(s.Active)
}

func TestGsrpcClient_QueryStorage(t *testing.T) {
	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, AliceKey, tlog, stop)
	assert.NoError(t, err)

	//pool, err := hexutil.Decode("0xbeb93b63149fd6a1b69f2d50492fe01983be085cbf09f689219838f6322763d8")
	//assert.NoError(t, err)
	//
	//key := submodel.PoolKey{Rsymbol: core.RDOT, Pool: pool}
	//keybz, err := types.EncodeToBytes(key)
	//assert.NoError(t, err)
	type poolkey struct {
		Rsymbol core.RSymbol
		Pool    types.Bytes
		Era     types.U32
	}
	pool, err := hexutil.Decode("0xeda331e37bf66b2393c4c271e384dfaa2bfcdd35")
	bz, err := types.EncodeToBytes(poolkey{core.RATOM, types.NewBytes(pool), types.U32(11386)})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hexutil.Encode(bz))
	type Unbonding struct {
		Who       types.AccountID
		Value     types.U128
		Recipient types.Bytes
	}
	unbonds := make([]Unbonding, 0)
	exist, err := gc.QueryStorage(config.RTokenLedgerModuleId, config.StoragePoolUnbonds, bz, nil, &unbonds)
	assert.NoError(t, err)
	assert.True(t, exist)
	for _, ac := range unbonds {
		t.Log(fmt.Println(hexutil.Encode(ac.Recipient)))
		addr:=types.NewAddressFromAccountID(ac.Recipient)
		t.Log(hexutil.Encode(addr.AsAccountID[:]))
	}
}

func TestGsrpcClient_Multisig(t *testing.T) {
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

	pk := &submodel.PoolKey{core.RDOT, pool}
	pkBz, err := types.EncodeToBytes(pk)
	assert.NoError(t, err)

	var threshold uint16
	exist, err := gc.QueryStorage(config.RTokenLedgerModuleId, config.StorageMultiThresholds, pkBz, nil, &threshold)
	assert.NoError(t, err)
	assert.True(t, exist)
	fmt.Println(threshold)

	subs := make([]types.Bytes, 0)
	exist, err = gc.QueryStorage(config.RTokenLedgerModuleId, config.StorageSubAccounts, pkBz, nil, &subs)
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
}

func TestStorages(t *testing.T) {
	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, AliceKey, tlog, stop)
	assert.NoError(t, err)

	addr, err := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
	assert.NoError(t, err)

	targets := make([]types.AccountID, 0)
	exist, err := gc.QueryStorage(config.StakingModuleId, config.StorageNominators, addr.AsAccountID[:], nil, &targets)
	assert.NoError(t, err)
	fmt.Println(exist)
	//fmt.Println(targets)
	for _, t := range targets {
		fmt.Println(hexutil.Encode(t[:]))
	}

	era := uint32(5)
	bz, err := types.EncodeToBytes(era)
	assert.NoError(t, err)

	re := new(submodel.EraRewardPoints)
	exist1, err := gc.QueryStorage(config.StakingModuleId, config.StorageErasRewardPoints, bz, nil, re)
	assert.NoError(t, err)
	fmt.Println(exist1)
	fmt.Println(re)

	ledger := new(submodel.StakingLedger)
	exist, err = gc.QueryStorage(config.StakingModuleId, config.StorageLedger, addr.AsAccountID[:], nil, ledger)
	assert.NoError(t, err)
	fmt.Println(exist)
	fmt.Println(ledger)
}

func TestBatch(t *testing.T) {
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

	addr, err := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
	assert.NoError(t, err)

	amount, _ := utils.StringToBigint("1000000000000")
	value := types.NewUCompact(amount)

	calls := make([]types.Call, 0)
	meta, err := gc.GetLatestMetadata()
	assert.NoError(t, err)

	for i := 0; i <= 4000; i++ {
		call, err := types.NewCall(
			meta,
			config.MethodTransferKeepAlive,
			addr,
			value,
		)
		assert.NoError(t, err)
		calls = append(calls, call)
	}

	ext, err := gc.NewUnsignedExtrinsic(config.MethodBatch, calls)
	assert.NoError(t, err)

	err = gc.SignAndSubmitTx(ext)
	assert.NoError(t, err)
}

func TestBatchTransfer(t *testing.T) {
	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, AliceKey, tlog, stop)
	assert.NoError(t, err)

	less, _ := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
	jun, _ := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
	wen, _ := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	bao, _ := types.NewAddressFromHexAccountID("0x9c4189297ad2140c85861f64656d1d1318994599130d98b75ff094176d2ca31e")

	addrs := []types.Address{less, jun, wen, bao}

	amount, _ := utils.StringToBigint("3000" + "000000000000")
	value := types.NewUCompact(amount)

	calls := make([]types.Call, 0)
	meta, err := gc.GetLatestMetadata()
	assert.NoError(t, err)

	for _, addr := range addrs {
		call, err := types.NewCall(
			meta,
			config.MethodTransferKeepAlive,
			addr,
			value,
		)
		assert.NoError(t, err)
		calls = append(calls, call)
	}

	ext, err := gc.NewUnsignedExtrinsic(config.MethodBatch, calls)
	assert.NoError(t, err)

	err = gc.SignAndSubmitTx(ext)
	assert.NoError(t, err)
}

func TestMultiBatch(t *testing.T) {
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

	//less, _ := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
	jun, _ := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
	wen, _ := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	bao, _ := types.NewAddressFromHexAccountID("0x9c4189297ad2140c85861f64656d1d1318994599130d98b75ff094176d2ca31e")
	//mul, _ := types.NewAddressFromHexAccountID("0xbeb93b63149fd6a1b69f2d50492fe01983be085cbf09f689219838f6322763d8")
	threshold := uint16(2)
	others := []types.AccountID{wen.AsAccountID, jun.AsAccountID}
	addrs := []types.Address{bao, jun, bao}

	amount, _ := utils.StringToBigint("1000000000000")
	value := types.NewUCompact(amount)

	calls := make([]types.Call, 0)

	var info *rpc.PaymentQueryInfo
	tp := submodel.NewOptionTimePointEmpty()
	for _, addr := range addrs {
		c, err := gc.TransferCall(addr, value)
		fmt.Println(c.CallHash)
		assert.NoError(t, err)
		if info == nil {
			sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, tlog)
			assert.NoError(t, err)
			info, err = sc.GetPaymentQueryInfo(c.Extrinsic)
			assert.NoError(t, err)
		}

		ext, err := gc.NewUnsignedExtrinsic(config.MethodAsMulti, threshold, others, tp, c.Opaque, false, info.Weight)
		xt, ok := ext.(*types.Extrinsic)
		assert.True(t, ok)
		calls = append(calls, xt.Method)
	}

	ext, err := gc.NewUnsignedExtrinsic(config.MethodBatch, calls)
	assert.NoError(t, err)

	err = gc.SignAndSubmitTx(ext)
	assert.NoError(t, err)
}

func TestBalance(t *testing.T) {
	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, AliceKey, tlog, stop)
	assert.NoError(t, err)

	less, _ := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
	ac := new(types.AccountInfo)
	exist, err := gc.QueryStorage(config.SystemModuleId, config.StorageAccount, less.AsAccountID[:], nil, ac)
	assert.NoError(t, err)

	fmt.Println(exist)
	fmt.Println(ac.Data.Free)
}

func TestGetConst(t *testing.T) {
	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, AliceKey, tlog, stop)
	assert.NoError(t, err)
	e, err := gc.ExistentialDeposit()
	assert.NoError(t, err)
	fmt.Println(e)
}

func TestBondExtra(t *testing.T) {
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

	amount, _ := utils.StringToBigint("10000000000000")
	value := types.NewUCompact(amount)
	ext, err := gc.NewUnsignedExtrinsic(config.MethodBondExtra, value)
	assert.NoError(t, err)

	err = gc.SignAndSubmitTx(ext)
	assert.NoError(t, err)
}

func TestPolkaBondExtra(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kp, err := keystore.KeypairFromAddress(LessPolka, keystore.SubChain, KeystorePath, false)
	if err != nil {
		t.Fatal(err)
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()
	stop := make(chan int)
	gc, err := NewGsrpcClient("wss://polkadot-test-rpc.stafi.io", AddressTypeMultiAddress, krp, tlog, stop)
	assert.NoError(t, err)

	amount, _ := utils.StringToBigint("10000000000000")
	value := types.NewUCompact(amount)
	ext, err := gc.NewUnsignedExtrinsic(config.MethodBondExtra, value)
	assert.NoError(t, err)

	err = gc.SignAndSubmitTx(ext)
	assert.NoError(t, err)
}

func TestFreeBalance(t *testing.T) {
	less, _ := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")

	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, AliceKey, tlog, stop)
	assert.NoError(t, err)
	free, err := gc.FreeBalance(less.AsAccountID[:])
	assert.NoError(t, err)
	fmt.Println(free)

	lessPolka, _ := types.NewMultiAddressFromHexAccountID("0x5a0c23479ba36898acb44e163fe58a9155d7b508041cc1b5d5ad6bbd3d5a6360")
	gc1, err := NewGsrpcClient("wss://polkadot-test-rpc.stafi.io", AddressTypeMultiAddress, AliceKey, tlog, stop)
	assert.NoError(t, err)

	free1, err := gc1.FreeBalance(lessPolka.AsID[:])
	assert.NoError(t, err)
	fmt.Println(free1)
}
