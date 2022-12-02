package substrate

import (
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
)

var (
	ErrormultiEnd             = errors.New("multiEnd")
	ErrorEraSnapShotsNotExist = errors.New("eraSnapShots not exist")

	ErrorEventEraIsOld                = errors.New("ErrorEventEraIsOld")
	ErrorBondStateNotEraUpdated       = errors.New("ErrorBondStateNotEraUpdated")
	ErrorBondStateNotBondReported     = errors.New("ErrorBondStateNotBondReported")
	ErrorBondStateNotActiveReported   = errors.New("ErrorBondStateNotActiveReported")
	ErrorBondStateNotWithdrawReported = errors.New("ErrorBondStateNotWithdrawReported")
	ErrorBondStateNotTransferReported = errors.New("ErrorBondStateNotTransferReported")
	ErrNotCared                       = errors.New("not care this symbol")
)

func (l *listener) processNewMultisigEvt(evt *submodel.ChainEvent) (*submodel.EventNewMultisig, error) {
	data, err := submodel.EventNewMultisigData(evt)
	if err != nil {
		return nil, err
	}

	mul := new(submodel.Multisig)
	exist, err := l.conn.QueryStorage(config.MultisigModuleId, config.StorageMultisigs, data.ID[:], data.CallHash[:], mul)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, ErrormultiEnd
	}

	data.TimePoint = submodel.NewOptionTimePoint(mul.When)
	data.Approvals = mul.Approvals
	data.CallHashStr = hexutil.Encode(data.CallHash[:])
	return data, nil
}

func (l *listener) processMultisigExecutedEvt(evt *submodel.ChainEvent) (*submodel.EventMultisigExecuted, error) {
	data, err := submodel.EventMultisigExecutedData(evt)
	if err != nil {
		return nil, err
	}
	data.CallHashStr = hexutil.Encode(data.CallHash[:])
	return data, nil
}
