package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"os"
)

//multisig
// cosmos1em384d8ek3y8nlugapz7p5k5skg58j66je3las

// cosmosvaloper1tgdut5ldng5gr9dh4ypxvdvznptgn6q8txjc46 //self
// cosmosvaloper14xj7e0aqraavut998j08eg6x7nw3r4ts0qhfpw
// cosmosvaloper1u22lut8qgqg8znxam72pwgqp8c09rnvmummr4w //target
// cosmosvaloper1u7m4j26ukn293latjtnv2pjrtadzu9s805g6pg

var client *rpc.Client
var keyDir = "/Users/tpkeeper/.gaia"
var addrMultiSig1, _ = types.AccAddressFromBech32("cosmos1em384d8ek3y8nlugapz7p5k5skg58j66je3las")
var adrValidatorEverStake, _ = types.ValAddressFromBech32("cosmosvaloper1u7m4j26ukn293latjtnv2pjrtadzu9s805g6pg")

func init() {
	fmt.Printf("Will open cosmos wallet from <%s>. \nPlease ", keyDir)
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, keyDir, os.Stdin)
	if err != nil {
		panic(err)
	}

	client, _ = rpc.NewClient(key, "chain-AALfXF", "key0key1key2", "0.00001stake", "stake", "http://127.0.0.1:26657")
}

func main() {
	err := GenMultiSigRawDelegateTx()
	if err != nil {
		fmt.Println(err)
	}
}

func GenMultiSigRawDelegateTx() error {
	err := client.SetFromName("key0key1key2")
	if err != nil {
		return err
	}
	rawTx, err := client.GenMultiSigRawDelegateTx(addrMultiSig1, []types.ValAddress{adrValidatorEverStake},
		types.NewCoin(client.GetDenom(), types.NewInt(10000000)))

	if err != nil {
		return err
	}

	signature1, err := client.SignMultiSigRawTx(rawTx, "key0")
	if err != nil {
		return err
	}
	signature2, err := client.SignMultiSigRawTx(rawTx, "key1")
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
