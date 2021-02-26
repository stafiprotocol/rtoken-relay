package substrate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

const (
	StakingModuleId  = "Staking"
	StorageActiveEra = "ActiveEra"
	BondExtraMethod  = "bond_extra"
	UnBondMethod     = "bond_extra"
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

func (fsc *FullSubClient) BondWork(ck *conn.ChunkKey) error {
	if len(fsc.SubClients) == 0 {
		return errors.New("FullSubClient BondWork has no subclient")
	}

	ckbz, err := types.EncodeToBytes(ck)
	if err != nil {
		fsc.Gc.log.Error("FullSubClient BondWork", "EncodeToBytes error", err)
		return err
	}

	lcs := make([]*conn.LinkChunk, 0)
	exist, err := fsc.Gc.QueryStorage(config.RTokenLedgerModuleId, config.StorageTotalLinking, ckbz, nil, lcs)
	if err != nil {
		return fmt.Errorf("BondWork QueryStorage error: %s", err)
	}

	if !exist {
		fsc.Gc.log.Info("FullSubClient BondWork Nothing to bond")
		return nil
	}

	for _, lc := range lcs {
		key := fsc.foundKey(lc.Pool)
		if key == nil {
			continue
		}

		gc := fsc.SubClients[key]
		err := gc.TryToBondOrUnbond(lc)
		if err != nil {
			return err
		}

		err = gc.TryToClaim(lc)
		if err != nil {
			return err
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
