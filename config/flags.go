// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	log "github.com/ChainSafe/log15"
	"github.com/urfave/cli/v2"
)

var (
	ConfigFileFlag = &cli.StringFlag{
		Name:  "config",
		Usage: "json configuration file",
		Value: DefaultConfigPath,
	}

	VerbosityFlag = &cli.StringFlag{
		Name:  "verbosity",
		Usage: "supports levels crit (silent) to trce (trace)",
		Value: log.LvlInfo.String(),
	}

	KeystorePathFlag = &cli.StringFlag{
		Name:  "keystore",
		Usage: "path to keystore directory",
		Value: DefaultKeystorePath,
	}

	BncNetwork = &cli.StringFlag{
		Name:  "bncnetwork",
		Usage: "specify network for bc chain, set test for TestNetwork, others will be ProdNetwork",
		Value: "",
	}

	NetworkFlag = &cli.StringFlag{
		Name:  "network",
		Usage: "specify network for subkey",
		Value: "stafi",
	}
)
