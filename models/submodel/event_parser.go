package submodel

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
	ValueNotStringError      = errors.New("value not string")
	ValueNotMapError         = errors.New("value not map")
	ValueNotU32              = errors.New("value not u32")
	ValueNotStringSliceError = errors.New("value not string slice")
)

func LiquidityBondEventData(evt *ChainEvent) (*EvtLiquidityBond, error) {
	if len(evt.Params) != 3 {
		return nil, fmt.Errorf("EvtLiquidityBond params number not right: %d, expected: 3", len(evt.Params))
	}
	accountId, err := parseAccountId(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EvtLiquidityBond params[0] -> AccountId error: %s", err)
	}
	symbol, err := parseRsymbol(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EvtLiquidityBond params[1] -> RSymbol error: %s", err)
	}
	bondId, err := parseHash(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EvtLiquidityBond params[2] -> BondId error: %s", err)
	}

	return &EvtLiquidityBond{
		AccountId: accountId,
		Symbol:    symbol,
		BondId:    bondId,
	}, nil
}

func EraPoolUpdatedData(evt *ChainEvent) (*EraPoolUpdatedFlow, error) {
	if len(evt.Params) != 4 {
		return nil, fmt.Errorf("EraPoolUpdatedData params number not right: %d, expected: 4", len(evt.Params))
	}

	symbol, err := parseRsymbol(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[0] -> RSymbol error: %s", err)
	}

	era, err := parseEra(evt.Params[1])
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[1] -> era error: %s", err)
	}

	shot, err := parseHash(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[2] -> shot_id error: %s", err)
	}

	voter, err := parseAccountId(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[3] -> lastVoter error: %s", err)
	}

	return &EraPoolUpdatedFlow{
		Symbol:    symbol,
		Era:       era.Value,
		ShotId:    shot,
		LastVoter: voter,
		Active:    big.NewInt(0),
	}, nil
}

func EventNewMultisigData(evt *ChainEvent) (*EventNewMultisig, error) {
	if len(evt.Params) != 3 {
		return nil, fmt.Errorf("EventNewMultisigData params number not right: %d, expected: 3", len(evt.Params))
	}
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

	return &EventNewMultisig{
		Who:      who,
		ID:       id,
		CallHash: hash,
	}, nil
}

func EventMultisigExecutedData(evt *ChainEvent) (*EventMultisigExecuted, error) {
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

	return &EventMultisigExecuted{
		Who:       approving,
		TimePoint: tp,
		ID:        id,
		CallHash:  hash,
		Result:    ok,
	}, nil
}

func EventBondReported(evt *ChainEvent) (*BondReportedFlow, error) {
	if len(evt.Params) != 3 {
		return nil, fmt.Errorf("EventBondReported params number not right: %d, expected: 3", len(evt.Params))
	}

	symbol, err := parseRsymbol(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondReported params[0] -> RSymbol error: %s", err)
	}

	shot, err := parseHash(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondReported params[1] -> shot_id error: %s", err)
	}

	voter, err := parseAccountId(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventBondReported params[2] -> lastVoter error: %s", err)
	}

	return &BondReportedFlow{
		Symbol:    symbol,
		ShotId:    shot,
		LastVoter: voter,
	}, nil
}

func EventActiveReported(evt *ChainEvent) (*ActiveReportedFlow, error) {
	if len(evt.Params) != 3 {
		return nil, fmt.Errorf("EventActiveReported params number not right: %d, expected: 3", len(evt.Params))
	}

	symbol, err := parseRsymbol(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventActiveReported params[0] -> RSymbol error: %s", err)
	}

	shot, err := parseHash(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventActiveReported params[1] -> shot_id error: %s", err)
	}

	voter, err := parseAccountId(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventActiveReported params[2] -> lastVoter error: %s", err)
	}

	return &ActiveReportedFlow{
		Symbol:    symbol,
		ShotId:    shot,
		LastVoter: voter,
	}, nil
}

func EventWithdrawReported(evt *ChainEvent) (*WithdrawReportedFlow, error) {
	if len(evt.Params) != 3 {
		return nil, fmt.Errorf("EventWithdrawReported params number not right: %d, expected: 3", len(evt.Params))
	}

	symbol, err := parseRsymbol(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventWithdrawReported params[0] -> RSymbol error: %s", err)
	}

	shot, err := parseHash(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventWithdrawReported params[1] -> shot_id error: %s", err)
	}

	voter, err := parseAccountId(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventWithdrawReported params[2] -> lastVoter error: %s", err)
	}

	return &WithdrawReportedFlow{
		Symbol:    symbol,
		ShotId:    shot,
		LastVoter: voter,
	}, nil
}

