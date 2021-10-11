package solana

import (
	"context"
	"errors"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mr-tron/base58"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
	solCommon "github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/multisigprog"
	"github.com/stafiprotocol/solana-go-sdk/stakeprog"
	"github.com/stafiprotocol/solana-go-sdk/sysprog"
	solTypes "github.com/stafiprotocol/solana-go-sdk/types"
)

//process eraPoolUpdate event
//0 check bond/unbond
//  (1)if no need just report bond to stafi
//  (2)if need go next
//1 query stake acount is created
//  (1)if no then created and go next
//  (2)if created go next
//2 query multisig tx acount is created
//  (1)if no then created and go next
//  (2)if created go next
//3 approve mutisig tx
//4 query multisig tx executed result
//  (1)if executed then report bond result to stafi
//  (2)if reach retry limit return false
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
	bondCmpUnbondResult := snap.Bond.Int.Cmp(snap.Unbond.Int)
	//stake
	switch {
	case bondCmpUnbondResult == 0: //bond report if no need
		w.log.Info("EvtEraPoolUpdated bond equal to unbond, no need to bond/unbond")
	case bondCmpUnbondResult > 0: //stake
		ok := w.processStake(snap)
		if !ok {
			w.log.Error("processStake failed")
			return false
		}
	case bondCmpUnbondResult < 0: //unstake
		ok := w.processUnStake(snap)
		if !ok {
			w.log.Error("processUnStake failed")
			return false
		}
	}

	//bond report
	callHash := utils.BlakeTwo256([]byte{})
	mFlow.OpaqueCalls = []*submodel.MultiOpaqueCall{
		{CallHash: hexutil.Encode(callHash[:])}}
	return w.informChain(m.Destination, m.Source, mFlow)
}

