package substrate

import (
	"context"
	"fmt"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stafiprotocol/chainbridge/utils/keystore"
)

var (
	AliceKey     = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
	From         = "31yavGB5CVb8EwpqKQaS9XY7JZcfbK6QpWPn5kkweHVpqcov"
	KeystorePath = "/Users/fwj/Go/stafi/chainbridge/keys"
)

func TestGsrpcClient_QueryStorage(t *testing.T) {
	//gc, err := NewGsrpcClient(context.Background(), "wss://mainnet-rpc.stafi.io", AliceKey, tlog)
	gc, err := NewGsrpcClient(context.Background(), "ws://127.0.0.1:9944", AliceKey, tlog)
	assert.NoError(t, err)

	bondId, err := types.NewHashFromHexString("0x354f4c6ce7cd428336d7b93ed684a7c694a69818d338169dbad0312a7c15685d")
	assert.NoError(t, err)

	bk := &conn.BondKey{Symbol: conn.RDOT, BondId: bondId}

	bondKey, err := types.EncodeToBytes(bk)
	assert.NoError(t, err)

	var br conn.BondRecord

	re, err := gc.QueryStorage("RTokenSeries", "BondRecords", bondKey, nil, &br)
	assert.NoError(t, err)
	fmt.Println(re)
	fmt.Println(br)
}
