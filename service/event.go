package service

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/substrate"
)

func (l *listener) processLiquidityBondEvents(evts []*substrate.ChainEvent) error {
	for _, evt := range evts {
		err := l.processLiquidityBondEvent(evt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *listener) processLiquidityBondEvent(evt *substrate.ChainEvent) error {
	lb, err := liquidityBondEventData(evt)
	if err != nil {
		return err
	}

	bondKey := lb.bondKey()
	bk, err := types.EncodeToBytes(bondKey)
	if err != nil {
		return err
	}

	br := new(conn.BondRecord)
	exist, err := l.gsrpc.QueryStorage(config.LiquidityBondModuleId, config.StorageBondRecords, bk, nil, br)
	if err != nil {
		return err
	}

	l.log.Info("BondRecord", "bonder", hexutil.Encode(br.Bonder[:]), "symbol", br.Symbol,
		"pubkey", hexutil.Encode(br.Pubkey), "pool", hexutil.Encode(br.Pool), "blockhash", hexutil.Encode(br.Blockhash),
		"txhash", hexutil.Encode(br.Txhash), "amount", br.Amount.Int)

	if !exist {
		return fmt.Errorf("unable to get bondrecord by bondkey: %+v", lb)
	}

	if br.Bonder != lb.accountId {
		return fmt.Errorf("bonder not matched: %s, %s", hexutil.Encode(br.Bonder[:]), hexutil.Encode(lb.accountId[:]))
	}

	chain, ok := l.chains[br.Symbol]
	if !ok {
		return fmt.Errorf("no validator for symbol: %s", br.Symbol)
	}

	reason, err := chain.TransferVerify(br)
	if err != nil {
		return err
	}
	l.log.Info("TransferVerify result", "reason", reason)

	bondProp, err := l.newLiquidityBondProposal(bondKey, reason)
	if err != nil {
		return err
	}

	result := l.resolveProposal(bondProp, reason == conn.Pass)
	l.log.Info("processLiquidityBondEvent", "result", result)

	return nil
}

func (l *listener) processEraUpdatedEvts(evts []*substrate.ChainEvent) error {
	for _, evt := range evts {
		err := l.processEraUpdatedEvt(evt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *listener) processEraUpdatedEvt(evt *substrate.ChainEvent) error {
	eu, err := eraUpdatedEventData(evt)
	if err != nil {
		return err
	}

	symBz, err := types.EncodeToBytes(eu.symbol)
	if err != nil {
		return err
	}

	plcs := make([]*conn.PoolLinkChunk, 0)
	exist, err := l.gsrpc.QueryStorage(config.RTokenLedgerModuleId, config.StorageBondFaucets, symBz, nil, plcs)
	if err != nil {
		return err
	}

	if !exist {
		l.log.Warn("get era updated event, but there were no pool link chunks", "evt", eu)
	}

	chain, ok := l.chains[eu.symbol]
	if !ok {
		return fmt.Errorf("no validator for symbol: %s", eu.symbol)
	}

	err = chain.BondWork(plcs)
	if err != nil {
		return err
	}

	return nil

	//key, err := types.EncodeToBytes(eu)
	//if err != nil {
	//	return err
	//}

	//lc := new(conn.LinkChunk)
	//exist, err := l.gsrpc.QueryStorage(config.LiquidityBondModuleId, config.StorageTotalBonding, bk, nil, br)
	//if err != nil {
	//	return err
	//}
	//
	//if !exist {
	//	return fmt.Errorf("unable to get bondrecord by bondkey: %+v", lb)
	//}

}

//func (l *listener) processPoolSubAccountAddedEvents(evts []*substrate.ChainEvent) error {
//	for _, evt := range evts {
//		err := l.processPoolSubAccountAddedEvent(evt)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

//func (l *listener) processPoolSubAccountAddedEvent(evt *substrate.ChainEvent) error {
//	evtData, err := poolSubAccountAddedEventData(evt)
//	if err != nil {
//		return err
//	}
//	l.log.Trace("processPoolSubAccountAddedEvent", "evtData", evtData)
//
//	var bondFlag bool
//	pk := evtData.poolKey()
//	pkbz, err := types.EncodeToBytes(pk)
//	if err != nil {
//		return err
//	}
//
//	exist, err := l.gsrpc.QueryStorage(config.RTokenLedgerModuleId, config.StoragePoolBonded, pkbz, nil, &bondFlag)
//	if err != nil {
//		return err
//	}
//
//	if !exist {
//		return fmt.Errorf("unable to get PoolBonded by poolkey: %+v", evtData)
//	}
//
//	if bondFlag {
//		l.log.Info("PoolAlreadyBonded", "poolkey", pk)
//		return nil
//	}
//
//	chain, ok := l.chains[pk.Symbol]
//	if !ok {
//		return fmt.Errorf("no validator for symbol: %s", pk.Symbol)
//	}
//
//	err = chain.TryToBondForPool(pk.Pool)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
