package substrate

import (
	"errors"
	"fmt"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
)

func (w *writer) rebuildEraUpdateEventUseStafiConn(shotId types.Hash, snap *submodel.PoolSnapshot) error {
	th, sub, err := w.stafiConn.thresholdAndSubAccounts(snap.Symbol, snap.Pool)
	if err != nil {
		return err
	}
	key, others := w.stafiConn.FoundFirstSubAccount(sub)
	if key == nil {
		return errors.New("EraPoolUpdated FoundFirstSubAccount have no sub key")
	}

	flow := &submodel.EraPoolUpdatedFlow{
		Symbol:    snap.Symbol,
		Era:       snap.Era,
		ShotId:    shotId,
		LastVoter: snap.LastVoter,
	}
	flow.LastVoterFlag = w.stafiConn.IsLastVoter(flow.LastVoter)
	flow.Snap = snap

	mef := &submodel.MultiEventFlow{
		EventId:     config.EraPoolUpdatedEventId,
		Symbol:      snap.Symbol,
		EventData:   flow,
		Threshold:   th,
		SubAccounts: sub,
	}

	mef.Key, mef.Others = key, others
	call, err := w.stafiConn.BondOrUnbondCall(snap)
	if err != nil {
		//no tx no need to rebuild
		if err.Error() == substrate.BondEqualToUnbondError.Error() {
			return err
		}
		w.log.Error("BondOrUnbondCall error", "err", err)
		return err
	}

	info, err := w.stafiConn.PaymentQueryInfo(call.Extrinsic)
	if err != nil {
		w.log.Error("PaymentQueryInfo error", "err", err, "callHash", call.CallHash, "Extrinsic", call.Extrinsic)
		return err
	}
	mef.PaymentInfo = info
	mef.OpaqueCalls = []*submodel.MultiOpaqueCall{call}
	callhash := call.CallHash
	mef.NewMulCallHashs = map[string]bool{callhash: true}
	mef.MulExeCallHashs = map[string]bool{callhash: true}

	w.setEvents(call.CallHash, mef)
	return nil
}

func (w *writer) getLatestBondedPoolsUseStafiConn(symbol core.RSymbol) ([]types.Bytes, error) {
	symbz, err := types.EncodeToBytes(symbol)
	if err != nil {
		w.sysErr <- err
		return nil, err
	}
	bondedPools := make([]types.Bytes, 0)
	exist, err := w.stafiConn.QueryStorage(config.RTokenLedgerModuleId, config.StorageBondedPools, symbz, nil, &bondedPools)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("bonded pools not extis: %s", symbol)
	}

	return bondedPools, nil
}
