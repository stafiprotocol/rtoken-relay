package substrate

import (
	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var (
	AliceKey     = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
	From         = "31yavGB5CVb8EwpqKQaS9XY7JZcfbK6QpWPn5kkweHVpqcov"
	LessPolka    = "1334v66HrtqQndbugYxX9m56V6222m97LbavB4KAMmqgjsas"
	From1        = "31d96Cq9idWQqPq3Ch5BFY84zrThVE3r98M7vG4xYaSWHwsX"
	From2        = "1TgYb5x8xjsZRyL5bwvxUoAWBn36psr4viSMHbRXA8bkB2h"
	Wen          = "1swvN162p1siDjm63UhhWoa59bpPZTSNKGVcbCwHUYkfRRW"
	Jun          = "33RQ73d9XfPTaE2SV7dzdhQQ17YaeMQ4yzhzAQhhFVenxMuJ"
	KeystorePath = "/Users/fwj/Go/stafi/rtoken-relay/keys"
)

var (
	tlog = log15.Root()
)

const (
	stafiTypesFile  = "/Users/tpkeeper/gowork/stafi/rtoken-relay/network/stafi.json"
	polkaTypesFile  = "/Users/fwj/Go/stafi/rtoken-relay/network/polkadot.json"
	kusamaTypesFile = "/Users/fwj/Go/stafi/rtoken-relay/network/kusama.json"
)

var sc *substrate.SarpcClient

func init() {
	var err error
	sc, err = substrate.NewSarpcClient(substrate.ChainTypePolkadot, "wss://kusama-rpc.polkadot.io", polkaTypesFile, tlog)
	if err != nil {
		panic(err)
	}
	stop := make(chan int)
	gc, err = substrate.NewGsrpcClient("wss://kusama-rpc.polkadot.io", substrate.AddressTypeMultiAddress, AliceKey, tlog, stop)
	if err != nil {
		panic(err)
	}
}

func TestConnection_TransferVerify(t *testing.T) {
	conn := Connection{sc: sc, gc: gc, log: tlog, symbol: core.RKSM, stop: make(chan int)}
	pool, _ := hexutil.Decode("0x0777f13da6fead588d8662fd63336b0008c4fbe4749e18779c1a4bd89ea50141")
	bonder, _ := hexutil.Decode("0xb2a90bcb80498a3286b57fadef0f1d2d0dabba41f977fe8f5117da9983a2592a")
	pubkey, _ := hexutil.Decode("0xb2a90bcb80498a3286b57fadef0f1d2d0dabba41f977fe8f5117da9983a2592a")
	blockHash, _ := hexutil.Decode("0x6df4292f19e8bbdb1d2563d877b262dac22e4307f98b29b249f7281bf971e72e")
	txHash, _ := hexutil.Decode("0xc8042e4a7e539b0dee10ad57dc2da483e0fd6f836224bf4acbe6e40547876f22")
	bond := &submodel.BondRecord{
		Pool:      types.NewBytes(pool),
		Bonder:    types.NewAccountID(bonder),
		Symbol:    core.RKSM,
		Pubkey:    types.NewBytes(pubkey),
		Blockhash: types.NewBytes(blockHash),
		Txhash:    types.NewBytes(txHash),
		Amount:    types.NewU128(*big.NewInt(1000000000)),
	}

	result, err := conn.TransferVerify(bond)
	assert.NoError(t, err)
	t.Log(result)
	assert.Equal(t, result, submodel.Pass)

	txHashErr, _ := hexutil.Decode("0xc8042e4a7e539b0dee10ad57dc2da483e0fd6f836224bf4acbe6e40547876f24")
	bondTxHashErr := &submodel.BondRecord{
		Pool:      types.NewBytes(pool),
		Bonder:    types.NewAccountID(bonder),
		Symbol:    core.RKSM,
		Pubkey:    types.NewBytes(pubkey),
		Blockhash: types.NewBytes(blockHash),
		Txhash:    types.NewBytes(txHashErr),
		Amount:    types.NewU128(*big.NewInt(1000000000)),
	}

	resultTxHashErr, err := conn.TransferVerify(bondTxHashErr)
	assert.NoError(t, err)
	t.Log(resultTxHashErr)
	assert.Equal(t, resultTxHashErr, submodel.TxhashUnmatch)

	bondAmountErr := &submodel.BondRecord{
		Pool:      types.NewBytes(pool),
		Bonder:    types.NewAccountID(bonder),
		Symbol:    core.RKSM,
		Pubkey:    types.NewBytes(pubkey),
		Blockhash: types.NewBytes(blockHash),
		Txhash:    types.NewBytes(txHash),
		Amount:    types.NewU128(*big.NewInt(100000000)),
	}
	resultBondAmountErr, err := conn.TransferVerify(bondAmountErr)
	assert.NoError(t, err)
	t.Log(resultBondAmountErr)
	assert.Equal(t, resultBondAmountErr, submodel.AmountUnmatch)

	blockHashErr, _ := hexutil.Decode("0x6df4292f19e8bbdb1d2563d877b262dac22e4307f98b29b249f7281bf971e78e")
	bondBlockHashErr := &submodel.BondRecord{
		Pool:      types.NewBytes(pool),
		Bonder:    types.NewAccountID(bonder),
		Symbol:    core.RKSM,
		Pubkey:    types.NewBytes(pubkey),
		Blockhash: types.NewBytes(blockHashErr),
		Txhash:    types.NewBytes(txHash),
		Amount:    types.NewU128(*big.NewInt(1000000000)),
	}

	resultBlockHasherr, err := conn.TransferVerify(bondBlockHashErr)
	assert.NoError(t, err)
	t.Log(resultBlockHasherr)
	assert.Equal(t, resultBlockHasherr, submodel.BlockhashUnmatch)

	pubkeyErr, _ := hexutil.Decode("0xb2a90bcb80498a3286b57fadef0f1d2d0dabba41f977fe8f5117da9983a2592e")
	bondPubkeyErr := &submodel.BondRecord{
		Pool:      types.NewBytes(pool),
		Bonder:    types.NewAccountID(bonder),
		Symbol:    core.RKSM,
		Pubkey:    types.NewBytes(pubkeyErr),
		Blockhash: types.NewBytes(blockHash),
		Txhash:    types.NewBytes(txHash),
		Amount:    types.NewU128(*big.NewInt(1000000000)),
	}

	resultPubkeyErr, err := conn.TransferVerify(bondPubkeyErr)
	assert.NoError(t, err)
	t.Log(resultPubkeyErr)
	assert.Equal(t, resultPubkeyErr, submodel.PubkeyUnmatch)

	poolErr, _ := hexutil.Decode("0x0777f13da6fead588d8662fd63336b0008c4fbe4749e18779c1a4bd89ea5014e")
	bondPoolErr := &submodel.BondRecord{
		Pool:      types.NewBytes(poolErr),
		Bonder:    types.NewAccountID(bonder),
		Symbol:    core.RKSM,
		Pubkey:    types.NewBytes(pubkey),
		Blockhash: types.NewBytes(blockHash),
		Txhash:    types.NewBytes(txHash),
		Amount:    types.NewU128(*big.NewInt(1000000000)),
	}

	resultPoolErr, err := conn.TransferVerify(bondPoolErr)
	assert.NoError(t, err)
	t.Log(resultPoolErr)
	assert.Equal(t, resultPoolErr, submodel.PoolUnmatch)
}

func TestConnection_GetEvents(t *testing.T) {
	_, err := sc.GetEvents(7111716)
	assert.NoError(t, err)

	for i := 7111000; i < 7112000; i++ {
		evts, err := sc.GetEvents(uint64(i))
		assert.NoError(t, err)
		for _, evt := range evts {
			t.Log("eventId", evt.EventId)
			t.Log("moduleId", evt.ModuleId)
			t.Log("params", evt.Params)
		}
	}
}

func TestConnection_GetExtrinsics(t *testing.T) {

	for i := 7111000; i < 7112000; i++ {
		bh, err := sc.GetBlockHash(uint64(i))
		assert.NoError(t, err)
		exts, err := sc.GetExtrinsics(bh)
		assert.NoError(t, err)
		for _, ext := range exts {
			t.Log("\n=============", i)
			t.Log(ext.ExtrinsicHash)
			t.Log(ext.Address)
			t.Log(ext.CallModuleName)
			t.Log(ext.CallName)
			t.Log(ext.Params)
		}

	}
}

func TestConnection_PaymentQueryInfo(t *testing.T) {
	info, err := sc.GetPaymentQueryInfo("0x74b58b2d6fbd6319e4cc7927f1b789d48fc2629437cf9373bf8224934b831f58")
	assert.NoError(t, err)
	t.Log(info)
}
