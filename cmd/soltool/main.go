package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stafiprotocol/rtoken-relay/core"
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
		&migrateCommand,
		&withdrawCommand,
		&stakeCommand,
	}

}

func main() {
	if err := app.Run(os.Args); err != nil {
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

var migrateCommand = cli.Command{
	Name:   "migrateStakeAccount",
	Usage:  "migrate stake account",
	Action: migrateStakeAccount,
	Flags:  cliFlags,
}

var withdrawCommand = cli.Command{
	Name:   "withdrawStakeAccount",
	Usage:  "withdraw stake account",
	Action: withdrawStakeAccount,
	Flags:  cliFlags,
}

var stakeCommand = cli.Command{
	Name:   "stake",
	Usage:  "stake to validator",
	Action: stake,
	Flags:  cliFlags,
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

type PoolAccountsForMigrate struct {
	KeystorePath          string   `json:"keystorePath"`
	FeeAccount            string   `json:"feeAccount"`
	MultisigTxBaseAccount string   `json:"multisigTxBaseAccount"`
	MultisigInfoPubkey    string   `json:"multisigInfoPubkey"`
	MultisigProgramId     string   `json:"multisigProgramId"`
	Endpoint              string   `json:"endpoint"`
	OtherFeeAccount       []string `json:"otherFeeAccount"`

	RSolProgramID string `json:"rSolProgramID"`
	StakeManager  string `json:"stakeManager"`
	StakePool     string `json:"stakePool"`
	StakeAccount  string `json:"stakeAccount"`
	Validator     string `json:"validator"`
}

func loadConfig(file string, config *PoolAccounts) (err error) {
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
	return json.NewDecoder(f).Decode(&config)
}

func loadConfigForMigrate(file string, config *PoolAccountsForMigrate) (err error) {
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
	return json.NewDecoder(f).Decode(&config)
}
