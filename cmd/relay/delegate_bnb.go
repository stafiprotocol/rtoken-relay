package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	multisig_onchain "github.com/stafiprotocol/rtoken-relay/bindings/MultisigOnchain"
	staking "github.com/stafiprotocol/rtoken-relay/bindings/Staking"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/core"
	"github.com/stafiprotocol/rtoken-relay/models/ethmodel"
	"github.com/stafiprotocol/rtoken-relay/shared/ethereum"
	"github.com/urfave/cli/v2"
)

var (
	StakingAbi                      abi.ABI
	ErrNoAvailableValsForUnDelegate = errors.New("ErrNoAvailableValsForUnDelegate")
)

func init() {

	stakingAbi, err := abi.JSON(strings.NewReader(staking.StakingABI))
	if err != nil {
		panic(err)
	}
	StakingAbi = stakingAbi

}

var delegateBnbCommand = cli.Command{
	Name:   "delegate-bnb",
	Usage:  "delegate bnb",
	Action: handleDelegateBnb,
	Flags: []cli.Flag{
		config.ConfigFileFlag,
	},
}

type DelegateBnbConfig struct {
	Endpoint         string   `json:"endpoint"`
	KeystorePath     string   `json:"keystorePath"`
	Accounts         []string `json:"accounts"`
	Amount           string   `json:"amount"`
	Validator        string   `json:"validator"`
	StakingContract  string   `json:"StakingContract"`
	MultiSigContract string   `json:"MultiSigContract"`
	ProposalFactor   uint64   `json:"proposalFactor"`
}

type Pool struct {
	poolAddress     common.Address
	bscClient       *ethereum.Client
	multisigOnchain *multisig_onchain.MultisigOnchain
}

func handleDelegateBnb(ctx *cli.Context) error {
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

	client, err := ethclient.Dial(cfg.Endpoint)
	if err != nil {
		return err
	}

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

	stakingContract, err := staking.NewStaking(stakingContractAddr, client)
	if err != nil {
		return err
	}
	relayerFee, err := stakingContract.GetRelayerFee(&bind.CallOpts{
		From:    multisigContractAddr,
		Context: context.Background(),
	})
	if err != nil {
		return errors.Wrap(err, "stakingContract.GetRelayerFee")
	}
	relayerFeeDeci := decimal.NewFromBigInt(relayerFee, 0)
	totalAmount, err := decimal.NewFromString(cfg.Amount)
	if err != nil {
		return err
	}
	proposal, err := getDelegateProposal(totalAmount, relayerFeeDeci, stakingContractAddr, common.HexToAddress(cfg.Validator))
	if err != nil {
		return err
	}

	proposalId := crypto.Keccak256Hash([]byte(fmt.Sprintf("delegate-bnb-%d", cfg.ProposalFactor)))
	for _, p := range pools {
		error := submitProposal(p, proposalId, proposal)
		if err != nil {
			return error
		}
	}
	return nil
}

func submitProposal(pool *Pool, proposalId [32]byte, proposalBts []byte) error {
	err := pool.bscClient.LockAndUpdateOpts(big.NewInt(0), big.NewInt(0))
	if err != nil {
		return errors.Wrap(err, "LockAndUpdateOpts")
	}
	defer pool.bscClient.UnlockOpts()

	tx, err := pool.multisigOnchain.ExecTransactions(pool.bscClient.Opts(), proposalId, proposalBts)
	if err != nil {
		return errors.Wrap(err, "multisigOnchain.ExecTransactions")
	}
	retry := 0
	for {
		if retry > 50*2 {
			return fmt.Errorf("multisigOnchain.ExecTransactions tx reach retry limit")
		}
		_, pending, err := pool.bscClient.Client().TransactionByHash(context.Background(), tx.Hash())
		if err == nil && !pending {
			break
		} else {
			if err != nil {
				logrus.Warn("tx status ", " hash ", tx.Hash(), " err ", err.Error())
			} else {
				logrus.Warn("tx status ", " hash ", tx.Hash(), " status ", "pending")
			}
			time.Sleep(time.Second * 3)
			retry++
			continue
		}
	}
	logrus.Info("submitProposal ok ", " pool ", pool.poolAddress, " proposalId ", hexutil.Encode(proposalId[:]), " txHash ", tx.Hash())
	return nil
}

func getDelegateProposal(totalAmount, relayerFee decimal.Decimal, stakingAddress common.Address, validator common.Address) ([]byte, error) {

	txs := make(ethmodel.BatchTransactions, 0)
	inputData, err := StakingAbi.Pack("delegate", validator, totalAmount.BigInt())
	if err != nil {
		return nil, errors.Wrap(err, "staking abi pack failed")
	}

	tx := &ethmodel.BatchTransaction{
		Operation:  uint8(ethmodel.Call),
		To:         stakingAddress,
		Value:      totalAmount.Add(relayerFee).BigInt(),
		DataLength: big.NewInt(int64(len(inputData))),
		Data:       inputData,
	}
	txs = append(txs, tx)

	return txs.Encode(), nil
}

func getClaimUndelegateProposal(stakingAddress common.Address) ([]byte, error) {
	txs := make(ethmodel.BatchTransactions, 0)
	inputData, err := StakingAbi.Pack("claimUndelegated")
	if err != nil {
		return nil, errors.Wrap(err, "staking abi pack failed")
	}

	tx := &ethmodel.BatchTransaction{
		Operation:  uint8(ethmodel.Call),
		To:         stakingAddress,
		Value:      big.NewInt(0),
		DataLength: big.NewInt(int64(len(inputData))),
		Data:       inputData,
	}
	txs = append(txs, tx)

	return txs.Encode(), nil
}

func loadConfig(file string, config *DelegateBnbConfig) (err error) {
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
