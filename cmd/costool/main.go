package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
)

//multisig
// cosmos1em384d8ek3y8nlugapz7p5k5skg58j66je3las

// cosmosvaloper1tgdut5ldng5gr9dh4ypxvdvznptgn6q8txjc46 //self
// cosmosvaloper14xj7e0aqraavut998j08eg6x7nw3r4ts0qhfpw
// cosmosvaloper1u22lut8qgqg8znxam72pwgqp8c09rnvmummr4w //target
// cosmosvaloper1u7m4j26ukn293latjtnv2pjrtadzu9s805g6pg

var client *rpc.Client
var keyDir = "/Users/tpkeeper/gowork/stafi/rtoken-relay/keys/cosmos"
var addrMultiSig1 types.AccAddress
var adrValidatorEverStake types.ValAddress

func init() {
	UseSdkConfigContext("swth")
	var err error
	addrMultiSig1, err = types.AccAddressFromBech32("swth10j3yjvgzm7r22us3tqz49gkgtkj0rt3pg6w8z7")
	if err != nil {
		panic(err)
	}
	adrValidatorEverStake, err = types.ValAddressFromBech32("swthvaloper1fdqkq5gc5x8h6a0j9hamc30stlvea6zlsh4y7s")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Will open cosmos wallet from <%s>. \nPlease ", keyDir)
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, keyDir, os.Stdin)
	if err != nil {
		panic(err)
	}

	client, err = rpc.NewClient(key, "carbon-1", "multisign", "1000swth", "swth", "https://tm-api.carbon.network:443")
	if err != nil {
		panic(err)
	}
}

func main() {
	err := GenMultiSigRawDelegateTx()
	if err != nil {
		fmt.Println(err)
	}
}

func GenMultiSigRawDelegateTx() error {
	err := client.SetFromName("multisign")
	if err != nil {
		return err
	}
	rawTx, err := client.GenMultiSigRawDelegateTx(addrMultiSig1, []types.ValAddress{adrValidatorEverStake},
		types.NewCoin(client.GetDenom(), types.NewInt(100000000)))

	if err != nil {
		return err
	}

	signature1, err := client.SignMultiSigRawTx(rawTx, "key1")
	if err != nil {
		return err
	}
	signature2, err := client.SignMultiSigRawTx(rawTx, "key2")
	if err != nil {
		return err
	}
	_, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2}, 2)
	if err != nil {
		return err
	}

	hash, err := client.BroadcastTx(tx)
	if err != nil {
		return err
	}
	fmt.Println("hash ", hash)
	return nil
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
	config := types.GetConfig()
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

	config := types.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}
