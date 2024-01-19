package main

import (
	"context"
	"encoding/hex"
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

func stake(ctx *cli.Context) (err error) {
	defer func() {
		if err != nil {
			fmt.Println("err: ", err)
		}
	}()

	path := ctx.String(configFlag.Name)
	pc := PoolAccountsForMigrate{}
	err = loadConfigForMigrate(path, &pc)
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
	if len(pc.StakeManager) == 0 {
		return fmt.Errorf("stakemanger zero")
	}
	if len(pc.StakePool) == 0 {
		return fmt.Errorf("StakePool zero")
	}
	if len(pc.StakeAccount) == 0 {
		return fmt.Errorf("StakeAccount zero")
	}
	if len(pc.RSolProgramID) == 0 {
		return fmt.Errorf("RSolProgramID zero")
	}

	if _, exist := privKeyMap[pc.StakeAccount]; !exist {
		return fmt.Errorf("StakeAccount not exist")
	}
	FeeAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.FeeAccount])
	MultisigTxBaseAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.MultisigTxBaseAccount])
	MultisigInfoAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.MultisigInfoPubkey])
	stakeAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.StakeAccount])
	MultisigProgramId := solCommon.PublicKeyFromString(pc.MultisigProgramId)
	ValidatorPubkey := solCommon.PublicKeyFromString(pc.Validator)
	// RSolProgramID := solCommon.PublicKeyFromString(pc.RSolProgramID)
	// StakeManager := solCommon.PublicKeyFromString(pc.StakeManager)
	// StakePool := solCommon.PublicKeyFromString(pc.StakePool)

	otherFeeAccount := make([]solTypes.Account, 0)
	for _, account := range pc.OtherFeeAccount {
		a := solTypes.AccountFromPrivateKeyBytes(privKeyMap[account])
		otherFeeAccount = append(otherFeeAccount, a)
	}
	multisignerPubkey, _, err := solCommon.FindProgramAddress(
		[][]byte{MultisigInfoAccount.PublicKey.Bytes()}, MultisigProgramId)
	if err != nil {
		return err
	}
	fmt.Println("multiSigner(old pool):", multisignerPubkey.ToBase58())
	fmt.Println("ValidatorPubkey:", ValidatorPubkey.ToBase58())
	// fmt.Println("stakePool(new pool):", StakePool.ToBase58())

Out:
	for {
		fmt.Println("\ncheck (old pool / new pool) again, then press (y/n) to continue:")
		var input string
		fmt.Scanln(&input)

		switch input {
		case "y":
			break Out
		case "n":
			return nil
		default:
			fmt.Println("press `y` or `n`")
			continue
		}
	}

	c := solClient.NewClient([]string{pc.Endpoint})

	multisigAccountMiniMum, err := c.GetMinimumBalanceForRentExemption(context.Background(),
		solClient.MultisigInfoAccountLengthDefault)
	if err != nil {
		return err
	}

	//create derived multisig tx account if not exist onchain
	multisigTxAccountPubkey, multisigTxAccountSeed := getMultisigTxAccountPubkeyForMigrate(
		MultisigTxBaseAccount.PublicKey,
		MultisigProgramId,
		3)

	// withdrawInstruction := stakeprog.Withdraw(StakeAccount, multisignerPubkey,
	// 	multisignerPubkey, 1594321715540, solCommon.PublicKey{})

	fmt.Printf("stake account: %s,%s", stakeAccount.PublicKey.ToBase58(), hex.EncodeToString(stakeAccount.PrivateKey.Seed()))

	transferInstruction := sysprog.Transfer(multisignerPubkey, stakeAccount.PublicKey, 1594321715540)
	stakeInstruction := stakeprog.DelegateStake(stakeAccount.PublicKey, multisignerPubkey, ValidatorPubkey)

	stakeAccountMiniMum, err := c.GetMinimumBalanceForRentExemption(context.Background(),
		solClient.StakeAccountInfoLengthDefault)
	if err != nil {
		return err
	}
	_, err = c.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
	if err != nil {
		if err != solClient.ErrAccountNotFound {
			return err
		}

		res, err := c.GetLatestBlockhash(context.Background(), solClient.GetLatestBlockhashConfig{
			Commitment: solClient.CommitmentConfirmed,
		})
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
				sysprog.CreateAccount(
					FeeAccount.PublicKey,
					stakeAccount.PublicKey,
					solCommon.StakeProgramID,
					2*stakeAccountMiniMum,
					solClient.StakeAccountInfoLengthDefault,
				),
				stakeprog.Initialize(
					stakeAccount.PublicKey,
					stakeprog.Authorized{
						Staker:     multisignerPubkey,
						Withdrawer: multisignerPubkey,
					},
					stakeprog.Lockup{},
				),
				multisigprog.CreateTransaction(
					MultisigProgramId,
					[]solCommon.PublicKey{transferInstruction.ProgramID, stakeInstruction.ProgramID},
					[][]solTypes.AccountMeta{transferInstruction.Accounts, stakeInstruction.Accounts},
					[][]byte{transferInstruction.Data, stakeInstruction.Data},
					MultisigInfoAccount.PublicKey,
					multisigTxAccountPubkey,
					FeeAccount.PublicKey,
				),
			},
			Signers:         []solTypes.Account{FeeAccount, MultisigTxBaseAccount, stakeAccount},
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
		retry := 0
		for {
			if retry > retryLimit {
				return err
			}
			_, err := c.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
			if err != nil {
				fmt.Println("GetMultisigTxAccountInfo failed will retry ...", err)
				retry++
				time.Sleep(3 * time.Second)
				continue
			}

			fmt.Printf("multisig tx Account create success\n\n")
			break
		}

	} else {
		fmt.Printf("multisig tx Account: %s exist on chain and will not create\n", multisigTxAccountPubkey.ToBase58())
	}

	//other fee account approve
	for i := 0; i < len(otherFeeAccount); i++ {
		res, err := c.GetLatestBlockhash(context.Background(), solClient.GetLatestBlockhashConfig{
			Commitment: solClient.CommitmentConfirmed,
		})
		if err != nil {
			return fmt.Errorf("get recent block hash error, err: %v", err)
		}
		remainingAccounts := multisigprog.GetRemainAccounts([]solTypes.Instruction{transferInstruction, stakeInstruction})
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

	retry := 0
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
			fmt.Printf("stake success\n\n")
			break
		} else {
			fmt.Printf("stake not success yet, waiting...\n")
			retry++
			time.Sleep(3 * time.Second)
			continue
		}
	}

	fmt.Println("stake success")
	return nil
}
