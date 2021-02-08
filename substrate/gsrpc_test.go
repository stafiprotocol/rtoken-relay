package substrate

import (
	"context"
	"fmt"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
)

var (
	AliceKey     = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
	From         = "31yavGB5CVb8EwpqKQaS9XY7JZcfbK6QpWPn5kkweHVpqcov"
	KeystorePath = "/Users/fwj/Go/stafi/chainbridge/keys"
)

func TestGsrpcClient_QueryStorage(t *testing.T) {
	tlog := log15.Root()

	gc, err := NewGsrpcClient(context.Background(), "wss://mainnet-rpc.stafi.io", AliceKey, tlog)
	//gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AliceKey, tlog)
	assert.NoError(t, err)

	era, err := types.EncodeToBytes(types.U32(545))
	assert.NoError(t, err)

	before := new(types.U128)
	after := new(types.U128)
	re, err := gc.QueryStorage("RFis", "TotalBondedBeforePayout", era, nil, before)
	assert.NoError(t, err)
	fmt.Println(re)
	fmt.Println(before)

	re1, err := gc.QueryStorage("RFis", "TotalBondedAfterPayout", era, nil, after)
	assert.NoError(t, err)
	fmt.Println(re1)
	fmt.Println(after)

	fmt.Println(before.Cmp(after.Int))

}
