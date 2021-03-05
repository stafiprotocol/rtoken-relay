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
				dest, _ := param.Value.(string)
				fsc.Gc.log.Debug("cmp dest", "pool", hexutil.Encode(r.Pool), "dest", param.Value)
				if hexutil.Encode(r.Pool) != utiles.AddHex(dest) {
					return conn.PoolUnmatch, nil
				}
			} else if param.Name == ParamValue && param.Type == ParamValueType {
				fsc.Gc.log.Debug("cmp amount", "amount", r.Amount, "paramAmount", param.Value)
				if fmt.Sprint(r.Amount) != fmt.Sprint(param.Value) {
					return conn.AmountUnmatch, nil
				}
			}
		}
		return conn.Pass, nil
	}

	return conn.TxhashUnmatch, nil
}

func (fsc *FullSubClient) CurrentEra() (types.U32, error) {
	if len(fsc.Keys) == 0 || len(fsc.SubClients) == 0 {
		return 0, fmt.Errorf("no keys, no subclients, no current era")
	}

	gc := fsc.SubClients[fsc.Keys[0]]
	var index types.U32
	exist, err := gc.QueryStorage(StakingModuleId, StorageActiveEra, nil, nil, &index)
	if err != nil {
		return 0, err
	}

	if !exist {
		return 0, fmt.Errorf("unable to get activeEraInfo for: %s", gc.endpoint)
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
	if err != nil {
		return nil, err
	}

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
