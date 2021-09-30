package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/ChainSafe/log15"

	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

var configFlag = &cli.StringFlag{
	Name:  "config",
	Usage: "json configuration file",
	Value: "./config_init.json",
}
var cliFlags = []cli.Flag{
	configFlag,
}

// init initializes CLI
func init() {
	app.Action = run
	app.Copyright = "Copyright 2020 Stafi Protocol Authors"
	app.Name = "soltool"
	app.Usage = "solTool"
	app.Authors = []*cli.Author{{Name: "Stafi Protocol 2020"}}
	app.Version = "0.0.1"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		&initCommand,
	}

}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	return nil
}

var initCommand = cli.Command{
	Name:        "initAccount",
	Usage:       "manage account",
	Description: "The init command is used to init account",
	Action:      initAction,
	Flags:       cliFlags,
}

type PoolAccounts struct {
	KeystorePath                string            `json:"keystorePath"`
	FeeAccount                  string            `json:"feeAccount"`
	StakeBaseAccountToValidator map[string]string `json:"stakeBaseAccountToValidator"`
	MultisigTxBaseAccount       string            `json:"multisigTxBaseAccount"`
	MultisigInfoPubkey          string            `json:"multisigInfoPubkey"`
	MultisigProgramId           string            `json:"multisigProgramId"`
	Endpoint                    string            `json:"endpoint"`
	OtherFeeAccount             []string          `json:"otherFeeAccount"`
	Threshold                   int64             `json:"threshold"`
}

func loadConfig(file string, config *PoolAccounts) (err error) {
	ext := filepath.Ext(file)
	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	log.Debug("Loading configuration", "path", filepath.Clean(fp))

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
	return json.NewDecoder(f).Decode(&config)
}
