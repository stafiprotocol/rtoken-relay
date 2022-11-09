package main

import (
	"context"
	"fmt"
	"time"

	"github.com/stafiprotocol/rtoken-relay/shared/solana/vault"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
	solCommon "github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/multisigprog"
	"github.com/stafiprotocol/solana-go-sdk/stakeprog"
	"github.com/stafiprotocol/solana-go-sdk/sysprog"
	solTypes "github.com/stafiprotocol/solana-go-sdk/types"
	"github.com/urfave/cli/v2"
)

var retryLimit = 200
var initStakeAmount = uint64(10000)

// 1 create stakeBaseAccount if not exist on chain
// 2 create multisigInfo account if not exist on chain
// 3 init stake for stakeBaseAccount if stakeBaseAccount`s stake amount is zero
func initAction(ctx *cli.Context) error {
	path := ctx.String(configFlag.Name)
	pc := PoolAccounts{}
	err := loadConfig(path, &pc)
	if err != nil {
		return err
	}
	fmt.Printf("accounts info:\n %+v\n\n", pc)
	v, err := vault.NewVaultFromWalletFile(pc.KeystorePath)
	if err != nil {
		return err
	}
	boxer, err := vault.SecretBoxerForType(v.SecretBoxWrap)
	if err != nil {
		return fmt.Errorf("secret boxer: %w", err)
	}

	if err := v.Open(boxer); err != nil {
		return fmt.Errorf("opening: %w", err)
	}

	privKeyMap := make(map[string]vault.PrivateKey)
	for _, privKey := range v.KeyBag {
		privKeyMap[privKey.PublicKey().String()] = privKey
	}

	stakeBaseStrToValidator := make(map[string]solCommon.PublicKey)
	stakeBaseStrToAccount := make(map[string]solTypes.Account)
	for stakeBaseStr, validatorStr := range pc.StakeBaseAccountToValidator {
		if privateKey, exist := privKeyMap[stakeBaseStr]; exist {
			stakeBaseStrToAccount[stakeBaseStr] = solTypes.AccountFromPrivateKeyBytes(privateKey)
			stakeBaseStrToValidator[stakeBaseStr] = solCommon.PublicKeyFromString(validatorStr)
		} else {
			return fmt.Errorf("stakeBaseAccount: %s, doesn't have privateKey", stakeBaseStr)
		}
	}

	if len(pc.StakeBaseAccountToValidator) != len(stakeBaseStrToAccount) ||
		len(stakeBaseStrToAccount) != len(stakeBaseStrToValidator) {
		return fmt.Errorf("stakeBaseAccountToValidator config not right")
	}

	FeeAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.FeeAccount])
	MultisigTxBaseAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.MultisigTxBaseAccount])
	MultisigInfoAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.MultisigInfoPubkey])
	MultisigProgramId := solCommon.PublicKeyFromString(pc.MultisigProgramId)

	otherFeeAccount := make([]solTypes.Account, 0)
	owners := make([]solCommon.PublicKey, 0)
	owners = append(owners, FeeAccount.PublicKey)
	for _, account := range pc.OtherFeeAccount {
		a := solTypes.AccountFromPrivateKeyBytes(privKeyMap[account])
		otherFeeAccount = append(otherFeeAccount, a)
		owners = append(owners, a.PublicKey)
	}
	multisignerPubkey, nonce, err := solCommon.FindProgramAddress(
		[][]byte{MultisigInfoAccount.PublicKey.Bytes()}, MultisigProgramId)
	if err != nil {
		return err
	}
	fmt.Println("multisigner:", multisignerPubkey.ToBase58())

	c := solClient.NewClient([]string{pc.Endpoint})

	stakeAccountMiniMum, err := c.GetMinimumBalanceForRentExemption(context.Background(),
		solClient.StakeAccountInfoLengthDefault)
	if err != nil {
		return err
	}

	multisigAccountMiniMum, err := c.GetMinimumBalanceForRentExemption(context.Background(),
		solClient.MultisigInfoAccountLengthDefault)
	if err != nil {
		return err
	}

	//1 create stakeBaseAccount if not exist on chain
	for stakeAccoountStr, account := range stakeBaseStrToAccount {
		res, err := c.GetRecentBlockhash(context.Background())
		if err != nil {
			return err
		}
		_, err = c.GetStakeAccountInfo(context.Background(), stakeAccoountStr)
		if err != nil && err != solClient.ErrAccountNotFound {
			return err
		}
		if err == nil {
			fmt.Printf("stakeBaseAccount: %s exist on chain and will not create\n", stakeAccoountStr)
			continue
		}

		//create stakeBaseAccount and transfer 100+rent lamport to stakeBaseAccount
		rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
			Instructions: []solTypes.Instruction{
				sysprog.CreateAccount(
					FeeAccount.PublicKey,
					account.PublicKey,
					solCommon.StakeProgramID,
					initStakeAmount+stakeAccountMiniMum,
					solClient.StakeAccountInfoLengthDefault,
				),
				stakeprog.Initialize(
					account.PublicKey,
					stakeprog.Authorized{
						Staker:     multisignerPubkey,
						Withdrawer: multisignerPubkey,
					},
					stakeprog.Lockup{},
				),
			},
			Signers:         []solTypes.Account{FeeAccount, account},
			FeePayer:        FeeAccount.PublicKey,
			RecentBlockHash: res.Blockhash,
		})
		if err != nil {
			return fmt.Errorf("generate tx error, err: %v", err)
		}

		txHash, err := c.SendRawTransaction(context.Background(), rawTx)
		if err != nil {
			return fmt.Errorf("send tx error, err: %v", err)
		}
		fmt.Println("create stakeBaseAccount txHash:", txHash, stakeAccoountStr)
		time.Sleep(time.Second * 2)

		retry := 0
		for {
			if retry > retryLimit {
				return err
			}
			_, err = c.GetStakeAccountInfo(context.Background(), stakeAccoountStr)
			if err != nil {
				fmt.Printf("stakeBaseAccount %s not init yet, waiting...\n", MultisigInfoAccount.PublicKey.ToBase58())
				retry++
				time.Sleep(3 * time.Second)
				continue
			}
			break
		}
	}

	//2 create multisigInfo account if not exist on chain
	_, err = c.GetMultisigInfoAccountInfo(ctx.Context, MultisigInfoAccount.PublicKey.ToBase58())
	if err != nil {
		if err != solClient.ErrAccountNotFound {
			return err
		}
		res, err := c.GetRecentBlockhash(context.Background())
		if err != nil {
			return fmt.Errorf("get recent block hash error, err: %v", err)
		}
		rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
			Instructions: []solTypes.Instruction{
				sysprog.CreateAccount(
					FeeAccount.PublicKey,
					MultisigInfoAccount.PublicKey,
					MultisigProgramId,
					multisigAccountMiniMum*2,
					solClient.MultisigInfoAccountLengthDefault,
				),
				multisigprog.CreateMultisig(
					MultisigProgramId,
					MultisigInfoAccount.PublicKey,
					owners,
					uint64(pc.Threshold),
					uint8(nonce),
				),
			},
			Signers:         []solTypes.Account{FeeAccount, MultisigInfoAccount},
			FeePayer:        FeeAccount.PublicKey,
			RecentBlockHash: res.Blockhash,
		})
		if err != nil {
			return fmt.Errorf("generate tx error, err: %v", err)
		}
		txHash, err := c.SendRawTransaction(context.Background(), rawTx)
		if err != nil {
			return fmt.Errorf("send tx error, err: %v", err)
		}
		fmt.Println("create multisigInfoAccount txHash:", txHash)
		time.Sleep(time.Second * 2)

		retry := 0
		for {
			if retry > retryLimit {
				return err
			}
			_, err = c.GetMultisigInfoAccountInfo(ctx.Context, MultisigInfoAccount.PublicKey.ToBase58())
			if err != nil {
				fmt.Printf("multisigInfoAccount %s not init yet, waiting...\n", MultisigInfoAccount.PublicKey.ToBase58())
				retry++
				time.Sleep(3 * time.Second)
				continue
			}
			break
		}

	} else {
		fmt.Printf("multisigInfoAccount: %s exist on chain and will not create\n", MultisigInfoAccount.PublicKey.ToBase58())
	}

	//3 init stake for stakeBaseAccount if stakeBaseAccount`s stake amount is zero
	for stakeAccoountStr, account := range stakeBaseStrToAccount {
		var accountInfo *solClient.StakeAccountRsp
		var err error
		retry := 0
		for {
			if retry > retryLimit {
				return err
			}
			accountInfo, err = c.GetStakeAccountInfo(context.Background(), stakeAccoountStr)
			if err != nil {
				retry++
				time.Sleep(time.Second * 3)
				continue
			}
			break
		}
		stakeAmount := accountInfo.StakeAccount.Info.Stake.Delegation.Stake
		if stakeAmount > 0 {
			fmt.Printf("stakeBaseAccount: %s has stake %d will skip init\n", stakeAccoountStr, stakeAmount)
			continue
		}

		fmt.Printf("\nstart init stake for stakeBaseAccount: %s\n", stakeAccoountStr)
		//create derived multisig tx account if not exist onchain
		multisigTxAccountPubkey, multisigTxAccountSeed := getMultisigTxAccountPubkey(
			MultisigTxBaseAccount.PublicKey,
			MultisigProgramId,
			account.PublicKey,
			0)
		validatorPubkey := stakeBaseStrToValidator[stakeAccoountStr]
		stakeInstruction := stakeprog.DelegateStake(account.PublicKey, multisignerPubkey, validatorPubkey)

		_, err = c.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
		if err != nil {
			if err != solClient.ErrAccountNotFound {
				return err
			}

			res, err := c.GetRecentBlockhash(context.Background())
			if err != nil {
				return fmt.Errorf("get recent block hash error, err: %v", err)
			}

			rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
				Instructions: []solTypes.Instruction{
					sysprog.CreateAccountWithSeed(
						FeeAccount.PublicKey,
						multisigTxAccountPubkey,
						MultisigTxBaseAccount.PublicKey,
						MultisigProgramId,
						multisigTxAccountSeed,
						multisigAccountMiniMum*2,
						solClient.MultisigTxAccountLengthDefault,
					),
					multisigprog.CreateTransaction(
						MultisigProgramId,
						[]solCommon.PublicKey{solCommon.StakeProgramID},
						[][]solTypes.AccountMeta{stakeInstruction.Accounts},
						[][]byte{stakeInstruction.Data},
						MultisigInfoAccount.PublicKey,
						multisigTxAccountPubkey,
						FeeAccount.PublicKey,
					),
				},
				Signers:         []solTypes.Account{FeeAccount, MultisigTxBaseAccount},
				FeePayer:        FeeAccount.PublicKey,
				RecentBlockHash: res.Blockhash,
			})

			if err != nil {
				return fmt.Errorf("generate createTransaction tx error, err: %v", err)
			}

			txHash, err := c.SendRawTransaction(context.Background(), rawTx)
			if err != nil {
				return fmt.Errorf("send tx error, err: %v", err)
			}
			fmt.Printf("create multisig tx account: %s Transaction txHash: %s\n", multisigTxAccountPubkey.ToBase58(), txHash)
			time.Sleep(time.Second * 2)

		} else {
			fmt.Printf("multisig tx Account: %s exist on chain and will not create\n", multisigTxAccountPubkey.ToBase58())
		}

		//other fee account approve
		for i := 0; i < len(otherFeeAccount); i++ {
			res, err := c.GetRecentBlockhash(context.Background())
			if err != nil {
				return fmt.Errorf("get recent block hash error, err: %v", err)
			}
			remainingAccounts := multisigprog.GetRemainAccounts([]solTypes.Instruction{stakeInstruction})
			rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
				Instructions: []solTypes.Instruction{
					multisigprog.Approve(
						MultisigProgramId,
						MultisigInfoAccount.PublicKey,
						multisignerPubkey,
						multisigTxAccountPubkey,
						otherFeeAccount[i].PublicKey,
						remainingAccounts,
					),
				},
				Signers:         []solTypes.Account{otherFeeAccount[i], FeeAccount},
				FeePayer:        FeeAccount.PublicKey,
				RecentBlockHash: res.Blockhash,
			})

			if err != nil {
				return fmt.Errorf("generate Approve tx error, err: %v", err)
			}

			txHash, err := c.SendRawTransaction(context.Background(), rawTx)
			if err != nil {
				return fmt.Errorf("send tx error, err: %v", err)
			}
			fmt.Printf("Approve txHash: %s otherfeeAccount: %s\n", txHash, otherFeeAccount[i].PublicKey.ToBase58())
			time.Sleep(time.Second * 5)
		}

		retry = 0
		for {
			if retry > retryLimit {
				return err
			}
			txInfo, err := c.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
			if err != nil {
				fmt.Println("GetMultisigTxAccountInfo failed will retry ...", err)
				retry++
				time.Sleep(3 * time.Second)
				continue
			}

			if txInfo.DidExecute == 1 {
				fmt.Printf("stakeBaseAccount %s init success\n\n", stakeAccoountStr)
				break
			} else {
				fmt.Printf("stakeBaseAccount %s not init yet, waiting...\n", stakeAccoountStr)
				retry++
				time.Sleep(3 * time.Second)
				continue
			}
		}
	}

	fmt.Println("all account init success")
	return nil
}

func getMultisigTxAccountPubkey(baseAccount, programID, stakeBaseAccount solCommon.PublicKey, index int) (solCommon.PublicKey, string) {
	seed := fmt.Sprintf("initAccount:%s:%d", stakeBaseAccount.ToBase58()[:4], index)
	return solCommon.CreateWithSeed(baseAccount, seed, programID), seed
}
