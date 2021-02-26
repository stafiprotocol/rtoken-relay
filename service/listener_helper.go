package service

import (
	"bytes"
	"errors"

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
		default:
			continue
		}
	}

	return lb, nil
}

func eraUpdatedEventData(evt *substrate.ChainEvent) (*evtEraUpdated, error) {
	//eu := new(evtEraUpdated)
	sym, ok := evt.Params[0].Value.(string)
	if !ok {
		return nil, errors.New("eraUpdatedEventData symbol error")
	}

	oldEra, ok := evt.Params[1].Value.(uint32)
	if !ok {
		return nil, errors.New("eraUpdatedEventData oldEra error")
	}

	newEra, ok := evt.Params[2].Value.(uint32)
	if !ok {
		return nil, errors.New("eraUpdatedEventData newEra error")
	}

	return &evtEraUpdated{conn.RSymbol(sym), types.NewU32(oldEra), types.NewU32(newEra)}, nil
}

//
//func poolSubAccountAddedEventData(evt *substrate.ChainEvent) (*evtPoolSubAccountAdded, error) {
//	sym, ok := evt.Params[0].Value.(string)
//
//	if !ok {
//		return nil, errors.New("poolSubAccountAddedEventData symbol error")
//	}
//
//	poolStr := utiles.AddHex(evt.Params[1].Value.(string))
//	pool, err := hexutil.Decode(poolStr)
//	if err != nil {
//		return nil, err
//	}
//
//	subStr := utiles.AddHex(evt.Params[2].Value.(string))
//	sub, err := hexutil.Decode(subStr)
//	if err != nil {
//		return nil, err
//	}
//
//	return &evtPoolSubAccountAdded{conn.RSymbol(sym), pool, sub}, nil
//}

func containsVote(votes []types.AccountID, voter types.AccountID) bool {
	for _, v := range votes {
		if bytes.Equal(v[:], voter[:]) {
			return true
		}
	}
	return false
}
