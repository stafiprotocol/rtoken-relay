// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
)

type listener struct {
	name          string
	symbol        core.RSymbol
	cares         []core.RSymbol
	startBlock    uint64
	blockstore    blockstore.Blockstorer
	conn          *Connection
	subscriptions map[eventName]eventHandler // Handlers for specific events
	router        chains.Router
	log           core.Logger
	stop          <-chan int
	sysErr        chan<- error
	lastEraBlock  uint64
}

// Frequency of polling for a new block
var (
	OneBlockTime              = 5 * time.Second
	BlockRetryInterval        = time.Second * 6
	BlockRetryLimit           = 100
	BlockIntervalToProcessEra = uint64(10)
	EventRetryLimit           = 6000
	EventRetryInterval        = 100 * time.Millisecond
)

func NewListener(name string, symbol core.RSymbol, opts map[string]interface{}, startBlock uint64, bs blockstore.Blockstorer, conn *Connection, log core.Logger, stop <-chan int, sysErr chan<- error) *listener {
	cares := make([]core.RSymbol, 0)
	optCares := opts["cares"]
	log.Info("NewListener", "optCares", optCares)
	if optCares != nil {
		if tmpCares, ok := optCares.([]interface{}); ok {
			for _, tc := range tmpCares {
				care, ok := tc.(string)
				if !ok {
					panic("care not string")
				}
				cares = append(cares, core.RSymbol(care))
			}
		} else {
			panic("opt cares not string array")
		}
	}
	return &listener{
		name:          name,
		symbol:        symbol,
		cares:         cares,
		startBlock:    startBlock,
		blockstore:    bs,
		conn:          conn,
		subscriptions: make(map[eventName]eventHandler),
		log:           log,
		stop:          stop,
		sysErr:        sysErr,
		lastEraBlock:  0,
	}
}

func (l *listener) setRouter(r chains.Router) {
	l.router = r
}

