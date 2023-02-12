// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

const (
	bondFlowLimit = 2048
)

type writer struct {
	symbol                core.RSymbol
	conn                  *Connection
	router                chains.Router
	log                   core.Logger
	sysErr                chan<- error
	eventMtx              sync.RWMutex
	newMulTicsMtx         sync.RWMutex
	bondedPoolsMtx        sync.RWMutex
	events                map[string]*submodel.MultiEventFlow
	newMultics            map[string]*submodel.EventNewMultisig
	bondedPools           map[string]bool
	liquidityBonds        chan *core.Message
	currentChainEra       uint32
	stop                  <-chan int
	transferRecord        string
	transferRecordHistory string
}

func NewWriter(symbol core.RSymbol, opts map[string]interface{}, conn *Connection, log core.Logger, sysErr chan<- error, stop <-chan int) *writer {
	transferRecord := ""
	transferRecordHistory := ""
	if symbol != core.RFIS {
		path, ok := opts["transferRecord"].(string)
		if !ok {
			panic("no filepath to save transferRecord")
		}

		historyPath, ok := opts["transferRecordHistory"].(string)
		if !ok {
			panic("no filepath to save transferRecord history")
		}
		transferRecord = path
		transferRecordHistory = historyPath
		if _, err := os.Stat(transferRecord); os.IsNotExist(err) {
			err = utils.WriteCSV(transferRecord, [][]string{})
			if err != nil {
				panic(err)
			}
		}

		if _, err := os.Stat(transferRecordHistory); os.IsNotExist(err) {
			err = utils.WriteCSV(transferRecordHistory, [][]string{})
			if err != nil {
				panic(err)
			}
		}
	}

	return &writer{
		symbol:                symbol,
		conn:                  conn,
		log:                   log,
		sysErr:                sysErr,
		events:                make(map[string]*submodel.MultiEventFlow),
		newMultics:            make(map[string]*submodel.EventNewMultisig),
		bondedPools:           make(map[string]bool),
		liquidityBonds:        make(chan *core.Message, bondFlowLimit),
		currentChainEra:       0,
		stop:                  stop,
		transferRecord:        transferRecord,
		transferRecordHistory: transferRecordHistory,
	}
}

func (w *writer) setRouter(r chains.Router) {
	w.router = r
}

func (w *writer) ResolveMessage(m *core.Message) (processOk bool) {
	defer func() {
		if !processOk {
			w.sysErr <- fmt.Errorf("ResolveMessage failed, message: %+v", m)
		}
	}()

	switch m.Reason {
	// handle by substrate
	case core.LiquidityBondEvent:
		return w.processLiquidityBond(m)
	case core.BondedPools:
		return w.processBondedPools(m)
	case core.EraPoolUpdatedEvent:
		return w.processEraPoolUpdatedEvent(m)
	case core.BondReportedEvent:
		return w.processBondReportEvent(m)
	case core.ActiveReportedEvent:
		return w.processActiveReportedEvent(m)
	case core.WithdrawReportedEvent:
		return w.processWithdrawReportedEvent(m)
	case core.TransferReportedEvent:
		return w.processTransferReportedEvent(m)
	case core.NominationUpdatedEvent:
		return w.processNominationUpdatedEvent(m)

	// multisig tx event on substrate, handle by substarte self
	case core.NewMultisig:
		return w.processNewMultisig(m)
	case core.MultisigExecuted:
		return w.processMultisigExecuted(m)

	// handle by stafi chain
	case core.LiquidityBondResult:
		return w.processLiquidityBondResult(m)
	case core.NewEra:
		return w.processNewEra(m)
	case core.InformChain:
		return w.processInformChain(m)
	case core.ActiveReport:
		return w.processActiveReport(m)
	case core.SubmitSignature:
		return w.processSubmitSignature(m)
	case core.ExeLiquidityBondAndSwap:
		return w.processExeLiquidityBondAndSwap(m)

	// get state, handle by stafi chain
	case core.GetEraNominated:
		go w.processGetEraNominated(m)
		return true
	case core.GetBondState:
		go w.processGetBondState(m)
		return true
	case core.GetSubmitSignatures:
		go w.processGetSubmitSignatures(m)
		return true
	case core.GetPoolThreshold:
		go w.processGetPoolThreshold(m)
		return true
	default:
		w.log.Warn("message reason unsupported", "reason", m.Reason)
		return false
	}
}

