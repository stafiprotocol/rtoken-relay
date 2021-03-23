package substrate

import (
	"fmt"
	"github.com/stafiprotocol/rtoken-relay/utils"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/utils/crypto/sr25519"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stretchr/testify/assert"
)

var (
	AliceKey     = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
	From         = "31yavGB5CVb8EwpqKQaS9XY7JZcfbK6QpWPn5kkweHVpqcov"
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
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", krp, tlog, stop)
	assert.NoError(t, err)

	bId, err := types.NewHashFromHexString("0xd2f787195c3498f941653fe542e62b397988eeb3e2b867cf39a93c1ae5127b41")
	assert.NoError(t, err)

	bondKey := &core.BondKey{
		Rsymbol: core.RDOT,
		BondId:  bId,
	}

	fmt.Println(bondKey)

	bk, err := types.EncodeToBytes(bondKey)
	assert.NoError(t, err)

	br := new(core.BondRecord)
	exist, err := gc.QueryStorage(config.LiquidityBondModuleId, config.StorageBondRecords, bk, nil, br)
	assert.NoError(t, err)
	fmt.Println("exist:", exist)

	meta, err := gc.GetLatestMetadata()
	assert.NoError(t, err)

	call, err := types.NewCall(
		meta,
		config.ExecuteBondRecord,
		bondKey,
		core.Pass,
	)
	assert.NoError(t, err)

	//prop := &conn.Proposal{call, bondKey}
	//symbol: RSymbol, prop_id: T::Hash, in_favour: bool
	ext, err := gc.NewUnsignedExtrinsic(config.RacknowledgeProposal, bondKey.Rsymbol, bondKey.BondId, true, call)
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
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", krp, tlog, stop)
	assert.NoError(t, err)

	bId, err := types.NewHashFromHexString("0x7942f736614444159dba03bc5fef684dcebb7c5edb16d4e57315e148ea4525be")
	assert.NoError(t, err)

	bondKey := &core.BondKey{
		Rsymbol: core.RDOT,
		BondId:  bId,
	}

	//prop := &conn.Proposal{call, bondKey}
	//symbol: RSymbol, prop_id: T::Hash, in_favour: bool
	ext, err := gc.NewUnsignedExtrinsic("RTokenSeries.add_hashes", bondKey, core.Pass)
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
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", krp, tlog, stop)
	assert.NoError(t, err)

	b, _ := hexutil.Decode("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")

	jun := types.NewAddressFromAccountID(b)
	s, err := gc.StakingLedger(jun.AsAccountID)
	fmt.Println(s.Active)
}

func TestGsrpcClient_QueryStorage(t *testing.T) {
	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AliceKey, tlog, stop)
	assert.NoError(t, err)

	pool, err := hexutil.Decode("0xbeb93b63149fd6a1b69f2d50492fe01983be085cbf09f689219838f6322763d8")
	assert.NoError(t, err)

	key := core.PoolKey{Rsymbol: core.RDOT, Pool: pool}
	keybz, err := types.EncodeToBytes(key)
	assert.NoError(t, err)

	subAcs := make([]types.Bytes, 0)
	exist, err := gc.QueryStorage(config.RTokenLedgerModuleId, config.StorageSubAccounts, keybz, nil, &subAcs)
	assert.NoError(t, err)
	assert.True(t, exist)
	for _, ac := range subAcs {
		fmt.Println(hexutil.Encode(ac))
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
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", krp, tlog, stop)
	assert.NoError(t, err)

	pool, err := hexutil.Decode("0xbeb93b63149fd6a1b69f2d50492fe01983be085cbf09f689219838f6322763d8")
	assert.NoError(t, err)

	pk := &core.PoolKey{core.RDOT, pool}
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

	extStr, opaque, err := gc.BondOrUnbondCall(bond, unbond)
	assert.NoError(t, err)
	fmt.Println(extStr)
	fmt.Println(hexutil.Encode(opaque))
	h := utils.BlakeTwo256(opaque)
	fmt.Println("callHash", hexutil.Encode(h[:]))

	sc, err := NewSarpcClient("ws://127.0.0.1:9944", stafiTypesFile, tlog)
	assert.NoError(t, err)

	info, err := sc.GetPaymentQueryInfo(extStr)
	assert.NoError(t, err)
	fmt.Println("info", info.Class, info.PartialFee, info.Weight)

	//tp := types.TimePoint{3010, 2}
	//callHash, _ := types.NewHashFromHexString("0xba6c8ec1798285f8f312523e2353ebe8468fab4b55afe1a788a64a65f8bcc72c")
	//ext, err := gc.NewUnsignedExtrinsic(config.MethodApproveAsMulti, threshold, others, tp, callHash, types.Weight(uint64(info.Weight)))

	//fmt.Println(opaque)

	tp := core.NewOptionTimePointEmpty()
	ext, err := gc.NewUnsignedExtrinsic(config.MethodAsMulti, threshold, others, tp, opaque, false, info.Weight)
	err = gc.SignAndSubmitTx(ext)
	assert.NoError(t, err)
}

func TestEncodedExtrinsic(t *testing.T) {
	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AliceKey, tlog, stop)
	assert.NoError(t, err)

	bond, _ := utils.StringToBigint("10000000000000")
	unbond := big.NewInt(0)

	extStr, opaque, err := gc.BondOrUnbondCall(bond, unbond)
	assert.NoError(t, err)

	//e, err := gc.NewUnsignedExtrinsic(config.MethodBondExtra, types.NewUCompact(bond))
	//assert.NoError(t, err)
	//bz, err := types.EncodeToBytes(e)
	//assert.NoError(t, err)

	sc, err := NewSarpcClient("ws://127.0.0.1:9944", stafiTypesFile, tlog)
	assert.NoError(t, err)

	info, err := sc.GetPaymentQueryInfo(extStr)
	assert.NoError(t, err)
	fmt.Println("info", info.Class, info.PartialFee, info.Weight)
	fmt.Println(opaque)
}

func TestStorages(t *testing.T) {
	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AliceKey, tlog, stop)
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

	re := new(EraRewardPoints)
	exist1, err := gc.QueryStorage(config.StakingModuleId, config.StorageErasRewardPoints, bz, nil, re)
	assert.NoError(t, err)
	fmt.Println(exist1)
	fmt.Println(re)

	ledger := new(StakingLedger)
	exist, err = gc.QueryStorage(config.StakingModuleId, config.StorageLedger, addr.AsAccountID[:], nil, ledger)
	assert.NoError(t, err)
	fmt.Println(exist)
	fmt.Println(ledger)
}