func EventTransferReported(evt *ChainEvent) (*TransferReportedFlow, error) {
	if len(evt.Params) != 2 {
		return nil, fmt.Errorf("EventTransferReported params number not right: %d, expected: 3", len(evt.Params))
	}

	symbol, err := parseRsymbol(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventTransferReported params[0] -> RSymbol error: %s", err)
	}

	shot, err := parseHash(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventTransferReported params[1] -> shot_id error: %s", err)
	}

	return &TransferReportedFlow{
		Symbol: symbol,
		ShotId: shot,
	}, nil
}

func EventNominationUpdated(evt *ChainEvent) (*NominationUpdatedFlow, error) {
	if len(evt.Params) != 5 {
		return nil, fmt.Errorf("EventNominationUpdated params number not right: %d, expected: 5", len(evt.Params))
	}

	symbol, err := parseRsymbol(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNominationUpdated params[0] -> RSymbol error: %s", err)
	}

	pool, err := parseBytes(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNominationUpdated params[1] -> pool error: %s", err)
	}

	vals, err := parseVecBytes(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNominationUpdated params[2] -> new_validators error: %s", err)
	}

	era, err := parseEra(evt.Params[3])
	if err != nil {
		return nil, fmt.Errorf("EventNominationUpdated params[3] -> era error: %s", err)
	}

	voter, err := parseAccountId(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNominationUpdated params[4] -> lastVoter error: %s", err)
	}

	return &NominationUpdatedFlow{
		Symbol:        symbol,
		Pool:          pool,
		NewValidators: vals,
		Era:           era.Value,
		LastVoter:     voter,
	}, nil
}

func EventValidatorUpdated(evt *ChainEvent) (*ValidatorUpdatedFlow, error) {
	if len(evt.Params) != 5 {
		return nil, fmt.Errorf("EventNominationUpdated params number not right: %d, expected: 5", len(evt.Params))
	}

	symbol, err := parseRsymbol(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNominationUpdated params[0] -> RSymbol error: %s", err)
	}

	pool, err := parseBytes(evt.Params[1].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNominationUpdated params[1] -> pool error: %s", err)
	}

	oldVal, err := parseBytes(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNominationUpdated params[2] -> old_validators error: %s", err)
	}

	newVal, err := parseBytes(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("EventNominationUpdated params[3] -> new_validators error: %s", err)
	}

	era, err := parseEra(evt.Params[4])
	if err != nil {
		return nil, fmt.Errorf("EventNominationUpdated params[4] -> era error: %s", err)
	}

	return &ValidatorUpdatedFlow{
		Symbol:       symbol,
		Pool:         pool,
		OldValidator: oldVal,
		NewValidator: newVal,
		Era:          era.Value,
	}, nil
}

func SignatureEnoughData(evt *ChainEvent) (*EvtSignatureEnough, error) {
	if len(evt.Params) != 5 {
		return nil, fmt.Errorf("params number not right: %d, expected: 6", len(evt.Params))
	}

	sym, err := parseRsymbol(evt.Params[0].Value)
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[0] -> rsymbol error: %s", err)
	}

	era, err := parseEra(evt.Params[1])
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[1] -> era error: %s", err)
	}

	pool, err := parseBytes(evt.Params[2].Value)
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[2] -> pool error: %s", err)
	}

	tx, err := parseOriginTx(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[3] -> bond error: %s", err)
	}

	proposalId, err := parseBytes(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("EraPoolUpdatedData params[4] -> unbond error: %s", err)
	}

	return &EvtSignatureEnough{
		RSymbol:    sym,
		Era:        era.Value,
		Pool:       pool,
		TxType:     tx,
		ProposalId: proposalId,
	}, nil
}

func parseOriginTx(value interface{}) (OriginalTx, error) {
	sym, ok := value.(string)
	if !ok {
		return OriginalTx(""), ValueNotStringError
	}

	return OriginalTx(sym), nil
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
		if err.Error() == hexutil.ErrSyntax.Error() {
			return []byte(val), nil
		}
		return nil, err
	}

	return bz, nil
}

func parseVecBytes(value interface{}) ([]types.Bytes, error) {
	vals, ok := value.([]interface{})
	if !ok {
		return nil, ValueNotStringSliceError
	}
	result := make([]types.Bytes, 0)
	for _, val := range vals {
		bz, err := parseBytes(val)
		if err != nil {
			return nil, err
		}

		result = append(result, bz)
	}

	return result, nil
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
