package cosmos_test

import (
	"testing"
	"time"

	"github.com/stafiprotocol/rtoken-relay/shared/cosmos"
	"github.com/stafiprotocol/rtoken-relay/shared/cosmos/rpc"
)

func TestGetCurrentEra(t *testing.T) {
	client, err := rpc.NewClient(nil, "cosmoshub-4", "", "0.00001uatom", "uatom", "https://cosmos-rpc1.stafi.io:443")
	if err != nil {
		t.Fatal(err)
	}
	poolClient := cosmos.NewPoolClient(client, "", 87000, 18245)

	era, err := poolClient.GetCurrentEra()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(era)

	t.Log(time.Now().Unix()/87000 - 18245)

	height, err := poolClient.GetHeightByEra(era)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(height)

	block, err := client.QueryBlock(height)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(block.Block.Header.Time.Unix()/87000 - 18245)
}
