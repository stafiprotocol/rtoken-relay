// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/ChainSafe/log15"
	"github.com/urfave/cli/v2"
)

const DefaultConfigPath = "./config.json"
const DefaultKeystorePath = "./keys"

type Config struct {
	MainConf   MainConfig    `json:"main"`
	OtherConfs []OtherConfig `json:"others"`
}

type MainConfig struct {
	Name           string            `json:"name"`
	Endpoint       string            `json:"endpoint"`
	From           string            `json:"from"`
	KeystorePath   string            `json:"keystorePath"`
	BlockstorePath string            `json:"blockstorePath"`
	TypesPath      string            `json:"typesPath"`
	Opts           map[string]string `json:"opts"`
}

// RawChainConfig is parsed directly from the config file and should be using to construct the core.ChainConfig
type OtherConfig struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	TypesPath    string   `json:"typesPath"`
	Endpoint     string   `json:"endpoint"`
	Symbol       string   `json:"symbol"`
	Accounts     []string `json:"accounts"`
	KeystorePath string   `json:"keystorePath"`
}

func GetConfig(ctx *cli.Context) (*Config, error) {
	var fig Config
	path := DefaultConfigPath
	if file := ctx.String(ConfigFileFlag.Name); file != "" {
		path = file
	}
	err := loadConfig(path, &fig)
	if err != nil {
		log.Warn("err loading json file", "err", err.Error())
		return &fig, err
	}
	log.Debug("Loaded config", "path", path)

	return &fig, nil
}

func loadConfig(file string, config *Config) error {
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

	if ext == ".json" {
		if err = json.NewDecoder(f).Decode(&config); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unrecognized extention: %s", ext)
	}

	return nil
}
