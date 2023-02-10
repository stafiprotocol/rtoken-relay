package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	multisig_onchain "github.com/stafiprotocol/rtoken-relay/bindings/MultisigOnchain"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
	"github.com/urfave/cli/v2"
)

var claimUndelegateCommand = cli.Command{
	Name:   "claim-undelegate",
	Usage:  "claim undelegate",
	Action: handleClaimUndelegate,
	Flags: []cli.Flag{
		config.ConfigFileFlag,
	},
}

func handleClaimUndelegate(ctx *cli.Context) error {
	path := "./config_delegate_bnb.json"
	if file := ctx.String(config.ConfigFileFlag.Name); file != "" {
		path = file
	}
	cfg := DelegateBnbConfig{}
	err := loadConfig(path, &cfg)
	if err != nil {
		return err
	}
	lvl := ctx.String(config.VerbosityFlag.Name)

	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	if !common.IsHexAddress(cfg.StakingContract) {
		return fmt.Errorf("StakingContract is not hex address")
	}
	stakingContractAddr := common.HexToAddress(cfg.StakingContract)

	if !common.IsHexAddress(cfg.MultiSigContract) {
		return fmt.Errorf("MultiSigContract is not hex address")
	}
	multisigContractAddr := common.HexToAddress(cfg.MultiSigContract)

	pools := make([]*Pool, 0)
	for i := 0; i < len(cfg.Accounts); i++ {
		kpI, err := keystore.KeypairFromAddress(cfg.Accounts[i], keystore.EthChain, cfg.KeystorePath, false)
		if err != nil {
			return err
		}
		kp, _ := kpI.(*secp256k1.Keypair)

		client := ethereum.NewClient(cfg.Endpoint, kp, log, big.NewInt(0), big.NewInt(0))
		if err := client.Connect(); err != nil {
			return err
		}

		multisigOnchain, err := multisig_onchain.NewMultisigOnchain(multisigContractAddr, client.Client())
		if err != nil {
			return err
		}
		pool := Pool{
			poolAddress:     multisigContractAddr,
			bscClient:       client,
			multisigOnchain: multisigOnchain,
		}
		pools = append(pools, &pool)
	}

	proposal, err := getClaimUndelegateProposal(stakingContractAddr)
	if err != nil {
		return err
	}

	proposalId := crypto.Keccak256Hash([]byte(fmt.Sprintf("claim-undelegate-%d", cfg.ProposalFactor)))
	for _, p := range pools {
		error := submitProposal(p, proposalId, proposal)
		if err != nil {
			return error
		}
	}
	return nil
}
