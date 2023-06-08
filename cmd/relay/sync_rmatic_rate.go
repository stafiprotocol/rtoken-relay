package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/rtoken-relay/bindings/RMaticStakeManager"
	stake_portal_rate "github.com/stafiprotocol/rtoken-relay/bindings/StakePortalRate"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
	"github.com/urfave/cli/v2"
)

var syncRMaticRateCommand = cli.Command{
	Name:   "sync-rmatic-rate",
	Usage:  "sync rmatic rate from eth to polygon",
	Action: handleSyncRMaticRate,
	Flags: []cli.Flag{
		config.ConfigFileFlag,
	},
}

type SyncRMaticRateConfig struct {
	EthEndpoint                    string `json:"ethEndpoint"`
	PolygonEndpoint                string `json:"polygonEndpoint"`
	KeystorePath                   string `json:"keystorePath"`
	Account                        string `json:"account"`
	EthStakeManagerContract        string `json:"ethStakeManagerContract"`
	PolygonStakePortalRateContract string `json:"polygonStakePortalRateContract"`
}

func loadSyncRMaticRateConfig(file string, config *SyncRMaticRateConfig) (err error) {
	ext := filepath.Ext(file)
	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	core.NewLog().Debug("Loading configuration", "path", filepath.Clean(fp))

	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()

	if ext != ".json" {
		return fmt.Errorf("unrecognized extention: %s", ext)
	}
	return json.NewDecoder(f).Decode(&config)
}

func handleSyncRMaticRate(ctx *cli.Context) error {
	path := "./config_sync_rmatic_rate.json"
	if file := ctx.String(config.ConfigFileFlag.Name); file != "" {
		path = file
	}
	cfg := SyncRMaticRateConfig{}
	err := loadSyncRMaticRateConfig(path, &cfg)
	if err != nil {
		return err
	}
	lvl := ctx.String(config.VerbosityFlag.Name)

	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	if !common.IsHexAddress(cfg.EthStakeManagerContract) {
		return fmt.Errorf("StakeMangerContract is not hex address")
	}
	ethStakManagerContractAddr := common.HexToAddress(cfg.EthStakeManagerContract)

	if !common.IsHexAddress(cfg.PolygonStakePortalRateContract) {
		return fmt.Errorf("PolygonStakePortalRateContract is not hex address")
	}
	polygonStakePortalRateContractAddr := common.HexToAddress(cfg.PolygonStakePortalRateContract)

	ethClient := ethereum.NewClient(cfg.EthEndpoint, nil, log, big.NewInt(0), big.NewInt(0))
	if err := ethClient.Connect(); err != nil {
		return err
	}

	kpI, err := keystore.KeypairFromAddress(cfg.Account, keystore.EthChain, cfg.KeystorePath, false)
	if err != nil {
		return err
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	ethStakeManager, err := rmatic_stake_manager.NewRMaticStakeManager(ethStakManagerContractAddr, ethClient.Client())
	if err != nil {
		return err
	}

	polygonClient := ethereum.NewClient(cfg.PolygonEndpoint, kp, log, big.NewInt(0), big.NewInt(0))
	if err := polygonClient.Connect(); err != nil {
		return err
	}
	polygonStakePortalRate, err := stake_portal_rate.NewStakePortalRate(polygonStakePortalRateContractAddr, polygonClient.Client())
	if err != nil {
		return err
	}

	for {
		rateOnEth, err := ethStakeManager.GetRate(nil)
		if err != nil {
			logrus.Warnf("ethStakeManager.GetRate failed, err: %s", err.Error())
			continue
		}

		rateOnPolygon, err := polygonStakePortalRate.GetRate(nil)
		if err != nil {
			logrus.Warnf("polygonStakePortalRate.GetRate failed, err: %s", err.Error())
			continue
		}

		if rateOnEth.Cmp(rateOnPolygon) == 0 {
			time.Sleep(time.Minute * 10)
			continue
		}

		latestEra, err := ethStakeManager.LatestEra(nil)
		if err != nil {
			logrus.Warnf("ethStakeManager.LatestEra failed, err: %s", err.Error())
			continue
		}
		proposalId := getProposalId(uint32(latestEra.Uint64()), rateOnEth, 0)
		err = polygonVoteRate(polygonStakePortalRate, proposalId, rateOnEth, polygonClient)
		if err != nil {
			logrus.Warnf("polygonVoteRate failed, err: %s", err.Error())
			continue
		}
	}
}

func polygonVoteRate(polygonStakePortalRateContract *stake_portal_rate.StakePortalRate, proposalId [32]byte, evmRate *big.Int, polygonConn *ethereum.Client) error {
	proposal, err := polygonStakePortalRateContract.Proposals(nil, proposalId)
	if err != nil {
		return fmt.Errorf("processSignatureEnough Proposals error %s ", err)
	}
	if proposal.Status == 2 { // success status
		return nil
	}
	hasVoted, err := polygonStakePortalRateContract.HasVoted(&bind.CallOpts{}, proposalId, polygonConn.Opts().From)
	if err != nil {
		return fmt.Errorf("processSignatureEnough HasVoted error %s", err)
	}
	if hasVoted {
		return nil
	}

	// send tx
	err = polygonConn.LockAndUpdateOpts(big.NewInt(0), big.NewInt(0))
	if err != nil {
		return fmt.Errorf("processSignatureEnough LockAndUpdateOpts error %s", err)
	}
	polygonConn.UnlockOpts()

	voteTx, err := polygonStakePortalRateContract.VoteRate(polygonConn.Opts(), proposalId, evmRate)
	if err != nil {
		return fmt.Errorf("processSignatureEnough VoteRate error %s", err)
	}

	err = waitPolygonTxOk(voteTx.Hash(), polygonConn)
	if err != nil {
		return fmt.Errorf("processSignatureEnough waitTxOk error %s", err)
	}

	err = waitPolygonRateUpdated(polygonStakePortalRateContract, proposalId)
	if err != nil {
		return fmt.Errorf("processSignatureEnough waitRateUpdated error %s", err)
	}
	return nil
}

func getProposalId(era uint32, rate *big.Int, factor int) common.Hash {
	return crypto.Keccak256Hash([]byte(fmt.Sprintf("era-%d-%s-%s-%d", era, "voteRate", rate.String(), factor)))
}

func waitPolygonTxOk(txHash common.Hash, polygonConn *ethereum.Client) error {
	retry := 0
	for {
		if retry > 300 {
			return fmt.Errorf("waitPolygonTxOk tx reach retry limit")
		}
		_, pending, err := polygonConn.TransactionByHash(context.Background(), txHash)
		if err == nil && !pending {
			break
		} else {
			if err != nil {
				logrus.Warn("tx status", "hash", txHash, "err", err.Error())
			} else {
				logrus.Warn("tx status", "hash", txHash, "status", "pending")
			}
			time.Sleep(6 * time.Second)
			retry++
			continue
		}

	}
	logrus.Info("tx send ok", "tx", txHash.String())
	return nil
}

func waitPolygonRateUpdated(polygonStakePortalRateContract *stake_portal_rate.StakePortalRate, proposalId [32]byte) error {
	retry := 0
	for {
		if retry > 300 {
			return fmt.Errorf("waitPolygonRateUpdated tx reach retry limit")
		}

		proposal, err := polygonStakePortalRateContract.Proposals(&bind.CallOpts{}, proposalId)
		if err != nil {
			time.Sleep(6 * time.Second)
			retry++
			continue
		}
		if proposal.Status != 2 {
			time.Sleep(6 * time.Second)
			retry++
			continue
		}
		break
	}
	return nil
}