// Start creates the initial subscription for all events
func (l *listener) start() error {
	latestBlk, err := l.conn.LatestBlockNumber()
	if err != nil {
		return err
	}

	if latestBlk < l.startBlock {
		return fmt.Errorf("starting block (%d) is greater than latest known block (%d)", l.startBlock, latestBlk)
	}

	if l.symbol == core.RFIS {
		for _, sub := range StafiSubscriptions {
			err := l.registerEventHandler(sub.name, sub.handler)
			if err != nil {
				return err
			}
		}
	} else {
		for _, sub := range OtherSubscriptions {
			err := l.registerEventHandler(sub.name, sub.handler)
			if err != nil {
				return err
			}
		}
	}

	go func() {
		err = l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

// registerEventHandler enables a handler for a given event. This cannot be used after Start is called.
func (l *listener) registerEventHandler(name eventName, handler eventHandler) error {
	if l.subscriptions[name] != nil {
		return fmt.Errorf("event %s already registered", name)
	}
	l.subscriptions[name] = handler
	return nil
}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before returning with an error.
func (l *listener) pollBlocks() error {
	var currentBlock = l.startBlock
	var retry = BlockRetryLimit
	for {
		select {
		case <-l.stop:
			return ErrorTerminated
		default:
			// No more retries, goto next block
			if retry == 0 {
				l.sysErr <- fmt.Errorf("event polling retries exceeded: %s", l.symbol)
				return nil
			}

			finalBlk, err := l.conn.FinalizedBlockNumber()
			if err != nil {
				l.log.Error("Failed to fetch latest blockNumber", "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Sleep if the block we want comes after the most recently finalized block
			if currentBlock+l.blockDelay() > finalBlk {
				if currentBlock%100 == 0 {
					l.log.Trace("Block not yet finalized", "target", currentBlock, "finalBlk", finalBlk)
				}
				time.Sleep(OneBlockTime)
				continue
			}

			if l.symbol != core.RFIS && l.lastEraBlock+BlockIntervalToProcessEra < currentBlock {
				l.lastEraBlock = currentBlock
				err = l.processEra()
				if err != nil {
					l.log.Error("Failed to processEra", "err", err)
					retry--
					continue
				}
			}

			err = l.processEvents(currentBlock)
			if err != nil {
				l.log.Error("Failed to process events in block", "block", currentBlock, "err", err)
				retry--
				continue
			}

			// Write to blockstore
			err = l.blockstore.StoreBlock(big.NewInt(0).SetUint64(currentBlock))
			if err != nil {
				l.log.Error("Failed to write to blockstore", "err", err)
			}
			currentBlock++
			retry = BlockRetryLimit
		}
	}
}

func (l *listener) processEra() error {
	era, err := l.conn.CurrentEra()
	if err != nil {
		return err
	}

	msg := &core.Message{Destination: core.RFIS, Reason: core.NewEra, Content: era}
	return l.submitMessage(msg, nil)
}

// processEvents fetches a block and parses out the events, calling Listener.handleEvents()
func (l *listener) processEvents(blockNum uint64) error {
	if blockNum%100 == 0 {
		l.log.Debug("processEvents", "blockNum", blockNum)
	}
	evts, err := l.conn.GetEvents(blockNum)
	if err != nil {
		for i := 0; i < EventRetryLimit; i++ {
			time.Sleep(EventRetryInterval)
			evts, err = l.conn.GetEvents(blockNum)
			if err == nil {
				break
			}
		}
		if err != nil {
			return err
		}
	}

	for _, evt := range evts {
		switch l.symbol {
		case core.RFIS:
			switch {
			case evt.ModuleId == config.RTokenSeriesModuleId && evt.EventId == config.LiquidityBondEventId:
				l.log.Trace("Handling LiquidityBondEvent", "block", blockNum)
				flow, err := l.processLiquidityBondEvent(evt)
				if err != nil {
					if err.Error() == ErrNotCared.Error() {
						continue
					}
					return err
				}
				if flow.State == submodel.Success {
					l.log.Info("bond flow already succeeded, will ignore", "symbol", flow.Symbol, "bondId", flow.BondId.Hex())
					continue
				}
				if l.cared(flow.Record.Symbol) && l.subscriptions[LiquidityBond] != nil {
					err := l.submitMessage(l.subscriptions[LiquidityBond](flow))
					if err != nil {
						return err
					}

					// here we wait until bondstate change to another
					// so we can continuely process this event when restart
					if flow.Symbol == core.RMATIC || flow.Symbol == core.RBNB || flow.Symbol == core.RSOL {
						symBz, err := types.EncodeToBytes(flow.Symbol)
						if err != nil {
							return err
						}
						bsk := submodel.BondStateKey{BlockHash: flow.Record.Blockhash, TxHash: flow.Record.Txhash}
						bskBz, err := types.EncodeToBytes(bsk)
						if err != nil {
							return err
						}

						for {
							var bs submodel.BondState
							exist, err := l.conn.QueryStorage(config.RTokenSeriesModuleId, config.StorageBondStates, symBz, bskBz, &bs)
							if err != nil || !exist {
								info := fmt.Sprintf("failed to get liquidity bondstate, symbol: %s, bondId: %s, BlockHash: %s, TxHash: %s",
									flow.Symbol, flow.BondId.Hex(), hexutil.Encode(flow.Record.Blockhash), hexutil.Encode(flow.Record.Txhash))

								l.log.Warn("get liquidity bondstate faile, will wait", "err", err, "liquidityBondInfo", info)
								time.Sleep(BlockRetryInterval)
								continue
							}

							if bs == submodel.Dealing || bs == submodel.Default {
								l.log.Debug("liquidity bondstate not change will wait", "bondState", bs)
								time.Sleep(BlockRetryInterval)
								continue
							}
							break
						}

					}
				}
			case evt.ModuleId == config.RTokenLedgerModuleId && evt.EventId == config.EraPoolUpdatedEventId:
				l.log.Trace("Handling EraPoolUpdatedEvent", "block", blockNum)
				flow, err := l.processEraPoolUpdatedEvt(evt)
				if err != nil {
					if err.Error() == ErrorBondStateNotEraUpdated.Error() ||
						err.Error() == ErrorEventEraIsOld.Error() ||
						err.Error() == ErrNotCared.Error() {
						continue
					}
					return err
				}

				if l.cared(flow.Symbol) && l.subscriptions[EraPoolUpdated] != nil {
					err = l.submitMessage(l.subscriptions[EraPoolUpdated](flow))
					if err != nil {
						return err
					}

					if flow.Symbol == core.RMATIC || flow.Symbol == core.RBNB || flow.Symbol == core.RSOL {
						eraPoolUpdatedFlow, ok := flow.EventData.(*submodel.EraPoolUpdatedFlow)
						if !ok {
							return fmt.Errorf("cast err")
						}
						// here we wait until snapshot's bondstate change to another
						// so we can continuely process this event when restart
						for {
							snap, err := l.snapshot(eraPoolUpdatedFlow.Symbol, eraPoolUpdatedFlow.ShotId)
							if err != nil {
								l.log.Warn("l.snapshot", "err", err)
								time.Sleep(BlockRetryInterval)
								continue
							}
							if snap.BondState == submodel.EraUpdated {
								l.log.Debug("bondstate not change will wait", "bondState", snap.BondState)
								time.Sleep(BlockRetryInterval)
								continue
							}
							break
						}

					}

				}
			case evt.ModuleId == config.RTokenLedgerModuleId && evt.EventId == config.BondReportedEventId:
				l.log.Trace("Handling BondReportedEvent", "block", blockNum)
				flow, err := l.processBondReportedEvt(evt)
				if err != nil {
					if err.Error() == ErrorBondStateNotBondReported.Error() ||
						err.Error() == ErrorEventEraIsOld.Error() ||
						err.Error() == ErrNotCared.Error() {
						continue
					}
					return err
				}
				if l.cared(flow.Snap.Symbol) && l.subscriptions[BondReported] != nil {
					err = l.submitMessage(l.subscriptions[BondReported](flow))
					if err != nil {
						return err
					}

					if flow.Symbol == core.RMATIC || flow.Symbol == core.RBNB || flow.Symbol == core.RSOL {
						// here we wait until snapshot's bondstate change to another
						// so we can continuely process this event when restart
						for {
							snap, err := l.snapshot(flow.Symbol, flow.ShotId)
							if err != nil {
								l.log.Warn("l.snapshot", "err", err)
								time.Sleep(BlockRetryInterval)
								continue
							}
							if snap.BondState == submodel.BondReported {
								l.log.Debug("bondstate not change will wait", "bondState", snap.BondState)
								time.Sleep(BlockRetryInterval)
								continue
							}
							break
						}

					}
				}
			case evt.ModuleId == config.RTokenLedgerModuleId && evt.EventId == config.ActiveReportedEventId:
				l.log.Trace("Handling ActiveReportEvent", "block", blockNum)
				flow, err := l.processActiveReportedEvt(evt)
				if err != nil {
					if err.Error() == ErrorBondStateNotActiveReported.Error() ||
						err.Error() == ErrorEventEraIsOld.Error() ||
						err.Error() == ErrNotCared.Error() {
						continue
					}
					return err
				}
				if l.cared(flow.Symbol) && l.subscriptions[ActiveReported] != nil {
					err = l.submitMessage(l.subscriptions[ActiveReported](flow))
					if err != nil {
						return err
					}

					if flow.Symbol == core.RMATIC || flow.Symbol == core.RSOL {
						activeReportedFlow, ok := flow.EventData.(*submodel.ActiveReportedFlow)
						if !ok {
							return fmt.Errorf("cast err")
						}
						// here we wait until snapshot's bondstate change to another
						// so we can continuely process this event when restart
						for {
							snap, err := l.snapshot(activeReportedFlow.Symbol, activeReportedFlow.ShotId)
							if err != nil {
								l.log.Warn("l.snapshot", "err", err)
								time.Sleep(BlockRetryInterval)
								continue
							}
							if snap.BondState == submodel.ActiveReported {
								l.log.Debug("bondstate not change will wait", "bondState", snap.BondState)
								time.Sleep(BlockRetryInterval)
								continue
							}
							break
						}

					}
					if flow.Symbol == core.RBNB {
						withdrawReportedFlow, ok := flow.EventData.(*submodel.WithdrawReportedFlow)
						if !ok {
							return fmt.Errorf("cast err")
						}
						// here we wait until snapshot's bondstate change to another
						// so we can continuely process this event when restart
						for {
							snap, err := l.snapshot(withdrawReportedFlow.Symbol, withdrawReportedFlow.ShotId)
							if err != nil {
								l.log.Warn("l.snapshot", "err", err)
								time.Sleep(BlockRetryInterval)
								continue
							}
							if snap.BondState == submodel.ActiveReported {
								l.log.Debug("bondstate not change will wait", "bondState", snap.BondState)
								time.Sleep(BlockRetryInterval)
								continue
							}
							break
						}

					}
				}
			case evt.ModuleId == config.RTokenLedgerModuleId && evt.EventId == config.WithdrawReportedEventId:
				l.log.Trace("Handling WithdrawReportedEvent", "block", blockNum)
				flow, err := l.processWithdrawReportedEvt(evt)
				if err != nil {
					if err.Error() == ErrorBondStateNotWithdrawReported.Error() ||
						err.Error() == ErrorEventEraIsOld.Error() ||
						err.Error() == ErrNotCared.Error() {
						continue
					}
					return err
				}
				if l.cared(flow.Symbol) && l.subscriptions[WithdrawReported] != nil {
					err = l.submitMessage(l.subscriptions[WithdrawReported](flow))
					if err != nil {
						return err
					}

					if flow.Symbol == core.RMATIC || flow.Symbol == core.RBNB || flow.Symbol == core.RSOL {
						withdrawReportedFlow, ok := flow.EventData.(*submodel.WithdrawReportedFlow)
						if !ok {
							return fmt.Errorf("cast err")
						}
						// here we wait until snapshot's bondstate change to another
						// so we can continuely process this event when restart
						for {
							snap, err := l.snapshot(withdrawReportedFlow.Symbol, withdrawReportedFlow.ShotId)
							if err != nil {
								l.log.Warn("l.snapshot", "err", err)
								time.Sleep(BlockRetryInterval)
								continue
							}
							if snap.BondState == submodel.WithdrawReported {
								l.log.Debug("bondstate not change will wait", "bondState", snap.BondState)
								time.Sleep(BlockRetryInterval)
								continue
							}
							break
						}

					}
				}
			case evt.ModuleId == config.RTokenLedgerModuleId && evt.EventId == config.TransferReportedEventId:
				l.log.Trace("Handling TransferReportedEvent", "block", blockNum)
				flow, err := l.processTransferReportedEvt(evt)
				if err != nil {
					if err.Error() == ErrorBondStateNotTransferReported.Error() ||
						err.Error() == ErrNotCared.Error() {
						continue
					}
					return err
				}
				if l.cared(flow.Symbol) && l.subscriptions[TransferReported] != nil {
					err = l.submitMessage(l.subscriptions[TransferReported](flow))
					if err != nil {
						return err
					}
				}
			case evt.ModuleId == config.RTokenSeriesModuleId && evt.EventId == config.NominationUpdatedEventId:
				l.log.Trace("Handling NominationUpdatedEvent", "block", blockNum)
				flow, err := l.processNominationUpdated(evt)
				if err != nil {
					if err.Error() == ErrNotCared.Error() {
						continue
					}
					return err
				}
				if l.cared(flow.Symbol) && l.subscriptions[NominationUpdated] != nil {
					err = l.submitMessage(l.subscriptions[NominationUpdated](flow))
					if err != nil {
						return err
					}
				}
			case evt.ModuleId == config.RTokenSeriesModuleId && evt.EventId == config.ValidatorUpdatedEventId:
				l.log.Trace("Handling ValidatorUpdatedEvent", "block", blockNum)
				flow, err := l.processValidatorUpdated(evt)
				if err != nil {
					if err.Error() == ErrNotCared.Error() {
						continue
					}
					return err
				}
				if l.cared(flow.Symbol) && l.subscriptions[ValidaterUpdated] != nil {
					err = l.submitMessage(l.subscriptions[ValidaterUpdated](flow))
					if err != nil {
						return err
					}
				}
			case evt.ModuleId == config.RTokenSeriesModuleId && evt.EventId == config.SignaturesEnoughEventId:
				l.log.Trace("Handling SignaturesEnoughEventId event", "block", blockNum)
				sigs, err := l.processSignatureEnoughEvt(evt)
				if err != nil {
					if err.Error() == ErrNotCared.Error() {
						continue
					}
					return err
				}
				if l.cared(sigs.Symbol) && l.subscriptions[SignatureEnough] != nil {
					err = l.submitMessage(l.subscriptions[SignatureEnough](sigs))
					if err != nil {
						return err
					}
				}
			}
		case core.RDOT, core.RKSM:
			switch {
			case evt.ModuleId == config.MultisigModuleId && evt.EventId == config.NewMultisigEventId:
				l.log.Trace("Handling NewMultisigEvent", "block", blockNum)
				flow, err := l.processNewMultisigEvt(evt)
				if err != nil {
					if err.Error() == ErrormultiEnd.Error() {
						l.log.Info("listener received an ended NewMultisig event, ignored")
						continue
					}
					return err
				}
				if l.subscriptions[NewMultisig] != nil {
					err = l.submitMessage(l.subscriptions[NewMultisig](flow))
					if err != nil {
						return err
					}
				}
			case evt.ModuleId == config.MultisigModuleId && evt.EventId == config.MultisigExecutedEventId:
				l.log.Trace("Handling MultisigExecutedEvent", "block", blockNum)
				flow, err := l.processMultisigExecutedEvt(evt)
				if err != nil {
					return err
				}
				if l.subscriptions[MultisigExecuted] != nil {
					err = l.submitMessage(l.subscriptions[MultisigExecuted](flow))
					if err != nil {
						return err
					}
				}
			}
		default:
			l.log.Error("process event unsupport symbol", "symbol", l.symbol)
		}
	}

	return nil
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (l *listener) submitMessage(m *core.Message, err error) error {
	if err != nil {
		l.log.Error("Critical error before sending message", "err", err)
		return err
	}
	m.Source = l.symbol
	if m.Destination == "" {
		m.Destination = m.Source
	}
	return l.router.Send(m)
}

func (l *listener) blockDelay() uint64 {
	switch l.symbol {
	case core.RFIS:
		return 1
	default:
		return 0
	}
}

func (l *listener) cared(symbol core.RSymbol) bool {
	if len(l.cares) == 0 {
		return true
	}

	for _, care := range l.cares {
		if care == symbol {
			return true
		}
	}

	return false
}
