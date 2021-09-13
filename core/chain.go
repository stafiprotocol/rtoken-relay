// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package core

type Chain interface {
	Start() error // Start chain
	SetRouter(*Router)
	Rsymbol() RSymbol
	Name() string
	Stop()
	InitBondedPools(symbols []RSymbol) error
}

type ChainConfig struct {
	Name            string                 // Human-readable chain name
	Symbol          RSymbol                // symbol
	Endpoint        string                 // url for rpc endpoint
	Accounts        []string               // addresses of key to use
	KeystorePath    string                 // Location of key files
	Insecure        bool                   // Indicated whether the test keyring should be used
	LatestBlockFlag bool                   // If true, overrides blockstore or latest block in config and starts from current block
	Opts            map[string]interface{} // Per chain options
}
