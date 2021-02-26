package substrate

import (
	"context"
	"fmt"
	"github.com/stafiprotocol/chainbridge/utils/crypto/sr25519"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/stafiprotocol/chainbridge/utils/keystore"
)

var (
	AliceKey     = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
	From         = "31yavGB5CVb8EwpqKQaS9XY7JZcfbK6QpWPn5kkweHVpqcov"
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

	gc, err := NewGsrpcClient(context.Background(), "ws://127.0.0.1:9944", krp, tlog)
	assert.NoError(t, err)

	bId, err := types.NewHashFromHexString("0xd2f787195c3498f941653fe542e62b397988eeb3e2b867cf39a93c1ae5127b41")
	assert.NoError(t, err)

	bondKey := &conn.BondKey{
		Symbol: conn.RDOT,
		BondId: bId,
	}

	fmt.Println(bondKey)

	bk, err := types.EncodeToBytes(bondKey)
	assert.NoError(t, err)

	br := new(conn.BondRecord)
	exist, err := gc.QueryStorage(config.LiquidityBondModuleId, config.StorageBondRecords, bk, nil, br)
	assert.NoError(t, err)
	fmt.Println("exist:", exist)

	meta, err := gc.GetLatestMetadata()
	assert.NoError(t, err)

	call, err := types.NewCall(
		meta,
		config.ExecuteBondRecord,
		bondKey,
		conn.Pass,
	)
	assert.NoError(t, err)

	//prop := &conn.Proposal{call, bondKey}
	//symbol: RSymbol, prop_id: T::Hash, in_favour: bool
	ext, err := gc.NewUnsignedExtrinsic(config.RacknowledgeProposal, bondKey.Symbol, bondKey.BondId, true, call)
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

	gc, err := NewGsrpcClient(context.Background(), "ws://127.0.0.1:9944", krp, tlog)
	assert.NoError(t, err)

	bId, err := types.NewHashFromHexString("0x7942f736614444159dba03bc5fef684dcebb7c5edb16d4e57315e148ea4525be")
	assert.NoError(t, err)

	bondKey := &conn.BondKey{
		Symbol: conn.RDOT,
		BondId: bId,
	}
	//
	//fmt.Println(bondKey)
	//
	//bk, err := types.EncodeToBytes(bondKey)
	//assert.NoError(t, err)

	//br := new(conn.BondRecord)
	//exist, err := gc.QueryStorage(config.RTokenLedgerModuleId, config.StorageBondRecords, bk, nil, br)
	//assert.NoError(t, err)
	//fmt.Println("exist:", exist)

	//meta, err := gc.GetLatestMetadata()
	//assert.NoError(t, err)
	//
	//call, err := types.NewCall(
	//	meta,
	//	config.ExecuteBondRecord,
	//	bondKey.Symbol,
	//	bondKey.BondId[:],
	//	conn.Pass,
	//)
	//assert.NoError(t, err)

	//prop := &conn.Proposal{call, bondKey}
	//symbol: RSymbol, prop_id: T::Hash, in_favour: bool
	ext, err := gc.NewUnsignedExtrinsic("RTokenSeries.add_hashes", bondKey, conn.Pass)
	err = gc.SignAndSubmitTx(ext)
	assert.NoError(t, err)
}

//amount := uint64(1000000000000 * 10)
////// Create a call, transferring 12345 units to wen
//wen, err := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
//if err != nil {
//	panic(err)
//}
//
//ext, err := gc.NewUnsignedExtrinsic("Balances.transfer", wen, types.NewUCompactFromUInt(amount))
//
//err = gc.SignAndSubmitTx(ext)
//if err != nil {
//	t.Fatal(err)
//}
//x, _ := ext.MarshalJSON()
//fmt.Println(hex.EncodeToString(x))
