package substrate

import (
	"fmt"
	"math/big"

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
	sc := fsc.Sc

	if !sc.IsConnected() {
		if err := sc.WebsocketReconnect(); err != nil {
			panic(err)
		}
	}

	exts, err := sc.GetExtrinsics(bh)
	if err != nil {
		return conn.BlockhashUnmatch, nil
	}

	th := hexutil.Encode(r.Txhash)
	for _, ext := range exts {
		if th != utiles.AddHex(ext.ExtrinsicHash) {
			continue
		}

		if ext.CallModule.Name != TransferModuleName || (ext.Call.Name != TransferKeepAlive && ext.Call.Name != Transfer) {
			return conn.TxhashUnmatch, nil
		}

		addr, ok := ext.Address.(string)
		if !ok {
			sc.log.Error("TransferVerify: address not string", "address", ext.Address)
			return conn.PubkeyUnmatch, nil
		}

		if hexutil.Encode(r.Pubkey) != utiles.AddHex(addr) {
			return conn.PubkeyUnmatch, nil
		}

		for _, p := range ext.Params {
			if p.Name == ParamDest && p.Type == ParamDestType {
				sc.log.Debug("cmp dest", "pool", hexutil.Encode(r.Pool), "dest", p.Value)

				dest, ok := p.Value.(string)
				if !ok {
					dest, ok := p.Value.(map[string]interface{})
					if !ok {
						return conn.PoolUnmatch, nil
					}

					destId, ok := dest["Id"]
					if !ok {
						return conn.PoolUnmatch, nil
					}

					d, ok := destId.(string)
					if !ok {
						return conn.PoolUnmatch, nil
					}

					if hexutil.Encode(r.Pool) != utiles.AddHex(d) {
						return conn.PoolUnmatch, nil
					}
				} else {
					if hexutil.Encode(r.Pool) != utiles.AddHex(dest) {
						return conn.PoolUnmatch, nil
					}
				}
			} else if p.Name == ParamValue && p.Type == ParamValueType {
				sc.log.Debug("cmp amount", "amount", r.Amount, "paramAmount", p.Value)
				if fmt.Sprint(r.Amount) != fmt.Sprint(p.Value) {
					return conn.AmountUnmatch, nil
				}
			} else {
				sc.log.Error("TransferVerify unexpected param", "name", p.Name, "value", p.Value, "type", p.Type)
				return conn.TxhashUnmatch, nil
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
