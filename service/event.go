package service

import (
	"fmt"
	"github.com/stafiprotocol/rtoken-relay/utils"

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

func (l *listener) processEraPoolUpdatedEvts(evts []*substrate.ChainEvent) error {
	for _, evt := range evts {

		err := l.processEraPoolUpdatedEvt(evt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *listener) processEraPoolUpdatedEvt(evt *substrate.ChainEvent) error {
	data, err := eraPoolUpdatedData(evt)
	if err != nil {
		return err
	}

	if data.NewEra < l.chainEras[data.Symbol] {
		l.log.Info("era_pool_updated_event of past era, ignored", "current",
			l.chainEras[data.Symbol], "eventEra", data.NewEra, "symbol", data.Symbol)
		return nil
	}

	l.log.Info("processEraPoolUpdatedEvt", "data", data)

	chain, ok := l.chains[data.Symbol]
	if !ok {
		return fmt.Errorf("no validator for symbol: %s", data.Symbol)
	}

	active, err := chain.BondWork(data)
	if err != nil {
		return err
	}

	if active == nil {
		l.log.Info("no need to bond")
		return nil
	}

	h, err := utils.Blake2Hash(active)
	if err != nil {
		return err
	}

	bk := &conn.BondKey{data.Symbol, h}
	prop, err := l.newSetPoolActiveProposal(bk, data.Symbol, data.NewEra, data.Pool, types.NewU128(*active))
	if err != nil {
		return err
	}

	result := l.resolveProposal(prop, true)
	l.log.Info("processEraPoolUpdatedEvt", "result", result)

	return nil
}