func (w *writer) processUnStake(snap *submodel.PoolSnapshot) bool {
	//get poolClient of this pool address
	poolAddrBase58Str := base58.Encode(snap.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrBase58Str)
	if err != nil {
		w.log.Error("EraPoolUpdated pool failed",
			"pool address", poolAddrBase58Str,
			"err", err)
		return false
	}
	rpcClient := poolClient.GetRpcClient()

	//collect stake amount
	//will check stakeAccounts derived in this era for preventing the case of interruption
	//if exist add the stake amount of derived stakeAccount to stakeBaseAccount
	stakeBaseAccountPubkeyToStakeAmount := make(map[solCommon.PublicKey]int64)
	willUseStakeBaseAccounts := make([]solCommon.PublicKey, 0)
	for _, stakeBaseAccountPubKey := range poolClient.StakeBaseAccountPubkeys {
		stakeBaseAccountInfo, err := rpcClient.GetStakeAccountInfo(context.Background(), stakeBaseAccountPubKey.ToBase58())
		if err != nil {
			w.log.Error("processUnStake GetStakeAccountInfo failed",
				"stake account", stakeBaseAccountPubKey.ToBase58(),
				"err", err)
			return false
		}
		stakeAmount := stakeBaseAccountInfo.StakeAccount.Info.Stake.Delegation.Stake

		//add stake amount of  derived stake account if exist
		derivedStakeAccountPubkey, _ := GetStakeAccountPubkey(stakeBaseAccountPubKey, snap.Era)
		derivedStakeAccountInfo, err := rpcClient.GetStakeAccountInfo(context.Background(), derivedStakeAccountPubkey.ToBase58())
		if err != nil {
			if err != solClient.ErrAccountNotFound {
				w.log.Error("processUnStake GetStakeAccountInfo failed",
					"stake account", stakeBaseAccountPubKey.ToBase58(),
					"err", err)
				return false
			}
		} else {
			stakeAmount += derivedStakeAccountInfo.StakeAccount.Info.Stake.Delegation.Stake
		}

		stakeBaseAccountPubkeyToStakeAmount[stakeBaseAccountPubKey] = stakeAmount
		willUseStakeBaseAccounts = append(willUseStakeBaseAccounts, stakeBaseAccountPubKey)
		time.Sleep(500 * time.Millisecond)
	}

	//sort by amount
	//sort stakeAccount pubkey by delegate amount
	sort.Slice(willUseStakeBaseAccounts, func(i int, j int) bool {
		return stakeBaseAccountPubkeyToStakeAmount[willUseStakeBaseAccounts[i]] >
			stakeBaseAccountPubkeyToStakeAmount[willUseStakeBaseAccounts[j]]
	})
	//choose validators to be undelegated
	choosedStakeBaseAccount := make([]solCommon.PublicKey, 0)
	choosedAmount := make(map[solCommon.PublicKey]int64)
	willUseTotalVal := new(big.Int).Sub(snap.Unbond.Int, snap.Bond.Int).Int64()
	selectedAmount := int64(0)
	for _, baseAccount := range willUseStakeBaseAccounts {
		nowValMaxUnDeleAmount := stakeBaseAccountPubkeyToStakeAmount[baseAccount] - int64(initStakeAmount)
		if selectedAmount+nowValMaxUnDeleAmount >= willUseTotalVal {
			willUseChoosedAmount := willUseTotalVal - selectedAmount

			choosedStakeBaseAccount = append(choosedStakeBaseAccount, baseAccount)
			choosedAmount[baseAccount] = willUseChoosedAmount
			selectedAmount = selectedAmount + willUseChoosedAmount
			break
		}

		choosedStakeBaseAccount = append(choosedStakeBaseAccount, baseAccount)
		choosedAmount[baseAccount] = nowValMaxUnDeleAmount
		selectedAmount = selectedAmount + nowValMaxUnDeleAmount
	}

	// must deal every selected takeBaseAccounts
	for stakeBaseAccountIndex, useStakeBaseAccountPubKey := range choosedStakeBaseAccount {
		w.log.Info("processEraPoolUpdatedEvt is dealing stakeBaseAccounts", "index", stakeBaseAccountIndex,
			"stakeBaseAccount", useStakeBaseAccountPubKey.ToBase58())
		//check exist and create
		stakeAccountPubkey, stakeAccountSeed := GetStakeAccountPubkey(useStakeBaseAccountPubKey, snap.Era)
		multisigTxtype := MultisigTxUnStakeType
		multisigTxAccountPubkey, multisigTxAccountSeed := GetMultisigTxAccountPubkey(
			poolClient.MultisigTxBaseAccountPubkey,
			poolClient.MultisigProgramId,
			multisigTxtype,
			snap.Era,
			stakeBaseAccountIndex)

		miniMumBalanceForStake, err := rpcClient.GetMinimumBalanceForRentExemption(context.Background(),
			solClient.StakeAccountInfoLengthDefault)
		if err != nil {
			w.log.Error("processEraPoolUpdatedEvt GetMinimumBalanceForRentExemption failed",
				"pool address", poolAddrBase58Str,
				"err", err)
			return false
		}
		miniMumBalanceForStake += initStakeAmount // stake account balance must > minRent

		if poolClient.HasBaseAccountAuth {
			useStakeBaseAccount, exist := poolClient.StakeBasePubkeyToAccounts[useStakeBaseAccountPubKey]
			if !exist {
				w.log.Error("stakeBaseAccount privkey not exist", "stakeBaseAccount", useStakeBaseAccountPubKey.ToBase58())
				return false
			}
			_, err = rpcClient.GetStakeAccountInfo(context.Background(), stakeAccountPubkey.ToBase58())
			if err != nil && err == solClient.ErrAccountNotFound {
				//send from  relayers no need multisig
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
							useStakeBaseAccountPubKey,
							solCommon.StakeProgramID,
							stakeAccountSeed,
							miniMumBalanceForStake,
							solClient.StakeAccountInfoLengthDefault,
						),
					},
					Signers:         []solTypes.Account{poolClient.FeeAccount, useStakeBaseAccount},
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
				w.log.Info("processEraPoolUpdatedEvt create stake account",
					"tx hash", txHash,
					"stake account", stakeAccountPubkey.ToBase58())

			}

			if err != nil && err != solClient.ErrAccountNotFound {
				w.log.Error("processEraPoolUpdatedEvt GetStakeAccountInfo err",
					"pool  address", poolAddrBase58Str,
					"stake address", stakeAccountPubkey.ToBase58(),
					"err", err)
				return false
			}
		}

		//check stakeaccount is created
		create := w.waitingForStakeAccountCreate(rpcClient, poolAddrBase58Str, stakeAccountPubkey.ToBase58(), "processEraPoolUpdatedEvt")
		if !create {
			return false
		}
		w.log.Info("processEraPoolUpdatedEvt stakeAccount has create", "stakeAccount", stakeAccountPubkey.ToBase58())

		var splitInstruction solTypes.Instruction
		var deactiveInstruction solTypes.Instruction

		var remainingAccounts []solTypes.AccountMeta
		var programsIds []solCommon.PublicKey
		var accountMetas [][]solTypes.AccountMeta
		var datas [][]byte

		//unstake val
		val := choosedAmount[useStakeBaseAccountPubKey]
		splitInstruction = stakeprog.Split(useStakeBaseAccountPubKey,
			poolClient.MultisignerPubkey, stakeAccountPubkey, uint64(val))

		deactiveInstruction = stakeprog.Deactivate(stakeAccountPubkey, poolClient.MultisignerPubkey)
		remainingAccounts = multisigprog.GetRemainAccounts([]solTypes.Instruction{splitInstruction, deactiveInstruction})
		programsIds = []solCommon.PublicKey{splitInstruction.ProgramID, deactiveInstruction.ProgramID}
		accountMetas = [][]solTypes.AccountMeta{splitInstruction.Accounts, deactiveInstruction.Accounts}
		datas = [][]byte{splitInstruction.Data, deactiveInstruction.Data}

		if poolClient.HasBaseAccountAuth {
			if poolClient.MultisigTxBaseAccount == nil {
				w.log.Error("MultisigTxBaseAccount privkey not exist", "MultisigTxBaseAccount", poolClient.MultisigTxBaseAccountPubkey)
				return false
			}

			_, err = rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
			if err != nil && err == solClient.ErrAccountNotFound {
				sendOk := w.createMultisigTxAccount(rpcClient, poolClient, poolAddrBase58Str, programsIds, accountMetas, datas,
					multisigTxAccountPubkey, multisigTxAccountSeed, "processEraPoolUpdatedEvt")
				if !sendOk {
					return false
				}
			}

			if err != nil && err != solClient.ErrAccountNotFound {
				w.log.Error("GetMultisigTxAccountInfo err",
					"pool  address", poolAddrBase58Str,
					"multisig tx account address", multisigTxAccountPubkey.ToBase58(),
					"err", err)
				return false
			}
		}

		//check multisig tx account is created
		create = w.waitingForMultisigTxCreate(rpcClient, poolAddrBase58Str, multisigTxAccountPubkey.ToBase58(), "processEraPoolUpdatedEvt")
		if !create {
			return false
		}
		w.log.Info("processEraPoolUpdatedEvt multisigTxAccount has create", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())

		valid := w.CheckMultisigTx(rpcClient, multisigTxAccountPubkey, programsIds, accountMetas, datas)
		if !valid {
			w.log.Info("processEraPoolUpdatedEvt CheckMultisigTx failed", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())
			return false
		}
		//if has exe just bond report
		isExe := w.IsMultisigTxExe(rpcClient, multisigTxAccountPubkey)
		if isExe {
			w.log.Info("processEraPoolUpdatedEvt multisigTxAccount has execute", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())
			continue
		}

		//approve tx
		send := w.approveMultisigTx(rpcClient, poolClient, poolAddrBase58Str, multisigTxAccountPubkey, remainingAccounts, "processEraPoolUpdatedEvt")
		if !send {
			return false
		}

		//check multisig exe result
		exe := w.waitingForMultisigTxExe(rpcClient, poolAddrBase58Str, multisigTxAccountPubkey.ToBase58(), "processEraPoolUpdatedEvt")
		if !exe {
			return false
		}
		w.log.Info("processEraPoolUpdatedEvt multisigTxAccount has execute", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())

		//check splitAccount
		stakeAccountValid := w.CheckStakeAccount(rpcClient, stakeAccountPubkey.ToBase58(), poolClient.MultisignerPubkey.ToBase58())
		if !stakeAccountValid {
			return false
		}
	}
	return true
}

