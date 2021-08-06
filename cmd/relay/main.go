package main

import (
	"errors"
	"os"
	"strconv"

	log "github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/chains/cosmos"
	"github.com/stafiprotocol/rtoken-relay/chains/matic"
	"github.com/stafiprotocol/rtoken-relay/chains/substrate"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

var cliFlags = []cli.Flag{
	config.ConfigFileFlag,
	config.VerbosityFlag,
	config.KeystorePathFlag,
}

var generateFlags = []cli.Flag{
	config.KeystorePathFlag,
}

var bncGenerateFlags = []cli.Flag{
	config.KeystorePathFlag,
	config.BncNetwork,
}

var accountCommand = cli.Command{
	Name:  "accounts",
	Usage: "manage reth keystore",
	Description: "The accounts command is used to manage the relay keystore.\n" +
		"\tMake sure the keystore dir is exist before generating\n" +
		"\tTo generate a substrate keystore: relay accounts gensub\n" +
		"\tTo generate a ethereum keystore: relay accounts geneth\n" +
		"\tTo generate a bc chain keystore: relay accounts genbc\n" +
		"\tTo list keys: chainbridge accounts list",
	Subcommands: []*cli.Command{
		{
			Action: wrapHandler(handleGenerateSubCmd),
			Name:   "gensub",
			Usage:  "generate subsrate keystore",
			Flags:  generateFlags,
			Description: "The generate subcommand is used to generate the substrate keystore.",
		},
		{
			Action: wrapHandler(handleGenerateEthCmd),
			Name:   "geneth",
			Usage:  "generate ethereum keystore",
			Flags:  generateFlags,
			Description: "The generate subcommand is used to generate the ethereum keystore.",
		},
		{
			Action: wrapHandler(handleGenerateBcCmd),
			Name:   "genbc",
			Usage:  "generate bc chain keystore",
			Flags:  bncGenerateFlags,
			Description: "The generate subcommand is used to generate the bc chain keystore.",
		},
	},
}

// init initializes CLI
func init() {
	app.Action = run
	app.Copyright = "Copyright 2020 Stafi Protocol Authors"
	app.Name = "reley"
	app.Usage = "relay"
	app.Authors = []*cli.Author{{Name: "Stafi Protocol 2020"}}
	app.Version = "0.0.1"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		&accountCommand,
	}

	app.Flags = append(app.Flags, cliFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startLogger(ctx *cli.Context) error {
	logger := log.Root()
	var lvl log.Lvl

	if lvlToInt, err := strconv.Atoi(ctx.String(config.VerbosityFlag.Name)); err == nil {
		lvl = log.Lvl(lvlToInt)
	} else if lvl, err = log.LvlFromString(ctx.String(config.VerbosityFlag.Name)); err != nil {
		return err
	}

	logger.SetHandler(log.MultiHandler(
		log.LvlFilterHandler(
			lvl,
			log.StreamHandler(os.Stdout, log.LogfmtFormat())),
		log.Must.FileHandler("relay_log.json", log.JsonFormat()),
		log.LvlFilterHandler(
			log.LvlError,
			log.Must.FileHandler("relay_log_errors.json", log.JsonFormat()))))

	return nil
}

func run(ctx *cli.Context) error {
	err := startLogger(ctx)
	if err != nil {
		return err
	}

	cfg, err := config.GetConfig(ctx)
	if err != nil {
		return err
	}

	// Used to signal core shutdown due to fatal error
	sysErr := make(chan error)
	c := core.NewCore(sysErr)

	for _, chain := range cfg.Chains {
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
		var newChain core.Chain
		logger := log.Root().New("chain", chainConfig.Name)

		if chain.Type == "substrate" {
			newChain, err = substrate.InitializeChain(chainConfig, logger, sysErr)
			if err != nil {
				return err
			}
		} else if chain.Type == "cosmos" {
			newChain, err = cosmos.InitializeChain(chainConfig, logger, sysErr)
			if err != nil {
				return err
			}
		} else if chain.Type == "ethereum" {
			newChain, err = matic.InitializeChain(chainConfig, logger, sysErr)
			if err != nil {
				return err
			}
		} else {
			return errors.New("unrecognized Chain Type")
		}
		c.AddChain(newChain)
	}

	c.Start()

	return nil
}
