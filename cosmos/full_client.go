package cosmos

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	subClientTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/cosmos/rpc"
	"math/big"
)

var _ conn.Chain = &FullClient{}
var eraBlockNumber = int64(3600) //6hours 6*60*60/6

//FullClient implement conn.Chain interface
type FullClient struct {
	Keys       []keyring.Info
	SubClients map[keyring.Info]*rpc.Client
}

func (fc *FullClient) TransferVerify(r *conn.BondRecord) (conn.BondReason, error) {
	if len(fc.SubClients) == 0 || len(fc.Keys) == 0 {
		return "", fmt.Errorf("no subClient")
	}
	hashStr := hex.EncodeToString(r.Txhash)
	client := fc.SubClients[fc.Keys[0]]

	//check tx hash
	res, err := client.QueryTxByHash(hashStr)
	if err != nil {
		return conn.TxhashUnmatch, err
	}
	if res.Empty() {
		return conn.TxhashUnmatch, nil
	}

	//check block hash
	blockRes, err := client.QueryBlock(res.Height)
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
						sendMsg.Amount.AmountOf(client.GetDenom()).Equal(types.NewIntFromBigInt(r.Amount.Int)) {
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
	accountRes, err := client.QueryAccount(fromAddress)
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
	status, err := client.GetStatus()
	if err != nil {
		return 0, err
	}
	era := status.SyncInfo.LatestBlockHeight / eraBlockNumber
	return subClientTypes.NewU32(uint32(era)), nil
}

func (fc *FullClient) BondWork(e *conn.EvtEraPoolUpdated) (*big.Int, error) {

	return nil, nil
}