func (w *writer) start() error {
	if w.symbol == core.RFIS {
		return nil
	}

	go func() {
		for {
			select {
			case <-w.stop:
				close(w.liquidityBonds)
				w.log.Info("writer stopped")
				return
			case msg := <-w.liquidityBonds:
				result := w.processLiquidityBond(msg)
				w.log.Info("retry processLiquidityBond", "result", result)
			}
		}
	}()

	return nil
}

func (w *writer) printContentError(m *core.Message) {
	w.log.Error("msg resolve failed", "source", m.Source, "dest", m.Destination, "reason", m.Reason)
}

func (w *writer) getLatestBondedPools(symbol core.RSymbol) ([]types.Bytes, error) {
	symbz, err := types.EncodeToBytes(symbol)
	if err != nil {
		return nil, err
	}
	bondedPools := make([]types.Bytes, 0)
	exist, err := w.conn.QueryStorage(config.RTokenLedgerModuleId, config.StorageBondedPools, symbz, nil, &bondedPools)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("bonded pools not extis: %s", symbol)
	}

	return bondedPools, nil
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (w *writer) submitMessage(m *core.Message) bool {
	if m.Destination == "" {
		m.Destination = core.RFIS
	}
	err := w.router.Send(m)
	if err != nil {
		w.log.Error("failed to process event", "err", err)
		return false
	}

	return true
}

func (w *writer) getEvents(key string) (*submodel.MultiEventFlow, bool) {
	w.eventMtx.RLock()
	defer w.eventMtx.RUnlock()
	value, exist := w.events[key]
	return value, exist
}

func (w *writer) setEvents(key string, value *submodel.MultiEventFlow) {
	w.eventMtx.Lock()
	defer w.eventMtx.Unlock()
	w.events[key] = value
}

func (w *writer) deleteEvents(key string) {
	w.eventMtx.Lock()
	defer w.eventMtx.Unlock()
	delete(w.events, key)
}

func (w *writer) getNewMultics(key string) (*submodel.EventNewMultisig, bool) {
	w.newMulTicsMtx.RLock()
	defer w.newMulTicsMtx.RUnlock()
	value, exist := w.newMultics[key]
	return value, exist
}

func (w *writer) setNewMultics(key string, value *submodel.EventNewMultisig) {
	w.newMulTicsMtx.Lock()
	defer w.newMulTicsMtx.Unlock()
	w.newMultics[key] = value
}

func (w *writer) deleteNewMultics(key string) {
	w.newMulTicsMtx.Lock()
	defer w.newMulTicsMtx.Unlock()
	delete(w.newMultics, key)
}

func (w *writer) getBondedPools(key string) (bool, bool) {
	w.bondedPoolsMtx.RLock()
	defer w.bondedPoolsMtx.RUnlock()
	value, exist := w.bondedPools[key]
	return value, exist
}

func (w *writer) setBondedPools(key string, value bool) {
	w.bondedPoolsMtx.Lock()
	defer w.bondedPoolsMtx.Unlock()
	w.bondedPools[key] = value
}

func (w *writer) informChain(source, dest core.RSymbol, flow *submodel.MultiEventFlow) bool {
	msg := &core.Message{Source: source, Destination: dest, Reason: core.InformChain, Content: flow}
	return w.submitMessage(msg)
}

func (w *writer) queryAndReportActive(m *core.Message) {
	flow, ok := m.Content.(*submodel.BondReportedFlow)
	if !ok {
		w.printContentError(m)
		return
	}

	ledger := new(submodel.StakingLedger)
	exist, err := w.conn.QueryStorage(config.StakingModuleId, config.StorageLedger, flow.Snap.Pool, nil, ledger)
	if err != nil {
		w.log.Error("queryAndReportActive get ledger error", "error", err, "pool", hexutil.Encode(flow.Snap.Pool))
		return
	}
	if !exist {
		w.log.Error("queryAndReportActive ledger not exist", "pool", hexutil.Encode(flow.Snap.Pool))
		return
	}

	flow.ReportActive = types.NewU128(big.Int(ledger.Active))
	w.log.Info("queryAndReportActive", "pool", hexutil.Encode(flow.Snap.Pool), "active", flow.Snap.Active)
	m.Content = flow
	m.Reason = core.ActiveReport

	w.submitMessage(m)
}
