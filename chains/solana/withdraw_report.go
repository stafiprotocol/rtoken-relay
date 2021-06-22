package solana

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"

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
	if len(flow.Receives) == 0 {
		w.log.Error("processWithdrawReportedEvent Receives len is zero")
		return false
	}
	//sort outPuts for the same rawTx from different relayer
	sort.SliceStable(flow.Receives, func(i, j int) bool {
		return bytes.Compare(flow.Receives[i].Recipient, flow.Receives[j].Recipient) < 0
	})

	poolAddrBase58Str := base58.Encode(flow.Snap.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrBase58Str)
	if err != nil {
		w.log.Error("processWithdrawReportedEvent failed",
			"pool  address", poolAddrBase58Str,
			"error", err)
		return false
	}
	rpcClient := poolClient.GetRpcClient()

	for i := 0; i <= len(flow.Receives)/5; i++ {
		multisigTxAccountPubkey, multisigTxAccountSeed := GetMultisigTxAccountPubkeyForTransfer(
			poolClient.MultisigTxBaseAccount.PublicKey,
			poolClient.MultisigProgramId,
			flow.Snap.Era,
			i)

		transferInstructions := make([]solTypes.Instruction, 0)
		programIds := make([]solCommon.PublicKey, 0)
		accountMetas := make([][]solTypes.AccountMeta, 0)
		txDatas := make([][]byte, 0)

		for j := 0; j < 5; j++ {
			index := i*5 + j
			//check overflow
			if index > len(flow.Receives)-1 {
				break
			}
			receive := flow.Receives[index]
			to := solCommon.PublicKeyFromBytes(receive.Recipient)
			value := big.Int(receive.Value)
			transferInstruction := sysprog.Transfer(poolClient.MultisignerPubkey, to, value.Uint64())
			transferInstructions = append(transferInstructions, transferInstruction)

			programIds = append(programIds, transferInstruction.ProgramID)
			accountMetas = append(accountMetas, transferInstruction.Accounts)
			txDatas = append(txDatas, transferInstruction.Data)
			fmt.Println("will transfer to ", "index ", index, " addr ", to.ToBase58(), " value ", value.Int64())
		}
		remainingAccounts := multisigprog.GetRemainAccounts(transferInstructions)

		_, err = rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
		if err != nil && err == solClient.ErrAccountNotFound {
			sendOk := w.createMultisigTxAccount(rpcClient, poolClient, poolAddrBase58Str, programIds, accountMetas, txDatas,
				multisigTxAccountPubkey, multisigTxAccountSeed, "processWithdrawReportedEvent")
			if !sendOk {
				return false
			}
		}

		if err != nil && err != solClient.ErrAccountNotFound {
			w.log.Error("processWithdrawReportedEvent GetMultisigTxAccountInfo err",
				"pool  address", poolAddrBase58Str,
				"multisig tx account address", multisigTxAccountPubkey.ToBase58(),
				"err", err)
			return false
		}

		//check multisig tx account is created
		create := w.waitingForMultisigTxCreate(rpcClient, poolAddrBase58Str, multisigTxAccountPubkey.ToBase58(), "processWithdrawReportedEvent")
		if !create {
			return false
		}
		w.log.Info("processWithdrawReportedEvent multisigTxAccount has create", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())

		res, err := rpcClient.GetRecentBlockhash(context.Background())
		if err != nil {
			w.log.Error("processWithdrawReportedEvent GetRecentBlockhash failed",
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
			w.log.Error("processWithdrawReportedEvent approve CreateRawTransaction failed",
				"pool address", poolAddrBase58Str,
				"err", err)
			return false
		}

		txHash, err := rpcClient.SendRawTransaction(context.Background(), rawTx)
		if err != nil {
			w.log.Error("processWithdrawReportedEvent approve SendRawTransaction failed",
				"pool address", poolAddrBase58Str,
				"err", err)
			return false
		}

		w.log.Info("processWithdrawReportedEvent approve multisig tx account",
			"tx hash", txHash,
			"multisig tx account", multisigTxAccountPubkey.ToBase58())

		//check multisig exe result
		exe := w.waitingForMultisigTxExe(rpcClient, poolAddrBase58Str, multisigTxAccountPubkey.ToBase58(), "processWithdrawReportedEvent")
		if !exe {
			return false
		}
		w.log.Info("processWithdrawReportedEvent multisigTxAccount has execute", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())
	}

	callHash := utils.BlakeTwo256(flow.Snap.Pool)
	mef.OpaqueCalls = []*submodel.MultiOpaqueCall{
		{CallHash: hexutil.Encode(callHash[:])}}
	return w.informChain(m.Destination, m.Source, mef)
}
