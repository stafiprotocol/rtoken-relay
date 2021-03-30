package substrate

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/utils"
	"math/big"
	"time"
)

func (c *Connection) submitSignature(param *submodel.SubmitSignatureParams) bool {
	for i := 0; i < BlockRetryLimit; i++ {
		c.log.Info("submitSignature on chain...")
		ext, err := c.gc.NewUnsignedExtrinsic(config.SubmitSignatures, param.Symbol, param.Era, param.Pool,
			param.TxType, param.ProposalId, param.Signature)
		err = c.gc.SignAndSubmitTx(ext)
		if err != nil {
			if err.Error() == TerminatedError.Error() {
				c.log.Error("submitSignature  met TerminatedError")
				return false
			}
			c.log.Error("submitSignature error", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}
		return true
	}
	return true
}

func (w *writer) unbondings(symbol core.RSymbol, pool []byte, era uint32) ([]*submodel.Receive, types.U128, error) {
	bz, err := types.EncodeToBytes(struct {
		core.RSymbol
		types.Bytes
		uint32
	}{symbol, pool, era})
	if err != nil {
		return nil, types.U128{}, err
	}

	unbonds := make([]*submodel.Unbonding, 0)
	exist, err := w.conn.QueryStorage(config.RTokenLedgerModuleId, config.StoragePoolUnbonds, bz, nil, &unbonds)
	if err != nil {

	}
	if !exist {
		return nil, types.U128{}, fmt.Errorf("pool unbonds not exist, symbol: %s, pool: %s, era: %d", symbol, hexutil.Encode(pool), era)
	}

	amounts := make(map[string]types.U128)
	for _, ub := range unbonds {
		rec := hexutil.Encode(ub.Recipient)
		acc, ok := amounts[rec]
		if !ok {
			amounts[rec] = ub.Value
		} else {
			amounts[rec] = utils.AddU128(acc, ub.Value)
		}
	}

	receives := make([]*submodel.Receive, 0)
	total := types.NewU128(*big.NewInt(0))
	for k, v := range amounts {
		r, _ := hexutil.Decode(k)
		rec := &submodel.Receive{Recipient: types.NewAddressFromAccountID(r), Value: types.NewUCompact(v.Int)}
		receives = append(receives, rec)
		total = utils.AddU128(total, v)
	}

	return receives, total, nil
}
