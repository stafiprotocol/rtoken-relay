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

func TestSarpcClient_GetChainEvents(t *testing.T) {
	sc, err := NewSarpcClient("wss://stafi-seiya.stafi.io", DefaultTypeFilePath, tlog)
	assert.NoError(t, err)

	evts, err := sc.GetEvents(973215)
	assert.NoError(t, err)
	for _, evt := range evts {
		fmt.Println(evt.ModuleId)
		fmt.Println(evt.EventId)
		fmt.Println(evt.Params)
	}
}
