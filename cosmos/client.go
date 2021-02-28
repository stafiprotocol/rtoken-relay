package cosmos

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	xAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	"os"
)

//cosmos client
type Client struct {
	clientCtx client.Context
	rpcClient rpcClient.Client
}

func NewClient(rpcClient rpcClient.Client, k keyring.Keyring, chainId, fromKeyName string) *Client {
	encodingConfig := MakeEncodingConfig()
	info, err := k.Key(fromKeyName)
	if err != nil {
		panic(fmt.Sprintf("keyring get address from fromKeyname err: %s", err))
	}

	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(xAuthTypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithChainID(chainId).
		WithClient(rpcClient).
		WithSkipConfirmation(true).
		WithFromName(fromKeyName). //keyBase need keyName to  find key info
		WithFromAddress(info.GetAddress()). //accountRetriever need fromAddress
		WithKeyring(k)

	return &Client{
		clientCtx: initClientCtx,
		rpcClient: rpcClient,
	}
}

//update from key
func (c *Client) SetFromKey(fromKey string) error {
	info, err := c.clientCtx.Keyring.Key(fromKey)
	if err != nil {
		return fmt.Errorf("keyring get address from fromKeyname err: %s", err)
	}

	c.clientCtx = c.clientCtx.WithFromName(fromKey).WithFromAddress(info.GetAddress())
	return nil
}

