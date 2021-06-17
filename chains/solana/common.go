package solana

import (
	"context"
	"fmt"
	"strings"
	"time"

	subTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/solana"
	solClient "github.com/tpkeeper/solana-go-sdk/client"
	solCommon "github.com/tpkeeper/solana-go-sdk/common"
	"github.com/tpkeeper/solana-go-sdk/multisigprog"
	"github.com/tpkeeper/solana-go-sdk/stakeprog"
	"github.com/tpkeeper/solana-go-sdk/sysprog"
	solTypes "github.com/tpkeeper/solana-go-sdk/types"
)

func (w *writer) printContentError(m *core.Message, err error) {
	w.log.Error("msg resolve failed", "source", m.Source, "dest", m.Destination, "reason", m.Reason, "err", err)
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (w *writer) submitMessage(m *core.Message) bool {
	err := w.router.Send(m)
	if err != nil {
		w.log.Error("failed to process event", "err", err)
		return false
	}

	return true
}

func (w *writer) informChain(source, dest core.RSymbol, flow *submodel.MultiEventFlow) bool {
	msg := &core.Message{Source: source, Destination: dest, Reason: core.InformChain, Content: flow}
	return w.submitMessage(msg)
}

func (w *writer) activeReport(source, dest core.RSymbol, flow *submodel.BondReportedFlow) bool {
	msg := &core.Message{Source: source, Destination: dest, Reason: core.ActiveReport, Content: flow}
	return w.submitMessage(msg)
}

type MultisigTxType string

var MultisigTxStakeType = MultisigTxType("stake")
var MultisigTxUnStakeType = MultisigTxType("unstake")
var MultisigTxWithdrawType = MultisigTxType("withdraw")
var MultisigTxTransferType = MultisigTxType("transfer")

func GetMultisigTxAccountPubkey(baseAccount, programID solCommon.PublicKey, txType MultisigTxType, era uint32) (solCommon.PublicKey, string) {
	seed := fmt.Sprintf("multisig2:%s:%d", txType, era)
	return solCommon.CreateWithSeed(baseAccount, seed, programID), seed
}

func GetStakeAccountPubkey(baseAccount solCommon.PublicKey, era uint32) (solCommon.PublicKey, string) {
	seed := fmt.Sprintf("stake:%d", era)
	return solCommon.CreateWithSeed(baseAccount, seed, solCommon.StakeProgramID), seed
}

func (w *writer) MergeAndWithdraw(poolClient *solana.PoolClient,
	poolAddrBase58Str string, currentEra uint32, shotId subTypes.Hash, pool []byte) bool {
	rpcClient := poolClient.GetRpcClient()
	//get derived account
	canWithdrawAccounts := make(map[solCommon.PublicKey]solClient.GetStakeActivationResponse)
	canMergeAccounts := make(map[solCommon.PublicKey]solClient.GetStakeActivationResponse)
	stakeBaseAccountInfo, err := rpcClient.GetStakeAccountInfo(context.Background(), poolClient.StakeBaseAccount.PublicKey.ToBase58())
	if err != nil {
		w.log.Error("MergeAndWithdraw GetStakeAccountInfo failed",
			"pool  address", poolAddrBase58Str,
			"stake account", poolClient.StakeBaseAccount.PublicKey.ToBase58(),
			"error", err)
		return false
	}
	creditsStakeBaseAccount := stakeBaseAccountInfo.StakeAccount.Info.Stake.CreditsObserved
	for i := uint32(0); i < 10; i++ {
		stakeAccountPubkey, _ := GetStakeAccountPubkey(poolClient.StakeBaseAccount.PublicKey, currentEra-i)
		accountInfo, err := rpcClient.GetStakeActivation(
			context.Background(),
			stakeAccountPubkey.ToBase58(),
			solClient.GetStakeActivationConfig{})

		if err != nil {
			if strings.Contains(err.Error(), "account not found") {
				continue
			} else {
				w.log.Error("MergeAndWithdraw GetStakeAccountInfo failed",
					"pool  address", poolAddrBase58Str,
					"stake account", stakeAccountPubkey.ToBase58(),
					"error", err)
				return false
			}
		}
		//filter credite observed not equal to stakeBaseAccount
		stakeAccountInfo, err := rpcClient.GetStakeAccountInfo(context.Background(), stakeAccountPubkey.ToBase58())
		if err != nil {
			w.log.Error("MergeAndWithdraw GetStakeAccountInfo failed",
				"pool  address", poolAddrBase58Str,
				"stake account", stakeAccountPubkey.ToBase58(),
				"error", err)
			return false
		}
		if stakeAccountInfo.StakeAccount.Info.Stake.CreditsObserved != creditsStakeBaseAccount {
			continue
		}

		//filter account
		if accountInfo.State == solClient.StakeActivationStateInactive {
			canWithdrawAccounts[stakeAccountPubkey] = accountInfo
		} else if accountInfo.State == solClient.StakeActivationStateActive {
			canMergeAccounts[stakeAccountPubkey] = accountInfo
		}
	}

	//no need withdraw,just report to stafi
	if len(canWithdrawAccounts) == 0 && len(canMergeAccounts) == 0 {
		w.log.Info("MergeAndWithdraw no need merge and withdraw ",
			"pool address", poolAddrBase58Str,
			"era", currentEra,
			"snapId", shotId.Hex())

		return true
	}
	w.log.Info("canWithdrawAccounts", "accouts", mapToString(canWithdrawAccounts))
	w.log.Info("canMergeAccounts", "accounts", mapToString(canMergeAccounts))

	// miniMumBalanceForTx, err := rpcClient.GetMinimumBalanceForRentExemption(context.Background(), 1000)
	// if err != nil {
	// 	w.log.Error("MergeAndWithdraw GetMinimumBalanceForRentExemption failed",
	// 		"pool address", poolAddrBase58Str,
	// 		"err", err)
	// 	return false
	// }
	miniMumBalanceForTx := uint64(1e9)

	//create multisig withdraw tx account
	multisigTxAccountPubkey, multisigTxAccountSeed := GetMultisigTxAccountPubkey(
		poolClient.MultisigTxBaseAccount.PublicKey,
		poolClient.MultisigProgramId,
		MultisigTxWithdrawType,
		currentEra)

	withdrawAndMergeInstructions := make([]solTypes.Instruction, 0)

	programIds := make([]solCommon.PublicKey, 0)
	accountMetas := make([][]solTypes.AccountMeta, 0)
	txDatas := make([][]byte, 0)

	for stakeAccountPubkey, accountInfo := range canWithdrawAccounts {
		withdrawInstruction := stakeprog.Withdraw(stakeAccountPubkey, poolClient.MultisignerPubkey,
			poolClient.MultisignerPubkey, accountInfo.Inactive, solCommon.PublicKey{})

		withdrawAndMergeInstructions = append(withdrawAndMergeInstructions, withdrawInstruction)

		programIds = append(programIds, withdrawInstruction.ProgramID)
		accountMetas = append(accountMetas, withdrawInstruction.Accounts)
		txDatas = append(txDatas, withdrawInstruction.Data)
	}

	for stakeAccountPubkey, _ := range canMergeAccounts {
		mergeInstruction := stakeprog.Merge(
			poolClient.StakeBaseAccount.PublicKey,
			stakeAccountPubkey,
			poolClient.MultisignerPubkey)

		withdrawAndMergeInstructions = append(withdrawAndMergeInstructions, mergeInstruction)

		programIds = append(programIds, mergeInstruction.ProgramID)
		accountMetas = append(accountMetas, mergeInstruction.Accounts)
		txDatas = append(txDatas, mergeInstruction.Data)
	}

	multisigTxAccountInfo, err := rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
	if err != nil {
		if err == solClient.ErrAccountNotFound {
			res, err := rpcClient.GetRecentBlockhash(context.Background())
			if err != nil {
				w.log.Error("MergeAndWithdraw GetRecentBlockhash failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}

			//send from  relayers
			//create multisig withdraw transaction account of this era
			rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
				Instructions: []solTypes.Instruction{
					sysprog.CreateAccountWithSeed(
						poolClient.FeeAccount.PublicKey,
						multisigTxAccountPubkey,
						poolClient.MultisigTxBaseAccount.PublicKey,
						poolClient.MultisigProgramId,
						multisigTxAccountSeed,
						miniMumBalanceForTx,
						1000,
					),
					multisigprog.CreateTransaction(
						poolClient.MultisigProgramId,
						programIds,
						accountMetas,
						txDatas,
						poolClient.MultisigInfoPubkey,
						multisigTxAccountPubkey,
						poolClient.FeeAccount.PublicKey,
					),
				},
				Signers:         []solTypes.Account{poolClient.FeeAccount, poolClient.MultisigTxBaseAccount},
				FeePayer:        poolClient.FeeAccount.PublicKey,
				RecentBlockHash: res.Blockhash,
			})

			if err != nil {
				w.log.Error("MergeAndWithdraw CreateTransaction CreateRawTransaction failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}

			txHash, err := rpcClient.SendRawTransaction(context.Background(), rawTx)
			if err != nil {
				w.log.Error("MergeAndWithdraw createTransaction SendRawTransaction failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}
			w.log.Info("create multisig tx account",
				"tx hash", txHash,
				"multisig tx account", multisigTxAccountPubkey.ToBase58())

		} else {
			w.log.Error("MergeAndWithdraw GetMultisigTxAccountInfo err",
				"pool  address", poolAddrBase58Str,
				"multisig tx account address", multisigTxAccountPubkey.ToBase58(),
				"err", err)
			return false
		}
	}
	//no need approve
	if multisigTxAccountInfo != nil && multisigTxAccountInfo.DidExecute == 1 {
		w.log.Info("MergeAndWithdraw no need approve multisig tx account",
			"multisig tx account", multisigTxAccountPubkey.ToBase58())
		return true
	}

	//approve multisig tx
	remainingAccounts := multisigprog.GetRemainAccounts(withdrawAndMergeInstructions)
	res, err := rpcClient.GetRecentBlockhash(context.Background())
	if err != nil {
		w.log.Error("MergeAndWithdraw GetRecentBlockhash failed",
			"pool address", poolAddrBase58Str,
			"err", err)
		return false
	}
	rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
		Instructions: []solTypes.Instruction{
			multisigprog.Approve(
				poolClient.MultisigProgramId,
				poolClient.MultisigInfoPubkey,
				poolClient.MultisignerPubkey,
				multisigTxAccountPubkey,
				poolClient.FeeAccount.PublicKey,
				remainingAccounts,
			),
		},
		Signers:         []solTypes.Account{poolClient.FeeAccount},
		FeePayer:        poolClient.FeeAccount.PublicKey,
		RecentBlockHash: res.Blockhash,
	})

	if err != nil {
		w.log.Error("MergeAndWithdraw approve CreateRawTransaction failed",
			"pool address", poolAddrBase58Str,
			"err", err)
		return false
	}

	txHash, err := rpcClient.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		w.log.Error("MergeAndWithdraw approve SendRawTransaction failed",
			"pool address", poolAddrBase58Str,
			"err", err)
		return false
	}

	w.log.Info("MergeAndWithdraw approve multisig tx account",
		"tx hash", txHash,
		"multisig tx account", multisigTxAccountPubkey.ToBase58())

	exe := false
	for i := 0; i < 30; i++ {
		multisigTxAccountInfo, err := rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
		if err == nil && multisigTxAccountInfo.DidExecute == 1 {
			exe = true
			break
		}
		w.log.Info("multisigTxAccount not execute", "multisigTxAccount", multisigTxAccountInfo)
		time.Sleep(time.Second * 2)
	}
	if !exe {
		return false
	}
	return true
}

func mapToString(accountsMap map[solCommon.PublicKey]solClient.GetStakeActivationResponse) string {
	ret := ""
	for account, active := range accountsMap {
		ret = ret + account.ToBase58() + fmt.Sprintf(" : %+v", active) + "\n"
	}
	return ret
}
