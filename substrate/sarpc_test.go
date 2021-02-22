package substrate

import (
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	tlog = log15.Root()
)

const DefaultTypeFilePath = "../network/stafi.json"

func TestSarpcClient_GetBlock(t *testing.T) {
	bh := "0xeb2e9ba35b9009475d93def5e951173af8161f953d7f3a95417d582bf64986ac"
	//gc, err := NewGsrpcClient(context.Background(), , AliceKey, tlog)
	//sc, err := NewSarpcClient("ws://127.0.0.1:9944", DefaultTypeFilePath, tlog)
	sc, err := NewSarpcClient("wss://mainnet-rpc.stafi.io", DefaultTypeFilePath, tlog)
	assert.NoError(t, err)
	//blk, err := sc.GetBlock(bh)
	//assert.NoError(t, err)
	//fmt.Println(blk)

	ehs, err := sc.GetExtrinsics(bh)
	assert.NoError(t, err)
	for _, eh := range ehs {
		fmt.Println(eh.ExtrinsicHash)
	}
}

//func TestSarpcClient_GetChainEvents(t *testing.T) {
//	sc, err := NewSarpcClient("ws://127.0.0.1:9944", tlog)
//	assert.NoError(t, err)
//	//sc.GetEventsByModuleIdAndEventId()
//}
