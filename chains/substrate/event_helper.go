package substrate

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

func liquidityBondEventData(evt *substrate.ChainEvent) (*core.EvtLiquidityBond, error) {
	lb := new(core.EvtLiquidityBond)
	for _, p := range evt.Params {
		switch p.Type {
		case "AccountId":
			a := utiles.AddHex(p.Value.(string))
			aid, err := hexutil.Decode(a)
			if err != nil {
				return nil, err
			}
			lb.AccountId = types.NewAccountID(aid)
		case "RSymbol":
			sym := p.Value.(string)
			lb.Rsymbol = core.RSymbol(sym)
		case "Hash":
			b := p.Value.(string)
			bid, err := types.NewHashFromHexString(b)
			if err != nil {
				return nil, err
			}
			lb.BondId = bid
		default:
			continue
		}
	}

	return lb, nil
}

func eraPoolUpdatedData(evt *substrate.ChainEvent) (*core.EvtEraPoolUpdated, error) {
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
	bond, ok := utils.StringToBigint(bondStr)
	if !ok {
		return nil, errors.New("eraPoolUpdatedData bond error")
	}

	unbondStr, _ := evt.Params[4].Value.(string)
	unbond, ok := utils.StringToBigint(unbondStr)
	if !ok {
		return nil, errors.New("eraPoolUpdatedData bond error")
	}

	voterStr, _ := evt.Params[5].Value.(string)
	voter, err := hexutil.Decode(utiles.AddHex(voterStr))
	if err != nil {
		return nil, err
	}

	return &core.EvtEraPoolUpdated{
		Rsymbol:   core.RSymbol(sym),
		NewEra:    era.Value,
		Pool:      types.NewBytes(pool),
		Bond:      types.NewU128(*bond),
		Unbond:    types.NewU128(*unbond),
		LastVoter: voter,
	}, nil
}
