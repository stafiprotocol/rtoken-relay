package substrate

import (
	"context"
	"fmt"
	"testing"

	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stretchr/testify/assert"
)

func TestFullSubClient_TransferVerify(t *testing.T) {
	gc, err := NewGsrpcClient(context.Background(), "wss://stafi-seiya.stafi.io", AliceKey, tlog)
	assert.NoError(t, err)

	bId, err := types.NewHashFromHexString("0x875fb65bf1f6f2ea3369c8849cc91a88a5784ac0cb903800f4cf6a3e6c8dc0a7")
	assert.NoError(t, err)

	bondKey := &conn.BondKey{Symbol: conn.RDOT, BondId: bId}

	bk, err := types.EncodeToBytes(bondKey)
	assert.NoError(t, err)

	br := new(conn.BondRecord)
	exist, err := gc.QueryStorage(config.LiquidityBondModuleId, config.StorageBondRecords, bk, nil, br)
	assert.NoError(t, err)
	assert.True(t, exist)

	sc1, err := NewSarpcClient("wss://stafi-seiya.stafi.io", stafiTypesFile, tlog)
	assert.NoError(t, err)
	evts, err := sc1.GetEvents(973215)
	assert.NoError(t, err)
	fmt.Println(len(evts))

	sc, err := NewSarpcClient("wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	assert.NoError(t, err)

	fsc := &FullSubClient{sc, gc, nil, nil}
	reason, err := fsc.TransferVerify(br)
	assert.NoError(t, err)
	assert.Equal(t, reason, conn.Pass)
}
