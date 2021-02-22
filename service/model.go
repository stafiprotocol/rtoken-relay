package service

import (
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/conn"
)

type evtLiquidityBond struct {
	accountId types.AccountID
	symbol    conn.RSymbol
	bondId    types.Hash
}

func (lb *evtLiquidityBond) bondKey() *conn.BondKey {
	return &conn.BondKey{
		Symbol: lb.symbol,
		BondId: lb.bondId,
	}
}
