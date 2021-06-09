package solana

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mr-tron/base58"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
	solClient "github.com/tpkeeper/solana-go-sdk/client"
	solCommon "github.com/tpkeeper/solana-go-sdk/common"
	"github.com/tpkeeper/solana-go-sdk/multisigprog"
	"github.com/tpkeeper/solana-go-sdk/stakeprog"
	"github.com/tpkeeper/solana-go-sdk/sysprog"
	solTypes "github.com/tpkeeper/solana-go-sdk/types"
)

func (w *writer) processActiveReportedEvent(m *core.Message) bool {
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m, errors.New("msg cast to MultiEventFlow not ok"))
		return false
	}

	flow, ok := mef.EventData.(*submodel.ActiveReportedFlow)
	if !ok {
		w.log.Error("processActiveReportedEvent eventData is not TransferFlow")
		return false
	}

	poolAddrBase58Str := base58.Encode(flow.Snap.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrBase58Str)
	if err != nil {
		w.log.Error("processBondReportEvent failed",
			"pool address", poolAddrBase58Str,
			"error", err)
		return false
	}

	currentEra := w.conn.GetCurrentEra()
	rpcClient := poolClient.GetRpcClient()
	//get derived account
	canWithdrawAccounts := make(map[solCommon.PublicKey]*solClient.StakeAccountRsp, 0)
	for i := uint32(0); i < 10; i++ {
		stakeAccountPubkey, _ := GetStakeAccountPubkey(poolClient.StakeBaseAccount.PublicKey, currentEra-i)
		accountInfo, err := rpcClient.GetStakeAccountInfo(context.Background(), stakeAccountPubkey.ToBase58())
		if err != nil {
			if err == solClient.ErrAccountNotFound {
				continue
			} else {
				w.log.Error("processBondReportEvent GetStakeAccountInfo failed",
					"pool  address", poolAddrBase58Str,
					"stake account", stakeAccountPubkey.ToBase58(),
					"error", err)
				return false
			}
		}

		//filter account that has cool down
		if accountInfo.StakeAccount.Type == 1 {
			canWithdrawAccounts[stakeAccountPubkey] = accountInfo
		}
	}
	//no need withdraw
	if len(canWithdrawAccounts) == 0 {
		w.log.Info("processActiveReportedEvent no need withdraw Tx",
			"pool address", poolAddrBase58Str,
			"era", flow.Snap.Era,
			"snapId", flow.ShotId)

		callHash := utils.BlakeTwo256(flow.Snap.Pool)
		mef.OpaqueCalls = []*submodel.MultiOpaqueCall{
			{
				CallHash: hexutil.Encode(callHash[:])}}

		return w.informChain(m.Destination, m.Source, mef)
	}

	miniMumBalanceForTx, err := rpcClient.GetMinimumBalanceForRentExemption(context.Background(), 1000)
	if err != nil {
		w.log.Error("processEraPoolUpdatedEvt GetMinimumBalanceForRentExemption failed",
			"pool address", poolAddrBase58Str,
			"err", err)
		return false
	}
	//create multisig withdraw tx account
	multisigTxAccountPubkey, multisigTxAccountSeed := GetMultisigTxAccountPubkey(
		poolClient.MultisigTxBaseAccount.PublicKey,
		w.multisigProgId,
		MultisigTxWithdrawType,
		flow.Snap.Era)

	withdrawInstructions := make([]solTypes.Instruction, 0)
	multisigTxAccountInfo, err := rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
	if err != nil {
		if err == solClient.ErrAccountNotFound {
			res, err := rpcClient.GetRecentBlockhash(context.Background())
			if err != nil {
				w.log.Error("processActiveReportedEvent GetRecentBlockhash failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}

			//send from  relayers
			//create multisig withdraw transaction account of this era
			programIds := make([]solCommon.PublicKey, 0)
			accountMetas := make([][]solTypes.AccountMeta, 0)
			txDatas := make([][]byte, 0)
			for stakeAccountPubkey, accountInfo := range canWithdrawAccounts {
				withdrawInstruction := stakeprog.Withdraw(stakeAccountPubkey, poolClient.MultisignerPubkey,
					poolClient.MultisignerPubkey, accountInfo.Lamports, solCommon.PublicKey{})

				withdrawInstructions = append(withdrawInstructions, withdrawInstruction)

				programIds = append(programIds, withdrawInstruction.ProgramID)
				accountMetas = append(accountMetas, withdrawInstruction.Accounts)
				txDatas = append(txDatas, withdrawInstruction.Data)
			}

			rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
				Instructions: []solTypes.Instruction{
					sysprog.CreateAccountWithSeed(
						poolClient.FeeAccount.PublicKey,
						multisigTxAccountPubkey,
						poolClient.MultisigTxBaseAccount.PublicKey,
						w.multisigProgId,
						multisigTxAccountSeed,
						miniMumBalanceForTx,
						1000,
					),
					multisigprog.CreateTransaction(
						w.multisigProgId,
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
				w.log.Error("processActiveReportedEvent CreateTransaction CreateRawTransaction failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}

			txHash, err := rpcClient.SendRawTransaction(context.Background(), rawTx)
			if err != nil {
				w.log.Error("processActiveReportedEvent createTransaction SendRawTransaction failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}
			w.log.Info("create multisig tx account",
				"tx hash", txHash,
				"multisig tx account", multisigTxAccountPubkey.ToBase58())

		} else {
			w.log.Error("processActiveReportedEvent GetMultisigTxAccountInfo err",
				"pool  address", poolAddrBase58Str,
				"multisig tx account address", multisigTxAccountPubkey.ToBase58(),
				"err", err)
			return false
		}
	}
	//no need approve
	if multisigTxAccountInfo.DidExecute == 1 {
		callHash := utils.BlakeTwo256(flow.Snap.Pool)
		mef.OpaqueCalls = []*submodel.MultiOpaqueCall{
			{CallHash: hexutil.Encode(callHash[:])}}

		return w.informChain(m.Destination, m.Source, mef)
	}

	remainingAccounts := multisigprog.GetRemainAccounts(withdrawInstructions)

	res, err := rpcClient.GetRecentBlockhash(context.Background())
	if err != nil {
		w.log.Error("processActiveReportedEvent GetRecentBlockhash failed",
			"pool address", poolAddrBase58Str,
			"err", err)
		return false
	}
	rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
		Instructions: []solTypes.Instruction{
			multisigprog.Approve(
				w.multisigProgId,
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
		w.log.Error("processActiveReportedEvent approve CreateRawTransaction failed",
			"pool address", poolAddrBase58Str,
			"err", err)
		return false
	}

	// t.Log("rawtx base58:", base58.Encode(rawTx))
	txHash, err := rpcClient.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		w.log.Error("processActiveReportedEvent approve SendRawTransaction failed",
			"pool address", poolAddrBase58Str,
			"err", err)
		return false
	}

	w.log.Info("approve multisig tx account",
		"tx hash", txHash,
		"multisig tx account", multisigTxAccountPubkey.ToBase58())

	callHash := utils.BlakeTwo256(flow.Snap.Pool)
	mef.OpaqueCalls = []*submodel.MultiOpaqueCall{
		{CallHash: hexutil.Encode(callHash[:])}}

	return w.informChain(m.Destination, m.Source, mef)

}
