package ethereum

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func NewSimpleClient(endpoint string) (*ethclient.Client, error) {
	ctx := context.Background()
	rpcClient, err := rpc.DialWebsocket(ctx, endpoint, "/ws")
	if err != nil {
		return nil, err
	}

	client := ethclient.NewClient(rpcClient)
	return client, nil
}
