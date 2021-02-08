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

func TestSarpcClient_GetBlock(t *testing.T) {
	bh := "0x6b2d076ea22c96accfd4afa2e6e4befdea35c77a73e256d2b76e66f884a6b996"
	sc, err := NewSarpcClient("ws://127.0.0.1:9944", tlog)
	assert.NoError(t, err)
	//blk, err := sc.GetBlock(bh)
	//assert.NoError(t, err)
	//fmt.Println(blk)

	err = sc.UpdateMeta(bh)
	assert.NoError(t, err)
	fmt.Println(sc.currentSpecVersion)

	ehs, err := sc.GetExtrinsicHashs(bh)
	assert.NoError(t, err)
	for _, eh := range ehs {
		if eh == "" {
			continue
		}
		fmt.Println(eh)
	}
}

//func TestSarpcClient_GetChainEvents(t *testing.T) {
//	sc, err := NewSarpcClient("ws://127.0.0.1:9944", tlog)
//	assert.NoError(t, err)
//	//sc.GetEventsByModuleIdAndEventId()
//}
