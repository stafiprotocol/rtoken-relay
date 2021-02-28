package cosmos

import (
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

func NewClient(rpcClient rpcClient.Client, k keyring.Keyring, chainId, fromName string) *Client {
	encodingConfig := MakeEncodingConfig()
	info, err := k.Key(fromName)
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
		WithSkipConfirmation(true). //skip password confirm
		WithFromName(fromName). //keyBase need FromName to find key info
		WithFromAddress(info.GetAddress()). //accountRetriever need FromAddress
		WithKeyring(k)

	return &Client{
		clientCtx: initClientCtx,
		rpcClient: rpcClient,
	}
}

//update clientCtx.FromName and clientCtx.FromAddress
func (c *Client) SetFromName(fromName string) error {
	info, err := c.clientCtx.Keyring.Key(fromName)
	if err != nil {
		return fmt.Errorf("keyring get address from fromKeyname err: %s", err)
	}

	c.clientCtx = c.clientCtx.WithFromName(fromName).WithFromAddress(info.GetAddress())
	return nil
}
