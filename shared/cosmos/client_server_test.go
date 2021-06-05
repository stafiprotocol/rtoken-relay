package cosmos_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
	"github.com/stretchr/testify/assert"
)

var client *rpc.Client

var addrMultiKey1, _ = types.AccAddressFromBech32("cosmos12yprrdprzat35zhqxe2fcnn3u26gwlt6xcq0pj")

var addrAccount1, _ = types.AccAddressFromBech32("cosmos12wrv225462drlz4dk3yg9hc8vavwjkmckshz7c")

var adrValidatorTestnetTecos, _ = types.ValAddressFromBech32("cosmosvaloper1p7e37nztj62mmra8xhgqde7sql3llhhu6hvcx8")
var adrValidatorEverStake, _ = types.ValAddressFromBech32("cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3")

func init() {
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, "/home/stafi/ratom/keys/keys/cosmos", strings.NewReader("12345678\n"))
	if err != nil {
		panic(err)
	}

	client, err = rpc.NewClient(key, "stargate-final", "recipient", "0.04umuon", "umuon", "https://testcosmosrpc.wetez.io:443")
	if err != nil {
		panic(err)
	}
}

func TestClient_SendTo(t *testing.T) {
	err := client.SingleTransferTo(addrAccount1, types.NewCoins(types.NewInt64Coin(client.GetDenom(), 100000000)))
	assert.NoError(t, err)
}

func TestClient_GenMultiSigRawDelegateTx(t *testing.T) {
	err := client.SetFromName("multikey1")
	assert.NoError(t, err)
	rawTx, err := client.GenMultiSigRawDelegateTx(addrMultiKey1, []types.ValAddress{adrValidatorTestnetTecos, adrValidatorEverStake},
		types.NewCoin(client.GetDenom(), types.NewInt(10000000)))

	assert.NoError(t, err)

	signature1, err := client.SignMultiSigRawTx(rawTx, "multisubkey1")
	assert.NoError(t, err)
	signature2, err := client.SignMultiSigRawTx(rawTx, "multisubkey2")
	assert.NoError(t, err)

	_, tx, err := client.AssembleMultiSigTx(rawTx, [][]byte{signature1, signature2},2)
	assert.NoError(t, err)

	hash, err := client.BroadcastTx(tx)
	assert.NoError(t, err)
	t.Log("hash", hash)
}
