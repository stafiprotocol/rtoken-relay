package solana

import (
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/solana"
	"github.com/stafiprotocol/rtoken-relay/shared/solana/vault"
	solClient "github.com/tpkeeper/solana-go-sdk/client"
	solCommon "github.com/tpkeeper/solana-go-sdk/common"
	solTypes "github.com/tpkeeper/solana-go-sdk/types"
)

type Connection struct {
	endpoint    string
	symbol      core.RSymbol
	poolClients map[string]*solana.PoolClient //map[poolAddressHexStr]poolClient
	log         log15.Logger
	stop        <-chan int
}

type PoolAccounts struct {
	FeeAccount            string
	StakeBaseAccount      string
	MultisigTxBaseAccount string
	MultisigInfoPubkey    string
	MultisigProgramId     string
}

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {
	poolAccounts := make(map[string]PoolAccounts)
	for _, account := range cfg.Accounts {
		accounts, ok := cfg.Opts[account].(PoolAccounts)
		if !ok {
			return nil, fmt.Errorf("account %s has no poolAccounts", account)
		}
		poolAccounts[account] = accounts
	}

	poolClientMap := make(map[string]*solana.PoolClient)

	v, err := vault.NewVaultFromWalletFile(cfg.KeystorePath)
	if err != nil {
		return nil, err
	}
	boxer, err := vault.SecretBoxerForType(v.SecretBoxWrap)
	if err != nil {
		return nil, fmt.Errorf("secret boxer: %w", err)
	}

	if err := v.Open(boxer); err != nil {
		return nil, fmt.Errorf("opening: %w", err)
	}

	privKeyMap := make(map[string]vault.PrivateKey)
	for _, privKey := range v.KeyBag {
		privKeyMap[privKey.PublicKey().String()] = privKey
	}

	for _, pool := range cfg.Accounts {

		pAccounts := poolAccounts[pool]

		poolAccounts := solana.PoolAccounts{
			FeeAccount:            solTypes.AccountFromPrivateKeyBytes(privKeyMap[pAccounts.FeeAccount]),
			StakeBaseAccount:      solTypes.AccountFromPrivateKeyBytes(privKeyMap[pAccounts.StakeBaseAccount]),
			MultisigTxBaseAccount: solTypes.AccountFromPrivateKeyBytes(privKeyMap[pAccounts.MultisigTxBaseAccount]),
			MultisigInfoPubkey:    solCommon.PublicKeyFromString(pAccounts.MultisigInfoPubkey),
			MultisignerPubkey:     solCommon.PublicKeyFromString(pool),
			MultisigProgramId:     solCommon.PublicKeyFromString(pAccounts.MultisigProgramId),
		}
		poolClientMap[pool] = solana.NewPoolClient(log, solClient.NewClient(cfg.Endpoint), poolAccounts)

	}

	return &Connection{
		endpoint:    cfg.Endpoint,
		symbol:      cfg.Symbol,
		log:         log,
		stop:        stop,
		poolClients: poolClientMap,
	}, nil
}
