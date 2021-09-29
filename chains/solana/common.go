package solana

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	subTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/solana"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
	solCommon "github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/multisigprog"
	"github.com/stafiprotocol/solana-go-sdk/stakeprog"
	"github.com/stafiprotocol/solana-go-sdk/sysprog"
	solTypes "github.com/stafiprotocol/solana-go-sdk/types"
)

var retryLimit = 50
var waitTime = time.Second * 5
var backCheckLen = 6

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

func GetMultisigTxAccountPubkey(baseAccount, programID solCommon.PublicKey, txType MultisigTxType, era uint32, stakeBaseAccountIndex int) (solCommon.PublicKey, string) {
	seed := fmt.Sprintf("multisig:%s:%d:%d", txType, era, stakeBaseAccountIndex)
	return solCommon.CreateWithSeed(baseAccount, seed, programID), seed
}

func GetMultisigTxAccountPubkeyForTransfer(baseAccount, programID solCommon.PublicKey, era uint32, batchTimes int) (solCommon.PublicKey, string) {
	seed := fmt.Sprintf("multisig:%s:%d:%d", MultisigTxTransferType, era, batchTimes)
	return solCommon.CreateWithSeed(baseAccount, seed, programID), seed
}

func GetStakeAccountPubkey(baseAccount solCommon.PublicKey, era uint32) (solCommon.PublicKey, string) {
	seed := fmt.Sprintf("stake:%d", era)
	return solCommon.CreateWithSeed(baseAccount, seed, solCommon.StakeProgramID), seed
}

//1 get stake derived accounts which state is active and merge to base account
//2 get stake derived accounts which state is inactive and withdraw to pool address
func (w *writer) MergeAndWithdraw(poolClient *solana.PoolClient,
	poolAddrBase58Str string, currentEra uint32, shotId subTypes.Hash, pool []byte) bool {
	rpcClient := poolClient.GetRpcClient()

	// must deal every stakeBaseAccounts
	for stakeBaseAccountIndex, useStakeBaseAccount := range poolClient.StakeBaseAccounts {

		//get derived account
		canWithdrawAccounts := make(map[solCommon.PublicKey]solClient.GetStakeActivationResponse)
		canMergeAccounts := make(map[solCommon.PublicKey]solClient.GetStakeActivationResponse)
		stakeBaseAccountInfo, err := rpcClient.GetStakeAccountInfo(context.Background(), useStakeBaseAccount.PublicKey.ToBase58())
		if err != nil {
			w.log.Error("MergeAndWithdraw GetStakeAccountInfo failed",
				"pool  address", poolAddrBase58Str,
				"stake account", useStakeBaseAccount.PublicKey.ToBase58(),
				"error", err)
			return false
		}
		creditsStakeBaseAccount := stakeBaseAccountInfo.StakeAccount.Info.Stake.CreditsObserved
		for i := uint32(0); i < uint32(backCheckLen); i++ {
			stakeAccountPubkey, _ := GetStakeAccountPubkey(useStakeBaseAccount.PublicKey, currentEra-i)
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
				//withdraw all balance
				accountInfo.Inactive = stakeAccountInfo.Lamports
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

		//create multisig withdraw tx account
		multisigTxAccountPubkey, multisigTxAccountSeed := GetMultisigTxAccountPubkey(
			poolClient.MultisigTxBaseAccount.PublicKey,
			poolClient.MultisigProgramId,
			MultisigTxWithdrawType,
			currentEra,
			stakeBaseAccountIndex)

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
				useStakeBaseAccount.PublicKey,
				stakeAccountPubkey,
				poolClient.MultisignerPubkey)

			withdrawAndMergeInstructions = append(withdrawAndMergeInstructions, mergeInstruction)

			programIds = append(programIds, mergeInstruction.ProgramID)
			accountMetas = append(accountMetas, mergeInstruction.Accounts)
			txDatas = append(txDatas, mergeInstruction.Data)
		}
		remainingAccounts := multisigprog.GetRemainAccounts(withdrawAndMergeInstructions)

		_, err = rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
		if err != nil && err == solClient.ErrAccountNotFound {
			sendOk := w.createMultisigTxAccount(rpcClient, poolClient, poolAddrBase58Str, programIds, accountMetas, txDatas,
				multisigTxAccountPubkey, multisigTxAccountSeed, "MergeAndWithdraw")
			if !sendOk {
				return false
			}
		}
		if err != nil && err != solClient.ErrAccountNotFound {
			w.log.Error("MergeAndWithdraw GetMultisigTxAccountInfo err",
				"pool  address", poolAddrBase58Str,
				"multisig tx account address", multisigTxAccountPubkey.ToBase58(),
				"err", err)
			return false
		}

		//check multisig tx account is created
		create := w.waitingForMultisigTxCreate(rpcClient, poolAddrBase58Str, multisigTxAccountPubkey.ToBase58(), "MergeAndWithdraw")
		if !create {
			return false
		}
		w.log.Info("MergeAndWithdraw multisigTxAccount has create", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())

		valid := w.CheckMultisigTx(rpcClient, multisigTxAccountPubkey, programIds, accountMetas, txDatas)
		if !valid {
			w.log.Info("MergeAndWithdraw CheckMultisigTx failed", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())
			return false
		}
		//if has exe just return
		isExe := w.IsMultisigTxExe(rpcClient, multisigTxAccountPubkey)
		if isExe {
			w.log.Info("MergeAndWithdraw multisigTxAccount has execute", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())
			return true
		}
		//approve multisig tx
		send := w.approveMultisigTx(rpcClient, poolClient, poolAddrBase58Str, multisigTxAccountPubkey, remainingAccounts, "MergeAndWithdraw")
		if !send {
			return false
		}

		//check multisig exe result
		exe := w.waitingForMultisigTxExe(rpcClient, poolAddrBase58Str, multisigTxAccountPubkey.ToBase58(), "MergeAndWithdraw")
		if !exe {
			return false
		}
		w.log.Info("MergeAndWithdraw multisigTxAccount has execute", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())
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

func (w *writer) waitingForMultisigTxExe(rpcClient *solClient.Client, poolAddress, multisigTxAddress, processName string) bool {
	retry := 0
	for {
		if retry >= retryLimit {
			w.log.Error(fmt.Sprintf("[%s] GetMultisigTxAccountInfo reach retry limit", processName),
				"pool  address", poolAddress,
				"multisig tx account address", multisigTxAddress)
			return false
		}
		multisigTxAccountInfo, err := rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAddress)
		if err == nil && multisigTxAccountInfo.DidExecute == 1 {
			break
		} else {
			w.log.Warn(fmt.Sprintf("[%s] multisigTxAccount not execute yet, waiting...", processName),
				"pool  address", poolAddress,
				"multisig tx Account", multisigTxAddress)
			time.Sleep(waitTime)
			retry++
			continue
		}
	}
	return true
}

