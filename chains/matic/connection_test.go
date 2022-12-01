// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package matic

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/urfave/cli/v2"
)

var (
	testLogger = newTestLogger("test")
)

func TestConnection(t *testing.T) {
	ctx, err := createCliContext("", []string{config.ConfigFileFlag.Name}, []interface{}{"/Users/fwj/Go/stafi/rtoken-relay/matic.json"})
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.GetConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}

	//t.Log(cfg)
	chain := cfg.Chains[0]
	chainConfig := &core.ChainConfig{
		Name:            chain.Name,
		Endpoint:        chain.Endpoint,
		KeystorePath:    chain.KeystorePath,
		Symbol:          core.RSymbol(chain.Rsymbol),
		Accounts:        chain.Accounts,
		LatestBlockFlag: chain.LatestBlockFlag,
		Insecure:        false,
		Opts:            chain.Opts,
	}

	_ = os.Setenv(keystore.EnvPassword, "123456")
	stop := make(chan int)
	conn, err := NewConnection(chainConfig, testLogger, stop)
	if err != nil {
		t.Fatal(err)
	}

	blk, err := conn.LatestBlock()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(blk)

	valId := big.NewInt(9)
	valFlag, err := conn.stakeManager.IsValidator(conn.conn.CallOpts(), valId)
	if err != nil {
		t.Fatal(err)
	}

	if !valFlag {
		t.Fatal("valFlag is false")
	}

	shareAddr, err := conn.GetValidator(valId)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("shareAddr", shareAddr)

	poolAddr := common.HexToAddress("0x1Cb8b55cB11152E74D34Be1961E4FFe169F5B99A")
	staked, err := conn.TotalStaked(shareAddr, poolAddr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("staked", staked)

	//bond := big.NewInt(0).Mul(&config.AmountBase, big.NewInt(1))
	//unbond := big.NewInt(0)
	//
	//method, tx, err := conn.BondOrUnbondCall(shareAddr, bond, unbond)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Log("method", method)
	//t.Log("tx", tx)
	//
	//msg, err := conn.MessageToSign(tx)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//key := conn.keys[0]
	//signature, err := crypto.Sign(msg[:], key.PrivateKey())
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log("signature", hexutil.Encode(signature))
	//param := submodel.SubmitSignatureParams{
	//	Symbol:     core.RMATIC,
	//	Era:        types.NewU32(0),
	//	Pool:       []byte("test"),
	//	TxType:     method,
	//	Signature:  signature,
	//}
	//
	//txHash, err := param.EncodeToHash()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Log("txHash", txHash.Hex())
	//safeTxGas := BuyVoucherSafeTxGas
	//
	//signatures := [][]byte{signature}
	//vs, rs, ss := utils.DecomposeSignature(signatures)
	//err = conn.AsMulti(tx.To, tx.Value, tx.CallData, config.Call, safeTxGas, txHash, vs, rs, ss)
	//if err != nil {
	//	t.Fatal(err)
	//}
}

func createCliContext(description string, flags []string, values []interface{}) (*cli.Context, error) {
	set := flag.NewFlagSet(description, 0)
	for i := range values {
		switch v := values[i].(type) {
		case bool:
			set.Bool(flags[i], v, "")
		case string:
			set.String(flags[i], v, "")
		case uint:
			set.Uint(flags[i], v, "")
		default:
			return nil, fmt.Errorf("unexpected cli value type: %T", values[i])
		}
	}
	context := cli.NewContext(nil, set, nil)
	return context, nil
}

func newTestLogger(name string) core.Logger {
	tLog := core.NewLog("chain", name)
	return tLog
}
