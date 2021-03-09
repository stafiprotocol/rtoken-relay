package substrate

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/utils/crypto/sr25519"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"

	"github.com/stafiprotocol/chainbridge/utils/keystore"
)

var (
	AliceKey     = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
	From         = "31yavGB5CVb8EwpqKQaS9XY7JZcfbK6QpWPn5kkweHVpqcov"
	From1        = "31d96Cq9idWQqPq3Ch5BFY84zrThVE3r98M7vG4xYaSWHwsX"
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

	//prop := &conn.Proposal{call, bondKey}
	//symbol: RSymbol, prop_id: T::Hash, in_favour: bool
	ext, err := gc.NewUnsignedExtrinsic("RTokenSeries.add_hashes", bondKey, conn.Pass)
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

	gc, err := NewGsrpcClient(context.Background(), "ws://127.0.0.1:9944", krp, tlog)
	assert.NoError(t, err)

	b, _ := hexutil.Decode("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")

	jun := types.NewAddressFromAccountID(b)
	s, err := gc.StakingActive(jun.AsAccountID)
	fmt.Println(s.Active)
}

func TestGsrpcClient_Bond(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kp, err := keystore.KeypairFromAddress(From1, keystore.SubChain, KeystorePath, false)
	if err != nil {
		t.Fatal(err)
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()

	gc, err := NewGsrpcClient(context.Background(), "ws://127.0.0.1:9944", krp, tlog)
	assert.NoError(t, err)

	err = gc.bond(big.NewInt(10000000000000))
	assert.NoError(t, err)

	//bob, err := types.NewAddressFromHexAccountID("0xf6241901b8e0048421427ef6cd3513865c2b6d2ad3ca2c3d95d28dfca2b4f722")
	//if err != nil {
	//	panic(err)
	//}
	//
	//amount := uint64(1000000000000 * 10)
	//ext, err := gc.NewUnsignedExtrinsic("Balances.transfer", bob, types.NewUCompactFromUInt(amount))
}
