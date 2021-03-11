package cosmos

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	subClientTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"math/big"
)

var _ conn.Chain = &FullClient{}

const eraBlockNumber = int64(3600) //6hours 6*60*60/6

//FullClient implement conn.Chain interface
type FullClient struct {
	Keys       []keyring.Info
	SubClients map[keyring.Info]*SubClient
	Log        log15.Logger
}

func (fc *FullClient) TransferVerify(r *conn.BondRecord) (conn.BondReason, error) {
	if len(fc.SubClients) == 0 || len(fc.Keys) == 0 {
		return "", fmt.Errorf("no subClient")
	}
	hashStr := hex.EncodeToString(r.Txhash)
	client := fc.SubClients[fc.Keys[0]]

	//check tx hash
	res, err := client.RpcClient.QueryTxByHash(hashStr)
	if err != nil {
		return conn.TxhashUnmatch, err
	}
	if res.Empty() {
		return conn.TxhashUnmatch, nil
	}

	//check block hash
	blockRes, err := client.RpcClient.QueryBlock(res.Height)
	if err != nil {
		return conn.BlockhashUnmatch, err
	}
	if !bytes.Equal(blockRes.BlockID.Hash, r.Blockhash) {
		return conn.BlockhashUnmatch, nil
	}

	//check amount and pool
	amountIsMatch := false
	poolIsMatch := false
	fromAddressStr := ""
	msgs := res.GetTx().GetMsgs()
	for i, _ := range msgs {
		if msgs[i].Type() == xBankTypes.TypeMsgSend {
			sendMsg, ok := msgs[i].(*xBankTypes.MsgSend)
			if ok {
				toAddr, err := types.AccAddressFromBech32(sendMsg.ToAddress)
				if err == nil {
					//amount and pool address must in one message
					if bytes.Equal(toAddr.Bytes(), r.Pool) &&
						sendMsg.Amount.AmountOf(client.RpcClient.GetDenom()).Equal(types.NewIntFromBigInt(r.Amount.Int)) {
						poolIsMatch = true
						amountIsMatch = true
						fromAddressStr = sendMsg.FromAddress
					}
				}

			}

		}
	}
	if !amountIsMatch {
		return conn.AmountUnmatch, nil
	}
	if !poolIsMatch {
		return conn.PoolUnmatch, nil
	}

	//check pubkey
	fromAddress, err := types.AccAddressFromBech32(fromAddressStr)
	if err != nil {
		return conn.PubkeyUnmatch, err
	}
	accountRes, err := client.RpcClient.QueryAccount(fromAddress)
	if err != nil {
		return conn.PubkeyUnmatch, err
	}

	if !bytes.Equal(accountRes.GetPubKey().Bytes(), r.Pubkey) {
		return conn.PubkeyUnmatch, nil
	}
	return conn.Pass, nil
}

func (fc *FullClient) CurrentEra() (subClientTypes.U32, error) {
	if len(fc.SubClients) == 0 || len(fc.Keys) == 0 {
		return 0, fmt.Errorf("no subClient")
	}
	client := fc.SubClients[fc.Keys[0]]
	status, err := client.RpcClient.GetStatus()
	if err != nil {
		return 0, err
	}
	era := status.SyncInfo.LatestBlockHeight / eraBlockNumber
	return subClientTypes.NewU32(uint32(era)), nil
}

func (fc *FullClient) BondWork(e *conn.EvtEraPoolUpdated) (*big.Int, error) {
	zero := big.NewInt(0)
	bond := e.Bond.Int
	unbond := e.Unbond.Int

	key := fc.foundKey(e.Pool)
	if key == nil {
		fc.Log.Info("no key for pool", "pool", hexutil.Encode(e.Pool))
		return nil, nil
	}

	if bond.Cmp(zero) == 0 && unbond.Cmp(zero) == 0 {
		fc.Log.Info("BondWork: bond and unbond are both zero")
		return nil, nil
	}

	subClient := fc.SubClients[key]

	fc.Log.Info("BondOrUnbond", "bond", bond, "unbond", unbond)
	if bond.Cmp(unbond) < 0 {
		diff := big.NewInt(0).Sub(unbond, bond)
		err := subClient.unbond(diff)
		if err != nil {
			return nil, err
		}
	} else if bond.Cmp(unbond) > 0 {
		diff := big.NewInt(0).Sub(bond, unbond)
		err := subClient.bond(diff)
		if err != nil {
			return nil, err
		}
	} else {
		fc.Log.Info("EvtEraPoolUpdated: bond is equal to unbond")
	}
	return nil, nil
}

func (fc *FullClient) foundKey(pool subClientTypes.Bytes) keyring.Info {
	for _, key := range fc.Keys {
		if bytes.Equal(key.GetAddress().Bytes(), pool) {
			return key
		}
	}
	return nil
}
