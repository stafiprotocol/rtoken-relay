package solana

import (
	"github.com/ChainSafe/log15"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
	solCommon "github.com/stafiprotocol/solana-go-sdk/common"
	solTypes "github.com/stafiprotocol/solana-go-sdk/types"
)

//one pool address with one poolClient
type PoolClient struct {
	log       log15.Logger
	rpcClient *solClient.Client
	PoolAccounts
}

type PoolAccounts struct {
	FeeAccount            solTypes.Account
	StakeBaseAccount      solTypes.Account
	StakeBaseAccounts     []solTypes.Account
	MultisigTxBaseAccount solTypes.Account
	MultisigInfoPubkey    solCommon.PublicKey
	MultisignerPubkey     solCommon.PublicKey
	MultisigProgramId     solCommon.PublicKey
}

func NewPoolClient(log log15.Logger, rpcClient *solClient.Client, poolAccount PoolAccounts) *PoolClient {
	return &PoolClient{
		log:          log,
		rpcClient:    rpcClient,
		PoolAccounts: poolAccount,
	}
}

func (p *PoolClient) GetRpcClient() *solClient.Client {
	return p.rpcClient
}
