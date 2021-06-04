// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/spf13/cobra"
	"github.com/stafiprotocol/rtoken-relay/shared/solana/vault"
)

// vaultExportCommand represents the export command
var vaultExportCommand = &cobra.Command{
	Use:   "export",
	Short: "Export private keys (and corresponding public keys) inside a Solana vault.",
	Run: func(cmd *cobra.Command, args []string) {
		vault := vault.MustGetWallet()

		vault.PrintPrivateKeys()
	},
}

func init() {
	VaultCmd.AddCommand(vaultExportCommand)
}