func (w *writer) waitingForMultisigTxCreate(rpcClient *solClient.Client, poolAddress, multisigTxAddress, processName string) bool {
	retry := 0
	for {
		if retry >= retryLimit {
			w.log.Error(fmt.Sprintf("[%s] GetMultisigTxAccountInfo reach retry limit", processName),
				"pool  address", poolAddress,
				"multisig tx account address", multisigTxAddress)
			return false
		}
		_, err := rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAddress)
		if err != nil {
			w.log.Warn(fmt.Sprintf("[%s] GetMultisigTxAccountInfo failed, waiting...", processName),
				"pool  address", poolAddress,
				"multisig tx account", multisigTxAddress,
				"err", err)
			time.Sleep(waitTime)
			retry++
			continue
		} else {
			break
		}
	}
	return true
}

func (w *writer) waitingForStakeAccountCreate(rpcClient *solClient.Client, poolAddress, stakeAccountAddress, processName string) bool {
	retry := 0
	for {
		if retry >= retryLimit {
			w.log.Error(fmt.Sprintf("[%s] GetStakeAccountInfo reach retry limit", processName),
				"pool  address", poolAddress,
				"stake address", stakeAccountAddress)
			return false
		}
		_, err := rpcClient.GetStakeAccountInfo(context.Background(), stakeAccountAddress)
		if err != nil {
			w.log.Warn(fmt.Sprintf("[%s] GetStakeAccountInfo failed, waiting...", processName),
				"pool  address", poolAddress,
				"stake address", stakeAccountAddress,
				"err", err)
			time.Sleep(waitTime)
			retry++
			continue
		} else {
			break
		}
	}
	return true
}

func (w *writer) createMultisigTxAccount(
	rpcClient *solClient.Client,
	poolClient *solana.PoolClient,
	poolAddress string,
	programsIds []solCommon.PublicKey,
	accountMetas [][]solTypes.AccountMeta,
	datas [][]byte,
	multisigTxAccountPubkey solCommon.PublicKey,
	multisigTxAccountSeed string,
	processName string,
) bool {
	res, err := rpcClient.GetRecentBlockhash(context.Background())
	if err != nil {
		w.log.Error(fmt.Sprintf("[%s] GetRecentBlockhash failed", processName),
			"pool address", poolAddress,
			"err", err)
		return false
	}
	miniMumBalanceForTx, err := rpcClient.GetMinimumBalanceForRentExemption(context.Background(), 1000)
	if err != nil {
		w.log.Error(fmt.Sprintf("[%s] GetMinimumBalanceForRentExemption failed", processName),
			"pool address", poolAddress,
			"err", err)
		return false
	}
	miniMumBalanceForTx *= 2
	//send from one relayers
	//create multisig tx account of this era
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
				programsIds,
				accountMetas,
				datas,
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
		w.log.Error(fmt.Sprintf("[%s] CreateTransaction CreateRawTransaction failed", processName),
			"pool address", poolAddress,
			"err", err)
		return false
	}

	txHash, err := rpcClient.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		w.log.Error(fmt.Sprintf("[%s] createTransaction SendRawTransaction failed", processName),
			"pool address", poolAddress,
			"err", err)
		return false
	}
	w.log.Info(fmt.Sprintf("[%s] create multisig tx account  has send", processName),
		"tx hash", txHash,
		"multisig tx account", multisigTxAccountPubkey.ToBase58())
	return true
}

