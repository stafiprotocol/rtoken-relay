package substrate

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/utils"
	"strings"
)

func (sc *SarpcClient) TransferVerify(r *conn.BondRecord) (conn.BondReason, error) {
	bh := hexutil.Encode(r.Blockhash)
	exts, err := sc.GetExtrinsics(bh)
	if err != nil {
		if strings.Contains(err.Error(), "websocket send error") {
			return conn.BondReason(""), err
		}
		return conn.BlockhashUnmatch, nil
	}

	th := hexutil.Encode(r.Txhash)
	for _, ext := range exts {
		if th != utiles.AddHex(ext.ExtrinsicHash) || !ext.ContainsTransaction ||
			ext.CallModule.Name != TransferModuleName || ext.Call.Name != TransferKeepAlive {
			continue
		}

		if hexutil.Encode(r.Pubkey) != utiles.AddHex(ext.Address) {
			return conn.PubkeyUnmatch, nil
		}

		for _, param := range ext.Params {
			if param.Name == ParamDest && param.Type == ParamDestType {
				if hexutil.Encode(r.Pool) != utiles.AddHex(param.ValueRaw) {
					return conn.PoolUnmatch, nil
				}
			} else if param.Name == ParamValue && param.Type == ParamValueType {
				v, ok := utils.FromString(param.ValueRaw)
				if !ok {
					return conn.AmountUnmatch, nil
				}
				if r.Amount.Cmp(v) != 0 {
					return conn.AmountUnmatch, nil
				}
			} else {
				return conn.BondReason(""), fmt.Errorf("got unexpected param: %+v", param)
			}
		}

	}

	return conn.TxhashUnmatch, nil
}
