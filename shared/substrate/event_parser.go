package substrate

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

var (
	ValueNotStringError = errors.New("value not string")
	ValueNotMapError    = errors.New("value not map")
	ValueNotU32         = errors.New("value not u32")
)

func LiquidityBondEventData(evt *ChainEvent) (*core.EvtLiquidityBond, error) {
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

func EraPoolUpdatedData(evt *ChainEvent) (*core.EraPoolUpdatedFlow, error) {
	if len(evt.Params) != 2 {
		return nil, fmt.Errorf("EraPoolUpdatedData params number not right: %d, expected: 2", len(evt.Params))
	}

	shot, err := parseHash(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[0] -> shot_id error: %s", err)
	}

	voter, err := parseAccountId(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[1] -> lastVoter error: %s", err)
	}

	return &core.EraPoolUpdatedFlow{
		ShotId:    shot,
		LastVoter: voter,
	}, nil
}

func EventNewMultisig(evt *ChainEvent) (*core.EventNewMultisig, error) {
	who, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNewMultisig params[0] -> who error: %s", err)
	}

	id, err := parseAccountId(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNewMultisig params[1] -> id error: %s", err)
	}

	hash, err := parseHash(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNewMultisig params[2] -> hash error: %s", err)
	}

	return &core.EventNewMultisig{
		Who:      who,
		ID:       id,
		CallHash: hash,
	}, nil
}

func EventMultisigExecuted(evt *ChainEvent) (*core.EventMultisigExecuted, error) {
	if len(evt.Params) != 5 {
		return nil, fmt.Errorf("EventMultisigExecuted params number not right: %d, expected: 5", len(evt.Params))
	}

	approving, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMultisigExecuted params[0] -> approving error: %s", err)
	}

	tp, err := parseTimePoint(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMultisigExecuted params[1] -> timepoint error: %s", err)
	}

	id, err := parseAccountId(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMultisigExecuted params[2] -> id error: %s", err)
	}

	hash, err := parseHash(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMultisigExecuted params[3] -> hash error: %s", err)
	}

	ok, err := parseDispatchResult(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("EventMultisigExecuted params[4] -> dispatchresult error: %s", err)
	}

	return &core.EventMultisigExecuted{
		Who:       approving,
		TimePoint: tp,
		ID:        id,
		CallHash:  hash,
		Result:    ok,
	}, nil
}

func EventBondReport(evt *ChainEvent) (*core.BondReportFlow, error) {
	if len(evt.Params) != 5 {
		return nil, fmt.Errorf("EventBondReport params number not right: %d, expected: 5", len(evt.Params))
	}

	shot, err := parseHash(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondReport params[0] -> shot_id error: %s", err)
	}

	symbol, err := parseRsymbol(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondReport params[1] -> rsymbol error: %s", err)
	}

	pool, err := parseBytes(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondReport params[2] -> pool error: %s", err)
	}

	era, err := parseEra(evt.Params[3])
	if err != nil {
		return nil, fmt.Errorf("EventBondReport params[3] -> era error: %s", err)
	}

	voter, err := parseAccountId(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondReport params[4] -> lastVoter error: %s", err)
	}

	return &core.BondReportFlow{
		ShotId:    shot,
		Rsymbol:   symbol,
		Pool:      pool,
		Era:       era.Value,
		LastVoter: voter,
	}, nil
}

func EventWithdrawUnbond(evt *ChainEvent) (*core.WithdrawUnbondFlow, error) {
	if len(evt.Params) != 5 {
		return nil, fmt.Errorf("EventBondReport params number not right: %d, expected: 5", len(evt.Params))
	}

	shot, err := parseHash(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondReport params[0] -> shot_id error: %s", err)
	}

	symbol, err := parseRsymbol(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondReport params[1] -> rsymbol error: %s", err)
	}

	pool, err := parseBytes(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondReport params[2] -> pool error: %s", err)
	}

	era, err := parseEra(evt.Params[3])
	if err != nil {
		return nil, fmt.Errorf("EventBondReport params[3] -> era error: %s", err)
	}

	voter, err := parseAccountId(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondReport params[4] -> lastVoter error: %s", err)
	}

	return &core.WithdrawUnbondFlow{
		ShotId:    shot,
		Rsymbol:   symbol,
		Pool:      pool,
		Era:       era.Value,
		LastVoter: voter,
	}, nil
}

func parseRsymbol(value interface{}) (core.RSymbol, error) {
	sym, ok := value.(string)
	if !ok {
		return core.RSymbol(""), ValueNotStringError
	}

	return core.RSymbol(sym), nil
}

func parseEra(param scalecodec.EventParam) (*Era, error) {
	bz, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	era := new(Era)
	err = json.Unmarshal(bz, era)
	if err != nil {
		return nil, err
	}

	return era, nil
}

func parseBytes(value interface{}) ([]byte, error) {
	val, ok := value.(string)
	if !ok {
		return nil, ValueNotStringError
	}

	bz, err := hexutil.Decode(utiles.AddHex(val))
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func parseBigint(value interface{}) (*big.Int, error) {
	val, ok := value.(string)
	if !ok {
		return nil, ValueNotStringError
	}

	i, ok := utils.StringToBigint(val)
	if !ok {
		return nil, fmt.Errorf("string to bigint error: %s", val)
	}

	return i, nil
}

func parseAccountId(value interface{}) (types.AccountID, error) {
	val, ok := value.(string)
	if !ok {
		return types.NewAccountID([]byte{}), ValueNotStringError
	}
	ac, err := hexutil.Decode(utiles.AddHex(val))
	if err != nil {
		return types.NewAccountID([]byte{}), err
	}

	return types.NewAccountID(ac), nil
}

func parseHash(value interface{}) (types.Hash, error) {
	val, ok := value.(string)
	if !ok {
		return types.NewHash([]byte{}), ValueNotStringError
	}

	hash, err := types.NewHashFromHexString(utiles.AddHex(val))
	if err != nil {
		return types.NewHash([]byte{}), err
	}

	return hash, err
}

func parseTimePoint(value interface{}) (types.TimePoint, error) {
	bz, err := json.Marshal(value)
	if err != nil {
		return types.TimePoint{}, err
	}

	var tp types.TimePoint
	err = json.Unmarshal(bz, &tp)
	if err != nil {
		return types.TimePoint{}, err
	}

	return tp, nil
}

func parseDispatchResult(value interface{}) (bool, error) {
	result, ok := value.(map[string]interface{})
	if !ok {
		return false, ValueNotMapError
	}
	_, ok = result["Ok"]
	return ok, nil
}
