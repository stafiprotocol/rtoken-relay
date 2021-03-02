package substrate

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

const (
	StakingModuleId  = "Staking"
	StorageActiveEra = "ActiveEra"
	StorageLegder    = "Ledger"
	MethodUnbond     = "Staking.unbond"
	MethodBondExtra  = "Staking.bond_extra"
)

func (fsc *FullSubClient) TransferVerify(r *conn.BondRecord) (conn.BondReason, error) {
	bh := hexutil.Encode(r.Blockhash)
	exts, err := fsc.Sc.GetExtrinsics(bh)
	if err != nil {
		if strings.Contains(err.Error(), "websocket send error") {
			return conn.BondReason(""), err
		}
		return conn.BlockhashUnmatch, nil
	}

	th := hexutil.Encode(r.Txhash)
	for _, ext := range exts {
		if th != utiles.AddHex(ext.ExtrinsicHash) || ext.CallModule.Name != TransferModuleName ||
			ext.Call.Name != TransferKeepAlive {
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

func (fsc *FullSubClient) CurrentEra() (types.U32, error) {
	var index types.U32
	if len(fsc.SubClients) == 0 {
		return 0, fmt.Errorf("no subclients")
	}

	exist, err := fsc.Gc.QueryStorage(StakingModuleId, StorageActiveEra, nil, nil, &index)
	if err != nil {
		return 0, err
	}

	if !exist {
		return 0, fmt.Errorf("unable to get activeEraInfo")
	}

	return index, nil
}

func (fsc *FullSubClient) BondWork(e *conn.EvtEraPoolUpdated) (*big.Int, error) {
	key := fsc.foundKey(e.Pool)
	if key == nil {
		fsc.Gc.log.Info("no key for pool", "pool", hexutil.Encode(e.Pool))
		return nil, nil
	}

	gc := fsc.SubClients[key]
	err := gc.BondOrUnbond(e.Bond.Int, e.Unbond.Int)
	if err != nil {
		return nil, err
	}

	addr := types.NewAddressFromAccountID(e.Pool)
	s, err := gc.StakingActive(addr.AsAccountID)

	return (*big.Int)(&s.Active), nil
}

func (fsc *FullSubClient) foundKey(pool types.Bytes) *signature.KeyringPair {
	for _, key := range fsc.Keys {
		if hexutil.Encode(key.PublicKey) == hexutil.Encode(pool) {
			return key
		}
	}

	return nil
}
