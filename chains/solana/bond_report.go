package solana

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/mr-tron/base58"
	substrateTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
)

// 1 get derived accounts used to stake
// 2 merge and withdraw
// 3 active report total stake value
func (w *writer) processBondReportEvent(m *core.Message) bool {
	flow, ok := m.Content.(*submodel.BondReportedFlow)
	if !ok {
		w.printContentError(m, errors.New("msg cast to BondReportFlow not ok"))
		return false
	}
	poolAddrBase58Str := base58.Encode(flow.Snap.Pool)
	poolClient, err := w.conn.GetPoolClient(poolAddrBase58Str)
	if err != nil {
		w.log.Error("processBondReportEvent GetPoolClient failed",
			"pool  address", poolAddrBase58Str,
			"error", err)
		return false
	}

	currentEra := flow.Snap.Era
	rpcClient := poolClient.GetRpcClient()

	ok = w.MergeAndWithdraw(poolClient, poolAddrBase58Str, flow.Snap.Era, flow.ShotId, flow.Snap.Pool)
	if !ok {
		return false
	}

	activeTotal := int64(0)
	// must deal every stakeBaseAccounts
	for _, useStakeBaseAccountPubKey := range poolClient.StakeBaseAccountPubkeys {

		//get base account
		accountInfo, err := rpcClient.GetStakeAccountInfo(context.Background(),
			useStakeBaseAccountPubKey.ToBase58())

		if err != nil {
			w.log.Error("processBondReportEvent GetStakeAccountInfo failed",
				"pool  address", poolAddrBase58Str,
				"stake base account", useStakeBaseAccountPubKey.ToBase58(),
				"error", err)
			return false
		}
		w.log.Info("active total will add base", "account", fmt.Sprintf(" %#v", accountInfo))
		activeTotal += accountInfo.StakeAccount.Info.Stake.Delegation.Stake
		//get derived account
		for i := uint32(0); i < uint32(backCheckLen); i++ {
			//this use current era not snap era
			stakeAccountPubkey, _ := GetStakeAccountPubkey(useStakeBaseAccountPubKey, currentEra-i)
			accountInfo, err := rpcClient.GetStakeAccountInfo(context.Background(), stakeAccountPubkey.ToBase58())
			if err != nil {
				if err == solClient.ErrAccountNotFound {
					w.log.Info("fetchSubStakeAccount not found", "account", stakeAccountPubkey.ToBase58(), "currentEra", currentEra, "index", i)
					continue
				} else {
					w.log.Error("processBondReportEvent GetStakeAccountInfo failed",
						"pool  address", poolAddrBase58Str,
						"stake account", stakeAccountPubkey.ToBase58(),
						"error", err)
					return false
				}
			}
			w.log.Info("fetchSubStakeAccount found", "account", stakeAccountPubkey.ToBase58(), "currentEra", currentEra, "index", i)

			//filter account used to stake
			if accountInfo.StakeAccount.IsStakeAndNoDeactive() {
				w.log.Info("active total will add", "account", fmt.Sprintf(" %#v", accountInfo))
				activeTotal += accountInfo.StakeAccount.Info.Stake.Delegation.Stake
			}
		}

	}

	flow.Snap.Active = substrateTypes.NewU128(*big.NewInt(activeTotal))

	w.log.Info("active report", "pool address", poolAddrBase58Str,
		"era", flow.Snap.Era, "active", activeTotal, "symbol", flow.Symbol)
	return w.activeReport(flow.Symbol, core.RFIS, flow)
}
