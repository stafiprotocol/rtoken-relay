package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
	"github.com/urfave/cli/v2"
)

var transferCommand = cli.Command{
	Name:   "transfer",
	Usage:  "tranfer eth",
	Action: handleTransfer,
	Flags: []cli.Flag{
		config.ConfigFileFlag,
	},
}

type TransferConfig struct {
	EthEndpoint  string `json:"ethEndpoint"`
	KeystorePath string `json:"keystorePath"`
	FromAccount  string `json:"fromAccount"`
	ToAccount    string `json:"toAccount"`
	Amount       string `json:"amount"`
}

func loadTransferConfig(file string, config *TransferConfig) (err error) {
	ext := filepath.Ext(file)
	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	core.NewLog().Debug("Loading configuration", "path", filepath.Clean(fp))

	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()

	if ext != ".json" {
		return fmt.Errorf("unrecognized extention: %s", ext)
	}
	return json.NewDecoder(f).Decode(config)
}

func handleTransfer(ctx *cli.Context) error {
	path := "./config_transfer.json"
	if file := ctx.String(config.ConfigFileFlag.Name); file != "" {
		path = file
	}
	cfg := TransferConfig{}
	err := loadTransferConfig(path, &cfg)
	if err != nil {
		return err
	}
	lvl := ctx.String(config.VerbosityFlag.Name)

	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	if !common.IsHexAddress(cfg.FromAccount) {
		return fmt.Errorf("fromAccount is not hex address")
	}
	if !common.IsHexAddress(cfg.ToAccount) {
		return fmt.Errorf("toAccount is not hex address")
	}

	kpI, err := keystore.KeypairFromAddress(cfg.FromAccount, keystore.EthChain, cfg.KeystorePath, false)
	if err != nil {
		return err
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	ethClient := ethereum.NewClient(cfg.EthEndpoint, kp, log, big.NewInt(0), big.NewInt(0))
	if err := ethClient.Connect(); err != nil {
		return err
	}

	// send tx
	err = ethClient.LockAndUpdateOpts(big.NewInt(0), big.NewInt(0))
	if err != nil {
		return fmt.Errorf("LockAndUpdateOpts error %s", err)
	}
	ethClient.UnlockOpts()

	amount, err := decimal.NewFromString(cfg.Amount)
	if err != nil {
		return fmt.Errorf("parase amount failed, err: %s", err.Error())
	}
	nonce := ethClient.Opts().Nonce.Uint64()
	to := common.HexToAddress(cfg.ToAccount)
	gasLimit := uint64(22000)
	gasPrice := ethClient.Opts().GasPrice
	tx := types.NewTransaction(nonce, to, amount.BigInt(), gasLimit, gasPrice, nil)
	chainId, err := ethClient.Client().ChainID(ctx.Context)
	if err != nil {
		return fmt.Errorf("get chainId failed, err: %s", err.Error())
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), kp.PrivateKey())
	if err != nil {
		return fmt.Errorf("sign tx failed, err: %s", err.Error())
	}

	err = ethClient.Client().SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("send tx failed, err: %s", err.Error())
	}

	err = waitTransferTxOk(signedTx.Hash(), ethClient)
	if err != nil {
		return fmt.Errorf("waitTx error %s", err)
	}

	return nil
}

func waitTransferTxOk(txHash common.Hash, polygonConn *ethereum.Client) error {
	retry := 0
	for {
		if retry > 300 {
			return fmt.Errorf("waitTxOk tx reach retry limit")
		}
		_, pending, err := polygonConn.TransactionByHash(context.Background(), txHash)
		if err == nil && !pending {
			break
		} else {
			if err != nil {
				logrus.Warn("tx query ", "hash ", txHash, "err ", err.Error())
			} else {
				logrus.Warn("tx pending ", "hash ", txHash)
			}
			time.Sleep(6 * time.Second)
			retry++
			continue
		}

	}
	logrus.Info("tx send ok ", "tx ", txHash.String())
	return nil
}
