// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

const DefaultConfigPath = "./config.json"
const DefaultKeystorePath = "./keys"

type Config struct {
	Chains []RawChainConfig `json:"chains"`
}

// RawChainConfig is parsed directly from the config file and should be using to construct the core.ChainConfig
type RawChainConfig struct {
	Name            string                 `json:"name"`
	Type            string                 `json:"type"`
	Rsymbol         string                 `json:"rsymbol"`
	Endpoint        string                 `json:"endpoint"` // url for rpc endpoint
	KeystorePath    string                 `json:"keystorePath"`
	Accounts        []string               `json:"accounts"` // address of key to use
	LatestBlockFlag bool                   `json:"latestBlockFlag"`
	Opts            map[string]interface{} `json:"opts"`
}

func NewConfig() *Config {
	return &Config{
		Chains: []RawChainConfig{},
	}
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
