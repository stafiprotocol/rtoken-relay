package solana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/mr-tron/base58"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/solana"
	"github.com/stafiprotocol/rtoken-relay/shared/solana/vault"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
	solCommon "github.com/stafiprotocol/solana-go-sdk/common"
	solTypes "github.com/stafiprotocol/solana-go-sdk/types"
)

const MaxBase58SignatureLen = 88

type Connection struct {
	endpoint    string
	queryClient *solClient.Client
	symbol      core.RSymbol
	poolClients map[string]*solana.PoolClient //map[poolAddressHexStr]poolClient
	log         core.Logger
	stop        <-chan int
}

type PoolAccounts struct {
	FeeAccount            string   `json:"feeAccount"`
	StakeBaseAccounts     []string `json:"stakeBaseAccounts"`
	MultisigTxBaseAccount string   `json:"multisigTxBaseAccount"`
	MultisigInfoPubkey    string   `json:"multisigInfoPubkey"`
	MultisigProgramId     string   `json:"multisigProgramId"`
}

func NewConnection(cfg *core.ChainConfig, log core.Logger, stop <-chan int) (*Connection, error) {
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

	//collect privkey
	pubKeyStrToPrivKey := make(map[string]vault.PrivateKey)
	for _, privKey := range v.KeyBag {
		pubKeyStrToPrivKey[privKey.PublicKey().String()] = privKey
	}

	//create pool client
	for _, pool := range cfg.Accounts {
		pAccounts := poolAccounts[pool]
		stakeBasePubkeyToAccounts := make(map[solCommon.PublicKey]solTypes.Account)
		stakeBaseAccountPubkeys := make([]solCommon.PublicKey, 0)
		for _, account := range pAccounts.StakeBaseAccounts {
			accountPubkey := solCommon.PublicKeyFromString(account)
			//stake base account privkey may not exist
			if pKey, exist := pubKeyStrToPrivKey[account]; exist {
				stakeBasePubkeyToAccounts[accountPubkey] = solTypes.AccountFromPrivateKeyBytes(pKey)
			}
			stakeBaseAccountPubkeys = append(stakeBaseAccountPubkeys, accountPubkey)
		}

		//sort stake base account pubkeys
		sort.SliceStable(stakeBaseAccountPubkeys, func(i, j int) bool {
			return stakeBaseAccountPubkeys[i].ToBase58() < stakeBaseAccountPubkeys[j].ToBase58()
		})

		if _, exist := pubKeyStrToPrivKey[pAccounts.FeeAccount]; !exist {
			return nil, fmt.Errorf("feeAccount privkey must exist")
		}

		var multisigTxBaseAccount *solTypes.Account
		if pKey, exist := pubKeyStrToPrivKey[pAccounts.MultisigTxBaseAccount]; exist {
			account := solTypes.AccountFromPrivateKeyBytes(pKey)
			multisigTxBaseAccount = &account
		}

		//check auth
		hasBaseAccountAuth := false
		if len(stakeBaseAccountPubkeys) == len(stakeBasePubkeyToAccounts) && multisigTxBaseAccount != nil {
			hasBaseAccountAuth = true
		}

		poolAccounts := solana.PoolAccounts{
			FeeAccount:                  solTypes.AccountFromPrivateKeyBytes(pubKeyStrToPrivKey[pAccounts.FeeAccount]),
			StakeBaseAccountPubkeys:     stakeBaseAccountPubkeys,
			StakeBasePubkeyToAccounts:   stakeBasePubkeyToAccounts,
			MultisigTxBaseAccount:       multisigTxBaseAccount,
			MultisigTxBaseAccountPubkey: solCommon.PublicKeyFromString(pAccounts.MultisigTxBaseAccount),
			MultisigInfoPubkey:          solCommon.PublicKeyFromString(pAccounts.MultisigInfoPubkey),
			MultisignerPubkey:           solCommon.PublicKeyFromString(pool),
			MultisigProgramId:           solCommon.PublicKeyFromString(pAccounts.MultisigProgramId),
			HasBaseAccountAuth:          hasBaseAccountAuth,
		}
		poolClientMap[pool] = solana.NewPoolClient(log, solClient.NewClient([]string{cfg.Endpoint}), poolAccounts)

	}

	return &Connection{
		endpoint:    cfg.Endpoint,
		symbol:      cfg.Symbol,
		queryClient: solClient.NewClient([]string{cfg.Endpoint}),
		log:         log,
		stop:        stop,
		poolClients: poolClientMap,
	}, nil
}

func (c *Connection) TransferVerify(r *submodel.BondRecord) (submodel.BondReason, error) {
	hashBase58Str := base58.Encode(r.Txhash)
	if len(hashBase58Str) > MaxBase58SignatureLen {
		return submodel.TxhashUnmatch, nil
	}
	//check tx hash
	tx, err := c.queryClient.GetTransactionV2(context.Background(), hashBase58Str)
	if err != nil {
		return submodel.BondReasonDefault, err
	}
	//check block hash
	block, err := c.queryClient.GetBlock(context.Background(), tx.Slot, solClient.GetBlockConfig{
		Commitment:                     solClient.CommitmentFinalized,
		MaxSupportedTransactionVersion: &solClient.DefaultMaxSupportedTransactionVersion,
	})
	if err != nil {
		return submodel.BondReasonDefault, err
	}

	if !strings.EqualFold(block.Blockhash, base58.Encode(r.Blockhash)) {
		return submodel.BlockhashUnmatch, nil
	}

	//check  pool
	poolAccountIndex := -1
	for i, key := range tx.Transaction.Message.AccountKeys {
		if strings.EqualFold(key, base58.Encode(r.Pool)) {
			poolAccountIndex = i
			break
		}
	}
	if poolAccountIndex < 0 {
		return submodel.PoolUnmatch, nil
	}
	//check pubkey
	userAccountIndex := -1
	for i, key := range tx.Transaction.Message.AccountKeys {
		if strings.EqualFold(key, base58.Encode(r.Pubkey)) {
			userAccountIndex = i
			break
		}
	}
	if userAccountIndex < 0 {
		return submodel.PubkeyUnmatch, nil
	}

	//check amount
	if len(tx.Meta.PostBalances) <= poolAccountIndex ||
		len(tx.Meta.PreBalances) <= poolAccountIndex ||
		len(tx.Meta.PostBalances) <= userAccountIndex ||
		len(tx.Meta.PreBalances) <= userAccountIndex {
		return submodel.BondReasonDefault, fmt.Errorf("solana api postBalances or preBalances not right. hash %s",
			hashBase58Str)
	}
	amount := tx.Meta.PostBalances[poolAccountIndex] - tx.Meta.PreBalances[poolAccountIndex]
	if amount != r.Amount.Int64() {
		return submodel.AmountUnmatch, nil
	}
	//may cost fee
	amount = tx.Meta.PreBalances[userAccountIndex] - tx.Meta.PostBalances[userAccountIndex]
	if amount < r.Amount.Int64() {
		return submodel.AmountUnmatch, nil
	}

	return submodel.Pass, nil
}

func (c *Connection) GetPoolClient(poolAddrHexStr string) (*solana.PoolClient, error) {
	if sub, exist := c.poolClients[poolAddrHexStr]; exist {
		return sub, nil
	}
	return nil, errors.New("subClient of this pool not exist")
}

func (c *Connection) GetQueryClient() *solClient.Client {
	return c.queryClient
}
