package solana

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mr-tron/base58"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
	solClient "github.com/tpkeeper/solana-go-sdk/client"
	solCommon "github.com/tpkeeper/solana-go-sdk/common"
	"github.com/tpkeeper/solana-go-sdk/multisigprog"
	"github.com/tpkeeper/solana-go-sdk/sysprog"
	solTypes "github.com/tpkeeper/solana-go-sdk/types"
)

func (w *writer) processWithdrawReportedEvent(m *core.Message) bool {
	mef, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m, errors.New("msg cast to MultiEventFlow not ok"))
		return false
	}

	flow, ok := mef.EventData.(*submodel.WithdrawReportedFlow)
	if !ok {
		w.log.Error("processWithdrawReportedEvent eventData is not TransferFlow")
		return false
	}

	poolAddrBase58Str := base58.Encode(flow.Snap.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrBase58Str)
	if err != nil {
		w.log.Error("processWithdrawReportedEvent failed",
			"pool  address", poolAddrBase58Str,
			"error", err)
		return false
	}

	multisigTxAccountPubkey, multisigTxAccountSeed := GetMultisigTxAccountPubkey(
		poolClient.MultisigTxBaseAccount.PublicKey,
		poolClient.MultisigProgramId,
		MultisigTxTransferType,
		flow.Snap.Era)
	rpcClient := poolClient.GetRpcClient()
	transferInstructions := make([]solTypes.Instruction, 0)
	multisigTxAccountInfo, err := rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
	if err != nil {
		if err == solClient.ErrAccountNotFound {
			res, err := rpcClient.GetRecentBlockhash(context.Background())
			if err != nil {
				w.log.Error("processEraPoolUpdatedEvt GetRecentBlockhash failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}
			miniMumBalanceForTx, err := rpcClient.GetMinimumBalanceForRentExemption(context.Background(), 1000)
			if err != nil {
				w.log.Error("processEraPoolUpdatedEvt GetMinimumBalanceForRentExemption failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}

			//send from o relayers
			//create transaction account of this era
			//todo get validator from chain
			programIds := make([]solCommon.PublicKey, 0)
			accountMetas := make([][]solTypes.AccountMeta, 0)
			txDatas := make([][]byte, 0)
			for _, receive := range flow.Receives {
				to := solCommon.PublicKeyFromBytes(receive.Recipient)
				value := big.Int(receive.Value)
				transferInstruction := sysprog.Transfer(poolClient.MultisignerPubkey, to, value.Uint64())
				transferInstructions = append(transferInstructions, transferInstruction)

				programIds = append(programIds, transferInstruction.ProgramID)
				accountMetas = append(accountMetas, transferInstruction.Accounts)
				txDatas = append(txDatas, transferInstruction.Data)

			}

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
				w.log.Error("processEraPoolUpdatedEvt CreateTransaction CreateRawTransaction failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}

			txHash, err := rpcClient.SendRawTransaction(context.Background(), rawTx)
			if err != nil {
				w.log.Error("processEraPoolUpdatedEvt createTransaction SendRawTransaction failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}
			w.log.Info("create multisig tx account",
				"tx hash", txHash,
				"multisig tx account", multisigTxAccountPubkey.ToBase58())

		} else {
			w.log.Error("GetMultisigTxAccountInfo err",
				"pool  address", poolAddrBase58Str,
				"multisig tx account address", multisigTxAccountPubkey.ToBase58(),
				"err", err)
			return false
		}
	}
	if multisigTxAccountInfo.DidExecute == 1 {
		callHash := utils.BlakeTwo256(flow.Snap.Pool)
		mef.OpaqueCalls = []*submodel.MultiOpaqueCall{
			{CallHash: hexutil.Encode(callHash[:])}}
		return w.informChain(m.Destination, m.Source, mef)
	}

	remainingAccounts := multisigprog.GetRemainAccounts(transferInstructions)

	res, err := rpcClient.GetRecentBlockhash(context.Background())
	if err != nil {
		w.log.Error("processEraPoolUpdatedEvt GetRecentBlockhash failed",
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
		w.log.Error("processEraPoolUpdatedEvt approve CreateRawTransaction failed",
			"pool address", poolAddrBase58Str,
			"err", err)
		return false
	}

	// t.Log("rawtx base58:", base58.Encode(rawTx))
	txHash, err := rpcClient.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		w.log.Error("processEraPoolUpdatedEvt approve SendRawTransaction failed",
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
