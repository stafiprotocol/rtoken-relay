package main

import (
	"context"
	"encoding/json"
	"fmt"
	multisig_onchain "github.com/stafiprotocol/rtoken-relay/bindings/MultisigOnchain"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	staking "github.com/stafiprotocol/rtoken-relay/bindings/Staking"
	"github.com/stafiprotocol/rtoken-relay/chains/bnb"
	"github.com/stafiprotocol/rtoken-relay/chains/matic"
	"github.com/stafiprotocol/rtoken-relay/chains/solana"
	"github.com/stafiprotocol/rtoken-relay/chains/substrate"
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

var app = cli.NewApp()

var mainFlags = []cli.Flag{
	config.ConfigFileFlag,
	config.VerbosityFlag,
}

var generateFlags = []cli.Flag{
	config.KeystorePathFlag,
	config.NetworkFlag,
}

var accountCommand = cli.Command{
	Name:  "accounts",
	Usage: "manage reth keystore",
	Description: "The accounts command is used to manage the relay keystore.\n" +
		"\tMake sure the keystore dir is exist before generating\n" +
		"\tTo generate a substrate keystore: relay accounts gensub\n" +
		"\tTo generate a ethereum keystore: relay accounts geneth\n" +
		"\tTo generate a bc chain keystore: relay accounts genbc\n" +
		"\tTo list keys: chainbridge accounts list",
	Subcommands: []*cli.Command{
		{
			Action:      handleGenerateSubCmd,
			Name:        "gensub",
			Usage:       "generate subsrate keystore",
			Flags:       generateFlags,
			Description: "The generate subcommand is used to generate the substrate keystore.",
		},
		{
			Action:      handleGenerateEthCmd,
			Name:        "geneth",
			Usage:       "generate ethereum keystore",
			Flags:       generateFlags,
			Description: "The generate subcommand is used to generate the ethereum keystore.",
		},
	},
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

		multisigOnchain, err := multisig_onchain.NewMultisigOnchain(common.HexToAddress(cfg.MultiSigContract), client.Client())
		if err != nil {
			return err
		}
		pool := Pool{
			poolAddress:     common.HexToAddress(cfg.MultiSigContract),
			bscClient:       client,
			multisigOnchain: multisigOnchain,
		}
		pools = append(pools, &pool)
	}

	stakingContract, err := staking.NewStaking(common.HexToAddress(cfg.StakingContract), client)
	if err != nil {
		return err
	}
	relayerFee, err := stakingContract.GetRelayerFee(&bind.CallOpts{
		From:    common.HexToAddress(cfg.MultiSigContract),
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
	proposal, err := getDelegateProposal(totalAmount, relayerFeeDeci, common.HexToAddress(cfg.StakingContract), common.HexToAddress(cfg.Validator))
	if err != nil {
		return err
	}

	for _, p := range pools {
		error := submitProposal(p, [32]byte{1, 1, 1}, proposal)
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
	logrus.Warn("submitProposal ok ", " pool ", pool.poolAddress, " proposalId ", hexutil.Encode(proposalId[:]), " txHash ", tx.Hash())
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

// init initializes CLI
func init() {
	app.Action = run
	app.Copyright = "Copyright 2020 Stafi Protocol Authors"
	app.Name = "reley"
	app.Usage = "relay"
	app.Authors = []*cli.Author{{Name: "Stafi Protocol 2020"}}
	app.Version = "0.1.0"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		&accountCommand,
		&delegateBnbCommand,
	}

	app.Flags = append(app.Flags, mainFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startLogger(ctx *cli.Context) error {
	lvl := ctx.String(config.VerbosityFlag.Name)

	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	err = core.InitLogFile("./log_file")
	if err != nil {
		return err
	}

	return nil
}

func run(ctx *cli.Context) error {
	err := startLogger(ctx)
	if err != nil {
		return err
	}

	cfg, err := config.GetConfig(ctx)
	if err != nil {
		return err
	}

	// Used to signal core shutdown due to fatal error
	sysErr := make(chan error)
	c := core.NewCore(sysErr)

	for _, chain := range cfg.Chains {
		//check symbol
		switch chain.Rsymbol {
		case string(core.RBNB), string(core.RDOT), string(core.RFIS), string(core.RKSM), string(core.RMATIC), string(core.RSOL):
		default:
			return fmt.Errorf("rsymbol not match: %s", chain.Rsymbol)
		}

		chainConfig := &core.ChainConfig{
			Name:            chain.Name,
			Endpoint:        chain.Endpoint,
			KeystorePath:    chain.KeystorePath,
			Symbol:          core.RSymbol(chain.Rsymbol),
			Accounts:        chain.Accounts,
			LatestBlockFlag: chain.LatestBlockFlag,
			Insecure:        false,
			Opts:            chain.Opts,
		}
		var newChain core.Chain
		logger := core.NewLog("chain", chainConfig.Name)

		switch chain.Type {
		case "substrate":
			newChain, err = substrate.InitializeChain(chainConfig, logger, sysErr)
			if err != nil {
				return err
			}
		case "solana":
			newChain, err = solana.InitializeChain(chainConfig, logger, sysErr)
			if err != nil {
				return err
			}

		case "ethereum":
			newChain, err = matic.InitializeChain(chainConfig, logger, sysErr)
			if err != nil {
				return err
			}
		case "bnb":
			newChain, err = bnb.InitializeChain(chainConfig, logger, sysErr)
			if err != nil {
				return err
			}
		default:
			return errors.New("unrecognized Chain Type")
		}

		c.AddChain(newChain)
	}

	c.Start()

	return nil
}
