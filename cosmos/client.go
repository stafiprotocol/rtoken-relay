package cosmos

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	xAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	"os"
)

//cosmos client
type Client struct {
	clientCtx   client.Context
	rpcClient   rpcClient.Client
	fromKeyInfo keyring.Info
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
		clientCtx:   initClientCtx,
		rpcClient:   rpcClient,
		fromKeyInfo: info,
	}
}

func (c *Client) SendTo(toAddr types.AccAddress, amount types.Coins) error {
	msg := xBankTypes.NewMsgSend(c.clientCtx.GetFromAddress(), toAddr, amount)
	cmd := cobra.Command{}
	return clientTx.GenerateOrBroadcastTxCLI(c.clientCtx, cmd.Flags(), msg)
}
