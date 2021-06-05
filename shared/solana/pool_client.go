package solana

import (
	"github.com/ChainSafe/log15"
	solClient "github.com/tpkeeper/solana-go-sdk/client"
)

//one pool address with one poolClient
type PoolClient struct {
	eraBlockNumber int64
	log            log15.Logger
	rpcClient      *solClient.Client
	subKeyName     string //subKey is one of pubKeys of multiSig pool address,subKeyName is subKey`s name in keyring
}
