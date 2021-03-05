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

type Era struct {
	Type  string `json:"type"`
	Value uint32 `json:"value"`
}

//type evtPoolSubAccountAdded struct {
//	symbol     conn.RSymbol
//	pool       types.Bytes
//	subAccount types.Bytes
//}
//
//func (ps *evtPoolSubAccountAdded) poolKey() *conn.PoolBondKey {
//	return &conn.PoolBondKey{
//		Symbol: ps.symbol,
//		Pool:   ps.pool,
//	}
//}

//type SubClients []*SubClient
//
//type SubClient struct {
//	KeyPair *signature.KeyringPair
//	ChainClient substrate.ChainClient
//}
