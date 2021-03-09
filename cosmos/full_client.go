package cosmos

import (
	"bytes"
	"encoding/hex"
	"fmt"
	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	subClientTypes "github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/cosmos/rpc"
	"math/big"
)

var _ conn.Chain = &FullClient{}

//FullClient implement conn.Chain interface
type FullClient struct {
	keys       []keyring.Info
	subClients map[keyring.Info]rpc.Client
}

func (fc *FullClient) TransferVerify(r *conn.BondRecord) (conn.BondReason, error) {
	if len(fc.subClients) == 0 || len(fc.keys) == 0 {
		return "", fmt.Errorf("no subClient")
	}
	hashStr := hex.EncodeToString(r.Txhash)
	client := fc.subClients[fc.keys[0]]

	//check tx hash
	res, err := client.QueryTxByHash(hashStr)
	if err != nil {
		return conn.TxhashUnmatch, nil
	}
	if res.Empty() {
		return conn.TxhashUnmatch, nil
	}

	//check block hash
	blockRes, err := client.QueryBlock(res.Height)
	if err != nil {
		return conn.BlockhashUnmatch, nil
	}
	if !bytes.Equal(blockRes.BlockID.Hash, r.Blockhash) {
		return conn.BlockhashUnmatch, nil
	}

	//todo use another way to check pubkey
	//check pubkey
	pubkeyIsMatch := false
	txI, err := client.GetTxConfig().TxDecoder()(res.Tx.Value)
	if err == nil {
		tx, ok := txI.(signing.Tx)
		if ok {
			stdTx, err := clientTx.ConvertTxToStdTx(client.GetLegacyAmino(), tx)
			if err == nil && len(stdTx.Signatures) != 0 {
				pubKeyBts := stdTx.Signatures[0].GetPubKey().Bytes()
				if bytes.Equal(pubKeyBts, r.Pubkey) {
					pubkeyIsMatch = true
				}
			}
		}
	}
	if !pubkeyIsMatch {
		return conn.PubkeyUnmatch, nil
	}

	//check amount and pool
	amountIsMatch := false
	poolIsMatch := false
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

	return conn.Pass, nil
}

func (fc *FullClient) CurrentEra() (subClientTypes.U32, error) {
	return 0, nil
}

func (fc *FullClient) BondWork(e *conn.EvtEraPoolUpdated) (*big.Int, error) {
	return nil, nil
}
