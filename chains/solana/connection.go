package solana

import (
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/solana"
)

type Connection struct {
	url           string
	symbol        core.RSymbol
	currentHeight int64
	poolClients   map[string]*solana.PoolClient //map[addressHexStr]solClient
	log           log15.Logger
	stop          <-chan int
}
