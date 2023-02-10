package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/rtoken-relay/chains/bnb"
	"github.com/stafiprotocol/rtoken-relay/chains/matic"
	"github.com/stafiprotocol/rtoken-relay/chains/solana"
	"github.com/stafiprotocol/rtoken-relay/chains/substrate"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

var mainFlags = []cli.Flag{
	config.ConfigFileFlag,
	config.VerbosityFlag,
}

var generateFlags = []cli.Flag{
	config.KeystorePathFlag,
	config.NetworkFlag,
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
			Action:      handleGenerateSubCmd,
			Name:        "gensub",
			Usage:       "generate subsrate keystore",
			Flags:       generateFlags,
			Description: "The generate subcommand is used to generate the substrate keystore.",
		},
		{
			Action:      handleGenerateEthCmd,
			Name:        "geneth",
			Usage:       "generate ethereum keystore",
			Flags:       generateFlags,
			Description: "The generate subcommand is used to generate the ethereum keystore.",
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
	app.Version = "0.1.0"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		&accountCommand,
		&delegateBnbCommand,
		&claimUndelegateCommand,
	}

	app.Flags = append(app.Flags, mainFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startLogger(ctx *cli.Context) error {
	lvl := ctx.String(config.VerbosityFlag.Name)

	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	err = core.InitLogFile("./log_file")
	if err != nil {
		return err
	}

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
		//check symbol
		switch chain.Rsymbol {
		case string(core.RBNB), string(core.RDOT), string(core.RFIS), string(core.RKSM), string(core.RMATIC), string(core.RSOL):
		default:
			return fmt.Errorf("rsymbol not match: %s", chain.Rsymbol)
		}

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
		logger := core.NewLog("chain", chainConfig.Name)

		switch chain.Type {
		case "substrate":
			newChain, err = substrate.InitializeChain(chainConfig, logger, sysErr)
			if err != nil {
				return err
			}
		case "solana":
			newChain, err = solana.InitializeChain(chainConfig, logger, sysErr)
			if err != nil {
				return err
			}

		case "ethereum":
			newChain, err = matic.InitializeChain(chainConfig, logger, sysErr)
			if err != nil {
				return err
			}
		case "bnb":
			newChain, err = bnb.InitializeChain(chainConfig, logger, sysErr)
			if err != nil {
				return err
			}
		default:
			return errors.New("unrecognized Chain Type")
		}

		c.AddChain(newChain)
	}

	c.Start()

	return nil
}
