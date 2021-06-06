package solana

import (
	"encoding/json"
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
	FeeAccount            string `json:"feeAccount"`
	StakeBaseAccount      string `json:"stakeBaseAccount"`
	MultisigTxBaseAccount string `json:"multisigTxBaseAccount"`
	MultisigInfoPubkey    string `json:"multisigInfoPubkey"`
	MultisigProgramId     string `json:"multisigProgramId"`
}

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {
	poolAccounts := make(map[string]PoolAccounts)
	for _, account := range cfg.Accounts {
		paBts, err := json.Marshal(cfg.Opts[account])
		if err != nil {
			return nil, err
		}
		accounts := PoolAccounts{}
		err = json.Unmarshal(paBts, &accounts)
		if err != nil {
			return nil, fmt.Errorf("account %s unmarshal poolAccounts err %s", account, err)
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