func (w *writer) approveMultisigTx(
	rpcClient *solClient.Client,
	poolClient *solana.PoolClient,
	poolAddress string,
	multisigTxAccountPubkey solCommon.PublicKey,
	remainingAccounts []solTypes.AccountMeta,
	processName string) bool {
	res, err := rpcClient.GetRecentBlockhash(context.Background())
	if err != nil {
		w.log.Error(fmt.Sprintf("[%s] GetRecentBlockhash failed", processName),
			"pool address", poolAddress,
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
		w.log.Error(fmt.Sprintf("[%s] approve CreateRawTransaction failed", processName),
			"pool address", poolAddress,
			"err", err)
		return false
	}

	txHash, err := rpcClient.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		w.log.Error(fmt.Sprintf("[%s] approve SendRawTransaction failed", processName),
			"pool address", poolAddress,
			"err", err)
		return false
	}

	w.log.Info(fmt.Sprintf("[%s] approve multisig tx account has send", processName),
		"tx hash", txHash,
		"multisig tx account", multisigTxAccountPubkey.ToBase58())

	return true
}

func (w *writer) IsMultisigTxExe(
	rpcClient *solClient.Client,
	multisigTxAccountPubkey solCommon.PublicKey) bool {
	accountInfo, err := rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
	if err == nil && accountInfo.DidExecute == 1 {
		return true
	}
	return false
}

func (w *writer) CheckMultisigTx(
	rpcClient *solClient.Client,
	multisigTxAccountPubkey solCommon.PublicKey,
	programsIds []solCommon.PublicKey,
	accountMetas [][]solTypes.AccountMeta,
	datas [][]byte) bool {
	accountInfo, err := rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
	if err == nil {
		thisProgramsIdsBts, err := solCommon.SerializeData(programsIds)
		if err != nil {
			return false
		}
		thisAccountMetasBts, err := solCommon.SerializeData(accountMetas)
		if err != nil {
			return false
		}
		thisDatasBts, err := solCommon.SerializeData(datas)
		if err != nil {
			return false
		}
		onchainProgramsIdsBts, err := solCommon.SerializeData(accountInfo.ProgramID)
		if err != nil {
			return false
		}
		onchainAccountMetasBts, err := solCommon.SerializeData(accountInfo.Accounts)
		if err != nil {
			return false
		}
		onchainDatasBts, err := solCommon.SerializeData(accountInfo.Data)
		if err != nil {
			return false
		}
		if bytes.Equal(thisProgramsIdsBts, onchainProgramsIdsBts) &&
			bytes.Equal(thisAccountMetasBts, onchainAccountMetasBts) &&
			bytes.Equal(thisDatasBts, onchainDatasBts) {
			return true
		}
		w.log.Error("CheckMultisigTx not equal ",
			"thisprogramsIds", hex.EncodeToString(thisProgramsIdsBts),
			"onchainProgramnsIdsBts", hex.EncodeToString(onchainProgramsIdsBts),
			"thisAccountMetasBts", hex.EncodeToString(thisAccountMetasBts),
			"onchainAccountMetasBts", hex.EncodeToString(onchainAccountMetasBts),
			"thisDatasBts", hex.EncodeToString(thisDatasBts),
			"onchainDatasBts", hex.EncodeToString(onchainDatasBts))
	}

	return false
}

func (w *writer) CheckStakeAccount(rpcClient *solClient.Client, stakeAccount, multisigner string) bool {
	stakeAccountInfo, err := rpcClient.GetStakeAccountInfo(context.Background(), stakeAccount)
	if err != nil {
		w.log.Error("CheckStakeAccount failed", "err", err)
		return false
	}
	if strings.EqualFold(stakeAccountInfo.StakeAccount.Info.Meta.Authorized.Staker.ToBase58(), multisigner) &&
		strings.EqualFold(stakeAccountInfo.StakeAccount.Info.Meta.Authorized.Withdrawer.ToBase58(), multisigner) {
		return true
	}
	w.log.Error("CheckStakeAccount failed",
		"multisigner", multisigner,
		"staker", stakeAccountInfo.StakeAccount.Info.Meta.Authorized.Staker.ToBase58(),
		"withdrawer", stakeAccountInfo.StakeAccount.Info.Meta.Authorized.Withdrawer.ToBase58())

	return false
}
