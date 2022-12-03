package substrate

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"sort"
	"testing"

	sr25519 "github.com/ChainSafe/go-schnorrkel"
	"github.com/JFJun/go-substrate-crypto/ss58"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/submodel"
	"github.com/stafiprotocol/rtoken-relay/shared/substrate"
	"github.com/stretchr/testify/assert"
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
	tlog = core.NewLog()
)

const (
	stafiTypesFile  = "/Users/tpkeeper/gowork/stafi/rtoken-relay/network/stafi.json"
	polkaTypesFile  = "/Users/fwj/Go/stafi/rtoken-relay/network/polkadot.json"
	kusamaTypesFile = "/Users/fwj/Go/stafi/rtoken-relay/network/kusama.json"
)

var sc *substrate.SarpcClient

func init() {
	var err error
	stop := make(chan int)
	// sc, err = substrate.NewSarpcClient(substrate.ChainTypePolkadot, "wss://kusama-rpc.polkadot.io", polkaTypesFile, substrate.AddressTypeMultiAddress, AliceKey, tlog, stop)
	sc, err = substrate.NewSarpcClient(substrate.ChainTypeStafi, "wss://mainnet-rpc.stafi.io", stafiTypesFile, substrate.AddressTypeAccountId, AliceKey, tlog, stop)
	if err != nil {
		panic(err)
	}
}

func TestConnection_TransferVerify(t *testing.T) {
	conn := Connection{sc: sc, log: tlog, symbol: core.RKSM, stop: make(chan int)}
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

func TestConnection_TransferVerify1(t *testing.T) {
	conn := Connection{sc: sc, log: tlog, symbol: core.RKSM, stop: make(chan int)}
	pool, _ := hexutil.Decode("0x0777f13da6fead588d8662fd63336b0008c4fbe4749e18779c1a4bd89ea50141")
	bonder, _ := hexutil.Decode("0xf698bb9fe10437d9d8b3b27d4bb4bdc683228330cf77f992c2c4f0181c1c5d72")
	pubkey, _ := hexutil.Decode("0xf698bb9fe10437d9d8b3b27d4bb4bdc683228330cf77f992c2c4f0181c1c5d72")
	blockHash, _ := hexutil.Decode("0x6157da60a188b3f31d250afe5acb2da786417fec00973f1c7f863504fbca4642")
	txHash, _ := hexutil.Decode("0xde0f3f2159f82cdfcc8ae161da48ebee347ba2d18e4f4564a636bedb495b15bd")
	bond := &submodel.BondRecord{
		Pool:      types.NewBytes(pool),
		Bonder:    types.NewAccountID(bonder),
		Symbol:    core.RKSM,
		Pubkey:    types.NewBytes(pubkey),
		Blockhash: types.NewBytes(blockHash),
		Txhash:    types.NewBytes(txHash),
		Amount:    types.NewU128(*big.NewInt(1000000000000)),
	}

	result, err := conn.TransferVerify(bond)
	assert.NoError(t, err)
	t.Log(result)
	assert.Equal(t, result, submodel.Pass)
}

func TestProposalVoters(t *testing.T) {
	conn := Connection{sc: sc, log: tlog, symbol: core.RFIS, stop: make(chan int)}
	accounts, err := conn.GetNewChainEraProposalVoters(core.RDOT, 889)
	assert.NoError(t, err)
	t.Log(accounts)
	voter, err := conn.GetSelectedVoters(core.RDOT, 889)
	assert.NoError(t, err)
	t.Log(voter)

}

func TestSortAddr(t *testing.T) {
	ss58AddressList := []string{
		"36BwqjgT8MkuwMwJBNLBohLci7mTqKLoC7YKMLeP7kzRASwC",
		"351cAhFpbcjBboEsAjeFzDv4nLzPgSggVyT2Pf56Z9nmNb6F",
		"35VeijRPg5zVt4kKZsMUkornNMZn7DrJrmJawKLxgkpRFprs",
		"32G1VrGGoJStYgVFv6n8qXaoAhRVHbmM4UHPvkQSsBkSQtom",
		"33v9bvusE56vhPuc7PT9GbAfhwkVS2e2N1ThhvRVxfw2U6XE",
		"33Cae8pSE2DLeVwL2Ugb3PtdXGuwMsoHC8mHDmtntNTh4cit",
		"34un2Kxb5UPq4DzBEu4N4FG7evfaSDUzcVgtXRdQWjeDsGGt",
		"32FjxDWcLGehAW4QZueWHVUHDfZe6L1ZvJFewMgCkKMZuBKC",
	}

	addressBtsList := make([][]byte, 0)
	for _, s58Addr := range ss58AddressList {
		bts, err := ss58.DecodeToPub(s58Addr)
		if err != nil {
			t.Fatal(err)
		}
		addressBtsList = append(addressBtsList, bts)
	}

	sort.SliceStable(addressBtsList, func(i, j int) bool {
		return bytes.Compare(addressBtsList[i][:], addressBtsList[j][:]) < 0
	})

	for _, addrBts := range addressBtsList {
		addrHexStr := hex.EncodeToString(addrBts)
		addrSs58Str, err := ss58.EncodeByPubHex(addrHexStr, ss58.StafiPrefix)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(addrHexStr, addrSs58Str)
	}

	pbk,err:=sr25519.NewPublicKey([32]byte{})
	if err!=nil{
		t.Fatal(err)
	}
	t.Log(pbk)

}
