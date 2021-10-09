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
	FeeAccount                  solTypes.Account
	StakeBasePubkeyToAccounts   map[solCommon.PublicKey]solTypes.Account //auth relay have
	StakeBaseAccountPubkeys     []solCommon.PublicKey                    //all relay must have
	MultisigTxBaseAccount       *solTypes.Account                        //auth relay have
	MultisigTxBaseAccountPubkey solCommon.PublicKey                      //all relay must have
	MultisigInfoPubkey          solCommon.PublicKey
	MultisignerPubkey           solCommon.PublicKey
	MultisigProgramId           solCommon.PublicKey
	HasBaseAccountAuth          bool
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
