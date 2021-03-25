package rpc

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	xAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	rpcCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
	"os"
)

//cosmos client
type Client struct {
	clientCtx client.Context
	rpcClient rpcClient.Client
	gasPrice  string
	denom     string
}

func NewClient(rpcClient rpcClient.Client, k keyring.Keyring, chainId, fromName string) (*Client, error) {
	encodingConfig := MakeEncodingConfig()
	info, err := k.Key(fromName)
	if err != nil {
		return nil, fmt.Errorf("keyring get address from name:%s err: %s", fromName, err)
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
	}, nil
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

func (c *Client) GetFromName() string {
	return c.clientCtx.FromName
}

func (c *Client) SetGasPrice(gasPrice string) error {
	//todo check value
	c.gasPrice = gasPrice
	return nil
}

func (c *Client) SetDenom(denom string) {
	c.denom = denom
}

func (c *Client) GetDenom() string {
	return c.denom
}

func (c *Client) GetTxConfig() client.TxConfig {
	return c.clientCtx.TxConfig
}

func (c *Client) GetLegacyAmino() *codec.LegacyAmino {
	return c.clientCtx.LegacyAmino
}

func (c *Client) GetStatus() (*rpcCoreTypes.ResultStatus, error) {
	return c.clientCtx.Client.Status(context.Background())
}

func (c *Client) Sign(fromName string, toBeSigned []byte) ([]byte, types.PubKey, error) {
	return c.clientCtx.Keyring.Sign(fromName, toBeSigned)
}