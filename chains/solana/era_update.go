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
	"github.com/tpkeeper/solana-go-sdk/stakeprog"
	"github.com/tpkeeper/solana-go-sdk/sysprog"
	solTypes "github.com/tpkeeper/solana-go-sdk/types"
)

//process eraPoolUpdate event
//0 check bond/unbond
//  (1)if no need just report bond to stafi
//  (2)if need go next
//1 query multisig tx acount is created
//  (1)if no then created and go next
//  (2)if created go next
//2 approve mutisig tx
//3 query multisig tx executed
//  (1)if executed  report bond result to stafi
//  (2)if reach retry limit return fase
func (w *writer) processEraPoolUpdatedEvt(m *core.Message) bool {
	mFlow, ok := m.Content.(*submodel.MultiEventFlow)
	if !ok {
		w.printContentError(m, errors.New("msg cast to MultisigFlow not ok"))
		return false
	}
	flow, ok := mFlow.EventData.(*submodel.EraPoolUpdatedFlow)
	if !ok {
		w.log.Error("processEraPoolUpdated HeadFlow is not EraPoolUpdatedFlow")
		return false
	}
	w.log.Trace("processEraPoolUpdate", "source", m.Source, "dest", m.Destination,
		"era", flow.Era, "shotId", flow.ShotId.Hex(), "symbol", flow.Symbol)

	snap := flow.Snap

	//check bond/unbond is needed
	//bond report if no need
	bondCmpUnbondResult := snap.Bond.Int.Cmp(snap.Unbond.Int)
	if bondCmpUnbondResult == 0 {
		w.log.Info("EvtEraPoolUpdated bond equal to unbond, no need to bond/unbond")
		callHash := utils.BlakeTwo256([]byte{})
		mFlow.OpaqueCalls = []*submodel.MultiOpaqueCall{
			{CallHash: hexutil.Encode(callHash[:])}}
		return w.informChain(m.Destination, m.Source, mFlow)
	}

	//get poolClient of this pool address
	poolAddrBase58Str := base58.Encode(snap.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrBase58Str)
	if err != nil {
		w.log.Error("EraPoolUpdated pool failed",
			"pool address", poolAddrBase58Str,
			"err", err)
		return false
	}

	//check exist and create
	stakeAccountSeed, stakeAccountPubkey := GetStakeAccountPubkey(poolClient.StakeBaseAccount.PublicKey, snap.Era)
	multisigTxtype := MultisigTxStakeType
	if bondCmpUnbondResult < 0 {
		multisigTxtype = MultisigTxUnStakeType
	}
	multisigTxAccountSeed, multisigTxAccountPubkey := GetMultisigTxAccountPubkey(poolClient.MultisigTxBaseAccount.PublicKey, multisigTxtype, snap.Era)
	rpcClient := poolClient.GetRpcClient()
	miniMumBalanceForStake, err := rpcClient.GetMinimumBalanceForRentExemption(context.Background(), 200)
	if err != nil {
		w.log.Error("processEraPoolUpdatedEvt GetMinimumBalanceForRentExemption failed",
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

	_, err = rpcClient.GetStakeAccountInfo(context.Background(), stakeAccountPubkey.ToBase58(), GetAccountInfoConfigDefault)
	if err != nil {
		if err == solClient.ErrAccountNotFound {
			//send from  relayers
			//create new stake acount of this era
			res, err := rpcClient.GetRecentBlockhash(context.Background())
			if err != nil {
				w.log.Error("processEraPoolUpdatedEvt GetRecentBlockhash failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}
			rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
				Instructions: []solTypes.Instruction{
					sysprog.CreateAccountWithSeed(
						poolClient.FeeAccount.PublicKey,
						stakeAccountPubkey,
						poolClient.StakeBaseAccount.PublicKey,
						solCommon.StakeProgramID,
						stakeAccountSeed,
						miniMumBalanceForStake,
						200,
					),
					stakeprog.Initialize(
						stakeAccountPubkey,
						stakeprog.Authorized{
							Staker:     poolClient.MultisignerPubkey,
							Withdrawer: poolClient.MultisignerPubkey,
						},
						stakeprog.Lockup{},
					),
				},
				Signers:         []solTypes.Account{poolClient.FeeAccount, poolClient.StakeBaseAccount},
				FeePayer:        poolClient.FeeAccount.PublicKey,
				RecentBlockHash: res.Blockhash,
			})
			if err != nil {
				w.log.Error("processEraPoolUpdatedEvt CreateRawTransaction failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}
			txHash, err := rpcClient.SendRawTransaction(context.Background(), rawTx)
			if err != nil {
				w.log.Error("processEraPoolUpdatedEvt SendRawTransaction failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}
			w.log.Info("create stake account",
				"tx hash", txHash,
				"stake account", stakeAccountPubkey.ToBase58())
		} else {
			w.log.Error("GetStakeAccountInfo err",
				"pool  address", poolAddrBase58Str,
				"stake address", stakeAccountPubkey.ToBase58(),
				"err", err)
			return false
		}
	}

	var transferInstruction solTypes.Instruction
	var stakeInstruction solTypes.Instruction
	var splitInstruction solTypes.Instruction
	var deactiveInstruction solTypes.Instruction

	validatorPubkey := solCommon.PublicKeyFromString("5MMCR4NbTZqjthjLGywmeT66iwE9J9f7kjtxzJjwfUx2")
	multisigTxAccountInfo, err := rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58(), GetAccountInfoConfigDefault)
	if err != nil {
		if err == solClient.ErrAccountNotFound {
			res, err := rpcClient.GetRecentBlockhash(context.Background())
			if err != nil {
				w.log.Error("processEraPoolUpdatedEvt GetRecentBlockhash failed",
					"pool address", poolAddrBase58Str,
					"err", err)
				return false
			}

			if bondCmpUnbondResult > 0 { //stake

				//send from o relayers
				//create transaction account of this era
				//todo get validator from chain

				val := new(big.Int).Sub(snap.Bond.Int, snap.Unbond.Int)
				transferInstruction = sysprog.Transfer(poolClient.MultisignerPubkey, stakeAccountPubkey, val.Uint64())
				stakeInstruction = stakeprog.DelegateStake(stakeAccountPubkey, poolClient.MultisignerPubkey, validatorPubkey)

				rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
					Instructions: []solTypes.Instruction{
						sysprog.CreateAccountWithSeed(
							poolClient.FeeAccount.PublicKey,
							multisigTxAccountPubkey,
							poolClient.MultisigTxBaseAccount.PublicKey,
							solCommon.MultisigProgramID,
							multisigTxAccountSeed,
							miniMumBalanceForTx,
							1000,
						),
						multisigprog.CreateTransaction(
							[]solCommon.PublicKey{solCommon.SystemProgramID, solCommon.StakeProgramID},
							[][]solTypes.AccountMeta{transferInstruction.Accounts, stakeInstruction.Accounts},
							[][]byte{transferInstruction.Data, stakeInstruction.Data},
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
			} else { //unstake
				val := new(big.Int).Sub(snap.Unbond.Int, snap.Bond.Int)
				splitInstruction = stakeprog.Split(poolClient.StakeBaseAccount.PublicKey,
					poolClient.MultisignerPubkey, stakeAccountPubkey, val.Uint64())

				deactiveInstruction = stakeprog.Deactivate(stakeAccountPubkey, poolClient.MultisignerPubkey)

				rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
					Instructions: []solTypes.Instruction{
						sysprog.CreateAccountWithSeed(
							poolClient.FeeAccount.PublicKey,
							multisigTxAccountPubkey,
							poolClient.MultisigTxBaseAccount.PublicKey,
							solCommon.MultisigProgramID,
							multisigTxAccountSeed,
							miniMumBalanceForTx,
							1000,
						),
						multisigprog.CreateTransaction(
							[]solCommon.PublicKey{solCommon.StakeProgramID, solCommon.StakeProgramID},
							[][]solTypes.AccountMeta{splitInstruction.Accounts, deactiveInstruction.Accounts},
							[][]byte{splitInstruction.Data, deactiveInstruction.Data},
							poolClient.MultisigInfoPubkey,
							multisigTxAccountPubkey,
							poolClient.FeeAccount.PublicKey,
						),
					},
					Signers:         []solTypes.Account{poolClient.FeeAccount},
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
			}

		} else {
			w.log.Error("GetMultisigTxAccountInfo err",
				"pool  address", poolAddrBase58Str,
				"multisig tx account address", multisigTxAccountPubkey.ToBase58(),
				"err", err)
			return false
		}
	}

	if multisigTxAccountInfo.DidExecute == 1 {
		callHash := utils.BlakeTwo256([]byte{})
		mFlow.OpaqueCalls = []*submodel.MultiOpaqueCall{
			{CallHash: hexutil.Encode(callHash[:])}}
		return w.informChain(m.Destination, m.Source, mFlow)
	}

	var remainingAccounts []solTypes.AccountMeta
	if bondCmpUnbondResult > 0 {
		remainingAccounts = multisigprog.GetRemainAccounts([]solTypes.Instruction{transferInstruction, stakeInstruction})
	} else {
		remainingAccounts = multisigprog.GetRemainAccounts([]solTypes.Instruction{splitInstruction, deactiveInstruction})
	}

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

	callHash := utils.BlakeTwo256([]byte{})
	mFlow.OpaqueCalls = []*submodel.MultiOpaqueCall{
		{CallHash: hexutil.Encode(callHash[:])}}
	return w.informChain(m.Destination, m.Source, mFlow)
}
