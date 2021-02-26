package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/crypto/sr25519"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"github.com/stafiprotocol/rtoken-relay/substrate"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

type ChainType string

const (
	Substrate = ChainType("substrate")
)

func Start(cfg *config.Config, log log15.Logger) {
	ctx, cancel := context.WithCancel(context.Background())
	sysErr := make(chan error)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(ctx context.Context) {
		mc := cfg.MainConf
		sblk, ok := mc.Opts["startBlock"]
		if !ok {
			sblk = "0"
		}
		krp, bs, blk, err := utils.BlockstoreAndKeyring(mc.From, mc.KeystorePath, mc.BlockstorePath, sblk)
		if err != nil {
			sysErr <- fmt.Errorf("BlockstoreAndKeyring error: %s", err)
			return
		}

		sc, err := substrate.NewSarpcClient(cfg.MainConf.Endpoint, cfg.MainConf.TypesPath, log.New("Sarpc", cfg.MainConf.Name))
		if err != nil {
			sysErr <- fmt.Errorf("NewSarpcClient error: %s", err)
			return
		}

		gc, err := substrate.NewGsrpcClient(ctx, cfg.MainConf.Endpoint, krp, log.New("Gsrpc", cfg.MainConf.Name))
		if err != nil {
			sysErr <- fmt.Errorf("NewGsrpcClient error: %s", err)
			return
		}

		chainEras := make(map[conn.RSymbol]types.U32)
		err = ChainEras(chainEras, gc, cfg, log)
		if err != nil {
			sysErr <- err
			return
		}

		chains := make(map[conn.RSymbol]conn.Chain)
		err = ChainClient(ctx, gc, chains, cfg, log)
		if err != nil {
			sysErr <- err
			return
		}

		listener := NewListener(ctx, sc, gc, bs, blk, chainEras, chains, sysErr, log)
		err = listener.Start()
		if err != nil {
			sysErr <- fmt.Errorf("listener start error: %s", err)
			return
		}
	}(ctx)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)

	// Block here and wait for a signal
	select {
	case err := <-sysErr:
		log.Error("FATAL ERROR. Shutting down.", "err", err)
	case <-sigc:
		log.Warn("Interrupt received, shutting down now.")
	}

	cancel()
}

func ChainClient(ctx context.Context, gc *substrate.GsrpcClient, chains map[conn.RSymbol]conn.Chain, cfg *config.Config, log log15.Logger) error {
	for _, chainConf := range cfg.OtherConfs {
		ctype := ChainType(chainConf.Type)
		switch ctype {
		case Substrate:
			//cc := new(substrate.ChainClient)
			sc, err := substrate.NewSarpcClient(chainConf.Endpoint, chainConf.TypesPath, log.New("Sarpc", chainConf.Name))
			if err != nil {
				return fmt.Errorf("NewSarpcClient error: %s for chain: %s", err, chainConf.Name)
			}
			sym := conn.RSymbol(chainConf.Symbol)
			_, err = types.EncodeToHexString(sym)
			if err != nil {
				return err
			}

			keys := make([]*signature.KeyringPair, 0)
			cls := make(map[*signature.KeyringPair]*substrate.GsrpcClient)
			for _, account := range chainConf.Accounts {
				kp, err := keystore.KeypairFromAddress(account, keystore.SubChain, chainConf.KeystorePath, false)
				if err != nil {
					return err
				}

				krp := kp.(*sr25519.Keypair).AsKeyringPair()

				gc, err := substrate.NewGsrpcClient(ctx, chainConf.Endpoint, krp, log.New("Gsrpc", chainConf.Name))
				if err != nil {
					return fmt.Errorf("ChainClient NewGsrpcClient error: %s", err)
				}

				keys = append(keys, krp)
				cls[krp] = gc
			}

			chains[sym] = &substrate.FullSubClient{sc, gc, keys, cls}
		}
	}

	return nil
}

func ChainEras(chainEras map[conn.RSymbol]types.U32, gc *substrate.GsrpcClient, cfg *config.Config, log log15.Logger) error {
	for _, chainConf := range cfg.OtherConfs {
		sym := conn.RSymbol(chainConf.Symbol)
		symBz, err := types.EncodeToBytes(sym)
		if err != nil {
			return err
		}

		var era types.U32
		exists, err := gc.QueryStorage(config.RTokenLedgerModuleId, config.StorageChainEras, symBz, nil, &era)
		if err != nil {
			return err
		}
		if !exists {
			era = 0
		}
		chainEras[sym] = era
	}

	return nil
}
