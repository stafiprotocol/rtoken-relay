package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"os"
)

var client *rpc.Client
var keyDir = "/Users/tpkeeper/.gaia"
var addrMultiSig1, _ = types.AccAddressFromBech32("cosmos1ak3nrcmm7e4j8y7ycfc78pxl4g4lehf43vw6wu")
var adrValidatorEverStake, _ = types.ValAddressFromBech32("cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3")

func init() {
	fmt.Printf("Will open cosmos wallet from <%s>. \nPlease ", keyDir)
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, keyDir, os.Stdin)
	if err != nil {
		panic(err)
	}

	client, _ = rpc.NewClient(key, "stargate-final", "recipient", "0.04umuon", "umuon", "http://127.0.0.1:26657")
}

func main() {
	err := GenMultiSigRawDelegateTx()
	if err != nil {
		fmt.Println(err)
	}
}
func GenMultiSigRawDelegateTx() error {
	err := client.SetFromName("multiSign1")
	if err != nil {
		return err
	}
	rawTx, err := client.GenMultiSigRawDelegateTx(addrMultiSig1, []types.ValAddress{adrValidatorEverStake},
		types.NewCoin(client.GetDenom(), types.NewInt(1)))

	if err != nil {
		return err
	}

	signature1, err := client.SignMultiSigRawTx(rawTx, "key1")
	signature2, err := client.SignMultiSigRawTx(rawTx, "key2")
	signature3, err := client.SignMultiSigRawTx(rawTx, "key3")

	_, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2, signature3})
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
