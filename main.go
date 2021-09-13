package main

import (
	"fmt"
	bncRpc "github.com/stafiprotocol/go-sdk/client/rpc"
	bncCmnTypes "github.com/stafiprotocol/go-sdk/common/types"
)

func main() {
	client := bncRpc.NewRPCClient("tcp://data-seed-pre-1-s3.binance.org:80", bncCmnTypes.TestNetwork)
	fmt.Println("IsActive", client.IsActive())
}
