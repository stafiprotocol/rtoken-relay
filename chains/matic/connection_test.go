// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package matic

import (
	"flag"
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/urfave/cli/v2"
	"math/big"
	"os"
	"testing"
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

	os.Setenv(keystore.EnvPassword, "123456")
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

	valId := big.NewInt(111)
	valFlag, err := conn.stakeManager.IsValidator(conn.conn.CallOpts(), valId)
	if err != nil {
		t.Fatal(err)
	}

	if !valFlag {
		t.Log("valFlag is false")
	}

	//shareAddr, err := conn.GetValidator(valId)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Log(shareAddr)
}

//import (
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/stafiprotocol/rtoken-relay/core"
//)
//
//var (
//
//	keystorePath = "/Users/fwj/Go/stafi/rtoken-relay/keys/ethereum/"
//
//
//	config = &core.ChainConfig{
//		Name: "ethereum",
//		Symbol: core.RMATIC,
//		Endpoint:goerliEndPoint,
//		Accounts: []string{
//			"0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6",
//			"0xBd39f5936969828eD9315220659cD11129071814",
//		},
//		KeystorePath: keystorePath,
//		Insecure: false,
//		LatestBlockFlag: true,
//		Opts: map[string]interface{}{
//			"StakeManagerContract": "0x4864d89DCE4e24b2eDF64735E014a7E4154bfA7A",
//			"MultisigContract": "multisigContract",
//			""
//
//		},
//
//
//	}
//)

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

func newTestLogger(name string) log15.Logger {
	tLog := log15.New("chain", name)
	tLog.SetHandler(log15.LvlFilterHandler(log15.LvlError, tLog.GetHandler()))
	return tLog
}
