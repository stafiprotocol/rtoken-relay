package bnb_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/rtoken-relay/chains/bnb"
	"github.com/stafiprotocol/rtoken-relay/core"
)

func TestConnection(t *testing.T) {
	stop := make(chan int)
	connection, err := bnb.NewConnection(&core.ChainConfig{
		Name:            "bnb test",
		Symbol:          "RBNB",
		Endpoint:        "https://rpc.ankr.com/bsc",
		Accounts:        []string{},
		KeystorePath:    "",
		Insecure:        false,
		LatestBlockFlag: false,
		Opts: map[string]interface{}{
			"bcEndpoint":          "https://api.binance.org",
			"bscSideChainId":      "bsc",
			"StakePortalContract": "0x7498ee6c2bf4dbd6af1a7c3eb085f7595abc0491",
			"StakingContract":     "0x0000000000000000000000000000000000002001",
		},
	}, core.NewLog("module", "test"), stop, true)
	if err != nil {
		t.Fatal(err)
	}

	logrus.SetLevel(logrus.TraceLevel)
	height, timestamp, err := connection.GetHeightTimestampByEra(1370, 86400, -18033)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(height, timestamp)

	_, lastEraTimestamp, err := connection.GetHeightTimestampByEra(1367, 86400, -18033)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lastEraTimestamp)

	poolAddr := common.HexToAddress("0x992701C853301B120A6dB0E4767AEEcC737c515a")
	rewardOnBsc, err := connection.GetStakingContract().GetDistributedReward(&bind.CallOpts{
		BlockNumber: big.NewInt(height),
		Context:     context.Background(),
	}, poolAddr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rewardOnBsc)

	undelegatedAmount, err := connection.GetStakingContract().GetUndelegated(&bind.CallOpts{
		BlockNumber: big.NewInt(height),
		Context:     context.Background(),
	}, poolAddr)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(undelegatedAmount)

	poolBalance, err := connection.GetQueryClient().Client().BalanceAt(context.Background(), poolAddr, big.NewInt(height))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(poolBalance)

	// total, timestamp, err := connection.RewardTotalTimesAndLastRewardTimestamp("bnb14jvtqvzarzflcwz74vz0y8r6nezqhuxhdf2k4y")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(total, timestamp)
	reward, err := connection.RewardOnBcDuTimes(poolAddr, timestamp, lastEraTimestamp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reward)
}
