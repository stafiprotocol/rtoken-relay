package solana

import (
	"context"
	"errors"
	"math/big"

	"github.com/mr-tron/base58"
	substrateTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	solClient "github.com/tpkeeper/solana-go-sdk/client"
)

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

	currentEra := w.conn.GetCurrentEra()
	rpcClient := poolClient.GetRpcClient()
	activeTotal := int64(0)
	//get base account
	accountInfo, err := rpcClient.GetStakeAccountInfo(context.Background(),
		poolClient.StakeBaseAccount.PublicKey.ToBase58())

	if err != nil {
		w.log.Error("processBondReportEvent GetStakeAccountInfo failed",
			"pool  address", poolAddrBase58Str,
			"stake account", poolClient.StakeBaseAccount.PublicKey.ToBase58(),
			"error", err)
		return false
	}

	activeTotal += accountInfo.StakeAccount.Info.Stake.Delegation.Stake
	//get derived account
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

		//filter account used to stake
		if accountInfo.StakeAccount.Type == 2 {
			activeTotal += accountInfo.StakeAccount.Info.Stake.Delegation.Stake
		}
	}


	flow.Snap.Active = substrateTypes.NewU128(*big.NewInt(activeTotal))

	w.log.Info("active report", "pool address", poolAddrBase58Str,
		"era", flow.Snap.Era, "active", activeTotal, "symbol", flow.Symbol)
	return w.activeReport(flow.Symbol, core.RFIS, flow)
}
