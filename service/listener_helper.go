package service

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/substrate"
)

func liquidityBondEventData(evt *substrate.ChainEvent) (*evtLiquidityBond, error) {
	lb := new(evtLiquidityBond)
	for _, p := range evt.Params {
		switch p.Type {
		case "AccountId":
			a := utiles.AddHex(p.Value.(string))
			aid, err := hexutil.Decode(a)
			if err != nil {
				return nil, err
			}
			lb.accountId = types.NewAccountID(aid)
		case "RSymbol":
			sym := p.Value.(string)
			lb.symbol = conn.RSymbol(sym)
		case "Hash":
			b := p.Value.(string)
			bid, err := types.NewHashFromHexString(b)
			if err != nil {
				return nil, err
			}
			lb.bondId = bid
		}
	}

	return lb, nil
}

func containsVote(votes []types.AccountID, voter types.AccountID) bool {
	for _, v := range votes {
		if bytes.Equal(v[:], voter[:]) {
			return true
		}
	}
	return false
}
