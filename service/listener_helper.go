package service

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/substrate"
	"github.com/stafiprotocol/rtoken-relay/utils"
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

func eraPoolUpdatedData(evt *substrate.ChainEvent) (*conn.EvtEraPoolUpdated, error) {
	sym, ok := evt.Params[0].Value.(string)
	if !ok {
		return nil, errors.New("eraPoolUpdatedData symbol error")
	}

	era := new(Era)
	x, _ := json.Marshal(evt.Params[1])
	err := json.Unmarshal(x, &era)
	if err != nil {
		return nil, err
	}

	poolStr := utiles.AddHex(evt.Params[2].Value.(string))
	pool, err := hexutil.Decode(poolStr)
	if err != nil {
		return nil, err
	}

	bondStr, _ := evt.Params[3].Value.(string)
	bond, ok := utils.FromString(bondStr)
	if !ok {
		return nil, errors.New("eraPoolUpdatedData bond error")
	}

	unbondStr, _ := evt.Params[4].Value.(string)
	unbond, ok := utils.FromString(unbondStr)
	if !ok {
		return nil, errors.New("eraPoolUpdatedData bond error")
	}

	return &conn.EvtEraPoolUpdated{
		Symbol: conn.RSymbol(sym),
		NewEra: types.NewU32(era.Value),
		Pool:   types.NewBytes(pool),
		Bond:   types.NewU128(*bond),
		Unbond: types.NewU128(*unbond),
	}, nil
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