func (w *writer) processStake(snap *submodel.PoolSnapshot) bool {
	//get poolClient of this pool address
	poolAddrBase58Str := base58.Encode(snap.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrBase58Str)
	if err != nil {
		w.log.Error("EraPoolUpdated pool failed",
			"pool address", poolAddrBase58Str,
			"err", err)
		return false
	}
	stakeBaseAccountLen := len(poolClient.StakeBaseAccountPubkeys)

	// must deal every stakeBaseAccounts
	for stakeBaseAccountIndex, useStakeBaseAccountPubKey := range poolClient.StakeBaseAccountPubkeys {
		w.log.Info("processEraPoolUpdatedEvt is dealing stakeBaseAccounts", "index", stakeBaseAccountIndex,
			"stakeBaseAccount", useStakeBaseAccountPubKey.ToBase58())
		//check exist and create
		stakeAccountPubkey, stakeAccountSeed := GetStakeAccountPubkey(useStakeBaseAccountPubKey, snap.Era)
		multisigTxtype := MultisigTxStakeType

		multisigTxAccountPubkey, multisigTxAccountSeed := GetMultisigTxAccountPubkey(
			poolClient.MultisigTxBaseAccountPubkey,
			poolClient.MultisigProgramId,
			multisigTxtype,
			snap.Era,
			stakeBaseAccountIndex)

		rpcClient := poolClient.GetRpcClient()
		miniMumBalanceForStake, err := rpcClient.GetMinimumBalanceForRentExemption(context.Background(),
			solClient.StakeAccountInfoLengthDefault)
		if err != nil {
			w.log.Error("processEraPoolUpdatedEvt GetMinimumBalanceForRentExemption failed",
				"pool address", poolAddrBase58Str,
				"err", err)
			return false
		}
		miniMumBalanceForStake += initStakeAmount // stake account balance must > minRent

		if poolClient.HasBaseAccountAuth {
			useStakeBaseAccount, exist := poolClient.StakeBasePubkeyToAccounts[useStakeBaseAccountPubKey]
			if !exist {
				w.log.Error("stakeBaseAccount privkey not exist", "stakeBaseAccount", useStakeBaseAccountPubKey.ToBase58())
				return false
			}
			_, err = rpcClient.GetStakeAccountInfo(context.Background(), stakeAccountPubkey.ToBase58())
			if err != nil && err == solClient.ErrAccountNotFound {
				//send from  relayers no need multisig
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
							useStakeBaseAccountPubKey,
							solCommon.StakeProgramID,
							stakeAccountSeed,
							miniMumBalanceForStake,
							solClient.StakeAccountInfoLengthDefault,
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
					Signers:         []solTypes.Account{poolClient.FeeAccount, useStakeBaseAccount},
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
				w.log.Info("processEraPoolUpdatedEvt create stake account",
					"tx hash", txHash,
					"stake account", stakeAccountPubkey.ToBase58())

			}

			if err != nil && err != solClient.ErrAccountNotFound {
				w.log.Error("processEraPoolUpdatedEvt GetStakeAccountInfo err",
					"pool  address", poolAddrBase58Str,
					"stake address", stakeAccountPubkey.ToBase58(),
					"err", err)
				return false
			}
		}

		//check stakeaccount is created
		create := w.waitingForStakeAccountCreate(rpcClient, poolAddrBase58Str, stakeAccountPubkey.ToBase58(), "processEraPoolUpdatedEvt")
		if !create {
			return false
		}
		w.log.Info("processEraPoolUpdatedEvt stakeAccount has create", "stakeAccount", stakeAccountPubkey.ToBase58())

		stakeAccountValid := w.CheckStakeAccount(rpcClient, stakeAccountPubkey.ToBase58(), poolClient.MultisignerPubkey.ToBase58())
		if !stakeAccountValid {
			return false
		}

		var transferInstruction solTypes.Instruction
		var stakeInstruction solTypes.Instruction

		var remainingAccounts []solTypes.AccountMeta
		var programsIds []solCommon.PublicKey
		var accountMetas [][]solTypes.AccountMeta
		var datas [][]byte

		stakeBaseAccountInfo, err := rpcClient.GetStakeAccountInfo(context.Background(), useStakeBaseAccountPubKey.ToBase58())
		if err != nil {
			w.log.Error("processEraPoolUpdatedEvt GetStakeAccountInfo err",
				"pool  address", poolAddrBase58Str,
				"stake base address", useStakeBaseAccountPubKey.ToBase58(),
				"err", err)
			return false
		}
		validatorPubkey := stakeBaseAccountInfo.StakeAccount.Info.Stake.Delegation.Voter

		//stake
		totalVal := new(big.Int).Sub(snap.Bond.Int, snap.Unbond.Int)
		//stake average to stakeBaseAccounts
		val := new(big.Int).Div(totalVal, big.NewInt(int64(stakeBaseAccountLen)))

		transferInstruction = sysprog.Transfer(poolClient.MultisignerPubkey, stakeAccountPubkey, val.Uint64())
		stakeInstruction = stakeprog.DelegateStake(stakeAccountPubkey, poolClient.MultisignerPubkey, validatorPubkey)
		remainingAccounts = multisigprog.GetRemainAccounts([]solTypes.Instruction{transferInstruction, stakeInstruction})
		programsIds = []solCommon.PublicKey{transferInstruction.ProgramID, stakeInstruction.ProgramID}
		accountMetas = [][]solTypes.AccountMeta{transferInstruction.Accounts, stakeInstruction.Accounts}
		datas = [][]byte{transferInstruction.Data, stakeInstruction.Data}

		if poolClient.HasBaseAccountAuth {
			if poolClient.MultisigTxBaseAccount == nil {
				w.log.Error("MultisigTxBaseAccount privkey not exist", "MultisigTxBaseAccount", poolClient.MultisigTxBaseAccountPubkey)
				return false
			}

			_, err = rpcClient.GetMultisigTxAccountInfo(context.Background(), multisigTxAccountPubkey.ToBase58())
			if err != nil && err == solClient.ErrAccountNotFound {
				sendOk := w.createMultisigTxAccount(rpcClient, poolClient, poolAddrBase58Str, programsIds, accountMetas, datas,
					multisigTxAccountPubkey, multisigTxAccountSeed, "processEraPoolUpdatedEvt")
				if !sendOk {
					return false
				}
			}

			if err != nil && err != solClient.ErrAccountNotFound {
				w.log.Error("GetMultisigTxAccountInfo err",
					"pool  address", poolAddrBase58Str,
					"multisig tx account address", multisigTxAccountPubkey.ToBase58(),
					"err", err)
				return false
			}
		}

		//check multisig tx account is created
		create = w.waitingForMultisigTxCreate(rpcClient, poolAddrBase58Str, multisigTxAccountPubkey.ToBase58(), "processEraPoolUpdatedEvt")
		if !create {
			return false
		}
		w.log.Info("processEraPoolUpdatedEvt multisigTxAccount has create", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())

		valid := w.CheckMultisigTx(rpcClient, multisigTxAccountPubkey, programsIds, accountMetas, datas)
		if !valid {
			w.log.Info("processEraPoolUpdatedEvt CheckMultisigTx failed", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())
			return false
		}
		//if has exe just bond report
		isExe := w.IsMultisigTxExe(rpcClient, multisigTxAccountPubkey)
		if isExe {
			w.log.Info("processEraPoolUpdatedEvt multisigTxAccount has execute", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())
			continue
		}

		//approve tx
		send := w.approveMultisigTx(rpcClient, poolClient, poolAddrBase58Str, multisigTxAccountPubkey, remainingAccounts, "processEraPoolUpdatedEvt")
		if !send {
			return false
		}

		//check multisig exe result
		exe := w.waitingForMultisigTxExe(rpcClient, poolAddrBase58Str, multisigTxAccountPubkey.ToBase58(), "processEraPoolUpdatedEvt")
		if !exe {
			return false
		}
		w.log.Info("processEraPoolUpdatedEvt multisigTxAccount has execute", "multisigTxAccount", multisigTxAccountPubkey.ToBase58())

	}
	return true
}
