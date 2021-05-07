package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"os"
)

var client *rpc.Client
var keyDir = "/home/stafi/ratom/keys/keys/cosmos"
var addrMultiSig1, _ = types.AccAddressFromBech32("cosmos12yprrdprzat35zhqxe2fcnn3u26gwlt6xcq0pj")
var adrValidatorEverStake, _ = types.ValAddressFromBech32("cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3")

func init() {
	fmt.Printf("Will open cosmos wallet from <%s>. \nPlease ", keyDir)
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, keyDir, os.Stdin)
	if err != nil {
		panic(err)
	}

	client, _ = rpc.NewClient(key, "stargate-final", "recipient", "0.04umuon", "umuon", "https://testcosmosrpc.wetez.io:443")
}

func main() {
	err := GenMultiSigRawDelegateTx()
	if err != nil {
		fmt.Println(err)
	}
}
func GenMultiSigRawDelegateTx() error {
	err := client.SetFromName("multikey1")
	if err != nil {
		return err
	}
	rawTx, err := client.GenMultiSigRawDelegateTx(addrMultiSig1, []types.ValAddress{adrValidatorEverStake},
		types.NewCoin(client.GetDenom(), types.NewInt(5000000)))

	if err != nil {
		return err
	}

	signature1, err := client.SignMultiSigRawTx(rawTx, "multisubkey1")
	signature2, err := client.SignMultiSigRawTx(rawTx, "multisubkey2")

	_, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2})
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
