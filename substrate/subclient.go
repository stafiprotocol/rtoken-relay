package substrate

import (
	"errors"
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
	BondExtraMethod  = "bond_extra"
	UnBondMethod     = "unbond"
	StorageLedger    = "Ledger"
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

func (fsc *FullSubClient) BondWork(plcs []*conn.PoolLinkChunk) error {
	if len(fsc.SubClients) == 0 {
		return errors.New("FullSubClient BondWork has no subclient")
	}

	calls := make([]types.Call, 0)
	var meta *types.Metadata
	for _, plc := range plcs {
		bond := plc.Bond.Int
		unbond := plc.Unbond.Int

		zero := big.NewInt(0)
		if bond.Cmp(zero) == 0 && unbond.Cmp(zero) == 0 {
			continue
		}

		key := fsc.foundKey(plc.Pool)
		if key == nil {
			continue
		}

		gc := fsc.SubClients[key]
		if meta == nil {
			latestMeta, err := gc.GetLatestMetadata()
			if err != nil {
				return err
			}
			meta = latestMeta
		}

		if bond.Cmp(unbond) < 0 {
			diff := big.NewInt(0).Sub(unbond, bond)
			realUnbond := types.NewU128(*diff)
			call, err := types.NewCall(meta, StakingModuleId, UnBondMethod, realUnbond)
			if err != nil {
				return err
			}
			calls = append(calls, call)
			continue
		} else if bond.Cmp(unbond) > 0 {
			diff := big.NewInt(0).Sub(bond, unbond)
			realBond := types.NewU128(*diff)
			call, err := types.NewCall(meta, StakingModuleId, BondExtraMethod, realBond)
			if err != nil {
				return err
			}
			calls = append(calls, call)
			continue
		} else {
			gc.log.Info("BondWork: bond is equal to unbond", "bond", bond, "unbond", unbond)
		}
	}


	return nil
}

func (fsc *FullSubClient) foundKey(pool types.Bytes) *signature.KeyringPair {
	for _, key := range fsc.Keys {
		if hexutil.Encode(key.PublicKey) == hexutil.Encode(pool) {
			return key
		}
	}

	return nil
}
