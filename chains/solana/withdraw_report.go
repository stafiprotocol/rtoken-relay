package solana

import (
	"bytes"
	"context"
	"errors"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mr-tron/base58"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
	solCommon "github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/multisigprog"
	"github.com/stafiprotocol/solana-go-sdk/sysprog"
	solTypes "github.com/stafiprotocol/solana-go-sdk/types"
)

var (
	multisigSendMaxNumber = 5
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

	for i := 0; i <= len(flow.Receives)/multisigSendMaxNumber; i++ {
		multisigTxAccountPubkey, multisigTxAccountSeed := GetMultisigTxAccountPubkeyForTransfer(
			poolClient.MultisigTxBaseAccountPubkey,
			poolClient.MultisigProgramId,
			flow.Snap.Era,
			i)

		transferInstructions := make([]solTypes.Instruction, 0)
		programIds := make([]solCommon.PublicKey, 0)
		accountMetas := make([][]solTypes.AccountMeta, 0)
		txDatas := make([][]byte, 0)

		for j := 0; j < multisigSendMaxNumber; j++ {
			index := i*multisigSendMaxNumber + j
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
			w.log.Info("will transfer to ", "index ", index, " addr ", to.ToBase58(), " value ", value.Int64())
		}
		remainingAccounts := multisigprog.GetRemainAccounts(transferInstructions)

		if poolClient.HasBaseAccountAuth {
			if poolClient.MultisigTxBaseAccount == nil {
				w.log.Error("MultisigTxBaseAccount privkey not exist", "MultisigTxBaseAccount", poolClient.MultisigTxBaseAccountPubkey)
				return false
			}

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
		}
		//check multisig tx account is created
		create := w.waitingForMultisigTxCreate(rpcClient, poolAddrBase58Str, multisigTxAccountPubkey.ToBase58(), "processWithdrawReportedEvent")
		if !create {
			return false
		}
		w.log.Info("processWithdrawReportedEvent multisigTxAccount has create", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())

		valid := w.CheckMultisigTx(rpcClient, multisigTxAccountPubkey, programIds, accountMetas, txDatas)
		if !valid {
			w.log.Info("processWithdrawReportedEvent CheckMultisigTx failed", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())
			return false
		}

		//if has exe just continue
		isExe := w.IsMultisigTxExe(rpcClient, multisigTxAccountPubkey)
		if isExe {
			w.log.Info("processWithdrawReportedEvent multisigTxAccount has execute", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())
			continue
		}
		//approve tx
		send := w.approveMultisigTx(rpcClient, poolClient, poolAddrBase58Str, multisigTxAccountPubkey, remainingAccounts, "processWithdrawReportedEvent")
		if !send {
			return false
		}

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
