package main

import (
	"os"
	"sync"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	xAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
)

const DefaultHomeDir = "./keys/cosmos"

func main() {
	encodingConfig := rpc.MakeEncodingConfig()

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(xAuthTypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(DefaultHomeDir)

	UseSdkConfigContext("swth")
	rootCmd := &cobra.Command{
		Use:   "gencos",
		Short: "tool to manage cosmos keys",
		Long:  "Notice: The keyring supports os|file|test backends, but relay now only support the file backend",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd, "", nil)
		},
	}

	rootCmd.AddCommand(keys.Commands(DefaultHomeDir))
	err := svrcmd.Execute(rootCmd, DefaultHomeDir)
	if err != nil {
		panic(err)
	}
}

var sdkContextMutex sync.Mutex

// UseSDKContext uses a custom Bech32 account prefix and returns a restore func
// CONTRACT: When using this function, caller must ensure that lock contention
// doesn't cause program to hang. This function is only for use in codec calls
func UseSdkConfigContext(accountPrefix string) func() {
	// Ensure we're the only one using the global context,
	// lock context to begin function
	sdkContextMutex.Lock()

	// Mutate the sdkConf

	if accountPrefix == "iaa" {
		setIrisPrefix()
	} else {
		setCommonPrefixes(accountPrefix)
	}

	// Return the unlock function, caller must lock and ensure that lock is released
	// before any other function needs to use c.UseSDKContext
	return sdkContextMutex.Unlock
}

func setCommonPrefixes(accountAddressPrefix string) {
	// Set prefixes
	accountPubKeyPrefix := accountAddressPrefix + "pub"
	validatorAddressPrefix := accountAddressPrefix + "valoper"
	validatorPubKeyPrefix := accountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := accountAddressPrefix + "valcons"
	consNodePubKeyPrefix := accountAddressPrefix + "valconspub"

	// Set config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(accountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
}

func setIrisPrefix() {

	// Bech32ChainPrefix defines the prefix of this chain
	Bech32ChainPrefix := "i"

	// PrefixAcc is the prefix for account
	PrefixAcc := "a"

	// PrefixValidator is the prefix for validator keys
	PrefixValidator := "v"

	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus := "c"

	// PrefixPublic is the prefix for public
	PrefixPublic := "p"

	// PrefixAddress is the prefix for address
	PrefixAddress := "a"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr := Bech32ChainPrefix + PrefixAcc + PrefixAddress
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub := Bech32ChainPrefix + PrefixAcc + PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr := Bech32ChainPrefix + PrefixValidator + PrefixAddress
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub := Bech32ChainPrefix + PrefixValidator + PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr := Bech32ChainPrefix + PrefixConsensus + PrefixAddress
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub := Bech32ChainPrefix + PrefixConsensus + PrefixPublic

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}
