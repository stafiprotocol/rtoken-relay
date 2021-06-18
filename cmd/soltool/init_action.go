package main

import (
	"context"
	"fmt"
	"time"

	"github.com/stafiprotocol/rtoken-relay/shared/solana/vault"
	solClient "github.com/tpkeeper/solana-go-sdk/client"
	solCommon "github.com/tpkeeper/solana-go-sdk/common"
	"github.com/tpkeeper/solana-go-sdk/multisigprog"
	"github.com/tpkeeper/solana-go-sdk/stakeprog"
	"github.com/tpkeeper/solana-go-sdk/sysprog"
	solTypes "github.com/tpkeeper/solana-go-sdk/types"
	"github.com/urfave/cli/v2"
)

func initAction(ctx *cli.Context) error {
	path := ctx.String(configFlag.Name)
	pc := PoolAccounts{}
	err := loadConfig(path, &pc)
	if err != nil {
		return err
	}
	fmt.Printf("accounts info %+v\n\n", pc)
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

	FeeAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.FeeAccount])
	StakeBaseAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.StakeBaseAccount])
	MultisigTxBaseAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.MultisigTxBaseAccount])
	MultisigInfoAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.MultisigInfoPubkey])
	MultisigProgramId := solCommon.PublicKeyFromString(pc.MultisigProgramId)
	ValidatorPubkey := solCommon.PublicKeyFromString(pc.Validator)
	multisignerPubkey, nonce, err := solCommon.FindProgramAddress([][]byte{MultisigInfoAccount.PublicKey.Bytes()}, MultisigProgramId)
	if err != nil {
		return err
	}
	fmt.Println("multisigner:", multisignerPubkey.ToBase58())

	otherFeeAccount := make([]solTypes.Account, 0)
	owners := make([]solCommon.PublicKey, 0)
	owners = append(owners, FeeAccount.PublicKey)
	for _, account := range pc.OtherFeeAccount {
		a := solTypes.AccountFromPrivateKeyBytes(privKeyMap[account])
		otherFeeAccount = append(otherFeeAccount, a)
		owners = append(owners, a.PublicKey)
	}

	c := solClient.NewClient(pc.Endpoint)
	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		return err
	}
	//create stakeBaseAccount
	rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
		Instructions: []solTypes.Instruction{
			sysprog.CreateAccount(
				FeeAccount.PublicKey,
				StakeBaseAccount.PublicKey,
				solCommon.StakeProgramID,
				2000000000,
				200,
			),
			stakeprog.Initialize(
				StakeBaseAccount.PublicKey,
				stakeprog.Authorized{
					Staker:     multisignerPubkey,
					Withdrawer: multisignerPubkey,
				},
				stakeprog.Lockup{},
			),
		},
		Signers:         []solTypes.Account{FeeAccount, StakeBaseAccount},
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
	fmt.Println("createStakeAccount txHash:", txHash)
	time.Sleep(time.Second * 2)

	res, err = c.GetRecentBlockhash(context.Background())
	if err != nil {
		return fmt.Errorf("get recent block hash error, err: %v", err)
	}

	//create multisigInfo account
	rawTx, err = solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
		Instructions: []solTypes.Instruction{
			sysprog.CreateAccount(
				FeeAccount.PublicKey,
				MultisigInfoAccount.PublicKey,
				MultisigProgramId,
				1000000000,
				500,
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
	txHash, err = c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		return fmt.Errorf("send tx error, err: %v", err)
	}
	fmt.Println("createMultisig txHash:", txHash)
	time.Sleep(time.Second * 2)

	//create multisig transaction for stake
	res, err = c.GetRecentBlockhash(context.Background())
	if err != nil {
		return fmt.Errorf("get recent block hash error, err: %v", err)
	}
	stakeInstruction := stakeprog.DelegateStake(StakeBaseAccount.PublicKey, multisignerPubkey, ValidatorPubkey)
	rawTx, err = solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
		Instructions: []solTypes.Instruction{
			sysprog.CreateAccount(
				FeeAccount.PublicKey,
				MultisigTxBaseAccount.PublicKey,
				MultisigProgramId,
				1000000000,
				1000,
			),
			multisigprog.CreateTransaction(
				MultisigProgramId,
				[]solCommon.PublicKey{solCommon.StakeProgramID},
				[][]solTypes.AccountMeta{stakeInstruction.Accounts},
				[][]byte{stakeInstruction.Data},
				MultisigInfoAccount.PublicKey,
				MultisigTxBaseAccount.PublicKey,
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

	txHash, err = c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		return fmt.Errorf("send tx error, err: %v", err)
	}
	fmt.Println("Create Transaction txHash:", txHash)
	time.Sleep(time.Second * 2)

	//other fee account approve
	for i := 0; i < len(otherFeeAccount); i++ {

		res, err = c.GetRecentBlockhash(context.Background())
		if err != nil {
			return fmt.Errorf("get recent block hash error, err: %v", err)
		}
		remainingAccounts := multisigprog.GetRemainAccounts([]solTypes.Instruction{stakeInstruction})

		rawTx, err = solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
			Instructions: []solTypes.Instruction{
				multisigprog.Approve(
					MultisigProgramId,
					MultisigInfoAccount.PublicKey,
					multisignerPubkey,
					MultisigTxBaseAccount.PublicKey,
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

		txHash, err = c.SendRawTransaction(context.Background(), rawTx)
		if err != nil {
			return fmt.Errorf("send tx error, err: %v", err)
		}
		fmt.Printf("Approve txHash: %s otherfeeAccount:%s\n", txHash, otherFeeAccount[i].PublicKey.ToBase58())
		time.Sleep(time.Second * 5)
	}

	for i := 0; i < 10; i++ {
		time.Sleep(5 * time.Second)
		txInfo, err := c.GetMultisigTxAccountInfo(context.Background(), MultisigTxBaseAccount.PublicKey.ToBase58())
		if err != nil {
			fmt.Println("GetMultisigTxAccountInfo failed will retry ...", err)
			continue
		}

		if txInfo.DidExecute == 1 {
			fmt.Println("init success")
			return nil
		} else {
			fmt.Println("not init yet, waiting...")
			continue
		}
	}
	fmt.Println("sorry init failed")
	return nil
}
